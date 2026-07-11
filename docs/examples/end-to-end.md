# End-to-End Example

> **Architecture Validation:** This example demonstrates the complete Forge architecture and validates that all components work together correctly.

## Purpose

The end-to-end example validates the Forge architecture by running a complete simulation with:

- Simulation models representing the physical world
- Virtual devices with behaviors that observe models
- Device memory as the source of truth
- Proper execution lifecycle (init → run → shutdown)

## Location

```
examples/complete/main.go
```

## What This Example Demonstrates

### Architecture Layers

```
┌─────────────────────────────────────────────────────────────────────┐
│                      Simulation Runtime                                │
│  Scheduler / Clock / Device Registry / Model Registry                │
└─────────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      Simulation Models                               │
│  Grid │ Sun │ Weather │ Wind │ Reservoir                            │
│  Physical world - private RAM - observed by firmware                  │
└─────────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      Virtual Devices                                  │
│  Weather Station │ PV Inverter │ Revenue Meter                       │
│  Own memory, execute behaviors, observe models                       │
└─────────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                   Communication Interfaces                             │
│  (Would serialize memory for external systems)                      │
└─────────────────────────────────────────────────────────────────────┘
```

### Data Flow

```
Simulation Models (physical truth)
        │
        │ Behaviors observe models
        ▼
Virtual Firmware (behaviors)
        │
        │ Behaviors write measurements
        ▼
Device Memory (source of truth)
        │
        │ Protocols expose memory
        ▼
MMA2 / SCADA / External Systems
```

### Component Responsibilities

| Component | Responsibility | Owns |
|-----------|---------------|------|
| Runtime | Hosting, scheduling | Scheduler, Clock, Registries |
| Grid Model | Electrical grid physics | Voltage, Frequency |
| Sun Model | Solar position and irradiance | Irradiance, Elevation |
| Weather Model | Ambient conditions | Temperature, Humidity |
| Weather Station | Measure weather | Weather Memory |
| PV Inverter | Convert solar to AC | Power Memory |
| Revenue Meter | Measure power | Power Memory |

## Running the Example

```bash
cd examples/complete
go run main.go
```

## Expected Output

```
=== Forge End-to-End Example ===

Step 1: Initialize Runtime
  - Tick interval: 100ms
  - Max devices: 100

Step 2: Create Simulation Models
  - Grid model created: main-grid
  - Sun model created: solar-sun
  - Weather model created: ambient-weather
  - Wind model created: wind-field

Step 3: Create Virtual Devices
  - Weather station created: weather-station-001
  - PV Inverter created: pv-inverter-001
  - Revenue Meter created: revenue-meter-001

Step 4: Initial Model State
  Grid:    Voltage=480.0V, Frequency=60.00Hz
  Sun:     Irradiance=1000.0W/m², Elevation=45.0°
  Weather: Temperature=25.0°C, Humidity=50.0%
  Wind:    Speed=5.0m/s, Direction=0°

Step 5: Run Simulation
  Running 10 ticks...

Step 6: Final Device Memory State
  Weather Station:
    - Temperature: 25.x°C
    - Humidity: 50.x%
    - Pressure: 1013.x hPa
  PV Inverter:
    - DC Power: xxx.xW
    - AC Power: xxx.xW
    - Efficiency: 95.0%
  Revenue Meter:
    - Voltage: 480.0V
    - Current: 10.0A
    - Power: 4800.0W

Step 7: Verify Determinism
  Determinism means: same inputs → same outputs
  The architecture guarantees:
    - Devices tick in registration order
    - Behaviors tick in registration order
    - Models tick in registration order
    - No unseeded randomness

=== Architecture Validated ===
```

## Key Architecture Principles Demonstrated

### 1. Memory as Source of Truth

Behaviors write to device memory, not directly to external systems:

```go
// Observe model
temperature := weatherModel.Temperature()

// Write to memory (source of truth)
b.device.Memory().WriteFloat32("input_registers", 0, temperature)
```

### 2. Devices Observe Models

Devices access models through the `Model()` method:

```go
// Get model through device
sun := b.device.Model("solar-sun")
sunModel := sun.(*models.SunModel)

// Use model
irradiance := sunModel.Irradiance()
```

### 3. Clean Separation of Concerns

| Layer | Does | Does Not |
|-------|------|----------|
| Models | Physics calculations | Protocols, memory |
| Devices | Memory ownership | Direct communication |
| Behaviors | Logic execution | Own state outside memory |
| Protocols | Memory serialization | Engineering calculations |

### 4. Deterministic Execution

The architecture guarantees reproducibility:

```go
// Same initial conditions
grid.SetVoltage(480.0)
sun.SetIrradiance(1000.0)

// Same execution
rt.Run(ctx)

// Same results (within floating point tolerance)
```

## Extending the Example

### Adding a New Model

```go
// 1. Create the model
reservoir := rt.CreateReservoirModel("tank-001", 1000.0)

// 2. Model automatically registers with scheduler
// 3. Devices can now observe it
```

### Adding a New Device

```go
// 1. Define memory regions
memRegions := map[string]uint32{
    "input_registers": 20,
}

// 2. Create the device
tank := rt.CreateDevice("tank-001", "storage_tank", memRegions)

// 3. Add behaviors
tank.AddBehavior(&TankLevelBehavior{})
```

### Adding a New Behavior

```go
// 1. Implement the Behavior interface
type TankLevelBehavior struct {
    device interface{}
}

func (b *TankLevelBehavior) ID() string { return "tank_level" }
func (b *TankLevelBehavior) Attach(d interface{}) { b.device = d }
func (b *TankLevelBehavior) Detach()         { b.device = nil }

func (b *TankLevelBehavior) Tick() {
    // 2. Observe models
    reservoir := b.device.Model("tank-001")
    
    // 3. Write to memory
    b.device.Memory().WriteFloat32("input_registers", 0, reservoir.Level())
}
```

## Architecture Validation Checklist

This example validates:

- [x] Runtime initialization and shutdown
- [x] Simulation model creation and ticking
- [x] Device creation with memory regions
- [x] Behaviors observing models
- [x] Behaviors writing device memory
- [x] Execution loop with proper shutdown
- [x] Deterministic execution

## References

- [Architecture Overview](../architecture/overview.md)
- [Device Model](../architecture/device-model.md)
- [Memory Model](../architecture/memory-model.md)
- [Simulation Models](../architecture/simulation-models.md)
- [Behavior Model](../architecture/behavior-model.md)

---

*This example is the authoritative reference for Forge architecture usage.*
