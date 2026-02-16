import { AlertCircle, Download, GitBranch, GitCommit, RefreshCw } from "lucide-react";
import { useCallback, useEffect, useState } from "react";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { api } from "../services/wails";
import type { GitInfo, Project } from "../types/project";

interface GitControlsProps {
	project: Project;
	onGitInfoUpdate?: (gitInfo: GitInfo) => void;
}

export function GitControls({ project, onGitInfoUpdate }: GitControlsProps) {
	const [gitInfo, setGitInfo] = useState<GitInfo | null>(null);
	const [loading, setLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);
	const [syncing, setSyncing] = useState(false);

	const loadGitInfo = useCallback(async () => {
		try {
			setLoading(true);
			setError(null);
			const info = await api.getProjectGitInfo(project.id);
			setGitInfo(info);
			onGitInfoUpdate?.(info);
		} catch (err) {
			const errorMsg = err instanceof Error ? err.message : "Erro ao carregar informações Git";
			setError(errorMsg);
			setGitInfo(null);
		} finally {
			setLoading(false);
		}
	}, [project.id, onGitInfoUpdate]);

	useEffect(() => {
		loadGitInfo();
	}, [loadGitInfo]);

	const handleSync = async () => {
		try {
			setSyncing(true);
			setError(null);
			await api.syncProjectBranch(project.id);
			await loadGitInfo();
		} catch (err) {
			const errorMsg = err instanceof Error ? err.message : "Erro ao sincronizar branch";
			setError(errorMsg);
		} finally {
			setSyncing(false);
		}
	};

	if (!gitInfo?.is_repository) {
		return null;
	}

	const currentBranch = gitInfo.current_branch || "unknown";
	const hasChanges = gitInfo.has_changes || false;

	return (
		<div className="space-y-2">
			<div className="flex items-center gap-2 text-xs">
				<div className="flex items-center gap-1">
					<GitBranch className="h-3 w-3 text-blue-400" />
					<span className="text-gray-300">Branch:</span>
				</div>

				<Badge variant="secondary" className="text-xs bg-zinc-800 text-blue-400">
					{currentBranch}
				</Badge>

				{hasChanges && (
					<Badge
						variant="secondary"
						className="text-xs bg-orange-500/20 text-orange-400 border-orange-500/30"
					>
						Modificado
					</Badge>
				)}
			</div>

			<div className="flex items-center gap-2">
				<Button
					onClick={handleSync}
					disabled={loading || syncing}
					size="sm"
					variant="secondary"
					className="h-7 px-2 text-xs bg-zinc-800 hover:bg-zinc-700 text-gray-200 border-zinc-700"
				>
					{syncing ? (
						<RefreshCw className="h-3 w-3 mr-1 animate-spin" />
					) : (
						<Download className="h-3 w-3 mr-1" />
					)}
					{syncing ? "Sincronizando..." : "Pull"}
				</Button>

				{gitInfo.last_commit && (
					<div className="flex items-center gap-1 text-xs text-gray-400">
						<GitCommit className="h-3 w-3" />
						<span className="truncate max-w-30" title={gitInfo.last_commit}>
							{gitInfo.last_commit}
						</span>
					</div>
				)}
			</div>

			{error && (
				<Alert variant="destructive" className="py-2">
					<AlertCircle className="h-3 w-3" />
					<AlertDescription className="text-xs">{error}</AlertDescription>
				</Alert>
			)}
		</div>
	);
}
