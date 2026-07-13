import type { CanvasEntity } from '../../types/editor';
import styles from './EditorInspector.module.css';

interface EditorInspectorProps {
  entity: CanvasEntity | null | undefined;
  onPropertyChange: (entityId: string, key: string, value: unknown) => void;
}

export function EditorInspector({ entity, onPropertyChange }: EditorInspectorProps) {
  if (!entity) {
    return (
      <div className={styles.inspector}>
        <div className={styles.header}>
          <h3>Inspector</h3>
        </div>
        <div className={styles.empty}>
          <p>Select an entity to view properties</p>
        </div>
      </div>
    );
  }

  const handleChange = (key: string, value: unknown) => {
    onPropertyChange(entity.id, key, value);
  };

  return (
    <div className={styles.inspector}>
      <div className={styles.header}>
        <h3>{entity.entity_type.toUpperCase()}</h3>
        <span className={styles.name}>{entity.name}</span>
      </div>

      <div className={styles.sections}>
        {/* Identity Section */}
        <div className={styles.section}>
          <h4>Identity</h4>
          <div className={styles.property}>
            <label>Name</label>
            <input
              type="text"
              value={entity.name}
              onChange={(e) => handleChange('name', e.target.value)}
            />
          </div>
          <div className={styles.property}>
            <label>ID</label>
            <input type="text" value={entity.id} readOnly />
          </div>
        </div>

        {/* Properties Section */}
        <div className={styles.section}>
          <h4>Properties</h4>
          {Object.entries(entity.properties).map(([key, prop]) => (
            <div key={key} className={styles.property}>
              <label>
                {key.replace(/_/g, ' ')}
                {prop.unit && <span className={styles.unit}>({prop.unit})</span>}
              </label>
              {prop.type === 'string' && (
                <input
                  type="text"
                  value={String(prop.value)}
                  onChange={(e) => handleChange(key, e.target.value)}
                  readOnly={prop.readonly}
                />
              )}
              {prop.type === 'number' && (
                <input
                  type="number"
                  value={Number(prop.value)}
                  onChange={(e) => handleChange(key, parseFloat(e.target.value))}
                  readOnly={prop.readonly}
                  min={prop.min}
                  max={prop.max}
                />
              )}
              {prop.type === 'boolean' && (
                <input
                  type="checkbox"
                  checked={Boolean(prop.value)}
                  onChange={(e) => handleChange(key, e.target.checked)}
                  disabled={prop.readonly}
                />
              )}
              {prop.type === 'enum' && (
                <select
                  value={String(prop.value)}
                  onChange={(e) => handleChange(key, e.target.value)}
                  disabled={prop.readonly}
                >
                  {prop.options?.map((opt) => (
                    <option key={opt} value={opt}>
                      {opt}
                    </option>
                  ))}
                </select>
              )}
            </div>
          ))}
        </div>

        {/* Position Section */}
        <div className={styles.section}>
          <h4>Position</h4>
          <div className={styles.property}>
            <label>X</label>
            <input type="number" value={Math.round(entity.position.x)} readOnly />
          </div>
          <div className={styles.property}>
            <label>Y</label>
            <input type="number" value={Math.round(entity.position.y)} readOnly />
          </div>
        </div>
      </div>
    </div>
  );
}
