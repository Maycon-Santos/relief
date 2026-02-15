// Package app fornece os bindings Wails para o frontend.
package app

import (
	"context"
	"fmt"

	"github.com/omelete/relief/internal/config"
	"github.com/omelete/relief/internal/dependency"
	"github.com/omelete/relief/internal/domain"
	"github.com/omelete/relief/internal/proxy"
	"github.com/omelete/relief/internal/runner"
	"github.com/omelete/relief/internal/storage"
	"github.com/omelete/relief/pkg/logger"
)

// App é a estrutura principal da aplicação
type App struct {
	ctx              context.Context
	logger           *logger.Logger
	config           *config.Config
	configLoader     *config.Loader
	db               *storage.DB
	projectRepo      *storage.ProjectRepository
	logRepo          *storage.LogRepository
	runnerFactory    *runner.Factory
	runners          map[string]runner.ProjectRunner
	dependencyMgr    *dependency.Manager
	traefikMgr       *proxy.TraefikManager
	hostsMgr         *proxy.HostsManager
}

// NewApp cria uma nova instância da aplicação
func NewApp() *App {
	return &App{
		runners: make(map[string]runner.ProjectRunner),
	}
}

// Startup é chamado quando a aplicação inicia
func (a *App) Startup(ctx context.Context) error {
	a.ctx = ctx

	// Inicializar logger
	a.logger = logger.Default()
	a.logger.Info("Initializing Relief Orchestrator", nil)

	// Inicializar banco de dados
	db, err := storage.NewDB(a.logger)
	if err != nil {
		return fmt.Errorf("erro ao inicializar banco de dados: %w", err)
	}
	a.db = db
	a.projectRepo = storage.NewProjectRepository(db)
	a.logRepo = storage.NewLogRepository(db)

	// Inicializar config loader
	a.configLoader = config.NewLoader()

	// Carregar configuração
	configPath, _ := config.GetConfigPath()
	localConfigPath, _ := config.GetLocalConfigPath()
	
	cfg, err := a.configLoader.LoadConfig("", configPath)
	if err != nil {
		a.logger.Warn("Usando configuração padrão", map[string]interface{}{
			"error": err.Error(),
		})
		cfg = &config.Config{}
	}
	
	// Tentar merge com local
	if localCfg, err := a.configLoader.LoadConfig("", localConfigPath); err == nil {
		cfg.MergeWith(localCfg)
	}
	
	a.config = cfg

	// Inicializar runner factory
	a.runnerFactory = runner.NewFactory(a.logger)

	// Inicializar dependency manager
	a.dependencyMgr = dependency.NewManager(a.logger)

	// Inicializar Traefik manager
	traefikMgr, err := proxy.NewTraefikManager(
		a.config.Proxy.HTTPPort,
		a.config.Proxy.HTTPSPort,
		a.logger,
	)
	if err != nil {
		a.logger.Warn("Erro ao inicializar Traefik manager", map[string]interface{}{
			"error": err.Error(),
		})
	}
	a.traefikMgr = traefikMgr

	// Inicializar hosts manager
	a. hostsMgr = proxy.NewHostsManager(a.logger)

	a.logger.Info("Relief Orchestrator started successfully", nil)
	return nil
}

// Shutdown é chamado quando a aplicação fecha
func (a *App) Shutdown(ctx context.Context) error {
	a.logger.Info("Shutting down Relief Orchestrator", nil)

	// Parar todos os projetos em execução
	projects, _ := a.projectRepo.List()
	for _, project := range projects {
		if project.IsRunning() {
			a.StopProject(project.ID)
		}
	}

	// Fechar banco de dados
	if a.db != nil {
		a.db.Close()
	}

	return nil
}

// GetProjects retorna todos os projetos
func (a *App) GetProjects() ([]*domain.Project, error) {
	projects, err := a.projectRepo.List()
	if err != nil {
		return nil, fmt.Errorf("erro ao listar projetos: %w", err)
	}
	return projects, nil
}

// GetProject retorna um projeto específico
func (a *App) GetProject(id string) (*domain.Project, error) {
	project, err := a.projectRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar projeto: %w", err)
	}
	return project, nil
}

// StartProject inicia um projeto
func (a *App) StartProject(id string) error {
	a.logger.Info("Iniciando projeto", map[string]interface{}{"id": id})

	// Buscar projeto
	project, err := a.projectRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("projeto não encontrado: %w", err)
	}

	// Verificar dependências
	if err := a.dependencyMgr.CheckDependencies(a.ctx, project); err != nil {
		return fmt.Errorf("erro ao verificar dependências: %w", err)
	}

	// Atualizar projeto com resultados da verificação
	if err := a.projectRepo.Update(project); err != nil {
		return fmt.Errorf("erro ao atualizar projeto: %w", err)
	}

	// Se há dependências não satisfeitas, retornar erro
	if project.HasUnsatisfiedDependencies() {
		unsatisfied := project.GetUnsatisfiedDependencies()
		return fmt.Errorf("dependências não satisfeitas: %v", unsatisfied)
	}

	// Criar runner apropriado
	projectRunner, err := a.runnerFactory.CreateRunner(project)
	if err != nil {
		return fmt.Errorf("erro ao criar runner: %w", err)
	}

	// Iniciar projeto
	if err := projectRunner.Start(a.ctx, project); err != nil {
		project.SetError(err)
		a.projectRepo.Update(project)
		return fmt.Errorf("erro ao iniciar projeto: %w", err)
	}

	// Armazenar runner
	a.runners[project.ID] = projectRunner

	// Adicionar ao Traefik
	if a.traefikMgr != nil && project.Domain != "" {
		a.traefikMgr.AddProject(project)
	}

	// Adicionar ao hosts
	if a.hostsMgr != nil && project.Domain != "" {
		a.hostsMgr.AddEntry(project.Domain)
	}

	// Atualizar projeto no banco
	project.UpdateStatus(domain.StatusRunning)
	if err := a.projectRepo.Update(project); err != nil {
		return fmt.Errorf("erro ao atualizar status: %w", err)
	}

	return nil
}

// StopProject para um projeto
func (a *App) StopProject(id string) error {
	a.logger.Info("Parando projeto", map[string]interface{}{"id": id})

	// Buscar projeto
	project, err := a.projectRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("projeto não encontrado: %w", err)
	}

	// Buscar runner
	projectRunner, exists := a.runners[id]
	if !exists {
		return fmt.Errorf("projeto não está em execução")
	}

	// Parar projeto
	if err := projectRunner.Stop(a.ctx, id); err != nil {
		return fmt.Errorf("erro ao parar projeto: %w", err)
	}

	// Remover runner
	delete(a.runners, id)

	// Remover do Traefik
	if a.traefikMgr != nil {
		a.traefikMgr.RemoveProject(id)
	}

	// Atualizar status
	project.UpdateStatus(domain.StatusStopped)
	project.PID = 0
	if err := a.projectRepo.Update(project); err != nil {
		return fmt.Errorf("erro ao atualizar status: %w", err)
	}

	return nil
}

// RestartProject reinicia um projeto
func (a *App) RestartProject(id string) error {
	if err := a.StopProject(id); err != nil {
		// Se já estava parado, ignorar erro
		a.logger.Debug("Projeto já estava parado", map[string]interface{}{"id": id})
	}
	return a.StartProject(id)
}

// GetProjectLogs retorna logs de um projeto
func (a *App) GetProjectLogs(id string, tail int) ([]domain.LogEntry, error) {
	// Buscar runner
	projectRunner, exists := a.runners[id]
	if exists {
		// Retornar logs do buffer do runner
		return projectRunner.GetLogs(id, tail)
	}

	// Se não estiver rodando, buscar do banco
	return a.logRepo.GetByProjectID(id, tail)
}

// AddLocalProject adiciona um projeto local (fora da config)
func (a *App) AddLocalProject(path string) error {
	a.logger.Info("Adicionando projeto local", map[string]interface{}{"path": path})

	// Parsear manifest
	manifest, err := domain.ParseManifest(path)
	if err != nil {
		return fmt.Errorf("erro ao ler manifest: %w", err)
	}

	// Converter para projeto
	project := manifest.ToProject(path)

	// Verificar se já existe
	existing, _ := a.projectRepo.GetByName(project.Name)
	if existing != nil {
		return fmt.Errorf("projeto '%s' já existe", project.Name)
	}

	// Salvar no banco
	if err := a.projectRepo.Create(project); err != nil {
		return fmt.Errorf("erro ao criar projeto: %w", err)
	}

	return nil
}

// RemoveProject remove um projeto
func (a *App) RemoveProject(id string) error {
	// Parar se estiver rodando
	if _, exists := a.runners[id]; exists {
		if err := a.StopProject(id); err != nil {
			return err
		}
	}

	// Buscar projeto para pegar o domínio
	project, err := a.projectRepo.GetByID(id)
	if err == nil && project.Domain != "" {
		// Remover do hosts
		if a.hostsMgr != nil {
			a.hostsMgr.RemoveEntry(project.Domain)
		}
	}

	// Deletar do banco
	return a.projectRepo.Delete(id)
}

// RefreshConfig recarrega a configuração
func (a *App) RefreshConfig() error {
	a.logger.Info("Recarregando configuração", nil)

	configPath, _ := config.GetConfigPath()
	localConfigPath, _ := config.GetLocalConfigPath()
	
	cfg, err := a.configLoader.LoadConfig("", configPath)
	if err != nil {
		return fmt.Errorf("erro ao carregar config: %w", err)
	}
	
	// Merge com local
	if localCfg, err := a.configLoader.LoadConfig("", localConfigPath); err == nil {
		cfg.MergeWith(localCfg)
	}
	
	a.config = cfg
	return nil
}

// GetStatus retorna o status geral da aplicação
func (a *App) GetStatus() (map[string]interface{}, error) {
	projects, _ := a.projectRepo.List()
	
	running := 0
	stopped := 0
	errors := 0
	
	for _, p := range projects {
		switch p.Status {
		case domain.StatusRunning:
			running++
		case domain.StatusStopped:
			stopped++
		case domain.StatusError:
			errors++
		}
	}

	return map[string]interface{}{
		"total_projects": len(projects),
		"running":        running,
		"stopped":        stopped,
		"errors":         errors,
		"traefik_running": a.traefikMgr != nil && a.traefikMgr.IsRunning(),
	}, nil
}
