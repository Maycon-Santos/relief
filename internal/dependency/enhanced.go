package dependency

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Maycon-Santos/relief/internal/config"
	"github.com/Maycon-Santos/relief/internal/domain"
	"github.com/Maycon-Santos/relief/pkg/logger"
	"github.com/Maycon-Santos/relief/pkg/shellenv"
)



type LogFunc func(level, message string)

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

func (m *EnhancedManager) StartManagedDependencies(ctx context.Context, project *domain.Project, logFn LogFunc) error {
	for _, dep := range project.Dependencies {
		if !dep.Managed {
			continue
		}

		if m.runningServices[dep.Name] {
			m.logger.Info("Dependência já está executando", map[string]interface{}{
				"dependency": dep.Name,
			})
			if logFn != nil {
				logFn("info", fmt.Sprintf("[dep:%s] já está em execução, pulando", dep.Name))
			}
			continue
		}

		managedDep, exists := m.config.ManagedDependencies[dep.Name]
		if !exists {
			msg := fmt.Sprintf("[dep:%s] configuração não encontrada em managed_dependencies", dep.Name)
			m.logger.Warn("Configuração não encontrada para dependência gerenciada", map[string]interface{}{
				"dependency": dep.Name,
			})
			if logFn != nil {
				logFn("warn", msg)
			}
			continue
		}

		if err := m.checkAndInstallDependency(ctx, dep.Name, dep.Version, managedDep, logFn); err != nil {
			if logFn != nil {
				logFn("error", fmt.Sprintf("[dep:%s] falha ao verificar/instalar: %s", dep.Name, err.Error()))
			}
			return fmt.Errorf("erro ao verificar/instalar dependência %s: %w", dep.Name, err)
		}

		if err := m.startService(ctx, dep.Name, managedDep, logFn); err != nil {
			if logFn != nil {
				logFn("error", fmt.Sprintf("[dep:%s] falha ao iniciar serviço: %s", dep.Name, err.Error()))
			}
			return fmt.Errorf("erro ao iniciar serviço %s: %w", dep.Name, err)
		}

		m.startHealthCheck(ctx, dep.Name)

		m.runningServices[dep.Name] = true
	}

	return nil
}

func (m *EnhancedManager) StopManagedDependencies(ctx context.Context, project *domain.Project, depsInUse map[string]bool) error {
	for _, dep := range project.Dependencies {
		if !dep.Managed || !m.runningServices[dep.Name] {
			continue
		}

		if depsInUse[dep.Name] {
			m.logger.Info("Dependência ainda em uso por outro projeto, mantendo ativa", map[string]interface{}{
				"dependency": dep.Name,
			})
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

func (m *EnhancedManager) checkAndInstallDependency(ctx context.Context, name, version string, managedDep config.ManagedDependency, logFn LogFunc) error {
	if err := m.checkDependency(ctx, name); err != nil {
		m.logger.Info("Dependência não encontrada, instalando...", map[string]interface{}{
			"dependency": name,
		})
		if logFn != nil {
			logFn("info", fmt.Sprintf("[dep:%s] não encontrado, iniciando instalação...", name))
		}

		if err := m.installDependency(ctx, name, managedDep, logFn); err != nil {
			return fmt.Errorf("erro ao instalar dependência: %w", err)
		}
	}

	if len(managedDep.InitDatabases) > 0 {
		if err := m.initializeDatabases(ctx, name, managedDep.InitDatabases, logFn); err != nil {
			return fmt.Errorf("erro ao inicializar bancos de dados: %w", err)
		}
	}

	return nil
}

func (m *EnhancedManager) checkDependency(ctx context.Context, name string) error {
	var probe string
	switch name {
	case "postgres":
		probe = "psql --version"
	case "redis":
		probe = "redis-cli --version"
	case "mongodb":
		probe = "mongosh --version"
	case "localstack":
		probe = "localstack --version"
	default:
		return nil
	}
	cmd := shellenv.CommandContext(ctx, probe)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s não encontrado: %w", name, err)
	}
	return nil
}

func (m *EnhancedManager) installDependency(ctx context.Context, name string, managedDep config.ManagedDependency, logFn LogFunc) error {
	if managedDep.InstallCommand == "" {
		err := fmt.Errorf("comando de instalação não definido para %s", name)
		if logFn != nil {
			logFn("error", fmt.Sprintf("[dep:%s] %s", name, err.Error()))
		}
		return err
	}

	m.logger.Info("Instalando dependência", map[string]interface{}{
		"dependency": name,
		"command":    managedDep.InstallCommand,
	})
	if logFn != nil {
		logFn("info", fmt.Sprintf("[dep:%s] instalando: %s", name, managedDep.InstallCommand))
	}

	cmd := shellenv.CommandContext(ctx, managedDep.InstallCommand)
	output, err := cmd.CombinedOutput()
	if err != nil {
		msg := fmt.Sprintf("[dep:%s] falha na instalação: %s\n%s", name, err.Error(), strings.TrimSpace(string(output)))
		if logFn != nil {
			logFn("error", msg)
		}
		return fmt.Errorf("erro ao executar comando de instalação: %w (output: %s)", err, string(output))
	}

	m.logger.Info("Dependência instalada com sucesso", map[string]interface{}{
		"dependency": name,
	})
	if logFn != nil {
		logFn("info", fmt.Sprintf("[dep:%s] instalada com sucesso", name))
	}

	return nil
}

func (m *EnhancedManager) startService(ctx context.Context, name string, managedDep config.ManagedDependency, logFn LogFunc) error {
	if managedDep.StartCommand == "" {
		err := fmt.Errorf("comando de início não definido para %s", name)
		if logFn != nil {
			logFn("error", fmt.Sprintf("[dep:%s] %s", name, err.Error()))
		}
		return err
	}

	m.logger.Info("Iniciando serviço", map[string]interface{}{
		"service": name,
		"command": managedDep.StartCommand,
	})
	if logFn != nil {
		logFn("info", fmt.Sprintf("[dep:%s] iniciando: %s", name, managedDep.StartCommand))
	}

	
	
	cmd := shellenv.CommandContext(ctx, managedDep.StartCommand)

	if len(managedDep.Environment) > 0 {
		cmd.Env = os.Environ()
		for key, value := range managedDep.Environment {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
		}
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		msg := fmt.Sprintf("[dep:%s] falha ao iniciar: %s", name, err.Error())
		if out := strings.TrimSpace(string(output)); out != "" {
			msg += "\n" + out
		}
		if logFn != nil {
			logFn("error", msg)
		}
		return fmt.Errorf("erro ao iniciar serviço: %w (output: %s)", err, string(output))
	}

	m.logger.Info("Serviço iniciado com sucesso", map[string]interface{}{
		"service": name,
	})
	if logFn != nil {
		logFn("info", fmt.Sprintf("[dep:%s] iniciado com sucesso", name))
	}

	if managedDep.PostStartCommand != "" {
		postCmd := managedDep.PostStartCommand
		go func() {
			out, err := shellenv.Command(postCmd).CombinedOutput()
			if err != nil {
				m.logger.Warn("post_start_command falhou", map[string]interface{}{
					"service": name,
					"error":   err.Error(),
					"output":  strings.TrimSpace(string(out)),
				})
				if logFn != nil {
					logFn("warn", fmt.Sprintf("[dep:%s] post_start_command falhou: %s", name, err.Error()))
				}
			} else {
				m.logger.Info("post_start_command executado com sucesso", map[string]interface{}{"service": name})
				if logFn != nil {
					logFn("info", fmt.Sprintf("[dep:%s] recursos inicializados", name))
				}
			}
		}()
	}

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

	cmd := shellenv.CommandContext(ctx, managedDep.StopCommand)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao parar serviço: %w (output: %s)", err, string(output))
	}

	m.logger.Info("Serviço parado com sucesso", map[string]interface{}{
		"service": name,
	})

	return nil
}

func (m *EnhancedManager) initializeDatabases(ctx context.Context, serviceName string, databases []config.DatabaseConfig, logFn LogFunc) error {
	if serviceName != "postgres" {
		return nil
	}

	for _, db := range databases {
		m.logger.Info("Criando banco de dados", map[string]interface{}{
			"database": db.Name,
			"owner":    db.Owner,
		})
		if logFn != nil {
			logFn("info", fmt.Sprintf("[dep:%s] criando banco de dados: %s", serviceName, db.Name))
		}

		createCmd := fmt.Sprintf("CREATE DATABASE %s;", db.Name)
		if db.Owner != "" {
			createCmd = fmt.Sprintf("CREATE DATABASE %s OWNER %s;", db.Name, db.Owner)
		}

		cmd := shellenv.CommandContext(ctx, fmt.Sprintf("psql -U postgres -c %q", createCmd))
		output, err := cmd.CombinedOutput()
		if err != nil {
			if strings.Contains(string(output), "already exists") {
				m.logger.Info("Banco de dados já existe", map[string]interface{}{
					"database": db.Name,
				})
				if logFn != nil {
					logFn("info", fmt.Sprintf("[dep:%s] banco %s já existe, pulando", serviceName, db.Name))
				}
				continue
			}
			if logFn != nil {
				logFn("error", fmt.Sprintf("[dep:%s] falha ao criar banco %s: %s\n%s", serviceName, db.Name, err.Error(), strings.TrimSpace(string(output))))
			}
			return fmt.Errorf("erro ao criar banco de dados %s: %w (output: %s)", db.Name, err, string(output))
		}

		m.logger.Info("Banco de dados criado com sucesso", map[string]interface{}{
			"database": db.Name,
		})
		if logFn != nil {
			logFn("info", fmt.Sprintf("[dep:%s] banco %s criado com sucesso", serviceName, db.Name))
		}
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

	cmd := shellenv.CommandContext(ctx, healthCheck.Command)

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


func (m *EnhancedManager) checkBrewService(serviceName string) bool {
	cmd := shellenv.Command("brew services list")
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

func (m *EnhancedManager) checkLocalstackStatus() bool {
	cmd := shellenv.Command("curl -sf --max-time 2 http://localhost:4566/_localstack/health")
	err := cmd.Run()
	return err == nil
}

func (m *EnhancedManager) StartService(ctx context.Context, serviceName string) error {
	if m.runningServices[serviceName] {
		return fmt.Errorf("serviço %s já está executando", serviceName)
	}

	managedDep, exists := m.config.ManagedDependencies[serviceName]
	if !exists {
		return fmt.Errorf("serviço %s não configurado", serviceName)
	}

	if err := m.startService(ctx, serviceName, managedDep, nil); err != nil {
		return err
	}

	m.runningServices[serviceName] = true
	m.startHealthCheck(ctx, serviceName)

	return nil
}

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

type ManagedServiceInfo struct {
	Name    string `json:"name"`
	Running bool   `json:"running"`
}
