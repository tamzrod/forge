# Editor Architecture Knowledge Base

## Purpose

This document defines the Forge Editor architecture. The Editor is responsible for editing the simulation world, not executing simulation logic.

## Design Principle

**Separation of Concerns**

```
┌─────────────────────────────────────────────────────────┐
│                    FORGE EDITOR                          │
│  Responsible for editing, visualization, and interaction  │
└─────────────────────────────────────────────────────────┘
                           │
                           │ Edits
                           ▼
┌─────────────────────────────────────────────────────────┐
│                    FORGE RUNTIME                        │
│  Responsible for simulation execution                    │
└─────────────────────────────────────────────────────────┘
```

The Editor and Runtime are completely independent. The Editor reads and writes to the Runtime model but never executes simulation logic.

## Editor Responsibilities

The Editor OWNS:
- Visual representation of entities
- User interaction (drag, drop, select, connect)
- Property editing
- Project management (save, load)
- Layout and positioning
- Canvas rendering

The Editor does NOT own:
- Simulation execution
- Physics calculations
- Entity behavior
- Time progression
- Solver execution

## Core Components

### 1. Editor

The Editor is the root container that manages the editing session.

```go
type Editor struct {
    canvas    *Canvas
    palette   *Palette
    inspector *Inspector
    explorer *ProjectExplorer
    runtime  *Runtime
}
```

**Responsibilities:**
- Initialize editor components
- Manage editor state
- Coordinate between components
- Handle file operations (save, load)

### 2. Canvas

The Canvas is the main editing area where entities are placed and connected.

```go
type Canvas struct {
    viewport     Viewport
    entities     map[string]*CanvasEntity
    connections  []*CanvasConnection
    selection    *Selection
    grid         Grid
}
```

**Responsibilities:**
- Render entities at their positions
- Handle pan and zoom
- Manage entity selection
- Render connections between entities
- Coordinate drag operations

**Interactions:**
| Action | Behavior |
|--------|----------|
| Pan | Middle-click drag or two-finger drag |
| Zoom | Scroll wheel or pinch |
| Select | Click on entity |
| Multi-select | Shift+click or drag rectangle |
| Move | Drag selected entity |
| Connect | Drag from terminal to terminal |

### 3. Palette

The Palette provides entity templates that can be dragged onto the canvas.

```go
type Palette struct {
    categories []PaletteCategory
    items     []PaletteItem
}

type PaletteItem struct {
    ID          string
    name        string
    category    string
    icon        string
    entityType  string
    defaults    map[string]interface{}
}
```

**Categories:**

| Category | Items |
|----------|-------|
| Electrical | Grid, Bus, Breaker, Transformer, VirtualGenerator, VirtualLoad, Meter |
| Environment | Sun, Weather, Wind |
| Simulation | Scenario, Clock |

### 4. Selection

Selection manages which entities are currently selected.

```go
type Selection struct {
    selectedIDs map[string]bool
    anchor     *SelectionAnchor
}

type SelectionAnchor struct {
    entityID string
    offset   Point
}
```

**Behaviors:**
- Single selection: Click to select
- Multi-selection: Shift+click or drag rectangle
- Range selection: Shift+arrow keys
- Select all: Ctrl+A

### 5. Inspector

The Inspector displays and edits properties of selected entities.

```go
type Inspector struct {
    entity    Entity
    sections  []InspectorSection
    isVisible bool
}

type InspectorSection struct {
    title    string
    properties []Property
}
```

**Property Types:**

| Type | Editor |
|------|--------|
| String | Text input |
| Number | Number input with optional unit |
| Boolean | Toggle switch |
| Enum | Dropdown select |
| ReadOnly | Display text |

### 6. Project Explorer

The Project Explorer shows the project structure as a tree.

```go
type ProjectExplorer struct {
    tree       TreeNode
    onSelect   func(nodeID string)
}
```

**Tree Structure:**

```
Project
├── World
│   ├── Topology
│   │   ├── Buses
│   │   ├── Branches
│   │   └── Switches
│   ├── Entities
│   │   ├── Generators
│   │   ├── Loads
│   │   └── Meters
│   └── Scenarios
├── Simulation
│   ├── Clock
│   └── Solver
└── Settings
```

### 7. Connection

Connections represent electrical links between entities.

```go
type Connection struct {
    ID          string
    fromEntity  string
    fromTerminal string
    toEntity    string
    toTerminal  string
    busID       string
}
```

**Connection Validation:**
- Source terminals can only connect to buses
- Destination terminals can only connect to buses
- Through terminals connect to buses on both sides
- Observation terminals can connect to any bus
- Voltage must match (or transformer in between)

### 8. Simulation Controls

Basic controls for running simulations.

```go
type SimulationControls struct {
    isRunning   bool
    isPaused    bool
    speed       float64
    currentTime time.Time
}
```

**Controls:**

| Control | Action |
|---------|--------|
| Run | Start simulation |
| Pause | Pause simulation |
| Reset | Reset to initial state |
| Speed | Adjust simulation speed |
| Time Display | Show current simulation time |

## Layout

The Editor uses a standard engineering CAD layout.

```
┌─────────────────────────────────────────────────────────────┐
│                        TOOLBAR                               │
│  [New] [Open] [Save] | [Run] [Pause] [Reset] | Speed: [1x] │
├────────────┬────────────────────────────────┬───────────────┤
│            │                                │               │
│  PALETTE   │           CANVAS               │   INSPECTOR   │
│            │                                │               │
│  ┌──────┐  │    ┌─────┐                     │   Properties │
│  │Grid  │  │    │  ●  │──┐                  │               │
│  ├──────┤  │    └─────┘  │  ┌─────┐         │   Name: [___] │
│  │Bus   │  │              ├──│     │         │   V: [___]   │
│  ├──────┤  │              │  └─────┘         │   Type: [▼]  │
│  │Breaker│ │              │                  │               │
│  ├──────┤  │              ▼                  │               │
│  │TX    │  │           ┌─────┐               │               │
│  ├──────┤  │           │     │               │               │
│  │Gen   │  │           └─────┘               │               │
│  ├──────┤  │                                │               │
│  │Load  │  │                                │               │
│  └──────┘  │                                │               │
├────────────┴────────────────────────────────┴───────────────┤
│                     PROJECT EXPLORER                        │
│  ▼ Project                                                   │
│    ▼ World                                                   │
│      ▼ Topology                                             │
│        ● 69kV Bus                                           │
│        ● 480V Bus                                           │
└─────────────────────────────────────────────────────────────┘
```

## Entity Representation

Each entity type has a visual representation on the canvas.

### Grid

```
    ┌──────────────┐
    │    GRID      │
    │    69kV      │
    └──────┬───────┘
           │
           ▼ (Terminal)
```

### Bus

```
    ┌──────────────┐
    │     ●        │
    │   69kV       │
    │   PCC Bus    │
    └──────────────┘
```

### Breaker

```
    ┌──────────────┐
    │   ╱    ╲     │
    │  ╱  CB  ╲    │
    │  ╲       ╱   │
    │   ╲     ╱    │
    └────╲───╱─────┘
```

### Transformer

```
    ┌──────────────┐
    │  HV ││ LV    │
    │  69kV││480V  │
    │ ════╪════    │
    │      │       │
    └──────┴───────┘
```

### Virtual Generator

```
    ┌──────────────┐
    │  ☀ SOLAR     │
    │  500 kW      │
    │  [OUTPUT]────│
    └──────────────┘
```

### Virtual Load

```
    ┌──────────────┐
    │  ⚡ FACTORY  │
    │  400 kW      │
    │────────[IN]  │
    └──────────────┘
```

### Meter

```
    ┌──────────────┐
    │  📊 METER    │
    │  PCC         │
    │  ●────[OBS]  │
    └──────────────┘
```

## Connection Types

### Terminal-to-Bus Connection

```
Entity Terminal          Bus
    │                    │
    ▼                    ▼
 [OUTPUT]───────────────[●]────
```

### Visual Connection Lines

```
    Source Bus              Destination Bus
         │                        │
         ▼                        ▼
    ┌─────────┐             ┌─────────┐
    │   ●     │─────────────│    ●    │
    └─────────┘             └─────────┘
         │                        │
         ▼                        ▼
    [OUTPUT]                 [INPUT]
```

## State Management

```go
type EditorState struct {
    project     *Project
    canvas      *CanvasState
    selection   []string
    inspector   *InspectorState
    explorer    *ExplorerState
    simulation  *SimulationState
}

type Project struct {
    id          string
    name        string
    world       *World
    topology    *Topology
    entities    []*Entity
    scenarios   []*Scenario
    metadata    ProjectMetadata
}

type CanvasState struct {
    zoom         float64
    panX         float64
    panY         float64
    gridVisible  bool
    snapToGrid   bool
}
```

## User Flows

### Creating a New Entity

1. User drags item from Palette to Canvas
2. Editor creates entity from template
3. Entity appears at drop position
4. Entity is selected
5. Inspector shows entity properties

### Connecting Entities

1. User clicks on entity terminal
2. User drags connection line
3. Editor highlights valid connection targets
4. User drops on target terminal
5. Editor validates connection using topology rules
6. Connection is rendered

### Editing Properties

1. User selects entity
2. Inspector displays entity properties
3. User modifies property value
4. Editor validates change
5. Editor updates entity model
6. Canvas updates entity visual

### Running Simulation

1. User clicks Run button
2. Editor sends world state to Runtime
3. Runtime executes simulation
4. Runtime sends results back to Editor
5. Editor updates visualizations
6. User can pause, step, or reset

## File Format

Projects are saved as JSON.

```json
{
  "version": "1.0",
  "project": {
    "id": "proj-001",
    "name": "Solar Farm Project",
    "world": {
      "topology": {
        "buses": [...],
        "branches": [...]
      },
      "entities": [...],
      "scenarios": [...]
    },
    "canvas": {
      "entities": [...],
      "connections": [...],
      "view": {...}
    }
  }
}
```

## Glossary

| Term | Definition |
|------|------------|
| Editor | UI application for editing simulation models |
| Canvas | Main editing area where entities are placed |
| Palette | List of entity templates available for dragging |
| Inspector | Panel for editing entity properties |
| Project Explorer | Tree view of project structure |
| Selection | Set of currently selected entities |
| Connection | Visual link between entity terminals |
| Terminal | Connection point on an entity |
| Runtime | Backend that executes simulation |

---

*Last Updated: 2026-07-13*
