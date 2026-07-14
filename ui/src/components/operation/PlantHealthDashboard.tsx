import { useMemo } from 'react';
import type { CanvasEntity } from '../../types/editor';
import type { State } from '../../types';
import styles from './PlantHealthDashboard.module.css';

interface PlantHealthDashboardProps {
  entities: CanvasEntity[];
  simulationState: State;
  onSelectEntity: (id: string) => void;
}

/**
 * Plant Health Dashboard (P0-2, P0-3)
 * 
 * High-level summary dashboard showing:
 * - Total generation (MW)
 * - Daily energy (MWh)
 * - Capacity factor (%)
 * - Performance ratio (%)
 * 
 * Design constraint: Must remain high-level summary, not SCADA dashboard.
 */
export function PlantHealthDashboard({
  entities,
  simulationState,
  onSelectEntity,
}: PlantHealthDashboardProps) {
  const { sun, weather, grid, clock } = simulationState;

  // Calculate plant-level metrics (P0-3)
  const metrics = useMemo(() => {
    // Get total plant capacity from PV arrays
    const pvArrays = entities.filter((e) => e.entity_type === 'generator');
    const totalCapacity = pvArrays.reduce((sum, pv) => {
      const capacity = parseFloat(pv.properties.rated_capacity?.value as string) || 0;
      return sum + capacity;
    }, 0);

    // Calculate current total generation
    let totalGeneration = 0;
    let expectedGeneration = 0;

    pvArrays.forEach((pv) => {
      const capacity = parseFloat(pv.properties.rated_capacity?.value as string) || 0;
      // Simulated generation based on irradiance
      const irradianceFactor = sun.is_daytime ? sun.irradiance / 1000 : 0;
      const tempDerating = Math.max(0, 1 - (weather.temperature - 25) * 0.004);
      const generation = capacity * irradianceFactor * 0.98 * tempDerating;
      totalGeneration += generation;
      expectedGeneration += capacity * irradianceFactor;
    });

    // Capacity factor: current generation / total capacity
    const capacityFactor = totalCapacity > 0
      ? (totalGeneration / totalCapacity) * 100
      : 0;

    // Daily energy (MWh) - estimate based on elapsed time
    const hoursElapsed = clock.elapsed_ms / 3600000;
    const avgGeneration = totalGeneration;
    const dailyEnergyMWh = (avgGeneration * hoursElapsed) / 1000;

    // Performance Ratio: actual vs theoretical (simplified)
    // PR = (actual yield) / (irradiated energy on panel plane)
    const theoreticalMax = totalCapacity * hoursElapsed / 1000;
    const performanceRatio = theoreticalMax > 0
      ? (dailyEnergyMWh / theoreticalMax) * 100
      : 0;

    // Plant name (from first PV array or default)
    const plantName = pvArrays[0]?.name?.match(/^(.+?)\s*\(/)?.[1] || 'Utility Solar Farm';
    const plantSize = totalCapacity > 0 ? `${totalCapacity / 1000} MWdc` : '10 MWdc';

    return {
      plantName,
      plantSize,
      totalCapacity,
      totalGeneration,
      expectedGeneration,
      capacityFactor,
      dailyEnergyMWh,
      performanceRatio,
      irradiance: sun.irradiance,
      isDaytime: sun.is_daytime,
      temperature: weather.temperature,
      gridVoltage: grid.voltage,
      gridFrequency: grid.frequency,
      gridStable: grid.is_stable,
    };
  }, [entities, sun, weather, grid, clock]);

  // Get individual PV array status
  const pvArrayStatus = useMemo(() => {
    return entities
      .filter((e) => e.entity_type === 'generator')
      .map((pv) => {
        const capacity = parseFloat(pv.properties.rated_capacity?.value as string) || 0;
        const irradianceFactor = sun.is_daytime ? sun.irradiance / 1000 : 0;
        const tempDerating = Math.max(0, 1 - (weather.temperature - 25) * 0.004);
        const generation = capacity * irradianceFactor * 0.98 * tempDerating;
        const expected = capacity * irradianceFactor;
        const deviation = expected > 0 ? ((generation - expected) / expected) * 100 : 0;

        return {
          id: pv.id,
          name: pv.name,
          designation: `PV-${pv.id.match(/(\d+)$/)?.[1]?.padStart(2, '0') || '01'}`,
          capacity,
          generation,
          deviation,
          status: Math.abs(deviation) < 5 ? 'healthy' : Math.abs(deviation) < 15 ? 'warning' : 'fault',
        };
      });
  }, [entities, sun, weather]);

  return (
    <div className={styles.dashboard}>
      {/* Plant Header */}
      <div className={styles.header}>
        <div className={styles.plantInfo}>
          <h2 className={styles.plantName}>{metrics.plantName}</h2>
          <span className={styles.plantSize}>{metrics.plantSize}</span>
        </div>
        <div className={styles.statusBadge}>
          <span className={`${styles.statusDot} ${metrics.isDaytime ? styles.daytime : styles.nighttime}`} />
          {metrics.isDaytime ? 'Daytime' : 'Night'}
        </div>
      </div>

      {/* Primary KPIs (P0-3) */}
      <div className={styles.primaryKpis}>
        <div className={styles.kpiCard}>
          <div className={styles.kpiLabel}>Total Generation</div>
          <div className={styles.kpiValue}>
            {(metrics.totalGeneration / 1000).toFixed(1)}
            <span className={styles.kpiUnit}>MW</span>
          </div>
          <div className={styles.kpiSubtext}>
            of {metrics.totalCapacity.toLocaleString()} kW
          </div>
        </div>

        <div className={styles.kpiCard}>
          <div className={styles.kpiLabel}>Today's Energy</div>
          <div className={styles.kpiValue}>
            {metrics.dailyEnergyMWh.toFixed(1)}
            <span className={styles.kpiUnit}>MWh</span>
          </div>
          <div className={styles.kpiSubtext}>
            {clock.elapsed} elapsed
          </div>
        </div>

        <div className={styles.kpiCard}>
          <div className={styles.kpiLabel}>Capacity Factor</div>
          <div className={styles.kpiValue}>
            {metrics.capacityFactor.toFixed(0)}
            <span className={styles.kpiUnit}>%</span>
          </div>
          <div className={styles.kpiSubtext}>
            Current output ratio
          </div>
        </div>

        <div className={styles.kpiCard}>
          <div className={styles.kpiLabel}>Performance Ratio</div>
          <div className={styles.kpiValue}>
            {metrics.performanceRatio.toFixed(0)}
            <span className={styles.kpiUnit}>%</span>
          </div>
          <div className={styles.kpiSubtext}>
            Actual vs theoretical
          </div>
        </div>
      </div>

      {/* Environment Summary */}
      <div className={styles.envSummary}>
        <div className={styles.envItem}>
          <span className={styles.envLabel}>Irradiance</span>
          <span className={styles.envValue}>{metrics.irradiance.toFixed(0)} W/m²</span>
        </div>
        <div className={styles.envItem}>
          <span className={styles.envLabel}>Temperature</span>
          <span className={styles.envValue}>{metrics.temperature.toFixed(1)}°C</span>
        </div>
        <div className={styles.envItem}>
          <span className={styles.envLabel}>Grid</span>
          <span className={`${styles.envValue} ${metrics.gridStable ? styles.stable : styles.unstable}`}>
            {metrics.gridVoltage.toFixed(0)}V @ {metrics.gridFrequency.toFixed(1)}Hz
          </span>
        </div>
      </div>

      {/* PV Array Status */}
      <div className={styles.arrayStatus}>
        <div className={styles.sectionTitle}>Array Status</div>
        <div className={styles.arrayGrid}>
          {pvArrayStatus.map((pv) => (
            <div
              key={pv.id}
              className={`${styles.arrayCard} ${styles[pv.status]}`}
              onClick={() => onSelectEntity(pv.id)}
            >
              <div className={styles.arrayDesignation}>{pv.designation}</div>
              <div className={styles.arrayGeneration}>
                {(pv.generation / 1000).toFixed(1)} MW
              </div>
              <div className={styles.arrayDeviation}>
                {pv.deviation.toFixed(0)}% vs expected
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
