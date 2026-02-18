package dependency

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/Maycon-Santos/relief/internal/config"
	"github.com/Maycon-Santos/relief/internal/domain"
	"github.com/Maycon-Santos/relief/pkg/logger"
)

type EnhancedManager struct {
	logger          *logger.Logger
	config          *config.Config
	runningServices map[string]bool
	healthCheckers  map[string]*time.Ticker
}

func NewEnhancedManager(log *logger.Logger, cfg *config.Config) *EnhancedManager {
	return &EnhancedManager{
		logger:          log,
		config:          cfg,
		runningServices: make(map[string]bool),
		healthCheckers:  make(map[string]*time.Ticker),
	}
}

func (m *EnhancedManager) StartManagedDependencies(ctx context.Context, project *domain.Project) error {
	for _, dep := range project.Dependencies {
		if !dep.Managed {
			continue
		}

		if m.runningServices[dep.Name] {
			m.logger.Info("Dependência já está executando", map[string]interface{}{
				"dependency": dep.Name,
			})
			continue
		}

		managedDep, exists := m.config.ManagedDependencies[dep.Name]
		if !exists {
			m.logger.Warn("Configuração não encontrada para dependência gerenciada", map[string]interface{}{
				"dependency": dep.Name,
			})
			continue
		}

		if err := m.checkAndInstallDependency(ctx, dep.Name, dep.Version, managedDep); err != nil {
			return fmt.Errorf("erro ao verificar/instalar dependência %s: %w", dep.Name, err)
		}

		if err := m.startService(ctx, dep.Name, managedDep); err != nil {
			return fmt.Errorf("erro ao iniciar serviço %s: %w", dep.Name, err)
		}

		m.startHealthCheck(ctx, dep.Name)

		m.runningServices[dep.Name] = true
	}

	return nil
}

func (m *EnhancedManager) StopManagedDependencies(ctx context.Context, project *domain.Project) error {
	for _, dep := range project.Dependencies {
		if !dep.Managed || !m.runningServices[dep.Name] {
			continue
		}

		managedDep, exists := m.config.ManagedDependencies[dep.Name]
		if !exists {
			continue
		}

		m.stopHealthCheck(dep.Name)

		if err := m.stopService(ctx, dep.Name, managedDep); err != nil {
			m.logger.Warn("Erro ao parar serviço", map[string]interface{}{
				"dependency": dep.Name,
				"error":      err.Error(),
			})
		}

		m.runningServices[dep.Name] = false
	}

	return nil
}

func (m *EnhancedManager) checkAndInstallDependency(ctx context.Context, name, version string, managedDep config.ManagedDependency) error {
	if err := m.checkDependency(ctx, name); err != nil {
		m.logger.Info("Dependência não encontrada, instalando...", map[string]interface{}{
			"dependency": name,
		})

		if err := m.installDependency(ctx, name, managedDep); err != nil {
			return fmt.Errorf("erro ao instalar dependência: %w", err)
		}
	}

	if len(managedDep.InitDatabases) > 0 {
		if err := m.initializeDatabases(ctx, name, managedDep.InitDatabases); err != nil {
			return fmt.Errorf("erro ao inicializar bancos de dados: %w", err)
		}
	}

	return nil
}

func (m *EnhancedManager) checkDependency(ctx context.Context, name string) error {
	switch name {
	case "postgres":
		cmd := exec.CommandContext(ctx, "psql", "--version")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("postgres não encontrado: %w", err)
		}
	case "redis":
		cmd := exec.CommandContext(ctx, "redis-cli", "--version")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("redis não encontrado: %w", err)
		}
	case "mongodb":
		cmd := exec.CommandContext(ctx, "mongosh", "--version")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("mongodb não encontrado: %w", err)
		}
	case "localstack":
		cmd := exec.CommandContext(ctx, "localstack", "--version")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("localstack não encontrado: %w", err)
		}
	}
	return nil
}

func (m *EnhancedManager) installDependency(ctx context.Context, name string, managedDep config.ManagedDependency) error {
	if managedDep.InstallCommand == "" {
		return fmt.Errorf("comando de instalação não definido para %s", name)
	}

	m.logger.Info("Instalando dependência", map[string]interface{}{
		"dependency": name,
		"command":    managedDep.InstallCommand,
	})

	parts := strings.Fields(managedDep.InstallCommand)
	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao executar comando de instalação: %w (output: %s)", err, string(output))
	}

	m.logger.Info("Dependência instalada com sucesso", map[string]interface{}{
		"dependency": name,
	})

	return nil
}

func (m *EnhancedManager) startService(ctx context.Context, name string, managedDep config.ManagedDependency) error {
	if managedDep.StartCommand == "" {
		return fmt.Errorf("comando de início não definido para %s", name)
	}

	m.logger.Info("Iniciando serviço", map[string]interface{}{
		"service": name,
		"command": managedDep.StartCommand,
	})

	// Executa via sh -c para suportar operadores shell (&, >, |, &&, etc.)
	cmd := exec.CommandContext(ctx, "sh", "-c", managedDep.StartCommand)

	if len(managedDep.Environment) > 0 {
		cmd.Env = os.Environ()
		for key, value := range managedDep.Environment {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
		}
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao iniciar serviço: %w (output: %s)", err, string(output))
	}

	m.logger.Info("Serviço iniciado com sucesso", map[string]interface{}{
		"service": name,
	})

	time.Sleep(2 * time.Second)

	return nil
}

func (m *EnhancedManager) stopService(ctx context.Context, name string, managedDep config.ManagedDependency) error {
	if managedDep.StopCommand == "" {
		m.logger.Info("Comando de parada não definido", map[string]interface{}{
			"service": name,
		})
		return nil
	}

	m.logger.Info("Parando serviço", map[string]interface{}{
		"service": name,
		"command": managedDep.StopCommand,
	})

	parts := strings.Fields(managedDep.StopCommand)
	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao parar serviço: %w (output: %s)", err, string(output))
	}

	m.logger.Info("Serviço parado com sucesso", map[string]interface{}{
		"service": name,
	})

	return nil
}

func (m *EnhancedManager) initializeDatabases(ctx context.Context, serviceName string, databases []config.DatabaseConfig) error {
	if serviceName != "postgres" {
		return nil
	}

	for _, db := range databases {
		m.logger.Info("Criando banco de dados", map[string]interface{}{
			"database": db.Name,
			"owner":    db.Owner,
		})

		createCmd := fmt.Sprintf("CREATE DATABASE %s;", db.Name)
		if db.Owner != "" {
			createCmd = fmt.Sprintf("CREATE DATABASE %s OWNER %s;", db.Name, db.Owner)
		}

		cmd := exec.CommandContext(ctx, "psql", "-U", "postgres", "-c", createCmd)
		output, err := cmd.CombinedOutput()
		if err != nil {
			if strings.Contains(string(output), "already exists") {
				m.logger.Info("Banco de dados já existe", map[string]interface{}{
					"database": db.Name,
				})
				continue
			}
			return fmt.Errorf("erro ao criar banco de dados %s: %w (output: %s)", db.Name, err, string(output))
		}

		m.logger.Info("Banco de dados criado com sucesso", map[string]interface{}{
			"database": db.Name,
		})
	}

	return nil
}

func (m *EnhancedManager) startHealthCheck(ctx context.Context, serviceName string) {
	healthCheck, exists := m.config.HealthChecks[serviceName]
	if !exists {
		return
	}

	interval, err := time.ParseDuration(healthCheck.Interval)
	if err != nil {
		m.logger.Warn("Interval inválido para health check", map[string]interface{}{
			"service":  serviceName,
			"interval": healthCheck.Interval,
			"error":    err.Error(),
		})
		return
	}

	ticker := time.NewTicker(interval)
	m.healthCheckers[serviceName] = ticker

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				m.performHealthCheck(ctx, serviceName, healthCheck)
			}
		}
	}()
}

func (m *EnhancedManager) stopHealthCheck(serviceName string) {
	if ticker, exists := m.healthCheckers[serviceName]; exists {
		ticker.Stop()
		delete(m.healthCheckers, serviceName)
	}
}

func (m *EnhancedManager) performHealthCheck(ctx context.Context, serviceName string, healthCheck config.HealthCheckConfig) {
	timeout, err := time.ParseDuration(healthCheck.Timeout)
	if err != nil {
		timeout = 5 * time.Second
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	parts := strings.Fields(healthCheck.Command)
	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)

	if err := cmd.Run(); err != nil {
		m.logger.Warn("Health check falhou", map[string]interface{}{
			"service": serviceName,
			"error":   err.Error(),
		})
	} else {
		m.logger.Debug("Health check bem-sucedido", map[string]interface{}{
			"service": serviceName,
		})
	}
}

// GetManagedServices retorna a lista de todos os serviços gerenciados disponíveis
func (m *EnhancedManager) GetManagedServices() []ManagedServiceInfo {
	names := make([]string, 0, len(m.config.ManagedDependencies))
	for name := range m.config.ManagedDependencies {
		names = append(names, name)
	}
	sort.Strings(names)

	services := make([]ManagedServiceInfo, 0, len(names))
	for _, name := range names {
		running := m.checkServiceStatus(name)
		m.runningServices[name] = running

		services = append(services, ManagedServiceInfo{
			Name:    name,
			Running: running,
		})
	}

	return services
}

// checkServiceStatus verifica se um serviço está realmente rodando no sistema
func (m *EnhancedManager) checkServiceStatus(serviceName string) bool {
	switch serviceName {
	case "postgres":
		return m.checkBrewService("postgresql@16")
	case "redis":
		return m.checkBrewService("redis")
	case "mongodb":
		return m.checkBrewService("mongodb-community")
	case "localstack":
		return m.checkLocalstackStatus()
	default:
		return m.runningServices[serviceName]
	}
}

// checkBrewService verifica se um serviço do homebrew está rodando
func (m *EnhancedManager) checkBrewService(serviceName string) bool {
	cmd := exec.Command("brew", "services", "list")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, serviceName) && strings.Contains(line, "started") {
			return true
		}
	}
	return false
}

// checkLocalstackStatus verifica se o LocalStack está rodando
func (m *EnhancedManager) checkLocalstackStatus() bool {
	cmd := exec.Command("localstack", "status")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.Contains(string(output), "running")
}

// StartService inicia um serviço gerenciado específico
func (m *EnhancedManager) StartService(ctx context.Context, serviceName string) error {
	if m.runningServices[serviceName] {
		return fmt.Errorf("serviço %s já está executando", serviceName)
	}

	managedDep, exists := m.config.ManagedDependencies[serviceName]
	if !exists {
		return fmt.Errorf("serviço %s não configurado", serviceName)
	}

	if err := m.startService(ctx, serviceName, managedDep); err != nil {
		return err
	}

	m.runningServices[serviceName] = true
	m.startHealthCheck(ctx, serviceName)

	return nil
}

// StopService para um serviço gerenciado específico
func (m *EnhancedManager) StopService(ctx context.Context, serviceName string) error {
	if !m.runningServices[serviceName] {
		return fmt.Errorf("serviço %s não está executando", serviceName)
	}

	managedDep, exists := m.config.ManagedDependencies[serviceName]
	if !exists {
		return fmt.Errorf("serviço %s não configurado", serviceName)
	}

	m.stopHealthCheck(serviceName)

	if err := m.stopService(ctx, serviceName, managedDep); err != nil {
		return err
	}

	m.runningServices[serviceName] = false

	return nil
}

// ManagedServiceInfo representa informações sobre um serviço gerenciado
type ManagedServiceInfo struct {
	Name    string `json:"name"`
	Running bool   `json:"running"`
}
