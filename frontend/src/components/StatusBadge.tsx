import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";

interface StatusBadgeProps {
	status: string;
}

export function StatusBadge({ status }: StatusBadgeProps) {
	const getStatusConfig = () => {
		switch (status) {
			case "running":
				return {
					variant: "default" as const,
					className: "bg-emerald-500/15 text-emerald-400 border-emerald-500/30",
					label: "Rodando",
					pulse: true,
				};
			case "stopped":
				return {
					variant: "secondary" as const,
					className: "bg-muted text-muted-foreground border-border",
					label: "Parado",
					pulse: false,
				};
			case "starting":
				return {
					variant: "default" as const,
					className: "bg-amber-500/15 text-amber-400 border-amber-500/30",
					label: "Iniciando",
					pulse: true,
				};
			case "error":
				return {
					variant: "destructive" as const,
					className: "bg-destructive/15 text-destructive border-destructive/30",
					label: "Error",
					pulse: false,
				};
			default:
				return {
					variant: "secondary" as const,
					className: "bg-muted text-muted-foreground border-border",
					label: "Desconhecido",
					pulse: false,
				};
		}
	};

	const { className, label, pulse } = getStatusConfig();

	return (
		<Badge variant="outline" className={cn("flex items-center gap-1.5", className)}>
			<span className={cn("h-1.5 w-1.5 rounded-full bg-current", pulse && "animate-pulse")} />
			{label}
		</Badge>
	);
}
