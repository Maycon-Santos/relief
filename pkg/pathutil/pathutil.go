package pathutil

import (
	"os"
	"path/filepath"
	"strings"
)

// ToRelativeHome converte um caminho absoluto para um caminho relativo ao home (~/)
func ToRelativeHome(path string) string {
	if path == "" {
		return path
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return path
	}

	if strings.HasPrefix(absPath, homeDir) {
		relativePath := strings.TrimPrefix(absPath, homeDir)
		if relativePath == "" {
			return "~"
		}
		if !strings.HasPrefix(relativePath, "/") {
			relativePath = "/" + relativePath
		}
		return "~" + relativePath
	}

	return path
}

// FromRelativeHome converte um caminho relativo ao home (~/) para um caminho absoluto
func FromRelativeHome(path string) string {
	if path == "" {
		return path
	}

	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(homeDir, path[2:])
	}

	if path == "~" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return homeDir
	}

	return path
}

// ConvertYAMLPathsToRelative converte caminhos absolutos em um YAML para relativos ao home
func ConvertYAMLPathsToRelative(yamlContent string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return yamlContent
	}

	return strings.ReplaceAll(yamlContent, homeDir, "~")
}

// ConvertYAMLPathsToAbsolute converte caminhos relativos ao home em um YAML para absolutos
func ConvertYAMLPathsToAbsolute(yamlContent string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return yamlContent
	}

	patterns := []struct {
		old string
		new string
	}{
		{`"~/`, `"` + homeDir + `/`},
		{`'~/`, `'` + homeDir + `/`},
		{`: ~/`, `: ` + homeDir + `/`},
	}

	result := yamlContent
	for _, pattern := range patterns {
		result = strings.ReplaceAll(result, pattern.old, pattern.new)
	}

	return result
}
