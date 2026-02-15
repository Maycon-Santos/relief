# Hello World - Exemplo do SofredorOrchestrator

Este é um projeto de exemplo simples para testar o **SofredorOrchestrator**.

## 📋 O que faz?

Um servidor HTTP básico em Node.js que:
- Escuta na porta configurada (padrão: 3000)
- Retorna um JSON com informações sobre a requisição
- Demonstra como usar `sofredor.yaml` para configurar um projeto

## 🚀 Como usar com o Orchestrator

### 1. Adicionar o projeto

No SofredorOrchestrator, clique em "Add Local Project" e selecione esta pasta.

Ou adicione manualmente ao arquivo de configuração:

```yaml
projects:
  - name: "hello-world"
    path: "./examples/hello-world"
    domain: "hello.sofredor.local"
    type: "node"
```

### 2. Iniciar o projeto

No painel do Orchestrator:
1. Encontre o projeto "hello-world"
2. Clique no botão "Start"
3. Aguarde o status mudar para "Running"

### 3. Testar

Acesse no navegador:
```
http://hello.sofredor.local
```

Ou via curl:
```bash
curl http://hello.sofredor.local
```

Você deve ver uma resposta JSON como:
```json
{
  "message": "Hello from SofredorOrchestrator!",
  "project": "hello-world",
  "timestamp": "2026-02-14T10:30:00.000Z",
  "environment": "development",
  "path": "/",
  "method": "GET"
}
```

## 🔍 Estrutura

- `sofredor.yaml` - Manifesto do projeto (configuração)
- `index.js` - Servidor HTTP simples
- `package.json` - Metadados do projeto Node.js

## ⚙️ Requisitos

- Node.js >= 18.0.0 (verificado automaticamente pelo Orchestrator)

## 📝 Modificando

Experimente modificar:

1. **Porta:** Altere `PORT` no `sofredor.yaml`
2. **Domínio:** Mude `domain` para `teste.sofredor.local`
3. **Resposta:** Edite o objeto `response` em `index.js`

Após modificar, reinicie o projeto no Orchestrator.

## ❓ Problemas comuns

### Porta já em uso
Se a porta 3000 estiver ocupada, altere no `sofredor.yaml`:
```yaml
env:
  PORT: "3001"
```

### Domínio não resolve
Verifique se o Orchestrator adicionou a entrada em `/etc/hosts`:
```bash
cat /etc/hosts | grep sofredor
```

Deve ter:
```
127.0.0.1 hello.sofredor.local # SOFREDOR
```

## 🎓 Próximos passos

Agora que você testou o exemplo básico:

1. Crie seu próprio `sofredor.yaml` em outros projetos
2. Explore dependências gerenciadas
3. Configure múltiplos serviços
4. Use Docker para projetos mais complexos
