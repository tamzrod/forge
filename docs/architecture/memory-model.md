# Memory Model

## Philosophy

**Memory is the source of truth.**

There are THREE distinct memory domains. These must never be confused.

## Three Memory Domains

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         Simulation Runtime                               │
│                                                                         │
│  ┌───────────────────────────────────────────────────────────────┐    │
│  │                     Simulation Models                           │    │
│  │                                                               │    │
│  │   ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐         │    │
│  │   │  Grid   │  │   Sun   │  │  Wind   │  │ Weather │         │    │
│  │   │  Model  │  │  Model  │  │  Model  │  │  Model  │         │    │
│  │   │  State  │  │  State  │  │  State  │  │  State  │         │    │
│  │   └─────────┘  └─────────┘  └─────────┘  └─────────┘         │    │
│  │                        (private RAM)                             │    │
│  └───────────────────────────────────────────────────────────────┘    │
│                                    │                                  │
│                                    ▼                                  │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                │
│  │   Device A  │  │   Device B  │  │   Device C  │                │
│  │  ┌───────┐  │  │  ┌───────┐  │  │  ┌───────┐  │                │
│  │  │Device │  │  │  │Device │  │  │  │Device │  │                │
│  │  │Memory │  │  │  │Memory │  │  │  │Memory │  │                │
│  │  └───────┘  │  │  └───────┘  │  │  └───────┘  │                │
│  │   (private) │  │   (private) │  │   (private) │                │
│  └─────────────┘  └─────────────┘  └─────────────┘                │
│         │                │                │                         │
│         └────────────────┼────────────────┘                         │
│                          │                                            │
│                          ▼                                            │
│               ┌───────────────────┐                                 │
│               │   Raw Ingest       │                                 │
│               │   Publisher        │                                 │
│               └───────────────────┘                                 │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                           MMA2                                        │
│                                                                     │
│  ┌───────────────────────────────────────────────────────────────┐ │
│  │                   Operational Memory                            │ │
│  │  (shared, visible to Atlas-PPC, SCADA, HMIs, Historians)        │ │
│  │                                                                   │ │
│  │  Exposed via Modbus, DNP3, REST, MQTT...                        │ │
│  └───────────────────────────────────────────────────────────────┘ │
│                                                                     │
│  ┌───────────────────────────────────────────────────────────────┐ │
│  │   Real Devices via Replicator                                  │ │
│  │   Real Device → Replicator → Raw Ingest → MMA2                │ │
│  └───────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────┘
```

### Domain 1: Model State (RAM)

- **Owned by**: Each simulation model
- **Scope**: Private, internal to model
- **Access**: Model methods only
- **Purpose**: Physical state (voltage, irradiance, temperature)
- **Lifetime**: Tied to simulation session
- **Exposure**: Never exposed through protocols

### Domain 2: Device Memory

- **Owned by**: Each virtual device
- **Scope**: Private, internal to device
- **Access**: Device behaviors only
- **Purpose**: Internal device state and simulation logic
- **Lifetime**: Tied to device lifecycle
- **Exposure**: May be published to MMA2 via Raw Ingest

### Domain 3: Operational Memory (MMA2)

- **Owned by**: MMA2 appliance
- **Scope**: Shared, visible to all
- **Access**: Atlas-PPC, SCADA, HMIs, Historians
- **Purpose**: Plant-wide operational state
- **Lifetime**: Persistent across simulation sessions
- **Exposure**: Exposed via Modbus, DNP3, REST, MQTT

## Why This Separation Exists

```
┌──────────────────────────────────────────────────────────────────────┐
│                 Three Domains, Three Purposes                          │
└──────────────────────────────────────────────────────────────────────┘

Simulation Model                     Virtual Device                      MMA2
     │                                    │                               │
     ▼                                    ▼                               ▼
 Private RAM                        Private Memory                    Shared Memory
     │                                    │                               │
     ▼                                    │                               │
 Physics (Grid, Sun, Wind)                │                               │
     │                                    │                               │
     └────────────── observes ────────────┘                               │
                                        │                               │
                                        ▼                               ▼
                              Device publishes                          External
                              to MMA2                                  Systems
```

### Key Insight

**Model State is NOT Device Memory. Device Memory is NOT MMA2. MMA2 is NOT Model State.**

| Domain | Never does |
|--------|-----------|
| **Model State** | Expose protocols, publish to MMA2 |
| **Device Memory** | Access other devices, expose protocols directly |
| **MMA2** | Know about physics, modify device memory |

A device observes models and maintains its own internal memory. When appropriate, the device publishes selected operational values into MMA2 via Raw Ingest.

The simulator never writes Modbus registers directly. Raw Ingest is the official interface for publishing simulated operational data.

## Model State Storage

Simulation Models store their state **directly in RAM**. This is different from devices.

```go
// GridModel stores physics in private RAM
type GridModel struct {
    // Physical state - private to the model
    voltage         float32 // V
    frequency       float32 // Hz
    theveninZReal   float32 // Ω
    theveninZImag   float32 // Ω
    shortCircuitMVA float32 // MVA
    
    // Power injection - reset each tick
    activePowerInjection   float32 // MW
    reactivePowerInjection float32 // MVAr
}

// Access through methods only
func (g *GridModel) Voltage() float32 {
    return g.voltage
}

func (g *GridModel) InjectActivePower(mw float32) {
    g.activePowerInjection += mw
}
```

### Why Models Don't Use MMA2

| Aspect | Model State | MMA2 |
|--------|-------------|------|
| **Protocols** | None | Modbus, DNP3 |
| **External Access** | Never | SCADA, PPC |
| **Engineering Units** | Native | Scaled |
| **Access Pattern** | Method calls | Register reads |

Models represent **physics**, not equipment. Physics doesn't expose Modbus.

### Why Models Don't Use Raw Ingest

| Aspect | Model State | Raw Ingest |
|--------|-------------|------------|
| **Purpose** | Internal simulation | External publishing |
| **Clients** | Devices (internal) | MMA2 (external) |
| **Protocol** | None | Custom wire format |

Raw Ingest is for **operational telemetry**, not physics simulation.

## Benefits of Memory as Foundation (Device Memory)

| Benefit | Explanation |
|---------|-------------|
| **Deterministic execution** | Same memory state produces same results, always |
| **Simple serialization** | Memory is already structured; no object graph to serialize |
| **Easy snapshots** | Freeze entire device state by copying memory |
| **Replay capability** | Record memory writes for deterministic replay |
| **Low coupling** | Behaviors don't know about MMA2 |
| **Cache-friendly** | Sequential memory access patterns |
| **Internal state** | Device can track private simulation state |

### Comparison to Object Graphs

Traditional simulation frameworks use object graphs:

```
Object A ──references──▶ Object B ──references──▶ Object C
```

This creates:
- Complex serialization
- Circular dependency risks
- Hidden state
- Difficult snapshots

Memory-centric design avoids these problems.

## Memory Ownership

**Each device owns its memory. There is no global memory.**

```
Device A owns its memory
Device B owns its memory
Device C owns its memory

No device can access another device's memory.
```

## Memory Structure

```go
type MemoryImage struct {
    regions map[string]*MemoryRegion
}

type MemoryRegion struct {
    Name   string
    Size   uint32
    Values []byte
}

type MemoryLocation struct {
    Region string
    Offset uint32
    Value  []byte
    Quality Quality
}
```

## Memory Regions

| Region | Access | Use |
|--------|--------|-----|
| **Holding Register** | Read/Write | Configuration parameters |
| **Input Register** | Read | Measured values |
| **Coil** | Read/Write | Binary control |
| **Discrete Input** | Read | Binary status |

## Device Memory Access

Behaviors read and write device memory:

```go
behavior.Tick() {
    // Read internal state
    irradiance := device.Memory().ReadFloat32("input_registers", irradianceAddr)
    
    // Compute
    power := irradiance * efficiency
    
    // Write internal state
    device.Memory().WriteFloat32("input_registers", powerAddr, power)
}
```

Behaviors may also publish to MMA2 via Raw Ingest:

```go
behavior.Tick() {
    // Read from device memory
    irradiance := device.Memory().ReadFloat32("input_registers", irradianceAddr)
    
    // Publish to operational memory (MMA2)
    publisher.Publish("weather/irradiance", irradiance, QualityGood)
}
```

## Quality Flags

Each memory location has a quality flag:

```go
type Quality uint8

const (
    QualityGood      Quality = 0x00
    QualityUncertain Quality = 0x40
    QualityBad      Quality = 0x80
    QualityOffline  Quality = 0x84
)
```

Quality flags indicate data validity. Faults set quality flags.

## Memory Isolation

Devices cannot access each other's memory:

```go
// This does not exist
func (b *Behavior) AccessOtherDevice(otherDevice *Device) {
    otherDevice.Memory().Read(...)  // Not possible
}
```

## Device Communication

Devices communicate only through the runtime's execution order:

```
Runtime executes Device A
└── Device A behavior writes irradiance → Device A Memory
└── Device A behavior publishes → Raw Ingest → MMA2

Runtime executes Device B
├── Device B behavior reads irradiance ← Device A Memory
└── Device B behavior writes power → Device B Memory
└── Device B behavior publishes → Raw Ingest → MMA2

Runtime executes Device C
├── Device C behavior reads power ← Device B Memory
└── Device C behavior writes energy → Device C Memory
└── Device C behavior publishes → Raw Ingest → MMA2
```

The runtime controls execution order. Devices pass data through their own memory. Publishing to MMA2 is optional and controlled by each device.

## No Global Memory System

There is no centralized memory:

```go
// Does not exist
runtime.GlobalMemory()
runtime.MemoryManager()
runtime.SharedState()
```

## Memory and Serialization

Memory is already structured for serialization:

```go
type Snapshot struct {
    DeviceID  string
    Timestamp time.Time
    Memory    []byte  // Direct memory dump
}
```

Snapshot and restore is a simple memory copy.

## Memory and Faults

Faults modify memory behavior:

- **Frozen**: Writes are ignored
- **Noise**: Reads return corrupted values
- **Offline**: All locations set to QualityOffline

Faults don't own state. They only intercept memory access.

## Encoding and Scaling (Device Responsibility)

Device memory stores raw bytes. Encoding and scaling are device responsibilities.

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         Device Memory (raw bytes)                          │
│                                                                         │
│  Regions: "sensors", "computed", "output"                              │
│  Types: uint8, uint16, uint32, float32 (raw bytes)                     │
│                                                                         │
│  Device memory does NOT know:                                           │
│  - Engineering units (°C, W/m², V, MW)                                 │
│  - Scaling factors                                                      │
│  - Register maps in MMA2                                               │
└─────────────────────────────────────────────────────────────────────────┘
                              │
                              │ Device encodes
                              ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                         MMA2 (raw uint16)                                │
│                                                                         │
│  Areas: HoldingRegisters, InputRegisters, Coils, DiscreteInputs         │
│  Values: raw uint16 or bits                                            │
│                                                                         │
│  MMA2 does NOT know:                                                   │
│  - Temperature, irradiance, power                                       │
│  - Scaling                                                             │
│  - What the values represent                                            │
└─────────────────────────────────────────────────────────────────────────┘
```

### Example: Weather Device Encoding

```go
// Device reads engineering value from memory
irradiance := device.Memory().ReadFloat32("sensors", 0) // 850.5 W/m²

// Device encodes for MMA2 (0-2000 W/m² → 0-65535)
scaled := (irradiance / 2000.0) * 65535.0
raw := encodeFloat32(scaled)

// Device writes to MMA2
rawIngest.WriteInputRegisters(unitID, 0, raw)
```

The runtime knows nothing about encoding, scaling, or engineering units. Only the device knows.
