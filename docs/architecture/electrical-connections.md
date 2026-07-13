# Electrical Connections Knowledge Base

## Purpose

This document establishes the authoritative engineering concepts for electrical connections in Forge. It defines how entities connect to the electrical network through terminals.

## Core Concepts

### Terminal

A **Terminal** is a connection point on an entity where electrical conductors connect.

```
Entity
  ├── Terminal "grid"      → Connects to utility grid
  ├── Terminal "hv"        → High voltage side
  └── Terminal "lv"        → Low voltage side
```

**Terminal Types by Role:**

| Terminal Role | Purpose | Example |
|---------------|---------|---------|
| **Source** | Injects power into network | Generator terminal |
| **Destination** | Consumes power from network | Load terminal |
| **Through** | Passes power through | Transformer, cable |
| **Observation** | Measures without affecting | Meter terminal |

**Terminal Types by Voltage:**

| Terminal Type | Purpose | Example |
|---------------|---------|---------|
| **HV** | High voltage side | Transformer primary |
| **LV** | Low voltage side | Transformer secondary |
| **Grid** | Utility grid connection | PCC terminal |

### Connection

A **Connection** represents the electrical link between an entity terminal and a bus.

```
Terminal ──Connection──> Bus
   │                      │
   │                      ├── Terminal
   │                      ├── Branch
   └── Entity             └── Other Connections
```

**Connection Types:**

| Type | Direction | Purpose |
|------|-----------|---------|
| **Injection** | Terminal → Bus | Power flowing from entity to network |
| **Withdrawal** | Bus → Terminal | Power flowing from network to entity |
| **Through** | Bus → Terminal → Bus | Power passing through entity |
| **Observation** | Bus → Terminal | Measurement without power flow |

### Bus Connection

A **Bus Connection** is the junction point where conductors meet.

```
Bus (69kV PCC)
  │
  ├── Grid Terminal (Source)
  ├── Transformer HV Terminal (Through)
  └── Meter Terminal (Observation)
```

### Valid Connection Matrix

| Entity Type | Source Terminal | Dest Terminal | Through Terminal | Observation Terminal |
|-------------|-----------------|---------------|------------------|---------------------|
| VirtualGenerator | ✓ (required) | ✗ | ✗ | ○ |
| VirtualLoad | ✗ | ✓ (required) | ✗ | ○ |
| VirtualMeter | ✗ | ✗ | ✗ | ✓ (required) |
| Transformer | ○ | ○ | ✓ (required) | ○ |
| Breaker | ○ | ○ | ✓ (optional) | ○ |
| Switch | ○ | ○ | ✓ (optional) | ○ |

✓ = Allowed, ✗ = Not allowed, ○ = Optional

## Connection Rules

### Rule 1: Every Power-Injecting Entity Requires a Source Terminal

Generators, inverters, and other sources must have a source terminal.

```
Solar Generator
  └── Terminal "output" (Source)
            │
            ▼
        480V Bus
```

### Rule 2: Every Power-Consuming Entity Requires a Destination Terminal

Loads must have a destination terminal.

```
480V Bus
    │
    ▼
Load (Factory)
  └── Terminal "input" (Destination)
```

### Rule 3: Transformers Have Through Terminals

Transformers connect two different voltage levels through their terminals.

```
69kV Bus ──Terminal "hv"──> Transformer ──Terminal "lv"──> 480V Bus
                Through                          Through
```

### Rule 4: Meters Have Observation Terminals

Meters measure power flow without injecting or withdrawing.

```
69kV Bus
    │
    ▼
Meter (PCC)
  └── Terminal "meter" (Observation)
            │
            ▼
        (Measurement only)
```

### Rule 5: Terminal Voltage Must Match Bus Voltage

A terminal can only connect to a bus with matching nominal voltage.

```
Valid:   Terminal (480V) → Bus (480V) ✓
Invalid: Terminal (480V) → Bus (69000V) ✗
```

Exception: Transformers with through terminals can connect different voltages.

### Rule 6: Each Terminal Connects to One Bus

A terminal can only be connected to a single bus at a time.

```
Terminal "output" ──> Bus A ✓
Terminal "output" ──> Bus B ✗
```

## Connection Validation

### Voltage Compatibility

```go
func CanConnect(terminal *Terminal, bus *Bus) bool {
    // Observation terminals can connect to any voltage
    if terminal.Role == TerminalRoleObservation {
        return true
    }
    
    // Through terminals on transformers can connect to any voltage
    if terminal.Role == TerminalRoleThrough {
        return terminal.Voltage == bus.NominalVoltage
    }
    
    // Source and Destination terminals must match voltage
    return terminal.Voltage == bus.NominalVoltage
}
```

### Terminal Uniqueness

```go
func IsTerminalAvailable(terminal *Terminal) bool {
    return terminal.Bus() == nil
}
```

## Network Construction

### Step 1: Define Buses

```go
// Create buses at different voltage levels
hvBus := topology.NewBus("hv-bus", "69kV PCC", 69000)
lvBus := topology.NewBus("lv-bus", "480V Collector", 480)
```

### Step 2: Create Branches

```go
// Create transformer branch
txBranch := topology.NewBranch("tx", "Transformer", hvBus, lvBus)
```

### Step 3: Connect Entities to Buses

```go
// Connect generator to LV bus
genTerminal := topology.NewTerminal(
    "gen-t1",
    "solar-gen-1",
    "output",
    topology.TerminalRoleSource,
    topology.TerminalTypeLV,
    480,
)
net.AddTerminal(genTerminal)
lvBus.AddTerminal(genTerminal)

// Connect load to LV bus
loadTerminal := topology.NewTerminal(
    "load-t1",
    "factory",
    "input",
    topology.TerminalRoleDestination,
    topology.TerminalTypeLV,
    480,
)
net.AddTerminal(loadTerminal)
lvBus.AddTerminal(loadTerminal)

// Connect meter to HV bus
meterTerminal := topology.NewTerminal(
    "meter-t1",
    "pcc-meter",
    "meter",
    topology.TerminalRoleObservation,
    topology.TerminalTypeHV,
    69000,
)
net.AddTerminal(meterTerminal)
hvBus.AddTerminal(meterTerminal)
```

## Connection Diagram

```
                    ┌─────────────────────────────────────┐
                    │         UTILITY GRID                │
                    │  (Source Terminal, 69kV)             │
                    └──────────────┬──────────────────────┘
                                   │
                    ┌──────────────▼──────────────────────┐
                    │         69kV PCC BUS                │
                    │  Nominal Voltage: 69000V            │
                    │                                     │
                    │  ┌─ Grid Terminal (Source)         │
                    │  ├─ Meter Terminal (Observation)    │
                    │  └─ Transformer HV Terminal (Through)│
                    └──────────────┬──────────────────────┘
                                   │
                    ┌──────────────▼──────────────────────┐
                    │         TRANSFORMER                 │
                    │  HV Terminal ──> LV Terminal      │
                    │  69kV            480V              │
                    └──────────────┬──────────────────────┘
                                   │
                    ┌──────────────▼──────────────────────┐
                    │         480V COLLECTOR BUS         │
                    │  Nominal Voltage: 480V             │
                    │                                     │
                    │  ┌─ PV Array Terminal (Source)     │
                    │  ├─ PV Array Terminal (Source)     │
                    │  ├─ Factory Terminal (Destination)  │
                    │  └─ Aux Load Terminal (Destination)│
                    └─────────────────────────────────────┘
```

## Entity Integration

### VirtualGenerator Terminals

Each VirtualGenerator must have at least one source terminal.

```go
type VirtualGenerator struct {
    BaseEntity
    terminals []*Terminal
}

// GetSourceTerminal returns the primary source terminal.
func (g *VirtualGenerator) GetSourceTerminal() *Terminal {
    for _, t := range g.terminals {
        if t.Role == TerminalRoleSource {
            return t
        }
    }
    return nil
}
```

### VirtualLoad Terminals

Each VirtualLoad must have at least one destination terminal.

```go
type VirtualLoad struct {
    BaseEntity
    terminals []*Terminal
}

// GetDestinationTerminal returns the primary destination terminal.
func (l *VirtualLoad) GetDestinationTerminal() *Terminal {
    for _, t := range l.terminals {
        if t.Role == TerminalRoleDestination {
            return t
        }
    }
    return nil
}
```

### VirtualMeter Terminals

Each VirtualMeter must have at least one observation terminal.

```go
type VirtualMeter struct {
    BaseEntity
    terminals []*Terminal
}

// GetObservationTerminal returns the observation terminal.
func (m *VirtualMeter) GetObservationTerminal() *Terminal {
    for _, t := range m.terminals {
        if t.Role == TerminalRoleObservation {
            return t
        }
    }
    return nil
}
```

## Glossary

| Term | Definition |
|------|------------|
| Terminal | Connection point on an entity |
| Source Terminal | Terminal that injects power |
| Destination Terminal | Terminal that withdraws power |
| Through Terminal | Terminal that passes power through |
| Observation Terminal | Terminal that measures without affecting |
| Connection | Link between terminal and bus |
| Injection | Power flowing from entity to network |
| Withdrawal | Power flowing from network to entity |

---

*Last Updated: 2026-07-13*
