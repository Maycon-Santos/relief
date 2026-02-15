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
        console.error('Error loading status:', err);
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
            <h1 className="app-title">⚡ Relief Orchestrator</h1>
            <p className="app-subtitle">Hybrid local development orchestration</p>
          </div>
          <div className="header-actions">
            <button type="button" onClick={refresh} className="btn-secondary" disabled={loading}>
              🔄 Refresh
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
              <span className="status-label">Running:</span>
              <span className="status-value status-running">{status.running}</span>
            </div>
            <div className="status-item">
              <span className="status-label">Stopped:</span>
              <span className="status-value status-stopped">{status.stopped}</span>
            </div>
            {status.errors > 0 && (
              <div className="status-item">
                <span className="status-label">Errors:</span>
                <span className="status-value status-error">{status.errors}</span>
              </div>
            )}
            <div className="status-item">
              <span className="status-label">Traefik:</span>
              <span className={`status-value ${status.traefik_running ? 'status-running' : 'status-stopped'}`}>
                {status.traefik_running ? 'Active' : 'Inactive'}
              </span>
            </div>
          </div>
        )}
      </header>

      <main className="app-main">
        {loading && projects.length === 0 ? (
          <div className="loading">
            <div className="spinner"></div>
            <p>Loading projects...</p>
          </div>
        ) : error ? (
          <div className="error-box">
            <h3>❌ Error</h3>
            <p>{error}</p>
            <button type="button" onClick={refresh} className="btn-primary">
              Try again
            </button>
          </div>
        ) : projects.length === 0 ? (
          <div className="empty-state">
            <h2>📦 No projects found</h2>
            <p>Add a project to get started</p>
            <div className="empty-state-actions">
              <button type="button" className="btn-primary">+ Add Local Project</button>
              <button type="button" className="btn-secondary" onClick={refresh}>
                🔄 Reload
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
