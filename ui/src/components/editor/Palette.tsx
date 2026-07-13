import { useState } from 'react';
import type { Palette as PaletteType, PaletteItem, EntityCategory } from '../../types/editor';
import styles from './Palette.module.css';

interface PaletteProps {
  palette: PaletteType;
  onDragStart: (item: PaletteItem) => void;
}

export function Palette({ palette, onDragStart }: PaletteProps) {
  const [expandedCategories, setExpandedCategories] = useState<Set<EntityCategory>>(
    new Set(['electrical', 'environment', 'simulation'])
  );

  const toggleCategory = (categoryId: EntityCategory) => {
    const newExpanded = new Set(expandedCategories);
    if (newExpanded.has(categoryId)) {
      newExpanded.delete(categoryId);
    } else {
      newExpanded.add(categoryId);
    }
    setExpandedCategories(newExpanded);
  };

  const handleDragStart = (e: React.DragEvent, item: PaletteItem) => {
    e.dataTransfer.setData('component-id', item.component_id);
    e.dataTransfer.effectAllowed = 'copy';
    onDragStart(item);
  };

  const sortedCategories = [...palette.categories].sort((a, b) => a.order - b.order);

  const getCategoryIcon = (category: EntityCategory) => {
    switch (category) {
      case 'electrical':
        return '⚡';
      case 'environment':
        return '🌤️';
      case 'simulation':
        return '🎬';
      default:
        return '📦';
    }
  };

  return (
    <div className={styles.palette}>
      <div className={styles.header}>
        <h3>Palette</h3>
      </div>

      <div className={styles.categories}>
        {sortedCategories.map((category) => (
          <div key={category.id} className={styles.category}>
            <button
              className={styles.categoryHeader}
              onClick={() => toggleCategory(category.id)}
            >
              <span className={styles.categoryIcon}>{getCategoryIcon(category.id)}</span>
              <span className={styles.categoryName}>{category.name}</span>
              <span className={styles.categoryArrow}>
                {expandedCategories.has(category.id) ? '▼' : '▶'}
              </span>
            </button>

            {expandedCategories.has(category.id) && (
              <div className={styles.items}>
                {palette.items
                  .filter((item) => item.category === category.id)
                  .map((item) => (
                    <div
                      key={item.id}
                      className={styles.item}
                      draggable
                      onDragStart={(e) => handleDragStart(e, item)}
                      title={item.description || item.name}
                    >
                      <span className={styles.itemIcon}>{item.icon}</span>
                      <span className={styles.itemName}>{item.name}</span>
                    </div>
                  ))}
              </div>
            )}
          </div>
        ))}
      </div>

      <div className={styles.footer}>
        <p className={styles.hint}>Drag items to canvas</p>
      </div>
    </div>
  );
}
