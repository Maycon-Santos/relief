// Package runner fornece a interface e implementações para executar projetos.
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

	"github.com/relief-org/relief/internal/domain"
	"github.com/relief-org/relief/pkg/logger"
)

// NativeRunner executa projetos como processos nativos do SO
type NativeRunner struct {
	*BaseRunner
	processes map[string]*ProcessInfo
	mu        sync.RWMutex
	logger    *logger.Logger
}

// ProcessInfo armazena informações sobre um processo em execução
type ProcessInfo struct {
	Project   *domain.Project
	Cmd       *exec.Cmd
	PID       int
	StartedAt time.Time
	Stdout    io.ReadCloser
	Stderr    io.ReadCloser
	Cancel    context.CancelFunc
}

// NewNativeRunner cria uma nova instância de NativeRunner
func NewNativeRunner(log *logger.Logger) *NativeRunner {
	return &NativeRunner{
		BaseRunner: NewBaseRunner(RunnerTypeNative),
		processes:  make(map[string]*ProcessInfo),
		logger:     log,
	}
}

// Start inicia o projeto como processo nativo
func (r *NativeRunner) Start(ctx context.Context, project *domain.Project) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Verificar se já está rodando
	if _, exists := r.processes[project.ID]; exists {
		return fmt.Errorf("projeto %s já está em execução", project.Name)
	}

	// Verificar se o manifest está presente
	if project.Manifest == nil {
		return fmt.Errorf("manifest não carregado para o projeto %s", project.Name)
	}

	// Obter script de dev
	devScript := project.Manifest.GetDevScript()
	if devScript == "" {
		return fmt.Errorf("script 'dev' não encontrado no manifest")
	}

	r.logger.Info("Starting project with script", map[string]interface{}{
		"project": project.Name,
		"script":  devScript,
	})

	// Criar contexto cancelável para o processo
	processCtx, cancel := context.WithCancel(ctx)

	// Preparar comando (executar via shell para suportar comandos complexos)
	cmd := exec.CommandContext(processCtx, "sh", "-c", devScript)
	cmd.Dir = project.Path

	// Configurar variáveis de ambiente
	cmd.Env = os.Environ()
	for key, value := range project.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	// Adicionar PORT se não existir
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

	// Capturar stdout e stderr
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

	// Iniciar processo
	if err := cmd.Start(); err != nil {
		cancel()
		return fmt.Errorf("erro ao iniciar processo: %w", err)
	}

	// Armazenar informações do processo
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

	// Atualizar projeto
	project.PID = cmd.Process.Pid
	project.UpdateStatus(domain.StatusRunning)

	// Iniciar goroutines para capturar logs
	go r.captureOutput(project.ID, stdout, "info")
	go r.captureOutput(project.ID, stderr, "error")

	// Monitorar processo
	go r.monitorProcess(project.ID)

	r.logger.Info("Projeto iniciado", map[string]interface{}{
		"project": project.Name,
		"pid":     cmd.Process.Pid,
	})

	return nil
}

// Stop para o projeto
func (r *NativeRunner) Stop(ctx context.Context, projectID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	processInfo, exists := r.processes[projectID]
	if !exists {
		return fmt.Errorf("projeto não está em execução")
	}

	// Cancelar contexto (envia SIGTERM)
	processInfo.Cancel()

	// Aguardar término gracioso por 5 segundos
	done := make(chan error, 1)
	go func() {
		done <- processInfo.Cmd.Wait()
	}()

	select {
	case <-time.After(5 * time.Second):
		// Force kill se não terminar
		if processInfo.Cmd.Process != nil {
			r.logger.Warn("Forçando término do processo", map[string]interface{}{
				"project": processInfo.Project.Name,
				"pid":     processInfo.PID,
			})
			processInfo.Cmd.Process.Signal(syscall.SIGKILL)
		}
	case <-done:
		// Processo terminou graciosamente
	}

	// Remover do mapa de processos
	delete(r.processes, projectID)

	r.logger.Info("Projeto parado", map[string]interface{}{
		"project": processInfo.Project.Name,
	})

	return nil
}

// Status retorna o status do projeto
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

// GetLogs retorna os logs do projeto
func (r *NativeRunner) GetLogs(projectID string, tail int) ([]domain.LogEntry, error) {
	return r.GetLogsFromBuffer(projectID, tail), nil
}

// Restart reinicia o projeto
func (r *NativeRunner) Restart(ctx context.Context, project *domain.Project) error {
	// Parar se estiver rodando
	if _, exists := r.processes[project.ID]; exists {
		if err := r.Stop(ctx, project.ID); err != nil {
			return fmt.Errorf("erro ao parar projeto: %w", err)
		}

		// Aguardar um pouco
		time.Sleep(1 * time.Second)
	}

	// Iniciar novamente
	return r.Start(ctx, project)
}

// captureOutput captura a saída do processo e adiciona aos logs
func (r *NativeRunner) captureOutput(projectID string, reader io.ReadCloser, level string) {
	defer reader.Close()

	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			message := string(buf[:n])
			// Remover newlines extras
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

// monitorProcess monitora o processo e atualiza o status
func (r *NativeRunner) monitorProcess(projectID string) {
	r.mu.RLock()
	processInfo, exists := r.processes[projectID]
	r.mu.RUnlock()

	if !exists {
		return
	}

	// Aguardar término do processo
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

// GetRunningProcesses retorna a lista de processos em execução
func (r *NativeRunner) GetRunningProcesses() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	projectIDs := make([]string, 0, len(r.processes))
	for id := range r.processes {
		projectIDs = append(projectIDs, id)
	}
	return projectIDs
}
