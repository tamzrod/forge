import { useState, useEffect, useCallback, useRef } from 'react';
import type { State } from '../types';

const DEFAULT_STATE: State = {
  clock: {
    elapsed: '0s',
    elapsed_ms: 0,
    tick_count: 0,
    mode: 'Stopped',
    is_paused: true,
  },
  sun: {
    elevation: 0,
    azimuth: 0,
    irradiance: 0,
    direct_normal: 0,
    diffuse: 0,
    is_daytime: false,
    latitude: 0,
    longitude: 0,
  },
  weather: {
    temperature: 0,
    humidity: 0,
    pressure: 0,
    cloud_cover: 0,
    wind_speed: 0,
    wind_direction: 0,
    is_raining: false,
  },
  grid: {
    voltage: 0,
    frequency: 0,
    voltage_pu: 1,
    frequency_pu: 1,
    active_balance: 0,
    reactive_balance: 0,
    is_stable: true,
    nominal_voltage: 480,
    nominal_frequency: 60,
  },
  devices: {
    count: 0,
    devices: [],
  },
};

function formatDuration(ns: number): string {
  const totalSeconds = Math.floor(ns / 1_000_000_000);
  const days = Math.floor(totalSeconds / 86400);
  const hours = Math.floor((totalSeconds % 86400) / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  const seconds = totalSeconds % 60;

  if (days > 0) {
    return `${days}d ${hours}h ${minutes}m`;
  } else if (hours > 0) {
    return `${hours}h ${minutes}m ${seconds}s`;
  } else if (minutes > 0) {
    return `${minutes}m ${seconds}s`;
  } else {
    return `${seconds}s`;
  }
}

function transformState(raw: Record<string, unknown>): State {
  const clock = raw.clock as Record<string, unknown> || {};
  const sun = raw.sun as Record<string, unknown> || {};
  const weather = raw.weather as Record<string, unknown> || {};
  const grid = raw.grid as Record<string, unknown> || {};
  const devices = raw.devices as Record<string, unknown> || {};

  return {
    clock: {
      elapsed: formatDuration((clock.elapsed as number) || 0),
      elapsed_ms: (clock.elapsed as number) || 0,
      tick_count: (clock.tick_count as number) || 0,
      mode: (clock.mode as string) || 'Stopped',
      is_paused: (clock.is_paused as boolean) ?? true,
    },
    sun: {
      elevation: (sun.elevation as number) || 0,
      azimuth: (sun.azimuth as number) || 0,
      irradiance: (sun.irradiance as number) || 0,
      direct_normal: (sun.direct_normal as number) || 0,
      diffuse: (sun.diffuse as number) || 0,
      is_daytime: (sun.is_daytime as boolean) || false,
      latitude: (sun.latitude as number) || 0,
      longitude: (sun.longitude as number) || 0,
    },
    weather: {
      temperature: (weather.temperature as number) || 0,
      humidity: (weather.humidity as number) || 0,
      pressure: (weather.pressure as number) || 0,
      cloud_cover: (weather.cloud_cover as number) || 0,
      wind_speed: (weather.wind_speed as number) || 0,
      wind_direction: (weather.wind_direction as number) || 0,
      is_raining: (weather.is_raining as boolean) || false,
    },
    grid: {
      voltage: (grid.voltage as number) || 0,
      frequency: (grid.frequency as number) || 0,
      voltage_pu: (grid.voltage_pu as number) || 1,
      frequency_pu: (grid.frequency_pu as number) || 1,
      active_balance: (grid.active_balance as number) || 0,
      reactive_balance: (grid.reactive_balance as number) || 0,
      is_stable: (grid.is_stable as boolean) ?? true,
      nominal_voltage: (grid.nominal_voltage as number) || 480,
      nominal_frequency: (grid.nominal_frequency as number) || 60,
    },
    devices: {
      count: (devices.count as number) || 0,
      devices: (devices.devices as State['devices']['devices']) || [],
    },
  };
}

interface UseSimulationReturn {
  state: State;
  connected: boolean;
  error: string | null;
  lastUpdate: Date | null;
}

export function useSimulation(wsUrl?: string): UseSimulationReturn {
  const [state, setState] = useState<State>(DEFAULT_STATE);
  const [connected, setConnected] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [lastUpdate, setLastUpdate] = useState<Date | null>(null);
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<number | null>(null);

  const connect = useCallback(() => {
    // Determine WebSocket URL
    const url = wsUrl || `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}/ws`;
    
    try {
      const ws = new WebSocket(url);
      wsRef.current = ws;

      ws.onopen = () => {
        setConnected(true);
        setError(null);
      };

      ws.onmessage = (event) => {
        try {
          const raw = JSON.parse(event.data);
          const transformed = transformState(raw);
          setState(transformed);
          setLastUpdate(new Date());
        } catch {
          console.error('Failed to parse WebSocket message');
        }
      };

      ws.onclose = () => {
        setConnected(false);
        wsRef.current = null;
        // Reconnect after 2 seconds
        reconnectTimeoutRef.current = window.setTimeout(() => {
          connect();
        }, 2000);
      };

      ws.onerror = () => {
        setError('WebSocket connection error');
        setConnected(false);
      };
    } catch {
      setError('Failed to connect');
      // Retry connection
      reconnectTimeoutRef.current = window.setTimeout(() => {
        connect();
      }, 2000);
    }
  }, [wsUrl]);

  useEffect(() => {
    connect();

    return () => {
      if (wsRef.current) {
        wsRef.current.close();
      }
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
    };
  }, [connect]);

  return { state, connected, error, lastUpdate };
}

// REST API fetch for initial state
export async function fetchState(apiUrl?: string): Promise<State> {
  const url = apiUrl || '/api/state';
  const response = await fetch(url);
  if (!response.ok) {
    throw new Error(`Failed to fetch state: ${response.statusText}`);
  }
  const raw = await response.json();
  return transformState(raw);
}
