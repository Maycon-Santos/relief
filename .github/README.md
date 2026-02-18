# 🤖 GitHub Actions Workflows

Este diretório contém os workflows de CI/CD do Relief.

## 📋 Workflows Disponíveis

### 🚀 Release (`release.yml`)

**Trigger:**
- Push de tags `v*` (ex: `v1.0.0`)
- Manualmente via Actions interface

**Funcionalidade:**
- Faz build de binários para todas as plataformas (macOS, Linux, Windows)
- Cria release no GitHub com binários e checksums
- Gera release notes automaticamente

**Como usar:**
```bash
# Via script
./scripts/release.sh v1.0.0

# Ou manualmente
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

**Documentação completa:** [docs/RELEASE.md](../docs/RELEASE.md)

---

## 🛠️ Estrutura dos Workflows

```
.github/
└── workflows/
    └── release.yml          # Build e release de binários
```

---

## 📝 Adicionando Novos Workflows

Para adicionar um novo workflow:

1. Crie um arquivo `.yml` nesta pasta
2. Siga a sintaxe do GitHub Actions
3. Documente aqui neste README

**Exemplo básico:**

```yaml
name: Nome do Workflow

on:
  push:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Seu passo
        run: echo "Hello"
```

---

## 🔍 Monitorando Workflows

- **Visualizar execuções:** https://github.com/Maycon-Santos/relief/actions
- **Logs em tempo real:** Clique em qualquer workflow em execução
- **Re-executar:** Clique em "Re-run jobs" se um workflow falhar

---

## 🐛 Troubleshooting

### Workflow não dispara

**Causa comum:** Permissões ou sintaxe incorreta.

**Solução:**
1. Verifique a sintaxe YAML em https://www.yamllint.com/
2. Garanta que o trigger está correto
3. Verifique permissões do repositório

### Build falha

1. Clique no workflow que falhou
2. Expanda o passo que deu erro
3. Veja os logs completos
4. Corrija o problema e faça novo commit/tag

---

## 📚 Recursos

- [GitHub Actions Docs](https://docs.github.com/actions)
- [Workflow Syntax](https://docs.github.com/actions/reference/workflow-syntax-for-github-actions)
- [Guia de Release](../docs/RELEASE.md)
