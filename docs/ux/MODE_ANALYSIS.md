# Mode Analysis: Design Mode vs. Operation Mode

**Document ID:** UX-MODE-001  
**Version:** 1.0  
**Date:** 2026-07-14  
**Status:** For Review  

---

## Executive Summary

This document analyzes whether Forge should operate in two distinct user modes: **Design Mode** and **Operation Mode**. The analysis is grounded in the Solar Farm Reference World mission and real power plant engineering workflows.

### Key Finding

**Forge should support two modes, but with a critical distinction:**

1. **Design Mode** (optional for MVP) — For engineers building new configurations
2. **Operation Mode** (primary for MVP) — For all users observing and understanding the Reference World

The primary workflow for the Solar Farm Reference World is **NOT design** but **observation and learning**.

---

## Phase 1: Knowledge Audit

### 1.1 Current Architecture Glossary Assessment

The current glossary contains software architecture terms but lacks engineering lifecycle concepts:

| Missing Concept | Status | Priority |
|----------------|--------|----------|
| Engineering Design | NOT DEFINED | HIGH |
| Commissioning | NOT DEFINED | HIGH |
| Operation | NOT DEFINED | HIGH |
| Maintenance | NOT DEFINED | MEDIUM |
| Digital Twin | NOT DEFINED | HIGH |
| Single Line Diagram | NOT DEFINED | HIGH |
| Control Room | NOT DEFINED | MEDIUM |
| Plant Overview | NOT DEFINED | HIGH |

### 1.2 Required Glossary Additions

The following terms should be added to the Knowledge Base:

```
Engineering Design
────────────────
The process of creating a solar farm configuration. Includes selecting
equipment, defining electrical topology, and configuring settings.

Related: Single Line Diagram, Commissioning, Operation

Commissioning
─────────────
The process of verifying that a solar farm operates correctly. Includes
running simulations, observing behavior, and verifying protection systems.

Related: Engineering Design, Operation, Digital Twin

Operation
─────────
The state when a solar farm is running and producing power. Operators
monitor performance and respond to conditions.

Related: Commissioning, Control Room, Digital Twin

Digital Twin
────────────
A digital representation of a physical system. In Forge, the simulation
world is the digital twin of the solar farm.

Related: Engineering Design, Commissioning, Operation

Single Line Diagram
──────────────────
A simplified graphical representation of an electrical system showing
key components and connections. The primary view in Forge.

Related: Engineering Design, Plant Overview

Control Room
────────────
The facility from which operators monitor and control a plant. In Forge,
this maps to the Operation Mode interface.

Related: Operation, Plant Overview

Plant Overview
─────────────
A high-level view of the entire plant showing key metrics and status.
The primary view in Operation Mode.

Related: Control Room, Single Line Diagram
```

---

## Phase 2: Workflow Discovery

### 2.1 Real Power Plant Lifecycle

A real utility-scale solar farm follows this lifecycle:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                             │
│  ┌─────────────┐     ┌──────────────┐     ┌─────────────┐     ┌───────────┐ │
│  │   DESIGN   │────▶│COMMISSIONING│────▶│  OPERATION │────▶│MAINTENANCE│ │
│  └─────────────┘     └──────────────┘     └─────────────┘     └───────────┘ │
│                                                                             │
│       │                    │                    │                    │        │
│       ▼                    ▼                    ▼                    ▼        │
│  ┌────────────┐      ┌────────────┐      ┌────────────┐      ┌─────────┐  │
│  │  CAD Tools │      │ Test &     │      │  Control   │      │  Plant   │  │
│  │  Electrical│      │ Verify     │      │   Room     │      │   Data   │  │
│  │  Design   │      │ Systems    │      │  Monitor   │      │ Analysis │  │
│  └────────────┘      └────────────┘      └────────────┘      └─────────┘  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 Forge's Role in the Lifecycle

Forge is a **Digital Twin** that serves multiple purposes:

| Lifecycle Phase | Forge's Role | Primary Mode |
|----------------|---------------|--------------|
| Design | Validate new configurations | Design Mode |
| Commissioning | Verify plant behavior | Commissioning Mode |
| Operation | Train operators, test scenarios | Operation Mode |
| Maintenance | Analyze performance data | Analysis Mode |

### 2.3 User Goals Analysis

**For the Solar Farm Reference World:**

| User Goal | Primary Mode | Rationale |
|-----------|--------------|-----------|
| Load a reference world | Operation | Primary workflow |
| Observe how a plant works | Operation | Primary workflow |
| Understand power flow | Operation | Primary workflow |
| Test different scenarios | Operation | Primary workflow |
| Modify equipment settings | Commissioning | Secondary |
| Build a new plant | Design | Post-MVP |

### 2.4 Discovery Conclusion

**The Solar Farm Reference World primarily serves Operation and Commissioning needs.**

- Users load a pre-built plant, they don't design from scratch
- The primary workflow is observation and understanding
- Modifying settings is secondary
- Building new plants is post-MVP

---

## Phase 3: Design Mode

### 3.1 Design Mode Definition

**Design Mode** is the mode for creating and editing solar farm configurations.

### 3.2 Design Mode Responsibilities

| Responsibility | Required | Notes |
|---------------|----------|-------|
| Create Project | YES | Starting point |
| Select Reference World | YES | For experienced users |
| Place Equipment | YES | Drag from palette |
| Connect Equipment | YES | Wire terminals |
| Edit Properties | YES | Configure settings |
| Build Electrical Topology | YES | Single Line Diagram |
| Save Project | YES | Persist configuration |
| Validate Topology | YES | Check electrical rules |
| Simulate | NO | This is Commissioning |

### 3.3 Design Mode Scope (Post-MVP)

Design Mode is **NOT** part of the MVP scope. The Reference World is pre-built.

However, for post-MVP:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          DESIGN MODE                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │                                                                 │   │
│  │                        CANVAS                                      │   │
│  │                    (Single Line Diagram)                           │   │
│  │                                                                 │   │
│  │    ┌─────────┐         ┌─────────┐                               │   │
│  │    │   PV    │────────│   INV   │                               │   │
│  │    │  Array  │         │ 500kW  │                               │   │
│  │    └─────────┘         └────┬────┘                               │   │
│  │                              │                                       │   │
│  │                              ▼                                       │   │
│  │                         ┌─────────┐                                 │   │
│  │                         │   TX    │                                 │   │
│  │                         │  5MVA   │                                 │   │
│  │                         └─────────┘                                 │   │
│  │                                                                 │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  PALETTE                                                           │   │
│  │  ───────                                                           │   │
│  │  PV Block  │  Inverter  │  Transformer  │  Meter  │  Breaker    │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  EQUIPMENT PROPERTIES                                             │   │
│  │  ──────────────────────────                                        │   │
│  │  Name: PV Block 1  │  Rating: 5 MW  │  Efficiency: 18%       │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.4 Design Mode Design Principles

1. **Visual Layout** — Focus on placing equipment visually
2. **Wiring** — Connect terminals by dragging wires
3. **Properties** — Configure each equipment's settings
4. **Validation** — Check electrical rules in real-time
5. **Save** — Persist configuration to file

---

## Phase 4: Operation Mode

### 4.1 Operation Mode Definition

**Operation Mode** is the mode for observing, understanding, and experimenting with a solar farm simulation.

### 4.2 Operation Mode Responsibilities

| Responsibility | Required | Notes |
|---------------|----------|-------|
| Plant Overview | YES | Primary view |
| Single Line Diagram | YES | Live values |
| Run Simulation | YES | Start/pause/stop |
| Clock Control | YES | Time of day |
| Scenario Selection | YES | Weather, time |
| Equipment Inspection | YES | Click to see details |
| Engineering Explainability | YES | WHY? answers |
| Event Log | YES | What happened |
| Timeline | YES | Time scrubbing |

### 4.3 Operation Mode Design (Primary for MVP)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │ PLANT OVERVIEW                                                     │  │
│  ├───────────────────────────────────────────────────────────────────┤  │
│  │                                                                   │  │
│  │  Total Output: 49.0 MW              PCC Export: 48.95 MW        │  │
│  │  ████████████████████████████░░░░ 98% capacity                  │  │
│  │                                                                   │  │
│  │  [▶ Run] [⏸ Pause] [⏹ Stop]        Time: 12:34:56          │  │
│  │                                                                   │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│                                                                             │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │ SINGLE LINE DIAGRAM (Live Values)                                  │  │
│  ├───────────────────────────────────────────────────────────────────┤  │
│  │                                                                   │  │
│  │  ┌────────┐   ┌────────┐   ┌────────┐   ┌────────┐            │  │
│  │  │ PV 1   │──▶│ PV 2   │──▶│   TX    │──▶│  Grid  │            │  │
│  │  │ 4.9MW │   │ 4.9MW │   │ 9.8MW │   │ 34.5kV│            │  │
│  │  │   ▲    │   │   ▲    │   │   ●    │   │   ●    │            │  │
│  │  └────────┘   └────────┘   └────────┘   └────────┘            │  │
│  │                                                                   │  │
│  │  Legend: ▲ increasing  ▼ decreasing  ● stable                    │  │
│  │                                                                   │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│                                                                             │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │ ANALYSIS                                    [Timeline][Events][Why?] │  │
│  ├───────────────────────────────────────────────────────────────────┤  │
│  │                                                                   │  │
│  │  WHY IS PCC EXPORT = 48.95 MW?                                  │  │
│  │  ────────────────────────────────────────────────                │  │
│  │                                                                   │  │
│  │  Because:                                                         │  │
│  │  PV Block 1: 4.90 MW                                           │  │
│  │  PV Block 2: 4.90 MW                                           │  │
│  │  ...                                                             │  │
│  │  Total Gross: 49.00 MW                                          │  │
│  │  - Auxiliary: 0.05 MW                                          │  │
│  │  = Net Export: 48.95 MW                                        │  │
│  │                                                                   │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 4.4 Operation Mode Design Principles

1. **Observation First** — Live values always visible
2. **Explainability** — Every measurement answers "WHY?"
3. **Plant Overview** — High-level metrics always visible
4. **Timeline** — Scrub through time to see history
5. **Events** — What happened and when

---

## Phase 5: Engineering Explainability

### 5.1 Explainability as First-Class Feature

Every measurement in Operation Mode should answer:

**WHY?**

### 5.2 Measurement Explainability Examples

#### Power Flow Explainability

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ WHY IS PCC EXPORT = 48.95 MW?                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  PCC Export = 48.95 MW                                                 │
│                                                                             │
│  Because:                                                                 │
│  ┌─────────────────────────────────────────────────────────────────────┐ │
│  │  Generation                                                         │ │
│  │  ───────────                                                       │ │
│  │  PV Block 1:    4.90 MW    ▲ increasing                         │ │
│  │  PV Block 2:    4.90 MW    ▲ increasing                         │ │
│  │  PV Block 3:    4.90 MW    ▲ increasing                         │ │
│  │  PV Block 4:    4.90 MW    ▲ increasing                         │ │
│  │  ...                                                             │ │
│  │  ─────────────────────────────────────                          │ │
│  │  Total Generation:    49.00 MW                                   │ │
│  └─────────────────────────────────────────────────────────────────────┘ │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐ │
│  │  Station Service                                                  │ │
│  │  ─────────────────                                              │ │
│  │  Auxiliary Load:     -0.05 MW    ● stable                      │ │
│  │  Transformer Loss:   -0.00 MW    ● negligible                   │ │
│  │  ─────────────────────────────────────                          │ │
│  │  Total Service:      -0.05 MW                                   │ │
│  └─────────────────────────────────────────────────────────────────────┘ │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐ │
│  │  Net Export = Generation - Service                               │ │
│  │  Net Export = 49.00 MW - 0.05 MW                               │ │
│  │  Net Export = 48.95 MW                                         │ │
│  └─────────────────────────────────────────────────────────────────────┘ │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### Voltage Explainability

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ WHY IS ARRAY VOLTAGE = 850V?                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Array Voltage = 850V                                                 │
│                                                                             │
│  Because:                                                               │
│  ┌─────────────────────────────────────────────────────────────────────┐ │
│  │  Module Voltage (Vmp)                                              │ │
│  │  ────────────────────────                                          │ │
│  │  Cell Voltage:        0.65V                                       │ │
│  │  Temperature Coeff:  -0.3%/°C                                     │ │
│  │  Cell Temperature:   45°C                                         │ │
│  │  ─────────────────────────────────────                            │ │
│  │  Adjusted Vmp:       0.65 × (1 - 0.003 × (45-25)) = 0.61V     │ │
│  └─────────────────────────────────────────────────────────────────────┘ │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐ │
│  │  Array Configuration                                              │ │
│  │  ─────────────────────                                          │ │
│  │  Modules/String:     20 modules                                   │ │
│  │  Strings/Array:      50 strings                                   │ │
│  │  ─────────────────────────────────────                            │ │
│  │  Array Voltage:      0.61V × 20 × 50 = 610V DC                 │ │
│  │  After MPPT boost:    610V × 1.39 = 850V                        │ │
│  └─────────────────────────────────────────────────────────────────────┘ │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### Fault Explainability

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ WHY DID MAIN BREAKER TRIP?                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Main Breaker opened at 14:30:05                                       │
│                                                                             │
│  Sequence:                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐ │
│  │  14:30:00.000 - Grid fault detected                               │ │
│  │  14:30:00.050 - Grid voltage dropped to 55%                     │ │
│  │  14:30:00.100 - Undervoltage relay issued trip command           │ │
│  │  14:30:00.105 - Main breaker received trip command             │ │
│  │  14:30:00.150 - Main breaker opened                            │ │
│  │  14:30:00.200 - Inverters detected islanding                   │ │
│  │  14:30:00.250 - Inverters entered standby mode                 │ │
│  └─────────────────────────────────────────────────────────────────────┘ │
│                                                                             │
│  Protection Settings:                                                   │
│  ┌─────────────────────────────────────────────────────────────────────┐ │
│  │  Undervoltage Setting:    80% (52V)                             │ │
│  │  Time Delay:              0.1 seconds                             │ │
│  │  ─────────────────────────────────────                            │ │
│  │  Grid Voltage:            55% (35.2V)                           │ │
│  │  Fault Duration:          0.15 seconds                           │ │
│  │  ─────────────────────────────────────                            │ │
│  │  Result: Trip issued ✓ (55% < 80% for 0.15s > 0.1s)         │ │
│  └─────────────────────────────────────────────────────────────────────┘ │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Phase 6: Terminology Audit

### 6.1 Software vs. Engineering Terminology

| Software Term | Engineering Term | Rationale |
|---------------|-------------------|-----------|
| Canvas | Single Line Diagram | Engineers think in SLDs |
| Inspector | Equipment Details | More familiar |
| Project Explorer | Plant Hierarchy | Reflects engineering |
| Debug | Explain | Not debugging, understanding |
| Property Panel | Equipment Properties | More familiar |
| Component | Equipment | Standard engineering term |
| Connection | Wire | Physical reality |
| Entity | Equipment | Standard engineering term |

### 6.2 Updated Terminology Mapping

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ FORGE TERMINOLOGY → ENGINEERING TERMINOLOGY                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Canvas              →  Single Line Diagram (SLD)                       │
│  Inspector           →  Equipment Details                               │
│  Project Explorer    →  Plant Hierarchy                                 │
│  Explorer           →  Plant Navigator                                 │
│  Property Panel     →  Equipment Properties                            │
│  Component          →  Equipment                                      │
│  Connection         →  Wire / Conductor                               │
│  Entity            →  Equipment                                      │
│  Debug Mode        →  Explain Mode                                   │
│  Signal Trace      →  Power Flow Analysis                             │
│  Data Watch        →  Measurement Monitor                             │
│  Memory View       →  Device State                                   │
│  Console           →  (Remove - developer concern)                    │
│  Editor            →  Designer (Design Mode) / Operator (Operation Mode)│
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Phase 7: Mode Separation Analysis

### 7.1 Arguments For Mode Separation

| Argument | Evidence |
|-----------|----------|
| Real engineering tools have distinct modes | AutoCAD has layouts, simulators have edit/run |
| Reduces cognitive load | Users know what they're doing |
| Optimizes UI for each task | Design needs palette, Operation needs values |
| Clear mental model | "Am I building or operating?" |

### 7.2 Arguments Against Mode Separation

| Argument | Evidence |
|-----------|----------|
| Reference World users don't design | Primary workflow is observation |
| Mode switching adds friction | May not need to design at all |
| Single interface can support both | Can show live values while editable |
| Additional complexity | Two modes means more code paths |

### 7.3 Recommended Approach: Seamless with Mode Indicator

**Decision: Support both modes in a single interface with clear mode indication.**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │  [DESIGN]  [OPERATION]           Simulation: Running  12:34:56     │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  In DESIGN mode:                                                        │
│  - Palette visible                                                     │
│  - Equipment draggable                                                 │
│  - Properties editable                                                 │
│  - Live values shown (grayed when simulation not running)              │
│                                                                             │
│  In OPERATION mode:                                                     │
│  - Palette hidden                                                      │
│  - Equipment not draggable                                             │
│  - Live values prominent                                               │
│  - WHY? explainability visible                                         │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 7.4 Mode Transition Rules

| From | To | Trigger | Behavior |
|------|-----|---------|----------|
| None | Operation | Load Reference World | Automatic |
| Operation | Design | Click [Design] | Palette appears |
| Design | Operation | Click [Operation] | Palette hides |
| Any | None | Close Project | Return to Welcome |

---

## Phase 8: Deliverables

### 8.1 UX Audit Summary

| Item | Finding | Decision |
|------|---------|----------|
| Mode Structure | Two modes recommended | Seamless with indicator |
| Primary Mode | Operation | For Reference World |
| Secondary Mode | Design | Post-MVP |
| Explainability | First-class feature | Required |
| Terminology | Software → Engineering | Replace throughout |

### 8.2 Updated UX Philosophy

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ FORGE UX PHILOSOPHY                                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  1. OPERATION MODE IS PRIMARY                                          │
│     For the Reference World, users OBSERVE not DESIGN                   │
│     Every measurement answers "WHY?"                                   │
│                                                                             │
│  2. DESIGN MODE IS OPTIONAL                                           │
│     Post-MVP feature for building new configurations                    │
│     Not required for learning the Reference World                      │
│                                                                             │
│  3. ENGINEERING TERMINOLOGY                                           │
│     Use terms familiar to electrical engineers                         │
│     Single Line Diagram, not Canvas                                    │
│     Equipment, not Components                                          │
│                                                                             │
│  4. SEAMLESS TRANSITION                                              │
│     Single interface with mode indicator                                │
│     Users switch between observing and editing                          │
│                                                                             │
│  5. EXPLAINABILITY FIRST                                             │
│     Every value has context                                            │
│     Show cause and effect                                              │
│     Answer "WHY?" before showing numbers                              │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 8.3 Design Mode Specification

**For post-MVP only.**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ DESIGN MODE (Post-MVP)                                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Purpose: Create new solar farm configurations                         │
│                                                                             │
│  Primary View: Single Line Diagram                                     │
│  Equipment Palette: Always visible                                     │
│  Properties: Always visible                                             │
│  Live Values: Shown but grayed (simulation may not be running)          │
│                                                                             │
│  Actions:                                                              │
│  - Drag equipment from palette to SLD                                 │
│  - Wire terminals by dragging                                          │
│  - Edit equipment properties                                           │
│  - Validate topology (check electrical rules)                         │
│  - Save project                                                       │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 8.4 Operation Mode Specification

**For MVP.**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ OPERATION MODE (MVP)                                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Purpose: Observe, understand, and experiment with solar farms          │
│                                                                             │
│  Primary Views:                                                        │
│  - Plant Overview (always visible)                                     │
│  - Single Line Diagram with live values                                │
│  - Analysis Panel (Timeline, Events, Why?)                             │
│                                                                             │
│  Equipment Palette: Hidden                                             │
│  Properties: Visible on equipment selection                             │
│  Live Values: Prominent and updating                                    │
│                                                                             │
│  Actions:                                                              │
│  - Run/Pause/Stop simulation                                          │
│  - Select equipment to see details                                     │
│  - Click "WHY?" to understand values                                  │
│  - Scrub timeline to replay                                           │
│  - View events to see what happened                                    │
│  - Modify equipment settings (enters Commissioning sub-mode)            │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 8.5 Updated Navigation Model

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ NAVIGATION FLOW                                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────┐                                                      │
│  │  WELCOME   │                                                      │
│  └──────┬──────┘                                                      │
│         │                                                               │
│         ├────────────────────────────────────────────────────┐           │
│         │                                        │                │           │
│         ▼                                        ▼                ▼           │
│  ┌─────────────┐                      ┌─────────────┐      ┌──────────┐ │
│  │   LOAD     │                      │   OPEN    │      │  CREATE  │ │
│  │ REFERENCE  │────────────────────▶│  PROJECT  │      │  NEW     │ │
│  │   WORLD    │                      └───────────┘      │  PROJECT │ │
│  └──────┬──────┘                                          └────┬─────┘ │
│         │                                                    │          │
│         ▼                                                    │          │
│  ┌────────────────────────────────────────────────────────────┐  │          │
│  │                    OPERATION MODE                         │◀─┘          │
│  │  ┌────────────────────────────────────────────────────┐  │             │
│  │  │              PLANT OVERVIEW                        │  │             │
│  │  └────────────────────────────────────────────────────┘  │             │
│  │  ┌────────────────────────────────────────────────────┐  │             │
│  │  │         SINGLE LINE DIAGRAM (Live)                │  │             │
│  │  └────────────────────────────────────────────────────┘  │             │
│  │  ┌────────────────────────────────────────────────────┐  │             │
│  │  │  ANALYSIS  │  Timeline  │  Events  │  Why?  │  │             │
│  │  └────────────────────────────────────────────────────┘  │             │
│  │  ┌────────────────────────────────────────────────────┐  │             │
│  │  │  EQUIPMENT DETAILS (when selected)                 │  │             │
│  │  └────────────────────────────────────────────────────┘  │             │
│  └────────────────────────────────────────────────────────────┘             │
│         │                                                               │
│         │ (optional - for experienced users)                            │
│         ▼                                                               │
│  ┌────────────────────────────────────────────────────────────┐           │
│  │                      DESIGN MODE                           │           │
│  │  ┌──────────┐  ┌────────────────────────────┐            │           │
│  │  │ PALETTE  │  │    SINGLE LINE DIAGRAM    │            │           │
│  │  │          │  │       (Editable)          │            │           │
│  │  │ PV Block │  │                         │            │           │
│  │  │ Inverter │  │    ┌─────┐  ┌─────┐    │            │           │
│  │  │ TX       │  │    │Equipment│──│Equipment│──▶│           │           │
│  │  │ Meter   │  │    └─────┘  └─────┘    │            │           │
│  │  │ Breaker │  │                         │            │           │
│  │  └──────────┘  └────────────────────────────┘            │           │
│  │  ┌────────────────────────────────────────────────────┐  │           │
│  │  │  EQUIPMENT PROPERTIES                            │  │           │
│  │  └────────────────────────────────────────────────────┘  │           │
│  └────────────────────────────────────────────────────────────┘           │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 8.6 Engineering Terminology Guide

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ ENGINEERING TERMINOLOGY GUIDE                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  FROM                    TO                    CONTEXT                  │
│  ──────────────────────  ────────────────────  ─────────────────────────  │
│                                                                             │
│  Canvas                  Single Line Diagram  Primary design view       │
│                          (SLD)                                          │
│                                                                             │
│  Inspector               Equipment Details     Equipment panel           │
│                                                                             │
│  Project Explorer        Plant Navigator      Left navigation panel      │
│                                                                             │
│  Debug Mode              Explain Mode         Understanding values       │
│                                                                             │
│  Signal Trace           Power Flow Analysis  Tracing power paths       │
│                                                                             │
│  Data Watch             Measurement Monitor  Watching live values       │
│                                                                             │
│  Memory View            Device State         Internal equipment state   │
│                                                                             │
│  Component              Equipment            Plant components           │
│                                                                             │
│  Connection             Wire                 Physical conductor         │
│                                                                             │
│  Entity                 Equipment            Same as Component         │
│                                                                             │
│  Properties             Settings            Equipment configuration     │
│                                                                             │
│  Palette                Equipment Library    Drag-drop items           │
│                                                                             │
│  Grid                   Canvas Grid        Drawing grid               │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 8.7 Migration Plan from Current Editor

| Current Element | New Element | Migration Notes |
|-----------------|-------------|-----------------|
| Canvas | Single Line Diagram | Rename + SLD rendering |
| Inspector | Equipment Details | Rename + engineering properties |
| Project Explorer | Plant Navigator | Rename + engineering hierarchy |
| Debug tools | Explain tools | Replace with Why? panel |
| Component palette | Equipment Library | Rename + solar-specific items |
| Entity types | Equipment types | No change, just naming |
| Connection | Wire | Visual update |
| Property panel | Settings | Rename |

### 8.8 Final Recommendation

**Forge should implement a single interface with two modes, but with Operation Mode as the primary and Design Mode as secondary (post-MVP).**

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ FINAL RECOMMENDATION                                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  1. PRIMARY MODE: OPERATION (MVP)                                        │
│     - Plant Overview always visible                                       │
│     - Single Line Diagram with live values                                │
│     - Analysis Panel (Timeline, Events, Why?)                             │
│     - Equipment Details on selection                                      │
│     - Explainability as first-class feature                              │
│                                                                             │
│  2. SECONDARY MODE: DESIGN (Post-MVP)                                    │
│     - Equipment Library visible                                            │
│     - Editable Single Line Diagram                                        │
│     - Properties panel                                                   │
│     - Topology validation                                                │
│                                                                             │
│  3. TERMINOLOGY: ENGINEERING-FIRST                                      │
│     - Replace software terms with engineering terms                       │
│     - Single Line Diagram, not Canvas                                      │
│     - Equipment, not Component                                            │
│     - Explain, not Debug                                                 │
│                                                                             │
│  4. UX PRINCIPLES: OBSERVE FIRST                                        │
│     - Users learn by observing, then modify                              │
│     - Every value answers "WHY?"                                         │
│     - Plant Overview provides context                                     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Acceptance Criteria

The resulting UX makes Forge feel like operating a real utility-scale solar farm.

A new user naturally progresses through:

```
OBSERVE
   ↓
Load Reference World → Run Simulation → Watch Values Change

UNDERSTAND
   ↓
Click Equipment → Read Why? → Trace Cause and Effect

MODIFY
   ↓
Change Settings → Run Simulation → Observe Impact

EXPERIMENT
   ↓
Try Multiple Changes → Compare Results → Learn Trade-offs

DESIGN (post-MVP)
   ↓
Start with Reference → Modify to Requirements → Save as New
```

The interface clearly separates plant construction from plant operation, producing a simpler and more realistic engineering workflow.

---

*Document Version: 1.0*  
*Status: For Review*  
*Mission: Solar Farm Reference World*
