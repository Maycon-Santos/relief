package config

import "time"

type Config struct {
	Remote              RemoteConfig                 `yaml:"remote"`
	Projects            []ProjectConfig              `yaml:"projects"`
	Tools               map[string]ToolVersion       `yaml:"tools"`
	Proxy               ProxyConfig                  `yaml:"proxy"`
	ManagedDependencies map[string]ManagedDependency `yaml:"managed_dependencies"`
	Development         DevelopmentConfig            `yaml:"development"`
	Logging             LoggingConfig                `yaml:"logging"`
	HealthChecks        map[string]HealthCheckConfig `yaml:"health_checks"`
	Environment         EnvironmentConfig            `yaml:"environment"`
}

type EnvironmentConfig struct {
	CompanyName             string `yaml:"company_name"`
	WorkspacePath           string `yaml:"workspace_path"`
	ExternalWorkspaceConfig string `yaml:"external_workspace_config"`
}

type RemoteConfig struct {
	URL             string        `yaml:"url"`
	RefreshInterval time.Duration `yaml:"refresh_interval"`
	Enabled         bool          `yaml:"enabled"`
}

type ProjectConfig struct {
	Name         string            `yaml:"name"`
	Path         string            `yaml:"path"`
	Repository   *RepositoryConfig `yaml:"repository,omitempty"`
	Domain       string            `yaml:"domain"`
	Type         string            `yaml:"type"`
	Dependencies []DependencySpec  `yaml:"dependencies"`
	Scripts      map[string]string `yaml:"scripts"`
	Env          map[string]string `yaml:"env"`
	Port         int               `yaml:"port,omitempty"`
	AutoStart    bool              `yaml:"auto_start"`
}

type RepositoryConfig struct {
	URL       string `yaml:"url"`
	Branch    string `yaml:"branch"`
	AutoClone bool   `yaml:"auto_clone"`
}

type DependencySpec struct {
	Name    string                 `yaml:"name"`
	Version string                 `yaml:"version"`
	Managed bool                   `yaml:"managed"`
	Config  map[string]interface{} `yaml:"config,omitempty"`
}

type ManagedDependency struct {
	InstallCommand string            `yaml:"install_command"`
	StartCommand   string            `yaml:"start_command"`
	StopCommand    string            `yaml:"stop_command"`
	ConfigFile     string            `yaml:"config_file,omitempty"`
	DataDir        string            `yaml:"data_dir,omitempty"`
	InitDatabases  []DatabaseConfig  `yaml:"init_databases,omitempty"`
	Environment    map[string]string `yaml:"environment,omitempty"`
}

type DatabaseConfig struct {
	Name  string `yaml:"name"`
	Owner string `yaml:"owner,omitempty"`
}

type DevelopmentConfig struct {
	StartupOrder  []string          `yaml:"startup_order"`
	GlobalScripts map[string]string `yaml:"global_scripts"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	Output string `yaml:"output"`
}

type HealthCheckConfig struct {
	Command  string `yaml:"command"`
	Interval string `yaml:"interval"`
	Timeout  string `yaml:"timeout"`
	Retries  int    `yaml:"retries"`
}

type ToolVersion struct {
	Version     string `yaml:"version"`
	DownloadURL string `yaml:"download_url,omitempty"`
}

type ProxyConfig struct {
	HTTPPort   int  `yaml:"http_port"`
	HTTPSPort  int  `yaml:"https_port"`
	Dashboard  bool `yaml:"dashboard"`
	AutoManage bool `yaml:"auto_manage"`
}

func (c *Config) Validate() error {
	if c.Proxy.HTTPPort <= 0 {
		c.Proxy.HTTPPort = 80
	}
	if c.Proxy.HTTPSPort <= 0 {
		c.Proxy.HTTPSPort = 443
	}

	for i := range c.Projects {
		if c.Projects[i].Name == "" {
			return &ValidationError{Field: "projects[].name", Message: "nome do projeto é obrigatório"}
		}
		if c.Projects[i].Type == "" {
			c.Projects[i].Type = "docker"
		}
	}

	return nil
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

func (c *Config) GetProjectByName(name string) *ProjectConfig {
	for i := range c.Projects {
		if c.Projects[i].Name == name {
			return &c.Projects[i]
		}
	}
	return nil
}

func (c *Config) GetProjectByDomain(domain string) *ProjectConfig {
	for i := range c.Projects {
		if c.Projects[i].Domain == domain {
			return &c.Projects[i]
		}
	}
	return nil
}

func (c *Config) MergeWith(other *Config) {
	for _, otherProject := range other.Projects {
		found := false
		for i := range c.Projects {
			if c.Projects[i].Name == otherProject.Name {
				c.Projects[i] = otherProject
				found = true
				break
			}
		}
		if !found {
			c.Projects = append(c.Projects, otherProject)
		}
	}

	if other.Tools != nil {
		if c.Tools == nil {
			c.Tools = make(map[string]ToolVersion)
		}
		for name, version := range other.Tools {
			c.Tools[name] = version
		}
	}

	if other.Proxy.HTTPPort != 0 {
		c.Proxy.HTTPPort = other.Proxy.HTTPPort
	}
	if other.Proxy.HTTPSPort != 0 {
		c.Proxy.HTTPSPort = other.Proxy.HTTPSPort
	}
	if other.Proxy.Dashboard {
		c.Proxy.Dashboard = other.Proxy.Dashboard
	}

	if other.Remote.URL != "" {
		c.Remote = other.Remote
	}
}
