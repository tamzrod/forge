# Layout

## Primary Layout

The default Forge layout consists of five main regions:

```
┌─────────────────────────────────────────────────────────────┐
│ Toolbar                                                       │
├──────────────┬──────────────────────┬─────────────────────────┤
│ Navigation   │ World Explorer       │ Inspector               │
│              │                      │                         │
│              │                      │                         │
├──────────────┴──────────────────────┴─────────────────────────┤
│ Console / Logs / Events                                       │
└─────────────────────────────────────────────────────────────┘
```

### Panel Responsibilities

| Panel | Purpose |
|-------|---------|
| **Toolbar** | Global actions, search, user menu |
| **Navigation** | Section navigation (VS Code style) |
| **World Explorer** | Hierarchy tree of simulation objects |
| **Inspector** | Current state of selected object |
| **Console** | Developer output, logs, events |

---

## Toolbar

**Purpose:** Global actions and navigation

**Contents:**
- Application title / logo
- Workspace name
- Global search
- Run/Stop simulation controls
- Settings access
- User menu

**Behavior:**
- Fixed at top
- Always visible
- Minimal height (48-56px)

---

## Navigation

**Purpose:** Section navigation, similar to VS Code

**Structure:**

```
┌─────────────────┐
│ Forge           │
├─────────────────┤
│ 📊 Dashboard    │
├─────────────────┤
│ 🌍 World        │
│   Models        │
│   Devices       │
│   Network       │
├─────────────────┤
│ 📡 Protocols   │
├─────────────────┤
│ 🎬 Scenarios   │
├─────────────────┤
│ 📈 Data        │
├─────────────────┤
│ 📚 Library     │
├─────────────────┤
│ ⚙️ Settings    │
├─────────────────┤
│ 🔧 Developer   │
│   Console       │
│   Logs          │
│   Events        │
│   Alarms        │
└─────────────────┘
```

### Sections

#### Dashboard
Overview of simulation state, recent activity

#### World
Contains the simulation hierarchy:
- Models
  - Clock
  - Sun
  - Weather
  - Grid
  - Wind
  - Cloud
- Devices
  - Weather Station
  - Revenue Meter
  - PV Inverter
  - Battery
  - Relay
- Network
- Scenarios

#### Protocols
Protocol adapter configuration and monitoring

#### Scenarios
Scenario definitions and playback

#### Data
Data explorer for historical analysis (separate from Inspector)

#### Library
Device templates, model configurations

#### Settings
Application preferences

#### Developer
Debugging tools (Console, Logs, Events, Alarms)

**Width:** 200-250px (collapsible to 48px icon strip)

---

## World Explorer

**Purpose:** Hierarchical view of the simulation world

**Content:** Tree view of all simulation objects

**Example Hierarchy:**
```
🌐 Simulation World
├── 📊 Models
│   ├── ⏱️ Clock
│   ├── ☀️ Sun
│   ├── 🌡️ Weather
│   ├── ⚡ Grid
│   ├── 💨 Wind
│   └── ☁️ Cloud
├── 📱 Devices
│   ├── 🌤️ Weather Station
│   ├── 💰 Revenue Meter
│   ├── ☀️ PV Inverter
│   ├── 🔋 Battery
│   └── ⚠️ Relay
├── 🌐 Network
│   ├── Ethernet
│   └── Modbus
└── 🎬 Scenarios
    └── Day-Night Cycle
```

**Behavior:**
- Expandable/collapsible tree
- Single selection
- Icons indicate object type
- Selection populates Inspector
- Drag-and-drop for reordering (future)

**Width:** Resizable, default ~300px

---

## Inspector

**Purpose:** Display current state of the selected object

**IMPORTANT:** The Inspector represents **current state only**. It is NOT:
- A SCADA
- A historian
- A trend viewer

Historical data belongs in Data Explorer.

### Content by Object Type

#### Simulation Model Inspector
```
┌────────────────────────────────────┐
│ Sun Model                      [⚙] │
├────────────────────────────────────┤
│ Overview                             │
│ ──────────────────────────────────  │
│ Type: Sun                            │
│ Location: 40°N, 105°W               │
│ Status: Running                      │
├────────────────────────────────────┤
│ Current State                        │
│ ──────────────────────────────────  │
│ Elevation:    45.2°                 │
│ Azimuth:      180.5°                │
│ Irradiance:   850 W/m²              │
│ DNI:          920 W/m²              │
│ Diffuse:      85 W/m²               │
│ Daytime:      Yes                    │
├────────────────────────────────────┤
│ Configuration                        │
│ ──────────────────────────────────  │
│ Latitude:     40.0                  │
│ Longitude:    -105.0                │
└────────────────────────────────────┘
```

#### Device Inspector
```
┌────────────────────────────────────┐
│ Revenue Meter                  [⚙] │
├────────────────────────────────────┤
│ Overview                             │
│ ──────────────────────────────────  │
│ Type: Revenue Meter                   │
│ Model: ABB Totalflow                 │
│ Address: 192.168.1.100              │
│ Protocol: Modbus TCP                  │
│ Status: Connected                     │
├────────────────────────────────────┤
│ Current Measurements                 │
│ ──────────────────────────────────  │
│ Voltage:    478.5 V                 │
│ Current:    125.3 A                 │
│ Frequency:  60.01 Hz                │
│ Power:      89.2 kW                 │
├────────────────────────────────────┤
│ Counters                             │
│ ──────────────────────────────────  │
│ Energy:     45,234 kWh              │
│ Demand:     92.1 kW                 │
└────────────────────────────────────┘
```

### Tabs
- Overview
- State (current values)
- Configuration
- Properties
- Diagnostics

**Width:** Resizable, default fills remaining space

---

## Console Panel

**Purpose:** Developer debugging and output

**Location:** Bottom of window, collapsible

**Tabs:**

### Console
Standard output from simulation and runtime

### Logs
Timestamped log entries with severity levels

### Events
Simulation events (device connections, faults, etc.)

### Alarms
Active alarms with acknowledgment

**Height:** Collapsible, default ~200px, resizable

---

## Responsive Behavior

### Minimum Window Size
- Width: 1024px
- Height: 768px

### Panel Behavior
- All panels resizable via drag handles
- Panels can be collapsed to minimum
- Layout persists in user preferences
- Reset to default layout option

### Overflow Handling
- Horizontal scroll for wide content
- Vertical scroll for tall lists
- Table columns resizable
- Text truncation with tooltips

---

## Future Workspaces

The layout supports future expansion:

### World Editor
Full-screen editor for simulation configuration

### Device Editor
Modal or side panel for device configuration

### Scenario Editor
Timeline-based scenario definition

### Protocol Monitor
Real-time protocol message viewer

### Data Explorer
Full-featured historian and analytics

---

*Layout should remain stable. New features extend existing panels rather than replacing them.*
