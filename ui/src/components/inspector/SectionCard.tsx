import { 
  Tag, 
  Eye, 
  Activity, 
  Settings, 
  AlertTriangle, 
  Radio, 
  Database,
  Layers,
  Cpu
} from 'lucide-react';
import type { Section, Property, ObjectRef } from '../../types/inspector';
import { PropertyValue } from './PropertyValue';
import { SECTION_META } from '../../types/inspector';
import styles from '../Inspector.module.css';

interface SectionCardProps {
  section: Section;
  onChildClick?: (child: ObjectRef) => void;
}

// Get icon component by name
function getIconByName(iconName?: string): React.ReactNode {
  switch (iconName) {
    case 'tag':
      return <Tag size={14} />;
    case 'eye':
      return <Eye size={14} />;
    case 'activity':
      return <Activity size={14} />;
    case 'settings':
      return <Settings size={14} />;
    case 'alert-triangle':
      return <AlertTriangle size={14} />;
    case 'radio':
      return <Radio size={14} />;
    case 'database':
      return <Database size={14} />;
    case 'layers':
      return <Layers size={14} />;
    case 'cpu':
      return <Cpu size={14} />;
    default:
      return <Activity size={14} />;
  }
}

// PropertyRow component
interface PropertyRowProps {
  property: Property;
}

function PropertyRow({ property }: PropertyRowProps) {
  return (
    <div className={styles.propertyRow}>
      <span className={styles.propertyLabel}>{property.name}</span>
      <PropertyValue property={property} />
    </div>
  );
}

// NestedSection component for nested properties
interface NestedSectionProps {
  property: Property;
}

function NestedSection({ property }: NestedSectionProps) {
  if (!property.children || property.children.length === 0) {
    return null;
  }

  return (
    <div className={styles.nestedSection}>
      <div className={styles.nestedSectionHeader}>
        {property.name}
      </div>
      <div className={styles.nestedSectionContent}>
        {property.children.map((child, index) => (
          <PropertyRow key={`${child.name}-${index}`} property={child} />
        ))}
      </div>
    </div>
  );
}

// ChildItem component
interface ChildItemProps {
  child: ObjectRef;
  onClick?: () => void;
}

function ChildItem({ child, onClick }: ChildItemProps) {
  const getTypeIcon = () => {
    switch (child.type) {
      case 'simulation':
        return <Activity size={14} />;
      case 'device':
        return <Cpu size={14} />;
      case 'firmware':
        return <Cpu size={14} />;
      case 'memory':
        return <Database size={14} />;
      case 'interface':
        return <Radio size={14} />;
      default:
        return <Eye size={14} />;
    }
  };

  return (
    <div className={styles.childItem} onClick={onClick}>
      <span className={styles.childIcon}>{getTypeIcon()}</span>
      <span className={styles.childName}>{child.name}</span>
      <span className={styles.childType}>{child.type}</span>
      <ChevronRightIcon />
    </div>
  );
}

function ChevronRightIcon() {
  return (
    <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <polyline points="9 18 15 12 9 6" />
    </svg>
  );
}

// SectionCard component
export function SectionCard({ section, onChildClick }: SectionCardProps) {
  const meta = SECTION_META[section.id];
  const icon = getIconByName(section.icon || meta?.icon);

  const handleChildClick = (child: ObjectRef) => {
    onChildClick?.(child);
  };

  return (
    <div className={styles.sectionCard}>
      <div className={styles.sectionHeader}>
        <span className={styles.sectionIcon}>{icon}</span>
        <span className={styles.sectionTitle}>{section.title}</span>
      </div>
      <div className={styles.sectionContent}>
        {/* Render regular properties */}
        {section.properties?.map((property, index) => {
          // Handle nested properties specially
          if (property.type === 'nested' && property.children && property.children.length > 0) {
            return <NestedSection key={`${property.name}-${index}`} property={property} />;
          }
          return <PropertyRow key={`${property.name}-${index}`} property={property} />;
        })}
        
        {/* Render child references */}
        {section.children && section.children.length > 0 && (
          <div className={styles.childList}>
            {section.children.map((child) => (
              <ChildItem 
                key={child.id} 
                child={child} 
                onClick={() => handleChildClick(child)} 
              />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
