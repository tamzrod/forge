# Behavior Model

## Philosophy

**Behaviors are logic that reads and writes device memory, and observes simulation models.**

A device owns its behaviors. Behaviors may optionally publish to MMA2 via Raw Ingest.

## Behavior Contract

A behavior:

1. Observes simulation models (Sun, Grid, Wind, etc.)
2. Reads from device memory
3. Computes new values
4. Writes to device memory
5. Optionally publishes to MMA2 via Raw Ingest
6. Never accesses other devices
7. Never calls protocols

## Data Flow

```
┌─────────────────┐
│ Simulation Model │ (Grid, Sun, Wind)
└────────┬────────┘
         │ observes
         ▼
┌─────────────────┐
│    Behavior     │
└────────┬────────┘
         │ reads/writes
         ▼
┌─────────────────┐
│  Device Memory  │
└────────┬────────┘
         │ publishes
         ▼
┌─────────────────┐
│   Raw Ingest    │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│      MMA2       │
└─────────────────┘
```

## Device Owns Behaviors

```go
type Device struct {
    behaviors []Behavior
    modelGetter func(ModelID) Model  // Access to models
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

// Device behavior observes models
type DeviceBehavior struct {
    device     *Device
    rawIngest  RawIngestClient
    unitID     uint16
}

// Observe models through device
func (b *DeviceBehavior) observeModel(id models.ModelID) models.Model {
    return b.device.Model(id)
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
    // Observe simulation models
    sun := b.device.Model("weather-sun")
    irradiance := sun.Irradiance()
    
    weather := b.device.Model("ambient-weather")
    temperature := weather.Temperature()
    
    wind := b.device.Model("weather-wind")
    windSpeed := wind.Speed()
    
    // Add measurement noise (simulated sensor)
    measured := b.addNoise(irradiance, temperature, windSpeed)
    
    // Write to device memory
    b.device.Memory().WriteFloat32("sensors", irradianceAddr, measured.irradiance)
    b.device.Memory().WriteFloat32("sensors", temperatureAddr, measured.temperature)
    b.device.Memory().WriteFloat32("sensors", windAddr, measured.windSpeed)
    
    // Device encodes: float32 → uint16 with scaling
    b.rawIngest.WriteInputRegisters(b.unitID, 0, encodeFloat32(measured.irradiance, 0, 2000))
    b.rawIngest.WriteInputRegisters(b.unitID, 2, encodeFloat32(measured.temperature, -50, 50))
}
```

## Example: PV Inverter Behavior

```go
type PVInverterBehavior struct {
    device    *Device
    rawIngest RawIngestClient
    unitID    uint16
}

func (b *PVInverterBehavior) Tick() {
    // Observe sun model
    sun := b.device.Model("solar-sun")
    irradiance := sun.Irradiance()
    temperature := sun.Elevation() // Affects efficiency
    
    // Observe ambient temperature from weather model
    weather := b.device.Model("ambient-weather")
    ambientTemp := weather.Temperature()
    
    // Compute DC power from sun
    dcPower := b.calculatePower(irradiance, ambientTemp)
    
    // Inject power into grid model
    grid := b.device.Model("main-grid")
    grid.InjectActivePower(dcPower)
    
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

## Model Observation Pattern

Behaviors observe models, not other devices:

```go
// Correct: Observe model
func (b *PVBehavior) Tick() {
    sun := b.device.Model("solar-sun")
    irradiance := sun.Irradiance()
}

// Incorrect: Access other device
func (b *BadBehavior) Tick() {
    // This does not exist
    meter := b.runtime.Device("meter-001")
    meter.Memory().Read(...)  // Not possible
}
```

## Determinism

Behaviors are deterministic:

- Same model state → same device behavior
- No randomness without seeded RNG
- No external system calls (except Raw Ingest publish)
- No time-of-day dependencies (use simulation clock)

## Raw Ingest Publishing

Publishing to MMA2 is optional. A behavior may:
- Write only to device memory (internal simulation)
- Write to device memory AND publish to MMA2 (operational data)

This allows devices to maintain private simulation state while exposing only relevant operational values.

## Model Injection Pattern

Some devices can inject back into models:

```go
// PV Inverter injects power into grid
func (b *PVInverterBehavior) Tick() {
    // Observe sun model
    sun := b.device.Model("solar-sun")
    power := b.calculatePower(sun.Irradiance())
    
    // Inject power into grid model
    grid := b.device.Model("main-grid")
    grid.InjectActivePower(power)  // Bidirectional
    
    // Publish to MMA2
    b.rawIngest.WriteInputRegisters(b.unitID, 0, encode(power))
}
```
