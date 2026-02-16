package domain

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Manifest struct {
	Name         string                 `yaml:"name"`
	Domain       string                 `yaml:"domain"`
	Type         string                 `yaml:"type"`
	Dependencies []ManifestDependency   `yaml:"dependencies"`
	Scripts      map[string]string      `yaml:"scripts"`
	Env          map[string]string      `yaml:"env"`
	Ports        map[string]int         `yaml:"ports,omitempty"`
	Volumes      []string               `yaml:"volumes,omitempty"`
	Networks     []string               `yaml:"networks,omitempty"`
	Extra        map[string]interface{} `yaml:",inline"`
}

type ManifestDependency struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Managed bool   `yaml:"managed"`
}

func ParseManifest(projectPath string) (*Manifest, error) {
	manifestPath := filepath.Join(projectPath, "relief.yaml")

	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("relief.yaml not found. Please create a relief.yaml file in the project directory")
	}

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("error reading relief.yaml: %w", err)
	}

	var manifest Manifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("invalid YAML format in relief.yaml: %w", err)
	}

	if err := manifest.Validate(); err != nil {
		return nil, err
	}

	return &manifest, nil
}

func (m *Manifest) Validate() error {
	if m.Name == "" {
		return fmt.Errorf("'name' field is required")
	}
	if m.Type == "" {
		return fmt.Errorf("'type' field is required")
	}

	validTypes := map[string]bool{
		"docker": true,
		"node":   true,
		"python": true,
		"java":   true,
		"go":     true,
		"ruby":   true,
	}
	if !validTypes[m.Type] {
		return fmt.Errorf("type '%s' is not valid", m.Type)
	}

	return nil
}

func (m *Manifest) GetDevScript() string {
	if script, ok := m.Scripts["dev"]; ok {
		return script
	}
	return ""
}

func (m *Manifest) GetInstallScript() string {
	if script, ok := m.Scripts["install"]; ok {
		return script
	}
	return ""
}

func (m *Manifest) HasDependency(name string) bool {
	for _, dep := range m.Dependencies {
		if dep.Name == name {
			return true
		}
	}
	return false
}

func (m *Manifest) GetDependency(name string) *ManifestDependency {
	for i := range m.Dependencies {
		if m.Dependencies[i].Name == name {
			return &m.Dependencies[i]
		}
	}
	return nil
}

func (m *Manifest) ToProject(path string) *Project {
	projectType := ProjectType(m.Type)
	project := NewProject(m.Name, path, m.Domain, projectType)

	project.Scripts = m.Scripts
	project.Env = m.Env
	project.Manifest = m

	if portStr, ok := m.Env["PORT"]; ok {
		if port, err := strconv.Atoi(portStr); err == nil {
			project.Port = port
		}
	} else if m.Ports != nil {
		if mainPort, ok := m.Ports["main"]; ok {
			project.Port = mainPort
		}
	}

	for _, dep := range m.Dependencies {
		project.Dependencies = append(project.Dependencies, Dependency{
			Name:            dep.Name,
			RequiredVersion: dep.Version,
			Managed:         dep.Managed,
			Satisfied:       false,
		})
	}

	return project
}

func (m *Manifest) SaveManifest(projectPath string) error {
	manifestPath := filepath.Join(projectPath, "relief.yaml")

	data, err := yaml.Marshal(m)
	if err != nil {
		return fmt.Errorf("error serializing manifest: %w", err)
	}

	if err := os.WriteFile(manifestPath, data, 0644); err != nil {
		return fmt.Errorf("error saving manifest: %w", err)
	}

	return nil
}

func CreateDefaultManifest(name, projectType string) *Manifest {
	manifest := &Manifest{
		Name:         name,
		Domain:       name + ".local.test",
		Type:         projectType,
		Dependencies: []ManifestDependency{},
		Scripts:      make(map[string]string),
		Env:          make(map[string]string),
	}

	switch projectType {
	case "node":
		manifest.Scripts["dev"] = "npm run dev"
		manifest.Scripts["install"] = "npm install"
		manifest.Dependencies = append(manifest.Dependencies, ManifestDependency{
			Name:    "node",
			Version: ">=18.0.0",
			Managed: false,
		})
	case "python":
		manifest.Scripts["dev"] = "python main.py"
		manifest.Scripts["install"] = "pip install -r requirements.txt"
		manifest.Dependencies = append(manifest.Dependencies, ManifestDependency{
			Name:    "python",
			Version: ">=3.9",
			Managed: false,
		})
	case "docker":
		manifest.Scripts["dev"] = "docker-compose up"
		manifest.Scripts["install"] = "docker-compose pull"
	}

	return manifest
}
