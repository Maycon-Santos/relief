// Package config gerencia a configuração da aplicação.
package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/omelete/relief/pkg/fileutil"
	"github.com/omelete/relief/pkg/httputil"
	"gopkg.in/yaml.v3"
)

// Loader é responsável por carregar e fazer merge de configurações
type Loader struct {
	httpClient *httputil.Client
}

// NewLoader cria uma nova instância de Loader
func NewLoader() *Loader {
	return &Loader{
		httpClient: httputil.NewClient(10 * time.Second),
	}
}

// LoadConfig carrega a configuração completa (remote + local merge)
func (l *Loader) LoadConfig(remoteURL, localPath string) (*Config, error) {
	var finalConfig *Config

	// 1. Carregar configuração remota (se URL fornecida)
	if remoteURL != "" {
		remoteConfig, err := l.loadRemoteConfig(remoteURL)
		if err != nil {
			// Não falhar se remote não disponível, apenas logar
			fmt.Printf("Aviso: não foi possível carregar config remota: %v\n", err)
		} else {
			finalConfig = remoteConfig
		}
	}

	// 2. Carregar configuração local (se existir)
	if fileutil.Exists(localPath) {
		localConfig, err := l.loadLocalConfig(localPath)
		if err != nil {
			return nil, fmt.Errorf("erro ao carregar config local: %w", err)
		}

		if finalConfig == nil {
			finalConfig = localConfig
		} else {
			// Fazer merge: local sobrescreve remote
			finalConfig.MergeWith(localConfig)
		}
	}

	// 3. Se não houver nenhuma config, usar defaults
	if finalConfig == nil {
		finalConfig = defaultConfig()
	}

	// 4. Validar configuração final
	if err := finalConfig.Validate(); err != nil {
		return nil, fmt.Errorf("erro ao validar configuração: %w", err)
	}

	return finalConfig, nil
}

// loadRemoteConfig baixa e faz parse da config remota
func (l *Loader) loadRemoteConfig(url string) (*Config, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Download do YAML
	data, err := l.httpClient.Get(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("erro ao baixar config remota: %w", err)
	}

	// Parse do YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse da config remota: %w", err)
	}

	return &config, nil
}

// loadLocalConfig carrega e faz parse da config local
func (l *Loader) loadLocalConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo local: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse da config local: %w", err)
	}

	return &config, nil
}

// SaveConfig salva a configuração em um arquivo
func (l *Loader) SaveConfig(config *Config, path string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("erro ao serializar config: %w", err)
	}

	if err := fileutil.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("erro ao salvar config: %w", err)
	}

	return nil
}

// defaultConfig retorna uma configuração padrão
func defaultConfig() *Config {
	return &Config{
		Remote: RemoteConfig{
			Enabled:         false,
			RefreshInterval: 1 * time.Hour,
		},
		Projects: []ProjectConfig{},
		Tools: map[string]ToolVersion{
			"node": {
				Version: "18.19.0",
			},
			"traefik": {
				Version: "2.10.7",
			},
		},
		Proxy: ProxyConfig{
			HTTPPort:   80,
			HTTPSPort:  443,
			Dashboard:  true,
			AutoManage: true,
		},
	}
}

// GetConfigPath retorna o caminho do arquivo de configuração principal
func GetConfigPath() (string, error) {
	reliefDir, err := fileutil.GetReliefDir()
	if err != nil {
		return "", err
	}
	return reliefDir + "/config.yaml", nil
}

// GetLocalConfigPath retorna o caminho do arquivo de configuração local (override)
func GetLocalConfigPath() (string, error) {
	reliefDir, err := fileutil.GetReliefDir()
	if err != nil {
		return "", err
	}
	return reliefDir + "/config.local.yaml", nil
}
