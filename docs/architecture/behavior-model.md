# Behavior Model

## Philosophy

**Behaviors are logic that reads and writes device memory.**

A device owns its behaviors. Behaviors never access other devices. Behaviors never call protocols.

## Behavior Contract

A behavior:

1. Reads from device memory
2. Computes new values
3. Writes to device memory
4. Never accesses other devices
5. Never calls protocols

## Data Flow

```
Behavior → Device Memory → Behavior
           (reads)           (writes)
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

## Example

```go
type PVModelBehavior struct {
    device *Device
}

func (b *PVModelBehavior) Tick() {
    irradiance := b.device.Memory().ReadFloat32("input_registers", 0)
    temperature := b.device.Memory().ReadFloat32("input_registers", 4)
    
    dcPower := irradiance * b.scaleFactor * (1 - 0.004*(temperature-25))
    
    b.device.Memory().WriteFloat32("input_registers", 8, dcPower)
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
- No external system calls
- No time-of-day dependencies
