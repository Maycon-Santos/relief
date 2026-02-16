import { Folder, Loader2, Package, Plus, RefreshCw, Zap } from "lucide-react";
import { useEffect, useState } from "react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { LogsViewer } from "./components/LogsViewer";
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

	useEffect(() => {
		const loadStatus = async () => {
			try {
				const data = await api.getStatus();
				setStatus(data);
			} catch (err) {
				console.error("Error loading status:", err);
			}
		};

		loadStatus();
		const interval = setInterval(loadStatus, 5000);

		return () => clearInterval(interval);
	}, []);

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
								onClick={handleAddProject}
								disabled={loading}
								size="default"
								className="bg-white hover:bg-gray-200 text-black font-semibold transition-colors"
							>
								<Plus className="h-4 w-4 mr-2" />
								Add Project
							</Button>
							<Button
								onClick={refresh}
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
								</div>
							</div>
						</Card>
					)}
				</div>
			</header>

			<main className="container mx-auto max-w-7xl px-6 py-8">
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
		</div>
	);
}

export default App;
