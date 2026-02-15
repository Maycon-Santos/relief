import { useState } from "react";
import type { PortConflict } from "../services/wails";
import "./PortConflictModal.css";

interface PortConflictModalProps {
  conflict: PortConflict;
  onKill: () => Promise<void>;
  onCancel: () => void;
}

export function PortConflictModal({
  conflict,
  onKill,
  onCancel,
}: PortConflictModalProps) {
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
    <div className="modal-overlay">
      <div className="modal-content port-conflict-modal">
        <div className="modal-header">
          <h2>⚠️ Port Already in Use</h2>
        </div>

        <div className="modal-body">
          <div className="conflict-info">
            <div className="conflict-item">
              <strong>Port:</strong>
              <span className="conflict-port">{conflict.port}</span>
            </div>
            <div className="conflict-item">
              <strong>Process ID:</strong>
              <span className="conflict-pid">{conflict.pid}</span>
            </div>
            {conflict.command && (
              <div className="conflict-item">
                <strong>Command:</strong>
                <code className="conflict-command">{conflict.command}</code>
              </div>
            )}
          </div>

          <div className="conflict-message">
            <p>
              Another application is already using port <strong>{conflict.port}</strong>.
            </p>
            <p>
              You can force stop the conflicting process to free up the port.
            </p>
          </div>
        </div>

        <div className="modal-footer">
          <button
            type="button"
            onClick={onCancel}
            className="btn-secondary"
            disabled={killing}
          >
            Cancel
          </button>
          <button
            type="button"
            onClick={handleKill}
            className="btn-danger"
            disabled={killing}
          >
            {killing ? "Stopping..." : "Force Stop Process"}
          </button>
        </div>
      </div>
    </div>
  );
}
