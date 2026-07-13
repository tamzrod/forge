// Package solver provides simulation solvers that advance the simulation state.
package solver

import (
	"fmt"
	"time"

	"github.com/tamzrod/forge/topology"
)

// ElectricalSolver solves electrical power balance.
type ElectricalSolver struct {
	BaseSolver

	topology *topology.Network

	// Aggregated values
	totalGeneration    float32
	totalConsumption  float32
	netPower          float32
	totalCapacity      float32
	availableGeneration float32

	// Entity collections
	generators []GeneratorSource
	loads      []LoadSink
	meters     []MeterMeasurement
}

// GeneratorSource is the interface for power generators.
type GeneratorSource interface {
	ActivePower() float32
	ReactivePower() float32
	RatedCapacity() float32
	AvailableCapacity() float32
	IsOnline() bool
	IsDispatchable() bool
	Name() string
}

// LoadSink is the interface for power loads.
type LoadSink interface {
	ActivePowerDemand() float32
	ReactivePowerDemand() float32
	IsConnected() bool
	Priority() int
	Name() string
}

// MeterMeasurement is the interface for meters.
type MeterMeasurement interface {
	SetMeasurements(p, q, v, f float32)
	RecordEnergy(dt time.Duration)
	Name() string
}

// NewElectricalSolver creates a new electrical solver.
func NewElectricalSolver() *ElectricalSolver {
	return &ElectricalSolver{
		BaseSolver: NewBaseSolver("Electrical", "power-balance"),
	}
}

// SetTopology sets the electrical topology.
func (s *ElectricalSolver) SetTopology(t interface{}) {
	if net, ok := t.(*topology.Network); ok {
		s.topology = net
	}
}

// Tick advances the electrical simulation.
func (s *ElectricalSolver) Tick(dt time.Duration) {
	if s.World() == nil {
		return
	}

	// Collect entities first
	s.collectEntities()

	// Tick all entities directly
	for _, e := range s.World().Entities() {
		e.Tick(dt)
	}

	// Calculate generation
	s.calculateGeneration()

	// Calculate consumption
	s.calculateConsumption()

	// Determine net power
	s.calculateNetPower()

	// Propagate to meters
	s.propagateToMeters(dt)
}

// Reset clears solver state.
func (s *ElectricalSolver) Reset() {
	s.totalGeneration = 0
	s.totalConsumption = 0
	s.netPower = 0
	s.totalCapacity = 0
	s.availableGeneration = 0
	s.generators = nil
	s.loads = nil
	s.meters = nil
}

// collectEntities gathers all electrical entities from the world.
func (s *ElectricalSolver) collectEntities() {
	s.generators = nil
	s.loads = nil
	s.meters = nil

	w := s.World()
	if w == nil {
		return
	}

	// Collect generators
	genEntities := w.EntitiesByType("virtual-generator")
	for _, e := range genEntities {
		if g, ok := e.(GeneratorSource); ok {
			s.generators = append(s.generators, g)
		}
	}

	// Collect loads
	loadEntities := w.EntitiesByType("virtual-load")
	for _, e := range loadEntities {
		if l, ok := e.(LoadSink); ok {
			s.loads = append(s.loads, l)
		}
	}

	// Collect meters
	meterEntities := w.EntitiesByType("virtual-meter")
	for _, e := range meterEntities {
		if m, ok := e.(MeterMeasurement); ok {
			s.meters = append(s.meters, m)
		}
	}
}

// calculateGeneration sums generator output.
func (s *ElectricalSolver) calculateGeneration() {
	s.totalGeneration = 0
	s.totalCapacity = 0
	s.availableGeneration = 0

	for _, g := range s.generators {
		if g.IsOnline() {
			s.totalGeneration += g.ActivePower()
			s.totalCapacity += g.RatedCapacity()
			s.availableGeneration += g.AvailableCapacity()
		}
	}
}

// calculateConsumption sums load demand.
func (s *ElectricalSolver) calculateConsumption() {
	s.totalConsumption = 0

	for _, l := range s.loads {
		if l.IsConnected() {
			s.totalConsumption += l.ActivePowerDemand()
		}
	}
}

// calculateNetPower determines the power balance.
func (s *ElectricalSolver) calculateNetPower() {
	s.netPower = s.totalGeneration - s.totalConsumption
}

// propagateToMeters updates all meters with current power flow.
func (s *ElectricalSolver) propagateToMeters(dt time.Duration) {
	for _, m := range s.meters {
		// For now, all meters see the same net power
		// Future: meters at specific topology locations
		m.SetMeasurements(s.netPower, 0, 480, 60)
		m.RecordEnergy(dt)
	}
}

// TotalGeneration returns the total generation in kW.
func (s *ElectricalSolver) TotalGeneration() float32 {
	return s.totalGeneration
}

// TotalConsumption returns the total consumption in kW.
func (s *ElectricalSolver) TotalConsumption() float32 {
	return s.totalConsumption
}

// NetPower returns the net power in kW (positive = export).
func (s *ElectricalSolver) NetPower() float32 {
	return s.netPower
}

// TotalCapacity returns the total generator capacity in kW.
func (s *ElectricalSolver) TotalCapacity() float32 {
	return s.totalCapacity
}

// AvailableGeneration returns the available generation in kW.
func (s *ElectricalSolver) AvailableGeneration() float32 {
	return s.availableGeneration
}

// GeneratorCount returns the number of generators.
func (s *ElectricalSolver) GeneratorCount() int {
	return len(s.generators)
}

// LoadCount returns the number of loads.
func (s *ElectricalSolver) LoadCount() int {
	return len(s.loads)
}

// MeterCount returns the number of meters.
func (s *ElectricalSolver) MeterCount() int {
	return len(s.meters)
}

// Generators returns the generator list.
func (s *ElectricalSolver) Generators() []GeneratorSource {
	return s.generators
}

// Loads returns the load list.
func (s *ElectricalSolver) Loads() []LoadSink {
	return s.loads
}

// String returns a summary of the electrical state.
func (s *ElectricalSolver) String() string {
	return fmt.Sprintf(
		"Electrical: Gen=%.1fkW, Load=%.1fkW, Net=%.1fkW (%s), Cap=%.1fkW, Avail=%.1fkW",
		s.totalGeneration,
		s.totalConsumption,
		s.netPower,
		s.powerDirection(),
		s.totalCapacity,
		s.availableGeneration,
	)
}

func (s *ElectricalSolver) powerDirection() string {
	if s.netPower > 0.1 {
		return "EXPORT"
	} else if s.netPower < -0.1 {
		return "IMPORT"
	}
	return "BALANCED"
}

// PowerBalance returns generation - consumption.
func (s *ElectricalSolver) PowerBalance() float32 {
	return s.totalGeneration - s.totalConsumption
}

// CapacityUtilization returns generation / capacity.
func (s *ElectricalSolver) CapacityUtilization() float32 {
	if s.totalCapacity == 0 {
		return 0
	}
	return s.totalGeneration / s.totalCapacity
}

// LoadServed returns consumption / demand.
func (s *ElectricalSolver) LoadServed() float32 {
	// Calculate total demand (including disconnected loads)
	totalDemand := float32(0)
	for _, l := range s.loads {
		totalDemand += l.ActivePowerDemand()
	}
	if totalDemand == 0 {
		return 1.0
	}
	return s.totalConsumption / totalDemand
}

// ExcessCapacity returns available - consumption.
func (s *ElectricalSolver) ExcessCapacity() float32 {
	return s.availableGeneration - s.totalConsumption
}
