// Package examples demonstrates usage of the forge runtime.
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/tamzrod/forge/device"
	"github.com/tamzrod/forge/runtime"
)

func main() {
	// Create runtime with 250ms tick interval
	cfg := runtime.Config{
		TickInterval: 250 * time.Millisecond,
		MaxDevices:   100,
	}
	rt := runtime.New(cfg)

	// Create a revenue meter device
	memRegions := map[string]uint32{
		"holding_registers": 100,
		"input_registers":   200,
		"coils":            20,
		"discrete_inputs":  20,
	}
	meter := rt.CreateDevice("meter-001", "revenue_meter", memRegions)

	// Add a behavior that measures power
	meter.AddBehavior(&powerMeasurement{
		device:    meter,
		voltage:   230.0,
		current:   10.0,
	})

	// Start the simulation
	fmt.Println("Starting simulation...")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- rt.Run(ctx)
	}()

	// Run for 2 seconds
	select {
	case <-ctx.Done():
		fmt.Println("Context cancelled")
	case err := <-done:
		if err != nil {
			fmt.Printf("Runtime error: %v\n", err)
		}
	}

	// Shutdown
	rt.Shutdown()
	fmt.Println("Simulation stopped")

	// Print final memory values
	v, _ := meter.Memory().ReadFloat32("input_registers", 0) // voltage
	p, _ := meter.Memory().ReadFloat32("input_registers", 4) // power
	fmt.Printf("Final values - Voltage: %.2f V, Power: %.2f W\n", v, p)
}

// powerMeasurement is a simple behavior that measures power.
type powerMeasurement struct {
	device  *device.Device
	voltage float32
	current float32
}

func (b *powerMeasurement) ID() string { return "power_measurement" }

func (b *powerMeasurement) Attach(d *device.Device) {
	b.device = d
}

func (b *powerMeasurement) Detach() {
	b.device = nil
}

func (b *powerMeasurement) Tick() {
	// Read from holding registers (configurable)
	scaleFactor, _ := b.device.Memory().ReadFloat32("holding_registers", 0)
	if scaleFactor == 0 {
		scaleFactor = 1.0
	}

	// Write to input registers (measurements)
	b.device.Memory().WriteFloat32("input_registers", 0, b.voltage)
	b.device.Memory().WriteFloat32("input_registers", 4, b.current)

	power := b.voltage * b.current * scaleFactor
	b.device.Memory().WriteFloat32("input_registers", 8, power)
}
