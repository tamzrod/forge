# Memory Model

## Philosophy

**Memory is the source of truth.**

Memory-centric design is an architectural decision, not an implementation detail.

## Why Memory as Foundation

A virtual device is fundamentally a memory image. Behaviors read and write memory. Protocols expose memory. This creates a simple, predictable system.

### Benefits

| Benefit | Explanation |
|---------|-------------|
| **Deterministic execution** | Same memory state produces same results, always |
| **Simple serialization** | Memory is already structured; no object graph to serialize |
| **Easy snapshots** | Freeze entire device state by copying memory |
| **Replay capability** | Record memory writes for deterministic replay |
| **Protocol independence** | Any protocol maps naturally to memory regions |
| **Modbus compatibility** | Registers map directly to memory regions |
| **DNP3 compatibility** | Points map directly to memory addresses |
| **Low coupling** | Behaviors and protocols don't know about each other |
| **Cache-friendly** | Sequential memory access patterns |
| **Single source of truth** | One memory image, many protocol views |

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

## Memory Access

Behaviors write to memory:

```go
behavior.Tick() {
    power := voltage * current
    device.Memory().WriteFloat32("input_registers", powerAddr, power)
}
```

Protocols read from memory:

```go
adapter.HandleRead(address) {
    return device.Memory().Read("input_registers", address)
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

## Device Communication Through Memory

Devices communicate only through the runtime's execution order:

```
Runtime executes Device A
└── Device A behavior writes irradiance → Device A Memory

Runtime executes Device B
├── Device B behavior reads irradiance ← Device A Memory
└── Device B behavior writes power → Device B Memory

Runtime executes Device C
├── Device C behavior reads power ← Device B Memory
└── Device C behavior writes energy → Device C Memory
```

The runtime controls execution order. Devices pass data through their own memory.

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
