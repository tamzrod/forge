# Runtime

## Philosophy

**The runtime hosts devices. That's all.**

The runtime provides common infrastructure. It contains no domain knowledge.

## Runtime Responsibilities

The runtime provides:

| Component | Purpose |
|-----------|---------|
| **Scheduler** | Advances simulation time |
| **Simulation Clock** | Tracks elapsed time |
| **Device Registry** | Holds loaded devices |
| **Plugin Loader** | Loads device types |
| **Configuration** | Provides settings |

That's the entire runtime.

## What the Runtime Does Not Do

The runtime explicitly does not:

- Own device memory
- Execute device behaviors
- Handle protocols
- Manage faults
- Own business concepts

## Runtime Structure

```go
type Runtime struct {
    scheduler Scheduler
    clock    SimulationClock
    devices  DeviceRegistry
    plugins  PluginLoader
    config   Config
}
```

The runtime is intentionally small.

## No Memory Management

There is no memory manager. Memory belongs to devices.

```go
// Runtime does not manage memory
type Runtime struct {
    scheduler   Scheduler
    clock       SimulationClock
    devices     DeviceRegistry
    plugins     PluginLoader
    // No memory manager
}

// Devices own memory
type Device struct {
    memory *MemoryImage  // Device owns this
}
```

## No Behavior Execution

The runtime does not execute behaviors. Devices execute their own behaviors.

```go
// Runtime just asks
func (r *Runtime) tick() {
    for _, device := range r.devices.All() {
        device.Tick()  // Device executes its behaviors
    }
    r.clock.Advance()
}

// Device executes
func (d *Device) Tick() {
    for _, behavior := range d.behaviors {
        behavior.Tick()
    }
}
```

## No Protocol Handling

The runtime does not handle protocols. Devices expose their own protocols.

```go
// Device exposes protocols
meter := runtime.Device("meter-001")
meter.ExposeProtocol("modbus", NewModbusAdapter())

// Runtime knows nothing about protocols
```

## No Fault Management

The runtime does not manage faults. Devices own their faults.

```go
device.AddFault(NewCommunicationLossFault())
```

## Domain Independence

The runtime knows nothing about:

- Energy
- Water
- Manufacturing
- Any industrial domain

All domain knowledge lives in plugins.

## Adding New Domains

Adding a new domain requires only:

1. New plugins with device types
2. No runtime changes

```
New Domain Plugin
├── Device Type A
├── Device Type B
└── Device Type C

Runtime: unchanged
```

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Simulation Runtime                         │
│                                                               │
│  Scheduler ──▶ Simulation Clock                              │
│                                                               │
│  Plugin Loader ──▶ Device Registry                          │
│                                                               │
│  Configuration                                               │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

The runtime is the smallest component. Devices are the system.

## Key Principle

**The runtime hosts devices. Devices own memory, behaviors, and protocols.**

The runtime disappears into the background.
