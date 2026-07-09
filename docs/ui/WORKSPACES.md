# Workspaces

## Philosophy

A **workspace** is a focused environment dedicated to one engineering activity. Forge is not a collection of unrelated pages—it is an engineering workbench composed of specialized workspaces.

### Core Principles

1. **One Primary Responsibility**
   - Each workspace has a clearly defined purpose
   - Avoid overlapping responsibilities between workspaces

2. **Workspace Over Pages**
   - Navigation exposes workspaces, not individual pages
   - Pages belong inside workspaces
   - Future capabilities extend existing workspaces

3. **Clarity Over Convenience**
   - The World workspace never configures protocols
   - The Protocols workspace never edits simulation physics
   - The Data Explorer never becomes the Simulation Inspector

4. **Professional Tool Aesthetic**
   - Resembles VS Code, JetBrains IDEs, Unity/Unreal Editors
   - Left sidebar navigation
   - Content area with workspace-specific tools

---

## Workspace Hierarchy

```
┌─────────────────────────────────────────────────────────────────┐
│ Toolbar                                                           │
├──────────────┬───────────────────────────────────────────────────┤
│ Navigation   │ Workspace Content                                  │
│              │                                                    │
│ Dashboard    │ [Content area specific to each workspace]           │
│ World        │                                                    │
│ Devices      │                                                    │
│ Network      │                                                    │
│ Protocols    │                                                    │
│ Scenarios    │                                                    │
│ Data         │                                                    │
│ Library      │                                                    │
│ Settings     │                                                    │
│ Developer    │                                                    │
└──────────────┴───────────────────────────────────────────────────┘
```

---

## Dashboard

**Purpose:** High-level overview of the current simulation

**Contains:**
- Runtime status (running, paused, stopped)
- Active scenario indicator
- Simulation clock display
- Quick actions (start, stop, reset)
- Recent activity feed
- Resource usage indicators

**NOT a SCADA Dashboard:**
- No real-time gauges or meters
- No control widgets
- No alarming visualizations
- Focus is overview, not monitoring

---

## World

**Purpose:** Explore and configure the simulated physical world

**Contains:**

### Simulation Models
- Clock model
- Sun model
- Weather model
- Grid model
- Wind model
- Cloud model
- Hydrology model
- Market model

### Configuration
- Location settings (latitude, longitude)
- Environmental parameters
- Grid parameters
- Time zone settings

**Responsibilities:**
- Configure physics parameters
- View model states
- Set initial conditions

**Does NOT:**
- Configure devices
- Set up protocols
- Define network topology

---

## Devices

**Purpose:** Instantiate and configure virtual industrial devices

**Contains:**

### Device Instances
- Weather Stations
- Revenue Meters
- PV Inverters
- Relays
- Breakers
- PLCs
- RTUs
- IEDs
- Transformers
- Batteries

### Device Management
- Create/delete devices
- Configure device parameters
- Assign to network
- Set communication protocols
- View device state

**Future:**
- Datasheet-generated devices
- Device cloning
- Bulk configuration

**Responsibilities:**
- Device lifecycle management
- Device parameters
- Device-to-network assignment

---

## Network

**Purpose:** Construct and inspect the industrial network topology

**Contains:**

### Network Objects
- Electrical buses
- Transformers
- Transmission lines
- Feeders
- Switches
- Protection devices
- Communication topology

### Network Visualization
- Single-line diagram
- Network hierarchy view
- Connection mapping
- Impedance/parameter display

**Responsibilities:**
- Network topology
- Electrical connections
- Communication paths
- Network parameters

**Does NOT:**
- Configure device parameters
- Set up protocols
- Edit simulation physics

---

## Protocols

**Purpose:** Configure and monitor protocol interfaces

**Contains:**

### Protocol Adapters
- Modbus (RTU, TCP)
- DNP3 (Serial, TCP)
- IEC 61850
- MQTT
- REST
- Raw Ingest

### Protocol Configuration
- Protocol parameters
- Register mappings
- Polling intervals
- Connection settings

### Protocol Diagnostics
- Message counters
- Error rates
- Connection status
- Protocol monitor

**Responsibilities:**
- Protocol configuration
- Connection management
- Protocol monitoring

---

## Scenarios

**Purpose:** Define and execute simulation events

**Contains:**

### Scenario Definitions
- Time acceleration profiles
- Weather event profiles
- Fault injection sequences
- Scheduled operations
- Operator action sequences
- Commissioning procedures

### Scenario Playback
- Event timeline
- Play/pause/rewind
- Variable speed
- Scenario recording

**Responsibilities:**
- Event sequencing
- Time manipulation
- State injection
- Scenario playback

---

## Data Explorer

**Purpose:** Inspect operational data produced by virtual devices

**Contains:**

### Operational Data
- MMA2 memory view
- Device telemetry
- Register values
- Measurement streams
- Protocol values

### Data Analysis
- Real-time value display
- Value trending (external historian)
- Export capabilities

**Important:**
- Displays operational data from MMA2
- Shows device measurements
- Does NOT inspect internal Simulation Models
- Simulation Models are viewed in World workspace

**Responsibilities:**
- MMA2 memory inspection
- Device telemetry viewing
- Operational data analysis

---

## Library

**Purpose:** Manage reusable engineering assets

**Contains:**

### Device Types
- Manufacturer device definitions
- Device type templates
- Protocol capability definitions

### Configuration Profiles
- Weather profiles
- Grid profiles
- Market profiles

### Components
- PV module libraries
- Battery models
- Relay templates
- Cable specifications

### Documentation
- Manufacturer datasheets
- Specification documents
- User guides

**Responsibilities:**
- Asset management
- Template storage
- Library organization

---

## Settings

**Purpose:** Application-wide configuration

**Contains:**

### Application Settings
- Theme preferences
- Layout customization
- Keyboard shortcuts
- Window behavior

### Runtime Options
- Default tick interval
- Memory limits
- Logging preferences

### Developer Options
- Debug mode
- Performance metrics
- Developer tools visibility

### Plugin Management
- Plugin installation
- Plugin configuration
- Plugin status

---

## Developer

**Purpose:** Debugging and development tools

**Contains:**

### Console
- Runtime output
- Debug messages
- Error logs

### Logs
- Application logs
- System events
- Filter and search

### Events
- Simulation events
- Device events
- Protocol events

### Alarms
- Active alarms
- Alarm acknowledgment
- Alarm history

**Note:** This workspace contains tools, not simulation data. Simulation data belongs in appropriate workspaces (World, Data Explorer, etc.)

---

## Future Workspaces

Reserved for future expansion when justified:

| Workspace | Purpose | Trigger |
|-----------|---------|---------|
| World Editor | Full-screen physics configuration | Complex model management |
| Device Editor | Modal device configuration | Advanced device setup |
| Protocol Monitor | Real-time message inspection | Deep protocol debugging |
| Historian | Historical data storage | Long-term trend analysis |
| Replay | Simulation recording/playback | Incident investigation |
| Analytics | Data analysis tools | Advanced reporting |
| AI Assistant | AI-powered assistance | Natural language queries |
| Documentation | Inline documentation | User guidance |
| Plugin Manager | Plugin development | Plugin ecosystem |

**Rule:** Add new workspaces only when justified by implementation needs. Prefer extending existing workspaces.

---

## Workspace Independence

Each workspace has clear ownership:

| Workspace | Owns | Does NOT Own |
|-----------|------|--------------|
| Dashboard | Overview | Device config, Physics |
| World | Simulation physics | Protocols, Devices |
| Devices | Device instances | Physics, Network topology |
| Network | Electrical/comm topology | Physics, Protocols |
| Protocols | Protocol configuration | Physics, Devices |
| Scenarios | Event sequences | Device config, Physics |
| Data Explorer | MMA2 telemetry | Simulation state |
| Library | Asset templates | Runtime configuration |
| Settings | App preferences | Simulation content |

---

## Navigation Guidelines

### Adding New Capabilities

When adding new features, follow this decision tree:

1. **Does it fit an existing workspace?**
   - Yes → Add to that workspace
   - No → Continue

2. **Is it a natural extension of multiple workspaces?**
   - Yes → Create new workspace
   - No → Find the best fit

3. **Is it justified by real use cases?**
   - Yes → Propose new workspace
   - No → Defer

### Navigation Structure

```
Dashboard
World
  └─ Models
  └─ Configuration
Devices
  └─ Device List
  └─ Add Device
Network
  └─ Topology
  └─ Add Connection
Protocols
  └─ Adapter List
  └─ Add Adapter
Scenarios
  └─ Scenario List
  └─ Playback
Data
  └─ MMA2 View
  └─ Export
Library
  └─ Device Types
  └─ Profiles
Settings
  └─ Application
  └─ Runtime
  └─ Developer
Developer
  └─ Console
  └─ Logs
  └─ Events
  └─ Alarms
```

---

## Long-Term Vision

As Forge grows:

1. **Extensibility** - New capabilities extend existing workspaces
2. **Clarity** - Each workspace has one clear purpose
3. **Discoverability** - Users find features in expected locations
4. **Maintainability** - Developers understand ownership boundaries

A developer should be able to understand the entire application by reading this document before exploring the implementation.

---

*Workspaces define the major functional areas of Forge. New features should naturally extend existing workspaces rather than introducing new navigation items.*
