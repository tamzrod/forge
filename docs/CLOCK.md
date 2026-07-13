# Simulation Clock Architecture

## Why Simulation Clock Exists

The Simulation Clock is the single authoritative source of time for the entire simulation. It decouples simulated time from wall-clock time, enabling:

1. **Deterministic Simulations** - Same inputs produce same outputs regardless of execution speed
2. **Time Scaling** - Run simulations faster or slower than real-time
3. **Reproducibility** - Start simulations at any point in time
4. **Testing** - Simulate years of operation in seconds
5. **Architecture** - No module calls wall clock directly

## Wall Clock vs Simulation Clock

| Aspect | Wall Clock | Simulation Clock |
|--------|------------|-----------------|
| Source | `time.Now()` | `simulation.Clock.Now()` |
| Behavior | Real-world progression | Configurable |
| Speed | Fixed (1x) | Variable (0.1x to 100x) |
| Start | Always now | Configurable |
| Use | Only for I/O timing | All simulation logic |

## Core Interface

```go
type Clock interface {
    Now() time.Time      // Current simulation datetime
    Elapsed() duration   // Time since simulation start
    Tick() uint64        // Current tick number
    StartTime() time.Time
    Speed() float64      // Speed multiplier
    SetSpeed(float64)
    Mode() Mode          // Simulation mode
    IsRunning() bool
    IsPaused() bool
    Start(time.Time) error
    Pause()
    Resume()
    Stop()
    Advance(duration)    // For manual mode
    Update()             // Advance based on wall time
}
```

## Simulation Modes

### Realtime
Simulation follows wall clock, optionally scaled.
```
1 second wall → 1 second sim (1x)
1 second wall → 2 seconds sim (2x)
```

### Simulated  
Simulation advances independently of wall clock.
```
100 ticks × 100ms = 10 seconds sim time
```

### Manual (Placeholder)
Simulation advances only when `Advance()` is called.
```go
clock.Advance(1 * time.Second)  // Step forward 1 second
```

### Replay (Placeholder)
Simulation follows a recorded timeline.

## Time Scaling

Supported speeds:
```
0.1x  - 10 seconds wall → 1 second sim
0.5x  - 2 seconds wall → 1 second sim  
1.0x  - 1 second wall → 1 second sim
2.0x  - 1 second wall → 2 seconds sim
5.0x  - 1 second wall → 5 seconds sim
10x   - 1 second wall → 10 seconds sim
20x   - 1 second wall → 20 seconds sim
50x   - 1 second wall → 50 seconds sim
100x  - 1 second wall → 100 seconds sim
```

Speed changes take effect immediately without restarting.

## Simulation DateTime

The simulation starts at a configurable DateTime:

```go
startTime := time.Date(2026, 6, 21, 12, 0, 0, 0, time.UTC)
clock.Start(startTime)

// Now() returns: 2026-06-21T12:00:00Z
```

This DateTime drives entity behavior:
- Sun position based on time of day
- Weather patterns tied to date
- Seasonal variations
- Day/night cycles

## Architecture Integration

```
┌─────────────────────────────────────────┐
│           Simulation Clock               │
│  ┌─────────────────────────────────┐   │
│  │ Mode: Realtime/Simulated/Manual │   │
│  │ Speed: 0.1x to 100x             │   │
│  │ Start: 2026-06-21 12:00:00 UTC  │   │
│  └─────────────────────────────────┘   │
└──────────────────┬────────────────────┘
                   │
     ┌─────────────┼─────────────┐
     │             │             │
     ▼             ▼             ▼
┌─────────┐ ┌─────────┐ ┌─────────┐
│Scheduler│ │  World  │ │Scenarios│
└────┬────┘ └────┬────┘ └────┬────┘
     │            │            │
     ▼            ▼            ▼
┌─────────┐ ┌─────────┐ ┌─────────┐
│ Devices │ │Entities │ │ Events  │
└─────────┘ └─────────┘ └─────────┘
```

## Clock Access

All simulation components access time through the clock:

```go
// World
simTime := world.Time()
simTime = world.Clock().Now()

// Entities  
func (e *SunEntity) Tick(dt time.Duration) {
    // Use sim time for calculations
    hour := e.world.Time().Hour()
}

// Scenarios
elapsed := scenario.Elapsed()
```

## Implementation

```go
type SimClock struct {
    mode       Mode
    speed      float64
    running    bool
    startTime  time.Time  // Sim start
    startWall  time.Time   // Wall start
    current    time.Time  // Current sim time
    tick       uint64
}

func (c *SimClock) Update() {
    switch c.mode {
    case ModeRealtime:
        elapsed := time.Since(c.startWall)
        c.current = c.startTime.Add(time.Duration(float64(elapsed) * c.speed))
        c.tick++
    case ModeSimulated:
        wallDelta := time.Since(c.lastWall)
        c.current = c.current.Add(time.Duration(float64(wallDelta) * c.speed))
        c.tick++
    }
}
```

## Future Expansion

### Replay Mode
```go
// Load recorded timeline
replayClock := NewReplayClock(timelineFile)
scheduler.SetClock(replayClock)
```

### Manual Mode
```go
// Step through simulation
for !done {
    clock.Advance(1 * time.Hour)  // Advance 1 hour
    scheduler.Tick()
}
```

### Multi-Speed Regions
```go
// Speed up night, slow down day
clock.SetSpeedFunction(func(t time.Time) float64 {
    hour := t.Hour()
    if hour >= 20 || hour < 6 {
        return 10.0  // Fast at night
    }
    return 1.0  // Real-time during day
})
```

## Best Practices

1. **Never call `time.Now()`** - Use `clock.Now()`
2. **Use duration for intervals** - `clock.Elapsed()` not wall duration
3. **Speed affects simulation** - 2x speed means entities see 2x time per tick
4. **DateTime drives behavior** - Sun, weather depend on simulated date
5. **Test at multiple speeds** - Verify behavior is consistent
