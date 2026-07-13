# Gap Analysis: Component Registry → Plugin Architecture

**Date:** 2026-07-13  
**Purpose:** Detailed gap analysis between current implementation and proposed plugin architecture

---

## Gap Summary

| Gap | Severity | Location | Fix Complexity |
|-----|----------|----------|----------------|
| Missing Plugin Interface | High | N/A | Low |
| Missing Plugin Manager | High | N/A | Medium |
| Domain Logic in Registry | Medium | `registry/registry.go:289-307` | Medium |
| init() Pattern Only | Medium | `registry/init.go` | Low |
| No Domain Validators | Medium | N/A | Medium |
| No Plugin Configuration | Low | N/A | Low |

---

## Gap 1: Missing Plugin Interface

### Current State
No `Plugin` interface exists in the codebase.

### Desired State
```go
type Plugin interface {
    ID() string
    Name() string
    Version() string
    Components() []*ComponentDescriptor
    Categories() []*ComponentCategory
    RegisterEntities(registry EntityRegistry)
    Validators() []ConnectionValidator
}
```

### Gap Impact
- No standard way to define domain plugins
- Plugin loading impossible without interface
- External plugins cannot be developed

### Fix
Create `registry/plugin.go` with interface definitions.

---

## Gap 2: Missing Plugin Manager

### Current State
Components register themselves via `init()` functions in domain packages.

### Desired State
```go
type PluginManager interface {
    Load(paths []string) error
    Unload(id string) error
    Get(id string) Plugin
    List() []Plugin
    Validate() []error
}

var manager = NewDefaultPluginManager()

func init() {
    // Existing packages auto-register
    manager.AutoRegister()
}
```

### Gap Impact
- No dynamic plugin loading
- All plugins must be imported at compile time
- No way to disable domains at runtime

### Fix
Create `registry/manager.go` with manager implementation.

---

## Gap 3: Domain Logic in Registry

### Current State
`registry/registry.go` line 289-307 contains electrical-specific validation:

```go
func (r *Registry) CanConnect(source, target *TerminalDescriptor) bool {
    if source.Role == TerminalRoleThrough || target.Role == TerminalRoleThrough {
        return true
    }
    if source.Role == TerminalRoleObservation || target.Role == TerminalRoleObservation {
        return true
    }
    if source.Voltage != nil && target.Voltage != nil {
        return *source.Voltage == *target.Voltage
    }
    return true
}
```

### Problems
1. **Voltage matching** is electrical-specific
2. Water domain needs flow-rate matching
3. Temperature domain needs thermal compatibility
4. Registry shouldn't know about domain physics

### Desired State
```go
type ConnectionValidator interface {
    CanConnect(source, target *TerminalDescriptor) (bool, error)
}

type Registry struct {
    // ... existing fields ...
    validators map[string][]ConnectionValidator
}

func (r *Registry) CanConnect(source, target *TerminalDescriptor) bool {
    for _, validator := range r.validators[source.Domain] {
        if valid, _ := validator.CanConnect(source, target); !valid {
            return false
        }
    }
    return true
}
```

### Gap Impact
- Cannot properly support multi-domain simulations
- Validation logic will accumulate domain-specific code
- Registry becomes coupled to electrical domain

### Fix
1. Add `ConnectionValidator` interface
2. Move electrical validation to electrical plugin
3. Registry delegates to registered validators

---

## Gap 4: init() Pattern Only

### Current State
Domain packages register via `init()`:
```go
// registry/init.go
func init() {
    r := GetRegistry()
    for _, comp := range forgeelectrical.Components {
        r.Register(componentToDescriptor(comp))
    }
}
```

### Desired State
Both `init()` and Plugin interface supported:
```go
// Backwards compatible
func init() {
    registry.RegisterPlugin(&MyElectricalPlugin{})
}

// New plugin interface
type MyElectricalPlugin struct{}

func (p *MyElectricalPlugin) ID() string { return "forge-electrical" }
func (p *MyElectricalPlugin) Components() []*ComponentDescriptor {
    return electrical.Components
}
```

### Gap Impact
- No migration path from `init()` to plugins
- External plugins impossible
- All domains must be imported at compile time

### Fix
1. Add Plugin interface
2. Implement manager that wraps existing `init()`
3. Deprecate direct registry access (future)

---

## Gap 5: No Domain Validators

### Current State
Only one validation method exists in registry.

### Desired State
Each domain provides validators:
```go
// Electrical plugin
func (p *ElectricalPlugin) Validators() []ConnectionValidator {
    return []ConnectionValidator{
        &VoltageValidator{},
        &ImpedanceValidator{},
    }
}

// Water plugin
func (p *WaterPlugin) Validators() []ConnectionValidator {
    return []ConnectionValidator{
        &FlowRateValidator{},
        &PressureValidator{},
    }
}
```

### Gap Impact
- Cannot add domain-specific validation
- Registry accumulates electrical-specific logic
- Multi-domain validation impossible

### Fix
1. Define `ConnectionValidator` interface
2. Each plugin implements validators
3. Registry aggregates validators by domain

---

## Gap 6: No Plugin Configuration

### Current State
All plugins loaded at compile time via imports.

### Desired State
```yaml
# forge.yaml
plugins:
  enabled:
    - forge-electrical
    - forge-environment
    - forge-simulation
  paths:
    - ./plugins/custom
```

### Gap Impact
- Cannot disable domains
- No plugin discovery mechanism
- External plugins cannot be loaded

### Fix
1. Add plugin configuration to `runtime.Config`
2. Implement plugin discovery from paths
3. Support enable/disable per domain

---

## Implementation Priority

### P0 - Critical (Enable Milestone 1.4)
1. Plugin interface definition
2. PluginManager interface
3. Backwards-compatible auto-registration

### P1 - Important (Enable Multi-Domain)
4. ConnectionValidator interface
5. Move electrical validation to plugin
6. Registry delegates to validators

### P2 - Nice to Have (Polish)
7. Plugin configuration
8. Plugin discovery
9. External plugin loading

---

## Risk Assessment

| Gap | Risk if Not Fixed | Risk if Fixed |
|-----|------------------|---------------|
| Missing Plugin Interface | Blocks plugin ecosystem | Interface instability |
| Missing Plugin Manager | Inflexible architecture | Complexity overhead |
| Domain Logic in Registry | Registry bloat | Migration effort |
| init() Pattern Only | Limited extensibility | Backwards compatibility |
| No Domain Validators | Poor multi-domain support | Validator coordination |
| No Plugin Configuration | Inflexible loading | Configuration complexity |

---

## Dependencies

```
Gap 1 (Plugin Interface)
    ↓
Gap 2 (Plugin Manager)
    ↓
Gap 4 (init() Pattern) ← Already handled by Gap 2
    ↓
Gap 5 (Domain Validators)
    ↓
Gap 3 (Domain Logic in Registry) ← Solved by Gap 5
    ↓
Gap 6 (Plugin Configuration)
```

---

*Gap Analysis Date: 2026-07-13*
*Part of: ADR-005 Plugin Architecture Audit*
