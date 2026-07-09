import { useState, useEffect, useCallback } from 'react';
import type { GenericInspectorData } from '../types/inspector';

interface UseInspectorReturn {
  data: GenericInspectorData | null;
  loading: boolean;
  error: string | null;
  refetch: () => void;
}

export function useInspector(objectId: string | null, apiBase?: string): UseInspectorReturn {
  const [data, setData] = useState<GenericInspectorData | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [trigger, setTrigger] = useState(0);

  const apiUrl = apiBase || '/api';

  const refetch = useCallback(() => {
    setTrigger(t => t + 1);
  }, []);

  useEffect(() => {
    if (!objectId) {
      setData(null);
      setLoading(false);
      setError(null);
      return;
    }

    let cancelled = false;

    async function fetchData() {
      setLoading(true);
      setError(null);

      try {
        const response = await fetch(`${apiUrl}/inspect/${objectId}`);
        
        if (!response.ok) {
          throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }
        
        const inspectionData: GenericInspectorData = await response.json();
        
        if (!cancelled) {
          setData(inspectionData);
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : 'Unknown error');
          setData(null);
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    fetchData();

    return () => {
      cancelled = true;
    };
  }, [objectId, apiUrl, trigger]);

  return { data, loading, error, refetch };
}

// REST API fetch for inspection data
export async function fetchInspectionData(objectId: string, apiBase?: string): Promise<GenericInspectorData> {
  const url = apiBase ? `${apiBase}/inspect/${objectId}` : `/api/inspect/${objectId}`;
  const response = await fetch(url);
  if (!response.ok) {
    throw new Error(`Failed to fetch inspection data: ${response.statusText}`);
  }
  return response.json();
}
