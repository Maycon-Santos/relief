package checkers

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/relief-org/relief/pkg/fileutil"
	"github.com/relief-org/relief/pkg/logger"
)

type PythonChecker struct {
	logger *logger.Logger
	path   string
}

func NewPythonChecker(log *logger.Logger) *PythonChecker {
	return &PythonChecker{
		logger: log,
	}
}

func (c *PythonChecker) Check(ctx context.Context) (string, error) {
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

func (c *PythonChecker) Install(ctx context.Context, version string) error {
	depsDir, err := fileutil.GetReliefSubDir(filepath.Join("deps", "python", version))
	if err != nil {
		return fmt.Errorf("erro ao criar diretório de deps: %w", err)
	}

	c.path = depsDir


	c.logger.Info("Python instalação completa", map[string]interface{}{
		"version": version,
		"path":    depsDir,
	})

	return fmt.Errorf("instalação automática de Python ainda não implementada - por favor instale manualmente")
}

func (c *PythonChecker) GetPath() string {
	if c.path == "" {
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
