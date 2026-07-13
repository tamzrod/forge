# Runtime Sequence Diagram

## Purpose

This diagram shows the **tick execution flow** of the Industrial Simulation Runtime. It answers: *"How does time advance through the system?"*

---

## Tick Loop Overview

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                                                                                     │
│                              SIMULATION TICK LOOP                                    │
│                                                                                     │
│  ┌─────────────┐      ┌─────────────┐      ┌─────────────┐      ┌─────────────┐ │
│  │   Ticker    │ ───▶ │ Scheduler  │ ───▶ │   Models    │ ───▶ │   Devices   │ │
│  │  (Timer)    │ tick │   .tick()  │ tick │   .Tick()   │ tick │   .Tick()   │ │
│  └─────────────┘      └─────────────┘      └─────────────┘      └─────────────┘ │
│                                                                          │         │
│                                                                          ▼         │
│                                                                  ┌─────────────┐ │
│                                                                  │    Clock    │ │
│                                                                  │  Advance    │ │
│                                                                  └─────────────┘ │
│                                                                                     │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

---

## Detailed Sequence: Single Tick

```
Actor       System                      Scheduler                    Models                      Devices                     Memory
   │           │                            │                           │                           │                          │
   │           │                            │                           │                           │                          │
   │           │                            │                           │                           │                          │
   │           │    ┌──────────────────────┴───────────────────────────┴───────────────────────────┴───────────────────────┐ │
   │           │    │                                    T I C K   L O O P                                              │ │
   │           │    └────────────────────────────────────────────────────────────────────────────────────────────────────────┘ │
   │           │                            │                                                                              │
   │           │                            │                                                                              │
   │  [Timer]  │                            │                                                                              │
   │    fires   │                            │                                                                              │
   │───────────▶│                            │                                                                              │
   │           │                            │                                                                              │
   │           │                            │ tick()                                                                       │
   │           │                            │────────┐                                                                     │
   │           │                            │        │ Acquire lock                                                      │
   │           │                            │◀───────┘                                                                     │
   │           │                            │                                                                              │
   │           │                            │                                                                              │
   │           │                            │ foreach model                                                                 │
   │           │                            │────────▶──────────────────────                                               │
   │           │                            │        │              │ model.Tick()                                         │
   │           │                            │◀────────│──────────────┘                                                      │
   │           │                            │                                                                              │
   │           │                            │                                                                              │
   │           │                            │ foreach device                                                               │
   │           │                            │────────▶────────────────────────────────────────────────────────────────▶    │
   │           │                            │        │              │              │            │            │           │    │
   │           │                            │        │              │              │            │            │           │    │
   │           │                            │        │              │              │            │            │           │    │
   │           │                            │        │              │              │            │            │           │    │
   │           │                            │        │              │    Observe   │            │            │           │    │
   │           │                            │        │              │◀─────────────│────────────▶│            │           │    │
   │           │                            │        │              │   models     │            │            │           │    │
   │           │                            │        │              │              │            │            │           │    │
   │           │                            │        │              │              │    Update   │            │           │    │
   │           │                            │        │              │              │◀───────────┼────────────▶│           │    │
   │           │                            │        │              │              │   memory   │            │           │    │
   │           │                            │        │              │              │            │            │           │    │
   │           │                            │        │              │              │            │            │           │    │
   │           │                            │◀────────│──────────────│──────────────│────────────│────────────│───────────┘    │
   │           │                            │                                                                              │
   │           │                            │ clock.elapsed += tickInterval                                                 │
   │           │                            │─────────────────────────────────────────────────────────────────────────────▶[Clock]│
   │           │                            │                                                                              │
   │           │                            │ clock.tickCount++                                                            │
   │           │                            │─────────────────────────────────────────────────────────────────────────────▶[Clock]│
   │           │                            │                                                                              │
   │           │                            │ Release lock                                                                  │
   │           │                            │────────┐                                                                     │
   │           │                            │        │ Release                                                             │
   │           │                            │◀───────┘                                                                     │
   │           │                            │                                                                              │
   │           │                            │                                                                              │
   │           │                            │                                                                              │
   │           │    ┌──────────────────────┴───────────────────────────┴───────────────────────────┴───────────────────────┐ │
   │           │    │                              E N D   T I C K                                              │ │
   │           │    └────────────────────────────────────────────────────────────────────────────────────────────────────────┘ │
   │           │                            │                                                                              │
   │           │                            │                                                                              │
   │  [Return to wait]                     │                                                                              │
   │◀──────────│                            │                                                                              │
   │           │                            │                                                                              │
```

---

## State Transitions

### Runtime State Machine

```
                    ┌──────────────┐
                    │   Created    │
                    └──────┬───────┘
                           │ NewRuntime()
                           ▼
                    ┌──────────────┐
          ┌───────▶│  Configured  │◀──────┐
          │         └──────┬───────┘       │
          │                │               │
          │    CreateDevice/              │ LoadConfig
          │    CreateModel                 │
          │                │               │
          │                ▼               │
          │         ┌──────────────┐      │
          │         │   Running     │──────┘
          │         │  (ticking)    │
          │         └──────┬───────┘
          │                │
          │     Run(ctx)   │ ctx cancelled
          │                │ or Stop()
          │                ▼
          │         ┌──────────────┐
          │         │   Stopped    │
          │         └──────┬───────┘
          │                │
          │    Shutdown()  │
          │                ▼
          │         ┌──────────────┐
          └─────────│  Destroyed   │
                    └──────────────┘
```

### Device State Machine

```
                    ┌──────────────┐
                    │  StateCreated │
                    └──────┬───────┘
                           │ New()
                           ▼
                    ┌──────────────┐
                    │StateInitialized│
                    └──────┬───────┘
                           │ Start()
                           ▼
                    ┌──────────────┐
                    │ StateRunning │
                    │   (ticking)  │
                    └──────┬───────┘
                           │ Stop()
                           ▼
                    ┌──────────────┐
                    │ StateStopped │
                    └──────┬───────┘
                           │ Destroy()
                           ▼
                    ┌──────────────┐
                    │StateDestroyed│
                    └──────────────┘
```

---

## Clock Advancement

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                              CLOCK STRUCTURE                                       │
│                                                                                     │
│  type SimulationClock struct {                                                      │
│      elapsed   time.Duration  // Total simulated time elapsed                      │
│      tickCount uint64         // Number of ticks executed                          │
│  }                                                                                 │
│                                                                                     │
│  type Scheduler struct {                                                           │
│      tickInterval time.Duration  // Time between ticks (e.g., 100ms)              │
│      clock      SimulationClock                                                    │
│  }                                                                                 │
│                                                                                     │
└─────────────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────────────┐
│                              TIME PROGRESSION                                      │
│                                                                                     │
│  tickInterval: 100ms                                                              │
│                                                                                     │
│  Tick 0:  elapsed =    0ms  ──┬──▶  clock = {elapsed: 0ms, tickCount: 0}     │
│                                  │                                                │
│  Tick 1:  elapsed =  100ms  ──┼──▶  clock = {elapsed: 100ms, tickCount: 1}    │
│                                  │                                                │
│  Tick 2:  elapsed =  200ms  ──┼──▶  clock = {elapsed: 200ms, tickCount: 2}    │
│                                  │                                                │
│  Tick 3:  elapsed =  300ms  ──┼──▶  clock = {elapsed: 300ms, tickCount: 3}    │
│                                  │                                                │
│  Tick N:  elapsed = N*100ms  ──┘  clock = {elapsed: N*100ms, tickCount: N}     │
│                                                                                     │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

---

## Tick Order Guarantee

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                            TICK ORDER GUARANTEE                                     │
│                                                                                     │
│  The Scheduler guarantees:                                                         │
│                                                                                     │
│  1. Models tick BEFORE devices (physics → equipment)                                │
│  2. Devices tick in registration order                                             │
│  3. Behaviors tick in attachment order within a device                             │
│                                                                                     │
│  This ensures deterministic, reproducible execution.                                  │
│                                                                                     │
└─────────────────────────────────────────────────────────────────────────────────────┘

Execution Order:

  Order  │  Component              │  Reason
  ────────┼────────────────────────┼──────────────────────────────────────
  1       │  GridModel.Tick()       │  Physics first
  2       │  SunModel.Tick()        │  Physics first
  3       │  WindModel.Tick()       │  Physics first
  4       │  WeatherModel.Tick()    │  Physics first
  5       │  ReservoirModel.Tick()  │  Physics first
  ────────┼────────────────────────┼──────────────────────────────────────
  6       │  WeatherStation.Tick()   │  Equipment observes physics
  7       │  PVInverter.Tick()      │  Equipment observes physics
  8       │  RevenueMeter.Tick()    │  Equipment observes physics
  ────────┼────────────────────────┼──────────────────────────────────────
  N+1     │  Clock Advance          │  Time advances last
```

---

## Behavior Execution Within Device

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                          DEVICE BEHAVIOR EXECUTION                                   │
│                                                                                     │
│  Device.Tick()                                                                      │
│       │                                                                             │
│       ├──▶ Behavior[0].Tick()                                                       │
│       │         │                                                                   │
│       │         ├──▶ Read Weather Model                                            │
│       │         │         │                                                         │
│       │         │         ├──▶ Calculate irradiance                                 │
│       │         │         │                                                         │
│       │         │         └──▶ Write to Memory (sensors region)                   │
│       │         │                                                                   │
│       │         └──▶ Return                                                       │
│       │                                                                             │
│       ├──▶ Behavior[1].Tick()                                                       │
│       │         │                                                                   │
│       │         ├──▶ Read Memory (sensors region)                                │
│       │         │         │                                                         │
│       │         │         ├──▶ Calculate power output                               │
│       │         │         │                                                         │
│       │         │         └──▶ Write to Memory (output region)                     │
│       │         │                                                                   │
│       │         └──▶ Return                                                       │
│       │                                                                             │
│       └──▶ Return                                                                  │
│                                                                                     │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

---

## Initialization Sequence

```
Actor       Developer              Runtime                   Scheduler                  Device
   │           │                      │                          │                         │
   │           │                      │                          │                         │
   │  NewRuntime(cfg)                │                          │                         │
   │──────────▶│                      │                          │                         │
   │           │   Create Runtime     │                          │                         │
   │           │────────────────────▶│                          │                         │
   │           │                      │                          │                         │
   │           │                      │   Scheduler.New()        │                         │
   │           │                      │────────────────────────▶│                         │
   │           │                      │                          │                         │
   │           │                      │                          │   clock.elapsed = 0      │
   │           │                      │                          │   clock.tickCount = 0    │
   │           │                      │                          │                         │
   │           │   return Runtime     │                          │                         │
   │◀──────────│◀────────────────────│                          │                         │
   │           │                      │                          │                         │
   │           │                      │                          │                         │
   │  CreateDevice(...)              │                          │                         │
   │──────────▶│                      │                          │                         │
   │           │                      │                          │                         │
   │           │                      │   device.New()           │                         │
   │           │                      │─────────────────────────│────────────────────────▶│
   │           │                      │                          │                         │
   │           │                      │                          │   mem = MemoryImage.New()│
   │           │                      │                          │   behaviors = []        │
   │           │                      │                          │   running = false       │
   │           │                      │                          │                         │
   │           │                      │   AddDevice()           │                         │
   │           │                      │─────────────────────────▶│                         │
   │           │                      │                          │                         │
   │           │                      │                          │   devices = append(...)  │
   │           │                      │                          │                         │
   │           │   return Device      │                          │                         │
   │◀──────────│◀────────────────────│                          │                         │
   │           │                      │                          │                         │
   │           │                      │                          │                         │
   │  runtime.Run(ctx)              │                          │                         │
   │──────────▶│                      │                          │                         │
   │           │                      │                          │                         │
   │           │                      │   Start()               │                         │
   │           │                      │─────────────────────────│────────────────────────▶│
   │           │                      │                          │                         │
   │           │                      │                          │   running = true         │
   │           │                      │                          │                         │
   │           │                      │                          │                         │
   │           │                      │   scheduler.Run(ctx)     │                         │
   │           │                      │────────────────────────▶│                         │
   │           │                      │                          │                         │
   │           │                      │                          │  [Tick Loop Begins]     │
   │           │                      │                          │                         │
   │           │                      │                          │◀────── tick() ──────    │
   │           │                      │                          │                         │
   │           │                      │   [Ticker fires]        │                         │
   │           │                      │◀─────────────────────────│                         │
   │           │                      │                          │                         │
   │           │                      │                          │                         │
   │           │                      │◀────── tick() ──────────│                         │
   │           │                      │                          │                         │
   │           │                      │                          │                         │
   │           │                      │                          │                         │
   │           │                      │                          │                         │
   │  [Context cancelled]             │                          │                         │
   │──────────▶│                      │                          │                         │
   │           │                      │                          │                         │
   │           │                      │   [Return]              │                         │
   │◀──────────│◀────────────────────│◀─────────────────────────▶│                         │
   │           │                      │                          │                         │
```

---

## Determinism Guarantee

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                           DETERMINISM GUARANTEE                                    │
│                                                                                     │
│  Same inputs → Same outputs, every time                                            │
│                                                                                     │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                     │
│  What Makes It Deterministic:                                                       │
│                                                                                     │
│  ✓ Fixed tick interval                                                             │
│  ✓ Models tick in registration order                                               │
│  ✓ Devices tick in registration order                                              │
│  ✓ Behaviors tick in attachment order                                                │
│  ✓ No external dependencies (time, network, randomness)                              │
│                                                                                     │
│  How to Break Determinism:                                                         │
│                                                                                     │
│  ✗ Use time.Now() instead of simulation clock                                     │
│  ✗ Read from network/sensors                                                       │
│  ✗ Use math/rand without seeding                                                  │
│  ✗ Access shared global state                                                     │
│                                                                                     │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

---

## Verification Points

To verify this diagram is accurate, check:

1. **Tick order is preserved** - Models before devices
2. **Clock advances last** - After all ticks complete
3. **No concurrent modifications** - Lock held during tick
4. **State transitions are valid** - Only valid transitions allowed
5. **Determinism maintained** - Same inputs produce same outputs

---

## Related Documents

| Document | Purpose |
|----------|---------|
| [Context Diagram](context-diagram.md) | System boundaries |
| [Component Diagram](component-diagram.md) | Module structure |
| [Scheduler](scheduler.md) | Scheduler implementation |
| [Runtime](runtime.md) | Runtime implementation |

---

*Created: 2026-07-13*  
*Type: Architecture Artifact*  
*Status: Initial*
