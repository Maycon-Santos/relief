# Contributing to Relief Orchestrator

Thank you for considering contributing to **Relief Orchestrator**! 🎉

## How to Contribute

### 🐛 Reporting Bugs

Found a bug? Please [open an issue](https://github.com/omelete/relief/issues/new) with:
- Clear description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Environment (OS, Go version, Node version)
- Relevant logs or screenshots

### 💡 Suggesting Features

Have an idea? [Open an issue](https://github.com/omelete/relief/issues/new) with:
- Clear description of the feature
- Use case / problem it solves
- Examples or mockups (if applicable)

### 🔧 Setting Up Development Environment

```bash
# Clone the repository
git clone https://github.com/omelete/relief.git
cd relief

# Install dependencies
go mod download
cd frontend && npm install && cd ..

# Run in dev mode
wails dev
```

**Requirements:**
- Go 1.22+
- Node.js 18+
- Wails CLI v2

## Code Conventions

### Go Code

- **Follow Standard Go Layout**
- **Use `gofmt`** for formatting
- **Run linter:** `golangci-lint run`
- **Write tests** for new features
- **Document public APIs** with comments
- **Use descriptive variable names**

**Example:**
```go
// GetProject retrieves a project by ID
func (r *ProjectRepository) GetProject(id string) (*domain.Project, error) {
    // Implementation...
}
```

### TypeScript/React Code

- **Use TypeScript** for type safety
- **Functional components** with hooks
- **Props interfaces** for components
- **No `any` types** (use specific types)
- **Format with Prettier**

**Example:**
```typescript
interface ProjectCardProps {
  project: Project;
  onStart: () => Promise<void>;
  onStop: () => Promise<void>;
}

export const ProjectCard: React.FC<ProjectCardProps> = ({ project, onStart, onStop }) => {
  // Implementation...
};
```

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `style`: Code style (formatting, no logic change)
- `refactor`: Code refactor (no feat/fix)
- `test`: Adding tests
- `chore`: Maintenance tasks

**Examples:**
```
feat(runner): add docker support
fix(proxy): handle missing hosts file gracefully
docs(readme): update installation instructions
```

## Pull Request Process

1. **Fork the repository** and create a branch from `main`
2. **Make your changes** following code conventions
3. **Write/update tests** to cover changes
4. **Update documentation** if needed
5. **Run tests:** `./build/ci/test.sh`
6. **Run linter:** `golangci-lint run`
7. **Create PR** with clear description

### PR Title

Follow commit message convention:
```
feat: add support for kubernetes runner
fix: prevent crash on missing manifest
```

### PR Description Template

```markdown
## Description
Brief description of changes

## Motivation
Why is this change needed?

## Changes
- List of changes made
- Another change

## Testing
How was this tested?

## Screenshots (if applicable)
Add screenshots for UI changes

## Checklist
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] Linter passing
- [ ] No breaking changes (or documented)
```

## Project Structure

```
/relief
├── cmd/app/              # Main entrypoint
├── internal/             # Private code
│   ├── app/              # Wails bindings
│   ├── config/           # Configuration management
│   ├── domain/           # Business entities
│   ├── runner/           # Execution strategies
│   ├── dependency/       # Dependency checkers
│   ├── proxy/            # Network management
│   └── storage/          # Database layer
├── pkg/                  # Reusable utilities
├── frontend/             # React app
├── examples/             # Example projects
├── configs/              # Config examples
├── docs/                 # Documentation
└── build/                # Build scripts
```

## Testing

### Running Tests

```bash
# All tests
./build/ci/test.sh

# Specific package
go test ./internal/config/...

# With coverage
go test -cover ./...

# Verbose
go test -v ./...
```

### Writing Tests

```go
func TestConfigLoader_LoadConfig(t *testing.T) {
    // Arrange
    loader := config.NewLoader()
    
    // Act
    cfg, err := loader.LoadConfig()
    
    // Assert
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if cfg == nil {
        t.Fatal("expected config, got nil")
    }
}
```

## Documentation

### Code Documentation

- Document all exported functions/types
- Use GoDoc format
- Include examples for complex APIs

### User Documentation

- Update README.md for major features
- Update docs/ for architectural changes
- Keep manifest-schema.md current

## Areas Needing Contributions

### 🚀 High Priority

- **DockerRunner:** Complete implementation
- **Auto-installers:** Node.js/Python portable installers
- **Tests:** Increase coverage (currently ~40%)
- **Documentation:** More examples and tutorials

### 🎨 UI/UX

- Improve interface design
- Add dark/light theme toggle
- Better error messages
- Loading states

### 📚 Documentation

- Video tutorials
- Blog posts / articles
- Translations (Spanish, Portuguese, etc.)
- Architecture diagrams

### 🔧 Features

- Health checks for projects
- Metrics / monitoring
- Desktop notifications
- Plugin system
- Remote project support (SSH)

## Code Review

All PRs require:
- ✅ At least one approval
- ✅ Passing CI/CD checks
- ✅ No merge conflicts
- ✅ Updated documentation

## Questions?

- Open an [Issue](https://github.com/omelete/relief/issues) with `question` tag
- Join discussions on GitHub

## Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Focus on the code, not the person
- Help others learn and grow

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

**Thank you for making Relief Orchestrator better!** ✨
