# 🆘 Guia de Solução de Problemas do Relief

Este guia ajuda você a resolver problemas comuns ao instalar, configurar e usar o Relief.

## 📑 Índice

- [🔴 Problemas de Instalação](#-problemas-de-instalação)
- [🔴 Problemas ao Executar o Relief](#-problemas-ao-executar-o-relief)
- [🔴 Problemas com Projetos](#-problemas-com-projetos)
- [🔴 Problemas de Rede e Domínios](#-problemas-de-rede-e-domínios)
- [🔴 Problemas de Performance](#-problemas-de-performance)
- [🔴 Problemas Específicos por Sistema](#-problemas-específicos-por-sistema)
- [🔧 Ferramentas de Diagnóstico](#-ferramentas-de-diagnóstico)
- [📞 Obtendo Suporte](#-obtendo-suporte)

---

## 🔴 Problemas de Instalação

### "go: command not found"

**📝 Sintoma:**  
Ao executar `go version`, aparece erro: "command not found" ou "não é reconhecido".

**🔍 Causa:**  
O Go não está instalado ou não foi adicionado ao PATH do sistema.

**✅ Solução:**

**macOS:**
```bash
brew install go
# Após instalar, feche e abra o terminal
go version
```

**Linux:**
```bash
# Baixe e instale o Go
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz

# Configure o PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Teste
go version
```

**Windows:**
1. Baixe o instalador de https://go.dev/dl/
2. Execute e certifique-se de marcar "Add to PATH"
3. **Importante:** Feche e abra um novo PowerShell
4. Teste: `go version`

---

### "wails: command not found"

**📝 Sintoma:**  
Ao executar `wails version`, aparece "command not found".

**🔍 Causa:**  
O Wails não foi instalado ou o PATH do Go não inclui a pasta de binários.

**✅ Solução:**

**Passo 1: Verifique o Go**
```bash
go version  # Deve funcionar
```

**Passo 2: Instale o Wails**
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

**Passo 3: Configure o PATH**

**macOS/Linux (zsh):**
```bash
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc
```

**macOS/Linux (bash):**
```bash
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc
```

**Windows (PowerShell como Administrador):**
```powershell
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";$env:USERPROFILE\go\bin", "User")
```

**Passo 4: Reinicie o Terminal**  
**IMPORTANTE:** Sempre feche e abra um novo terminal após mudanças no PATH.

**Passo 5: Teste**
```bash
wails version
```

---

### "node: command not found" ou "npm: command not found"

**📝 Sintoma:**  
Node.js ou npm não são reconhecidos.

**✅ Solução:**

**macOS:**
```bash
brew install node
node --version && npm --version
```

**Linux (Ubuntu/Debian):**
```bash
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs
node --version && npm --version
```

**Windows:**
1. Baixe de https://nodejs.org/ (versão LTS)
2. Instale com as opções padrão
3. **Feche e abra novo terminal**
4. Teste: `node --version`

---

### Erro durante "npm install"

**📝 Sintoma:**  
Ao executar `npm install` na pasta `frontend`, aparecem erros diversos.

**🔍 Causas Possíveis:**

#### 1. Versão antiga do Node.js

```bash
# Verifique sua versão
node --version

# Deve ser >= v18.0.0
# Se for menor, atualize o Node.js
```

#### 2. Cache corrompido do npm

```bash
# Limpe o cache
npm cache clean --force

# Tente instalar novamente
npm install
```

#### 3. Problemas de permissão (Linux/macOS)

**❌ NÃO USE `sudo npm install`**

```bash
# Corrija as permissões
sudo chown -R $USER:$USER ~/.npm
sudo chown -R $USER:$USER node_modules

# Tente novamente SEM sudo
npm install
```

#### 4. Dependências do sistema faltando (Linux)

```bash
# Ubuntu/Debian
sudo apt-get install -y build-essential

# Fedora
sudo dnf install gcc-c++ make
```

---

### "gcc: command not found" (Linux)

**📝 Sintoma:**  
Ao instalar o Wails ou executar `wails dev`, aparece erro relacionado ao gcc.

**✅ Solução:**

```bash
# Ubuntu/Debian
sudo apt-get install -y gcc g++ build-essential

# Fedora/RHEL/CentOS
sudo dnf install gcc gcc-c++

# Arch
sudo pacman -S gcc
```

---

### "Package webkit2gtk-4.0 was not found" (Linux)

**📝 Sintoma:**  
Erro ao tentar executar `wails dev` no Linux.

**✅ Solução:**

```bash
# Ubuntu/Debian
sudo apt-get install -y libgtk-3-dev libwebkit2gtk-4.0-dev

# Fedora
sudo dnf install gtk3-devel webkit2gtk3-devel

# Arch
sudo pacman -S gtk3 webkit2gtk
```

---

## 🔴 Problemas ao Executar o Relief

### "Failed to build frontend"

**📝 Sintoma:**  
Ao executar `wails dev`, aparece erro dizendo que falhou ao fazer build do frontend.

**✅ Solução:**

```bash
# Entre na pasta do frontend
cd frontend

# Remova instalação anterior
rm -rf node_modules package-lock.json

# Reinstale do zero
npm install

# Volte para a raiz
cd ..

# Tente novamente
wails dev
```

Se o erro persistir:

```bash
# Limpe tudo
cd frontend
rm -rf node_modules package-lock.json dist .vite
npm cache clean --force
npm install
cd ..

wails dev
```

---

### "Port already in use" / "Porta já está em uso"

**📝 Sintoma:**  
Mensagem dizendo que a porta já está sendo usada.

**🔍 Causa:**  
Outra instância do Relief ou outro programa está usando a porta.

**✅ Solução:**

**Descubra qual porta está em uso (geralmente será exibida na mensagem de erro)**

**macOS/Linux:**
```bash
# Exemplo: se a porta for 34115
lsof -i :34115

# Isso mostrará algo como:
# COMMAND    PID   USER
# wails    12345  Maycon-Santos

# Mate o processo
kill -9 12345  # substitua 12345 pelo PID real
```

**Windows (PowerShell):**
```powershell
# Descubra o processo
netstat -ano | findstr :34115

# Isso mostrará algo como:
# TCP    0.0.0.0:34115    0.0.0.0:0    LISTENING    12345

# Mate o processo (substitua 12345 pelo PID real)
taskkill /PID 12345 /F
```

---

### Relief abre mas a tela fica em branco

**📝 Sintoma:**  
A janela do Relief abre mas não mostra nada, fica branca ou preta.

**🔍 Causa:**  
Problema com o WebView2 (Windows) ou with bibliotecas gráficas (Linux).

**✅ Solução:**

**Windows:**
1. Verifique se tem o WebView2 instalado:
   - Vá em "Configurações → Apps → Apps e recursos"
   - Procure por "WebView2"
2. Se não tiver, baixe de: https://developer.microsoft.com/microsoft-edge/webview2/
3. Instale o "Evergreen Standalone Installer"

**Linux:**
```bash
# Ubuntu/Debian
sudo apt-get install -y libgtk-3-dev libwebkit2gtk-4.0-dev

# Fedora
sudo dnf install gtk3-devel webkit2gtk3-devel
```

**Todos os sistemas:**

Tente executar em modo debug para ver os logs:
```bash
wails dev -debug
```

Verifique os logs para mais detalhes.

---

### "Failed to start Traefik" / Erro com proxy

**📝 Sintoma:**  
Mensagem de erro relacionada ao Traefik ou proxy.

**✅ Solução:**

1. Verifique se a porta 80 está livre:

```bash
# macOS/Linux
sudo lsof -i :80

# Windows
netstat -ano | findstr :80
```

2. Se já houver algo na porta 80 (como Apache, Nginx), você precisa pará-lo temporariamente:

```bash
# macOS
sudo apachectl stop

# Linux
sudo systemctl stop apache2  # ou nginx

# Windows
# Pare pelo Gerenciador de Serviços
```

---

## 🔴 Problemas com Projetos

### Projeto não aparece após adicionar

**📝 Sintoma:**  
Você seleciona a pasta do projeto mas ele não aparece na interface.

**🔍 Diagnóstico:**

1. **Verifique se existe o arquivo `relief.yaml` na raiz da pasta**

```bash
# No terminal, entre na pasta do projeto
cd /caminho/do/projeto

# Liste os arquivos
ls -la  # macOS/Linux
dir     # Windows

# Deve aparecer: relief.yaml
```

2. **Verifique se o YAML está válido**

```bash
# Veja o conteúdo
cat relief.yaml  # macOS/Linux
type relief.yaml  # Windows
```

**✅ Soluções Comuns:**

#### ❌ Falta o arquivo relief.yaml
Crie o arquivo `relief.yaml` na pasta raiz do projeto. Veja exemplos em: [Como Adicionar Seus Projetos](README.md#-como-adicionar-seus-próprios-projetos)

#### ❌ YAML com erros de sintaxe

YAML é muito sensível a indentação!

**Errado:**
```yaml
name: "meu-projeto"
  domain: "projeto.local.dev"  # ❌ indentação incorreta
```

**Correto:**
```yaml
name: "meu-projeto"
domain: "projeto.local.dev"  # ✅ mesma indentação
```

**Dicas:**
- Use **2 espaços** para indentar (não tabs)
- Sempre use aspas em strings: `"meu-valor"`
- Dois pontos devem ser seguidos de espaço: `chave: valor` (não `chave:valor`)

Valide seu YAML em: https://www.yamllint.com/

#### ❌ Selecionou a pasta errada

Certifique-se de selecionar a pasta raiz do projeto (onde está o `relief.yaml`), não uma subpasta.

---

### Projeto não inicia (fica em "Iniciando...")

**📝 Sintoma:**  
Você clica em "Iniciar" mas o projeto nunca muda para "Rodando".

**🔍 Diagnóstico:**

**Passo 1:** Clique em **"Ver Logs"** no card do projeto  
Os logs vão mostrar o que está acontecendo.

**✅ Soluções por erro comum:**

#### ❌ "command not found" nos logs

**Causa:** A dependência necessária não está instalada.

**Exemplo de log:**
```
npm: command not found
```

**Solução:** Instale a dependência (Node.js, Python, etc.)

#### ❌ "EADDRINUSE" ou "Address already in use"

**Causa:** A porta definida já está em uso por outro programa.

**Solução:**

Opção 1 - Mude a porta no `relief.yaml`:
```yaml
env:
  PORT: "3001"  # mude para uma porta diferente
```

Opção 2 - Mate o processo que está usando a porta:
```bash
# macOS/Linux
lsof -i :3000
kill -9 PID

# Windows
netstat -ano | findstr :3000
taskkill /PID <PID> /F
```

#### ❌ Erro "Cannot find module" ou "Module não encontrado"

**Causa:** Dependências do projeto não foram instaladas.

**Solução:**

```bash
# Entre na pasta do projeto
cd /caminho/do/projeto

# Para projetos Node.js
npm install

# Para projetos Python
pip install -r requirements.txt

# Depois, no Relief, pare e inicie o projeto novamente
```

#### ❌ Comando `dev` não existe ou está incorreto

**Causa:** O comando definido em `scripts.dev` não existe no projeto.

**Exemplo de erro:**
```
npm ERR! Missing script: "dev"
```

**Solução:**

1. Verifique qual script seu projeto usa:

```bash
# Para Node.js, veja o package.json
cat package.json
```

Procure pela seção `"scripts"`:
```json
{
  "scripts": {
    "start": "node index.js",   // ← seu comando pode ser este
    "dev": "nodemon index.js"    // ← ou este
  }
}
```

2. Atualize o `relief.yaml` com o comando correto:

```yaml
scripts:
  dev: "npm run start"  # ou o comando que realmente existe
```

---

### "Failed to modify /etc/hosts" (macOS/Linux)

**📝 Sintoma:**  
Erro ao tentar modif icar o arquivo `/etc/hosts`.

**🔍 Causa:**  
O Relief precisa de permissões de administrador para modificar este arquivo.

**✅ Solução:**

**Opção 1: Digite sua senha**

Quando o Relief pedir, digite sua senha de administrador do sistema.

**Opção 2: Modifique manualmente**

```bash
# Abra o arquivo
sudo nano /etc/hosts

# Adicione esta linha no final (substitua pelo seu domínio):
127.0.0.1 meu-projeto.local.dev

# Salve e saia:
# Ctrl+O (salvar)
# Enter (confirmar)
# Ctrl+X (sair)
```

**Opção 3: Dê permissões permanentes (não recomendado)**

```bash
sudo chmod 666 /etc/hosts
```

⚠️ **Atenção:** Isso deixa o arquivo menos seguro.

---

## 🔴 Problemas de Rede e Domínios

### Domínio .local.dev não abre no navegador

**📝 Sintoma:**  
Você acessa `http://projeto.local.dev` mas o navegador diz "Site não encontrado" ou "ERR_NAME_NOT_RESOLVED".

**🔍 Diagnóstico Completo:**

#### Passo 1: Verifique se o projeto está rodando

No Relief, o status deve estar **🟢 "Rodando"**.

Se não estiver, veja: [Projeto não inicia](#projeto-não-inicia-fica-em-iniciando)

---

#### Passo 2: Teste acesso direto pela porta

```bash
# Descubra qual porta o projeto usa (veja no relief.yaml em env.PORT)
# Exemplo: se PORT=3000
```

Abra no navegador: `http://localhost:3000`

**Se funcionar:** O problema é com o proxy/hosts  
**Se não funcionar:** O problema é com o projeto em si

---

#### Passo 3: Verifique o arquivo hosts

**macOS/Linux:**
```bash
cat /etc/hosts | grep local.dev
```

**Windows:**
```powershell
type C:\Windows\System32\drivers\etc\hosts | findstr local.dev
```

**Deve aparecer algo assim:**
```
127.0.0.1 projeto.local.dev
```

**Se não aparecer:** O Relief não conseguiu modificar o arquivo. Veja: ["Failed to modify /etc/hosts"](#failed-to-modify-etchosts-macoslinux)

---

#### Passo 4: Verifique o Traefik

No Relief, deve haver um indicador de que o proxy (Traefik) está ativo.

Se não estiver, veja: ["Failed to start Traefik"](#failed-to-start-traefik--erro-com-proxy)

---

**✅ Soluções:**

### Solução 1: Reinicie tudo

```bash
# Pare o projeto no Relief
# Feche o Relief
# Abra o Relief novamente
# Inicie o projeto
```

### Solução 2: Limpe o cache de DNS

**macOS:**
```bash
sudo dscacheutil -flushcache
sudo killall -HUP mDNSResponder
```

**Windows:**
```powershell
ipconfig /flushdns
```

**Linux:**
```bash
# Ubuntu/Debian com systemd-resolved
sudo systemd-resolve --flush-caches

# Outras distros
sudo /etc/init.d/nscd restart
```

### Solução 3: Limpe o cache do navegador

1. Pressione `Ctrl+Shift+Delete` (Windows/Linux) ou `Cmd+Shift+Delete` (macOS)
2. Marque "Cache" e "Cookies"
3. Clique em "Limpar dados"

Ou force refresh na página:
- Windows/Linux: `Ctrl+Shift+R` ou `Ctrl+F5`
- macOS: `Cmd+Shift+R`

### Solução 4: Tente outro navegador

Às vezes o problema é específico de um navegador. Tente Chrome, Firefox ou Safari.

---

### Erro 502 Bad Gateway ou 503 Service Unavailable

**📝 Sintoma:**  
O domínio abre, mas aparece "502 Bad Gateway" ou "503 Service Unavailable".

**🔍 Causa:**  
O Traefik está funcionando, mas seu projeto não está respondendo na porta correta.

**✅ Solução:**

1. **Verifique os logs do projeto:**
   - Clique em "Ver Logs" no Relief
   - Veja se há erros

2. **Confirme a porta:**
   - Verifique se a variável `PORT` no `relief.yaml` está correta
   - Certifique-se de que seu código usa essa porta

3. **Para projetos Node.js/Express:**

Seu código deve escutar em `0.0.0.0`, não em `localhost`:

```javascript
// ❌ Errado
app.listen(3000, 'localhost', () => { ... })

// ✅ Correto
app.listen(3000, '0.0.0.0', () => { ... })

// ✅ Ou simplesmente
app.listen(3000, () => { ... })
```

4. **Reinicie o projeto:**
   - Pare e inicie novamente no Relief

---

## 🔴 Problemas de Performance

### Relief está lento ou travando

**📝 Sintoma:**  
A interface do Relief demora para responder ou trava.

**🔍 Causas Comuns:**

#### 1. Muitos projetos rodando ao mesmo tempo

**Solução:** Pare projetos que você não está usando no momento.

#### 2. Logs muito grandes

**Solução:** Reinicie o projeto para limpar os logs acumulados.

#### 3. Muita saída no console

Se seu projeto imprime muitas mensagens, pode deixar o Relief lento.

**Solução:** Reduza logs desnecessários:

```javascript
// Node.js - remova console.logs desnecessários
// ou use níveis de log
if (process.env.NODE_ENV !== 'production') {
  console.log('...')
}
```

#### 4. Recurso do sistema insuficientes

**Solução:** Feche outros programas pesados.

---

## 🔴 Problemas Específicos por Sistema

### Windows

#### "The system cannot find the path specified"

**Causa:** Caminhos com espaços ou caracteres especiais.

**Solução:**

- Evite espaços nos caminhos (use `C:\projetos` ao invés de `C:\Meus Projetos`)
- Se não puder evitar, use aspas duplas e escape barras:

```yaml
scripts:
  dev: "\"C:\\Program Files\\node\\node.exe\" index.js"
```

---

#### Antivírus/Windows Defender bloqueia o Relief

**Causa:** Falso positivo (comum em apps desenvolvidos localmente).

**Solução:**

1. Abra o **Windows Security**
2. Vá em **Proteção contra vírus e ameaças**
3. **"Gerenciar configurações"**
4. Adicione a pasta do Relief em **"Exclusões"**

---

### macOS

#### "Relief.app is damaged and can't be opened"

**Causa:** macOS Gatekeeper bloqueia apps não assinados.

**Solução Rápida:**

```bash
xattr -cr /caminho/para/Relief.app
```

**Solução pela Interface:**

1. Vá em **Preferências do Sistema**
2. **Segurança e Privacidade**
3. Aba **"Geral"**
4. Clique em **"Abrir Mesmo Assim"**

---

#### "Permission denied" ao modificar /etc/hosts

**Solução:**

```bash
sudo chmod 644 /etc/hosts
```

Se pedir senha, digite sua senha de administrador.

---

### Linux

#### "failed to initialize glfw: API unavailable"

**Causa:** Faltam bibliotecas gráficas.

**Solução:**

```bash
# Ubuntu/Debian
sudo apt-get install -y libgl1-mesa-dev xorg-dev

# Fedora
sudo dnf install mesa-libGL-devel libX11-devel

# Arch
sudo pacman -S mesa libx11
```

---

## 🔧 Ferramentas de Diagnóstico

### wails doctor

Execute para ver um relatório completo do ambiente:

```bash
wails doctor
```

Isso mostra:
- ✓ ou ✗ para cada dependência instalada
- Versões de todas as ferramentas
- Configurações do sistema

---

### Modo Debug

Execute o Relief em modo debug para ver logs detalhados:

```bash
wails dev -debug
```

Isso mostra:
- Logs do backend (Go)
- Logs do frontend (React)
- Comunicação entre eles
- Erros detalhados

---

### Verificar Portas em Uso

**macOS/Linux:**
```bash
# Ver todas as portas em uso
sudo lsof -i -P -n | grep LISTEN

# Ver porta específica
sudo lsof -i :3000
```

**Windows:**
```powershell
# Ver todas as portas em uso
netstat -ano | findstr LISTENING

# Ver porta específica
netstat -ano | findstr :3000
```

---

### Verificar Processos do Relief

**macOS/Linux:**
```bash
ps aux | grep relief
ps aux | grep wails
```

**Windows:**
```powershell
tasklist | findstr relief
tasklist | findstr wails
```

---

## 📞 Obtendo Suporte

Se nenhuma dessas soluções funcionou:

### 1. Colete Informações

Antes de pedir ajuda, colete estas informações:

```bash
# Sistema operacional e versão
# macOS
sw_vers

# Linux
cat /etc/os-release

# Windows
systeminfo | findstr /B /C:"OS Name" /C:"OS Version"

# Versões das ferramentas
go version
node --version
npm --version
wails version

# Diagnóstico do Wails
wails doctor
```

---

### 2. Execute em Modo Debug

```bash
wails dev -debug > relief-debug.log 2>&1
```

Isso salva todos os logs em `relief-debug.log`.

---

### 3. Abra uma Issue no GitHub

1. Acesse: https://github.com/Maycon-Santos/relief/issues/new
2. Use o template de bug report
3. Inclua:
   - ✅ Descrição do problema
   - ✅ O que você esperava que acontecesse
   - ✅ O que realmente aconteceu
   - ✅ Passos para reproduzir
   - ✅ Sistema operacional e versões (do passo 1)
   - ✅ Output do `wails doctor`
   - ✅ Logs completos (do passo 2)
   - ✅ Screenshots (se aplicável)
   - ✅ Arquivo `relief.yaml` do projeto (se for problema com um projeto específico)

---

### 4. Faça Perguntas na Comunidade

Para dúvidas gerais:
- **GitHub Discussions**: https://github.com/Maycon-Santos/relief/discussions
- Faça pergunâsa, compartilhe dicas, ajude outros usuários

---

### 5. Recursos de Suporte

- 📖 **Documentação Completa**: [README.md](README.md)
- 🚀 **Guia Rápido**: [QUICKSTART.md](QUICKSTART.md)
- 📦 **Guia de Instalação**: [INSTALLATION.md](INSTALLATION.md)
- 💬 **Discussions**: https://github.com/Maycon-Santos/relief/discussions
- 🐞 **Issues**: https://github.com/Maycon-Santos/relief/issues

---

## 🎯 Checklist de Troubleshooting

Use esta lista quando tiver problemas:

- [ ] Li a mensagem de erro completa
- [ ] Verifiquei se todas as dependências estão instaladas (`wails doctor`)
- [ ] Fechei e abri novamente o terminal após instalar algo
- [ ] Tentei executar em modo debug (`wails dev -debug`)
- [ ] Verifiquei os logs do projeto no Relief
- [ ] Procurei o erro neste guia de troubleshooting
- [ ] Procurei issues similares no GitHub
- [ ] Testei os comandos manualmente no terminal
- [ ] (Para projetos) Validei o `relief.yaml` em https://www.yamllint.com/
- [ ] Reiniciei o Relief
- [ ] Reiniciei o computador (quando tudo mais falhar!)

---

<p align="center">
  <b>Não encontrou solução?</b><br>
  Não hesite em abrir uma <a href="https://github.com/Maycon-Santos/relief/issues/new">issue</a> ou fazer uma pergunta nas <a href="https://github.com/Maycon-Santos/relief/discussions">discussions</a>!
</p>

<p align="center">
  <sub>Estamos aqui para ajudar! 💙</sub>
</p>
