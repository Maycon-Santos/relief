import { useState, useEffect, useRef } from 'react';
import { api } from '../services/wails';
import type { LogEntry } from '../types/project';

interface LogsViewerProps {
  projectId: string;
  projectName: string;
  onClose: () => void;
}

export function LogsViewer({ projectId, projectName, onClose }: LogsViewerProps) {
  const [logs, setLogs] = useState<LogEntry[]>([]);
  const [autoScroll, setAutoScroll] = useState(true);
  const logsEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const loadLogs = async () => {
      try {
        const data = await api.getProjectLogs(projectId, 500);
        setLogs(data);
      } catch (err) {
        console.error('Erro ao carregar logs:', err);
      }
    };

    loadLogs();
    const interval = setInterval(loadLogs, 2000);

    return () => clearInterval(interval);
  }, [projectId]);

  useEffect(() => {
    if (autoScroll) {
      logsEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    }
  }, [logs, autoScroll]);

  const getLevelColor = (level: string) => {
    switch (level.toLowerCase()) {
      case 'error':
        return '#ef4444';
      case 'warn':
        return '#f59e0b';
      case 'info':
        return '#3b82f6';
      case 'debug':
        return '#6b7280';
      default:
        return '#9ca3af';
    }
  };

  return (
    <div
      style={{
        position: 'fixed',
        top: '0',
        left: '0',
        right: '0',
        bottom: '0',
        backgroundColor: 'rgba(0, 0, 0, 0.8)',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        zIndex: 1000,
      }}
      onClick={onClose}
    >
      <div
        style={{
          backgroundColor: '#111827',
          borderRadius: '12px',
          width: '90%',
          maxWidth: '1200px',
          height: '80%',
          display: 'flex',
          flexDirection: 'column',
          border: '1px solid #374151',
        }}
        onClick={(e) => e.stopPropagation()}
      >
        <div
          style={{
            padding: '20px',
            borderBottom: '1px solid #374151',
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
          }}
        >
          <h2 style={{ margin: '0', fontSize: '20px', color: '#f9fafb' }}>Logs: {projectName}</h2>
          <div style={{ display: 'flex', gap: '12px', alignItems: 'center' }}>
            <label style={{ display: 'flex', alignItems: 'center', gap: '8px', fontSize: '14px', color: '#d1d5db' }}>
              <input
                type="checkbox"
                checked={autoScroll}
                onChange={(e) => setAutoScroll(e.target.checked)}
              />
              Auto-scroll
            </label>
            <button
              onClick={onClose}
              style={{
                padding: '8px 16px',
                backgroundColor: '#374151',
                color: 'white',
                border: 'none',
                borderRadius: '8px',
                fontSize: '14px',
                cursor: 'pointer',
              }}
            >
              Fechar
            </button>
          </div>
        </div>

        <div
          style={{
            flex: 1,
            overflowY: 'auto',
            padding: '20px',
            fontFamily: 'monospace',
            fontSize: '13px',
            backgroundColor: '#0f172a',
          }}
        >
          {logs.length === 0 ? (
            <p style={{ color: '#6b7280' }}>Nenhum log disponível</p>
          ) : (
            logs.map((log) => (
              <div
                key={log.id}
                style={{
                  marginBottom: '8px',
                  display: 'flex',
                  gap: '12px',
                  color: '#d1d5db',
                }}
              >
                <span style={{ color: '#6b7280', flexShrink: 0 }}>
                  {new Date(log.timestamp).toLocaleTimeString()}
                </span>
                <span
                  style={{
                    color: getLevelColor(log.level),
                    fontWeight: '600',
                    flexShrink: 0,
                    width: '60px',
                  }}
                >
                  [{log.level.toUpperCase()}]
                </span>
                <span style={{ flex: 1, wordBreak: 'break-word' }}>{log.message}</span>
              </div>
            ))
          )}
          <div ref={logsEndRef} />
        </div>
      </div>
    </div>
  );
}
