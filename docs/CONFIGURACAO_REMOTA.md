# Configuração Remota para Organizações

## 🌐 Como Funciona a Configuração Remota

O Relief pode baixar automaticamente configurações de uma URL, permitindo que organizações mantenham padrões centralizados.

### 1. Configuração do Servidor

Disponibilize sua configuração via HTTP/HTTPS:

```yaml
# config-servidor.yaml (hospedado em servidor/CDN)
remote:
  enabled: true
  refresh_interval: "30m"

tools:
  node:
    version: "20.11.0"
  docker:
    version: "24.0.0"

projects:
  - name: "template-api"
    git:
      url: "https://github.com/company/api-template.git"
      branch: "main"
    runner: "docker"
    dependencies:
      postgres:
        enabled: true
        database: "api_template"

environment:
  company_name: "Minha Empresa"
  registry_url: "registry.empresa.com"
```

### 2. Configuração no Relief

Configure o Relief para usar a URL remota:

```yaml
# config.global.yaml (local)
remote:
  enabled: true
  url: "https://config.empresa.com/relief/config.yaml"
  refresh_interval: "1h"
  
# Configurações locais extras são mergeadas
proxy:
  http_port: 80
```

### 3. Comportamento Automático

1. **Startup**: Relief baixa config remota
2. **Merge**: Combina remote → global → local  
3. **Cache**: Mantém cache local em caso de falha de rede
4. **Refresh**: Atualiza automaticamente no intervalo configurado

## 📡 Casos de Uso Avançados

### Configuração por Ambiente

```bash
# Prod
export RELIEF_REMOTE_CONFIG="https://config.empresa.com/prod/relief.yaml"

# Dev  
export RELIEF_REMOTE_CONFIG="https://config.empresa.com/dev/relief.yaml"

# Local
# Usa config.global.yaml local apenas
```

### Versionamento de Configuração

```yaml
# URL com versioning
remote:
  url: "https://config.empresa.com/relief/v2.3/config.yaml"
  refresh_interval: "1h"
```

### Fallback Automático  

```yaml
remote:
  enabled: true
  urls:
    - "https://config.empresa.com/relief.yaml"     # Principal
    - "https://backup.empresa.com/relief.yaml"     # Backup
    - "file:///etc/relief/config.yaml"             # Local fallback
```

## 🔄 Comandos CLI

```bash
# Forçar reload da configuração remota
relief config reload --remote

# Ver configuração atual (merged)
relief config show

# Testar URL remota
relief config test-remote https://exemplo.com/config.yaml
```

## 🎯 Vantagens

- **Padronização**: Todos usam mesmas versões/configurações
- **Updates Centralizados**: Atualiza todos de uma vez
- **Flexibilidade**: Cada dev pode customizar localmente
- **Versionamento**: Controle de versão das configurações
- **Rollback**: Fácil voltar para versões anteriores