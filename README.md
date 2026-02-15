# ⚡ SofredorOrchestrator

> **Ferramenta de orquestração de desenvolvimento local híbrida**

Gerencie múltiplos projetos simultaneamente com suporte a Node.js, Python, Docker e mais. Interface gráfica moderna, configuração em camadas e zero secrets hardcoded.

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Wails](https://img.shields.io/badge/Wails-v2-DF5B00?style=flat)](https://wails.io)
[![React](https://img.shields.io/badge/React-18-61DAFB?style=flat&logo=react)](https://reactjs.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

---

## 🎯 Features

- ✅ **Multi-Runtime:** Node.js, Python, Go, Ruby, Java, Docker
- ✅ **Interface Gráfica:** GUI moderna com Wails + React
- ✅ **Proxy Reverso:** Traefik integrado para roteamento de domínios
- ✅ **Config em Camadas:** Remote (empresa) + Local (dev) com deep merge
- ✅ **Dependency Manager:** Verificação automática de versões
- ✅ **Zero Secrets:** Nenhuma credencial hardcoded
- ✅ **Cross-Platform:** Linux, macOS, Windows
- ✅ **Logs em Tempo Real:** Visualizador de logs integrado
- ✅ **Open Source:** MIT License, contribuições bem-vindas

---

## 🚀 Quick Start

### Pré-requisitos

- Go 1.22+
- Node.js 18+
- Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

### Instalação

```bash
# Clone o repositório
git clone https://github.com/omelete/sofredor-orchestrator.git
cd sofredor-orchestrator

# Instale dependências
go mod download
cd frontend && npm install && cd ..

# Execute em modo dev
wails dev
```

### Testando com o Exemplo Hello World

1. Inicie o Sofredor Orchestrator
2. Clique em "Add Local Project"
3. Selecione a pasta `examples/hello-world`
4. Clique em "Start" no projeto
5. Acesse: `http://hello.sofredor.local`

Você verá:
```json
{
  "message": "Hello from SofredorOrchestrator!",
  "project": "hello-world",
  ...
}
```

---

## 📋 Uso

### 1. Criar um `sofredor.yaml` no seu projeto

```yaml
name: "my-api"
domain: "api.sofredor.local"
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

### 2. Adicionar ao Orchestrator

Na interface:
1. Clique em "Add Local Project"
2. Selecione a pasta do projeto
3. Visualize o status das dependências
4. Clique em "Start"

### 3. Acessar o serviço

```bash
curl http://api.sofredor.local
```

---

## 🏗️ Arquitetura

### Visão Geral

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

### Componentes Principais

- **Config Loader:** Merge de config remota + local (YAML)
- **Runner Factory:** Strategy Pattern (Native, Docker)
- **Dependency Manager:** Checkers para Node, Python, PostgreSQL
- **Proxy Manager:** Traefik + manipulação de `/etc/hosts`
- **Storage:** SQLite para estado e logs

📖 [Documentação completa de arquitetura](docs/architecture.md)

---

## 📁 Estrutura do Projeto

```
/sofredor-orchestrator
├── cmd/app/              # Entrypoint principal (Wails)
├── internal/             # Código privado
│   ├── app/              # Wails bindings
│   ├── config/           # Gerenciamento de configuração
│   ├── domain/           # Entidades de negócio
│   ├── runner/           # Strategy Pattern (Native, Docker)
│   ├── dependency/       # Checkers de dependências
│   ├── proxy/            # Traefik + Hosts manager
│   └── storage/          # SQLite + Repositories
├── pkg/                  # Código reutilizável (Logger, Utils)
├── frontend/             # React + TypeScript
│   └── src/
│       ├── components/   # ProjectCard, StatusBadge, LogsViewer
│       ├── hooks/        # useProjects
│       └── services/     # Wails API wrapper
├── examples/             # Projetos de exemplo
│   └── hello-world/      # Exemplo Node.js pronto para uso
├── configs/              # Configurações de exemplo
└── docs/                 # Documentação
```

---

## 🔧 Configuração

### Arquivo Principal: `~/.sofredor/config.yaml`

```yaml
projects:
  - name: "my-api"
    path: "~/projects/my-api"
    domain: "api.sofredor.local"
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

### Configuração Local (Override): `~/.sofredor/config.local.yaml`

```yaml
# Sobrescreve configuração remota
projects:
  - name: "my-api"
    path: "/custom/path"    # Override do path
    env:
      PORT: "4000"           # Override da porta
```

🔗 [Schema completo do sofredor.yaml](docs/manifest-schema.md)

---

## 🌐 Networking

### Traefik (Proxy Reverso)

O Orchestrator configura automaticamente o Traefik para rotear domínios `*.sofredor.local` para as portas dos projetos.

**Exemplo de roteamento:**
```
hello.sofredor.local  →  localhost:3000
api.sofredor.local    →  localhost:4000
app.sofredor.local    →  localhost:5173
```

### /etc/hosts

O Orchestrator adiciona entradas automaticamente:
```
127.0.0.1 hello.sofredor.local # SOFREDOR
127.0.0.1 api.sofredor.local   # SOFREDOR
```

⚠️ **Nota:** Requer privilégios elevados no Linux/Mac. O app solicitará permissão.

---

## 🛠️ Desenvolvimento

### Executar Testes

```bash
./build/ci/test.sh
```

### Executar Linter

```bash
golangci-lint run
```

### Build para Produção

```bash
./build/ci/build.sh
```

Gera binários em `build/bin/` para:
- macOS (Intel e ARM)
- Linux (AMD64)
- Windows (AMD64)

---

## 🤝 Contribuindo

Contribuições são muito bem-vindas! Por favor, leia o [Guia de Contribuição](docs/contributing.md).

### Áreas que precisam de ajuda:

- 🐳 **DockerRunner:** Implementação completa
- 📦 **Instaladores:** Node.js, Python portáteis
- 🎨 **UI/UX:** Melhorias no design
- 📚 **Docs:** Tutoriais, exemplos, traduções
- 🧪 **Testes:** Aumentar cobertura

---

## 📖 Documentação

- [Arquitetura](docs/architecture.md) - Decisões de design e fluxos
- [Guia de Contribuição](docs/contributing.md) - Como contribuir
- [Schema do Manifest](docs/manifest-schema.md) - Referência do `sofredor.yaml`
- [Exemplo Hello World](examples/hello-world/README.md) - Tutorial prático

---

## 🐛 Issues & Bugs

Encontrou um bug? [Abra uma issue](https://github.com/omelete/sofredor-orchestrator/issues/new)

---

## 📜 Licença

Este projeto é licenciado sob a [MIT License](LICENSE).

---

## 🙏 Agradecimentos

- [Wails](https://wails.io) - Framework Go + Web GUI
- [Traefik](https://traefik.io) - Proxy reverso moderno
- [React](https://reactjs.org) - Biblioteca UI
- Todos os [contribuidores](https://github.com/omelete/sofredor-orchestrator/graphs/contributors)

---

## 📞 Contato

- **Email:** dev@omelete.com
- **Issues:** [GitHub Issues](https://github.com/omelete/sofredor-orchestrator/issues)

---

<p align="center">
  Feito com ❤️ pela equipe Omelete
</p>

<p align="center">
  <sub>⭐ Se você gostou, dê uma estrela no repositório!</sub>
</p>
