#!/bin/bash
# Script de build para múltiplas plataformas

set -e

echo "🔨 Building Sofredor Orchestrator..."

# Cores
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Diretório de saída
BUILD_DIR="build/bin"
mkdir -p "$BUILD_DIR"

# Versão
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
echo -e "${BLUE}Version: ${VERSION}${NC}"

# Plataformas
PLATFORMS=("darwin/amd64" "darwin/arm64" "linux/amd64" "windows/amd64")

for platform in "${PLATFORMS[@]}"; do
    IFS='/' read -r -a parts <<< "$platform"
    GOOS="${parts[0]}"
    GOARCH="${parts[1]}"
    
    OUTPUT_NAME="sofredor-orchestrator-${GOOS}-${GOARCH}"
    if [ "$GOOS" = "windows" ]; then
        OUTPUT_NAME="${OUTPUT_NAME}.exe"
    fi
    
    echo -e "${BLUE}Building for ${GOOS}/${GOARCH}...${NC}"
    
    # Build usando Wails
    wails build \
        -platform "${GOOS}/${GOARCH}" \
        -o "${BUILD_DIR}/${OUTPUT_NAME}" \
        -ldflags "-X main.version=${VERSION}"
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Built ${OUTPUT_NAME}${NC}"
    else
        echo -e "❌ Failed to build for ${GOOS}/${GOARCH}"
        exit 1
    fi
done

echo -e "${GREEN}✅ Build completed successfully!${NC}"
echo -e "Binaries available in: ${BUILD_DIR}"
