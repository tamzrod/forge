import type { SimulationState } from '../../types/editor';
import styles from './SimulationControls.module.css';

interface SimulationControlsProps {
  state: SimulationState;
  onRun: () => void;
  onPause: () => void;
  onStop: () => void;
  onReset: () => void;
  onSpeedChange: (speed: number) => void;
}

export function SimulationControls({
  state,
  onRun,
  onPause,
  onStop,
  onReset,
  onSpeedChange,
}: SimulationControlsProps) {
  const speedOptions = [0.1, 0.25, 0.5, 1, 2, 4, 8];

  return (
    <div className={styles.controls}>
      <div className={styles.buttons}>
        {!state.is_running ? (
          <button
            className={`${styles.button} ${styles.run}`}
            onClick={onRun}
            title="Run"
          >
            ▶
          </button>
        ) : state.is_paused ? (
          <button
            className={`${styles.button} ${styles.run}`}
            onClick={onPause}
            title="Resume"
          >
            ▶
          </button>
        ) : (
          <button
            className={`${styles.button} ${styles.pause}`}
            onClick={onPause}
            title="Pause"
          >
            ⏸
          </button>
        )}

        <button
          className={styles.button}
          onClick={onStop}
          disabled={!state.is_running}
          title="Stop"
        >
          ⏹
        </button>

        <button
          className={styles.button}
          onClick={onReset}
          title="Reset"
        >
          ↺
        </button>
      </div>

      <div className={styles.speed}>
        <label>Speed:</label>
        <select
          value={state.speed}
          onChange={(e) => onSpeedChange(parseFloat(e.target.value))}
        >
          {speedOptions.map((speed) => (
            <option key={speed} value={speed}>
              {speed}x
            </option>
          ))}
        </select>
      </div>

      <div className={styles.time}>
        <span className={styles.timeLabel}>Time:</span>
        <span className={styles.timeValue}>{state.current_time}</span>
      </div>

      <div className={styles.status}>
        <span
          className={`${styles.statusIndicator} ${
            state.is_running ? (state.is_paused ? styles.paused : styles.running) : ''
          }`}
        />
        <span className={styles.statusText}>
          {state.is_running ? (state.is_paused ? 'Paused' : 'Running') : 'Stopped'}
        </span>
      </div>
    </div>
  );
}
