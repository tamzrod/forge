# Industrial Simulation Runtime

## Architectural Philosophy

**The devices are the system. The runtime only hosts them.**

This architecture describes an ecosystem of virtual industrial devices. Every design decision reinforces a single principle:

**Memory is the source of truth.**

### Core Principles

1. **Devices own memory** - Every virtual device owns its memory image
2. **Behaviors modify memory** - Logic reads from and writes to device memory
3. **Protocols expose memory** - External systems read device memory through protocols
4. **Devices never communicate directly** - Devices communicate only by reading and writing memory
5. **Runtime provides infrastructure** - The runtime only hosts, schedules, and coordinates
6. **Plugins provide domain knowledge** - New domains add device types, not runtime changes

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
                            │
        ┌───────────────────┼───────────────────┐
        │                   │                   │
        ▼                   ▼                   ▼
┌───────────────┐   ┌───────────────┐   ┌───────────────┐
│ Revenue Meter │   │Weather Station│   │    Relay     │
│               │   │               │   │               │
│ ┌───────────┐ │   │ ┌───────────┐ │   │ ┌───────────┐│
│ │  Memory   │ │   │ │  Memory   │ │   │ │  Memory   ││
│ │           │ │   │ │           │ │   │ │           ││
│ │  (core)   │ │   │ │  (core)   │ │   │ │  (core)   ││
│ └───────────┘ │   │ └───────────┘ │   │ └───────────┘│
│               │   │               │   │               │
│ ┌───────────┐ │   │ ┌───────────┐ │   │ ┌───────────┐│
│ │ Behaviors │ │   │ │ Behaviors │ │   │ │ Behaviors ││
│ └───────────┘ │   │ └───────────┘ │   │ └───────────┘│
│               │   │               │   │               │
│ ┌───────────┐ │   │ ┌───────────┐ │   │ ┌───────────┐│
│ │ Protocols │ │   │ │ Protocols │ │   │ │ Protocols ││
│ │(external) │ │   │ │(external) │ │   │ │(external) ││
│ └───────────┘ │   │ └───────────┘ │   │ └───────────┘│
└───────────────┘   └───────────────┘   └───────────────┘
```

Memory is the core of each device. Protocols are external views attached to devices.

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

## Device Types (from plugins)

```
Energy Plugin
├── Revenue Meter
├── Weather Station
├── PV Inverter
├── Relay
└── Grid

Water Plugin
├── Pump
├── Valve
├── Tank
└── Flow Meter
```

Adding a new domain requires only adding new device types through plugins. The runtime never changes.

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
| [Runtime](runtime.md) | Infrastructure that hosts devices |
| [Device Model](device-model.md) | Virtual device structure |
| [Memory Model](memory-model.md) | Device memory ownership |
| [Behavior Model](behavior-model.md) | Device-owned logic |
| [Protocol Architecture](protocol-architecture.md) | External memory views |
| [Scheduler](scheduler.md) | Time advancement |
| [Plugin Architecture](plugin-architecture.md) | Device type contribution |
| [Fault Model](fault-model.md) | Memory behavior modification |
| [Scenario Engine](scenario-engine.md) | Event injection |
| [Execution Model](execution-model.md) | End-to-end flow |
