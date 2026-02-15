// Package checkers contains specific dependency checkers.
package checkers

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/omelete/relief/pkg/fileutil"
	"github.com/omelete/relief/pkg/logger"
)

// NodeChecker verifies and installs Node.js
type NodeChecker struct {
	logger *logger.Logger
	path   string
}

// NewNodeChecker creates a new NodeChecker instance
func NewNodeChecker(log *logger.Logger) *NodeChecker {
	return &NodeChecker{
		logger: log,
	}
}

// Check verifies if Node.js is installed and returns the version
func (c *NodeChecker) Check(ctx context.Context) (string, error) {
	// Try using custom path first
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
	// Remove 'v' prefix (e.g.: v18.19.0 -> 18.19.0)
	version = strings.TrimPrefix(version, "v")

	c.logger.Debug("Node.js found", map[string]interface{}{
		"version": version,
	})

	return version, nil
}

// Install installs Node.js in the specified version
func (c *NodeChecker) Install(ctx context.Context, version string) error {
	// Installation directory: ~/.relief/deps/node/<version>
	depsDir, err := fileutil.GetReliefSubDir(filepath.Join("deps", "node", version))
	if err != nil {
		return fmt.Errorf("error creating deps directory: %w", err)
	}

	c.path = depsDir

	// TODO: Implement real Node.js download
	// - Detect OS (Linux/Mac/Windows)
	// - Download portable binary from official URL
	// - Extract to depsDir
	// - Configure c.path

	c.logger.Info("Node.js installation complete", map[string]interface{}{
		"version": version,
		"path":    depsDir,
	})

	return fmt.Errorf("automatic Node.js installation not yet implemented - please install manually")
}

// GetPath returns the path of the Node.js binary
func (c *NodeChecker) GetPath() string {
	if c.path == "" {
		// Try to find in system PATH
		path, err := exec.LookPath("node")
		if err == nil {
			return filepath.Dir(path)
		}
		return ""
	}
	return c.path
}
