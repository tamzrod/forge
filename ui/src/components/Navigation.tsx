import { 
  LayoutDashboard, 
  Globe, 
  Cpu, 
  Network, 
  FileCode, 
  Play, 
  Database, 
  Library, 
  Settings, 
  Terminal 
} from 'lucide-react';
import type { Workspace } from '../types';
import styles from './Navigation.module.css';

interface NavigationProps {
  activeWorkspace: Workspace;
  onNavigate: (workspace: Workspace) => void;
}

interface NavItem {
  id: Workspace;
  label: string;
  icon: React.ReactNode;
  disabled?: boolean;
}

const mainNavItems: NavItem[] = [
  { id: 'dashboard', label: 'Dashboard', icon: <LayoutDashboard size={18} /> },
  { id: 'world', label: 'World', icon: <Globe size={18} /> },
  { id: 'devices', label: 'Devices', icon: <Cpu size={18} />, disabled: true },
  { id: 'network', label: 'Network', icon: <Network size={18} />, disabled: true },
];

const secondaryNavItems: NavItem[] = [
  { id: 'protocols', label: 'Protocols', icon: <FileCode size={18} />, disabled: true },
  { id: 'scenarios', label: 'Scenarios', icon: <Play size={18} />, disabled: true },
  { id: 'data', label: 'Data', icon: <Database size={18} />, disabled: true },
  { id: 'library', label: 'Library', icon: <Library size={18} />, disabled: true },
];

const bottomNavItems: NavItem[] = [
  { id: 'settings', label: 'Settings', icon: <Settings size={18} />, disabled: true },
  { id: 'developer', label: 'Developer', icon: <Terminal size={18} />, disabled: true },
];

export function Navigation({ activeWorkspace, onNavigate }: NavigationProps) {
  const renderNavItem = (item: NavItem) => (
    <button
      key={item.id}
      className={`${styles.navItem} ${activeWorkspace === item.id ? styles.active : ''} ${item.disabled ? styles.disabled : ''}`}
      onClick={() => !item.disabled && onNavigate(item.id)}
      disabled={item.disabled}
      title={item.disabled ? `${item.label} (coming soon)` : item.label}
    >
      <span className={styles.icon}>{item.icon}</span>
      <span className={styles.label}>{item.label}</span>
    </button>
  );

  return (
    <nav className={styles.navigation}>
      <div className={styles.section}>
        {mainNavItems.map(renderNavItem)}
      </div>
      
      <div className={styles.divider} />
      
      <div className={styles.section}>
        {secondaryNavItems.map(renderNavItem)}
      </div>

      <div className={styles.spacer} />
      
      <div className={styles.section}>
        {bottomNavItems.map(renderNavItem)}
      </div>
    </nav>
  );
}
