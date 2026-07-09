package runtime

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/tamzrod/forge/device"
)

func TestNew(t *testing.T) {
	cfg := Config{
		TickInterval: 100 * time.Millisecond,
		MaxDevices:   100,
	}

	r := New(cfg)

	if r.config.TickInterval != 100*time.Millisecond {
		t.Errorf("expected tick interval 100ms, got %v", r.config.TickInterval)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.TickInterval != 100*time.Millisecond {
		t.Errorf("expected default tick interval 100ms, got %v", cfg.TickInterval)
	}
	if cfg.MaxDevices != 1000 {
		t.Errorf("expected default max devices 1000, got %d", cfg.MaxDevices)
	}
}

func TestLoadConfig(t *testing.T) {
	// Create temp config file
	content := `
tick_interval: 250ms
max_devices: 500
`
	tmpfile, err := os.CreateTemp("", "config*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.TickInterval != 250*time.Millisecond {
		t.Errorf("expected tick interval 250ms, got %v", cfg.TickInterval)
	}
	if cfg.MaxDevices != 500 {
		t.Errorf("expected max devices 500, got %d", cfg.MaxDevices)
	}
}

func TestCreateDevice(t *testing.T) {
	r := New(DefaultConfig())

	memRegions := map[string]uint32{
		"holding_registers": 100,
		"input_registers":  200,
	}

	d := r.CreateDevice("meter-001", "revenue_meter", memRegions)

	if d == nil {
		t.Fatal("expected device, got nil")
	}
	if d.ID() != "meter-001" {
		t.Errorf("expected ID 'meter-001', got '%s'", d.ID())
	}

	// Check device is in registry
	found := r.Device("meter-001")
	if found != d {
		t.Error("device not found in registry")
	}
}

func TestDevices(t *testing.T) {
	r := New(DefaultConfig())

	r.CreateDevice("d1", "test", map[string]uint32{"input_registers": 10})
	r.CreateDevice("d2", "test", map[string]uint32{"input_registers": 10})
	r.CreateDevice("d3", "test", map[string]uint32{"input_registers": 10})

	devices := r.Devices()
	if len(devices) != 3 {
		t.Errorf("expected 3 devices, got %d", len(devices))
	}
}

func TestStartStop(t *testing.T) {
	r := New(DefaultConfig())

	d := r.CreateDevice("test", "test", map[string]uint32{"input_registers": 10})

	if d.Running() {
		t.Error("device should not be running initially")
	}

	r.Start()

	// Give scheduler a moment
	time.Sleep(time.Millisecond)

	if !d.Running() {
		t.Error("device should be running after Start")
	}

	r.Stop()

	if d.Running() {
		t.Error("device should not be running after Stop")
	}
}

func TestRun(t *testing.T) {
	r := New(Config{TickInterval: time.Millisecond})

	d := r.CreateDevice("test", "test", map[string]uint32{"input_registers": 10})

	var tickCount int
	d.AddBehavior(&testBehavior{
		id:   "test",
		tick: func() { tickCount++ },
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- r.Run(ctx)
	}()

	select {
	case <-ctx.Done():
		r.Shutdown()
	case err := <-done:
		if err != nil && err != context.DeadlineExceeded {
			t.Errorf("unexpected error: %v", err)
		}
	}

	if tickCount < 2 {
		t.Errorf("expected at least 2 ticks, got %d", tickCount)
	}
}

func TestShutdown(t *testing.T) {
	r := New(DefaultConfig())

	r.CreateDevice("test", "test", map[string]uint32{"input_registers": 10})

	err := r.Shutdown()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSchedulerAccess(t *testing.T) {
	r := New(DefaultConfig())

	sched := r.Scheduler()
	if sched == nil {
		t.Error("expected scheduler, got nil")
	}

	if sched.TickInterval() != 100*time.Millisecond {
		t.Errorf("expected tick interval 100ms, got %v", sched.TickInterval())
	}
}

// testBehavior is a simple behavior implementation for testing.
type testBehavior struct {
	id   string
	tick func()
}

func (b *testBehavior) ID() string { return b.id }
func (b *testBehavior) Attach(d *device.Device) {
}
func (b *testBehavior) Detach() {}
func (b *testBehavior) Tick()   { b.tick() }
