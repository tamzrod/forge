// Package world provides the simulation world and entity models.
// Forge simulates reality - clients observe or influence that state.
package world

import (
	"fmt"
	"sync"
	"time"

	"github.com/tamzrod/forge/simulation"
)

// EntityID uniquely identifies an entity in the world.
type EntityID string

// Capability represents what an entity can do in the simulation.
// Capabilities are domain-independent abstractions.
type Capability string

const (
	// CapabilityProduce generates or injects energy/material into the network.
	CapabilityProduce Capability = "produce"
	// CapabilityConsume uses or withdraws energy/material from the network.
	CapabilityConsume Capability = "consume"
	// CapabilityStore holds energy/material in storage.
	CapabilityStore Capability = "store"
	// CapabilityTransform changes energy/material form (e.g., AC/DC conversion).
	CapabilityTransform Capability = "transform"
	// CapabilityTransport moves energy/material through the network.
	CapabilityTransport Capability = "transport"
	// CapabilitySwitch can interrupt flow in the network.
	CapabilitySwitch Capability = "switch"
	// CapabilityMeasure observes and records values without affecting the network.
	CapabilityMeasure Capability = "measure"
	// CapabilityProtect responds to abnormal conditions.
	CapabilityProtect Capability = "protect"
	// CapabilityCommunicate sends or receives data.
	CapabilityCommunicate Capability = "communicate"
)

// Solver is the interface for simulation solvers.
// Solvers advance the simulation state.
type Solver interface {
	Name() string
	Type() string
	Tick(dt time.Duration)
	Reset()
	SetWorld(w World)
}

// World is the container for all simulated entities.
// It delegates simulation evolution to a Solver.
type World interface {
	// AddEntity adds an entity to the world.
	AddEntity(e Entity)

	// RemoveEntity removes an entity by ID.
	RemoveEntity(id EntityID)

	// Entity returns an entity by ID.
	Entity(id EntityID) Entity

	// Entities returns all entities.
	Entities() []Entity

	// EntitiesByType returns entities matching a type.
	EntitiesByType(entityType string) []Entity

	// EntitiesByCapability returns entities that have the given capability.
	EntitiesByCapability(capability Capability) []Entity

	// EntitiesByCapabilities returns entities that have all of the given capabilities.
	EntitiesByCapabilities(capabilities []Capability) []Entity

	// Tick advances the world by one time step.
	// Delegates to the configured Solver.
	Tick(dt time.Duration)

	// Measurement returns a measurement by entity ID and name.
	Measurement(entityID EntityID, name string) Measurement

	// AllMeasurements returns all current measurements.
	AllMeasurements() []Measurement

	// PublishEvent publishes an event to the world.
	PublishEvent(evt Event)

	// Events returns all events since the last call.
	Events() []Event

	// Time returns the current simulation time.
	Time() time.Time

	// Clock returns the simulation clock.
	Clock() simulation.Clock

	// SetClock sets the simulation clock.
	SetClock(clock simulation.Clock)

	// SetSolver sets the solver for this world.
	SetSolver(s Solver)

	// Solver returns the configured solver.
	Solver() Solver

	// Close cleans up the world.
	Close()
}

// Entity is the base interface for all simulated entities.
type Entity interface {
	// ID returns the entity's unique identifier.
	ID() EntityID

	// Type returns the entity type name (e.g., "grid", "battery", "pump").
	Type() string

	// Capabilities returns the entity's capabilities.
	// These are domain-independent abstractions.
	Capabilities() []Capability

	// Tick updates the entity state for one time step.
	// dt is the elapsed time since the last tick.
	Tick(dt time.Duration)

	// Measurements returns all current measurements for this entity.
	Measurements() []Measurement

	// Inputs returns the entity's input channels.
	Inputs() []Input

	// Outputs returns the entity's output channels.
	Outputs() []Output

	// HandleEvent processes an event.
	HandleEvent(evt Event)

	// Connect connects this entity to another via an input.
	Connect(inputName string, source EntityID, outputName string)

	// HasCapability returns true if the entity has the given capability.
	HasCapability(capability Capability) bool
}

// HasCapability returns true if the entity has the given capability.
func (e *BaseEntity) HasCapability(capability Capability) bool {
	for _, c := range e.capabilities {
		if c == capability {
			return true
		}
	}
	return false
}

// Input represents an input channel to an entity.
type Input struct {
	Name  string
	Value interface{}
}

// Output represents an output channel from an entity.
type Output struct {
	Name  string
	Value interface{}
}

// Measurement represents an observable value from an entity.
type Measurement struct {
	EntityID  EntityID
	Name      string
	Value     interface{}
	Unit      string
	Timestamp time.Time
}

// Event represents something that happened in the simulation.
type Event struct {
	ID        string
	Type      string
	Source    EntityID
	Time      time.Time
	Data      map[string]interface{}
}

// BaseEntity provides common functionality for entities.
type BaseEntity struct {
	id          EntityID
	entityType  string
	capabilities []Capability
	inputs      map[string]Input
	outputs     map[string]Output
	connections map[string]struct {
		source EntityID
		output string
	}
}

// NewBaseEntity creates a new base entity with no capabilities.
func NewBaseEntity(id EntityID, entityType string) BaseEntity {
	return BaseEntity{
		id:          id,
		entityType:  entityType,
		capabilities: []Capability{},
		inputs:      make(map[string]Input),
		outputs:     make(map[string]Output),
		connections: make(map[string]struct{ source EntityID; output string }),
	}
}

// NewBaseEntityWithCapabilities creates a new base entity with capabilities.
func NewBaseEntityWithCapabilities(id EntityID, entityType string, capabilities []Capability) BaseEntity {
	return BaseEntity{
		id:          id,
		entityType:  entityType,
		capabilities: capabilities,
		inputs:      make(map[string]Input),
		outputs:     make(map[string]Output),
		connections: make(map[string]struct{ source EntityID; output string }),
	}
}

// ID returns the entity ID.
func (e *BaseEntity) ID() EntityID {
	return e.id
}

// Type returns the entity type.
func (e *BaseEntity) Type() string {
	return e.entityType
}

// Capabilities returns the entity's capabilities.
func (e *BaseEntity) Capabilities() []Capability {
	return e.capabilities
}

// Inputs returns all inputs.
func (e *BaseEntity) Inputs() []Input {
	result := make([]Input, 0, len(e.inputs))
	for _, inp := range e.inputs {
		result = append(result, inp)
	}
	return result
}

// Outputs returns all outputs.
func (e *BaseEntity) Outputs() []Output {
	result := make([]Output, 0, len(e.outputs))
	for _, out := range e.outputs {
		result = append(result, out)
	}
	return result
}

// SetInput sets an input value.
func (e *BaseEntity) SetInput(name string, value interface{}) {
	e.inputs[name] = Input{Name: name, Value: value}
}

// GetInput gets an input value.
func (e *BaseEntity) GetInput(name string) interface{} {
	if inp, ok := e.inputs[name]; ok {
		return inp.Value
	}
	return nil
}

// SetOutput sets an output value.
func (e *BaseEntity) SetOutput(name string, value interface{}) {
	e.outputs[name] = Output{Name: name, Value: value}
}

// GetOutput gets an output value.
func (e *BaseEntity) GetOutput(name string) interface{} {
	if out, ok := e.outputs[name]; ok {
		return out.Value
	}
	return nil
}

// Connect connects this entity's input to another entity's output.
func (e *BaseEntity) Connect(inputName string, source EntityID, outputName string) {
	e.connections[inputName] = struct {
		source EntityID
		output string
	}{source: source, output: outputName}
}

// Measurements returns empty measurements (override in subtypes).
func (e *BaseEntity) Measurements() []Measurement {
	return []Measurement{}
}

// HandleEvent handles events (override in subtypes).
func (e *BaseEntity) HandleEvent(evt Event) {
}

// Tick updates state (override in subtypes).
func (e *BaseEntity) Tick(dt time.Duration) {
}

// simpleWorld is a basic implementation of World.
type simpleWorld struct {
	mu       sync.RWMutex
	entities map[EntityID]Entity
	clock    simulation.Clock
	solver   Solver
	events   []Event
	eventID  int
}

// NewWorld creates a new simulation world.
func NewWorld() World {
	return &simpleWorld{
		entities: make(map[EntityID]Entity),
		clock:    simulation.NewClock(),
		events:   make([]Event, 0),
	}
}

// SetSolver sets the solver for this world.
func (w *simpleWorld) SetSolver(s Solver) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.solver = s
	if s != nil {
		s.SetWorld(w)
	}
}

// Solver returns the configured solver.
func (w *simpleWorld) Solver() Solver {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.solver
}

func (w *simpleWorld) AddEntity(e Entity) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.entities[e.ID()] = e
}

func (w *simpleWorld) RemoveEntity(id EntityID) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.entities, id)
}

func (w *simpleWorld) Entity(id EntityID) Entity {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.entities[id]
}

func (w *simpleWorld) Entities() []Entity {
	w.mu.RLock()
	defer w.mu.RUnlock()
	result := make([]Entity, 0, len(w.entities))
	for _, e := range w.entities {
		result = append(result, e)
	}
	return result
}

func (w *simpleWorld) EntitiesByType(entityType string) []Entity {
	w.mu.RLock()
	defer w.mu.RUnlock()
	result := make([]Entity, 0)
	for _, e := range w.entities {
		if e.Type() == entityType {
			result = append(result, e)
		}
	}
	return result
}

func (w *simpleWorld) EntitiesByCapability(capability Capability) []Entity {
	w.mu.RLock()
	defer w.mu.RUnlock()
	result := make([]Entity, 0)
	for _, e := range w.entities {
		if e.HasCapability(capability) {
			result = append(result, e)
		}
	}
	return result
}

func (w *simpleWorld) EntitiesByCapabilities(capabilities []Capability) []Entity {
	w.mu.RLock()
	defer w.mu.RUnlock()
	result := make([]Entity, 0)
	for _, e := range w.entities {
		hasAll := true
		for _, cap := range capabilities {
			if !e.HasCapability(cap) {
				hasAll = false
				break
			}
		}
		if hasAll {
			result = append(result, e)
		}
	}
	return result
}

func (w *simpleWorld) Tick(dt time.Duration) {
	// If a solver is configured, delegate to it
	w.mu.RLock()
	s := w.solver
	w.mu.RUnlock()

	if s != nil {
		s.Tick(dt)
		return
	}

	// Otherwise, use the default behavior
	w.mu.Lock()
	entities := make([]Entity, len(w.entities))
	i := 0
	for _, e := range w.entities {
		entities[i] = e
		i++
	}
	w.mu.Unlock()

	// First, propagate outputs to inputs
	w.propagateSignals(entities)

	// Then tick all entities
	for _, e := range entities {
		e.Tick(dt)
	}
}

func (w *simpleWorld) propagateSignals(entities []Entity) {
	// Build a map of outputs by entity
	outputs := make(map[EntityID]map[string]interface{})
	for _, e := range entities {
		if base, ok := e.(*BaseEntity); ok {
			outputs[e.ID()] = make(map[string]interface{})
			for name, out := range base.outputs {
				outputs[e.ID()][name] = out.Value
			}
		}
	}

	// Propagate connected signals
	for _, e := range entities {
		if base, ok := e.(*BaseEntity); ok {
			for inputName, conn := range base.connections {
				if srcOutputs, ok := outputs[conn.source]; ok {
					if val, ok := srcOutputs[conn.output]; ok {
						base.SetInput(inputName, val)
					}
				}
			}
		}
	}
}

func (w *simpleWorld) Measurement(entityID EntityID, name string) Measurement {
	w.mu.RLock()
	defer w.mu.RUnlock()
	entity := w.entities[entityID]
	if entity == nil {
		return Measurement{}
	}
	for _, m := range entity.Measurements() {
		if m.Name == name {
			return m
		}
	}
	return Measurement{}
}

func (w *simpleWorld) AllMeasurements() []Measurement {
	w.mu.RLock()
	defer w.mu.RUnlock()
	result := make([]Measurement, 0)
	for _, e := range w.entities {
		result = append(result, e.Measurements()...)
	}
	return result
}

func (w *simpleWorld) PublishEvent(evt Event) {
	w.mu.Lock()
	defer w.mu.Unlock()
	evt.ID = fmt.Sprintf("evt-%d", w.eventID)
	w.eventID++
	evt.Time = w.Clock().Now()
	w.events = append(w.events, evt)

	// Dispatch event to entities
	for _, e := range w.entities {
		e.HandleEvent(evt)
	}
}

func (w *simpleWorld) Events() []Event {
	w.mu.Lock()
	defer w.mu.Unlock()
	events := w.events
	w.events = make([]Event, 0)
	return events
}

func (w *simpleWorld) Time() time.Time {
	return w.Clock().Now()
}

func (w *simpleWorld) Clock() simulation.Clock {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.clock
}

func (w *simpleWorld) SetClock(clock simulation.Clock) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.clock = clock
}

func (w *simpleWorld) Close() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.entities = nil
	w.events = nil
}
