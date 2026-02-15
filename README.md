# ⚡ Relief Orchestrator

> **Hybrid local development orchestration tool**

Manage multiple projects simultaneously with support for Node.js, Python, Docker and more. Modern GUI, layered configuration, and zero hardcoded secrets.

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Wails](https://img.shields.io/badge/Wails-v2-DF5B00?style=flat)](https://wails.io)
[![React](https://img.shields.io/badge/React-18-61DAFB?style=flat&logo=react)](https://reactjs.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

---

## 🎯 Features

- ✅ **Multi-Runtime:** Node.js, Python, Go, Ruby, Java, Docker
- ✅ **Graphical Interface:** Modern GUI with Wails + React
- ✅ **Reverse Proxy:** Integrated Traefik for domain routing
- ✅ **Layered Config:** Remote (company) + Local (dev) with deep merge
- ✅ **Dependency Manager:** Automatic version verification
- ✅ **Zero Secrets:** No hardcoded credentials
- ✅ **Cross-Platform:** Linux, macOS, Windows
- ✅ **Real-time Logs:** Integrated log viewer
- ✅ **Open Source:** MIT License, contributions welcome

---

## 🚀 Quick Start

### Prerequisites

- Go 1.22+
- Node.js 18+
- Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

### Installation

```bash
# Clone the repository
git clone https://github.com/omelete/relief.git
cd relief

# Install dependencies
go mod download
cd frontend && npm install && cd ..

# Run in dev mode
wails dev
```

### Testing with Hello World Example

1. Start Relief Orchestrator
2. Click "Add Local Project"
3. Select the `examples/hello-world` folder
4. Click "Start" on the project
5. Access: `http://hello.local.dev`

You will see:
```json
{
  "message": "Hello from Relief Orchestrator!",
  "project": "hello-world",
  ...
}
```

---

## 📋 Usage

### 1. Create a `relief.yaml` in your project

```yaml
name: "my-api"
domain: "api.local.dev"
type: "node"

dependencies:
  - name: "node"
    version: ">=18.0.0"

scripts:
  dev: "npm run dev"
  install: "npm ci"

env:
  PORT: "3000"
  NODE_ENV: "development"
```

### 2. Add to Orchestrator

In the interface:
1. Click "Add Local Project"
2. Select the project folder
3. View dependency status
4. Click "Start"

### 3. Access the service

```bash
curl http://api.local.dev
```

---

## 🏗️ Architecture

### Overview

```
┌─────────────────────────────────────────────────┐
│              GUI (React + Wails)                 │
├─────────────────────────────────────────────────┤
│                 App Layer (Go)                   │
├──────────┬──────────┬──────────┬────────────────┤
│  Config  │  Runner  │ Depend.  │     Proxy      │
│  Loader  │ Factory  │ Manager  │  (Traefik)     │
├──────────┴──────────┴──────────┴────────────────┤
│           Storage (SQLite) + Logger             │
└─────────────────────────────────────────────────┘
```

### Main Components

- **Config Loader:** Remote + local config merge (YAML)
- **Runner Factory:** Strategy Pattern (Native, Docker)
- **Dependency Manager:** Checkers for Node, Python, PostgreSQL
- **Proxy Manager:** Traefik + `/etc/hosts` manipulation
- **Storage:** SQLite for state and logs

📖 [Complete architecture documentation](docs/architecture.md)

---

## 📁 Project Structure

```
/relief
├── cmd/app/              # Main entrypoint (Wails)
├── internal/             # Private code
│   ├── app/              # Wails bindings
│   ├── config/           # Configuration management
│   ├── domain/           # Business entities
│   ├── runner/           # Strategy Pattern (Native, Docker)
│   ├── dependency/       # Dependency checkers
│   ├── proxy/            # Traefik + Hosts manager
│   └── storage/          # SQLite + Repositories
├── pkg/                  # Reusable code (Logger, Utils)
├── frontend/             # React + TypeScript
│   └── src/
│       ├── components/   # ProjectCard, StatusBadge, LogsViewer
│       ├── hooks/        # useProjects
│       └── services/     # Wails API wrapper
├── examples/             # Example projects
│   └── hello-world/      # Ready-to-use Node.js example
├── configs/              # Example configurations
└── docs/                 # Documentation
```

---

## 🔧 Configuration

### Main File: `~/.relief/config.yaml`

```yaml
projects:
  - name: "my-api"
    path: "~/projects/my-api"
    domain: "api.local.dev"
    type: "node"
    auto_start: true

tools:
  node:
    version: "18.19.0"
  traefik:
    version: "2.10.7"

proxy:
  http_port: 80
  https_port: 443
  dashboard: true
```

### Local Configuration (Override): `~/.relief/config.local.yaml`

```yaml
# Overrides remote configuration
projects:
  - name: "my-api"
    path: "/custom/path"    # Path override
    env:
      PORT: "4000"           # Port override
```

🔗 [Complete relief.yaml schema](docs/manifest-schema.md)

---

## 🌐 Networking

### Traefik (Reverse Proxy)

The Orchestrator automatically configures Traefik to route `*.local.dev` domains to project ports.

**Routing example:**
```
hello.local.dev  →  localhost:3000
api.local.dev    →  localhost:4000
app.local.dev    →  localhost:5173
```

### /etc/hosts

The Orchestrator adds entries automatically:
```
127.0.0.1 hello.local.dev # RELIEF
127.0.0.1 api.local.dev   # RELIEF
```

⚠️ **Note:** Requires elevated privileges on Linux/Mac. The app will request permission.

---

## 🛠️ Development

### Run Tests

```bash
./build/ci/test.sh
```

### Run Linter

```bash
golangci-lint run
```

### Build for Production

```bash
./build/ci/build.sh
```

Generates binaries in `build/bin/` for:
- macOS (Intel e ARM)
- Linux (AMD64)
- Windows (AMD64)

---

## 🤝 Contributing

Contributions are very welcome! Please read the [Contributing Guide](docs/contributing.md).

### Areas that need help:

- 🐳 **DockerRunner:** Complete implementation
- 📦 **Installers:** Portable Node.js, Python
- 🎨 **UI/UX:** Design improvements
- 📚 **Docs:** Tutorials, examples, translations
- 🧪 **Tests:** Increase coverage

---

## 📚 Documentation

- [Architecture](docs/architecture.md) - Design decisions and flows
- [Contributing Guide](docs/contributing.md) - How to contribute
- [Manifest Schema](docs/manifest-schema.md) - `relief.yaml` reference
- [Hello World Example](examples/hello-world/README.md) - Practical tutorial

---

## 🐛 Issues & Bugs

Found a bug? [Open an issue](https://github.com/omelete/relief/issues/new)

---

## 📜 License

This project is licensed under the [MIT License](LICENSE).

---

## 🙏 Acknowledgments

- [Wails](https://wails.io) - Go + Web GUI framework
- [Traefik](https://traefik.io) - Modern reverse proxy
- [React](https://reactjs.org) - UI library
- All [contributors](https://github.com/omelete/relief/graphs/contributors)

---

## 📞 Contact

- **Email:** dev@omelete.com
- **Issues:** [GitHub Issues](https://github.com/omelete/relief/issues)

---

<p align="center">
  Made with ❤️ by the Omelete team
</p>

<p align="center">
  <sub>⭐ If you liked it, star the repository!</sub>
</p>
