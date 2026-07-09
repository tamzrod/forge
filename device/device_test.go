package device

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	memRegions := map[string]uint32{
		"holding_registers": 100,
		"input_registers":  200,
	}

	d := New("meter-001", "revenue_meter", memRegions)

	if d.ID() != "meter-001" {
		t.Errorf("expected ID 'meter-001', got '%s'", d.ID())
	}
	if d.Type() != "revenue_meter" {
		t.Errorf("expected Type 'revenue_meter', got '%s'", d.Type())
	}
	if d.Memory() == nil {
		t.Error("expected memory to be non-nil")
	}
	if d.Running() {
		t.Error("expected device to not be running initially")
	}
}

func TestAddBehavior(t *testing.T) {
	d := New("test", "test", map[string]uint32{"input_registers": 10})

	var tickCount int
	b := &testBehavior{
		id:   "test",
		tick: func() { tickCount++ },
	}

	d.AddBehavior(b)

	if len(d.Behaviors()) != 1 {
		t.Errorf("expected 1 behavior, got %d", len(d.Behaviors()))
	}
}

func TestTick(t *testing.T) {
	d := New("test", "test", map[string]uint32{"input_registers": 10})

	var tickCount int
	b := &testBehavior{
		id:   "test",
		tick: func() { tickCount++ },
	}

	d.AddBehavior(b)
	d.Tick()

	if tickCount != 1 {
		t.Errorf("expected tick count 1, got %d", tickCount)
	}

	d.Tick()
	if tickCount != 2 {
		t.Errorf("expected tick count 2, got %d", tickCount)
	}
}

func TestBehaviorAttach(t *testing.T) {
	d := New("test", "test", map[string]uint32{"input_registers": 10})

	var attachedDevice *Device
	b := &testBehavior{
		id: "test",
		attachFn: func(dev *Device) {
			attachedDevice = dev
		},
	}

	d.AddBehavior(b)

	if attachedDevice != d {
		t.Error("behavior was not attached to device")
	}
}

func TestStartStop(t *testing.T) {
	d := New("test", "test", map[string]uint32{"input_registers": 10})

	if d.Running() {
		t.Error("device should not be running initially")
	}

	d.Start()
	if !d.Running() {
		t.Error("device should be running after Start")
	}

	d.Stop()
	if d.Running() {
		t.Error("device should not be running after Stop")
	}
}

func TestMemoryAccess(t *testing.T) {
	memRegions := map[string]uint32{
		"input_registers": 10,
	}
	d := New("test", "test", memRegions)

	// Write to memory
	d.Memory().Write("input_registers", 0, []byte{0x12, 0x34})

	// Read from memory
	val, err := d.Memory().ReadUint16("input_registers", 0)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}

	if val != 0x1234 {
		t.Errorf("expected 0x1234, got 0x%04x", val)
	}
}

func TestConcurrentTick(t *testing.T) {
	d := New("test", "test", map[string]uint32{"input_registers": 10})

	var tickCount int
	b := &testBehavior{
		id: "test",
		tick: func() {
			time.Sleep(time.Millisecond)
			tickCount++
		},
	}

	d.AddBehavior(b)

	// Tick concurrently
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			d.Tick()
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	// Each tick increments by 1
	if tickCount != 10 {
		t.Errorf("expected tick count 10, got %d", tickCount)
	}
}

// testBehavior is a simple behavior implementation for testing.
type testBehavior struct {
	id       string
	attachFn func(*Device)
	detachFn func()
	tick     func()
}

func (b *testBehavior) ID() string              { return b.id }
func (b *testBehavior) Attach(d *Device)          { b.attachFn(d) }
func (b *testBehavior) Detach()                   { b.detachFn() }
func (b *testBehavior) Tick()                     { b.tick() }
