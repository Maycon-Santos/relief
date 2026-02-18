package app

import (
	"context"
	"fmt"

	"github.com/Maycon-Santos/relief/internal/config"
	"github.com/Maycon-Santos/relief/internal/dependency"
	"github.com/Maycon-Santos/relief/internal/domain"
	"github.com/Maycon-Santos/relief/internal/git"
	"github.com/Maycon-Santos/relief/internal/proxy"
	"github.com/Maycon-Santos/relief/internal/runner"
	"github.com/Maycon-Santos/relief/internal/storage"
	"github.com/Maycon-Santos/relief/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

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

func NewApp() *App {
	return &App{
		runners: make(map[string]runner.ProjectRunner),
	}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	a.logger = logger.Default()
	a.logger.Info("Initializing Relief Orchestrator", nil)

	db, err := storage.NewDB(a.logger)
	if err != nil {
		a.logger.Fatal("Failed to initialize database", err, nil)
	}
	a.db = db
	a.projectRepo = storage.NewProjectRepository(db)
	a.logRepo = storage.NewLogRepository(db)

	a.configLoader = config.NewLoader()

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

	if cfg.Proxy.HTTPPort == 0 {
		cfg.Proxy.HTTPPort = 80
	}
	if cfg.Proxy.HTTPSPort == 0 {
		cfg.Proxy.HTTPSPort = 443
	}

	if localCfg, err := a.configLoader.LoadConfig("", localConfigPath); err == nil {
		cfg.MergeWith(localCfg)
	}

	a.config = cfg

	a.gitManager = git.NewManager(a.logger)

	a.runnerFactory = runner.NewFactory(a.logger)

	a.dependencyMgr = dependency.NewManager(a.logger)

	a.enhancedDepMgr = dependency.NewEnhancedManager(a.logger, cfg)

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
		if err := traefikMgr.Start(a.ctx); err != nil {
			a.logger.Warn("Erro ao iniciar Traefik", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}
	a.traefikMgr = traefikMgr

	a.hostsMgr = proxy.NewHostsManager(a.logger)

	a.cleanupOrphanProcesses()

	a.syncConfigProjects()

	a.logger.Info("Relief Orchestrator started successfully", nil)
}

func (a *App) Shutdown(ctx context.Context) {
	a.logger.Info("Shutting down Relief Orchestrator", nil)

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

	if a.traefikMgr != nil {
		if err := a.traefikMgr.Stop(); err != nil {
			a.logger.Warn("Erro ao parar Traefik", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}

	if a.db != nil {
		a.db.Close()
	}
}

func (a *App) GetProjects() ([]*domain.Project, error) {
	projects, err := a.projectRepo.List()
	if err != nil {
		return nil, fmt.Errorf("erro ao listar projetos: %w", err)
	}
	return projects, nil
}

func (a *App) GetProject(id string) (*domain.Project, error) {
	project, err := a.projectRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar projeto: %w", err)
	}
	return project, nil
}

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

	project, err := a.projectRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("projeto não encontrado: %w", err)
	}

	if project == nil {
		return fmt.Errorf("project is nil")
	}

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

	if err := a.dependencyMgr.CheckDependencies(a.ctx, project); err != nil {
		return fmt.Errorf("erro ao verificar dependências: %w", err)
	}

	if err := a.enhancedDepMgr.StartManagedDependencies(a.ctx, project); err != nil {
		return fmt.Errorf("erro ao iniciar dependências gerenciadas: %w", err)
	}

	if err := a.projectRepo.Update(project); err != nil {
		return fmt.Errorf("erro ao atualizar projeto: %w", err)
	}

	if project.HasUnsatisfiedDependencies() {
		unsatisfied := project.GetUnsatisfiedDependencies()
		return fmt.Errorf("dependências não satisfeitas: %v", unsatisfied)
	}

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

	projectRunner, err := a.runnerFactory.CreateRunner(project)
	if err != nil {
		return fmt.Errorf("erro ao criar runner: %w", err)
	}

	a.logger.Info("Runner created successfully", map[string]interface{}{
		"project": project.Name,
	})

	if err := projectRunner.Start(a.ctx, project); err != nil {
		project.SetError(err)
		a.projectRepo.Update(project)
		return fmt.Errorf("erro ao iniciar projeto: %w", err)
	}

	a.logger.Info("Project started successfully", map[string]interface{}{
		"project": project.Name,
	})

	a.runners[project.ID] = projectRunner

	if a.traefikMgr != nil && project.Domain != "" {
		a.traefikMgr.AddProject(project)
	}

	if a.hostsMgr != nil && project.Domain != "" {
		a.hostsMgr.AddEntry(project.Domain)
	}

	project.UpdateStatus(domain.StatusRunning)
	if err := a.projectRepo.Update(project); err != nil {
		return fmt.Errorf("erro ao atualizar status: %w", err)
	}

	return nil
}

func (a *App) StopProject(id string) error {
	a.logger.Info("Parando projeto", map[string]interface{}{"id": id})

	project, err := a.projectRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("projeto não encontrado: %w", err)
	}

	projectRunner, exists := a.runners[id]
	if exists {
		if err := projectRunner.Stop(a.ctx, id); err != nil {
			a.logger.Warn("Erro ao parar via runner", map[string]interface{}{
				"error": err.Error(),
			})
		}
		delete(a.runners, id)
	} else if project.PID > 0 {
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

	if a.traefikMgr != nil {
		a.traefikMgr.RemoveProject(id)
	}

	if err := a.enhancedDepMgr.StopManagedDependencies(a.ctx, project); err != nil {
		a.logger.Warn("Erro ao parar dependências gerenciadas", map[string]interface{}{
			"project": project.Name,
			"error":   err.Error(),
		})
	}

	project.UpdateStatus(domain.StatusStopped)
	project.PID = 0
	if err := a.projectRepo.Update(project); err != nil {
		return fmt.Errorf("erro ao atualizar status: %w", err)
	}

	return nil
}

func (a *App) RestartProject(id string) error {
	if err := a.StopProject(id); err != nil {
		a.logger.Debug("Projeto já estava parado", map[string]interface{}{"id": id})
	}
	return a.StartProject(id)
}

func (a *App) GetProjectLogs(id string, tail int) ([]domain.LogEntry, error) {
	projectRunner, exists := a.runners[id]
	if exists {
		return projectRunner.GetLogs(id, tail)
	}

	return a.logRepo.GetByProjectID(id, tail)
}

func (a *App) AddLocalProject(path string) error {
	a.logger.Info("Adicionando projeto local", map[string]interface{}{"path": path})

	if path == "" {
		return fmt.Errorf("no directory selected")
	}

	manifest, err := domain.ParseManifest(path)
	if err != nil {
		return fmt.Errorf("failed to read relief.yaml in selected directory: %w", err)
	}

	project := manifest.ToProject(path)

	existing, _ := a.projectRepo.GetByName(project.Name)
	if existing != nil {
		return fmt.Errorf("project '%s' already exists", project.Name)
	}

	if gitInfo, err := a.gitManager.GetGitInfo(a.ctx, path); err != nil {
		a.logger.Warn("Erro ao obter informações Git para novo projeto", map[string]interface{}{
			"path":  path,
			"error": err.Error(),
		})
	} else {
		project.UpdateGitInfo(gitInfo)
	}

	if err := a.projectRepo.Create(project); err != nil {
		return fmt.Errorf("failed to save project: %w", err)
	}

	a.logger.Info("Projeto adicionado com sucesso", map[string]interface{}{
		"name": project.Name,
		"path": path,
	})

	return nil
}

func (a *App) RemoveProject(id string) error {
	if _, exists := a.runners[id]; exists {
		if err := a.StopProject(id); err != nil {
			return err
		}
	}

	project, err := a.projectRepo.GetByID(id)
	if err == nil && project.Domain != "" {
		if a.hostsMgr != nil {
			a.hostsMgr.RemoveEntry(project.Domain)
		}
	}

	return a.projectRepo.Delete(id)
}

func (a *App) RefreshConfig() error {
	a.logger.Info("Recarregando configuração", nil)

	configPath, _ := config.GetConfigPath()
	localConfigPath, _ := config.GetLocalConfigPath()

	cfg, err := a.configLoader.LoadConfig("", configPath)
	if err != nil {
		return fmt.Errorf("erro ao carregar config: %w", err)
	}

	if localCfg, err := a.configLoader.LoadConfig("", localConfigPath); err == nil {
		cfg.MergeWith(localCfg)
	}

	a.config = cfg
	return nil
}

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

func (a *App) SelectProjectDirectory() (string, error) {
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Project Directory",
	})
	if err != nil {
		return "", fmt.Errorf("error opening directory dialog: %w", err)
	}
	return path, nil
}

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

		if project.PID > 0 {
			a.logger.Info("Encontrado processo órfão com PID registrado", map[string]interface{}{
				"project": project.Name,
				"pid":     project.PID,
				"status":  project.Status,
			})

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

		if project.Port > 0 {
			conflict, err := a.CheckPortInUse(project.Port)
			if err == nil && conflict != nil {
				a.logger.Info("Encontrado processo órfão usando porta do projeto", map[string]interface{}{
					"project": project.Name,
					"port":    project.Port,
					"pid":     conflict.PID,
					"command": conflict.Command,
				})

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

		if cleaned {
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

func (a *App) syncConfigProjects() {
	a.logger.Info("Sincronizando projetos da configuração", nil)

	for _, projectConfig := range a.config.Projects {
		existingProject, err := a.projectRepo.GetByName(projectConfig.Name)
		if err != nil && err.Error() != "project not found" {
			a.logger.Warn("Erro ao buscar projeto existente", map[string]interface{}{
				"project": projectConfig.Name,
				"error":   err.Error(),
			})
			continue
		}

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

		if existingProject == nil {
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

func (a *App) createProjectFromConfig(projectConfig config.ProjectConfig) *domain.Project {
	project := domain.NewProject(
		projectConfig.Name,
		projectConfig.Path,
		projectConfig.Domain,
		domain.ProjectType(projectConfig.Type),
	)

	project.Port = projectConfig.Port

	project.Scripts = make(map[string]string)
	for k, v := range projectConfig.Scripts {
		project.Scripts[k] = v
	}

	project.Env = make(map[string]string)
	for k, v := range projectConfig.Env {
		project.Env[k] = v
	}

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

func (a *App) updateProjectFromConfig(project *domain.Project, projectConfig config.ProjectConfig) {
	project.Path = projectConfig.Path
	project.Domain = projectConfig.Domain
	project.Type = domain.ProjectType(projectConfig.Type)
	project.Port = projectConfig.Port

	project.Scripts = make(map[string]string)
	for k, v := range projectConfig.Scripts {
		project.Scripts[k] = v
	}

	project.Env = make(map[string]string)
	for k, v := range projectConfig.Env {
		project.Env[k] = v
	}

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

func (a *App) StartProjectDependencies(id string) error {
	project, err := a.projectRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("projeto não encontrado: %w", err)
	}

	return a.enhancedDepMgr.StartManagedDependencies(a.ctx, project)
}

func (a *App) StopProjectDependencies(id string) error {
	project, err := a.projectRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("projeto não encontrado: %w", err)
	}

	return a.enhancedDepMgr.StopManagedDependencies(a.ctx, project)
}

func (a *App) SyncRepository(id string) error {
	project, err := a.projectRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("projeto não encontrado: %w", err)
	}

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

func (a *App) GetProjectGitInfo(id string) (*domain.GitInfo, error) {
	project, err := a.projectRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("projeto não encontrado: %w", err)
	}

	return a.gitManager.GetGitInfo(a.ctx, project.Path)
}

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
