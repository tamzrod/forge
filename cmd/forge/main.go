// Package main is the entry point for the forge simulation runtime.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tamzrod/forge/device"
	"github.com/tamzrod/forge/models"
	"github.com/tamzrod/forge/runtime"
)

var (
	flagTickInterval = flag.Duration("tick", 100*time.Millisecond, "Simulation tick interval")
	flagMaxDevices  = flag.Int("devices", 100, "Maximum number of devices")
	flagDuration    = flag.Duration("duration", 0, "Run duration (0 = indefinite)")
	flagVerbose     = flag.Bool("v", false, "Verbose output")
)

func main() {
	flag.Parse()

	if *flagVerbose {
		fmt.Println("Forge Industrial Simulation Runtime")
		fmt.Printf("Tick Interval: %v\n", *flagTickInterval)
		fmt.Printf("Max Devices: %d\n", *flagMaxDevices)
		fmt.Println()
	}

	// Create runtime configuration
	cfg := runtime.Config{
		TickInterval: *flagTickInterval,
		MaxDevices:   *flagMaxDevices,
	}

	// Create runtime
	rt := runtime.New(cfg)

	// Create simulation models
	grid := rt.CreateGridModel("main-grid")
	_ = rt.CreateSunModel("solar-sun")
	_ = rt.CreateWindModel("wind-farm")
	_ = rt.CreateWeatherModel("ambient-weather")

	if *flagVerbose {
		fmt.Println("Created models:")
		for _, m := range rt.Models() {
			fmt.Printf("  - %s (%s)\n", m.ID(), m.Type())
		}
		fmt.Println()
	}

	// Create a sample weather station device
	memRegions := map[string]uint32{
		"sensors":    64, // Temperature, humidity, etc.
		"computed":   64, // Calculated values
		"status":     16, // Device status
	}
	weatherStation := rt.CreateDevice("ws-001", "weather_station", memRegions)

	// Add behavior to observe weather model
	weatherStation.AddBehavior(&weatherBehavior{
		device: weatherStation,
	})

	// Create a sample PV inverter device
	memRegions = map[string]uint32{
		"input":   32,  // DC input
		"output":  32,  // AC output
		"config":  32,  // Configuration
		"status": 16,  // Status flags
	}
	pvInverter := rt.CreateDevice("pv-001", "pv_inverter", memRegions)

	// Add behavior to calculate power from sun model
	pvInverter.AddBehavior(&pvBehavior{
		device: pvInverter,
	})

	// Create a sample revenue meter
	memRegions = map[string]uint32{
		"holding_registers": 100,
		"input_registers":   200,
		"coils":            20,
		"discrete_inputs":   20,
	}
	meter := rt.CreateDevice("meter-001", "revenue_meter", memRegions)

	// Add power measurement behavior
	meter.AddBehavior(&powerMeasurement{
		device:    meter,
		voltage:   230.0,
		current:   0.0,
	})

	if *flagVerbose {
		fmt.Println("Created devices:")
		for _, d := range rt.Devices() {
			fmt.Printf("  - %s (%s)\n", d.ID(), d.Type())
		}
		fmt.Println()
	}

	// Setup context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\nShutting down...")
		cancel()
	}()

	// Start runtime in background
	done := make(chan error, 1)
	go func() {
		done <- rt.Run(ctx)
	}()

	// Print tick information
	tickCount := 0
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	fmt.Println("Simulation running...")
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println()

	// Run for specified duration or until cancelled
	var runDuration time.Duration
	if *flagDuration > 0 {
		runDuration = *flagDuration
	} else {
		runDuration = 10 * time.Second // Default 10 seconds
	}
	timeout := time.After(runDuration)

	for {
		select {
		case <-timeout:
			fmt.Println("\nDuration reached, stopping...")
			cancel()
		case <-ctx.Done():
			goto shutdown
		case <-ticker.C:
			tickCount++
			if *flagVerbose {
				clock := rt.Scheduler().Clock()
				fmt.Printf("[Tick %d] elapsed=%v\n", tickCount, clock.elapsed)
				fmt.Printf("  Grid: %.1fV @ %.2fHz\n", grid.Voltage(), grid.Frequency())
			} else {
				fmt.Printf("\rRunning... Tick %d", tickCount)
			}
		}
	}

shutdown:
	// Wait for runtime to stop
	<-done

	// Shutdown
	rt.Shutdown()

	fmt.Println("\nSimulation complete")
}

// weatherBehavior samples the weather model and updates device memory.
type weatherBehavior struct {
	device *device.Device
}

func (b *weatherBehavior) ID() string { return "weather_sampling" }
func (b *weatherBehavior) Attach(d *device.Device) { b.device = d }
func (b *weatherBehavior) Detach() { b.device = nil }

func (b *weatherBehavior) Tick() {
	if b.device == nil {
		return
	}

	// Get weather model
	weatherModel := b.device.Model("ambient-weather")
	if weatherModel == nil {
		return
	}

	wm, ok := weatherModel.(*models.WeatherModel)
	if !ok {
		return
	}

	// Sample weather into device memory
	b.device.Memory().WriteFloat32("sensors", 0, wm.Temperature())
	b.device.Memory().WriteFloat32("sensors", 4, wm.Humidity())
	b.device.Memory().WriteFloat32("sensors", 8, wm.Pressure())
}

// pvBehavior calculates PV output based on sun model.
type pvBehavior struct {
	device *device.Device
}

func (b *pvBehavior) ID() string { return "pv_power_calc" }
func (b *pvBehavior) Attach(d *device.Device) { b.device = d }
func (b *pvBehavior) Detach() { b.device = nil }

func (b *pvBehavior) Tick() {
	if b.device == nil {
		return
	}

	// Get sun model
	sunModel := b.device.Model("solar-sun")
	if sunModel == nil {
		return
	}

	sm, ok := sunModel.(*models.SunModel)
	if !ok {
		return
	}

	// Calculate power based on irradiance
	irradiance := sm.Irradiance()
	power := irradiance * 0.15 // Simplified efficiency

	b.device.Memory().WriteFloat32("input", 0, irradiance)
	b.device.Memory().WriteFloat32("output", 0, power)
}

// powerMeasurement measures power and updates memory.
type powerMeasurement struct {
	device  *device.Device
	voltage float32
	current float32
}

func (b *powerMeasurement) ID() string { return "power_measurement" }
func (b *powerMeasurement) Attach(d *device.Device) { b.device = d }
func (b *powerMeasurement) Detach() { b.device = nil }

func (b *powerMeasurement) Tick() {
	if b.device == nil {
		return
	}

	// Read power from PV inverter if available via model
	sunModel := b.device.Model("solar-sun")
	if sunModel != nil {
		if sm, ok := sunModel.(*models.SunModel); ok {
			irradiance := sm.Irradiance()
			power := irradiance * 0.15
			b.current = power / b.voltage
		}
	}

	// Calculate power
	power := b.voltage * b.current

	// Write to memory
	b.device.Memory().WriteFloat32("input_registers", 0, b.voltage)
	b.device.Memory().WriteFloat32("input_registers", 4, b.current)
	b.device.Memory().WriteFloat32("input_registers", 8, power)
}
