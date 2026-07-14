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

const ENTITY_ICONS: Record<string, string> = {
  grid: '🔌',
  bus: '⚫',
  breaker: '🔀',
  transformer: '🔄',
  generator: '☀️',
  load: '🏭',
  meter: '📊',
};

const ENTITY_LABELS: Record<string, string> = {
  grid: 'GRID',
  bus: 'BUS',
  breaker: 'CB',
  transformer: 'TX',
  generator: 'PV',
  load: 'LOAD',
  meter: 'MTR',
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

  // Get entity icon
  const getEntityIcon = (type: string) => ENTITY_ICONS[type] || '📦';
  const getEntityLabel = (type: string) => ENTITY_LABELS[type] || type.toUpperCase();

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

            return (
              <g
                key={entity.id}
                className={`${styles.entity} ${isSelected ? styles.selected : ''}`}
                transform={`translate(${entity.position.x}, ${entity.position.y})`}
                onClick={(e) => handleEntityClick(entity.id, e)}
              >
                {/* Background */}
                <rect
                  width={entity.size.width}
                  height={entity.size.height}
                  className={styles.entityBackground}
                />

                {/* Icon */}
                <text
                  x={entity.size.width / 2}
                  y={entity.size.height / 2 - (measurement ? 6 : 0)}
                  className={styles.entityIcon}
                >
                  {getEntityIcon(entity.entity_type)}
                </text>

                {/* Label */}
                <text
                  x={entity.size.width / 2}
                  y={entity.size.height / 2 + 14}
                  className={styles.entityLabel}
                >
                  {getEntityLabel(entity.entity_type)}
                </text>

                {/* Measurement Value - units already included in measurement string */}
                {measurement && (
                  <text
                    x={entity.size.width / 2}
                    y={entity.size.height / 2 + 28}
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
