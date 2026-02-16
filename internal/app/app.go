// Package app fornece os bindings Wails para o frontend.
package app

import (
	"context"
	"fmt"

	"github.com/relief-org/relief/internal/config"
	"github.com/relief-org/relief/internal/dependency"
	"github.com/relief-org/relief/internal/domain"
	"github.com/relief-org/relief/internal/git"
	"github.com/relief-org/relief/internal/proxy"
	"github.com/relief-org/relief/internal/runner"
	"github.com/relief-org/relief/internal/storage"
	"github.com/relief-org/relief/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App é a estrutura principal da aplicação
type App struct {
	ctx            context.Context
	logger         *logger.Logger
	config         *config.Config
	configLoader   *config.Loader
	db             *storage.DB
	projectRepo    *storage.ProjectRepository
	logRepo        *storage.LogRepository
	runnerFactory  *runner.Factory
	runners        map[string]runner.ProjectRunner
	dependencyMgr  *dependency.Manager
	enhancedDepMgr *dependency.EnhancedManager
	gitManager     *git.Manager
	traefikMgr     *proxy.TraefikManager
	hostsMgr       *proxy.HostsManager
}

// NewApp cria uma nova instância da aplicação
func NewApp() *App {
	return &App{
		runners: make(map[string]runner.ProjectRunner),
	}
}

// Startup é chamado quando a aplicação inicia
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	// Inicializar logger
	a.logger = logger.Default()
	a.logger.Info("Initializing Relief Orchestrator", nil)

	// Inicializar banco de dados
	db, err := storage.NewDB(a.logger)
	if err != nil {
		a.logger.Fatal("Failed to initialize database", err, nil)
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
		cfg = &config.Config{
			Proxy: config.ProxyConfig{
				HTTPPort:  80,
				HTTPSPort: 443,
			},
		}
	}

	// Garantir valores padrão para proxy se não estiverem definidos
	if cfg.Proxy.HTTPPort == 0 {
		cfg.Proxy.HTTPPort = 80
	}
	if cfg.Proxy.HTTPSPort == 0 {
		cfg.Proxy.HTTPSPort = 443
	}

	// Tentar merge com local
	if localCfg, err := a.configLoader.LoadConfig("", localConfigPath); err == nil {
		cfg.MergeWith(localCfg)
	}

	a.config = cfg

	// Inicializar Git manager
	a.gitManager = git.NewManager(a.logger)

	// Inicializar runner factory
	a.runnerFactory = runner.NewFactory(a.logger)

	// Inicializar dependency manager
	a.dependencyMgr = dependency.NewManager(a.logger)

	// Inicializar enhanced dependency manager
	a.enhancedDepMgr = dependency.NewEnhancedManager(a.logger, cfg)

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
	} else {
		// Iniciar Traefik
		if err := traefikMgr.Start(a.ctx); err != nil {
			a.logger.Warn("Erro ao iniciar Traefik", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}
	a.traefikMgr = traefikMgr

	// Inicializar hosts manager
	a.hostsMgr = proxy.NewHostsManager(a.logger)

	// Limpar processos órfãos de execuções anteriores
	a.cleanupOrphanProcesses()

	// Sincronizar projetos da configuração
	a.syncConfigProjects()

	a.logger.Info("Relief Orchestrator started successfully", nil)
}

// Shutdown é chamado quando a aplicação fecha
func (a *App) Shutdown(ctx context.Context) {
	a.logger.Info("Shutting down Relief Orchestrator", nil)

	// Parar todos os projetos em execução
	projects, _ := a.projectRepo.List()
	for _, project := range projects {
		if project.IsRunning() || project.PID > 0 {
			a.logger.Info("Parando projeto no shutdown", map[string]interface{}{
				"project": project.Name,
				"pid":     project.PID,
			})
			a.StopProject(project.ID)
		}
	}

	// Parar Traefik
	if a.traefikMgr != nil {
		if err := a.traefikMgr.Stop(); err != nil {
			a.logger.Warn("Erro ao parar Traefik", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

	// Fechar banco de dados
	if a.db != nil {
		a.db.Close()
	}
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
	defer func() {
		if r := recover(); r != nil {
			panicErr := fmt.Errorf("panic: %v", r)
			a.logger.Error("Panic in StartProject", panicErr, map[string]interface{}{
				"id": id,
			})
		}
	}()

	if a.logger == nil || a.projectRepo == nil || a.dependencyMgr == nil || a.runnerFactory == nil {
		return fmt.Errorf("application not fully initialized")
	}

	a.logger.Info("Iniciando projeto", map[string]interface{}{"id": id})

	// Buscar projeto
	project, err := a.projectRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("projeto não encontrado: %w", err)
	}

	if project == nil {
		return fmt.Errorf("project is nil")
	}

	// Garantir que o manifest está carregado
	if project.Manifest == nil {
		a.logger.Warn("Manifest not loaded, attempting to load", map[string]interface{}{
			"project": project.Name,
			"path":    project.Path,
		})
		manifest, err := domain.ParseManifest(project.Path)
		if err != nil {
			return fmt.Errorf("erro ao carregar manifest: %w", err)
		}
		project.Manifest = manifest
	}

	// Verificar dependências
	if err := a.dependencyMgr.CheckDependencies(a.ctx, project); err != nil {
		return fmt.Errorf("erro ao verificar dependências: %w", err)
	}

	// Iniciar dependências gerenciadas
	if err := a.enhancedDepMgr.StartManagedDependencies(a.ctx, project); err != nil {
		return fmt.Errorf("erro ao iniciar dependências gerenciadas: %w", err)
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

	// Verificar se a porta está em uso
	if project.Port > 0 {
		conflict, err := a.CheckPortInUse(project.Port)
		if err != nil {
			a.logger.Warn("Erro ao verificar porta", map[string]interface{}{
				"port":  project.Port,
				"error": err.Error(),
			})
		} else if conflict != nil {
			return fmt.Errorf("PORT_IN_USE:%d:%d:%s", conflict.Port, conflict.PID, conflict.Command)
		}
	}

	// Criar runner apropriado
	projectRunner, err := a.runnerFactory.CreateRunner(project)
	if err != nil {
		return fmt.Errorf("erro ao criar runner: %w", err)
	}

	a.logger.Info("Runner created successfully", map[string]interface{}{
		"project": project.Name,
	})

	// Iniciar projeto
	if err := projectRunner.Start(a.ctx, project); err != nil {
		project.SetError(err)
		a.projectRepo.Update(project)
		return fmt.Errorf("erro ao iniciar projeto: %w", err)
	}

	a.logger.Info("Project started successfully", map[string]interface{}{
		"project": project.Name,
	})

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
	if exists {
		// Parar projeto via runner
		if err := projectRunner.Stop(a.ctx, id); err != nil {
			a.logger.Warn("Erro ao parar via runner", map[string]interface{}{
				"error": err.Error(),
			})
		}
		// Remover runner
		delete(a.runners, id)
	} else if project.PID > 0 {
		// Se não tem runner mas tem PID, matar processo diretamente
		a.logger.Info("Runner não encontrado, matando processo pelo PID", map[string]interface{}{
			"pid": project.PID,
		})
		if err := a.KillProcessByPID(project.PID); err != nil {
			a.logger.Warn("Erro ao matar processo", map[string]interface{}{
				"pid":   project.PID,
				"error": err.Error(),
			})
		}
	}

	// Remover do Traefik
	if a.traefikMgr != nil {
		a.traefikMgr.RemoveProject(id)
	}

	// Parar dependências gerenciadas
	if err := a.enhancedDepMgr.StopManagedDependencies(a.ctx, project); err != nil {
		a.logger.Warn("Erro ao parar dependências gerenciadas", map[string]interface{}{
			"project": project.Name,
			"error":   err.Error(),
		})
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

	// Verificar se o path foi fornecido
	if path == "" {
		return fmt.Errorf("no directory selected")
	}

	// Parsear manifest
	manifest, err := domain.ParseManifest(path)
	if err != nil {
		return fmt.Errorf("failed to read relief.yaml in selected directory: %w", err)
	}

	// Converter para projeto
	project := manifest.ToProject(path)

	// Verificar se já existe
	existing, _ := a.projectRepo.GetByName(project.Name)
	if existing != nil {
		return fmt.Errorf("project '%s' already exists", project.Name)
	}

	// Obter informações Git
	if gitInfo, err := a.gitManager.GetGitInfo(a.ctx, path); err != nil {
		a.logger.Warn("Erro ao obter informações Git para novo projeto", map[string]interface{}{
			"path":  path,
			"error": err.Error(),
		})
	} else {
		project.UpdateGitInfo(gitInfo)
	}

	// Salvar no banco
	if err := a.projectRepo.Create(project); err != nil {
		return fmt.Errorf("failed to save project: %w", err)
	}

	a.logger.Info("Projeto adicionado com sucesso", map[string]interface{}{
		"name": project.Name,
		"path": path,
	})

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
		"total_projects":  len(projects),
		"running":         running,
		"stopped":         stopped,
		"errors":          errors,
		"traefik_running": a.traefikMgr != nil && a.traefikMgr.IsRunning(),
	}, nil
}

// SelectProjectDirectory abre um diálogo para selecionar um diretório
func (a *App) SelectProjectDirectory() (string, error) {
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Project Directory",
	})
	if err != nil {
		return "", fmt.Errorf("error opening directory dialog: %w", err)
	}
	return path, nil
}

// cleanupOrphanProcesses mata processos órfãos de execuções anteriores
func (a *App) cleanupOrphanProcesses() {
	a.logger.Info("Verificando processos órfãos...", nil)

	projects, err := a.projectRepo.List()
	if err != nil {
		a.logger.Warn("Erro ao listar projetos para limpeza", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	for _, project := range projects {
		cleaned := false

		// Se o projeto tem um PID registrado
		if project.PID > 0 {
			a.logger.Info("Encontrado processo órfão com PID registrado", map[string]interface{}{
				"project": project.Name,
				"pid":     project.PID,
				"status":  project.Status,
			})

			// Tentar matar o processo
			if err := a.KillProcessByPID(project.PID); err != nil {
				a.logger.Warn("Erro ao matar processo órfão", map[string]interface{}{
					"project": project.Name,
					"pid":     project.PID,
					"error":   err.Error(),
				})
			} else {
				a.logger.Info("Processo órfão encerrado", map[string]interface{}{
					"project": project.Name,
					"pid":     project.PID,
				})
				cleaned = true
			}
		}

		// Verificar se a porta do projeto está em uso (processo órfão sem PID registrado)
		if project.Port > 0 {
			conflict, err := a.CheckPortInUse(project.Port)
			if err == nil && conflict != nil {
				a.logger.Info("Encontrado processo órfão usando porta do projeto", map[string]interface{}{
					"project": project.Name,
					"port":    project.Port,
					"pid":     conflict.PID,
					"command": conflict.Command,
				})

				// Tentar matar o processo
				if err := a.KillProcessByPID(conflict.PID); err != nil {
					a.logger.Warn("Erro ao matar processo órfão pela porta", map[string]interface{}{
						"project": project.Name,
						"port":    project.Port,
						"pid":     conflict.PID,
						"error":   err.Error(),
					})
				} else {
					a.logger.Info("Processo órfão pela porta encerrado", map[string]interface{}{
						"project": project.Name,
						"port":    project.Port,
						"pid":     conflict.PID,
					})
					cleaned = true
				}
			}
		}

		// Se limpou algo, atualizar o projeto
		if cleaned {
			// Limpar PID e status do projeto
			project.PID = 0
			project.UpdateStatus(domain.StatusStopped)
			if err := a.projectRepo.Update(project); err != nil {
				a.logger.Warn("Erro ao atualizar projeto após limpeza", map[string]interface{}{
					"project": project.Name,
					"error":   err.Error(),
				})
			}
		}
	}

	a.logger.Info("Limpeza de processos órfãos concluída", nil)
}

// syncConfigProjects sincroniza projetos da configuração com o banco de dados e clona repositórios se necessário
func (a *App) syncConfigProjects() {
	a.logger.Info("Sincronizando projetos da configuração", nil)

	for _, projectConfig := range a.config.Projects {
		// Verificar se o projeto já existe no banco
		existingProject, err := a.projectRepo.GetByName(projectConfig.Name)
		if err != nil && err.Error() != "project not found" {
			a.logger.Warn("Erro ao buscar projeto existente", map[string]interface{}{
				"project": projectConfig.Name,
				"error":   err.Error(),
			})
			continue
		}

		// Clonar repositório se necessário
		if projectConfig.Repository != nil && projectConfig.Repository.AutoClone {
			if err := a.ensureRepositoryCloned(projectConfig); err != nil {
				a.logger.Warn("Erro ao clonar repositório", map[string]interface{}{
					"project": projectConfig.Name,
					"repo":    projectConfig.Repository.URL,
					"error":   err.Error(),
				})
				continue
			}
		}

		// Criar ou atualizar projeto
		if existingProject == nil {
			// Criar novo projeto
			project := a.createProjectFromConfig(projectConfig)
			if err := a.projectRepo.Create(project); err != nil {
				a.logger.Warn("Erro ao salvar novo projeto", map[string]interface{}{
					"project": projectConfig.Name,
					"error":   err.Error(),
				})
			} else {
				a.logger.Info("Projeto criado da configuração", map[string]interface{}{
					"project": projectConfig.Name,
				})
			}
		} else {
			// Atualizar projeto existente
			a.updateProjectFromConfig(existingProject, projectConfig)
			if err := a.projectRepo.Update(existingProject); err != nil {
				a.logger.Warn("Erro ao atualizar projeto", map[string]interface{}{
					"project": projectConfig.Name,
					"error":   err.Error(),
				})
			} else {
				a.logger.Info("Projeto atualizado da configuração", map[string]interface{}{
					"project": projectConfig.Name,
				})
			}
		}
	}
}

// ensureRepositoryCloned clona um repositório se ele não existir
func (a *App) ensureRepositoryCloned(projectConfig config.ProjectConfig) error {
	if projectConfig.Repository == nil {
		return nil
	}

	a.logger.Info("Verificando repositório", map[string]interface{}{
		"project": projectConfig.Name,
		"path":    projectConfig.Path,
		"repo":    projectConfig.Repository.URL,
	})

	return a.gitManager.CloneOrUpdate(
		a.ctx,
		projectConfig.Repository.URL,
		projectConfig.Path,
		projectConfig.Repository.Branch,
	)
}

// createProjectFromConfig cria um novo projeto a partir da configuração
func (a *App) createProjectFromConfig(projectConfig config.ProjectConfig) *domain.Project {
	project := domain.NewProject(
		projectConfig.Name,
		projectConfig.Path,
		projectConfig.Domain,
		domain.ProjectType(projectConfig.Type),
	)

	// Configurar porta
	project.Port = projectConfig.Port

	// Configurar scripts
	project.Scripts = make(map[string]string)
	for k, v := range projectConfig.Scripts {
		project.Scripts[k] = v
	}

	// Configurar env
	project.Env = make(map[string]string)
	for k, v := range projectConfig.Env {
		project.Env[k] = v
	}

	// Configurar dependências
	project.Dependencies = make([]domain.Dependency, 0, len(projectConfig.Dependencies))
	for _, dep := range projectConfig.Dependencies {
		project.Dependencies = append(project.Dependencies, domain.Dependency{
			Name:            dep.Name,
			Version:         "",
			RequiredVersion: dep.Version,
			Managed:         dep.Managed,
			Satisfied:       false,
		})
	}

	// Obter informações Git se o diretório existir
	if a.gitManager != nil {
		if gitInfo, err := a.gitManager.GetGitInfo(a.ctx, projectConfig.Path); err != nil {
			a.logger.Debug("Erro ao obter informações Git para projeto da config", map[string]interface{}{
				"project": projectConfig.Name,
				"path":    projectConfig.Path,
				"error":   err.Error(),
			})
		} else {
			project.UpdateGitInfo(gitInfo)
		}
	}

	return project
}

// updateProjectFromConfig atualiza um projeto existente a partir da configuração
func (a *App) updateProjectFromConfig(project *domain.Project, projectConfig config.ProjectConfig) {
	// Atualizar configurações básicas
	project.Path = projectConfig.Path
	project.Domain = projectConfig.Domain
	project.Type = domain.ProjectType(projectConfig.Type)
	project.Port = projectConfig.Port

	// Atualizar scripts
	project.Scripts = make(map[string]string)
	for k, v := range projectConfig.Scripts {
		project.Scripts[k] = v
	}

	// Atualizar env
	project.Env = make(map[string]string)
	for k, v := range projectConfig.Env {
		project.Env[k] = v
	}

	// Atualizar dependências
	project.Dependencies = make([]domain.Dependency, 0, len(projectConfig.Dependencies))
	for _, dep := range projectConfig.Dependencies {
		project.Dependencies = append(project.Dependencies, domain.Dependency{
			Name:            dep.Name,
			Version:         "",
			RequiredVersion: dep.Version,
			Managed:         dep.Managed,
			Satisfied:       false,
		})
	}

	// Atualizar informações Git
	if a.gitManager != nil {
		if gitInfo, err := a.gitManager.GetGitInfo(a.ctx, projectConfig.Path); err != nil {
			a.logger.Debug("Erro ao obter informações Git para projeto atualizado", map[string]interface{}{
				"project": projectConfig.Name,
				"path":    projectConfig.Path,
				"error":   err.Error(),
			})
		} else {
			project.UpdateGitInfo(gitInfo)
		}
	}
}

// StartProjectDependencies inicia as dependências gerenciadas de um projeto
func (a *App) StartProjectDependencies(id string) error {
	project, err := a.projectRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("projeto não encontrado: %w", err)
	}

	return a.enhancedDepMgr.StartManagedDependencies(a.ctx, project)
}

// StopProjectDependencies para as dependências gerenciadas de um projeto
func (a *App) StopProjectDependencies(id string) error {
	project, err := a.projectRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("projeto não encontrado: %w", err)
	}

	return a.enhancedDepMgr.StopManagedDependencies(a.ctx, project)
}

// SyncRepository sincroniza um repositório (git pull/fetch)
func (a *App) SyncRepository(id string) error {
	project, err := a.projectRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("projeto não encontrado: %w", err)
	}

	// Buscar configuração do projeto
	projectConfig := a.config.GetProjectByName(project.Name)
	if projectConfig == nil || projectConfig.Repository == nil {
		return fmt.Errorf("projeto não possui configuração de repositório")
	}

	a.logger.Info("Sincronizando repositório", map[string]interface{}{
		"project": project.Name,
		"path":    project.Path,
	})

	return a.gitManager.CloneOrUpdate(
		a.ctx,
		projectConfig.Repository.URL,
		project.Path,
		projectConfig.Repository.Branch,
	)
}

// GetProjectGitInfo obtém informações Git de um projeto
func (a *App) GetProjectGitInfo(id string) (*domain.GitInfo, error) {
	project, err := a.projectRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("projeto não encontrado: %w", err)
	}

	return a.gitManager.GetGitInfo(a.ctx, project.Path)
}

// CheckoutProjectBranch faz checkout para uma branch específica em um projeto
func (a *App) CheckoutProjectBranch(id, branch string) error {
	project, err := a.projectRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("projeto não encontrado: %w", err)
	}

	a.logger.Info("Fazendo checkout da branch", map[string]interface{}{
		"project": project.Name,
		"branch":  branch,
		"path":    project.Path,
	})

	err = a.gitManager.CheckoutBranch(a.ctx, project.Path, branch)
	if err != nil {
		return err
	}

	// Atualizar informações Git do projeto
	gitInfo, err := a.gitManager.GetGitInfo(a.ctx, project.Path)
	if err != nil {
		a.logger.Warn("Erro ao atualizar informações Git após checkout", map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		project.UpdateGitInfo(gitInfo)
		if err := a.projectRepo.Update(project); err != nil {
			a.logger.Warn("Erro ao salvar informações Git atualizadas", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

	return nil
}

// SyncProjectBranch sincroniza a branch atual de um projeto (git pull)
func (a *App) SyncProjectBranch(id string) error {
	project, err := a.projectRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("projeto não encontrado: %w", err)
	}

	a.logger.Info("Sincronizando branch do projeto", map[string]interface{}{
		"project": project.Name,
		"path":    project.Path,
	})

	err = a.gitManager.SyncBranch(a.ctx, project.Path)
	if err != nil {
		return err
	}

	// Atualizar informações Git do projeto
	gitInfo, err := a.gitManager.GetGitInfo(a.ctx, project.Path)
	if err != nil {
		a.logger.Warn("Erro ao atualizar informações Git após sync", map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		project.UpdateGitInfo(gitInfo)
		if err := a.projectRepo.Update(project); err != nil {
			a.logger.Warn("Erro ao salvar informações Git atualizadas", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

	return nil
}

// RefreshProjectGitInfo atualiza as informações Git de um projeto
func (a *App) RefreshProjectGitInfo(id string) error {
	project, err := a.projectRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("projeto não encontrado: %w", err)
	}

	gitInfo, err := a.gitManager.GetGitInfo(a.ctx, project.Path)
	if err != nil {
		return fmt.Errorf("erro ao obter informações Git: %w", err)
	}

	project.UpdateGitInfo(gitInfo)
	return a.projectRepo.Update(project)
}
