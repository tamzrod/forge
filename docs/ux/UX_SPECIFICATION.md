# UX Specification: Solar Farm Reference World Engineering Workbench

**Document ID:** UX-SPEC-SOLAR-001  
**Version:** 1.0  
**Date:** 2026-07-13  
**Status:** Draft  

---

## 1. Overview

### 1.1 Purpose

This document defines the user experience for the **Solar Farm Reference World Engineering Workbench** — a tool for designing, commissioning, and operating virtual utility-scale solar farms.

### 1.2 What This Is NOT

- **NOT a Dashboard** — Real-time status boards are for operators watching a running plant
- **NOT SCADA** — Supervisory control systems are for plant operators
- **NOT Node-RED** — Flow-based programming is for integration logic
- **NOT a Dashboard** — This is for engineers designing and commissioning the plant

### 1.3 What This IS

An **Engineering Workbench** where the user feels like they are:

1. **Designing** the electrical topology of a solar farm
2. **Configuring** equipment with real specifications
3. **Commissioning** the plant by running simulations and verifying behavior
4. **Debugging** issues by tracing signals and observing measurements

The user should feel like they are working with a real solar farm, not a toy simulation.

---

## 2. User Profile

### 2.1 Primary User: Electrical Engineer

**Background:**
- Understands single-line diagrams
- Knows what a revenue meter is
- Understands PV inverter operating modes
- Can read protection coordination curves
- Familiar with utility interconnection requirements

**Goals:**
- Design a solar farm electrical topology
- Verify power flow under different conditions
- Commission protective devices
- Trace faults through the system
- Verify control system behavior

### 2.2 Secondary User: Software Developer

**Background:**
- Writes SCADA or control software
- Tests integration with Modbus/OPC-UA
- Validates telemetry systems

**Goals:**
- Develop against realistic device models
- Test fault scenarios
- Verify protocol implementations

---

## 3. First Screen

### 3.1 Welcome View

When the user launches Forge, they see:

```
+------------------------------------------------------------------------------+
|                                                                              |
|                          +------------------+                                 |
|                          |                  |                                 |
|                          |      FORGE       |                                 |
|                          |   Solar Farm     |                                 |
|                          |   Workbench      |                                 |
|                          |                  |                                 |
|                          +------------------+                                 |
|                                                                              |
|                      +-----------------------------+                         |
|                      |   [New Solar Farm Project]  |                         |
|                      +-----------------------------+                         |
|                                                                              |
|                      +-----------------------------+                         |
|                      |   [Open Existing Project]   |                         |
|                      +-----------------------------+                         |
|                                                                              |
|                      +-----------------------------+                         |
|                      |   [Load Reference World]    |                         |
|                      +-----------------------------+                         |
|                                                                              |
|  +--------------------------------------------------------------------------+  |
|  | RECENT PROJECTS                                                          |  |
|  | -----------------------------------------------------------------------  |  |
|  |                                                                          |  |
|  | Project A     Last modified: 2026-07-10                                 |  |
|  | Project B     Last modified: 2026-07-09                                 |  |
|  |                                                                          |  |
|  +--------------------------------------------------------------------------+  |
|                                                                              |
+------------------------------------------------------------------------------+
```

### 3.2 Design Rationale

The first screen establishes context:
- **New Solar Farm Project** — User creates a blank plant from scratch
- **Open Existing Project** — User loads a saved project
- **Load Reference World** — User loads a pre-configured reference solar farm

The Reference World is the primary workflow. New projects start from reference templates.

---

## 4. Primary Workflow

### 4.1 The Solar Farm Commissioning Workflow

An electrical engineer commissioning a solar farm follows this sequence:

```
+------------------------------------------------------------------------------+
|                                                                              |
|  PHASE 1: DESIGN          PHASE 2: BUILD           PHASE 3: COMMISSION      |
|  --------------------        ----------------          --------------------   |
|                                                                              |
|  1. Create project         5. Add equipment        9. Run simulation        |
|  2. Configure site         6. Connect equipment     10. Observe measurements  |
|  3. Set up environment    7. Configure settings     11. Inject faults        |
|  4. Define grid connection 8. Validate topology     12. Verify behavior       |
|                                                                              |
+------------------------------------------------------------------------------+
```

### 4.2 Phase 1: Design

**Goal:** Define the physical parameters of the solar farm site.

1. **Create Project**
   - Project name
   - Location (lat/lon for sun position)
   - Plant capacity (MW)
   - Grid connection voltage

2. **Configure Site**
   - Array orientation (azimuth)
   - Tracker type (fixed, single-axis)
   - Meteorological station location

3. **Set Up Environment**
   - Sun model configuration
   - Weather model (clear, cloudy, storm)
   - Time of day / date

4. **Define Grid Connection**
   - Utility grid parameters
   - Point of common coupling (PCC)
   - Interconnection requirements

### 4.3 Phase 2: Build

**Goal:** Assemble the electrical topology.

5. **Add Equipment**
   - PV arrays
   - Inverters
   - Transformers
   - Switchgear
   - Meters
   - Protection relays

6. **Connect Equipment**
   - Wire components in single-line diagram
   - Assign to voltage levels
   - Route collection circuits

7. **Configure Settings**
   - Inverter setpoints
   - Protection settings
   - Meter configuration
   - Protocol bindings

8. **Validate Topology**
   - Check for open circuits
   - Verify voltage compatibility
   - Confirm protection coordination

### 4.4 Phase 3: Commission

**Goal:** Verify the plant operates correctly under simulation.

9. **Run Simulation**
   - Start/pause/stop
   - Adjust time scale
   - Step through events

10. **Observe Measurements**
    - Real-time values on canvas
    - Detailed readings in inspector
    - Power flow visualization

11. **Inject Faults**
    - Open breaker
    - Apply grid disturbance
    - Simulate equipment fault

12. **Verify Behavior**
    - Check protection operation
    - Validate control response
    - Confirm telemetry accuracy

---

## 5. Core Interactions

### 5.1 How Does a User Create a Project?

**Action:** Click "New Solar Farm Project" from welcome screen

**Flow:**
```
1. User clicks [New Solar Farm Project]
        |
        v
2. Dialog appears: "New Solar Farm Project"
   - Name: [________________]  (default: "Solar Farm 1")
   - Location: [Select from map or enter lat/lon]
   - Capacity: [___] MW
   - Grid Voltage: [34.5 kV v]
        |
        v
3. User fills form and clicks [Create]
        |
        v
4. Editor opens with blank canvas
   - Site configured
   - Grid model created
   - Sun/Weather models created
        |
        v
5. Project appears in Recent Projects list
```

### 5.2 How Does a User Build a Plant?

**Action:** Drag equipment from palette to canvas

**Flow:**
```
1. User views palette (left panel)
   +------------------+
   | v Substation     |
   |   Grid Feed      |
   |   Main Breaker   |
   |   Transformer    |
   | v Collection     |
   |   Feeder Breaker |
   |   Combiner       |
   | v Generation     |
   |   PV Array       |
   |   Inverter       |
   | v Protection     |
   |   Relay          |
   |   Meter          |
   | v Environment    |
   |   Weather Station|
   +------------------+
        |
        v
2. User drags "Inverter" to canvas
        |
        v
3. Inverter appears at drop location
        |
        v
4. Inspector shows inverter properties:
   - Name: [Inverter 1]
   - Rating: [500] kW
   - AC Voltage: [480] V
   - Efficiency: [98.0] %
   - Operating Mode: [Grid-Follow v]
        |
        v
5. User configures properties
        |
        v
6. User drags another inverter
        |
        v
7. Repeat until plant is complete
```

### 5.3 How Does a User Connect Equipment?

**Action:** Wire terminals between components

**Flow:**
```
1. User hovers over inverter
   - Terminal points appear (highlighted circles)
   +-----------------+
   |   INVERTER      |
   |   500 kW        |
   |   *------------<- DC Input (left terminal)
   |   ---------------> AC Output (right terminal)
   +-----------------+
        |
        v
2. User clicks DC Input terminal
        |
        v
3. User drags wire to PV Array terminal
   - Wire follows cursor
   - Valid targets highlight
   - Invalid targets gray out
        |
        v
4. User releases on PV Array terminal
        |
        v
5. Connection established
   - Wire appears on canvas
   - Connection validated (compatible voltages)
        |
        v
6. User continues connecting:
   - Inverter AC Output -> Transformer LV
   - Transformer HV -> Main Breaker
   - Main Breaker -> Grid Feed
```

### 5.4 How Does a User Inspect Equipment?

**Action:** Select component on canvas or explorer

**Flow:**
```
1. User clicks on Inverter in canvas
        |
        v
2. Inspector panel updates (right side)
   +----------------------------------------+
   | v Overview                             |
   |   Type: PV Inverter                   |
   |   Name: INV-001                        |
   |   Status: ONLINE                       |
   +----------------------------------------+
   | v Measurements                         |
   |   DC Voltage:    600 V     ^           |
   |   DC Current:    833 A     ^           |
   |   DC Power:      500 kW    ^           |
   |   AC Voltage:    480 V     *           |
   |   AC Current:    833 A     *           |
   |   AC Power:      490 kW    *           |
   |   Efficiency:    98.0 %    *           |
   |   Temperature:   45.2 C   *           |
   +----------------------------------------+
   | v Configuration                        |
   |   Rated Power:   500 kW                 |
   |   Max Power:     510 kW                 |
   |   Min Power:     0 kW                   |
   |   Mode:          MPPT                   |
   |   Grid Mode:     Grid-Forming           |
   +----------------------------------------+
   | v Diagnostics                          |
   |   Operating Hours: 1,234 h              |
   |   Grid Events:     0                    |
   |   Fault Count:      0                    |
   +----------------------------------------+
        |
        v
3. Values update in real-time during simulation
        |
        v
4. User can modify Configuration values during edit mode
```

### 5.5 How Does a User Run the Simulation?

**Action:** Control simulation from toolbar or keyboard

**Controls:**
```
+-------------------------------------------------------------------------+
|  [>] Run]  [|| Pause]  [o Stop]  [O Reset]    Speed: [1x v]           |
|                                      Time: 12:34:56                     |
+-------------------------------------------------------------------------+
```

**Flow:**
```
1. User clicks [>] Run]
        |
        v
2. Simulation starts
   - Status changes to "RUNNING"
   - Clock starts advancing
   - All models begin ticking
   - Devices sample models and update memory
        |
        v
3. User can:
   - Click [|| Pause] to freeze time
   - Click [o Stop] to halt and reset
   - Change Speed to 10x for fast-forward
        |
        v
4. User clicks [|| Pause]
        |
        v
5. Simulation pauses
   - All values frozen
   - User can inspect state
   - User can modify configuration
        |
        v
6. User clicks [>] Run] again
        |
        v
7. Simulation resumes from paused time
```

### 5.6 How Does a User Observe Measurements?

**Action:** View real-time values on canvas or in inspector

**Canvas Display:**
```
During simulation, values appear on canvas:

        +--------------+
        |   INVERTER   |      +--------------+
        |   INV-001    |      | TRANSFORMER  |
        |   500 kW  ^  |------|   TX-001     |
        |   98.0%  *  |      |   2.5 MVA    |
        +--------------+      +--------------+

   ^ = Value increasing
   v = Value decreasing
   * = Value stable
```

**Multi-Value Display:**
```
User can enable floating measurements:

+-----------------------------------------+
|          PLANT SUMMARY                  |
+-----------------------------------------+
|  Grid Import:     450 kW    v         |
|  Grid Export:     0 kW                |
|  Array Power:     500 kW    ^         |
|  Plant Output:    490 kW    ^         |
|  Capacity Factor: 98.0 %     *        |
|  Irradiance:      950 W/m^  ^         |
|  Ambient Temp:    28.5 C     *        |
+-----------------------------------------+
```

### 5.7 How Does a User Replay Events?

**Action:** Use timeline to scrub through simulation

**Timeline View:**
```
+-----------------------------------------------------------------------------+
|  [< <] [<] [>] [> >]                                    [Event Log]         |
|  ----------------------------------------------------------------------    |
|  06:00    08:00    10:00    12:00    14:00    16:00    18:00    20:00     |
|     |        |        |        |        |        |        |        |        |
|     *--------*--------*--------*--------*--------*--------*--------*        |
|                                              A                             |
|                                         Current Position                    |
+-----------------------------------------------------------------------------+
|  [Sunrise]  [Peak]  [Sunset]  [Night]  [+ Add Marker]                     |
+-----------------------------------------------------------------------------+
```

**Flow:**
```
1. User clicks [>] on timeline
        |
        v
2. Simulation plays from current time
        |
        v
3. User clicks [||] on timeline
        |
        v
4. Simulation pauses
        |
        v
5. User drags timeline slider to 10:00
        |
        v
6. Simulation state jumps to 10:00
        |
        v
7. All measurements reflect 10:00 state
        |
        v
8. User clicks [Event Log] button
        |
        v
9. Event Log panel opens:
   +----------------------------------------+
   | EVENT LOG                               |
   +----------------------------------------+
   | 08:30:00  [INFO]  Inverter INV-001     |
   |                     connected to grid   |
   | 09:15:00  [INFO]  Array power reached  |
   |                     rated capacity      |
   | 10:00:00  [WARN]  Grid voltage         |
   |                     dropped 2%          |
   | 10:00:00  [DEBUG] Meas: V=678V        |
   |                          I=737A         |
   +----------------------------------------+
```

### 5.8 How Does a User Debug the Simulation?

**Action:** Use debug tools to trace signal flow

**Debug Panel:**
```
+-----------------------------------------------------------------------------+
|  DEBUG TOOLS                                                               |
+-----------------------------------------------------------------------------+
|                                                                             |
|  [* Signal Trace]  [= Data Watch]  [! Breakpoints]  [= Memory View]       |
|                                                                             |
+-----------------------------------------------------------------------------+
|  SIGNAL TRACE                                                              |
|  -----------------------------------------------------------------------    |
|  Path: Sun -> PV Array -> Inverter -> Transformer -> Main Breaker -> Grid    |
|                                                                             |
|  Step 1: Sun Model                                                         |
|          irradiance = 950 W/m^                                              |
|                 v                                                           |
|  Step 2: PV Array (5,000 m^ @ 18%)                                        |
|          dc_power = 855 kW                                                  |
|                 v                                                           |
|  Step 3: Inverter (98% eff)                                                |
|          ac_power = 838 kW                                                  |
|                 v                                                           |
|  Step 4: Transformer (98.5% eff)                                           |
|          hv_power = 825 kW                                                  |
|                 v                                                           |
|  Step 5: Main Breaker (closed)                                              |
|          output = 825 kW                                                    |
|                 v                                                           |
|  Step 6: Grid                                                              |
|          Grid receives 825 kW                                                |
|                                                                             |
+-----------------------------------------------------------------------------+
```

**Flow:**
```
1. User clicks [* Signal Trace] in debug panel
        |
        v
2. User clicks on "Inverter" in canvas
        |
        v
3. System traces signal path backward:
   - Inverter output -> Transformer input
   - Transformer output -> Breaker input
   - Breaker output -> Grid
        |
        v
4. Trace appears in panel showing signal flow
        |
        v
5. User clicks [= Data Watch]
        |
        v
6. User adds variables to watch list:
   - INV-001.DC_Voltage
   - INV-001.DC_Current
   - INV-001.AC_Power
        |
        v
7. Values update in real-time during simulation
        |
        v
8. User sets breakpoint on "Inverter Tick"
        |
        v
9. Simulation pauses when inverter ticks
   - User can inspect all variables
   - User can step through tick logic
```

---

## 6. Screen Inventory

### 6.1 Welcome Screen

| Element | Description |
|---------|-------------|
| Logo | Forge branding |
| New Project Button | Creates new solar farm project |
| Open Project Button | Opens existing project file |
| Load Reference World | Loads pre-configured reference |
| Recent Projects List | Shows last 5 opened projects |

### 6.2 Project Setup Dialog

| Element | Description |
|---------|-------------|
| Project Name Field | Text input for project name |
| Location Selector | Map or lat/lon input |
| Capacity Input | MW rating |
| Grid Voltage Dropdown | 34.5kV, 69kV, 138kV options |
| Create Button | Creates project |

### 6.3 Main Editor

| Panel | Purpose |
|-------|---------|
| Menu Bar | File, Edit, View, Simulation, Debug, Help |
| Toolbar | Quick actions, simulation controls |
| Palette | Equipment library (left) |
| Canvas | Wiring diagram (center) |
| Inspector | Selected component details (right) |
| Explorer | Project hierarchy (bottom-left) |
| Timeline | Simulation time control (bottom) |
| Console | Debug output, logs |

### 6.4 Canvas Views

| View | Description |
|------|-------------|
| Single-Line Diagram | Primary electrical schematic |
| Physical Layout | Geographic plant view |
| Equipment Rack | Cabinet/inverter lineup view |

### 6.5 Inspector Tabs

| Tab | Content |
|-----|---------|
| Overview | Name, type, status |
| Measurements | Real-time values |
| Configuration | User-editable settings |
| Diagnostics | Operating hours, events |
| Protocols | Modbus, DNP3 bindings |
| Datasheet | Equipment specifications |

### 6.6 Debug Views

| View | Description |
|------|-------------|
| Signal Trace | Path from source to measurement |
| Data Watch | Live variable values |
| Memory View | Raw device memory |
| Breakpoints | Tick breakpoints |

### 6.7 Event Log

| Column | Description |
|--------|-------------|
| Time | Simulation timestamp |
| Level | INFO, WARN, ERROR, DEBUG |
| Source | Device or model |
| Message | Event description |
| Details | Expanded event data |

---

## 7. Navigation Flow

### 7.1 Top-Level Navigation

```
+------------------------------------------------------------------------------+
|  File  |  Edit  |  View  |  Simulation  |  Debug  |  Help                   |
+------------------------------------------------------------------------------+
```

### 7.2 Menu Hierarchy

**File**
- New Project (Ctrl+N)
- Open Project (Ctrl+O)
- Save Project (Ctrl+S)
- Save As (Ctrl+Shift+S)
- Export (SVG, JSON, PDF)
- Recent Projects
- Exit (Alt+F4)

**Edit**
- Undo (Ctrl+Z)
- Redo (Ctrl+Y)
- Cut (Ctrl+X)
- Copy (Ctrl+C)
- Paste (Ctrl+V)
- Delete (Del)
- Select All (Ctrl+A)
- Find (Ctrl+F)

**View**
- Palette (Ctrl+1)
- Inspector (Ctrl+2)
- Explorer (Ctrl+3)
- Console (Ctrl+4)
- Timeline (Ctrl+5)
- Zoom In (Ctrl++)
- Zoom Out (Ctrl+-)
- Fit to Window (Ctrl+0)
- Toggle Full Screen (F11)

**Simulation**
- Run (F5)
- Pause (F6)
- Stop (Shift+F5)
- Reset (Shift+F6)
- Step Forward (F7)
- Speed (submenu)
- Jump to Time...

**Debug**
- Signal Trace
- Data Watch
- Breakpoints
- Memory View
- Event Log
- Performance Monitor

**Help**
- Documentation
- Keyboard Shortcuts
- About Forge

### 7.3 Panel Navigation

```
+------------------------------------------------------------------------------+
|                              TOOLBAR                                         |
+------------------------------------------------------------------------------+
|          |                                                  |               |
|  PALETTE |              CANVAS                              |  INSPECTOR     |
|          |                                                  |               |
|  [Elec]  |  +---------+     +---------+                   |  [Overview]   |
|  [Env]   |  |   PV    |-----|   INV   |                   |  [Measure]    |
|  [Sim]   |  |  Array  |     |  500kW  |                   |  [Config]     |
|  [Prot]  |  +---------+     +----+----+                   |  [Diag]       |
|          |                        |                        |               |
+----------+------------------------+------------------------+---------------+
|                                                                             |
|  EXPLORER                                         TIMELINE                 |
|  ---------                                        -------                   |
|  v Project                                       [>] 12:00:00  [1x v]     |
|    v Solar Farm                                                                |
|      v Substation                                                         |
|        Grid Feed                                                           |
|        Main Breaker                                                        |
|        Transformer                                                         |
|      v Collection                                                          |
|        Feeder 1                                                            |
|        Feeder 2                                                            |
|      v Generation                                                          |
|        Array 1                                                             |
|        INV-001                                                             |
|                                                                             |
+------------------------------------------------------------------------------+
|  CONSOLE                                                                  |
|  -------                                                                  |
|  [Events] [Logs] [Debug]                                                   |
|                                                                             |
+------------------------------------------------------------------------------+
```

---

## 8. Wireframes

### 8.1 Welcome Screen Wireframe

```
+------------------------------------------------------------------------------+
|                                                                              |
|                                                                              |
|                              +------------------+                           |
|                              |                  |                           |
|                              |      FORGE       |                           |
|                              |   Solar Farm     |                           |
|                              |   Workbench      |                           |
|                              |                  |                           |
|                              +------------------+                           |
|                                                                              |
|                         +-----------------------------+                     |
|                         |     + New Solar Farm +      |                     |
|                         +-----------------------------+                     |
|                                                                              |
|                         +-----------------------------+                     |
|                         |   + Open Existing Project +  |                     |
|                         +-----------------------------+                     |
|                                                                              |
|                         +-----------------------------+                     |
|                         |   + Load Reference World  +  |                     |
|                         +-----------------------------+                     |
|                                                                              |
|  +--------------------------------------------------------------------------+  |
|  | RECENT PROJECTS                                                          |  |
|  | -----------------------------------------------------------------------  |  |
|  |                                                                          |  |
|  |  [Solar Farm Alpha]           Last modified: 2026-07-10                 |  |
|  |  [Utility-Scale 100MW]        Last modified: 2026-07-09                 |  |
|  |  [Commercial Rooftop]         Last modified: 2026-07-08                 |  |
|  |                                                                          |  |
|  +--------------------------------------------------------------------------+  |
|                                                                              |
+------------------------------------------------------------------------------+
```

### 8.2 Main Editor Wireframe

```
+------------------------------------------------------------------------------+
| File  Edit  View  Simulation  Debug  Help              [Forge] [User]     |
+------------------------------------------------------------------------------+
|  [New] [Open] [Save]  |  [>] [||] [o] [O]  |  Speed: [1x v]  |  12:34:56  |
+------------------------------------------------------------------------------+
|        |                                            |                        |
|        |                                            |  +------------------+  |
|        |                                            |  | INSPECTOR        |  |
|        |     +-----------+     +-----------+       |  +------------------+  |
|        |     |    PV     |     |   INV     |       |  | [Overview]       |  |
|        |     |  ARRAY    |---->|  500kW   |       |  | [Measurements]   |  |
|        |     | 5,000 m2 |     |  98.0%   |       |  | [Configuration]  |  |
|        |     +-----------+     +-----+-----+       |  | [Diagnostics]    |  |
|        |                             |              |  +------------------+  |
|        |                             |              |  |                   |  |
|        |                             v              |  |  Name: INV-001    |  |
|        |                       +-----------+        |  |  Status: ONLINE   |  |
|        |                       |   TXFRM   |        |  |  ----------------  |  |
|        |     +-----------+     |   2.5MVA  |        |  |  DC Power: 500kW  |  |
|        |     |   GRID    |---->|  98.5%    |        |  |  AC Power: 490kW  |  |
|        |     |   69kV   |     +-----------+        |  |  Efficiency: 98%  |  |
|        |     +-----------+                         |  |  ----------------  |  |
|        |                                            |  |  DC Voltage: 600V |  |
|        |                                            |  |  DC Current: 833A |  |
|  PAL   |                                            |  |                   |  |
| ETTE   |                                            |  |  [Apply Changes]  |  |
|        |                                            |  +------------------+  |
|  [v]   +--------------------------------------------+                        |
|  Elec  |                                                                     |
|  [v]   |                                                                     |
|  Env   +-------------------------------------------------------------------+ |
|  [v]   | EXPLORER                                    | TIMELINE           | |
|  Sim   +--------------------------------------------+--------------------+ |
|        | v Project                                  |                    | |
|        |   v Solar Farm Alpha                        | [< <] [<] [>] [> >] | |
|        |     v Substation                           |                    | |
|        |       Grid Feed                            | 06:00  12:00  18:00| |
|        |       Main Breaker                          |  |      |      |    | |
|        |       Transformer                          |  *-----*------*    | |
|        |     v Collection                           |                    | |
|        |       Feeder 1                             | [Markers: Sun Peak Sunset] | |
|        |       Feeder 2                             |                    | |
|        |     v Generation                           +--------------------+ |
|        |       PV Array 1                           |                    |
|        |       INV-001                              +--------------------+ |
|        |                                             | CONSOLE            | |
|        |                                             +--------------------+ |
|        |                                             | [Events][Logs][Dbg] | |
|        |                                             |                    | |
|        |                                             | [INFO]  12:00:00   | |
|        |                                             | Simulation started  | |
|        |                                             |                    | |
+--------+---------------------------------------------+--------------------+ |
```

### 8.3 Inspector Wireframe

```
+------------------------------------------------------------------------------+
|                                                                              |
|  +------------------------------------------------------------------------------+
|  | INSPECTOR                                                [Pin] [Close]     |
|  +------------------------------------------------------------------------------+
|  |                                                                              |
|  |  +------------------------------------------------------------------------+  |
|  |  | [Overview] [Measurements] [Configuration] [Diagnostics] [Protocols]   |  |
|  |  +------------------------------------------------------------------------+  |
|  |                                                                              |
|  |  OVERVIEW                                          +--------------------+  |
|  |  ----------------------------------------------    |  Visual Status     |  |
|  |                                                    |  +--------------+  |  |
|  |  Name:        [INV-001____________]               |  |     Sun       |  |  |
|  |  Type:        PV Inverter                           |  |   RUNNING     |  |  |
|  |  Manufacturer: [SMA______________]                   |  |   490 kW      |  |  |
|  |  Model:       [Sunny Central________]               |  +--------------+  |  |
|  |  Serial:      [SC5000-2024-001___]                  +--------------------+  |
|  |  Status:      [* Online__________]                                       |  |
|  |                                                                              |  |
|  |  MEASUREMENTS                                                              |  |
|  |  ----------------------------------------------------------------------------  |
|  |                                                                              |  |
|  |  DC Side                              AC Side                               |  |
|  |  ------                              -----                               |  |
|  |  Voltage:    600.0 V  v              Voltage:    480.0 V    *              |  |
|  |  Current:    833.3 A  ^              Current:    833.3 A    *              |  |
|  |  Power:      500.0 kW ^              Power:      490.0 kW   *              |  |
|  |                                                  Efficiency: 98.0%   *      |  |
|  |                                                                              |  |
|  |  Environment                           Grid                                  |  |
|  |  ----------                           ----                                  |  |
|  |  Irradiance:  950 W/m^  ^             Grid Freq:  60.00 Hz    *            |  |
|  |  Cell Temp:   45.2 C   *             Grid V:    69.0 kV     *            |  |
|  |  Ambient:    28.5 C   *                                                  |  |
|  |                                                                              |  |
|  |  ----------------------------------------------------------------------------  |
|  |                                                                              |  |
|  |  [Mini Trend: Power v]                                                    |  |
|  |  +-----------------------------------------------------------------------+  |
|  |  | 500|                          xxxxx                                   |  |
|  |  | 400|               xxxxxx          xxxx                              |  |
|  |  | 300|      xxxxxx                    xxxx                             |  |
|  |  | 200| xxxx                              xxxx                          |  |
|  |  | 100|xx                                    xxxxxx                     |  |
|  |  |   0|xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx         |  |
|  |  |     06:00   08:00   10:00   12:00   14:00   16:00   18:00   20:00     |  |
|  |  +-----------------------------------------------------------------------+  |
|  |                                                                              |  |
|  +------------------------------------------------------------------------------+  |
|                                                                              |
+------------------------------------------------------------------------------+
```

### 8.4 Timeline Wireframe

```
+------------------------------------------------------------------------------+
|                                                                              |
|  +------------------------------------------------------------------------------+
|  | TIMELINE                                    [< <] [<] [>] [>] [>] [o Reset] |
|  +------------------------------------------------------------------------------+
|  |                                                                              |
|  |  06:00    08:00    10:00    12:00    14:00    16:00    18:00    20:00       |
|  |   |        |        |        |        |        |        |        |           |
|  |   *--------*--------*--------*--------*--------*--------*--------*           |
|  |                                           A                                 |
|  |                                      [CURRENT]                              |
|  |                                                                              |
|  |  +-------------+  +-------------+  +-------------+  +-------------+           |
|  |  | Sun Sunrise|  | Sun Peak   |  | Sun Peak   |  | Sun Sunset |           |
|  |  |   06:23   |  |   12:34   |  |   13:45   |  |   18:47   |           |
|  |  +-------------+  +-------------+  +-------------+  +-------------+           |
|  |                                                                              |
|  |  [+ Add Marker]                                                             |
|  |                                                                              |
|  +------------------------------------------------------------------------------+
|                                                                              |
+------------------------------------------------------------------------------+
```

### 8.5 Event Log Wireframe

```
+------------------------------------------------------------------------------+
|                                                                              |
|  +------------------------------------------------------------------------------+
|  | EVENT LOG                                    [Filter v] [Export] [Clear]    |
|  +------------------------------------------------------------------------------+
|  |                                                                              |
|  |  [All] [Info] [Warnings] [Errors] [Debug]                                    |
|  |                                                                              |
|  +------------------------------------------------------------------------------+
|  |                                                                              |
|  |  TIME        LEVEL   SOURCE        MESSAGE                                   |
|  |  ----------------------------------------------------------------------------  |
|  |                                                                              |
|  |  10:00:00.000  [INFO]   Runtime      Simulation started                     |
|  |  10:00:00.100  [INFO]   Sun          Sunrise detected                       |
|  |  10:00:00.200  [INFO]   INV-001      Inverter starting                     |
|  |  10:00:00.500  [INFO]   INV-001      Grid connection established            |
|  |  10:00:01.000  [INFO]   INV-001      MPPT mode activated                   |
|  |  10:00:02.000  [INFO]   MTR-001      Power flow detected                   |
|  |  10:05:00.000  [INFO]   INV-001      Operating at rated capacity           |
|  |  10:30:00.000  [WARN]   Weather      Cloud passing, irradiance drop         |
|  |  10:30:00.500  [INFO]   INV-001      Power derating due to low irradiance  |
|  |  11:00:00.000  [INFO]   Weather      Cloud cleared                          |
|  |  11:00:00.500  [INFO]   INV-001      Returning to MPPT mode                |
|  |  12:00:00.000  [DEBUG]  INV-001      DC_V=600.2V DC_I=833.5A              |
|  |  12:00:00.100  [DEBUG]  INV-001      AC_V=480.1V AC_I=833.0A              |
|  |  12:00:00.200  [DEBUG]  INV-001      Efficiency=98.01%                     |
|  |                                                                              |
|  +------------------------------------------------------------------------------+
|  |                                                                              |
|  |  DETAILS                                                                   |
|  |  ----------------------------------------------------------------------------  |
|  |                                                                              |
|  |  10:30:00.500  [INFO]   INV-001      Power derating                          |
|  |                                                                              |
|  |  Irradiance dropped from 950 W/m^ to 420 W/m^                               |
|  |  Inverter reducing power output to maintain efficiency                       |
|  |                                                                              |
|  |  Current State:                                                             |
|  |    DC Power:   210 kW (derated from 450 kW)                                 |
|  |    AC Power:   206 kW                                                       |
|  |    Efficiency: 98.0%                                                         |
|  |                                                                              |
|  +------------------------------------------------------------------------------+
|                                                                              |
+------------------------------------------------------------------------------+
```

### 8.6 Signal Trace Wireframe

```
+------------------------------------------------------------------------------+
|                                                                              |
|  +------------------------------------------------------------------------------+
|  | SIGNAL TRACE                                        [Clear] [Export]         |
|  +------------------------------------------------------------------------------+
|  |                                                                              |
|  |  Select source:  [INV-001 v]                           [Trace ->]           |
|  |                                                                              |
|  +------------------------------------------------------------------------------+
|  |                                                                              |
|  |  TRACE PATH                                                                 |
|  |  ----------------------------------------------------------------------------  |
|  |                                                                              |
|  |  +-------------+                                                            |
|  |  | SUN MODEL   |  irradiance = 950 W/m^                                      |
|  |  +------+------+                                                            |
|  |         |                                                                   |
|  |         v                                                                   |
|  |  +-------------+                                                            |
|  |  |  PV ARRAY   |  dc_power = 855 kW (5000 m^ x 0.18 x 950)                |
|  |  |  5,000 m^   |  panel_temp = 45.2 C                                       |
|  |  +------+------+                                                            |
|  |         |                                                                   |
|  |         v                                                                   |
|  |  +-------------+                                                            |
|  |  |  INVERTER   |  dc_power = 855 kW                                           |
|  |  |   INV-001   |  ac_power = 838 kW (855 x 0.98)                             |
|  |  |   98.0%     |  dc_voltage = 600 V                                         |
|  |  +------+------+                                                            |
|  |         |                                                                   |
|  |         v                                                                   |
|  |  +-------------+                                                            |
|  |  | TRANSFORMER |  primary_power = 838 kW                                     |
|  |  |   TX-001    |  secondary_power = 825 kW (838 x 0.985)                     |
|  |  |   98.5%     |  hv_voltage = 69.0 kV                                      |
|  |  +------+------+                                                            |
|  |         |                                                                   |
|  |         v                                                                   |
|  |  +-------------+                                                            |
|  |  |MAIN BREAKER |  state = CLOSED                                             |
|  |  |   CB-001    |  output_power = 825 kW                                      |
|  |  +------+------+                                                            |
|  |         |                                                                   |
|  |         v                                                                   |
|  |  +-------------+                                                            |
|  |  |    GRID     |  Grid receives 825 kW from plant                            |
|  |  |   69 kV     |  Grid frequency: 60.00 Hz                                  |
|  |  +-------------+                                                            |
|  |                                                                              |
|  +------------------------------------------------------------------------------+
|  |                                                                              |
|  |  CALCULATION DETAILS                                                        |
|  |  ----------------------------------------------------------------------------  |
|  |                                                                              |
|  |  DC Power = Irradiance x Area x Efficiency                                   |
|  |           = 950 W/m^ x 5,000 m^ x 0.18                                      |
|  |           = 855,000 W = 855 kW                                               |
|  |                                                                              |
|  |  AC Power = DC Power x Inverter Efficiency                                   |
|  |           = 855 kW x 0.98                                                    |
|  |           = 837.9 kW ~ 838 kW                                                |
|  |                                                                              |
|  |  HV Power = AC Power x Transformer Efficiency                                |
|  |           = 838 kW x 0.985                                                   |
|  |           = 825.4 kW ~ 825 kW                                                |
|  |                                                                              |
|  +------------------------------------------------------------------------------+
|                                                                              |
+------------------------------------------------------------------------------+
```

---

## 9. User Journey

### 9.1 Scenario: First-Time User

```
+------------------------------------------------------------------------------+
|                                                                              |
|  STEP 1: Launch Application                                                  |
|  ----------------------------------------------------------------------------  |
|  User double-clicks Forge icon                                               |
|                                                                              |
|  +----------------------------------------------------------------------+    |
|  |                                                                       |    |
|  |                     FORGE                                            |    |
|  |                  Solar Farm Workbench                                |    |
|  |                                                                       |    |
|  |                 [+ New Solar Farm Project]                            |    |
|  |                 [Open Existing Project]                              |    |
|  |                 [Load Reference World]                               |    |
|  |                                                                       |    |
|  +----------------------------------------------------------------------+    |
|                                                                              |
|  User sees clear options. They might start with Reference World to learn.    |
|                                                                              |
+------------------------------------------------------------------------------+

                                    |

+------------------------------------------------------------------------------+
|                                                                              |
|  STEP 2: Load Reference World                                               |
|  ----------------------------------------------------------------------------  |
|  User clicks "Load Reference World"                                          |
|                                                                              |
|  +----------------------------------------------------------------------+    |
|  |  Reference Worlds                                    [Load] [Cancel]  |    |
|  |  -------------------------------------------------------------------   |    |
|  |                                                                       |    |
|  |  ( ) 50 MW Utility-Scale Solar Farm                                   |    |
|  |      Single-axis trackers, 34.5 kV collection, PCC metering           |    |
|  |                                                                       |    |
|  |  (*) 10 MW Commercial Solar + Storage                                 |    |
|  |      Fixed tilt, 480 V collection, battery storage                    |    |
|  |                                                                       |    |
|  |  ( ) 5 MW Community Solar Garden                                      |    |
|  |      Ground-mount fixed tilt, single inverter                          |    |
|  |                                                                       |    |
|  |  ( ) Residential Rooftop (Educational)                                |    |
|  |      Small-scale single inverter demo                                  |    |
|  |                                                                       |    |
|  +----------------------------------------------------------------------+    |
|                                                                              |
|  User selects "50 MW Utility-Scale" (the reference world)                    |
|                                                                              |
+------------------------------------------------------------------------------+

                                    |

+------------------------------------------------------------------------------+
|                                                                              |
|  STEP 3: Explore Pre-Built Plant                                             |
|  ----------------------------------------------------------------------------  |
|  Editor loads with complete solar farm                                       |
|                                                                              |
|  +---------------------------------------------------------------------+    |
|  |                                                                      |    |
|  |  [Grid]--------[Main Breaker]--------[Transformer]                  |    |
|  |                                              |                      |    |
|  |                                              |                      |    |
|  |                        +---------------------+--------------------+  |    |
|  |                        | Feeder 1            | Feeder 2           |  |    |
|  |                        | [CB]-[INV]-[Array]  | [CB]-[INV]-[Array] |  |    |
|  |                        | [CB]-[INV]-[Array]  | [CB]-[INV]-[Array] |  |    |
|  |                        | [CB]-[INV]-[Array]  | [CB]-[INV]-[Array] |  |    |
|  |                        | [CB]-[INV]-[Array]  | [CB]-[INV]-[Array] |  |    |
|  |                        +---------------------+--------------------+  |    |
|  |                                                                      |    |
|  |  [Weather Station]  [Revenue Meter at PCC]                            |    |
|  |                                                                      |    |
|  +---------------------------------------------------------------------+    |
|                                                                              |
|  User clicks on an inverter to see its properties                            |
|                                                                              |
+------------------------------------------------------------------------------+

                                    |

+------------------------------------------------------------------------------+
|                                                                              |
|  STEP 4: Run Simulation                                                     |
|  ----------------------------------------------------------------------------  |
|  User clicks [>] Run]                                                       |
|                                                                              |
|  +---------------------------------------------------------------------+    |
|  |  Simulation running...                                               |    |
|  |                                                                      |    |
|  |  Time: 12:34:56  |  Speed: 1x  |  [|| Pause] [o Stop]              |    |
|  |                                                                      |    |
|  |  User can see power flowing from arrays through inverters to grid    |    |
|  |                                                                      |    |
|  |  +----------+   +----------+   +----------+   +----------+         |    |
|  |  | Array 1  |--->| Inverter |--->|Transformer|--->| Grid     |         |    |
|  |  |  850 kW  |   |  833 kW  |   |  820 kW  |   |  820 kW  |         |    |
|  |  |    ^     |   |    ^     |   |    ^     |   |    *     |         |    |
|  |  +----------+   +----------+   +----------+   +----------+         |    |
|  |                                                                      |    |
|  |  ^ increasing  v decreasing  * stable                               |    |
|  |                                                                      |    |
|  +---------------------------------------------------------------------+    |
|                                                                              |
|  User pauses to inspect values                                               |
|                                                                              |
+------------------------------------------------------------------------------+

                                    |

+------------------------------------------------------------------------------+
|                                                                              |
|  STEP 5: Inject Fault and Observe Response                                   |
|  ----------------------------------------------------------------------------  |
|  User selects main breaker and opens it                                      |
|                                                                              |
|  +---------------------------------------------------------------------+    |
|  |  Context Menu:                                                       |    |
|  |  ----------------                                                  |    |
|  |  Open Breaker                                                        |    |
|  |  Close Breaker                                                       |    |
|  |  View Settings                                                        |    |
|  |  ----------------                                                  |    |
|  |  Add to Watch                                                        |    |
|  |                                                                      |    |
|  +---------------------------------------------------------------------+    |
|                                                                              |
|  User clicks "Open Breaker"                                                 |
|                                                                              |
|  Simulation shows:                                                          |
|  - Breaker opens                                                            |
|  - Inverters detect islanding                                               |
|  - Inverters shut down for protection                                       |
|  - Event log shows sequence                                                 |
|                                                                              |
|  +---------------------------------------------------------------------+    |
|  |  EVENT LOG:                                                          |    |
|  |  -------------------------------------------------------------------  |    |
|  |  14:30:00.100  [INFO]    CB-001   Breaker OPEN command received     |    |
|  |  14:30:00.200  [INFO]    CB-001   Breaker opened                   |    |
|  |  14:30:00.300  [WARN]    INV-001  Islanding detected               |    |
|  |  14:30:00.400  [WARN]    INV-001  Grid voltage lost                |    |
|  |  14:30:00.500  [INFO]    INV-001  Entering standby mode            |    |
|  |  14:30:00.600  [INFO]    MTR-001  Power flow stopped               |    |
|  |                                                                      |    |
|  +---------------------------------------------------------------------+    |
|                                                                              |
|  User has successfully simulated a fault and observed response               |
|                                                                              |
+------------------------------------------------------------------------------+

                                    |

+------------------------------------------------------------------------------+
|                                                                              |
|  STEP 6: Save and Exit                                                      |
|  ----------------------------------------------------------------------------  |
|  User saves project and closes application                                  |
|                                                                              |
|  Next time: User can open this project and continue work                    |
|                                                                              |
+------------------------------------------------------------------------------+
```

### 9.2 Scenario: Creating a New Project

```
+------------------------------------------------------------------------------+
|                                                                              |
|  STEP 1: Start New Project                                                  |
|  ----------------------------------------------------------------------------  |
|  User clicks "New Solar Farm Project"                                        |
|                                                                              |
|  +----------------------------------------------------------------------+    |
|  |  NEW SOLAR FARM PROJECT                                       [X]    |    |
|  |  ------------------------------------------------------------------    |    |
|  |                                                                       |    |
|  |  Project Name                                                         |    |
|  |  [________________________________]                                  |    |
|  |                                                                       |    |
|  |  Location (for sun position calculation)                             |    |
|  |  Latitude:   [40.0] N                                               |    |
|  |  Longitude:  [-105.0] W                                              |    |
|  |                                                                       |    |
|  |  Plant Capacity                                                      |    |
|  |  [_____50___] MW                                                     |    |
|  |                                                                       |    |
|  |  Grid Connection Voltage                                             |    |
|  |  [34.5 kV v]                                                        |    |
|  |                                                                       |    |
|  |  Project Template                                                    |    |
|  |  ( ) Start from scratch                                              |    |
|  |  (*) Based on reference world                                        |    |
|  |                                                                       |    |
|  |                                                       [Cancel] [Create]|    |
|  |                                                                       |    |
|  +----------------------------------------------------------------------+    |
|                                                                              |
|  User fills in details and clicks "Create"                                  |
|                                                                              |
+------------------------------------------------------------------------------+

                                    |

+------------------------------------------------------------------------------+
|                                                                              |
|  STEP 2: Empty Canvas with Site Configured                                  |
|  ----------------------------------------------------------------------------  |
|  Editor opens with:                                                         |
|  - Site location set                                                        |
|  - Grid model created (69kV)                                               |
|  - Sun model configured (based on location)                                 |
|  - Weather model ready                                                     |
|  - Empty canvas for building                                               |
|                                                                              |
|  +---------------------------------------------------------------------+    |
|  |                                                                      |    |
|  |                              CANVAS                                  |    |
|  |                                                                      |    |
|  |                              (empty)                                 |    |
|  |                                                                      |    |
|  |                              +-------------+                        |    |
|  |                              |   GRID      |                        |    |
|  |                              |   69 kV    |                        |    |
|  |                              +-------------+                        |    |
|  |                                                                      |    |
|  |                                                                      |    |
|  |  HINT: Drag components from the Palette to start building           |    |
|  |                                                                      |    |
|  +---------------------------------------------------------------------+    |
|                                                                              |
|  User is ready to build their plant                                          |
|                                                                              |
+------------------------------------------------------------------------------+
```

---

## 10. Implementation Roadmap

### 10.1 Phase 1: Foundation (Week 1-2)

**Goal:** Enable basic project workflow

| Task | Description | Priority |
|------|-------------|----------|
| T1.1 | Project creation dialog | High |
| T1.2 | Project save/load (JSON format) | High |
| T1.3 | Welcome screen with recent projects | Medium |
| T1.4 | Reference world templates | High |
| T1.5 | Basic canvas rendering | High |

### 10.2 Phase 2: Equipment (Week 3-4)

**Goal:** Enable building electrical topology

| Task | Description | Priority |
|------|-------------|----------|
| T2.1 | Equipment palette (solar-specific) | High |
| T2.2 | Drag-and-drop to canvas | High |
| T2.3 | Equipment configuration | High |
| T2.4 | Basic wiring/connections | High |
| T2.5 | Single-line diagram rendering | Medium |

### 10.3 Phase 3: Simulation (Week 5-6)

**Goal:** Enable running and observing simulations

| Task | Description | Priority |
|------|-------------|----------|
| T3.1 | Simulation run/pause/stop | High |
| T3.2 | Real-time measurements display | High |
| T3.3 | Timeline scrubbing | Medium |
| T3.4 | Speed control | Medium |
| T3.5 | Event log | Medium |

### 10.4 Phase 4: Inspection (Week 7-8)

**Goal:** Enable detailed equipment inspection

| Task | Description | Priority |
|------|-------------|----------|
| T4.1 | Inspector panel with tabs | High |
| T4.2 | Real-time value updates | High |
| T4.3 | Configuration editing | Medium |
| T4.4 | Mini trend charts | Low |
| T4.5 | Export datasheets | Low |

### 10.5 Phase 5: Debugging (Week 9-10)

**Goal:** Enable simulation debugging

| Task | Description | Priority |
|------|-------------|----------|
| T5.1 | Signal trace viewer | Medium |
| T5.2 | Data watch list | Medium |
| T5.3 | Fault injection (breaker) | Medium |
| T5.4 | Simulation breakpoints | Low |
| T5.5 | Memory inspection | Low |

### 10.6 Phase 6: Polish (Week 11-12)

**Goal:** Refine and complete

| Task | Description | Priority |
|------|-------------|----------|
| T6.1 | Keyboard shortcuts | Low |
| T6.2 | Undo/redo | Medium |
| T6.3 | Validation and warnings | Medium |
| T6.4 | Performance optimization | Low |
| T6.5 | Documentation | Low |

---

## 11. Design Principles

### 11.1 Visual Language

| Principle | Application |
|-----------|-------------|
| **Dark First** | Professional engineering tool aesthetic |
| **Semantic Colors** | Green=healthy, Red=fault, Yellow=warning |
| **Grid-Based** | Canvas uses 20px grid for alignment |
| **Engineering Notation** | Use kW, MW, kV, Hz, not watts, volts |
| **No Decoration** | No gratuitous animations or effects |

### 11.2 Interaction Principles

| Principle | Application |
|-----------|-------------|
| **Direct Manipulation** | Drag components, wire terminals |
| **Context Sensitivity** | Right-click shows relevant actions |
| **Keyboard First** | Power users can work without mouse |
| **Progressive Disclosure** | Advanced options in expandable panels |
| **Predictable** | Standard patterns (File>Save, Ctrl+S) |

### 11.3 Information Architecture

| Principle | Application |
|-----------|-------------|
| **Hierarchy** | Explorer shows project structure |
| **Overview + Detail** | Canvas shows plant, Inspector shows component |
| **Real-Time** | Measurements update during simulation |
| **Historical** | Event log shows what happened |
| **Traceable** | Signal trace shows why values are what they are |

---

## 12. Open Questions

The following items require further investigation:

1. **Multi-User** — Is collaboration needed? Version control integration?
2. **Cloud Simulation** — Should simulation run remotely for performance?
3. **Hardware Integration** — Support for physical hardware-in-the-loop?
4. **Report Generation** — Automated commissioning reports?
5. **Standards Compliance** — IEC 61850 configuration support?

These items are marked for future consideration and do not block initial implementation.

---

## 13. Appendix

### A. Glossary of Terms

| Term | Definition |
|------|------------|
| PCC | Point of Common Coupling — where plant connects to grid |
| MPPT | Maximum Power Point Tracking — inverter optimization |
| Array | Collection of solar panels |
| Inverter | Converts DC from arrays to AC for grid |
| Revenue Meter | Measures power flow for billing purposes |
| SCADA | Supervisory Control and Data Acquisition |
| IED | Intelligent Electronic Device |

### B. Reference Standards

- IEEE 1547 — Interconnection and Interoperability of Distributed Energy Resources
- IEC 61850 — Communication Networks for Power Systems
- NEC Article 690 — Solar Photovoltaic Systems

### C. Color Palette

| Usage | Color | Hex |
|-------|-------|-----|
| Background (dark) | Near-black | #1a1a2e |
| Panel background | Dark gray | #252536 |
| Border | Medium gray | #3a3a4c |
| Text (primary) | White | #ffffff |
| Text (secondary) | Light gray | #a0a0b0 |
| Healthy/Online | Green | #4ade80 |
| Warning | Yellow | #fbbf24 |
| Error/Fault | Red | #ef4444 |
| Selection | Blue | #3b82f6 |
| Link/Wire | Cyan | #22d3ee |

---

*Document Version: 1.0*  
*Status: Draft - Pending Review*
