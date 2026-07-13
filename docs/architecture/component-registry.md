# Component Registry Architecture

## Purpose

The Component Registry is the authoritative source for all engineering components in Forge. It enables a plugin-based architecture where new engineering domains can be added without modifying the core Editor.

## Design Principle

**Registry as Single Source of Truth**

```
┌─────────────────────────────────────────────────────────────┐
│                    COMPONENT REGISTRY                        │
│  Central authority for all engineering components            │
└─────────────────────────────────────────────────────────────┘
                           │
           ┌───────────────┼───────────────┐
           │               │               │
           ▼               ▼               ▼
      ┌─────────┐    ┌──────────┐   ┌──────────┐
      │Palette  │    │Inspector │   │ Canvas   │
      │         │    │          │   │          │
      │Queries  │    │Queries   │   │Queries   │
      │Registry │    │Registry  │   │Registry  │
      └─────────┘    └──────────┘   └──────────┘
           │               │               │
           └───────────────┼───────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                    RUNTIME                                  │
│  Factory creates entities from descriptors                   │
└─────────────────────────────────────────────────────────────┘
```

## Core Concepts

### 1. Component Registry

The global registry that holds all registered components.

```go
type ComponentRegistry struct {
    components map[string]*ComponentDescriptor
    categories map[string]*ComponentCategory
    mu        sync.RWMutex
}

// Global singleton
var registry = NewComponentRegistry()

// Register adds a component to the registry
func (r *ComponentRegistry) Register(desc *ComponentDescriptor) error

// Get retrieves a component by ID
func (r *ComponentRegistry) Get(id string) *ComponentDescriptor

// List returns all registered components
func (r *ComponentRegistry) List() []*ComponentDescriptor

// ListByCategory returns components in a category
func (r *ComponentRegistry) ListByCategory(categoryID string) []*ComponentDescriptor

// Categories returns all categories
func (r *ComponentRegistry) Categories() []*ComponentCategory
```

### 2. Component Descriptor

Describes a component's editing and runtime properties.

```go
type ComponentDescriptor struct {
    ID          string                 // Unique identifier (e.g., "forge-electrical:grid")
    Name        string                 // Display name (e.g., "Utility Grid")
    Category    string                 // Category ID (e.g., "electrical")
    Icon        string                 // Emoji or icon identifier (e.g., "🔌")
    Description string                 // Human-readable description
    
    // Editing properties
    Properties  []PropertyDescriptor   // Editable properties
    Terminals   []TerminalDescriptor    // Connection points
    
    // Canvas properties
    Width       float64                // Default canvas width
    Height      float64                // Default canvas height
    
    // Runtime factory
    Factory     ComponentFactory        // Creates runtime entity
    
    // Capabilities
    Capabilities []string               // What the component can do
    
    // Domain identifier
    Domain      string                 // e.g., "forge-electrical"
}
```

### 3. Property Descriptor

Defines an editable property.

```go
type PropertyDescriptor struct {
    Key         string           // Property key (e.g., "nominal_voltage")
    Label       string           // Display label (e.g., "Nominal Voltage")
    Type        PropertyType     // string, number, boolean, enum
    Default     interface{}      // Default value
    Unit        string           // Optional unit (e.g., "V", "kW")
    Min         *float64         // Optional minimum
    Max         *float64         // Optional maximum
    Options     []string         // For enum type
    ReadOnly    bool             // Read-only property
    Required    bool             // Required property
}

type PropertyType string

const (
    PropertyTypeString  PropertyType = "string"
    PropertyTypeNumber  PropertyType = "number"
    PropertyTypeBoolean PropertyType = "boolean"
    PropertyTypeEnum    PropertyType = "enum"
)
```

### 4. Terminal Descriptor

Defines a connection point.

```go
type TerminalDescriptor struct {
    ID       string       // Terminal ID (e.g., "output")
    Name     string       // Display name (e.g., "Output")
    Role     TerminalRole // source, destination, through, observation
    Voltage  *float64     // Optional voltage level
    Direction TerminalDirection // input, output, bidirectional
}

type TerminalRole string

const (
    TerminalRoleSource       TerminalRole = "source"       // Injects power
    TerminalRoleDestination  TerminalRole = "destination"  // Withdraws power
    TerminalRoleThrough      TerminalRole = "through"     // Passes through
    TerminalRoleObservation  TerminalRole = "observation" // Measures
)

type TerminalDirection string

const (
    TerminalDirectionInput  TerminalDirection = "input"
    TerminalDirectionOutput TerminalDirection = "output"
    TerminalDirectionBoth  TerminalDirection = "bidirectional"
)
```

### 5. Component Category

Organizes components into groups.

```go
type ComponentCategory struct {
    ID        string // Category ID (e.g., "electrical")
    Name      string // Display name (e.g., "Electrical")
    Icon      string // Category icon (e.g., "⚡")
    Order     int    // Display order
    Domain    string // Owning domain (e.g., "forge-electrical")
    Expandable bool  // Can be collapsed
}
```

### 6. Component Factory

Creates runtime entities from component instances.

```go
type ComponentFactory func(instance *ComponentInstance) (Entity, error)

type ComponentInstance struct {
    ID         string
    ComponentID string
    Properties map[string]interface{}
    Position   Point
}
```

### 7. Component Instance

An actual placed component on the canvas.

```go
type ComponentInstance struct {
    ID          string
    ComponentID string           // Reference to descriptor
    Name        string           // User-assigned name
    Position    Point            // Canvas position
    Properties  map[string]interface{}  // Property values
    Connections []Connection     // Terminal connections
}
```

## Registration Pattern

### Registering a New Component

```go
// In forge-electrical plugin
func init() {
    registry := GetRegistry()
    
    registry.Register(&ComponentDescriptor{
        ID:          "forge-electrical:grid",
        Name:        "Utility Grid",
        Category:    "electrical",
        Icon:        "🔌",
        Description: "Utility grid connection point",
        Properties: []PropertyDescriptor{
            {
                Key:     "nominal_voltage",
                Label:   "Nominal Voltage",
                Type:    PropertyTypeNumber,
                Default: float64(69000),
                Unit:    "V",
                Min:     floatPtr(1),
            },
            {
                Key:     "nominal_frequency",
                Label:   "Frequency",
                Type:    PropertyTypeNumber,
                Default: float64(60),
                Unit:    "Hz",
                Options: []string{"50", "60"},
            },
        },
        Terminals: []TerminalDescriptor{
            {
                ID:        "output",
                Name:      "Output",
                Role:      TerminalRoleSource,
                Direction: TerminalDirectionOutput,
            },
        },
        Width:   80,
        Height:  60,
        Domain:  "forge-electrical",
        Factory: createGridEntity,
    })
    
    // Register category if not exists
    registry.RegisterCategory(&ComponentCategory{
        ID:       "electrical",
        Name:     "Electrical",
        Icon:     "⚡",
        Order:    1,
        Domain:   "forge-electrical",
    })
}
```

### Using the Registry

```go
// In Palette component
func GetPaletteItems() []PaletteItem {
    registry := GetRegistry()
    items := make([]PaletteItem, 0)
    
    for _, category := range registry.Categories() {
        for _, component := range registry.ListByCategory(category.ID) {
            items = append(items, PaletteItem{
                ID:         component.ID,
                Name:       component.Name,
                Category:   component.Category,
                Icon:       component.Icon,
                EntityType: component.ID,  // Use component ID
            })
        }
    }
    
    return items
}

// In Canvas component
func RenderEntity(instance *ComponentInstance) {
    desc := registry.Get(instance.ComponentID)
    
    // Use descriptor properties for rendering
    renderIcon(desc.Icon)
    renderLabel(instance.Name)
    renderTerminals(desc.Terminals)
}

// In Inspector component
func GetPropertyEditors(instance *ComponentInstance) []PropertyEditor {
    desc := registry.Get(instance.ComponentID)
    
    editors := make([]PropertyEditor, 0)
    for _, prop := range desc.Properties {
        editors = append(editors, PropertyEditor{
            Key:    prop.Key,
            Label:  prop.Label,
            Type:   prop.Type,
            Value:  instance.Properties[prop.Key],
            Unit:   prop.Unit,
            Min:    prop.Min,
            Max:    prop.Max,
            Options: prop.Options,
        })
    }
    
    return editors
}
```

## Domain Plugins

### forge-electrical

```go
package forgeelectrical

func init() {
    // Grid
    RegisterComponent("grid", ...)
    
    // Bus
    RegisterComponent("bus", ...)
    
    // Breaker
    RegisterComponent("breaker", ...)
    
    // Transformer
    RegisterComponent("transformer", ...)
    
    // Generator
    RegisterComponent("generator", ...)
    
    // Load
    RegisterComponent("load", ...)
    
    // Meter
    RegisterComponent("meter", ...)
}
```

### forge-environment

```go
package forgeenvironment

func init() {
    RegisterComponent("sun", ...)
    RegisterComponent("weather", ...)
    RegisterComponent("wind", ...)
}
```

### forge-simulation

```go
package forgesimulation

func init() {
    RegisterComponent("scenario", ...)
    RegisterComponent("clock", ...)
}
```

## Registry API (REST)

### GET /api/registry/components

Returns all registered components.

```json
{
  "components": [
    {
      "id": "forge-electrical:grid",
      "name": "Utility Grid",
      "category": "electrical",
      "icon": "🔌",
      "properties": [...],
      "terminals": [...]
    }
  ]
}
```

### GET /api/registry/categories

Returns all categories.

```json
{
  "categories": [
    {
      "id": "electrical",
      "name": "Electrical",
      "icon": "⚡",
      "order": 1
    }
  ]
}
```

### GET /api/registry/component/:id

Returns a single component.

## Validation

### Connection Validation

```go
func CanConnect(source, target *TerminalDescriptor) bool {
    // Source -> Bus always valid
    if target.Role == "bus" {
        return true
    }
    
    // Bus -> Source always valid
    if source.Role == "bus" {
        return true
    }
    
    // Voltage must match
    if source.Voltage != nil && target.Voltage != nil {
        return *source.Voltage == *target.Voltage
    }
    
    return true
}
```

## Benefits

1. **Decoupling**: Editor knows nothing about electrical/environment/simulation
2. **Extensibility**: Add new domains by registering components
3. **Consistency**: All components use the same property system
4. **Validation**: Centralized connection and property validation
5. **Discovery**: Registry enables dynamic palette generation

## Acceptance Criteria

Adding a new component requires only:

1. **Register component**: `registry.Register(&ComponentDescriptor{...})`
2. **Provide descriptor**: Define properties, terminals, icons
3. **Provide factory**: Create runtime entity
4. **Provide icon**: Specify emoji or icon identifier

No editor modifications required.

## Example: Adding a Water Domain

```go
// In forge-water plugin
func init() {
    registry.Register(&ComponentDescriptor{
        ID:       "forge-water:pump",
        Name:     "Water Pump",
        Category: "water",
        Icon:     "💧",
        Properties: []PropertyDescriptor{
            {
                Key:     "flow_rate",
                Label:   "Flow Rate",
                Type:    PropertyTypeNumber,
                Default: float64(100),
                Unit:    "L/min",
            },
        },
        Terminals: []TerminalDescriptor{
            {ID: "input", Name: "Input", Direction: TerminalDirectionInput},
            {ID: "output", Name: "Output", Direction: TerminalDirectionOutput},
        },
        Width:   60,
        Height:  60,
        Domain:  "forge-water",
    })
    
    registry.RegisterCategory(&ComponentCategory{
        ID:    "water",
        Name:  "Water",
        Icon:  "💧",
        Order: 4,
        Domain: "forge-water",
    })
}
```

The Editor automatically shows "Water Pump" in the Palette under a new "💧 Water" category.

---

*Last Updated: 2026-07-13*
