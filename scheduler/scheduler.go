// Package scheduler provides the simulation scheduler.
package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/tamzrod/forge/device"
	"github.com/tamzrod/forge/models"
)

// Scheduler advances simulation time and tells devices to tick.
type Scheduler struct {
	mu          sync.Mutex
	devices     []*device.Device
	models      []models.Model
	tickInterval time.Duration
	clock        SimulationClock
	running      bool
	stopCh       chan struct{}
}

// SimulationClock tracks elapsed simulation time.
type SimulationClock struct {
	elapsed   time.Duration
	tickCount uint64
}

// Elapsed returns the total elapsed simulation time.
func (c SimulationClock) Elapsed() time.Duration {
	return c.elapsed
}

// TickCount returns the number of ticks executed.
func (c SimulationClock) TickCount() uint64 {
	return c.tickCount
}

// New creates a new Scheduler.
func New(tickInterval time.Duration) *Scheduler {
	return &Scheduler{
		devices:     make([]*device.Device, 0),
		models:      make([]models.Model, 0),
		tickInterval: tickInterval,
		clock: SimulationClock{
			elapsed:   0,
			tickCount: 0,
		},
		running: false,
		stopCh:  make(chan struct{}),
	}
}

// AddDevice adds a device to the scheduler.
func (s *Scheduler) AddDevice(d *device.Device) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.devices = append(s.devices, d)
}

// RemoveDevice removes a device from the scheduler.
func (s *Scheduler) RemoveDevice(id device.DeviceID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, d := range s.devices {
		if d.ID() == id {
			s.devices = append(s.devices[:i], s.devices[i+1:]...)
			return
		}
	}
}

// Device returns a device by ID.
func (s *Scheduler) Device(id device.DeviceID) *device.Device {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, d := range s.devices {
		if d.ID() == id {
			return d
		}
	}
	return nil
}

// Devices returns all devices.
func (s *Scheduler) Devices() []*device.Device {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]*device.Device, len(s.devices))
	copy(result, s.devices)
	return result
}

// AddModel adds a simulation model to the scheduler.
func (s *Scheduler) AddModel(m models.Model) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.models = append(s.models, m)
}

// RemoveModel removes a simulation model from the scheduler.
func (s *Scheduler) RemoveModel(id models.ModelID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, m := range s.models {
		if m.ID() == id {
			s.models = append(s.models[:i], s.models[i+1:]...)
			return
		}
	}
}

// Model returns a model by ID.
func (s *Scheduler) Model(id models.ModelID) models.Model {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, m := range s.models {
		if m.ID() == id {
			return m
		}
	}
	return nil
}

// Models returns all simulation models.
func (s *Scheduler) Models() []models.Model {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]models.Model, len(s.models))
	copy(result, s.models)
	return result
}

// Run starts the scheduler.
// It ticks devices at the configured interval until the context is cancelled.
func (s *Scheduler) Run(ctx context.Context) error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return nil
	}
	s.running = true
	s.mu.Unlock()

	ticker := time.NewTicker(s.tickInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.mu.Lock()
			s.running = false
			s.mu.Unlock()
			return ctx.Err()
		case <-s.stopCh:
			s.mu.Lock()
			s.running = false
			s.mu.Unlock()
			return nil
		case <-ticker.C:
			s.tick()
		}
	}
}

// Stop stops the scheduler.
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	select {
	case <-s.stopCh:
	default:
		close(s.stopCh)
	}
}

// Pause pauses the scheduler.
func (s *Scheduler) Pause() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.running = false
}

// Resume resumes the scheduler.
func (s *Scheduler) Resume() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.running = true
}

// Running returns true if the scheduler is running.
func (s *Scheduler) Running() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}

// tick executes one simulation tick.
// Models tick first, then devices tick.
func (s *Scheduler) tick() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	devices := make([]*device.Device, len(s.devices))
	copy(devices, s.devices)
	models := make([]models.Model, len(s.models))
	copy(models, s.models)
	s.mu.Unlock()

	// 1. Devices read models and SET power injections/withdrawals
	//    (Devices sample current model state and report power flow)
	for _, d := range devices {
		if d.Running() {
			d.Tick()
		}
	}

	// 2. Models calculate new state based on power
	//    (Bus voltages update based on P/Q injections)
	for _, m := range models {
		m.Tick()
	}

	// 3. Advance the clock
	s.mu.Lock()
	s.clock.elapsed += s.tickInterval
	s.clock.tickCount++
	s.mu.Unlock()
}

// Clock returns the current simulation clock.
func (s *Scheduler) Clock() SimulationClock {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.clock
}

// TickInterval returns the configured tick interval.
func (s *Scheduler) TickInterval() time.Duration {
	return s.tickInterval
}
