// Package runner fornece a interface e implementações para executar projetos.
package runner

import (
	"context"
	"time"

	"github.com/relief-org/relief/internal/domain"
)

// ProjectRunner define a interface para executores de projeto
type ProjectRunner interface {
	// Start inicia o projeto
	Start(ctx context.Context, project *domain.Project) error

	// Stop para o projeto
	Stop(ctx context.Context, projectID string) error

	// Status retorna o status atual do projeto
	Status(projectID string) (*RunnerStatus, error)

	// GetLogs retorna os logs do projeto
	GetLogs(projectID string, tail int) ([]domain.LogEntry, error)

	// Restart reinicia o projeto
	Restart(ctx context.Context, project *domain.Project) error
}

// RunnerStatus representa o status de um runner
type RunnerStatus struct {
	ProjectID  string
	Status     domain.Status
	PID        int
	Port       int
	Uptime     time.Duration
	MemoryUsed int64   // em bytes
	CPUUsed    float64 // percentual
	Message    string
}

// RunnerType define o tipo de runner
type RunnerType string

const (
	RunnerTypeDocker RunnerType = "docker"
	RunnerTypeNative RunnerType = "native"
)

// BaseRunner contém campos comuns a todos os runners
type BaseRunner struct {
	Type          RunnerType
	LogBuffer     []domain.LogEntry
	MaxLogEntries int
}

// NewBaseRunner cria uma nova instância de BaseRunner
func NewBaseRunner(runnerType RunnerType) *BaseRunner {
	return &BaseRunner{
		Type:          runnerType,
		LogBuffer:     make([]domain.LogEntry, 0),
		MaxLogEntries: 1000, // Manter últimas 1000 linhas
	}
}

// AddLog adiciona uma entrada de log ao buffer
func (b *BaseRunner) AddLog(projectID, level, message string) {
	entry := domain.LogEntry{
		ProjectID: projectID,
		Level:     level,
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	b.LogBuffer = append(b.LogBuffer, entry)

	// Limitar tamanho do buffer
	if len(b.LogBuffer) > b.MaxLogEntries {
		b.LogBuffer = b.LogBuffer[len(b.LogBuffer)-b.MaxLogEntries:]
	}
}

// GetLogsFromBuffer retorna logs do buffer
func (b *BaseRunner) GetLogsFromBuffer(projectID string, tail int) []domain.LogEntry {
	// Filtrar logs do projeto
	projectLogs := []domain.LogEntry{}
	for _, log := range b.LogBuffer {
		if log.ProjectID == projectID {
			projectLogs = append(projectLogs, log)
		}
	}

	// Retornar últimas N entradas
	if tail > 0 && len(projectLogs) > tail {
		return projectLogs[len(projectLogs)-tail:]
	}

	return projectLogs
}

// ClearLogs limpa o buffer de logs de um projeto
func (b *BaseRunner) ClearLogs(projectID string) {
	filtered := []domain.LogEntry{}
	for _, log := range b.LogBuffer {
		if log.ProjectID != projectID {
			filtered = append(filtered, log)
		}
	}
	b.LogBuffer = filtered
}
