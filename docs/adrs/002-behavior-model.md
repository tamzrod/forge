# ADR-002: Behavior Model Design

**ADR ID:** ADR-002  
**Title:** Behavior Model for Device Logic Execution  
**Date:** 2026-07-11  
**Status:** Accepted  
**Deciders:** Engineering Team  
**Repository:** https://github.com/tamzrod/forge

---

## Context

Devices in the Industrial Simulation Runtime need to execute logic on each simulation tick. The behavior model provides:

1. A simple interface for device logic
2. Lifecycle management (attach/detach)
3. Access to device memory and simulation context
4. Execution on fixed tick intervals

---

## Decision

We adopt a **Function-Based Behavior Model** with an interface:

```go
type Behavior interface {
    ID() string
    Attach(device *Device)
    Detach(device *Device)
    Tick(device *Device)
}
```

### Design Principles

1. **Simple Interface**: Only 4 methods required
2. **Lifecycle Hooks**: Attach/Detach for setup/teardown
3. **Device Access**: Behaviors can read/write device memory
4. **Context Access**: Behaviors can read simulation models

---

## Implementation Details

### Behavior Interface

```go
// Behavior defines device behavior execution.
type Behavior interface {
    ID() string
    Attach(device *Device)
    Detach(device *Device)
    Tick(device *Device)
}
```

### Lifecycle

1. **Attach**: Called when behavior is added to device
   - Initialize resources
   - Register memory regions
   - Set up callbacks

2. **Detach**: Called when behavior is removed
   - Clean up resources
   - Release memory regions

3. **Tick**: Called on each simulation tick
   - Execute device logic
   - Read/write device memory
   - Sample simulation models

### Example: Weather Station Behavior

```go
type weatherBehavior struct {
    attachFn func(*Device)
    detachFn func(*Device)
}

func (b *weatherBehavior) ID() string { return "weather_behavior" }

func (b *weatherBehavior) Attach(d *Device) {
    d.AddMemoryRegion("input_registers", 0, 9, ReadOnly)
    if b.attachFn != nil {
        b.attachFn(d)
    }
}

func (b *weatherBehavior) Detach(d *Device) {
    if b.detachFn != nil {
        b.detachFn(d)
    }
}

func (b *weatherBehavior) Tick(d *Device) {
    // Sample weather model
    weather := d.ctx.ReadWeather()
    
    // Update memory
    d.Memory().WriteFloat32("input_registers", 0, weather.Temperature)
    d.Memory().WriteFloat32("input_registers", 1, weather.Humidity)
    d.Memory().WriteFloat32("input_registers", 2, weather.Pressure)
}
```

---

## Consequences

### Positive
- Simple, intuitive interface
- Clear lifecycle management
- Flexible for different device types
- Easy to test with mock behaviors

### Negative
- Behaviors cannot spawn goroutines (single-threaded)
- No behavior composition built-in
- No priority/ordering control

### Risks
- Behaviors sharing state need external synchronization
- Long-running ticks may block simulation

---

## Alternatives Considered

### Alternative 1: Method-Based Behaviors
```go
type DeviceBehavior interface {
    OnAttach()
    OnDetach()
    OnTick()
}
```
**Rejected**: Less flexible, no device reference.

### Alternative 2: Event-Based Behaviors
```go
device.On("tick", func() { ... })
```
**Rejected**: Too complex for initial implementation.

---

## References

- [Behavior Model](docs/architecture/behavior-model.md)
- [Device Model](docs/architecture/device-model.md)

---

## Related ADRs

- ADR-001: Runtime Architecture
- ADR-003: Memory Model Design

---

## Testing

Behaviors should be tested with:
- Unit tests for lifecycle methods
- Integration tests with mock devices
- Concurrent access tests (if sharing state)

---

## Milestone Traceability

| Milestone | Status |
|-----------|--------|
| 1.2.1 Behavior Interface | ✅ Complete |
| 1.2.2 Attach/Detach Lifecycle | ✅ Complete |
| 1.2.3 Tick Execution | ✅ Complete |
| 1.2.4 Behavior Testing | ✅ Complete |
