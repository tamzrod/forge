import { useState, useCallback } from 'react';
import { WelcomeScreen } from './components/welcome';
import { OperationWorkspace } from './components/operation';
import { createSolarFarmWorld } from './data/solarFarmWorld';
import type { CanvasEntity } from './types/editor';
import type { State } from './types';
import styles from './App.module.css';

// Default simulation state
const DEFAULT_STATE: State = {
  clock: {
    elapsed: '00:00:00',
    elapsed_ms: 0,
    tick_count: 0,
    mode: 'Stopped',
    is_paused: true,
  },
  sun: {
    elevation: 45,
    azimuth: 180,
    irradiance: 850,
    direct_normal: 750,
    diffuse: 100,
    is_daytime: true,
    latitude: 33.9,
    longitude: -117.5,
  },
  weather: {
    temperature: 28,
    humidity: 45,
    pressure: 1013,
    cloud_cover: 15,
    wind_speed: 3.5,
    wind_direction: 225,
    is_raining: false,
  },
  grid: {
    voltage: 480,
    frequency: 60.0,
    voltage_pu: 1.0,
    frequency_pu: 1.0,
    active_balance: 0,
    reactive_balance: 0,
    is_stable: true,
    nominal_voltage: 480,
    nominal_frequency: 60,
  },
  devices: {
    count: 4,
    devices: [],
  },
};

function App() {
  // App state
  const [view, setView] = useState<'welcome' | 'operation'>('welcome');
  const [entities, setEntities] = useState<CanvasEntity[]>([]);
  const [connections, setConnections] = useState<Array<{ id: string; from_entity: string; to_entity: string }>>([]);
  const [simulationState, setSimulationState] = useState<State>(DEFAULT_STATE);

  // Handle load solar farm
  const handleLoadSolarFarm = useCallback(() => {
    const world = createSolarFarmWorld();
    setEntities(world.entities);
    setConnections(world.connections);
    setSimulationState(DEFAULT_STATE);
    setView('operation');
  }, []);

  // Handle open existing project
  const handleOpenExisting = useCallback(() => {
    // TODO: Implement file picker
    console.log('Open existing project...');
  }, []);

  // Handle simulation run
  const handleRun = useCallback(() => {
    setSimulationState((prev) => ({
      ...prev,
      clock: {
        ...prev.clock,
        mode: 'Running',
        is_paused: false,
      },
    }));

    // Start simulation loop
    let tickCount = 0;
    const interval = setInterval(() => {
      tickCount++;
      setSimulationState((prev) => {
        // Update sun position (simplified)
        const elapsed = prev.clock.elapsed_ms + 1000;
        const hours = Math.floor(elapsed / 3600000) % 24;
        const elevation = hours >= 6 && hours <= 18 
          ? Math.sin((hours - 6) / 12 * Math.PI) * 75 
          : 0;
        const irradiance = elevation > 0 ? (elevation / 75) * 1000 : 0;
        const cloudCover = prev.weather.cloud_cover + (Math.random() - 0.5) * 2;
        
        return {
          ...prev,
          clock: {
            ...prev.clock,
            elapsed_ms: elapsed,
            elapsed: formatDuration(elapsed),
            tick_count: prev.clock.tick_count + 1,
          },
          sun: {
            ...prev.sun,
            elevation: Math.max(0, elevation),
            irradiance: Math.max(0, irradiance),
            is_daytime: elevation > 0,
          },
          weather: {
            ...prev.weather,
            cloud_cover: Math.max(0, Math.min(100, cloudCover)),
            temperature: prev.weather.temperature + (Math.random() - 0.5) * 0.5,
          },
        };
      });
    }, 1000);

    // Store interval for cleanup
    (window as unknown as { simulationInterval: number }).simulationInterval = interval as unknown as number;
  }, []);

  // Handle simulation pause
  const handlePause = useCallback(() => {
    const interval = (window as unknown as { simulationInterval: number }).simulationInterval;
    if (interval) {
      clearInterval(interval);
    }
    setSimulationState((prev) => ({
      ...prev,
      clock: {
        ...prev.clock,
        is_paused: !prev.clock.is_paused,
      },
    }));
  }, []);

  // Handle simulation reset
  const handleReset = useCallback(() => {
    const interval = (window as unknown as { simulationInterval: number }).simulationInterval;
    if (interval) {
      clearInterval(interval);
    }
    setSimulationState(DEFAULT_STATE);
  }, []);

  // Handle speed change
  const handleSpeedChange = useCallback((speed: number) => {
    // TODO: Implement speed adjustment
    console.log('Speed:', speed);
  }, []);

  // Handle scenario change
  const handleScenarioChange = useCallback((scenarioId: string) => {
    // TODO: Implement scenario loading
    console.log('Scenario:', scenarioId);
  }, []);

  // Handle property change
  const handlePropertyChange = useCallback((entityId: string, key: string, value: unknown) => {
    setEntities((prev) =>
      prev.map((e) =>
        e.id === entityId
          ? {
              ...e,
              properties: { ...e.properties, [key]: { ...e.properties[key], value } },
            }
          : e
      )
    );
  }, []);

  return (
    <div className={styles.app}>
      {view === 'welcome' ? (
        <WelcomeScreen
          onLoadSolarFarm={handleLoadSolarFarm}
          onOpenExisting={handleOpenExisting}
        />
      ) : (
        <OperationWorkspace
          entities={entities}
          connections={connections}
          simulationState={simulationState}
          onPropertyChange={handlePropertyChange}
          onRun={handleRun}
          onPause={handlePause}
          onReset={handleReset}
          onSpeedChange={handleSpeedChange}
          onScenarioChange={handleScenarioChange}
        />
      )}
    </div>
  );
}

function formatDuration(ms: number): string {
  const totalSeconds = Math.floor(ms / 1000);
  const hours = Math.floor(totalSeconds / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  const seconds = totalSeconds % 60;
  return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
}

export default App;
