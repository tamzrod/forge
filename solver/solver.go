// Package solver provides simulation solvers that advance the simulation state.
// Solvers determine how the simulated world evolves.
package solver

import (
	"time"

	"github.com/tamzrod/forge/world"
)

// Solver advances the simulation state.
// It owns evaluation order, dependency resolution, and state propagation.
type Solver interface {
	// Name returns the solver name.
	Name() string

	// Type returns the solver type.
	Type() string

	// Tick advances the simulation by one step.
	// The solver evaluates entities and propagates state.
	Tick(dt time.Duration)

	// Reset clears solver state.
	Reset()

	// SetWorld sets the world to solve.
	SetWorld(w world.World)

	// SetTopology sets the topology for spatial calculations.
	SetTopology(t interface{})
}

// BaseSolver provides common solver functionality.
type BaseSolver struct {
	name  string
	typ   string
	world world.World
}

// NewBaseSolver creates a new base solver.
func NewBaseSolver(name, typ string) BaseSolver {
	return BaseSolver{
		name: name,
		typ:  typ,
	}
}

// Name returns the solver name.
func (s *BaseSolver) Name() string {
	return s.name
}

// Type returns the solver type.
func (s *BaseSolver) Type() string {
	return s.typ
}

// SetWorld sets the world.
func (s *BaseSolver) SetWorld(w world.World) {
	s.world = w
}

// World returns the world.
func (s *BaseSolver) World() world.World {
	return s.world
}

// Reset clears solver state.
func (s *BaseSolver) Reset() {
	// Default implementation does nothing
}
