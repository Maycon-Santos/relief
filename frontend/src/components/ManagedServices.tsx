import { Database, PlayCircle, Server, StopCircle } from "lucide-react";
import { useState } from "react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader } from "@/components/ui/card";

interface ManagedService {
	name: string;
	running: boolean;
}

interface ManagedServicesProps {
	services: ManagedService[];
	onStartService: (serviceName: string) => Promise<void>;
	onStopService: (serviceName: string) => Promise<void>;
}

export function ManagedServices({ services, onStartService, onStopService }: ManagedServicesProps) {
	const [loadingServices, setLoadingServices] = useState<Set<string>>(new Set());

	const handleServiceAction = async (
		serviceName: string,
		action: (name: string) => Promise<void>,
	) => {
		setLoadingServices((prev) => new Set(prev).add(serviceName));
		try {
			await action(serviceName);
		} catch (err) {
			console.error(`Error with service ${serviceName}:`, err);
			const message = err instanceof Error ? err.message : "Unknown error";
			alert(`Failed to manage service:\n\n${message}`);
		} finally {
			setLoadingServices((prev) => {
				const next = new Set(prev);
				next.delete(serviceName);
				return next;
			});
		}
	};

	const getServiceIcon = (name: string) => {
		const lowerName = name.toLowerCase();
		if (lowerName.includes("postgres") || lowerName.includes("mongo")) {
			return <Database className="h-5 w-5" />;
		}
		return <Server className="h-5 w-5" />;
	};

	if (services.length === 0) {
		return null;
	}

	return (
		<Card className="border-border/50 bg-zinc-900/50">
			<CardHeader className="pb-4">
				<div className="flex items-center gap-2 text-gray-300">
					<Server className="h-5 w-5" />
					<h2 className="text-lg font-semibold">Managed Services</h2>
				</div>
			</CardHeader>
			<CardContent>
				<div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
					{services.map((service) => (
						<div
							key={service.name}
							className="flex items-center justify-between p-3 bg-zinc-800/50 rounded-lg border border-zinc-700/50"
						>
							<div className="flex items-center gap-2 flex-1 min-w-0">
								{getServiceIcon(service.name)}
								<div className="flex flex-col min-w-0">
									<span className="font-medium text-sm text-white capitalize truncate">
										{service.name}
									</span>
									<Badge
										variant={service.running ? "default" : "secondary"}
										className={`text-xs w-fit ${
											service.running
												? "bg-green-500/20 text-green-400 border-green-500/30"
												: "bg-zinc-800/80 text-gray-400 border-zinc-700/50"
										}`}
									>
										{service.running ? "Running" : "Stopped"}
									</Badge>
								</div>
							</div>
							<div className="flex gap-1">
								{service.running ? (
									<Button
										size="sm"
										variant="ghost"
										onClick={() => handleServiceAction(service.name, onStopService)}
										disabled={loadingServices.has(service.name)}
										className="h-8 w-8 p-0 hover:bg-zinc-700"
									>
										<StopCircle className="h-4 w-4 text-red-400" />
									</Button>
								) : (
									<Button
										size="sm"
										variant="ghost"
										onClick={() => handleServiceAction(service.name, onStartService)}
										disabled={loadingServices.has(service.name)}
										className="h-8 w-8 p-0 hover:bg-zinc-700"
									>
										<PlayCircle className="h-4 w-4 text-green-400" />
									</Button>
								)}
							</div>
						</div>
					))}
				</div>
			</CardContent>
		</Card>
	);
}
