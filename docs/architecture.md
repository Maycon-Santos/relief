# Relief Orchestrator Architecture

## Overview

**Relief Orchestrator** is a local development orchestration tool that allows managing multiple projects simultaneously, using a modular and extensible architecture.

## Design Principles

1. **Modularity:** Each component has a single responsibility
2. **Extensibility:** Easy to add new runners, checkers, and providers
3. **Zero Secrets:** No hardcoded credentials
4. **Configuration as Code:** Everything defined in YAML
5. **Standard Go Layout:** Follows community best practices

## System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                   Frontend (React + TypeScript)              │
│  Components: ProjectCard, StatusBadge, LogsViewer, etc.     │
└──────────────────────────┬──────────────────────────────────┘
                           │ Wails Bindings
┌──────────────────────────┴──────────────────────────────────┐
│                      App Layer (Go)                          │
│  - GetProjects(), StartProject(), StopProject()             │
│  - AddLocalProject(), RefreshConfig()                       │
└──────────┬─────────┬─────────┬──────────┬───────────────────┘
           │         │         │          │
    ┌──────▼────┐ ┌─▼────────┐ ┌────────▼┐  ┌──────────▼──────┐
    │  Config   │ │  Runner  │ │Dependency│  │     Proxy       │
    │  Loader   │ │ Factory  │ │ Manager  │  │   Manager       │
    └───────────┘ └──────────┘ └──────────┘  └─────────────────┘
           │          │              │               │
           │          │              │               │
    ┌──────▼──────────▼──────────────▼───────────────▼─────────┐
    │               Storage (SQLite) + Logger                   │
    │  Tables: projects, dependencies, logs                     │
    └───────────────────────────────────────────────────────────┘
```

## Core Components

### 1. Domain Layer (`internal/domain/`)

Defines business entities:

- **Project:** Represents a development project
  - ID, Name, Path, Type, Status, Port
  - Domain, Dependencies, Environment
  
- **Manifest:** Parser for `relief.yaml` file
  - Name, Domain, Type
  - Dependencies, Scripts, Environment
  - Validation and default generation

### 2. Config Layer (`internal/config/`)

Manages configuration with layered approach:

- **Remote Config:** Fetches company configuration (HTTP)
- **Local Config:** Reads `~/.relief/config.yaml`
- **Local Override:** Reads `~/.relief/config.local.yaml`
- **Deep Merge:** Combines all configs (using `github.com/imdario/mergo`)
- **Validation:** Checks structure and required fields

**Flow:**
1. Fetch remote config (if URL configured)
2. Read local config
3. Read local override
4. Deep merge (override > local > remote)
5. Validate and return

### 3. Runner Layer (`internal/runner/`)

**Strategy Pattern** for project execution. Each project type has its own runner implementation.

#### Interface:
```go
type Runner interface {
    Start(ctx context.Context, project *domain.Project) error
    Stop(project *domain.Project) error
    Restart(ctx context.Context, project *domain.Project) error
    GetLogs(project *domain.Project, tail int) ([]string, error)
    IsRunning(project *domain.Project) bool
}
```

#### Implementations:

- **NativeRunner:** Executes projects directly on the OS
  - Uses `os/exec.CommandContext`
  - Captures stdout/stderr via pipes
  - Graceful shutdown with SIGTERM
  - Log buffering (last 1000 lines)
  - Environment variable injection

- **DockerRunner:** (stub for future implementation)
  - Use Docker SDK
  - Create/start containers
  - Network management

### 4. Dependency Layer (`internal/dependency/`)

Manages project dependencies with pluggable checkers:

**Current checkers:**
- **NodeChecker:** Verifies Node.js installation and version
- **PythonChecker:** Verifies Python installation and version
- **PostgresChecker:** Verifies PostgreSQL availability

**Adding new checkers:** Implement the `DependencyChecker` interface

### 5. Proxy Layer (`internal/proxy/`)

#### TraefikManager:
- Generates dynamic Traefik config (YAML)
- Routes `*.local.dev` to project ports
- Updates on project start/stop

#### HostsManager:
- Manipulates `/etc/hosts` file
- Adds entries `127.0.0.1 project.local.dev`
- Marked with `# BEGIN RELIEF` / `# END RELIEF`
- Requires elevated privileges (Linux/Mac)

### 6. Storage Layer (`internal/storage/`)

- **Location:** `~/.relief/data/orchestrator.db`
- **Driver:** SQLite3
- **Repositories:** ProjectRepository, LogRepository
- **Migrations:** Automatic on startup

### 7. App Layer (`internal/app/`)

Wails bindings exposing Go methods to React frontend:

```go
- GetProjects() - List all projects
- StartProject(id) - Start a project
- StopProject(id) - Stop a project
- RestartProject(id) - Restart a project
- GetProjectLogs(id, tail) - Get logs
- AddLocalProject(path) - Add project from path
- RemoveProject(id) - Remove a project
- RefreshConfig() - Reload configuration
- GetStatus() - Get orchestrator status
```

## Key Flows

### Starting a Project

1. User clicks "Start" button
2. App calls `StartProject(id)`
3. Runner executes dev script
4. Proxy adds hosts entry
5. Traefik configuration updated
6. Status saved to database

### Adding Local Project

1. User selects project folder
2. App parses `relief.yaml`
3. Validates manifest
4. Creates project entity
5. Checks dependencies
6. Saves to database

## Design Decisions

### Why Strategy Pattern for Runners?
- Easy to add Docker, Podman, Kubernetes support
- Can mock runners in tests
- Each runner is independent

### Why SQLite?
- Local-only, no network required
- Zero-config, no external database
- Sufficient performance for hundreds of projects

### Why Traefik?
- Dynamic configuration reloading
- Modern features (HTTP/2, gRPC)
- Proven in production

### Why Deep Merge?
- Company can set defaults
- Developers can customize locally
- Clear hierarchy (override > local > remote)

## Extension Points

- **New Project Types:** Create new runner in `internal/runner/`
- **New Dependency Checkers:** Add to `internal/dependency/checkers/`
- **Custom Storage:** Implement `Storage` interface

## Security

- No hardcoded secrets
- Elevated privileges only for `/etc/hosts`
- Input validation on all external inputs
- Parameterized SQL queries

## Future Improvements

- Auto-installers for Node.js/Python
- Full DockerRunner implementation
- Health checks for projects
- Desktop notifications
- Plugin system
- Team synchronization
