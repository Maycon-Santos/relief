import {
	AlertCircle,
	ExternalLink,
	FileText,
	Play,
	RotateCw,
	Square,
	Trash2,
} from "lucide-react";
import { useState } from "react";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
	Card,
	CardContent,
	CardFooter,
	CardHeader,
} from "@/components/ui/card";
import { BrowserOpenURL } from "../../wailsjs/runtime/runtime";
import { api, type PortConflict } from "../services/wails";
import type { Project } from "../types/project";
import { DependencyAlert } from "./DependencyAlert";
import { GitControls } from "./GitControls";
import { PortConflictModal } from "./PortConflictModal";

interface ProjectCardProps {
	project: Project;
	onStart: () => Promise<void>;
	onStop: () => Promise<void>;
	onRestart: () => Promise<void>;
	onRemove: () => Promise<void>;
	onViewLogs: () => void;
}

export function ProjectCard({
	project,
	onStart,
	onStop,
	onRestart,
	onRemove,
	onViewLogs,
}: ProjectCardProps) {
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);
	const [portConflict, setPortConflict] = useState<PortConflict | null>(null);

	const handleGitInfoUpdate = (gitInfo: any) => {
		// Callback para quando as informações Git são atualizadas
		console.log("Git info updated:", gitInfo);
	};

	const handleAction = async (
		action: () => Promise<void>,
		actionName: string,
	) => {
		try {
			setLoading(true);
			setError(null);
			await action();
		} catch (err) {
			const errorMsg =
				err instanceof Error ? err.message : `Error on ${actionName}`;

			// Detectar erro de porta em uso
			if (errorMsg.startsWith("PORT_IN_USE:")) {
				const parts = errorMsg.split(":");
				if (parts.length >= 4) {
					setPortConflict({
						port: parseInt(parts[1]),
						pid: parseInt(parts[2]),
						command: parts.slice(3).join(":"),
					});
				}
			} else {
				setError(errorMsg);
			}
		} finally {
			setLoading(false);
		}
	};

	const handleKillProcess = async () => {
		if (!portConflict) return;

		try {
			await api.killProcessByPID(portConflict.pid);
			setPortConflict(null);
			// Tentar iniciar novamente
			await handleAction(onStart, "start");
		} catch (err) {
			throw err;
		}
	};

	const unsatisfiedDeps = project.dependencies.filter((d) => !d.satisfied);
	const isRunning = project.status === "running";
	const isStopped = project.status === "stopped";

	const getStatusBadge = () => {
		if (isRunning) {
			return (
				<Badge className="bg-green-500/20 text-green-400 border-green-500/30">
					Running
				</Badge>
			);
		}
		return (
			<Badge variant="secondary" className="text-gray-400 bg-zinc-800/50">
				Stopped
			</Badge>
		);
	};

	return (
		<Card className="overflow-hidden border-border/50 bg-zinc-900/50 hover:border-primary/40 transition-colors flex flex-col">
			<CardHeader className="pb-4">
				<div className="flex items-start justify-between gap-3">
					<div className="flex-1 min-w-0">
						<h3 className="text-xl font-bold mb-2 text-white">
							{project.name}
						</h3>
						<p className="text-sm text-gray-400 truncate">{project.path}</p>
					</div>
					{getStatusBadge()}
				</div>
			</CardHeader>

			<CardContent className="space-y-3 flex-1 pb-4">
				{project.domain && (
					<button
						type="button"
						onClick={() => BrowserOpenURL(`http://${project.domain}`)}
						className="flex items-center gap-2 text-sm text-blue-400 hover:text-blue-300 transition-colors font-medium cursor-pointer"
					>
						<ExternalLink className="h-4 w-4" />
						{project.domain}
					</button>
				)}

				<div className="flex flex-wrap gap-2">
					<Badge
						variant="secondary"
						className="text-xs font-normal bg-zinc-800/80 text-gray-300 border-zinc-700/50"
					>
						{project.type === "node" ? "Node.js" : project.type}
					</Badge>
					{project.port > 0 && (
						<Badge
							variant="secondary"
							className="text-xs font-normal bg-zinc-800/80 text-gray-300 border-zinc-700/50"
						>
							Port {project.port}
						</Badge>
					)}
					{project.pid && (
						<Badge
							variant="secondary"
							className="text-xs font-normal bg-zinc-800/80 text-gray-300 border-zinc-700/50"
						>
							PID: {project.pid}
						</Badge>
					)}
				</div>

				{/* Git Controls */}
				<GitControls project={project} onGitInfoUpdate={handleGitInfoUpdate} />

				{unsatisfiedDeps.length > 0 && (
					<DependencyAlert dependencies={unsatisfiedDeps} />
				)}

				{error && (
					<Alert variant="destructive">
						<AlertCircle className="h-4 w-4" />
						<AlertDescription className="text-xs">{error}</AlertDescription>
					</Alert>
				)}

				{project.last_error && (
					<Alert variant="destructive">
						<AlertCircle className="h-4 w-4" />
						<AlertDescription className="text-xs">
							{project.last_error}
						</AlertDescription>
					</Alert>
				)}
			</CardContent>

			<CardFooter className="bg-zinc-900/80 border-t border-zinc-800 px-4 py-3">
				<div className="flex items-center gap-2 w-full flex-wrap">
					{/* Action Buttons - Left Group */}
					{isStopped && (
						<Button
							onClick={() => handleAction(onStart, "iniciar")}
							disabled={loading}
							size="sm"
							variant="secondary"
							className="bg-zinc-800 hover:bg-zinc-700 text-gray-200 border-zinc-700"
						>
							<Play className="h-4 w-4 mr-1.5" />
							Start
						</Button>
					)}

					{isRunning && (
						<>
							<Button
								onClick={() => handleAction(onStop, "parar")}
								disabled={loading}
								size="sm"
								variant="secondary"
								className="bg-zinc-800 hover:bg-zinc-700 text-gray-200 border-zinc-700"
							>
								<Square className="h-4 w-4 mr-1.5" />
								Stop
							</Button>

							<Button
								onClick={() => handleAction(onRestart, "reiniciar")}
								disabled={loading}
								size="sm"
								variant="secondary"
								className="bg-zinc-800 hover:bg-zinc-700 text-gray-200 border-zinc-700"
							>
								<RotateCw className="h-4 w-4 mr-1.5" />
								Restart
							</Button>

							<Button
								onClick={onViewLogs}
								size="sm"
								variant="secondary"
								className="bg-zinc-800 hover:bg-zinc-700 text-gray-200 border-zinc-700"
							>
								<FileText className="h-4 w-4 mr-1.5" />
								View Logs
							</Button>
						</>
					)}

					{/* Remove Button - Far Right */}
					<Button
						onClick={() => handleAction(onRemove, "remover")}
						disabled={loading || isRunning}
						size="sm"
						variant="ghost"
						className="ml-auto text-red-400 hover:text-red-300 hover:bg-red-500/10"
					>
						<Trash2 className="h-4 w-4 mr-1.5" />
						Remove
					</Button>
				</div>
			</CardFooter>

			{portConflict && (
				<PortConflictModal
					conflict={portConflict}
					onKill={handleKillProcess}
					onCancel={() => setPortConflict(null)}
				/>
			)}
		</Card>
	);
}
