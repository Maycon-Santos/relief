package app

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// PortConflict representa um conflito de porta
type PortConflict struct {
	Port    int    `json:"port"`
	PID     int    `json:"pid"`
	Command string `json:"command"`
}

// CheckPortInUse verifica se uma porta está em uso e retorna informações sobre o processo
func (a *App) CheckPortInUse(port int) (*PortConflict, error) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin", "linux":
		// lsof -ti:PORT retorna o PID
		cmd = exec.Command("lsof", "-ti:"+strconv.Itoa(port))
	case "windows":
		// netstat -ano | findstr :PORT
		cmd = exec.Command("netstat", "-ano")
	default:
		return nil, fmt.Errorf("sistema operacional não suportado")
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Se lsof retornar erro, significa que a porta está livre
		return nil, nil
	}

	pidStr := strings.TrimSpace(string(output))
	if pidStr == "" {
		return nil, nil
	}

	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear PID: %w", err)
	}

	// Obter comando do processo
	command := ""
	psCmd := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "command=")
	if psOutput, err := psCmd.CombinedOutput(); err == nil {
		command = strings.TrimSpace(string(psOutput))
	}

	return &PortConflict{
		Port:    port,
		PID:     pid,
		Command: command,
	}, nil
}

// KillProcessByPID mata um processo específico
func (a *App) KillProcessByPID(pid int) error {
	a.logger.Info("Encerrando processo", map[string]interface{}{
		"pid": pid,
	})

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin", "linux":
		cmd = exec.Command("kill", "-9", strconv.Itoa(pid))
	case "windows":
		cmd = exec.Command("taskkill", "/F", "/PID", strconv.Itoa(pid))
	default:
		return fmt.Errorf("sistema operacional não suportado")
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("erro ao encerrar processo: %w", err)
	}

	a.logger.Info("Processo encerrado com sucesso", map[string]interface{}{
		"pid": pid,
	})

	return nil
}
