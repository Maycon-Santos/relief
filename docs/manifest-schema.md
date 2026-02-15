# relief.yaml Schema

The `relief.yaml` file is the configuration manifest for each project. It should be placed in the project root.

---

## Complete Example

```yaml
# Project name (required)
name: "my-api"

# Local domain (optional, default: <name>.local.dev)
domain: "api.local.dev"

# Project type (required): node, python, go, ruby, java, docker, etc.
type: "node"

# Port (optional, auto-assigned if not specified)
port: 3000

# Auto-start on orchestrator startup (optional, default: false)
auto_start: true

# Dependencies (optional)
dependencies:
  - name: "node"
    version: ">=18.0.0"
    managed: false  # If true, Relief will install
  - name: "postgres"
    version: "14"
    managed: true

# Execution scripts (required)
scripts:
  dev: "npm run dev"              # Development script
  install: "npm ci"               # Dependency installation
  build: "npm run build"          # Build (optional)
  test: "npm test"                # Tests (optional)

# Environment variables (optional)
env:
  PORT: "3000"
  NODE_ENV: "development"
  DATABASE_URL: "postgresql://localhost:5432/mydb"

# Docker-specific (only for type: docker)
docker:
  image: "node:18-alpine"
  dockerfile: "./Dockerfile"       # Or custom Dockerfile
  build_args:
    NODE_VERSION: "18"
  volumes:
    - "./src:/app/src"
    - "/app/node_modules"
  ports:
    - "3000:3000"
  networks:
    - "relief-network"
```

---

## Field Reference

### `name` (required)
- **Type:** `string`
- **Description:** Unique project identifier
- **Example:** `"my-api"`

### `domain` (optional)
- **Type:** `string`
- **Default:** `<name>.local.dev`
- **Description:** Local domain for proxied access
- **Example:** `"api.local.dev"`

### `type` (required)
- **Type:** `string`
- **Options:** `node`, `python`, `go`, `ruby`, `java`, `docker`, `static`
- **Description:** Project technology type
- **Example:** `"node"`

### `port` (optional)
- **Type:** `integer`
- **Description:** Port where the project will run
- **Default:** Auto-assigned (3000+)
- **Example:** `3000`

### `auto_start` (optional)
- **Type:** `boolean`
- **Default:** `false`
- **Description:** If true, starts automatically with Relief
- **Example:** `true`

### `dependencies` (optional)
- **Type:** `array of Dependency`
- **Description:** List of required dependencies
- **Example:**
  ```yaml
  dependencies:
    - name: "node"
      version: ">=18.0.0"
      managed: false
  ```

#### Dependency Object
```yaml
name: string        # Dependency name (node, python, postgres, etc.)
version: string     # Semantic version (>=X.Y.Z, ~X.Y, ^X.Y, =X.Y.Z)
managed: boolean    # If true, Relief manages installation
```

### `scripts` (required)
- **Type:** `object`
- **Description:** Execution commands
- **Required fields:**
  - `dev`: Development script
  - `install`: Dependency installation
- **Optional fields:**
  - `build`: Build script
  - `test`: Test script
  - `lint`: Linting script

**Example:**
```yaml
scripts:
  dev: "npm run dev"
  install: "npm ci"
  build: "npm run build"
  test: "npm test"
```

### `env` (optional)
- **Type:** `object (key-value)`
- **Description:** Environment variables injected at runtime
- **Example:**
  ```yaml
  env:
    PORT: "3000"
    NODE_ENV: "development"
    API_KEY: "${API_KEY}"  # Interpolation from system env
  ```

### `docker` (optional, only for `type: docker`)
- **Type:** `object`
- **Description:** Docker-specific configuration

#### Docker Fields
```yaml
docker:
  image: string          # Docker image (if not using Dockerfile)
  dockerfile: string     # Path to custom Dockerfile
  build_args: object     # Build arguments
  volumes: array         # Volume mounts
  ports: array           # Port mappings
  networks: array        # Docker networks
  environment: object    # Docker env vars (in addition to `env`)
```

---

## Examples by Type

### Node.js API

```yaml
name: "users-api"
domain: "users.local.dev"
type: "node"

dependencies:
  - name: "node"
    version: ">=18.0.0"
  - name: "postgres"
    version: "15"
    managed: true

scripts:
  dev: "npm run dev"
  install: "npm ci"
  test: "npm test"

env:
  PORT: "3001"
  DATABASE_URL: "postgresql://localhost:5432/users"
  JWT_SECRET: "${JWT_SECRET}"
```

### Python Django

```yaml
name: "admin-panel"
domain: "admin.local.dev"
type: "python"

dependencies:
  - name: "python"
    version: ">=3.9"

scripts:
  dev: "python manage.py runserver 0.0.0.0:8000"
  install: "pip install -r requirements.txt"
  test: "python manage.py test"

env:
  PORT: "8000"
  DJANGO_SETTINGS_MODULE: "project.settings.dev"
  DATABASE_URL: "postgresql://localhost:5432/admin"
```

### Go API

```yaml
name: "ml-service"
domain: "ml.local.dev"
type: "go"

dependencies:
  - name: "go"
    version: ">=1.22"

scripts:
  dev: "go run cmd/server/main.go"
  install: "go mod download"
  build: "go build -o bin/server cmd/server/main.go"
  test: "go test ./..."

env:
  PORT: "8080"
  ENV: "development"
```

### Ruby on Rails

```yaml
name: "legacy-app"
domain: "legacy.local.dev"
type: "ruby"

dependencies:
  - name: "ruby"
    version: "~3.2"
  - name: "postgres"
    version: "14"

scripts:
  dev: "rails server -b 0.0.0.0 -p 3000"
  install: "bundle install"
  test: "rspec"

env:
  PORT: "3000"
  RAILS_ENV: "development"
```

### Frontend (Vite/React)

```yaml
name: "frontend-app"
domain: "app.local.dev"
type: "node"

dependencies:
  - name: "node"
    version: ">=18.0.0"

scripts:
  dev: "npm run dev"
  install: "npm install"
  build: "npm run build"

env:
  PORT: "5173"
  VITE_API_URL: "http://api.local.dev"
  VITE_ENV: "development"
```

### Docker Compose

```yaml
name: "microservices"
domain: "micro.local.dev"
type: "docker"

scripts:
  dev: "docker-compose up"
  install: "docker-compose pull"
  build: "docker-compose build"

docker:
  dockerfile: "./docker-compose.yml"
```

---

## Best Practices

### 1. Version Control
- **Commit `relief.yaml`** to git
- **DO NOT commit `config.local.yaml`** (secrets)
- Use `.gitignore` for sensitive configs

### 2. Environment Variables
- Use `${VAR}` for interpolation from system
- Store secrets in `~/.relief/config.local.yaml`
- Never hardcode credentials

### 3. Dependencies
- **Specify exact versions** in production (`=X.Y.Z`)
- **Use ranges** in development (`>=X.Y.Z`)
- **Prefer `managed: false`** for system tools

### 4. Scripts
- **Keep scripts simple** (single command)
- **Use npm scripts / Makefile** for complex workflows
- **Ensure `dev` script runs continuously**

### 5. Domain Naming
- Use `.local.dev` suffix
- Keep names short and descriptive
- Avoid special characters

---

## Validation Rules

Relief validates the manifest automatically:

- ✅ `name` must be present and non-empty
- ✅ `type` must be a valid type
- ✅ `scripts.dev` must be present
- ✅ `scripts.install` must be present
- ✅ `domain` must be valid domain format
- ✅ `port` must be between 1024-65535
- ✅ Dependency versions must be valid semver

---

## Migration from Other Tools

### From Docker Compose

```yaml
# docker-compose.yml
services:
  api:
    image: node:18
    ports:
      - "3000:3000"
    environment:
      NODE_ENV: development
```

**relief.yaml equivalent:**
```yaml
name: "api"
domain: "api.local.dev"
type: "docker"
port: 3000

docker:
  image: "node:18"
  ports:
    - "3000:3000"

env:
  NODE_ENV: "development"

scripts:
  dev: "docker run -it --rm -p 3000:3000 node:18"
  install: "docker pull node:18"
```

### From Procfile (Heroku)

```
# Procfile
web: npm start
worker: npm run worker
```

**relief.yaml:**
```yaml
name: "my-app"
type: "node"

scripts:
  dev: "npm start"    # Only web process
  install: "npm ci"
```

---

## Troubleshooting

### Manifest not found

**Error:** `relief.yaml file not found`

**Solution:** Ensure file is in project root:
```bash
cd my-project
ls -la relief.yaml  # Must exist
```

### Invalid YAML syntax

**Error:** `error parsing relief.yaml`

**Solution:** Validate YAML:
```bash
yamllint relief.yaml
```

### Dependency not found

**Error:** `dependency 'node' not found`

**Solution:**
1. Install Node.js on your system
2. Or set `managed: true` (future feature)

### Port already in use

**Error:** `port 3000 already in use`

**Solution:**
1. Change port in `relief.yaml`:
   ```yaml
   env:
     PORT: "3001"
   ```
2. Or stop conflicting process

---

## Advanced

### Multi-service Projects

For projects with multiple services, create separate `relief.yaml` files:

```
/my-app
├── backend/
│   └── relief.yaml      # API service
├── frontend/
│   └── relief.yaml      # Web app
└── worker/
    └── relief.yaml      # Background jobs
```

### Dynamic Configuration

Use environment variable interpolation:

```yaml
env:
  DATABASE_URL: "${DATABASE_URL}"
  API_KEY: "${MY_API_KEY}"
  DEBUG: "${DEBUG:-false}"  # Default value
```

### Conditional Environments

Relief doesn't support conditionals directly. Use different config files:

```
relief.yaml           # Development
relief.prod.yaml      # Production (manually copy)
```

---

## Schema Version

Current schema version: **1.0.0**

Future versions will maintain backward compatibility.

---

## Support

Questions about the schema? [Open an issue](https://github.com/omelete/relief/issues/new)
