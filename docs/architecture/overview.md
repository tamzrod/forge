# Industrial Simulation Runtime

## Vision

> **A Virtual Industrial Laboratory for industrial software development.**

This project provides a deterministic virtual industrial environment for developing, testing, commissioning, and training industrial software through realistic virtual industrial environments.

**Use cases:** Software development, Controller development, SCADA development, Protocol integration, Factory Acceptance Testing, Commissioning, Training, Education, Demonstrations.

See [Vision](vision.md) for the complete vision statement.

---

## Terminology

**This project has precise, authoritative terminology.**

See the **[Architecture Glossary](GLOSSARY.md)** for definitions of all major concepts.

Key terms used in this document:

| Term | Definition |
|------|------------|
| **Simulation Model** | Represents physical world (Grid, Sun, Weather) |
| **Virtual Firmware** | Samples models, owns Device Memory |
| **Device Memory** | Internal RAM owned by firmware |
| **Communication Interface** | Serializes Device Memory |

---

## Architectural Philosophy

**The devices are the system. Simulation Models represent the physical world. The runtime hosts both.**

This architecture describes an ecosystem of virtual industrial devices within a simulated world. Every design decision reinforces a single principle:

**Memory is the source of truth.**

### Core Principles

1. **Devices own memory** - Every virtual device owns its memory image
2. **Behaviors modify memory** - Logic reads from and writes to device memory
3. **Protocols expose memory** - External systems read device memory through protocols
4. **Devices never communicate directly** - Devices observe models and publish results
5. **Simulation Models represent the physical world** - Grid, Sun, Wind, Weather
6. **Runtime provides hosting** - The runtime hosts models, devices, schedules, and coordinates
7. **Plugins provide domain knowledge** - New domains add model types and device types, not runtime changes

### Why Memory as Foundation

Memory-centric design provides:

- **Deterministic execution** - Same memory state produces same results
- **Simple serialization** - Memory is already structured for storage
- **Easy snapshots** - Freeze memory state at any point
- **Replay capability** - Record memory changes for debugging
- **Protocol independence** - Any protocol can expose memory
- **Direct compatibility** - Native Modbus register and DNP3 point mapping
- **Low coupling** - Behaviors and protocols don't know about each other
- **Cache-friendly** - Sequential memory access patterns
- **Shared source of truth** - Single memory image, multiple protocol views

### One-Minute Summary

```
A virtual device owns deterministic memory.
Behaviors modify that memory.
Protocols expose that memory.
The runtime simply hosts and schedules devices.
Everything else is a plugin.
```

### Fitness for Purpose

Every feature is evaluated against this question:

> *"Does this improve the ability to develop, test, commission, or train industrial software?"*

If the answer is no, it should probably not be part of the Runtime.

### Believe Before Sophisticate

Models should be **credible** before they become **sophisticated**. Simple deterministic models are preferred over highly accurate but complex models unless additional fidelity clearly benefits industrial software development.

---

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    Simulation Runtime                             │
│                     (intentionally small)                       │
│                                                               │
│  Scheduler / Simulation Clock / Plugin Loader                  │
└─────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Simulation Models                           │
│                                                               │
│  Grid │ Sun │ Wind │ Weather │ Reservoir │ Hydraulic            │
│                                                               │
│  Physical world - private RAM - observed by firmware           │
└─────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Virtual Firmware                            │
│                                                               │
│  Weather Station │ PV Inverter │ Revenue Meter │ Relay         │
│                                                               │
│  Samples models, owns device memory, exposes via interfaces    │
└─────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                   Communication Interfaces                        │
│                                                               │
│  Raw Ingest │ Modbus │ DNP3 │ IEC 61850 │ MQTT │ REST         │
│                                                               │
│  Serialize device memory - never access models                  │
└─────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                              MMA2                                 │
└─────────────────────────────────────────────────────────────────┘
```

Four layers: Runtime → Models → Firmware → Interfaces → MMA2

### Data Flow

```
Simulation Models (external physical world)
        ↓
Virtual Firmware samples models, updates Device Memory
        ↓
Communication Interfaces serialize Device Memory
        ↓
MMA2 / SCADA / Historians / Atlas-PPC
```

### Firmware Architecture

Each Virtual Device represents firmware running inside an industrial device:

```
Simulation Models (external world - read-only)
        ↓
Virtual Firmware
├── Identity (ID, Name, Type)
├── Configuration
├── Firmware Logic (samples models, updates memory)
├── Device Memory (owned by firmware)
└── Communication Interfaces (serialize memory)
```

Communication Interfaces:
- Are attached to firmware
- Read Device Memory only
- Never access Simulation Models
- Never perform engineering calculations

---

## Virtual Firmware Anatomy

Every virtual device firmware owns:

| Component | Role |
|-----------|------|
| **Device Memory** | Source of truth. Internal RAM/register space owned by firmware. |
| **Firmware Logic** | Samples models, updates memory. |
| **Communication Interfaces** | Serialize memory for external systems. |
| **Fault Injection Points** | Where faults modify behavior. |

Device Memory stores:
- Temperature, Humidity, Wind Speed, Pressure
- Status, Quality, Timestamp
- Internal state
- Engineering values (already converted by firmware)

---

## Plugin Types (from plugins)

Plugins provide both **Simulation Models** and **Devices**:

```
Energy Plugin

Simulation Models:
├── Grid Model
├── Sun Model
├── Wind Model
└── Weather Model

Devices:
├── Weather Station
├── PV Inverter
├── Revenue Meter
└── Relay

Water Plugin

Simulation Models:
├── Reservoir Model
├── River Model
└── Hydraulic Network Model

Devices:
├── Pump
├── Valve
├── Tank
└── Flow Meter
```

Simulation Models provide the physical world. Devices observe models and publish to MMA2.

---

## Design Philosophy

- **Deterministic** - Same inputs produce same outputs, every time
- **Memory-driven** - Memory is the single source of truth
- **Device-centric** - Devices own their memory and behavior
- **Believable before sophisticated** - Credible models, then enhance
- **Fitness for purpose** - Every feature must improve software development
- **Minimal** - The smallest runtime that achieves the mission

---

## Non-Goals

This project is **not** intended to become:

| Not A | Why Not |
|--------|---------|
| Power system analysis package | Focus is on software behavior, not power flow studies |
| Electromagnetic transient simulator | Not needed for industrial software development |
| Finite element solver | Out of scope for software testing |
| CFD package | Not relevant to industrial protocols |
| Generic physics engine | Domain-specific models are sufficient |
| Digital twin platform | Focus is on software integration, not plant fidelity |

These may inspire future plugins but are not the mission of the Runtime.

---

## Document Map

| Document | Description |
|----------|-------------|
| [Glossary](GLOSSARY.md) | **Authoritative terminology** |
| [Vision](vision.md) | Project purpose, audience, and philosophy |
| [Runtime](runtime.md) | Hosts models and devices |
| [Simulation Models](simulation-models.md) | Physical world representation |
| [Device Model](device-model.md) | Virtual device structure |
| [Memory Model](memory-model.md) | Device memory ownership |
| [Behavior Model](behavior-model.md) | Device-owned logic |
| [Protocol Architecture](protocol-architecture.md) | External memory views |
| [Scheduler](scheduler.md) | Time advancement |
| [Plugin Architecture](plugin-architecture.md) | Domain contribution |
| [Fault Model](fault-model.md) | Memory behavior modification |
| [Scenario Engine](scenario-engine.md) | Event injection |
| [Execution Model](execution-model.md) | End-to-end flow |
