// Registry types - matching the Go backend registry types

export type PropertyType = 'string' | 'number' | 'boolean' | 'enum';

export type TerminalRole = 'source' | 'destination' | 'through' | 'observation';

export type TerminalDirection = 'input' | 'output' | 'bidirectional';

export interface PropertyDescriptor {
  key: string;
  label: string;
  type: PropertyType;
  default?: unknown;
  unit?: string;
  min?: number;
  max?: number;
  options?: string[];
  readonly?: boolean;
  required?: boolean;
}

export interface TerminalDescriptor {
  id: string;
  name: string;
  role: TerminalRole;
  voltage?: number;
  direction: TerminalDirection;
}

export interface ComponentDescriptor {
  id: string;
  name: string;
  category: string;
  icon: string;
  description?: string;
  properties: PropertyDescriptor[];
  terminals: TerminalDescriptor[];
  width: number;
  height: number;
  domain: string;
  capabilities?: string[];
}

export interface ComponentCategory {
  id: string;
  name: string;
  icon: string;
  order: number;
  domain: string;
  expandable?: boolean;
}
