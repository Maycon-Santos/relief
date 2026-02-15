import type { ProjectStatus } from '../types/project';

interface StatusBadgeProps {
  status: ProjectStatus;
}

export function StatusBadge({ status }: StatusBadgeProps) {
  const getStatusStyle = () => {
    switch (status) {
      case 'running':
        return { color: '#10b981', label: 'Rodando' };
      case 'stopped':
        return { color: '#6b7280', label: 'Parado' };
      case 'starting':
        return { color: '#f59e0b', label: 'Iniciando' };
      case 'error':
        return { color: '#ef4444', label: 'Erro' };
      default:
        return { color: '#9ca3af', label: 'Desconhecido' };
    }
  };

  const { color, label } = getStatusStyle();

  return (
    <span
      style={{
        display: 'inline-flex',
        alignItems: 'center',
        padding: '4px 12px',
        borderRadius: '12px',
        fontSize: '12px',
        fontWeight: '600',
        backgroundColor: `${color}20`,
        color: color,
      }}
    >
      <span
        style={{
          width: '6px',
          height: '6px',
          borderRadius: '50%',
          backgroundColor: color,
          marginRight: '6px',
          animation: status === 'running' ? 'pulse 2s infinite' : 'none',
        }}
      />
      {label}
    </span>
  );
}
