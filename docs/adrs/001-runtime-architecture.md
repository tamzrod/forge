# ADR-001: Runtime Architecture Decision

**ADR ID:** ADR-001  
**Title:** Runtime Architecture for Industrial Simulation  
**Date:** 2026-07-11  
**Status:** Accepted  
**Deciders:** Engineering Team  
**Repository:** https://github.com/tamzrod/forge

---

## Context

The Industrial Simulation Runtime requires a deterministic virtual industrial environment for developing, testing, commissioning, and training industrial software. The runtime must:

1. Initialize and manage multiple virtual devices
2. Execute behaviors on a fixed tick interval
3. Sample simulation models (Grid, Sun, Weather, Wind, Reservoir)
4. Expose device memory through communication interfaces
5. Support deterministic execution for reproducible testing

---

## Decision

We adopt a **Scheduler-Centric Architecture** with the following components:

```
┌─────────────────────────────────────────────────────────────┐
│                     Runtime                                  │
│  ┌───────────────┐    ┌───────────────┐    ┌─────────────┐ │
│  │   Scheduler   │───▶│    Device     │───▶│  Behavior   │ │
│  │  (tick loop)  │    │   (memory)    │    │ (execution) │ │
│  └───────────────┘    └───────────────┘    └─────────────┘ │
│         │                    │                    │         │
│         ▼                    ▼                    ▼         │
│  ┌───────────────┐    ┌───────────────┐    ┌─────────────┐ │
│  │    Clock      │    │    Models     │    │  Interface   │ │
│  │ (simulation)  │    │  (physical)   │    │  (protocol)  │ │
│  └───────────────┘    └───────────────┘    └─────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

---

## Architecture Components

### 1. Scheduler

The scheduler is the central coordinator that:
- Manages the simulation tick loop
- Maintains a list of registered devices
- Advances the simulation clock on each tick
- Triggers device behavior execution

**Implementation:**
- `scheduler/scheduler.go`: Scheduler struct with fixed tick interval
- `scheduler.New()`: Creates scheduler with configurable tick duration
- `scheduler.Run(ctx)`: Starts the tick loop
- `scheduler.AddDevice()`: Registers devices

### 2. Device

Devices represent virtual firmware:
- Own device memory (internal RAM)
- Register behaviors for execution
- Sample simulation models through context
- Expose memory through interfaces

**Implementation:**
- `device/device.go`: Device interface and implementation
- `device.New()`: Creates device with memory regions
- `device.AddBehavior()`: Registers behavior
- `device.Tick()`: Executes all behaviors

### 3. Behavior

Behaviors encapsulate device logic:
- Simple function-based interface
- Access to device memory
- Executed on each tick

**Implementation:**
- `device/device.go`: Behavior interface
- `Behavior` interface: `ID()`, `Attach()`, `Detach()`, `Tick()`

### 4. Memory

Device memory provides:
- Named memory regions (e.g., "holding_registers", "input_registers")
- Typed read/write operations
- Quality flag support

**Implementation:**
- `memory/memory.go`: Memory struct with regions
- `Memory.Read()`, `Memory.Write()`
- Quality: Good, Uncertain, Bad, Offline

### 5. Models

Simulation models represent the physical world:
- Grid: Electrical grid (voltage, frequency)
- Sun: Solar position and irradiance
- Weather: Temperature, humidity, pressure, wind
- Wind: Wind speed and direction
- Reservoir: Water storage levels

**Implementation:**
- `models/models.go`: Grid, Sun, Weather, Wind, Reservoir models
- `internal/models/sun/`, `internal/models/weather/`: Detailed implementations

### 6. Interface

Interfaces serialize device memory for protocols:
- Raw Ingest protocol
- Modbus TCP (future)

**Implementation:**
- `internal/publishers/rawingest/`: Raw Ingest protocol
- `internal/devices/weatherstation/`: Weather station device

---

## Consequences

### Positive
- Deterministic execution for reproducible testing
- Clear separation of concerns
- Extensible behavior system
- Model-agnostic device design

### Negative
- Fixed tick interval may not suit all use cases
- Single-threaded execution by default
- Memory model is simplified

### Risks
- Performance at scale (many devices) not yet tested
- Concurrency model needs verification

---

## References

- [Runtime Architecture](docs/architecture/runtime.md)
- [Scheduler Model](docs/architecture/scheduler.md)
- [Device Model](docs/architecture/device-model.md)
- [Memory Model](docs/architecture/memory-model.md)
- [Simulation Models](docs/architecture/simulation-models.md)

---

## Related ADRs

- ADR-002: Behavior Model Design
- ADR-003: Memory Model Design

---

## Milestone Traceability

| Milestone | Status |
|-----------|--------|
| 1.1 Runtime Core | ✅ Complete |
| 1.2 Behaviors | ✅ Complete |
| 1.3 Protocols | ⏳ Pending |
| 1.4 Plugin System | ⏳ Pending |
| 1.5 Fault System | ⏳ Pending |
