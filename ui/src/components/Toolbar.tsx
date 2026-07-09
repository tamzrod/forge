import { Search, Settings, Play, Pause, Square, Zap } from 'lucide-react';
import styles from './Toolbar.module.css';

interface ToolbarProps {
  workspaceName: string;
  connected: boolean;
  isRunning: boolean;
  onToggleRun?: () => void;
  onStop?: () => void;
}

export function Toolbar({ workspaceName, connected, isRunning, onToggleRun, onStop }: ToolbarProps) {
  return (
    <header className={styles.toolbar}>
      <div className={styles.left}>
        <div className={styles.logo}>
          <Zap size={20} className={styles.logoIcon} />
          <span className={styles.logoText}>Forge</span>
        </div>
        <div className={styles.divider} />
        <span className={styles.workspace}>{workspaceName}</span>
      </div>

      <div className={styles.center}>
        <div className={styles.searchContainer}>
          <Search size={16} className={styles.searchIcon} />
          <input
            type="text"
            placeholder="Search..."
            className={styles.searchInput}
          />
        </div>
      </div>

      <div className={styles.right}>
        <div className={styles.controls}>
          {isRunning ? (
            <button className={styles.controlButton} onClick={onToggleRun} title="Pause">
              <Pause size={16} />
            </button>
          ) : (
            <button className={styles.controlButton} onClick={onToggleRun} title="Run">
              <Play size={16} />
            </button>
          )}
          <button className={styles.controlButton} onClick={onStop} title="Stop">
            <Square size={16} />
          </button>
        </div>
        
        <div className={styles.divider} />

        <div className={`${styles.connectionStatus} ${connected ? styles.connected : styles.disconnected}`}>
          <span className={styles.statusDot} />
          <span className={styles.statusText}>{connected ? 'Connected' : 'Disconnected'}</span>
        </div>

        <button className={styles.settingsButton} title="Settings">
          <Settings size={18} />
        </button>
      </div>
    </header>
  );
}
