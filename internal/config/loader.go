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

func (l *Loader) LoadConfig(remoteURL, localPath string) (*Config, error) {
	var finalConfig *Config

	if remoteURL != "" {
		remoteConfig, err := l.loadRemoteConfig(remoteURL)
		if err != nil {
			fmt.Printf("Aviso: não foi possível carregar config remota: %v\n", err)
		} else {
			finalConfig = remoteConfig
		}
	}

	globalPath, err := GetGlobalConfigPath()
	if err == nil && fileutil.Exists(globalPath) {
		globalConfig, err := l.loadLocalConfig(globalPath)
		if err != nil {
			fmt.Printf("Aviso: não foi possível carregar config global: %v\n", err)
		} else {
			if finalConfig == nil {
				finalConfig = globalConfig
			} else {
				finalConfig.MergeWith(globalConfig)
			}
		}
	}

	if fileutil.Exists(localPath) {
		localConfig, err := l.loadLocalConfig(localPath)
		if err != nil {
			return nil, fmt.Errorf("erro ao carregar config local: %w", err)
		}

		if finalConfig == nil {
			finalConfig = localConfig
		} else {
			finalConfig.MergeWith(localConfig)
		}
	}

	if finalConfig == nil {
		finalConfig = defaultConfig()
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

func GetLocalConfigPath() (string, error) {
	reliefDir, err := fileutil.GetReliefDir()
	if err != nil {
		return "", err
	}
	return reliefDir + "/config.local.yaml", nil
}

func GetGlobalConfigPath() (string, error) {
	reliefDir, err := fileutil.GetReliefDir()
	if err != nil {
		return "", err
	}
	return reliefDir + "/config.global.yaml", nil
}
