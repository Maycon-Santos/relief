# 🤝 Guia de Contribuição - Relief

Obrigado por considerar contribuir com o **Relief**! Este guia vai te ajudar a fazer sua primeira contribuição, mesmo se você nunca contribuiu com projetos open source antes.

---

## 🎯 Como Posso Ajudar?

Existem várias formas de contribuir, e você não precisa ser um expert:

### Para Iniciantes

- 📝 **Melhorar a documentação**: Encontrou algo confuso? Ajude a clarificar!
- 🐛 **Reportar bugs**: Encontrou algo que não funciona? Nos avise!
- 💡 **Sugerir melhorias**: Tem ideias? Compartilhe!
- ✅ **Testar o projeto**: Use e nos diga o que achou
- 🌍 **Traduzir**: Ajude a traduzir a documentação

### Para Desenvolvedores

- 🔧 **Corrigir bugs**: Escolha uma issue e resolva
- ✨ **Implementar features**: Adicione novas funcionalidades
- 🧪 **Escrever testes**: Aumente a cobertura de testes
- 🎨 **Melhorar UI**: Faça a interface mais bonita
- 📚 **Criar exemplos**: Adicione novos projetos de exemplo

---

## 📝 Reportando Bugs

### Passo 1: Verifique se já foi reportado

Antes de criar uma nova issue, [procure nas issues existentes](https://github.com/Maycon-Santos/relief/issues) para ver se alguém já reportou o mesmo problema.

### Passo 2: Crie uma Issue

Se é um bug novo, [crie uma issue](https://github.com/Maycon-Santos/relief/issues/new) incluindo:

```markdown
**Descrição do Bug**
O que aconteceu? Seja claro e objetivo.

**Como Reproduzir**
1. Faça isso...
2. Depois faça aquilo...
3. Veja o erro

**Comportamento Esperado**
O que você esperava que acontecesse?

**Comportamento Atual**
O que aconteceu de fato?

**Ambiente**
- OS: macOS 14.0 / Windows 11 / Ubuntu 22.04
- Versão do Go: 1.22.0
- Versão do Node: 18.19.0
- Versão do Relief: 1.0.0

**Screenshots**
Se possível, adicione prints ou vídeos

**Logs**
Cole aqui os logs relevantes
```

---

## 💡 Sugerindo Novas Funcionalidades

### Passo 1: Descreva sua ideia

[Abra uma issue](https://github.com/Maycon-Santos/relief/issues/new) com:

```markdown
**Funcionalidade Sugerida**
Descrição clara do que você quer

**Problema que Resolve**
Por que isso seria útil?

**Como Deveria Funcionar**
Descreva o comportamento esperado

**Alternativas Consideradas**
Você pensou em outras formas de resolver isso?

**Informações Adicionais**
Mockups, exemplos, screenshots, etc.
```

### Passo 2: Aguarde Feedback

A comunidade vai discutir e se aprovado, você ou outra pessoa pode implementar!

---

## 🔧 Configurando Ambiente de Desenvolvimento

### Requisitos

Antes de começar, instale:

- **Go 1.22+**: [Download](https://go.dev/dl/)
- **Node.js 18+**: [Download](https://nodejs.org/)
- **Git**: [Download](https://git-scm.com/)
- **Wails CLI**: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- **golangci-lint**: [Instruções](https://golangci-lint.run/usage/install/)

### Passo a Passo

```bash
# 1. Faça um fork do repositório no GitHub
# Clique em "Fork" no canto superior direito

# 2. Clone SEU fork (não o repositório original)
git clone https://github.com/SEU-USUARIO/relief.git
cd relief

# 3. Adicione o repositório original como "upstream"
git remote add upstream https://github.com/usuario-original/relief.git

# 4. Instale as dependências
go mod download
cd frontend && npm install && cd ..

# 5. Execute em modo desenvolvimento
wails dev
```

Se tudo funcionou, uma janela do Relief vai abrir! 🎉

---

## ✍️ Fazendo Sua Primeira Contribuição

### Passo 1: Crie uma Branch

```bash
# Certifique-se de estar na branch main e atualizado
git checkout main
git pull upstream main

# Crie uma branch para sua mudança
# Use um nome descritivo!
git checkout -b corrige-bug-porta
```

### Passo 2: Faça Suas Mudanças

Edite os arquivos necessários. Algumas dicas:

**Para código Go:**
```bash
# Execute o formatador
gofmt -w .

# Execute o linter
golangci-lint run

# Execute os testes
go test ./...
```

**Para código Frontend:**
```bash
cd frontend

# Execute o linter
npm run lint:fix

# Execute os testes (se houver)
npm test
```

### Passo 3: Teste Suas Mudanças

```bash
# Execute o Relief e teste manualmente
wails dev

# Execute os testes automatizados
./build/ci/test.sh
```

### Passo 4: Commit Suas Mudanças

Use mensagens de commit claras:

```bash
# Adicione os arquivos modificados
git add .

# Faça o commit com uma mensagem descritiva
git commit -m "fix: corrige erro de porta já em uso"
```

**Formato de mensagens de commit:**
```
tipo: descrição curta

Tipo pode ser:
- feat: nova funcionalidade
- fix: correção de bug
- docs: apenas documentação
- style: formatação, ponto e vírgula, etc
- refactor: melhoria de código sem mudar funcionalidade
- test: adicionar ou corrigir testes
- chore: tarefas de manutenção
```

**Exemplos:**
```bash
git commit -m "feat: adiciona suporte para Python 3.11"
git commit -m "fix: previne crash quando relief.yaml está ausente"
git commit -m "docs: atualiza guia de instalação"
git commit -m "style: formata código com gofmt"
```

### Passo 5: Envie para o GitHub

```bash
# Envie sua branch para SEU fork no GitHub
git push origin corrige-bug-porta
```

### Passo 6: Crie um Pull Request

1. Vá para o GitHub no SEU fork
2. Clique em **"Compare & pull request"**
3. Preencha o template:

```markdown
## Descrição
Breve descrição do que mudou

## Motivação
Por que essa mudança é necessária?

## Mudanças
- Mudança 1
- Mudança 2

## Como Testar
1. Faça isso
2. Faça aquilo
3. Verifique que...

## Screenshots (se aplicável)
Cole aqui prints de antes e depois

## Checklist
- [ ] Testei localmente
- [ ] Adicionei/atualizei testes
- [ ] Atualizei a documentação
- [ ] O linter passou
- [ ] Não quebra nada existente
```

4. Clique em **"Create pull request"**

### Passo 7: Aguarde Revisão

- Um mantenedor vai revisar seu código
- Eles podem pedir mudanças
- Faça as mudanças solicitadas na mesma branch
- Assim que aprovado, será feito o merge! 🎉

---

## 📚 Padrões de Código

### Código Go

**Formatação:**
```bash
# Use gofmt antes de commitar
gofmt -w .
```

**Boas práticas:**
- Use nomes descritivos de variáveis
- Documente funções públicas
- Trate erros adequadamente
- Escreva testes

**Exemplo:**
```go
func GetProject(id string) (*domain.Project, error) {
    if id == "" {
        return nil, errors.New("id cannot be empty")
    }
    
    return nil, nil
}
```

### Código TypeScript/React

**Formatação:**
```bash
cd frontend
npm run format
```

**Boas práticas:**
- Use TypeScript (não JavaScript)
- Componentes funcionais com hooks
- Defina tipos para props
- Evite `any`

**Exemplo:**
```typescript
interface ProjectCardProps {
  project: Project;
  onStart: () => Promise<void>;
}

export function ProjectCard({ project, onStart }: ProjectCardProps) {
  return (
    <Card>
      <h2>{project.name}</h2>
    </Card>
  );
}
```

---

## 🧪 Testes

### Executando Testes

```bash
# Todos os testes
./build/ci/test.sh

# Apenas um pacote
go test ./internal/config/

# Com output verbose
go test -v ./...

# Com cobertura
go test -cover ./...
```

### Escrevendo Testes

```go
package config

import "testing"

func TestLoadConfig(t *testing.T) {
    loader := NewLoader()
    
    config, err := loader.LoadConfig()
    if err != nil {
        t.Fatalf("erro inesperado: %v", err)
    }
    
    if config == nil {
        t.Fatal("esperava config, recebeu nil")
    }
}
```

---

## 📂 Estrutura do Projeto

Entenda onde cada coisa fica:

```
relief/
├── main.go                 # Ponto de entrada da aplicação
├── internal/               # Código Go interno
│   ├── app/                # Conecta Go com a interface
│   ├── config/             # Carrega configurações
│   ├── domain/             # Modelos (Project, Manifest, etc)
│   ├── runner/             # Inicia projetos (Native, Docker)
│   ├── dependency/         # Verifica Node, Python, etc
│   ├── proxy/              # Gerencia Traefik e /etc/hosts
│   ├── git/                # Operações Git
│   └── storage/            # Banco de dados SQLite
├── pkg/                    # Código reutilizável
│   ├── logger/             # Sistema de logs
│   ├── fileutil/           # Utilitários de arquivos
│   └── httputil/           # Utilitários HTTP
├── frontend/               # Interface (React + TypeScript)
│   └── src/
│       ├── components/     # Componentes visuais
│       ├── hooks/          # Hooks React personalizados
│       ├── services/       # Comunicação com backend
│       └── types/          # Tipos TypeScript
├── examples/               # Projetos de exemplo
├── configs/                # Exemplos de configuração
├── docs/                   # Documentação
└── build/                  # Scripts de compilação
```

---

## 🎯 Áreas que Precisam de Ajuda

### 🔥 Alta Prioridade

- **DockerRunner**: Implementação completa do suporte Docker
- **Auto-instaladores**: Instalação automática de Node.js/Python
- **Testes**: Aumentar cobertura (atualmente ~40%)
- **Documentação**: Mais exemplos e tutoriais

### 🎨 Interface

- Melhorar design visual
- Adicionar tema escuro
- Mensagens de erro mais claras
- Estados de loading melhores

### 📚 Documentação

- Vídeos tutoriais
- Artigos e blog posts
- Traduções (inglês, espanhol, etc)
- Diagramas de arquitetura

### ✨ Features

- Health checks para projetos
- Métricas e monitoramento
- Notificações desktop
- Sistema de plugins
- Suporte para projetos remotos (SSH)

---

## 👀 Processo de Code Review

Todos os Pull Requests passam por revisão:

**O que os revisores verificam:**
- ✅ O código funciona?
- ✅ Está bem testado?
- ✅ A documentação foi atualizada?
- ✅ Segue os padrões do projeto?
- ✅ Não quebra nada existente?

**Dicas para aprovação rápida:**
- Mantenha PRs pequenos e focados
- Escreva testes
- Documente suas mudanças
- Seja receptivo ao feedback

---

## ❓ Dúvidas Frequentes

### Nunca contribuí com open source antes. É difícil?

Não! Todo mundo começou algum dia. Este guia foi feito justamente para ajudar iniciantes. Se tiver dúvidas, pergunte nas issues!

### Não sei por onde começar

1. Leia a documentação
2. Execute o projeto localmente
3. Procure issues com a tag `good first issue`
4. Pergunte se precisar de ajuda

### Meu PR foi rejeitado. E agora?

Isso é normal! Leia o feedback, faça os ajustes sugeridos e atualize o PR. É um processo de aprendizado.

### Posso trabalhar em qualquer issue?

Issues abertas estão disponíveis. Se alguém já comentou que está trabalhando, escolha outra ou pergunte se pode ajudar.

### Quanto tempo leva para meu PR ser revisado?

Geralmente alguns dias. Projetos open source dependem de voluntários, então pode demorar um pouco.

---

## 🤔 Precisa de Ajuda?

- **Issues**: [GitHub Issues](https://github.com/Maycon-Santos/relief/issues)
- **Discussões**: [GitHub Discussions](https://github.com/Maycon-Santos/relief/discussions)
- **Email**: Veja o README principal

---

## 📜 Código de Conduta

Seja gentil e respeitoso:

- ✅ Seja inclusivo e acolhedor
- ✅ Dê feedback construtivo
- ✅ Foque no código, não na pessoa
- ✅ Ajude outros a aprender
- ❌ Não seja rude ou desrespeitoso
- ❌ Não faça ataques pessoais

---

## 📝 Licença

Ao contribuir, você concorda que suas contribuições serão licenciadas sob a **Licença MIT**.

---

<p align="center">
  <b>Obrigado por tornar o Relief melhor! ✨</b>
</p>

<p align="center">
  Toda contribuição, por menor que seja, faz diferença! 🚀
</p>
