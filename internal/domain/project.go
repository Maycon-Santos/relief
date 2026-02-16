package domain

import (
	"time"
)

type ProjectType string

const (
	ProjectTypeDocker ProjectType = "docker"
	ProjectTypeNode   ProjectType = "node"
	ProjectTypePython ProjectType = "python"
	ProjectTypeJava   ProjectType = "java"
	ProjectTypeGo     ProjectType = "go"
	ProjectTypeRuby   ProjectType = "ruby"
)

type Status string

const (
	StatusStopped  Status = "stopped"
	StatusStarting Status = "starting"
	StatusRunning  Status = "running"
	StatusError    Status = "error"
	StatusUnknown  Status = "unknown"
)

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
	CreatedAt    string            `json:"created_at"`
	UpdatedAt    string            `json:"updated_at"`
	LastError    string            `json:"last_error,omitempty"`
	GitInfo *GitInfo `json:"git_info,omitempty"`
}

type GitInfo struct {
	IsRepository      bool     `json:"is_repository"`
	CurrentBranch     string   `json:"current_branch,omitempty"`
	AvailableBranches []string `json:"available_branches,omitempty"`
	RemoteURL         string   `json:"remote_url,omitempty"`
	HasChanges        bool     `json:"has_changes,omitempty"`
	LastCommit        string   `json:"last_commit,omitempty"`
}

type Dependency struct {
	Name            string `json:"name"`
	Version         string `json:"version"`
	RequiredVersion string `json:"required_version"`
	Managed         bool   `json:"managed"`
	Satisfied       bool   `json:"satisfied"`
	Message         string `json:"message,omitempty"`
}

type LogEntry struct {
	ID        int64  `json:"id"`
	ProjectID string `json:"project_id"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

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

func (p *Project) IsRunning() bool {
	return p.Status == StatusRunning
}

func (p *Project) IsStopped() bool {
	return p.Status == StatusStopped
}

func (p *Project) HasError() bool {
	return p.Status == StatusError
}

func (p *Project) UpdateStatus(status Status) {
	p.Status = status
	p.UpdatedAt = time.Now().Format(time.RFC3339)
}

func (p *Project) SetError(err error) {
	p.Status = StatusError
	p.LastError = err.Error()
	p.UpdatedAt = time.Now().Format(time.RFC3339)
}

func (p *Project) ClearError() {
	p.LastError = ""
	p.UpdatedAt = time.Now().Format(time.RFC3339)
}

func (p *Project) HasUnsatisfiedDependencies() bool {
	for _, dep := range p.Dependencies {
		if !dep.Satisfied {
			return true
		}
	}
	return false
}

func (p *Project) GetUnsatisfiedDependencies() []Dependency {
	unsatisfied := []Dependency{}
	for _, dep := range p.Dependencies {
		if !dep.Satisfied {
			unsatisfied = append(unsatisfied, dep)
		}
	}
	return unsatisfied
}

func (p *Project) UpdateGitInfo(gitInfo *GitInfo) {
	p.GitInfo = gitInfo
	p.UpdatedAt = time.Now().Format(time.RFC3339)
}

func (p *Project) HasGitRepository() bool {
	return p.GitInfo != nil && p.GitInfo.IsRepository
}

func (p *Project) GetCurrentBranch() string {
	if p.GitInfo != nil {
		return p.GitInfo.CurrentBranch
	}
	return ""
}

func generateID(name string) string {
	return name + "-" + time.Now().Format("20060102150405")
}
