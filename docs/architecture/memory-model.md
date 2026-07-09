# Memory Model

## Philosophy

**Memory is the source of truth.**

There are TWO distinct memory domains. These must never be confused.

## Two Memory Domains

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         Simulation Runtime                               │
│                                                                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                   │
│  │   Device A  │  │   Device B  │  │   Device C  │                   │
│  │  ┌───────┐  │  │  ┌───────┐  │  │  ┌───────┐  │                   │
│  │  │Device │  │  │  │Device │  │  │  │Device │  │                   │
│  │  │Memory │  │  │  │Memory │  │  │  │Memory │  │                   │
│  │  └───────┘  │  │  └───────┘  │  │  └───────┘  │                   │
│  │   (private) │  │   (private) │  │   (private) │                   │
│  └─────────────┘  └─────────────┘  └─────────────┘                   │
│         │                │                │                             │
│         └────────────────┼────────────────┘                             │
│                          │                                              │
│                          ▼                                              │
│               ┌───────────────────┐                                    │
│               │   Raw Ingest       │                                    │
│               │   Publisher        │                                    │
│               └───────────────────┘                                    │
└─────────────────────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                           MMA2                                           │
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐  │
│  │                   Operational Memory                              │  │
│  │  (shared, visible to Atlas-PPC, SCADA, HMIs, Historians)         │  │
│  │                                                                   │  │
│  │  Exposed via Modbus, DNP3, REST, MQTT...                         │  │
│  └─────────────────────────────────────────────────────────────────┘  │
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐  │
│  │   Real Devices via Replicator                                    │  │
│  │   Real Device → Replicator → Raw Ingest → MMA2                   │  │
│  └─────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────┘
```

### Domain 1: Device Memory

- **Owned by**: Each virtual device
- **Scope**: Private, internal
- **Access**: Device behaviors only
- **Purpose**: Internal device state and simulation logic
- **Lifetime**: Tied to device lifecycle

### Domain 2: Operational Memory (MMA2)

- **Owned by**: MMA2 appliance
- **Scope**: Shared, visible to all
- **Access**: Atlas-PPC, SCADA, HMIs, Historians
- **Purpose**: Plant-wide operational state
- **Lifetime**: Persistent across simulation sessions

## Why This Separation Exists

```
┌──────────────────────────────────────────────────────────────────────┐
│                    Real Device vs Virtual Device                       │
└──────────────────────────────────────────────────────────────────────┘

Real Device                          Virtual Device
     │                                   │
     ▼                                   ▼
Replicator                        Raw Ingest
     │                                   │
     └───────────────┬───────────────────┘
                     │
                     ▼
              ┌─────────────┐
              │    MMA2     │
              │  (shared    │
              │   state)    │
              └─────────────┘
                     │
                     ▼
         Atlas-PPC cannot distinguish
         real from virtual origin
```

### Key Insight

**Device Memory is NOT MMA2. MMA2 is NOT Device Memory.**

A device maintains its own internal memory. When appropriate, the device publishes selected operational values into MMA2 via Raw Ingest.

The simulator never writes Modbus registers directly. Raw Ingest is the official interface for publishing simulated operational data.

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
