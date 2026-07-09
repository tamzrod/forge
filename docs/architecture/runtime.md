# Runtime

## Purpose

The runtime hosts the Virtual Industrial Laboratory. It provides scheduling, model management, device management, and Raw Ingest publishing.

The runtime contains no domain knowledge—all industrial specifics live in plugins.

## Philosophy

**The runtime hosts simulation models and devices. That's all.**

## Runtime Responsibilities

The runtime provides:

| Component | Purpose |
|-----------|---------|
| **Scheduler** | Advances simulation time deterministically |
| **Simulation Clock** | Tracks elapsed time |
| **Model Registry** | Holds simulation models |
| **Device Registry** | Holds loaded devices |
| **Plugin Loader** | Loads device types and model types |
| **Raw Ingest Publisher** | Publishes to MMA2 |
| **Configuration** | Provides settings |

That's the entire runtime—small enough to be easily understood and maintained.

## What the Runtime Does Not Do

The runtime explicitly does not:

- Own model state (models own their state)
- Own device memory (devices own their memory)
- Execute model behaviors (models execute their behaviors)
- Execute device behaviors (devices execute their behaviors)
- Expose protocols (MMA2 exposes protocols)
- Own operational memory (MMA2 owns operational memory)

## Runtime Structure

```go
type Runtime struct {
    scheduler    Scheduler
    clock        SimulationClock
    models       map[ModelID]Model
    devices      map[DeviceID]*Device
    plugins      PluginLoader
    rawIngest    RawIngestPublisher  // Publishes to MMA2
    config       Config
}
```

The runtime is intentionally small.

## Model Management

The runtime manages simulation models. Models represent the physical world.

```go
// Runtime manages models
type Runtime struct {
    scheduler    Scheduler
    models       map[ModelID]Model  // Simulation Models
    devices      map[DeviceID]*Device
}

// Create models
runtime.CreateGridModel("main-grid")
runtime.CreateSunModel("solar-sun")
runtime.CreateWindModel("wind-farm")

// Access models
grid := runtime.Model("main-grid")
```

## Device Memory Management

There is no memory manager. Memory belongs to devices.

```go
// Devices own memory
type Device struct {
    memory *MemoryImage  // Device owns this
}
```

## Execution Order

The scheduler executes models before devices:

```go
func (s *Scheduler) tick() {
    // 1. Models evolve first (physics)
    for _, model := range s.models {
        model.Tick()
    }
    
    // 2. Devices observe models and update memory
    for _, device := range s.devices {
        device.Tick()
    }
    
    // 3. Advance clock
    s.clock.Advance()
}
```

## Domain Independence

The runtime knows nothing about:

- Energy, Water, Manufacturing
- Grid, Sun, Wind, Weather
- Any industrial domain
- float32, int32, or any data type
- Engineering units
- Scaling
- Register layouts
- Application semantics

All domain knowledge lives in plugins.

## Separation of Responsibilities

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           Simulation Runtime                              │
│                                                                         │
│  Scheduler ──▶ Simulation Clock                                        │
│  Plugin Loader ──▶ Model Registry │ Device Registry                    │
│  Raw Ingest Publisher ──▶ MMA2                                        │
│  Configuration                                                         │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                       Simulation Models                                   │
│                                                                         │
│  Grid Model │ Sun Model │ Wind Model │ Weather Model                   │
│                                                                         │
│  - Represent physics (not equipment)                                   │
│  - Store state in private RAM                                          │
│  - Never expose protocols                                              │
│  - Never publish to MMA2                                               │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                          Virtual Devices                                  │
│                                                                         │
│  Weather Device │ PV Device │ Meter Device                             │
│                                                                         │
│  - Represent equipment (not physics)                                   │
│  - Observe models through behaviors                                     │
│  - Own device memory                                                   │
│  - Publish to MMA2 via Raw Ingest                                      │
└─────────────────────────────────────────────────────────────────────────┘
```

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Simulation Runtime                         │
│                                                               │
│  Scheduler ──▶ Simulation Clock                              │
│                                                               │
│  Plugin Loader ──▶ Model Registry ──▶ Simulation Models     │
│                   │                                         │
│                   └──▶ Device Registry ──▶ Virtual Devices  │
│                                                               │
│  Raw Ingest Publisher ──▶ MMA2 ──▶ Protocols               │
│                                                               │
│  Configuration                                               │
└─────────────────────────────────────────────────────────────┘
```

The runtime is the smallest component. Models represent physics. Devices observe physics. MMA2 exposes operational telemetry.

## Key Principle

**The runtime hosts models and devices. Models represent physics. Devices observe physics and publish to MMA2. MMA2 owns protocols.**

The runtime disappears into the background.

## Model Types

| Domain | Models |
|--------|--------|
| **Energy** | Grid, Sun, Wind, Weather |
| **Water** | Reservoir, Hydraulic Network, River |
| **Manufacturing** | Factory Power, Compressed Air, Conveyor Physics |

## Device Types

| Domain | Devices |
|--------|---------|
| **Energy** | Weather Station, PV Inverter, Revenue Meter, Relay |
| **Water** | Pump, Valve, Flow Meter, Tank Sensor |
| **Manufacturing** | PLC, Robot, Conveyor, Sensor |
