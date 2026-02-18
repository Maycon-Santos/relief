# 📦 Guia Completo de Instalação do Relief

Este guia fornece instruções detalhadas de instalação para cada sistema operacional.

## 📑 Índice

- [⚡ Instalação Rápida (Binário Pré-compilado)](#-instalação-rápida-binário-pré-compilado) **← Recomendado!**
  - [macOS](#macos-binário)
  - [Linux](#linux-binário)
  - [Windows](#windows-binário)
- [🛠️ Instalação via Código-fonte](#️-instalação-via-código-fonte)
  - [🍎 macOS](#-instalação-no-macos)
  - [🐧 Linux](#-instalação-no-linux)
  - [🪟 Windows](#-instalação-no-windows)
- [✅ Verificação da Instalação](#-verificação-da-instalação)
- [🔄 Atualizando o Relief](#-atualizando-o-relief)
- [🗑️ Desinstalando](#-desinstalando)

---

## ⚡ Instalação Rápida (Binário Pré-compilado)

**Esta é a forma mais fácil e rápida de instalar o Relief!**

Baixe o binário pré-compilado da [página de releases](https://github.com/Maycon-Santos/relief/releases/latest) e adicione ao PATH do seu sistema.

### macOS (Binário)

#### Opção 1: Download e Instalação Automática (Recomendado)

```bash
# Para Intel (x86_64)
curl -L https://github.com/Maycon-Santos/relief/releases/latest/download/relief-darwin-amd64.tar.gz | tar xz
sudo mv Relief.app /Applications/
echo 'export PATH="/Applications/Relief.app/Contents/MacOS:$PATH"' >> ~/.zshrc
source ~/.zshrc

# Para Apple Silicon (M1/M2/M3)
curl -L https://github.com/Maycon-Santos/relief/releases/latest/download/relief-darwin-arm64.tar.gz | tar xz
sudo mv Relief.app /Applications/
echo 'export PATH="/Applications/Relief.app/Contents/MacOS:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

#### Opção 2: Download Manual

1. **Baixe o binário:**
   - Intel: [relief-darwin-amd64.tar.gz](https://github.com/Maycon-Santos/relief/releases/latest/download/relief-darwin-amd64.tar.gz)
   - Apple Silicon: [relief-darwin-arm64.tar.gz](https://github.com/Maycon-Santos/relief/releases/latest/download/relief-darwin-arm64.tar.gz)

2. **Extraia e instale:**
   ```bash
   # Vá até a pasta de Downloads
   cd ~/Downloads
   
   # Extraia o arquivo
   tar -xzf relief-darwin-*.tar.gz
   
   # Mova para Applications
   sudo mv Relief.app /Applications/
   ```

3. **Adicione ao PATH:**
   ```bash
   # Para zsh (padrão no macOS moderno)
   echo 'export PATH="/Applications/Relief.app/Contents/MacOS:$PATH"' >> ~/.zshrc
   source ~/.zshrc
   
   # Para bash
   echo 'export PATH="/Applications/Relief.app/Contents/MacOS:$PATH"' >> ~/.bash_profile
   source ~/.bash_profile
   ```

4. **Remova a quarentena do macOS:**
   ```bash
   xattr -cr /Applications/Relief.app
   ```

5. **Verifique:**
   ```bash
   relief --version
   ```

---

### Linux (Binário)

#### Opção 1: Download e Instalação Automática (Recomendado)

```bash
# Para x86_64 (AMD64)
curl -L https://github.com/Maycon-Santos/relief/releases/latest/download/relief-linux-amd64.tar.gz | tar xz
sudo mv relief /usr/local/bin/
sudo chmod +x /usr/local/bin/relief

# Para ARM64
curl -L https://github.com/Maycon-Santos/relief/releases/latest/download/relief-linux-arm64.tar.gz | tar xz
sudo mv relief /usr/local/bin/
sudo chmod +x /usr/local/bin/relief
```

#### Opção 2: Download Manual

1. **Baixe o binário:**
   - x86_64: [relief-linux-amd64.tar.gz](https://github.com/Maycon-Santos/relief/releases/latest/download/relief-linux-amd64.tar.gz)
   - ARM64: [relief-linux-arm64.tar.gz](https://github.com/Maycon-Santos/relief/releases/latest/download/relief-linux-arm64.tar.gz)

2. **Extraia e instale:**
   ```bash
   # Vá até a pasta de downloads
   cd ~/Downloads
   
   # Extraia
   tar -xzf relief-linux-*.tar.gz
   
   # Mova para o PATH do sistema
   sudo mv relief /usr/local/bin/
   
   # Dê permissão de execução
   sudo chmod +x /usr/local/bin/relief
   ```

3. **Verifique:**
   ```bash
   relief --version
   ```

**Alternativa:** Instalar no diretório do usuário (sem sudo):

```bash
# Crie o diretório bin no seu home
mkdir -p ~/.local/bin

# Extraia e mova
tar -xzf relief-linux-*.tar.gz
mv relief ~/.local/bin/
chmod +x ~/.local/bin/relief

# Adicione ao PATH
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc

# Verifique
relief --version
```

---

### Windows (Binário)

#### Opção 1: PowerShell (Recomendado)

```powershell
# Crie uma pasta para o Relief
New-Item -ItemType Directory -Force -Path "$env:LOCALAPPDATA\Relief"

# Baixe o binário
Invoke-WebRequest -Uri "https://github.com/Maycon-Santos/relief/releases/latest/download/relief-windows-amd64.exe.zip" -OutFile "$env:TEMP\relief.zip"

# Extraia
Expand-Archive -Path "$env:TEMP\relief.zip" -DestinationPath "$env:LOCALAPPDATA\Relief" -Force

# Adicione ao PATH
$oldPath = [Environment]::GetEnvironmentVariable("Path", "User")
$newPath = "$oldPath;$env:LOCALAPPDATA\Relief"
[Environment]::SetEnvironmentVariable("Path", $newPath, "User")

# Limpe arquivos temporários
Remove-Item "$env:TEMP\relief.zip"

# Recarregue o PATH na sessão atual
$env:Path = [Environment]::GetEnvironmentVariable("Path", "User")

Write-Host "✅ Relief instalado com sucesso!"
Write-Host "Feche e abra um novo PowerShell para usar o comando 'relief'"
```

#### Opção 2: Download Manual

1. **Baixe o binário:**
   - [relief-windows-amd64.exe.zip](https://github.com/Maycon-Santos/relief/releases/latest/download/relief-windows-amd64.exe.zip)

2. **Extraia o arquivo:**
   - Clique com botão direito no arquivo `.zip`
   - Selecione **"Extrair Tudo..."**
   - Escolha uma pasta (ex: `C:\Program Files\Relief`)

3. **Adicione ao PATH:**
   
   **Via Interface Gráfica:**
   1. Pressione `Win + R`, digite `sysdm.cpl` e pressione Enter
   2. Vá na aba **"Avançado"**
   3. Clique em **"Variáveis de Ambiente"**
   4. Em **"Variáveis do usuário"**, selecione **"Path"** e clique em **"Editar"**
   5. Clique em **"Novo"**
   6. Cole o caminho da pasta onde está o `relief.exe` (ex: `C:\Program Files\Relief`)
   7. Clique em **"OK"** em todas as janelas
   
   **Via PowerShell (como Administrador):**
   ```powershell
   # Substitua o caminho pela pasta onde você extraiu
   $reliefPath = "C:\Program Files\Relief"
   $oldPath = [Environment]::GetEnvironmentVariable("Path", "User")
   [Environment]::SetEnvironmentVariable("Path", "$oldPath;$reliefPath", "User")
   ```

4. **Verifique:**
   
   Abra um **novo** PowerShell ou Prompt de Comando:
   ```powershell
   relief --version
   ```

---

## 🛠️ Instalação via Código-fonte

Se você preferir compilar do código-fonte (para desenvolvimento ou customização), siga as instruções abaixo.

---

## 🍎 Instalação no macOS

### Pré-requisitos

- macOS 10.15 (Catalina) ou superior
- Pelo menos 2GB de espaço livre em disco
- Acesso à internet

### Passo 1: Instale o Homebrew

O Homebrew é um gerenciador de pacotes que facilita a instalação de ferramentas no macOS.

**Verifique se você já tem:**
```bash
brew --version
```

Se o comando funcionar, você já tem Homebrew. Pule para o Passo 2.

**Se não tiver, instale:**
```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

O instalador pedirá:
1. Sua senha do macOS (não aparece nada na tela ao digitar, é normal)
2. Pressione **Enter** para confirmar

Aguarde a instalação (pode levar 5-10 minutos).

Após terminar, você pode ver uma mensagem para adicionar o Homebrew ao PATH:

**Para Mac Intel (x86_64):**
```bash
echo 'eval "$(/usr/local/bin/brew shellenv)"' >> ~/.zprofile
eval "$(/usr/local/bin/brew shellenv)"
```

**Para Mac Apple Silicon (M1/M2/M3):**
```bash
echo 'eval "$(/opt/homebrew/bin/brew shellenv)"' >> ~/.zprofile
eval "$(/opt/homebrew/bin/brew shellenv)"
```

Verifique se funcionou:
```bash
brew --version
# Deve mostrar: Homebrew 4.x.x
```

---

### Passo 2: Instale o Go

```bash
# Instale o Go
brew install go

# Verifique a instalação
go version
# Deve mostrar: go version go1.22.x darwin/amd64 (ou darwin/arm64 para M1/M2/M3)
```

---

### Passo 3: Instale o Node.js

```bash
# Instale o Node.js e npm
brew install node

# Verifique a instalação
node --version
# Deve mostrar: v20.x.x ou superior

npm --version
# Deve mostrar: 10.x.x ou superior
```

---

### Passo 4: Configure o PATH do Go

```bash
# Descubra qual shell você usa
echo $SHELL

# Se for /bin/zsh (padrão no macOS moderno):
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc

# Se for /bin/bash:
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bash_profile
source ~/.bash_profile
```

---

### Passo 5: Instale o Wails

```bash
# Instale o Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Verifique a instalação
wails version
# Deve mostrar: Wails v2.x.x
```

Se o comando `wails` não for encontrado, feche o terminal e abra novamente.

---

### Passo 6: Instale o Git (se necessário)

O macOS geralmente já vem com Git instalado.

```bash
# Verifique se você tem o Git
git --version

# Se não tiver, instale:
brew install git
```

---

### Passo 7: Baixe o Relief

```bash
# Navegue até onde você quer guardar o Relief
cd ~/Developer  # ou qualquer pasta de sua preferência
# Se a pasta não existe, crie: mkdir -p ~/Developer

# Clone o repositório
git clone https://github.com/Maycon-Santos/relief.git

# Entre na pasta
cd relief
```

---

### Passo 8: Instale as Dependências

```bash
# Baixe as bibliotecas Go
go mod download

# Entre na pasta do frontend
cd frontend

# Instale as dependências do Node
npm install

# Volte para a pasta principal
cd ..
```

Isso pode levar alguns minutos na primeira vez.

---

### Passo 9: Execute o Relief

```bash
# Inicie em modo desenvolvimento
wails dev
```

Uma janela vai abrir com a interface do Relief! 🎉

**Primeira execução:**
- O macOS pode pedir permissão para o Relief acessar arquivos
- Clique em **"OK"** ou **"Permitir"**
- Pode pedir sua senha para modificar o arquivo `/etc/hosts` - isso é normal e seguro

---

### Compilar uma Versão Standalone (Opcional)

Se você quiser criar um app Relief.app que pode ser executado sem o terminal:

```bash
# Compile o aplicativo
wails build

# O aplicativo estará em:
# build/bin/Relief.app

# Para executar:
open build/bin/Relief.app
```

---

## 🐧 Instalação no Linux

### Ubuntu/Debian

#### Requisitos
- Ubuntu 20.04+ ou Debian 11+
- Pelo menos 2GB de espaço livre
- Acesso sudo

#### Passo a Passo Completo

**1. Atualize o sistema**

```bash
sudo apt update && sudo apt upgrade -y
```

**2. Instale o Git**

```bash
# Instale o Git
sudo apt install -y git

# Verifique
git --version
```

**3. Instale o Go**

```bash
# Remova instalações antigas (se houver)
sudo rm -rf /usr/local/go

# Baixe o Go (verifique a versão mais recente em https://go.dev/dl/)
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz

# Extraia
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz

# Limpe o arquivo de instalação
rm go1.22.0.linux-amd64.tar.gz

# Configure o PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc

# Verifique
go version
# Deve mostrar: go version go1.22.0 linux/amd64
```

**4. Instale o Node.js**

```bash
# Adicione o repositório do Node.js (versão 18 LTS)
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -

# Instale o Node.js
sudo apt install -y nodejs

# Verifique
node --version  # Deve mostrar v18.x.x ou maior
npm --version   # Deve mostrar 9.x.x ou maior
```

**5. Instale dependências do sistema para Wails**

```bash
sudo apt install -y \
  gcc \
  g++ \
  libgtk-3-dev \
  libwebkit2gtk-4.1-dev \
  build-essential \
  pkg-config
```

**6. Instale o Wails**

```bash
# Instale o Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Verifique
wails version
# Deve mostrar: Wails v2.x.x
```

Se o comando não for encontrado, feche e abra o terminal novamente.

**7. Baixe o Relief**

```bash
# Navegue até onde quer guardar
cd ~  # ou outra pasta de preferência

# Clone o repositório
git clone https://github.com/Maycon-Santos/relief.git

# Entre na pasta
cd relief
```

**8. Instale as dependências**

```bash
# Dependências do Go
go mod download

# Dependências do frontend
cd frontend
npm install
cd ..
```

**9. Execute o Relief**

```bash
wails dev
```

Uma janela vai abrir com o Relief! 🎉

---

### Fedora/RHEL/CentOS

**1. Atualize o sistema**

```bash
sudo dnf update -y
```

**2. Instale o Git**

```bash
sudo dnf install -y git
```

**3. Instale o Go**

```bash
# Baixe o Go
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz

# Extraia
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz

# Configure o PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc

# Verifique
go version
```

**4. Instale o Node.js**

```bash
# Adicione o repositório
curl -fsSL https://rpm.nodesource.com/setup_18.x | sudo bash -

# Instale
sudo dnf install -y nodejs
```

**5. Instale dependências para Wails**

```bash
sudo dnf install -y \
  gcc \
  gcc-c++ \
  gtk3-devel \
  webkit2gtk3-devel
```

**6. Continue a partir do passo 6 da seção Ubuntu/Debian**

---

### Arch Linux

**1. Atualize o sistema**

```bash
sudo pacman -Syu
```

**2. Instale as ferramentas necessárias**

```bash
sudo pacman -S go nodejs npm git gcc gtk3 webkit2gtk base-devel
```

**3. Configure o PATH do Go**

```bash
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc
```

**4. Instale o Wails**

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

**5. Continue a partir do passo 7 da seção Ubuntu/Debian**

---

## 🪟 Instalação no Windows

### Requisitos
- Windows 10 versão 1903 ou superior / Windows 11
- Pelo menos 2GB de espaço livre
- Acesso de administrador (para algumas etapas)

### Método de Instalação

Vamos usar o PowerShell para a instalação.

#### Passo 1: Abra o PowerShell

1. Pressione `Win + X`
2. Clique em **"Windows PowerShell"** ou **"Terminal"**

> **Nota:** Alguns comandos precisarão de privilégios de administrador. Quando necessário, feche o PowerShell e abra como administrador (botão direito → "Executar como administrador").

---

#### Passo 2: Instale o Chocolatey (Gerenciador de Pacotes)

O Chocolatey facilita a instalação de ferramentas no Windows.

**Abra o PowerShell como Administrador** e execute:

```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
```

Após a instalação, feche e abra novamente o PowerShell como Administrador.

Verifique:
```powershell
choco --version
```

**Alternativa sem Chocolatey:** Você pode instalar manualmente cada ferramenta baixando dos sites oficiais (instruções abaixo).

---

#### Passo 3: Instale o Go

**Opção A: Com Chocolatey (Recomendado)**

```powershell
choco install golang -y
```

**Opção B: Instalação Manual**

1. Acesse: https://go.dev/dl/
2. Baixe o arquivo **Windows installer** (`.msi`)
3. Execute o instalador
4. Mantenha todas as opções padrão
5. Clique em **Next** → **Install** → **Finish**

**Verifique:**

Abra um **novo** PowerShell:
```powershell
go version
# Deve mostrar: go version go1.22.x windows/amd64
```

---

#### Passo 4: Instale o Node.js

**Opção A: Com Chocolatey**

```powershell
choco install nodejs-lts -y
```

**Opção B: Instalação Manual**

1. Acesse: https://nodejs.org/
2. Baixe a versão **LTS** (Long Term Support)
3. Execute o instalador
4. Mantenha todas as opções padrão marcadas
5. Complete a instalação

**Verifique:**

Abra um **novo** PowerShell:
```powershell
node --version  # Deve mostrar v18.x.x ou superior
npm --version   # Deve mostrar 9.x.x ou superior
```

---

#### Passo 5: Instale o Git

**Opção A: Com Chocolatey**

```powershell
choco install git -y
```

**Opção B: Instalação Manual**

1. Acesse: https://git-scm.com/download/win
2. Baixe o instalador
3. Execute e mantenha as opções padrão
4. Complete a instalação

**Verifique:**

Abra um **novo** PowerShell:
```powershell
git --version
```

---

#### Passo 6: Configure o PATH do Go

No PowerShell (não precisa ser como administrador):

```powershell
# Adicione o diretório de binários do Go ao PATH do usuário
$goPath = "$env:USERPROFILE\go\bin"
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";$goPath", "User")
```

**Feche e abra novamente o PowerShell** para as mudanças terem efeito.

---

#### Passo 7: Instale o Wails

```powershell
# Instale o Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Verifique (pode precisar fechar e abrir o PowerShell novamente)
wails version
# Deve mostrar: Wails v2.x.x
```

---

#### Passo 8: Instale/Verifique o WebView2

O WebView2 é necessário para renderizar a interface do Relief.

**Windows 11:** Já vem instalado por padrão.

**Windows 10:** 
1. Verifique se já tem: Vá em **"Configurações → Apps → Apps e recursos"** e procure por "WebView2"
2. Se não tiver, baixe de: https://developer.microsoft.com/microsoft-edge/webview2/
3. Baixe o **"Evergreen Standalone Installer"**
4. Execute e instale

---

#### Passo 9: Baixe o Relief

No PowerShell:

```powershell
# Navegue até onde quer guardar (exemplo: Documentos)
cd $env:USERPROFILE\Documents

# Clone o repositório
git clone https://github.com/Maycon-Santos/relief.git

# Entre na pasta
cd relief
```

---

#### Passo 10: Instale as Dependências

```powershell
# Dependências do Go
go mod download

# Entre na pasta do frontend
cd frontend

# Instale dependências do Node
npm install

# Volte para a pasta principal
cd ..
```

Isso pode levar alguns minutos.

---

#### Passo 11: Execute o Relief

```powershell
wails dev
```

Uma janela vai abrir com o Relief! 🎉

**Na primeira execução:**
- O Windows Defender pode perguntar se quer permitir o app
- Clique em **"Permitir acesso"**

---

### Compilar uma Versão Standalone (Opcional)

```powershell
# Compile o aplicativo
wails build

# O executável estará em:
# build\bin\Relief.exe

# Para executar:
.\build\bin\Relief.exe
```

---

## ✅ Verificação da Instalação

Após seguir os passos acima, você pode verificar se tudo está instalado corretamente:

### Todas as Plataformas

Execute estes comandos no terminal:

```bash
# Verifique o Go
go version
# Esperado: go version go1.22.x ...

# Verifique o Node
node --version
# Esperado: v18.x.x ou superior

npm --version
# Esperado: 9.x.x ou superior

# Verifique o Git
git --version
# Esperado: git version 2.x.x

# Verifique o Wails
wails version
# Esperado: Wails v2.x.x

# Verifique o Wails Doctor
wails doctor
```

O comando `wails doctor` mostra um relatório completo de todas as dependências. Idealmente, tudo deve estar marcado como ✓ (check).

---

## 🔄 Atualizando o Relief

### Se você instalou via binário pré-compilado:

#### macOS

```bash
# Para Intel (x86_64)
curl -L https://github.com/Maycon-Santos/relief/releases/latest/download/relief-darwin-amd64.tar.gz | tar xz
sudo mv Relief.app /Applications/
xattr -cr /Applications/Relief.app

# Para Apple Silicon (M1/M2/M3)
curl -L https://github.com/Maycon-Santos/relief/releases/latest/download/relief-darwin-arm64.tar.gz | tar xz
sudo mv Relief.app /Applications/
xattr -cr /Applications/Relief.app
```

#### Linux

```bash
# Para x86_64
curl -L https://github.com/Maycon-Santos/relief/releases/latest/download/relief-linux-amd64.tar.gz | tar xz
sudo mv relief /usr/local/bin/
sudo chmod +x /usr/local/bin/relief

# Para ARM64
curl -L https://github.com/Maycon-Santos/relief/releases/latest/download/relief-linux-arm64.tar.gz | tar xz
sudo mv relief /usr/local/bin/
sudo chmod +x /usr/local/bin/relief

# Verifique a versão
relief --version
```

#### Windows (PowerShell)

```powershell
# Baixe a nova versão
Invoke-WebRequest -Uri "https://github.com/Maycon-Santos/relief/releases/latest/download/relief-windows-amd64.exe.zip" -OutFile "$env:TEMP\relief.zip"

# Extraia sobrescrevendo a versão antiga
Expand-Archive -Path "$env:TEMP\relief.zip" -DestinationPath "$env:LOCALAPPDATA\Relief" -Force

# Limpe
Remove-Item "$env:TEMP\relief.zip"

Write-Host "✅ Relief atualizado com sucesso!"
```

---

### Se você instalou via código-fonte:

```bash
# Entre na pasta do Relief
cd caminho/para/relief

# Baixe as últimas alterações
git pull origin main

# Atualize dependências do Go
go mod download

# Atualize dependências do frontend
cd frontend
npm install
cd ..

# Execute novamente
wails dev
```

---

## 🗑️ Desinstalando

### Remover o Relief

```bash
# Simplesmente delete a pasta
rm -rf caminho/para/relief  # macOS/Linux
rmdir /s caminho\para\relief  # Windows
```

### Remover as Ferramentas

Se você quiser remover completamente todas as ferramentas instaladas:

#### macOS

```bash
# Desinstalar com Homebrew
brew uninstall go node git

# Remover o Wails
rm $(which wails)

# Limpar cache do Go
rm -rf ~/go

# (Opcional) Remover Homebrew completamente
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/uninstall.sh)"
```

#### Linux (Ubuntu/Debian)

```bash
# Remover pacotes
sudo apt remove nodejs go git

# Remover dependências órfãs
sudo apt autoremove

# Remover o Go instalado manualmente (se foi instalado assim)
sudo rm -rf /usr/local/go

# Remover o Wails
rm $(which wails)

# Limpar cache do Go
rm -rf ~/go
```

#### Windows

**Com Chocolatey:**
```powershell
choco uninstall golang nodejs git -y
```

**Manualmente:**
1. Vá em "Configurações → Apps → Apps e recursos"
2. Procure e desinstale: "Go", "Node.js", "Git"

**Remover o Wails:**
```powershell
Remove-Item "$env:USERPROFILE\go\bin\wails.exe"
```

---

## 📞 Precisa de Ajuda?

Se você encontrou algum problema durante a instalação:

1. **Verifique a seção** [🆘 Problemas Comuns](README.md#-problemas-comuns-e-soluções) no README
2. **Execute** `wails doctor` para diagnóstico
3. **Abra uma issue** em: https://github.com/Maycon-Santos/relief/issues
4. **Pergunte na comunidade**: https://github.com/Maycon-Santos/relief/discussions

---

## 🎓 Próximos Passos

Agora que você tem o Relief instalado:

1. 📖 Leia o [Guia de Primeiro Uso](README.md#-primeiro-uso)
2. 🚀 Veja o [Guia Rápido](QUICKSTART.md)
3. 🔧 Aprenda a [Adicionar Seus Projetos](README.md#-como-adicionar-seus-próprios-projetos)

---

<p align="center">
  <b>Instalação completa! Agora você está pronto para usar o Relief! 🎉</b>
</p>
