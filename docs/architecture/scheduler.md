# Scheduler

## Role

The scheduler advances simulation time deterministically. It tells models and devices to tick.

## Tick Model

```
Simulation Tick
      │
      ▼
┌─────────────────────────────┐
│  1. Tell each model to tick │
│     (Models evolve physics)   │
└─────────────────────────────┘
      │
      ▼
┌─────────────────────────────┐
│  2. Tell each device to tick │
│     (Devices observe models)  │
└─────────────────────────────┘
      │
      ▼
┌─────────────────────────────┐
│  3. Advance simulation clock │
└─────────────────────────────┘
```

## Implementation

```go
func (s *Scheduler) tick() {
    // 1. Models evolve first (physics)
    for _, model := range s.models {
        model.Tick()
    }
    
    // 2. Devices observe models and update memory
    for _, device := range s.devices {
        device.Tick()
    }
    
    // 3. Advance the clock
    s.clock.Advance(s.tickInterval)
}
```

Models tick first to ensure devices see consistent, updated physical state.

## Determinism

Execution is deterministic:

1. Models tick in registration order
2. Devices tick in registration order
3. Behaviors tick in registration order
4. Same inputs → same outputs, every time

This enables reproducible testing and training scenarios.

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
