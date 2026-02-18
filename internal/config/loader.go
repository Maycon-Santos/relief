package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Maycon-Santos/relief/pkg/fileutil"
	"github.com/Maycon-Santos/relief/pkg/httputil"
	"gopkg.in/yaml.v3"
)

type Loader struct {
	httpClient *httputil.Client
}

func NewLoader() *Loader {
	return &Loader{
		httpClient: httputil.NewClient(10 * time.Second),
	}
}

func (l *Loader) LoadConfig(remoteURL, configPath string) (*Config, error) {
	var finalConfig *Config

	if remoteURL != "" {
		remoteConfig, err := l.loadRemoteConfig(remoteURL)
		if err != nil {
			fmt.Printf("Aviso: não foi possível carregar config remota: %v\n", err)
		} else {
			finalConfig = remoteConfig
			fmt.Printf("Configuração remota carregada de: %s\n", remoteURL)
		}
	}

	if fileutil.Exists(configPath) {
		localConfig, err := l.loadLocalConfig(configPath)
		if err != nil {
			return nil, fmt.Errorf("erro ao carregar config: %w", err)
		}

		if finalConfig == nil {
			finalConfig = localConfig
		} else {
			finalConfig.MergeWith(localConfig)
		}
		fmt.Printf("Configuração local carregada de: %s\n", configPath)
	}

	if finalConfig == nil {
		fmt.Println("Usando configuração padrão")
		finalConfig = defaultConfig()
	}

	if finalConfig.Environment.ExternalWorkspaceConfig != "" {
		if err := l.loadExternalWorkspaceProjects(finalConfig); err != nil {
			fmt.Printf("Aviso: não foi possível carregar projetos do workspace externo: %v\n", err)
		} else {
			fmt.Printf("Projetos carregados do workspace externo: %s\n", finalConfig.Environment.ExternalWorkspaceConfig)
		}
	}

	if err := finalConfig.Validate(); err != nil {
		return nil, fmt.Errorf("erro ao validar configuração: %w", err)
	}

	return finalConfig, nil
}

func (l *Loader) loadRemoteConfig(url string) (*Config, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	data, err := l.httpClient.Get(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("erro ao baixar config remota: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse da config remota: %w", err)
	}

	return &config, nil
}

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

func (l *Loader) loadExternalWorkspaceProjects(cfg *Config) error {
	externalPath := cfg.Environment.ExternalWorkspaceConfig

	fmt.Printf("🔍 Tentando carregar projetos externos...\n")
	fmt.Printf("   Caminho configurado: %s\n", externalPath)
	fmt.Printf("   WorkspacePath: %s\n", cfg.Environment.WorkspacePath)

	if cfg.Environment.WorkspacePath != "" {
		externalPath = cfg.Environment.WorkspacePath + "/" + externalPath
		fmt.Printf("   Caminho completo: %s\n", externalPath)
	}

	if len(externalPath) > 0 && externalPath[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			externalPath = homeDir + externalPath[1:]
			fmt.Printf("   Caminho expandido: %s\n", externalPath)
		}
	}

	if !fileutil.Exists(externalPath) {
		return fmt.Errorf("arquivo de workspace externo não encontrado: %s", externalPath)
	}

	fmt.Printf("   ✅ Arquivo encontrado!\n")

	data, err := os.ReadFile(externalPath)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo de workspace externo: %w", err)
	}

	var externalConfig Config
	if err := yaml.Unmarshal(data, &externalConfig); err != nil {
		return fmt.Errorf("erro ao fazer parse do arquivo de workspace externo: %w", err)
	}

	if len(externalConfig.Projects) > 0 {
		fmt.Printf("Carregando %d projetos do workspace externo: %s\n", len(externalConfig.Projects), externalPath)
		cfg.Projects = append(cfg.Projects, externalConfig.Projects...)
	}

	return nil
}

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

func GetConfigPath() (string, error) {
	reliefDir, err := fileutil.GetReliefDir()
	if err != nil {
		return "", err
	}
	return reliefDir + "/config.yaml", nil
}
