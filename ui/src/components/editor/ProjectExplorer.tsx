import { useState, useCallback } from 'react';
import type { TreeNode } from '../../types/editor';
import styles from './ProjectExplorer.module.css';

interface ProjectExplorerProps {
  tree: TreeNode | null;
  selectedId: string | null;
  onSelect: (nodeId: string) => void;
}

export function ProjectExplorer({ tree, selectedId, onSelect }: ProjectExplorerProps) {
  const [expandedNodes, setExpandedNodes] = useState<Set<string>>(new Set(['project']));

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

    return (
      <div key={node.id} className={styles.nodeWrapper}>
        <div
          className={`${styles.node} ${isSelected ? styles.selected : ''}`}
          style={{ paddingLeft: `${level * 16 + 8}px` }}
          onClick={(e) => handleNodeClick(node, e)}
        >
          {hasChildren && (
            <span className={styles.expandIcon}>
              {isExpanded ? '▼' : '▶'}
            </span>
          )}
          {!hasChildren && <span className={styles.expandIcon} />}
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

  return (
    <div className={styles.explorer}>
      <div className={styles.header}>
        <h3>Project Explorer</h3>
      </div>

      <div className={styles.tree}>
        {tree ? renderNode(tree) : (
          <div className={styles.empty}>
            <p>No project loaded</p>
          </div>
        )}
      </div>
    </div>
  );
}
