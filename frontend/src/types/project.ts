import type { domain } from "../../wailsjs/go/models";

export type Project = domain.Project;
export type LogEntry = domain.LogEntry;
export type Dependency = domain.Dependency;

export type ProjectType = "docker" | "node" | "python" | "java" | "go" | "ruby";

export type ProjectStatus = "stopped" | "starting" | "running" | "error" | "unknown";

export interface GitInfo {
	is_repository: boolean;
	current_branch?: string;
	available_branches?: string[];
	remote_url?: string;
	has_changes?: boolean;
	last_commit?: string;
}

export interface AppStatus {
	total_projects: number;
	running: number;
	stopped: number;
	errors: number;
	traefik_running: boolean;
}
