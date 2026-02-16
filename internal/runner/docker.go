// Package runner fornece a interface e implementações para executar projetos.
package runner

import (
	"context"
	"fmt"

	"github.com/relief-org/relief/internal/domain"
	"github.com/relief-org/relief/pkg/logger"
)

// DockerRunner executa projetos usando Docker/Docker Compose
type DockerRunner struct {
	*BaseRunner
	logger *logger.Logger
}

// NewDockerRunner cria uma nova instância de DockerRunner
func NewDockerRunner(log *logger.Logger) *DockerRunner {
	return &DockerRunner{
		BaseRunner: NewBaseRunner(RunnerTypeDocker),
		logger:     log,
	}
}

// Start inicia o projeto usando Docker
func (r *DockerRunner) Start(ctx context.Context, project *domain.Project) error {
	// TODO: Implementar usando Docker SDK
	// - Verificar se docker-compose.yml existe
	// - Gerar docker-compose.yml dinâmico se necessário
	// - Usar github.com/docker/docker/client para gerenciar containers
	// - Add relief.project=<name> labels to containers

	r.logger.Info("DockerRunner.Start chamado", map[string]interface{}{
		"project": project.Name,
	})

	return fmt.Errorf("DockerRunner ainda não implementado - use NativeRunner")
}

// Stop para o projeto
func (r *DockerRunner) Stop(ctx context.Context, projectID string) error {
	// TODO: Implementar
	return fmt.Errorf("DockerRunner ainda não implementado")
}

// Status retorna o status do projeto
func (r *DockerRunner) Status(projectID string) (*RunnerStatus, error) {
	// TODO: Implementar
	return nil, fmt.Errorf("DockerRunner ainda não implementado")
}

// GetLogs retorna os logs do projeto
func (r *DockerRunner) GetLogs(projectID string, tail int) ([]domain.LogEntry, error) {
	// TODO: Implementar
	return nil, fmt.Errorf("DockerRunner ainda não implementado")
}

// Restart reinicia o projeto
func (r *DockerRunner) Restart(ctx context.Context, project *domain.Project) error {
	// TODO: Implementar
	return fmt.Errorf("DockerRunner ainda não implementado")
}
