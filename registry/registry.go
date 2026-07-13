// Package registry provides the Component Registry for Forge.
// The registry is the authoritative source for all engineering components.
package registry

import (
	"fmt"
	"sort"
	"sync"
)

// PropertyType defines the type of a property.
type PropertyType string

const (
	PropertyTypeString  PropertyType = "string"
	PropertyTypeNumber  PropertyType = "number"
	PropertyTypeBoolean PropertyType = "boolean"
	PropertyTypeEnum   PropertyType = "enum"
)

// TerminalRole defines the role of a terminal.
type TerminalRole string

const (
	TerminalRoleSource      TerminalRole = "source"
	TerminalRoleDestination TerminalRole = "destination"
	TerminalRoleThrough    TerminalRole = "through"
	TerminalRoleObservation TerminalRole = "observation"
)

// TerminalDirection defines the direction of a terminal.
type TerminalDirection string

const (
	TerminalDirectionInput       TerminalDirection = "input"
	TerminalDirectionOutput      TerminalDirection = "output"
	TerminalDirectionBidirectional TerminalDirection = "bidirectional"
)

// PropertyDescriptor describes an editable property.
type PropertyDescriptor struct {
	Key      string       `json:"key"`
	Label    string       `json:"label"`
	Type     PropertyType `json:"type"`
	Default  interface{}  `json:"default"`
	Unit     string       `json:"unit,omitempty"`
	Min      *float64     `json:"min,omitempty"`
	Max      *float64     `json:"max,omitempty"`
	Options  []string     `json:"options,omitempty"`
	ReadOnly bool         `json:"readonly,omitempty"`
	Required bool         `json:"required,omitempty"`
}

// TerminalDescriptor describes a connection point.
type TerminalDescriptor struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Role      TerminalRole      `json:"role"`
	Voltage   *float64          `json:"voltage,omitempty"`
	Direction TerminalDirection `json:"direction"`
}

// ComponentDescriptor describes a component.
type ComponentDescriptor struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Category     string                 `json:"category"`
	Icon         string                 `json:"icon"`
	Description  string                 `json:"description,omitempty"`
	Properties   []PropertyDescriptor   `json:"properties"`
	Terminals    []TerminalDescriptor    `json:"terminals"`
	Width        float64                `json:"width"`
	Height       float64                `json:"height"`
	Domain       string                 `json:"domain"`
	Capabilities []string               `json:"capabilities,omitempty"`
}

// ComponentCategory organizes components into groups.
type ComponentCategory struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Order       int    `json:"order"`
	Domain      string `json:"domain"`
	Expandable  bool   `json:"expandable,omitempty"`
}

// ComponentFactory creates runtime entities.
type ComponentFactory func(instance *ComponentInstance) (interface{}, error)

// ComponentInstance is an actual placed component.
type ComponentInstance struct {
	ID          string                 `json:"id"`
	ComponentID string                 `json:"component_id"`
	Name        string                 `json:"name"`
	Position    Point                  `json:"position"`
	Properties  map[string]interface{} `json:"properties"`
	Connections []Connection           `json:"connections,omitempty"`
}

// Point represents a 2D coordinate.
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Size represents dimensions.
type Size struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// Connection represents a connection between terminals.
type Connection struct {
	ID           string `json:"id"`
	FromInstance string `json:"from_instance"`
	FromTerminal string `json:"from_terminal"`
	ToInstance   string `json:"to_instance"`
	ToTerminal   string `json:"to_terminal"`
	BusID        string `json:"bus_id,omitempty"`
}

// Registry is the component registry.
type Registry struct {
	mu         sync.RWMutex
	components map[string]*ComponentDescriptor
	categories map[string]*ComponentCategory
	factories  map[string]ComponentFactory
}

// New creates a new registry.
func New() *Registry {
	return &Registry{
		components: make(map[string]*ComponentDescriptor),
		categories: make(map[string]*ComponentCategory),
		factories:  make(map[string]ComponentFactory),
	}
}

// Global registry instance
var globalRegistry = New()

// GetRegistry returns the global registry.
func GetRegistry() *Registry {
	return globalRegistry
}

// Register registers a component.
func (r *Registry) Register(desc *ComponentDescriptor) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if desc.ID == "" {
		return fmt.Errorf("component ID is required")
	}
	if desc.Name == "" {
		return fmt.Errorf("component name is required")
	}
	if desc.Category == "" {
		return fmt.Errorf("component category is required")
	}

	r.components[desc.ID] = desc
	return nil
}

// RegisterFactory registers a factory for a component.
func (r *Registry) RegisterFactory(componentID string, factory ComponentFactory) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.factories[componentID] = factory
}

// Get retrieves a component by ID.
func (r *Registry) Get(id string) *ComponentDescriptor {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.components[id]
}

// List returns all registered components.
func (r *Registry) List() []*ComponentDescriptor {
	r.mu.RLock()
	defer r.mu.RUnlock()

	components := make([]*ComponentDescriptor, 0, len(r.components))
	for _, c := range r.components {
		components = append(components, c)
	}

	sort.Slice(components, func(i, j int) bool {
		return components[i].Name < components[j].Name
	})

	return components
}

// ListByCategory returns components in a category.
func (r *Registry) ListByCategory(categoryID string) []*ComponentDescriptor {
	r.mu.RLock()
	defer r.mu.RUnlock()

	components := make([]*ComponentDescriptor, 0)
	for _, c := range r.components {
		if c.Category == categoryID {
			components = append(components, c)
		}
	}

	sort.Slice(components, func(i, j int) bool {
		return components[i].Name < components[j].Name
	})

	return components
}

// RegisterCategory registers a category.
func (r *Registry) RegisterCategory(cat *ComponentCategory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if cat.ID == "" {
		return fmt.Errorf("category ID is required")
	}
	if cat.Name == "" {
		return fmt.Errorf("category name is required")
	}

	r.categories[cat.ID] = cat
	return nil
}

// Categories returns all categories.
func (r *Registry) Categories() []*ComponentCategory {
	r.mu.RLock()
	defer r.mu.RUnlock()

	categories := make([]*ComponentCategory, 0, len(r.categories))
	for _, c := range r.categories {
		categories = append(categories, c)
	}

	sort.Slice(categories, func(i, j int) bool {
		return categories[i].Order < categories[j].Order
	})

	return categories
}

// CreateInstance creates a new component instance.
func (r *Registry) CreateInstance(componentID string, name string, position Point) (*ComponentInstance, error) {
	desc := r.Get(componentID)
	if desc == nil {
		return nil, fmt.Errorf("component not found: %s", componentID)
	}

	// Create instance with default properties
	properties := make(map[string]interface{})
	for _, prop := range desc.Properties {
		properties[prop.Key] = prop.Default
	}

	instance := &ComponentInstance{
		ID:          fmt.Sprintf("instance-%d", len(properties)),
		ComponentID: componentID,
		Name:        name,
		Position:    position,
		Properties:  properties,
	}

	return instance, nil
}

// CreateFromFactory creates an entity using the registered factory.
func (r *Registry) CreateFromFactory(instance *ComponentInstance) (interface{}, error) {
	r.mu.RLock()
	factory := r.factories[instance.ComponentID]
	r.mu.RUnlock()

	if factory == nil {
		return nil, fmt.Errorf("no factory registered for component: %s", instance.ComponentID)
	}

	return factory(instance)
}

// CanConnect checks if two terminals can be connected.
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

// PaletteItem represents an item in the palette.
type PaletteItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Icon        string `json:"icon"`
	ComponentID string `json:"component_id"`
	Description string `json:"description,omitempty"`
}

// GetPaletteItems returns all items for the palette.
func (r *Registry) GetPaletteItems() []PaletteItem {
	items := make([]PaletteItem, 0)

	for _, cat := range r.Categories() {
		for _, comp := range r.ListByCategory(cat.ID) {
			items = append(items, PaletteItem{
				ID:          comp.ID,
				Name:        comp.Name,
				Category:    comp.Category,
				Icon:        comp.Icon,
				ComponentID: comp.ID,
				Description: comp.Description,
			})
		}
	}

	return items
}

// Reset clears all registrations (for testing).
func (r *Registry) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.components = make(map[string]*ComponentDescriptor)
	r.categories = make(map[string]*ComponentCategory)
	r.factories = make(map[string]ComponentFactory)
}
