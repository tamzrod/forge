import styles from './WelcomeScreen.module.css';

export interface WelcomeScreenProps {
  onLoadSolarFarm: () => void;
  onOpenRecent?: (projectId: string) => void;
  onOpenExisting?: () => void;
  recentProjects?: Array<{
    id: string;
    name: string;
    path: string;
    lastOpened: string;
  }>;
}

export function WelcomeScreen({
  onLoadSolarFarm,
  onOpenRecent,
  onOpenExisting,
  recentProjects = [],
}: WelcomeScreenProps) {
  const handleRecentClick = (projectId: string) => {
    onOpenRecent?.(projectId);
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    });
  };

  return (
    <div className={styles.welcome}>
      <div className={styles.welcomeContent}>
        <div className={styles.logo}>⚡</div>
        <h1 className={styles.title}>Forge</h1>
        <p className={styles.subtitle}>Utility-Scale Solar Farm Operations</p>

        <div className={styles.actions}>
          {/* Primary Action: Load Solar Farm */}
          <button className={styles.primaryAction} onClick={onLoadSolarFarm}>
            <span>☀️</span>
            <span>Load Utility-Scale Solar Farm</span>
          </button>

          {/* Secondary Actions */}
          <div className={styles.secondaryActions}>
            <button className={styles.secondaryAction} onClick={onOpenExisting}>
              <span>📂</span>
              <span>Open Existing Project</span>
            </button>
          </div>
        </div>

        {/* Recent Projects */}
        {recentProjects.length > 0 && (
          <div className={styles.recentProjects}>
            <h3 className={styles.recentTitle}>Recent Projects</h3>
            <div className={styles.recentList}>
              {recentProjects.map((project) => (
                <div
                  key={project.id}
                  className={styles.recentItem}
                  onClick={() => handleRecentClick(project.id)}
                >
                  <span className={styles.recentIcon}>☀️</span>
                  <div className={styles.recentInfo}>
                    <div className={styles.recentName}>{project.name}</div>
                    <div className={styles.recentDate}>
                      {formatDate(project.lastOpened)}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        <div className={styles.version}>Forge v1.0.0</div>
      </div>
    </div>
  );
}
