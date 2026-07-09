// Generic Inspector Types
// These types mirror the backend inspector.GenericInspectorData structure

export type ObjectType = 
  | 'simulation' 
  | 'device' 
  | 'firmware' 
  | 'memory' 
  | 'interface' 
  | 'unknown';

export type SectionID = 
  | 'identity' 
  | 'overview' 
  | 'state' 
  | 'configuration' 
  | 'diagnostics' 
  | 'communications' 
  | 'children' 
  | 'memory';

export type PropertyType = 
  | 'text' 
  | 'number' 
  | 'boolean' 
  | 'status' 
  | 'timestamp' 
  | 'duration' 
  | 'quality' 
  | 'enum' 
  | 'nested' 
  | 'list' 
  | 'angle' 
  | 'percentage';

export type Quality = 0 | 1 | 2 | 3;

export interface ObjectIdentity {
  id: string;
  type: ObjectType;
  name: string;
}

export interface Property {
  name: string;
  value: unknown;
  type: PropertyType;
  unit?: string;
  quality?: Quality;
  precision?: number;
  options?: string[];
  children?: Property[];
  items?: Property[];
  sort_order?: number;
  color_func?: string;
}

export interface ObjectRef {
  id: string;
  type: ObjectType;
  name: string;
  path?: string;
}

export interface Section {
  id: SectionID;
  title: string;
  icon?: string;
  properties?: Property[];
  children?: ObjectRef[];
}

export interface GenericInspectorData {
  object: ObjectIdentity;
  sections: Section[];
}

// Section display metadata
export interface SectionMeta {
  id: SectionID;
  label: string;
  icon: string;
  order: number;
}

export const SECTION_META: Record<SectionID, SectionMeta> = {
  identity: { id: 'identity', label: 'Identity', icon: 'tag', order: 1 },
  overview: { id: 'overview', label: 'Overview', icon: 'eye', order: 2 },
  state: { id: 'state', label: 'State', icon: 'activity', order: 3 },
  configuration: { id: 'configuration', label: 'Configuration', icon: 'settings', order: 4 },
  diagnostics: { id: 'diagnostics', label: 'Diagnostics', icon: 'alert-triangle', order: 5 },
  communications: { id: 'communications', label: 'Communications', icon: 'radio', order: 6 },
  memory: { id: 'memory', label: 'Memory', icon: 'database', order: 7 },
  children: { id: 'children', label: 'Children', icon: 'layers', order: 8 },
};

// Type guards
export function isValidSectionID(id: string): id is SectionID {
  return id in SECTION_META;
}

export function isValidPropertyType(typeStr: string): typeStr is PropertyType {
  const validTypes: PropertyType[] = [
    'text', 'number', 'boolean', 'status', 'timestamp', 
    'duration', 'quality', 'enum', 'nested', 'list', 
    'angle', 'percentage'
  ];
  return validTypes.includes(typeStr as PropertyType);
}
