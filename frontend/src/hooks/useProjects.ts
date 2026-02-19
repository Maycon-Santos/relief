import { useCallback, useEffect, useRef, useState } from "react";
import { api } from "../services/wails";
import type { Project } from "../types/project";

export function useProjects() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const projectsRef = useRef<Project[]>([]);

  const loadProjects = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await api.getProjects();
      setProjects(data);
      projectsRef.current = data;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Error loading projects");
      console.error("Error loading projects:", err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadProjects();

    // Intervalo dinâmico: 3 s quando há projetos em estado transiente, 15 s caso contrário.
    const schedule = () => {
      const transient = projectsRef.current.some(
        (p) => p.status === "starting" || p.status === "running",
      );
      return setTimeout(
        async () => {
          await loadProjects();
          timer = schedule();
        },
        transient ? 3000 : 15000,
      );
    };

    let timer = schedule();
    return () => clearTimeout(timer);
  }, [loadProjects]);

  const startProject = async (id: string) => {
    try {
      await api.startProject(id);
      await loadProjects();
    } catch (err) {
      const msg =
        err instanceof Error
          ? err.message
          : typeof err === "string"
            ? err
            : "Error starting project";
      throw new Error(msg);
    }
  };

  const stopProject = async (id: string) => {
    try {
      await api.stopProject(id);
      await loadProjects();
    } catch (err) {
      const msg =
        err instanceof Error
          ? err.message
          : typeof err === "string"
            ? err
            : "Error stopping project";
      throw new Error(msg);
    }
  };

  const restartProject = async (id: string) => {
    try {
      await api.restartProject(id);
      await loadProjects();
    } catch (err) {
      const msg =
        err instanceof Error
          ? err.message
          : typeof err === "string"
            ? err
            : "Error restarting project";
      throw new Error(msg);
    }
  };

  const removeProject = async (id: string) => {
    try {
      await api.removeProject(id);
      await loadProjects();
    } catch (err) {
      throw new Error(
        err instanceof Error ? err.message : "Error removing project",
      );
    }
  };

  const addProject = async () => {
    try {
      const path = await api.selectProjectDirectory();
      if (path) {
        await api.addLocalProject(path);
        await loadProjects();
      }
    } catch (err) {
      throw new Error(
        err instanceof Error ? err.message : "Error adding project",
      );
    }
  };

  const updateProject = useCallback((updatedProject: Project) => {
    setProjects((prevProjects) =>
      prevProjects.map((project) =>
        project.id === updatedProject.id ? updatedProject : project,
      ),
    );
  }, []);

  return {
    projects,
    loading,
    error,
    startProject,
    stopProject,
    restartProject,
    removeProject,
    addProject,
    updateProject,
    refresh: loadProjects,
  };
}
