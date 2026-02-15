// Package proxy gerencia proxy reverso e configuração de rede.
package proxy

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/omelete/relief/pkg/fileutil"
	"github.com/omelete/relief/pkg/logger"
)

// HostsManager gerencia entradas no arquivo /etc/hosts
type HostsManager struct {
	hostsPath string
	logger    *logger.Logger
}

// NewHostsManager cria uma nova instância de HostsManager
func NewHostsManager(log *logger.Logger) *HostsManager {
	hostsPath := getHostsPath()
	return &HostsManager{
		hostsPath: hostsPath,
		logger:    log,
	}
}

// AddEntry adiciona uma entrada ao arquivo hosts
func (h *HostsManager) AddEntry(domain string) error {
	h.logger.Info("Adicionando entrada ao hosts", map[string]interface{}{
		"domain": domain,
	})

	// Verificar se já existe
	exists, err := h.HasEntry(domain)
	if err != nil {
		return fmt.Errorf("erro ao verificar entrada existente: %w", err)
	}

	if exists {
		h.logger.Debug("Entrada já existe no hosts", map[string]interface{}{
			"domain": domain,
		})
		return nil
	}

	// Ler arquivo atual
	content, err := os.ReadFile(h.hostsPath)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo hosts: %w", err)
	}

	// Adicionar nova entrada
	newEntry := fmt.Sprintf("127.0.0.1 %s # SOFREDOR\n", domain)
	
	// Procurar bloco SOFREDOR
	contentStr := string(content)
	if strings.Contains(contentStr, "# BEGIN SOFREDOR") {
		// Adicionar antes do END
		contentStr = strings.Replace(contentStr, "# END SOFREDOR", newEntry+"# END SOFREDOR", 1)
	} else {
		// Criar bloco
		contentStr += "\n# BEGIN SOFREDOR\n" + newEntry + "# END SOFREDOR\n"
	}

	// Escrever de volta
	if err := os.WriteFile(h.hostsPath, []byte(contentStr), 0644); err != nil {
		return fmt.Errorf("erro ao escrever arquivo hosts (permissão necessária): %w", err)
	}

	h.logger.Info("Entrada adicionada ao hosts com sucesso", map[string]interface{}{
		"domain": domain,
	})

	return nil
}

// RemoveEntry remove uma entrada do arquivo hosts
func (h *HostsManager) RemoveEntry(domain string) error {
	h.logger.Info("Removendo entrada do hosts", map[string]interface{}{
		"domain": domain,
	})

	// Ler arquivo atual
	content, err := os.ReadFile(h.hostsPath)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo hosts: %w", err)
	}

	// Remover linha com o domínio
	lines := strings.Split(string(content), "\n")
	newLines := []string{}
	
	for _, line := range lines {
		// Skip lines containing domain with RELIEF marker
		if strings.Contains(line, domain) && strings.Contains(line, "# RELIEF") {
			continue
		}
		newLines = append(newLines, line)
	}

	// Escrever de volta
	newContent := strings.Join(newLines, "\n")
	if err := os.WriteFile(h.hostsPath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("erro ao escrever arquivo hosts: %w", err)
	}

	h.logger.Info("Entrada removida do hosts", map[string]interface{}{
		"domain": domain,
	})

	return nil
}

// HasEntry verifica se uma entrada existe no hosts
func (h *HostsManager) HasEntry(domain string) (bool, error) {
	content, err := os.ReadFile(h.hostsPath)
	if err != nil {
		return false, fmt.Errorf("erro ao ler arquivo hosts: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.Contains(line, domain) && strings.Contains(line, "127.0.0.1") {
			return true, nil
		}
	}

	return false, nil
}

// ListEntries returns all Relief entries
func (h *HostsManager) ListEntries() ([]string, error) {
	content, err := os.ReadFile(h.hostsPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo hosts: %w", err)
	}

	entries := []string{}
	lines := strings.Split(string(content), "\n")
	
	for _, line := range lines {
		if strings.Contains(line, "# RELIEF") {
			// Extrair domínio
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				entries = append(entries, parts[1])
			}
		}
	}

	return entries, nil
}

// CleanupAll removes all Relief entries
func (h *HostsManager) CleanupAll() error {
	h.logger.Info("Limpando todas as entradas do hosts", nil)

	// Ler arquivo atual
	content, err := os.ReadFile(h.hostsPath)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo hosts: %w", err)
	}

	// Remove RELIEF block
	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")
	newLines := []string{}
	inReliefBlock := false

	for _, line := range lines {
		if strings.Contains(line, "# BEGIN RELIEF") {
			inReliefBlock = true
			continue
		}
		if strings.Contains(line, "# END RELIEF") {
			inReliefBlock = false
			continue
		}
		if !inReliefBlock {
			newLines = append(newLines, line)
		}
	}

	// Escrever de volta
	newContent := strings.Join(newLines, "\n")
	if err := os.WriteFile(h.hostsPath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("erro ao escrever arquivo hosts: %w", err)
	}

	h.logger.Info("Todas as entradas removidas do hosts", nil)
	return nil
}

// RequiresElevation verifica se é necessário privilégios elevados
func (h *HostsManager) RequiresElevation() bool {
	// Testar se podemos escrever no arquivo hosts
	file, err := os.OpenFile(h.hostsPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return true // Precisa de elevação
	}
	file.Close()
	return false
}

// GetHostsPath retorna o caminho do arquivo hosts
func (h *HostsManager) GetHostsPath() string {
	return h.hostsPath
}

// getHostsPath retorna o caminho do arquivo hosts baseado no SO
func getHostsPath() string {
	switch runtime.GOOS {
	case "windows":
		return "C:\\Windows\\System32\\drivers\\etc\\hosts"
	case "darwin", "linux":
		return "/etc/hosts"
	default:
		return "/etc/hosts"
	}
}
