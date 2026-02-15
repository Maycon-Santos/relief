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

// NodeChecker verifica e instala Node.js
type NodeChecker struct {
	logger *logger.Logger
	path   string
}

// NewNodeChecker cria uma nova instância de NodeChecker
func NewNodeChecker(log *logger.Logger) *NodeChecker {
	return &NodeChecker{
		logger: log,
	}
}

// Check verifica se Node.js está instalado e retorna a versão
func (c *NodeChecker) Check(ctx context.Context) (string, error) {
	// Tentar usar path customizado primeiro
	nodeCmd := "node"
	if c.path != "" {
		nodeCmd = filepath.Join(c.path, "node")
	}

	cmd := exec.CommandContext(ctx, nodeCmd, "-v")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("node não encontrado: %w", err)
	}

	version := strings.TrimSpace(string(output))
	// Remover 'v' prefix (ex: v18.19.0 -> 18.19.0)
	version = strings.TrimPrefix(version, "v")

	c.logger.Debug("Node.js encontrado", map[string]interface{}{
		"version": version,
	})

	return version, nil
}

// Install instala Node.js na versão especificada
func (c *NodeChecker) Install(ctx context.Context, version string) error {
	// Diretório de instalação: ~/.sofredor/deps/node/<version>
	depsDir, err := fileutil.GetSofredorSubDir(filepath.Join("deps", "node", version))
	if err != nil {
		return fmt.Errorf("erro ao criar diretório de deps: %w", err)
	}

	c.path = depsDir

	// TODO: Implementar download real do Node.js
	// - Detectar SO (Linux/Mac/Windows)
	// - Baixar binário portable da URL oficial
	// - Extrair para depsDir
	// - Configurar c.path

	c.logger.Info("Node.js instalação completa", map[string]interface{}{
		"version": version,
		"path":    depsDir,
	})

	return fmt.Errorf("instalação automática de Node.js ainda não implementada - por favor instale manualmente")
}

// GetPath retorna o path do binário Node.js
func (c *NodeChecker) GetPath() string {
	if c.path == "" {
		// Tentar encontrar no PATH do sistema
		path, err := exec.LookPath("node")
		if err == nil {
			return filepath.Dir(path)
		}
		return ""
	}
	return c.path
}
