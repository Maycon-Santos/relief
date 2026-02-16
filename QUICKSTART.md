# 🚀 Início Rápido - Relief

Guia prático de 5 minutos para começar a usar o Relief.

---

## ⚡ Instalação Expressa

### macOS

```bash
# 1. Instale as ferramentas necessárias
brew install go node

# 2. Instale o Wails
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 3. Clone o Relief
git clone https://github.com/seu-usuario/relief.git
cd relief

# 4. Instale dependências
go mod download
cd frontend && npm install && cd ..

# 5. Execute!
wails dev
```

### Linux (Ubuntu/Debian)

```bash
# 1. Instale Go
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# 2. Instale Node.js
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# 3. Instale dependências do Wails
sudo apt-get install -y gcc libgtk-3-dev libwebkit2gtk-4.0-dev

# 4. Instale Wails
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 5. Clone e configure
git clone https://github.com/seu-usuario/relief.git
cd relief
go mod download
cd frontend && npm install && cd ..

# 6. Execute!
wails dev
```

### Windows

```powershell
# 1. Instale Go
# Baixe e instale: https://go.dev/dl/

# 2. Instale Node.js
# Baixe e instale: https://nodejs.org/

# 3. Instale Wails
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 4. Clone e configure
git clone https://github.com/seu-usuario/relief.git
cd relief
go mod download
cd frontend
npm install
cd ..

# 5. Execute!
wails dev
```

---

## 🎯 Primeiro Uso (2 minutos)

### 1. Adicione o Projeto de Exemplo

1. Abra o Relief
2. Clique em **"Adicionar Projeto Local"**
3. Navegue até `relief/examples/hello-world`
4. Clique em **"Selecionar"**

### 2. Inicie o Projeto

1. No card do projeto, clique no botão verde **"Iniciar"** ▶️
2. Aguarde o status mudar para 🟢 **Rodando**

### 3. Acesse no Navegador

Abra: [http://hello.local.dev](http://hello.local.dev)

Você verá:
```json
{
  "message": "Hello from Relief Orchestrator!",
  "project": "hello-world"
}
```

🎉 **Sucesso!** Você está rodando um projeto com o Relief!

---

## 📝 Configure Seu Próprio Projeto

### Passo 1: Crie o Arquivo `relief.yaml`

Na raiz do seu projeto, crie:

**Para Node.js:**
```yaml
name: "meu-projeto"
domain: "meu-projeto.local.dev"
type: "node"

dependencies:
  - name: "node"
    version: ">=18.0.0"

scripts:
  install: "npm install"
  dev: "npm run dev"

env:
  PORT: "3000"
  NODE_ENV: "development"
```

**Para Python:**
```yaml
name: "meu-projeto"
domain: "meu-projeto.local.dev"
type: "python"

dependencies:
  - name: "python"
    version: ">=3.9"

scripts:
  install: "pip install -r requirements.txt"
  dev: "python app.py"

env:
  FLASK_ENV: "development"
```

**Para Docker:**
```yaml
name: "meu-projeto"
domain: "meu-projeto.local.dev"
type: "docker"

scripts:
  dev: "docker-compose up"
  stop: "docker-compose down"
```

### Passo 2: Adicione ao Relief

1. Clique em **"Adicionar Projeto Local"**
2. Selecione a pasta do seu projeto
3. Clique em **"Iniciar"** ▶️
4. Acesse `http://meu-projeto.local.dev`

---

## 🔧 Comandos Úteis

### Desenvolvimento

```bash
# Executar em modo desenvolvimento
wails dev

# Executar testes
./build/ci/test.sh

# Verificar código (linter)
golangci-lint run
```

### Build

```bash
# Compilar para produção
./build/ci/build.sh

# Compilar apenas para seu sistema
wails build

# O executável estará em: build/bin/relief
```

---

## 💡 Dicas Rápidas

### Ver Logs de um Projeto

1. Clique em **"Ver Logs"** 📋 no card do projeto
2. Os logs aparecem em tempo real

### Reiniciar um Projeto

1. Clique em **"Reiniciar"** 🔄
2. O projeto para e inicia automaticamente

### Verificar Dependências

Se um projeto não iniciar, verifique se as dependências estão instaladas:
- Node.js: `node --version`
- Python: `python --version` ou `python3 --version`
- Go: `go version`

### Problema com Portas

Se aparecer erro de "porta em uso":
1. O Relief detectará automaticamente
2. Mostrará qual processo está usando a porta
3. Você pode matar o processo direto pelo Relief

---

## ❓ Problemas Comuns

### "command not found: wails"

**Solução:** Adicione o Go bin ao PATH:
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Adicione essa linha ao seu `~/.bashrc` ou `~/.zshrc` para tornar permanente.

### "No files were processed"

Se você estiver vendo erros do Biome, ignore - são apenas avisos de formatação e não impedem o uso.

### Permissão negada no `/etc/hosts`

O Relief precisa modificar o arquivo `/etc/hosts` para configurar os domínios `.local.dev`.

**No macOS/Linux:** Digite sua senha quando solicitado.

### Porta 80 já está em uso

Outro serviço está usando a porta 80 (como Apache ou Nginx).

**Solução:**
1. Pare o outro serviço: `sudo systemctl stop apache2` (Linux)
2. Ou configure o Relief para usar outra porta no arquivo de configuração

---

## 📚 Próximos Passos

Agora que você tem o básico funcionando:

1. **Explore a Interface**: Clique em todos os botões e veja o que acontece
2. **Adicione Seus Projetos**: Crie arquivos `relief.yaml` para seus projetos reais
3. **Leia a Documentação Completa**: [README.md](README.md)
4. **Configure Múltiplos Projetos**: Veja como rodar vários ao mesmo tempo
5. **Contribua**: Melhorias são sempre bem-vindas!

---

## 🆘 Precisa de Ajuda?

- **Issues**: [GitHub Issues](https://github.com/seu-usuario/relief/issues)
- **Discussões**: [GitHub Discussions](https://github.com/seu-usuario/relief/discussions)
- **Documentação**: [README.md](README.md)

---

<p align="center">
  <b>Divirta-se usando o Relief! 🎉</b>
</p>
