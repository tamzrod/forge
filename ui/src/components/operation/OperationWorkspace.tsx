import { useState, useCallback } from 'react';
import type { CanvasEntity, TreeNode } from '../../types/editor';
import type { State } from '../../types';
import { PlantExplorer } from './PlantExplorer';
import { SingleLineDiagram } from './SingleLineDiagram';
import { EquipmentDetails } from './EquipmentDetails';
import { AnalysisPanel } from './AnalysisPanel';
import styles from './OperationWorkspace.module.css';

export interface OperationWorkspaceProps {
  entities: CanvasEntity[];
  connections: Array<{
    id: string;
    from_entity: string;
    to_entity: string;
  }>;
  simulationState: State;
  onPropertyChange: (entityId: string, key: string, value: unknown) => void;
  onRun: () => void;
  onPause: () => void;
  onReset: () => void;
  onSpeedChange: (speed: number) => void;
  onScenarioChange: (scenarioId: string) => void;
  scenarios?: Array<{ id: string; name: string }>;
}

type SidebarTab = 'plant' | 'analysis';
type EquipmentTab = 'explain' | 'measurements' | 'identity' | 'status' | 'properties';
type AnalysisTab = 'timeline' | 'events' | 'why';

export function OperationWorkspace({
  entities,
  connections,
  simulationState,
  onPropertyChange,
  onRun,
  onPause,
  onReset,
  onSpeedChange,
  onScenarioChange,
  scenarios = [],
}: OperationWorkspaceProps) {
  // State
  const [selection, setSelection] = useState<string | null>(null);
  const [sidebarTab, setSidebarTab] = useState<SidebarTab>('plant');
  const [equipmentTab, setEquipmentTab] = useState<EquipmentTab>('explain');
  const [analysisTab, setAnalysisTab] = useState<AnalysisTab>('timeline');

  // Selection
  const handleSelect = useCallback((id: string | null) => {
    setSelection(id);
    if (id) {
      setEquipmentTab('explain');
    }
  }, []);

  // Get selected entity
  const selectedEntity = selection
    ? entities.find((e) => e.id === selection)
    : null;

  // Build plant explorer tree from entities
  const buildPlantTree = useCallback((): TreeNode => {
    const buses = entities.filter((e) => e.entity_type === 'bus');
    const generators = entities.filter((e) => e.entity_type === 'generator');
    const transformers = entities.filter((e) => e.entity_type === 'transformer');
    const meters = entities.filter((e) => e.entity_type === 'meter');
    const loads = entities.filter((e) => e.entity_type === 'load');
    const grids = entities.filter((e) => e.entity_type === 'grid');

    return {
      id: 'plant',
      label: 'Utility Solar Farm',
      icon: '☀️',
      type: 'plant',
      expanded: true,
      children: [
        {
          id: 'world',
          label: 'Environment',
          icon: '🌍',
          type: 'environment',
          children: [
            {
              id: 'sun',
              label: 'Sun',
              icon: '🌞',
              type: 'model',
            },
            {
              id: 'weather',
              label: 'Weather',
              icon: '🌤️',
              type: 'model',
            },
          ],
        },
        {
          id: 'grid-section',
          label: 'Grid Connection',
          icon: '🔌',
          type: 'section',
          children: grids.map((e) => ({
            id: e.id,
            label: e.name,
            icon: '🔌',
            type: 'grid',
          })),
        },
        {
          id: 'substation',
          label: 'Substation',
          icon: '⚡',
          type: 'substation',
          children: transformers.map((e) => ({
            id: e.id,
            label: e.name,
            icon: '🔄',
            type: 'transformer',
          })),
        },
        {
          id: 'switchyard',
          label: 'Switchyard',
          icon: '🔀',
          type: 'switchyard',
          children: buses.map((e) => ({
            id: e.id,
            label: e.name,
            icon: '⚫',
            type: 'bus',
          })),
        },
        {
          id: 'arrays',
          label: 'PV Arrays',
          icon: '📦',
          type: 'arrays',
          children: generators.map((e) => ({
            id: e.id,
            label: e.name,
            icon: '☀️',
            type: 'generator',
          })),
        },
        {
          id: 'meters',
          label: 'Revenue Meters',
          icon: '📊',
          type: 'meters',
          children: meters.map((e) => ({
            id: e.id,
            label: e.name,
            icon: '📊',
            type: 'meter',
          })),
        },
        {
          id: 'loads',
          label: 'Station Loads',
          icon: '🏭',
          type: 'loads',
          children: loads.map((e) => ({
            id: e.id,
            label: e.name,
            icon: '🏭',
            type: 'load',
          })),
        },
      ],
    };
  }, [entities]);

  // Determine simulation status
  const isRunning = simulationState.clock.mode === 'Running' && !simulationState.clock.is_paused;
  const isPaused = simulationState.clock.mode === 'Running' && simulationState.clock.is_paused;

  return (
    <div className={styles.workspace}>
      {/* Header */}
      <div className={styles.header}>
        <div className={styles.headerLeft}>
          <span className={styles.logo}>⚡</span>
          <span className={styles.title}>Utility-Scale Solar Farm</span>
          <div className={styles.status}>
            <span
              className={`${styles.statusDot} ${
                isRunning ? styles.running : isPaused ? styles.paused : styles.stopped
              }`}
            />
            <span>
              {isRunning ? 'Running' : isPaused ? 'Paused' : 'Stopped'}
            </span>
          </div>
        </div>
        <div className={styles.headerRight}>
          <button className={styles.headerButton}>⚙️ Settings</button>
        </div>
      </div>

      {/* Main Content */}
      <div className={styles.main}>
        {/* Sidebar - Plant/Analysis */}
        <div className={styles.sidebar}>
          <div className={styles.sidebarTabs}>
            <button
              className={`${styles.sidebarTab} ${sidebarTab === 'plant' ? styles.active : ''}`}
              onClick={() => setSidebarTab('plant')}
            >
              Plant
            </button>
            <button
              className={`${styles.sidebarTab} ${sidebarTab === 'analysis' ? styles.active : ''}`}
              onClick={() => setSidebarTab('analysis')}
            >
              Analysis
            </button>
          </div>
          <div className={styles.sidebarContent}>
            {sidebarTab === 'plant' ? (
              <PlantExplorer
                tree={buildPlantTree()}
                selectedId={selection}
                onSelect={handleSelect}
              />
            ) : (
              <AnalysisPanel
                simulationState={simulationState}
                activeTab={analysisTab}
                onTabChange={(tab) => setAnalysisTab(tab as AnalysisTab)}
              />
            )}
          </div>
        </div>

        {/* Single Line Diagram */}
        <div className={styles.content}>
          <SingleLineDiagram
            entities={entities}
            connections={connections}
            simulationState={simulationState}
            selection={selection}
            onSelect={handleSelect}
          />
        </div>

        {/* Equipment Details */}
        <div className={styles.rightPanel}>
          <div className={styles.panelTabs}>
            <button
              className={`${styles.panelTab} ${equipmentTab === 'explain' ? styles.active : ''}`}
              onClick={() => setEquipmentTab('explain')}
            >
              Explain
            </button>
            <button
              className={`${styles.panelTab} ${equipmentTab === 'measurements' ? styles.active : ''}`}
              onClick={() => setEquipmentTab('measurements')}
            >
              Measurements
            </button>
            <button
              className={`${styles.panelTab} ${equipmentTab === 'identity' ? styles.active : ''}`}
              onClick={() => setEquipmentTab('identity')}
            >
              Identity
            </button>
            <button
              className={`${styles.panelTab} ${equipmentTab === 'status' ? styles.active : ''}`}
              onClick={() => setEquipmentTab('status')}
            >
              Status
            </button>
            <button
              className={`${styles.panelTab} ${equipmentTab === 'properties' ? styles.active : ''}`}
              onClick={() => setEquipmentTab('properties')}
            >
              Properties
            </button>
          </div>
          <div className={styles.panelContent}>
            <EquipmentDetails
              entity={selectedEntity}
              simulationState={simulationState}
              activeTab={equipmentTab}
              onPropertyChange={onPropertyChange}
            />
          </div>
        </div>
      </div>

      {/* Simulation Controls */}
      <div className={styles.simControls}>
        <div className={styles.simButtons}>
          {!isRunning ? (
            <button className={`${styles.simButton} ${styles.run}`} onClick={onRun} title="Run">
              ▶
            </button>
          ) : (
            <button className={`${styles.simButton} ${styles.pause}`} onClick={onPause} title="Pause">
              ⏸
            </button>
          )}
          <button className={`${styles.simButton} ${styles.reset}`} onClick={onReset} title="Reset">
            ↺
          </button>
        </div>

        <div className={styles.simTime}>
          {simulationState.clock.elapsed}
        </div>

        <div className={styles.simSpeed}>
          <label>Speed:</label>
          <select
            value={1}
            onChange={(e) => onSpeedChange(parseFloat(e.target.value))}
          >
            <option value={0.1}>0.1x</option>
            <option value={0.25}>0.25x</option>
            <option value={0.5}>0.5x</option>
            <option value={1}>1x</option>
            <option value={2}>2x</option>
            <option value={4}>4x</option>
            <option value={8}>8x</option>
          </select>
        </div>

        <div className={styles.simScenario}>
          <label>Scenario:</label>
          <select
            value=""
            onChange={(e) => onScenarioChange(e.target.value)}
          >
            <option value="">Select Scenario</option>
            {scenarios.map((scenario) => (
              <option key={scenario.id} value={scenario.id}>
                {scenario.name}
              </option>
            ))}
          </select>
        </div>
      </div>
    </div>
  );
}
