// Package plugin provides the Forge Plugin System.
//
// The Plugin System enables Forge Core to remain domain-independent while
// allowing engineering domains to contribute components, solvers, validators,
// scenarios, and other domain-specific assets through plugins.
//
// Core Principle: Forge Core must never contain engineering-domain knowledge.
// All domain logic lives in plugins.
package plugin

import "github.com/tamzrod/forge/world"

// Plugin is the contract that all Forge plugins must implement.
//
// Plugins contribute capabilities to Forge without modifying Core.
// Each plugin represents an engineering domain or runtime extension.
type Plugin interface {
	// ID returns the unique identifier for this plugin.
	// Format: "namespace/name" (e.g., "forge-electrical", "forge-water")
	ID() string

	// Name returns the human-readable name of the plugin.
	Name() string

	// Version returns the plugin version in semver format.
	Version() string

	// Description returns a description of the plugin's purpose.
	Description() string

	// Dependencies returns the plugin IDs this plugin depends on.
	// Plugins are initialized after their dependencies.
	Dependencies() []string

	// OnInit is called when the plugin is initialized.
	// Use this to register components, factories, and validators with Core Services.
	OnInit(ctx Context) error

	// OnShutdown is called when the plugin is being unloaded.
	// Use this to clean up resources.
	OnShutdown() error
}

// Context provides access to Core Services during plugin initialization.
type Context interface {
	// ComponentCatalog returns the component catalog service.
	ComponentCatalog() ComponentCatalog

	// FactoryRegistry returns the factory registry service.
	FactoryRegistry() FactoryRegistry

	// EventBus returns the event bus service.
	EventBus() EventBus

	// World returns the simulation world.
	World() world.World

	// Logger returns the plugin logger.
	Logger() Logger

	// Config returns the plugin configuration.
	Config() Config
}

// ComponentCatalog stores component metadata.
//
// The catalog is metadata-only: it knows what components exist and their
// properties, but not how to simulate them (that's the plugin's domain).
type ComponentCatalog interface {
	// Register registers a component descriptor.
	Register(desc *ComponentDescriptor) error

	// Get retrieves a component by ID.
	Get(id string) *ComponentDescriptor

	// List returns all registered components.
	List() []*ComponentDescriptor

	// ListByCategory returns components in a specific category.
	ListByCategory(categoryID string) []*ComponentDescriptor

	// RegisterCategory registers a component category.
	RegisterCategory(cat *ComponentCategory) error

	// Categories returns all registered categories.
	Categories() []*ComponentCategory
}

// FactoryRegistry stores entity factories.
//
// Plugins register factories to create runtime entities from component instances.
type FactoryRegistry interface {
	// Register registers a factory for a component.
	Register(componentID string, factory ComponentFactory)

	// Get retrieves a factory by component ID.
	Get(componentID string) ComponentFactory

	// Create creates an entity using the registered factory.
	Create(componentID string, instance *ComponentInstance) (interface{}, error)
}

// ComponentFactory creates runtime entities from component instances.
type ComponentFactory func(instance *ComponentInstance) (interface{}, error)

// ComponentInstance represents an instance of a component placed in the simulation.
type ComponentInstance struct {
	ID          string
	ComponentID string
	Name        string
	Position    Point
	Properties  map[string]interface{}
}

// Point represents a 2D coordinate.
type Point struct {
	X float64
	Y float64
}

// ComponentDescriptor describes a component's editing and runtime properties.
type ComponentDescriptor struct {
	ID          string
	Name        string
	Category    string
	Icon        string
	Description string
	Properties  []PropertyDescriptor
	Terminals   []TerminalDescriptor
	Width       float64
	Height      float64
	Domain      string
}

// PropertyDescriptor defines an editable property.
type PropertyDescriptor struct {
	Key      string
	Label    string
	Type     PropertyType
	Default  interface{}
	Unit     string
	Min      *float64
	Max      *float64
	Options  []string
	ReadOnly bool
	Required bool
}

// PropertyType defines the type of a property.
type PropertyType string

const (
	PropertyTypeString  PropertyType = "string"
	PropertyTypeNumber PropertyType = "number"
	PropertyTypeBool   PropertyType = "boolean"
	PropertyTypeEnum   PropertyType = "enum"
)

// TerminalDescriptor describes a connection point.
type TerminalDescriptor struct {
	ID        string
	Name      string
	Role      TerminalRole
	Voltage   *float64
	Direction TerminalDirection
}

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
	TerminalDirectionInput  TerminalDirection = "input"
	TerminalDirectionOutput TerminalDirection = "output"
	TerminalDirectionBoth  TerminalDirection = "bidirectional"
)

// ComponentCategory organizes components into groups.
type ComponentCategory struct {
	ID         string
	Name       string
	Icon       string
	Order      int
	Domain     string
	Expandable bool
}

// ConnectionValidator validates domain-specific connections.
//
// Each domain plugin provides validators for its connection rules.
type ConnectionValidator interface {
	// Domain returns the domain this validator applies to.
	Domain() string

	// CanConnect validates whether two terminals can be connected.
	CanConnect(source, target *TerminalDescriptor) (bool, error)
}

// EventBus provides event pub/sub capabilities.
type EventBus interface {
	// Publish publishes an event.
	Publish(event Event)

	// Subscribe subscribes to events.
	Subscribe(handler EventHandler) UnsubscribeFunc

	// Events returns all events since the last call.
	Events() []Event
}

// Event represents a simulation event.
type Event struct {
	ID        string
	Type      string
	Source    string
	Data      map[string]interface{}
}

// EventHandler handles events.
type EventHandler func(Event)

// UnsubscribeFunc cancels an event subscription.
type UnsubscribeFunc func()

// Logger provides logging capabilities.
type Logger interface {
	// Debug logs a debug message.
	Debug(msg string, args ...interface{})

	// Info logs an info message.
	Info(msg string, args ...interface{})

	// Warn logs a warning message.
	Warn(msg string, args ...interface{})

	// Error logs an error message.
	Error(msg string, args ...interface{})
}

// Config provides plugin configuration.
type Config interface {
	// Get returns a configuration value by key.
	Get(key string) interface{}

	// GetString returns a string configuration value.
	GetString(key string, defaultValue string) string

	// GetInt returns an int configuration value.
	GetInt(key string, defaultValue int) int

	// GetFloat returns a float64 configuration value.
	GetFloat(key string, defaultValue float64) float64

	// GetBool returns a bool configuration value.
	GetBool(key string, defaultValue bool) bool
}
