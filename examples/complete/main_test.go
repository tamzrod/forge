package main

import (
	"context"
	"testing"
	"time"

	"github.com/tamzrod/forge/device"
	"github.com/tamzrod/forge/models"
	"github.com/tamzrod/forge/runtime"
)

// TestCompleteExample_RuntimeInitialization verifies that the runtime
// can be initialized with models and devices.
func TestCompleteExample_RuntimeInitialization(t *testing.T) {
	cfg := runtime.Config{
		TickInterval: 100 * time.Millisecond,
		MaxDevices:   100,
	}
	rt := runtime.New(cfg)

	// Create models
	rt.CreateGridModel("test-grid")
	rt.CreateSunModel("test-sun")
	rt.CreateWeatherModel("test-weather")

	// Create device
	memRegions := map[string]uint32{
		"input_registers": 10,
	}
	device := rt.CreateDevice("test-device", "test_type", memRegions)
	device.AddBehavior(&testBehavior{})

	// Verify models exist
	if rt.Model("test-grid") == nil {
		t.Error("expected grid model to exist")
	}
	if rt.Model("test-sun") == nil {
		t.Error("expected sun model to exist")
	}
	if rt.Model("test-weather") == nil {
		t.Error("expected weather model to exist")
	}

	// Verify device exists
	if rt.Device("test-device") == nil {
		t.Error("expected device to exist")
	}
}

// TestCompleteExample_ExecutionLoop verifies that the simulation
// runs for the expected number of ticks.
func TestCompleteExample_ExecutionLoop(t *testing.T) {
	cfg := runtime.Config{
		TickInterval: 10 * time.Millisecond,
		MaxDevices:   100,
	}
	rt := runtime.New(cfg)

	// Create model
	rt.CreateGridModel("test-grid")

	// Create device with counting behavior
	memRegions := map[string]uint32{
		"input_registers": 10,
	}
	device := rt.CreateDevice("test-device", "test_type", memRegions)
	var tickCount int
	device.AddBehavior(&countingBehavior{count: &tickCount})

	// Run simulation
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- rt.Run(ctx)
	}()

	select {
	case <-ctx.Done():
		// Timeout reached
	case err := <-done:
		if err != nil && err != context.DeadlineExceeded {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	rt.Shutdown()

	// Verify behavior was called
	if tickCount < 3 {
		t.Errorf("expected at least 3 ticks, got %d", tickCount)
	}
}

// TestCompleteExample_ModelObservation verifies that behaviors can
// observe models through the ModelProvider interface.
func TestCompleteExample_ModelObservation(t *testing.T) {
	cfg := runtime.Config{
		TickInterval: 10 * time.Millisecond,
		MaxDevices:   100,
	}
	rt := runtime.New(cfg)

	// Create model with specific state
	grid := rt.CreateGridModel("test-grid")
	grid.SetVoltage(480.0)

	// Create device with observation behavior
	memRegions := map[string]uint32{
		"input_registers": 10,
	}
	device := rt.CreateDevice("test-device", "test_type", memRegions)
	device.AddBehavior(&observationBehavior{})

	// Run for one tick
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- rt.Run(ctx)
	}()

	select {
	case <-ctx.Done():
		// Timeout reached
	case err := <-done:
		if err != nil && err != context.DeadlineExceeded {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	rt.Shutdown()

	// Verify the behavior observed the model
	observedVoltage, _ := device.Memory().ReadFloat32("input_registers", 0)
	if observedVoltage != 480.0 {
		t.Errorf("expected observed voltage 480.0, got %f", observedVoltage)
	}
}

// testBehavior is a simple behavior for testing.
type testBehavior struct{}

func (b *testBehavior) ID() string { return "test" }
func (b *testBehavior) Attach(d *device.Device) {}
func (b *testBehavior) Detach()         {}
func (b *testBehavior) Tick()           {}

// countingBehavior counts the number of ticks.
type countingBehavior struct {
	count *int
}

func (b *countingBehavior) ID() string { return "counting" }
func (b *countingBehavior) Attach(d *device.Device) {}
func (b *countingBehavior) Detach()         {}
func (b *countingBehavior) Tick() {
	*b.count++
}

// observationBehavior observes a model and writes to memory.
type observationBehavior struct {
	device *device.Device
}

func (b *observationBehavior) ID() string { return "observation" }
func (b *observationBehavior) Attach(d *device.Device) {
	b.device = d
}
func (b *observationBehavior) Detach() {}
func (b *observationBehavior) Tick() {
	// Observe the grid model
	grid := b.device.Model("test-grid")
	if grid == nil {
		return
	}

	gridModel, ok := grid.(*models.GridModel)
	if !ok {
		return
	}

	// Read grid state and write to memory
	voltage := gridModel.Voltage()
	b.device.Memory().WriteFloat32("input_registers", 0, voltage)
}
