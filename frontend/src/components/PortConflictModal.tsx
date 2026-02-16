import { AlertTriangle, Hash, Terminal } from "lucide-react";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog";
import type { PortConflict } from "../services/wails";

interface PortConflictModalProps {
	conflict: PortConflict;
	onKill: () => Promise<void>;
	onCancel: () => void;
}

export function PortConflictModal({ conflict, onKill, onCancel }: PortConflictModalProps) {
	const [killing, setKilling] = useState(false);

	const handleKill = async () => {
		setKilling(true);
		try {
			await onKill();
		} catch (err) {
			console.error("Error killing process:", err);
			alert(`Failed to kill process: ${err instanceof Error ? err.message : "Unknown error"}`);
		} finally {
			setKilling(false);
		}
	};

	return (
		<Dialog open onOpenChange={onCancel}>
			<DialogContent className="sm:max-w-md">
				<DialogHeader>
					<DialogTitle className="flex items-center gap-2">
						<AlertTriangle className="h-5 w-5 text-amber-500" />
						Port Already in Use
					</DialogTitle>
					<DialogDescription>
						Another application is already using this port. You can force stop the conflicting
						process.
					</DialogDescription>
				</DialogHeader>

				<div className="space-y-3 py-4">
					<div className="flex items-center gap-3 rounded-lg border border-border bg-muted/50 p-3">
						<Hash className="h-4 w-4 text-muted-foreground" />
						<div className="flex-1">
							<div className="text-xs text-muted-foreground">Port</div>
							<div className="font-mono font-semibold">{conflict.port}</div>
						</div>
					</div>

					<div className="flex items-center gap-3 rounded-lg border border-border bg-muted/50 p-3">
						<Terminal className="h-4 w-4 text-muted-foreground" />
						<div className="flex-1">
							<div className="text-xs text-muted-foreground">Process ID</div>
							<div className="font-mono font-semibold">{conflict.pid}</div>
						</div>
					</div>

					{conflict.command && (
						<div className="rounded-lg border border-border bg-muted/50 p-3">
							<div className="text-xs text-muted-foreground mb-1">Command</div>
							<code className="text-xs font-mono break-all">{conflict.command}</code>
						</div>
					)}
				</div>

				<DialogFooter className="gap-2 sm:gap-0">
					<Button variant="outline" onClick={onCancel} disabled={killing}>
						Cancel
					</Button>
					<Button variant="destructive" onClick={handleKill} disabled={killing}>
						{killing ? "Stopping..." : "Force Stop Process"}
					</Button>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	);
}
