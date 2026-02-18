#!/bin/bash

# Script para criar uma nova release do Relief
# Uso: ./scripts/release.sh v1.0.0

set -e

VERSION=$1

if [ -z "$VERSION" ]; then
    echo "❌ Erro: Versão não especificada"
    echo "Uso: ./scripts/release.sh v1.0.0"
    exit 1
fi

# Validar formato da versão
if ! [[ $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "❌ Erro: Versão deve estar no formato v1.0.0"
    exit 1
fi

echo "🚀 Preparando release $VERSION"

# Verificar se o repositório está limpo
if [[ -n $(git status -s) ]]; then
    echo "❌ Erro: Existem alterações não commitadas"
    echo "Commite ou descarte as alterações antes de criar uma release"
    git status -s
    exit 1
fi

# Verificar se estamos na branch main
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" != "main" ]; then
    echo "⚠️  Aviso: Você não está na branch main (atual: $CURRENT_BRANCH)"
    read -p "Deseja continuar mesmo assim? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Verificar se a tag já existe
if git rev-parse "$VERSION" >/dev/null 2>&1; then
    echo "❌ Erro: Tag $VERSION já existe"
    exit 1
fi

# Atualizar dependências
echo "📦 Atualizando dependências..."
go mod tidy
cd frontend && npm install && cd ..

# Executar testes
echo "🧪 Executando testes..."
go test ./... || {
    echo "❌ Testes falharam"
    exit 1
}

# Compilar localmente para verificar
echo "🔨 Verificando compilação..."
wails build -clean || {
    echo "❌ Build falhou"
    exit 1
}

# Criar tag
echo "🏷️  Criando tag $VERSION..."
git tag -a "$VERSION" -m "Release $VERSION"

# Push da tag
echo "⬆️  Enviando tag para o GitHub..."
git push origin "$VERSION"

echo ""
echo "✅ Release $VERSION criada com sucesso!"
echo ""
echo "📍 Próximos passos:"
echo "  1. Acesse: https://github.com/$(git config --get remote.origin.url | sed 's/.*github.com[:/]\(.*\)\.git/\1/')/actions"
echo "  2. Aguarde a conclusão do workflow de release"
echo "  3. Verifique a release em: https://github.com/$(git config --get remote.origin.url | sed 's/.*github.com[:/]\(.*\)\.git/\1/')/releases"
echo ""
echo "Para reverter (caso necessário):"
echo "  git tag -d $VERSION"
echo "  git push origin :refs/tags/$VERSION"
