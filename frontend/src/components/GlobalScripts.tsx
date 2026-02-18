import { Code2, Play } from "lucide-react";
import { useCallback, useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { api } from "@/services/wails";

export function GlobalScripts() {
	const [globalScripts, setGlobalScripts] = useState<Record<string, string>>({});
	const [executingScript, setExecutingScript] = useState<string | null>(null);
	const [error, setError] = useState<string | null>(null);
	const [loading, setLoading] = useState(true);

	const loadGlobalScripts = useCallback(async () => {
		try {
			setLoading(true);
			const scripts = await api.getGlobalScripts();
			setGlobalScripts(scripts || {});
		} catch (err) {
			console.error("Error loading global scripts:", err);
			setError(err instanceof Error ? err.message : "Failed to load global scripts");
		} finally {
			setLoading(false);
		}
	}, []);

	useEffect(() => {
		loadGlobalScripts();
	}, [loadGlobalScripts]);

	const handleExecuteScript = async (scriptName: string) => {
		try {
			setExecutingScript(scriptName);
			setError(null);
			await api.executeGlobalScript(scriptName);
			alert(`Script "${scriptName}" executed successfully!`);
		} catch (err) {
			const message = err instanceof Error ? err.message : "Error executing script";
			setError(`Failed to execute "${scriptName}": ${message}`);
			console.error("Error executing script:", err);
		} finally {
			setExecutingScript(null);
		}
	};

	const scripts = Object.entries(globalScripts);

	if (loading) {
		return null;
	}

	if (scripts.length === 0) {
		return null;
	}

	return (
		<div className="mb-8">
			<div className="mb-4 flex items-center gap-2 text-gray-500">
				<Code2 className="h-5 w-5" />
				<h2 className="text-lg font-semibold">Global Scripts</h2>
			</div>

			{error && (
				<div className="mb-4 bg-red-500/10 border border-red-500/30 text-red-400 px-4 py-3 rounded-lg text-sm">
					{error}
				</div>
			)}

			<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
				{scripts.map(([name, command]) => (
					<Card
						key={name}
						className="border border-zinc-800/50 bg-zinc-950/30 hover:bg-zinc-950/50 transition-colors"
					>
						<div className="p-5">
							<div className="mb-3">
								<h3 className="font-semibold text-white text-base mb-2">{name}</h3>
								<p className="text-xs text-gray-400 font-mono leading-relaxed wrap-break-word">
									{command.length > 100 ? `${command.substring(0, 100)}...` : command}
								</p>
							</div>
							<Button
								size="sm"
								onClick={() => handleExecuteScript(name)}
								disabled={executingScript !== null}
								className="w-full h-9"
								variant={executingScript === name ? "outline" : "default"}
							>
								{executingScript === name ? (
									<>Running...</>
								) : (
									<>
										<Play className="h-3.5 w-3.5 mr-2" />
										Run
									</>
								)}
							</Button>
						</div>
					</Card>
				))}
			</div>
		</div>
	);
}
