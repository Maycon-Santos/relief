import { AlertTriangle } from "lucide-react";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import type { Dependency } from "../types/project";

interface DependencyAlertProps {
	dependencies: Dependency[];
}

export function DependencyAlert({ dependencies }: DependencyAlertProps) {
	if (dependencies.length === 0) return null;

	return (
		<Alert variant="destructive" className="mb-3">
			<AlertTriangle className="h-4 w-4" />
			<AlertTitle>Unsatisfied dependencies</AlertTitle>
			<AlertDescription>
				<ul className="mt-2 space-y-1 list-disc list-inside">
					{dependencies.map((dep) => (
						<li key={`${dep.name}-${dep.required_version}`} className="text-xs">
							<strong>{dep.name}</strong>: {dep.message || `Required: ${dep.required_version}`}
						</li>
					))}
				</ul>
			</AlertDescription>
		</Alert>
	);
}
