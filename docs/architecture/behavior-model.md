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
    publisher RawIngestPublisher
}

func (b *WeatherBehavior) Tick() {
    // Read from device memory
    irradiance := b.device.Memory().ReadFloat32("input_registers", irradianceAddr)
    temperature := b.device.Memory().ReadFloat32("input_registers", temperatureAddr)
    
    // Compute (internal simulation logic)
    // ...
    
    // Write to device memory
    b.device.Memory().WriteFloat32("input_registers", computedAddr, computedValue)
    
    // Publish to MMA2 via Raw Ingest
    b.publisher.Publish("weather/irradiance", irradiance, QualityGood)
    b.publisher.Publish("weather/temperature", temperature, QualityGood)
}
```

## Example: PV Model Behavior

```go
type PVModelBehavior struct {
    device    *Device
    publisher RawIngestPublisher
}

func (b *PVModelBehavior) Tick() {
    irradiance := b.device.Memory().ReadFloat32("input_registers", irradianceAddr)
    temperature := b.device.Memory().ReadFloat32("input_registers", temperatureAddr)
    
    dcPower := irradiance * b.scaleFactor * (1 - 0.004*(temperature-25))
    
    // Write to device memory
    b.device.Memory().WriteFloat32("input_registers", powerAddr, dcPower)
    
    // Publish to MMA2
    b.publisher.Publish("pv/dc_power", dcPower, QualityGood)
}
```

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
