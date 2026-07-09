/**
 * Generic Inspector Component
 * 
 * A data-driven inspector that displays inspection data from the backend.
 * This replaces the hardcoded Inspector with a generic, extensible framework.
 * 
 * Supported Objects:
 * - Simulation Models (Clock, Sun, Weather, Grid, Wind)
 * - Virtual Devices (Weather Station, Revenue Meter, PV Inverter, etc.)
 * - Virtual Firmware
 * - Device Memory
 * - Communication Interfaces
 * 
 * Sections:
 * - Identity: Name, Type, ID
 * - Overview: High-level summary
 * - State: Current operational state
 * - Configuration: Setup parameters
 * - Diagnostics: Health and error information
 * - Communications: Interface statistics
 * - Memory: Device Memory contents
 * - Children: Nested inspectable objects
 */

import { useState } from 'react';
import type { State } from '../types';
import { GenericInspector } from './inspector';

interface InspectorProps {
  state: State;
  selectedNode: string | null;
}

// Legacy component wrapper - forwards to GenericInspector
// This maintains backward compatibility with existing code
export function Inspector({ selectedNode }: InspectorProps) {
  const [currentObjectId, setCurrentObjectId] = useState<string | null>(null);

  // Update current object when selectedNode changes
  const effectiveObjectId = selectedNode || currentObjectId;

  const handleSelectChild = (childId: string) => {
    setCurrentObjectId(childId);
  };

  return (
    <GenericInspector 
      objectId={effectiveObjectId}
      onSelectChild={handleSelectChild}
    />
  );
}
