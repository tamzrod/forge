# Scheduler

## Role

The scheduler advances simulation time. It tells devices to tick.

## Tick Model

```
Simulation Tick
      │
      ▼
┌─────────────────────────────┐
│    Tell each device to tick   │
└─────────────────────────────┘
      │
      ▼
┌─────────────────────────────┐
│    Advance simulation clock   │
└─────────────────────────────┘
```

## Implementation

```go
func (s *Scheduler) tick() {
    for _, device := range s.devices {
        device.Tick()
    }
    s.clock.Advance(s.tickInterval)
}
```

The scheduler tells devices to tick. Devices execute their own behaviors.

## Determinism

Execution is deterministic:

1. Devices tick in registration order
2. Behaviors tick in registration order
3. Same inputs → same outputs

## Configuration

```yaml
runtime:
  tick_interval: 250ms
  time_multiplier: 1.0
```

## Pause and Resume

```go
scheduler.Pause()
scheduler.Resume()
```
