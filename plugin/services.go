// Package plugin provides the Forge Plugin System.
//
// This file contains default implementations of Core Services.
package plugin

import (
	"fmt"
	"sort"
	"sync"
)

// DefaultComponentCatalog is the default implementation of ComponentCatalog.
type DefaultComponentCatalog struct {
	mu         sync.RWMutex
	components map[string]*ComponentDescriptor
	categories map[string]*ComponentCategory
}

// NewComponentCatalog creates a new component catalog.
func NewComponentCatalog() *DefaultComponentCatalog {
	return &DefaultComponentCatalog{
		components: make(map[string]*ComponentDescriptor),
		categories: make(map[string]*ComponentCategory),
	}
}

// Register implements ComponentCatalog.
func (c *DefaultComponentCatalog) Register(desc *ComponentDescriptor) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if desc.ID == "" {
		return fmt.Errorf("component ID is required")
	}
	if desc.Name == "" {
		return fmt.Errorf("component name is required")
	}
	if desc.Category == "" {
		return fmt.Errorf("component category is required")
	}

	c.components[desc.ID] = desc
	return nil
}

// Get implements ComponentCatalog.
func (c *DefaultComponentCatalog) Get(id string) *ComponentDescriptor {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.components[id]
}

// List implements ComponentCatalog.
func (c *DefaultComponentCatalog) List() []*ComponentDescriptor {
	c.mu.RLock()
	defer c.mu.RUnlock()

	components := make([]*ComponentDescriptor, 0, len(c.components))
	for _, comp := range c.components {
		components = append(components, comp)
	}

	sort.Slice(components, func(i, j int) bool {
		return components[i].Name < components[j].Name
	})

	return components
}

// ListByCategory implements ComponentCatalog.
func (c *DefaultComponentCatalog) ListByCategory(categoryID string) []*ComponentDescriptor {
	c.mu.RLock()
	defer c.mu.RUnlock()

	components := make([]*ComponentDescriptor, 0)
	for _, comp := range c.components {
		if comp.Category == categoryID {
			components = append(components, comp)
		}
	}

	sort.Slice(components, func(i, j int) bool {
		return components[i].Name < components[j].Name
	})

	return components
}

// RegisterCategory implements ComponentCatalog.
func (c *DefaultComponentCatalog) RegisterCategory(cat *ComponentCategory) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if cat.ID == "" {
		return fmt.Errorf("category ID is required")
	}
	if cat.Name == "" {
		return fmt.Errorf("category name is required")
	}

	c.categories[cat.ID] = cat
	return nil
}

// Categories implements ComponentCatalog.
func (c *DefaultComponentCatalog) Categories() []*ComponentCategory {
	c.mu.RLock()
	defer c.mu.RUnlock()

	categories := make([]*ComponentCategory, 0, len(c.categories))
	for _, cat := range c.categories {
		categories = append(categories, cat)
	}

	sort.Slice(categories, func(i, j int) bool {
		return categories[i].Order < categories[j].Order
	})

	return categories
}

// Reset clears all registrations (for testing).
func (c *DefaultComponentCatalog) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.components = make(map[string]*ComponentDescriptor)
	c.categories = make(map[string]*ComponentCategory)
}

// DefaultFactoryRegistry is the default implementation of FactoryRegistry.
type DefaultFactoryRegistry struct {
	mu       sync.RWMutex
	factories map[string]ComponentFactory
}

// NewFactoryRegistry creates a new factory registry.
func NewFactoryRegistry() *DefaultFactoryRegistry {
	return &DefaultFactoryRegistry{
		factories: make(map[string]ComponentFactory),
	}
}

// Register implements FactoryRegistry.
func (r *DefaultFactoryRegistry) Register(componentID string, factory ComponentFactory) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.factories[componentID] = factory
}

// Get implements FactoryRegistry.
func (r *DefaultFactoryRegistry) Get(componentID string) ComponentFactory {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.factories[componentID]
}

// Create implements FactoryRegistry.
func (r *DefaultFactoryRegistry) Create(componentID string, instance *ComponentInstance) (interface{}, error) {
	r.mu.RLock()
	factory := r.factories[componentID]
	r.mu.RUnlock()

	if factory == nil {
		return nil, fmt.Errorf("no factory registered for component: %s", componentID)
	}

	return factory(instance)
}

// Reset clears all registrations (for testing).
func (r *DefaultFactoryRegistry) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.factories = make(map[string]ComponentFactory)
}

// DefaultEventBus is the default implementation of EventBus.
type DefaultEventBus struct {
	mu       sync.RWMutex
	handlers []EventHandler
	events   []Event
}

// NewEventBus creates a new event bus.
func NewEventBus() *DefaultEventBus {
	return &DefaultEventBus{
		handlers: make([]EventHandler, 0),
		events:   make([]Event, 0),
	}
}

// Publish implements EventBus.
func (b *DefaultEventBus) Publish(event Event) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.events = append(b.events, event)

	// Deliver to handlers immediately
	for _, handler := range b.handlers {
		handler(event)
	}
}

// Subscribe implements EventBus.
func (b *DefaultEventBus) Subscribe(handler EventHandler) UnsubscribeFunc {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers = append(b.handlers, handler)

	return func() {
		b.mu.Lock()
		defer b.mu.Unlock()

		for i, h := range b.handlers {
			if h == handler {
				b.handlers = append(b.handlers[:i], b.handlers[i+1:]...)
				return
			}
		}
	}
}

// Events implements EventBus.
func (b *DefaultEventBus) Events() []Event {
	b.mu.Lock()
	defer b.mu.Unlock()

	events := b.events
	b.events = make([]Event, 0)
	return events
}
