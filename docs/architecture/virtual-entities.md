# Virtual Electrical Entities

## Purpose

Virtual Electrical Entities are behavioral placeholders for electrical components. They expose electrical behavior without being technology-specific.

## Design Philosophy

**Virtual entities are behavioral placeholders, not simplified equipment.**

Technology-specific implementations will extend or compose these entities later.

```
Virtual Generator
    ↓
Solar Generator / Wind Generator / Diesel Generator / Battery Generator

Virtual Load
    ↓
Residential Load / Industrial Load / City Load / Station Service
```

Forge reasons about **electrical behavior** rather than equipment names.

## Core Entities

### VirtualGenerator

Represents any source capable of injecting electrical power.

```go
type VirtualGenerator struct {
    // Power output
    activePower   float32 // kW
    reactivePower float32 // kVAr

    // Capacity
    ratedCapacity     float32 // kW
    availableCapacity float32 // kW

    // Status
    isOnline       bool
    isDispatchable bool
    rampRate       float32 // kW per minute
}
```

**Interface:**
```go
// Power
ActivePower() float32       // Current output (kW)
ReactivePower() float32     // Current Q (kVAr)
RatedCapacity() float32      // Nameplate capacity (kW)
AvailableCapacity() float32  // Current available (kW)

// Status
IsOnline() bool             // Is generator running?
IsDispatchable() bool       // Can it be dispatched?
RampRate() float32          // kW per minute

// Control
SetOnline(bool)             // Start/stop
SetTargetPower(float32)     // Dispatch target (kW)
SetAvailableCapacity(float32) // Available for dispatch (kW)
```

**Examples of what VirtualGenerator represents:**
- Solar Plant
- Wind Farm
- Diesel Generator
- Hydro Plant
- Battery (Discharging)
- Fuel Cell
- Gas Turbine

### VirtualLoad

Represents any consumer of electrical power.

```go
type VirtualLoad struct {
    // Power demand
    activePowerDemand   float32 // kW
    reactivePowerDemand float32 // kVAr

    // Status
    isConnected bool
    priority   int  // 1-10, higher = more important
}
```

**Interface:**
```go
// Power
ActivePowerDemand() float32    // Current demand (kW)
ReactivePowerDemand() float32  // Current Q demand (kVAr)

// Status
IsConnected() bool            // Is load connected?
Priority() int                // Load priority (1-10)

// Control
SetConnected(bool)            // Connect/disconnect
SetPowerDemand(float32)       // Set demand (kW)
SetPowerScale(float32)        // Scale factor (0.0-1.0)
SetPriority(int)              // Set priority (1-10)
```

**Examples of what VirtualLoad represents:**
- Residential Load
- Industrial Plant
- Factory
- Building
- City
- Station Service
- Auxiliary Load
- EV Charging
- Process Equipment

### VirtualMeter

Represents a virtual power meter at a point in the network.

```go
type VirtualMeter struct {
    // Measured values
    activePower    float32 // kW
    reactivePower  float32 // kVAr
    apparentPower float32 // kVA
    powerFactor   float32

    // Energy
    energyImport float32 // kWh
    energyExport float32 // kWh
}
```

**Interface:**
```go
// Power
ActivePower() float32       // Measured P (kW)
ReactivePower() float32      // Measured Q (kVAr)
ApparentPower() float32      // Measured S (kVA)
PowerFactor() float32        // Calculated PF

// Energy
EnergyImport() float32       // Total imported (kWh)
EnergyExport() float32       // Total exported (kWh)

// Control
SetMeasurements(P, Q, V, F)  // Update measurements
RecordEnergy(dt)             // Accumulate energy
```

## World Integration

### Generators Inject Power

Generators call `InjectPower()` to report injection:

```go
// Generator injects 100 kW
actual := generator.InjectPower(100.0)
```

### Loads Withdraw Power

Loads call `WithdrawPower()` to report consumption:

```go
// Load wants 50 kW
actual := load.WithdrawPower(50.0)
```

### Meter Measures Net Flow

```go
// Calculate net power
netPower := totalGen - totalLoad
meter.SetMeasurements(netPower, 0, 69000, 60)
```

## Usage Example

```go
// Create entities
solar := world.NewVirtualGenerator("gen-1", "Solar A", 500)
factory := world.NewVirtualLoad("load-1", "Factory", 400)
meter := world.NewVirtualMeter("pcc", "PCC Meter")

// Dispatch
solar.SetTargetPower(400) // Dispatch 400 kW
solar.SetAvailableCapacity(400) // Cloud-free

// Simulate
for {
    world.Tick(100 * time.Millisecond)
    
    // Calculate totals
    totalGen := solar.ActivePower()
    totalLoad := factory.ActivePowerDemand()
    netPower := totalGen - totalLoad
    
    // Update meter
    meter.SetMeasurements(netPower, 0, 69000, 60)
}
```

## Behavior vs Technology

### Why Separate Behavior from Technology?

| Aspect | Behavior (Virtual) | Technology (Specific) |
|--------|-------------------|---------------------|
| Focus | Power injection/consumption | How power is produced/consumed |
| Interface | Generic (P, Q, status) | Domain-specific (irradiance, fuel) |
| Use | Dispatch, optimization | Engineering, design |
| Extension | Compose with technology | Extend behavior |

### Example Extension

```go
// VirtualGenerator provides behavioral interface
type SolarGenerator struct {
    *world.VirtualGenerator
    irradiance float32
    efficiency float32
}

// Solar-specific logic
func (s *SolarGenerator) SetIrradiance(wm2 float32) {
    s.irradiance = wm2
    // Calculate available capacity from irradiance
    available := s.irradiance * s.area * s.efficiency
    s.SetAvailableCapacity(available)
}
```

## Design Rules

### Virtual Entities Know Nothing About:
- Topology (they don't know what they connect to)
- Other entities (they don't reference each other)
- Protocols (they don't communicate externally)
- Time (they don't own the clock)

### Virtual Entities Only Know:
- Their own electrical behavior
- Power output/demand
- Status (online/offline, connected/disconnected)
- Capacity limits

### The World Knows:
- How to aggregate multiple entities
- How to calculate net power
- How to determine import/export

## Future Specialization

### Generator Specializations

```
VirtualGenerator
├── SolarGenerator (irradiance, temperature, degradation)
├── WindGenerator (wind speed, turbulence, cut-out)
├── DieselGenerator (fuel, emissions, governor)
├── HydroGenerator (water flow, head, efficiency)
├── BatteryGenerator (SOC, temperature, cycle life)
└── FuelCellGenerator (fuel flow, efficiency)
```

### Load Specializations

```
VirtualLoad
├── ResidentialLoad (time-of-day profiles)
├── IndustrialLoad (process profiles)
├── CityLoad (weather correlation)
├── StationServiceLoad (auxiliary systems)
└── EVChargingLoad (arrival patterns, SOC)
```

## Summary

| Entity | Purpose | Interface |
|--------|---------|-----------|
| **VirtualGenerator** | Power source | Inject power, online/offline |
| **VirtualLoad** | Power consumer | Withdraw power, connected/disconnected |
| **VirtualMeter** | Power measurement | Measure P, Q, energy |

**The World determines:**
- Total Generation
- Total Consumption
- Net Power
- Import / Export

**Without knowing underlying technology.**

---

*Last Updated: 2026-07-13*
