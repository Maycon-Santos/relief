import type { Dependency } from '../types/project';

interface DependencyAlertProps {
  dependencies: Dependency[];
}

export function DependencyAlert({ dependencies }: DependencyAlertProps) {
  if (dependencies.length === 0) return null;

  return (
    <div
      style={{
        padding: '12px',
        backgroundColor: '#78350f',
        borderRadius: '8px',
        marginBottom: '12px',
        border: '1px solid #f59e0b',
      }}
    >
      <div style={{ display: 'flex', alignItems: 'center', marginBottom: '8px' }}>
        <span style={{ fontSize: '16px', marginRight: '8px' }}>⚠️</span>
        <strong style={{ fontSize: '14px', color: '#fbbf24' }}>Unsatisfied dependencies</strong>
      </div>
      <ul style={{ margin: '0', paddingLeft: '24px', color: '#fde68a' }}>
        {dependencies.map((dep, index) => (
          <li key={index} style={{ fontSize: '13px', marginBottom: '4px' }}>
            <strong>{dep.name}</strong>: {dep.message || `Required: ${dep.required_version}`}
          </li>
        ))}
      </ul>
    </div>
  );
}
