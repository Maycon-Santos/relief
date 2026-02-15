// Package fileutil fornece utilitários para operações de filesystem.
package fileutil

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Exists verifica se um arquivo ou diretório existe
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsDir verifica se o path é um diretório
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// EnsureDir cria um diretório se não existir (incluindo pais)
func EnsureDir(path string) error {
	if Exists(path) {
		return nil
	}
	return os.MkdirAll(path, 0755)
}

// CopyFile copia um arquivo de src para dst
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo fonte: %w", err)
	}
	defer sourceFile.Close()

	// Criar diretório de destino se não existir
	if err := EnsureDir(filepath.Dir(dst)); err != nil {
		return fmt.Errorf("erro ao criar diretório destino: %w", err)
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo destino: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("erro ao copiar arquivo: %w", err)
	}

	// Copiar permissões
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("erro ao obter info do arquivo fonte: %w", err)
	}

	return os.Chmod(dst, sourceInfo.Mode())
}

// WriteFile escreve conteúdo em um arquivo, criando diretórios se necessário
func WriteFile(path string, content []byte, perm os.FileMode) error {
	if err := EnsureDir(filepath.Dir(path)); err != nil {
		return fmt.Errorf("erro ao criar diretório: %w", err)
	}
	return os.WriteFile(path, content, perm)
}

// ReadFile lê o conteúdo de um arquivo
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// GetHomeDir retorna o diretório home do usuário
func GetHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("erro ao obter diretório home: %w", err)
	}
	return home, nil
}

// GetSofredorDir retorna o diretório ~/.sofredor
func GetSofredorDir() (string, error) {
	home, err := GetHomeDir()
	if err != nil {
		return "", err
	}
	sofredorDir := filepath.Join(home, ".sofredor")
	if err := EnsureDir(sofredorDir); err != nil {
		return "", fmt.Errorf("erro ao criar diretório .sofredor: %w", err)
	}
	return sofredorDir, nil
}

// GetSofredorSubDir retorna um subdiretório dentro de ~/.sofredor
func GetSofredorSubDir(subdir string) (string, error) {
	sofredorDir, err := GetSofredorDir()
	if err != nil {
		return "", err
	}
	subPath := filepath.Join(sofredorDir, subdir)
	if err := EnsureDir(subPath); err != nil {
		return "", fmt.Errorf("erro ao criar subdiretório: %w", err)
	}
	return subPath, nil
}
