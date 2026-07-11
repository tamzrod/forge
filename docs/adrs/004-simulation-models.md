# ADR-004: Simulation Models Design

**ADR ID:** ADR-004  
**Title:** Physical World Simulation Models  
**Date:** 2026-07-11  
**Status:** Accepted  
**Deciders:** Engineering Team  
**Repository:** https://github.com/tamzrod/forge

---

## Context

The Industrial Simulation Runtime needs to simulate the physical world:
- Electrical grid (voltage, frequency)
- Solar position and irradiance
- Weather conditions (temperature, humidity, pressure, wind)
- Water reservoir levels

These models must:
1. Advance deterministically with simulation time
2. Be observable by devices
3. Support fault injection
4. Maintain physical validity

---

## Decision

We adopt a **Model-Context Observation Pattern**:

```go
type ModelContext interface {
    ReadGrid() GridState
    ReadSun() SunState
    ReadWeather() WeatherState
    ReadWind() WindState
    ReadReservoir() ReservoirState
}

type DeviceContext interface {
    ModelContext
    Clock() clock.Clock
}
```

### Grid Model

Simulates electrical grid parameters:
- Voltage: 450-520V
- Frequency: 59.5-60.5 Hz
- Power injection/absorption affects values

```go
type GridModel struct {
    voltage   float32  // Volts
    frequency float32  // Hertz
}

func (g *GridModel) SetVoltage(v float32)
func (g *GridModel) SetFrequency(f float32)
func (g *GridModel) InjectPower(p float32)  // Positive: injection, Negative: absorption
```

### Sun Model

Simulates solar position and irradiance:
- Elevation: -90° to 90°
- Azimuth: 0° to 360°
- Irradiance: 0 to 1000 W/m²
- IsDaytime: boolean

```go
type SunModel struct {
    elevation  float64  // Degrees
    azimuth    float64  // Degrees [0, 360)
    irradiance float64  // W/m²
}

func (s *SunModel) Elevation() float64
func (s *SunModel) Azimuth() float64
func (s *SunModel) Irradiance() float64
func (s *SunModel) IsDaytime() bool
```

### Weather Model

Simulates atmospheric conditions:
- Temperature: -40°C to 60°C
- Humidity: 0% to 100%
- Pressure: 800 to 1200 hPa
- Cloud Cover: 0% to 100%
- Wind Speed: 0 to 50 m/s
- Rain Status: boolean

```go
type WeatherState struct {
    Temperature  float64  // °C
    Humidity     float64  // %
    Pressure     float64  // hPa
    CloudCover   float64  // %
    WindSpeed    float64  // m/s
    WindDirection float64 // Degrees
    RainStatus   bool
}
```

### Wind Model

Simulates wind conditions:
- Speed: 0 to 50 m/s
- Direction: 0° to 360°

```go
type WindModel struct {
    speed     float64
    direction float64
}
```

### Reservoir Model

Simulates water storage:
- Flow In: m³/s
- Flow Out: m³/s
- Level: computed from flows

```go
type ReservoirModel struct {
    flowIn    float64  // m³/s
    flowOut   float64  // m³/s
    capacity  float64  // m³
}
```

---

## Deterministic Time

All models advance based on simulation clock:

```go
func (m *SunModel) Tick() {
    elapsed := m.clock.Elapsed()
    // Compute position from elapsed time
}
```

### Clock Modes

| Mode | Description |
|------|-------------|
| Realtime | Advances at real-time speed |
| Manual | Advances only on explicit step |
| Accelerated | Advances faster than realtime |

---

## Consequences

### Positive
- Deterministic execution for reproducible tests
- Physical validity enforced by bounds
- Models are decoupled from devices
- Easy to mock for testing

### Negative
- Fixed physics (no configurable parameters)
- Simplified models (no complex interactions)
- Clock-dependent (not suitable for real-time IoT)

### Risks
- Models may drift from physical reality over long simulations
- Grid/weather coupling not implemented

---

## Alternatives Considered

### Alternative 1: Real-Time Sampling
Sample from real sensors instead of models.

**Rejected**: Loses determinism, requires hardware.

### Alternative 2: Physics Engine
Use external physics engine (e.g., EnergyPlus).

**Rejected**: Too complex for initial implementation.

---

## References

- [Simulation Models](docs/architecture/simulation-models.md)
- [Grid Model](docs/architecture/infrastructure-model.md)
- [Sun Model](docs/architecture/infrastructure-model.md)
- [Weather Model](docs/architecture/infrastructure-model.md)

---

## Related ADRs

- ADR-001: Runtime Architecture
- ADR-005: Simulation Clock Design

---

## Testing

Models are tested for:
- Bounds enforcement
- Deterministic behavior
- Clock advancement
- Fault injection support

---

## Milestone Traceability

| Milestone | Status |
|-----------|--------|
| 2.1 Grid Model | ✅ Complete |
| 2.2 Sun Model | ✅ Complete |
| 2.3 Weather Model | ✅ Complete |
| 2.4 Wind Model | ✅ Complete |
| 2.5 Reservoir Model | ✅ Complete |
| 2.6 Fault Injection | ⏳ Pending |
