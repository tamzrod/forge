# Electrical Topology Model

## Purpose

The Electrical Topology model represents the structure of the electrical network. It is the authoritative source of network connectivity.

## Design Principle

**Topology owns connectivity. Entities own behavior.**

```
Topology                    Entities
├─ Buses                   ├─ Internal behavior
├─ Branches                ├─ Measurements
├─ Switches                ├─ State
└─ Terminals               └─ Event handling
```

Topology knows **what connects to what**. Entities know **how they work**.

## Core Concepts

### Bus

A **Bus** is an electrical node where conductors connect.

```go
type Bus struct {
    ID             ID
    Name           string
    NominalVoltage float32  // V
}
```

Properties:
- Each bus has a nominal voltage (e.g., 69000V, 480V)
- Multiple terminals and branches connect to a bus
- Buses are identified by unique IDs

### Branch

A **Branch** is a connection between two buses.

```go
type Branch struct {
    ID       ID
    Name     string
    FromBus *Bus  // Source bus
    ToBus   *Bus  // Destination bus
}
```

Properties:
- Branches contain switching devices (breakers)
- Power flows through branches
- Branches connect two buses

### Terminal

A **Terminal** is a connection point on an entity.

```go
type Terminal struct {
    ID       ID
    EntityID world.EntityID  // Owner entity
    Name     string          // e.g., "primary", "secondary"
}
```

Properties:
- Each entity can have multiple terminals
- Terminals connect entities to buses
- Example: A breaker has "line" and "load" terminals

### Switch

A **Switch** is a switching device (breaker) in a branch.

```go
type Switch struct {
    ID       ID
    Name     string
    isOpen   bool
}
```

Operations:
- `Open()` - Opens the switch, interrupting current flow
- `Close()` - Closes the switch, allowing current flow
- `IsOpen()` - Returns switch state

### Network

The **Network** contains all topology elements.

```go
type Network struct {
    buses    map[ID]*Bus
    branches map[ID]*Branch
    switches map[ID]*Switch
}
```

## Responsibilities

### Topology Owns

| Responsibility | Description |
|----------------|-------------|
| Connectivity | Which entities connect to which buses |
| Upstream/Downstream | Direction relationships |
| Network traversal | Finding paths between buses |
| Island detection | Identifying disconnected subgraphs |
| Connected components | Grouping related elements |

### Entities Own

| Responsibility | Description |
|----------------|-------------|
| Internal behavior | How the entity operates |
| Measurements | What the entity measures |
| State | Internal state (not connectivity) |
| Event handling | Response to commands |

## Topology Queries

### What is connected to this bus?

```go
entities := net.EntitiesConnectedTo(bus)
```

Returns all entities attached to a bus.

### What is upstream/downstream?

```go
// Entities toward sources (generators)
upstream := net.EntitiesUpstream(bus)

// Entities away from sources (loads)
downstream := net.EntitiesDownstream(bus)
```

### Which entities become isolated?

```go
// Returns entities that would lose grid connection
isolated := net.IsolatedIf("grid-breaker")
```

### Is this bus energized?

```go
energized := net.IsBusEnergized(bus)
```

### Which island does this entity belong to?

```go
island := net.IslandFor(bus)
```

### Find path between buses

```go
path := net.PathBetween(fromBus, toBus)
```

## Network Structure

Example: Simple PV Plant

```
Utility Grid (69kV)
        ↓
    [Grid Breaker]
        ↓
    69kV Bus (PCC)
        ↓
    [Transformer]
        ↓
    480V Bus (Collector)
        ↓
    [CB Breaker]
        ↓
    480V Bus (Array)
        ↓
    PV Array
        ↓
    Aux Load
```

## Builder

The `Builder` simplifies network construction:

```go
builder := topology.NewBuilder()

// Build simple radial network
builder.BuildSimpleRadial()

net := builder.Network()
```

### Builder Methods

```go
AddBus(id, name, voltage)
AddBranch(id, name, fromID, toID)
AddBreaker(breakerID, branchID, name)
AddTransformer(id, name, highBusID, lowBusID, ratio)
ConnectEntity(entityID, busID, terminalName)
```

## Voltage Levels

```go
type VoltageLevel int

const (
    VoltageLevelLow       // < 1kV
    VoltageLevelMedium    // 1kV - 35kV
    VoltageLevelHigh      // 35kV - 230kV
    VoltageLevelExtraHigh // > 230kV
)
```

## Island Detection

When breakers open, the network splits into islands:

```
Initial (all closed):
  Island 0: [Grid] - [69kV] - [480V] - [Array]

After Grid Breaker opens:
  Island 0: [Grid] - [69kV]
  Island 1: [480V] - [Array]
```

Islands are calculated automatically based on switch states.

## Design Rules

### Topology Never Contains Electrical Behavior

Topology only describes structure. It does not:
- Calculate power flow
- Model voltage drop
- Track power injection

### Entities Never Own Network Connectivity

Entities do not:
- Know which other entities they connect to
- Access other entities directly
- Modify topology structure

## Future Expansion

### Power Flow Integration

Power flow calculations will use topology:

```go
// Future: Calculate power flow
powerFlow := topology.CalculateFlow(net)
```

### Protection Coordination

Protection decisions will use topology:

```go
// Future: Find affected relays
relays := topology.RelaysDownstream(breaker)
```

### State Estimation

State estimation will use topology:

```go
// Future: Build measurement model
model := topology.BuildMeasurementModel(net)
```

## Example: Creating a Network

```go
net := topology.NewNetwork()

// Add buses
net.AddBus(topology.NewBus("hv", "69kV Bus", 69000))
net.AddBus(topology.NewBus("lv", "480V Bus", 480))

// Add branch
branch := topology.NewBranch("tx", "Transformer", hvBus, lvBus)
net.AddBranch(branch)

// Add breaker
breaker := topology.NewSwitch("cb", "Circuit Breaker", topology.SwitchTypeBreaker)
breaker.SetBranch(branch)
branch.SetSwitch(breaker)
net.AddSwitch(breaker)

// Connect entities
gridTerminal := topology.NewTerminal("grid-t", "utility-grid", "grid")
gridTerminal.bus = hvBus
hvBus.AddTerminal(gridTerminal)
net.AddTerminal(gridTerminal)

// Query
connected := net.EntitiesConnectedTo(hvBus)
// Returns: ["utility-grid", ...]
```

## Glossary

| Term | Definition |
|------|------------|
| Bus | Electrical node where conductors connect |
| Branch | Connection between two buses |
| Terminal | Connection point on an entity |
| Switch | Device that interrupts current (breaker) |
| Island | Disconnected subgraph of the network |
| Upstream | Toward sources (generators) |
| Downstream | Away from sources (loads) |

---

*Last Updated: 2026-07-13*
