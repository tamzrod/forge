import type { Property, Quality } from '../../types/inspector';
import styles from '../Inspector.module.css';

// Format a numeric value with optional precision
function formatNumber(value: number, precision?: number, unit?: string): string {
  const formatted = precision !== undefined 
    ? value.toFixed(precision) 
    : (Number.isInteger(value) ? value.toString() : value.toFixed(2));
  return unit ? `${formatted} ${unit}` : formatted;
}

// Format a timestamp (unix timestamp in seconds)
function formatTimestamp(value: number): string {
  const date = new Date(value * 1000);
  return date.toLocaleString();
}

// Format a duration (nanoseconds)
function formatDuration(value: number): string {
  const ns = BigInt(value);
  const seconds = Number(ns / BigInt(1e9));
  const minutes = Math.floor(seconds / 60);
  const hours = Math.floor(minutes / 60);
  const days = Math.floor(hours / 24);

  if (days > 0) {
    return `${days}d ${hours % 24}h ${minutes % 60}m`;
  }
  if (hours > 0) {
    return `${hours}h ${minutes % 60}m ${Math.floor(seconds % 60)}s`;
  }
  if (minutes > 0) {
    return `${minutes}m ${Math.floor(seconds % 60)}s`;
  }
  return `${seconds.toFixed(3)}s`;
}

// Get status color class
function getStatusColorClass(status: string): string {
  const lowerStatus = status.toLowerCase();
  if (['running', 'healthy', 'connected', 'stable', 'ok', 'good', 'active'].includes(lowerStatus)) {
    return styles.statusHealthy;
  }
  if (['warning', 'transition', 'degraded', 'uncertain'].includes(lowerStatus)) {
    return styles.statusWarning;
  }
  if (['fault', 'error', 'critical', 'disconnected', 'bad', 'stopped', 'offline', 'disabled'].includes(lowerStatus)) {
    return styles.statusFault;
  }
  return styles.statusDefault;
}

// Quality display helper
function qualityToString(quality: Quality): string {
  switch (quality) {
    case 0: return 'Good';
    case 1: return 'Uncertain';
    case 2: return 'Bad';
    case 3: return 'Offline';
    default: return 'Unknown';
  }
}

// Boolean display helper
function formatBoolean(value: boolean): string {
  return value ? 'Yes' : 'No';
}

// Color based on value ranges
function getValueColor(name: string, value: number): string | undefined {
  const lowerName = name.toLowerCase();
  
  // Temperature coloring
  if (lowerName.includes('temperature')) {
    if (value < 0 || value > 40) return 'var(--color-fault)';
    if (value < 10 || value > 35) return 'var(--color-warning)';
    return 'var(--color-environmental)';
  }
  
  // Voltage coloring (per-unit)
  if (lowerName.includes('voltage') && lowerName.includes('pu')) {
    if (value >= 0.95 && value <= 1.05) return 'var(--color-healthy)';
    if (value >= 0.9 && value <= 1.1) return 'var(--color-warning)';
    return 'var(--color-fault)';
  }
  
  // Frequency coloring (per-unit)
  if (lowerName.includes('frequency') && lowerName.includes('pu')) {
    if (value >= 0.99 && value <= 1.01) return 'var(--color-healthy)';
    if (value >= 0.98 && value <= 1.02) return 'var(--color-warning)';
    return 'var(--color-fault)';
  }

  // Pressure coloring
  if (lowerName.includes('pressure')) {
    if (value >= 980 && value <= 1050) return 'var(--color-engineering)';
    return 'var(--color-warning)';
  }
  
  return undefined;
}

// PropertyValue component
interface PropertyValueProps {
  property: Property;
}

export function PropertyValue({ property }: PropertyValueProps) {
  const { name, value, type, unit, quality, precision } = property;

  const renderValue = () => {
    if (value === null || value === undefined) {
      return <span className={styles.propertyValue}>—</span>;
    }

    switch (type) {
      case 'text':
        return <span className={styles.propertyValue}>{String(value)}</span>;

      case 'number': {
        const numValue = typeof value === 'number' ? value : parseFloat(String(value));
        const color = getValueColor(name, numValue);
        return (
          <span className={styles.propertyValue} style={color ? { color } : undefined}>
            {formatNumber(numValue, precision, unit)}
          </span>
        );
      }

      case 'angle': {
        const numValue = typeof value === 'number' ? value : parseFloat(String(value));
        return (
          <span className={styles.propertyValue} style={{ color: 'var(--color-environmental)' }}>
            {formatNumber(numValue, precision || 1, unit || '°')}
          </span>
        );
      }

      case 'percentage': {
        const numValue = typeof value === 'number' ? value : parseFloat(String(value));
        return (
          <span className={styles.propertyValue}>
            {formatNumber(numValue, precision, unit || '%')}
          </span>
        );
      }

      case 'boolean':
        return <span className={styles.propertyValue}>{formatBoolean(Boolean(value))}</span>;

      case 'status': {
        const statusStr = String(value);
        return (
          <span className={`${styles.statusBadge} ${getStatusColorClass(statusStr)}`}>
            <span className={styles.statusDot} />
            {statusStr}
          </span>
        );
      }

      case 'timestamp':
        return <span className={styles.propertyValue}>{formatTimestamp(Number(value))}</span>;

      case 'duration':
        return <span className={styles.propertyValue}>{formatDuration(Number(value))}</span>;

      case 'quality':
        return (
          <span className={`${styles.statusBadge} ${quality === 0 ? styles.statusHealthy : styles.statusWarning}`}>
            <span className={styles.statusDot} />
            {qualityToString(value as Quality)}
          </span>
        );

      case 'enum': {
        const enumStr = String(value);
        return <span className={styles.propertyValue}>{enumStr}</span>;
      }

      case 'nested':
        return <span className={styles.propertyValue}>Object ({property.children?.length || 0} properties)</span>;

      case 'list':
        return <span className={styles.propertyValue}>List ({property.items?.length || 0} items)</span>;

      default:
        return <span className={styles.propertyValue}>{String(value)}</span>;
    }
  };

  return renderValue();
}
