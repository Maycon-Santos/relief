// Package dependency manages dependency verification and installation.
package dependency

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/relief-org/relief/internal/dependency/checkers"
	"github.com/relief-org/relief/internal/domain"
	"github.com/relief-org/relief/pkg/logger"
)

// Manager manages dependency verification and installation
type Manager struct {
	checkers map[string]Checker
	logger   *logger.Logger
}

// Checker defines the interface for dependency checkers
type Checker interface {
	// Check verifies if the dependency is installed and returns the version
	Check(ctx context.Context) (string, error)

	// Install installs the dependency
	Install(ctx context.Context, version string) error

	// GetPath returns the path of the dependency binary
	GetPath() string
}

// NewManager creates a new Manager instance
func NewManager(log *logger.Logger) *Manager {
	m := &Manager{
		checkers: make(map[string]Checker),
		logger:   log,
	}

	// Register checkers
	m.checkers["node"] = checkers.NewNodeChecker(log)
	m.checkers["python"] = checkers.NewPythonChecker(log)
	m.checkers["postgres"] = checkers.NewPostgresChecker(log)

	return m
}

// CheckDependencies verifies all dependencies of a project
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

// checkDependency verifies a specific dependency
func (m *Manager) checkDependency(ctx context.Context, dep *domain.Dependency) error {
	checker, exists := m.checkers[dep.Name]
	if !exists {
		// If there's no specific checker, try to verify via command
		return m.checkGenericCommand(ctx, dep)
	}

	// Verify installed version
	installedVersion, err := checker.Check(ctx)
	if err != nil {
		if dep.Managed {
			// Try to install
			m.logger.Info("Trying to install dependency", map[string]interface{}{
				"dependency": dep.Name,
				"version":    dep.RequiredVersion,
			})

			if err := checker.Install(ctx, dep.RequiredVersion); err != nil {
				return fmt.Errorf("error installing: %w", err)
			}

			// Verify again
			installedVersion, err = checker.Check(ctx)
			if err != nil {
				return fmt.Errorf("dependency not found after installation: %w", err)
			}
		} else {
			return fmt.Errorf("not installed: %w", err)
		}
	}

	// Validate version
	if err := m.validateVersion(installedVersion, dep.RequiredVersion); err != nil {
		return err
	}

	dep.Version = installedVersion
	return nil
}

// validateVersion validates if the installed version satisfies the requirement
func (m *Manager) validateVersion(installed, required string) error {
	// Clean common prefixes
	installed = strings.TrimPrefix(installed, "v")
	required = strings.TrimPrefix(required, "v")

	// Remove operators from required
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

	// Remove spaces
	installed = strings.TrimSpace(installed)
	required = strings.TrimSpace(required)

	// If there's no operator and versions are equal, OK
	if operator == "" && installed == required {
		return nil
	}

	// Parse versions
	vInstalled, err := version.NewVersion(installed)
	if err != nil {
		return fmt.Errorf("error parsing installed version: %w", err)
	}

	vRequired, err := version.NewVersion(required)
	if err != nil {
		return fmt.Errorf("error parsing required version: %w", err)
	}

	// Compare
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

// checkGenericCommand verifies a generic dependency via command
func (m *Manager) checkGenericCommand(ctx context.Context, dep *domain.Dependency) error {
	// Try to execute <name> --version
	cmd := exec.CommandContext(ctx, dep.Name, "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command not found")
	}

	// Extract version from output (simple heuristic)
	version := strings.TrimSpace(string(output))
	lines := strings.Split(version, "\n")
	if len(lines) > 0 {
		version = lines[0]
	}

	dep.Version = version
	dep.Satisfied = true

	m.logger.Info("Generic dependency verified", map[string]interface{}{
		"dependency": dep.Name,
		"version":    version,
	})

	return nil
}

// InstallDependency installs a specific dependency
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

// GetDependencyPath returns the path of an installed dependency
func (m *Manager) GetDependencyPath(name string) (string, error) {
	checker, exists := m.checkers[name]
	if !exists {
		return "", fmt.Errorf("checker not found for %s", name)
	}

	return checker.GetPath(), nil
}
