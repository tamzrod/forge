package clock

import (
	"testing"
	"time"
)

func TestClock_Creation(t *testing.T) {
	clock := New(DefaultConfig())

	if clock.Elapsed() != 0 {
		t.Errorf("expected elapsed 0, got %v", clock.Elapsed())
	}

	if clock.TickCount() != 0 {
		t.Errorf("expected tick count 0, got %d", clock.TickCount())
	}

	if clock.IsPaused() {
		t.Error("expected clock to not be paused")
	}
}

func TestClock_Tick(t *testing.T) {
	cfg := Config{
		Mode:        ModeManual,
		TickInterval: 100 * time.Millisecond,
	}
	clock := New(cfg)

	// Tick should advance time
	clock.Tick()

	if clock.Elapsed() != 100*time.Millisecond {
		t.Errorf("expected elapsed 100ms, got %v", clock.Elapsed())
	}

	if clock.TickCount() != 1 {
		t.Errorf("expected tick count 1, got %d", clock.TickCount())
	}

	// Multiple ticks
	clock.Tick()
	clock.Tick()
	clock.Tick()

	if clock.Elapsed() != 400*time.Millisecond {
		t.Errorf("expected elapsed 400ms, got %v", clock.Elapsed())
	}

	if clock.TickCount() != 4 {
		t.Errorf("expected tick count 4, got %d", clock.TickCount())
	}
}

func TestClock_PauseResume(t *testing.T) {
	cfg := Config{
		Mode:        ModeManual,
		TickInterval: 100 * time.Millisecond,
	}
	clock := New(cfg)

	// Advance some time
	clock.Tick()
	clock.Tick()

	// Pause
	clock.Pause()
	if !clock.IsPaused() {
		t.Error("expected clock to be paused")
	}

	// Ticks should not advance when paused
	elapsedBefore := clock.Elapsed()
	clock.Tick()
	clock.Tick()
	if clock.Elapsed() != elapsedBefore {
		t.Errorf("expected elapsed not to change when paused, got %v", clock.Elapsed())
	}

	// Resume
	clock.Resume()
	if clock.IsPaused() {
		t.Error("expected clock to not be paused")
	}

	// Ticks should advance again
	clock.Tick()
	if clock.Elapsed() == elapsedBefore {
		t.Error("expected elapsed to advance after resume")
	}
}

func TestClock_Advance(t *testing.T) {
	cfg := Config{
		Mode:        ModeManual,
		TickInterval: 100 * time.Millisecond,
	}
	clock := New(cfg)

	// Advance by specific duration
	clock.Advance(500 * time.Millisecond)

	if clock.Elapsed() != 500*time.Millisecond {
		t.Errorf("expected elapsed 500ms, got %v", clock.Elapsed())
	}

	if clock.TickCount() != 1 {
		t.Errorf("expected tick count 1, got %d", clock.TickCount())
	}
}

func TestClock_Reset(t *testing.T) {
	cfg := Config{
		Mode:        ModeManual,
		TickInterval: 100 * time.Millisecond,
	}
	clock := New(cfg)

	// Advance time
	clock.Tick()
	clock.Tick()
	clock.Tick()

	// Reset
	clock.Reset()

	if clock.Elapsed() != 0 {
		t.Errorf("expected elapsed 0 after reset, got %v", clock.Elapsed())
	}

	if clock.TickCount() != 0 {
		t.Errorf("expected tick count 0 after reset, got %d", clock.TickCount())
	}

	if clock.IsPaused() {
		t.Error("expected clock to not be paused after reset")
	}
}

func TestClock_Mode(t *testing.T) {
	clock := New(DefaultConfig())

	if clock.Mode() != ModeRealtime {
		t.Errorf("expected mode Realtime, got %v", clock.Mode())
	}

	clock.SetMode(ModeManual)
	if clock.Mode() != ModeManual {
		t.Errorf("expected mode Manual, got %v", clock.Mode())
	}

	clock.SetMode(ModeAccelerated)
	if clock.Mode() != ModeAccelerated {
		t.Errorf("expected mode Accelerated, got %v", clock.Mode())
	}
}

func TestClock_Rate(t *testing.T) {
	cfg := Config{
		Mode:        ModeAccelerated,
		Rate:        2.0,
		TickInterval: 100 * time.Millisecond,
	}
	clock := New(cfg)

	if clock.Rate() != 2.0 {
		t.Errorf("expected rate 2.0, got %f", clock.Rate())
	}

	clock.SetRate(5.0)
	if clock.Rate() != 5.0 {
		t.Errorf("expected rate 5.0, got %f", clock.Rate())
	}
}

func TestClock_Deterministic(t *testing.T) {
	// Two clocks with the same initial state should tick identically
	cfg := Config{
		Mode:        ModeManual,
		TickInterval: 100 * time.Millisecond,
	}

	clock1 := New(cfg)
	clock2 := New(cfg)

	// Perform same operations
	for i := 0; i < 100; i++ {
		clock1.Tick()
		clock2.Tick()
	}

	// Results should be identical
	if clock1.Elapsed() != clock2.Elapsed() {
		t.Errorf("elapsed mismatch: %v vs %v", clock1.Elapsed(), clock2.Elapsed())
	}

	if clock1.TickCount() != clock2.TickCount() {
		t.Errorf("tick count mismatch: %d vs %d", clock1.TickCount(), clock2.TickCount())
	}
}
