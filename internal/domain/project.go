// Package domain contém as entidades de negócio da aplicação.
package domain

import (
	"time"
)

// ProjectType representa o tipo de um projeto
type ProjectType string

const (
	ProjectTypeDocker ProjectType = "docker"
	ProjectTypeNode   ProjectType = "node"
	ProjectTypePython ProjectType = "python"
	ProjectTypeJava   ProjectType = "java"
	ProjectTypeGo     ProjectType = "go"
	ProjectTypeRuby   ProjectType = "ruby"
)

// Status representa o estado atual de um projeto
type Status string

const (
	StatusStopped Status = "stopped"
	StatusStarting Status = "starting"
	StatusRunning Status = "running"
	StatusError   Status = "error"
	StatusUnknown Status = "unknown"
)

// Project representa um projeto gerenciado pelo orquestrador
type Project struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Path         string            `json:"path"`
	Domain       string            `json:"domain"`
	Type         ProjectType       `json:"type"`
	Status       Status            `json:"status"`
	Port         int               `json:"port"`
	PID          int               `json:"pid,omitempty"`
	Dependencies []Dependency      `json:"dependencies"`
	Scripts      map[string]string `json:"scripts"`
	Env          map[string]string `json:"env"`
	Manifest     *Manifest         `json:"manifest,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	LastError    string            `json:"last_error,omitempty"`
}

// Dependency representa uma dependência de um projeto
type Dependency struct {
	Name            string `json:"name"`
	Version         string `json:"version"`
	RequiredVersion string `json:"required_version"`
	Managed         bool   `json:"managed"`
	Satisfied       bool   `json:"satisfied"`
	Message         string `json:"message,omitempty"`
}

// LogEntry representa uma entrada de log de um projeto
type LogEntry struct {
	ID        int64     `json:"id"`
	ProjectID string    `json:"project_id"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// NewProject cria uma nova instância de Project
func NewProject(name, path, domain string, projectType ProjectType) *Project {
	now := time.Now()
	return &Project{
		ID:           generateID(name),
		Name:         name,
		Path:         path,
		Domain:       domain,
		Type:         projectType,
		Status:       StatusStopped,
		Dependencies: []Dependency{},
		Scripts:      make(map[string]string),
		Env:          make(map[string]string),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// IsRunning verifica se o projeto está em execução
func (p *Project) IsRunning() bool {
	return p.Status == StatusRunning
}

// IsStopped verifica se o projeto está parado
func (p *Project) IsStopped() bool {
	return p.Status == StatusStopped
}

// HasError verifica se o projeto está em erro
func (p *Project) HasError() bool {
	return p.Status == StatusError
}

// UpdateStatus atualiza o status do projeto
func (p *Project) UpdateStatus(status Status) {
	p.Status = status
	p.UpdatedAt = time.Now()
}

// SetError define um erro no projeto
func (p *Project) SetError(err error) {
	p.Status = StatusError
	p.LastError = err.Error()
	p.UpdatedAt = time.Now()
}

// ClearError limpa o erro do projeto
func (p *Project) ClearError() {
	p.LastError = ""
	p.UpdatedAt = time.Now()
}

// HasUnsatisfiedDependencies verifica se há dependências não satisfeitas
func (p *Project) HasUnsatisfiedDependencies() bool {
	for _, dep := range p.Dependencies {
		if !dep.Satisfied {
			return true
		}
	}
	return false
}

// GetUnsatisfiedDependencies retorna as dependências não satisfeitas
func (p *Project) GetUnsatisfiedDependencies() []Dependency {
	unsatisfied := []Dependency{}
	for _, dep := range p.Dependencies {
		if !dep.Satisfied {
			unsatisfied = append(unsatisfied, dep)
		}
	}
	return unsatisfied
}

// generateID gera um ID único para o projeto
func generateID(name string) string {
	// Simplificado: usar nome + timestamp
	// Em produção, usar UUID ou hash
	return name + "-" + time.Now().Format("20060102150405")
}
