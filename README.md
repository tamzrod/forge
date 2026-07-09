# Industrial Simulation Runtime

**A Virtual Industrial Laboratory for industrial software development.**

This project provides a deterministic virtual industrial environment for developing, testing, commissioning, and training industrial software through realistic virtual industrial environments.

## One-Minute Summary

```
Simulation Models (physics) → Virtual Devices (equipment) → MMA2 (telemetry) → Applications
```

A virtual device owns deterministic memory.
Behaviors observe models and modify memory.
Protocols expose memory through MMA2.
The runtime hosts models, devices, and schedules.
Everything else is a plugin.

## Use Cases

- Software development
- Controller development
- SCADA development
- Protocol integration
- Factory Acceptance Testing (FAT)
- Commissioning
- Training
- Education
- Demonstrations

## Core Principles

1. **Devices own memory** - Every device owns its memory image
2. **Behaviors modify memory** - Logic reads from models, writes to device memory
3. **Protocols expose memory** - External systems read device memory through MMA2
4. **Devices never communicate directly** - Devices observe models, publish results
5. **Simulation Models represent physics** - Grid, Sun, Wind, Weather
6. **Runtime provides hosting** - Hosting, scheduling, plugin loading
7. **Plugins provide domain knowledge** - New domains add model types and device types

## Architecture

```
Simulation Runtime
      Scheduler / Clock / Plugin Loader / Model Registry
                         |
                         v
Simulation Models
      Grid | Sun | Wind | Weather | Reservoir | etc.
      Physical world - observed by devices
                         |
                         v
Virtual Devices
      Revenue Meter | Weather Station | PV Inverter | etc.
      Observe models, publish to MMA2
                         |
                         v
MMA2
      Operational memory - Modbus, DNP3, REST, MQTT
```

## Quick Start

```go
runtime, _ := forge.NewRuntime(forge.Config{TickInterval: 250 * time.Millisecond})
runtime.LoadPlugins("./plugins/energy")

// Create simulation models (physical world)
runtime.CreateGridModel("main-grid")
runtime.CreateSunModel("solar-sun")

// Create devices
runtime.CreateDevice("meter-001", "revenue_meter", memRegions)

runtime.Run(context.Background())
```

## Determinism

Execution is deterministic: same inputs produce same outputs, every time. This enables reproducible testing and training scenarios.

## Documentation

See `docs/architecture/` for full documentation:

- [Vision](docs/architecture/vision.md) - Project purpose and philosophy
- [Architecture Overview](docs/architecture/overview.md) - System architecture
- [Simulation Models](docs/architecture/simulation-models.md) - Physical world
- [Device Model](docs/architecture/device-model.md) - Virtual devices
- [Runtime](docs/architecture/runtime.md) - Runtime architecture

## License

See [LICENSE](LICENSE)
