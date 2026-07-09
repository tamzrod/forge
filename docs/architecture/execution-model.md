# Execution Model

## Overview

The execution model describes how a simulation runs.

## Execution Flow

```
Startup: Load Plugins → Create Devices → Connect Raw Ingest
       ↓
Simulation Loop: Tick Devices → Publish to MMA2 → Advance Clock
       ↓
Shutdown: Stop Devices → Disconnect → Unload Plugins
```

## Startup

```go
runtime, _ := forge.NewRuntime(Config{TickInterval: 250 * time.Millisecond})
runtime.LoadPlugins("./plugins/energy")

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
    for _, device := range r.devices {
        device.Tick()
    }
    // Raw Ingest publishes happen during device tick
    r.clock.Advance(r.tickInterval)
}
```

## Device Tick

```go
func (d *Device) Tick() {
    for _, behavior := range d.behaviors {
        behavior.Tick()  // May call publisher.Publish()
    }
}
```

## Raw Ingest Publishing

During tick, behaviors may publish to MMA2:

```go
func (b *WeatherBehavior) Tick() {
    // Compute weather values
    irradiance := b.computeIrradiance()
    
    // Write to device memory
    b.device.Memory().WriteFloat32("input_registers", addr, irradiance)
    
    // Publish to MMA2
    b.publisher.Publish("weather/irradiance", irradiance, QualityGood)
}
```

## Shutdown

```go
func (r *Runtime) Shutdown() error {
    r.Stop()
    r.rawIngest.Disconnect()
    r.devices.Clear()
    r.plugins.Shutdown()
    return nil
}
```

## Summary

- Runtime: hosts, schedules, advances time, manages Raw Ingest
- Devices: own memory, execute behaviors
- Behaviors: read and write memory, publish to MMA2
- MMA2: owns operational memory, exposes protocols
