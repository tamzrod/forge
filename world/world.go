// Package world provides the simulation world and entity models.
// Forge simulates reality - clients observe or influence that state.
package world

import (
	"fmt"
	"sync"
	"time"
)

// EntityID uniquely identifies an entity in the world.
type EntityID string

// World is the container for all simulated entities.
// It owns the simulation state and advances all entities.
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

	// Tick advances the world by one time step.
	// All entities update their state.
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

	// Close cleans up the world.
	Close()
}

// Entity is the base interface for all simulated entities.
type Entity interface {
	// ID returns the entity's unique identifier.
	ID() EntityID

	// Type returns the entity type name.
	Type() string

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
	id         EntityID
	entityType string
	inputs     map[string]Input
	outputs    map[string]Output
	connections map[string]struct {
		source EntityID
		output string
	}
}

// NewBaseEntity creates a new base entity.
func NewBaseEntity(id EntityID, entityType string) BaseEntity {
	return BaseEntity{
		id:          id,
		entityType:  entityType,
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
	time     time.Time
	events   []Event
	eventID  int
}

// NewWorld creates a new simulation world.
func NewWorld() World {
	return &simpleWorld{
		entities: make(map[EntityID]Entity),
		time:     time.Time{},
		events:   make([]Event, 0),
	}
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

func (w *simpleWorld) Tick(dt time.Duration) {
	w.mu.Lock()
	w.time = w.time.Add(dt)
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
	evt.Time = w.time
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
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.time
}

func (w *simpleWorld) Close() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.entities = nil
	w.events = nil
}
