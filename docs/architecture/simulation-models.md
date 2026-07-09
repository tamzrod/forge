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
│  Shared physical world                                                 │
│  No protocols, no external clients, no MMA2 publishing                 │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                          Virtual Devices                                  │
│                                                                         │
│  Revenue Meter │ Weather Station │ PV Inverter │ Relay                │
│                                                                         │
│  Observe models, publish to MMA2                                       │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                      Operational Publisher                                │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                              MMA2                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## Data Flow

```
Simulation Truth (Models)
        ↓
Device Observation (Behaviors read models)
        ↓
Operational Telemetry (Devices publish to MMA2)
        ↓
Control Applications (Atlas-PPC, SCADA)
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

Virtual Devices **observe** Simulation Models through their behaviors.

### Revenue Meter observes Grid

```go
type RevenueMeterBehavior struct {
    device *Device
}

func (b *RevenueMeterBehavior) Tick() {
    // Observe grid model
    grid := b.device.Model("main-grid")
    if grid == nil {
        return
    }
    gridModel := grid.(*models.GridModel)

    // Read grid values
    voltage := gridModel.Voltage()
    frequency := gridModel.Frequency()

    // Write to device memory
    b.device.Memory().WriteFloat32("input_registers", voltageAddr, voltage)
    b.device.Memory().WriteFloat32("input_registers", freqAddr, frequency)

    // Publish to MMA2
    b.rawIngest.WriteInputRegisters(unitID, voltageReg, encode(voltage))
}
```

### PV Inverter observes Sun and Grid

```go
type PVInverterBehavior struct {
    device *Device
}

func (b *PVInverterBehavior) Tick() {
    // Observe sun model
    sun := b.device.Model("solar-sun")
    irradiance := sun.(*models.SunModel).Irradiance()

    // Observe grid model
    grid := b.device.Model("main-grid")
    gridModel := grid.(*models.GridModel)

    // Compute DC power from irradiance
    dcPower := b.calculatePower(irradiance)

    // Inject power into grid
    gridModel.InjectActivePower(dcPower)

    // Write to device memory
    b.device.Memory().WriteFloat32("output", powerAddr, dcPower)

    // Publish to MMA2
    b.rawIngest.WriteInputRegisters(unitID, powerReg, encode(dcPower))
}
```

### Weather Station observes multiple models

```go
type WeatherStationBehavior struct {
    device *Device
}

func (b *WeatherStationBehavior) Tick() {
    // Observe sun model
    sun := b.device.Model("weather-sun")
    irradiance := sun.(*models.SunModel).Irradiance()

    // Observe wind model
    wind := b.device.Model("weather-wind")
    windSpeed := wind.(*models.WindModel).Speed()

    // Observe weather model
    weather := b.device.Model("ambient-weather")
    temperature := weather.(*models.WeatherModel).Temperature()

    // Add measurement noise (simulated sensor characteristics)
    measured := b.addNoise(irradiance, windSpeed, temperature)

    // Write to device memory
    b.device.Memory().WriteFloat32("sensors", irradianceAddr, measured.irradiance)
    b.device.Memory().WriteFloat32("sensors", windAddr, measured.windSpeed)
    b.device.Memory().WriteFloat32("sensors", tempAddr, measured.temperature)

    // Publish to MMA2
    b.rawIngest.WriteInputRegisters(unitID, 0, encode(measured.irradiance))
    b.rawIngest.WriteInputRegisters(unitID, 2, encode(measured.windSpeed))
    b.rawIngest.WriteInputRegisters(unitID, 4, encode(measured.temperature))
}
```

## Publishing

Only **Virtual Devices** are responsible for publishing operational measurements.

```
Grid Model
    ↓
Revenue Meter
    ↓
Measurement
    ↓
Raw Ingest
    ↓
MMA2
```

**The Grid Model never publishes directly.** Only devices publish operational information.

## Execution Order

Simulation Models execute **before** Virtual Devices:

```
Tick 1:
    ┌─────────────────────────────────────────────────────────┐
    │ 1. Sun Model evolves (position, irradiance)             │
    │ 2. Wind Model evolves (speed, direction)                │
    │ 3. Weather Model evolves (temperature, pressure)        │
    │ 4. Grid Model evolves (voltage, frequency)              │
    └─────────────────────────────────────────────────────────┘
                        ↓
    ┌─────────────────────────────────────────────────────────┐
    │ 5. Weather Station observes models, updates memory      │
    │ 6. PV Inverter observes sun, injects power             │
    │ 7. Revenue Meter observes grid, records measurements    │
    │ 8. Devices publish to MMA2                              │
    └─────────────────────────────────────────────────────────┘
```

This ensures devices always see consistent, updated model state.

## Architectural Principles

| Layer | Responsibility |
|-------|----------------|
| **Simulation Models** | Represent physics |
| **Virtual Devices** | Represent equipment |
| **MMA2** | Represent operational telemetry |
| **Atlas-PPC** | Represent control logic |

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
| **Virtual Devices** | Represent equipment that observes physics |
| **Operational Publishing** | Devices expose telemetry to MMA2 |
| **MMA2** | Owns operational memory and protocols |

The architecture mirrors a real industrial system:

```
Physics → Devices → Operational Telemetry → Control Applications
```

without conflating these responsibilities.

---

*Last Updated: 2026-07-09*
