import { Check, ExternalLink, Save, X } from "lucide-react";
import { useCallback, useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog";
import { api } from "@/services/wails";

const EDITOR_OPTIONS = [
	{ value: "code", label: "VS Code" },
	{ value: "cursor", label: "Cursor" },
	{ value: "zed", label: "Zed" },
	{ value: "idea", label: "IntelliJ IDEA" },
	{ value: "webstorm", label: "WebStorm" },
	{ value: "sublime", label: "Sublime Text" },
	{ value: "atom", label: "Atom" },
	{ value: "nvim", label: "Neovim" },
];

interface ConfigEditorProps {
	open: boolean;
	onOpenChange: (open: boolean) => void;
}

export function ConfigEditor({ open, onOpenChange }: ConfigEditorProps) {
	const [configYAML, setConfigYAML] = useState("");
	const [originalYAML, setOriginalYAML] = useState("");
	const [loading, setLoading] = useState(true);
	const [saving, setSaving] = useState(false);
	const [error, setError] = useState<string | null>(null);
	const [success, setSuccess] = useState(false);
	const [editor, setEditor] = useState("code");

	const loadConfig = useCallback(async () => {
		try {
			setLoading(true);
			setError(null);
			const yaml = await api.getConfigYAML();
			setConfigYAML(yaml);
			setOriginalYAML(yaml);
			const currentEditor = await api.getEditorConfig();
			setEditor(currentEditor);
		} catch (err) {
			setError(err instanceof Error ? err.message : "Error loading config");
			console.error("Error loading config:", err);
		} finally {
			setLoading(false);
		}
	}, []);

	useEffect(() => {
		if (open) {
			loadConfig();
		}
	}, [open, loadConfig]);

	const handleSave = async () => {
		try {
			setSaving(true);
			setError(null);
			setSuccess(false);
			await api.saveGlobalConfig(configYAML);
			setOriginalYAML(configYAML);
			setSuccess(true);
			setTimeout(() => setSuccess(false), 3000);
		} catch (err) {
			setError(err instanceof Error ? err.message : "Error saving config");
			console.error("Error saving config:", err);
		} finally {
			setSaving(false);
		}
	};

	const handleOpenInEditor = async () => {
		try {
			await api.openConfigInEditor();
		} catch (err) {
			console.error("Error opening config in editor:", err);
			alert("Could not open config in editor");
		}
	};

	const hasChanges = configYAML !== originalYAML;

	return (
		<Dialog open={open} onOpenChange={onOpenChange}>
			<DialogContent className="max-w-4xl h-[85vh] flex flex-col gap-0 p-0 bg-zinc-900 border-zinc-800">
				<DialogHeader className="px-6 pt-6 pb-4">
					<DialogTitle className="flex items-center gap-2">
						Global Configuration
						<Button
							variant="ghost"
							size="sm"
							onClick={handleOpenInEditor}
							className="h-6 px-2 text-xs"
						>
							<ExternalLink className="h-3 w-3 mr-1" />
							Open in Editor
						</Button>
					</DialogTitle>
					<DialogDescription>
						Edit the global configuration file (~/.relief/config.yaml). Paths use ~/ for home
						directory. Changes are applied immediately after saving.
					</DialogDescription>
					<div className="flex items-center gap-3 pt-2">
						<label className="text-sm text-gray-400">Preferred Editor:</label>
						<select
							value={editor}
							onChange={async (e) => {
								const newEditor = e.target.value;
								setEditor(newEditor);
								try {
									await api.setEditorConfig(newEditor);
								} catch (err) {
									console.error("Error saving editor config:", err);
								}
							}}
							className="bg-zinc-800 border border-zinc-700 text-gray-200 text-sm rounded px-2 py-1 focus:outline-none focus:ring-2 focus:ring-blue-500"
						>
							{EDITOR_OPTIONS.map((opt) => (
								<option key={opt.value} value={opt.value}>
									{opt.label}
								</option>
							))}
						</select>
					</div>
				</DialogHeader>

				<div className="flex-1 flex flex-col gap-3 overflow-hidden px-6 min-h-0">
					{error && (
						<div className="bg-red-500/10 border border-red-500/30 text-red-400 px-4 py-2 rounded-lg text-sm shrink-0">
							{error}
						</div>
					)}

					{success && (
						<div className="bg-green-500/10 border border-green-500/30 text-green-400 px-4 py-2 rounded-lg text-sm flex items-center gap-2 shrink-0">
							<Check className="h-4 w-4" />
							Configuration saved successfully!
						</div>
					)}

					{loading ? (
						<div className="flex-1 flex items-center justify-center">
							<div className="text-gray-400">Loading configuration...</div>
						</div>
					) : (
						<textarea
							value={configYAML}
							onChange={(e) => setConfigYAML(e.target.value)}
							className="flex-1 w-full p-4 bg-zinc-950 border border-zinc-700 rounded-lg text-sm font-mono text-gray-200 focus:outline-none focus:ring-2 focus:ring-blue-500 resize-none min-h-0"
							spellCheck={false}
						/>
					)}
				</div>

				<DialogFooter className="px-6 py-4 border-t border-zinc-800 shrink-0">
					<Button
						variant="outline"
						onClick={() => {
							setConfigYAML(originalYAML);
							setError(null);
							setSuccess(false);
						}}
						disabled={!hasChanges || saving}
						className="border-zinc-700"
					>
						<X className="h-4 w-4 mr-2" />
						Reset
					</Button>
					<Button onClick={handleSave} disabled={!hasChanges || saving}>
						<Save className="h-4 w-4 mr-2" />
						{saving ? "Saving..." : "Save Changes"}
					</Button>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	);
}
