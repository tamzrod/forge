# Architecture Freeze

## Purpose

This document establishes the architectural laws and freezes the architecture from further redesign.

The architecture has reached a stable state. Future work prioritizes implementation over speculation.

## Why the Architecture Was Frozen

The architecture was refined through multiple iterations until it satisfied these criteria:

1. **Simple** - Few concepts, clear responsibilities
2. **Memory-centric** - Memory as single source of truth
3. **Device-owned** - Everything belongs to devices
4. **Runtime-minimal** - Infrastructure only, no business logic
5. **Deterministic** - Reproducible execution
6. **Extensible** - New domains through plugins, no runtime changes
7. **MMA2 Integration** - Clear separation between simulation and plant integration
8. **Infrastructure Separation** - Shared simulated world distinct from devices

The architecture is now considered a contract for implementation.

## Two Memory Domains

There are TWO distinct memory domains. These must never be confused.

### Domain 1: Device Memory

- **Owned by**: Each virtual device
- **Scope**: Private, internal
- **Access**: Device behaviors only
- **Purpose**: Internal device state and simulation logic

### Domain 2: Operational Memory (MMA2)

- **Owned by**: MMA2 appliance
- **Scope**: Shared, visible to all
- **Access**: Atlas-PPC, SCADA, HMIs, Historians
- **Purpose**: Plant-wide operational state

**Device Memory is NOT MMA2. MMA2 is NOT Device Memory.**

## Architectural Laws

These principles are no longer under discussion.

### 1. Device Definition

A virtual industrial device is fundamentally:
- Deterministic memory (private, internal)
- Executable behaviors
- Raw Ingest publisher (to MMA2)

### 2. Device Memory is Source of Truth

Device memory is the single source of truth within the simulation. There is no state outside device memory.

### 3. MMA2 is Operational Memory

MMA2 owns the shared operational memory. The simulation runtime publishes to MMA2 via Raw Ingest.

### 4. Behaviors Modify Device Memory

Behaviors read from and write to device memory. Behaviors may publish to MMA2 via Raw Ingest.

### 5. Raw Ingest is the Integration Point

The simulation runtime publishes data to MMA2 via Raw Ingest. The runtime does NOT expose protocols.

### 6. MMA2 Exposes Protocols

MMA2 owns operational memory and exposes protocols (Modbus, DNP3, REST, MQTT).

### 7. Devices Never Communicate Directly

Devices communicate only by reading and writing memory. There are no direct device references, message buses, callbacks, or service calls between devices.

### 8. Runtime Provides Infrastructure

The runtime hosts devices. It provides scheduling, time advancement, plugin loading, and Raw Ingest connection. It contains no business logic.

### 9. Plugins Provide Domain Knowledge

Plugins provide device types. The runtime remains domain-independent. New domains require no runtime changes.

### 10. Infrastructure is the Shared World

Infrastructure represents the simulated world (Grid, Sun, Wind, Ambient Temperature). Infrastructure is observed by devices. Infrastructure has no protocols and no device identity.

## Three-Layer Architecture

The architecture consists of three layers:

```
Simulation Runtime
        │
        ▼
Simulation Infrastructure
        │
        ▼
Virtual Devices
```

### Layer 1: Simulation Runtime

The runtime hosts infrastructure and devices. It provides scheduling, time advancement, and plugin loading.

### Layer 2: Simulation Infrastructure

Infrastructure represents the shared simulated world. Examples:
- Grid (electrical)
- Sun, Wind, Ambient Temperature
- Reservoir, River
- Factory Utilities

Infrastructure has no protocols, no external clients, and no device identity. It is observed by devices.

### Layer 3: Virtual Devices

Devices observe infrastructure. Devices own memory and publish operational data via Raw Ingest.

## Ownership Rules

```
Simulation Runtime owns:
├── Scheduler
├── Simulation Clock
├── Device Registry
├── Plugin Loader
├── Infrastructure Registry
└── Raw Ingest Publisher

Infrastructure owns:
├── Grid State
├── Sun State
├── Wind State
├── Ambient State
└── Other Shared World State

Device owns:
├── Device Memory (private)
├── Behaviors
└── Faults

MMA2 owns:
├── Operational Memory (shared)
└── Protocols (Modbus, DNP3, REST, MQTT)
```

### What Devices Own

- Device Memory (private, internal)
- Behaviors
- Faults

### What Devices Do Not Own

- Scheduling
- Time management
- Plugin loading
- Infrastructure state
- Operational memory
- Protocols

### What the Runtime Owns

- Scheduling
- Time advancement
- Device lifecycle
- Infrastructure lifecycle
- Plugin loading
- Raw Ingest connection

### What the Runtime Does Not Own

- Device Memory
- Infrastructure State
- Behaviors
- Operational memory
- Protocols
- Domain logic

### What Infrastructure Owns

- Shared world state (Grid, Sun, Wind, etc.)
- Infrastructure behaviors

### What Infrastructure Does Not Own

- Device Memory
- Protocols
- Device identity

### What MMA2 Owns

- Operational Memory
- Protocols (Modbus, DNP3, REST, MQTT)

### What MMA2 Does Not Own

- Device Memory
- Infrastructure State
- Behaviors
- Simulation logic

## When Architecture May Be Revisited

Future architectural modifications require evidence, not speculation.

### Valid Reasons to Revisit Architecture

1. **Repeated code duplication** - The same pattern appears in multiple plugins without a clean solution
2. **Performance bottlenecks** - The architecture prevents meeting measurable performance requirements
3. **Ownership confusion** - Unclear where a responsibility belongs despite clear rules
4. **Inability to support a real use case** - A legitimate use case cannot be modeled
5. **Architectural contradiction** - Implementation reveals an inherent conflict in the principles

### Invalid Reasons to Revisit Architecture

- Personal preference
- "This could be cleaner"
- "I would have done it differently"
- Hypothetical future requirements
- Speculation about scale

## Implementation Rule

**Implementation adapts to architecture. Architecture does not adapt to implementation.**

If implementation is difficult, consider:
1. Am I implementing correctly?
2. Is the problem in my code, not the architecture?
3. Can this be solved with a plugin, not a runtime change?

Only genuine architectural limitations justify revisiting the architecture.

## Future Work

Priority order:

1. **Implement Weather Device** - Device with weather behavior, publishes to MMA2 via Raw Ingest
2. **Validate Raw Ingest integration** - Can simulation publish to MMA2?
3. **Validate execution model** - Does the tick loop work as specified?
4. **Implement PV Device** - Device with PV model behavior, publishes to MMA2
5. **Measure performance** - What are actual tick times?

**Important**: Do NOT implement Modbus servers inside the simulation. The simulation publishes to MMA2. MMA2 exposes protocols.

Avoid architecture discussions unless implementation demonstrates a real problem.

## Project Principles

- **Prefer implementation over speculation** - Write code that works
- **Prefer measured evidence over assumptions** - Profile before optimizing
- **Prefer small evolutionary improvements over rewrites** - Incrementally improve working code

## Summary

The architecture is complete. It is a contract for implementation.

The goal is now to prove the architecture through working software rather than continuing to redesign it.

---

**Last Updated:** 2024
**Status:** Frozen
**Rationale:** Architecture satisfies simplicity, memory-centricity, device-ownership, and extensibility requirements.
