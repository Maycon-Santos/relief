// Package checkers contém verificadores de dependências específicos.
package checkers

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/omelete/sofredor-orchestrator/pkg/fileutil"
	"github.com/omelete/sofredor-orchestrator/pkg/logger"
)

// PythonChecker verifica e instala Python
type PythonChecker struct {
	logger *logger.Logger
	path   string
}

// NewPythonChecker cria uma nova instância de PythonChecker
func NewPythonChecker(log *logger.Logger) *PythonChecker {
	return &PythonChecker{
		logger: log,
	}
}

// Check verifica se Python está instalado e retorna a versão
func (c *PythonChecker) Check(ctx context.Context) (string, error) {
	// Tentar python3 primeiro, depois python
	pythonCmds := []string{"python3", "python"}
	
	if c.path != "" {
		pythonCmds = []string{
			filepath.Join(c.path, "python3"),
			filepath.Join(c.path, "python"),
		}
	}

	for _, pythonCmd := range pythonCmds {
		cmd := exec.CommandContext(ctx, pythonCmd, "--version")
		output, err := cmd.CombinedOutput()
		if err == nil {
			version := strings.TrimSpace(string(output))
			// Output é "Python 3.9.7", extrair apenas a versão
			parts := strings.Fields(version)
			if len(parts) >= 2 {
				version = parts[1]
			}

			c.logger.Debug("Python encontrado", map[string]interface{}{
				"version": version,
				"command": pythonCmd,
			})

			return version, nil
		}
	}

	return "", fmt.Errorf("python não encontrado")
}

// Install instala Python na versão especificada
func (c *PythonChecker) Install(ctx context.Context, version string) error {
	// Diretório de instalação: ~/.sofredor/deps/python/<version>
	depsDir, err := fileutil.GetSofredorSubDir(filepath.Join("deps", "python", version))
	if err != nil {
		return fmt.Errorf("erro ao criar diretório de deps: %w", err)
	}

	c.path = depsDir

	// TODO: Implementar download real do Python
	// - Detectar SO
	// - Baixar Python portable/standalone
	// - Extrair para depsDir

	c.logger.Info("Python instalação completa", map[string]interface{}{
		"version": version,
		"path":    depsDir,
	})

	return fmt.Errorf("instalação automática de Python ainda não implementada - por favor instale manualmente")
}

// GetPath retorna o path do binário Python
func (c *PythonChecker) GetPath() string {
	if c.path == "" {
		// Tentar encontrar no PATH do sistema
		for _, cmd := range []string{"python3", "python"} {
			path, err := exec.LookPath(cmd)
			if err == nil {
				return filepath.Dir(path)
			}
		}
		return ""
	}
	return c.path
}
