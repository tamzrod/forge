// Editor Types

export interface Point {
  x: number;
  y: number;
}

export interface Size {
  width: number;
  height: number;
}

export type EntityType =
  | 'grid'
  | 'bus'
  | 'breaker'
  | 'transformer'
  | 'generator'
  | 'load'
  | 'meter'
  | 'sun'
  | 'weather'
  | 'wind'
  | 'scenario'
  | 'clock';

export type EntityCategory = 'electrical' | 'environment' | 'simulation';

export interface PropertyValue {
  value: unknown;
  type: string;
  readonly?: boolean;
  unit?: string;
  min?: number;
  max?: number;
  options?: string[];
}

export interface Properties {
  [key: string]: PropertyValue;
}

export interface CanvasEntity {
  id: string;
  entity_type: EntityType;
  component_id?: string;
  name: string;
  position: Point;
  size: Size;
  world_id?: string;
  properties: Properties;
}

export interface Connection {
  id: string;
  from_entity: string;
  from_terminal: string;
  to_entity: string;
  to_terminal: string;
  bus_id?: string;
}

export interface CanvasState {
  zoom: number;
  pan_x: number;
  pan_y: number;
  grid_visible: boolean;
  snap_to_grid: boolean;
  grid_size: number;
}

export interface Selection {
  entity_ids: string[];
  anchor?: string;
}

export interface PaletteItem {
  id: string;
  name: string;
  category: EntityCategory;
  entity_type: EntityType;
  component_id: string;
  icon: string;
  description?: string;
  default_properties?: Properties;
}

export interface PaletteCategory {
  id: EntityCategory;
  name: string;
  icon: string;
  order: number;
}

export interface Palette {
  categories: PaletteCategory[];
  items: PaletteItem[];
}

export interface InspectorSection {
  title: string;
  properties: InspectorProperty[];
}

export interface InspectorProperty {
  key: string;
  label: string;
  value: unknown;
  type: string;
  readonly?: boolean;
  unit?: string;
  min?: number;
  max?: number;
  options?: string[];
}

export interface InspectorState {
  selected_entity?: CanvasEntity;
  sections: InspectorSection[];
}

export interface TreeNode {
  id: string;
  label: string;
  icon?: string;
  type: string;
  children?: TreeNode[];
  data?: unknown;
  expanded?: boolean;
}

export interface ProjectMetadata {
  created_at: string;
  modified_at: string;
  author: string;
  description: string;
}

export interface Project {
  id: string;
  name: string;
  entities: CanvasEntity[];
  connections: Connection[];
  canvas: CanvasState;
  metadata: ProjectMetadata;
}

export interface EditorState {
  project: Project | null;
  selection: Selection;
  inspector: InspectorState;
  is_modified: boolean;
  is_running: boolean;
  is_paused: boolean;
  speed: number;
  current_time: string;
}

export interface SimulationState {
  is_running: boolean;
  is_paused: boolean;
  speed: number;
  current_time: string;
}

// Canvas interaction types
export type DragState = {
  type: 'none' | 'pan' | 'select' | 'move' | 'connect';
  startPoint?: Point;
  currentPoint?: Point;
  entityId?: string;
  fromTerminal?: string;
};

export type ConnectionDraft = {
  fromEntity: string;
  fromTerminal: string;
  currentPoint: Point;
};
