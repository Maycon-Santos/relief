package dependency

import (
	"context"
	"fmt"
	"strings"

	"github.com/Maycon-Santos/relief/internal/dependency/checkers"
	"github.com/Maycon-Santos/relief/internal/domain"
	"github.com/Maycon-Santos/relief/pkg/logger"
	"github.com/Maycon-Santos/relief/pkg/shellenv"
	"github.com/hashicorp/go-version"
)

type Manager struct {
	checkers map[string]Checker
	logger   *logger.Logger
}

type Checker interface {
	Check(ctx context.Context) (string, error)

	Install(ctx context.Context, version string) error

	GetPath() string
}

func NewManager(log *logger.Logger) *Manager {
	m := &Manager{
		checkers: make(map[string]Checker),
		logger:   log,
	}

	m.checkers["node"] = checkers.NewNodeChecker(log)
	m.checkers["python"] = checkers.NewPythonChecker(log)
	m.checkers["postgres"] = checkers.NewPostgresChecker(log)

	return m
}

func (m *Manager) CheckDependencies(ctx context.Context, project *domain.Project) error {
	m.logger.Info("Verificando dependências", map[string]interface{}{
		"project": project.Name,
		"count":   len(project.Dependencies),
	})

	for i := range project.Dependencies {
		dep := &project.Dependencies[i]

		if err := m.checkDependency(ctx, dep); err != nil {
			m.logger.Warn("Dependency not satisfied", map[string]interface{}{
				"dependency": dep.Name,
				"error":      err.Error(),
			})
			dep.Satisfied = false
			dep.Message = err.Error()
		} else {
			dep.Satisfied = true
			m.logger.Info("Dependency satisfied", map[string]interface{}{
				"dependency": dep.Name,
				"version":    dep.Version,
			})
		}
	}

	return nil
}

func (m *Manager) checkDependency(ctx context.Context, dep *domain.Dependency) error {
	if dep.Managed {
		m.logger.Debug("Dependência gerenciada, pulando verificação", map[string]interface{}{
			"dependency": dep.Name,
		})
		dep.Satisfied = true
		dep.Version = dep.RequiredVersion
		return nil
	}

	checker, exists := m.checkers[dep.Name]
	if !exists {
		return m.checkGenericCommand(ctx, dep)
	}

	installedVersion, err := checker.Check(ctx)
	if err != nil {
		return fmt.Errorf("not installed: %w", err)
	}

	if err := m.validateVersion(installedVersion, dep.RequiredVersion); err != nil {
		return err
	}

	dep.Version = installedVersion
	return nil
}

func (m *Manager) validateVersion(installed, required string) error {
	installed = strings.TrimPrefix(installed, "v")
	required = strings.TrimPrefix(required, "v")

	operator := ""
	if strings.HasPrefix(required, ">=") {
		operator = ">="
		required = strings.TrimPrefix(required, ">=")
	} else if strings.HasPrefix(required, ">") {
		operator = ">"
		required = strings.TrimPrefix(required, ">")
	} else if strings.HasPrefix(required, "<=") {
		operator = "<="
		required = strings.TrimPrefix(required, "<=")
	} else if strings.HasPrefix(required, "<") {
		operator = "<"
		required = strings.TrimPrefix(required, "<")
	} else if strings.HasPrefix(required, "=") {
		operator = "="
		required = strings.TrimPrefix(required, "=")
	}

	installed = strings.TrimSpace(installed)
	required = strings.TrimSpace(required)

	if operator == "" && installed == required {
		return nil
	}

	vInstalled, err := version.NewVersion(installed)
	if err != nil {
		return fmt.Errorf("error parsing installed version: %w", err)
	}

	vRequired, err := version.NewVersion(required)
	if err != nil {
		return fmt.Errorf("error parsing required version: %w", err)
	}

	switch operator {
	case ">=":
		if vInstalled.LessThan(vRequired) {
			return fmt.Errorf("version %s does not satisfy requirement >=%s", installed, required)
		}
	case ">":
		if vInstalled.LessThanOrEqual(vRequired) {
			return fmt.Errorf("version %s does not satisfy requirement >%s", installed, required)
		}
	case "<=":
		if vInstalled.GreaterThan(vRequired) {
			return fmt.Errorf("version %s does not satisfy requirement <=%s", installed, required)
		}
	case "<":
		if vInstalled.GreaterThanOrEqual(vRequired) {
			return fmt.Errorf("version %s does not satisfy requirement <%s", installed, required)
		}
	case "=", "":
		if !vInstalled.Equal(vRequired) {
			return fmt.Errorf("version %s does not match %s", installed, required)
		}
	}

	return nil
}

func (m *Manager) checkGenericCommand(ctx context.Context, dep *domain.Dependency) error {
	commandMap := map[string]string{
		"redis":      "redis-server",
		"mongodb":    "mongod",
		"postgres":   "psql",
		"postgresql": "psql",
	}

	cmdName := dep.Name
	if mapped, exists := commandMap[dep.Name]; exists {
		cmdName = mapped
	}

	cmd := shellenv.CommandContext(ctx, cmdName+" --version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command not found or not in PATH")
	}

	version := strings.TrimSpace(string(output))
	lines := strings.Split(version, "\n")
	if len(lines) > 0 {
		version = lines[0]
	}

	dep.Version = version
	dep.Satisfied = true

	m.logger.Info("Generic dependency verified", map[string]interface{}{
		"dependency": dep.Name,
		"command":    cmdName,
		"version":    version,
	})

	return nil
}

func (m *Manager) InstallDependency(ctx context.Context, name, version string) error {
	checker, exists := m.checkers[name]
	if !exists {
		return fmt.Errorf("installer not available for %s", name)
	}

	m.logger.Info("Installing dependency", map[string]interface{}{
		"name":    name,
		"version": version,
	})

	if err := checker.Install(ctx, version); err != nil {
		return fmt.Errorf("error installing %s: %w", name, err)
	}

	m.logger.Info("Dependency installed successfully", map[string]interface{}{
		"name":    name,
		"version": version,
	})

	return nil
}

func (m *Manager) GetDependencyPath(name string) (string, error) {
	checker, exists := m.checkers[name]
	if !exists {
		return "", fmt.Errorf("checker not found for %s", name)
	}

	return checker.GetPath(), nil
}
