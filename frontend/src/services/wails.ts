// Wrapper para os bindings Wails

import * as App from "../../wailsjs/go/app/App";
import type { AppStatus, LogEntry, Project } from "../types/project";

export interface PortConflict {
  port: number;
  pid: number;
  command: string;
}

export const api = {
  async getProjects(): Promise<Project[]> {
    return await App.GetProjects();
  },

  async getProject(id: string): Promise<Project> {
    return await App.GetProject(id);
  },

  async startProject(id: string): Promise<void> {
    return await App.StartProject(id);
  },

  async stopProject(id: string): Promise<void> {
    return await App.StopProject(id);
  },

  async restartProject(id: string): Promise<void> {
    return await App.RestartProject(id);
  },

  async getProjectLogs(id: string, tail: number = 100): Promise<LogEntry[]> {
    return await App.GetProjectLogs(id, tail);
  },

  async addLocalProject(path: string): Promise<void> {
    return await App.AddLocalProject(path);
  },

  async removeProject(id: string): Promise<void> {
    return await App.RemoveProject(id);
  },

  async refreshConfig(): Promise<void> {
    return await App.RefreshConfig();
  },

  async getStatus(): Promise<AppStatus> {
    return (await App.GetStatus()) as AppStatus;
  },

  async selectProjectDirectory(): Promise<string> {
    return await App.SelectProjectDirectory();
  },

  async checkPortInUse(port: number): Promise<PortConflict | null> {
    return (await App.CheckPortInUse(port)) as PortConflict | null;
  },

  async killProcessByPID(pid: number): Promise<void> {
    return await App.KillProcessByPID(pid);
  },
};
