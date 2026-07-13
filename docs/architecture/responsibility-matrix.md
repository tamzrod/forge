# Responsibility Matrix: Component Registry vs Plugin Architecture

**Date:** 2026-07-13  
**Purpose:** Define clear ownership boundaries for architectural components

---

## Current State (Component Registry)

| Component | Responsibility | Owner | Leakage |
|-----------|----------------|-------|---------|

### Forge Core
| Component | Responsibility | Owner | Leakage |
|-----------|----------------|-------|---------|
| World | Entity container, tick coordination | `world/world.go` | None |
| Clock | Simulation time advancement | `simulation/clock.go` | None |
| Topology | Spatial relationships | `topology/` | None |
| Solver | State advancement algorithms | `solver/` | None |

### Component Registry
| Component | Responsibility | Owner | Leakage |
|-----------|----------------|-------|---------|
| Component Registration | Accept and store component descriptors | `registry/registry.go` | None |
| Category Management | Organize components into categories | `registry/registry.go` | None |
| Factory Registration | Store component factories | `registry/registry.go` | None |
| Connection Validation | Validate terminal connections | `registry/registry.go` | ⚠️ Electrical-specific |
| Domain Derivation | Extract domain from category | `registry/init.go` | ⚠️ Hardcoded |

### Domain Packages
| Component | Responsibility | Owner | Leakage |
|-----------|----------------|-------|---------|
| Electrical Components | Define electrical component descriptors | `forge-electrical/` | None |
| Electrical Entities | Runtime electrical simulation | `world/electrical/` | None |
| Electrical Solver | Power balance calculations | `solver/electrical.go` | None |
| Electrical Topology | Electrical network modeling | `topology/electrical.go` | None |
| Environment Components | Define environment component descriptors | `forge-environment/` | None |
| Environment Entities | Runtime environment simulation | `world/environmental/` | None |
| Simulation Components | Define simulation component descriptors | `forge-simulation/` | None |
| Simulation Scenarios | Scenario definitions | `scenarios/` | None |

---

## Target State (Plugin Architecture)

### Forge Core
| Component | Responsibility | Owner | Leakage |
|-----------|----------------|-------|---------|
| World | Entity container, tick coordination | `world/world.go` | None |
| Clock | Simulation time advancement | `simulation/clock.go` | None |
| Topology | Spatial relationships | `topology/` | None |
| Solver | State advancement algorithms | `solver/` | None |
| Plugin Manager | Plugin lifecycle management | `registry/manager.go` | None |

### Plugin Manager
| Component | Responsibility | Owner | Leakage |
|-----------|----------------|-------|---------|
| Plugin Loading | Load plugins from paths | `registry/manager.go` | None |
| Plugin Registry | Track loaded plugins | `registry/manager.go` | None |
| Component Aggregation | Combine components from plugins | `registry/manager.go` | None |
| Validator Aggregation | Combine validators from plugins | `registry/manager.go` | None |

### Plugins
| Component | Responsibility | Owner | Leakage |
|-----------|----------------|-------|---------|
| Plugin Interface | Define plugin contract | `registry/plugin.go` | None |
| Electrical Plugin | All electrical domain concerns | `forge-electrical/` | None |
| Environment Plugin | All environment domain concerns | `forge-environment/` | None |
| Simulation Plugin | All simulation domain concerns | `forge-simulation/` | None |

### Electrical Plugin
| Component | Responsibility | Owner | Leakage |
|-----------|----------------|-------|---------|
| Components | Electrical component descriptors | `forge-electrical/` | None |
| Entities | Electrical runtime entities | `world/electrical/` | None |
| Solver | Electrical power balance | `solver/electrical.go` | None |
| Topology | Electrical network | `topology/electrical.go` | None |
| Validators | Electrical connection rules | `forge-electrical/` | None |

---

## Responsibility Transitions

### Registry → Plugin Manager

| Responsibility | Current Owner | Target Owner | Migration |
|---------------|---------------|--------------|-----------|
| Plugin lifecycle | N/A (doesn't exist) | Plugin Manager | Create |
| Component aggregation | Registry | Plugin Manager | Delegate |
| Validator aggregation | N/A (doesn't exist) | Plugin Manager | Create |

### Registry → Domain Plugin

| Responsibility | Current Owner | Target Owner | Migration |
|---------------|---------------|--------------|-----------|
| Connection validation | Registry | Domain Plugin | Move |
| Domain-specific types | Registry | Domain Plugin | Move |
| Domain icons/docs | N/A | Domain Plugin | Add |

### Registry (Stays)

| Responsibility | Owner | Justification |
|---------------|-------|---------------|
| Component storage | Registry | Single source of truth |
| Category storage | Registry | Editor needs unified view |
| Factory storage | Registry | Runtime entity creation |
| Palette generation | Registry | Editor integration point |

---

## Boundary Definition

### What Registry OWNS (Core)
- Component descriptor storage
- Category storage
- Factory storage
- Palette item generation
- Inspector property generation

### What Registry DELEGATES (To Plugins)
- Connection validation
- Domain-specific properties
- Terminal compatibility rules
- Icon associations
- Documentation URLs

### What Plugin OWNS
- Component definitions
- Entity implementations
- Solver implementations
- Domain validators
- Domain icons
- Domain documentation

### What Plugin DELEGATES (To Core)
- Component registration (to Registry)
- Entity lifecycle (to World)
- Solver coordination (to World)
- Time advancement (to Clock)

---

## Single Owner Verification

### Connection Validation

**Current State:**
```
Registry.CanConnect() ← Electrical-specific logic
```

**Problem:** Voltage matching is electrical; water needs pressure matching.

**Target State:**
```
Registry.CanConnect()
    ↓
For each domain's validators:
    ↓
DomainPlugin.Validator.CanConnect()
```

**Verification:** Each domain owns its validation rules.

### Entity Creation

**Current State:**
```
Registry.CreateFromFactory() ← Generic factory
    ↓
ComponentFactory() ← Returns interface{}
```

**Target State:**
```
Registry.CreateFromFactory()
    ↓
Plugin.CreateEntity() ← Domain-specific creation
```

**Verification:** Each plugin owns entity creation.

### Solver Selection

**Current State:**
```
World.SetSolver() ← Accepts any Solver
```

**Target State:**
```
World.SetSolver() ← Accepts domain solver
    ↓
DomainPlugin.RegisterSolver() ← Registers with World
```

**Verification:** Each plugin registers its solver.

---

## Interface Contracts

### Registry Interface (Core)

```go
type Registry interface {
    // Storage (Core)
    Register(desc *ComponentDescriptor) error
    Get(id string) *ComponentDescriptor
    List() []*ComponentDescriptor
    ListByCategory(categoryID string) []*ComponentDescriptor
    
    // Storage (Core)
    RegisterCategory(cat *ComponentCategory) error
    Categories() []*ComponentCategory
    
    // Storage (Core)
    RegisterFactory(componentID string, factory ComponentFactory)
    CreateFromFactory(instance *ComponentInstance) (interface{}, error)
    
    // Delegation (To Plugins)
    CanConnect(source, target *TerminalDescriptor) bool
    
    // Delegation (To Plugins)
    GetPaletteItems() []PaletteItem
}
```

### Plugin Interface (Domain)

```go
type Plugin interface {
    // Identity
    ID() string
    Name() string
    Version() string
    
    // Content
    Components() []*ComponentDescriptor
    Categories() []*ComponentCategory
    
    // Runtime Integration
    RegisterEntities(world World)
    RegisterSolver(solver Solver)
    
    // Delegation
    Validators() []ConnectionValidator
}
```

### ConnectionValidator Interface (Domain)

```go
type ConnectionValidator interface {
    Domain() string
    CanConnect(source, target *TerminalDescriptor) (bool, error)
}
```

---

## Migration Checklist

- [ ] Define `Plugin` interface
- [ ] Define `PluginManager` interface
- [ ] Define `ConnectionValidator` interface
- [ ] Implement `DefaultPluginManager`
- [ ] Implement electrical validator as plugin
- [ ] Remove electrical logic from registry
- [ ] Update `init()` to use plugin pattern
- [ ] Add plugin configuration support
- [ ] Update documentation

---

*Responsibility Matrix Date: 2026-07-13*
*Part of: ADR-005 Plugin Architecture Audit*
