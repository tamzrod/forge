// Package topology provides the domain-independent topology framework.
//
// Topology defines the structural relationships between entities.
// Domain-specific topology implementations (electrical, water, etc.)
// are provided by plugins.
package topology

import "github.com/tamzrod/forge/world"

// Topology is the domain-independent topology interface.
type Topology interface {
	// Entities returns all entities in the topology.
	Entities() []world.EntityID

	// Connections returns all connections in the topology.
	Connections() []Connection

	// AddEntity adds an entity to the topology.
	AddEntity(id world.EntityID)

	// RemoveEntity removes an entity from the topology.
	RemoveEntity(id world.EntityID)
}

// Connection represents a connection between two entities.
type Connection struct {
	From world.EntityID
	To   world.EntityID
}

// Terminal defines a connection point on an entity.
// Terminals are domain-specific and defined by plugins.
type Terminal interface {
	// EntityID returns the entity this terminal belongs to.
	EntityID() world.EntityID

	// Name returns the terminal name.
	Name() string

	// Role returns the terminal role.
	Role() string
}
