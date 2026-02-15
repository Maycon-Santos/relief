// Wrapper para os bindings Wails
// Este arquivo será atualizado pelo Wails após a primeira build

import type { Project, LogEntry, AppStatus } from "../types/project";

// @ts-ignore - Wails injeta estas funções globalmente
const {
  GetProjects,
  GetProject,
  StartProject,
  StopProject,
  RestartProject,
  GetProjectLogs,
  AddLocalProject,
  RemoveProject,
  RefreshConfig,
  GetStatus,
} = window.go?.main?.App || {};

export const api = {
  async getProjects(): Promise<Project[]> {
    if (!GetProjects) return [];
    return await GetProjects();
  },

  async getProject(id: string): Promise<Project> {
    if (!GetProject) throw new Error("API not available");
    return await GetProject(id);
  },

  async startProject(id: string): Promise<void> {
    if (!StartProject) throw new Error("API not available");
    return await StartProject(id);
  },

  async stopProject(id: string): Promise<void> {
    if (!StopProject) throw new Error("API not available");
    return await StopProject(id);
  },

  async restartProject(id: string): Promise<void> {
    if (!RestartProject) throw new Error("API not available");
    return await RestartProject(id);
  },

  async getProjectLogs(id: string, tail: number = 100): Promise<LogEntry[]> {
    if (!GetProjectLogs) return [];
    return await GetProjectLogs(id, tail);
  },

  async addLocalProject(path: string): Promise<void> {
    if (!AddLocalProject) throw new Error("API not available");
    return await AddLocalProject(path);
  },

  async removeProject(id: string): Promise<void> {
    if (!RemoveProject) throw new Error("API not available");
    return await RemoveProject(id);
  },

  async refreshConfig(): Promise<void> {
    if (!RefreshConfig) throw new Error("API not available");
    return await RefreshConfig();
  },

  async getStatus(): Promise<AppStatus> {
    if (!GetStatus) {
      return {
        total_projects: 0,
        running: 0,
        stopped: 0,
        errors: 0,
        traefik_running: false,
      };
    }
    return await GetStatus();
  },
};

// Declaração global para TypeScript
declare global {
  interface Window {
    go?: {
      main?: {
        App?: any;
      };
    };
  }
}
