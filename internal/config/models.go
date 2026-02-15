// Package config gerencia a configuração da aplicação.
package config

import "time"

// Config é a estrutura principal de configuração
type Config struct {
	Remote   RemoteConfig            `yaml:"remote"`
	Projects []ProjectConfig         `yaml:"projects"`
	Tools    map[string]ToolVersion  `yaml:"tools"`
	Proxy    ProxyConfig             `yaml:"proxy"`
}

// RemoteConfig contém configurações para carregar config remota
type RemoteConfig struct {
	URL             string        `yaml:"url"`
	RefreshInterval time.Duration `yaml:"refresh_interval"`
	Enabled         bool          `yaml:"enabled"`
}

// ProjectConfig define a configuração de um projeto
type ProjectConfig struct {
	Name         string            `yaml:"name"`
	Path         string            `yaml:"path"`
	Domain       string            `yaml:"domain"`
	Type         string            `yaml:"type"` // node, python, docker, java
	Dependencies []DependencySpec  `yaml:"dependencies"`
	Scripts      map[string]string `yaml:"scripts"`
	Env          map[string]string `yaml:"env"`
	Port         int               `yaml:"port,omitempty"`
	AutoStart    bool              `yaml:"auto_start"`
}

// DependencySpec especifica uma dependência de um projeto
type DependencySpec struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Managed bool   `yaml:"managed"` // Se o orquestrador deve prover a dependência
}

// ToolVersion especifica a versão de uma ferramenta
type ToolVersion struct {
	Version     string `yaml:"version"`
	DownloadURL string `yaml:"download_url,omitempty"`
}

// ProxyConfig contém configurações do proxy (Traefik)
type ProxyConfig struct {
	HTTPPort   int  `yaml:"http_port"`
	HTTPSPort  int  `yaml:"https_port"`
	Dashboard  bool `yaml:"dashboard"`
	AutoManage bool `yaml:"auto_manage"` // Gerenciar automaticamente o Traefik
}

// Validate valida a configuração
func (c *Config) Validate() error {
	// Validações básicas
	if c.Proxy.HTTPPort <= 0 {
		c.Proxy.HTTPPort = 80
	}
	if c.Proxy.HTTPSPort <= 0 {
		c.Proxy.HTTPSPort = 443
	}

	// Validar projetos
	for i := range c.Projects {
		if c.Projects[i].Name == "" {
			return &ValidationError{Field: "projects[].name", Message: "nome do projeto é obrigatório"}
		}
		if c.Projects[i].Type == "" {
			c.Projects[i].Type = "docker" // Padrão
		}
	}

	return nil
}

// ValidationError representa um erro de validação
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// GetProjectByName retorna um projeto por nome
func (c *Config) GetProjectByName(name string) *ProjectConfig {
	for i := range c.Projects {
		if c.Projects[i].Name == name {
			return &c.Projects[i]
		}
	}
	return nil
}

// GetProjectByDomain retorna um projeto por domínio
func (c *Config) GetProjectByDomain(domain string) *ProjectConfig {
	for i := range c.Projects {
		if c.Projects[i].Domain == domain {
			return &c.Projects[i]
		}
	}
	return nil
}

// MergeWith faz merge desta config com outra (override)
func (c *Config) MergeWith(other *Config) {
	// Merge de projetos (por nome)
	for _, otherProject := range other.Projects {
		found := false
		for i := range c.Projects {
			if c.Projects[i].Name == otherProject.Name {
				// Override do projeto existente
				c.Projects[i] = otherProject
				found = true
				break
			}
		}
		if !found {
			// Adicionar novo projeto
			c.Projects = append(c.Projects, otherProject)
		}
	}

	// Merge de tools
	if other.Tools != nil {
		if c.Tools == nil {
			c.Tools = make(map[string]ToolVersion)
		}
		for name, version := range other.Tools {
			c.Tools[name] = version
		}
	}

	// Override de proxy config se especificado
	if other.Proxy.HTTPPort != 0 {
		c.Proxy.HTTPPort = other.Proxy.HTTPPort
	}
	if other.Proxy.HTTPSPort != 0 {
		c.Proxy.HTTPSPort = other.Proxy.HTTPSPort
	}
	if other.Proxy.Dashboard {
		c.Proxy.Dashboard = other.Proxy.Dashboard
	}

	// Override de remote config
	if other.Remote.URL != "" {
		c.Remote = other.Remote
	}
}
