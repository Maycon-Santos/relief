import { useEffect, useRef, useState } from "react";
import { Badge } from "@/components/ui/badge";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { ScrollArea } from "@/components/ui/scroll-area";
import { cn } from "@/lib/utils";
import { api } from "../services/wails";
import type { LogEntry } from "../types/project";

interface LogsViewerProps {
	projectId: string;
	projectName: string;
	onClose: () => void;
}

export function LogsViewer({ projectId, projectName, onClose }: LogsViewerProps) {
	const [logs, setLogs] = useState<LogEntry[]>([]);
	const [autoScroll, setAutoScroll] = useState(true);
	const logsEndRef = useRef<HTMLDivElement>(null);

	useEffect(() => {
		const loadLogs = async () => {
			try {
				const data = await api.getProjectLogs(projectId, 500);
				setLogs(data);
			} catch (err) {
				console.error("Error loading logs:", err);
			}
		};

		loadLogs();
		const interval = setInterval(loadLogs, 2000);

		return () => clearInterval(interval);
	}, [projectId]);

	useEffect(() => {
		if (autoScroll) {
			logsEndRef.current?.scrollIntoView({ behavior: "smooth" });
		}
	}, [autoScroll]);

	const getLevelConfig = (level: string) => {
		const normalizedLevel = level.toLowerCase();
		switch (normalizedLevel) {
			case "error":
				return {
					color: "text-red-400",
					bg: "bg-red-500/20",
					border: "border-red-500/30",
				};
			case "warn":
				return {
					color: "text-yellow-400",
					bg: "bg-yellow-500/20",
					border: "border-yellow-500/30",
				};
			case "info":
				return {
					color: "text-blue-400",
					bg: "bg-blue-600/20",
					border: "border-blue-500/40",
				};
			case "debug":
				return {
					color: "text-gray-400",
					bg: "bg-gray-600/20",
					border: "border-gray-500/30",
				};
			default:
				return {
					color: "text-gray-400",
					bg: "bg-gray-600/20",
					border: "border-gray-500/30",
				};
		}
	};

	return (
		<Dialog open onOpenChange={onClose}>
			<DialogContent className="max-w-5xl h-[80vh] flex flex-col p-0 bg-zinc-950 border-zinc-800">
				<DialogHeader className="p-6 pb-4 border-b border-zinc-800">
					<div className="flex items-center justify-between pr-6">
						<DialogTitle className="text-xl font-bold text-white">Logs: {projectName}</DialogTitle>
						<div className="flex items-center gap-3">
							<label className="flex items-center gap-2 text-sm text-gray-300 cursor-pointer select-none">
								<input
									type="checkbox"
									checked={autoScroll}
									onChange={(e) => setAutoScroll(e.target.checked)}
									className="w-4 h-4 rounded border-zinc-600 bg-zinc-800 checked:bg-blue-500 checked:border-blue-500 cursor-pointer"
								/>
								Auto-scroll
							</label>
						</div>
					</div>
				</DialogHeader>
				<ScrollArea className="flex-1 p-6">
					<div className="font-mono text-sm space-y-0.5 bg-[#0a1628] rounded-lg p-5 border border-blue-950/50 shadow-inner">
						{logs.length === 0 ? (
							<p className="text-gray-500 text-center py-8">Nenhum log disponível</p>
						) : (
							logs.map((log) => {
								const { color, bg, border } = getLevelConfig(log.level);
								return (
									<div
										key={log.id}
										className="flex gap-3 items-start py-2 hover:bg-blue-950/40 px-3 -mx-3 rounded transition-colors group"
									>
										<span className="text-gray-500 shrink-0 min-w-17.5 text-xs font-medium">
											{new Date(log.timestamp).toLocaleTimeString("en-US", {
												hour12: false,
											})}
										</span>
										<Badge
											variant="outline"
											className={cn(
												"shrink-0 min-w-17.5 justify-center font-bold text-[11px] uppercase tracking-wide",
												color,
												bg,
												border,
											)}
										>
											{log.level}
										</Badge>
										<span className="flex-1 text-gray-100 leading-relaxed">{log.message}</span>
									</div>
								);
							})
						)}
						<div ref={logsEndRef} />
					</div>
				</ScrollArea>
			</DialogContent>
		</Dialog>
	);
}
