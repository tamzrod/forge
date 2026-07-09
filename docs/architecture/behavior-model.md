# Behavior Model

## Philosophy

**Behaviors are logic that reads and writes device memory.**

A device owns its behaviors. Behaviors may optionally publish to MMA2 via Raw Ingest.

## Behavior Contract

A behavior:

1. Reads from device memory
2. Computes new values
3. Writes to device memory
4. Optionally publishes to MMA2 via Raw Ingest
5. Never accesses other devices
6. Never calls protocols

## Data Flow

```
Behavior → Device Memory → Behavior
           (reads)           (writes)
               │
               ▼
         Raw Ingest
               │
               ▼
             MMA2
```

## Device Owns Behaviors

```go
type Device struct {
    behaviors []Behavior
}

func (d *Device) Tick() {
    for _, behavior := range d.behaviors {
        behavior.Tick()
    }
}
```

The device executes its own behaviors. The runtime tells the device to tick.

## Behavior Interface

```go
type Behavior interface {
    ID() string
    Attach(device *Device)
    Detach()
    Tick()
}
```

## Example: Weather Behavior

```go
type WeatherBehavior struct {
    device    *Device
    rawIngest RawIngestClient
    unitID    uint16
}

func (b *WeatherBehavior) Tick() {
    // Read from device memory (engineering values)
    irradiance := b.device.Memory().ReadFloat32("sensors", irradianceAddr)
    temperature := b.device.Memory().ReadFloat32("sensors", temperatureAddr)
    
    // Compute (internal simulation logic)
    // ...
    
    // Write to device memory
    b.device.Memory().WriteFloat32("computed", computedAddr, computedValue)
    
    // Device encodes: float32 → uint16 with scaling
    // Device publishes via low-level Raw Ingest
    b.rawIngest.WriteInputRegisters(b.unitID, 0, encodeFloat32(irradiance, 0, 2000))
    b.rawIngest.WriteInputRegisters(b.unitID, 2, encodeFloat32(temperature, -50, 50))
}
```

## Example: PV Model Behavior

```go
type PVModelBehavior struct {
    device    *Device
    rawIngest RawIngestClient
    unitID    uint16
}

func (b *PVModelBehavior) Tick() {
    irradiance := b.device.Memory().ReadFloat32("input", irradianceAddr)
    temperature := b.device.Memory().ReadFloat32("input", temperatureAddr)
    
    dcPower := irradiance * b.scaleFactor * (1 - 0.004*(temperature-25))
    
    // Write to device memory
    b.device.Memory().WriteFloat32("output", powerAddr, dcPower)
    
    // Device encodes and publishes
    b.rawIngest.WriteInputRegisters(b.unitID, 10, encodeFloat32(dcPower, 0, 500))
}

// Device-owned encoding helper
func encodeFloat32(value float32, min, max float32) []byte {
    scaled := (value - min) / (max - min) * 65535.0
    raw := make([]byte, 2)
    binary.BigEndian.PutUint16(raw, uint16(scaled))
    return raw
}
```

**Note:** Encoding and scaling are device responsibilities. The behavior owns the engineering semantics.

## No Cross-Device Access

Behaviors cannot access other devices:

```go
// This does not exist
func (b *Behavior) AccessOther(other *Device) {
    other.Memory().Write(...)  // Not possible
}
```

## Determinism

Behaviors are deterministic:

- Same memory → same results
- No randomness without seeded RNG
- No external system calls (except Raw Ingest publish)
- No time-of-day dependencies

## Raw Ingest Publishing

Publishing to MMA2 is optional. A behavior may:
- Write only to device memory (internal simulation)
- Write to device memory AND publish to MMA2 (operational data)

This allows devices to maintain private simulation state while exposing only relevant operational values.
