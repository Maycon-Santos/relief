# Guia de Contribuição

Obrigado por considerar contribuir com o **SofredorOrchestrator**! 🎉

## Código de Conduta

- Seja respeitoso e inclusivo
- Feedback construtivo é sempre bem-vindo
- Foque no problema técnico, não na pessoa

## Como Contribuir

### 1. Setup do Ambiente de Desenvolvimento

#### Pré-requisitos
- Go 1.22+
- Node.js 18+
- Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- golangci-lint: `brew install golangci-lint` (Mac) ou equivalente

#### Clone e Instale
```bash
git clone https://github.com/omelete/sofredor-orchestrator.git
cd sofredor-orchestrator

# Instalar dependências Go
go mod download

# Instalar dependências Frontend
cd frontend
npm install
cd ..
```

#### Executar em Modo Dev
```bash
wails dev
```

### 2. Estrutura de Branches

- `main`: Branch principal (sempre estável)
- `develop`: Branch de desenvolvimento
- `feature/nome-da-feature`: Novas funcionalidades
- `fix/nome-do-bug`: Correções de bugs

### 3. Processo de Contribuição

1. **Fork** o repositório
2. Crie uma **branch** a partir de `develop`:
   ```bash
   git checkout -b feature/minha-feature develop
   ```
3. **Implemente** sua mudança
4. **Teste** localmente:
   ```bash
   ./build/ci/test.sh
   ```
5. **Commit** com mensagem descritiva:
   ```bash
   git commit -m "feat: adiciona suporte a Podman runner"
   ```
6. **Push** para seu fork:
   ```bash
   git push origin feature/minha-feature
   ```
7. Abra um **Pull Request** para `develop`

### 4. Convenções de Código

#### Go
- Siga o [Effective Go](https://golang.org/doc/effective_go)
- Use `gofmt` para formatação
- Adicione comentários em funções públicas
- Mantenha funções pequenas (< 50 linhas)

**Exemplo:**
```go
// CheckDependencies verifica todas as dependências de um projeto.
// Retorna erro se alguma dependência crítica não estiver satisfeita.
func (m *Manager) CheckDependencies(ctx context.Context, project *domain.Project) error {
    // Implementação
}
```

#### TypeScript/React
- Use TypeScript strict mode
- Componentes funcionais com hooks
- Props tipadas com interfaces
- Nomeie arquivos com PascalCase para componentes

**Exemplo:**
```typescript
interface ProjectCardProps {
  project: Project;
  onStart: () => Promise<void>;
}

export function ProjectCard({ project, onStart }: ProjectCardProps) {
  // Implementação
}
```

### 5. Testes

#### Testes Unitários (Go)
```bash
go test ./...
```

Estrutura:
```go
func TestManager_CheckDependencies(t *testing.T) {
    // Arrange
    manager := NewManager(logger)
    project := &domain.Project{...}
    
    // Act
    err := manager.CheckDependencies(ctx, project)
    
    // Assert
    if err != nil {
        t.Errorf("expected no error, got %v", err)
    }
}
```

#### Testes de Integração
- Coloque em arquivos `*_integration_test.go`
- Use build tag: `// +build integration`

### 6. Mensagens de Commit

Siga [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` Nova funcionalidade
- `fix:` Correção de bug
- `docs:` Mudanças em documentação
- `style:` Formatação, ponto e vírgula, etc
- `refactor:` Refatoração de código
- `test:` Adição/modificação de testes
- `chore:` Tarefas de manutenção

**Exemplos:**
```
feat: adiciona DockerRunner completo
fix: corrige race condition no NativeRunner
docs: atualiza README com instruções de build
refactor: simplifica lógica de merge de configs
test: adiciona testes para HostsManager
```

### 7. Pull Request Guidelines

#### Checklist antes de submeter:
- [ ] Código compila sem erros
- [ ] Testes passam (`./build/ci/test.sh`)
- [ ] Linter passa (`golangci-lint run`)
- [ ] Documentação atualizada (se aplicável)
- [ ] CHANGELOG.md atualizado (para features/fixes significativos)
- [ ] Commit messages seguem convenções

#### Descrição do PR:
```markdown
## Descrição
Breve descrição da mudança

## Motivação
Por que essa mudança é necessária?

## Mudanças
- Mudança 1
- Mudança 2

## Testes
Como você testou isso?

## Screenshots (se aplicável)
```

### 8. Áreas que Precisam de Ajuda

Procurando por contribuições em:

- **DockerRunner:** Implementação completa usando Docker SDK
- **Instaladores Automáticos:** Download e instalação de Node.js, Python
- **UI/UX:** Melhorias na interface React
- **Testes:** Aumentar cobertura de testes
- **Documentação:** Tutoriais, exemplos, traduções
- **Novos Runners:** Podman, systemd, PM2
- **Novos Checkers:** Ruby, PHP, Java

### 9. Dúvidas?

- Abra uma [Issue](https://github.com/omelete/sofredor-orchestrator/issues) com a tag `question`
- Entre no Discord da comunidade (link no README)
- Envie email para: dev@omelete.com

## Reconhecimento

Todos os contribuidores serão listados no README e terão nosso agradecimento eterno! 🙏

---

**Obrigado por tornar o SofredorOrchestrator melhor!** ✨
