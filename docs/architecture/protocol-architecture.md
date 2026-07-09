# Protocol Architecture

## Philosophy

**Protocols are external views of device memory.**

A protocol adapter is not part of the simulation itself. It is an external interface that exposes an existing memory image.

## Protocol Principles

A protocol adapter:

1. **Never owns state** - It only exposes memory
2. **Never owns behavior** - It only reads and writes memory
3. **Never synchronizes** - Multiple protocols expose the same memory
4. **Never transforms data** - It maps protocol concepts to memory

## Protocol vs Simulation

```
┌─────────────────────────────────────────────────────────────┐
│                      Simulation                              │
│                                                               │
│   ┌─────────────────────────────────────────────────────┐ │
│   │                 Device Memory                         │ │
│   │                                                      │ │
│   │  Behaviors write here                               │ │
│   └─────────────────────────────────────────────────────┘ │
│                          ▲                                  │
│                          │                                  │
│                    Behaviors                               │
└─────────────────────────────────────────────────────────────┘
                           │
                           │ exposed through
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                   External Systems                           │
│                                                               │
│   ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│   │  Modbus TCP   │  │     DNP3     │  │   REST API   │  │
│   └──────────────┘  └──────────────┘  └──────────────┘  │
│                                                               │
│   Protocols are external to the simulation                   │
└─────────────────────────────────────────────────────────────┘
```

## Memory Exposure

A protocol adapter exposes device memory:

```go
type ModbusAdapter struct {
    device *Device  // Exposes device memory
}

func (a *ModbusAdapter) ReadHoldingRegister(address uint16) uint16 {
    return a.device.Memory().ReadUint16("holding_registers", uint32(address))
}
```

All protocols expose the same memory. There is only one memory image.

## Multiple Protocols

A device can expose multiple protocol interfaces. All expose the same memory.

```
┌─────────────────────────────────────────────────────┐
│                   Device Memory                       │
│                                                     │
│   ┌───────────────────────────────────────────┐   │
│   │           Single Memory Image               │   │
│   └───────────────────────────────────────────┘   │
│                        ▲                            │
│   ┌────────────────────┼────────────────────┐     │
│   │                    │                    │     │
│   ▼                    ▼                    ▼     │
│ ┌──────┐        ┌──────┐         ┌──────┐     │
│ │Modbus│        │ DNP3 │         │ REST │     │
│ └──────┘        └──────┘         └──────┘     │
│                                                     │
│ Multiple views of the same memory                   │
└─────────────────────────────────────────────────────┘
```

## Supported Protocols

| Protocol | Domain | Mapping |
|----------|--------|---------|
| **Modbus TCP** | Industrial | Registers → Memory regions |
| **DNP3** | SCADA | Points → Memory locations |
| **REST** | Web | JSON → Memory read/write |
| **MQTT** | IoT | Topics → Memory publish |

## Protocol Interface

```go
type Protocol interface {
    Attach(device *Device)
    Detach()
}
```

The interface is minimal. A protocol only needs to know which device it exposes.

## No Protocol Synchronization

Protocols don't synchronize with each other. They all read from the same memory.

```
Protocol A reads memory
Protocol B reads memory
Protocol C reads memory

No synchronization.
Last write wins.
```

This is intentional. Synchronization would introduce coupling between protocols.

## No Protocol Caching

Protocols don't cache device state:

```go
// Wrong - caching creates inconsistency
type BadModbusAdapter struct {
    device *Device
    cache map[string][]byte  // No!
}

// Correct - always read memory
type GoodModbusAdapter struct {
    device *Device  // Exposes memory directly
}
```

Caching would create a second source of truth.

## Quality Propagation

Protocols propagate memory quality flags:

```go
func (a *ModbusAdapter) ReadInputRegister(address uint16) (uint16, error) {
    quality := a.device.Memory().Quality("input_registers", uint32(address))
    if quality != QualityGood {
        return 0, ModbusError{Code: quality}
    }
    return a.device.Memory().ReadUint16("input_registers", uint32(address)), nil
}
```

## Fault Reflection

Protocols reflect memory quality set by faults:

| Fault | Effect |
|-------|--------|
| **Offline** | Protocols return offline quality |
| **Bad Quality** | Protocols return bad quality |
| **Frozen** | Protocols return frozen values |

## Example: Device Exposes Protocols

```go
meter := runtime.CreateDevice(DeviceConfig{
    ID:   "meter-001",
    Type: "revenue_meter",
})

meter.ExposeProtocol("modbus", NewModbusAdapter("192.168.1.100:502"))
meter.ExposeProtocol("rest", NewRESTAdapter(":8080"))
```

The device decides which protocols to expose. The runtime knows nothing about protocols.

## Key Principle

**A protocol is an external view, not part of the simulation.**

Protocols expose memory. They don't own memory. They don't create state. They don't synchronize.
