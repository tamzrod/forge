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

// Inspector Tab Types
export type InspectorTab = 'overview' | 'state' | 'configuration' | 'diagnostics' | 'communications';

// Navigation Types
export type Workspace = 'dashboard' | 'world' | 'devices' | 'network' | 'protocols' | 'scenarios' | 'data' | 'library' | 'settings' | 'developer';
