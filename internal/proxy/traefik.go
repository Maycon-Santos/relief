// Package proxy gerencia proxy reverso e configuração de rede.
package proxy

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/omelete/relief/internal/domain"
	"github.com/omelete/relief/pkg/fileutil"
	"github.com/omelete/relief/pkg/logger"
	"gopkg.in/yaml.v3"
)

// TraefikManager gerencia o Traefik como proxy reverso
type TraefikManager struct {
	configPath  string
	binaryPath  string
	process     *exec.Cmd
	httpPort    int
	httpsPort   int
	running     bool
	mu          sync.RWMutex
	logger      *logger.Logger
	projects    map[string]*domain.Project
}

// NewTraefikManager cria uma nova instância de TraefikManager
func NewTraefikManager(httpPort, httpsPort int, log *logger.Logger) (*TraefikManager, error) {
	// Config directory: ~/.relief/traefik
	traefikDir, err := fileutil.GetReliefSubDir("traefik")
	if err != nil {
		return nil, fmt.Errorf("erro ao criar diretório traefik: %w", err)
	}

	configPath := filepath.Join(traefikDir, "dynamic.yaml")
	
	// Binary: ~/.relief/bin/traefik
	binDir, err := fileutil.GetReliefSubDir("bin")
	if err != nil {
		return nil, fmt.Errorf("erro ao criar diretório bin: %w", err)
	}
	binaryPath := filepath.Join(binDir, "traefik")

	return &TraefikManager{
		configPath:  configPath,
		binaryPath:  binaryPath,
		httpPort:    httpPort,
		httpsPort:   httpsPort,
		running:     false,
		logger:      log,
		projects:    make(map[string]*domain.Project),
	}, nil
}

// Start inicia o Traefik
func (t *TraefikManager) Start(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.running {
		return fmt.Errorf("traefik já está rodando")
	}

	// Verificar se binário existe
	if !fileutil.Exists(t.binaryPath) {
		return fmt.Errorf("binário do Traefik não encontrado em %s - execute o comando de instalação primeiro", t.binaryPath)
	}

	// Gerar configuração inicial
	if err := t.generateConfig(); err != nil {
		return fmt.Errorf("erro ao gerar configuração: %w", err)
	}

	// TODO: Iniciar Traefik
	// cmd := exec.CommandContext(ctx, t.binaryPath, 
	//     "--providers.file.filename=" + t.configPath,
	//     "--entrypoints.web.address=:"+strconv.Itoa(t.httpPort))

	t.logger.Info("Traefik iniciado", map[string]interface{}{
		"http_port":  t.httpPort,
		"https_port": t.httpsPort,
	})

	t.running = true
	return nil
}

// Stop para o Traefik
func (t *TraefikManager) Stop() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.running {
		return nil
	}

	if t.process != nil && t.process.Process != nil {
		if err := t.process.Process.Kill(); err != nil {
			return fmt.Errorf("erro ao parar traefik: %w", err)
		}
	}

	t.running = false
	t.logger.Info("Traefik parado", nil)

	return nil
}

// AddProject adiciona um projeto ao Traefik
func (t *TraefikManager) AddProject(project *domain.Project) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if project.Domain == "" {
		return fmt.Errorf("projeto não tem domínio configurado")
	}

	t.projects[project.ID] = project

	// Regenerar configuração
	if err := t.generateConfig(); err != nil {
		return fmt.Errorf("erro ao regenerar configuração: %w", err)
	}

	t.logger.Info("Projeto adicionado ao Traefik", map[string]interface{}{
		"project": project.Name,
		"domain":  project.Domain,
		"port":    project.Port,
	})

	return nil
}

// RemoveProject remove um projeto do Traefik
func (t *TraefikManager) RemoveProject(projectID string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	delete(t.projects, projectID)

	// Regenerar configuração
	if err := t.generateConfig(); err != nil {
		return fmt.Errorf("erro ao regenerar configuração: %w", err)
	}

	t.logger.Info("Projeto removido do Traefik", map[string]interface{}{
		"project_id": projectID,
	})

	return nil
}

// generateConfig gera o arquivo de configuração dinâmica do Traefik
func (t *TraefikManager) generateConfig() error {
	config := TraefikConfig{
		HTTP: HTTPConfig{
			Routers:  make(map[string]Router),
			Services: make(map[string]Service),
		},
	}

	// Adicionar rotas para cada projeto
	for _, project := range t.projects {
		routerName := fmt.Sprintf("%s-router", project.Name)
		serviceName := fmt.Sprintf("%s-service", project.Name)

		config.HTTP.Routers[routerName] = Router{
			Rule:    fmt.Sprintf("Host(`%s`)", project.Domain),
			Service: serviceName,
		}

		config.HTTP.Services[serviceName] = Service{
			LoadBalancer: LoadBalancer{
				Servers: []Server{
					{
						URL: fmt.Sprintf("http://localhost:%d", project.Port),
					},
				},
			},
		}
	}

	// Serializar para YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("erro ao serializar configuração: %w", err)
	}

	// Escrever arquivo
	if err := os.WriteFile(t.configPath, data, 0644); err != nil {
		return fmt.Errorf("erro ao escrever configuração: %w", err)
	}

	t.logger.Debug("Configuração do Traefik gerada", map[string]interface{}{
		"path":     t.configPath,
		"projects": len(t.projects),
	})

	return nil
}

// IsRunning verifica se o Traefik está rodando
func (t *TraefikManager) IsRunning() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.running
}

// GetConfigPath retorna o path do arquivo de configuração
func (t *TraefikManager) GetConfigPath() string {
	return t.configPath
}

// InstallTraefik baixa e instala o binário do Traefik
func (t *TraefikManager) InstallTraefik(ctx context.Context, version string) error {
	t.logger.Info("Instalando Traefik", map[string]interface{}{
		"version": version,
	})

	// TODO: Implementar download real do Traefik
	// - Detectar SO e arquitetura
	// - Baixar da URL oficial: https://github.com/traefik/traefik/releases/download/v{version}/traefik_{os}_{arch}.tar.gz
	// - Extrair para t.binaryPath
	// - Dar permissão de execução

	return fmt.Errorf("instalação automática do Traefik ainda não implementada - baixe de https://github.com/traefik/traefik/releases e coloque em %s", t.binaryPath)
}
