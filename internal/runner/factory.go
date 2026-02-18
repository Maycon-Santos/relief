package runner

import (
	"fmt"

	"github.com/Maycon-Santos/relief/internal/domain"
	"github.com/Maycon-Santos/relief/pkg/logger"
)

type Factory struct {
	logger *logger.Logger
}

func NewFactory(log *logger.Logger) *Factory {
	return &Factory{
		logger: log,
	}
}

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

func (f *Factory) GetAllRunners() map[string]ProjectRunner {
	return map[string]ProjectRunner{
		"docker": NewDockerRunner(f.logger),
		"native": NewNativeRunner(f.logger),
	}
}
