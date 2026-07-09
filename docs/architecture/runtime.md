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
| **Raw Ingest Publisher** | Publishes to MMA2 |
| **Configuration** | Provides settings |

That's the entire runtime.

## What the Runtime Does Not Do

The runtime explicitly does not:

- Own device memory (devices own their memory)
- Execute device behaviors (devices execute their behaviors)
- Expose protocols (MMA2 exposes protocols)
- Own operational memory (MMA2 owns operational memory)

## Runtime Structure

```go
type Runtime struct {
    scheduler    Scheduler
    clock        SimulationClock
    devices      DeviceRegistry
    plugins      PluginLoader
    rawIngest    RawIngestPublisher  // Publishes to MMA2
    config       Config
}
```

The runtime is intentionally small.

## No Device Memory Management

There is no memory manager. Memory belongs to devices.

```go
// Runtime does not manage memory
type Runtime struct {
    scheduler   Scheduler
    clock       SimulationClock
    devices     DeviceRegistry
    plugins     PluginLoader
    rawIngest   RawIngestPublisher
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

The runtime does not handle protocols. MMA2 handles protocols.

```go
// Runtime connects to MMA2
runtime.ConnectRawIngest("mma2:8080")

// Runtime knows nothing about protocols - MMA2 owns them
```

## Raw Ingest Connection

The runtime connects to MMA2 via Raw Ingest:

```go
// Connect
runtime.ConnectRawIngest(endpoint string) error

// Disconnect
runtime.Disconnect()
```

Behaviors access the Raw Ingest publisher through their device.

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
│  Raw Ingest Publisher ──▶ MMA2                             │
│                                                               │
│  Configuration                                               │
│                                                               │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                           MMA2                                │
│                                                               │
│  Operational Memory                                          │
│  Modbus, DNP3, REST, MQTT...                               │
└─────────────────────────────────────────────────────────────┘
```

The runtime is the smallest component. Devices are the system. MMA2 is separate.

## Key Principle

**The runtime hosts devices and publishes to MMA2. Devices own memory. MMA2 owns protocols.**

The runtime disappears into the background.
