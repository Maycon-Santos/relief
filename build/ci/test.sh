#!/bin/bash
# Script de teste

set -e

echo "🧪 Running tests..."

# Cores
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

# Executar linter
echo -e "${BLUE}Running linter...${NC}"
if command -v golangci-lint &> /dev/null; then
    golangci-lint run
    echo -e "${GREEN}✓ Linter passed${NC}"
else
    echo -e "${RED}⚠ golangci-lint not installed, skipping${NC}"
fi

# Executar testes
echo -e "${BLUE}Running unit tests...${NC}"
go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Tests passed${NC}"
else
    echo -e "${RED}❌ Tests failed${NC}"
    exit 1
fi

# Gerar relatório de coverage
echo -e "${BLUE}Generating coverage report...${NC}"
go tool cover -html=coverage.txt -o coverage.html
echo -e "${GREEN}✓ Coverage report generated: coverage.html${NC}"

# Exibir coverage summary
go tool cover -func=coverage.txt | tail -1

echo -e "${GREEN}✅ All tests completed!${NC}"
