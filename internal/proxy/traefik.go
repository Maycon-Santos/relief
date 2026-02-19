package proxy

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/Maycon-Santos/relief/internal/domain"
	"github.com/Maycon-Santos/relief/pkg/fileutil"
	"github.com/Maycon-Santos/relief/pkg/logger"
	"gopkg.in/yaml.v3"
)

type TraefikManager struct {
	configPath string
	binaryPath string
	process    *exec.Cmd
	httpPort   int
	httpsPort  int
	running    bool
	mu         sync.RWMutex
	logger     *logger.Logger
	projects   map[string]*domain.Project
}

func NewTraefikManager(httpPort, httpsPort int, log *logger.Logger) (*TraefikManager, error) {
	traefikDir, err := fileutil.GetReliefSubDir("traefik")
	if err != nil {
		return nil, fmt.Errorf("erro ao criar diretório traefik: %w", err)
	}

	configPath := filepath.Join(traefikDir, "dynamic.yaml")

	binDir, err := fileutil.GetReliefSubDir("bin")
	if err != nil {
		return nil, fmt.Errorf("erro ao criar diretório bin: %w", err)
	}
	binaryPath := filepath.Join(binDir, "traefik")

	return &TraefikManager{
		configPath: configPath,
		binaryPath: binaryPath,
		httpPort:   httpPort,
		httpsPort:  httpsPort,
		running:    false,
		logger:     log,
		projects:   make(map[string]*domain.Project),
	}, nil
}

func (t *TraefikManager) Start(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.running {
		return fmt.Errorf("traefik já está rodando")
	}

	if !fileutil.Exists(t.binaryPath) {
		t.logger.Info("Traefik não encontrado, instalando automaticamente...", nil)
		if err := t.InstallTraefik(ctx, "v3.0.0"); err != nil {
			return fmt.Errorf("erro ao instalar Traefik: %w", err)
		}
	}

	if err := t.generateConfig(); err != nil {
		return fmt.Errorf("erro ao gerar configuração: %w", err)
	}

	logDir := filepath.Dir(t.configPath)
	logFile := filepath.Join(logDir, "traefik.log")

	cmd := exec.CommandContext(ctx, t.binaryPath,
		"--providers.file.filename="+t.configPath,
		"--entrypoints.web.address=:"+fmt.Sprintf("%d", t.httpPort),
		"--log.level=INFO",
		"--log.filepath="+logFile,
		"--accesslog=false",
	)

	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("erro ao iniciar Traefik: %w", err)
	}

	t.process = cmd

	t.logger.Info("Traefik iniciado", map[string]interface{}{
		"http_port":  t.httpPort,
		"https_port": t.httpsPort,
		"pid":        cmd.Process.Pid,
		"config":     t.configPath,
	})

	t.running = true
	return nil
}

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

func (t *TraefikManager) AddProject(project *domain.Project) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if project.Domain == "" {
		return fmt.Errorf("projeto não tem domínio configurado")
	}

	t.projects[project.ID] = project

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

func (t *TraefikManager) RemoveProject(projectID string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	delete(t.projects, projectID)

	if err := t.generateConfig(); err != nil {
		return fmt.Errorf("erro ao regenerar configuração: %w", err)
	}

	t.logger.Info("Projeto removido do Traefik", map[string]interface{}{
		"project_id": projectID,
	})

	return nil
}

func (t *TraefikManager) generateConfig() error {
	config := TraefikConfig{
		HTTP: HTTPConfig{
			Routers:  make(map[string]Router),
			Services: make(map[string]Service),
		},
	}

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

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("erro ao serializar configuração: %w", err)
	}

	if err := os.WriteFile(t.configPath, data, 0644); err != nil {
		return fmt.Errorf("erro ao escrever configuração: %w", err)
	}

	t.logger.Debug("Configuração do Traefik gerada", map[string]interface{}{
		"path":     t.configPath,
		"projects": len(t.projects),
	})

	return nil
}

func (t *TraefikManager) IsRunning() bool {
	addr := fmt.Sprintf(":%d", t.httpPort)
	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		t.mu.Lock()
		t.running = false
		t.mu.Unlock()
		return false
	}
	conn.Close()
	return true
}

func (t *TraefikManager) Restart(ctx context.Context) error {
	t.logger.Info("Reiniciando Traefik", nil)
	_ = t.Stop()
	return t.Start(ctx)
}

func (t *TraefikManager) GetConfigPath() string {
	return t.configPath
}

func (t *TraefikManager) InstallTraefik(ctx context.Context, version string) error {
	t.logger.Info("Instalando Traefik", map[string]interface{}{
		"version": version,
	})

	osName := runtime.GOOS
	arch := runtime.GOARCH

	filename := fmt.Sprintf("traefik_%s_%s_%s.tar.gz", version, osName, arch)
	url := fmt.Sprintf("https://github.com/traefik/traefik/releases/download/%s/%s", version, filename)

	t.logger.Info("Baixando Traefik", map[string]interface{}{
		"url":  url,
		"os":   osName,
		"arch": arch,
	})

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("erro ao criar request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao baixar Traefik: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("erro ao baixar Traefik: status %d", resp.StatusCode)
	}

	gzr, err := gzip.NewReader(resp.Body)
	if err != nil {
		return fmt.Errorf("erro ao descompactar gzip: %w", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("erro ao ler tar: %w", err)
		}

		if header.Name == "traefik" || header.Name == "traefik.exe" {
			destFile, err := os.OpenFile(t.binaryPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
			if err != nil {
				return fmt.Errorf("erro ao criar arquivo: %w", err)
			}

			if _, err := io.Copy(destFile, tr); err != nil {
				destFile.Close()
				return fmt.Errorf("erro ao copiar binário: %w", err)
			}

			destFile.Close()

			t.logger.Info("Traefik instalado com sucesso", map[string]interface{}{
				"path": t.binaryPath,
			})

			return nil
		}
	}

	return fmt.Errorf("binário do Traefik não encontrado no arquivo tar")
}
