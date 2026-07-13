import { useRef, useState, useCallback } from 'react';
import type {
  CanvasEntity,
  Connection,
  CanvasState,
  Point,
  DragState,
  EntityType
} from '../../types/editor';
import styles from './Canvas.module.css';

interface CanvasProps {
  entities: CanvasEntity[];
  connections: Connection[];
  canvas: CanvasState;
  selection: string[];
  onSelect: (ids: string[], additive?: boolean) => void;
  onMove: (id: string, position: Point) => void;
  onDropEntity: (componentId: string, position: Point) => void;
}

const ENTITY_ICONS: Record<EntityType, string> = {
  grid: '🔌',
  bus: '⚫',
  breaker: '🔀',
  transformer: '🔄',
  generator: '☀️',
  load: '🏭',
  meter: '📊',
  sun: '🌞',
  weather: '🌤️',
  wind: '💨',
  scenario: '🎬',
  clock: '⏱️',
};

const ENTITY_LABELS: Record<EntityType, string> = {
  grid: 'GRID',
  bus: 'BUS',
  breaker: 'CB',
  transformer: 'TX',
  generator: 'GEN',
  load: 'LOAD',
  meter: 'MTR',
  sun: 'SUN',
  weather: 'WX',
  wind: 'WIND',
  scenario: 'SCEN',
  clock: 'CLK',
};

export function Canvas({
  entities,
  connections,
  canvas,
  selection,
  onSelect,
  onMove,
  onDropEntity,
}: CanvasProps) {
  const canvasRef = useRef<HTMLDivElement>(null);
  const [dragState, setDragState] = useState<DragState>({ type: 'none' });
  const [rectangleSelect, setRectangleSelect] = useState<{ start: Point; end: Point } | null>(null);

  // Convert screen coordinates to canvas coordinates
  const screenToCanvas = useCallback(
    (screenX: number, screenY: number): Point => {
      if (!canvasRef.current) return { x: 0, y: 0 };
      const rect = canvasRef.current.getBoundingClientRect();
      return {
        x: (screenX - rect.left - canvas.pan_x) / canvas.zoom,
        y: (screenY - rect.top - canvas.pan_y) / canvas.zoom,
      };
    },
    [canvas.pan_x, canvas.pan_y, canvas.zoom]
  );

  // Handle mouse down
  const handleMouseDown = useCallback(
    (e: React.MouseEvent) => {
      if (e.button === 1 || (e.button === 0 && e.altKey)) {
        // Middle click or Alt+left click - pan
        setDragState({ type: 'pan', startPoint: { x: e.clientX, y: e.clientY } });
      } else if (e.button === 0) {
        const pos = screenToCanvas(e.clientX, e.clientY);

        // Check if clicking on an entity
        const clickedEntity = entities.find((entity) =>
          pos.x >= entity.position.x &&
          pos.x <= entity.position.x + entity.size.width &&
          pos.y >= entity.position.y &&
          pos.y <= entity.position.y + entity.size.height
        );

        if (clickedEntity) {
          if (e.shiftKey) {
            // Additive selection
            onSelect([...selection, clickedEntity.id], true);
          } else if (!selection.includes(clickedEntity.id)) {
            onSelect([clickedEntity.id]);
          }
          // Start move drag
          setDragState({ type: 'move', startPoint: pos, entityId: clickedEntity.id });
        } else {
          // Start rectangle selection
          if (!e.shiftKey) {
            onSelect([]);
          }
          setRectangleSelect({ start: pos, end: pos });
          setDragState({ type: 'select', startPoint: pos });
        }
      }
    },
    [entities, selection, onSelect, screenToCanvas]
  );

  // Handle mouse move
  const handleMouseMove = useCallback(
    (e: React.MouseEvent) => {
      const pos = screenToCanvas(e.clientX, e.clientY);

      if (dragState.type === 'pan' && dragState.startPoint) {
        // Pan logic would update canvas state
      } else if (dragState.type === 'move' && dragState.startPoint && dragState.entityId) {
        const dx = pos.x - dragState.startPoint.x;
        const dy = pos.y - dragState.startPoint.y;
        const entity = entities.find((ent) => ent.id === dragState.entityId);
        if (entity) {
          onMove(entity.id, {
            x: entity.position.x + dx,
            y: entity.position.y + dy,
          });
          setDragState({ ...dragState, startPoint: pos });
        }
      } else if (dragState.type === 'select' && rectangleSelect) {
        setRectangleSelect({ ...rectangleSelect, end: pos });
      }
    },
    [dragState, rectangleSelect, entities, onMove, screenToCanvas]
  );

  // Handle mouse up
  const handleMouseUp = useCallback(
    (_e: React.MouseEvent) => {
      if (dragState.type === 'select' && rectangleSelect) {
        // Find entities within rectangle
        const minX = Math.min(rectangleSelect.start.x, rectangleSelect.end.x);
        const maxX = Math.max(rectangleSelect.start.x, rectangleSelect.end.x);
        const minY = Math.min(rectangleSelect.start.y, rectangleSelect.end.y);
        const maxY = Math.max(rectangleSelect.start.y, rectangleSelect.end.y);

        const selected = entities
          .filter(
            (entity) =>
              entity.position.x >= minX &&
              entity.position.x + entity.size.width <= maxX &&
              entity.position.y >= minY &&
              entity.position.y + entity.size.height <= maxY
          )
          .map((e) => e.id);

        if (selected.length > 0) {
          onSelect(selected);
        }
      }

      setDragState({ type: 'none' });
      setRectangleSelect(null);
    },
    [dragState, rectangleSelect, entities, onSelect]
  );

  // Handle wheel for zoom
  const handleWheel = useCallback((e: React.WheelEvent) => {
    e.preventDefault();
    // Zoom logic would go here - omitted for simplicity
  }, []);

  // Handle drop from palette
  const handleDrop = useCallback(
    (e: React.DragEvent) => {
      e.preventDefault();
      const componentId = e.dataTransfer.getData('component-id');
      if (componentId) {
        const pos = screenToCanvas(e.clientX, e.clientY);
        onDropEntity(componentId, pos);
      }
    },
    [onDropEntity, screenToCanvas]
  );

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault();
  }, []);

  // Get entity icon
  const getEntityIcon = (type: EntityType) => ENTITY_ICONS[type] || '📦';
  const getEntityLabel = (type: EntityType) => ENTITY_LABELS[type] || type.toUpperCase();

  // Render connection line
  const renderConnection = (conn: Connection) => {
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
  };

  // Render rectangle selection
  const renderRectangleSelect = () => {
    if (!rectangleSelect) return null;
    const x = Math.min(rectangleSelect.start.x, rectangleSelect.end.x);
    const y = Math.min(rectangleSelect.start.y, rectangleSelect.end.y);
    const width = Math.abs(rectangleSelect.end.x - rectangleSelect.start.x);
    const height = Math.abs(rectangleSelect.end.y - rectangleSelect.start.y);

    return (
      <rect
        x={x}
        y={y}
        width={width}
        height={height}
        className={styles.rectangleSelect}
      />
    );
  };

  return (
    <div
      ref={canvasRef}
      className={styles.canvas}
      onMouseDown={handleMouseDown}
      onMouseMove={handleMouseMove}
      onMouseUp={handleMouseUp}
      onWheel={handleWheel}
      onDrop={handleDrop}
      onDragOver={handleDragOver}
    >
      <svg
        className={styles.canvasSvg}
        style={{
          transform: `translate(${canvas.pan_x}px, ${canvas.pan_y}px) scale(${canvas.zoom})`,
        }}
      >
        {/* Grid */}
        {canvas.grid_visible && (
          <defs>
            <pattern
              id="grid"
              width={canvas.grid_size}
              height={canvas.grid_size}
              patternUnits="userSpaceOnUse"
            >
              <circle cx="1" cy="1" r="1" fill="#ddd" />
            </pattern>
          </defs>
        )}
        <rect width="100%" height="100%" fill="url(#grid)" />

        {/* Connections */}
        <g className={styles.connections}>
          {connections.map(renderConnection)}
        </g>

        {/* Rectangle selection */}
        {renderRectangleSelect()}

        {/* Entities */}
        <g className={styles.entities}>
          {entities.map((entity) => (
            <g
              key={entity.id}
              transform={`translate(${entity.position.x}, ${entity.position.y})`}
              className={`${styles.entity} ${selection.includes(entity.id) ? styles.selected : ''}`}
            >
              {/* Entity background */}
              <rect
                width={entity.size.width}
                height={entity.size.height}
                className={styles.entityBackground}
                rx={4}
              />

              {/* Entity icon */}
              <text
                x={entity.size.width / 2}
                y={entity.size.height / 2 - 8}
                textAnchor="middle"
                dominantBaseline="middle"
                className={styles.entityIcon}
              >
                {getEntityIcon(entity.entity_type)}
              </text>

              {/* Entity label */}
              <text
                x={entity.size.width / 2}
                y={entity.size.height / 2 + 10}
                textAnchor="middle"
                dominantBaseline="middle"
                className={styles.entityLabel}
              >
                {getEntityLabel(entity.entity_type)}
              </text>

              {/* Terminal (bottom) */}
              <circle
                cx={entity.size.width / 2}
                cy={entity.size.height + 6}
                r={6}
                className={styles.terminal}
              />

              {/* Terminal (top) */}
              <circle
                cx={entity.size.width / 2}
                cy={-6}
                r={6}
                className={styles.terminal}
              />
            </g>
          ))}
        </g>
      </svg>

      {/* Canvas controls */}
      <div className={styles.controls}>
        <button onClick={() => {}} title="Zoom In">+</button>
        <span>{Math.round(canvas.zoom * 100)}%</span>
        <button onClick={() => {}} title="Zoom Out">-</button>
        <button onClick={() => {}} title="Fit">⊡</button>
      </div>
    </div>
  );
}
