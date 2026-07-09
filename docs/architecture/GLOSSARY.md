# Architecture Glossary

> **Single source of truth for Forge terminology.**

This glossary defines every major architectural concept used throughout the project. All documentation, code, comments, pull requests, issues, and discussions must use these terms consistently.

---

## Terminology Rules

### Every Concept Has One Preferred Name

| Status | Policy |
|--------|--------|
| **Preferred** | Use this exact term |
| **Avoided** | Do not use these synonyms |

### Recent Renames

| Avoided | Preferred | Reason |
|---------|----------|--------|
| Operational Memory | Device Memory | Owned by firmware, not protocols |
| Publisher | Interface | Communication channel, not publishing |
| Device Behavior | Firmware Logic | Mirrors embedded systems |

---

## Core Architecture Terms

### Simulation Runtime

**Definition:** The execution environment that hosts and coordinates all simulation components.

**Ownership:** Runtime Layer

**Related Terms:**
- Simulation Clock
- Device Registry
- Plugin System

**Common Misunderstandings:**
- The Runtime does NOT implement device logic
- The Runtime does NOT own Simulation Models
- The Runtime is intentionally minimal

**Example:** The Simulation Runtime initializes the clock, loads plugins, and orchestrates the tick cycle.

---

### Simulation Clock

**Definition:** The source of simulation time. Advances deterministically and drives all simulation activity.

**Ownership:** Runtime Layer

**Related Terms:**
- Simulation Runtime
- Tick

**Common Misunderstandings:**
- The Simulation Clock is NOT wall-clock time
- The Simulation Clock does NOT run in real-time by default

**Example:** `simClock.Advance(12 * time.Hour)` moves simulation time forward.

---

### Simulation Model

**Definition:** Represents the physical world being simulated. Contains physics calculations and private state.

**Ownership:** Simulation Layer

**Related Terms:**
- Weather Model
- Grid Model
- Sun Model
- Simulation World

**Common Misunderstandings:**
- A Simulation Model is NOT a Virtual Device
- A Simulation Model does NOT publish to MMA2
- A Simulation Model does NOT have protocols

**Example:** The Weather Model computes atmospheric conditions like temperature and pressure.

---

### Simulation World

**Definition:** The complete collection of all Simulation Models representing the physical environment.

**Ownership:** Simulation Layer

**Related Terms:**
- Simulation Model
- Physical Truth

**Common Misunderstandings:**
- The Simulation World is NOT visible to external systems
- The Simulation World does NOT communicate through protocols

**Example:** The Simulation World contains Grid, Sun, Wind, and Weather models.

---

### Physical Truth

**Definition:** The ground-truth state of the simulated physical system, owned exclusively by Simulation Models.

**Ownership:** Simulation Layer

**Related Terms:**
- Simulation World
- Simulation Model
- Operational Truth

**Common Misunderstandings:**
- Physical Truth is NOT the same as Operational Truth
- Physical Truth is NOT visible to external systems

**Example:** Grid voltage at 240V is Physical Truth; what a Revenue Meter reports is Operational Truth.

---

### Virtual Device

**Definition:** A software representation of an industrial device. Contains Virtual Firmware that samples Simulation Models and owns Device Memory.

**Ownership:** Device Layer

**Related Terms:**
- Virtual Firmware
- Device Memory
- Device Registry
- Device Type

**Common Misunderstandings:**
- A Virtual Device is NOT the same as a Simulation Model
- A Virtual Device does NOT represent physics

**Example:** A Weather Station virtual device observes the Weather Model.

---

### Virtual Firmware

**Definition:** The software running inside a Virtual Device that mirrors embedded industrial firmware. Samples Simulation Models and updates Device Memory.

**Ownership:** Device Layer

**Related Terms:**
- Virtual Device
- Firmware Logic
- Device Memory
- Simulation Context

**Common Misunderstandings:**
- Virtual Firmware is NOT the same as Simulation Model logic
- Virtual Firmware does NOT have direct access to the Simulation World

**Example:** The Weather Station firmware samples the Weather Model and updates its Device Memory.

---

### Firmware Logic

**Definition:** The behavior code within Virtual Firmware that samples Simulation Models and updates Device Memory.

**Ownership:** Device Layer (within Virtual Firmware)

**Related Terms:**
- Virtual Firmware
- Device Memory
- Simulation Context

**Common Misunderstandings:**
- Firmware Logic is NOT responsible for communication protocols
- Firmware Logic is NOT the same as Simulation Model physics

**Example:** On each tick, the firmware logic reads weather values and writes them to Device Memory.

---

### Device Memory

**Definition:** The internal RAM/register space owned by Virtual Firmware. Stores measurements, status, and other device state.

**Ownership:** Device Layer (within Virtual Firmware)

**Related Terms:**
- Virtual Firmware
- Communication Interface

**Avoided Names:**
- Operational Memory
- Firmware Memory
- Internal Memory
- Working Memory

**Common Misunderstandings:**
- Device Memory is NOT owned by protocols
- Device Memory is NOT the same as Simulation Model state
- Device Memory is NOT directly accessible by external systems

**Example:** Device Memory stores temperature=25.5, humidity=60.0, pressure=1013.25.

---

### Communication Interface

**Definition:** A channel attached to Virtual Firmware that serializes Device Memory for transmission to external systems.

**Ownership:** Device Layer (within Virtual Firmware)

**Related Terms:**
- Device Memory
- Protocol
- Raw Ingest

**Avoided Names:**
- Publisher
- Protocol Handler
- Transport

**Common Misunderstandings:**
- A Communication Interface does NOT access Simulation Models
- A Communication Interface does NOT perform engineering calculations
- A Communication Interface does NOT own device state

**Example:** The Raw Ingest Interface serializes Device Memory and sends it to MMA2.

---

### Communication

**Definition:** The act of transmitting data from Device Memory through a Communication Interface to external systems.

**Ownership:** Device Layer

**Related Terms:**
- Communication Interface
- Protocol

**Common Misunderstandings:**
- Communication is NOT initiated by Simulation Models
- Communication does NOT expose Simulation Model state

**Example:** Communication occurs when the Weather Station pushes its memory to MMA2.

---

### Protocol

**Definition:** A defined format and procedure for encoding and transmitting data between systems.

**Ownership:** External

**Related Terms:**
- Communication Interface
- Protocol Adapter

**Common Misunderstandings:**
- A Protocol is NOT the same as a Communication Interface
- A Protocol is NOT owned by the simulation

**Example:** Modbus, DNP3, IEC 61850 are industrial protocols.

---

### Protocol Adapter

**Definition:** A component that translates between a Protocol and a Communication Interface.

**Ownership:** Integration Layer

**Related Terms:**
- Protocol
- Communication Interface

**Common Misunderstandings:**
- A Protocol Adapter is NOT the same as a Communication Interface
- A Protocol Adapter does NOT own device state

**Example:** A Modbus Adapter translates register reads to protocol frames.

---

### Raw Ingest

**Definition:** A specific Communication Interface that serializes Device Memory into a binary format for transmission to MMA2.

**Ownership:** Device Layer (concrete implementation)

**Related Terms:**
- Communication Interface
- MMA2

**Common Misunderstandings:**
- Raw Ingest is NOT a protocol in the traditional sense
- Raw Ingest does NOT access Simulation Models

**Example:** The Weather Station's Raw Ingest Interface sends temperature, humidity, and pressure to MMA2.

---

### Simulation Context

**Definition:** A read-only interface provided to Virtual Firmware for sampling Simulation Models.

**Ownership:** Device Layer

**Related Terms:**
- Virtual Firmware
- Simulation Model

**Common Misunderstandings:**
- The Simulation Context is NOT writeable by firmware
- The Simulation Context does NOT expose all simulation internals

**Example:** `ctx.ReadWeather()` returns a consistent snapshot of weather values.

---

### Device Registry

**Definition:** A container that manages all Virtual Devices in the simulation.

**Ownership:** Device Layer

**Related Terms:**
- Virtual Device
- Simulation Runtime

**Common Misunderstandings:**
- The Device Registry is NOT responsible for device communication
- The Device Registry does NOT own Device Memory

**Example:** `registry.Register(device)` adds a Weather Station to the simulation.

---

### Device Type

**Definition:** A category of Virtual Device with specific behavior and memory layout.

**Ownership:** Device Layer

**Related Terms:**
- Virtual Device
- Device Instance

**Common Misunderstandings:**
- Device Type is NOT the same as Device Instance
- Device Type defines the class, not the individual device

**Example:** `weather_station` is a Device Type; `weather-station-001` is a Device Instance.

---

### Device Instance

**Definition:** A specific, uniquely-identified Virtual Device of a given Device Type.

**Ownership:** Device Layer

**Related Terms:**
- Device Type
- Device Identity

**Common Misunderstandings:**
- Device Instance is NOT the same as Device Type
- Each Device Instance has unique identity

**Example:** `weather-station-001` is a Device Instance of type `weather_station`.

---

### Device Identity

**Definition:** The unique characteristics that identify a Device Instance.

**Ownership:** Device Layer

**Related Terms:**
- Device Instance
- Device Type

**Common Misunderstandings:**
- Device Identity is NOT the same as Device Type
- Device Identity must be unique within the simulation

**Example:** ID=`weather-station-001`, Name=`Weather Station 001`, Type=`weather_station`.

---

### Configuration

**Definition:** The set of parameters that define how a component operates.

**Ownership:** Varies by component

**Related Terms:**
- Device Identity

**Common Misunderstandings:**
- Configuration is NOT the same as state
- Configuration typically does not change during operation

**Example:** Weather Station Configuration includes latitude, sample interval, and publishing settings.

---

### Tick

**Definition:** A single iteration of the simulation loop where all models and devices advance by one time step.

**Ownership:** Runtime Layer

**Related Terms:**
- Simulation Clock
- Simulation Runtime

**Common Misunderstandings:**
- A Tick is NOT the same as wall-clock time
- All components Tick in a defined order

**Example:** On each Tick, models evolve first, then firmware samples models.

---

### Determinism

**Definition:** The property that the same simulation inputs always produce the same outputs, regardless of when or how many times the simulation runs.

**Ownership:** Runtime Layer

**Related Terms:**
- Simulation Clock
- Tick

**Common Misunderstandings:**
- Determinism is NOT the same as reproducibility in real-time
- Determinism requires ordered execution

**Example:** Running the simulation twice with the same seed produces identical results.

---

### MMA2

**Definition:** The external operational telemetry system that receives data from Virtual Devices.

**Ownership:** External System

**Related Terms:**
- Communication Interface
- Raw Ingest
- Operational Truth

**Common Misunderstandings:**
- MMA2 is NOT part of the simulation
- MMA2 does NOT access Simulation Models directly

**Example:** MMA2 receives weather data from the Weather Station's Raw Ingest Interface.

---

### Operational Truth

**Definition:** The state reported by Virtual Devices through Communication Interfaces, representing what external systems observe.

**Ownership:** Device Layer

**Related Terms:**
- Physical Truth
- Device Memory
- MMA2

**Common Misunderstandings:**
- Operational Truth is NOT the same as Physical Truth
- Operational Truth may differ from Physical Truth due to sensor characteristics

**Example:** The Weather Station reports 25.3°C (Operational Truth) while the Weather Model has 25.5°C (Physical Truth).

---

## Simulation Models

### Weather Model

**Definition:** Simulation Model that computes atmospheric conditions including temperature, humidity, pressure, wind, and precipitation.

**Ownership:** Simulation Layer

**Related Terms:**
- Simulation Model
- Sun Model
- Simulation World

**Common Misunderstandings:**
- The Weather Model is NOT a Virtual Device
- The Weather Model does NOT have protocols

**Example:** `weatherModel.Temperature()` returns the current simulated temperature.

---

### Grid Model

**Definition:** Simulation Model that computes electrical grid state including voltage, frequency, and power balance.

**Ownership:** Simulation Layer

**Related Terms:**
- Simulation Model
- Simulation World

**Common Misunderstandings:**
- The Grid Model is NOT a Virtual Device
- The Grid Model does NOT have protocols

**Example:** `gridModel.Voltage()` returns the current simulated grid voltage.

---

### Sun Model

**Definition:** Simulation Model that computes solar position and irradiance based on location and time.

**Ownership:** Simulation Layer

**Related Terms:**
- Simulation Model
- Simulation World

**Common Misunderstandings:**
- The Sun Model is NOT a Virtual Device
- The Sun Model does NOT have protocols

**Example:** `sunModel.Irradiance()` returns the current simulated solar irradiance.

---

## Supporting Systems

### Inspector

**Definition:** A developer tool that provides read-only visibility into the simulation state.

**Ownership:** Development Tools

**Related Terms:**
- World Explorer
- Data Explorer

**Common Misunderstandings:**
- The Inspector is NOT a runtime component
- The Inspector does NOT modify simulation state

**Example:** The web-based Inspector shows Sun, Weather, and Grid model values.

---

### World Explorer

**Definition:** An Inspector component that displays Simulation Models.

**Ownership:** Development Tools

**Related Terms:**
- Inspector
- Data Explorer

**Common Misunderstandings:**
- The World Explorer is NOT the same as Data Explorer
- The World Explorer shows physics, not device state

**Example:** World Explorer shows Weather Model temperature and pressure.

---

### Data Explorer

**Definition:** An Inspector component that displays Device Memory.

**Ownership:** Development Tools

**Related Terms:**
- Inspector
- World Explorer

**Common Misunderstandings:**
- Data Explorer is NOT the same as World Explorer
- Data Explorer shows device measurements

**Example:** Data Explorer shows Weather Station Device Memory values.

---

### Workspace

**Definition:** A development environment configuration that defines which plugins and scenarios are active.

**Ownership:** Development Tools

**Related Terms:**
- Library
- Scenario

**Common Misunderstandings:**
- A Workspace is NOT the same as a Simulation
- A Workspace defines what will be simulated

**Example:** A Workspace loads the Energy plugin with Weather Station and Grid scenarios.

---

### Library

**Definition:** A reusable collection of simulation configurations.

**Ownership:** Development Tools

**Related Terms:**
- Workspace
- Scenario

**Common Misunderstandings:**
- A Library is NOT the same as a Workspace
- A Library provides components for workspaces

**Example:** A Library contains reusable scenarios and device configurations.

---

### Scenario

**Definition:** A predefined sequence of simulation events and conditions.

**Ownership:** Simulation Layer

**Related Terms:**
- Simulation Model
- Workspace

**Common Misunderstandings:**
- A Scenario is NOT the same as a Simulation Model
- A Scenario defines conditions over time

**Example:** A Day/Night Cycle scenario changes Sun Model elevation over 24 hours.

---

## Plugin System

### Plugin

**Definition:** A modular component that provides Simulation Models and Device Types for a specific domain.

**Ownership:** Plugin System

**Related Terms:**
- Plugin Domain
- Simulation Runtime

**Common Misunderstandings:**
- A Plugin is NOT the same as a Runtime
- A Plugin does NOT modify core Runtime behavior

**Example:** The Energy Plugin provides Grid, Sun, Wind, and Weather models plus device implementations.

---

### Plugin Domain

**Definition:** A category of physical or industrial domain served by a Plugin.

**Ownership:** Plugin System

**Related Terms:**
- Plugin

**Common Misunderstandings:**
- Plugin Domain is NOT the same as Plugin
- Plugin Domain defines the scope of a Plugin

**Example:** Energy, Water, and Manufacturing are Plugin Domains.

---

### Simulation HAL

**Definition:** Hardware Abstraction Layer for simulation-specific device interfaces.

**Ownership:** Device Layer

**Related Terms:**
- Virtual Device
- Simulation Model

**Common Misunderstandings:**
- Simulation HAL is NOT the same as hardware drivers
- Simulation HAL provides model access, not hardware access

**Example:** Simulation HAL exposes Simulation Context to Virtual Firmware.

---

## Reserved Terms

These words have specific meanings in Forge and require careful use:

### Explorer

**Reserved because:**
- World Explorer (Simulation Models view)
- Data Explorer (Device Memory view)
- Potential future File Explorer

**Use:** Only for inspection/viewing concepts.

---

### Interface

**Reserved because:**
- Communication Interface (firmware communication channel)

**Use:** Only for describing how firmware communicates externally.

---

### Driver

**Reserved for:**
- Future hardware driver concepts (if needed)

**Avoid:** Using "driver" for software patterns that don't match hardware drivers.

---

### Manager

**Reserved for:**
- Device Registry (manages devices)
- Future resource managers

**Avoid:** Using "manager" for containers or registries.

---

### Controller

**Reserved for:**
- Control logic in Atlas-PPC (future)

**Avoid:** Using for runtime coordination.

---

### Engine

**Reserved for:**
- Scenario Engine (event injection)

**Avoid:** Using for runtime or processing components.

---

### Service

**Reserved for:**
- External service integrations (future)

**Avoid:** Using for internal simulation components.

---

### Publisher

**Avoid:** This term was previously used; use "Communication Interface" instead.

**Reason:** Publisher implies business logic that doesn't belong in firmware.

---

### Adapter

**Reserved for:**
- Protocol Adapters (translation layer)

**Avoid:** Using for generic translation or wrapping patterns.

---

## Naming Guidelines

### Prefer Nouns Over Verbs

| Correct | Avoid |
|---------|-------|
| Device Memory | Memory Management |
| Communication Interface | Communicating |
| Device Registry | Registry Handler |

### Prefer Domain Terminology Over Implementation

| Correct | Avoid |
|---------|-------|
| Firmware Logic | Tick Handler |
| Device Memory | Internal State Map |
| Communication Interface | TCP Writer |

### Avoid Abbreviations

| Correct | Avoid |
|---------|-------|
| Interface | IF |
| Configuration | CFG |
| Temperature | TEMP |

Exception: Industry-standard abbreviations like HVAC, SCADA, IED are acceptable.

### Avoid Overloaded Software Terms

| Correct | Avoid |
|---------|-------|
| Device Memory | Cache |
| Simulation Context | Service |
| Device Registry | Object Pool |

### Match Real Industrial Systems

| Correct | Avoid |
|---------|-------|
| Revenue Meter | Power Meter |
| IED | Intelligent Device |
| RTU | Remote Terminal |

---

## Architecture Review Rule

When adding a new architectural concept:

1. **Check the glossary** for existing terms
2. **Reuse existing terminology** whenever possible
3. **If a new concept is required:**
   - Add it to this glossary
   - Update related documentation
   - Verify no existing term already represents the concept
   - Document avoided names to prevent drift

---

## Glossary Maintenance

This glossary evolves with the architecture.

### Before Each Milestone

Review whether:
- New concepts require glossary entries
- Existing terminology remains accurate
- Duplicate terminology has appeared

### When Adding Terms

Include:
- Preferred name (exact)
- Definition
- Ownership (which layer)
- Related terms
- Avoided names (with reasons)
- Common misunderstandings
- Example usage

### When Renaming Terms

1. Add to "Avoided Names" section
2. Update all documentation references
3. Update code identifiers
4. Document the reason for the rename

---

## Quick Reference

### Layers

| Layer | Components |
|-------|------------|
| Runtime | Simulation Runtime, Simulation Clock, Device Registry |
| Simulation | Simulation Models (Weather, Grid, Sun), Simulation World |
| Device | Virtual Firmware, Device Memory, Communication Interfaces |
| External | MMA2, SCADA, Historians |

### Data Flow

```
Simulation World (Physical Truth)
        ↓
Virtual Firmware samples through Simulation Context
        ↓
Device Memory (firmware-owned)
        ↓
Communication Interfaces (serialize memory)
        ↓
MMA2 (Operational Truth)
```

---

*Last Updated: 2026-07-09*
*Maintained by: Architecture Review*
