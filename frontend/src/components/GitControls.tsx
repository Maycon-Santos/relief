import {
	AlertCircle,
	Check,
	ChevronsUpDown,
	Download,
	GitBranch,
	GitCommit,
	RefreshCw,
} from "lucide-react";
import { useCallback, useEffect, useRef, useState } from "react";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { EventsOff, EventsOn } from "../../wailsjs/runtime/runtime";
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
	const [checkingOut, setCheckingOut] = useState(false);
	const [branchDropdownOpen, setBranchDropdownOpen] = useState(false);
	const dropdownRef = useRef<HTMLDivElement>(null);

	const loadGitInfo = useCallback(async () => {
		try {
			setLoading(true);
			setError(null);
			const info = await api.getProjectGitInfo(project.id);
			setGitInfo(info);
			onGitInfoUpdate?.(info);
		} catch (err) {
			const errorMsg =
				err instanceof Error ? err.message : "Erro ao carregar informações Git";
			setError(errorMsg);
			setGitInfo(null);
		} finally {
			setLoading(false);
		}
	}, [project.id, onGitInfoUpdate]);

	useEffect(() => {
		loadGitInfo();
	}, [loadGitInfo]);

	// Escutar eventos de mudança de branch via watcher do backend
	useEffect(() => {
		const handler = (data: { projectId: string; branch: string }) => {
			if (data.projectId === project.id) {
				// Recarregar git info completo quando branch muda externamente
				loadGitInfo();
			}
		};

		EventsOn("git:branch-changed", handler);
		return () => {
			EventsOff("git:branch-changed");
		};
	}, [project.id, loadGitInfo]);

	// Fechar dropdown ao clicar fora
	useEffect(() => {
		const handleClickOutside = (event: MouseEvent) => {
			if (
				dropdownRef.current &&
				!dropdownRef.current.contains(event.target as Node)
			) {
				setBranchDropdownOpen(false);
			}
		};
		document.addEventListener("mousedown", handleClickOutside);
		return () => document.removeEventListener("mousedown", handleClickOutside);
	}, []);

	const handleSync = async () => {
		try {
			setSyncing(true);
			setError(null);
			await api.syncProjectBranch(project.id);
			await loadGitInfo();
		} catch (err) {
			const errorMsg =
				err instanceof Error ? err.message : "Erro ao sincronizar branch";
			setError(errorMsg);
		} finally {
			setSyncing(false);
		}
	};

	const handleCheckout = async (branch: string) => {
		try {
			setCheckingOut(true);
			setError(null);
			setBranchDropdownOpen(false);
			await api.checkoutProjectBranch(project.id, branch);
			await loadGitInfo();
		} catch (err) {
			const errorMsg =
				err instanceof Error ? err.message : "Erro ao trocar de branch";
			setError(errorMsg);
		} finally {
			setCheckingOut(false);
		}
	};

	if (!gitInfo?.is_repository) {
		return null;
	}

	const currentBranch = gitInfo.current_branch || "unknown";
	const hasChanges = gitInfo.has_changes || false;
	const branches = gitInfo.available_branches || [];

	return (
		<div className="space-y-2">
			<div className="flex items-center gap-2 text-xs">
				<div className="flex items-center gap-1">
					<GitBranch className="h-3 w-3 text-blue-400" />
					<span className="text-gray-300">Branch:</span>
				</div>

				<div className="relative" ref={dropdownRef}>
					<button
						type="button"
						onClick={() => setBranchDropdownOpen(!branchDropdownOpen)}
						disabled={checkingOut}
						className="flex items-center gap-1 text-xs bg-zinc-800 text-blue-400 px-2 py-0.5 rounded border border-zinc-700 hover:bg-zinc-700 hover:border-zinc-600 transition-colors cursor-pointer disabled:opacity-50"
					>
						{checkingOut ? (
							<RefreshCw className="h-3 w-3 animate-spin" />
						) : (
							<span className="truncate max-w-24">{currentBranch}</span>
						)}
						<ChevronsUpDown className="h-3 w-3 text-gray-400" />
					</button>

					{branchDropdownOpen && branches.length > 0 && (
						<div className="absolute z-50 top-full left-0 mt-1 w-56 max-h-48 overflow-y-auto bg-zinc-900 border border-zinc-700 rounded-lg shadow-xl">
							{branches.map((branch) => (
								<button
									key={branch}
									type="button"
									onClick={() => handleCheckout(branch)}
									className={`w-full text-left px-3 py-1.5 text-xs hover:bg-zinc-800 transition-colors flex items-center gap-2 ${
										branch === currentBranch ? "text-blue-400" : "text-gray-300"
									}`}
								>
									{branch === currentBranch && <Check className="h-3 w-3" />}
									<span className={branch === currentBranch ? "" : "ml-5"}>
										{branch}
									</span>
								</button>
							))}
						</div>
					)}
				</div>

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
