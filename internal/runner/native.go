package runner

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/Maycon-Santos/relief/internal/domain"
	"github.com/Maycon-Santos/relief/pkg/logger"
)

type NativeRunner struct {
	*BaseRunner
	processes map[string]*ProcessInfo
	mu        sync.RWMutex
	logger    *logger.Logger
}

type ProcessInfo struct {
	Project   *domain.Project
	Cmd       *exec.Cmd
	PID       int
	StartedAt time.Time
	Stdout    io.ReadCloser
	Stderr    io.ReadCloser
	Cancel    context.CancelFunc
}

func NewNativeRunner(log *logger.Logger) *NativeRunner {
	return &NativeRunner{
		BaseRunner: NewBaseRunner(RunnerTypeNative),
		processes:  make(map[string]*ProcessInfo),
		logger:     log,
	}
}

func (r *NativeRunner) Start(ctx context.Context, project *domain.Project) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.processes[project.ID]; exists {
		return fmt.Errorf("projeto %s já está em execução", project.Name)
	}

	if project.Manifest == nil {
		return fmt.Errorf("manifest não carregado para o projeto %s", project.Name)
	}

	devScript := project.Manifest.GetDevScript()
	if devScript == "" {
		return fmt.Errorf("script 'dev' não encontrado no manifest")
	}

	r.logger.Info("Starting project with script", map[string]interface{}{
		"project": project.Name,
		"script":  devScript,
	})

	processCtx, cancel := context.WithCancel(ctx)

	cmd := exec.CommandContext(processCtx, "sh", "-c", devScript)
	cmd.Dir = project.Path

	cmd.Env = os.Environ()
	for key, value := range project.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	if project.Port > 0 {
		hasPort := false
		for _, env := range cmd.Env {
			if strings.HasPrefix(env, "PORT=") {
				hasPort = true
				break
			}
		}
		if !hasPort {
			cmd.Env = append(cmd.Env, fmt.Sprintf("PORT=%d", project.Port))
		}
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		return fmt.Errorf("erro ao criar pipe stdout: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		cancel()
		return fmt.Errorf("erro ao criar pipe stderr: %w", err)
	}

	if err := cmd.Start(); err != nil {
		cancel()
		return fmt.Errorf("erro ao iniciar processo: %w", err)
	}

	processInfo := &ProcessInfo{
		Project:   project,
		Cmd:       cmd,
		PID:       cmd.Process.Pid,
		StartedAt: time.Now(),
		Stdout:    stdout,
		Stderr:    stderr,
		Cancel:    cancel,
	}
	r.processes[project.ID] = processInfo

	project.PID = cmd.Process.Pid
	project.UpdateStatus(domain.StatusRunning)

	go r.captureOutput(project.ID, stdout, "info")
	go r.captureOutput(project.ID, stderr, "error")

	go r.monitorProcess(project.ID)

	r.logger.Info("Projeto iniciado", map[string]interface{}{
		"project": project.Name,
		"pid":     cmd.Process.Pid,
	})

	return nil
}

func (r *NativeRunner) Stop(ctx context.Context, projectID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	processInfo, exists := r.processes[projectID]
	if !exists {
		return fmt.Errorf("projeto não está em execução")
	}

	processInfo.Cancel()

	done := make(chan error, 1)
	go func() {
		done <- processInfo.Cmd.Wait()
	}()

	select {
	case <-time.After(5 * time.Second):
		if processInfo.Cmd.Process != nil {
			r.logger.Warn("Forçando término do processo", map[string]interface{}{
				"project": processInfo.Project.Name,
				"pid":     processInfo.PID,
			})
			processInfo.Cmd.Process.Signal(syscall.SIGKILL)
		}
	case <-done:
	}

	delete(r.processes, projectID)

	r.logger.Info("Projeto parado", map[string]interface{}{
		"project": processInfo.Project.Name,
	})

	return nil
}

func (r *NativeRunner) Status(projectID string) (*RunnerStatus, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	processInfo, exists := r.processes[projectID]
	if !exists {
		return &RunnerStatus{
			ProjectID: projectID,
			Status:    domain.StatusStopped,
		}, nil
	}

	uptime := time.Since(processInfo.StartedAt)

	return &RunnerStatus{
		ProjectID: projectID,
		Status:    domain.StatusRunning,
		PID:       processInfo.PID,
		Port:      processInfo.Project.Port,
		Uptime:    uptime,
		Message:   fmt.Sprintf("Rodando há %s", uptime.Round(time.Second)),
	}, nil
}

func (r *NativeRunner) GetLogs(projectID string, tail int) ([]domain.LogEntry, error) {
	return r.GetLogsFromBuffer(projectID, tail), nil
}

func (r *NativeRunner) Restart(ctx context.Context, project *domain.Project) error {
	if _, exists := r.processes[project.ID]; exists {
		if err := r.Stop(ctx, project.ID); err != nil {
			return fmt.Errorf("erro ao parar projeto: %w", err)
		}

		time.Sleep(1 * time.Second)
	}

	return r.Start(ctx, project)
}

func (r *NativeRunner) captureOutput(projectID string, reader io.ReadCloser, level string) {
	defer reader.Close()

	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			message := string(buf[:n])
			message = strings.TrimSpace(message)
			if message != "" {
				r.AddLog(projectID, level, message)
			}
		}
		if err != nil {
			if err != io.EOF {
				r.logger.Error("Erro ao ler output", err, map[string]interface{}{
					"project_id": projectID,
				})
			}
			break
		}
	}
}

func (r *NativeRunner) monitorProcess(projectID string) {
	r.mu.RLock()
	processInfo, exists := r.processes[projectID]
	r.mu.RUnlock()

	if !exists {
		return
	}

	err := processInfo.Cmd.Wait()

	r.mu.Lock()
	delete(r.processes, projectID)
	r.mu.Unlock()

	if err != nil {
		exitCode := -1
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}

		r.logger.Warn("Processo terminou com erro", map[string]interface{}{
			"project":   processInfo.Project.Name,
			"exit_code": exitCode,
		})

		r.AddLog(projectID, "error", fmt.Sprintf("Processo terminou com código %d", exitCode))
	} else {
		r.logger.Info("Processo terminou", map[string]interface{}{
			"project": processInfo.Project.Name,
		})
	}
}

func (r *NativeRunner) GetRunningProcesses() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	projectIDs := make([]string, 0, len(r.processes))
	for id := range r.processes {
		projectIDs = append(projectIDs, id)
	}
	return projectIDs
}
