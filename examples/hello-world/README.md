# 👋 Hello World - Exemplo para Relief

Este é um projeto de exemplo super simples para você testar o **Relief** e entender como funciona.

---

## 🤔 O Que Este Projeto Faz?

É um servidor Node.js básico que:
- Escuta na porta configurada (padrão: 3000)
- Retorna uma mensagem JSON quando você acessa
- Serve como modelo para você criar seus próprios projetos

---

## 📖 Tutorial Passo a Passo

### 1️⃣ Adicione o Projeto ao Relief

**Opção A - Pela Interface (Recomendado):**
1. Abra o Relief
2. Clique no botão **"Adicionar Projeto Local"**
3. Navegue até a pasta `examples/hello-world` dentro do Relief
4. Clique em **"Selecionar"**

**Opção B - Manualmente no Arquivo de Configuração:**

Edite `~/.relief/config.yaml` e adicione:
```yaml
projects:
  - name: "hello-world"
    path: "/caminho/completo/para/relief/examples/hello-world"
    domain: "hello.local.dev"
    type: "node"
```

### 2️⃣ Inicie o Projeto

1. No Relief, encontre o card **"hello-world"**
2. Clique no botão verde **"Iniciar"** ▶️
3. Aguarde o status mudar para 🟢 **Rodando**

### 3️⃣ Teste no Navegador

Abra seu navegador e acesse:
```
http://hello.local.dev
```

Você verá algo assim:
```json
{
  "message": "Hello from Relief Orchestrator!",
  "project": "hello-world",
  "timestamp": "2026-02-15T10:30:00.000Z",
  "port": 3000,
  "env": "development"
}
```

🎉 **Funcionou!** Seu primeiro projeto está rodando no Relief!

### 4️⃣ Teste via Terminal (Opcional)

Se preferir testar pelo terminal:
```bash
curl http://hello.local.dev
```

---

## 📂 O Que Tem Nesta Pasta?

```
hello-world/
├── relief.yaml      # Configuração do projeto para o Relief
├── index.js         # Código do servidor (muito simples!)
├── package.json     # Informações do projeto Node.js
└── README.md        # Este arquivo que você está lendo :)
```

### 📄 Entendendo o `relief.yaml`

```yaml
# Nome que aparece no Relief
name: "hello-world"

# URL que você vai acessar
domain: "hello.local.dev"

# Tipo do projeto (node, python, docker, etc)
type: "node"

# Verificações de requisitos
dependencies:
  - name: "node"
    version: ">=18.0.0"  # Precisa do Node 18 ou superior

# Comandos que o Relief vai executar
scripts:
  install: "npm install"  # Instala dependências na primeira vez
  dev: "node index.js"    # Comando para rodar o projeto

# Variáveis de ambiente
env:
  PORT: "3000"           # Porta onde o servidor vai rodar
  NODE_ENV: "development"
```

### 📄 Entendendo o `index.js`

```javascript
const http = require('http');

const port = process.env.PORT || 3000;

const server = http.createServer((req, res) => {
  const response = {
    message: "Hello from Relief Orchestrator!",
    project: "hello-world",
    timestamp: new Date().toISOString(),
    port: port,
    env: process.env.NODE_ENV || 'development'
  };

  res.writeHead(200, { 'Content-Type': 'application/json' });
  res.end(JSON.stringify(response, null, 2));
});

server.listen(port, () => {
  console.log(`🚀 Hello World rodando em http://localhost:${port}`);
  console.log(`🌐 Acesse via Relief: http://hello.local.dev`);
});
```

**O que este código faz:**
1. Cria um servidor HTTP básico
2. Quando alguém acessa, retorna um JSON com informações
3. Usa a porta definida no `relief.yaml`

---

## 🎨 Experimente Modificar!

### Mudar a Porta

1. Abra `relief.yaml`
2. Mude a linha:
   ```yaml
   env:
     PORT: "3001"  # Nova porta
   ```
3. No Relief, clique em **"Reiniciar"** 🔄

### Mudar o Domínio

1. Abra `relief.yaml`
2. Mude a linha:
   ```yaml
   domain: "ola.local.dev"  # Novo domínio
   ```
3. No Relief, clique em **"Reiniciar"** 🔄
4. Acesse: `http://ola.local.dev`

### Mudar a Mensagem

1. Abra `index.js`
2. Mude o texto em `message`:
   ```javascript
   message: "Olá! Modifiquei meu primeiro projeto!",
   ```
3. No Relief, clique em **"Reiniciar"** 🔄
4. Atualize o navegador para ver a mudança

---

## 🐛 Problemas Comuns

### ❌ Erro: "Porta já está em uso"

**Problema:** Outro programa está usando a porta 3000.

**Soluções:**
1. Mude a porta no `relief.yaml` (veja acima)
2. Ou pare o outro programa que está usando a porta

### ❌ Erro: "Node.js não encontrado"

**Problema:** Node.js não está instalado ou não está no PATH.

**Solução:**
```bash
# macOS
brew install node

# Linux
sudo apt install nodejs npm

# Verifique se funcionou
node --version
```

### ❌ Site não carrega (ERR_NAME_NOT_RESOLVED)

**Problema:** O domínio `.local.dev` não foi configurado no `/etc/hosts`.

**Solução:**
1. O Relief deveria fazer isso automaticamente
2. Verifique se você deu permissão quando o Relief pediu
3. Manualmente, verifique:
   ```bash
   cat /etc/hosts | grep hello
   ```
   Deveria aparecer:
   ```
   127.0.0.1 hello.local.dev # RELIEF
   ```

---

## 🎓 Próximos Passos

Agora que você testou o exemplo básico:

### 1. Crie Seu Próprio Projeto

Crie uma pasta nova com:
- Um arquivo `relief.yaml` (copie e modifique o deste exemplo)
- Seu código (pode ser Node.js, Python, etc.)
- Adicione no Relief

### 2. Rode Múltiplos Projetos

Adicione vários projetos ao Relief e rode todos ao mesmo tempo:
- Um frontend
- Uma API
- Um banco de dados

### 3. Explore Tipos Diferentes

Tente criar projetos com:
- **Python** (`type: "python"`)
- **Docker** (`type: "docker"`)
- **Go** (`type: "go"`)

---

## 💡 Dicas

- **Ver Logs:** Clique em "Ver Logs" para ver o que está acontecendo
- **Reiniciar Rápido:** Use Ctrl+R (ou Cmd+R no Mac) no Relief
- **Múltiplas Instâncias:** Você pode ter várias cópias deste exemplo com domínios diferentes

---

## 📚 Documentação

- [README Principal do Relief](../../README.md) - Documentação completa
- [Guia de Início Rápido](../../QUICKSTART.md) - Tutorial de instalação
- [Schema do relief.yaml](../../docs/manifest-schema.md) - Todas as opções

---

## 🤝 Contribua

Tem ideias para melhorar este exemplo? Abra uma issue ou pull request!

---

<p align="center">
  <b>Divirta-se explorando o Relief! 🚀</b>
</p>

1. Create your own `relief.yaml` in other projects
2. Explore managed dependencies
3. Configure multiple services
4. Use Docker for more complex projects
