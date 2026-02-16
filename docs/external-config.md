# Configurações Externas no Relief

O Relief suporta configurações externas que permitem definir projetos e suas dependências fora do código principal, tornando-o útil para diferentes equipes e organizações.

## 📁 Tipos de Configuração

O Relief carrega configurações na seguinte ordem de prioridade:

1. **Configuração Remota** (opcional)
2. **Configuração Global** (`~/.relief/config.global.yaml`)
3. **Configuração Local** (`~/.relief/config.local.yaml`)

### Hierarquia de Merge

- **Remote** ← sobrescrito por **Global** ← sobrescrito por **Local**
- Configurações locais sempre têm prioridade máxima

## ⚙️ Estrutura da Configuração

### Configuração Básica de Projeto

```yaml
projects:
  - name: "my-backend"
    path: "/path/to/project"
    repository:
      url: "https://github.com/org/repo.git"
      branch: "main"
      auto_clone: true
    domain: "api.local.dev"
    type: "node"  # node, python, docker, java
    port: 3000
    auto_start: false
    dependencies:
      - name: "node"
        version: ">=18.0.0"
        managed: false
      - name: "postgres"
        version: "15"
        managed: true
        config:
          port: 5432
          database: "my_database"
          username: "postgres"
          password: "postgres"
    scripts:
      install: "npm ci"
      dev: "npm run dev"
      build: "npm run build"
      test: "npm test"
    env:
      NODE_ENV: "development"
      DATABASE_URL: "postgresql://postgres:postgres@localhost:5432/my_database"
```

### Dependências Gerenciadas

```yaml
managed_dependencies:
  postgres:
    install_command: "brew install postgresql@15"
    start_command: "brew services start postgresql@15"
    stop_command: "brew services stop postgresql@15"
    init_databases:
      - name: "my_database"
        owner: "postgres"
  
  redis:
    install_command: "brew install redis"
    start_command: "brew services start redis"
    stop_command: "brew services stop redis"
```

### Health Checks

```yaml
health_checks:
  postgres:
    command: "pg_isready -U postgres"
    interval: "5s"
    timeout: "5s"
    retries: 5
  
  redis:
    command: "redis-cli ping"
    interval: "5s"
    timeout: "5s"
    retries: 5
```

### Tools e Versões

```yaml
tools:
  node:
    version: "18.19.0"
  postgres:
    version: "15.5"
    install_method: "homebrew"
    homebrew_formula: "postgresql@15"
```

## 🚀 Como Usar

### 1. Criar Configuração Global

Crie um arquivo em `~/.relief/config.global.yaml`:

```yaml
projects:
  - name: "my-project"
    # ... configurações do projeto
```

### 2. Configuração Local (Override)

Para customizações específicas, crie `~/.relief/config.local.yaml`:

```yaml
projects:
  - name: "my-project"
    path: "/caminho/diferente"  # sobrescreve apenas o path
    env:
      DEBUG: "true"             # adiciona variável extra
```

### 3. Iniciar Relief

O Relief carregará automaticamente as configurações e:

- ✅ Clonará repositórios automaticamente se necessário
- ✅ Gerenciará dependências (PostgreSQL, Redis, etc.)
- ✅ Configurará health checks
- ✅ Aplicará configurações de proxy

## 📂 Exemplo de Workspace Externo

Para projetos específicos, crie um workspace separado:

```
my-workspace/
├── configs/
│   └── projects.yaml
├── scripts/
│   ├── setup.sh
│   └── test.sh
└── README.md
```

### Carregamento Manual

No Relief, use a opção "Load Config" para carregar `my-workspace/configs/projects.yaml`.

## 🔧 Funcionalidades Avançadas

### Git Integration

O Relief pode clonar e sincronizar repositórios automaticamente:

```yaml
projects:
  - name: "my-project"
    repository:
      url: "https://github.com/org/project.git"
      branch: "develop"
      auto_clone: true  # clona automaticamente se não existir
```

### Dependências Compartilhadas

Múltiplos projetos podem compartilhar dependências:

```yaml
projects:
  - name: "backend-api"
    dependencies:
      - name: "postgres"
        managed: true
        
  - name: "background-worker"
    dependencies:
      - name: "postgres"  # mesmo PostgreSQL
        managed: true
```

### Scripts Globais

Scripts que afetam o ambiente inteiro:

```yaml
development:
  global_scripts:
    setup_all: |
      echo "Configurando ambiente..."
      brew install postgresql redis
      
    start_all: |
      brew services start postgresql
      brew services start redis
      
    stop_all: |
      brew services stop postgresql  
      brew services stop redis
```

## 🎯 Casos de Uso

### 1. Equipe de Desenvolvimento

Cada desenvolvedor pode ter configurações locais específicas:

```yaml
# ~/.relief/config.local.yaml
projects:
  - name: "shared-project"
    path: "/Users/john/repos/project"  # path específico
    env:
      DEBUG_LEVEL: "verbose"           # preferência pessoal
```

### 2. Diferentes Ambientes

```yaml
# config.staging.yaml
projects:
  - name: "api"
    env:
      DATABASE_URL: "postgresql://staging-host:5432/db"
      
# config.local.yaml  
projects:
  - name: "api"
    env:
      DATABASE_URL: "postgresql://localhost:5432/db"
```

### 3. Microserviços

```yaml
projects:
  - name: "user-service"
    port: 3001
    dependencies: ["postgres", "redis"]
    
  - name: "auth-service"
    port: 3002
    dependencies: ["postgres"]
    
  - name: "notification-service"
    port: 3003
    dependencies: ["redis", "mongodb"]
```

## 📋 Templates

### Backend Node.js

```yaml
- name: "node-backend"
  type: "node"
  port: 3000
  dependencies:
    - name: "node"
      version: ">=18.0.0"
    - name: "postgres"
      version: "15"
      managed: true
  scripts:
    install: "npm ci"
    dev: "npm run dev"
    test: "npm test"
  env:
    NODE_ENV: "development"
```

### Frontend React/Next.js

```yaml
- name: "react-frontend"
  type: "node" 
  port: 3000
  dependencies:
    - name: "node"
      version: ">=18.0.0"
  scripts:
    install: "npm ci"
    dev: "npm start"
    build: "npm run build"
```

### Python API

```yaml
- name: "python-api"
  type: "python"
  port: 8000
  dependencies:
    - name: "python"
      version: ">=3.9"
    - name: "postgres"
      version: "15"
      managed: true
  scripts:
    install: "pip install -r requirements.txt"
    dev: "uvicorn main:app --reload"
```

## 🔍 Troubleshooting

### Configuração não carrega

1. Verifique sintaxe YAML
2. Confirme que arquivo está em `~/.relief/config.global.yaml`
3. Verifique logs do Relief

### Dependências não instalam

1. Verifique comandos de instalação no `managed_dependencies`
2. Confirme que Homebrew/pip estão instalados
3. Execute comandos manualmente para debug

### Projetos não aparecem

1. Verifique estrutura do arquivo de configuração
2. Confirme que seção `projects` existe
3. Recarregue configuração no Relief

## 📚 Referência Completa

### Campos de Projeto

- `name`: Nome único do projeto
- `path`: Caminho local do projeto
- `repository`: Configurações Git (opcional)
- `domain`: Domínio local para proxy
- `type`: Tipo do projeto (node, python, docker, java)
- `port`: Porta para executar o projeto
- `auto_start`: Iniciar automaticamente
- `dependencies`: Lista de dependências
- `scripts`: Scripts disponíveis
- `env`: Variáveis de ambiente

### Campos de Dependência

- `name`: Nome da dependência
- `version`: Versão requerida
- `managed`: Se Relief deve gerenciar
- `config`: Configurações específicas

Consulte o arquivo `config.example.yaml` para exemplos completos.