import { useState, useEffect } from 'react';
import { api } from '../services/wails';
import type { Project } from '../types/project';

export function useProjects() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadProjects = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await api.getProjects();
      setProjects(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error loading projects');
      console.error('Error loading projects:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadProjects();
    
    // Atualizar a cada 5 segundos
    const interval = setInterval(loadProjects, 5000);
    
    return () => clearInterval(interval);
  }, []);

  const startProject = async (id: string) => {
    try {
      await api.startProject(id);
      await loadProjects();
    } catch (err) {
      throw new Error(err instanceof Error ? err.message : 'Error starting project');
    }
  };

  const stopProject = async (id: string) => {
    try {
      await api.stopProject(id);
      await loadProjects();
    } catch (err) {
      throw new Error(err instanceof Error ? err.message : 'Error stopping project');
    }
  };

  const restartProject = async (id: string) => {
    try {
      await api.restartProject(id);
      await loadProjects();
    } catch (err) {
      throw new Error(err instanceof Error ? err.message : 'Error restarting project');
    }
  };

  const removeProject = async (id: string) => {
    try {
      await api.removeProject(id);
      await loadProjects();
    } catch (err) {
      throw new Error(err instanceof Error ? err.message : 'Error removing project');
    }
  };

  return {
    projects,
    loading,
    error,
    startProject,
    stopProject,
    restartProject,
    removeProject,
    refresh: loadProjects,
  };
}
