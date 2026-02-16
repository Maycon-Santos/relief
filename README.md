# ⚡ Relief

> **Gerenciador visual para rodar múltiplos projetos localmente ao mesmo tempo**

Relief é uma ferramenta que permite você iniciar, parar e monitorar vários projetos de desenvolvimento (Node.js, Python, Docker, etc.) de forma simples através de uma interface gráfica moderna.

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Wails](https://img.shields.io/badge/Wails-v2-DF5B00?style=flat)](https://wails.io)
[![React](https://img.shields.io/badge/React-18-61DAFB?style=flat&logo=react)](https://reactjs.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

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

## 🚀 Primeiros Passos

### Passo 1: Instalar os Pré-Requisitos

Antes de usar o Relief, você precisa ter instalado:

**1. Go (Linguagem de Programação)**
```bash
# macOS (usando Homebrew)
brew install go

# Verifique se instalou corretamente
go version
```

**2. Node.js (para o frontend do Relief)**
```bash
# macOS (usando Homebrew)
brew install node

# Verifique se instalou corretamente
node --version
npm --version
```

**3. Wails CLI (ferramenta para criar apps desktop com Go)**
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Verifique se instalou corretamente
wails version
```

### Passo 2: Baixar o Relief

```bash
# Clone o repositório
git clone https://github.com/seu-usuario/relief.git

# Entre na pasta
cd relief
```

### Passo 3: Instalar as Dependências

```bash
# Baixa as bibliotecas do Go
go mod download

# Entra na pasta do frontend e instala dependências do Node
cd frontend
npm install
cd ..
```

### Passo 4: Executar o Relief

```bash
# Inicia o Relief em modo desenvolvimento
wails dev
```

Uma janela vai abrir com o Relief funcionando! 🎉

---

## 📖 Tutorial: Seu Primeiro Projeto

Vamos adicionar um projeto de exemplo que vem com o Relief.

### Passo 1: Inicie o Relief

Se ainda não iniciou, execute:
```bash
wails dev
```

### Passo 2: Adicione o Projeto de Exemplo

1. Clique no botão **"Adicionar Projeto Local"**
2. Navegue até a pasta `examples/hello-world` dentro do Relief
3. Selecione a pasta
4. O projeto aparecerá na interface

### Passo 3: Inicie o Projeto

1. Clique no botão verde **"Iniciar"** no card do projeto
2. Aguarde alguns segundos enquanto o projeto inicia
3. O status mudará para "Rodando" 🟢

### Passo 4: Acesse no Navegador

Abra seu navegador e acesse:
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

**Parabéns! 🎉** Você rodou seu primeiro projeto com o Relief!

---

## 🔧 Como Adicionar Seus Próprios Projetos

Para que seus projetos funcionem com o Relief, você precisa criar um arquivo de configuração.

### Criando o arquivo `relief.yaml`

Na raiz do seu projeto, crie um arquivo chamado `relief.yaml`:

```yaml
# Nome do projeto (aparece na interface)
name: "minha-api"

# URL pela qual você vai acessar (sem http://)
domain: "api.local.dev"

# Tipo do projeto (node, python, docker, go, ruby, java)
type: "node"

# Dependências necessárias
dependencies:
  - name: "node"
    version: ">=18.0.0"

# Comandos para executar
scripts:
  # Comando para instalar dependências
  install: "npm install"
  
  # Comando para iniciar o projeto
  dev: "npm run dev"

# Variáveis de ambiente
env:
  PORT: "3000"
  NODE_ENV: "development"
```

### Exemplo para Projeto Python

```yaml
name: "api-python"
domain: "python-api.local.dev"
type: "python"

dependencies:
  - name: "python"
    version: ">=3.9"

scripts:
  install: "pip install -r requirements.txt"
  dev: "python app.py"

env:
  FLASK_APP: "app.py"
  FLASK_ENV: "development"
```

### Exemplo para Projeto Docker

```yaml
name: "meu-container"
domain: "container.local.dev"
type: "docker"

scripts:
  dev: "docker-compose up"
  stop: "docker-compose down"
```

---

## 🎨 Entendendo a Interface

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

1. Verifique se já não existe uma [issue aberta](https://github.com/seu-usuario/relief/issues)
2. Se não existe, [crie uma nova issue](https://github.com/seu-usuario/relief/issues/new)
3. Descreva o problema com o máximo de detalhes possível
4. Inclua screenshots se possível

### 💡 Tem uma Ideia?

1. Abra uma [issue](https://github.com/seu-usuario/relief/issues/new) descrevendo sua ideia
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

E especialmente a todos os [contribuidores](https://github.com/seu-usuario/relief/graphs/contributors) que ajudaram a melhorar o projeto! ❤️

---

## 📞 Suporte

- **Issues**: [GitHub Issues](https://github.com/seu-usuario/relief/issues)
- **Discussões**: [GitHub Discussions](https://github.com/seu-usuario/relief/discussions)

---

<p align="center">
  <b>Feito com ❤️ pela comunidade Relief</b>
</p>

<p align="center">
  <sub>⭐ Se você gostou, deixe uma estrela no repositório!</sub>
</p>
