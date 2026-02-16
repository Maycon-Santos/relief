// Package runner fornece a interface e implementações para executar projetos.
package runner

import (
	"fmt"

	"github.com/relief-org/relief/internal/domain"
	"github.com/relief-org/relief/pkg/logger"
)

// Factory cria runners baseado no tipo de projeto
type Factory struct {
	logger *logger.Logger
}

// NewFactory cria uma nova instância de Factory
func NewFactory(log *logger.Logger) *Factory {
	return &Factory{
		logger: log,
	}
}

// CreateRunner cria um runner apropriado para o projeto
func (f *Factory) CreateRunner(project *domain.Project) (ProjectRunner, error) {
	switch project.Type {
	case domain.ProjectTypeDocker:
		return NewDockerRunner(f.logger), nil

	case domain.ProjectTypeNode,
		domain.ProjectTypePython,
		domain.ProjectTypeGo,
		domain.ProjectTypeJava,
		domain.ProjectTypeRuby:
		return NewNativeRunner(f.logger), nil

	default:
		return nil, fmt.Errorf("tipo de projeto não suportado: %s", project.Type)
	}
}

// GetAllRunners retorna todos os runners disponíveis
func (f *Factory) GetAllRunners() map[string]ProjectRunner {
	return map[string]ProjectRunner{
		"docker": NewDockerRunner(f.logger),
		"native": NewNativeRunner(f.logger),
	}
}
