// Package domain contains the business entities of the application.
package domain

import (
	"time"
)

// ProjectType represents the type of a project
type ProjectType string

const (
	ProjectTypeDocker ProjectType = "docker"
	ProjectTypeNode   ProjectType = "node"
	ProjectTypePython ProjectType = "python"
	ProjectTypeJava   ProjectType = "java"
	ProjectTypeGo     ProjectType = "go"
	ProjectTypeRuby   ProjectType = "ruby"
)

// Status represents the current state of a project
type Status string

const (
	StatusStopped  Status = "stopped"
	StatusStarting Status = "starting"
	StatusRunning  Status = "running"
	StatusError    Status = "error"
	StatusUnknown  Status = "unknown"
)

// Project represents a project managed by the orchestrator
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
	CreatedAt    string            `json:"created_at"` // RFC3339 format
	UpdatedAt    string            `json:"updated_at"` // RFC3339 format
	LastError    string            `json:"last_error,omitempty"`
	// Git-related fields
	GitInfo *GitInfo `json:"git_info,omitempty"`
}

// GitInfo contém informações sobre o repositório Git do projeto
type GitInfo struct {
	IsRepository      bool     `json:"is_repository"`
	CurrentBranch     string   `json:"current_branch,omitempty"`
	AvailableBranches []string `json:"available_branches,omitempty"`
	RemoteURL         string   `json:"remote_url,omitempty"`
	HasChanges        bool     `json:"has_changes,omitempty"`
	LastCommit        string   `json:"last_commit,omitempty"`
}

// Dependency represents a dependency of a project
type Dependency struct {
	Name            string `json:"name"`
	Version         string `json:"version"`
	RequiredVersion string `json:"required_version"`
	Managed         bool   `json:"managed"`
	Satisfied       bool   `json:"satisfied"`
	Message         string `json:"message,omitempty"`
}

// LogEntry represents a log entry of a project
type LogEntry struct {
	ID        int64  `json:"id"`
	ProjectID string `json:"project_id"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"` // RFC3339 format
}

// NewProject creates a new Project instance
func NewProject(name, path, domain string, projectType ProjectType) *Project {
	now := time.Now().Format(time.RFC3339)
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

// IsRunning checks if the project is running
func (p *Project) IsRunning() bool {
	return p.Status == StatusRunning
}

// IsStopped checks if the project is stopped
func (p *Project) IsStopped() bool {
	return p.Status == StatusStopped
}

// HasError checks if the project is in error state
func (p *Project) HasError() bool {
	return p.Status == StatusError
}

// UpdateStatus updates the project status
func (p *Project) UpdateStatus(status Status) {
	p.Status = status
	p.UpdatedAt = time.Now().Format(time.RFC3339)
}

// SetError sets an error on the project
func (p *Project) SetError(err error) {
	p.Status = StatusError
	p.LastError = err.Error()
	p.UpdatedAt = time.Now().Format(time.RFC3339)
}

// ClearError clears the project error
func (p *Project) ClearError() {
	p.LastError = ""
	p.UpdatedAt = time.Now().Format(time.RFC3339)
}

// HasUnsatisfiedDependencies checks if there are unsatisfied dependencies
func (p *Project) HasUnsatisfiedDependencies() bool {
	for _, dep := range p.Dependencies {
		if !dep.Satisfied {
			return true
		}
	}
	return false
}

// GetUnsatisfiedDependencies returns the unsatisfied dependencies
func (p *Project) GetUnsatisfiedDependencies() []Dependency {
	unsatisfied := []Dependency{}
	for _, dep := range p.Dependencies {
		if !dep.Satisfied {
			unsatisfied = append(unsatisfied, dep)
		}
	}
	return unsatisfied
}

// UpdateGitInfo updates git-related information for the project
func (p *Project) UpdateGitInfo(gitInfo *GitInfo) {
	p.GitInfo = gitInfo
	p.UpdatedAt = time.Now().Format(time.RFC3339)
}

// HasGitRepository checks if the project has Git repository
func (p *Project) HasGitRepository() bool {
	return p.GitInfo != nil && p.GitInfo.IsRepository
}

// GetCurrentBranch returns the current git branch
func (p *Project) GetCurrentBranch() string {
	if p.GitInfo != nil {
		return p.GitInfo.CurrentBranch
	}
	return ""
}

// generateID generates a unique ID for the project
func generateID(name string) string {
	// Simplified: use name + timestamp
	// In production, use UUID or hash
	return name + "-" + time.Now().Format("20060102150405")
}
