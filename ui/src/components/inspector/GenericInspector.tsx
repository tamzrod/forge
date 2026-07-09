import { useState, useEffect, useCallback } from 'react';
import {
  Activity,
  Settings,
  AlertTriangle,
  Radio,
  Eye,
  Tag,
  Database,
  Layers
} from 'lucide-react';
import type { 
  GenericInspectorData, 
  Section, 
  SectionID, 
  ObjectRef
} from '../../types/inspector';
import { SectionCard } from './SectionCard';
import { SECTION_META } from '../../types/inspector';
import styles from '../Inspector.module.css';

// Tab type for the Generic Inspector
type TabID = 'identity' | 'overview' | 'state' | 'configuration' | 'diagnostics' | 'communications' | 'memory' | 'children';

// Icon mapping for tabs
const tabIcons: Record<TabID, React.ReactNode> = {
  identity: <Tag size={14} />,
  overview: <Eye size={14} />,
  state: <Activity size={14} />,
  configuration: <Settings size={14} />,
  diagnostics: <AlertTriangle size={14} />,
  communications: <Radio size={14} />,
  memory: <Database size={14} />,
  children: <Layers size={14} />,
};

// Section to tab mapping
const sectionToTab: Record<SectionID, TabID> = {
  identity: 'identity',
  overview: 'overview',
  state: 'state',
  configuration: 'configuration',
  diagnostics: 'diagnostics',
  communications: 'communications',
  memory: 'communications',
  children: 'children',
};

interface GenericInspectorProps {
  objectId: string | null;
  apiBase?: string;
  onSelectChild?: (childId: string) => void;
}

// Loading state component
function LoadingState() {
  return (
    <div className={styles.loadingState}>
      <div className={styles.loadingSpinner} />
      <span>Loading...</span>
    </div>
  );
}

// Error state component
function ErrorState({ message }: { message: string }) {
  return (
    <div className={styles.errorState}>
      <AlertTriangle className={styles.errorIcon} />
      <span>Failed to load inspection data</span>
      <span>{message}</span>
    </div>
  );
}

// Empty state component
function EmptyState() {
  return (
    <div className={styles.emptyState}>
      <Eye size={32} />
      <span>Select an object to view inspection data</span>
    </div>
  );
}

export function GenericInspector({ objectId, apiBase = '/api', onSelectChild }: GenericInspectorProps) {
  const [data, setData] = useState<GenericInspectorData | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<TabID>('overview');

  // Fetch inspection data
  const fetchData = useCallback(async () => {
    if (!objectId) {
      setData(null);
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const response = await fetch(`${apiBase}/inspect/${objectId}`);
      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }
      const inspectionData: GenericInspectorData = await response.json();
      setData(inspectionData);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
      setData(null);
    } finally {
      setLoading(false);
    }
  }, [objectId, apiBase]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  // Get available tabs based on sections present
  const getAvailableTabs = (): TabID[] => {
    if (!data) return [];
    
    const tabs = new Set<TabID>();
    for (const section of data.sections) {
      const tab = sectionToTab[section.id];
      if (tab) {
        tabs.add(tab);
      }
    }
    
    // Ensure overview is always first
    const result: TabID[] = [];
    if (tabs.has('overview')) result.push('overview');
    if (tabs.has('state')) result.push('state');
    if (tabs.has('configuration')) result.push('configuration');
    if (tabs.has('diagnostics')) result.push('diagnostics');
    if (tabs.has('communications')) result.push('communications');
    if (tabs.has('identity')) result.push('identity');
    if (tabs.has('children')) result.push('children');
    
    return result;
  };

  // Get sections for the active tab
  const getSectionsForTab = (tab: TabID): Section[] => {
    if (!data) return [];
    
    return data.sections.filter(section => {
      const sectionTab = sectionToTab[section.id];
      return sectionTab === tab && 
             (section.properties?.length || section.children?.length);
    });
  };

  // Handle child click
  const handleChildClick = (child: ObjectRef) => {
    onSelectChild?.(child.id);
  };

  // Determine title
  const getTitle = () => {
    if (!objectId) return 'Inspector';
    if (!data) return 'Inspector';
    return data.object.name || objectId;
  };

  // Render content based on state
  const renderContent = () => {
    if (!objectId) {
      return <EmptyState />;
    }

    if (loading) {
      return <LoadingState />;
    }

    if (error) {
      return <ErrorState message={error} />;
    }

    if (!data) {
      return <EmptyState />;
    }

    const sections = getSectionsForTab(activeTab);
    
    if (sections.length === 0) {
      return (
        <div className={styles.emptyState}>
          <Eye size={32} />
          <span>No data available for this section</span>
        </div>
      );
    }

    return (
      <>
        {sections.map((section, index) => (
          <SectionCard 
            key={`${section.id}-${index}`} 
            section={section}
            onChildClick={handleChildClick}
          />
        ))}
      </>
    );
  };

  const availableTabs = getAvailableTabs();

  return (
    <div className={styles.inspector}>
      <div className={styles.header}>
        <span className={styles.title}>{getTitle()}</span>
      </div>
      
      {availableTabs.length > 0 && (
        <div className={styles.tabs}>
          {availableTabs.map((tab) => {
            const meta = SECTION_META[tab];
            return (
              <button
                key={tab}
                className={`${styles.tab} ${activeTab === tab ? styles.active : ''}`}
                onClick={() => setActiveTab(tab)}
                title={meta.label}
              >
                {tabIcons[tab]}
                <span className={styles.tabLabel}>{meta.label}</span>
              </button>
            );
          })}
        </div>
      )}

      <div className={styles.content}>
        {renderContent()}
      </div>
    </div>
  );
}
