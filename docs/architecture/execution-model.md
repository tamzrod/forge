# Execution Model

## Overview

The execution model describes how a simulation runs.

## Execution Flow

```
Startup: Load Plugins → Create Devices → Expose Protocols
       ↓
Simulation Loop: Tick Devices → Advance Clock
       ↓
Shutdown: Stop Devices → Unload Plugins
```

## Startup

```go
runtime, _ := forge.NewRuntime(Config{TickInterval: 250 * time.Millisecond})
runtime.LoadPlugins("./plugins/energy")

runtime.CreateDevices([]DeviceConfig{
    {ID: "meter-001", Type: "revenue_meter"},
})

runtime.Device("meter-001").ExposeProtocol("modbus", NewModbusAdapter())
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
    r.clock.Advance(r.tickInterval)
}
```

## Device Tick

```go
func (d *Device) Tick() {
    for _, behavior := range d.behaviors {
        behavior.Tick()
    }
}
```

## Shutdown

```go
func (r *Runtime) Shutdown() error {
    for _, device := range r.devices {
        for _, protocol := range device.Protocols() {
            protocol.Stop()
        }
    }
    r.devices.Clear()
    r.plugins.Shutdown()
    return nil
}
```

## Summary

- Runtime: hosts, schedules, advances time
- Devices: own memory, execute behaviors, expose protocols
- Behaviors: read and write memory
- Protocols: expose memory externally
