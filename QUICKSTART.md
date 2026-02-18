# 🚀 Início Rápido - Relief

**Objetivo:** Ter o Relief funcionando em 10 minutos! ⏱️

Este guia oferece duas formas de instalação:
1. **⚡ Binário Pré-compilado** (mais rápido e fácil - recomendado!)
2. **🛠️ Compilar do Código-fonte** (se você quer desenvolver ou contribuir)

---

## ⚡ Instalação via Binário (MAIS RÁPIDO)

### 🍎 macOS - 2 minutos

```bash
# Intel (x86_64)
curl -L https://github.com/Maycon-Santos/relief/releases/latest/download/relief-darwin-amd64.tar.gz | tar xz
sudo mv Relief.app /Applications/
xattr -cr /Applications/Relief.app

# Apple Silicon (M1/M2/M3)
curl -L https://github.com/Maycon-Santos/relief/releases/latest/download/relief-darwin-arm64.tar.gz | tar xz
sudo mv Relief.app /Applications/
xattr -cr /Applications/Relief.app

# Adicione ao PATH (opcional, para usar no terminal)
echo 'export PATH="/Applications/Relief.app/Contents/MacOS:$PATH"' >> ~/.zshrc
source ~/.zshrc

# Execute!
open /Applications/Relief.app
# ou no terminal: relief
```

**Tempo aproximado:** 2 minutos ☕

---

### 🐧 Linux - 1 minuto

```bash
# x86_64 (AMD64)
curl -L https://github.com/Maycon-Santos/relief/releases/latest/download/relief-linux-amd64.tar.gz | tar xz
sudo mv relief /usr/local/bin/
sudo chmod +x /usr/local/bin/relief

# ARM64
curl -L https://github.com/Maycon-Santos/relief/releases/latest/download/relief-linux-arm64.tar.gz | tar xz
sudo mv relief /usr/local/bin/
sudo chmod +x /usr/local/bin/relief

# Execute!
relief
```

**Tempo aproximado:** 1 minuto ⚡

---

### 🪟 Windows - 2 minutos

```powershell
# PowerShell (abra normalmente, não precisa ser como Admin)
New-Item -ItemType Directory -Force -Path "$env:LOCALAPPDATA\Relief"
Invoke-WebRequest -Uri "https://github.com/Maycon-Santos/relief/releases/latest/download/relief-windows-amd64.exe.zip" -OutFile "$env:TEMP\relief.zip"
Expand-Archive -Path "$env:TEMP\relief.zip" -DestinationPath "$env:LOCALAPPDATA\Relief" -Force
$oldPath = [Environment]::GetEnvironmentVariable("Path", "User")
[Environment]::SetEnvironmentVariable("Path", "$oldPath;$env:LOCALAPPDATA\Relief", "User")
Remove-Item "$env:TEMP\relief.zip"

# Feche e abra um novo PowerShell, depois execute:
relief
```

**Tempo aproximado:** 2 minutos ☕

---

## 🛠️ Instalação Compilando do Código-fonte

Se você prefere compilar do zero (para desenvolvimento), continue lendo.

---

## 📋 Antes de Começar

**Você tem estas ferramentas instaladas?**

- [ ] **Go** (versão 1.22+)
- [ ] **Node.js** (versão 18+)  
- [ ] **Git**
- [ ] **Wails CLI**

**Não tem certeza?** Execute no terminal:

```bash
go version && node --version && git --version && wails version
```

Se todos mostrarem a versão, você está pronto! 🎉  
Se algum mostrar "command not found", veja as [instruções de instalação](#-não-tenho-as-ferramentas-instaladas).

---

## ⚡ Instalação Expressa

### 🍎 macOS

```bash
# 1. Instale as ferramentas (se não tiver)
brew install go node git
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 2. Configure o PATH
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc

# 3. Clone o Relief
git clone https://github.com/Maycon-Santos/relief.git
cd relief

# 4. Instale dependências
go mod download
cd frontend && npm install && cd ..

# 5. Execute!
wails dev
```

**Tempo aproximado:** 5-8 minutos (dependendo da internet) ☕

---

### 🐧 Linux (Ubuntu/Debian)

```bash
# 1. Instale o Go
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin:$(go env GOPATH)/bin' >> ~/.bashrc

# 2. Instale o Node.js
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# 3. Instale dependências do Wails
sudo apt-get install -y gcc libgtk-3-dev libwebkit2gtk-4 0-dev git

# 4. Instale o Wails
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 5. Clone o Relief
git clone https://github.com/Maycon-Santos/relief.git
cd relief

# 6. Instale dependências
go mod download
cd frontend && npm install && cd ..

# 7. Execute!
wails dev
```

**Tempo aproximado:** 8-12 minutos ☕☕

---

### 🪟 Windows (PowerShell)

```powershell
# 1. Instale o Chocolatey (gerenciador de pacotes)
# Abra PowerShell como ADMINISTRADOR e execute:
Set-ExecutionPolicy Bypass -Scope Process -Force
iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

# Feche e abra um novo PowerShell como Administrador

# 2. Instale as ferramentas
choco install golang nodejs git -y

# Feche e abra um novo PowerShell (não precisa mais ser como admin)

# 3. Configure o PATH do Go
$goPath = "$env:USERPROFILE\go\bin"
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";$goPath", "User")

# Feche e abra um novo PowerShell novamente

# 4. Instale o Wails
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 5. Clone o Relief
cd $env:USERPROFILE\Documents
git clone https://github.com/Maycon-Santos/relief.git
cd relief

# 6. Instale dependências
go mod download
cd frontend
npm install
cd ..

# 7. Execute!
wails dev
```

**Tempo aproximado:** 10-15 minutos ☕☕

---

## 🎯 Testando o Relief (3 minutos)

Agora que o Relief está rodando, vamos testar com um projeto de exemplo!

### 1️⃣ Adicione o Projeto de Exemplo

Na interface do Relief que abriu:

1. Clique em **"Adicionar Projeto Local"** (botão com ícone de "+")
   
2. **Navegue até a pasta do Relief** que você clonou
   - macOS/Linux: Provavelmente em `~/relief` ou onde você clonou
   - Windows: Provavelmente em `C:\Users\SeuNome\Documents\relief`

3. **Abra a pasta** `examples`

4. **Selecione a pasta** `hello-world`

5. Clique em **"Abrir"** ou **"Selecionar"**

✅ **Pronto!** O projeto "hello-world" apareceu na interface!

---

### 2️⃣ Inicie o Projeto

No card do projeto "hello-world":

1. Clique no **botão verde "Iniciar"** (▶️)

2. O status vai mudar:
   - 🟡 **"Iniciando..."** (aguarde 5-10 segundos)
   - 🟢 **"Rodando"** (pronto!)

> **💡 Primeira vez?** O Relief pode pedir sua senha para configurar o sistema. Digite sua senha normalmente (por segurança, nada aparece na tela while você digita).

---

### 3️⃣ Acesse no Navegador

1. Abra seu navegador (Chrome, Firefox, Safari, Edge...)

2. Acesse: **http://hello.local.dev**

3. Você deve ver:
   ```json
   {
     "message": "Hello from Relief Orchestrator!",
     "project": "hello-world"
   }
   ```

**🎉 PARABÉNS!** Você rodou seu primeiro projeto com o Relief!

---

### 4️⃣ Explore a Interface

**Ver logs do projeto:**
- Clique no botão **"Ver Logs"** (📋) para ver as mensagens do projeto

**Parar o projeto:**
- Clique no botão vermelho **"Parar"** (⏹️) quando terminar

**Remover o projeto:**
- Clique no botão **"Remover"** (🗑️) se quiser tirar o projeto do Relief

---

## 🎓 Próximos Passos

Agora que você já sabe o básico:

### 📚 Aprenda Mais

1. **Adicione seus próprios projetos**
   - Leia: [Como Adicionar Seus Projetos](README.md#-como-adicionar-seus-próprios-projetos)
   - Crie um arquivo `relief.yaml` nos seus projetos

2. **Entenda como funciona**
   - Leia: [Como Funcionam os Domínios .local.dev](README.md#-como-funcionam-os-domínios-localdev)

3. **Configurações avançadas**
   - Veja: [Schema do relief.yaml](docs/manifest-schema.md)

### 🛠️ Customize

Exemplos de `relief.yaml` para diferentes tecnologias:
- [Projeto Node.js](README.md#exemplo-1-projeto-nodejs)
- [Projeto Python/Flask](README.md#exemplo-2-projeto-pythonflask)
- [Projeto Next.js](README.md#exemplo-6-frontend-nextjs)
- [Projeto Docker](README.md#exemplo-4-projeto-dockerdocker-compose)

---

## 🆘 Não Tenho as Ferramentas Instaladas

Se você não tem Go, Node.js, Git ou Wails instalados, escolha seu sistema:

- 🍎 **macOS**: [Instruções completas para macOS](INSTALLATION.md#-instalação-no-macos)
- 🐧 **Linux**: [Instruções completas para Linux](INSTALLATION.md#-instalação-no-linux)
- 🪟 **Windows**: [Instruções completas para Windows](INSTALLATION.md#-instalação-no-windows)

---

## 🔴 Problemas Comuns

### "command not found" ao tentar executar wails, go, node, etc.

**Causa:** A ferramenta não está instalada ou não está no PATH.

**Solução:**
1. Verifique se a ferramenta está instalada
2. Se estiver, configure o PATH dela
3. **Sempre feche e abra novamente o terminal depois de instalar**

Veja: [Guia de Instalação](INSTALLATION.md)

---

### "Failed to modify /etc/hosts"

**Causa:** O Relief precisa de permissões para modificar o arquivo de hosts.

**Solução:**
- **macOS/Linux:** Digite sua senha quando solicitado
- **Windows:** Execute o PowerShell como Administrador

---

### "Port already in use"

**Causa:** Outra aplicação está usando a mesma porta.

**Solução:**

**macOS/Linux:**
```bash
# Descubra qual processo está usando a porta (ex: 34115)
lsof -i :34115

# Mate o processo
kill -9 PID
```

**Windows:**
```powershell
netstat -ano | findstr :34115
taskkill /PID <PID> /F
```

---

### Relief abre mas fica com tela em branco

**Solução:**

**Linux:**
```bash
sudo apt-get install -y libgtk-3-dev libwebkit2gtk-4.0-dev
```

**Windows:**
- Instale o WebView2: https://developer.microsoft.com/microsoft-edge/webview2/

**Todos:**
```bash
wails dev -debug  # Execute com logs detalhados
```

---

### Projeto não inicia (fica em "Iniciando...")

**Solução:**
1. Clique em **"Ver Logs"** no card do projeto
2. Veja qual a mensagem de erro
3. Geralmente é:
   - Dependência não instalada (ex: Node.js)
   - Comando `dev` incorreto no `relief.yaml`
   - Porta já em uso

---

### Domínio .local.dev não abre no navegador

**Solução rápida:**
1. Verifique se o projeto está 🟢 "Rodando"
2. Tente acessar pela porta direta: `http://localhost:PORTA`
3. Limpe o cache do navegador: `Ctrl+Shift+Delete`
4. Reinicie o Relief

**Solução completa:**
Veja: [Domínio .local.dev não funciona](README.md#domínio-localdev-não-funciona-no-navegador)

---

## ❓ Mais Ajuda?

- 📖 **Documentação completa**: [README.md](README.md)
- 🔧 **Guia de instalação detalhado**: [INSTALLATION.md](INSTALLATION.md)
- 🐛 **Problemas e soluções**: [Seção de Troubleshooting](README.md#-problemas-comuns-e-soluções)
- 💬 **Pergunte à comunidade**: [GitHub Discussions](https://github.com/Maycon-Santos/relief/discussions)
- 🐞 **Reporte bugs**: [GitHub Issues](https://github.com/Maycon-Santos/relief/issues)

---

## 📊 Checklist de Conclusão

Marque o que você já conseguiu fazer:

- [ ] Instalei todas as ferramentas necessárias
- [ ] Clonei o Relief do GitHub
- [ ] Instalei as dependências com sucesso
- [ ] Executei `wails dev` e a interface abriu
- [ ] Adicionei o projeto hello-world
- [ ] Iniciei o projeto com sucesso
- [ ] Acessei http://hello.local.dev no navegador e funcionou
- [ ] Vi os logs do projeto
- [ ] Parei o projeto

**Completou tudo?** 🎉 Você está pronto para usar o Relief!

**Teve algum problema?** Veja a seção de [Problemas Comuns](#-problemas-comuns) acima.

---

<p align="center">
  <b>Pronto para começar a adicionar seus próprios projetos?</b><br>
  Veja o <a href="README.md#-como-adicionar-seus-próprios-projetos">Guia Completo</a>!
</p>

<p align="center">
  <sub>⭐ Gostou do Relief? Deixe uma estrela no repositório!</sub>
</p>


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

- **Issues**: [GitHub Issues](https://github.com/Maycon-Santos/relief/issues)
- **Discussões**: [GitHub Discussions](https://github.com/Maycon-Santos/relief/discussions)
- **Documentação**: [README.md](README.md)

---

<p align="center">
  <b>Divirta-se usando o Relief! 🎉</b>
</p>
