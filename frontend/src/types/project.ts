// Tipos TypeScript para o frontend

export interface Project {
  id: string;
  name: string;
  path: string;
  domain: string;
  type: ProjectType;
  status: ProjectStatus;
  port: number;
  pid?: number;
  dependencies: Dependency[];
  scripts: Record<string, string>;
  env: Record<string, string>;
  created_at: string;
  updated_at: string;
  last_error?: string;
}

export type ProjectType = 'docker' | 'node' | 'python' | 'java' | 'go' | 'ruby';

export type ProjectStatus = 'stopped' | 'starting' | 'running' | 'error' | 'unknown';

export interface Dependency {
  name: string;
  version: string;
  required_version: string;
  managed: boolean;
  satisfied: boolean;
  message?: string;
}

export interface LogEntry {
  id: number;
  project_id: string;
  level: string;
  message: string;
  timestamp: string;
}

export interface AppStatus {
  total_projects: number;
  running: number;
  stopped: number;
  errors: number;
  traefik_running: boolean;
}
