package proxy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/relief-org/relief/pkg/fileutil"
	"github.com/relief-org/relief/pkg/logger"
)

type HostsManager struct {
	hostsPath string
	logger    *logger.Logger
}

func NewHostsManager(log *logger.Logger) *HostsManager {
	hostsPath := getHostsPath()
	return &HostsManager{
		hostsPath: hostsPath,
		logger:    log,
	}
}

func (h *HostsManager) AddEntry(domain string) error {
	h.logger.Info("Adicionando entrada ao hosts", map[string]interface{}{
		"domain": domain,
	})

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

	content, err := os.ReadFile(h.hostsPath)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo hosts: %w", err)
	}

	newEntry := fmt.Sprintf("127.0.0.1 %s # RELIEF\n", domain)

	contentStr := string(content)
	if strings.Contains(contentStr, "# BEGIN RELIEF") {
		contentStr = strings.Replace(contentStr, "# END RELIEF", newEntry+"# END RELIEF", 1)
	} else {
		contentStr += "\n# BEGIN RELIEF\n" + newEntry + "# END RELIEF\n"
	}

	if err := os.WriteFile(h.hostsPath, []byte(contentStr), 0644); err != nil {
		h.logger.Warn("Sem permissão para escrever em /etc/hosts, tentando com sudo...", nil)
		if err := h.writeWithSudo(contentStr); err != nil {
			return fmt.Errorf("erro ao escrever arquivo hosts: %w", err)
		}
	}

	h.logger.Info("Entrada adicionada ao hosts com sucesso", map[string]interface{}{
		"domain": domain,
	})

	return nil
}

func (h *HostsManager) RemoveEntry(domain string) error {
	h.logger.Info("Removendo entrada do hosts", map[string]interface{}{
		"domain": domain,
	})

	content, err := os.ReadFile(h.hostsPath)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo hosts: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	newLines := []string{}

	for _, line := range lines {
		if strings.Contains(line, domain) && strings.Contains(line, "# RELIEF") {
			continue
		}
		newLines = append(newLines, line)
	}

	newContent := strings.Join(newLines, "\n")
	if err := os.WriteFile(h.hostsPath, []byte(newContent), 0644); err != nil {
		h.logger.Warn("Sem permissão para escrever em /etc/hosts, tentando com sudo...", nil)
		if err := h.writeWithSudo(newContent); err != nil {
			return fmt.Errorf("erro ao escrever arquivo hosts: %w", err)
		}
	}

	h.logger.Info("Entrada removida do hosts", map[string]interface{}{
		"domain": domain,
	})

	return nil
}

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

func (h *HostsManager) ListEntries() ([]string, error) {
	content, err := os.ReadFile(h.hostsPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo hosts: %w", err)
	}

	entries := []string{}
	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		if strings.Contains(line, "# RELIEF") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				entries = append(entries, parts[1])
			}
		}
	}

	return entries, nil
}

func (h *HostsManager) CleanupAll() error {
	h.logger.Info("Limpando todas as entradas do hosts", nil)

	content, err := os.ReadFile(h.hostsPath)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo hosts: %w", err)
	}

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

	newContent := strings.Join(newLines, "\n")
	if err := os.WriteFile(h.hostsPath, []byte(newContent), 0644); err != nil {
		h.logger.Warn("Sem permissão para escrever em /etc/hosts, tentando com sudo...", nil)
		if err := h.writeWithSudo(newContent); err != nil {
			return fmt.Errorf("erro ao escrever arquivo hosts: %w", err)
		}
	}

	h.logger.Info("Todas as entradas removidas do hosts", nil)
	return nil
}

func (h *HostsManager) RequiresElevation() bool {
	file, err := os.OpenFile(h.hostsPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return true
	}
	file.Close()
	return false
}

func (h *HostsManager) GetHostsPath() string {
	return h.hostsPath
}

func (h *HostsManager) writeWithSudo(content string) error {
	reliefDir, err := fileutil.GetReliefDir()
	if err != nil {
		return fmt.Errorf("erro ao obter diretório relief: %w", err)
	}

	tempFile := filepath.Join(reliefDir, "hosts.tmp")
	if err := os.WriteFile(tempFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("erro ao criar arquivo temporário: %w", err)
	}
	defer os.Remove(tempFile)

	switch runtime.GOOS {
	case "darwin":
		script := fmt.Sprintf(`do shell script "cat %s > %s" with administrator privileges`, tempFile, h.hostsPath)
		cmd := exec.Command("osascript", "-e", script)

		h.logger.Info("Solicitando permissões administrativas...", nil)

		if output, err := cmd.CombinedOutput(); err != nil {
			h.logger.Error("Falha ao executar com sudo", err, map[string]interface{}{
				"output": string(output),
			})
			return fmt.Errorf("usuário cancelou ou erro ao executar com privilégios administrativos: %w", err)
		}

		h.logger.Info("Arquivo hosts atualizado com privilégios administrativos", nil)
		return nil

	case "linux":
		if _, err := exec.LookPath("pkexec"); err == nil {
			cmd := exec.Command("pkexec", "cp", tempFile, h.hostsPath)
			if output, err := cmd.CombinedOutput(); err != nil {
				h.logger.Error("Falha ao executar com pkexec", err, map[string]interface{}{
					"output": string(output),
				})
				return fmt.Errorf("erro ao executar com privilégios administrativos: %w", err)
			}
			return nil
		}

		return fmt.Errorf("permissão negada. Execute manualmente: sudo cp %s %s", tempFile, h.hostsPath)

	case "windows":
		return fmt.Errorf("elevação de privilégios no Windows não implementada. Execute como administrador")

	default:
		return fmt.Errorf("sistema operacional não suportado: %s", runtime.GOOS)
	}
}

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
