# Plugin Architecture

## Overview

Forge uses a Plugin Architecture to enable domain-independent core while supporting multiple engineering domains. Plugins contribute components, solvers, validators, scenarios, and other domain-specific assets.

**Core Principle:** Forge Core must never contain engineering-domain knowledge. All domain logic lives in plugins.

## Architecture

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
                              │
        ┌─────────────────────┼─────────────────────┐
        ▼                     ▼                     ▼
┌───────────────┐     ┌───────────────┐     ┌───────────────┐
│   Electrical  │     │  Environment  │     │  Simulation   │
│    Plugin     │     │    Plugin     │     │    Plugin     │
├───────────────┤     ├───────────────┤     ├───────────────┤
│ Components    │     │ Components    │     │ Components    │
│ Entities      │     │ Entities      │     │ Scenarios     │
│ Solver        │     │               │     │               │
│ Validators    │     │               │     │               │
│ Templates     │     │               │     │               │
└───────────────┘     └───────────────┘     └───────────────┘
```

## Plugin Contract

Every plugin must implement the Plugin interface:

```go
// Plugin is the contract that all Forge plugins must implement.
type Plugin interface {
    // ID returns the unique identifier for this plugin.
    ID() string
    
    // Name returns the human-readable name.
    Name() string
    
    // Version returns the plugin version (semver).
    Version() string
    
    // Description returns a description of the plugin.
    Description() string
    
    // Dependencies returns plugin IDs this plugin depends on.
    Dependencies() []string
    
    // OnInit initializes the plugin with access to Core Services.
    OnInit(ctx PluginContext) error
    
    // OnShutdown performs cleanup before the plugin is unloaded.
    OnShutdown() error
    
    // Components returns all component descriptors provided by this plugin.
    Components() []*ComponentDescriptor
    
    // Categories returns all component categories.
    Categories() []*ComponentCategory
    
    // Validators returns connection validators for this domain.
    Validators() []ConnectionValidator
    
    // RegisterEntities registers runtime entities with the world.
    RegisterEntities(registry EntityRegistry)
}
```

## Core Services

Plugins consume Core Services through the PluginContext:

```go
// PluginContext provides access to Core Services.
type PluginContext interface {
    // ComponentCatalog returns the component catalog service.
    ComponentCatalog() ComponentCatalog
    
    // FactoryRegistry returns the factory registry service.
    FactoryRegistry() FactoryRegistry
    
    // EventBus returns the event bus service.
    EventBus() EventBus
    
    // World returns the simulation world.
    World() World
    
    // Logger returns the plugin logger.
    Logger() Logger
    
    // Config returns plugin configuration.
    Config() PluginConfig
}

// ComponentCatalog stores component metadata.
type ComponentCatalog interface {
    Register(desc *ComponentDescriptor) error
    Get(id string) *ComponentDescriptor
    List() []*ComponentDescriptor
    ListByCategory(categoryID string) []*ComponentDescriptor
    Categories() []*ComponentCategory
}

// FactoryRegistry stores entity factories.
type FactoryRegistry interface {
    Register(componentID string, factory ComponentFactory)
    Create(componentID string, instance *ComponentInstance) (interface{}, error)
    Get(componentID string) ComponentFactory
}

// ConnectionValidator validates domain-specific connections.
type ConnectionValidator interface {
    Domain() string
    CanConnect(source, target *TerminalDescriptor) (bool, error)
}
```

## Plugin Lifecycle

```
┌────────────┐
│ Discovered │
└─────┬──────┘
      │ Register()
      ▼
┌────────────┐
│ Registered │
└─────┬──────┘
      │ Initialize all dependencies
      ▼
┌────────────┐
│Initialized │◄────── OnInit() called
└─────┬──────┘
      │ Ready to provide services
      ▼
┌────────────┐
│  Running   │
└─────┬──────┘
      │ OnShutdown()
      ▼
┌────────────┐
│  Shutdown  │
└─────┬──────┘
      │ Unregister()
      ▼
┌────────────┐
│Unregistered│
└────────────┘
```

## Electrical Plugin Example

```go
// Electrical plugin provides electrical power distribution components.
type ElectricalPlugin struct {
    ctx plugin.PluginContext
}

func (p *ElectricalPlugin) ID() string    { return "forge-electrical" }
func (p *ElectricalPlugin) Name() string  { return "Electrical Plugin" }
func (p *ElectricalPlugin) Version() string { return "1.0.0" }
func (p *ElectricalPlugin) Description() string { return "Electrical power distribution components" }

func (p *ElectricalPlugin) OnInit(ctx plugin.PluginContext) error {
    p.ctx = ctx
    
    // Register components with catalog
    catalog := ctx.ComponentCatalog()
    for _, comp := range p.Components() {
        if err := catalog.Register(comp); err != nil {
            return err
        }
    }
    
    // Register validators
    // ...
    
    return nil
}

func (p *ElectricalPlugin) Validators() []plugin.ConnectionValidator {
    return []plugin.ConnectionValidator{
        &VoltageValidator{},
        &ImpedanceValidator{},
    }
}
```

## Forge Core Boundaries

### What Forge Core OWNS

| Component | Purpose |
|-----------|---------|
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

## Adding a New Domain

To add a new engineering domain (e.g., Water):

1. **Create plugin package** (`plugins/forge-water/`)
2. **Implement Plugin interface** with domain components
3. **Register with Plugin Manager** via init()
4. **No Forge Core changes required**

Example water plugin structure:

```
forge-water/
├── water.go          # Plugin implementation
├── components/       # Component definitions
│   ├── pump.go
│   ├── valve.go
│   └── tank.go
├── entities/         # Runtime entities
│   ├── pump_entity.go
│   └── ...
├── solver/          # Domain solver
│   └── hydraulic_solver.go
└── validators/      # Connection validation
    └── flow_validator.go
```

## Configuration

```yaml
# forge.yaml
plugins:
  enabled:
    - forge-electrical
    - forge-environment
    - forge-simulation
  
  settings:
    forge-electrical:
      default_voltage: 480
```

## Backwards Compatibility

Existing domain packages (`forge-electrical`, `forge-environment`, `forge-simulation`) continue to work via static registration. The Plugin interface is additive—existing code does not need to change.

Migration path:
1. Plugins register via init() (current)
2. Plugins implement Plugin interface (gradual)
3. Plugin Manager supports both patterns (current)
4. Dynamic loading (future, optional)

## Glossary

See [GLOSSARY.md](GLOSSARY.md) for Plugin System terminology.

---

*Last Updated: 2026-07-13*
*Part of: ADR-006 Plugin Architecture*
