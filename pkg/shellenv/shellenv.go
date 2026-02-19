package shellenv

import (
	"context"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var knownPaths = []string{
	"/opt/homebrew/bin",
	"/opt/homebrew/sbin",
	"/usr/local/bin",
	"/usr/local/sbin",
	"/usr/local/opt/postgresql@16/bin",
	"/usr/local/opt/redis/bin",
	"/usr/bin",
	"/usr/sbin",
	"/bin",
	"/sbin",
	"/opt/local/bin",
	"/home/linuxbrew/.linuxbrew/bin",
}

func EnrichedEnv() []string {
	existing := os.Environ()

	currentPath := os.Getenv("PATH")
	extra := make([]string, 0, len(knownPaths))
	for _, p := range knownPaths {
		if !strings.Contains(currentPath, p) {
			extra = append(extra, p)
		}
	}

	newPath := currentPath
	if len(extra) > 0 {
		newPath = strings.Join(extra, ":") + ":" + currentPath
	}

	env := make([]string, 0, len(existing))
	for _, e := range existing {
		if strings.HasPrefix(e, "PATH=") {
			continue
		}
		env = append(env, e)
	}
	env = append(env, "PATH="+newPath)
	return env
}

func Command(commandLine string) *exec.Cmd {
	cmd := exec.Command("sh", "-c", commandLine)
	cmd.Env = EnrichedEnv()
	return cmd
}

func CommandContext(ctx context.Context, commandLine string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, "sh", "-c", commandLine)
	cmd.Env = EnrichedEnv()
	return cmd
}

func LookPath(name string) (string, error) {
	path, err := exec.LookPath(name)
	if err == nil {
		return path, nil
	}

	for _, dir := range knownPaths {
		candidate := dir + "/" + name
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			if runtime.GOOS != "windows" {
				if info.Mode()&0111 != 0 {
					return candidate, nil
				}
			} else {
				return candidate, nil
			}
		}
	}

	return "", &exec.Error{Name: name, Err: exec.ErrNotFound}
}
