/**
 * Simulation Data Service
 * Bridges between UI entities and simulation state.
 * All measurements originate from simulation - UI is a passive observer.
 */

import type { CanvasEntity } from '../types/editor';
import type { 
  State, 
  EntityMeasurements, 
  PVArrayMeasurement, 
  TransformerMeasurement,
  BusMeasurement,
  MeterMeasurement,
  BreakerMeasurement,
  LoadMeasurement 
} from '../types';

// Default measurement values (used when no simulation data is available)
const DEFAULT_PV_MEASUREMENT: PVArrayMeasurement = {
  dc_power: 0,
  ac_power: 0,
  dc_voltage: 0,
  dc_current: 0,
  efficiency: 0,
  inverter_temp: 0,
  operating_state: 'night',
};

/**
 * Physics constants for solar PV calculations
 * These should come from simulation but are derived from entity properties
 */
const INVERTER_EFFICIENCY = 0.98;
const DC_VOLTAGE_NOMINAL = 400;

/**
 * Computes PV array measurements from simulation state and entity properties.
 * All values originate from simulation state.
 */
export function computePVArrayMeasurement(
  entity: CanvasEntity, 
  sun: State['sun'], 
  weather: State['weather']
): PVArrayMeasurement {
  const capacity = parseFloat(entity.properties.rated_capacity?.value as string) || 0;
  const efficiency = 0.20; // Panel efficiency from entity/simulation
  
  // Only generate if sun is up
  if (!sun.is_daytime || sun.irradiance <= 0) {
    return {
      ...DEFAULT_PV_MEASUREMENT,
      operating_state: 'night',
    };
  }

  // DC power = irradiance * area * efficiency / 1000 (kW)
  // Area is derived from rated capacity (assuming 1000 W/m² and 20% efficiency)
  const panelArea = capacity / (1000 * efficiency); // m²
  const dcPower = (sun.irradiance / 1000) * panelArea * efficiency;
  
  // AC power = DC power * inverter efficiency
  const acPower = dcPower * INVERTER_EFFICIENCY;
  
  // DC voltage (nominal, simplified)
  const dcVoltage = DC_VOLTAGE_NOMINAL;
  
  // DC current = power / voltage
  const dcCurrent = dcVoltage > 0 ? dcPower * 1000 / dcVoltage : 0;
  
  // Temperature derating: efficiency decreases above 25°C
  const tempDerating = Math.max(0, 1 - (weather.temperature - 25) * 0.004);
  
  // Inverter temperature (simplified - ambient + losses)
  const inverterTemp = weather.temperature + 10;

  return {
    dc_power: Math.min(dcPower, capacity), // Clamp to rated capacity
    ac_power: Math.min(acPower, capacity * INVERTER_EFFICIENCY),
    dc_voltage: dcVoltage,
    dc_current: dcCurrent,
    efficiency: efficiency * INVERTER_EFFICIENCY * tempDerating,
    inverter_temp: inverterTemp,
    operating_state: acPower > capacity * 0.01 ? 'generating' : 'standby',
  };
}

/**
 * Computes transformer measurements from simulation state.
 */
export function computeTransformerMeasurement(
  entity: CanvasEntity,
  grid: State['grid'],
  totalGeneration: number,
  totalLoad: number
): TransformerMeasurement {
  const hvVoltage = parseFloat(entity.properties.hv_voltage?.value as string) || 34500;
  const lvVoltage = parseFloat(entity.properties.lv_voltage?.value as string) || 480;
  const rating = parseFloat(entity.properties.rating?.value as string) || 10000; // kVA
  const tapPosition = parseFloat(entity.properties.tap_position?.value as string) || 0;
  
  // Net power through transformer
  const netPower = totalGeneration - totalLoad;
  
  // Load percentage (assume rated at 480V LV side)
  const loadPercent = rating > 0 ? Math.min(100, (Math.abs(netPower) / rating) * 100) : 0;
  
  // Oil temperature (simplified - increases with load)
  const baseTemp = 30;
  const loadTemp = loadPercent * 0.5;
  const ambientTemp = 0; // Could come from weather
  const oilTemp = baseTemp + loadTemp + ambientTemp;

  return {
    primary_voltage: hvVoltage * grid.voltage_pu,
    secondary_voltage: lvVoltage * grid.voltage_pu,
    load_percent: loadPercent,
    oil_temp: oilTemp,
    tap_position: tapPosition,
  };
}

/**
 * Computes bus measurements from simulation state.
 */
export function computeBusMeasurement(
  entity: CanvasEntity,
  grid: State['grid']
): BusMeasurement {
  // Scale voltage to the bus's nominal voltage from entity properties
  const nominalVoltage = parseFloat(entity.properties.nominal_voltage?.value as string) || 480;
  
  return {
    voltage: grid.voltage * grid.voltage_pu * (nominalVoltage / 480), // Scale to bus nominal
    voltage_pu: grid.voltage_pu,
    frequency: grid.frequency,
  };
}

/**
 * Computes meter measurements from simulation state.
 */
export function computeMeterMeasurement(
  _entity: CanvasEntity,
  grid: State['grid'],
  totalGeneration: number,
  totalLoad: number,
  energyExport: number,
  energyImport: number
): MeterMeasurement {
  // Net power is computed from generation and load from simulation
  const netPower = totalGeneration - totalLoad;
  
  // Reactive power (simplified - assume 0.95 PF)
  const reactivePower = Math.abs(netPower) * Math.sqrt(1 - 0.95 * 0.95);
  
  // Power factor
  const powerFactor = netPower >= 0 ? 0.95 : -0.95;

  return {
    voltage: grid.voltage,
    frequency: grid.frequency,
    active_power: netPower,
    reactive_power: reactivePower * (netPower >= 0 ? 1 : -1),
    power_factor: powerFactor,
    energy_export: energyExport,
    energy_import: energyImport,
  };
}

/**
 * Computes breaker measurements from simulation state.
 */
export function computeBreakerMeasurement(
  _entity: CanvasEntity
): BreakerMeasurement {
  // For now, breaker state comes from entity properties
  // In a full implementation, this would come from simulation
  const isOpen = _entity.properties.is_open?.value as boolean || false;
  
  return {
    is_open: isOpen,
    trip_count: 0, // Would come from simulation events
  };
}

/**
 * Computes load measurements from simulation state.
 */
export function computeLoadMeasurement(
  entity: CanvasEntity
): LoadMeasurement {
  const activePower = parseFloat(entity.properties.active_power_demand?.value as string) || 0;
  const powerFactor = parseFloat(entity.properties.power_factor?.value as string) || 0.9;
  
  return {
    active_power: activePower,
    power_factor: powerFactor,
  };
}

/**
 * Computes all entity measurements from simulation state.
 * This is the single source of truth for equipment measurements.
 */
export function computeAllMeasurements(
  entities: CanvasEntity[],
  state: State
): EntityMeasurements {
  const measurements: EntityMeasurements = {};
  
  // First pass: compute total generation and load
  let totalGeneration = 0;
  let totalLoad = 0;
  
  for (const entity of entities) {
    if (entity.entity_type === 'generator') {
      const pv = computePVArrayMeasurement(entity, state.sun, state.weather);
      totalGeneration += pv.ac_power;
    } else if (entity.entity_type === 'load') {
      const load = computeLoadMeasurement(entity);
      totalLoad += load.active_power;
    }
  }
  
  // Track energy (simplified)
  const tickHours = 0.1 / 3600;
  let energyExport = 0;
  let energyImport = 0;
  if (totalGeneration > totalLoad) {
    energyExport = (totalGeneration - totalLoad) * tickHours;
  } else {
    energyImport = (totalLoad - totalGeneration) * tickHours;
  }
  
  // Second pass: compute all measurements
  for (const entity of entities) {
    switch (entity.entity_type) {
      case 'generator':
        measurements[entity.id] = computePVArrayMeasurement(entity, state.sun, state.weather);
        break;
        
      case 'transformer':
        measurements[entity.id] = computeTransformerMeasurement(
          entity, state.grid, totalGeneration, totalLoad
        );
        break;
        
      case 'bus':
        measurements[entity.id] = computeBusMeasurement(entity, state.grid);
        break;
        
      case 'meter':
        measurements[entity.id] = computeMeterMeasurement(
          entity, state.grid, totalGeneration, totalLoad, energyExport, energyImport
        );
        break;
        
      case 'breaker':
        measurements[entity.id] = computeBreakerMeasurement(entity);
        break;
        
      case 'load':
        measurements[entity.id] = computeLoadMeasurement(entity);
        break;
        
      case 'grid':
        // Grid doesn't have specific measurements, uses grid state directly
        measurements[entity.id] = {};
        break;
        
      default:
        measurements[entity.id] = {};
    }
  }
  
  return measurements;
}

/**
 * Gets a typed measurement for a specific entity.
 */
export function getEntityMeasurement<T extends keyof EntityMeasurements>(
  entityId: string,
  measurements: EntityMeasurements,
  _type: T
): EntityMeasurements[T] | undefined {
  const measurement = measurements[entityId];
  if (measurement && Object.keys(measurement).length > 0) {
    return measurement as EntityMeasurements[T];
  }
  return undefined;
}
