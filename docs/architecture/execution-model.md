# Execution Model

## Overview

The execution model describes how a simulation runs.

## Execution Flow

```
Startup: Load Plugins → Create Infrastructure → Create Devices → Connect Raw Ingest
       ↓
Simulation Loop: Tick Infrastructure → Tick Devices → Publish to MMA2 → Advance Clock
       ↓
Shutdown: Stop Devices → Disconnect → Unload Plugins
```

## Startup

```go
runtime, _ := forge.NewRuntime(Config{TickInterval: 250 * time.Millisecond})
runtime.LoadPlugins("./plugins/energy")

// Create shared infrastructure
runtime.CreateInfrastructure(Sun{}, Grid{})

// Connect to MMA2 via Raw Ingest
runtime.ConnectRawIngest(mma2Endpoint)

// Create devices
runtime.CreateDevices([]DeviceConfig{
    {ID: "weather-001", Type: "weather_station"},
})

runtime.CreateDevices([]DeviceConfig{
    {ID: "pv-001", Type: "pv_inverter"},
})
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
    // 1. Infrastructure evolves first
    for _, infra := range r.infrastructure {
        infra.Tick()
    }
    
    // 2. Devices observe infrastructure and update memory
    for _, device := range r.devices {
        device.Tick()
    }
    
    // 3. Raw Ingest publishes happen during device tick
    r.clock.Advance(r.tickInterval)
}
```

## Infrastructure Tick

Infrastructure evolves before devices:

```go
func (s *Sun) Tick() {
    // Calculate sun position from simulation time
    position := s.clock.SunPosition()
    
    // Calculate irradiance based on position
    irradiance := s.calculateIrradiance(position)
    
    // Update shared infrastructure state
    s.SetIrradiance(irradiance)
    s.SetPosition(position)
}
```

## Device Tick

Devices observe infrastructure and update memory:

```go
func (d *Device) Tick() {
    for _, behavior := range d.behaviors {
        behavior.Tick()  // May call publisher.Publish()
    }
}

// Example: PV Inverter observes sun
func (b *PVInverterBehavior) Tick() {
    // Observe infrastructure (not another device)
    irradiance := b.infrastructure.Sun().Irradiance()
    
    // Compute output based on infrastructure
    power := b.calculatePower(irradiance)
    
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
    // Weather behavior observes infrastructure
    irradiance := b.infrastructure.Sun().Irradiance()
    temperature := b.infrastructure.Ambient().Temperature()
    
    // Compute weather values
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
    r.infrastructure.Clear()
    r.plugins.Shutdown()
    return nil
}
```

## Summary

- Runtime: hosts, schedules, advances time, manages Raw Ingest
- Infrastructure: shared world (Grid, Sun, Wind), evolves first
- Devices: observe infrastructure, own memory, execute behaviors
- Behaviors: read infrastructure, write device memory, publish to MMA2
- MMA2: owns operational memory, exposes protocols
