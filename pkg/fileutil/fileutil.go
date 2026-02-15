// Package fileutil provides utilities for filesystem operations.
package fileutil

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Exists checks if a file or directory exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsDir checks if the path is a directory
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// EnsureDir creates a directory if it doesn't exist (including parents)
func EnsureDir(path string) error {
	if Exists(path) {
		return nil
	}
	return os.MkdirAll(path, 0755)
}

// CopyFile copies a file from src to dst
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error opening source file: %w", err)
	}
	defer sourceFile.Close()

	// Create destination directory if it doesn't exist
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

	// Copy permissions
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("error getting source file info: %w", err)
	}

	return os.Chmod(dst, sourceInfo.Mode())
}

// WriteFile writes content to a file, creating directories if necessary
func WriteFile(path string, content []byte, perm os.FileMode) error {
	if err := EnsureDir(filepath.Dir(path)); err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}
	return os.WriteFile(path, content, perm)
}

// ReadFile reads the content of a file
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// GetHomeDir returns the user's home directory
func GetHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting home directory: %w", err)
	}
	return home, nil
}

// GetSofredorDir returns the ~/.sofredor directory
func GetSofredorDir() (string, error) {
	home, err := GetHomeDir()
	if err != nil {
		return "", err
	}
	sofredorDir := filepath.Join(home, ".sofredor")
	if err := EnsureDir(sofredorDir); err != nil {
		return "", fmt.Errorf("error creating .sofredor directory: %w", err)
	}
	return sofredorDir, nil
}

// GetReliefSubDir returns a subdirectory inside ~/.relief
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
