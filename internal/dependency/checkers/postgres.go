package checkers

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/relief-org/relief/pkg/logger"
)

type PostgresChecker struct {
	logger *logger.Logger
	path   string
}

func NewPostgresChecker(log *logger.Logger) *PostgresChecker {
	return &PostgresChecker{
		logger: log,
	}
}

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
	parts := strings.Fields(version)
	for i, part := range parts {
		if strings.HasPrefix(part, "(PostgreSQL)") && i+1 < len(parts) {
			return parts[i+1], nil
		}
	}

	if len(parts) > 0 {
		return parts[len(parts)-1], nil
	}

	c.logger.Debug("PostgreSQL encontrado", map[string]interface{}{
		"raw_output": version,
	})

	return version, nil
}

func (c *PostgresChecker) Install(ctx context.Context, version string) error {
	c.logger.Info("Tentando iniciar PostgreSQL via Docker", map[string]interface{}{
		"version": version,
	})


	return fmt.Errorf("instalação automática de PostgreSQL ainda não implementada - use Docker ou instale manualmente")
}

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
