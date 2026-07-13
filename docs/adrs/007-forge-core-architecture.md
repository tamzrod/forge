# ADR-007: Forge Core Architecture

**ADR ID:** ADR-007  
**Title:** Forge Core Architecture Freeze  
**Date:** 2026-07-13  
**Status:** Accepted  
**Deciders:** Engineering Team  
**Repository:** https://github.com/tamzrod/forge

---

## Context

Forge has evolved from a simulation engine with domain-specific components toward an extensible engineering platform. This ADR freezes the permanent Forge Core architecture and establishes the boundaries that must remain stable for future engineering domains.

The previous architecture had accumulated domain-specific leakage:
- Electrical topology types in `topology/` package
- Voltage validation in `registry/` package
- Equipment-centric entity types

This ADR defines what belongs in Forge Core vs. what belongs in plugins.

---

## Decision

We adopt the following architecture principles:

1. **Forge Core is domain-independent** — No engineering-domain logic
2. **Plugins provide domain knowledge** — Components, validators, solvers, templates
3. **Capabilities replace equipment types** — Domain-independent abstractions
4. **World Templates are plugin-provided** — Starting points for simulations

---

## Forge Core Boundaries

### What Forge Core OWNS

| Component | Package | Purpose |
|-----------|---------|---------|
| World | `world/` | Entity container, tick coordination |
| Entity | `world/` | Base interface with capabilities |
| Capability | `world/` | Domain-independent entity abilities |
| Clock | `simulation/` | Simulation time advancement |
| Solver Interface | `world/` | State advancement contract |
| Topology Interface | `topology/` | Structural relationships contract |
| Event | `world/` | Simulation event notifications |
| Plugin Interface | `plugin/` | Plugin contract |
| Plugin Manager | `plugin/` | Plugin lifecycle management |
| Component Catalog | `plugin/` | Component metadata storage |
| Factory Registry | `plugin/` | Entity factory storage |
| ConnectionValidator Interface | `plugin/` | Validation contract |
| WorldTemplate Interface | `plugin/` | World initialization contract |

### What Forge Core Does NOT OWN

| Component | Owner | Rationale |
|-----------|-------|-----------|
| Electrical topology | `plugins/electrical/topology/` | Domain-specific |
| Electrical validators | `plugins/electrical/validators/` | Domain-specific |
| Electrical entities | `world/electrical/` | Domain-specific |
| Water topology | `plugins/water/topology/` | Domain-specific |
| Water validators | `plugins/water/validators/` | Domain-specific |
| Water entities | `world/water/` | Domain-specific |
| Component definitions | Domain Plugins | Domain-specific |
| Solvers | Domain Plugins | Domain-specific |

---

## Capability Framework

Capabilities are domain-independent abstractions that replace equipment-centric assumptions:

```go
type Capability string

const (
    CapabilityProduce     Capability = "produce"     // Generates energy/material
    CapabilityConsume    Capability = "consume"     // Uses energy/material
    CapabilityStore      Capability = "store"       // Holds energy/material
    CapabilityTransform  Capability = "transform"   // Changes form
    CapabilityTransport  Capability = "transport"  // Moves through network
    CapabilitySwitch     Capability = "switch"     // Interrupts flow
    CapabilityMeasure    Capability = "measure"      // Observes values
    CapabilityProtect    Capability = "protect"     // Responds to faults
    CapabilityCommunicate Capability = "communicate" // Sends/receives data
)
```

### Entity with Capabilities

```go
type Entity interface {
    ID() EntityID
    Type() string
    Capabilities() []Capability       // Domain-independent
    HasCapability(capability Capability) bool
    Tick(dt time.Duration)
    Measurements() []Measurement
    Inputs() []Input
    Outputs() []Output
    HandleEvent(evt Event)
    Connect(inputName string, source EntityID, outputName string)
}
```

### World Capability Queries

```go
type World interface {
    EntitiesByCapability(capability Capability) []Entity
    EntitiesByCapabilities(capabilities []Capability) []Entity
    // ... other methods
}
```

### Example: Battery Entity

A battery has multiple capabilities:

```go
// Battery has: Produce, Consume, Store
type BatteryEntity struct {
    world.BaseEntity
    // ...
}

func NewBattery(id world.EntityID) *BatteryEntity {
    return &BatteryEntity{
        BaseEntity: world.NewBaseEntityWithCapabilities(
            id,
            "battery",
            []world.Capability{
                world.CapabilityProduce,
                world.CapabilityConsume,
                world.CapabilityStore,
            },
        ),
    }
}
```

---

## World Template Framework

World Templates are predefined initializations provided by plugins:

```go
type WorldTemplate interface {
    ID() string
    Name() string
    Description() string
    Domain() string
    Build() (world.World, error)
}
```

### Example Templates

```go
// In plugins/electrical/
func (p *Plugin) Templates() []plugin.WorldTemplate {
    return []plugin.WorldTemplate{
        &EmptyWorldTemplate{},
        &SolarFarmTemplate{},
        &BatteryStorageTemplate{},
        &DistributionFeederTemplate{},
    }
}
```

### Template Categories

| Domain | Templates |
|--------|-----------|
| Electrical | Empty, Solar Farm, Battery Storage, Distribution Feeder, Microgrid |
| Water | Empty, Water Distribution, Wastewater Treatment |
| Building | Empty, HVAC, Lighting |
| Process | Empty, Manufacturing Line |

---

## Plugin Contract

Finalized Plugin interface:

```go
type Plugin interface {
    // Identity
    ID() string
    Name() string
    Version() string
    Description() string
    Dependencies() []string
    
    // Lifecycle
    OnInit(ctx Context) error
    OnShutdown() error
    
    // Contributions
    Components() []*ComponentDescriptor
    Categories() []*ComponentCategory
    Validators() []ConnectionValidator
    Templates() []WorldTemplate
}
```

---

## Domain Leakage Resolution

### Before (Domain Leakage)

```
registry/registry.go:
    CanConnect() { /* electrical voltage logic */ }
    
topology/electrical.go:
    TerminalRole, Bus, Branch, Switch  // Electrical types
```

### After (Domain Isolation)

```
plugin/validators/electrical.go:
    ElectricalValidator.CanConnect() { /* electrical logic */ }
    
plugins/electrical/topology/:
    Bus, Branch, Switch, Terminal  // Electrical types
```

### Verification Command

```bash
# Verify no domain logic in Forge Core
grep -r "voltage\|electrical\|water" world/ --include="*.go" | grep -v "_test.go"
grep -r "electrical\|water" topology/ --include="*.go" | grep -v "plugins/"
```

Expected output: No matches (except in tests).

---

## Adding a New Domain

To add a new engineering domain (e.g., Water):

### Step 1: Create Plugin

```go
// plugins/water/water.go
package water

type Plugin struct{}

func (p *Plugin) ID() string    { return "forge-water" }
func (p *Plugin) Name() string    { return "Water Plugin" }
func (p *Plugin) Version() string { return "1.0.0" }
// ... implement Plugin interface
```

### Step 2: Define Components

```go
// plugins/water/components/pump.go
var PumpComponent = types.Component{
    ID:       "forge-water:pump",
    Name:     "Water Pump",
    Category: "water",
    // ...
}
```

### Step 3: Implement Entities

```go
// world/water/pump.go
type PumpEntity struct {
    world.BaseEntity
    flowRate float32
}

func NewPump(id world.EntityID) *PumpEntity {
    return &PumpEntity{
        BaseEntity: world.NewBaseEntityWithCapabilities(
            id, "pump",
            []world.Capability{world.CapabilityProduce, world.CapabilityTransport},
        ),
    }
}
```

### Step 4: Add Validators

```go
// plugins/water/validators/hydraulic.go
type HydraulicValidator struct{}

func (v *HydraulicValidator) Domain() string { return "water" }
func (v *HydraulicValidator) CanConnect(source, target *TerminalDescriptor) (bool, error) {
    // Water-specific validation
}
```

### Step 5: Add Templates

```go
// plugins/water/water.go
func (p *Plugin) Templates() []plugin.WorldTemplate {
    return []plugin.WorldTemplate{
        &WaterDistributionTemplate{},
        &WastewaterTreatmentTemplate{},
    }
}
```

### Step 6: Register Plugin

```go
// plugins/water/init.go (or in plugin/main.go)
func init() {
    plugin.GetManager().Register(&water.Plugin{})
}
```

### Step 7: NO FORGE CORE CHANGES

---

## Verification Checklist

Before any PR that touches core packages, verify:

- [ ] No electrical-specific types in `world/` (except `world/electrical/`)
- [ ] No water-specific types in `world/` (except `world/water/`)
- [ ] No voltage/power references in `world/world.go`
- [ ] No pressure/flow references in `world/world.go`
- [ ] `topology/` contains only generic interfaces
- [ ] Electrical types in `plugins/electrical/topology/`
- [ ] `registry/` contains only metadata (no validation logic)
- [ ] Electrical validation in `plugins/electrical/validators/`

---

## Consequences

### Positive

- **Domain Independence** — Forge Core contains no domain logic
- **Stable Core** — Adding domains doesn't require core changes
- **Testable** — Core can be tested without domain implementations
- **Extensible** — Third parties can add domains
- **Clean Boundaries** — Clear ownership of responsibilities

### Negative

- **Indirection** — More interfaces and abstractions
- **Package Structure** — Domain packages are deeper in tree
- **Migration Effort** — Existing code needed updates

### Risks

- **Over-Abstraction** — May have more indirection than needed
- **Interface Stability** — Core interfaces must remain stable
- **Discovery** — Finding where things live may be harder

### Mitigations

- Keep interfaces minimal
- Document ownership clearly
- Provide migration guides

---

## References

- [Glossary](../architecture/GLOSSARY.md)
- [Plugin Architecture](../architecture/plugin-architecture.md)
- [World Package](../architecture/world-architecture.md)
- [ADR-006 Plugin Architecture](./006-plugin-architecture.md)

---

## Related ADRs

- ADR-001: Runtime Architecture
- ADR-002: Behavior Model Design
- ADR-003: Memory Model Design
- ADR-004: Simulation Models Design
- ADR-005: Plugin Architecture Audit
- ADR-006: Plugin Architecture

---

## Milestone Traceability

| Milestone | Status | Notes |
|-----------|--------|-------|
| Core Architecture Freeze | ✅ Complete | This ADR |
| Domain Isolation | ✅ Complete | Electrical in plugins/ |
| Capability Framework | ✅ Complete | Capabilities in world/ |
| Template Framework | ✅ Complete | Templates in plugins/ |

---

*Accepted: 2026-07-13*
*Audit: KDSE Methodology*
