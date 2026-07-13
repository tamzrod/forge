import { useState, useCallback, useEffect } from 'react';
import { Canvas } from './Canvas';
import { Palette } from './Palette';
import { ProjectExplorer } from './ProjectExplorer';
import { SimulationControls } from './SimulationControls';
import { EditorInspector } from './EditorInspector';
import { registryService, type PaletteItem } from '../../services/registry';
import type { ComponentDescriptor } from '../../types/registry';
import type {
  CanvasEntity,
  Project,
  Palette as PaletteType,
  TreeNode,
  Point,
  EntityType,
  EntityCategory,
  SimulationState,
} from '../../types/editor';
import styles from './Editor.module.css';

// Map component ID to entity type (extracts type from "forge-domain:type")
function componentIdToEntityType(componentId: string): string {
  const parts = componentId.split(':');
  return parts[parts.length - 1];
}

interface EditorProps {
  onProjectSave?: (project: Project) => void;
}

export function Editor({ onProjectSave }: EditorProps) {
  // Registry state
  const [components, setComponents] = useState<ComponentDescriptor[]>([]);
  const [paletteData, setPaletteData] = useState<{ categories: { id: string; name: string; icon: string; order: number }[]; items: PaletteItem[] } | null>(null);

  // Project state
  const [project, setProject] = useState<Project>({
    id: 'proj-1',
    name: 'New Project',
    entities: [],
    connections: [],
    canvas: {
      zoom: 1,
      pan_x: 0,
      pan_y: 0,
      grid_visible: true,
      snap_to_grid: true,
      grid_size: 20,
    },
    metadata: {
      created_at: new Date().toISOString(),
      modified_at: new Date().toISOString(),
      author: '',
      description: '',
    },
  });

  // Selection state
  const [selection, setSelection] = useState<string[]>([]);

  // Simulation state
  const [simulationState, setSimulationState] = useState<SimulationState>({
    is_running: false,
    is_paused: false,
    speed: 1,
    current_time: '00:00:00',
  });

  // Explorer tree
  const [explorerTree, setExplorerTree] = useState<TreeNode | null>(null);

  // Load registry data
  useEffect(() => {
    registryService.getPalette().then((data) => {
      setPaletteData(data);
      setComponents(data.items.map((item) => ({ id: item.component_id } as ComponentDescriptor)));
    });
  }, []);

  // Get component descriptor by ID
  const getComponent = useCallback((componentId: string): ComponentDescriptor | null => {
    return components.find((c) => c.id === componentId) || null;
  }, [components]);

  // Build explorer tree
  useEffect(() => {
    if (!paletteData) return;

    const buildTree = (): TreeNode => ({
      id: 'project',
      label: project.name,
      icon: '📁',
      type: 'project',
      children: [
        {
          id: 'world',
          label: 'World',
          icon: '🌍',
          type: 'world',
          children: [
            { id: 'world-clock', label: 'Clock', icon: '⏱️', type: 'clock' },
            { id: 'world-solver', label: 'Solver', icon: '⚙️', type: 'solver' },
          ],
        },
        {
          id: 'topology',
          label: 'Topology',
          icon: '🔗',
          type: 'topology',
          children: [
            {
              id: 'topology-buses',
              label: 'Buses',
              icon: '⚫',
              type: 'category',
              children: project.entities
                .filter((e) => e.entity_type === 'bus')
                .map((e) => ({
                  id: e.id,
                  label: e.name,
                  icon: '⚫',
                  type: 'bus',
                })),
            },
            {
              id: 'topology-branches',
              label: 'Branches',
              icon: '🔗',
              type: 'category',
              children: project.entities
                .filter((e) => ['breaker', 'transformer'].includes(e.entity_type))
                .map((e) => ({
                  id: e.id,
                  label: e.name,
                  icon: e.entity_type === 'breaker' ? '🔀' : '🔄',
                  type: e.entity_type,
                })),
            },
          ],
        },
        {
          id: 'entities',
          label: 'Entities',
          icon: '📦',
          type: 'entities',
          children: [
            {
              id: 'entities-generators',
              label: 'Generators',
              icon: '☀️',
              type: 'category',
              children: project.entities
                .filter((e) => e.entity_type === 'generator')
                .map((e) => ({
                  id: e.id,
                  label: e.name,
                  icon: '☀️',
                  type: 'generator',
                })),
            },
            {
              id: 'entities-loads',
              label: 'Loads',
              icon: '🏭',
              type: 'category',
              children: project.entities
                .filter((e) => e.entity_type === 'load')
                .map((e) => ({
                  id: e.id,
                  label: e.name,
                  icon: '🏭',
                  type: 'load',
                })),
            },
          ],
        },
        {
          id: 'simulation',
          label: 'Simulation',
          icon: '▶️',
          type: 'simulation',
          children: [
            { id: 'sim-controls', label: 'Controls', icon: '🎛️', type: 'controls' },
            { id: 'sim-clock', label: 'Clock', icon: '⏰', type: 'clock' },
          ],
        },
      ],
    });

    setExplorerTree(buildTree());
  }, [project, paletteData]);

  // Handle selection
  const handleSelect = useCallback((ids: string[], additive?: boolean) => {
    if (additive) {
      setSelection((prev) => [...new Set([...prev, ...ids])]);
    } else {
      setSelection(ids);
    }
  }, []);

  // Handle entity move
  const handleMove = useCallback((id: string, position: Point) => {
    setProject((prev) => ({
      ...prev,
      entities: prev.entities.map((e) =>
        e.id === id ? { ...e, position } : e
      ),
    }));
  }, []);

  // Handle drop entity from palette
  const handleDropEntity = useCallback((componentId: string, position: Point) => {
    const component = getComponent(componentId);
    if (!component) return;

    const newEntity: CanvasEntity = {
      id: `entity-${Date.now()}`,
      entity_type: componentIdToEntityType(componentId) as EntityType,
      component_id: componentId,
      name: component.name,
      position: project.canvas.snap_to_grid
        ? {
            x: Math.round(position.x / project.canvas.grid_size) * project.canvas.grid_size,
            y: Math.round(position.y / project.canvas.grid_size) * project.canvas.grid_size,
          }
        : position,
      size: { width: component.width, height: component.height },
      properties: component.properties.reduce((acc, prop) => {
        acc[prop.key] = { value: prop.default, type: prop.type, unit: prop.unit, options: prop.options };
        return acc;
      }, {} as Record<string, { value: unknown; type: string; unit?: string; options?: string[] }>),
    };

    setProject((prev) => ({
      ...prev,
      entities: [...prev.entities, newEntity],
    }));

    setSelection([newEntity.id]);
  }, [project.canvas.snap_to_grid, project.canvas.grid_size, getComponent]);

  // Handle palette drag start
  const handlePaletteDragStart = useCallback((_item: PaletteItem) => {
    // Could track this for analytics
  }, []);

  // Handle explorer select
  const handleExplorerSelect = useCallback((nodeId: string) => {
    const entity = project.entities.find((e) => e.id === nodeId);
    if (entity) {
      setSelection([entity.id]);
    }
  }, [project.entities]);

  // Handle simulation controls
  const handleRun = useCallback(() => {
    setSimulationState((prev) => ({
      ...prev,
      is_running: true,
      is_paused: false,
    }));
  }, []);

  const handlePause = useCallback(() => {
    setSimulationState((prev) => ({
      ...prev,
      is_paused: !prev.is_paused,
    }));
  }, []);

  const handleStop = useCallback(() => {
    setSimulationState((prev) => ({
      ...prev,
      is_running: false,
      is_paused: false,
      current_time: '00:00:00',
    }));
  }, []);

  const handleReset = useCallback(() => {
    setSimulationState((prev) => ({
      ...prev,
      is_running: false,
      is_paused: false,
      current_time: '00:00:00',
    }));
  }, []);

  const handleSpeedChange = useCallback((speed: number) => {
    setSimulationState((prev) => ({
      ...prev,
      speed,
    }));
  }, []);

  // Handle property update
  const handlePropertyUpdate = useCallback((entityId: string, key: string, value: unknown) => {
    setProject((prev) => ({
      ...prev,
      entities: prev.entities.map((e) =>
        e.id === entityId
          ? {
              ...e,
              properties: { ...e.properties, [key]: { ...e.properties[key], value } },
            }
          : e
      ),
    }));
  }, []);

  // Get selected entity for inspector
  const selectedEntity = selection.length === 1
    ? project.entities.find((e) => e.id === selection[0])
    : null;

  // Build palette from registry
  const palette: PaletteType | null = paletteData
    ? {
        categories: paletteData.categories.map((cat) => ({
          id: cat.id as EntityCategory,
          name: cat.name,
          icon: cat.icon,
          order: cat.order,
        })),
        items: paletteData.items.map((item) => ({
          id: item.id,
          name: item.name,
          category: item.category as EntityCategory,
          entity_type: componentIdToEntityType(item.component_id) as EntityType,
          component_id: item.component_id,
          icon: item.icon,
        })),
      }
    : null;

  return (
    <div className={styles.editor}>
      {/* Toolbar */}
      <div className={styles.toolbar}>
        <div className={styles.toolbarGroup}>
          <button title="New Project">📄</button>
          <button title="Open Project">📂</button>
          <button title="Save Project" onClick={() => onProjectSave?.(project)}>💾</button>
        </div>
        <div className={styles.toolbarTitle}>{project.name}</div>
        <div className={styles.toolbarGroup}>
          <button title="Undo">↩</button>
          <button title="Redo">↪</button>
        </div>
      </div>

      {/* Main content */}
      <div className={styles.main}>
        {/* Palette */}
        {palette && <Palette palette={palette} onDragStart={handlePaletteDragStart} />}

        {/* Canvas */}
        <Canvas
          entities={project.entities}
          connections={project.connections}
          canvas={project.canvas}
          selection={selection}
          onSelect={handleSelect}
          onMove={handleMove}
          onDropEntity={handleDropEntity}
        />

        {/* Inspector */}
        <EditorInspector
          entity={selectedEntity}
          onPropertyChange={handlePropertyUpdate}
        />
      </div>

      {/* Project Explorer */}
      <ProjectExplorer
        tree={explorerTree}
        selectedId={selection[0] || null}
        onSelect={handleExplorerSelect}
      />

      {/* Simulation Controls */}
      <SimulationControls
        state={simulationState}
        onRun={handleRun}
        onPause={handlePause}
        onStop={handleStop}
        onReset={handleReset}
        onSpeedChange={handleSpeedChange}
      />
    </div>
  );
}
