# Industrial Simulation Runtime

A generic runtime for virtual industrial devices. Written in Go.

**The devices are the system. The runtime only hosts them.**

## One-Minute Summary

```
A virtual device owns deterministic memory.
Behaviors modify that memory.
Protocols expose that memory.
The runtime simply hosts and schedules devices.
Everything else is a plugin.
```

## Core Principles

1. **Devices own memory** - Every device owns its memory image
2. **Behaviors modify memory** - Logic reads from and writes to device memory
3. **Protocols expose memory** - External systems read device memory
4. **Devices never communicate directly** - Communication through memory
5. **Runtime provides infrastructure** - Hosting, scheduling, plugin loading
6. **Plugins provide domain knowledge** - New domains add device types

## Why Memory as Foundation

- Deterministic execution
- Simple serialization
- Easy snapshots and replay
- Protocol independence
- Low coupling

## Architecture

```
┌─────────────────────────────────────────┐
│              Runtime                       │
│  Scheduler / Clock / Plugin Loader      │
└─────────────────────────────────────────┘
                    │
      ┌─────────────┼─────────────┐
      ▼             ▼             ▼
┌───────────┐ ┌───────────┐ ┌───────────┐
│  Device   │ │  Device   │ │  Device   │
│           │ │           │ │           │
│ ┌───────┐ │ │ ┌───────┐ │ │ ┌───────┐│
│ │ Memory │ │ │ │ Memory │ │ │ │ Memory ││
│ └───────┘ │ │ └───────┘ │ │ └───────┘│
│ ┌───────┐ │ │ ┌───────┐ │ │ ┌───────┐│
│ │Behavior│ │ │ │Behavior│ │ │ │Behavior││
│ └───────┘ │ │ └───────┘ │ │ └───────┘│
│ ┌───────┐ │ │ ┌───────┐ │ │ ┌───────┐│
│ │Protocol│ │ │ │Protocol│ │ │ │Protocol││
│ └───────┘ │ │ └───────┘ │ │ └───────┘│
└───────────┘ └───────────┘ └───────────┘
```

## Quick Start

```go
runtime, _ := forge.NewRuntime(forge.Config{TickInterval: 250 * time.Millisecond})
runtime.LoadPlugins("./plugins/energy")

runtime.CreateDevices([]forge.DeviceConfig{
    {ID: "meter-001", Type: "revenue_meter"},
})

runtime.Device("meter-001").ExposeProtocol("modbus", NewModbusAdapter())
runtime.Run(context.Background())
```

## Device Types

```
Energy Plugin
├── Revenue Meter
├── Weather Station
├── PV Inverter
├── Relay
└── Grid

Water Plugin
├── Pump
├── Valve
├── Tank
└── Flow Meter
```

## Documentation

See `docs/architecture/` for full documentation.

## License

See [LICENSE](LICENSE)
