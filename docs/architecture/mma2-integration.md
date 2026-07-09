# MMA2 Integration

## Overview

The Simulation Runtime integrates with MMA2 via Raw Ingest.

**MMA2** (Modbus Memory Appliance) is the operational memory appliance. It owns:
- Raw Modbus memory
- Protocol exposure (Modbus TCP, REST, MQTT)
- Access control
- State sealing

The Simulation Runtime publishes data to MMA2 via Raw Ingest. It does not implement Modbus servers.

## Two Memory Domains

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Simulation Runtime                                     │
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐ │
│  │                    Device Memory                                    │ │
│  │  (private, owned by each device)                                   │ │
│  └─────────────────────────────────────────────────────────────────┘ │
│                              │                                           │
│                              ▼                                           │
│   ┌─────────────────────────────────────────────────────────────────┐ │
│   │               Raw Ingest Publisher                                  │ │
│   └─────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────────┘
                              │
                              ▼ TCP (binary protocol)
┌─────────────────────────────────────────────────────────────────────────┐
│                              MMA2                                         │
│                                                                         │
│  ┌─────────────────────────────────────────────────────────────────┐ │
│  │                   Operational Memory                                │ │
│  │  (Port, UnitID) → Memory Areas                                    │ │
│  │  - Coils (bit)                                                    │ │
│  │  - Discrete Inputs (bit)                                          │ │
│  │  - Holding Registers (16-bit uint)                                │ │
│  │  - Input Registers (16-bit uint)                                  │ │
│  └─────────────────────────────────────────────────────────────────┘ │
│                              │                                           │
│                              ▼                                           │
│   Modbus TCP │ REST │ MQTT │ etc.                                      │
└─────────────────────────────────────────────────────────────────────────┘
```

## Raw Ingest Protocol

Raw Ingest is a **write-only TCP transport**.

### Protocol Format (v1)

```
Magic      (2 bytes)   = 'R' 'I' (0x52 0x49)
Version    (1 byte)    = 0x01
Area       (1 byte)    = 1 (Coils) | 2 (DiscreteInputs) | 3 (HoldingRegs) | 4 (InputRegs)
UnitID     (2 bytes)   = target unit (big-endian uint16)
Address    (2 bytes)   = start address (big-endian uint16)
Count      (2 bytes)   = number of values (big-endian uint16)
Payload    (variable)  = bit-packed or word-aligned data
```

### Response Codes

After each write, MMA2 returns **exactly 1 byte**:

| Code | Name | Meaning |
|------|------|---------|
| 0x00 | OK | Write committed |
| 0x10 | INVALID_MAGIC | Bad magic bytes |
| 0x11 | INVALID_VERSION | Unknown version |
| 0x12 | INVALID_AREA | Unknown area |
| 0x13 | INVALID_COUNT | Count is zero |
| 0x14 | INVALID_LENGTH | Payload length mismatch |
| 0x20 | MEMORY_NOT_FOUND | No memory for (Port, UnitID) |
| 0x21 | OUT_OF_BOUNDS | Address exceeds layout |
| 0x30 | INTERNAL_ERROR | Unexpected failure |

### Key Properties

- **Stateless**: Each connection carries exactly one write operation
- **Bypasses authority**: No policy enforcement
- **Bypasses state sealing**: Writes accepted even when sealed
- **Always available**: No special configuration needed on MMA2

## Integration Architecture

```go
// RawIngestPublisher publishes to MMA2
type RawIngestPublisher interface {
    Publish(area uint8, unitID uint16, address uint16, values []byte) error
}
```

### Device publishes to MMA2

```go
func (b *WeatherBehavior) Tick() {
    // Read from device memory
    irradiance := b.device.Memory().ReadFloat32("input_registers", irradianceAddr)
    
    // Write to device memory (internal)
    b.device.Memory().WriteFloat32("input_registers", powerAddr, power)
    
    // Publish to MMA2 via Raw Ingest
    // Convert float32 to uint16 (with scaling)
    rawValue := uint16(irradiance * 10) // Example scaling
    
    b.publisher.Publish(
        area: 4,           // Input Registers
        unitID: 1,         // Unit ID
        address: 0,        // Address
        values: toBigEndian(rawValue),
    )
}
```

## Configuration

MMA2 configuration defines memory layout:

```yaml
# MMA2 config (mma2.yaml)
listeners:
  - id: "simulation"
    listen: "0.0.0.0:5020"  # Separate port for simulation
    
    memory:
      - unit_id: 1
        input_registers:
          start: 0
          count: 100
        holding_registers:
          start: 0
          count: 100
```

Simulation Runtime connects to MMA2:

```go
// Runtime connects to MMA2
runtime.ConnectMMA2("localhost:5020")
```

## Implementation Notes

### No Compile-Time Dependency

The Simulation Runtime should NOT import MMA2 packages.

Implement Raw Ingest as a standalone TCP client:

```go
// rawingest/publisher.go
package rawingest

import (
    "encoding/binary"
    "net"
)

type Publisher struct {
    addr string
}

func NewPublisher(addr string) *Publisher {
    return &Publisher{addr: addr}
}

func (p *Publisher) Publish(area, unitID, address uint16, values []byte) error {
    conn, err := net.Dial("tcp", p.addr)
    if err != nil {
        return err
    }
    defer conn.Close()
    
    // Build packet
    pkt := make([]byte, 10+len(values))
    pkt[0] = 'R'
    pkt[1] = 'I'
    pkt[2] = 0x01        // version
    pkt[3] = uint8(area) // area
    binary.BigEndian.PutUint16(pkt[4:6], unitID)
    binary.BigEndian.PutUint16(pkt[6:8], address)
    binary.BigEndian.PutUint16(pkt[8:10], uint16(len(values)/2)) // count
    copy(pkt[10:], values)
    
    if _, err := conn.Write(pkt); err != nil {
        return err
    }
    
    resp := make([]byte, 1)
    conn.Read(resp)
    
    if resp[0] != 0x00 {
        return fmt.Errorf("raw ingest error: 0x%02x", resp[0])
    }
    return nil
}
```

### Area Constants

```go
const (
    AreaCoils          = 1
    AreaDiscreteInputs = 2
    AreaHoldingRegs    = 3
    AreaInputRegs     = 4
)
```

### Value Encoding

- **Bit areas** (Coils, Discrete Inputs): bits packed LSB-first, padded to byte boundary
- **Register areas** (Holding, Input): big-endian uint16 words, 2 bytes each

### Scaling

MMA2 stores raw 16-bit values. Simulation Runtime handles scaling:

```go
// Sensor outputs float32
irradiance := float32(850.5) // W/m²

// Scale to uint16
// Example: 0-2000 W/m² → 0-65535
scaled := (irradiance / 2000.0) * 65535.0
raw := uint16(scaled)

// Publish to MMA2
publisher.Publish(AreaInputRegs, unitID, address, toBigEndian(raw))
```

## Real vs Virtual Equivalence

```
Real Device                          Virtual Device
     │                                   │
     ▼                                   ▼
Replicator                        Simulation Runtime
     │                                   │
     └───────────────┬───────────────────┘
                     │
                     ▼ Raw Ingest
              ┌─────────────┐
              │    MMA2     │
              └─────────────┘
                     │
     ┌───────────────┼───────────────┐
     ▼               ▼               ▼
Modbus TCP        DNP3            REST
```

Atlas-PPC, SCADA, and HMIs cannot distinguish real from virtual origins.

## Summary

1. **MMA2 is external**: The Simulation Runtime connects to MMA2 as an external system
2. **Raw Ingest is the interface**: Binary TCP protocol for writing memory
3. **No protocol implementation**: Runtime does not implement Modbus servers
4. **No compile-time dependency**: Implement Raw Ingest client standalone
5. **Scaling is runtime responsibility**: MMA2 stores raw 16-bit values
