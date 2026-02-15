// Package checkers contém verificadores de dependências específicos.
package checkers

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/omelete/sofredor-orchestrator/pkg/logger"
)

// PostgresChecker verifica e gerencia PostgreSQL
type PostgresChecker struct {
	logger *logger.Logger
	path   string
}

// NewPostgresChecker cria uma nova instância de PostgresChecker
func NewPostgresChecker(log *logger.Logger) *PostgresChecker {
	return &PostgresChecker{
		logger: log,
	}
}

// Check verifica se PostgreSQL está instalado e retorna a versão
func (c *PostgresChecker) Check(ctx context.Context) (string, error) {
	psqlCmd := "psql"
	if c.path != "" {
		psqlCmd = filepath.Join(c.path, "psql")
	}

	cmd := exec.CommandContext(ctx, psqlCmd, "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("postgres não encontrado: %w", err)
	}

	version := strings.TrimSpace(string(output))
	// Output é "psql (PostgreSQL) 15.3", extrair versão
	parts := strings.Fields(version)
	for i, part := range parts {
		if strings.HasPrefix(part, "(PostgreSQL)") && i+1 < len(parts) {
			return parts[i+1], nil
		}
	}

	// Fallback: pegar último campo
	if len(parts) > 0 {
		return parts[len(parts)-1], nil
	}

	c.logger.Debug("PostgreSQL encontrado", map[string]interface{}{
		"raw_output": version,
	})

	return version, nil
}

// Install instala/inicia PostgreSQL
func (c *PostgresChecker) Install(ctx context.Context, version string) error {
	// Para PostgreSQL, podemos iniciar um container Docker ao invés de instalar
	c.logger.Info("Tentando iniciar PostgreSQL via Docker", map[string]interface{}{
		"version": version,
	})

	// TODO: Implementar container PostgreSQL via Docker
	// docker run -d --name sofredor-postgres-<version> \
	//   -e POSTGRES_PASSWORD=sofredor \
	//   -p 5432:5432 \
	//   postgres:<version>

	return fmt.Errorf("instalação automática de PostgreSQL ainda não implementada - use Docker ou instale manualmente")
}

// GetPath retorna o path do binário PostgreSQL
func (c *PostgresChecker) GetPath() string {
	if c.path == "" {
		path, err := exec.LookPath("psql")
		if err == nil {
			return filepath.Dir(path)
		}
		return ""
	}
	return c.path
}
