# Solver Architecture

## Purpose

The Solver is responsible for advancing the simulation state. It determines how the simulated world evolves.

## Design Philosophy

**Separation of Concerns**

```
World                    Solver                   Entities
  │                        │                        │
  │  delegates evolution   │  evaluates & updates   │
  ├───┐              ┌─────┘                   ┌────┘
  │   │              │                         │
  ▼   ▼              ▼                         ▼
 Container         Engine                    Behavior
 owns state        advances state             responds
```

**The Solver owns:**
- Evaluation order
- Dependency resolution
- State propagation
- Network traversal
- Simulation iteration

**The Solver does NOT own:**
- Time (World clock)
- Topology (Network)
- Entity behavior (Entities)
- Events (World)
- Measurements (Entities)

## Core Interfaces

### Solver Interface

```go
type Solver interface {
    Name() string              // Solver name
    Type() string              // Solver type
    Tick(dt time.Duration)     // Advance simulation
    Reset()                    // Clear state
    SetWorld(w World)          // Attach world
}
```

### ElectricalSolver

The ElectricalSolver implements power balance calculations.

```go
type ElectricalSolver struct {
    BaseSolver
    topology *topology.Network
    
    // Aggregated values
    totalGeneration    float32
    totalConsumption  float32
    netPower          float32
    totalCapacity      float32
    availableGeneration float32
}
```

**Responsibilities:**
1. Collect all Virtual Generators
2. Collect all Virtual Loads
3. Calculate total generation
4. Calculate total consumption
5. Determine net power
6. Propagate results through topology
7. Update measurement entities

**Tick Order:**
```
1. Collect entities
2. Tick all entities
3. Calculate generation
4. Calculate consumption
5. Determine net power
6. Propagate to meters
```

## World Integration

### Setting a Solver

```go
w := world.NewWorld()
solver := solver.NewElectricalSolver()
w.SetSolver(solver)
```

### Running the Simulation

```go
// World.Tick delegates to Solver.Tick
w.Tick(dt)
```

### Accessing Solver State

```go
solver := w.Solver().(*solver.ElectricalSolver)
fmt.Printf("Net Power: %.1f kW\n", solver.NetPower())
fmt.Printf("Generation: %.1f kW\n", solver.TotalGeneration())
fmt.Printf("Consumption: %.1f kW\n", solver.TotalConsumption())
```

## Solver Types

### Current

| Solver | Type | Purpose |
|--------|------|---------|
| ElectricalSolver | power-balance | Calculate power balance |

### Future

| Solver | Type | Purpose |
|--------|------|---------|
| ACPowerFlowSolver | ac-power-flow | Full AC power flow |
| DCPowerFlowSolver | dc-power-flow | DC power flow |
| HydraulicSolver | hydraulic | Water/steam flow |
| ThermalSolver | thermal | Heat transfer |

## Electrical Solver

### Entity Interfaces

The ElectricalSolver uses interface-based access to entities.

```go
type GeneratorSource interface {
    ActivePower() float32
    RatedCapacity() float32
    AvailableCapacity() float32
    IsOnline() bool
    Name() string
}

type LoadSink interface {
    ActivePowerDemand() float32
    IsConnected() bool
    Name() string
}

type MeterMeasurement interface {
    SetMeasurements(p, q, v, f float32)
    RecordEnergy(dt time.Duration)
    Name() string
}
```

### Power Calculations

```go
// Total generation
totalGen := sum(gen.ActivePower() for gen in onlineGenerators)

// Total consumption
totalLoad := sum(load.ActivePowerDemand() for load in connectedLoads)

// Net power (positive = export, negative = import)
netPower := totalGen - totalLoad
```

### Metrics

```go
// Capacity utilization
util := solver.CapacityUtilization()  // 0.0 - 1.0

// Load served
served := solver.LoadServed()  // 0.0 - 1.0

// Excess capacity
excess := solver.ExcessCapacity()  // kW
```

## Why Separate Solver from World?

### Benefits

1. **Specialization**
   - World: Container for state
   - Solver: Engine for computation

2. **Extensibility**
   - Swap solvers for different physics
   - Combine solvers for coupled systems

3. **Testability**
   - Test solvers independently
   - Mock world for solver tests

4. **Clarity**
   - Clear responsibility boundaries
   - Easier to understand data flow

### Without Solver

```
World.Tick():
  1. Propagate signals
  2. Tick entities
  3. Calculate power flow  ← Mixed concerns
  4. Update meters
  5. Check protection
  6. Run optimization       ← Too much responsibility
```

### With Solver

```
World.Tick():
  1. Delegate to Solver

ElectricalSolver.Tick():
  1. Tick entities
  2. Calculate power balance

ProtectionSolver.Tick():
  1. Check fault conditions
  2. Issue trips

OptimizationSolver.Tick():
  1. Run dispatch optimization
```

## Design Rules

### Solver Rules

1. **Solvers do not create entities**
   - World owns entity lifecycle
   - Solver operates on existing entities

2. **Solvers do not own time**
   - World clock provides time
   - Solver receives dt parameter

3. **Solvers do not modify topology**
   - Topology owns connectivity
   - Solver reads topology for calculations

4. **Solvers may create measurements**
   - Solvers compute derived quantities
   - Measurements are stored in entities

### World Rules

1. **World delegates evolution**
   - World.Tick() calls Solver.Tick()
   - World does not contain simulation mathematics

2. **World owns entities**
   - Add/remove entity lifecycle
   - Solver operates on entities in world

3. **World owns time**
   - Clock provides simulation time
   - Solver receives time step

## Future Solvers

### AC Power Flow Solver

```go
type ACPowerFlowSolver struct {
    BaseSolver
    tolerance float32
    maxIterations int
}

// Newton-Raphson or Fast Decoupled method
func (s *ACPowerFlowSolver) Tick(dt time.Duration) {
    // 1. Build Y-bus matrix
    // 2. Initialize voltages
    // 3. Iterate until convergence
    // 4. Calculate branch flows
    // 5. Update entity measurements
}
```

### Hydraulic Solver

```go
type HydraulicSolver struct {
    BaseSolver
    fluidDensity float32
    gravity float32
}

// Mass/energy balance
func (s *HydraulicSolver) Tick(dt time.Duration) {
    // 1. Calculate flows
    // 2. Update pressures
    // 3. Calculate power output
}
```

### Coupled Solvers

```go
// Run multiple solvers in sequence
type CoupledSolver struct {
    BaseSolver
    solvers []Solver
    order []int
}

func (s *CoupledSolver) Tick(dt time.Duration) {
    for _, idx := range s.order {
        s.solvers[idx].Tick(dt)
    }
}
```

## Example: Power Balance

```go
// Create world with electrical solver
w := world.NewWorld()
w.SetSolver(solver.NewElectricalSolver())

// Add entities
w.AddEntity(world.NewVirtualGenerator("gen-1", "Solar", 500))
w.AddEntity(world.NewVirtualLoad("load-1", "Factory", 400))
w.AddEntity(world.NewVirtualMeter("pcc", "PCC"))

// Simulate
for i := 0; i < 100; i++ {
    w.Tick(100 * time.Millisecond)
    
    // Access solver state
    s := w.Solver().(*solver.ElectricalSolver)
    fmt.Printf("Net: %.1f kW\n", s.NetPower())
}
```

## Glossary

| Term | Definition |
|------|------------|
| Solver | Engine that advances simulation state |
| Power Balance | Generation - Consumption |
| Net Power | Positive = Export, Negative = Import |
| Utilization | Generation / Capacity |
| Load Served | Consumption / Demand |

---

*Last Updated: 2026-07-13*
