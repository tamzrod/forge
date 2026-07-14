import { useRef, useState, useCallback, useMemo } from 'react';
import type { CanvasEntity } from '../../types/editor';
import type { State } from '../../types';
import { computePVArrayMeasurement, computeBusMeasurement, computeMeterMeasurement, computeLoadMeasurement, computeTransformerMeasurement } from '../../services/simulation';
import styles from './SingleLineDiagram.module.css';

interface SingleLineDiagramProps {
  entities: CanvasEntity[];
  connections: Array<{ id: string; from_entity: string; to_entity: string }>;
  simulationState: State;
  selection: string | null;
  onSelect: (id: string | null) => void;
}

interface PanState {
  isPanning: boolean;
  startX: number;
  startY: number;
  panX: number;
  panY: number;
}

// IEEE Std 315 / ANSI standard SVG symbols for electrical equipment
const IEEE_SYMBOLS: Record<string, {
  path: string;
  viewBox: string;
  typeLabel: string;
}> = {
  // Grid: Circle with G (utility source)
  grid: {
    path: 'M -25,-25 L 25,-25 L 25,25 L -25,25 Z M -18,-18 L 18,-18 L 18,18 L -18,18 Z',
    viewBox: '-30 -30 60 60',
    typeLabel: 'GND'
  },
  // Bus: Horizontal bar (common connection point)
  bus: {
    path: 'M -35,0 L 35,0',
    viewBox: '-40 -10 80 20',
    typeLabel: 'BUS'
  },
  // Circuit Breaker: Rectangle with X (ANSI style)
  breaker: {
    path: 'M -25,-20 L 25,-20 L 25,20 L -25,20 Z M -20,-15 L 20,-15 L 20,15 L -20,15 Z M -8,-8 L 8,8 M 8,-8 L -8,8',
    viewBox: '-30 -25 60 50',
    typeLabel: 'CB'
  },
  // Transformer: Circle with ratio notation
  transformer: {
    path: 'M 0,-25 A 25,25 0 1,1 0,25 A 25,25 0 1,1 0,-25 M -15,-15 A 15,15 0 1,0 -15,15 A 15,15 0 1,0 -15,-15',
    viewBox: '-30 -30 60 60',
    typeLabel: 'TX'
  },
  // Generator/PV Array: Circle with arrow (power source)
  generator: {
    path: 'M 0,-25 A 25,25 0 1,1 0,25 A 25,25 0 1,1 0,-25 M 0,-12 L 12,0 L 0,12 Z M -15,-15 L 15,15',
    viewBox: '-30 -30 60 60',
    typeLabel: 'GEN'
  },
  // Load: Filled rectangle with resistance lines
  load: {
    path: 'M -25,-15 L 25,-15 L 25,15 L -25,15 Z M -20,-5 L -10,-5 L -15,5 L -5,5 L -10,-5 L 0,5 L 5,-5 L 15,5 L 10,-5 L 20,5',
    viewBox: '-30 -20 60 40',
    typeLabel: 'LOAD'
  },
  // Meter/Revenue Meter: Circle with M
  meter: {
    path: 'M 0,-25 A 25,25 0 1,1 0,25 A 25,25 0 1,1 0,-25',
    viewBox: '-30 -30 60 60',
    typeLabel: 'MTR'
  }
};

// Get equipment designation from entity ID (P0-4: Equipment designations)
const getEquipmentDesignation = (entity: CanvasEntity): string => {
  const id = entity.id;
  
  // Extract designation patterns
  if (id.includes('grid')) return 'PCC';
  if (id.includes('meter')) return 'MTR';
  if (id.includes('tx') || id.includes('transformer')) return 'TX-1';
  if (id.includes('bus')) return 'BUS-A';
  if (id.includes('breaker')) {
    const match = id.match(/breaker[-_]?(\d+)/i);
    return `CB-${match ? match[1].padStart(3, '0') : '001'}`;
  }
  if (id.includes('pv') || id.includes('generator')) {
    const match = id.match(/[-_]?(\d+)$/i);
    return `PV-${match ? match[1].padStart(2, '0') : '01'}`;
  }
  if (id.includes('load')) return 'LOAD-1';
  
  return id.toUpperCase().slice(0, 6);
};

export function SingleLineDiagram({
  entities,
  connections,
  simulationState,
  selection,
  onSelect,
}: SingleLineDiagramProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const [zoom, setZoom] = useState(1);
  const [pan, setPan] = useState({ x: 0, y: 0 });
  const [panState, setPanState] = useState<PanState>({
    isPanning: false,
    startX: 0,
    startY: 0,
    panX: 0,
    panY: 0,
  });

  // Handle zoom
  const handleZoom = useCallback((delta: number) => {
    setZoom((prev) => {
      const newZoom = Math.max(0.25, Math.min(4, prev + delta));
      return newZoom;
    });
  }, []);

  // Handle mouse down for panning
  const handleMouseDown = useCallback(
    (e: React.MouseEvent) => {
      if (e.button === 1 || (e.button === 0 && e.altKey)) {
        e.preventDefault();
        setPanState({
          isPanning: true,
          startX: e.clientX - pan.x,
          startY: e.clientY - pan.y,
          panX: pan.x,
          panY: pan.y,
        });
      } else if (e.button === 0 && e.target === containerRef.current) {
        // Click on empty canvas - deselect
        onSelect(null);
      }
    },
    [pan, onSelect]
  );

  // Handle mouse move for panning
  const handleMouseMove = useCallback(
    (e: React.MouseEvent) => {
      if (panState.isPanning) {
        setPan({
          x: e.clientX - panState.startX,
          y: e.clientY - panState.startY,
        });
      }
    },
    [panState]
  );

  // Handle mouse up
  const handleMouseUp = useCallback(() => {
    setPanState((prev) => ({
      ...prev,
      isPanning: false,
    }));
  }, []);

  // Handle wheel for zoom
  const handleWheel = useCallback((e: React.WheelEvent) => {
    e.preventDefault();
    const delta = e.deltaY > 0 ? -0.1 : 0.1;
    handleZoom(delta);
  }, [handleZoom]);

  // Handle entity click
  const handleEntityClick = useCallback(
    (entityId: string, e: React.MouseEvent) => {
      e.stopPropagation();
      onSelect(entityId);
    },
    [onSelect]
  );

  // Handle fit to view
  const handleFit = useCallback(() => {
    if (entities.length === 0) return;

    const minX = Math.min(...entities.map((e) => e.position.x));
    const maxX = Math.max(...entities.map((e) => e.position.x + e.size.width));
    const minY = Math.min(...entities.map((e) => e.position.y));
    const maxY = Math.max(...entities.map((e) => e.position.y + e.size.height));

    const contentWidth = maxX - minX + 100;
    const contentHeight = maxY - minY + 100;

    const container = containerRef.current;
    if (!container) return;

    const containerWidth = container.clientWidth;
    const containerHeight = container.clientHeight;

    const scaleX = containerWidth / contentWidth;
    const scaleY = containerHeight / contentHeight;
    const newZoom = Math.min(scaleX, scaleY, 1);

    const centerX = (minX + maxX) / 2;
    const centerY = (minY + maxY) / 2;

    setZoom(newZoom);
    setPan({
      x: containerWidth / 2 - centerX * newZoom,
      y: containerHeight / 2 - centerY * newZoom,
    });
  }, [entities]);

  // Calculate viewBox
  const viewBox = useMemo(() => {
    const container = containerRef.current;
    if (!container) return { width: 1000, height: 600 };

    return {
      width: container.clientWidth,
      height: container.clientHeight,
    };
  }, []);

  // Get power measurement for entity based on type - all values from simulation
  const getEntityMeasurement = (entity: CanvasEntity): string | null => {
    const { sun, weather, grid } = simulationState;
    
    switch (entity.entity_type) {
      case 'generator':
        // Get PV measurement from simulation service
        const pv = computePVArrayMeasurement(entity, sun, weather);
        return pv.ac_power.toFixed(1);
      case 'meter':
        // Get meter measurement from simulation service
        const meter = computeMeterMeasurement(entity, grid, 0, 0, 0, 0);
        return `${meter.active_power.toFixed(0)} kW`;
      case 'bus':
        const bus = computeBusMeasurement(entity, grid);
        return `${bus.voltage.toFixed(0)} V`;
      case 'load':
        const load = computeLoadMeasurement(entity);
        return `${load.active_power.toFixed(0)} kW`;
      case 'transformer':
        const tx = computeTransformerMeasurement(entity, grid, 0, 0);
        return `${tx.load_percent.toFixed(0)}%`;
      default:
        return null;
    }
  };

  // Get IEEE symbol for entity type
  const getIEEESymbol = (type: string) => {
    return IEEE_SYMBOLS[type] || {
      path: 'M -20,-20 L 20,-20 L 20,20 L -20,20 Z',
      viewBox: '-25 -25 50 50',
      typeLabel: type.toUpperCase().slice(0, 4)
    };
  };

  return (
    <div
      ref={containerRef}
      className={styles.diagram}
      onMouseDown={handleMouseDown}
      onMouseMove={handleMouseMove}
      onMouseUp={handleMouseUp}
      onMouseLeave={handleMouseUp}
      onWheel={handleWheel}
    >
      <svg
        className={styles.svg}
        viewBox={`0 0 ${viewBox.width} ${viewBox.height}`}
        style={{
          transform: `translate(${pan.x}px, ${pan.y}px) scale(${zoom})`,
          transformOrigin: '0 0',
        }}
      >
        <defs>
          <pattern
            id="sldGrid"
            width="40"
            height="40"
            patternUnits="userSpaceOnUse"
          >
            <circle cx="20" cy="20" r="1" fill="#2d2d2d" />
          </pattern>

          <marker
            id="arrowGreen"
            viewBox="0 0 10 10"
            refX="9"
            refY="5"
            markerWidth="6"
            markerHeight="6"
            orient="auto-start-reverse"
          >
            <path d="M 0 0 L 10 5 L 0 10 z" fill="#4caf50" />
          </marker>

          <marker
            id="arrowOrange"
            viewBox="0 0 10 10"
            refX="9"
            refY="5"
            markerWidth="6"
            markerHeight="6"
            orient="auto-start-reverse"
          >
            <path d="M 0 0 L 10 5 L 0 10 z" fill="#ff9800" />
          </marker>
        </defs>

        {/* Grid Background */}
        <rect width="10000" height="10000" fill="url(#sldGrid)" x="-5000" y="-5000" />

        {/* Connections */}
        <g className={styles.connections}>
          {connections.map((conn) => {
            const fromEntity = entities.find((e) => e.id === conn.from_entity);
            const toEntity = entities.find((e) => e.id === conn.to_entity);
            if (!fromEntity || !toEntity) return null;

            const fromX = fromEntity.position.x + fromEntity.size.width / 2;
            const fromY = fromEntity.position.y + fromEntity.size.height;
            const toX = toEntity.position.x + toEntity.size.width / 2;
            const toY = toEntity.position.y;

            return (
              <line
                key={conn.id}
                x1={fromX}
                y1={fromY}
                x2={toX}
                y2={toY}
                className={styles.connection}
              />
            );
          })}
        </g>

        {/* Entities */}
        <g className={styles.entities}>
          {entities.map((entity) => {
            const isSelected = selection === entity.id;
            const measurement = getEntityMeasurement(entity);
            const ieeeSymbol = getIEEESymbol(entity.entity_type);
            const designation = getEquipmentDesignation(entity);
            const isBus = entity.entity_type === 'bus';

            return (
              <g
                key={entity.id}
                className={`${styles.entity} ${isSelected ? styles.selected : ''}`}
                transform={`translate(${entity.position.x}, ${entity.position.y})`}
                onClick={(e) => handleEntityClick(entity.id, e)}
              >
                {/* IEEE Symbol */}
                <svg
                  x={isBus ? entity.size.width / 2 - 35 : 5}
                  y={isBus ? entity.size.height / 2 - 5 : 5}
                  width={isBus ? 70 : entity.size.width - 10}
                  height={isBus ? 10 : entity.size.height - 10}
                  viewBox={ieeeSymbol.viewBox}
                  className={styles.ieeeSymbol}
                >
                  <path
                    d={ieeeSymbol.path}
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="2"
                    className={styles.symbolPath}
                  />
                </svg>

                {/* Equipment Designation (P0-4) */}
                <text
                  x={entity.size.width / 2}
                  y={isBus ? entity.size.height / 2 - 12 : entity.size.height / 2 - 20}
                  className={styles.designation}
                >
                  {designation}
                </text>

                {/* Entity Name */}
                <text
                  x={entity.size.width / 2}
                  y={isBus ? entity.size.height / 2 + 8 : entity.size.height / 2}
                  className={styles.entityName}
                >
                  {entity.name.split(' ').slice(0, 3).join(' ')}
                </text>

                {/* Measurement Value */}
                {measurement && (
                  <text
                    x={entity.size.width / 2}
                    y={isBus ? entity.size.height / 2 + 22 : entity.size.height / 2 + 20}
                    className={styles.entityValue}
                  >
                    {measurement}
                  </text>
                )}

                {/* Terminal (top) */}
                <circle
                  cx={entity.size.width / 2}
                  cy={-6}
                  r={6}
                  className={styles.terminal}
                />

                {/* Terminal (bottom) */}
                <circle
                  cx={entity.size.width / 2}
                  cy={entity.size.height + 6}
                  r={6}
                  className={styles.terminal}
                />
              </g>
            );
          })}
        </g>
      </svg>

      {/* Controls */}
      <div className={styles.controls}>
        <button
          className={styles.controlButton}
          onClick={() => handleZoom(0.25)}
          title="Zoom In"
        >
          +
        </button>
        <span className={styles.zoomLevel}>{Math.round(zoom * 100)}%</span>
        <button
          className={styles.controlButton}
          onClick={() => handleZoom(-0.25)}
          title="Zoom Out"
        >
          −
        </button>
        <button
          className={styles.controlButton}
          onClick={handleFit}
          title="Fit to View"
        >
          ⊡
        </button>
      </div>
    </div>
  );
}
