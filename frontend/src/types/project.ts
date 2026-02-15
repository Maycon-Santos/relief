// Re-export tipos gerados pelo Wails
import type { domain } from "../../wailsjs/go/models";

export type Project = domain.Project;
export type LogEntry = domain.LogEntry;
export type Dependency = domain.Dependency;

// Tipos auxiliares
export type ProjectType = "docker" | "node" | "python" | "java" | "go" | "ruby";

export type ProjectStatus =
  | "stopped"
  | "starting"
  | "running"
  | "error"
  | "unknown";

export interface AppStatus {
  total_projects: number;
  running: number;
  stopped: number;
  errors: number;
  traefik_running: boolean;
}
