# Guia de Configuração Global/Remota do Relief

## 🎯 Para que Serve?

O sistema de configuração global/remota permite que **organizações padronizem** o ambiente de desenvolvimento de todos os desenvolvedores automaticamente.

## 📊 Cenários de Uso

### Cenário 1: Organização com Equipe Distribuída

```yaml
# Servidor: https://config.empresa.com/relief.yaml
tools:
  node: { version: "20.11.0" }
  docker: { version: "24.0.0" }

project_defaults:
  dependencies:
    postgres: { enabled: true, version: "15" }
    redis: { enabled: true }
```

**Resultado**: Todos os devs automaticamente usam Node 20.11.0, Docker 24.0.0 e mesmas dependências.

### Cenário 2: Múltiplos Ambientes

```bash
# Desenvolvimento
relief start --remote-config=https://config.empresa.com/dev.yaml

# Produção  
relief start --remote-config=https://config.empresa.com/prod.yaml
```

### Cenário 3: Configuração Híbrida (Omelete)

```
┌─ Configuração Remota (empresa.com/relief.yaml)
│  ├─ Ferramentas padronizadas
│  ├─ Dependências obrigatórias  
│  └─ Templates de projeto
│
├─ Configuração Global Local (config.global.yaml)
│  ├─ Projetos específicos da organização
│  └─ Configurações de rede/proxy
│
└─ Configuração Local (config.local.yaml)
   ├─ Preferências do desenvolvedor
   └─ Projetos pessoais/experimentais
```

## 🚀 Setup Rápido

### Para Administradores:

1. **Crie configuração centralizada:**
   ```bash
   # Upload para servidor/CDN
   curl -X POST https://config.empresa.com/relief.yaml \
        -d @config-padrao.yaml
   ```

2. **Configure URL no Relief:**
   ```bash
   relief config set-remote https://config.empresa.com/relief.yaml
   ```

### Para Desenvolvedores:

```bash
# Instalação automática com config remota
relief init --remote https://config.empresa.com/relief.yaml

# Ou configurar depois
relief config remote enable --url https://config.empresa.com/relief.yaml
```

## 🔄 Comandos Úteis

```bash
# Ver configuração atual (merged)
relief config show

# Recarregar config remota
relief config refresh

# Testar nova URL remota  
relief config test-remote https://nova-url.com/config.yaml

# Verificar hierarquia
relief config status
```

## ✨ Vantagens Reais

### 1. **Onboarding Instantâneo**
Novo desenvolvedor: `relief init --remote URL` → ambiente pronto.

### 2. **Updates Sem Comunicação**
Admin atualiza Node.js → todos recebem automaticamente.

### 3. **Conformidade Automática**
Garante que toda equipe usa mesmas versões/configurações.

### 4. **Flexibilidade Preservada**
Cada dev pode customizar localmente sem afetar outros.

### 5. **Rollback Rápido**
Problema? Rollback da configuração remota → todos voltam à versão estável.

## 🎯 Exemplo Omelete

Com configuração remota, a Omelete pode:

1. **Padronizar** Node.js 20.11.0 para todos
2. **Forçar** PostgreSQL 15 em todos os projetos  
3. **Automatizar** setup do Traefik/proxy
4. **Distribuir** novos templates de projeto
5. **Atualizar** dependências centralizadamente

## 🔧 Implementação

O sistema já está implementado no Relief! Basta configurar:

```bash
cd relief/
echo 'remote:
  enabled: true  
  url: "https://sua-org.com/relief.yaml"' > config.global.yaml

relief start
```

Lindo e funcional! 🎉