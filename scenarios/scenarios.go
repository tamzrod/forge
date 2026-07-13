// Package scenarios provides engineering experiment scenarios.
// Scenarios orchestrate events into the simulation world.
package scenarios

import (
	"fmt"
	"sync"
	"time"

	"github.com/tamzrod/forge/world"
)

// Scenario represents a repeatable engineering experiment.
// Scenarios schedule events into the World without directly manipulating entities.
type Scenario interface {
	// Name returns the scenario name.
	Name() string

	// Start begins the scenario.
	Start(w world.World) error

	// Stop halts the scenario.
	Stop()

	// IsRunning returns true if the scenario is active.
	IsRunning() bool

	// IsComplete returns true if the scenario has finished.
	IsComplete() bool

	// Progress returns 0.0 to 1.0 progress through the scenario.
	Progress() float64

	// Duration returns the expected duration.
	Duration() time.Duration

	// Elapsed returns the time since start.
	Elapsed() time.Duration

	// Events returns all events this scenario has published.
	Events() []world.Event

	// String returns a human-readable status.
	String() string
}

// EventAction represents an action to perform at a specific time.
type EventAction struct {
	Time    time.Duration
	Type    string
	Source  world.EntityID
	Data    map[string]interface{}
}

// BaseScenario provides common scenario functionality.
type BaseScenario struct {
	name         string
	description string
	duration    time.Duration
	startTime   time.Time
	endTime     time.Time
	running     bool
	complete    bool
	mu          sync.RWMutex
	events      []world.Event
	actions     []EventAction
}

// NewBaseScenario creates a new base scenario.
func NewBaseScenario(name, description string, duration time.Duration) *BaseScenario {
	return &BaseScenario{
		name:         name,
		description: description,
		duration:    duration,
		events:      make([]world.Event, 0),
		actions:     make([]EventAction, 0),
	}
}

// Name returns the scenario name.
func (s *BaseScenario) Name() string {
	return s.name
}

// Duration returns the expected duration.
func (s *BaseScenario) Duration() time.Duration {
	return s.duration
}

// Elapsed returns the time since start.
func (s *BaseScenario) Elapsed() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.startTime.IsZero() {
		return 0
	}
	return time.Since(s.startTime)
}

// IsRunning returns true if the scenario is active.
func (s *BaseScenario) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// IsComplete returns true if the scenario has finished.
func (s *BaseScenario) IsComplete() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.complete
}

// Progress returns 0.0 to 1.0 progress through the scenario.
func (s *BaseScenario) Progress() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.duration == 0 {
		return 0
	}
	elapsed := time.Since(s.startTime)
	progress := float64(elapsed) / float64(s.duration)
	if progress > 1 {
		progress = 1
	}
	return progress
}

// Events returns all events this scenario has published.
func (s *BaseScenario) Events() []world.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.events
}

// AddAction adds an event action to the scenario.
func (s *BaseScenario) AddAction(action EventAction) {
	s.actions = append(s.actions, action)
}

// Start begins the scenario.
func (s *BaseScenario) Start(w world.World) error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return fmt.Errorf("scenario %s is already running", s.name)
	}
	s.running = true
	s.complete = false
	s.startTime = time.Now()
	s.endTime = s.startTime.Add(s.duration)
	s.events = make([]world.Event, 0)
	s.mu.Unlock()
	return nil
}

// Stop halts the scenario.
func (s *BaseScenario) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.running = false
}

// Complete marks the scenario as complete.
func (s *BaseScenario) Complete() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.complete = true
	s.running = false
}

// RecordEvent records an event that was published.
func (s *BaseScenario) RecordEvent(evt world.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, evt)
}

// ProcessActions processes scheduled actions based on elapsed time.
func (s *BaseScenario) ProcessActions(w world.World) {
	if !s.IsRunning() {
		return
	}

	s.mu.RLock()
	elapsed := time.Since(s.startTime)
	s.mu.RUnlock()

	for _, action := range s.actions {
		if elapsed >= action.Time && elapsed < action.Time+100*time.Millisecond {
			evt := world.Event{
				Type:   action.Type,
				Source: action.Source,
				Data:   action.Data,
			}
			w.PublishEvent(evt)
			s.RecordEvent(evt)
		}
	}

	// Check if scenario is complete
	if elapsed >= s.duration {
		s.Complete()
	}
}

// String returns a human-readable status.
func (s *BaseScenario) String() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	status := "STOPPED"
	if s.running {
		status = "RUNNING"
	} else if s.complete {
		status = "COMPLETE"
	}
	return fmt.Sprintf("%s [%s] %.0f%%", s.name, status, s.Progress()*100)
}

// NormalDayScenario simulates a normal operating day.
type NormalDayScenario struct {
	*BaseScenario
}

// NewNormalDay creates a normal day scenario.
func NewNormalDay() *NormalDayScenario {
	s := &NormalDayScenario{
		BaseScenario: NewBaseScenario("Normal Day", "Typical operating day with gradual sun movement", 2*time.Minute),
	}
	// No special events - just lets the world run
	return s
}

// CloudPassingScenario simulates a cloud passing over the plant.
type CloudPassingScenario struct {
	*BaseScenario
	cloudCoverage float32
}

// NewCloudPassing creates a cloud passing scenario.
func NewCloudPassing() *CloudPassingScenario {
	s := &CloudPassingScenario{
		BaseScenario:   NewBaseScenario("Cloud Passing", "Cloud passes over PV arrays reducing irradiance", 90*time.Second),
		cloudCoverage: 0.7,
	}
	s.AddAction(EventAction{Time: 15 * time.Second, Type: "cloud_cover", Data: map[string]interface{}{"coverage": s.cloudCoverage}})
	s.AddAction(EventAction{Time: 45 * time.Second, Type: "cloud_cover", Data: map[string]interface{}{"coverage": float32(0)}})
	return s
}

// GridVoltageSagScenario simulates a grid voltage sag.
type GridVoltageSagScenario struct {
	*BaseScenario
	sagVoltage float32
}

// NewGridVoltageSag creates a grid voltage sag scenario.
func NewGridVoltageSag() *GridVoltageSagScenario {
	s := &GridVoltageSagScenario{
		BaseScenario: NewBaseScenario("Grid Voltage Sag", "Grid voltage drops to 80% for 5 seconds", 30*time.Second),
		sagVoltage:   55200, // 80% of 69kV
	}
	s.AddAction(EventAction{Time: 10 * time.Second, Type: "voltage_sag", Source: "utility-grid", Data: map[string]interface{}{"voltage": s.sagVoltage}})
	s.AddAction(EventAction{Time: 15 * time.Second, Type: "voltage_sag_end", Source: "utility-grid", Data: map[string]interface{}{"voltage": float32(69000)}})
	return s
}

// FrequencyExcursionScenario simulates a grid frequency excursion.
type FrequencyExcursionScenario struct {
	*BaseScenario
	excursionFreq float32
}

// NewFrequencyExcursion creates a frequency excursion scenario.
func NewFrequencyExcursion() *FrequencyExcursionScenario {
	s := &FrequencyExcursionScenario{
		BaseScenario:  NewBaseScenario("Frequency Excursion", "Grid frequency drops to 59.5 Hz", 30*time.Second),
		excursionFreq: 59.5,
	}
	s.AddAction(EventAction{Time: 10 * time.Second, Type: "frequency_sag", Source: "utility-grid", Data: map[string]interface{}{"frequency": s.excursionFreq}})
	s.AddAction(EventAction{Time: 20 * time.Second, Type: "frequency_sag_end", Source: "utility-grid", Data: map[string]interface{}{"frequency": float32(60.0)}})
	return s
}

// BreakerTripScenario simulates a breaker trip event.
type BreakerTripScenario struct {
	*BaseScenario
}

// NewBreakerTrip creates a breaker trip scenario.
func NewBreakerTrip() *BreakerTripScenario {
	s := &BreakerTripScenario{
		BaseScenario: NewBaseScenario("Breaker Trip", "Protection relay trips breaker then recloses after fault", 60*time.Second),
	}
	s.AddAction(EventAction{Time: 10 * time.Second, Type: "fault", Source: "grid-breaker", Data: map[string]interface{}{"type": "overcurrent", "location": "pcc", "severity": float32(1.0)}})
	s.AddAction(EventAction{Time: 30 * time.Second, Type: "reconnect_command", Data: map[string]interface{}{"reason": "reclose_after_fault"}})
	return s
}

// IslandingScenario simulates islanding and reconnection.
type IslandingScenario struct {
	*BaseScenario
}

// NewIslanding creates an islanding scenario.
func NewIslanding() *IslandingScenario {
	s := &IslandingScenario{
		BaseScenario: NewBaseScenario("Islanding Test", "Planned plant isolation from grid", 2*time.Minute),
	}
	s.AddAction(EventAction{Time: 15 * time.Second, Type: "island_command", Source: "grid-breaker", Data: map[string]interface{}{"reason": "planned_test"}})
	s.AddAction(EventAction{Time: 75 * time.Second, Type: "reconnect_command", Data: map[string]interface{}{"reason": "test_complete"}})
	return s
}

// LoadStepScenario simulates a sudden load change.
type LoadStepScenario struct {
	*BaseScenario
	stepChange float32
}

// NewLoadStep creates a load step scenario.
func NewLoadStep() *LoadStepScenario {
	s := &LoadStepScenario{
		BaseScenario: NewBaseScenario("Load Step", "Auxiliary load increases by 50%", 45*time.Second),
		stepChange:  7.5, // 50% increase from 5kW to 7.5kW
	}
	s.AddAction(EventAction{Time: 15 * time.Second, Type: "load_change", Data: map[string]interface{}{"load_id": "aux-load", "power": s.stepChange}})
	s.AddAction(EventAction{Time: 30 * time.Second, Type: "load_change", Data: map[string]interface{}{"load_id": "aux-load", "power": float32(5.0)}})
	return s
}

// GeneratorTripScenario simulates a generator trip.
type GeneratorTripScenario struct {
	*BaseScenario
	generatorID world.EntityID
}

// NewGeneratorTrip creates a generator trip scenario.
func NewGeneratorTrip() *GeneratorTripScenario {
	s := &GeneratorTripScenario{
		BaseScenario: NewBaseScenario("Generator Trip", "One generator trips offline then restarts", 90*time.Second),
		generatorID: "gen-1",
	}
	s.AddAction(EventAction{Time: 20 * time.Second, Type: "generator_fault", Source: s.generatorID, Data: map[string]interface{}{"generator_id": s.generatorID}})
	s.AddAction(EventAction{Time: 60 * time.Second, Type: "generator_clear", Source: s.generatorID, Data: map[string]interface{}{"generator_id": s.generatorID}})
	return s
}

// StormScenario simulates a storm affecting generation.
type StormScenario struct {
	*BaseScenario
}

// NewStorm creates a storm scenario.
func NewStorm() *StormScenario {
	s := &StormScenario{
		BaseScenario: NewBaseScenario("Storm", "Storm with cloud cover and wind gusts", 3*time.Minute),
	}
	// Gradual cloud buildup
	s.AddAction(EventAction{Time: 30 * time.Second, Type: "cloud_cover", Data: map[string]interface{}{"coverage": float32(0.3)}})
	s.AddAction(EventAction{Time: 60 * time.Second, Type: "cloud_cover", Data: map[string]interface{}{"coverage": float32(0.6)}})
	s.AddAction(EventAction{Time: 60 * time.Second, Type: "wind_gust", Data: map[string]interface{}{"speed": float32(15.0)}})
	// Storm passes
	s.AddAction(EventAction{Time: 120 * time.Second, Type: "cloud_cover", Data: map[string]interface{}{"coverage": float32(0.2)}})
	s.AddAction(EventAction{Time: 120 * time.Second, Type: "wind_gust", Data: map[string]interface{}{"speed": float32(5.0)}})
	s.AddAction(EventAction{Time: 150 * time.Second, Type: "cloud_cover", Data: map[string]interface{}{"coverage": float32(0)}})
	return s
}

// ScenarioRunner runs scenarios against a world.
type ScenarioRunner struct {
	world world.World
	mu    sync.RWMutex
	scenario *BaseScenario
}

// NewScenarioRunner creates a new scenario runner.
func NewScenarioRunner(w world.World) *ScenarioRunner {
	return &ScenarioRunner{
		world: w,
	}
}

// Load loads a scenario.
func (r *ScenarioRunner) Load(s Scenario) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.scenario != nil && r.scenario.IsRunning() {
		return fmt.Errorf("a scenario is already running")
	}
	// Use the embedded BaseScenario
	r.scenario = &BaseScenario{}
	if base, ok := interface{}(s).(interface{ GetBase() *BaseScenario }); ok {
		r.scenario = base.GetBase()
	} else if ptr, ok := s.(*BaseScenario); ok {
		r.scenario = ptr
	} else {
		return fmt.Errorf("unsupported scenario type")
	}
	return nil
}

// GetBase returns the embedded BaseScenario.
func (s *NormalDayScenario) GetBase() *BaseScenario { return s.BaseScenario }
func (s *CloudPassingScenario) GetBase() *BaseScenario { return s.BaseScenario }
func (s *GridVoltageSagScenario) GetBase() *BaseScenario { return s.BaseScenario }
func (s *FrequencyExcursionScenario) GetBase() *BaseScenario { return s.BaseScenario }
func (s *BreakerTripScenario) GetBase() *BaseScenario { return s.BaseScenario }
func (s *IslandingScenario) GetBase() *BaseScenario { return s.BaseScenario }
func (s *LoadStepScenario) GetBase() *BaseScenario { return s.BaseScenario }
func (s *GeneratorTripScenario) GetBase() *BaseScenario { return s.BaseScenario }
func (s *StormScenario) GetBase() *BaseScenario { return s.BaseScenario }

// Start begins the loaded scenario.
func (r *ScenarioRunner) Start() error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.scenario == nil {
		return fmt.Errorf("no scenario loaded")
	}
	return r.scenario.Start(r.world)
}

// Stop halts the current scenario.
func (r *ScenarioRunner) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.scenario != nil {
		r.scenario.Stop()
	}
}

// Tick advances the scenario runner.
func (r *ScenarioRunner) Tick(dt time.Duration) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.scenario != nil && r.scenario.IsRunning() {
		r.scenario.ProcessActions(r.world)
	}
}

// Scenario returns the current scenario.
func (r *ScenarioRunner) Scenario() Scenario {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.scenario
}

// IsRunning returns true if a scenario is running.
func (r *ScenarioRunner) IsRunning() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.scenario != nil && r.scenario.IsRunning()
}
