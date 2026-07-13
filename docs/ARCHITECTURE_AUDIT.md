# Architecture Audit: Electrical Concepts

## Purpose

This audit compares the Forge implementation against the existing reference architecture to identify gaps in electrical engineering concept definitions.

## Reference Architecture Definitions

### From GLOSSARY.md

| Concept | Definition Source |
|---------|-------------------|
| Simulation Model | Physics (Grid, Sun, Wind) - represents physical world |
| Virtual Device | Equipment (Meter, Inverter, Relay) - mirrors industrial devices |
| Physical Truth | Ground-truth state owned by Simulation Models |
| Revenue Meter | **Explicitly defined** - measures voltage from Grid |

### From simulation-models.md

**Grid Model Properties:**
- voltage (V)
- frequency (Hz)
- Thevenin impedance (Ω)
- short circuit level (MVA)
- reactive sensitivity (PU MVAr per PU voltage)

**Energy Domain Models:**
- Grid - voltage, frequency, Thevenin impedance, short circuit level
- Sun - irradiance, position, azimuth, elevation
- Wind - speed, direction, gusts, turbulence
- Weather - temperature, humidity, pressure, cloud cover

**Energy Domain Devices:**
- Weather Station - observes Weather, Sun
- PV Inverter - observes Sun, injects power to Grid
- Revenue Meter - observes Grid
- Relay - protection device

## Forge Implementation Analysis

### world/electrical/entities.go

| Entity | Status | Reference Match |
|--------|--------|-----------------|
| GridEntity | ✅ Defined | Matches GridModel in reference |
| BusEntity | ⚠️ Gap | **Missing from reference** - electrical topology concept |
| BreakerEntity | ⚠️ Gap | **Missing from reference** - switchgear concept |
| LoadEntity | ⚠️ Gap | **Missing from reference** - power sink concept |
| GeneratorEntity | ⚠️ Gap | **Missing from reference** - power source concept |
| TransformerEntity | ⚠️ Gap | **Missing from reference** - voltage transformation concept |
| MeterEntity | ⚠️ Gap | Reference uses "Revenue Meter" |

### world/environmental/entities.go

| Entity | Status | Reference Match |
|--------|--------|-----------------|
| SunEntity | ✅ Defined | Matches Sun Model |
| WeatherEntity | ✅ Defined | Matches Weather Model |
| PVArrayEntity | ⚠️ Gap | Reference uses "PV Inverter" |

## Missing Electrical Concepts

### 1. Load (Power Sink)

**Definition Needed:**
```
A Load is a power sink that consumes electrical energy.
It withdraws active and/or reactive power from the network.

Properties:
- Rated power (kW)
- Power factor
- Load type (resistive, inductive, capacitive)

Relationship:
- Load → Bus (consumes power)
- Load is NOT a simulation model
- Load IS a virtual device type
```

**Current Implementation:**
```go
type LoadEntity struct {
    power       float32
    basePower   float32
}
```
⚠️ Deviation: Uses internal oscillation instead of physics-based modeling

### 2. Generator (Power Source)

**Definition Needed:**
```
A Generator is a power source that produces electrical energy.
It injects active and/or reactive power into the network.

Properties:
- Rated power (kW)
- Voltage setpoint
- Frequency response

Relationship:
- Generator → Bus (produces power)
- Generator IS a virtual device type
```

**Current Implementation:**
```go
type GeneratorEntity struct {
    powerOutput float32
    ratedPower  float32
    isRunning   bool
}
```
⚠️ Deviation: No distinction between source/sink in architecture

### 3. Breaker (Switchgear)

**Definition Needed:**
```
A Breaker is a switching device that interrupts current flow.
It connects or disconnects sections of the electrical network.

Properties:
- State (open/closed)
- Operating time
- Trip count

Relationship:
- Breaker → Bus (topology control)
- Breaker responds to trip/close commands
- Breaker IS a virtual device type
```

**Current Implementation:**
```go
type BreakerEntity struct {
    isOpen     bool
    tripCount  int
    closeCount int
}
```
⚠️ Deviation: No command interface for trip/close

### 4. Transformer

**Definition Needed:**
```
A Transformer transfers energy between voltage levels.
It changes voltage magnitude while preserving power.

Properties:
- Primary voltage (V)
- Secondary voltage (V)
- Rated power (kVA)
- Efficiency

Relationship:
- Transformer → Bus (connects voltage levels)
- Transformer IS a simulation model (physics)
```

**Current Implementation:**
```go
type TransformerEntity struct {
    primaryVoltage   float32
    secondaryVoltage float32
    ratio            float32
    loading          float32
}
```
⚠️ Deviation: Missing efficiency, impedance modeling

### 5. Meter (Measurement)

**Definition Needed:**
```
A Meter measures electrical quantities.
Reference architecture uses "Revenue Meter" specifically.

Properties:
- Voltage measurement
- Current measurement
- Power measurement (active, reactive)
- Energy accumulation

Relationship:
- Meter → Bus (observes)
- Meter IS a virtual device type
```

**Current Implementation:**
```go
type MeterEntity struct {
    voltage        float32
    activePower    float32
    reactivePower  float32
    powerFactor   float32
    energyExport   float32
    energyImport   float32
}
```
✅ Matches reference concept

### 6. Bus (Electrical Topology)

**Definition Needed:**
```
A Bus is a node in the electrical network where multiple conductors connect.
It represents a point in the network with a single voltage.

Properties:
- Nominal voltage (V)
- Current voltage (V)
- Connected elements

Relationship:
- Bus is part of electrical topology
- Multiple Breakers, Loads, Generators can connect to one Bus
- Bus IS a simulation model (topology)
```

**Current Implementation:**
```go
type BusEntity struct {
    nominalVoltage float32
    voltage       float32
    powerIn       float32
    powerOut      float32
}
```
⚠️ Gap: No connectivity information (which elements connect)

### 7. Electrical Topology

**Definition Needed:**
```
Electrical Topology describes how electrical elements connect.
It defines the network structure for power flow analysis.

Components:
- Buses (nodes)
- Branches (lines between buses)
- Breakers (switches in branches)

Relationship:
- Topology IS owned by Grid Model
- Devices connect to Buses
- Power flow depends on topology
```

**Gap:** No topology model in Forge

### 8. Power Flow Direction

**Definition Needed:**
```
Power Flow Direction defines how active and reactive power move through the network.

Convention:
- Positive P: Power flowing FROM source TO load
- Negative P: Power flowing FROM load TO source (generation)
- Import: Energy purchased from grid
- Export: Energy sold to grid

Relationship:
- Direction determined by sign of power
- Revenue Meters track import/export separately
```

**Current Implementation:**
```go
direction := "EXPORT"
if e.activePower < 0 {
    direction = "IMPORT"
}
```
✅ Matches reference

## Gap Summary

| Concept | Status | Priority |
|---------|--------|----------|
| Load | Implemented | Medium |
| Generator | Implemented | Medium |
| Breaker | Implemented | High |
| Transformer | Implemented | High |
| Meter | Implemented | Low |
| Bus | Implemented | High |
| Electrical Topology | **Missing** | High |
| Power Flow Direction | Implemented | Low |

## Architectural Alignment Issues

### Issue 1: World Entities vs Simulation Models

**Problem:**
The `world/electrical` entities mix concepts:

```go
// These are ENTITY types (like Virtual Devices)
- LoadEntity
- GeneratorEntity
- BreakerEntity
- MeterEntity

// These are MODEL types (Simulation Models)
- GridEntity (should be GridModel)
- BusEntity (topology model)
```

**Resolution:**
Reference architecture separates:
- Simulation Models (physics) - Grid, Sun, Weather
- Virtual Devices (equipment) - Meter, Inverter, Relay

Forge world package needs clarification:
- Keep electrical entities as World Entities (generic)
- Distinguish between physics and equipment

### Issue 2: Power Injection vs Power Flow

**Problem:**
Current implementation uses `InjectPower`/`WithdrawPower`:

```go
func (e *BusEntity) InjectPower(kw float32) {
    e.powerIn += kw
}
```

Reference architecture uses:
```go
// Grid Model records injections
grid.InjectActivePower(dcPower)

// Positive = injection, negative = withdrawal
```

**Resolution:**
Align on power sign convention:
- Positive: Injection (generation)
- Negative: Withdrawal (load)

## Recommendations

### 1. Add Electrical Topology Concept

Create topology model:
```go
type Topology struct {
    buses []Bus
    branches []Branch
}

type Bus struct {
    ID string
    NominalVoltage float32
    ConnectedElements []string
}

type Branch struct {
    From string  // bus ID
    To   string  // bus ID
    BreakerID string
}
```

### 2. Distinguish Simulation Model vs Entity

Current: All entities in world/
Should be:
- `models/electrical.go` - GridModel, TopologyModel
- `entities/electrical.go` - Load, Generator, Breaker, Meter

### 3. Add Power Source/Sink Interface

```go
type PowerInjection interface {
    InjectActivePower(kW float32)
    InjectReactivePower(kVAr float32)
}

type PowerExtraction interface {
    WithdrawActivePower(kW float32)
    WithdrawReactivePower(kVAr float32)
}
```

### 4. Update Glossary

Add definitions for:
- Load (Power Sink)
- Generator (Power Source)
- Breaker (Switchgear)
- Transformer (Voltage Transformation)
- Bus (Electrical Node)
- Electrical Topology

## Verification Checklist

- [ ] Load entity follows power sink convention
- [ ] Generator entity follows power source convention
- [ ] Breaker responds to trip/close commands, not direct state manipulation
- [ ] Transformer models voltage transformation with efficiency
- [ ] Meter matches Revenue Meter specification
- [ ] Bus represents electrical topology node
- [ ] Power flow direction uses standard convention
- [ ] Glossary updated with electrical terms

---

*Audit Date: 2026-07-13*
*Auditor: Architecture Review*
