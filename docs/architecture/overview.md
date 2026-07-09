# Industrial Simulation Runtime

## Architectural Philosophy

**The devices are the system. Infrastructure is the world they inhabit. The runtime hosts both.**

This architecture describes an ecosystem of virtual industrial devices within a simulated world. Every design decision reinforces a single principle:

**Memory is the source of truth.**

### Core Principles

1. **Devices own memory** - Every virtual device owns its memory image
2. **Behaviors modify memory** - Logic reads from and writes to device memory
3. **Protocols expose memory** - External systems read device memory through protocols
4. **Devices never communicate directly** - Devices observe infrastructure and publish results
5. **Infrastructure represents the shared world** - Grid, Sun, Wind, Environment
6. **Runtime provides infrastructure** - The runtime hosts, schedules, and coordinates
7. **Plugins provide domain knowledge** - New domains add device types and infrastructure, not runtime changes

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

---

## Project Vision

This is a **generic Industrial Simulation Runtime** written in Go. It hosts virtual industrial devices across multiple domains:

- Energy
- Manufacturing
- Water
- Building Automation
- Oil & Gas

The runtime is **not**:
- A power system simulator
- A protocol simulator
- A digital twin platform

Energy is only one possible plugin domain.

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
│                  Simulation Infrastructure                        │
│                                                               │
│  Grid │ Sun │ Wind │ Ambient Temperature │ Geography            │
│                                                               │
│  Shared world - no protocols - observed by devices            │
└─────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                       Virtual Devices                             │
│                                                               │
│  Revenue Meter │ Weather Station │ PV Inverter │ Relay         │
│                                                               │
│  Observe infrastructure, publish to MMA2                       │
└─────────────────────────────────────────────────────────────────┘
```

Three layers: Runtime → Infrastructure → Devices

---

## Device Anatomy

Every virtual device owns:

| Component | Role |
|-----------|------|
| **Memory** | Source of truth. Holds all device state. |
| **Behaviors** | Logic that reads and writes memory. |
| **Protocols** | External interfaces that expose memory. |
| **Faults** | Modifiers that alter memory behavior. |

---

## Plugin Types (from plugins)

Plugins provide both **Infrastructure** and **Devices**:

```
Energy Plugin

Infrastructure:
├── Grid
├── Sun
├── Wind
└── Ambient Temperature

Devices:
├── Weather Station
├── PV Inverter
├── Revenue Meter
└── Relay

Water Plugin

Infrastructure:
├── Reservoir
├── River
└── Ambient Temperature

Devices:
├── Pump
├── Valve
├── Tank
└── Flow Meter
```

Infrastructure provides the shared world. Devices observe infrastructure and publish to MMA2.

---

## Design Philosophy

- **Deterministic** - Same inputs produce same outputs
- **Memory-driven** - Memory is the single source of truth
- **Device-centric** - Devices own their memory and behavior
- **Simple** - No unnecessary abstractions
- **Opinionated** - Clear architectural decisions
- **Minimal** - The smallest runtime possible

---

## Document Map

| Document | Description |
|----------|-------------|
| [Runtime](runtime.md) | Hosts infrastructure and devices |
| [Infrastructure Model](infrastructure-model.md) | Shared simulated world |
| [Device Model](device-model.md) | Virtual device structure |
| [Memory Model](memory-model.md) | Device memory ownership |
| [Behavior Model](behavior-model.md) | Device-owned logic |
| [Protocol Architecture](protocol-architecture.md) | External memory views |
| [Scheduler](scheduler.md) | Time advancement |
| [Plugin Architecture](plugin-architecture.md) | Domain contribution |
| [Fault Model](fault-model.md) | Memory behavior modification |
| [Scenario Engine](scenario-engine.md) | Event injection |
| [Execution Model](execution-model.md) | End-to-end flow |
