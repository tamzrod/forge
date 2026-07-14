// Simulation State Types

export interface ClockState {
  elapsed: string;
  elapsed_ms: number;
  tick_count: number;
  mode: string;
  is_paused: boolean;
}

export interface SunState {
  elevation: number;
  azimuth: number;
  irradiance: number;
  direct_normal: number;
  diffuse: number;
  is_daytime: boolean;
  latitude: number;
  longitude: number;
}

export interface WeatherState {
  temperature: number;
  humidity: number;
  pressure: number;
  cloud_cover: number;
  wind_speed: number;
  wind_direction: number;
  is_raining: boolean;
}

export interface GridState {
  voltage: number;
  frequency: number;
  voltage_pu: number;
  frequency_pu: number;
  active_balance: number;
  reactive_balance: number;
  is_stable: boolean;
  nominal_voltage: number;
  nominal_frequency: number;
}

// Equipment-specific measurements from simulation
export interface PVArrayMeasurement {
  dc_power: number;         // DC power output in kW
  ac_power: number;         // AC power output after inverter in kW
  dc_voltage: number;       // DC voltage in V
  dc_current: number;       // DC current in A
  efficiency: number;       // Overall efficiency (0-1)
  inverter_temp: number;    // Inverter temperature in °C
  operating_state: 'generating' | 'standby' | 'fault' | 'night';
}

export interface TransformerMeasurement {
  primary_voltage: number;  // Primary side voltage in V
  secondary_voltage: number; // Secondary side voltage in V
  load_percent: number;     // Current load as percentage
  oil_temp: number;         // Oil temperature in °C
  tap_position: number;     // Tap changer position
}

export interface BusMeasurement {
  voltage: number;         // Bus voltage in V
  voltage_pu: number;      // Per-unit voltage
  frequency: number;        // Frequency in Hz
}

export interface MeterMeasurement {
  voltage: number;          // Voltage in V
  frequency: number;        // Frequency in Hz
  active_power: number;     // Active power in kW
  reactive_power: number;   // Reactive power in kVAr
  power_factor: number;     // Power factor
  energy_export: number;    // Exported energy in kWh
  energy_import: number;    // Imported energy in kWh
}

export interface BreakerMeasurement {
  is_open: boolean;        // Breaker status
  trip_count: number;      // Number of trips
}

export interface LoadMeasurement {
  active_power: number;     // Active power demand in kW
  power_factor: number;     // Power factor
}

// Entity measurements indexed by entity ID
export interface EntityMeasurements {
  [entityId: string]: 
    | PVArrayMeasurement 
    | TransformerMeasurement 
    | BusMeasurement 
    | MeterMeasurement 
    | BreakerMeasurement 
    | LoadMeasurement 
    | Record<string, never>; // Empty for unknown entities
}

export interface InterfaceInfo {
  enabled: boolean;
  connected: boolean;
  packets_sent: number;
  errors: number;
  last_error?: string;
}

export interface DeviceState {
  id: string;
  type: string;
  name: string;
  state: string;
  interface_enabled?: boolean;
  interface?: InterfaceInfo;
}

export interface DevicesState {
  count: number;
  devices: DeviceState[];
}

export interface State {
  clock: ClockState;
  sun: SunState;
  weather: WeatherState;
  grid: GridState;
  devices: DevicesState;
  measurements: EntityMeasurements;
}

// World Explorer Tree Types
export interface TreeNode {
  id: string;
  label: string;
  icon?: string;
  children?: TreeNode[];
  type: 'root' | 'models' | 'device' | 'model' | 'network' | 'scenarios';
  data?: unknown;
}

// Inspector Tab Types (legacy - for backward compatibility)
export type InspectorTab = 'overview' | 'state' | 'configuration' | 'diagnostics' | 'communications';

// Navigation Types
export type Workspace = 'dashboard' | 'world' | 'devices' | 'network' | 'protocols' | 'scenarios' | 'data' | 'library' | 'settings' | 'developer';

// Re-export Generic Inspector types
export * from './inspector';
