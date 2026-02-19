import { Folder, Loader2, Package, Plus, RefreshCw, Settings, Zap } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { ConfigEditor } from "./components/ConfigEditor";
import { GlobalScripts } from "./components/GlobalScripts";
import { LogsViewer } from "./components/LogsViewer";
import { ManagedServices } from "./components/ManagedServices";
import { ProjectCard } from "./components/ProjectCard";
import { useProjects } from "./hooks/useProjects";
import { api } from "./services/wails";
import type { AppStatus } from "./types/project";

function App() {
	const {
		projects,
		loading,
		error,
		startProject,
		stopProject,
		restartProject,
		removeProject,
		addProject,
		refresh,
	} = useProjects();
	const [selectedProjectId, setSelectedProjectId] = useState<string | null>(null);
	const [status, setStatus] = useState<AppStatus | null>(null);
	const [managedServices, setManagedServices] = useState<Array<{ name: string; running: boolean }>>(
		[],
	);
	const [configEditorOpen, setConfigEditorOpen] = useState(false);

	useEffect(() => {
		const loadStatus = async () => {
			try {
				const data = await api.getStatus();
				setStatus(data);
			} catch (err) {
				console.error("Error loading status:", err);
			}
		};

		const loadServices = async () => {
			try {
				const services = await api.getManagedServices();
				setManagedServices(services);
			} catch (err) {
				console.error("Error loading managed services:", err);
			}
		};

		loadStatus();
		loadServices();
		const interval = setInterval(() => {
			loadStatus();
			loadServices();
		}, 5000);

		return () => clearInterval(interval);
	}, []);

	const prevStatusesRef = useRef<Record<string, string>>({});
	useEffect(() => {
		const prev = prevStatusesRef.current;
		for (const project of projects) {
			if (project.status === "error" && prev[project.id] && prev[project.id] !== "error") {
				setSelectedProjectId(project.id);
			}
		}
		prevStatusesRef.current = Object.fromEntries(projects.map((p) => [p.id, p.status]));
	}, [projects]);

	const selectedProject = projects.find((p) => p.id === selectedProjectId);

	const handleAddProject = async () => {
		try {
			await addProject();
		} catch (err) {
			console.error("Error adding project:", err);
			const message = err instanceof Error ? err.message : "Error adding project";
			alert(`Failed to add project:\n\n${message}`);
		}
	};

	const handleStartService = async (serviceName: string) => {
		try {
			await api.startManagedService(serviceName);
			const services = await api.getManagedServices();
			setManagedServices(services);
		} catch (err) {
			console.error("Error starting service:", err);
			throw err;
		}
	};

	const handleStopService = async (serviceName: string) => {
		try {
			await api.stopManagedService(serviceName);
			const services = await api.getManagedServices();
			setManagedServices(services);
		} catch (err) {
			console.error("Error stopping service:", err);
			throw err;
		}
	};

	const handleRestartTraefik = async () => {
		try {
			await api.restartTraefik();
			const data = await api.getStatus();
			setStatus(data);
		} catch (err) {
			console.error("Error restarting traefik:", err);
		}
	};

	const handleRefresh = async () => {
		try {
			await api.reloadConfig();
		} catch (err) {
			console.error("Error reloading config:", err);
		}
		refresh();
	};

	return (
		<div className="min-h-screen bg-background text-foreground">
			<header className="border-b border-zinc-800 bg-black">
				<div className="container mx-auto max-w-7xl px-6 py-6">
					<div className="flex items-center justify-between mb-6">
						<div className="flex items-center gap-3">
							<Zap className="h-8 w-8 text-white" />
							<div>
								<h1 className="text-2xl font-bold tracking-tight text-white">
									Relief Orchestrator
								</h1>
								<p className="text-sm text-gray-400">Hybrid local development orchestration</p>
							</div>
						</div>
						<div className="flex gap-3">
							<Button
								onClick={() => setConfigEditorOpen(true)}
								disabled={loading}
								variant="outline"
								size="default"
								className="border-zinc-700 text-gray-300 hover:bg-zinc-800 hover:text-white transition-colors"
							>
								<Settings className="h-4 w-4 mr-2" />
								Settings
							</Button>
							<Button
								onClick={handleAddProject}
								disabled={loading}
								size="default"
								className="bg-white hover:bg-gray-200 text-black font-semibold transition-colors"
							>
								<Plus className="h-4 w-4 mr-2" />
								Add Project
							</Button>
							<Button
							onClick={handleRefresh}
								disabled={loading}
								variant="outline"
								size="default"
								className="border-zinc-700 text-gray-300 hover:bg-zinc-800 hover:text-white transition-colors"
							>
								<RefreshCw className={`h-4 w-4 mr-2 ${loading ? "animate-spin" : ""}`} />
								Refresh
							</Button>
						</div>
					</div>

					{status && (
						<Card className="border border-zinc-800/50 bg-zinc-950/50">
							<div className="flex items-center gap-6 px-6 py-4 flex-wrap">
								<div className="flex items-center gap-2 text-sm">
									<span className="text-gray-400">Total Projects:</span>
									<span className="font-semibold text-white">{status.total_projects}</span>
								</div>
								<div className="flex items-center gap-2 text-sm">
									<div className="w-2 h-2 rounded-full bg-green-500" />
									<span className="text-gray-400">Running:</span>
									<Badge
										variant="secondary"
										className="bg-green-500/20 text-green-400 border-green-500/30 font-semibold"
									>
										{status.running}
									</Badge>
								</div>
								<div className="flex items-center gap-2 text-sm">
									<span className="text-gray-400">Stopped:</span>
									<span className="font-semibold text-white">{status.stopped}</span>
								</div>
								<div className="flex items-center gap-2 text-sm">
									<span className="text-gray-400">Traefik:</span>
									<Badge
										variant="secondary"
										className={
											status.traefik_running
												? "bg-green-500/20 text-green-400 border-green-500/30 font-semibold"
												: "text-gray-400"
										}
									>
										{status.traefik_running ? "Active" : "Inactive"}
									</Badge>
									{!status.traefik_running && (
										<button
											type="button"
											onClick={handleRestartTraefik}
											className="text-xs text-yellow-400 underline hover:text-yellow-300 transition-colors"
										>
											Fix
										</button>
									)}
								</div>
							</div>
						</Card>
					)}
				</div>
			</header>

			<main className="container mx-auto max-w-7xl px-6 py-8">
				{managedServices.length > 0 && (
					<div className="mb-8">
						<ManagedServices
							services={managedServices}
							onStartService={handleStartService}
							onStopService={handleStopService}
						/>
					</div>
				)}

				<GlobalScripts />

				<div className="mb-6 flex items-center gap-2 text-gray-500">
					<Folder className="h-5 w-5" />
					<h2 className="text-lg font-semibold">Projects</h2>
				</div>
				{loading && projects.length === 0 ? (
					<div className="flex flex-col items-center justify-center py-20">
						<Loader2 className="h-12 w-12 animate-spin text-primary mb-4" />
						<p className="text-muted-foreground">Loading projects...</p>
					</div>
				) : error ? (
					<Card className="p-8 text-center max-w-md mx-auto">
						<div className="text-5xl mb-4">❌</div>
						<h3 className="text-xl font-semibold mb-2">Error</h3>
						<p className="text-muted-foreground mb-6">{error}</p>
						<Button onClick={refresh}>Try again</Button>
					</Card>
				) : projects.length === 0 ? (
					<Card className="p-12 text-center max-w-lg mx-auto">
						<Package className="h-16 w-16 mx-auto mb-4 text-muted-foreground" />
						<h2 className="text-2xl font-bold mb-2">No projects found</h2>
						<p className="text-muted-foreground mb-6">Add a project to get started</p>
						<div className="flex gap-3 justify-center">
							<Button onClick={handleAddProject} size="lg">
								<Plus className="h-4 w-4 mr-2" />
								Add Local Project
							</Button>
							<Button onClick={refresh} variant="outline" size="lg">
								<RefreshCw className="h-4 w-4 mr-2" />
								Reload
							</Button>
						</div>
					</Card>
				) : (
					<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
						{projects.map((project) => (
							<ProjectCard
								key={project.id}
								project={project}
								onStart={() => startProject(project.id)}
								onStop={() => stopProject(project.id)}
								onRestart={() => restartProject(project.id)}
								onRemove={() => removeProject(project.id)}
								onViewLogs={() => setSelectedProjectId(project.id)}
							/>
						))}
					</div>
				)}
			</main>
			{selectedProject && (
				<LogsViewer
					projectId={selectedProject.id}
					projectName={selectedProject.name}
					onClose={() => setSelectedProjectId(null)}
				/>
			)}
			<ConfigEditor open={configEditorOpen} onOpenChange={setConfigEditorOpen} />
		</div>
	);
}

export default App;
