import type { CanvasEntity } from '../../types/editor';
import type { State } from '../../types';
import styles from './EquipmentDetails.module.css';

type EquipmentTab = 'explain' | 'measurements' | 'identity' | 'status' | 'properties';

interface EquipmentDetailsProps {
  entity: CanvasEntity | null | undefined;
  simulationState: State;
  activeTab: EquipmentTab;
  onPropertyChange: (entityId: string, key: string, value: unknown) => void;
}

const ENTITY_ICONS: Record<string, string> = {
  grid: '🔌',
  bus: '⚫',
  breaker: '🔀',
  transformer: '🔄',
  generator: '☀️',
  load: '🏭',
  meter: '📊',
};

const ENTITY_NAMES: Record<string, string> = {
  grid: 'Utility Grid',
  bus: 'Electrical Bus',
  breaker: 'Circuit Breaker',
  transformer: 'Power Transformer',
  generator: 'PV Array',
  load: 'Station Load',
  meter: 'Revenue Meter',
};

export function EquipmentDetails({
  entity,
  simulationState,
  activeTab,
  onPropertyChange,
}: EquipmentDetailsProps) {
  if (!entity) {
    return (
      <div className={styles.details}>
        <div className={styles.empty}>
          <div className={styles.emptyIcon}>📡</div>
          <div className={styles.emptyText}>
            Select equipment from the diagram or plant explorer to view details
          </div>
        </div>
      </div>
    );
  }

  const icon = ENTITY_ICONS[entity.entity_type] || '📦';
  const name = ENTITY_NAMES[entity.entity_type] || entity.entity_type;

  // Generate explain content based on entity type and simulation state
  const getExplainContent = () => {
    const { sun, grid } = simulationState;

    switch (entity.entity_type) {
      case 'generator':
        const irradiance = sun.irradiance;
        const capacity = parseFloat(entity.properties.rated_capacity?.value as string) || 500;
        const panelArea = 10000; // 10,000 m²
        const efficiency = 0.20;
        const theoreticalPower = irradiance / 1000 * panelArea * efficiency;
        const actualPower = theoreticalPower * (capacity / 1000);
        const utilization = (irradiance / 1000) * 100;

        return {
          title: 'PV Array Power Generation',
          content: `This PV array is currently generating ${actualPower.toFixed(1)} kW.
          
The power output depends on three factors:
• Solar irradiance: ${irradiance.toFixed(0)} W/m² (${utilization.toFixed(0)}% of peak)
• Array capacity: ${capacity} kW
• Panel efficiency: 20%

The theoretical maximum power at current irradiance is ${theoreticalPower.toFixed(0)} kW. This array can produce up to ${capacity} kW under standard test conditions (1000 W/m²).

${sun.is_daytime ? 'Sun is above the horizon' : 'Sun is below the horizon - no generation'}`,
          highlight: `${actualPower.toFixed(1)} kW`,
        };

      case 'meter':
        return {
          title: 'Grid Interconnection Point',
          content: `This revenue meter monitors power flow at the grid interconnection point.

Current readings:
• Voltage: ${grid.voltage.toFixed(0)} V (${(grid.voltage_pu * 100).toFixed(1)}% of nominal)
• Frequency: ${grid.frequency.toFixed(2)} Hz
• Grid status: ${grid.is_stable ? 'Stable' : 'Unstable'}

${grid.is_stable ? 'The grid is operating within normal parameters.' : 'Grid frequency deviation detected.'}`,
          highlight: `${grid.voltage.toFixed(0)} V`,
        };

      case 'transformer':
        return {
          title: 'Power Transformer',
          content: `This transformer steps voltage between the PV array and the utility grid.

The transformer's tap position affects the secondary voltage. Current grid voltage is ${grid.voltage.toFixed(0)} V at the low voltage side.

Typical transformer losses: 1-2%
Operating temperature depends on load current and ambient conditions.`,
          highlight: `${grid.voltage.toFixed(0)} V`,
        };

      case 'bus':
        return {
          title: 'Electrical Bus',
          content: `This bus distributes power within the solar farm collection system.

Buses aggregate power from multiple PV arrays and route it through transformers to the grid connection point.

Current bus voltage: ${grid.voltage.toFixed(0)} V`,
          highlight: `${grid.voltage.toFixed(0)} V`,
        };

      case 'grid':
        return {
          title: 'Utility Grid Connection',
          content: `The utility grid provides the point of common coupling (PCC) for this solar farm.

Grid parameters:
• Nominal voltage: ${grid.nominal_voltage} V
• Nominal frequency: ${grid.nominal_frequency} Hz
• Current voltage: ${grid.voltage.toFixed(0)} V
• Current frequency: ${grid.frequency.toFixed(2)} Hz

The grid must remain stable for this plant to export power.`,
          highlight: `${grid.frequency.toFixed(2)} Hz`,
        };

      default:
        return {
          title: 'Equipment Status',
          content: 'Select different equipment to see detailed explanations.',
          highlight: '',
        };
    }
  };

  // Get measurements based on entity type
  const getMeasurements = () => {
    const { sun, weather, grid } = simulationState;

    switch (entity.entity_type) {
      case 'generator':
        const irradiance = sun.irradiance;
        const capacity = parseFloat(entity.properties.rated_capacity?.value as string) || 500;
        const panelArea = 10000;
        const efficiency = 0.20;
        const power = irradiance / 1000 * panelArea * efficiency * (capacity / 1000);
        const dcVoltage = 400;
        const dcCurrent = power / dcVoltage;
        const acPower = power * 0.98; // Inverter losses

        return [
          { name: 'Active Power', value: acPower.toFixed(1), unit: 'kW', max: capacity },
          { name: 'DC Power', value: power.toFixed(1), unit: 'kW', max: capacity },
          { name: 'Irradiance', value: irradiance.toFixed(0), unit: 'W/m²', max: 1000 },
          { name: 'DC Voltage', value: dcVoltage.toFixed(0), unit: 'V', max: 600 },
          { name: 'DC Current', value: dcCurrent.toFixed(1), unit: 'A', max: 2000 },
          { name: 'Temperature', value: weather.temperature.toFixed(1), unit: '°C', max: 60 },
        ];

      case 'meter':
        return [
          { name: 'Voltage', value: grid.voltage.toFixed(0), unit: 'V', max: 600, status: grid.voltage_pu < 0.95 ? 'warning' : 'normal' },
          { name: 'Frequency', value: grid.frequency.toFixed(2), unit: 'Hz', max: 62, status: !grid.is_stable ? 'error' : 'normal' },
          { name: 'Active Power', value: (simulationState.sun.irradiance * 10).toFixed(0), unit: 'kW', max: 1000 },
          { name: 'Power Factor', value: '0.99', unit: '', max: 1 },
        ];

      case 'transformer':
        return [
          { name: 'Primary Voltage', value: grid.voltage.toFixed(0), unit: 'V', max: 80000 },
          { name: 'Secondary Voltage', value: grid.voltage.toFixed(0), unit: 'V', max: 600 },
          { name: 'Load', value: '45', unit: '%', max: 100 },
          { name: 'Temperature', value: '45', unit: '°C', max: 120 },
        ];

      default:
        return [];
    }
  };

  // Render based on active tab
  switch (activeTab) {
    case 'explain':
      const explain = getExplainContent();
      return (
        <div className={styles.details}>
          <div className={styles.header}>
            <div className={styles.headerIcon}>{icon}</div>
            <div className={styles.headerTitle}>{entity.name}</div>
            <div className={styles.headerSubtitle}>{name}</div>
          </div>

          <div className={styles.explainCard}>
            <div className={styles.explainHeader}>
              <span className={styles.explainIcon}>💡</span>
              <span className={styles.explainTitle}>{explain.title}</span>
            </div>
            <div className={styles.explainContent}>
              {explain.content.split('\n\n').map((paragraph, i) => (
                <p key={i} style={{ marginBottom: paragraph.includes('\n') ? '0' : '12px' }}>
                  {paragraph.split('\n').map((line, j) => (
                    <span key={j}>
                      {line}
                      {j < paragraph.split('\n').length - 1 && <br />}
                    </span>
                  ))}
                </p>
              ))}
            </div>
          </div>

          {explain.highlight && (
            <div className={styles.whyBox}>
              <div className={styles.whyHeader}>
                <span>🎯</span>
                <span>Current Value</span>
              </div>
              <div className={styles.whyContent}>
                The <span className={styles.whyHighlight}>{explain.highlight}</span> reading is directly 
                influenced by the current environmental conditions and equipment configuration.
              </div>
            </div>
          )}
        </div>
      );

    case 'measurements':
      const measurements = getMeasurements();
      return (
        <div className={styles.details}>
          <div className={styles.header}>
            <div className={styles.headerIcon}>{icon}</div>
            <div className={styles.headerTitle}>{entity.name}</div>
            <div className={styles.headerSubtitle}>Live Measurements</div>
          </div>

          {measurements.map((m, i) => {
            const percentage = Math.min((parseFloat(m.value) / m.max) * 100, 100);
            const statusClass = m.status === 'warning' ? styles.warning : m.status === 'error' ? styles.error : '';
            const valueClass = m.status === 'warning' ? styles.warning : m.status === 'error' ? styles.error : '';

            return (
              <div key={i} className={`${styles.measurementCard} ${statusClass}`}>
                <div className={styles.measurementHeader}>
                  <span className={styles.measurementName}>{m.name}</span>
                  <span>
                    <span className={`${styles.measurementValue} ${valueClass}`}>{m.value}</span>
                    {m.unit && <span className={styles.measurementUnit}>{m.unit}</span>}
                  </span>
                </div>
                <div className={styles.measurementBar}>
                  <div
                    className={`${styles.measurementBarFill} ${statusClass}`}
                    style={{ width: `${percentage}%` }}
                  />
                </div>
              </div>
            );
          })}
        </div>
      );

    case 'identity':
      return (
        <div className={styles.details}>
          <div className={styles.header}>
            <div className={styles.headerIcon}>{icon}</div>
            <div className={styles.headerTitle}>{entity.name}</div>
            <div className={styles.headerSubtitle}>{name}</div>
          </div>

          <div className={styles.section}>
            <div className={styles.sectionTitle}>Identification</div>
            <div className={styles.listItem}>
              <span className={styles.listLabel}>Equipment ID</span>
              <span className={styles.listValue}>{entity.id}</span>
            </div>
            <div className={styles.listItem}>
              <span className={styles.listLabel}>Type</span>
              <span className={styles.listValue}>{entity.entity_type}</span>
            </div>
            <div className={styles.listItem}>
              <span className={styles.listLabel}>Component</span>
              <span className={styles.listValue}>{entity.component_id || 'N/A'}</span>
            </div>
          </div>

          <div className={styles.section}>
            <div className={styles.sectionTitle}>Location</div>
            <div className={styles.listItem}>
              <span className={styles.listLabel}>Position X</span>
              <span className={styles.listValue}>{entity.position.x.toFixed(0)}</span>
            </div>
            <div className={styles.listItem}>
              <span className={styles.listLabel}>Position Y</span>
              <span className={styles.listValue}>{entity.position.y.toFixed(0)}</span>
            </div>
          </div>
        </div>
      );

    case 'status':
      const isOnline = entity.entity_type === 'generator' ? simulationState.sun.is_daytime : true;
      return (
        <div className={styles.details}>
          <div className={styles.header}>
            <div className={styles.headerIcon}>{icon}</div>
            <div className={styles.headerTitle}>{entity.name}</div>
            <div className={styles.headerSubtitle}>Equipment Status</div>
          </div>

          <div className={styles.section}>
            <div className={styles.sectionTitle}>Operational Status</div>
            <div className={styles.statusIndicator}>
              <span className={`${styles.statusDot} ${isOnline ? styles.online : styles.offline}`} />
              <span className={styles.statusText}>
                {isOnline ? 'Online' : 'Offline'}
              </span>
            </div>
          </div>

          <div className={styles.section}>
            <div className={styles.sectionTitle}>Communication</div>
            <div className={styles.statusIndicator}>
              <span className={`${styles.statusDot} ${styles.online}`} />
              <span className={styles.statusText}>Connected to SCADA</span>
            </div>
          </div>

          <div className={styles.section}>
            <div className={styles.sectionTitle}>Alarms</div>
            <div className={styles.statusIndicator}>
              <span className={`${styles.statusDot} ${styles.online}`} />
              <span className={styles.statusText}>No active alarms</span>
            </div>
          </div>
        </div>
      );

    case 'properties':
      return (
        <div className={styles.details}>
          <div className={styles.header}>
            <div className={styles.headerIcon}>{icon}</div>
            <div className={styles.headerTitle}>{entity.name}</div>
            <div className={styles.headerSubtitle}>Configuration</div>
          </div>

          {Object.entries(entity.properties).map(([key, prop]) => (
            <div key={key} className={styles.inputGroup}>
              <label className={styles.inputLabel}>
                {key.replace(/_/g, ' ').replace(/\b\w/g, (l) => l.toUpperCase())}
              </label>
              <div className={styles.inputFieldWithUnit}>
                {prop.type === 'number' ? (
                  <input
                    type="number"
                    className={styles.inputField}
                    value={prop.value as number}
                    onChange={(e) => onPropertyChange(entity.id, key, parseFloat(e.target.value))}
                    readOnly={prop.readonly}
                    min={prop.min}
                    max={prop.max}
                  />
                ) : prop.type === 'boolean' ? (
                  <input
                    type="checkbox"
                    checked={prop.value as boolean}
                    onChange={(e) => onPropertyChange(entity.id, key, e.target.checked)}
                    disabled={prop.readonly}
                  />
                ) : (
                  <input
                    type="text"
                    className={styles.inputField}
                    value={prop.value as string}
                    onChange={(e) => onPropertyChange(entity.id, key, e.target.value)}
                    readOnly={prop.readonly}
                  />
                )}
                {prop.unit && <span className={styles.inputUnit}>{prop.unit}</span>}
              </div>
            </div>
          ))}
        </div>
      );

    default:
      return null;
  }
}
