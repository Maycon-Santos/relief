package fileutil

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func EnsureDir(path string) error {
	if Exists(path) {
		return nil
	}
	return os.MkdirAll(path, 0755)
}

func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error opening source file: %w", err)
	}
	defer sourceFile.Close()

	if err := EnsureDir(filepath.Dir(dst)); err != nil {
		return fmt.Errorf("error creating destination directory: %w", err)
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("error creating destination file: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("error copying file: %w", err)
	}

	sourceInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("error getting source file info: %w", err)
	}

	return os.Chmod(dst, sourceInfo.Mode())
}

func WriteFile(path string, content []byte, perm os.FileMode) error {
	if err := EnsureDir(filepath.Dir(path)); err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}
	return os.WriteFile(path, content, perm)
}

func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func GetHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting home directory: %w", err)
	}
	return home, nil
}

func GetReliefDir() (string, error) {
	home, err := GetHomeDir()
	if err != nil {
		return "", err
	}
	reliefDir := filepath.Join(home, ".relief")
	if err := EnsureDir(reliefDir); err != nil {
		return "", fmt.Errorf("error creating relief directory: %w", err)
	}
	return reliefDir, nil
}

func GetReliefSubDir(subdir string) (string, error) {
	reliefDir, err := GetReliefDir()
	if err != nil {
		return "", err
	}
	subPath := filepath.Join(reliefDir, subdir)
	if err := EnsureDir(subPath); err != nil {
		return "", fmt.Errorf("error creating subdirectory: %w", err)
	}
	return subPath, nil
}
