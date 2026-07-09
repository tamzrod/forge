import { useState } from 'react';
import {
  Activity,
  Settings,
  AlertTriangle,
  Radio,
  Eye,
  Wind,
  Sun,
  Zap,
  Clock,
  Thermometer,
  Gauge,
  Cpu
} from 'lucide-react';
import type { State, InspectorTab } from '../types';
import styles from './Inspector.module.css';

interface InspectorProps {
  state: State;
  selectedNode: string | null;
}

const tabs: { id: InspectorTab; label: string; icon: React.ReactNode }[] = [
  { id: 'overview', label: 'Overview', icon: <Eye size={14} /> },
  { id: 'state', label: 'State', icon: <Activity size={14} /> },
  { id: 'configuration', label: 'Configuration', icon: <Settings size={14} /> },
  { id: 'diagnostics', label: 'Diagnostics', icon: <AlertTriangle size={14} /> },
  { id: 'communications', label: 'Communications', icon: <Radio size={14} /> },
];

function StatusBadge({ status, size = 'default' }: { status: string; size?: 'small' | 'default' | 'large' }) {
  const getStatusClass = () => {
    switch (status.toLowerCase()) {
      case 'running':
      case 'healthy':
      case 'connected':
      case 'stable':
        return styles.statusHealthy;
      case 'warning':
      case 'transition':
        return styles.statusWarning;
      case 'fault':
      case 'error':
      case 'critical':
      case 'disconnected':
        return styles.statusFault;
      case 'offline':
      case 'stopped':
      case 'disabled':
        return styles.statusDisabled;
      default:
        return styles.statusDefault;
    }
  };

  return (
    <span className={`${styles.statusBadge} ${getStatusClass()} ${styles[`badge${size.charAt(0).toUpperCase() + size.slice(1)}`]}`}>
      <span className={styles.statusDot} />
      {status}
    </span>
  );
}

function PropertyRow({ label, value, valueColor }: { label: string; value: React.ReactNode; valueColor?: string }) {
  return (
    <div className={styles.propertyRow}>
      <span className={styles.propertyLabel}>{label}</span>
      <span className={styles.propertyValue} style={valueColor ? { color: valueColor } : undefined}>
        {value}
      </span>
    </div>
  );
}

function SectionCard({ title, icon, children }: { title: string; icon?: React.ReactNode; children: React.ReactNode }) {
  return (
    <div className={styles.sectionCard}>
      <div className={styles.sectionHeader}>
        {icon && <span className={styles.sectionIcon}>{icon}</span>}
        <span className={styles.sectionTitle}>{title}</span>
      </div>
      <div className={styles.sectionContent}>
        {children}
      </div>
    </div>
  );
}

function CompassWidget({ direction, label }: { direction: number; label: string }) {
  return (
    <div className={styles.compass}>
      <div className={styles.compassOuter}>
        <div 
          className={styles.compassNeedle} 
          style={{ transform: `rotate(${direction}deg)` }}
        />
        <div className={styles.compassN}>N</div>
        <div className={styles.compassE}>E</div>
        <div className={styles.compassS}>S</div>
        <div className={styles.compassW}>W</div>
      </div>
      <span className={styles.compassLabel}>{label}: {direction.toFixed(0)}°</span>
    </div>
  );
}

function OverviewContent({ nodeId, state }: { nodeId: string; state: State }) {
  if (nodeId === 'world') {
    return (
      <>
        <SectionCard title="Simulation Summary" icon={<Activity size={14} />}>
          <PropertyRow label="Status" value={<StatusBadge status={state.clock.is_paused ? 'Paused' : 'Running'} />} />
          <PropertyRow label="Elapsed" value={state.clock.elapsed} />
          <PropertyRow label="Tick Count" value={state.clock.tick_count.toLocaleString()} />
          <PropertyRow label="Devices" value={state.devices.count} />
        </SectionCard>
      </>
    );
  }

  if (nodeId === 'clock') {
    return (
      <SectionCard title="Clock Model" icon={<Clock size={14} />}>
        <PropertyRow label="Status" value={<StatusBadge status={state.clock.is_paused ? 'Paused' : 'Running'} />} />
        <PropertyRow label="Elapsed" value={state.clock.elapsed} />
        <PropertyRow label="Tick Count" value={state.clock.tick_count.toLocaleString()} />
        <PropertyRow label="Mode" value={state.clock.mode} />
      </SectionCard>
    );
  }

  if (nodeId === 'sun') {
    return (
      <>
        <SectionCard title="Sun Model" icon={<Sun size={14} />}>
          <PropertyRow label="Status" value={<StatusBadge status={state.sun.is_daytime ? 'Daytime' : 'Nighttime'} />} />
          <PropertyRow label="Elevation" value={`${state.sun.elevation.toFixed(1)}°`} valueColor="var(--color-environmental)" />
          <PropertyRow label="Azimuth" value={`${state.sun.azimuth.toFixed(1)}°`} />
          <PropertyRow label="Daytime" value={state.sun.is_daytime} />
        </SectionCard>
        <SectionCard title="Solar Irradiance" icon={<Sun size={14} />}>
          <PropertyRow label="GHI" value={`${state.sun.irradiance.toFixed(0)} W/m²`} valueColor="var(--color-environmental)" />
          <PropertyRow label="DNI" value={`${state.sun.direct_normal.toFixed(0)} W/m²`} valueColor="var(--color-environmental)" />
          <PropertyRow label="Diffuse" value={`${state.sun.diffuse.toFixed(0)} W/m²`} />
        </SectionCard>
      </>
    );
  }

  if (nodeId === 'weather') {
    return (
      <>
        <SectionCard title="Weather Model" icon={<Thermometer size={14} />}>
          <PropertyRow label="Temperature" value={`${state.weather.temperature.toFixed(1)} °C`} valueColor="var(--color-environmental)" />
          <PropertyRow label="Humidity" value={`${state.weather.humidity.toFixed(1)} %`} />
          <PropertyRow label="Pressure" value={`${state.weather.pressure.toFixed(1)} hPa`} valueColor="var(--color-engineering)" />
          <PropertyRow label="Cloud Cover" value={`${(state.weather.cloud_cover * 100).toFixed(0)} %`} />
          <PropertyRow label="Raining" value={state.weather.is_raining} />
        </SectionCard>
        <SectionCard title="Wind" icon={<Wind size={14} />}>
          <PropertyRow label="Speed" value={`${state.weather.wind_speed.toFixed(1)} m/s`} />
          <PropertyRow label="Direction" value={`${state.weather.wind_direction.toFixed(0)}°`} />
          <CompassWidget direction={state.weather.wind_direction} label="Wind" />
        </SectionCard>
      </>
    );
  }

  if (nodeId === 'grid') {
    const voltageClass = state.grid.voltage_pu >= 0.95 && state.grid.voltage_pu <= 1.05 
      ? 'var(--color-healthy)' 
      : state.grid.voltage_pu >= 0.9 && state.grid.voltage_pu <= 1.1 
        ? 'var(--color-warning)' 
        : 'var(--color-fault)';

    return (
      <>
        <SectionCard title="Grid Model" icon={<Zap size={14} />}>
          <PropertyRow label="Status" value={<StatusBadge status={state.grid.is_stable ? 'Stable' : 'Unstable'} />} />
          <PropertyRow label="Voltage" value={`${state.grid.voltage.toFixed(1)} V`} valueColor={voltageClass} />
          <PropertyRow label="Frequency" value={`${state.grid.frequency.toFixed(3)} Hz`} />
          <PropertyRow label="Voltage PU" value={state.grid.voltage_pu.toFixed(4)} valueColor={voltageClass} />
          <PropertyRow label="Frequency PU" value={state.grid.frequency_pu.toFixed(4)} />
        </SectionCard>
        <SectionCard title="Power Balance" icon={<Gauge size={14} />}>
          <PropertyRow label="Active" value={`${state.grid.active_balance.toFixed(2)} MW`} />
          <PropertyRow label="Reactive" value={`${state.grid.reactive_balance.toFixed(2)} MVAr`} />
        </SectionCard>
      </>
    );
  }

  // Weather Station device
  if (nodeId.startsWith('device-')) {
    const deviceId = nodeId.replace('device-', '');
    const device = state.devices.devices.find(d => d.id === deviceId);
    
    if (device) {
      return (
        <>
          <SectionCard title="Device Overview" icon={<Thermometer size={14} />}>
            <PropertyRow label="Name" value={device.name} />
            <PropertyRow label="Type" value={device.type} />
            <PropertyRow label="ID" value={device.id} />
            <PropertyRow label="Status" value={<StatusBadge status={device.state} />} />
          </SectionCard>
          <SectionCard title="Current Measurements" icon={<Activity size={14} />}>
            <PropertyRow label="Temperature" value={`${state.weather.temperature.toFixed(1)} °C`} valueColor="var(--color-environmental)" />
            <PropertyRow label="Humidity" value={`${state.weather.humidity.toFixed(1)} %`} />
            <PropertyRow label="Pressure" value={`${state.weather.pressure.toFixed(1)} hPa`} valueColor="var(--color-engineering)" />
            <PropertyRow label="Wind Speed" value={`${state.weather.wind_speed.toFixed(1)} m/s`} />
            <PropertyRow label="Wind Dir" value={`${state.weather.wind_direction.toFixed(0)}°`} />
          </SectionCard>
          <SectionCard title="Solar Context" icon={<Sun size={14} />}>
            <PropertyRow label="Irradiance" value={`${state.sun.irradiance.toFixed(0)} W/m²`} valueColor="var(--color-environmental)" />
            <PropertyRow label="Daytime" value={state.sun.is_daytime} />
            <CompassWidget direction={state.sun.azimuth} label="Sun Azimuth" />
          </SectionCard>
        </>
      );
    }
  }

  return (
    <div className={styles.emptyState}>
      <span>Select an item from the World Explorer</span>
    </div>
  );
}

function StateContent({ nodeId, state }: { nodeId: string; state: State }) {
  if (nodeId === 'sun') {
    return (
      <SectionCard title="Sun State" icon={<Sun size={14} />}>
        <PropertyRow label="Elevation" value={`${state.sun.elevation.toFixed(2)}°`} valueColor="var(--color-environmental)" />
        <PropertyRow label="Azimuth" value={`${state.sun.azimuth.toFixed(2)}°`} />
        <PropertyRow label="GHI" value={`${state.sun.irradiance.toFixed(2)} W/m²`} valueColor="var(--color-environmental)" />
        <PropertyRow label="DNI" value={`${state.sun.direct_normal.toFixed(2)} W/m²`} valueColor="var(--color-environmental)" />
        <PropertyRow label="Diffuse" value={`${state.sun.diffuse.toFixed(2)} W/m²`} />
        <PropertyRow label="Is Daytime" value={state.sun.is_daytime} />
      </SectionCard>
    );
  }

  if (nodeId === 'weather') {
    return (
      <SectionCard title="Weather State" icon={<Thermometer size={14} />}>
        <PropertyRow label="Temperature" value={`${state.weather.temperature.toFixed(2)} °C`} valueColor="var(--color-environmental)" />
        <PropertyRow label="Humidity" value={`${state.weather.humidity.toFixed(2)} %`} />
        <PropertyRow label="Pressure" value={`${state.weather.pressure.toFixed(2)} hPa`} valueColor="var(--color-engineering)" />
        <PropertyRow label="Cloud Cover" value={`${(state.weather.cloud_cover * 100).toFixed(2)} %`} />
        <PropertyRow label="Wind Speed" value={`${state.weather.wind_speed.toFixed(2)} m/s`} />
        <PropertyRow label="Wind Direction" value={`${state.weather.wind_direction.toFixed(2)}°`} />
        <PropertyRow label="Is Raining" value={state.weather.is_raining} />
      </SectionCard>
    );
  }

  if (nodeId === 'grid') {
    return (
      <SectionCard title="Grid State" icon={<Zap size={14} />}>
        <PropertyRow label="Voltage" value={`${state.grid.voltage.toFixed(4)} V`} />
        <PropertyRow label="Frequency" value={`${state.grid.frequency.toFixed(6)} Hz`} />
        <PropertyRow label="Voltage PU" value={state.grid.voltage_pu.toFixed(6)} />
        <PropertyRow label="Frequency PU" value={state.grid.frequency_pu.toFixed(6)} />
        <PropertyRow label="Active Balance" value={`${state.grid.active_balance.toFixed(4)} MW`} />
        <PropertyRow label="Reactive Balance" value={`${state.grid.reactive_balance.toFixed(4)} MVAr`} />
        <PropertyRow label="Is Stable" value={state.grid.is_stable} />
      </SectionCard>
    );
  }

  if (nodeId === 'clock') {
    return (
      <SectionCard title="Clock State" icon={<Clock size={14} />}>
        <PropertyRow label="Elapsed (ns)" value={state.clock.elapsed_ms.toLocaleString()} />
        <PropertyRow label="Elapsed" value={state.clock.elapsed} />
        <PropertyRow label="Tick Count" value={state.clock.tick_count.toLocaleString()} />
        <PropertyRow label="Mode" value={state.clock.mode} />
        <PropertyRow label="Is Paused" value={state.clock.is_paused} />
      </SectionCard>
    );
  }

  if (nodeId.startsWith('device-')) {
    const deviceId = nodeId.replace('device-', '');
    const device = state.devices.devices.find(d => d.id === deviceId);
    
    if (device) {
      return (
        <>
          <SectionCard title="Device State" icon={<Activity size={14} />}>
            <PropertyRow label="State" value={<StatusBadge status={device.state} />} />
            <PropertyRow label="Type" value={device.type} />
            <PropertyRow label="Interface Enabled" value={device.interface_enabled ? 'Yes' : 'No'} />
          </SectionCard>
          <SectionCard title="Measurements" icon={<Thermometer size={14} />}>
            <PropertyRow label="Temperature" value={`${state.weather.temperature.toFixed(2)} °C`} valueColor="var(--color-environmental)" />
            <PropertyRow label="Humidity" value={`${state.weather.humidity.toFixed(2)} %`} />
            <PropertyRow label="Pressure" value={`${state.weather.pressure.toFixed(2)} hPa`} valueColor="var(--color-engineering)" />
            <PropertyRow label="Wind Speed" value={`${state.weather.wind_speed.toFixed(2)} m/s`} />
          </SectionCard>
        </>
      );
    }
  }

  return (
    <div className={styles.emptyState}>
      <span>Select an item to view state</span>
    </div>
  );
}

function ConfigurationContent({ nodeId, state }: { nodeId: string; state: State }) {
  if (nodeId === 'sun') {
    return (
      <SectionCard title="Sun Configuration" icon={<Sun size={14} />}>
        <PropertyRow label="Latitude" value={`${state.sun.latitude.toFixed(4)}°`} />
        <PropertyRow label="Longitude" value={`${state.sun.longitude.toFixed(4)}°`} />
      </SectionCard>
    );
  }

  if (nodeId === 'weather') {
    return (
      <SectionCard title="Weather Configuration" icon={<Thermometer size={14} />}>
        <PropertyRow label="Location" value="40°N, 105°W" />
        <PropertyRow label="Elevation" value="1,640 m" />
      </SectionCard>
    );
  }

  if (nodeId === 'grid') {
    return (
      <SectionCard title="Grid Configuration" icon={<Zap size={14} />}>
        <PropertyRow label="Nominal Voltage" value={`${state.grid.nominal_voltage} V`} />
        <PropertyRow label="Nominal Frequency" value={`${state.grid.nominal_frequency} Hz`} />
      </SectionCard>
    );
  }

  if (nodeId === 'clock') {
    return (
      <SectionCard title="Clock Configuration" icon={<Clock size={14} />}>
        <PropertyRow label="Mode" value={state.clock.mode} />
      </SectionCard>
    );
  }

  if (nodeId.startsWith('device-')) {
    const deviceId = nodeId.replace('device-', '');
    const device = state.devices.devices.find(d => d.id === deviceId);
    
    if (device) {
      return (
        <SectionCard title="Device Configuration" icon={<Settings size={14} />}>
          <PropertyRow label="Device ID" value={device.id} />
          <PropertyRow label="Device Type" value={device.type} />
          <PropertyRow label="Sampling Interval" value="1,000 ms" />
        </SectionCard>
      );
    }
  }

  return (
    <div className={styles.emptyState}>
      <span>Select an item to view configuration</span>
    </div>
  );
}

function DiagnosticsContent({ nodeId, state }: { nodeId: string; state: State }) {
  if (nodeId === 'sun') {
    return (
      <SectionCard title="Sun Diagnostics" icon={<Sun size={14} />}>
        <PropertyRow label="Model Health" value="OK" valueColor="var(--color-healthy)" />
        <PropertyRow label="Last Update" value="100 ms ago" />
      </SectionCard>
    );
  }

  if (nodeId === 'weather') {
    return (
      <SectionCard title="Weather Diagnostics" icon={<Thermometer size={14} />}>
        <PropertyRow label="Model Health" value="OK" valueColor="var(--color-healthy)" />
        <PropertyRow label="Sensor Simulation" value="Active" />
        <PropertyRow label="Last Update" value="100 ms ago" />
      </SectionCard>
    );
  }

  if (nodeId === 'grid') {
    return (
      <SectionCard title="Grid Diagnostics" icon={<Zap size={14} />}>
        <PropertyRow label="Model Health" value={state.grid.is_stable ? 'OK' : 'Warning'} valueColor="var(--color-healthy)" />
        <PropertyRow label="Power Flow" value="Converged" />
        <PropertyRow label="Stability" value={state.grid.is_stable ? 'Stable' : 'Marginal'} valueColor={state.grid.is_stable ? 'var(--color-healthy)' : 'var(--color-warning)'} />
      </SectionCard>
    );
  }

  if (nodeId === 'clock') {
    return (
      <SectionCard title="Clock Diagnostics" icon={<Clock size={14} />}>
        <PropertyRow label="Tick Rate" value="10 Hz" />
        <PropertyRow label="Total Ticks" value={state.clock.tick_count.toLocaleString()} />
        <PropertyRow label="Health" value="OK" valueColor="var(--color-healthy)" />
      </SectionCard>
    );
  }

  if (nodeId.startsWith('device-')) {
    const deviceId = nodeId.replace('device-', '');
    const device = state.devices.devices.find(d => d.id === deviceId);
    
    if (device) {
      return (
        <>
          <SectionCard title="Device Diagnostics" icon={<AlertTriangle size={14} />}>
            <PropertyRow label="Firmware Status" value="Running" valueColor="var(--color-healthy)" />
            <PropertyRow label="Last Tick" value="50 ms ago" />
            <PropertyRow label="Health" value="OK" valueColor="var(--color-healthy)" />
          </SectionCard>
          <SectionCard title="Virtual Firmware" icon={<Cpu size={14} />}>
            <PropertyRow label="Firmware Version" value="1.0.0" />
            <PropertyRow label="Manufacturer" value="Forge Labs" />
            <PropertyRow label="Sampling Interval" value="1,000 ms" />
            <PropertyRow label="Memory Regions" value="4" />
          </SectionCard>
        </>
      );
    }
  }

  return (
    <div className={styles.emptyState}>
      <span>Select an item to view diagnostics</span>
    </div>
  );
}

function CommunicationsContent({ nodeId, state }: { nodeId: string; state: State }) {
  if (nodeId.startsWith('device-')) {
    const deviceId = nodeId.replace('device-', '');
    const device = state.devices.devices.find(d => d.id === deviceId);
    
    if (device) {
      const iface = device.interface;
      return (
        <>
          <SectionCard title="Communication Interface" icon={<Radio size={14} />}>
            <PropertyRow label="Enabled" value={iface?.enabled ? 'Yes' : 'No'} />
            <PropertyRow label="Status" value={<StatusBadge status={iface?.connected ? 'Connected' : 'Disconnected'} />} />
            <PropertyRow label="Interface Type" value="Raw Ingest" />
          </SectionCard>
          <SectionCard title="Traffic Statistics" icon={<Activity size={14} />}>
            <PropertyRow label="Packets Sent" value={iface?.packets_sent.toLocaleString() || '0'} />
            <PropertyRow label="Errors" value={iface?.errors.toLocaleString() || '0'} valueColor={iface?.errors ? 'var(--color-fault)' : 'var(--color-healthy)'} />
            {iface?.last_error && (
              <PropertyRow label="Last Error" value={iface.last_error} valueColor="var(--color-fault)" />
            )}
          </SectionCard>
          <SectionCard title="Device Memory" icon={<Cpu size={14} />}>
            <PropertyRow label="Temperature" value={`${state.weather.temperature.toFixed(1)} °C`} valueColor="var(--color-environmental)" />
            <PropertyRow label="Humidity" value={`${state.weather.humidity.toFixed(1)} %`} />
            <PropertyRow label="Pressure" value={`${state.weather.pressure.toFixed(1)} hPa`} valueColor="var(--color-engineering)" />
            <PropertyRow label="Wind Speed" value={`${state.weather.wind_speed.toFixed(1)} m/s`} />
            <PropertyRow label="Wind Direction" value={`${state.weather.wind_direction.toFixed(0)}°`} />
            <PropertyRow label="Quality" value="Good" valueColor="var(--color-healthy)" />
          </SectionCard>
        </>
      );
    }
  }

  // For non-device nodes, show aggregate communication info
  if (state.devices.devices.some(d => d.interface_enabled)) {
    return (
      <SectionCard title="Device Communications" icon={<Radio size={14} />}>
        <PropertyRow label="Active Interfaces" value={state.devices.devices.filter(d => d.interface_enabled).length} />
        <PropertyRow label="Total Packets" value={state.devices.devices.reduce((sum, d) => sum + (d.interface?.packets_sent || 0), 0).toLocaleString()} />
        <PropertyRow label="Total Errors" value={state.devices.devices.reduce((sum, d) => sum + (d.interface?.errors || 0), 0).toLocaleString()} />
      </SectionCard>
    );
  }

  return (
    <div className={styles.emptyState}>
      <span>Select a device to view communications</span>
    </div>
  );
}

function getNodeTitle(nodeId: string | null, state: State): string {
  if (!nodeId) return 'Inspector';
  
  const titles: Record<string, string> = {
    world: 'Simulation World',
    clock: 'Clock Model',
    sun: 'Sun Model',
    weather: 'Weather Model',
    grid: 'Grid Model',
    wind: 'Wind Model',
    models: 'Models',
    devices: 'Devices',
  };

  if (titles[nodeId]) return titles[nodeId];

  if (nodeId.startsWith('device-')) {
    const deviceId = nodeId.replace('device-', '');
    const device = state.devices.devices.find(d => d.id === deviceId);
    return device?.name || device?.type || 'Device';
  }

  return 'Inspector';
}

export function Inspector({ state, selectedNode }: InspectorProps) {
  const [activeTab, setActiveTab] = useState<InspectorTab>('overview');

  const renderTabContent = () => {
    switch (activeTab) {
      case 'overview':
        return <OverviewContent nodeId={selectedNode || ''} state={state} />;
      case 'state':
        return <StateContent nodeId={selectedNode || ''} state={state} />;
      case 'configuration':
        return <ConfigurationContent nodeId={selectedNode || ''} state={state} />;
      case 'diagnostics':
        return <DiagnosticsContent nodeId={selectedNode || ''} state={state} />;
      case 'communications':
        return <CommunicationsContent nodeId={selectedNode || ''} state={state} />;
      default:
        return null;
    }
  };

  return (
    <div className={styles.inspector}>
      <div className={styles.header}>
        <span className={styles.title}>{getNodeTitle(selectedNode, state)}</span>
      </div>
      
      <div className={styles.tabs}>
        {tabs.map((tab) => (
          <button
            key={tab.id}
            className={`${styles.tab} ${activeTab === tab.id ? styles.active : ''}`}
            onClick={() => setActiveTab(tab.id)}
            title={tab.label}
          >
            {tab.icon}
            <span className={styles.tabLabel}>{tab.label}</span>
          </button>
        ))}
      </div>

      <div className={styles.content}>
        {renderTabContent()}
      </div>
    </div>
  );
}
