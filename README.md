# ⚡ Relief

> **Gerenciador visual para rodar múltiplos projetos localmente ao mesmo tempo**

Relief é uma ferramenta que permite você iniciar, parar e monitorar vários projetos de desenvolvimento (Node.js, Python, Docker, etc.) de forma simples através de uma interface gráfica moderna.

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Wails](https://img.shields.io/badge/Wails-v2-DF5B00?style=flat)](https://wails.io)
[![React](https://img.shields.io/badge/React-18-61DAFB?style=flat&logo=react)](https://reactjs.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Latest Release](https://img.shields.io/github/v/release/Maycon-Santos/relief?style=flat)](https://github.com/Maycon-Santos/relief/releases/latest)
[![GitHub Issues](https://img.shields.io/github/issues/Maycon-Santos/relief?style=flat)](https://github.com/Maycon-Santos/relief/issues)

---

## 📑 Índice

- [💡 Por Que Usar?](#-por-que-usar)
- [✨ O Que o Relief Faz?](#-o-que-o-relief-faz)
- [🎯 Para Quem é Esta Ferramenta?](#-para-quem-é-esta-ferramenta)
- [🚀 Instalação](#-instalação)
- [📖 Primeiro Uso](#-primeiro-uso)
- [🔧 Como Adicionar Seus Projetos](#-como-adicionar-seus-projetos)
- [🆘 Problemas Comuns](#-problemas-comuns)
- [❓ Perguntas Frequentes](#-perguntas-frequentes)
- [📚 Documentação Adicional](#-documentação-adicional)

---

## 📚 Documentos Úteis

- 🚀 **[Guia de Início Rápido](QUICKSTART.md)** - Tenha o Relief funcionando em 10 minutos
- 📦 **[Guia Completo de Instalação](INSTALLATION.md)** - Instruções detalhadas para cada sistema operacional
- 🆘 **[Guia de Solução de Problemas](TROUBLESHOOTING.md)** - Solução para problemas comuns
- 🏗️ **[Arquitetura](docs/architecture.md)** - Como o Relief funciona por dentro
- 📋 **[Schema do relief.yaml](docs/manifest-schema.md)** - Referência completa de configuração

---

## 🎯 Para Quem é Esta Ferramenta?

### Você vai gostar do Relief se:

✅ **Você trabalha com múltiplos projetos** - Tem mais de um projeto local para rodar ao mesmo tempo  
✅ **Quer economizar tempo** - Não quer ficar abrindo vários terminais e digitando comandos  
✅ **Prefere interfaces visuais** - Gosta de ver tudo organizado numa tela  
✅ **Trabalha em equipe** - Precisa padronizar como os projetos são executados  
✅ **É iniciante** - Não se sente confortável com terminal e linha de comando

### O que você precisa saber antes de começar?

**Conhecimento Básico (Necessário):**
- Como instalar programas no seu computador
- Como abrir o terminal/prompt de comando
- O que é uma pasta/diretório no seu computador

**Não é necessário:**
- Saber programar em Go ou React
- Conhecimento avançado de terminal
- Experiência prévia com Docker ou containers

> **💡 Dica:** Se você já desenvolve projetos em Node.js, Python, ou outras linguagens, você já tem o conhecimento necessário para usar o Relief!

---

## 💡 Por Que Usar?

Imagine que você trabalha em uma empresa com vários projetos:
- Uma API em Node.js
- Um frontend em React
- Um backend em Python
- Um banco de dados PostgreSQL

Toda vez que você for trabalhar, precisa:
1. Abrir 4 terminais diferentes
2. Navegar até cada pasta
3. Executar o comando para iniciar cada projeto
4. Lembrar as portas de cada um
5. Abrir o navegador em URLs diferentes

**Com o Relief, você:**
1. Abre o app
2. Clica em "Iniciar" em cada projeto
3. Pronto! 🎉

Todos os projetos ficam acessíveis em URLs amigáveis como:
- `http://api.local.dev`
- `http://app.local.dev`
- `http://admin.local.dev`

---

## ✨ O Que o Relief Faz?

### 🎯 Funcionalidades Principais

- **Inicia e Para Projetos**: Um clique para iniciar ou parar qualquer projeto
- **Interface Visual**: Veja o status de todos os projetos de uma só vez
- **Logs em Tempo Real**: Acompanhe o que está acontecendo em cada projeto
- **URLs Amigáveis**: Acesse seus projetos com nomes fáceis de lembrar
- **Gerencia Dependências**: Verifica se você tem Node.js, Python, etc. instalados
- **Multi-Linguagem**: Suporta Node.js, Python, Go, Ruby, Java e Docker
- **Git Integrado**: Veja em qual branch você está e sincronize facilmente

---

## 🚀 Instalação

**Escolha o seu método preferido:**

### ⚡ Opção 1: Binário Pré-compilado (Recomendado - Mais Rápido)

Baixe e instale o Relief em menos de 2 minutos!

**🍎 macOS:**
```bash
# Intel (x86_64)
curl -L https://github.com/Maycon-Santos/relief/releases/latest/download/relief-darwin-amd64.tar.gz | tar xz
sudo mv Relief.app /Applications/

# Apple Silicon (M1/M2/M3)
curl -L https://github.com/Maycon-Santos/relief/releases/latest/download/relief-darwin-arm64.tar.gz | tar xz
sudo mv Relief.app /Applications/

# Abra o aplicativo
open /Applications/Relief.app
```

**🐧 Linux:**
```bash
# x86_64/AMD64
curl -L https://github.com/Maycon-Santos/relief/releases/latest/download/relief-linux-amd64.tar.gz | tar xz
sudo mv relief /usr/local/bin/
relief

# ARM64
curl -L https://github.com/Maycon-Santos/relief/releases/latest/download/relief-linux-arm64.tar.gz | tar xz
sudo mv relief /usr/local/bin/
relief
```

**🪟 Windows (PowerShell):**
```powershell
# Download e instalação automática
New-Item -ItemType Directory -Force -Path "$env:LOCALAPPDATA\Relief"
Invoke-WebRequest -Uri "https://github.com/Maycon-Santos/relief/releases/latest/download/relief-windows-amd64.exe.zip" -OutFile "$env:TEMP\relief.zip"
Expand-Archive -Path "$env:TEMP\relief.zip" -DestinationPath "$env:LOCALAPPDATA\Relief" -Force

# Adicione ao PATH
$path = [Environment]::GetEnvironmentVariable("Path", "User")
[Environment]::SetEnvironmentVariable("Path", "$path;$env:LOCALAPPDATA\Relief", "User")

# Execute (feche e abra novo PowerShell)
relief
```

📥 **Ou baixe manualmente:** [Todas as versões](https://github.com/Maycon-Santos/relief/releases/latest)

---

### 🛠️ Opção 2: Compilar do Código-fonte

Para desenvolvedores que querem compilar o Relief ou contribuir com o projeto.

**📖 [Veja o guia completo de instalação](INSTALLATION.md)** para instruções detalhadas de cada sistema operacional.

**Resumo rápido:**

Escolha seu sistema operacional e siga o passo a passo:

---

### 🍎 macOS

#### O que você vai precisar:

1. **Homebrew** - Gerenciador de pacotes para macOS (facilita instalar programas)
2. **Go** - Linguagem de programação (Relief foi feito com Go)
3. **Node.js** - Para executar o frontend do Relief
4. **Wails** - Ferramenta para criar aplicativos desktop

#### Passo a Passo Completo:

**1️⃣ Instale o Homebrew**

Se você ainda não tem o Homebrew instalado:

```bash
# Abra o Terminal (Cmd + Espaço, digite "Terminal" e pressione Enter)
# Cole este comando e pressione Enter:
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

O instalador vai pedir sua senha do macOS - digite normalmente (não vai aparecer nada na tela, mas está sendo digitado).

**2️⃣ Instale o Go**

```bash
# No Terminal, execute:
brew install go

# Aguarde finalizar, então verifique se funcionou:
go version
```

Você deve ver algo como: `go version go1.22.0 darwin/amd64`

**3️⃣ Instale o Node.js**

```bash
# No Terminal, execute:
brew install node

# Verifique se funcionou:
node --version
npm --version
```

Você deve ver versões como `v20.11.0` e `10.2.4`

**4️⃣ Configure o PATH do Go**

O Go precisa que você configure onde ele guarda ferramentas. Execute:

```bash
# Adicione estas linhas ao seu arquivo de configuração do shell
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc

# Recarregue o arquivo de configuração
source ~/.zshrc
```

>**Nota:** Se você usa bash ao invés de zsh, substitua `~/.zshrc` por `~/.bash_profile`

**5️⃣ Instale o Wails**

```bash
# Instale o Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Verifique se funcionou:
wails version
```

Você deve ver: `Wails v2.x.x`

**6️⃣ Baixe o Relief**

```bash
# Navegue até a pasta onde você guarda seus projetos
cd ~/Development  # ou a pasta que você preferir

# Clone o repositório do Relief
git clone https://github.com/Maycon-Santos/relief.git

# Entre na pasta do Relief
cd relief
```

>**Não tem o Git instalado?** Execute: `brew install git`

**7️⃣ Instale as Dependências do Relief**

```bash
# Baixe as bibliotecas Go necessárias
go mod download

# Entre na pasta do frontend
cd frontend

# Instale as dependências do Node
npm install

# Volte para a pasta principal
cd ..
```

Este processo pode demorar alguns minutos, é normal!

**8️⃣ Execute o Relief**

```bash
# Inicie o Relief em modo desenvolvimento
wails dev
```

Uma janela vai abrir com a interface do Relief! 🎉

Se aparecer uma mensagem pedindo permissões, clique em **"Permitir"**.

---

### 🐧 Linux (Ubuntu/Debian)

#### O que você vai precisar:

1. **Go** - Linguagem de programação
2. **Node.js** - Para executar o frontend
3. **Dependências do sistema** - Bibliotecas para criar interfaces gráficas
4. **Wails** - Ferramenta para criar aplicativos desktop

#### Passo a Passo Completo:

**1️⃣ Atualize o Sistema**

```bash
# Abra o Terminal (Ctrl + Alt + T)
# Atualize a lista de pacotes:
sudo apt update && sudo apt upgrade -y
```

Digite sua senha quando solicitado.

**2️⃣ Instale o Go**

```bash
# Baixe o Go (verifique a versão mais recente em https://go.dev/dl/)
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz

# Remova instalações antigas (se houver)
sudo rm -rf /usr/local/go

# Extraia o arquivo baixado
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz

# Remova o arquivo de instalação
rm go1.22.0.linux-amd64.tar.gz
```

**3️⃣ Configure o PATH do Go**

```bash
# Adicione o Go ao seu PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc

# Recarregue as configurações
source ~/.bashrc

# Verifique se funcionou:
go version
```

Você deve ver: `go version go1.22.0 linux/amd64`

**4️⃣ Instale o Node.js**

```bash
# Adicione o repositório do Node.js (versão 18 LTS)
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -

# Instale o Node.js
sudo apt-get install -y nodejs

# Verifique se funcionou:
node --version
npm --version
```

**5️⃣ Instale as Dependências para o Wails**

O Wails precisa de algumas bibliotecas do sistema para criar interfaces gráficas:

```bash
# Instale as dependências necessárias
sudo apt-get install -y \
  gcc \
  libgtk-3-dev \
  libwebkit2gtk-4.0-dev \
  build-essential \
  pkg-config
```

**6️⃣ Instale o Wails**

```bash
# Instale o Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Verifique se funcionou:
wails version
```

**7️⃣ Instale o Git (se não tiver)**

```bash
# Verifique se você já tem o Git
git --version

# Se não tiver, instale:
sudo apt-get install -y git
```

**8️⃣ Baixe o Relief**

```bash
# Navegue até a pasta onde você guarda seus projetos
cd ~  # ou cd ~/Documentos ou qualquer pasta de sua preferência

# Clone o repositório do Relief
git clone https://github.com/Maycon-Santos/relief.git

# Entre na pasta do Relief
cd relief
```

**9️⃣ Instale as Dependências do Relief**

```bash
# Baixe as bibliotecas Go necessárias
go mod download

# Entre na pasta do frontend
cd frontend

# Instale as dependências do Node
npm install

# Volte para a pasta principal
cd ..
```

**🔟 Execute o Relief**

```bash
# Inicie o Relief em modo desenvolvimento
wails dev
```

Uma janela vai abrir com a interface do Relief! 🎉

---

### 🪟 Windows

#### O que você vai precisar:

1. **Go** - Linguagem de programação
2. **Node.js** - Para executar o frontend
3. **Git** - Para baixar o código do Relief
4. **Wails** - Ferramenta para criar aplicativos desktop
5. **WebView2** - Para renderizar a interface (geralmente já vem com Windows 11)

#### Passo a Passo Completo:

**1️⃣ Instale o Go**

1. Acesse: https://go.dev/dl/
2. Baixe o arquivo **Windows installer** (algo como `go1.22.0.windows-amd64.msi`)
3. Execute o instalador
4. Clique em **Next** → **Next** → **Install**
5. Aguarde a instalação concluir
6. Clique em **Finish**

**Verifique se funcionou:**

1. Abra o **PowerShell** ou **Prompt de Comando**:
   - Pressione `Win + R`
   - Digite `powershell` e pressione Enter
2. Digite: `go version`
3. Você deve ver: `go version go1.22.0 windows/amd64`

**2️⃣ Instale o Node.js**

1. Acesse: https://nodejs.org/
2. Baixe a versão **LTS** (recomendada)
3. Execute o instalador
4. Mantenha todas as opções padrão marcadas
5. Clique em **Next** → **Next** → **Install**
6. Aguarde a instalação concluir
7. Clique em **Finish**

**Verifique se funcionou:**

1. Abra um novo PowerShell (feche o anterior e abra novamente)
2. Digite: `node --version`
3. Digite: `npm --version`
4. Você deve ver as versões instaladas

**3️⃣ Instale o Git**

1. Acesse: https://git-scm.com/download/win
2. Baixe o instalador
3. Execute o instalador
4. Mantenha todas as opções padrão
5. Clique em **Next** várias vezes e depois em **Install**
6. Clique em **Finish**

**4️⃣ Configure o PATH do Go**

No PowerShell, execute:

```powershell
# Adicione o diretório de binários do Go ao PATH
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";$env:USERPROFILE\go\bin", "User")
```

**Feche e abra novamente o PowerShell** para as mudanças terem efeito.

**5️⃣ Instale o Wails**

No PowerShell, execute:

```powershell
# Instale o Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Verifique se funcionou:
wails version
```

Se o comando `wails` não for encontrado, feche e abra o PowerShell novamente.

**6️⃣ Verifique/Instale o WebView2**

O Windows 11 já vem com o WebView2. Se você usa Windows 10:

1. Acesse: https://developer.microsoft.com/microsoft-edge/webview2/
2. Baixe o **Evergreen Runtime**
3. Instale

**7️⃣ Baixe o Relief**

No PowerShell:

```powershell
# Navegue até a pasta onde você guarda seus projetos
cd $env:USERPROFILE\Documents

# Clone o repositório do Relief
git clone https://github.com/Maycon-Santos/relief.git

# Entre na pasta do Relief
cd relief
```

**8️⃣ Instale as Dependências do Relief**

```powershell
# Baixe as bibliotecas Go necessárias
go mod download

# Entre na pasta do frontend
cd frontend

# Instale as dependências do Node
npm install

# Volte para a pasta principal
cd ..
```

Este processo pode demorar alguns minutos.

**9️⃣ Execute o Relief**

```powershell
# Inicie o Relief em modo desenvolvimento
wails dev
```

Uma janela vai abrir com a interface do Relief! 🎉

Se o **Windows Defender** perguntar se você quer permitir o aplicativo, clique em **"Permitir acesso"**.

---

## 📖 Primeiro Uso

Agora que você já tem o Relief instalado, vamos testar com um projeto de exemplo!

### Passo 1: Abra o Relief

Se você ainda não abriu, execute no terminal:

```bash
wails dev
```

Uma janela vai aparecer com uma interface limpa, sem nenhum projeto ainda.

### Passo 2: Adicione o Projeto de Exemplo

O Relief vem com um projeto de exemplo chamado "hello-world". Vamos adicioná-lo:

1. **Clique no botão** **"Adicionar Projeto Local"** (ou "+") no canto superior da tela
2. **Uma janela de seleção de pasta vai abrir**
3. **Navegue até a pasta do Relief** que você baixou (ex: `~/Development/relief` ou `C:\Users\SeuNome\Documents\relief`)
4. **Entre na pasta** `examples`
5. **Entre na pasta** `hello-world`
6. **Clique em "Selecionar"** ou "Abrir"

Pronto! O projeto "hello-world" vai aparecer na interface como um cartão (card).

### Passo 3: Entenda o Card do Projeto

Você vai ver um cartão com:
- **📦 Nome**: "hello-world"
- **🔴 Status**: "Parado" (vermelho)
- **🌐 Domínio**: "hello.local.dev"
- **Botões de ação**: ▶️ Iniciar, 📋 Ver Logs, 🗑️ Remover

### Passo 4: Verifique as Dependências

Antes de iniciar, o Relief verifica se você tem tudo que precisa:

- Se aparecer um **alerta amarelo** dizendo "Node.js não encontrado", você precisa instalar o Node.js
- Se estiver tudo ok, você pode prosseguir

### Passo 5: Inicie o Projeto

1. **Clique no botão verde** **"Iniciar"** (▶️) no card do projeto
2. O status vai mudar para **🟡 "Iniciando..."**
3. O Relief vai:
   - Verificar se as dependências estão instaladas
   - Executar `npm install` (se necessário)
   - Executar `npm run dev`
   - Configurar o proxy Traefik
4. Após alguns segundos, o status vai mudar para **🟢 "Rodando"**

> **💡 Nota:** Na primeira vez, pode pedir sua senha para configurar o arquivo `/etc/hosts` (macOS/Linux) ou permissões de administrador (Windows). Isso é normal e seguro!

### Passo 6: Veja os Logs (Opcional)

Para ver o que está acontecendo "por baixo dos panos":

1. Clique no botão **"Ver Logs"** (📋) no card do projeto
2. Uma janela vai abrir mostrando todas as mensagens do projeto
3. Você verá algo como:
   ```
   > hello-world@1.0.0 dev
   > node index.js
   
   Server running on port 3000
   ```

### Passo 7: Acesse no Navegador

Abra seu navegador favorito (Chrome, Firefox, Safari, Edge) e acesse:

```
http://hello.local.dev
```

Você verá uma mensagem JSON:

```json
{
  "message": "Hello from Relief Orchestrator!",
  "project": "hello-world"
}
```

**🎉 Parabéns!** Você rodou seu primeiro projeto com o Relief!

### Passo 8: Pare o Projeto

Quando terminar de testar:

1. Volte para a janela do Relief
2. Clique no botão vermelho **"Parar"** (⏹️) no card do projeto
3. O status vai mudar para **🔴 "Parado"**

---

## 🔧 Como Adicionar Seus Próprios Projetos

Agora que você já testou o exemplo, vamos configurar seus próprios projetos!

### O Que Você Precisa Fazer

Para que o Relief reconheça seu projeto, você precisa criar um arquivo de configuração chamado `relief.yaml` na pasta raiz do projeto.

### Estrutura do Arquivo relief.yaml

O arquivo tem 4 seções principais:

1. **Informações gerais** - Nome e domínio
2. **Tipo e dependências** - Que tecnologia usa
3. **Comandos** - Como instalar e rodar
4. **Variáveis de ambiente** - Configurações do projeto

---

### Exemplo 1: Projeto Node.js

Vamos supor que você tem um projeto Node.js/Express na pasta `/Users/voce/projetos/minha-api`.

**1. Abra a pasta do seu projeto no terminal ou editor de código**

**2. Crie um arquivo chamado `relief.yaml` na raiz do projeto**

**3. Adicione este conteúdo:**

```yaml
# =============================
# INFORMAÇÕES GERAIS
# =============================

# Nome que aparecerá na interface do Relief
name: "minha-api"

# URL pela qual você vai acessar (sem http://)
# Exemplo: se você colocar "api.local.dev", vai acessar em http://api.local.dev
domain: "api.local.dev"

# =============================
# TIPO E DEPENDÊNCIAS
# =============================

# Tipo do projeto: node, python, docker, go, ruby, java
type: "node"

# Lista de dependências necessárias (o Relief verifica se você tem instalado)
dependencies:
  - name: "node"           # Nome da ferramenta
    version: ">=18.0.0"     # Versão mínima necessária

# =============================
# COMANDOS
# =============================

scripts:
  # Comando para instalar dependências (executado quando você adiciona o projeto)
  install: "npm install"
  
  # Comando para iniciar o projeto (executado quando você clica em "Iniciar")
  dev: "npm run dev"
  
  # Comando para parar (opcional, o Relief para automaticamente)
  # stop: "npm run stop"

# =============================
# VARIÁVEIS DE AMBIENTE
# =============================

env:
  PORT: "3000"                    # Porta onde seu app vai rodar
  NODE_ENV: "development"         # Ambiente de execução
  DATABASE_URL: "postgresql://localhost:5432/mydb"  # Exemplo de conexão com banco
```

**4. Salve o arquivo**

**5. No Relief, clique em "Adicionar Projeto Local" e selecione a pasta `/Users/voce/projetos/minha-api`**

Pronto! Seu projeto vai aparecer no Relief e você pode iniciá-lo com um clique! 🎉

---

### Exemplo 2: Projeto Python/Flask

Para um projeto Python com Flask:

```yaml
name: "api-python"
domain: "python-api.local.dev"
type: "python"

dependencies:
  - name: "python"
    version: ">=3.9"
  - name: "pip"
    version: ">=20.0"

scripts:
  # Instala as dependências do requirements.txt
  install: "pip install -r requirements.txt"
  
  # Inicia o servidor Flask
  dev: "python app.py"

env:
  FLASK_APP: "app.py"
  FLASK_ENV: "development"
  PORT: "5000"
```

---

### Exemplo 3: Projeto Python/Django

Para um projeto Django:

```yaml
name: "meu-site-django"
domain: "django.local.dev"
type: "python"

dependencies:
  - name: "python"
    version: ">=3.10"

scripts:
  install: "pip install -r requirements.txt"
  
  # Django usa o comando manage.py runserver
  dev: "python manage.py runserver 0.0.0.0:8000"

env:
  DJANGO_SETTINGS_MODULE: "myproject.settings"
  DEBUG: "True"
```

---

### Exemplo 4: Projeto Docker/Docker Compose

Se seu projeto usa Docker:

```yaml
name: "app-dockerizado"
domain: "docker-app.local.dev"
type: "docker"

dependencies:
  - name: "docker"
    version: ">=20.0"

scripts:
  # Sobe os containers
  dev: "docker-compose up"
  
  # Para os containers
  stop: "docker-compose down"

# Para projetos Docker, geralmente as envs ficam no docker-compose.yml
# mas você pode adicionar aqui também se quiser
env:
  COMPOSE_PROJECT_NAME: "meu-projeto"
```

---

### Exemplo 5: Frontend React/Vite

Para um frontend React com Vite:

```yaml
name: "meu-frontend"
domain: "app.local.dev"
type: "node"

dependencies:
  - name: "node"
    version: ">=18.0.0"

scripts:
  install: "npm install"
  
  # Vite geralmente roda com 'npm run dev'
  dev: "npm run dev"

env:
  # URL da API que o frontend vai consumir
  VITE_API_URL: "http://api.local.dev"
  PORT: "5173"
```

---

### Exemplo 6: Frontend Next.js

Para projetos Next.js:

```yaml
name: "site-nextjs"
domain: "site.local.dev"
type: "node"

dependencies:
  - name: "node"
    version: ">=18.0.0"

scripts:
  install: "npm install"
  dev: "npm run dev"

env:
  # URL base da aplicação
  NEXT_PUBLIC_API_URL: "http://api.local.dev"
  PORT: "3000"
```

---

### 🎓 Entendendo Cada Campo

#### `name` (obrigatório)
- **O que é:** Nome do projeto que aparece na interface
- **Exemplo:** `"minha-api"`, `"frontend"`, `"backoffice"`
- **Dica:** Use nomes curtos e descritivos

#### `domain` (obrigatório)
- **O que é:** URL pela qual você vai acessar o projeto
- **Formato:** `"<nome>.local.dev"` (sem `http://`)
- **Exemplo:** `"api.local.dev"` → você acessará em `http://api.local.dev`
- **Dica:** Use domínios que façam sentido (`api`, `admin`, `app`, `web`)

#### `type` (obrigatório)
- **O que é:** Tipo de tecnologia do projeto
- **Opções:** `node`, `python`, `go`, `ruby`, `java`, `docker`
- **Por que importa:** O Relief usa isso para saber como executar seu projeto

#### `dependencies` (opcional mas recomendado)
- **O que é:** Lista de ferramentas que seu projeto precisa
- **Formato:**
  ```yaml
  dependencies:
    - name: "nome-da-ferramenta"
      version: ">=versao-minima"
  ```
- **Exemplo prático:** Se seu projeto precisa do Node.js versão 18 ou superior:
  ```yaml
  dependencies:
    - name: "node"
      version: ">=18.0.0"
  ```
- **O Relief vai:** Avisar você se não tiver a ferramenta instalada ou se a versão está antiga

#### `scripts` (obrigatório)
- **O que é:** Comandos para gerenciar o projeto
- **Campos principais:**
  - `install`: Comando para instalar dependências (ex: `npm install`, `pip install -r requirements.txt`)
  - `dev`: Comando para iniciar o projeto em desenvolvimento (ex: `npm run dev`, `python app.py`)
  - `stop`: (opcional) Comando para parar o projeto (geralmente não é necessário)
- **Importante:** Use exatamente os comandos que você normalmente digita no terminal

#### `env` (opcional)
- **O que é:** Variáveis de ambiente que seu projeto precisa
- **Formato:** Chave-valor
  ```yaml
  env:
    NOME_DA_VARIAVEL: "valor"
    OUTRA_VARIAVEL: "outro-valor"
  ```
- **Exemplo prático:**
  ```yaml
  env:
    PORT: "3000"
    DATABASE_URL: "postgresql://localhost:5432/mydb"
    API_KEY: "minha-chave-secreta"
  ```
- **Dica:** Variáveis sensíveis (senhas, tokens) podem ser definidas aqui, mas considere usar um arquivo `.env` também

---

### ✅ Checklist: Antes de Adicionar seu Projeto

Use esta lista para garantir que está tudo certo:

- [ ] Criei o arquivo `relief.yaml` na raiz do projeto?
- [ ] Defini o `name`, `domain` e `type`?
- [ ] Listei as `dependencies` necessárias?
- [ ] Configurei o comando `install` correto?
- [ ] Configurei o comando `dev` que realmente inicia meu projeto?
- [ ] Adicionei as `env` vars necessárias (se houver)?
- [ ] Testei os comandos manualmente no terminal antes?

---

### 🚨 Erros Comuns ao Configurar

#### ❌ "Projeto não inicia"
**Problema:** O comando `dev` está incorreto  
**Solução:** Teste o comando manualmente no terminal da pasta do projeto primeiro

#### ❌ "Domínio não funciona"
**Problema:** O domínio precisa terminar com `.local.dev`  
**Solução:** Use sempre o formato `<nome>.local.dev`

#### ❌ "Dependências não encontradas"
**Problema:** A ferramenta listada não está instalada ou o nome está errado  
**Solução:** Verifique se o nome está correto (`node`, não `nodejs`) e se está no PATH

#### ❌ "Porta já está em uso"
**Problema:** Outro projeto está usando a mesma porta  
**Solução:** Mude o `PORT` no `env` para uma porta diferente

---

## � Problemas Comuns e Soluções

### 🔴 Problemas na Instalação

#### "go: command not found"

**Sintoma:** Quando você digita `go version`, aparece erro dizendo que o comando não foi encontrado.

**Causa:** O Go não está instalado ou não está no PATH.

**Solução:**
- **macOS:** Execute `brew install go`
- **Linux:** Siga os passos de instalação do Go acima e configure o PATH
- **Windows:** Reinstale o Go e verifique se marcou a opção "Add to PATH"

Depois, **feche e abra novamente o terminal** e teste: `go version`

---

#### "wails: command not found"

**Sintoma:** Quando você digita `wails version`, aparece erro.

**Causa:** O Wails não foi instalado corretamente ou o PATH do Go não foi configurado.

**Solução:**

1. Primeiro, verifique se o Go está funcionando: `go version`
2. Se estiver, execute: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
3. Verifique se o diretório de binários do Go está no PATH:

**macOS/Linux:**
```bash
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc  # ou ~/.bashrc
source ~/.zshrc  # ou source ~/.bashrc
```

**Windows (PowerShell como Administrador):**
```powershell
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";$env:USERPROFILE\go\bin", "User")
```

4. **Feche e abra novamente o terminal**
5. Teste: `wails version`

---

#### "node: command not found"

**Sintoma:** O Node.js não é reconhecido.

**Solução:**
- **macOS:** `brew install node`
- **Linux:** Siga os passos de instalação do Node acima
- **Windows:** Baixe e instale de https://nodejs.org/

**Depois de instalar, feche e abra novamente o terminal.**

---

#### Erro durante `npm install`

**Sintoma:** Ao executar `npm install` na pasta `frontend`, aparecem erros.

**Causas possíveis:**

1. **Versão muito antiga do Node.js**
   - Verifique: `node --version`
   - Deve ser >= 16.0.0
   - Solução: Atualize o Node.js

2. **Cache corrompido do npm**
   - Solução: `npm cache clean --force` e depois `npm install` novamente

3. **Problema de permissão (Linux/macOS)**
   - **NÃO USE SUDO**
   - Solução: Corrija as permissões:
     ```bash
     sudo chown -R $USER:$USER ~/.npm
     sudo chown -R $USER:$USER node_modules
     ```

---

### 🔴 Problemas ao Executar o Relief

#### Erro: "Failed to build frontend"

**Sintoma:** Ao executar `wails dev`, aparece erro de build do frontend.

**Solução:**

1. Entre na pasta frontend: `cd frontend`
2. Remova dependências antigas: `rm -rf node_modules package-lock.json`
3. Instale novamente: `npm install`
4. Volte para a raiz: `cd ..`
5. Tente novamente: `wails dev`

---

#### Erro: "Port already in use"

**Sintoma:** Mensagem dizendo que a porta já está sendo usada.

**Causa:** Outra instância do Relief ou outro programa está usando a mesma porta.

**Solução:**

**macOS/Linux:**
```bash
# Descubra qual processo está usando a porta (exemplo: porta 34115)
lsof -i :34115

# Mate o processo (substitua PID pelo número que apareceu)
kill -9 PID
```

**Windows (PowerShell):**
```powershell
# Descubra qual processo está usando a porta
netstat -ano | findstr :34115

# Mate o processo (substitua PID)
taskkill /PID <PID> /F
```

---

#### Relief abre mas a tela fica em branco

**Causa:** Problema com o WebView2 (Windows) ou bibliotecas gráficas (Linux).

**Solução:**

**Windows:**
- Instale/reinstale o WebView2: https://developer.microsoft.com/microsoft-edge/webview2/

**Linux:**
- Instale as bibliotecas necessárias:
  ```bash
  sudo apt-get install -y libgtk-3-dev libwebkit2gtk-4.0-dev
  ```

**Todos os sistemas:**
- Tente executar no modo de desenvolvimento com logs:
  ```bash
  wails dev -debug
  ```

---

### 🔴 Problemas com Projetos

#### Projeto não aparece depois de adicionar

**Sintoma:** Você seleciona a pasta mas o projeto não aparece no Relief.

**Causas possíveis:**

1. **Não existe arquivo `relief.yaml` na raiz**
   - Solução: Crie o arquivo `relief.yaml` conforme os exemplos acima

2. **Arquivo `relief.yaml` tem erros de sintaxe**
   - YAML é sensível a indentação
   - Use 2 espaços para indentar, não tabs
   - Verifique se não faltam aspas ou dois-pontos
   - Use um validador online: https://www.yamllint.com/

3. **Você selecionou a pasta errada**
   - Certifique-se de selecionar a pasta raiz do projeto, onde está o `relief.yaml`

---

#### Projeto não inicia (fica em "Iniciando...")

**Sintoma:** Você clica em "Iniciar" mas o projeto nunca muda para "Rodando".

**Diagnóstico:** Clique em "Ver Logs" para ver o que está acontecendo.

**Causas comuns:**

1. **Dependência não instalada**
   - Verifique nos logs se aparece "command not found"
   - Solução: Instale a dependência necessária (Node.js, Python, etc.)

2. **Comando `dev` está incorreto**
   - Solução: Teste o comando manualmente no terminal:
     ```bash
     cd /caminho/do/projeto
     npm run dev  # ou o comando que está no relief.yaml
     ```
   - Se funcionar manualmente, copie exatamente o comando para o `relief.yaml`

3. **Porta já está em uso**
   - Solução: Mude a porta no `env` do `relief.yaml` ou mate o processo que está usando a porta

4. **Faltam variáveis de ambiente**
   - Alguns projetos precisam de variáveis específicas
   - Solução: Adicione todas as variáveis necessárias na seção `env` do `relief.yaml`

---

#### "Failed to modify /etc/hosts" (macOS/Linux)

**Sintoma:** Erro ao tentar modificar o arquivo `/etc/hosts`.

**Causa:** O Relief precisa de permissões de administrador para modificar este arquivo.

**Solução:**
1. Digite sua senha quando solicitado
2. Se continuar com erro, modifique manualmente:
   ```bash
   sudo nano /etc/hosts
   ```
3. Adicione esta linha no final:
   ```
   127.0.0.1 meu-projeto.local.dev
   ```
   (substitua `meu-projeto` pelo seu domínio)
4. Salve (Ctrl+O, Enter, Ctrl+X)

---

#### Domínio .local.dev não funciona no navegador

**Sintoma:** Você acessa `http://projeto.local.dev` mas o navegador diz que não encontrou.

**Diagnóstico:**

1. **Verifique se o projeto está rodando**
   - O status deve estar 🟢 "Rodando" no Relief
   
2. **Verifique o arquivo hosts**
   - **macOS/Linux:** `cat /etc/hosts | grep local.dev`
   - **Windows:** `type C:\Windows\System32\drivers\etc\hosts | findstr local.dev`
   - Deve aparecer algo como: `127.0.0.1 projeto.local.dev`
   
3. **Verifique se o Traefik está rodando**
   - No Relief, deve haver um indicador de que o proxy está ativo
   
4. **Tente acessar pela porta diretamente**
   - Se seu projeto roda na porta 3000: `http://localhost:3000`
   - Se funcionar, o problema é no proxy/hosts
   
**Soluções:**

- Reinicie o Relief
- Reinicie o projeto
- No navegador, forçe atualização: `Ctrl+F5` (Windows/Linux) ou `Cmd+Shift+R` (macOS)
- Limpe o cache de DNS:
  - **macOS:** `sudo dscacheutil -flushcache; sudo killall -HUP mDNSResponder`
  - **Windows:** `ipconfig /flushdns`
  - **Linux:** `sudo systemd-resolve --flush-caches`

---

#### Projeto roda mas retorna erro 502/503

**Sintoma:** O domínio abre mas aparece "Bad Gateway" ou "Service Unavailable".

**Causa:** O Traefik está funcionando, mas o seu projeto não está respondendo corretamente na porta especificada.

**Solução:**

1. Verifique os logs do projeto no Relief
2. Certifique-se de que o projeto está escutando na porta correta:
   - Verifique a variável `PORT` no `env` do `relief.yaml`
   - Certifique-se de que seu código usa essa porta
3. Para projetos Node.js, certifique-se de usar `0.0.0.0` ao invés de `localhost`:
   ```javascript
   app.listen(port, '0.0.0.0', () => { ... })
   ```

---

### 🔴 Problemas de Performance

#### Relief está lento

**Sintoma:** A interface trava ou responde devagar.

**Causas e soluções:**

1. **Muitos projetos rodando ao mesmo tempo**
   - Solução: Pare projetos que você não está usando no momento

2. **Logs muito grandes**
   - Logs acumulam ao longo do tempo
   - Solução: Reinicie o projeto para limpar os logs

3. **Muita saída no console**
   - Se seu projeto imprime muitas mensagens, pode deixar o Relief lento
   - Solução: Reduza logs desnecessários no seu projeto

---

### 🔴 Problemas no Windows Especificamente

#### "The system cannot find the path specified"

**Causa:** Problemas com caminhos que contêm espaços ou caracteres especiais.

**Solução:**
- Evite espaços nos caminhos (use `C:\projetos` ao invés de `C:\Meus Projetos`)
- Se não puder evitar, use aspas nos comandos do `relief.yaml`:
  ```yaml
  scripts:
    dev: "\"C:\\Program Files\\node\\node.exe\" index.js"
  ```

---

#### Antivírus bloqueia o Relief

**Sintoma:** O Windows Defender ou outro antivírus diz que o Relief é suspeito.

**Causa:** É um falso positivo comum em apps desenvolvidos localmente.

**Solução:**
1. Adicione uma exceção para a pasta do Relief no seu antivírus
2. Se compilar o Relief, assine digitalmente o executável

---

### 🔴 Problemas no macOS Especificamente

#### "Relief.app is damaged and can't be opened"

**Sintoma:** O macOS impede de abrir o Relief compilado.

**Causa:** O macOS Gatekeeper bloqueia apps não assinados.

**Solução:**
```bash
xattr -cr /caminho/para/Relief.app
```

Ou:
1. Vá em **Preferências do Sistema**
2. **Segurança e Privacidade**
3. Clique em **"Abrir Mesmo Assim"**

---

#### "Permission denied" ao modificar /etc/hosts

**Solução:**
```bash
sudo chmod 644 /etc/hosts
```

---

### 📞 Ainda com Problemas?

Se nenhuma dessas soluções funcionou:

1. **Habilite o modo debug:**
   ```bash
   wails dev -debug
   ```
   Isso vai mostrar logs mais detalhados

2. **Verifique os logs do sistema:**
   - **macOS:** Abra o Console.app
   - **Linux:** `journalctl -f`
   - **Windows:** Visualizador de Eventos

3. **Abra uma issue no GitHub:**
   - Acesse: https://github.com/Maycon-Santos/relief/issues/new
   - Inclua:
     - Seu sistema operacional e versão
     - Output do comando `wails doctor`
     - Logs completos do erro
     - Passos para reproduzir o problema
     - Screenshots se possível

4. **Faça perguntas nas Discussions:**
   - https://github.com/Maycon-Santos/relief/discussions

---

## �🎨 Entendendo a Interface

### Card de Projeto

Cada projeto aparece em um card com:

- **Nome do Projeto**: O nome que você definiu no `relief.yaml`
- **Status**: 
  - 🔴 Parado
  - 🟡 Iniciando
  - 🟢 Rodando
  - 🔴 Erro
- **Domínio**: A URL para acessar o projeto
- **Botões**:
  - ▶️ **Iniciar**: Inicia o projeto
  - ⏹️ **Parar**: Para o projeto
  - 🔄 **Reiniciar**: Para e inicia novamente
  - 📋 **Ver Logs**: Abre os logs do projeto
  - 🗑️ **Remover**: Remove o projeto do Relief

### Painel de Logs

Quando você clica em "Ver Logs", uma janela abre mostrando:
- Todas as mensagens que o projeto está imprimindo
- Erros que aconteceram
- Informações de inicialização

### Controles Git

Se seu projeto é um repositório Git, você verá:
- **Branch atual**: Em qual branch você está trabalhando
- **Mudanças**: Se há arquivos modificados
- **Botão de Sincronizar**: Para fazer pull das últimas alterações

---

## 🌐 Como Funcionam os Domínios `.local.dev`

### O Que é um Reverse Proxy?

Normalmente seus projetos rodam em portas específicas:
- API: `http://localhost:3000`
- Frontend: `http://localhost:5173`
- Admin: `http://localhost:4000`

Isso é confuso! É difícil lembrar qual porta é de qual projeto.

O Relief usa uma ferramenta chamada **Traefik** (um reverse proxy) que funciona assim:

```
Você acessa: http://api.local.dev
      ↓
Traefik redireciona para: http://localhost:3000
      ↓
Você vê sua API! 🎉
```

### Como Isso Funciona?

O Relief faz duas coisas automaticamente:

1. **Configura o Traefik**: Cria regras de roteamento
2. **Modifica o arquivo `/etc/hosts`**: Diz ao seu computador que `.local.dev` é o próprio computador

**Nota**: No macOS/Linux, você precisará digitar sua senha quando o Relief tentar modificar o `/etc/hosts`.

---

## 📂 Estrutura do Projeto Relief

Se você quiser contribuir ou entender melhor o código:

```
relief/
├── main.go                 # Arquivo principal que inicia tudo
├── internal/               # Código principal do Relief
│   ├── app/                # Lógica que conecta Go com a interface
│   ├── config/             # Carrega e gerencia configurações
│   ├── domain/             # Modelos de projeto, manifesto, etc.
│   ├── runner/             # Código que inicia projetos (Native, Docker)
│   ├── dependency/         # Verifica se Node, Python, etc. estão instalados
│   ├── proxy/              # Gerencia Traefik e /etc/hosts
│   ├── git/                # Operações Git (branches, sync)
│   └── storage/            # Banco de dados SQLite para guardar projetos
├── pkg/                    # Utilitários reutilizáveis
│   ├── logger/             # Sistema de logs
│   ├── fileutil/           # Funções para mexer com arquivos
│   └── httputil/           # Funções para HTTP
├── frontend/               # Interface visual (React + TypeScript)
│   └── src/
│       ├── components/     # Componentes da interface (cards, botões, etc.)
│       ├── hooks/          # Hooks React customizados
│       ├── services/       # Comunicação com o backend Go
│       └── types/          # Tipos TypeScript
├── examples/               # Projetos de exemplo
│   └── hello-world/        # Projeto Node.js de exemplo
├── configs/                # Exemplos de configuração
├── docs/                   # Documentação adicional
└── build/                  # Scripts de build e binários compilados
```

---

## 🛠️ Desenvolvimento

### Executar Testes

```bash
# Executa todos os testes Go
./build/ci/test.sh
```

### Verificar Qualidade do Código

```bash
# Executa o linter Go
golangci-lint run
```

### Compilar para Produção

```bash
# Cria binários executáveis
./build/ci/build.sh
```

Os executáveis são criados em `build/bin/` para:
- macOS (Intel e Apple Silicon)
- Linux
- Windows

---

## 🤝 Como Contribuir

Adoraríamos ter sua contribuição! Aqui estão algumas formas de ajudar:

### 🐛 Encontrou um Bug?

1. Verifique se já não existe uma [issue aberta](https://github.com/Maycon-Santos/relief/issues)
2. Se não existe, [crie uma nova issue](https://github.com/Maycon-Santos/relief/issues/new)
3. Descreva o problema com o máximo de detalhes possível
4. Inclua screenshots se possível

### 💡 Tem uma Ideia?

1. Abra uma [issue](https://github.com/Maycon-Santos/relief/issues/new) descrevendo sua ideia
2. Aguarde feedback da comunidade
3. Se aprovado, implemente e envie um Pull Request

### 🔧 Quer Contribuir com Código?

1. Faça um fork do repositório
2. Crie uma branch para sua feature: `git checkout -b minha-feature`
3. Faça suas alterações
4. Commit suas mudanças: `git commit -m 'Adiciona nova feature X'`
5. Push para a branch: `git push origin minha-feature`
6. Abra um Pull Request

### 📚 Melhorar Documentação

Documentação sempre pode melhorar! Sinta-se livre para:
- Corrigir erros de digitação
- Adicionar exemplos
- Escrever tutoriais
- Traduzir para outros idiomas

Leia nosso [Guia de Contribuição](docs/contributing.md) para mais detalhes.

---

## 📚 Documentação Adicional

- **[Arquitetura](docs/architecture.md)** - Como o Relief funciona por dentro
- **[Schema do relief.yaml](docs/manifest-schema.md)** - Todas as opções disponíveis
- **[Guia de Contribuição](docs/contributing.md)** - Como colaborar com o projeto
- **[Exemplo Hello World](examples/hello-world/README.md)** - Tutorial prático

---

## ❓ Perguntas Frequentes

### O Relief funciona no Windows?

Sim! O Relief é multiplataforma e funciona em Windows, macOS e Linux.

### Preciso saber programar para usar?

Não necessariamente. Se você já trabalha com desenvolvimento e usa Node.js, Python, etc., você consegue usar o Relief. Este README foi feito para ensinar o básico.

### Meus projetos precisam estar na mesma pasta?

Não! Cada projeto pode estar em qualquer lugar do seu computador. Você apenas adiciona o caminho no Relief.

### O que acontece se eu fechar o Relief?

Todos os projetos que estavam rodando serão parados automaticamente.

### Posso usar portas normais ao invés de domínios `.local.dev`?

Sim! Você pode acessar diretamente pela porta também (ex: `localhost:3000`). Os domínios são apenas para facilitar.

---

## 📜 Licença

Este projeto está licenciado sob a **Licença MIT** - veja o arquivo [LICENSE](LICENSE) para detalhes.

Isso significa que você pode:
- ✅ Usar comercialmente
- ✅ Modificar
- ✅ Distribuir
- ✅ Uso privado

---

## 🙏 Agradecimentos

O Relief foi construído usando ferramentas incríveis:

- **[Wails](https://wails.io)** - Framework para criar apps desktop com Go e React
- **[Traefik](https://traefik.io)** - Reverse proxy moderno
- **[React](https://reactjs.org)** - Biblioteca para interfaces
- **[Tailwind CSS](https://tailwindcss.com)** - Framework CSS
- **[shadcn/ui](https://ui.shadcn.com)** - Componentes de UI

E especialmente a todos os [contribuidores](https://github.com/Maycon-Santos/relief/graphs/contributors) que ajudaram a melhorar o projeto! ❤️

---

## 📞 Suporte

- **Issues**: [GitHub Issues](https://github.com/Maycon-Santos/relief/issues)
- **Discussões**: [GitHub Discussions](https://github.com/Maycon-Santos/relief/discussions)

---

<p align="center">
  <b>Feito com ❤️ pela comunidade Relief</b>
</p>

<p align="center">
  <sub>⭐ Se você gostou, deixe uma estrela no repositório!</sub>
</p>
