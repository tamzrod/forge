# UX Specification V2: Solar Farm Reference World Engineering Workbench

**Document ID:** UX-SPEC-SOLAR-002  
**Version:** 2.0  
**Date:** 2026-07-13  
**Status:** For Implementation  

> **This is an OPTIMIZED specification.**  
> See `OPTIMIZATION_AUDIT.md` for rationale behind all changes.

---

## 1. Mission Statement

### 1.1 Purpose

Forge is a **Solar Farm Reference World Engineering Workbench** — a tool for understanding, learning, and experimenting with utility-scale solar farms through realistic simulation.

### 1.2 Primary Mission

**Forge's first and only deliverable is the Solar Farm Reference World.**

Everything else is deferred to post-MVP.

### 1.3 What This IS

An **Engineering Workbench** where the user:

1. **Observes** a complete, working solar farm
2. **Understands** why the plant behaves as it does
3. **Modifies** equipment and settings
4. **Experiments** with different scenarios
5. **Learns** through cause and effect

### 1.4 What This Is NOT

- NOT a general-purpose simulation platform
- NOT a SCADA system
- NOT a monitoring dashboard
- NOT Node-RED
- NOT a configuration tool for generic simulations

---

## 2. User Profile

### 2.1 Primary User: Electrical Engineer (Learning)

**Background:**
- Understands single-line diagrams conceptually
- Knows what a revenue meter is
- May not know PV inverter operating modes yet
- Is here to **learn** how solar farms work

**Goals:**
- Understand how a solar farm operates
- Learn what affects power output
- Understand protection coordination
- Experiment with different configurations
- Apply learning to real projects

### 2.2 NOT In Scope (Post-MVP)

- Software developers testing integrations
- Engineers designing new plants from scratch
- Operators monitoring real plants
- Multi-user collaboration

---

## 3. Core Principle: Observe First

### 3.1 Learning Path

```
┌─────────────────────────────────────────────────────────────┐
│                                                             │
│   OBSERVE          UNDERSTAND        MODIFY         EXPERIMENT│
│      │                  │               │               │      │
│      ▼                  ▼               ▼               ▼      │
│  ┌────────┐        ┌────────┐      ┌────────┐      ┌────────┐ │
│  │Load    │   →    │Read    │   →  │Change  │   →  │Compare │ │
│  │Reference│        │"Why?"  │      │Settings│      │Results │ │
│  │World   │        │Panel   │      │        │      │        │ │
│  └────────┘        └────────┘      └────────┘      └────────┘ │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 3.2 Traditional vs. Forge Approach

| Traditional | Forge |
|------------|-------|
| Create empty project | Load Reference World |
| Add equipment (trial and error) | Observe working plant |
| Connect equipment (more trial and error) | Understand connections |
| Configure settings (even more trial) | Modify settings |
| Run simulation | Run modified simulation |
| Discover mistakes | Observe improved results |
| Repeat | Experiment further |

---

## 4. First Screen: Welcome

### 4.1 Design

```
+------------------------------------------------------------------------------+
|                                                                              |
|                                                                              |
|                              +------------------+                           |
|                              |                  |                           |
|                              |      FORGE       |                           |
|                              |   Solar Farm     |                           |
|                              |   Reference      |                           |
|                              |   World          |                           |
|                              |                  |                           |
|                              +------------------+                           |
|                                                                              |
|                         ┌─────────────────────────────┐                     |
|                         │                             │                     |
|                         │    LOAD SOLAR FARM          │                     |
|                         │      REFERENCE WORLD        │                     |
|                         │         [START]            │                     |
|                         │                             │                     |
|                         └─────────────────────────────┘                     |
|                                                                              |
|                                                                              |
|                         ┌─────────────────────────────┐                     |
|                         │                             │                     |
|                         │     CREATE NEW PROJECT      │                     |
|                         │      (for experts)          │                     |
|                         │                             │                     |
|                         └─────────────────────────────┘                     |
|                                                                              |
+------------------------------------------------------------------------------+
```

### 4.2 Key Changes from V1

| V1 (Before) | V2 (After) |
|-------------|-------------|
| Three options equal weight | Reference World PRIMARY |
| "New Project" first | "Load Reference World" first |
| "Open Existing" prominent | "Create New" secondary |
| Technical names | Simplified names |

### 4.3 Reference World Selector

After clicking "Load Solar Farm Reference World":

```
+------------------------------------------------------------------------------+
|                                                                              |
|  SELECT REFERENCE WORLD                                               [X]   |
|  ------------------------------------------------------------------------------  |
|                                                                              |
|  ┌────────────────────────────────────────────────────────────────────────┐ |
|  |                                                                         | |
|  |  ┌──────────────────┐                                                  | |
|  |  │                  │  50 MW Utility-Scale Solar Farm              | |
|  |  │   [Preview]      │  ──────────────────────────────────────────   | |
|  |  │                  │                                                  | |
|  |  │                  │  Single-axis tracking arrays, 34.5 kV         | |
|  |  │                  │  collection, revenue metering at PCC.          | |
|  |  │                  │                                                  | |
|  |  │                  │  Components: 100 PV Blocks, 10 Inverters      | |
|  |  │                  │                                                  | |
|  |  └──────────────────┘                                                  | |
|  |                                                                         | |
|  └────────────────────────────────────────────────────────────────────────┘ |
|                                                                              |
|                                                         [Cancel] [Load World]   |
|                                                                              |
+------------------------------------------------------------------------------+
```

---

## 5. Main Editor: Optimized Layout

### 5.1 Primary Layout

```
+------------------------------------------------------------------------------+
|  File  Edit  View  Simulation  Analyze                                    |
+------------------------------------------------------------------------------+
|  [▶ Run]  [⏸ Pause]  [⏹ Stop]  [↺ Reset]                    12:34:56    |
+------------------------------------------------------------------------------+
|        |                                            |                        |
|        |    ┌──────────────────────────────────────┴───────────────────┐   |
|        |    |                                                              |   |
|        |    |                        CANVAS                                 |   |
|        |    |                                                              |   |
|        |    |    ┌─────────┐      ┌─────────┐      ┌─────────┐         |   |
|        |    |    │   PV    │      │   PV    │      │   PV    │         |   |
|        |    |    │  Block  │      │  Block  │      │  Block  │         |   |
|        |    |    │ 1       │      │ 2       │      │ 3       │         |   |
|        |    |    └────┬────┘      └────┬────┘      └────┬────┘         |   |
|        |    |         │                 │                 │               |   |
|        |    |         └────────┬────────┴─────────────────┘               |   |
|        |    |                  │                                        |   |
|        |    |         ┌────────┴────────┐                              |   |
|  PAL   |    |         │    COLLECTOR  │                              |   |
| ETTE   |    |         │      BUS       │                              |   |
|        |    |         └────────┬────────┘                              |   |
|  Plant |    |                  │                                        |   |
|  Env   |    |         ┌───────┴───────┐                              |   |
|        |    |         │                │                              |   |
|        |    |    ┌────┴────┐    ┌─────┴─────┐                       |   |
|        |    |    │Revenue  │    │  Main      │                       |   |
|        |    |    │Meter    │    │  Breaker   │                       |   |
|        |    |    │ (PCC)   │    └─────┬─────┘                       |   |
|        |    |    └─────────┘          │                              |   |
|        |    |                         │                              |   |
|        |    |                   ┌─────┴─────┐                        |   |
|        |    |                   │   Utility  │                        |   |
|        |    |                   │   Grid     │                        |   |
|        |    |                   └────────────┘                        |   |
|        |    |                                                              |   |
|        |    └──────────────────────────────────────────────────────────┘   |
|        |                                                                    |
+--------+--------------------------------------------------------------------+ |
|                                                                             |
|  EXPLORER                                          ANALYSIS                |
|  ─────────                                         ────────                 |
|  ▼ Solar Farm                                     [Timeline][Events][Why?]  |
|    ▼ Plant                                                                  |
|      Grid Interconnection                                                   |
|      Main Breaker                                                          |
|      Revenue Meter                                                         |
|      Collector Bus                                                         |
|      PV Block 1                                                           |
|      PV Block 2                                                           |
|      ...                                                                  |
|    ▼ Environment                                                          |
|      Sun                                                                   |
|      Weather                                                               |
|                                                                             |
+-----------------------------------------------------------------------------+
```

### 5.2 Panel Breakdown

| Panel | Purpose | Keep/Remove |
|-------|---------|--------------|
| Menu Bar | Actions | Keep (minimal) |
| Toolbar | Simulation controls | Keep |
| Palette | Equipment library | Keep (reduced) |
| Canvas | Primary workspace | Keep |
| Inspector | Equipment details | Keep (tabbed) |
| Explorer | Plant hierarchy | Keep (redesigned) |
| Analysis | Timeline, Events, Why? | Keep (merged) |
| Console | Developer logs | **REMOVE** |

---

## 6. Equipment Palette: Solar Farm Only

### 6.1 Optimized Palette

```
+------------------------------------------------------------------------------+
|  PALETTE                                                    [Search]         |
+------------------------------------------------------------------------------+
|                                                                              |
|  ┌────────────────────────────────────────────────────────────────────────┐ |
|  │ PLANT                                                                    │ |
|  │ ------------------------------------------------------------------------│ |
|  │                                                                         │ |
|  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │ |
|  │  │    [GRID]     │  │   [  CB  ]    │  │    [MTR]    │              │ |
|  │  │   34.5 kV     │  │   Breaker     │  │   Revenue    │              │ |
|  │  │   Utility      │  │               │  │    Meter    │              │ |
|  │  └──────────────┘  └──────────────┘  └──────────────┘              │ |
|  │                                                                         │ |
|  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │ |
|  │  │    [BUS]     │  │  [=====]     │  │  [  ⚡  ]    │              │ |
|  │  │   Collector  │  │  PV Block     │  │  Transformer │              │ |
|  │  │     Bus      │  │   5 MW        │  │   Step-Up    │              │ |
|  │  └──────────────┘  └──────────────┘  └──────────────┘              │ |
|  │                                                                         │ |
|  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │ |
|  │  │  [  INV   ]  │  │  [ RELAY ]   │  │  [ 🏭 ]     │              │ |
|  │  │  Inverter    │  │  Protective  │  │   Auxiliary  │              │ |
|  │  │   5 MW       │  │    Relay     │  │    Load      │              │ |
|  │  └──────────────┘  └──────────────┘  └──────────────┘              │ |
|  │                                                                         │ |
|  │  ┌──────────────┐                                                     │ |
|  │  │  [  WX   ]   │                                                     │ |
|  │  │   Weather    │                                                     │ |
|  │  │  Station     │                                                     │ |
|  │  └──────────────┘                                                     │ |
|  │                                                                         │ |
|  └────────────────────────────────────────────────────────────────────────┘ |
|                                                                              |
|  ┌────────────────────────────────────────────────────────────────────────┐ |
|  │ ENVIRONMENT                                                              │ |
|  │ ------------------------------------------------------------------------│ |
|  │                                                                         │ |
|  │  ┌──────────────┐  ┌──────────────┐                                   │ |
|  │  │    (  )      │  │   🌤️        │                                   │ |
|  │  │   (    )     │  │   Weather    │                                   │ |
|  │  │  (  /\  )    │  │   Conditions │                                   │ |
|  │  │    SUN        │  │              │                                   │ |
|  │  └──────────────┘  └──────────────┘                                   │ |
|  │                                                                         │ |
|  └────────────────────────────────────────────────────────────────────────┘ |
|                                                                              |
+------------------------------------------------------------------------------+
```

### 6.2 Equipment List

| Category | Equipment | Purpose |
|----------|-----------|---------|
| **Plant** | Grid Connection | Utility grid interface |
| | Main Breaker | Plant protection |
| | Revenue Meter | PCC measurement |
| | Collector Bus | Combines PV blocks |
| | PV Block | Self-contained generation unit |
| | Inverter | DC to AC conversion |
| | Transformer | Voltage step-up |
| | Protection Relay | Fault detection |
| | Auxiliary Load | Station service |
| | Weather Station | Environmental monitoring |
| **Environment** | Sun | Solar irradiance |
| | Weather | Ambient conditions |

**Total: 12 items (reduced from 23)**

### 6.3 NEW CONCEPT: PV Block

A **PV Block** is a self-contained generation unit that includes:

```
┌─────────────────────────────────────┐
│           PV BLOCK                   │
│                                     │
│   ┌─────────┐   ┌─────────┐      │
│   │ PV Array│   │ PV Array│      │
│   │ (DC)    │   │ (DC)    │      │
│   └────┬────┘   └────┬────┘      │
│        │              │            │
│        └──────┬───────┘            │
│               │ Combiner (internal) │
│               ▼                    │
│        ┌────────────┐              │
│        │  Inverter  │              │
│        │   (AC)     │              │
│        └─────┬──────┘              │
│              │                      │
│              ▼                      │
│        ┌────────────┐              │
│        │Step-Up TX  │              │
│        └────────────┘              │
│                                     │
│  Visible terminals: AC Output       │
│  Internal: Array, Combiner, DC     │
└─────────────────────────────────────┘
```

This reduces palette complexity while maintaining simulation accuracy.

---

## 7. Project Explorer: Plant-Focused

### 7.1 Optimized Explorer

```
+------------------------------------------------------------------------------+
|  EXPLORER                                                                   |
+------------------------------------------------------------------------------+
|                                                                              |
|  ▼ Solar Farm Project                                                       |
|                                                                              |
|    ▼ Plant                                                                  |
|      ▼ Grid Interconnection                                                  │
|        ⚡ Utility Grid (34.5 kV)                                           │
|        ⏺ Main Breaker (1200A)                                              │
|        📊 Revenue Meter (PCC)                                               │
|                                                                              |
|      ▼ Collection System                                                     │
|        ⚡ Collector Bus                                                      │
|        ▼ PV Block 1                                                         │
|          📦 Inverter (5 MW, 98%)                                            │
|          ⚡ Step-Up Transformer (5 MVA)                                     │
|        ▼ PV Block 2                                                         │
|          📦 Inverter (5 MW, 98%)                                            │
|          ⚡ Step-Up Transformer (5 MVA)                                      │
|        ▼ PV Block 3                                                         │
|          ...                                                                │
|                                                                              |
|      ▼ Auxiliary Systems                                                    │
|        🏭 Auxiliary Load (50 kW)                                            │
|                                                                              |
|    ▼ Environment                                                             |
|      ☀️ Sun (Irradiance: 950 W/m²)                                          │
|      🌤️ Weather (Clear, 28°C)                                             │
|                                                                              |
+------------------------------------------------------------------------------+
```

### 7.2 Hierarchy Change

| Before (V1) | After (V2) |
|-------------|------------|
| Project | Project |
| → World | → Plant |
| → Topology | → Grid Interconnection |
| → Entities | → Collection System |
| → Simulation | → Auxiliary Systems |
| → Settings | Environment |

**Key Change:** "World" → "Plant" (engineering language)

---

## 8. Analysis Panel: Merged Functionality

### 8.1 Single Panel with Tabs

```
+------------------------------------------------------------------------------+
|  ANALYSIS                                              [Timeline][Events][Why?] |
+------------------------------------------------------------------------------+
|                                                                              |
|  ┌────────────────────────────────────────────────────────────────────────┐ |
|  │  [Timeline] [Events] [Why?]                                          │ |
|  └────────────────────────────────────────────────────────────────────────┘ |
|                                                                              |
|  ┌────────────────────────────────────────────────────────────────────────┐ |
|  │                                                                         │ |
|  │  (Content changes based on selected tab)                                │ |
|  │                                                                         │ |
|  └────────────────────────────────────────────────────────────────────────┘ |
|                                                                              |
+------------------------------------------------------------------------------+
```

### 8.2 Timeline Tab

```
+------------------------------------------------------------------------------+
|  ANALYSIS                                        [Timeline][Events][Why?]   |
+------------------------------------------------------------------------------+
|                                                                              |
|  ┌────────────────────────────────────────────────────────────────────────┐ |
|  │                                                                         │ |
|  │  06:00    08:00    10:00    12:00    14:00    16:00    18:00    20:00│ │
|  │   |        |        |        |        |        |        |        |   │
|  │   ●────────●────────●────────●────────●────────●────────●────────●────────● │
|  │                                    ▲                                    │   │
|  │                               [CURRENT]                              │   │
|  │                                                                       │   │
|  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌───────────┐ │   │
|  │  │ ☀ Sunrise │  │ ☀️ Peak    │  │ ☀️ Peak    │  │ ☽ Sunset │ │   │
|  │  │   06:23   │  │   12:34   │  │   13:45   │  │   18:47  │ │   │
|  │  └─────────────┘  └─────────────┘  └─────────────┘  └───────────┘ │   │
|  │                                                                       │   │
|  │  [◀◀] [◀] [▶] [▶▶] [↺ Reset]                                         │   │
|  │                                                                       │   │
|  └────────────────────────────────────────────────────────────────────────┘ |
|                                                                              |
+------------------------------------------------------------------------------+
```

### 8.3 Events Tab

```
+------------------------------------------------------------------------------+
|  ANALYSIS                                        [Timeline][Events][Why?]   |
+------------------------------------------------------------------------------+
|                                                                              |
|  ┌────────────────────────────────────────────────────────────────────────┐ |
|  │                                                                         │ |
|  │  ┌──────────────────────────────────────────────────────────────────┐ │ │
|  │  │  TIME        LEVEL   SOURCE        MESSAGE                       │ │ │
|  │  │  ─────────────────────────────────────────────────────────────── │ │ │
|  │  │  10:00:00    INFO    Runtime      Simulation started            │ │ │
|  │  │  10:00:00    INFO    Sun          Sunrise detected              │ │ │
|  │  │  10:00:00    INFO    INV-001      Inverter starting           │ │ │
|  │  │  10:00:05    INFO    INV-001      Grid connection established  │ │ │
|  │  │  10:00:10    INFO    INV-001      MPPT mode activated         │ │ │
|  │  │  10:30:00    WARN    Weather      Cloud cover increasing      │ │ │
|  │  │  10:30:05    INFO    INV-001      Power derating              │ │ │
|  │  │  11:00:00    INFO    Weather      Cloud cleared               │ │ │
|  │  │  11:00:05    INFO    INV-001      Returning to MPPT            │ │ │
|  │  └──────────────────────────────────────────────────────────────────┘ │ │
|  │                                                                         │ │
|  └────────────────────────────────────────────────────────────────────────┘ │
|                                                                              |
+------------------------------------------------------------------------------+
```

### 8.4 Why? Tab (Engineering Explainability)

```
+------------------------------------------------------------------------------+
|  ANALYSIS                                        [Timeline][Events][Why?]   |
+------------------------------------------------------------------------------+
|                                                                              |
|  ┌────────────────────────────────────────────────────────────────────────┐ |
|  │                                                                         │ |
|  │  WHY IS PCC EXPORT = 5.2 MW?                                          │ │
|  │  ─────────────────────────────────────────────────────────────────── │ │
|  │                                                                         │ │
|  │  PCC Export = 5.2 MW                                                  │ │
|  │                                                                         │ │
|  │  Because:                                                              │ │
|  │  ┌─────────────────────────────────────────────────────────────────┐ │ │
|  │  │  PV Block 1 Output:     2.1 MW                    ▲ increasing   │ │ │
|  │  │  PV Block 2 Output:     2.1 MW                    ▲ increasing   │ │ │
|  │  │  PV Block 3 Output:     1.0 MW                    ▲ increasing   │ │ │
|  │  │  ─────────────────────────────────────────────────────────────── │ │ │
|  │  │  Gross Generation:       5.2 MW                                  │ │ │
|  │  │  Auxiliary Load:       -0.0 MW                                  │ │ │
|  │  │  ─────────────────────────────────────────────────────────────── │ │ │
|  │  │  Net Export (PCC):       5.2 MW                                │ │ │
|  │  └─────────────────────────────────────────────────────────────────┘ │ │
|  │                                                                         │ │
|  │  ┌─────────────────────────────────────────────────────────────────┐ │ │
|  │  │  VERIFICATION                                                     │ │ │
|  │  │  Revenue Meter reading: 5.2 MW ✓                                │ │ │
|  │  │  (matches calculated PCC value)                                  │ │ │
|  │  └─────────────────────────────────────────────────────────────────┘ │ │
|  │                                                                         │ │
|  │  [Click on any value above to see calculation details]               │ │
|  │                                                                         │ │
|  └────────────────────────────────────────────────────────────────────────┘ │
|                                                                              |
+------------------------------------------------------------------------------+
```

---

## 9. Inspector: Equipment Details

### 9.1 Single Selection

```
+------------------------------------------------------------------------------+
|  INSPECTOR                                                       [Pin]       |
+------------------------------------------------------------------------------+
|                                                                              |
|  ┌────────────────────────────────────────────────────────────────────────┐ |
|  │  [Overview] [Measurements] [Why?]                                       │ │
|  └────────────────────────────────────────────────────────────────────────┘ |
|                                                                              |
|  ┌────────────────────────────────────────────────────────────────────────┐ |
|  │                                                                         │ |
|  │  OVERVIEW                                                               │ │
|  │  ────────────────────────────────────────────────────────────────────  │ │
|  │                                                                         │ |
|  │  Name:       PV Block 1                                                │ │
|  │  Type:       PV Generation Unit                                         │ │
|  │  Status:     ● Online                                                   │ │
|  │                                                                         │ │
|  │  ────────────────────────────────────────────────────────────────────  │ │
|  │                                                                         │ │
|  │  MEASUREMENTS                                                           │ │
|  │  ────────────────────────────────────────────────────────────────────  │ │
|  │                                                                         │ │
|  │  DC Side                              AC Side                           │ │
|  │  ━━━━━━━━                            ━━━━━━━━                          │ │
|  │  Voltage:    850 V  ▲               Voltage:    480 V    ●           │ │
|  │  Current:    5,882 A ▲               Current:    4,375 A   ●           │ │
|  │  Power:      5.0 MW ▲               Power:      4.9 MW   ▲           │ │
|  │  Irradiance: 950 W/m² ▲             Efficiency: 98.0%    ●           │ │
|  │                                                                         │ │
|  │  Temperature: 45°C                                                        │ │
|  │                                                                         │ │
|  │  ────────────────────────────────────────────────────────────────────  │ │
|  │                                                                         │ │
|  │  WHY? (Explainability)                                                  │ │
|  │  ────────────────────────────────────────────────────────────────────  │ │
|  │                                                                         │ │
|  │  AC Power = 4.9 MW                                                     │ │
|  │                                                                         │ │
|  │  Because:                                                               │ │
|  │  DC Power = 5.0 MW                                                    │ │
|  │  Inverter Efficiency = 98.0%                                            │ │
|  │  ─────────────────────────────────────                                │ │
|  │  AC Power = 5.0 MW × 0.98 = 4.9 MW                                   │ │
|  │                                                                         │ │
|  └────────────────────────────────────────────────────────────────────────┘ │
|                                                                              |
+------------------------------------------------------------------------------+
```

---

## 10. Simplified Menus

### 10.1 File Menu

```
File
├── New Project          Ctrl+N      (for experts)
├── Open Project        Ctrl+O
├── Save Project        Ctrl+S
├── ─────────────
└── Exit               Alt+F4
```

(Removed: Save As, Export, Recent, Properties)

### 10.2 Edit Menu

```
Edit
├── Cut                 Ctrl+X
├── Copy                Ctrl+C
├── Paste               Ctrl+V
└── Delete              Del
```

(Removed: Undo, Redo, Duplicate, Find)

### 10.3 View Menu

```
View
└── (Minimal - rely on default layout)
```

(Removed: All panel toggles, zoom controls, fullscreen)

### 10.4 Simulation Menu

```
Simulation
├── Run                 F5
├── Pause               F6
├── Stop                Shift+F5
└── Reset               Shift+F6
```

(Removed: Step, Speed, Jump to Time)

### 10.5 Analyze Menu

```
Analyze
├── Timeline            Ctrl+1
├── Events              Ctrl+2
└── Why?                Ctrl+3
```

(New menu replacing "Debug")

---

## 11. Keyboard Shortcuts (Simplified)

| Action | Shortcut | Notes |
|--------|----------|-------|
| Run Simulation | F5 | - |
| Pause Simulation | F6 | - |
| Stop Simulation | Shift+F5 | - |
| Reset Simulation | Shift+F6 | - |
| Cut | Ctrl+X | - |
| Copy | Ctrl+C | - |
| Paste | Ctrl+V | - |
| Delete | Del | - |
| New Project | Ctrl+N | Experts only |
| Open Project | Ctrl+O | - |
| Save Project | Ctrl+S | - |
| Timeline Tab | Ctrl+1 | - |
| Events Tab | Ctrl+2 | - |
| Why? Tab | Ctrl+3 | - |

**Total: 14 shortcuts (reduced from 31)**

---

## 12. User Journey: Optimized

### 12.1 First-Time User Journey

```
+------------------------------------------------------------------------------+
|                                                                              |
|  STEP 1: Launch Forge                                                       |
|  ─────────────────────────────────────────────────────────────────────────   |
|                                                                              |
|  User sees Welcome Screen with "Load Solar Farm Reference World" prominent.    |
|                                                                              |
|  +--------------------------------------------------------------------+     |
|  |                                                                     |     |
|  |                     LOAD SOLAR FARM                                  |     |
|  |                       REFERENCE WORLD                                |     |
|  |                          [START]                                     |     |
|  |                                                                     |     |
|  +--------------------------------------------------------------------+     |
|                                                                              |
+------------------------------------------------------------------------------+

                                    ↓

+------------------------------------------------------------------------------+
|                                                                              |
|  STEP 2: Load Reference World                                               |
|  ─────────────────────────────────────────────────────────────────────────   |
|                                                                              |
|  User selects "50 MW Utility-Scale Solar Farm" and clicks Load.               |
|                                                                              |
|  Editor opens with complete solar farm:                                       |
|                                                                              |
|  +---------------------------------------------------------------------+    |
|  |                                                                      |    |
|  |  [Grid]──[Main]──[Transformer]──[Collector Bus]──[PV Blocks]     |    |
|  |                                              │                        |    |
|  |                                        [Revenue Meter]               |    |
|  |                                                                      |    |
|  |  [Weather Station]                                                 |    |
|  |                                                                      |    |
|  +---------------------------------------------------------------------+    |
|                                                                              |
+------------------------------------------------------------------------------+

                                    ↓

+------------------------------------------------------------------------------+
|                                                                              |
|  STEP 3: Run Simulation                                                     |
|  ─────────────────────────────────────────────────────────────────────────   |
|                                                                              |
|  User clicks [▶ Run]. Simulation starts. User observes:                       |
|                                                                              |
|  - Time advancing (sunrise → peak → sunset)                                 |
|  - Power flowing from PV blocks to grid                                      |
|  - Revenue meter counting energy export                                     |
|  - Weather affecting irradiance                                             |
|                                                                              |
|  +---------------------------------------------------------------------+    |
|  |                                                                      |    |
|  |  Simulation running...                            [▶][⏸][⏹]  12:34:56 |    |
|  |                                                                      |    |
|  |  ┌────────┐   ┌────────┐   ┌────────┐   ┌────────┐                |    |
|  |  │ PV 1   │   │ PV 2   │   │ PV 3   │   │ PV 4   │                |    |
|  |  │ 4.9 MW │──▶│ 4.9 MW │──▶│ 4.9 MW │──▶│ 4.9 MW │                |    |
|  |  │   ▲    │   │   ▲    │   │   ▲    │   │   ▲    │                |    |
|  |  └────────┘   └────────┘   └────────┘   └────────┘                |    |
|  |                                                                     |    |
|  +---------------------------------------------------------------------+    |
|                                                                              |
+------------------------------------------------------------------------------+

                                    ↓

+------------------------------------------------------------------------------+
|                                                                              |
|  STEP 4: Inspect Equipment                                                 |
|  ─────────────────────────────────────────────────────────────────────────   |
|                                                                              |
|  User clicks on PV Block 1. Inspector shows:                                 |
|                                                                              |
|  +---------------------------------------------------------------------+    |
|  |  INSPECTOR                                                            |    |
|  |  ─────────────────────────────────────────────────────────────────── |    |
|  |  PV Block 1                                                          |    |
|  |  Status: ● Online                                                     |    |
|  |                                                                       |    |
|  |  DC Power: 5.0 MW  ▲    AC Power: 4.9 MW  ▲                         |    |
|  |  Efficiency: 98.0%      Voltage: 480V                                  |    |
|  +---------------------------------------------------------------------+    |
|                                                                              |
+------------------------------------------------------------------------------+

                                    ↓

+------------------------------------------------------------------------------+
|                                                                              |
|  STEP 5: Understand with "Why?"                                             |
|  ─────────────────────────────────────────────────────────────────────────   |
|                                                                              |
|  User clicks "Why?" tab in inspector:                                        |
|                                                                              |
|  +---------------------------------------------------------------------+    |
|  |  WHY?                                                                  |    |
|  |  ─────────────────────────────────────────────────────────────────── |    |
|  |  AC Power = 4.9 MW                                                   |    |
|  |                                                                       |    |
|  |  Because:                                                             |    |
|  |  DC Power = 5.0 MW                                                  |    |
|  |  × Inverter Efficiency = 98.0%                                        |    |
|  │  = 4.9 MW                                                           |    |
|  +---------------------------------------------------------------------+    |
|                                                                              |
+------------------------------------------------------------------------------+

                                    ↓

+------------------------------------------------------------------------------+
|                                                                              |
|  STEP 6: Modify Equipment                                                   |
|  ─────────────────────────────────────────────────────────────────────────   |
|                                                                              |
|  User pauses simulation, changes:                                           |
|  - PV Block 1 rated power: 5.0 MW → 4.0 MW                                |
|                                                                              |
|  User runs simulation again.                                                 |
|                                                                              |
+------------------------------------------------------------------------------+

                                    ↓

+------------------------------------------------------------------------------+
|                                                                              |
|  STEP 7: Observe Results                                                    |
|  ─────────────────────────────────────────────────────────────────────────   |
|                                                                              |
|  User observes:                                                             |
|  - PCC Export decreased from 5.2 MW to 4.7 MW                              |
|  - "Why?" explains the change                                               |
|                                                                              |
|  User has learned through experimentation!                                   |
|                                                                              |
+------------------------------------------------------------------------------+

                                    ↓

+------------------------------------------------------------------------------+
|                                                                              |
|  STEP 8: Save Project                                                       |
|  ─────────────────────────────────────────────────────────────────────────   |
|                                                                              |
|  User clicks File → Save Project.                                            |
|                                                                              |
|  Project saved. User can return later to continue experimenting.             |
|                                                                              |
+------------------------------------------------------------------------------+
```

---

## 13. MVP Scope Definition

### 13.1 IN SCOPE (Must Have)

| Feature | Priority | Rationale |
|---------|----------|-----------|
| Welcome Screen | High | Entry point |
| Reference World Selection | High | Primary workflow |
| Main Editor Layout | High | Core workspace |
| Equipment Palette (12 items) | High | Build solar farms |
| Canvas with Single-Line Diagram | High | Primary visualization |
| Inspector Panel | High | Equipment details |
| Project Explorer (Plant hierarchy) | High | Navigation |
| Analysis Panel (Timeline/Events/Why?) | High | Understanding |
| Simulation Controls (Run/Pause/Stop/Reset) | High | Core functionality |
| Engineering Explainability | High | Learning focus |
| Project Save/Load | Medium | Persistence |
| Real-time Measurements | Medium | Observation |
| Basic Event Log | Medium | Troubleshooting |

### 13.2 OUT OF SCOPE (Post-MVP)

| Feature | Priority | Rationale |
|---------|----------|-----------|
| Empty Project Creation | Low | Reference World preferred |
| Export (SVG/JSON/PDF) | Low | Not core to learning |
| Undo/Redo | Medium | Complexity vs. benefit |
| Step-by-Step Execution | Low | Not needed for learning |
| Speed Control | Low | Default speed works |
| Zoom Controls | Low | Basic zoom sufficient |
| Full Screen Mode | Low | Nice to have |
| Help System | Low | Post-launch |
| Advanced Protection Schemes | Medium | Post-solar basics |
| Battery Storage | Medium | Post-solar basics |
| IEC 61850 Support | Low | Post-MVP |
| Multi-User Collaboration | Low | Post-MVP |
| Cloud Simulation | Low | Post-MVP |
| Hardware-in-the-Loop | Low | Post-MVP |
| Report Generation | Low | Post-MVP |

---

## 14. Summary: Key Optimizations

| Category | V1 (Before) | V2 (After) | Change |
|----------|-------------|-------------|--------|
| Screens | 8 | 4 | -50% |
| Menu Items | 42 | 16 | -62% |
| Palette Items | 23 | 12 | -48% |
| Explorer Levels | 5 | 3 | -40% |
| Shortcuts | 31 | 14 | -55% |
| Debug Concepts | 5 | 1 (Why?) | -80% |

### Key Principles

1. **Reference World is PRIMARY** — First-time users start with a working plant
2. **Engineering Explainability** — "Why?" replaces debugging
3. **Single Analysis Panel** — Timeline, Events, and Why? merged
4. **Plant-Focused Explorer** — Engineering mental model, not software model
5. **Simplified Menus** — Only essential actions
6. **Observe First** — Users learn by observing, then modifying

---

*Specification Version: 2.0*  
*Status: Ready for Implementation*
