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

### Low-Level Raw Ingest Client

The Raw Ingest client exposes only memory operations:

```go
// RawIngestClient provides low-level memory operations
// It does NOT know about float32, temperature, voltage, or any engineering unit
type RawIngestClient interface {
    WriteHoldingRegisters(unitID uint16, address uint16, values []byte) error
    WriteInputRegisters(unitID uint16, address uint16, values []byte) error
    WriteCoils(unitID uint16, address uint16, values []byte) error
    WriteDiscreteInputs(unitID uint16, address uint16, values []byte) error
}
```

### Device Owns Encoding

The device is responsible for encoding:

```go
// WeatherDevice encodes engineering values
func (d *WeatherDevice) Tick() {
    // Read from device memory (engineering value)
    irradiance := d.Memory().ReadFloat32("sensors", irradianceAddr)
    
    // Write to device memory (internal)
    d.Memory().WriteFloat32("computed", powerAddr, power)
    
    // Device encodes: float32 → uint16 with scaling
    // 0-2000 W/m² → 0-65535
    rawValue := uint16((irradiance / 2000.0) * 65535.0)
    
    // Device publishes to MMA2 (low-level memory operation)
    d.rawIngest.WriteInputRegisters(
        unitID:  1,
        address: 0,
        values:  toBigEndian(rawValue),
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

Implement Raw Ingest as a standalone TCP client with low-level operations only:

```go
// rawingest/client.go
// Low-level memory operations only - no engineering knowledge
package rawingest

import (
    "encoding/binary"
    "fmt"
    "net"
)

// Area constants
const (
    AreaCoils          = 1
    AreaDiscreteInputs = 2
    AreaHoldingRegs    = 3
    AreaInputRegs     = 4
)

// Client provides low-level Modbus memory operations
// It does NOT know about float32, temperature, voltage, or any engineering unit
type Client struct {
    addr string
}

func NewClient(addr string) *Client {
    return &Client{addr: addr}
}

func (c *Client) WriteInputRegisters(unitID, address uint16, values []byte) error {
    return c.write(AreaInputRegs, unitID, address, values)
}

func (c *Client) WriteHoldingRegisters(unitID, address uint16, values []byte) error {
    return c.write(AreaHoldingRegs, unitID, address, values)
}

func (c *Client) WriteCoils(unitID, address uint16, values []byte) error {
    return c.write(AreaCoils, unitID, address, values)
}

func (c *Client) WriteDiscreteInputs(unitID, address uint16, values []byte) error {
    return c.write(AreaDiscreteInputs, unitID, address, values)
}

func (c *Client) write(area uint8, unitID, address uint16, values []byte) error {
    conn, err := net.Dial("tcp", c.addr)
    if err != nil {
        return err
    }
    defer conn.Close()
    
    // Build packet
    pkt := make([]byte, 10+len(values))
    pkt[0] = 'R'
    pkt[1] = 'I'
    pkt[2] = 0x01        // version
    pkt[3] = area
    binary.BigEndian.PutUint16(pkt[4:6], unitID)
    binary.BigEndian.PutUint16(pkt[6:8], address)
    binary.BigEndian.PutUint16(pkt[8:10], uint16(len(values)/2)) // count (uint16 words)
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

### Encoding is Device Responsibility

Devices handle encoding. The runtime does not know about float32, scaling, or engineering units:

```go
// WeatherDevice - device owns encoding
func (d *WeatherDevice) Tick() {
    // Read engineering value from device memory
    irradiance := d.Memory().ReadFloat32("sensors", irradianceAddr)
    
    // Device encodes: float32 → uint16
    // 0-2000 W/m² → 0-65535
    scaled := (irradiance / 2000.0) * 65535.0
    raw := make([]byte, 2)
    binary.BigEndian.PutUint16(raw, uint16(scaled))
    
    // Device publishes via low-level Raw Ingest
    d.rawIngest.WriteInputRegisters(1, 0, raw)
}
```

### Value Encoding Rules

- **Bit areas** (Coils, Discrete Inputs): bits packed LSB-first, padded to byte boundary
- **Register areas** (Holding, Input): big-endian uint16 words, 2 bytes each

### What This Is NOT

The Raw Ingest client does NOT provide:
- `WriteFloat32()` - Device encodes float32
- `WriteTemperature()` - Device owns temperature semantics
- `WriteVoltage()` - Device owns voltage semantics
- Any engineering unit conversion

These belong in device implementations.

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
