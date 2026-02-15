// Package dependency gerencia verificação e instalação de dependências.
package dependency

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/omelete/sofredor-orchestrator/internal/domain"
	"github.com/omelete/sofredor-orchestrator/internal/dependency/checkers"
	"github.com/omelete/sofredor-orchestrator/pkg/logger"
)

// Manager gerencia verificação e instalação de dependências
type Manager struct {
	checkers map[string]Checker
	logger   *logger.Logger
}

// Checker define a interface para verificadores de dependência
type Checker interface {
	// Check verifica se a dependência está instalada e retorna a versão
	Check(ctx context.Context) (string, error)
	
	// Install instala a dependência
	Install(ctx context.Context, version string) error
	
	// GetPath retorna o path do binário da dependência
	GetPath() string
}

// NewManager cria uma nova instância de Manager
func NewManager(log *logger.Logger) *Manager {
	m := &Manager{
		checkers: make(map[string]Checker),
		logger:   log,
	}

	// Registrar checkers
	m.checkers["node"] = checkers.NewNodeChecker(log)
	m.checkers["python"] = checkers.NewPythonChecker(log)
	m.checkers["postgres"] = checkers.NewPostgresChecker(log)

	return m
}

// CheckDependencies verifica todas as dependências de um projeto
func (m *Manager) CheckDependencies(ctx context.Context, project *domain.Project) error {
	m.logger.Info("Verificando dependências", map[string]interface{}{
		"project": project.Name,
		"count":   len(project.Dependencies),
	})

	for i := range project.Dependencies {
		dep := &project.Dependencies[i]
		
		if err := m.checkDependency(ctx, dep); err != nil {
			m.logger.Warn("Dependência não satisfeita", map[string]interface{}{
				"dependency": dep.Name,
				"error":      err.Error(),
			})
			dep.Satisfied = false
			dep.Message = err.Error()
		} else {
			dep.Satisfied = true
			m.logger.Info("Dependência satisfeita", map[string]interface{}{
				"dependency": dep.Name,
				"version":    dep.Version,
			})
		}
	}

	return nil
}

// checkDependency verifica uma dependência específica
func (m *Manager) checkDependency(ctx context.Context, dep *domain.Dependency) error {
	checker, exists := m.checkers[dep.Name]
	if !exists {
		// Se não há checker específico, tentar verificar via command
		return m.checkGenericCommand(ctx, dep)
	}

	// Verificar versão instalada
	installedVersion, err := checker.Check(ctx)
	if err != nil {
		if dep.Managed {
			// Tentar instalar
			m.logger.Info("Tentando instalar dependência", map[string]interface{}{
				"dependency": dep.Name,
				"version":    dep.RequiredVersion,
			})
			
			if err := checker.Install(ctx, dep.RequiredVersion); err != nil {
				return fmt.Errorf("erro ao instalar: %w", err)
			}

			// Verificar novamente
			installedVersion, err = checker.Check(ctx)
			if err != nil {
				return fmt.Errorf("dependência não encontrada após instalação: %w", err)
			}
		} else {
			return fmt.Errorf("não instalada: %w", err)
		}
	}

	// Validar versão
	if err := m.validateVersion(installedVersion, dep.RequiredVersion); err != nil {
		return err
	}

	dep.Version = installedVersion
	return nil
}

// validateVersion valida se a versão instalada satisfaz o requisito
func (m *Manager) validateVersion(installed, required string) error {
	// Limpar prefixos comuns
	installed = strings.TrimPrefix(installed, "v")
	required = strings.TrimPrefix(required, "v")

	// Remover operadores do required
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

	// Remover espaços
	installed = strings.TrimSpace(installed)
	required = strings.TrimSpace(required)

	// Se não houver operador e versions são iguais, OK
	if operator == "" && installed == required {
		return nil
	}

	// Parse versions
	vInstalled, err := version.NewVersion(installed)
	if err != nil {
		return fmt.Errorf("erro ao fazer parse da versão instalada: %w", err)
	}

	vRequired, err := version.NewVersion(required)
	if err != nil {
		return fmt.Errorf("erro ao fazer parse da versão requerida: %w", err)
	}

	// Comparar
	switch operator {
	case ">=":
		if vInstalled.LessThan(vRequired) {
			return fmt.Errorf("versão %s não satisfaz requisito >=%s", installed, required)
		}
	case ">":
		if vInstalled.LessThanOrEqual(vRequired) {
			return fmt.Errorf("versão %s não satisfaz requisito >%s", installed, required)
		}
	case "<=":
		if vInstalled.GreaterThan(vRequired) {
			return fmt.Errorf("versão %s não satisfaz requisito <=%s", installed, required)
		}
	case "<":
		if vInstalled.GreaterThanOrEqual(vRequired) {
			return fmt.Errorf("versão %s não satisfaz requisito <%s", installed, required)
		}
	case "=", "":
		if !vInstalled.Equal(vRequired) {
			return fmt.Errorf("versão %s não corresponde a %s", installed, required)
		}
	}

	return nil
}

// checkGenericCommand verifica uma dependência genérica via comando
func (m *Manager) checkGenericCommand(ctx context.Context, dep *domain.Dependency) error {
	// Tentar executar <name> --version
	cmd := exec.CommandContext(ctx, dep.Name, "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("comando não encontrado")
	}

	// Extrair versão da output (heurística simples)
	version := strings.TrimSpace(string(output))
	lines := strings.Split(version, "\n")
	if len(lines) > 0 {
		version = lines[0]
	}

	dep.Version = version
	dep.Satisfied = true
	
	m.logger.Info("Dependência genérica verificada", map[string]interface{}{
		"dependency": dep.Name,
		"version":    version,
	})

	return nil
}

// InstallDependency instala uma dependência específica
func (m *Manager) InstallDependency(ctx context.Context, name, version string) error {
	checker, exists := m.checkers[name]
	if !exists {
		return fmt.Errorf("instalador não disponível para %s", name)
	}

	m.logger.Info("Instalando dependência", map[string]interface{}{
		"name":    name,
		"version": version,
	})

	if err := checker.Install(ctx, version); err != nil {
		return fmt.Errorf("erro ao instalar %s: %w", name, err)
	}

	m.logger.Info("Dependência instalada com sucesso", map[string]interface{}{
		"name":    name,
		"version": version,
	})

	return nil
}

// GetDependencyPath retorna o path de uma dependência instalada
func (m *Manager) GetDependencyPath(name string) (string, error) {
	checker, exists := m.checkers[name]
	if !exists {
		return "", fmt.Errorf("checker não encontrado para %s", name)
	}

	return checker.GetPath(), nil
}
