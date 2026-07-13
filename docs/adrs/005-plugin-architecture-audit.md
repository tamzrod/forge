# ADR-005: Plugin Architecture Audit

**ADR ID:** ADR-005  
**Title:** Component Registry Architecture Audit - Plugin Evolution Analysis  
**Date:** 2026-07-13  
**Status:** Accepted  
**Deciders:** Engineering Team  
**Repository:** https://github.com/tamzrod/forge

---

## Executive Summary

This audit evaluates whether Forge should evolve from its current Component Registry architecture toward a first-class Plugin Architecture. The audit examines existing knowledge, current implementation, responsibility boundaries, and provides evidence-based recommendations.

**Decision:** Evolve toward Plugin Architecture (evolutionary, not disruptive).

---

## Phase 1 — Knowledge Audit

### Concept Inventory

| Concept | Status | Evidence |
|---------|--------|----------|
| **Plugin** | Partially Defined | `docs/architecture/plugin-architecture.md` defines concept; Glossary defines ownership |
| **Plugin Manager** | Missing | No implementation; referenced in roadmap (Milestone 1.4) |
| **Domain Plugin** | Implicitly Defined | Domain plugins exist as packages (`forge-electrical`, `forge-environment`) but no formal interface |
| **Component Registry** | Fully Defined | `registry/registry.go`, `docs/architecture/component-registry.md` |
| **Component Descriptor** | Fully Defined | `registry/registry.go` (lines 64-76) |
| **Component Factory** | Fully Defined | `registry/registry.go` (lines 88-89) |
| **Domain Package** | Implicitly Defined | Packages exist but no formal interface |
| **Runtime Extension** | Missing | No formal mechanism defined |
| **Editor Extension** | Implicitly Defined | Registry-based, but no formal extension interface |

### Detailed Analysis

#### Plugin (Partially Defined)

**Definition Location:** `docs/architecture/plugin-architecture.md`

```go
type Plugin interface {
    ID() string
    DeviceTypes() []DeviceType
}
```

**Current State:** The concept exists in documentation but lacks implementation. The plugin architecture doc states: "Plugins provide device types. The runtime knows only the generic Device interface."

**Gap:** No actual `Plugin` interface in code; no plugin loader.

---

#### Plugin Manager (Missing)

**Evidence:** 
- Roadmap Milestone 1.4: "Basic Plugin System" lists this as pending
- Glossary references "Plugin Loader" but no implementation
- `docs/development/design-principles.md` shows Plugin Loader in architecture

**Gap:** The runtime does not have a plugin manager or loader. Components are registered via `init()` functions in domain packages.

---

#### Domain Plugin (Implicitly Defined)

**Evidence:**
- `registry/forge-electrical/` - Contains component definitions
- `registry/forge-environment/` - Contains component definitions  
- `registry/forge-simulation/` - Contains component definitions

**Current Pattern:**
```go
// In init.go
func init() {
    r := GetRegistry()
    for _, comp := range forgeelectrical.Components {
        r.Register(componentToDescriptor(comp))
    }
}
```

**Gap:** No formal domain plugin interface; components register themselves via package-level `init()` functions.

---

## Phase 2 — Architecture Audit

### Current Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         Forge Core                          │
│  ┌─────────┐  ┌──────────┐  ┌─────────┐  ┌──────────┐     │
│  │  World  │  │  Solver  │  │ Clock   │  │Topology  │     │
│  └─────────┘  └──────────┘  └─────────┘  └──────────┘     │
│                                                              │
│  ┌─────────────────────────────────────────────────────┐    │
│  │              Component Registry                       │    │
│  │  Components | Categories | Factories | Validation   │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
                           │
        ┌──────────────────┼──────────────────┐
        ▼                  ▼                  ▼
┌───────────────┐  ┌───────────────┐  ┌───────────────┐
│ forge-electrical│ │forge-environment│ │forge-simulation│
│  Components    │  │  Components    │  │  Components    │
│  (7 types)     │  │  (3 types)     │  │  (2 types)     │
└───────────────┘  └───────────────┘  └───────────────┘
```

### What Each Domain Currently Contributes

#### Electrical Domain (`registry/forge-electrical/`)
- **Component Definitions:** Grid, Bus, Breaker, Transformer, Generator, Load, Meter (7 total)
- **Domain Logic:** None in registry package
- **Runtime Entities:** `world/electrical/entities.go` (6 entity types)
- **Solver:** `solver/electrical.go` (ElectricalSolver)
- **Topology:** `topology/electrical.go` (Network, Bus, Branch, Terminal)

#### Environment Domain (`registry/forge-environment/`)
- **Component Definitions:** Sun, Weather, Wind (3 total)
- **Domain Logic:** None in registry package
- **Runtime Entities:** `world/environmental/entities.go` (Sun, Weather, PVArray - 3 entity types)

#### Simulation Domain (`registry/forge-simulation/`)
- **Component Definitions:** Scenario, Clock (2 total)
- **Domain Logic:** Scenario behavior in `scenarios/scenarios.go`

### Registry Responsibilities Analysis

The Component Registry currently accumulates:

| Responsibility | Belongs in Registry? | Evidence |
|---------------|---------------------|----------|
| Component Registration | ✅ Yes | Core purpose |
| Category Management | ✅ Yes | Core purpose |
| Factory Registration | ✅ Yes | Core purpose |
| Property Descriptors | ✅ Yes | Editor metadata |
| Terminal Descriptors | ✅ Yes | Connection metadata |
| Connection Validation | ⚠️ Partial | Only basic electrical validation |
| Domain-Specific Validation | ❌ No | Electrical-specific logic in generic registry |

### Future Domain Contributions Analysis

If future domains (Water, Building, Process, Transportation) follow the same pattern:

| Domain | Components | Entities | Solvers | Validators | Templates | Docs |
|--------|-----------|----------|---------|------------|-----------|------|
| Electrical | 7 | 6 | 1 | ~5 | Future | Existing |
| Environment | 3 | 3 | 0 | 0 | Future | Existing |
| Simulation | 2 | 0 | 0 | 0 | ~3 | Existing |
| Water | ~5 | ~5 | 1 | 2 | Future | New |
| Building | ~8 | ~8 | 1 | 3 | Future | New |
| Process | ~6 | ~6 | 1 | 2 | Future | New |
| Transportation | ~4 | ~4 | 1 | 2 | Future | New |

**Observation:** Each domain contributes approximately:
- 5-8 component definitions
- 3-8 entity types  
- 1 solver
- 2-5 validators
- Domain-specific documentation

**Problem:** The registry `CanConnect()` function contains electrical-specific validation logic that won't scale to other domains.

---

## Phase 3 — Responsibility Analysis

### Current Responsibility Map

| Component | Owns | Should Own |
|-----------|------|-----------|
| **Forge Core** | World, Clock, Topology, Solver interfaces | World, Clock, Topology, Solver interfaces |
| **Component Registry** | Component descriptors, Categories, Factories, Basic validation | Component descriptors, Categories, Factory registration |
| **forge-electrical** | Component definitions, Entities, ElectricalSolver, Topology | Component definitions, Electrical entities, Electrical solver |
| **forge-environment** | Component definitions, Entities | Component definitions, Environment entities |
| **forge-simulation** | Component definitions, Scenarios | Component definitions, Scenario behavior |

### Responsibility Leakage Identified

#### Leakage #1: Registry Contains Domain Logic

**Location:** `registry/registry.go` lines 289-307

```go
func (r *Registry) CanConnect(source, target *TerminalDescriptor) bool {
    // Bus can connect to most things
    if source.Role == TerminalRoleThrough || target.Role == TerminalRoleThrough {
        return true
    }
    // Observation terminals can connect anywhere
    if source.Role == TerminalRoleObservation || target.Role == TerminalRoleObservation {
        return true
    }
    // Voltage must match
    if source.Voltage != nil && target.Voltage != nil {
        return *source.Voltage == *target.Voltage
    }
    return true
}
```

**Problem:** Voltage matching is electrical-specific; other domains (water flow, temperature) need different validation.

#### Leakage #2: Registry Imports Domain Packages

**Location:** `registry/init.go`

```go
import (
    "github.com/tamzrod/forge/registry/forge-electrical"
    "github.com/tamzrod/forge/registry/forge-environment"
    "github.com/tamzrod/forge/registry/forge-simulation"
)
```

**Problem:** Registry depends on specific domain packages; adding new domains requires modifying registry.

#### Leakage #3: Domain Hardcoded in Component Descriptor

**Location:** `registry/init.go` line 87

```go
Domain: fmt.Sprintf("forge-%s", comp.Category),
```

**Problem:** Domain is derived from category, not explicitly declared.

### Single Owner Analysis

| Responsibility | Current Owner | Should Be Owner | Leakage |
|---------------|---------------|-----------------|---------|
| Electrical topology | `topology/electrical.go` | Electrical Plugin | No |
| Electrical solver | `solver/electrical.go` | Electrical Plugin | No |
| Electrical entities | `world/electrical/` | Electrical Plugin | No |
| Connection validation | `registry/` | Domain Plugin | Yes |
| Component registration | `registry/` | Registry (correct) | No |
| Factory registration | `registry/` | Registry (correct) | No |

---

## Phase 4 — Evolution Proposal

### Justification for Plugin Architecture

**Evidence:**
1. Documentation already defines plugin concept (`docs/architecture/plugin-architecture.md`)
2. Roadmap includes Plugin System (Milestone 1.4)
3. Multiple domains are planned (Water, Building, Manufacturing)
4. Current registry has domain-specific leakage
5. ADR-001 Milestone Traceability explicitly mentions "Plugin System" as pending

### Proposed Target Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         Forge Core                          │
│  ┌─────────┐  ┌──────────┐  ┌─────────┐  ┌──────────┐     │
│  │  World  │  │  Solver  │  │ Clock   │  │Topology  │     │
│  └─────────┘  └──────────┘  └─────────┘  └──────────┘     │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              Plugin Manager                           │   │
│  │  load() | unload() | get() | list() | validate()   │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                           │
        ┌──────────────────┼──────────────────┐
        ▼                  ▼                  ▼
┌───────────────┐  ┌───────────────┐  ┌───────────────┐
│Electrical Plugin│ │Environment Plug│ │Simulation Plug│
│  Components    │  │  Components    │  │  Components    │
│  Entities      │  │  Entities      │  │  Entities      │
│  Solver        │  │                │  │  Scenarios     │
│  Validators    │  │                │  │                │
│  Icons         │  │                │  │                │
│  Templates     │  │                │  │                │
│  Docs          │  │                │  │                │
└───────────────┘  └───────────────┘  └───────────────┘
```

### Plugin Interface

```go
// Plugin is the interface for domain plugins.
type Plugin interface {
    // ID returns the plugin identifier.
    ID() string
    
    // Name returns the human-readable name.
    Name() string
    
    // Version returns the plugin version.
    Version() string
    
    // Components returns all component descriptors.
    Components() []*ComponentDescriptor
    
    // Categories returns all categories provided by this plugin.
    Categories() []*ComponentCategory
    
    // RegisterEntities registers runtime entities with the world.
    RegisterEntities(registry EntityRegistry)
    
    // Validators returns connection validators for this domain.
    Validators() []ConnectionValidator
}

// ConnectionValidator validates connections within a domain.
type ConnectionValidator interface {
    // CanConnect validates a connection between two terminals.
    CanConnect(source, target *TerminalDescriptor) (bool, error)
}
```

### Backwards Compatibility

**Preserve:**
1. Current `registry.Registry` continues to work
2. Current `init()` registration pattern remains valid
3. All existing components remain functional
4. Editor continues to query registry as before

**Evolution Path:**
1. Add Plugin interface (non-breaking)
2. Add Plugin Manager (non-breaking)
3. Existing plugins continue via `init()`
4. New plugins use Plugin interface
5. Registry gains ability to delegate to plugins

### Migration Strategy

#### Phase 1: Define Interface (Non-Breaking)
- Add `Plugin` interface to `registry/plugin.go`
- Add `PluginManager` interface
- No changes to existing code

#### Phase 2: Implement Manager (Non-Breaking)
- Add `DefaultPluginManager` implementation
- Existing `init()` functions auto-register as plugins
- Registry queries manager for components

#### Phase 3: Migrate Domains (Optional, Gradual)
- Migrate electrical domain to implement Plugin interface
- Add domain-specific validators
- Remove electrical-specific logic from registry

#### Phase 4: Deprecate init() Pattern (Future)
- Mark direct registry access as deprecated
- Encourage Plugin interface usage
- Full plugin architecture complete

---

## Phase 5 — Decision

### Recommendation: Evolve Toward Plugin Architecture

**Rationale:**

1. **Documentation Already Commits to Plugins**
   - `docs/architecture/plugin-architecture.md` exists and defines the concept
   - Glossary defines Plugin ownership
   - Design principles state "Plugins Provide Domain Knowledge"

2. **Roadmap Already Plans for Plugins**
   - Milestone 1.4: "Basic Plugin System" is explicitly planned
   - Version 2.0: "Multi-Domain" explicitly mentions "Water Plugin"

3. **Current Architecture Has Leakage**
   - Registry contains electrical-specific validation
   - Registry imports specific domain packages
   - Hardcoded domain derivation

4. **Scale Requires Formal Boundaries**
   - 7 future domains planned (Water, Building, Manufacturing, Transportation, etc.)
   - Each domain contributes solvers, validators, entities
   - Ad-hoc registration won't scale

5. **Minimal Disruption Path Exists**
   - Plugin interface can be additive
   - Existing `init()` pattern can coexist
   - No breaking changes required

### Conditions Not Met for Immediate Large-Scale Refactor

1. **Only one domain with solver** - Electrical is the only domain with a solver; other domains don't need this yet
2. **No external plugin ecosystem** - No third-party plugins to worry about
3. **Documentation already exists** - The design is defined; implementation is missing
4. **Working architecture** - Current registry works for current needs

### Final Recommendation

| Action | Timeline | Priority |
|--------|----------|----------|
| Define `Plugin` interface | Immediate | High |
| Implement `PluginManager` | Near-term | High |
| Migrate electrical validation to plugin | When Water domain starts | Medium |
| Deprecate `init()` registration | Post-2.0 | Low |

---

## Consequences

### Positive
- Clear domain boundaries
- Scalable to multiple domains
- Aligns implementation with existing documentation
- Enables third-party plugins in future
- Removes domain leakage from registry

### Negative
- Additional abstraction layer
- Plugin interface must be stable once defined
- Migration effort for existing code

### Risks
- Over-engineering if only 2-3 domains ever exist
- Interface instability during early iterations
- Plugin loading complexity

### Mitigations
- Start with minimal interface
- Allow `init()` fallback indefinitely
- Use semantic versioning for plugin API

---

## References

- [Plugin Architecture (Existing Doc)](docs/architecture/plugin-architecture.md)
- [Component Registry Architecture](docs/architecture/component-registry.md)
- [Glossary](docs/architecture/GLOSSARY.md)
- [Roadmap](docs/development/roadmap.md)
- [ADR-001 Runtime Architecture](docs/adrs/001-runtime-architecture.md)

---

## Related ADRs

- ADR-001: Runtime Architecture (references Plugin System)
- ADR-002: Behavior Model Design
- ADR-003: Memory Model Design
- ADR-004: Simulation Models Design

---

## Milestone Traceability

| Milestone | Status | Notes |
|-----------|--------|-------|
| 1.4 Basic Plugin System | ⏳ Pending | This ADR enables Milestone 1.4 |
| 2.0 Multi-Domain | ⏳ Pending | Requires plugin architecture |
| 2.0.1 Water Plugin | ⏳ Pending | First new-domain plugin |

---

*Audit conducted: 2026-07-13*
*Auditor: KDSE Methodology*
