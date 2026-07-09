# Simulation Models

## Purpose

Simulation Models represent the **physical world** in which virtual devices operate. They are first-class architectural components that model physics, not equipment.

This supports the Virtual Industrial Laboratory vision: providing a believable industrial environment where software behaves as it would with real hardware.

## Why Simulation Models Exist

Virtual devices require a simulated world in which to operate. A **Revenue Meter** does not generate voltage—it measures voltage from the Grid. A **Weather Station** does not generate sunlight—it measures the Weather Model. A **PV Inverter** does not create irradiance—it converts power from the Sun Model.

Therefore, these physical concepts should not be modeled as devices.

## Modeling Philosophy

**Believe before sophisticated.** Models should be credible before they become complex.

Simple deterministic models are preferred over highly accurate but complex models unless additional fidelity clearly benefits industrial software development.

For software testing, a Grid Model that produces believable voltage fluctuations is more valuable than an electromagnetic transient simulator.

**Examples of Physical Phenomena:**

| Domain | Simulation Models |
|--------|------------------|
| **Energy** | Grid, Sun, Wind, Weather |
| **Water** | Reservoir, Hydraulic Network, River |
| **Manufacturing** | Factory Power, Compressed Air, Conveyor Physics |

## Simulation Models vs Devices

| Property | Simulation Model | Virtual Device |
|----------|-----------------|----------------|
| **Represents** | Physics (Grid, Sun, Wind) | Equipment (Meter, Inverter, Relay) |
| **Identity** | None (type-based) | Unique device ID |
| **Memory** | Private RAM state | Private memory image |
| **Protocols** | None | Modbus, DNP3, REST |
| **External Clients** | Never connected | Atlas-PPC, SCADA, HMIs |
| **Publishes** | Never | Yes, to MMA2 |
| **Observers** | Many devices | None |

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         Simulation Runtime                                 │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                      Simulation Models                                    │
│                                                                         │
│  Grid │ Sun │ Wind │ Weather │ Reservoir │ Hydraulic Network           │
│                                                                         │
│  External physical world - private RAM                                  │
│  No protocols, no external clients, no MMA2 publishing                 │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                        Virtual Firmware                                   │
│                                                                         │
│  Weather Station │ PV Inverter │ Revenue Meter │ Relay                │
│                                                                         │
│  Samples models, owns Device Memory, exposes via interfaces             │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                    Communication Interfaces                                │
│                                                                         │
│  Raw Ingest │ Modbus │ DNP3 │ IEC 61850 │ MQTT │ REST                │
│                                                                         │
│  Serialize Device Memory only                                            │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                              MMA2                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Data Flow

```
Simulation Models (external physical world)
        ↓
Virtual Firmware samples models, updates Device Memory
        ↓
Communication Interfaces serialize Device Memory
        ↓
MMA2 / SCADA / Atlas-PPC
```

## Model API

Simulation Models expose behavior through **methods**, not public fields.

### Grid Model

```go
// GridModel represents an electrical grid
type GridModel struct {
    // Private state - use methods to access
}

// Voltage returns grid voltage in volts
func (g *GridModel) Voltage() float32

// Frequency returns grid frequency in Hz
func (g *GridModel) Frequency() float32

// InjectActivePower records power injection from a device
func (g *GridModel) InjectActivePower(mw float32)

// InjectReactivePower records reactive power injection
func (g *GridModel) InjectReactivePower(mvar float32)
```

### Sun Model

```go
// SunModel represents solar conditions
type SunModel struct{}

// Irradiance returns global horizontal irradiance in W/m²
func (s *SunModel) Irradiance() float32

// DirectNormalIrradiance returns DNI in W/m²
func (s *SunModel) DirectNormalIrradiance() float32

// Elevation returns sun elevation angle in degrees
func (s *SunModel) Elevation() float32
```

### Weather Model

```go
// WeatherModel represents ambient conditions
type WeatherModel struct{}

// Temperature returns ambient temperature in °C
func (w *WeatherModel) Temperature() float32

// Humidity returns relative humidity in %
func (w *WeatherModel) Humidity() float32

// Pressure returns atmospheric pressure in hPa
func (w *WeatherModel) Pressure() float32

// CloudCover returns cloud cover fraction (0-1)
func (w *WeatherModel) CloudCover() float32
```

## Storage

Simulation Models store their state **directly in RAM**.

| Aspect | Simulation Models |
|--------|------------------|
| **Storage** | RAM only |
| **MMA2** | Never published |
| **Modbus** | Never exposed |
| **Raw Ingest** | Not used |
| **Protocol Abstraction** | None |
| **Memory Appliance** | Not applicable |

Simulation Models are **internal implementation objects**. Their state is private to the simulation.

### Example: Grid Model State

```go
type GridModel struct {
    // Physical state stored in RAM
    voltage         float32 // V
    frequency       float32 // Hz
    theveninZReal   float32 // Ω
    theveninZImag   float32 // Ω
    shortCircuitMVA float32 // MVA
    reactiveSens    float32 // PU MVAr per PU voltage

    // Power injection from devices (reset each tick)
    activePowerInjection   float32 // MW
    reactivePowerInjection float32 // MVAr
}
```

## Interaction

Virtual Firmware **samples** Simulation Models through a read-only context.

### Revenue Meter observes Grid

```go
type RevenueMeter struct {
    *devices.BaseDevice
    ctx *devices.Context
    memory *devices.DeviceMemory
}

func (r *RevenueMeter) Tick() {
    // Sample grid model
    grid := r.ctx.ReadGrid()

    // Update device memory with observations
    r.memory.Set("voltage", grid.Voltage)
    r.memory.Set("frequency", grid.Frequency)
}
```

### PV Inverter observes Sun and Grid

```go
type PVInverter struct {
    *devices.BaseDevice
    ctx *devices.Context
    memory *devices.DeviceMemory
}

func (p *PVInverter) Tick() {
    // Sample sun model
    sun := p.ctx.ReadSun()
    irradiance := sun.Irradiance

    // Sample grid model
    grid := p.ctx.ReadGrid()

    // Compute DC power from irradiance
    dcPower := p.calculatePower(irradiance)

    // Inject power into grid (grid is writeable via context)
    p.ctx.Grid().InjectActivePower(dcPower)

    // Update device memory
    p.memory.Set("dc_power", dcPower)
}
```

### Weather Station observes multiple models

```go
type WeatherStation struct {
    *devices.BaseDevice
    ctx *devices.Context
    memory *devices.DeviceMemory
}

func (w *WeatherStation) Tick() {
    // Sample sun model
    sun := w.ctx.ReadSun()

    // Sample weather model
    weather := w.ctx.ReadWeather()

    // Update device memory with observations
    w.memory.Set("temperature", weather.Temperature)
    w.memory.Set("humidity", weather.Humidity)
    w.memory.Set("irradiance", sun.Irradiance)
}
```

Firmware owns Device Memory. Communication Interfaces serialize memory without accessing models.

## Publishing

Only **Virtual Firmware** is responsible for publishing operational measurements.

```
Simulation Models (external world)
        ↓
Virtual Firmware samples and updates Device Memory
        ↓
Communication Interface serializes Device Memory
        ↓
MMA2
```

**Simulation Models never publish.** Only firmware publishes through its communication interfaces.

## Execution Order

Simulation Models execute **before** Virtual Firmware:

```
Tick 1:
    ┌─────────────────────────────────────────────────────────┐
    │ 1. Sun Model evolves (position, irradiance)             │
    │ 2. Wind Model evolves (speed, direction)              │
    │ 3. Weather Model evolves (temperature, pressure)      │
    │ 4. Grid Model evolves (voltage, frequency)            │
    └─────────────────────────────────────────────────────────┘
                        ↓
    ┌─────────────────────────────────────────────────────────┐
    │ 5. Weather Station firmware samples models              │
    │ 6. PV Inverter firmware samples sun, injects power     │
    │ 7. Revenue Meter firmware samples grid                 │
    │ 8. Firmware pushes to communication interfaces        │
    └─────────────────────────────────────────────────────────┘
```

This ensures firmware always sees consistent, updated model state.

## Architectural Principles

| Layer | Responsibility |
|-------|----------------|
| **Simulation Models** | Represent physics (external world) |
| **Virtual Firmware** | Owns Device Memory, samples models |
| **Communication Interfaces** | Serialize Device Memory only |
| **MMA2** | Receive operational telemetry |
| **Atlas-PPC** | Control logic |

Each layer owns **one responsibility**.

## Model Types

### Energy Domain

| Model | Properties |
|-------|------------|
| Grid | voltage, frequency, Thevenin impedance, short circuit level |
| Sun | irradiance, position, azimuth, elevation |
| Wind | speed, direction, gusts, turbulence |
| Weather | temperature, humidity, pressure, cloud cover |

### Water Domain

| Model | Properties |
|-------|------------|
| Reservoir | level, flow in/out, temperature |
| Hydraulic Network | pressure, flow rates |
| River | flow, quality, temperature |

### Manufacturing Domain

| Model | Properties |
|-------|------------|
| Factory Power | voltage, frequency, power factor |
| Compressed Air | pressure, flow rate |
| Conveyor Physics | speed, load, position |

## Summary

| Concept | Purpose |
|---------|---------|
| **Simulation Models** | Represent the physical world (physics) |
| **Virtual Firmware** | Samples models, owns Device Memory |
| **Communication Interfaces** | Serialize Device Memory only |
| **MMA2** | Receive operational telemetry |

The architecture mirrors a real industrial system:

```
Physics → Firmware → Device Memory → Communication Interface → External Systems
```

without conflating these responsibilities.

---

*Last Updated: 2026-07-09*
