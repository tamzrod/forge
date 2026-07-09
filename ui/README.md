# Forge UI - Engineering Workbench

## Overview

This is the React-based Engineering Workbench UI for Forge, an Industrial Simulation Runtime.

## UI Milestone 1 - Engineering Workbench

This implementation delivers the first usable Engineering Workbench UI following the documented design system.

## Features

### Implemented

- **Application Shell**: Professional IDE-style layout with toolbar, navigation, explorer, and inspector
- **Toolbar**: Global controls, search, connection status, simulation controls
- **Left Navigation**: VS Code-style sidebar with workspace navigation (World, Dashboard, Devices, etc.)
- **World Explorer**: Hierarchical tree view of simulation models (Clock, Sun, Weather, Grid, Wind, Devices)
- **Inspector**: Tabbed property viewer with:
  - Overview: Summary information
  - State: Current values
  - Configuration: Device/model settings
  - Diagnostics: Health and status information
  - Communications: Interface and traffic statistics
- **Property Grid**: Consistent key-value display for all properties
- **Status Badges**: Semantic color-coded status indicators (Healthy, Warning, Fault, Offline)
- **Compass Widgets**: Visual wind direction and sun azimuth displays
- **Dark Theme**: Professional dark-first design with semantic colors
- **Live Updates**: WebSocket-based real-time state updates

### Architecture

```
┌─────────────────────────────────────────────────────────────┐
│ Toolbar: Logo, Workspace Name, Search, Run/Stop, Settings    │
├──────────────┬──────────────────────┬───────────────────────┤
│ Navigation   │ World Explorer       │ Inspector              │
│              │                      │                        │
│ Dashboard    │ Simulation World      │ [Overview]             │
│ World ●      │ ├── Models           │ [State]                │
│ Devices      │ │   ├── Clock        │ [Configuration]       │
│ Network      │ │   ├── Sun          │ [Diagnostics]          │
│ ...          │ │   ├── Weather      │ [Communications]       │
│              │ │   ├── Grid         │                        │
│              │ │   └── Wind        │                        │
│              │ └── Devices          │                        │
│              │     └── Weather Stn  │                        │
├──────────────┴──────────────────────┴───────────────────────┤
└─────────────────────────────────────────────────────────────┘
```

### Data Flow

```
Simulation Runtime
    │
    ├── Clock Model ──────────────────┐
    ├── Sun Model ────────────────────┤
    ├── Weather Model ───────────────┼──► Inspector.View
    ├── Grid Model ──────────────────┤    (via WebSocket)
    └── Devices ─────────────────────┘
         └── Weather Station
              ├── Virtual Firmware
              ├── Device Memory
              └── Communication Interface
```

### Terminology

This implementation uses documented terminology from the Architecture Glossary:

- **Simulation World**: The complete collection of Simulation Models
- **Simulation Models**: Clock, Sun, Weather, Grid, Wind (physics)
- **Virtual Firmware**: Device behavior code
- **Device Memory**: Internal device state
- **Communications**: Interface and traffic statistics
- **World Explorer**: Tree view of simulation objects
- **Inspector**: Property viewer for selected objects

### Semantic Colors

| Color | Meaning | Usage |
|-------|---------|-------|
| Green (#00ff88) | Healthy | Running, Connected |
| Yellow (#ffaa00) | Warning | Transition, Caution |
| Orange (#ff8c00) | Environmental | Temperature, Solar |
| Blue (#00d4ff) | Selection | Active, Navigation |
| Purple (#b388ff) | Engineering | Pressure, Flow |
| Gray (#666666) | Disabled | Offline |
| Red (#ff4444) | Fault | Alarm, Error |

## Getting Started

### Prerequisites

- Node.js 18+
- Go 1.21+ (for backend)
- Running Forge simulation backend

### Development

```bash
cd ui
npm install
npm run dev
```

The UI will be available at http://localhost:5173

### Production Build

```bash
npm run build
```

Output will be in the `dist/` directory.

## Design System

### Colors

All colors follow the documented semantic color system. See `src/styles/theme.css` for the complete color palette.

### Typography

- **Headings**: System UI / Inter (sans-serif)
- **Body**: System UI / Inter (sans-serif)
- **Values/Code**: JetBrains Mono, Fira Code, Consolas (monospace)

### Spacing

Based on 4px grid:
- xs: 4px
- sm: 8px
- md: 12px
- lg: 16px
- xl: 24px
- xxl: 32px

## Known Limitations

1. Only the World workspace is implemented
2. Console panel not implemented (UI Milestone 2)
3. Simulation controls (Run/Stop/Pause) are UI-only (no backend integration)
4. Settings workspace shows placeholder

## Future Milestones

- UI Milestone 2: Console, Simulation Controls
- Dashboard workspace
- Device workspace
- Network workspace
- Responsive layout improvements

## Project Structure

```
ui/
├── src/
│   ├── components/
│   │   ├── Toolbar.tsx         # Top toolbar
│   │   ├── Navigation.tsx       # Left sidebar
│   │   ├── WorldExplorer.tsx    # Tree view
│   │   └── Inspector.tsx       # Property viewer
│   ├── hooks/
│   │   └── useSimulation.ts     # WebSocket state hook
│   ├── styles/
│   │   └── theme.css            # Design system
│   ├── types/
│   │   └── index.ts             # TypeScript types
│   ├── App.tsx                  # Main app component
│   └── main.tsx                 # Entry point
├── public/
│   └── forge.svg                # App icon
├── index.html
├── package.json
├── vite.config.ts
└── tsconfig.json
```

## Validation Checklist

- [x] Navigation works
- [x] Explorer selection updates Inspector
- [x] Weather Station information is visible
- [x] Firmware information is visible
- [x] Device Memory is visible
- [x] Diagnostics are visible
- [x] Communications are visible
- [x] Live updates function correctly (WebSocket)
- [x] Responsive layout works
- [x] Terminology matches documentation
