// Package git provides Git repository management functionality.
package git

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/relief-org/relief/internal/domain"
	"github.com/relief-org/relief/pkg/fileutil"
	"github.com/relief-org/relief/pkg/logger"
)

// Manager gerencia operações Git para projetos
type Manager struct {
	logger *logger.Logger
}

// NewManager cria uma nova instância de Manager
func NewManager(log *logger.Logger) *Manager {
	return &Manager{
		logger: log,
	}
}

// CloneOrUpdate clona um repositório ou atualiza se já existir
func (m *Manager) CloneOrUpdate(ctx context.Context, repoURL, targetPath, branch string) error {
	// Verificar se o diretório já existe
	if fileutil.Exists(targetPath) {
		// Verificar se é um repositório Git
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

	// Criar diretório pai se não existir
	parentDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório pai: %w", err)
	}

	// Clonar repositório
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

// updateRepository atualiza um repositório existente
func (m *Manager) updateRepository(ctx context.Context, path, branch string) error {
	// Mudar para o diretório do repositório
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("erro ao obter diretório atual: %w", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(path); err != nil {
		return fmt.Errorf("erro ao mudar para diretório do repositório: %w", err)
	}

	// Fazer fetch
	cmd := exec.CommandContext(ctx, "git", "fetch")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("erro ao executar git fetch: %w (output: %s)", err, string(output))
	}

	// Trocar de branch se necessário
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

	// Fazer pull
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

// getCurrentBranch obtém o branch atual
func (m *Manager) getCurrentBranch(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("erro ao obter branch atual: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// IsRepository verifica se um diretório é um repositório Git
func (m *Manager) IsRepository(path string) bool {
	return fileutil.Exists(filepath.Join(path, ".git"))
}

// GetRemoteURL obtém a URL remota do repositório
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

// GetGitInfo obtém informações completas do git para um projeto
func (m *Manager) GetGitInfo(ctx context.Context, path string) (*domain.GitInfo, error) {
	gitInfo := &domain.GitInfo{
		IsRepository: m.IsRepository(path),
	}

	if !gitInfo.IsRepository {
		return gitInfo, nil
	}

	// Salvar diretório atual
	originalDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter diretório atual: %w", err)
	}
	defer os.Chdir(originalDir)

	// Mudar para o diretório do repositório
	if err := os.Chdir(path); err != nil {
		return nil, fmt.Errorf("erro ao mudar para diretório do repositório: %w", err)
	}

	// Obter branch atual
	if currentBranch, err := m.getCurrentBranch(ctx); err == nil {
		gitInfo.CurrentBranch = currentBranch
	}

	// Obter branches disponíveis
	if branches, err := m.getBranches(ctx); err == nil {
		gitInfo.AvailableBranches = branches
	}

	// Obter URL remota
	if remoteURL, err := m.getRemoteURL(ctx); err == nil {
		gitInfo.RemoteURL = remoteURL
	}

	// Verificar se há mudanças pendentes
	if hasChanges, err := m.hasUncommittedChanges(ctx); err == nil {
		gitInfo.HasChanges = hasChanges
	}

	// Obter último commit
	if lastCommit, err := m.getLastCommit(ctx); err == nil {
		gitInfo.LastCommit = lastCommit
	}

	return gitInfo, nil
}

// getBranches obtém todas as branches disponíveis
func (m *Manager) getBranches(ctx context.Context) ([]string, error) {
	// Branches locais
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

	// Branches remotas (sem duplicatas)
	cmd = exec.CommandContext(ctx, "git", "branch", "-r", "--format=%(refname:short)")
	remoteOutput, err := cmd.Output()
	if err == nil {
		remoteBranches := strings.Split(strings.TrimSpace(string(remoteOutput)), "\n")
		for _, branch := range remoteBranches {
			if branch = strings.TrimSpace(branch); branch != "" && strings.HasPrefix(branch, "origin/") {
				// Remove o prefixo "origin/"
				remoteBranch := strings.TrimPrefix(branch, "origin/")
				if remoteBranch != "HEAD" && !contains(branches, remoteBranch) {
					branches = append(branches, remoteBranch)
				}
			}
		}
	}

	return branches, nil
}

// CheckoutBranch faz checkout para uma branch específica
func (m *Manager) CheckoutBranch(ctx context.Context, path, branch string) error {
	if !m.IsRepository(path) {
		return fmt.Errorf("diretório %s não é um repositório Git", path)
	}

	// Salvar diretório atual
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("erro ao obter diretório atual: %w", err)
	}
	defer os.Chdir(originalDir)

	// Mudar para o diretório do repositório
	if err := os.Chdir(path); err != nil {
		return fmt.Errorf("erro ao mudar para diretório do repositório: %w", err)
	}

	// Verificar se há mudanças não commitadas
	hasChanges, err := m.hasUncommittedChanges(ctx)
	if err != nil {
		m.logger.Warn("Não foi possível verificar mudanças pendentes", map[string]interface{}{
			"error": err.Error(),
		})
	} else if hasChanges {
		return fmt.Errorf("existem mudanças não commitadas. Commit ou stash as mudanças antes de trocar de branch")
	}

	// Fazer fetch para garantir que temos as últimas informações
	cmd := exec.CommandContext(ctx, "git", "fetch")
	if output, err := cmd.CombinedOutput(); err != nil {
		m.logger.Warn("Erro ao fazer fetch", map[string]interface{}{
			"error":  err.Error(),
			"output": string(output),
		})
	}

	// Tentar checkout
	cmd = exec.CommandContext(ctx, "git", "checkout", branch)
	if output, err := cmd.CombinedOutput(); err != nil {
		// Se falhar, talvez seja uma branch remota que precisa ser criada localmente
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

// SyncBranch faz pull da branch atual
func (m *Manager) SyncBranch(ctx context.Context, path string) error {
	if !m.IsRepository(path) {
		return fmt.Errorf("diretório %s não é um repositório Git", path)
	}

	// Salvar diretório atual
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("erro ao obter diretório atual: %w", err)
	}
	defer os.Chdir(originalDir)

	// Mudar para o diretório do repositório
	if err := os.Chdir(path); err != nil {
		return fmt.Errorf("erro ao mudar para diretório do repositório: %w", err)
	}

	// Fazer fetch primeiro
	cmd := exec.CommandContext(ctx, "git", "fetch")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("erro ao executar git fetch: %w (output: %s)", err, string(output))
	}

	// Fazer pull
	cmd = exec.CommandContext(ctx, "git", "pull")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("erro ao executar git pull: %w (output: %s)", err, string(output))
	}

	m.logger.Info("Branch sincronizada com sucesso", map[string]interface{}{
		"path": path,
	})

	return nil
}

// hasUncommittedChanges verifica se há mudanças não commitadas
func (m *Manager) hasUncommittedChanges(ctx context.Context) (bool, error) {
	cmd := exec.CommandContext(ctx, "git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("erro ao verificar status do git: %w", err)
	}

	return len(strings.TrimSpace(string(output))) > 0, nil
}

// getRemoteURL obtém a URL remota do repositório (helper interno)
func (m *Manager) getRemoteURL(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("erro ao obter URL remota: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// getLastCommit obtém o hash e mensagem do último commit
func (m *Manager) getLastCommit(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "log", "-1", "--format=%h %s")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("erro ao obter último commit: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// contains verifica se um slice contém um item
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
