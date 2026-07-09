import { useState, useMemo } from 'react';
import { 
  ChevronRight, 
  ChevronDown, 
  Globe,
  Clock,
  Sun,
  Cloud,
  Zap,
  Wind,
  Cpu,
  Thermometer,
  Gauge,
  Radio
} from 'lucide-react';
import type { State, TreeNode } from '../types';
import styles from './WorldExplorer.module.css';

interface WorldExplorerProps {
  state: State;
  selectedNode: string | null;
  onSelectNode: (nodeId: string) => void;
}

interface TreeItemProps {
  node: TreeNode;
  depth: number;
  selectedNode: string | null;
  onSelect: (id: string) => void;
  expandedNodes: Set<string>;
  onToggleExpand: (id: string) => void;
}

const iconMap: Record<string, React.ReactNode> = {
  clock: <Clock size={16} />,
  sun: <Sun size={16} />,
  weather: <Cloud size={16} />,
  grid: <Zap size={16} />,
  wind: <Wind size={16} />,
  'weather-station': <Thermometer size={16} />,
  'revenue-meter': <Gauge size={16} />,
  'pv-inverter': <Zap size={16} />,
  battery: <Cpu size={16} />,
  relay: <Radio size={16} />,
};

function getIconForNode(nodeId: string): React.ReactNode {
  const id = nodeId.toLowerCase().replace(/\s+/g, '-');
  for (const [key, icon] of Object.entries(iconMap)) {
    if (id.includes(key)) {
      return icon;
    }
  }
  return <Globe size={16} />;
}

function TreeItem({ node, depth, selectedNode, onSelect, expandedNodes, onToggleExpand }: TreeItemProps) {
  const hasChildren = node.children && node.children.length > 0;
  const isExpanded = expandedNodes.has(node.id);
  const isSelected = selectedNode === node.id;

  const handleClick = () => {
    onSelect(node.id);
  };

  const handleToggle = (e: React.MouseEvent) => {
    e.stopPropagation();
    if (hasChildren) {
      onToggleExpand(node.id);
    }
  };

  return (
    <div className={styles.treeItemContainer}>
      <div
        className={`${styles.treeItem} ${isSelected ? styles.selected : ''}`}
        style={{ paddingLeft: `${depth * 16 + 8}px` }}
        onClick={handleClick}
      >
        <span className={styles.expandHandle} onClick={handleToggle}>
          {hasChildren ? (
            isExpanded ? <ChevronDown size={14} /> : <ChevronRight size={14} />
          ) : (
            <span className={styles.expandPlaceholder} />
          )}
        </span>
        <span className={styles.icon}>{getIconForNode(node.id)}</span>
        <span className={styles.label}>{node.label}</span>
      </div>
      {hasChildren && isExpanded && (
        <div className={styles.children}>
          {node.children!.map((child) => (
            <TreeItem
              key={child.id}
              node={child}
              depth={depth + 1}
              selectedNode={selectedNode}
              onSelect={onSelect}
              expandedNodes={expandedNodes}
              onToggleExpand={onToggleExpand}
            />
          ))}
        </div>
      )}
    </div>
  );
}

export function WorldExplorer({ state, selectedNode, onSelectNode }: WorldExplorerProps) {
  const [expandedNodes, setExpandedNodes] = useState<Set<string>>(new Set(['world', 'models', 'devices']));

  const treeData: TreeNode = useMemo(() => {
    const weatherStationDevices = state.devices.devices.filter(d => 
      d.type.toLowerCase().includes('weather')
    );
    const otherDevices = state.devices.devices.filter(d => 
      !d.type.toLowerCase().includes('weather')
    );

    return {
      id: 'world',
      label: 'Simulation World',
      type: 'root',
      icon: 'globe',
      children: [
        {
          id: 'models',
          label: 'Models',
          type: 'models',
          children: [
            { id: 'clock', label: 'Clock', type: 'model' },
            { id: 'sun', label: 'Sun', type: 'model' },
            { id: 'weather', label: 'Weather', type: 'model' },
            { id: 'grid', label: 'Grid', type: 'model' },
            { id: 'wind', label: 'Wind', type: 'model' },
          ],
        },
        {
          id: 'devices',
          label: `Devices (${state.devices.count})`,
          type: 'device',
          children: [
            ...weatherStationDevices.map(d => ({
              id: `device-${d.id}`,
              label: d.name || d.type,
              type: 'device' as const,
              data: d,
            })),
            ...otherDevices.map(d => ({
              id: `device-${d.id}`,
              label: d.name || d.type,
              type: 'device' as const,
              data: d,
            })),
          ],
        },
      ],
    };
  }, [state.devices]);

  const handleToggleExpand = (nodeId: string) => {
    setExpandedNodes((prev) => {
      const next = new Set(prev);
      if (next.has(nodeId)) {
        next.delete(nodeId);
      } else {
        next.add(nodeId);
      }
      return next;
    });
  };

  return (
    <div className={styles.explorer}>
      <div className={styles.header}>
        <Globe size={16} className={styles.headerIcon} />
        <span className={styles.headerTitle}>World Explorer</span>
      </div>
      <div className={styles.treeContainer}>
        <TreeItem
          node={treeData}
          depth={0}
          selectedNode={selectedNode}
          onSelect={onSelectNode}
          expandedNodes={expandedNodes}
          onToggleExpand={handleToggleExpand}
        />
      </div>
    </div>
  );
}
