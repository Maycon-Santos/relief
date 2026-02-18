package git

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Maycon-Santos/relief/internal/domain"
	"github.com/Maycon-Santos/relief/pkg/fileutil"
	"github.com/Maycon-Santos/relief/pkg/logger"
)

type Manager struct {
	logger *logger.Logger
}

func NewManager(log *logger.Logger) *Manager {
	return &Manager{
		logger: log,
	}
}

func (m *Manager) CloneOrUpdate(ctx context.Context, repoURL, targetPath, branch string) error {
	if fileutil.Exists(targetPath) {
		if fileutil.Exists(filepath.Join(targetPath, ".git")) {
			m.logger.Info("Repositório já existe, atualizando...", map[string]interface{}{
				"path": targetPath,
				"repo": repoURL,
			})
			return m.updateRepository(ctx, targetPath, branch)
		} else {
			return fmt.Errorf("diretório %s já existe mas não é um repositório Git", targetPath)
		}
	}

	parentDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório pai: %w", err)
	}

	m.logger.Info("Clonando repositório...", map[string]interface{}{
		"repo":   repoURL,
		"path":   targetPath,
		"branch": branch,
	})

	args := []string{"clone"}
	if branch != "" {
		args = append(args, "-b", branch)
	}
	args = append(args, repoURL, targetPath)

	cmd := exec.CommandContext(ctx, "git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao clonar repositório: %w (output: %s)", err, string(output))
	}

	m.logger.Info("Repositório clonado com sucesso", map[string]interface{}{
		"path": targetPath,
	})

	return nil
}

func (m *Manager) updateRepository(ctx context.Context, path, branch string) error {
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("erro ao obter diretório atual: %w", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(path); err != nil {
		return fmt.Errorf("erro ao mudar para diretório do repositório: %w", err)
	}

	cmd := exec.CommandContext(ctx, "git", "fetch")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("erro ao executar git fetch: %w (output: %s)", err, string(output))
	}

	if branch != "" {
		currentBranch, err := m.getCurrentBranch(ctx)
		if err != nil {
			m.logger.Warn("Não foi possível determinar branch atual", map[string]interface{}{
				"error": err.Error(),
			})
		}

		if currentBranch != branch {
			m.logger.Info("Mudando para branch especificado", map[string]interface{}{
				"current": currentBranch,
				"target":  branch,
			})

			cmd := exec.CommandContext(ctx, "git", "checkout", branch)
			if output, err := cmd.CombinedOutput(); err != nil {
				return fmt.Errorf("erro ao trocar de branch: %w (output: %s)", err, string(output))
			}
		}
	}

	cmd = exec.CommandContext(ctx, "git", "pull")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("erro ao executar git pull: %w (output: %s)", err, string(output))
	}

	m.logger.Info("Repositório atualizado", map[string]interface{}{
		"path":   path,
		"branch": branch,
	})

	return nil
}

func (m *Manager) getCurrentBranch(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("erro ao obter branch atual: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

func (m *Manager) IsRepository(path string) bool {
	return fileutil.Exists(filepath.Join(path, ".git"))
}

func (m *Manager) GetRemoteURL(ctx context.Context, path string) (string, error) {
	originalDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("erro ao obter diretório atual: %w", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(path); err != nil {
		return "", fmt.Errorf("erro ao mudar para diretório do repositório: %w", err)
	}

	cmd := exec.CommandContext(ctx, "git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("erro ao obter URL remota: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

func (m *Manager) GetGitInfo(ctx context.Context, path string) (*domain.GitInfo, error) {
	gitInfo := &domain.GitInfo{
		IsRepository: m.IsRepository(path),
	}

	if !gitInfo.IsRepository {
		return gitInfo, nil
	}

	originalDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter diretório atual: %w", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(path); err != nil {
		return nil, fmt.Errorf("erro ao mudar para diretório do repositório: %w", err)
	}

	if currentBranch, err := m.getCurrentBranch(ctx); err == nil {
		gitInfo.CurrentBranch = currentBranch
	}

	if branches, err := m.getBranches(ctx); err == nil {
		gitInfo.AvailableBranches = branches
	}

	if remoteURL, err := m.getRemoteURL(ctx); err == nil {
		gitInfo.RemoteURL = remoteURL
	}

	if hasChanges, err := m.hasUncommittedChanges(ctx); err == nil {
		gitInfo.HasChanges = hasChanges
	}

	if lastCommit, err := m.getLastCommit(ctx); err == nil {
		gitInfo.LastCommit = lastCommit
	}

	return gitInfo, nil
}

func (m *Manager) getBranches(ctx context.Context) ([]string, error) {
	cmd := exec.CommandContext(ctx, "git", "branch", "--format=%(refname:short)")
	localOutput, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter branches locais: %w", err)
	}

	branches := []string{}
	localBranches := strings.Split(strings.TrimSpace(string(localOutput)), "\n")
	for _, branch := range localBranches {
		if branch = strings.TrimSpace(branch); branch != "" {
			branches = append(branches, branch)
		}
	}

	cmd = exec.CommandContext(ctx, "git", "branch", "-r", "--format=%(refname:short)")
	remoteOutput, err := cmd.Output()
	if err == nil {
		remoteBranches := strings.Split(strings.TrimSpace(string(remoteOutput)), "\n")
		for _, branch := range remoteBranches {
			if branch = strings.TrimSpace(branch); branch != "" && strings.HasPrefix(branch, "origin/") {
				remoteBranch := strings.TrimPrefix(branch, "origin/")
				if remoteBranch != "HEAD" && !contains(branches, remoteBranch) {
					branches = append(branches, remoteBranch)
				}
			}
		}
	}

	return branches, nil
}

func (m *Manager) CheckoutBranch(ctx context.Context, path, branch string) error {
	if !m.IsRepository(path) {
		return fmt.Errorf("diretório %s não é um repositório Git", path)
	}

	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("erro ao obter diretório atual: %w", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(path); err != nil {
		return fmt.Errorf("erro ao mudar para diretório do repositório: %w", err)
	}

	hasChanges, err := m.hasUncommittedChanges(ctx)
	if err != nil {
		m.logger.Warn("Não foi possível verificar mudanças pendentes", map[string]interface{}{
			"error": err.Error(),
		})
	} else if hasChanges {
		return fmt.Errorf("existem mudanças não commitadas. Commit ou stash as mudanças antes de trocar de branch")
	}

	cmd := exec.CommandContext(ctx, "git", "fetch")
	if output, err := cmd.CombinedOutput(); err != nil {
		m.logger.Warn("Erro ao fazer fetch", map[string]interface{}{
			"error":  err.Error(),
			"output": string(output),
		})
	}

	cmd = exec.CommandContext(ctx, "git", "checkout", branch)
	if output, err := cmd.CombinedOutput(); err != nil {
		cmd = exec.CommandContext(ctx, "git", "checkout", "-b", branch, "origin/"+branch)
		if output2, err2 := cmd.CombinedOutput(); err2 != nil {
			return fmt.Errorf("erro ao fazer checkout para branch %s: %w (output: %s), erro ao criar branch local: %v (output: %s)", branch, err, string(output), err2, string(output2))
		}
	}

	m.logger.Info("Checkout realizado com sucesso", map[string]interface{}{
		"branch": branch,
		"path":   path,
	})

	return nil
}

func (m *Manager) SyncBranch(ctx context.Context, path string) error {
	if !m.IsRepository(path) {
		return fmt.Errorf("diretório %s não é um repositório Git", path)
	}

	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("erro ao obter diretório atual: %w", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(path); err != nil {
		return fmt.Errorf("erro ao mudar para diretório do repositório: %w", err)
	}

	cmd := exec.CommandContext(ctx, "git", "fetch")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("erro ao executar git fetch: %w (output: %s)", err, string(output))
	}

	cmd = exec.CommandContext(ctx, "git", "pull")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("erro ao executar git pull: %w (output: %s)", err, string(output))
	}

	m.logger.Info("Branch sincronizada com sucesso", map[string]interface{}{
		"path": path,
	})

	return nil
}

func (m *Manager) hasUncommittedChanges(ctx context.Context) (bool, error) {
	cmd := exec.CommandContext(ctx, "git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("erro ao verificar status do git: %w", err)
	}

	return len(strings.TrimSpace(string(output))) > 0, nil
}

func (m *Manager) getRemoteURL(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("erro ao obter URL remota: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

func (m *Manager) getLastCommit(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "log", "-1", "--format=%h %s")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("erro ao obter último commit: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
