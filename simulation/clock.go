// Package simulation provides simulation timing and orchestration.
// The Simulation Clock is the single authoritative source of time.
package simulation

import (
	"sync"
	"time"
)

// Mode represents the simulation mode.
type Mode int

const (
	// ModeRealtime runs simulation synchronized to wall clock.
	ModeRealtime Mode = iota
	// ModeSimulated runs simulation independently from wall clock.
	ModeSimulated
	// ModeManual advances simulation only when requested.
	ModeManual
	// ModeReplay follows a recorded timeline.
	ModeReplay
)

// String returns the mode name.
func (m Mode) String() string {
	switch m {
	case ModeRealtime:
		return "realtime"
	case ModeSimulated:
		return "simulated"
	case ModeManual:
		return "manual"
	case ModeReplay:
		return "replay"
	default:
		return "unknown"
	}
}

// Clock is the single authoritative source of time for the simulation.
// Nothing inside Forge should call the system clock directly.
type Clock interface {
	// Now returns the current simulation datetime.
	Now() time.Time

	// Elapsed returns the elapsed simulation time since start.
	Elapsed() time.Duration

	// Tick returns the current tick number.
	Tick() uint64

	// StartTime returns when the simulation started.
	StartTime() time.Time

	// Speed returns the simulation speed multiplier.
	Speed() float64

	// SetSpeed changes the simulation speed (1.0 = real-time).
	SetSpeed(speed float64)

	// Mode returns the current simulation mode.
	Mode() Mode

	// IsRunning returns true if the simulation is running.
	IsRunning() bool

	// IsPaused returns true if the simulation is paused.
	IsPaused() bool

	// Start begins the simulation.
	Start(startTime time.Time) error

	// Pause suspends the simulation.
	Pause()

	// Resume continues a paused simulation.
	Resume()

	// Stop halts the simulation.
	Stop()

	// Advance moves the simulation forward by dt (for Manual mode).
	Advance(dt time.Duration)

	// Update advances the clock based on elapsed wall time.
	Update()
}

// SimClock implements Clock with configurable mode and speed.
type SimClock struct {
	mu sync.RWMutex

	// Configuration
	mode  Mode
	speed float64

	// State
	running   bool
	paused    bool
	startTime time.Time
	startWall time.Time
	current   time.Time
	tick      uint64

	// For simulated mode
	lastWall  time.Time
	simOffset time.Duration
}

// NewClock creates a new simulation clock.
func NewClock() *SimClock {
	return &SimClock{
		mode:  ModeRealtime,
		speed: 1.0,
	}
}

// Now returns the current simulation datetime.
func (c *SimClock) Now() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.current
}

// Elapsed returns the elapsed simulation time since start.
func (c *SimClock) Elapsed() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.startTime.IsZero() {
		return 0
	}
	return c.current.Sub(c.startTime)
}

// Tick returns the current tick number.
func (c *SimClock) Tick() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.tick
}

// StartTime returns when the simulation started.
func (c *SimClock) StartTime() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.startTime
}

// Speed returns the simulation speed multiplier.
func (c *SimClock) Speed() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.speed
}

// SetSpeed changes the simulation speed.
func (c *SimClock) SetSpeed(speed float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if speed < 0 {
		speed = 0
	}
	c.speed = speed
}

// Mode returns the current simulation mode.
func (c *SimClock) Mode() Mode {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.mode
}

// SetMode changes the simulation mode.
func (c *SimClock) SetMode(mode Mode) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.mode = mode
}

// IsRunning returns true if the simulation is running.
func (c *SimClock) IsRunning() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.running
}

// IsPaused returns true if the simulation is paused.
func (c *SimClock) IsPaused() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.paused
}

// Start begins the simulation.
func (c *SimClock) Start(startTime time.Time) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return &Error{"simulation already running"}
	}

	c.startTime = startTime
	c.startWall = time.Now()
	c.current = startTime
	c.tick = 0
	c.running = true
	c.paused = false
	c.lastWall = time.Now()
	c.simOffset = 0

	return nil
}

// Pause suspends the simulation.
func (c *SimClock) Pause() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.running && !c.paused {
		c.paused = true
	}
}

// Resume continues a paused simulation.
func (c *SimClock) Resume() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.running && c.paused {
		c.paused = false
		c.lastWall = time.Now()
	}
}

// Stop halts the simulation.
func (c *SimClock) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.running = false
	c.paused = false
}

// Advance moves the simulation forward by dt (for Manual mode).
func (c *SimClock) Advance(dt time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.mode != ModeManual {
		return
	}
	c.tick++
	c.current = c.current.Add(dt)
}

// Update advances the clock based on elapsed wall time.
// Call this once per simulation tick.
func (c *SimClock) Update() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running || c.paused {
		return
	}

	now := time.Now()

	switch c.mode {
	case ModeRealtime:
		// Follow wall clock
		elapsed := now.Sub(c.startWall)
		simElapsed := time.Duration(float64(elapsed) * c.speed)
		c.current = c.startTime.Add(simElapsed)
		c.tick++

	case ModeSimulated:
		// Advance independently
		wallDelta := now.Sub(c.lastWall)
		simDelta := time.Duration(float64(wallDelta) * c.speed)
		c.current = c.current.Add(simDelta)
		c.tick++
		c.lastWall = now

	case ModeManual:
		// Do nothing - Advance() must be called manually

	case ModeReplay:
		// TODO: Follow recorded timeline
	}
}

// Error represents a simulation clock error.
type Error struct {
	msg string
}

func (e *Error) Error() string {
	return e.msg
}

// Predefined speeds for common use cases.
const (
	SpeedRealtime  = 1.0
	SpeedSlow      = 0.5
	SpeedFast      = 2.0
	SpeedVeryFast  = 5.0
	SpeedUltraFast = 10.0
	SpeedExtreme   = 20.0
	SpeedLudicrous = 50.0
	SpeedPlaid     = 100.0
)

// SpeedOptions returns common speed options.
func SpeedOptions() []float64 {
	return []float64{0.1, 0.5, 1.0, 2.0, 5.0, 10.0, 20.0, 50.0, 100.0}
}
