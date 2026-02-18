package runner

import (
	"context"
	"fmt"

	"github.com/Maycon-Santos/relief/internal/domain"
	"github.com/Maycon-Santos/relief/pkg/logger"
)

type DockerRunner struct {
	*BaseRunner
	logger *logger.Logger
}

func NewDockerRunner(log *logger.Logger) *DockerRunner {
	return &DockerRunner{
		BaseRunner: NewBaseRunner(RunnerTypeDocker),
		logger:     log,
	}
}

func (r *DockerRunner) Start(ctx context.Context, project *domain.Project) error {

	r.logger.Info("DockerRunner.Start chamado", map[string]interface{}{
		"project": project.Name,
	})

	return fmt.Errorf("DockerRunner ainda não implementado - use NativeRunner")
}

func (r *DockerRunner) Stop(ctx context.Context, projectID string) error {
	return fmt.Errorf("DockerRunner ainda não implementado")
}

func (r *DockerRunner) Status(projectID string) (*RunnerStatus, error) {
	return nil, fmt.Errorf("DockerRunner ainda não implementado")
}

func (r *DockerRunner) GetLogs(projectID string, tail int) ([]domain.LogEntry, error) {
	return nil, fmt.Errorf("DockerRunner ainda não implementado")
}

func (r *DockerRunner) Restart(ctx context.Context, project *domain.Project) error {
	return fmt.Errorf("DockerRunner ainda não implementado")
}
