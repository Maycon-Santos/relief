import { useState, useEffect } from 'react';
import { useProjects } from './hooks/useProjects';
import { ProjectCard } from './components/ProjectCard';
import { LogsViewer } from './components/LogsViewer';
import { api } from './services/wails';
import type { AppStatus } from './types/project';
import './App.css';

function App() {
  const { projects, loading, error, startProject, stopProject, restartProject, removeProject, refresh } = useProjects();
  const [selectedProjectId, setSelectedProjectId] = useState<string | null>(null);
  const [status, setStatus] = useState<AppStatus | null>(null);

  useEffect(() => {
    const loadStatus = async () => {
      try {
        const data = await api.getStatus();
        setStatus(data);
      } catch (err) {
        console.error('Erro ao carregar status:', err);
      }
    };

    loadStatus();
    const interval = setInterval(loadStatus, 5000);

    return () => clearInterval(interval);
  }, []);

  const selectedProject = projects.find((p) => p.id === selectedProjectId);

  return (
    <div className="app">
      <header className="app-header">
        <div className="header-content">
          <div>
            <h1 className="app-title">⚡ Sofredor Orchestrator</h1>
            <p className="app-subtitle">Orquestração de desenvolvimento local híbrida</p>
          </div>
          <div className="header-actions">
            <button onClick={refresh} className="btn-secondary" disabled={loading}>
              🔄 Atualizar
            </button>
          </div>
        </div>

        {status && (
          <div className="status-bar">
            <div className="status-item">
              <span className="status-label">Total:</span>
              <span className="status-value">{status.total_projects}</span>
            </div>
            <div className="status-item">
              <span className="status-label">Rodando:</span>
              <span className="status-value status-running">{status.running}</span>
            </div>
            <div className="status-item">
              <span className="status-label">Parados:</span>
              <span className="status-value status-stopped">{status.stopped}</span>
            </div>
            {status.errors > 0 && (
              <div className="status-item">
                <span className="status-label">Erros:</span>
                <span className="status-value status-error">{status.errors}</span>
              </div>
            )}
            <div className="status-item">
              <span className="status-label">Traefik:</span>
              <span className={`status-value ${status.traefik_running ? 'status-running' : 'status-stopped'}`}>
                {status.traefik_running ? 'Ativo' : 'Inativo'}
              </span>
            </div>
          </div>
        )}
      </header>

      <main className="app-main">
        {loading && projects.length === 0 ? (
          <div className="loading">
            <div className="spinner"></div>
            <p>Carregando projetos...</p>
          </div>
        ) : error ? (
          <div className="error-box">
            <h3>❌ Erro</h3>
            <p>{error}</p>
            <button onClick={refresh} className="btn-primary">
              Tentar novamente
            </button>
          </div>
        ) : projects.length === 0 ? (
          <div className="empty-state">
            <h2>📦 Nenhum projeto encontrado</h2>
            <p>Adicione um projeto para começar</p>
            <div className="empty-state-actions">
              <button className="btn-primary">+ Adicionar Projeto Local</button>
              <button className="btn-secondary" onClick={refresh}>
                🔄 Recarregar
              </button>
            </div>
          </div>
        ) : (
          <div className="projects-grid">
            {projects.map((project) => (
              <ProjectCard
                key={project.id}
                project={project}
                onStart={() => startProject(project.id)}
                onStop={() => stopProject(project.id)}
                onRestart={() => restartProject(project.id)}
                onRemove={() => removeProject(project.id)}
                onViewLogs={() => setSelectedProjectId(project.id)}
              />
            ))}
          </div>
        )}
      </main>

      {selectedProject && (
        <LogsViewer
          projectId={selectedProject.id}
          projectName={selectedProject.name}
          onClose={() => setSelectedProjectId(null)}
        />
      )}
    </div>
  );
}

export default App;
