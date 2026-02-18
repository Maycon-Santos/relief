# 🚀 Guia de Release - Relief

Este documento explica como criar uma nova release do Relief e gerar binários automaticamente.

## 📋 Índice

- [Pré-requisitos](#pré-requisitos)
- [Processo de Release](#processo-de-release)
- [Como Funciona o CI/CD](#como-funciona-o-cicd)
- [Testando Localmente](#testando-localmente)
- [Troubleshooting](#troubleshooting)

---

## Pré-requisitos

Para criar uma release você precisa:

1. **Permissões de Write** no repositório
2. **Git configurado** localmente com acesso ao repositório
3. **Código testado e funcionando** na branch `main`
4. **Changelog atualizado** (opcional mas recomendado)

---

## Processo de Release

### Método 1: Automático via Script (Recomendado)

```bash
# Execute o script de release
./scripts/release.sh v1.0.0
```

O script irá:
- ✅ Validar o formato da versão
- ✅ Verificar se há alterações não commitadas
- ✅ Atualizar dependências
- ✅ Executar testes
- ✅ Fazer build local para validação
- ✅ Criar e enviar a tag para o GitHub
- ✅ Automaticamente disparar o workflow de release

---

### Método 2: Manual via Git

1. **Certifique-se de estar na branch main e com código atualizado:**

```bash
git checkout main
git pull origin main
```

2. **Verifique se tudo está funcionando:**

```bash
# Execute os testes
go test ./...

# Faça um build local
wails build -clean
```

3. **Crie a tag de versão:**

```bash
# Formato: vMAJOR.MINOR.PATCH (ex: v1.0.0)
git tag -a v1.0.0 -m "Release v1.0.0"
```

4. **Envie a tag para o GitHub:**

```bash
git push origin v1.0.0
```

5. **Aguarde o build:**

- Acesse: https://github.com/Maycon-Santos/relief/actions
- O workflow "Release" será executado automaticamente
- Aguarde ~15-20 minutos para concluir

6. **Verifique a release:**

- Acesse: https://github.com/Maycon-Santos/relief/releases
- A nova release deve aparecer com todos os binários

---

### Método 3: Manual via GitHub Interface

Se você preferir usar a interface do GitHub:

1. Acesse: https://github.com/Maycon-Santos/relief/releases/new
2. Clique em **"Choose a tag"**
3. Digite a nova versão (ex: `v1.0.0`) e clique em **"Create new tag: v1.0.0 on publish"**
4. Preencha:
   - **Release title**: `Relief v1.0.0`
   - **Description**: Descreva as mudanças
5. Clique em **"Publish release"**
6. O workflow será disparado automaticamente

---

## Como Funciona o CI/CD

O Relief usa GitHub Actions para automatizar o processo de build e release.

### Workflow: `.github/workflows/release.yml`

**Trigger:**
- Push de tags que começam com `v*` (ex: `v1.0.0`, `v2.1.3`)
- Manualmente via GitHub Actions interface

**Jobs:**

#### 1. Build (Paralelo)

Cria binários para:
- **macOS:** Intel (amd64) e Apple Silicon (arm64)
- **Linux:** x86_64 (amd64) e ARM64
- **Windows:** x86_64 (amd64)

**Processo para cada plataforma:**
1. Setup Go 1.22
2. Setup Node.js 18
3. Instala dependências do sistema (Linux: libgtk, webkit)
4. Instala Wails CLI
5. Instala dependências do frontend
6. Faz build com Wails
7. Compacta os binários (tar.gz para Unix, zip para Windows)
8. Faz upload dos artifacts

#### 2. Release

Após todos os builds concluírem:
1. Baixa todos os artifacts
2. Gera checksums SHA256
3. Cria release no GitHub com:
   - Todos os binários
   - Arquivo de checksums
   - Release notes automáticas
   - Instruções de instalação

**Tempo total:** ~15-20 minutos

---

## Testando Localmente

Antes de fazer uma release, teste localmente:

### Build Local Completo

```bash
# macOS
wails build -clean

# Linux
wails build -clean -platform linux/amd64

# Windows (em Linux/macOS com cross-compilation)
wails build -clean -platform windows/amd64
```

### Testando o Binário Gerado

```bash
# macOS
./build/bin/Relief.app/Contents/MacOS/relief --version

# Linux
./build/bin/relief --version

# Windows
./build/bin/relief.exe --version
```

### Testando Scripts de Instalação

Teste os comandos de instalação em uma máquina limpa ou Docker:

```bash
# Teste o download e instalação
curl -L https://github.com/Maycon-Santos/relief/releases/download/v1.0.0/relief-linux-amd64.tar.gz | tar xz
```

---

## Versionamento Semântico

O Relief segue [Semantic Versioning](https://semver.org/):

**Formato:** `vMAJOR.MINOR.PATCH`

- **MAJOR:** Mudanças incompatíveis na API
- **MINOR:** Nova funcionalidade compatível com versões anteriores
- **PATCH:** Correções de bugs compatíveis

**Exemplos:**
- `v1.0.0` - Primeira release estável
- `v1.1.0` - Nova feature adicionada
- `v1.1.1` - Bug fix
- `v2.0.0` - Breaking change

**Pre-releases:**
- `v1.0.0-beta.1` - Beta release
- `v1.0.0-rc.1` - Release candidate
- `v1.0.0-alpha.1` - Alpha release

---

## Troubleshooting

### Erro: "tag already exists"

**Causa:** A tag já foi criada anteriormente.

**Solução:**

Se foi um erro e você quer recriar:
```bash
# Delete localmente
git tag -d v1.0.0

# Delete no GitHub
git push origin :refs/tags/v1.0.0

# Recrie
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

---

### Build Falha no CI

**Diagnóstico:**

1. Acesse o workflow com erro
2. Clique no job que falhou
3. Veja os logs

**Causas comuns:**

- **Teste falhou:** Corrija o código e faça nova release
- **Dependência faltando:** Atualize o workflow para instalar
- **Timeout:** Pode acontecer, reexecute o workflow

---

### Binário Não Funciona

**Problema:** Binário baixado não executa ou dá erro.

**Checklist:**

1. **Permissão de execução (Linux/macOS):**
   ```bash
   chmod +x relief
   ```

2. **Biblioteca faltando (Linux):**
   ```bash
   ldd relief  # Veja quais libs estão faltando
   sudo apt-get install libgtk-3-0 libwebkit2gtk-4.0-37
   ```

3. **Quarentena do macOS:**
   ```bash
   xattr -cr Relief.app
   ```

4. **Windows Defender bloqueando:**
   - Adicione exceção para o executável

---

### Release Notes Não Aparecem

**Causa:** GitHub gera automaticamente baseado em commits desde a última tag.

**Solução:**

Edite a release manualmente:
1. Vá em: https://github.com/Maycon-Santos/relief/releases
2. Clique em "Edit" na release
3. Adicione/edite a descrição
4. Salve

---

## Checklist de Release

Use esta lista antes de fazer uma release:

- [ ] Código está na branch `main`
- [ ] Todos os testes passando (`go test ./...`)
- [ ] Build local funciona (`wails build`)
- [ ] Changelog atualizado (se houver)
- [ ] Versão segue semântica (vMAJOR.MINOR.PATCH)
- [ ] Commits relevantes têm mensagens claras
- [ ] Documentação atualizada (se necessário)
- [ ] Testado em pelo menos um sistema operacional
- [ ] Breaking changes documentados (se houver)

---

## Processo Post-Release

Após criar uma release:

1. **Teste a instalação via binário:**
   ```bash
   # Baixe e teste em cada plataforma
   ```

2. **Anuncie a release:**
   - Discussions do GitHub
   - README com link para última versão
   - Redes sociais (se aplicável)

3. **Monitore issues:**
   - Fique atento a reports de bugs na nova versão
   - Prepare hotfix se necessário (vMAJOR.MINOR.PATCH+1)

---

## Recursos Adicionais

- **GitHub Actions Docs:** https://docs.github.com/actions
- **Wails Build Docs:** https://wails.io/docs/guides/building
- **Semantic Versioning:** https://semver.org/
- **Conventional Commits:** https://www.conventionalcommits.org/

---

<p align="center">
  <b>Dúvidas sobre o processo de release?</b><br>
  Abra uma <a href="https://github.com/Maycon-Santos/relief/discussions">discussão</a>!
</p>
