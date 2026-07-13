// Registry service - uses static registry data
import type { ComponentDescriptor, ComponentCategory } from '../types/registry';

export interface PaletteItem {
  id: string;
  name: string;
  category: string;
  icon: string;
  component_id: string;
  description?: string;
}

export interface PaletteResponse {
  items: PaletteItem[];
  categories: ComponentCategory[];
}

// Static registry data - matches the Go backend registry
const REGISTRY_CATEGORIES: ComponentCategory[] = [
  { id: 'electrical', name: 'Electrical', icon: '⚡', order: 1, domain: 'forge-electrical', expandable: true },
  { id: 'environment', name: 'Environment', icon: '🌤️', order: 2, domain: 'forge-environment', expandable: true },
  { id: 'simulation', name: 'Simulation', icon: '🎬', order: 3, domain: 'forge-simulation', expandable: true },
];

const REGISTRY_COMPONENTS: ComponentDescriptor[] = [
  // Electrical components
  {
    id: 'forge-electrical:grid',
    name: 'Utility Grid',
    category: 'electrical',
    icon: '🔌',
    description: 'Utility grid connection point',
    properties: [
      { key: 'name', label: 'Name', type: 'string', default: 'Utility Grid', required: true },
      { key: 'nominal_voltage', label: 'Nominal Voltage', type: 'number', default: 69000, unit: 'V' },
      { key: 'nominal_frequency', label: 'Frequency', type: 'number', default: 60, unit: 'Hz', options: ['50', '60'] },
    ],
    terminals: [{ id: 'output', name: 'Output', role: 'source', voltage: 69000, direction: 'output' }],
    width: 80,
    height: 60,
    domain: 'forge-electrical',
  },
  {
    id: 'forge-electrical:bus',
    name: 'Bus',
    category: 'electrical',
    icon: '⚫',
    description: 'Electrical bus node',
    properties: [
      { key: 'name', label: 'Name', type: 'string', default: 'New Bus', required: true },
      { key: 'nominal_voltage', label: 'Nominal Voltage', type: 'number', default: 480, unit: 'V' },
    ],
    terminals: [
      { id: 'input', name: 'Input', role: 'through', direction: 'input' },
      { id: 'output', name: 'Output', role: 'through', direction: 'output' },
    ],
    width: 60,
    height: 60,
    domain: 'forge-electrical',
  },
  {
    id: 'forge-electrical:breaker',
    name: 'Breaker',
    category: 'electrical',
    icon: '🔀',
    description: 'Circuit breaker switch',
    properties: [
      { key: 'name', label: 'Name', type: 'string', default: 'Circuit Breaker', required: true },
      { key: 'is_open', label: 'Open', type: 'boolean', default: false },
      { key: 'rating', label: 'Rating', type: 'number', default: 1200, unit: 'A' },
    ],
    terminals: [
      { id: 'input', name: 'Input', role: 'through', direction: 'input' },
      { id: 'output', name: 'Output', role: 'through', direction: 'output' },
    ],
    width: 50,
    height: 50,
    domain: 'forge-electrical',
  },
  {
    id: 'forge-electrical:transformer',
    name: 'Transformer',
    category: 'electrical',
    icon: '🔄',
    description: 'Power transformer',
    properties: [
      { key: 'name', label: 'Name', type: 'string', default: 'Transformer', required: true },
      { key: 'hv_voltage', label: 'HV Voltage', type: 'number', default: 69000, unit: 'V' },
      { key: 'lv_voltage', label: 'LV Voltage', type: 'number', default: 480, unit: 'V' },
      { key: 'rating', label: 'Rating', type: 'number', default: 1000, unit: 'kVA' },
      { key: 'tap_position', label: 'Tap Position', type: 'number', default: 0 },
    ],
    terminals: [
      { id: 'hv', name: 'HV', role: 'through', direction: 'input' },
      { id: 'lv', name: 'LV', role: 'through', direction: 'output' },
    ],
    width: 80,
    height: 60,
    domain: 'forge-electrical',
  },
  {
    id: 'forge-electrical:generator',
    name: 'Virtual Generator',
    category: 'electrical',
    icon: '☀️',
    description: 'Solar or wind generator',
    properties: [
      { key: 'name', label: 'Name', type: 'string', default: 'Solar Generator', required: true },
      { key: 'rated_capacity', label: 'Rated Capacity', type: 'number', default: 500, unit: 'kW' },
      { key: 'available_capacity', label: 'Available Capacity', type: 'number', default: 500, unit: 'kW' },
      { key: 'is_online', label: 'Online', type: 'boolean', default: true },
      { key: 'is_dispatchable', label: 'Dispatchable', type: 'boolean', default: true },
    ],
    terminals: [{ id: 'output', name: 'Output', role: 'source', direction: 'output' }],
    width: 80,
    height: 80,
    domain: 'forge-electrical',
  },
  {
    id: 'forge-electrical:load',
    name: 'Virtual Load',
    category: 'electrical',
    icon: '🏭',
    description: 'Factory or facility load',
    properties: [
      { key: 'name', label: 'Name', type: 'string', default: 'Factory Load', required: true },
      { key: 'active_power_demand', label: 'Active Power', type: 'number', default: 400, unit: 'kW' },
      { key: 'power_factor', label: 'Power Factor', type: 'number', default: 0.9 },
      { key: 'is_connected', label: 'Connected', type: 'boolean', default: true },
    ],
    terminals: [{ id: 'input', name: 'Input', role: 'destination', direction: 'input' }],
    width: 80,
    height: 80,
    domain: 'forge-electrical',
  },
  {
    id: 'forge-electrical:meter',
    name: 'Meter',
    category: 'electrical',
    icon: '📊',
    description: 'Power measurement meter',
    properties: [
      { key: 'name', label: 'Name', type: 'string', default: 'PCC Meter', required: true },
      { key: 'meter_type', label: 'Type', type: 'enum', default: 'pcc', options: ['pcc', 'array', 'feeder'] },
    ],
    terminals: [{ id: 'observation', name: 'Observation', role: 'observation', direction: 'bidirectional' }],
    width: 70,
    height: 70,
    domain: 'forge-electrical',
  },
  // Environment components
  {
    id: 'forge-environment:sun',
    name: 'Sun',
    category: 'environment',
    icon: '🌞',
    description: 'Solar position and irradiance',
    properties: [
      { key: 'name', label: 'Name', type: 'string', default: 'Sun', required: true },
      { key: 'latitude', label: 'Latitude', type: 'number', default: 35.2271 },
      { key: 'longitude', label: 'Longitude', type: 'number', default: -80.8431 },
      { key: 'tilt', label: 'Panel Tilt', type: 'number', default: 20, unit: '°' },
      { key: 'azimuth', label: 'Azimuth', type: 'number', default: 180, unit: '°' },
    ],
    terminals: [{ id: 'output', name: 'Irradiance', role: 'source', direction: 'output' }],
    width: 60,
    height: 60,
    domain: 'forge-environment',
  },
  {
    id: 'forge-environment:weather',
    name: 'Weather',
    category: 'environment',
    icon: '🌤️',
    description: 'Weather conditions',
    properties: [
      { key: 'name', label: 'Name', type: 'string', default: 'Weather', required: true },
      { key: 'temperature', label: 'Temperature', type: 'number', default: 25, unit: '°C' },
      { key: 'humidity', label: 'Humidity', type: 'number', default: 50, unit: '%' },
      { key: 'cloud_cover', label: 'Cloud Cover', type: 'number', default: 0, unit: '%' },
    ],
    terminals: [{ id: 'output', name: 'Conditions', role: 'observation', direction: 'output' }],
    width: 60,
    height: 60,
    domain: 'forge-environment',
  },
  {
    id: 'forge-environment:wind',
    name: 'Wind',
    category: 'environment',
    icon: '💨',
    description: 'Wind conditions',
    properties: [
      { key: 'name', label: 'Name', type: 'string', default: 'Wind', required: true },
      { key: 'speed', label: 'Speed', type: 'number', default: 5, unit: 'm/s' },
      { key: 'direction', label: 'Direction', type: 'number', default: 0, unit: '°' },
    ],
    terminals: [{ id: 'output', name: 'Wind Data', role: 'source', direction: 'output' }],
    width: 60,
    height: 60,
    domain: 'forge-environment',
  },
  // Simulation components
  {
    id: 'forge-simulation:scenario',
    name: 'Scenario',
    category: 'simulation',
    icon: '🎬',
    description: 'Test scenario',
    properties: [
      { key: 'name', label: 'Name', type: 'string', default: 'Test Scenario', required: true },
      { key: 'duration', label: 'Duration', type: 'number', default: 3600, unit: 's' },
      { key: 'description', label: 'Description', type: 'string', default: '' },
    ],
    terminals: [],
    width: 80,
    height: 60,
    domain: 'forge-simulation',
  },
  {
    id: 'forge-simulation:clock',
    name: 'Simulation Clock',
    category: 'simulation',
    icon: '⏱️',
    description: 'Simulation time control',
    properties: [
      { key: 'name', label: 'Name', type: 'string', default: 'Clock', required: true },
      { key: 'start_time', label: 'Start Time', type: 'string', default: '2024-01-01T08:00:00Z' },
      { key: 'end_time', label: 'End Time', type: 'string', default: '2024-01-01T20:00:00Z' },
      { key: 'time_step', label: 'Time Step', type: 'number', default: 100, unit: 'ms' },
    ],
    terminals: [{ id: 'output', name: 'Time', role: 'source', direction: 'output' }],
    width: 60,
    height: 60,
    domain: 'forge-simulation',
  },
];

class RegistryService {
  async getComponents(): Promise<ComponentDescriptor[]> {
    return REGISTRY_COMPONENTS;
  }

  async getCategories(): Promise<ComponentCategory[]> {
    return REGISTRY_CATEGORIES;
  }

  async getPalette(): Promise<PaletteResponse> {
    const items: PaletteItem[] = REGISTRY_COMPONENTS.map((comp) => ({
      id: comp.id,
      name: comp.name,
      category: comp.category,
      icon: comp.icon,
      component_id: comp.id,
      description: comp.description,
    }));

    return {
      items,
      categories: REGISTRY_CATEGORIES,
    };
  }

  async getComponent(id: string): Promise<ComponentDescriptor | null> {
    return REGISTRY_COMPONENTS.find((c) => c.id === id) || null;
  }

  clearCache(): void {
    // No-op for static data
  }
}

export const registryService = new RegistryService();
