# ADR-006: Plugin Architecture

**ADR ID:** ADR-006  
**Title:** Forge Plugin Architecture  
**Date:** 2026-07-13  
**Status:** Accepted  
**Deciders:** Engineering Team  
**Repository:** https://github.com/tamzrod/forge

---

## Context

Forge must support multiple engineering domains (Electrical, Water, Building, Process, Transportation) without accumulating domain logic in Forge Core. The Component Registry architecture, while functional, has accumulated domain-specific leakage (electrical voltage validation in the generic registry). We need a permanent architecture that:

1. Keeps Forge Core domain-independent
2. Enables adding new domains without modifying Forge Core
3. Provides stable contracts for plugin development
4. Supports future dynamic plugin loading

---

## Decision

We adopt a **Plugin Architecture** with the following principles:

1. **Forge Core provides infrastructure only** — World, Clock, Topology, Solver Framework, Plugin Host, and Core Services
2. **Plugins provide domain knowledge** — Components, validators, solvers, scenarios, and other domain-specific assets
3. **Stable Plugin Contract** — A formal interface that all plugins must implement
4. **Core Services are metadata-only** — The Component Catalog stores descriptors, not domain logic

---

## Architecture

### Forge Core

```
┌─────────────────────────────────────────────────────────────────────┐
│                            Forge Core                                 │
│                                                                      │
│  ┌──────────┐  ┌─────────────┐  ┌─────────────┐  ┌────────────┐  │
│  │  World   │  │    Clock    │  │  Topology   │  │   Solver   │  │
│  └──────────┘  └─────────────┘  └─────────────┘  └────────────┘  │
│                                                                      │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │                     Plugin Host                              │   │
│  │  Discovery | Registration | Initialization | Shutdown        │   │
│  └──────────────────────────────────────────────────────────────┘   │
│                              │                                       │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │                   Core Services                               │   │
│  │  Component Catalog | Factory Registry | Event Bus            │   │
│  └──────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
```

### Plugins

```
┌───────────────┐     ┌───────────────┐     ┌───────────────┐
│   Electrical  │     │  Environment  │     │  Simulation   │
│    Plugin     │     │    Plugin     │     │    Plugin     │
├───────────────┤     ├───────────────┤     ├───────────────┤
│ Components    │     │ Components    │     │ Components    │
│ Entities      │     │ Entities      │     │ Scenarios     │
│ Solver        │     │               │     │               │
│ Validators    │     │               │     │               │
└───────────────┘     └───────────────┘     └───────────────┘
```

---

## Core Responsibilities

### What Forge Core OWNS

| Component | Responsibility |
|-----------|----------------|
| World | Entity container, tick coordination |
| Clock | Simulation time advancement |
| Topology | Spatial relationships |
| Solver Framework | State advancement algorithms |
| Plugin Host | Plugin lifecycle management |
| Component Catalog | Component metadata storage |
| Factory Registry | Entity factory storage |
| Event Bus | Event pub/sub |

### What Forge Core Does NOT OWN

| Component | Owner |
|-----------|-------|
| Component definitions | Domain Plugins |
| Entity implementations | Domain Plugins |
| Solvers | Domain Plugins |
| Connection validators | Domain Plugins |
| Domain-specific validation | Domain Plugins |
| Scenarios | Domain Plugins |

---

## Plugin Contract

Every plugin must implement the `Plugin` interface:

```go
type Plugin interface {
    // Identity
    ID() string                    // Unique plugin identifier
    Name() string                  // Human-readable name
    Version() string               // Semver version
    Description() string           // Plugin description
    Dependencies() []string        // Plugin dependencies
    
    // Lifecycle
    OnInit(ctx Context) error     // Initialize with Core Services
    OnShutdown() error             // Cleanup before unload
    
    // Contributions
    Components() []*ComponentDescriptor
    Categories() []*ComponentCategory
    Validators() []ConnectionValidator
    RegisterEntities(registry EntityRegistry)
}
```

---

## Core Services

Plugins consume Core Services through `PluginContext`:

```go
type Context interface {
    ComponentCatalog() ComponentCatalog
    FactoryRegistry() FactoryRegistry
    EventBus() EventBus
    World() World
    Logger() Logger
    Config() Config
}
```

### Component Catalog

```go
type ComponentCatalog interface {
    Register(desc *ComponentDescriptor) error
    Get(id string) *ComponentDescriptor
    List() []*ComponentDescriptor
    ListByCategory(categoryID string) []*ComponentDescriptor
    RegisterCategory(cat *ComponentCategory) error
    Categories() []*ComponentCategory
}
```

### Factory Registry

```go
type FactoryRegistry interface {
    Register(componentID string, factory ComponentFactory)
    Get(componentID string) ComponentFactory
    Create(componentID string, instance *ComponentInstance) (interface{}, error)
}
```

### Connection Validator

```go
type ConnectionValidator interface {
    Domain() string
    CanConnect(source, target *TerminalDescriptor) (bool, error)
}
```

---

## Plugin Lifecycle

```
Discovered → Registered → Initialized → Running → Shutdown → Unregistered
```

| State | Description |
|-------|-------------|
| Discovered | Plugin found (dynamic loading future) |
| Registered | Plugin added to Plugin Manager |
| Initialized | `OnInit()` called successfully |
| Running | Plugin ready to provide services |
| Shutdown | `OnShutdown()` called |
| Unregistered | Plugin removed from Plugin Manager |

---

## Plugin Communication Model

Plugins communicate only through **published Core interfaces**:

1. **Plugin → Core Services** — Plugins register components, factories, and validators
2. **Core Services → Plugins** — Core queries plugins for their contributions
3. **Plugin → Plugin** — Through Core Services (e.g., Event Bus for events)

**Prohibited:**
- Plugins accessing other plugin's internal state
- Plugins modifying Forge Core
- Plugins bypassing Core Services

---

## Future Dynamic Loading

The architecture supports future dynamic loading:

```go
type Manager interface {
    // Static registration (current)
    Register(p Plugin)
    
    // Dynamic loading (future)
    Load(paths []string) error
    Unload(id string) error
    
    // Query
    Get(id string) Plugin
    List() []Plugin
}
```

The Plugin Contract is designed to support both patterns without changes.

---

## Backwards Compatibility

### Existing Domain Packages

Existing packages (`forge-electrical`, `forge-environment`, `forge-simulation`) continue to work via static registration. The Plugin interface is additive—existing code does not need to change.

### Migration Path

1. **Phase 1:** Plugins register via `init()` (current pattern)
2. **Phase 2:** Plugins implement `Plugin` interface (gradual)
3. **Phase 3:** Plugin Manager supports both patterns
4. **Phase 4:** Dynamic loading (future, optional)

### Registry Transition

The existing `registry.Registry` continues to work as a **Component Catalog** implementation. The `CanConnect()` method now delegates to registered `ConnectionValidator` instances.

---

## Adding a New Domain

To add a new engineering domain (e.g., Water):

1. **Create plugin package** (`plugins/forge-water/`)
2. **Implement Plugin interface** with domain components
3. **Implement ConnectionValidator** for domain-specific rules
4. **Register with Plugin Manager** via `init()`
5. **No Forge Core changes required**

Example structure:

```
forge-water/
├── water.go              # Plugin implementation
├── components/           # Component definitions
├── validators/           # Connection validators
└── entities/            # Runtime entities
```

---

## Consequences

### Positive

- **Domain Independence** — Forge Core contains no engineering logic
- **Extensibility** — New domains without core changes
- **Stability** — Stable contracts for plugin development
- **Testability** — Plugins can be tested independently
- **Scalability** — Supports many engineering domains

### Negative

- **Abstraction Overhead** — Additional interfaces and indirection
- **Contract Stability** — Plugin contract must remain stable once published
- **Migration Effort** — Existing code needs gradual updates

### Risks

- **Over-engineering** — Plugin system may be excessive for few domains
- **Interface Instability** — Early iterations may need breaking changes
- **Plugin Loading Complexity** — Dynamic loading adds infrastructure

### Mitigations

- Start with minimal interface
- Allow `init()` fallback indefinitely
- Use semantic versioning for plugin API
- Defer dynamic loading until needed

---

## Implementation

### Code Structure

```
plugin/
├── plugin.go        # Plugin interface and Core Service interfaces
├── manager.go      # Plugin Manager implementation
├── services.go     # Default Core Service implementations

plugin/validators/
├── electrical.go   # Electrical connection validator

plugins/electrical/
├── electrical.go   # Electrical plugin implementation
```

### Files Modified

| File | Change |
|------|--------|
| `registry/registry.go` | Added `ConnectionValidator` interface, delegating `CanConnect()` |
| `docs/architecture/plugin-architecture.md` | Updated with complete architecture |
| `docs/architecture/GLOSSARY.md` | Added Plugin System terminology |

### New Files

| File | Purpose |
|------|---------|
| `plugin/plugin.go` | Plugin interface and Core Service interfaces |
| `plugin/manager.go` | Plugin Manager with static registration |
| `plugin/services.go` | Default Core Service implementations |
| `plugin/validators/electrical.go` | Electrical connection validator |
| `plugins/electrical/electrical.go` | Electrical plugin implementation |
| `docs/adrs/006-plugin-architecture.md` | This ADR |

---

## References

- [Plugin Architecture](../architecture/plugin-architecture.md)
- [Component Registry Architecture](../architecture/component-registry.md)
- [Glossary](../architecture/GLOSSARY.md)
- [ADR-005 Plugin Architecture Audit](./005-plugin-architecture-audit.md)
- [Runtime Architecture](./001-runtime-architecture.md)

---

## Related ADRs

- ADR-001: Runtime Architecture (references Plugin System)
- ADR-002: Behavior Model Design
- ADR-003: Memory Model Design
- ADR-004: Simulation Models Design
- ADR-005: Plugin Architecture Audit

---

## Milestone Traceability

| Milestone | Status | Notes |
|-----------|--------|-------|
| 1.4 Basic Plugin System | ✅ Complete | This ADR completes Milestone 1.4 |
| 2.0 Multi-Domain | ⏳ Pending | Requires plugin architecture |
| 2.0.1 Water Plugin | ⏳ Pending | First new-domain plugin |

---

*Accepted: 2026-07-13*
*Audit: KDSE Methodology*
