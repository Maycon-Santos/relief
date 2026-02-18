package checkers

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Maycon-Santos/relief/pkg/fileutil"
	"github.com/Maycon-Santos/relief/pkg/logger"
)

type NodeChecker struct {
	logger *logger.Logger
	path   string
}

func NewNodeChecker(log *logger.Logger) *NodeChecker {
	return &NodeChecker{
		logger: log,
	}
}

func (c *NodeChecker) Check(ctx context.Context) (string, error) {
	nodeCmd := "node"
	if c.path != "" {
		nodeCmd = filepath.Join(c.path, "node")
	}

	cmd := exec.CommandContext(ctx, nodeCmd, "-v")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("node not found: %w", err)
	}

	version := strings.TrimSpace(string(output))
	version = strings.TrimPrefix(version, "v")

	c.logger.Debug("Node.js found", map[string]interface{}{
		"version": version,
	})

	return version, nil
}

func (c *NodeChecker) Install(ctx context.Context, version string) error {
	depsDir, err := fileutil.GetReliefSubDir(filepath.Join("deps", "node", version))
	if err != nil {
		return fmt.Errorf("error creating deps directory: %w", err)
	}

	c.path = depsDir


	c.logger.Info("Node.js installation complete", map[string]interface{}{
		"version": version,
		"path":    depsDir,
	})

	return fmt.Errorf("automatic Node.js installation not yet implemented - please install manually")
}

func (c *NodeChecker) GetPath() string {
	if c.path == "" {
		path, err := exec.LookPath("node")
		if err == nil {
			return filepath.Dir(path)
		}
		return ""
	}
	return c.path
}
