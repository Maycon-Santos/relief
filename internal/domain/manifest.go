// Package domain contém as entidades de negócio da aplicação.
package domain

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Manifest representa o arquivo sofredor.yaml de um projeto
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

// ManifestDependency representa uma dependência no manifest
type ManifestDependency struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Managed bool   `yaml:"managed"`
}

// ParseManifest lê e faz parse do arquivo sofredor.yaml
func ParseManifest(projectPath string) (*Manifest, error) {
	manifestPath := filepath.Join(projectPath, "sofredor.yaml")

	// Verificar se arquivo existe
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("arquivo sofredor.yaml não encontrado em %s", projectPath)
	}

	// Ler arquivo
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler sofredor.yaml: %w", err)
	}

	// Parse YAML
	var manifest Manifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse do sofredor.yaml: %w", err)
	}

	// Validação básica
	if err := manifest.Validate(); err != nil {
		return nil, fmt.Errorf("manifest inválido: %w", err)
	}

	return &manifest, nil
}

// Validate valida o manifest
func (m *Manifest) Validate() error {
	if m.Name == "" {
		return fmt.Errorf("campo 'name' é obrigatório")
	}
	if m.Type == "" {
		return fmt.Errorf("campo 'type' é obrigatório")
	}

	// Validar tipo
	validTypes := map[string]bool{
		"docker": true,
		"node":   true,
		"python": true,
		"java":   true,
		"go":     true,
		"ruby":   true,
	}
	if !validTypes[m.Type] {
		return fmt.Errorf("tipo '%s' não é válido", m.Type)
	}

	return nil
}

// GetDevScript retorna o script de desenvolvimento
func (m *Manifest) GetDevScript() string {
	if script, ok := m.Scripts["dev"]; ok {
		return script
	}
	return ""
}

// GetInstallScript retorna o script de instalação
func (m *Manifest) GetInstallScript() string {
	if script, ok := m.Scripts["install"]; ok {
		return script
	}
	return ""
}

// HasDependency verifica se o manifest tem uma dependência específica
func (m *Manifest) HasDependency(name string) bool {
	for _, dep := range m.Dependencies {
		if dep.Name == name {
			return true
		}
	}
	return false
}

// GetDependency retorna uma dependência específica
func (m *Manifest) GetDependency(name string) *ManifestDependency {
	for i := range m.Dependencies {
		if m.Dependencies[i].Name == name {
			return &m.Dependencies[i]
		}
	}
	return nil
}

// ToProject converte o manifest em um Project
func (m *Manifest) ToProject(path string) *Project {
	projectType := ProjectType(m.Type)
	project := NewProject(m.Name, path, m.Domain, projectType)
	
	project.Scripts = m.Scripts
	project.Env = m.Env
	project.Manifest = m

	// Converter dependências
	for _, dep := range m.Dependencies {
		project.Dependencies = append(project.Dependencies, Dependency{
			Name:            dep.Name,
			RequiredVersion: dep.Version,
			Managed:         dep.Managed,
			Satisfied:       false, // Será verificado posteriormente
		})
	}

	return project
}

// SaveManifest salva o manifest em um arquivo
func (m *Manifest) SaveManifest(projectPath string) error {
	manifestPath := filepath.Join(projectPath, "sofredor.yaml")

	data, err := yaml.Marshal(m)
	if err != nil {
		return fmt.Errorf("erro ao serializar manifest: %w", err)
	}

	if err := os.WriteFile(manifestPath, data, 0644); err != nil {
		return fmt.Errorf("erro ao salvar manifest: %w", err)
	}

	return nil
}

// CreateDefaultManifest cria um manifest padrão para um projeto
func CreateDefaultManifest(name, projectType string) *Manifest {
	manifest := &Manifest{
		Name:         name,
		Domain:       name + ".sofredor.local",
		Type:         projectType,
		Dependencies: []ManifestDependency{},
		Scripts:      make(map[string]string),
		Env:          make(map[string]string),
	}

	// Scripts padrão baseados no tipo
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
