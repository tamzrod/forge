# Architecture Freeze

## Purpose

This document establishes the architectural laws and freezes the architecture from further redesign.

The architecture has reached a stable state. Future work prioritizes implementation over speculation.

## Vision

This project provides a **Virtual Industrial Laboratory** for industrial software development.

The simulator provides a realistic industrial environment for:
- Software development
- Controller development
- SCADA development
- Protocol integration
- Factory Acceptance Testing (FAT)
- Commissioning
- Training
- Education
- Demonstrations

**Every feature is evaluated against this question:**

> *"Does this improve the ability to develop, test, commission, or train industrial software?"*

If the answer is no, it should probably not be part of the Runtime.

See [Vision](vision.md) for the complete vision statement.

## Why the Architecture Was Frozen

The architecture was refined through multiple iterations until it satisfied these criteria:

1. **Simple** - Few concepts, clear responsibilities
2. **Memory-centric** - Memory as single source of truth
3. **Device-owned** - Everything belongs to devices
4. **Runtime-minimal** - Infrastructure only, no business logic
5. **Deterministic** - Reproducible execution
6. **Extensible** - New domains through plugins, no runtime changes
7. **MMA2 Integration** - Clear separation between simulation and plant integration
8. **Simulation Models** - Physical world distinct from devices

The architecture is now considered a contract for implementation.

## Three Memory Domains

There are THREE distinct memory domains. These must never be confused.

### Domain 1: Model State (RAM)

- **Owned by**: Each simulation model
- **Scope**: Private, internal to model
- **Access**: Model methods only
- **Purpose**: Physical state (voltage, irradiance, temperature)
- **Exposure**: Never exposed through protocols

### Domain 2: Device Memory

- **Owned by**: Each virtual device
- **Scope**: Private, internal to device
- **Access**: Device behaviors only
- **Purpose**: Internal device state and simulation logic
- **Exposure**: May be published to MMA2 via Raw Ingest

### Domain 3: Operational Memory (MMA2)

- **Owned by**: MMA2 appliance
- **Scope**: Shared, visible to all
- **Access**: Atlas-PPC, SCADA, HMIs, Historians
- **Purpose**: Plant-wide operational state
- **Exposure**: Exposed via Modbus, DNP3, REST, MQTT

**Model State is NOT Device Memory. Device Memory is NOT MMA2. MMA2 is NOT Model State.**

## Architectural Laws

These principles are no longer under discussion.

### 1. Device Definition

A virtual industrial device is fundamentally:
- Deterministic memory (private, internal)
- Executable behaviors that observe models
- Raw Ingest publisher (to MMA2)

### 2. Device Memory is Source of Truth

Device memory is the single source of truth within the simulation. There is no state outside device memory.

### 3. Simulation Models Represent Physics

Simulation Models represent the physical world (Grid, Sun, Wind, Weather). Models are observed by devices. Models have no protocols and no device identity.

### 4. Devices Observe Models

Devices observe simulation models through their behaviors. Devices never access other devices directly.

### 5. MMA2 is Operational Memory

MMA2 owns the shared operational memory. The simulation runtime publishes to MMA2 via Raw Ingest.

### 6. Behaviors Modify Device Memory

Behaviors read from models, read device memory, and write device memory. Behaviors may publish to MMA2 via Raw Ingest.

### 7. Raw Ingest is the Integration Point

The simulation runtime publishes data to MMA2 via Raw Ingest. The runtime does NOT expose protocols.

### 8. MMA2 Exposes Protocols

MMA2 owns operational memory and exposes protocols (Modbus, DNP3, REST, MQTT).

### 9. Devices Never Communicate Directly

Devices communicate only by observing models and writing memory. There are no direct device references, message buses, callbacks, or service calls between devices.

### 10. Runtime Provides Hosting

The runtime hosts models and devices. It provides scheduling, time advancement, plugin loading, and Raw Ingest connection. It contains no business logic.

### 11. Plugins Provide Domain Knowledge

Plugins provide device types and model types. The runtime remains domain-independent. New domains require no runtime changes.

## Four-Layer Architecture

The architecture consists of four layers:

```
Simulation Runtime
        │
        ▼
Simulation Models
        │
        ▼
Virtual Devices
        │
        ▼
MMA2
```

### Layer 1: Simulation Runtime

The runtime hosts models and devices. It provides scheduling, time advancement, and plugin loading.

### Layer 2: Simulation Models

Models represent the physical world. Examples:
- Grid (electrical)
- Sun, Wind, Weather
- Reservoir, River
- Factory Utilities

Models have no protocols, no external clients, and no device identity. They are observed by devices.

### Layer 3: Virtual Devices

Devices observe models. Devices own memory and publish operational data via Raw Ingest.

### Layer 4: MMA2

MMA2 owns operational memory and exposes protocols to external systems.

## Ownership Rules

```
Simulation Runtime owns:
├── Scheduler
├── Simulation Clock
├── Model Registry
├── Device Registry
├── Plugin Loader
└── Raw Ingest Publisher

Simulation Models own:
├── Grid State
├── Sun State
├── Wind State
├── Weather State
└── Other Physical State

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
- Model state
- Operational memory
- Protocols

### What the Runtime Owns

- Scheduling
- Time advancement
- Model lifecycle
- Device lifecycle
- Plugin loading
- Raw Ingest connection

### What the Runtime Does Not Own

- Device Memory
- Model State
- Behaviors
- Operational memory
- Protocols
- Domain logic

### What Models Own

- Physical state (Grid, Sun, Wind, etc.)
- Model behaviors

### What Models Do Not Own

- Device Memory
- Protocols
- Device identity
- Operational memory

### What MMA2 Owns

- Operational Memory
- Protocols (Modbus, DNP3, REST, MQTT)

### What MMA2 Does Not Own

- Device Memory
- Model State
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

**The Virtual Industrial Laboratory** provides a deterministic, believable industrial environment for developing, testing, commissioning, and training industrial software.

---

**Last Updated:** 2026-07-09
**Status:** Frozen
**Rationale:** Architecture satisfies Virtual Industrial Laboratory vision, simplicity, memory-centricity, device-ownership, and extensibility requirements.
