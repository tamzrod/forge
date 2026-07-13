import { useState, useCallback } from 'react';
import type { TreeNode } from '../../types/editor';
import type { State } from '../../types';
import styles from './PlantExplorer.module.css';

interface PlantExplorerProps {
  tree: TreeNode;
  selectedId: string | null;
  onSelect: (nodeId: string | null) => void;
  simulationState?: State;
}

export function PlantExplorer({
  tree,
  selectedId,
  onSelect,
  simulationState,
}: PlantExplorerProps) {
  const [expandedNodes, setExpandedNodes] = useState<Set<string>>(
    new Set(['plant', 'world', 'arrays', 'switchyard'])
  );

  const toggleExpand = useCallback((nodeId: string) => {
    setExpandedNodes((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(nodeId)) {
        newSet.delete(nodeId);
      } else {
        newSet.add(nodeId);
      }
      return newSet;
    });
  }, []);

  const handleNodeClick = useCallback(
    (node: TreeNode, e: React.MouseEvent) => {
      e.stopPropagation();
      onSelect(node.id);

      if (node.children && node.children.length > 0) {
        toggleExpand(node.id);
      }
    },
    [onSelect, toggleExpand]
  );

  const renderNode = (node: TreeNode, level: number = 0): JSX.Element => {
    const hasChildren = node.children && node.children.length > 0;
    const isExpanded = expandedNodes.has(node.id);
    const isSelected = selectedId === node.id;
    const isEnvironmentSection = node.id === 'world';

    return (
      <div
        key={node.id}
        className={`${styles.nodeWrapper} ${isEnvironmentSection ? styles.environmentSection : ''}`}
      >
        <div
          className={`${styles.node} ${isSelected ? styles.selected : ''}`}
          style={{ paddingLeft: `${level * 16 + 8}px` }}
          onClick={(e) => handleNodeClick(node, e)}
        >
          {hasChildren ? (
            <span className={styles.expandIcon}>
              {isExpanded ? '▼' : '▶'}
            </span>
          ) : (
            <span className={styles.expandIcon} />
          )}
          <span className={styles.nodeIcon}>{node.icon || '📄'}</span>
          <span className={styles.nodeLabel}>{node.label}</span>
        </div>

        {hasChildren && isExpanded && (
          <div className={styles.children}>
            {node.children!.map((child) => renderNode(child, level + 1))}
          </div>
        )}
      </div>
    );
  };

  // Render sun status card
  const renderSunCard = () => {
    if (!simulationState) return null;
    const { sun } = simulationState;

    return (
      <div className={styles.sunCard}>
        <div className={styles.sunHeader}>
          <span className={styles.sunIcon}>🌞</span>
          <span className={styles.sunTitle}>Solar Position</span>
          <span className={`${styles.sunStatus} ${!sun.is_daytime ? styles.night : ''}`}>
            {sun.is_daytime ? 'Daytime' : 'Night'}
          </span>
        </div>
        <div className={styles.sunMetrics}>
          <div className={styles.sunMetric}>
            <div className={styles.sunMetricLabel}>Irradiance</div>
            <div className={styles.sunMetricValue}>
              {sun.irradiance.toFixed(0)} W/m²
            </div>
          </div>
          <div className={styles.sunMetric}>
            <div className={styles.sunMetricLabel}>Elevation</div>
            <div className={styles.sunMetricValue}>
              {sun.elevation.toFixed(1)}°
            </div>
          </div>
          <div className={styles.sunMetric}>
            <div className={styles.sunMetricLabel}>Azimuth</div>
            <div className={styles.sunMetricValue}>
              {sun.azimuth.toFixed(1)}°
            </div>
          </div>
          <div className={styles.sunMetric}>
            <div className={styles.sunMetricLabel}>Direct Normal</div>
            <div className={styles.sunMetricValue}>
              {sun.direct_normal.toFixed(0)} W/m²
            </div>
          </div>
        </div>
      </div>
    );
  };

  // Render weather card
  const renderWeatherCard = () => {
    if (!simulationState) return null;
    const { weather } = simulationState;

    return (
      <div className={styles.weatherCard}>
        <div className={styles.weatherHeader}>
          <span className={styles.weatherIcon}>🌤️</span>
          <span className={styles.weatherTitle}>Weather</span>
        </div>
        <div className={styles.weatherMetrics}>
          <div className={styles.weatherMetric}>
            <div className={styles.weatherMetricLabel}>Temperature</div>
            <div className={styles.weatherMetricValue}>
              {weather.temperature.toFixed(1)}°C
            </div>
          </div>
          <div className={styles.weatherMetric}>
            <div className={styles.weatherMetricLabel}>Humidity</div>
            <div className={styles.weatherMetricValue}>
              {weather.humidity.toFixed(0)}%
            </div>
          </div>
          <div className={styles.weatherMetric}>
            <div className={styles.weatherMetricLabel}>Wind Speed</div>
            <div className={styles.weatherMetricValue}>
              {weather.wind_speed.toFixed(1)} m/s
            </div>
          </div>
          <div className={styles.weatherMetric}>
            <div className={styles.weatherMetricLabel}>Cloud Cover</div>
            <div className={styles.weatherMetricValue}>
              {weather.cloud_cover.toFixed(0)}%
            </div>
          </div>
        </div>
      </div>
    );
  };

  return (
    <div className={styles.explorer}>
      {renderSunCard()}
      {renderWeatherCard()}
      <div className={styles.tree}>{renderNode(tree)}</div>
    </div>
  );
}
