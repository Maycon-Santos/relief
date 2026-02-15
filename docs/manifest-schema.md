# Schema do sofredor.yaml

O arquivo `sofredor.yaml` é o manifesto de configuração de cada projeto. Ele deve ficar na raiz do projeto.

## Estrutura Completa

```yaml
# Nome do projeto (obrigatório)
name: "my-service"

# Domínio local (opcional, padrão: <name>.sofredor.local)
domain: "api.sofredor.local"

# Tipo do projeto (obrigatório)
# Valores: docker, node, python, java, go, ruby
type: "node"

# Dependências de runtime (opcional)
dependencies:
  - name: "node"
    version: ">=18.0.0"  # Operadores: >=, >, <=, <, =
    managed: false        # Se true, orquestrador tenta instalar

  - name: "postgres"
    version: "15"
    managed: true         # Orquestrador gerencia via Docker

# Scripts de execução (obrigatório para tipos não-docker)
scripts:
  dev: "npm run dev"      # Script principal de desenvolvimento
  install: "npm ci"       # Script de instalação de dependências
  build: "npm run build"  # Script de build (opcional)
  test: "npm test"        # Script de testes (opcional)

# Variáveis de ambiente (opcional)
env:
  PORT: "3000"
  NODE_ENV: "development"
  DATABASE_URL: "postgresql://localhost:5432/mydb"

# Portas expostas (opcional, usado por DockerRunner)
ports:
  http: 3000
  metrics: 9090

# Volumes Docker (opcional, usado por DockerRunner)
volumes:
  - "./data:/app/data"
  - "/var/log/app:/logs"

# Networks Docker (opcional, usado por DockerRunner)
networks:
  - "sofredor-network"
```

## Campos Detalhados

### `name` (obrigatório)
- **Tipo:** string
- **Descrição:** Identificador único do projeto
- **Exemplo:** `"api-gateway"`

### `domain` (opcional)
- **Tipo:** string
- **Descrição:** Domínio local para acessar o serviço
- **Padrão:** `<name>.sofredor.local`
- **Exemplo:** `"api.sofredor.local"`

### `type` (obrigatório)
- **Tipo:** enum
- **Valores:**
  - `docker`: Usa Docker Compose
  - `node`: Projeto Node.js
  - `python`: Projeto Python
  - `java`: Projeto Java
  - `go`: Projeto Go
  - `ruby`: Projeto Ruby
- **Exemplo:** `"node"`

### `dependencies` (opcional)
Array de dependências de runtime.

#### Campos:
- `name`: Nome da ferramenta (`node`, `python`, `postgres`, etc)
- `version`: Versão requerida (suporta operadores semver)
- `managed`: Booleano indicando se o orquestrador deve gerenciar

**Operadores de Versão:**
```yaml
version: ">=18.0.0"  # Maior ou igual a 18.0.0
version: ">3.9"      # Maior que 3.9
version: "15"        # Exatamente 15.x
version: "=3.11.7"   # Exatamente 3.11.7
```

### `scripts` (obrigatório para non-docker)
Comandos de shell para diferentes operações.

**Scripts Comuns:**
- `dev`: Inicia o servidor de desenvolvimento
- `install`: Instala dependências
- `build`: Compila o projeto
- `test`: Executa testes

**Exemplos:**

Node.js:
```yaml
scripts:
  dev: "npm run dev"
  install: "npm ci"
```

Python:
```yaml
scripts:
  dev: "uvicorn main:app --reload"
  install: "pip install -r requirements.txt"
```

Go:
```yaml
scripts:
  dev: "air"  # Requer github.com/cosmtrek/air
  install: "go mod download"
```

### `env` (opcional)
Variáveis de ambiente injetadas no processo.

```yaml
env:
  PORT: "3000"
  NODE_ENV: "development"
  LOG_LEVEL: "debug"
  API_KEY: "${API_KEY}"  # Lê da env do sistema
```

### `ports` (opcional)
Mapeamento de portas (usado principalmente por Docker).

```yaml
ports:
  http: 3000
  https: 3443
  metrics: 9090
```

### `volumes` (opcional)
Volumes Docker (apenas para `type: docker`).

```yaml
volumes:
  - "./local:/container"  # Path relativo
  - "/absolute:/path"     # Path absoluto
```

### `networks` (opcional)
Networks Docker (apenas para `type: docker`).

```yaml
networks:
  - "sofredor-network"
  - "database-network"
```

## Exemplos por Tipo

### Node.js + PostgreSQL

```yaml
name: "user-api"
domain: "users.sofredor.local"
type: "node"

dependencies:
  - name: "node"
    version: ">=18.0.0"
    managed: false
  - name: "postgres"
    version: "15"
    managed: true

scripts:
  dev: "npm run dev"
  install: "npm ci"
  test: "npm test"

env:
  PORT: "3000"
  NODE_ENV: "development"
  DATABASE_URL: "postgresql://localhost:5432/users"
```

### Python FastAPI

```yaml
name: "ml-service"
domain: "ml.sofredor.local"
type: "python"

dependencies:
  - name: "python"
    version: ">=3.9"
    managed: false

scripts:
  dev: "uvicorn main:app --reload --port 8000"
  install: "pip install -r requirements.txt"

env:
  PORT: "8000"
  PYTHONUNBUFFERED: "1"
  MODEL_PATH: "./models"
```

### Docker Compose

```yaml
name: "legacy-monolith"
domain: "legacy.sofredor.local"
type: "docker"

scripts:
  dev: "docker-compose up"
  install: "docker-compose pull"

volumes:
  - "./data:/app/data"

networks:
  - "app-network"
```

### Frontend React + Vite

```yaml
name: "frontend-app"
domain: "app.sofredor.local"
type: "node"

dependencies:
  - name: "node"
    version: ">=18.0.0"

scripts:
  dev: "npm run dev"
  install: "npm install"
  build: "npm run build"

env:
  VITE_API_URL: "http://api.sofredor.local"
  VITE_PORT: "5173"
```

## Validação

O orquestrador valida o manifest ao carregar o projeto:

1. **Campos obrigatórios:** `name`, `type`
2. **Tipo válido:** Deve ser um dos tipos suportados
3. **Scripts presentes:** Para tipos non-docker, `dev` é obrigatório
4. **Formato de versões:** Deve seguir semver

Erros de validação são exibidos na UI com mensagens claras.

## Boas Práticas

1. **Use versões específicas:** Preferir `>=18.0.0` a `18`
2. **Configure managed=true:** Para dependências de dev (PostgreSQL, Redis)
3. **Documente env vars:** Adicione comentários no YAML
4. **Teste localmente:** Verifique que scripts funcionam antes de commitar
5. **Versionamento:** Commite o `sofredor.yaml` no git

## Migrando de Outras Ferramentas

### De docker-compose.yml
```yaml
# docker-compose.yml
services:
  api:
    build: .
    ports:
      - "3000:3000"
    environment:
      - NODE_ENV=development

# sofredor.yaml equivalente
name: "api"
type: "docker"
scripts:
  dev: "docker-compose up"
env:
  NODE_ENV: "development"
```

### De package.json scripts
```json
{
  "scripts": {
    "dev": "nodemon index.js",
    "start": "node index.js"
  }
}
```

```yaml
# sofredor.yaml
name: "my-app"
type: "node"
scripts:
  dev: "npm run dev"
```

## Troubleshooting

### Script não executa
- **Verifique o path:** Scripts executam na raiz do projeto
- **Permissões:** Binários devem ter permissão de execução
- **Shell:** Scripts são executados via `sh -c`

### Dependência não satisfeita
- **Versão instalada:** Verifique com `node -v`, `python --version`, etc
- **PATH correto:** Tool deve estar no PATH do sistema
- **managed=true:** Deixe o orquestrador instalar

### Porta em uso
- **Altere PORT:** Mude a variável `env.PORT`
- **Verifique conflitos:** `lsof -i :3000` (Linux/Mac)
