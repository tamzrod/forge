# Infrastructure Model

## Purpose

The Industrial Simulation Runtime simulates complete industrial environments, not just individual devices.

Many simulated entities are **not devices**:

- Electrical Grid
- Sun
- Wind
- Time
- Geography
- Ambient Environment
- Reservoir
- Factory Utilities

These entities are **shared** by many devices. They are not individually addressable. They do not expose industrial protocols. They simply represent the **simulated world** in which devices operate.

---

## Why Infrastructure Exists

The current architecture models devices well:

```
Device → Memory → Behavior → Protocol → MMA2
```

However, real industrial systems depend on shared environmental factors:

- A PV Inverter depends on the **Sun** (irradiance)
- A Revenue Meter depends on the **Grid** (voltage, frequency)
- A Weather Station measures **Ambient Environment** (temperature, wind)

These dependencies are **infrastructure**, not devices.

---

## Three-Layer Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         Simulation Runtime                                 │
│                                                                         │
│  Scheduler │ Simulation Clock │ Plugin Loader │ Device Registry          │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                      Simulation Infrastructure                             │
│                                                                         │
│  Grid │ Sun │ Wind │ Geography │ Ambient Environment │ Time            │
│                                                                         │
│  Shared simulated world                                                 │
│  No protocols, no external clients                                      │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                          Virtual Devices                                  │
│                                                                         │
│  Weather Station │ PV Inverter │ Revenue Meter │ Relay │ PLC            │
│                                                                         │
│  Observe infrastructure, publish to MMA2                                 │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Infrastructure vs Devices

### Devices

| Property | Description |
|----------|-------------|
| **Identity** | Unique device ID |
| **Memory** | Owns private memory image |
| **Protocols** | Exposes Modbus, DNP3, etc. |
| **External** | Connected by Atlas-PPC, SCADA, HMIs |
| **Publishes** | Operational data to MMA2 |

### Infrastructure

| Property | Description |
|----------|-------------|
| **Identity** | No device ID |
| **Memory** | Shared state (not device memory) |
| **Protocols** | None |
| **External** | Never connected by external systems |
| **Publishes** | Never publishes directly |

---

## Infrastructure Examples

### Energy Domain

```
Infrastructure:
├── Grid (voltage, frequency, power)
├── Sun (irradiance, position)
├── Wind (speed, direction)
└── Ambient Temperature

Devices:
├── Weather Station (measures irradiance, temperature)
├── PV Inverter (converts DC power, depends on irradiance)
├── Revenue Meter (measures grid power)
└── Relay (protects grid)
```

### Water Domain

```
Infrastructure:
├── Reservoir (level, flow rate)
├── River (flow, quality)
├── Gravity (head pressure)
└── Ambient Temperature

Devices:
├── Pump (moves water)
├── Valve (controls flow)
├── Flow Meter (measures flow)
└── Tank Sensor (monitors level)
```

### Manufacturing Domain

```
Infrastructure:
├── Factory Power (voltage, frequency)
├── Compressed Air (pressure, flow)
├── Ambient Temperature
└── Network (latency, availability)

Devices:
├── PLC (programmable controller)
├── Robot (automated motion)
├── Conveyor (material transport)
└── Sensor (quality inspection)
```

---

## Infrastructure Characteristics

### 1. Shared State

Infrastructure represents **shared simulation state**:

```go
// Infrastructure is shared, not owned by devices
type Infrastructure struct {
    grid     *GridState
    sun      *SunState
    ambient  *AmbientState
}
```

### 2. Observed by Devices

Devices observe infrastructure:

```go
type PVInverterBehavior struct {
    device       *Device
    infrastructure *Infrastructure
}

func (b *PVInverterBehavior) Tick() {
    // Observe infrastructure
    irradiance := b.infrastructure.Sun().Irradiance()
    temperature := b.infrastructure.Ambient().Temperature()
    
    // Compute output
    power := b.calculatePower(irradiance, temperature)
    
    // Write to device memory
    b.device.Memory().WriteFloat32("output", powerAddr, power)
}
```

### 3. No External Exposure

Infrastructure **never** exposes:
- Industrial protocols
- Modbus registers
- DNP3 points
- REST endpoints

External systems connect to **devices**, not infrastructure.

### 4. Many Observers

Multiple devices can observe the same infrastructure:

```
Sun (irradiance)
    │
    ├── Weather Station (measures)
    ├── PV Inverter (converts)
    └── Simulation Data Logger (records)
```

### 5. Evolves Over Time

Infrastructure changes according to simulation behaviors:

```go
type SunBehavior struct {
    infrastructure *Infrastructure
    timeProvider  TimeProvider
}

func (b *SunBehavior) Tick() {
    // Update sun position based on simulation time
    position := b.timeProvider.SunPosition()
    irradiance := b.calculateIrradiance(position)
    
    b.infrastructure.Sun().SetIrradiance(irradiance)
    b.infrastructure.Sun().SetPosition(position)
}
```

---

## Interaction Model

### Data Flow

```
Simulation Runtime
        │
        ▼
┌───────────────────┐
│  Infrastructure   │
│  (shared state)    │
└───────────────────┘
        │
        ▼
┌───────────────────┐     ┌───────────────────┐
│  Device A         │     │  Device B         │
│  Weather Station   │     │  PV Inverter      │
└───────────────────┘     └───────────────────┘
        │                         │
        ▼                         ▼
┌───────────────────────────────────────────────┐
│                  MMA2                          │
│         (operational memory)                   │
└───────────────────────────────────────────────┘
        │
        ▼
┌───────────────────────────────────────────────┐
│            Atlas-PPC / SCADA                   │
└───────────────────────────────────────────────┘
```

### Real-World Analogy

Infrastructure is like **physics** in a simulation:

- Physics doesn't expose Modbus
- Devices measure physics
- Control systems read device measurements

```
Physics (Infrastructure)
    │
    ├── Devices measure (sensors)
    │
    └── Controllers act (actuators)

Physics never talks to the controller directly.
```

---

## Infrastructure Types

### 1. Physical Environment

Represents the physical world:

| Type | Properties |
|------|------------|
| Sun | irradiance, position, azimuth, elevation |
| Wind | speed, direction, gusts |
| Ambient Temperature | temperature, humidity |
| Geography | latitude, longitude, altitude |

### 2. Utility Networks

Represents industrial utilities:

| Type | Properties |
|------|------------|
| Grid | voltage, frequency, power, islanded |
| Reservoir | level, flow rate, temperature |
| Compressed Air | pressure, flow rate |

### 3. Simulation State

Represents simulation parameters:

| Type | Properties |
|------|------------|
| Time | simulation_time, tick_count |
| Scenario | active_events, conditions |

---

## Infrastructure Behaviors

Infrastructure has **behaviors** that evolve its state:

```go
type InfrastructureBehavior interface {
    Tick()
}

// SunBehavior evolves the sun
type SunBehavior struct {
    infrastructure *Infrastructure
    clock           SimulationClock
}

func (b *SunBehavior) Tick() {
    // Calculate sun position from simulation time
    position := b.calculatePosition(b.clock.Time())
    
    // Calculate irradiance from position
    irradiance := b.calculateIrradiance(position)
    
    // Update infrastructure
    b.infrastructure.Sun().Update(position, irradiance)
}
```

---

## Why Not Model as Devices?

Infrastructure is **not** modeled as devices because:

| Aspect | Devices | Infrastructure |
|--------|---------|----------------|
| Identity | Unique ID | None |
| Memory | Private | Shared |
| Protocols | Yes | No |
| External Clients | Yes | No |
| Observers | None | Many |
| Purpose | Active participants | Environment |

### Key Differences

1. **No Identity**: Infrastructure has no address, no unit ID, no device ID
2. **No Protocols**: Infrastructure never exposes Modbus or DNP3
3. **Shared by Many**: Multiple devices observe the same infrastructure
4. **Environment, Not Participant**: Infrastructure is the stage; devices are the actors

---

## Execution Model

Infrastructure executes **before** devices:

```
┌─────────────────────────────────────────────────────────────────┐
│                          Tick Loop                                │
└─────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│  1. Infrastructure Behaviors                                    │
│     Sun evolves, Grid oscillates, Wind fluctuates               │
└─────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│  2. Device Behaviors                                            │
│     Devices observe infrastructure and update memory             │
└─────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│  3. Publish to MMA2                                            │
│     Devices publish operational data                             │
└─────────────────────────────────────────────────────────────────┘
```

### Determinism

Infrastructure evolution is **deterministic**:

- Same simulation time → same sun position
- Same tick count → same grid state
- No external dependencies

Devices can rely on reproducible infrastructure state.

---

## Summary

| Concept | Purpose |
|---------|---------|
| **Runtime** | Hosts infrastructure and devices |
| **Infrastructure** | Shared simulated world |
| **Devices** | Active participants that observe and publish |

Infrastructure is the **environment** in which devices operate. It is shared, observable, and evolves over time—but it never exposes protocols or connects to external systems.

---

*Last Updated: 2026-07-09*
