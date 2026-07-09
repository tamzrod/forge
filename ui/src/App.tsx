import { useState } from 'react';
import { Toolbar } from './components/Toolbar';
import { Navigation } from './components/Navigation';
import { WorldExplorer } from './components/WorldExplorer';
import { Inspector } from './components/Inspector';
import { useSimulation } from './hooks/useSimulation';
import type { Workspace } from './types';
import styles from './App.module.css';

function App() {
  const [activeWorkspace, setActiveWorkspace] = useState<Workspace>('world');
  const [selectedNode, setSelectedNode] = useState<string | null>('world');
  const [isRunning, setIsRunning] = useState(false);
  const { state, connected } = useSimulation();

  const workspaceNames: Record<Workspace, string> = {
    dashboard: 'Dashboard',
    world: 'World',
    devices: 'Devices',
    network: 'Network',
    protocols: 'Protocols',
    scenarios: 'Scenarios',
    data: 'Data Explorer',
    library: 'Library',
    settings: 'Settings',
    developer: 'Developer',
  };

  const handleToggleRun = () => {
    setIsRunning(!isRunning);
    // In a real implementation, this would control the simulation
  };

  const handleStop = () => {
    setIsRunning(false);
    // In a real implementation, this would stop the simulation
  };

  const handleNavigate = (workspace: Workspace) => {
    setActiveWorkspace(workspace);
  };

  const handleSelectNode = (nodeId: string) => {
    setSelectedNode(nodeId);
  };

  return (
    <div className={styles.app}>
      <Toolbar
        workspaceName={workspaceNames[activeWorkspace]}
        connected={connected}
        isRunning={isRunning}
        onToggleRun={handleToggleRun}
        onStop={handleStop}
      />
      
      <div className={styles.main}>
        <Navigation
          activeWorkspace={activeWorkspace}
          onNavigate={handleNavigate}
        />
        
        <div className={styles.content}>
          {activeWorkspace === 'world' ? (
            <>
              <WorldExplorer
                state={state}
                selectedNode={selectedNode}
                onSelectNode={handleSelectNode}
              />
              <Inspector
                state={state}
                selectedNode={selectedNode}
              />
            </>
          ) : (
            <div className={styles.workspacePlaceholder}>
              <span className={styles.placeholderIcon}>
                {activeWorkspace === 'dashboard' ? '📊' : 
                 activeWorkspace === 'devices' ? '🔌' : 
                 activeWorkspace === 'network' ? '🔗' : 
                 activeWorkspace === 'protocols' ? '📡' : 
                 activeWorkspace === 'scenarios' ? '🎬' : 
                 activeWorkspace === 'data' ? '📊' : 
                 activeWorkspace === 'library' ? '📚' : 
                 activeWorkspace === 'settings' ? '⚙️' : 
                 activeWorkspace === 'developer' ? '💻' : '📁'}
              </span>
              <span className={styles.placeholderText}>
                {workspaceNames[activeWorkspace]}
              </span>
              <span className={styles.placeholderSubtext}>
                Coming soon
              </span>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default App;
