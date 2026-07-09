# Device Model

## Philosophy

**A device is the fundamental unit of simulation.**

A virtual industrial device is:
- A deterministic memory image
- Executable behavior
- One or more protocol interfaces

The device owns all of these.

## Device Anatomy

```
┌─────────────────────────────────────────────────────────────┐
│                         Device                               │
│                                                               │
│  ┌─────────────────────────────────────────────────────┐ │
│  │                     Memory Image                       │ │
│  │                                                      │ │
│  │  (source of truth, owned by device)                   │ │
│  └─────────────────────────────────────────────────────┘ │
│                                                               │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐ │
│  │ Behaviors  │  │ Protocols  │  │     Faults      │ │
│  │ (internal) │  │ (external) │  │  (modifiers)    │ │
│  └─────────────┘  └─────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

Memory is the core. Behaviors and protocols attach to it.

## Memory Ownership

Each device owns its memory. Memory is the single source of truth.

```go
type Device struct {
    memory    *MemoryImage    // Device owns this
    behaviors []Behavior       // Operate on memory
    protocols []Protocol       // Expose memory
    faults   []Fault         // Modify memory behavior
}
```

## Behaviors

Behaviors are logic that reads and writes device memory.

```go
func (b *PVModelBehavior) Tick() {
    irradiance := b.device.Memory().ReadFloat32("input_registers", 0)
    dcPower := irradiance * b.scaleFactor
    b.device.Memory().WriteFloat32("input_registers", 8, dcPower)
}
```

Behaviors never access other devices. Behaviors never call protocols.

## Protocols

Protocols are external views of device memory.

```go
type ModbusAdapter struct {
    device *Device  // Exposes device memory
}
```

A protocol adapter exposes the same memory that behaviors modify. There is only one memory image.

## Faults

Faults modify how device memory behaves.

```go
device.AddFault(NewFrozenValuesFault())
```

Faults never own state. They only modify memory access.

## Device Communication

**Devices never communicate directly.**

Devices communicate only by reading and writing memory. The runtime controls execution order.

```
Device A (Weather Station)
└── Writes irradiance → Device A Memory
└── Publishes irradiance → Raw Ingest → MMA2

Device B (PV Inverter)
├── Reads irradiance ← Device A Memory (external input)
└── Writes power → Device B Memory
└── Publishes power → Raw Ingest → MMA2

Device C (Revenue Meter)
├── Reads power ← Device B Memory (external input)
└── Writes energy → Device C Memory
└── Publishes energy → Raw Ingest → MMA2
```

Publishing to MMA2 is optional. Devices can maintain private state while exposing operational values.

## Device Types

Device types define the structure of devices. Device types come from plugins.

```
Energy Plugin defines:
├── revenue_meter
├── weather_station
├── pv_inverter
├── relay
└── grid

Water Plugin defines:
├── pump
├── valve
├── tank
└── flow_meter
```

The runtime knows only the generic Device interface. Device types are plugin-specific.

## Example: Revenue Meter

```
┌─────────────────────────────────────────┐
│            Revenue Meter                  │
│                                         │
│  ┌─────────────────────────────────┐  │
│  │           Memory                  │  │
│  │  Holding Registers (config)      │  │
│  │  Input Registers (measurements)   │  │
│  │  Coils (binary control)          │  │
│  │  Discrete Inputs (status)        │  │
│  └─────────────────────────────────┘  │
│                                         │
│  Behaviors:                             │
│  ├── Power Measurement                 │
│  └── Demand Calculation                │
│                                         │
│  Protocols:                            │
│  ├── Modbus TCP                       │
│  ├── DNP3                             │
│  └── REST API                         │
└─────────────────────────────────────────┘
```

## Key Principle

**A device is self-contained.**

A device owns its memory. Behaviors read and write that memory. Protocols expose that memory. Nothing outside the device touches its memory directly.
