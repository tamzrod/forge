import type { CanvasEntity, Connection } from '../types/editor';

// Utility-Scale Solar Farm Reference World
// This represents a typical 10 MW solar PV plant with:
// - Grid connection at 34.5 kV
// - Main substation with transformer
// - Switchyard with bus bars
// - Multiple PV arrays with inverters
// - Revenue metering at PCC
// - Station service loads

export interface SolarFarmWorld {
  entities: CanvasEntity[];
  connections: Connection[];
}

export function createSolarFarmWorld(): SolarFarmWorld {
  const entities: CanvasEntity[] = [
    // Grid Connection Point
    {
      id: 'grid-pcc',
      entity_type: 'grid',
      name: 'Utility Grid (34.5 kV)',
      position: { x: 400, y: 50 },
      size: { width: 100, height: 60 },
      properties: {
        name: { value: 'Utility Grid (34.5 kV)', type: 'string' },
        nominal_voltage: { value: 34500, type: 'number', unit: 'V' },
        nominal_frequency: { value: 60, type: 'number', unit: 'Hz' },
      },
    },

    // Revenue Meter
    {
      id: 'meter-revenue',
      entity_type: 'meter',
      name: 'PCC Revenue Meter',
      position: { x: 400, y: 130 },
      size: { width: 70, height: 70 },
      properties: {
        name: { value: 'PCC Revenue Meter', type: 'string' },
        meter_type: { value: 'pcc', type: 'enum', options: ['pcc', 'array', 'feeder'] },
      },
    },

    // Main Substation Transformer
    {
      id: 'tx-main',
      entity_type: 'transformer',
      name: 'Main Transformer',
      position: { x: 400, y: 220 },
      size: { width: 80, height: 60 },
      properties: {
        name: { value: 'Main Transformer', type: 'string' },
        hv_voltage: { value: 34500, type: 'number', unit: 'V' },
        lv_voltage: { value: 480, type: 'number', unit: 'V' },
        rating: { value: 12500, type: 'number', unit: 'kVA' },
        tap_position: { value: 0, type: 'number' },
      },
    },

    // Switchyard Bus 1
    {
      id: 'bus-sy-1',
      entity_type: 'bus',
      name: 'Switchyard Bus A',
      position: { x: 400, y: 310 },
      size: { width: 60, height: 60 },
      properties: {
        name: { value: 'Switchyard Bus A', type: 'string' },
        nominal_voltage: { value: 480, type: 'number', unit: 'V' },
      },
    },

    // PV Array 1
    {
      id: 'pv-array-1',
      entity_type: 'generator',
      name: 'PV Array 1 (2.5 MW)',
      position: { x: 200, y: 450 },
      size: { width: 80, height: 80 },
      properties: {
        name: { value: 'PV Array 1 (2.5 MW)', type: 'string' },
        rated_capacity: { value: 2500, type: 'number', unit: 'kW' },
        available_capacity: { value: 2500, type: 'number', unit: 'kW' },
        is_online: { value: true, type: 'boolean' },
        is_dispatchable: { value: false, type: 'boolean' },
      },
    },

    // PV Array 2
    {
      id: 'pv-array-2',
      entity_type: 'generator',
      name: 'PV Array 2 (2.5 MW)',
      position: { x: 400, y: 450 },
      size: { width: 80, height: 80 },
      properties: {
        name: { value: 'PV Array 2 (2.5 MW)', type: 'string' },
        rated_capacity: { value: 2500, type: 'number', unit: 'kW' },
        available_capacity: { value: 2500, type: 'number', unit: 'kW' },
        is_online: { value: true, type: 'boolean' },
        is_dispatchable: { value: false, type: 'boolean' },
      },
    },

    // PV Array 3
    {
      id: 'pv-array-3',
      entity_type: 'generator',
      name: 'PV Array 3 (2.5 MW)',
      position: { x: 600, y: 450 },
      size: { width: 80, height: 80 },
      properties: {
        name: { value: 'PV Array 3 (2.5 MW)', type: 'string' },
        rated_capacity: { value: 2500, type: 'number', unit: 'kW' },
        available_capacity: { value: 2500, type: 'number', unit: 'kW' },
        is_online: { value: true, type: 'boolean' },
        is_dispatchable: { value: false, type: 'boolean' },
      },
    },

    // PV Array 4
    {
      id: 'pv-array-4',
      entity_type: 'generator',
      name: 'PV Array 4 (2.5 MW)',
      position: { x: 300, y: 600 },
      size: { width: 80, height: 80 },
      properties: {
        name: { value: 'PV Array 4 (2.5 MW)', type: 'string' },
        rated_capacity: { value: 2500, type: 'number', unit: 'kW' },
        available_capacity: { value: 2500, type: 'number', unit: 'kW' },
        is_online: { value: true, type: 'boolean' },
        is_dispatchable: { value: false, type: 'boolean' },
      },
    },

    // Station Load
    {
      id: 'load-station',
      entity_type: 'load',
      name: 'Station Service',
      position: { x: 600, y: 600 },
      size: { width: 80, height: 80 },
      properties: {
        name: { value: 'Station Service', type: 'string' },
        active_power_demand: { value: 50, type: 'number', unit: 'kW' },
        power_factor: { value: 0.9, type: 'number' },
        is_connected: { value: true, type: 'boolean' },
      },
    },

    // Array Feeders - Breakers
    {
      id: 'breaker-1',
      entity_type: 'breaker',
      name: 'Feeder 1 CB',
      position: { x: 200, y: 380 },
      size: { width: 50, height: 50 },
      properties: {
        name: { value: 'Feeder 1 CB', type: 'string' },
        is_open: { value: false, type: 'boolean' },
        rating: { value: 5000, type: 'number', unit: 'A' },
      },
    },

    {
      id: 'breaker-2',
      entity_type: 'breaker',
      name: 'Feeder 2 CB',
      position: { x: 400, y: 380 },
      size: { width: 50, height: 50 },
      properties: {
        name: { value: 'Feeder 2 CB', type: 'string' },
        is_open: { value: false, type: 'boolean' },
        rating: { value: 5000, type: 'number', unit: 'A' },
      },
    },

    {
      id: 'breaker-3',
      entity_type: 'breaker',
      name: 'Feeder 3 CB',
      position: { x: 600, y: 380 },
      size: { width: 50, height: 50 },
      properties: {
        name: { value: 'Feeder 3 CB', type: 'string' },
        is_open: { value: false, type: 'boolean' },
        rating: { value: 5000, type: 'number', unit: 'A' },
      },
    },
  ];

  const connections: Connection[] = [
    // Grid to Revenue Meter
    { id: 'conn-1', from_entity: 'grid-pcc', to_entity: 'meter-revenue', from_terminal: 'output', to_terminal: 'observation' },
    
    // Revenue Meter to Transformer
    { id: 'conn-2', from_entity: 'meter-revenue', to_entity: 'tx-main', from_terminal: 'observation', to_terminal: 'hv' },
    
    // Transformer to Switchyard Bus
    { id: 'conn-3', from_entity: 'tx-main', to_entity: 'bus-sy-1', from_terminal: 'lv', to_terminal: 'input' },
    
    // Bus to Feeders
    { id: 'conn-4', from_entity: 'bus-sy-1', to_entity: 'breaker-1', from_terminal: 'output', to_terminal: 'input' },
    { id: 'conn-5', from_entity: 'bus-sy-1', to_entity: 'breaker-2', from_terminal: 'output', to_terminal: 'input' },
    { id: 'conn-6', from_entity: 'bus-sy-1', to_entity: 'breaker-3', from_terminal: 'output', to_terminal: 'input' },
    
    // Feeders to PV Arrays
    { id: 'conn-7', from_entity: 'breaker-1', to_entity: 'pv-array-1', from_terminal: 'output', to_terminal: 'input' },
    { id: 'conn-8', from_entity: 'breaker-2', to_entity: 'pv-array-2', from_terminal: 'output', to_terminal: 'input' },
    { id: 'conn-9', from_entity: 'breaker-3', to_entity: 'pv-array-3', from_terminal: 'output', to_terminal: 'input' },
    
    // Station Load connected to Bus
    { id: 'conn-10', from_entity: 'bus-sy-1', to_entity: 'load-station', from_terminal: 'output', to_terminal: 'input' },
    
    // Array 4 connected via Array 1 feeder (simplified topology)
    { id: 'conn-11', from_entity: 'breaker-1', to_entity: 'pv-array-4', from_terminal: 'output', to_terminal: 'input' },
  ];

  return { entities, connections };
}
