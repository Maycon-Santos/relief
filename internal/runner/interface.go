package runner

import (
	"context"
	"time"

	"github.com/relief-org/relief/internal/domain"
)

type ProjectRunner interface {
	Start(ctx context.Context, project *domain.Project) error

	Stop(ctx context.Context, projectID string) error

	Status(projectID string) (*RunnerStatus, error)

	GetLogs(projectID string, tail int) ([]domain.LogEntry, error)

	Restart(ctx context.Context, project *domain.Project) error
}

type RunnerStatus struct {
	ProjectID  string
	Status     domain.Status
	PID        int
	Port       int
	Uptime     time.Duration
	MemoryUsed int64
	CPUUsed    float64
	Message    string
}

type RunnerType string

const (
	RunnerTypeDocker RunnerType = "docker"
	RunnerTypeNative RunnerType = "native"
)

type BaseRunner struct {
	Type          RunnerType
	LogBuffer     []domain.LogEntry
	MaxLogEntries int
}

func NewBaseRunner(runnerType RunnerType) *BaseRunner {
	return &BaseRunner{
		Type:          runnerType,
		LogBuffer:     make([]domain.LogEntry, 0),
		MaxLogEntries: 1000,
	}
}

func (b *BaseRunner) AddLog(projectID, level, message string) {
	entry := domain.LogEntry{
		ProjectID: projectID,
		Level:     level,
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	b.LogBuffer = append(b.LogBuffer, entry)

	if len(b.LogBuffer) > b.MaxLogEntries {
		b.LogBuffer = b.LogBuffer[len(b.LogBuffer)-b.MaxLogEntries:]
	}
}

func (b *BaseRunner) GetLogsFromBuffer(projectID string, tail int) []domain.LogEntry {
	projectLogs := []domain.LogEntry{}
	for _, log := range b.LogBuffer {
		if log.ProjectID == projectID {
			projectLogs = append(projectLogs, log)
		}
	}

	if tail > 0 && len(projectLogs) > tail {
		return projectLogs[len(projectLogs)-tail:]
	}

	return projectLogs
}

func (b *BaseRunner) ClearLogs(projectID string) {
	filtered := []domain.LogEntry{}
	for _, log := range b.LogBuffer {
		if log.ProjectID != projectID {
			filtered = append(filtered, log)
		}
	}
	b.LogBuffer = filtered
}
