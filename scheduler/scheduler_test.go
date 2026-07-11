package scheduler

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/tamzrod/forge/device"
)

func TestNew(t *testing.T) {
	s := New(100 * time.Millisecond)

	if s.tickInterval != 100*time.Millisecond {
		t.Errorf("expected tick interval 100ms, got %v", s.tickInterval)
	}
	if s.Running() {
		t.Error("scheduler should not be running initially")
	}
}

func TestAddDevice(t *testing.T) {
	s := New(time.Millisecond)
	d := device.New("test", "test", map[string]uint32{"input_registers": 10})

	s.AddDevice(d)

	if len(s.Devices()) != 1 {
		t.Errorf("expected 1 device, got %d", len(s.Devices()))
	}
}

func TestRemoveDevice(t *testing.T) {
	s := New(time.Millisecond)
	d := device.New("test", "test", map[string]uint32{"input_registers": 10})

	s.AddDevice(d)
	s.RemoveDevice("test")

	if len(s.Devices()) != 0 {
		t.Errorf("expected 0 devices, got %d", len(s.Devices()))
	}
}

func TestDeviceLookup(t *testing.T) {
	s := New(time.Millisecond)
	d := device.New("test", "test", map[string]uint32{"input_registers": 10})

	s.AddDevice(d)

	found := s.Device("test")
	if found != d {
		t.Error("device not found")
	}

	notFound := s.Device("nonexistent")
	if notFound != nil {
		t.Error("nonexistent device should return nil")
	}
}

func TestTickExecution(t *testing.T) {
	s := New(time.Millisecond)
	d := device.New("test", "test", map[string]uint32{"input_registers": 10})

	var tickCount int
	b := &testBehavior{
		id:   "test",
		tick: func() { tickCount++ },
	}
	d.AddBehavior(b)
	d.Start()

	s.AddDevice(d)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	go func() {
		s.Run(ctx)
	}()

	// Wait for some ticks
	time.Sleep(4 * time.Millisecond)

	if tickCount < 2 {
		t.Errorf("expected at least 2 ticks, got %d", tickCount)
	}
}

func TestPauseResume(t *testing.T) {
	s := New(time.Millisecond)
	d := device.New("test", "test", map[string]uint32{"input_registers": 10})

	var tickCount int
	b := &testBehavior{
		id:   "test",
		tick: func() { tickCount++ },
	}
	d.AddBehavior(b)
	d.Start()

	s.AddDevice(d)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	go func() {
		s.Run(ctx)
	}()

	time.Sleep(3 * time.Millisecond)

	pauseCount := tickCount

	s.Pause()
	time.Sleep(5 * time.Millisecond)

	if tickCount != pauseCount {
		t.Errorf("tick count should not change when paused, was %d, still %d", pauseCount, tickCount)
	}

	s.Resume()
	time.Sleep(3 * time.Millisecond)

	if tickCount <= pauseCount {
		t.Errorf("tick count should increase after resume, was %d, now %d", pauseCount, tickCount)
	}
}

func TestDeviceNotRunning(t *testing.T) {
	s := New(time.Millisecond)
	d := device.New("test", "test", map[string]uint32{"input_registers": 10})

	// Device not started
	var tickCount int
	b := &testBehavior{
		id:   "test",
		tick: func() { tickCount++ },
	}
	d.AddBehavior(b)

	s.AddDevice(d)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	go func() {
		s.Run(ctx)
	}()

	time.Sleep(3 * time.Millisecond)

	// Device is not running, so no ticks should occur
	if tickCount != 0 {
		t.Errorf("expected 0 ticks (device not started), got %d", tickCount)
	}
}

func TestClockAdvance(t *testing.T) {
	s := New(time.Millisecond)

	initialClock := s.Clock()
	if initialClock.elapsed != 0 {
		t.Errorf("expected 0 elapsed time, got %v", initialClock.elapsed)
	}
	if initialClock.tickCount != 0 {
		t.Errorf("expected 0 tick count, got %d", initialClock.tickCount)
	}
}

func TestConcurrentAccess(t *testing.T) {
	s := New(time.Millisecond)

	var wg sync.WaitGroup

	// Concurrent adds
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			d := device.New(device.DeviceID(string(rune('a'+id))), "test", map[string]uint32{"input_registers": 10})
			s.AddDevice(d)
		}(i)
	}

	wg.Wait()

	if len(s.Devices()) != 10 {
		t.Errorf("expected 10 devices, got %d", len(s.Devices()))
	}
}

// testBehavior is a simple behavior implementation for testing.
type testBehavior struct {
	id       string
	attachFn func(*device.Device)
	detachFn func()
	tick     func()
}

func (b *testBehavior) ID() string { return b.id }
func (b *testBehavior) Attach(d *device.Device) {
	if b.attachFn != nil {
		b.attachFn(d)
	}
}
func (b *testBehavior) Detach() {
	if b.detachFn != nil {
		b.detachFn()
	}
}
func (b *testBehavior) Tick() { b.tick() }
