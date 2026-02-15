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
      setError(err instanceof Error ? err.message : 'Erro ao carregar projetos');
      console.error('Erro ao carregar projetos:', err);
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
      throw new Error(err instanceof Error ? err.message : 'Erro ao iniciar projeto');
    }
  };

  const stopProject = async (id: string) => {
    try {
      await api.stopProject(id);
      await loadProjects();
    } catch (err) {
      throw new Error(err instanceof Error ? err.message : 'Erro ao parar projeto');
    }
  };

  const restartProject = async (id: string) => {
    try {
      await api.restartProject(id);
      await loadProjects();
    } catch (err) {
      throw new Error(err instanceof Error ? err.message : 'Erro ao reiniciar projeto');
    }
  };

  const removeProject = async (id: string) => {
    try {
      await api.removeProject(id);
      await loadProjects();
    } catch (err) {
      throw new Error(err instanceof Error ? err.message : 'Erro ao remover projeto');
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
