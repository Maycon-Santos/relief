import { useState } from "react";
import { api, type PortConflict } from "../services/wails";
import type { Project } from "../types/project";
import { DependencyAlert } from "./DependencyAlert";
import { PortConflictModal } from "./PortConflictModal";
import { StatusBadge } from "./StatusBadge";

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

	return (
		<div
			style={{
				backgroundColor: "#1f2937",
				borderRadius: "12px",
				padding: "20px",
				marginBottom: "16px",
				border: "1px solid #374151",
			}}
		>
			<div
				style={{
					display: "flex",
					justifyContent: "space-between",
					alignItems: "flex-start",
					marginBottom: "16px",
				}}
			>
				<div>
					<h3
						style={{
							margin: "0 0 8px 0",
							fontSize: "18px",
							fontWeight: "600",
							color: "#f9fafb",
						}}
					>
						{project.name}
					</h3>
					<p style={{ margin: "0", fontSize: "14px", color: "#9ca3af" }}>
						{project.path}
					</p>
				</div>
				<StatusBadge status={project.status} />
			</div>

			{project.domain && (
				<div style={{ marginBottom: "12px" }}>
					<span style={{ fontSize: "14px", color: "#9ca3af" }}>Domínio: </span>
					<a
						href={`http://${project.domain}`}
						target="_blank"
						rel="noopener noreferrer"
						style={{
							fontSize: "14px",
							color: "#3b82f6",
							textDecoration: "none",
						}}
					>
						{project.domain}
					</a>
				</div>
			)}

			<div
				style={{
					display: "flex",
					gap: "8px",
					marginBottom: "12px",
					flexWrap: "wrap",
				}}
			>
				<span
					style={{
						fontSize: "12px",
						padding: "4px 8px",
						backgroundColor: "#374151",
						borderRadius: "6px",
						color: "#d1d5db",
					}}
				>
					{project.type}
				</span>
				{project.port > 0 && (
					<span
						style={{
							fontSize: "12px",
							padding: "4px 8px",
							backgroundColor: "#374151",
							borderRadius: "6px",
							color: "#d1d5db",
						}}
					>
						Porta: {project.port}
					</span>
				)}
				{project.pid && (
					<span
						style={{
							fontSize: "12px",
							padding: "4px 8px",
							backgroundColor: "#374151",
							borderRadius: "6px",
							color: "#d1d5db",
						}}
					>
						PID: {project.pid}
					</span>
				)}
			</div>

			{unsatisfiedDeps.length > 0 && (
				<DependencyAlert dependencies={unsatisfiedDeps} />
			)}

			{error && (
				<div
					style={{
						padding: "12px",
						backgroundColor: "#7f1d1d",
						borderRadius: "8px",
						marginBottom: "12px",
					}}
				>
					<p style={{ margin: "0", fontSize: "14px", color: "#fca5a5" }}>
						{error}
					</p>
				</div>
			)}

			{project.last_error && (
				<div
					style={{
						padding: "12px",
						backgroundColor: "#7f1d1d",
						borderRadius: "8px",
						marginBottom: "12px",
					}}
				>
					<p style={{ margin: "0", fontSize: "14px", color: "#fca5a5" }}>
						{project.last_error}
					</p>
				</div>
			)}

			<div style={{ display: "flex", gap: "8px", flexWrap: "wrap" }}>
				{project.status === "stopped" && (
					<button
						type="button"
						onClick={() => handleAction(onStart, "iniciar")}
						disabled={loading}
						style={{
							padding: "8px 16px",
							backgroundColor: "#10b981",
							color: "white",
							border: "none",
							borderRadius: "8px",
							fontSize: "14px",
							fontWeight: "500",
							cursor: loading ? "not-allowed" : "pointer",
							opacity: loading ? 0.6 : 1,
						}}
					>
						{loading ? "Iniciando..." : "Iniciar"}
					</button>
				)}

				{project.status === "running" && (
					<>
						<button
							type="button"
							onClick={() => handleAction(onStop, "parar")}
							disabled={loading}
							style={{
								padding: "8px 16px",
								backgroundColor: "#ef4444",
								color: "white",
								border: "none",
								borderRadius: "8px",
								fontSize: "14px",
								fontWeight: "500",
								cursor: loading ? "not-allowed" : "pointer",
								opacity: loading ? 0.6 : 1,
							}}
						>
							{loading ? "Parando..." : "Parar"}
						</button>

						<button
							type="button"
							onClick={() => handleAction(onRestart, "reiniciar")}
							disabled={loading}
							style={{
								padding: "8px 16px",
								backgroundColor: "#f59e0b",
								color: "white",
								border: "none",
								borderRadius: "8px",
								fontSize: "14px",
								fontWeight: "500",
								cursor: loading ? "not-allowed" : "pointer",
								opacity: loading ? 0.6 : 1,
							}}
						>
							{loading ? "Reiniciando..." : "Reiniciar"}
						</button>

						<button
							type="button"
							onClick={onViewLogs}
							style={{
								padding: "8px 16px",
								backgroundColor: "#6b7280",
								color: "white",
								border: "none",
								borderRadius: "8px",
								fontSize: "14px",
								fontWeight: "500",
								cursor: "pointer",
							}}
						>
							View Logs
						</button>
					</>
				)}

				<button
					type="button"
					onClick={() => handleAction(onRemove, "remover")}
					disabled={loading || project.status === "running"}
					style={{
						padding: "8px 16px",
						backgroundColor: "transparent",
						color: "#ef4444",
						border: "1px solid #ef4444",
						borderRadius: "8px",
						fontSize: "14px",
						fontWeight: "500",
						cursor:
							loading || project.status === "running"
								? "not-allowed"
								: "pointer",
						opacity: loading || project.status === "running" ? 0.4 : 1,
						marginLeft: "auto",
					}}
				>
					Remove
				</button>
			</div>

			{portConflict && (
				<PortConflictModal
					conflict={portConflict}
					onKill={handleKillProcess}
					onCancel={() => setPortConflict(null)}
				/>
			)}
		</div>
	);
}
