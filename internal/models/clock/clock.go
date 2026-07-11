// Package clock provides the simulation clock.
//
// The Simulation Clock is the single source of time for the entire simulation.
// All models must use this clock—never call time.Now() from simulation code.
//
// Modes:
//   - Realtime: Advances at real-time speed
//   - Manual: Advances only when explicitly stepped
//   - Accelerated: Advances faster than realtime (future)
package clock

import (
	"context"
	"sync"
	"time"
)

// Mode defines how the clock advances time.
type Mode int

const (
	ModeRealtime Mode = iota
	ModeManual
	ModeAccelerated
)

func (m Mode) String() string {
	switch m {
	case ModeRealtime:
		return "Realtime"
	case ModeManual:
		return "Manual"
	case ModeAccelerated:
		return "Accelerated"
	default:
		return "Unknown"
	}
}

// Clock provides deterministic simulation time.
// All simulation models must use this clock instead of system time.
type Clock struct {
	mu sync.RWMutex

	// Configuration
	mode           Mode
	rate           float64 // Time multiplier for accelerated mode
	tickInterval   time.Duration

	// State
	elapsed        time.Duration
	tickCount      uint64
	startTime      time.Time
	paused         bool

	// Manual mode state
	manualElapsed  time.Duration
}

// Config holds clock configuration.
type Config struct {
	// Mode determines how the clock advances.
	Mode Mode

	// Rate is the time multiplier for accelerated mode.
	// For example, 2.0 means 2x realtime.
	Rate float64

	// TickInterval is the interval between ticks in realtime mode.
	TickInterval time.Duration
}

// DefaultConfig returns the default configuration.
func DefaultConfig() Config {
	return Config{
		Mode:        ModeRealtime,
		Rate:        1.0,
		TickInterval: 100 * time.Millisecond,
	}
}

// New creates a new simulation clock.
func New(cfg Config) *Clock {
	return &Clock{
		mode:         cfg.Mode,
		rate:         cfg.Rate,
		tickInterval: cfg.TickInterval,
		elapsed:      0,
		tickCount:    0,
		paused:       false,
	}
}

// Elapsed returns the total elapsed simulation time.
func (c *Clock) Elapsed() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.elapsed
}

// TickCount returns the number of simulation ticks.
func (c *Clock) TickCount() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.tickCount
}

// Tick advances the clock by one tick interval.
// This is the deterministic tick for simulation time advancement.
func (c *Clock) Tick() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.paused {
		return
	}

	switch c.mode {
	case ModeManual:
		// Manual mode: advance by tick interval
		c.elapsed += c.tickInterval
	case ModeRealtime, ModeAccelerated:
		// For realtime/accelerated, this is called by the ticker
		// In manual stepping, this advances simulation time
		c.elapsed += c.tickInterval
	}

	c.tickCount++
}

// Advance advances the clock by a specific duration.
// This is used in manual mode for explicit time stepping.
func (c *Clock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.elapsed += d
	c.tickCount++
}

// Pause pauses the clock.
func (c *Clock) Pause() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.paused = true
}

// Resume resumes the clock.
func (c *Clock) Resume() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.paused = false
}

// IsPaused returns whether the clock is paused.
func (c *Clock) IsPaused() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.paused
}

// Mode returns the current clock mode.
func (c *Clock) Mode() Mode {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.mode
}

// SetMode changes the clock mode.
func (c *Clock) SetMode(mode Mode) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.mode = mode
}

// SetRate sets the time multiplier for accelerated mode.
func (c *Clock) SetRate(rate float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.rate = rate
}

// Rate returns the current time multiplier.
func (c *Clock) Rate() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.rate
}

// Reset resets the clock to zero.
func (c *Clock) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.elapsed = 0
	c.tickCount = 0
	c.paused = false
}

// Now returns the current wall-clock time.
// This should only be used for non-simulation purposes.
func Now() time.Time {
	return time.Now()
}

// RealtimeTicker runs the clock in realtime mode.
// It advances the clock at real-time speed until ctx is cancelled.
func RealtimeTicker(ctx context.Context, clock *Clock) error {
	ticker := time.NewTicker(clock.tickInterval)
	defer ticker.Stop()

	clock.mu.Lock()
	clock.startTime = time.Now()
	clock.mu.Unlock()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			clock.Tick()
		}
	}
}
