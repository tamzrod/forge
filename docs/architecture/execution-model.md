# Execution Model

## Purpose

The execution model describes how a simulation runs within the Virtual Industrial Laboratory.

Deterministic execution enables reproducible testing and training scenarios—the same inputs always produce the same outputs.

## Execution Flow

```
Startup: Load Plugins → Create Models → Create Devices → Connect Raw Ingest
       ↓
Simulation Loop: Tick Models → Tick Devices → Publish to MMA2 → Advance Clock
       ↓
Shutdown: Stop Devices → Disconnect → Unload Plugins
```

## Startup

```go
runtime, _ := forge.NewRuntime(Config{TickInterval: 250 * time.Millisecond})
runtime.LoadPlugins("./plugins/energy")

// Create simulation models (physical world)
runtime.CreateGridModel("main-grid")
runtime.CreateSunModel("solar-sun")
runtime.CreateWindModel("wind-farm")
runtime.CreateWeatherModel("ambient-weather")

// Connect to MMA2 via Raw Ingest
runtime.ConnectRawIngest(mma2Endpoint)

// Create devices
runtime.CreateDevice("weather-001", "weather_station", memRegions)
runtime.CreateDevice("pv-001", "pv_inverter", memRegions)
```

## Simulation Loop

```go
func (r *Runtime) Run(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            return r.Shutdown()
        case <-ticker.C:
            r.tick()
        }
    }
}

func (r *Runtime) tick() {
    // 1. Models evolve first (physics)
    for _, model := range r.models {
        model.Tick()
    }
    
    // 2. Devices observe models and update memory
    for _, device := range r.devices {
        device.Tick()
    }
    
    // 3. Raw Ingest publishes happen during device tick
    r.clock.Advance(r.tickInterval)
}
```

## Model Tick

Models evolve before devices. This ensures devices see consistent, updated state:

```go
func (s *SunModel) Tick() {
    // Calculate sun position from simulation time
    position := s.clock.SunPosition()
    
    // Calculate irradiance based on position
    irradiance := s.calculateIrradiance(position)
    
    // Update model state (private RAM)
    s.irradiance = irradiance
    s.elevation = position.elevation
}
```

## Device Tick

Devices observe models and update memory:

```go
func (d *Device) Tick() {
    for _, behavior := range d.behaviors {
        behavior.Tick()  // May observe models and publish
    }
}

// Example: PV Inverter observes sun model
func (b *PVInverterBehavior) Tick() {
    // Observe sun model
    sun := b.device.Model("solar-sun")
    irradiance := sun.Irradiance()
    
    // Compute output based on sun model
    power := b.calculatePower(irradiance)
    
    // Inject power into grid model
    grid := b.device.Model("main-grid")
    grid.InjectActivePower(power)
    
    // Write to device memory
    b.device.Memory().WriteFloat32("output", addr, power)
    
    // Publish to MMA2
    b.publisher.Publish("pv/power", power, QualityGood)
}
```

## Raw Ingest Publishing

During tick, behaviors publish to MMA2:

```go
func (b *WeatherBehavior) Tick() {
    // Weather behavior observes models
    sun := b.device.Model("weather-sun")
    irradiance := sun.Irradiance()
    
    weather := b.device.Model("ambient-weather")
    temperature := weather.Temperature()
    
    // Compute weather values with sensor noise
    measured := b.measureWeather(irradiance, temperature)
    
    // Write to device memory
    b.device.Memory().WriteFloat32("input_registers", addr, measured)
    
    // Publish to MMA2
    b.publisher.Publish("weather/irradiance", measured, QualityGood)
}
```

## Shutdown

```go
func (r *Runtime) Shutdown() error {
    r.Stop()
    r.rawIngest.Disconnect()
    r.devices.Clear()
    r.models.Clear()
    r.plugins.Shutdown()
    return nil
}
```

## Summary

| Component | Purpose | Order |
|------------|---------|-------|
| **Runtime** | Hosts models and devices | - |
| **Simulation Models** | Represent physics | 1st (tick first) |
| **Devices** | Observe models, own memory | 2nd |
| **Behaviors** | Read models, write memory, publish | within device tick |
| **MMA2** | Owns operational memory, exposes protocols | after publish |

The execution order ensures:

```
Physics (Models) → Observation (Devices) → Telemetry (MMA2) → Control (Atlas-PPC)
```
