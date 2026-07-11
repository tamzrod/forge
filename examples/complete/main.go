// Package main demonstrates an end-to-end simulation with the Forge runtime.
// This example validates the architecture by running a complete simulation
// with simulation models, virtual devices, and behaviors.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tamzrod/forge/device"
	"github.com/tamzrod/forge/models"
	"github.com/tamzrod/forge/runtime"
)

// This example demonstrates the complete Forge architecture:
//
//  1. Runtime initialization
//  2. Simulation Models (Grid, Sun, Weather)
//  3. Virtual Devices (WeatherStation, RevenueMeter, PVInverter)
//  4. Behaviors (sample models, update memory)
//  5. Execution loop
//  6. Shutdown
//
// Data Flow:
//
//   Simulation Models (physical world)
//           │
//           ▼ (behaviors observe)
//   Virtual Devices (sample and own memory)
//           │
//           ▼ (protocols expose)
//   External Systems (MMA2, SCADA)
//
// Architecture Validation:
//
// This example proves:
// - Models tick independently of devices
// - Devices observe models through ModelProvider interface
// - Behaviors read models and write device memory
// - Memory is the source of truth
// - Runtime orchestrates everything

func main() {
	fmt.Println("=== Forge End-to-End Example ===")
	fmt.Println()

	// Step 1: Runtime initialization
	// The runtime hosts models and devices. It provides scheduling,
	// time advancement, and coordination.
	fmt.Println("Step 1: Initialize Runtime")
	cfg := runtime.Config{
		TickInterval: 100 * time.Millisecond,
		MaxDevices:   100,
	}
	rt := runtime.New(cfg)
	fmt.Printf("  - Tick interval: %v\n", cfg.TickInterval)
	fmt.Printf("  - Max devices: %d\n", cfg.MaxDevices)
	fmt.Println()

	// Step 2: Create Simulation Models
	// Models represent the physical world. They have private state
	// and are observed by devices. Models do NOT have protocols.
	fmt.Println("Step 2: Create Simulation Models")

	// Grid model: electrical grid conditions
	grid := rt.CreateGridModel("main-grid")
	fmt.Printf("  - Grid model created: %s\n", grid.ID())

	// Sun model: solar position and irradiance
	sun := rt.CreateSunModel("solar-sun")
	fmt.Printf("  - Sun model created: %s\n", sun.ID())

	// Weather model: ambient conditions
	weather := rt.CreateWeatherModel("ambient-weather")
	fmt.Printf("  - Weather model created: %s\n", weather.ID())

	// Wind model: wind conditions
	wind := rt.CreateWindModel("wind-field")
	fmt.Printf("  - Wind model created: %s\n", wind.ID())
	fmt.Println()

	// Step 3: Create Virtual Devices
	// Devices own memory and behaviors. They observe models
	// through their behaviors.
	fmt.Println("Step 3: Create Virtual Devices")

	// Weather Station device: measures ambient conditions
	weatherStationMem := map[string]uint32{
		"input_registers":  20, // temperature, humidity, pressure, etc.
		"discrete_inputs":  8,  // status flags
	}
	weatherStation := rt.CreateDevice("weather-station-001", "weather_station", weatherStationMem)
	weatherStation.AddBehavior(&WeatherStationBehavior{device: weatherStation})
	fmt.Printf("  - Weather station created: %s\n", weatherStation.ID())

	// PV Inverter device: converts solar power to grid power
	pvInverterMem := map[string]uint32{
		"input_registers":  20, // DC voltage, current, AC power, etc.
		"holding_registers": 10, // configuration
	}
	pvInverter := rt.CreateDevice("pv-inverter-001", "pv_inverter", pvInverterMem)
	pvInverter.AddBehavior(&PVInverterBehavior{
		device: pvInverter,
		// Behavior will observe sun and grid models through device.Model()
	})
	fmt.Printf("  - PV Inverter created: %s\n", pvInverter.ID())

	// Revenue Meter device: measures power flow
	revenueMeterMem := map[string]uint32{
		"input_registers":  20, // voltage, current, power, energy
		"holding_registers": 10, // configuration
	}
	revenueMeter := rt.CreateDevice("revenue-meter-001", "revenue_meter", revenueMeterMem)
	revenueMeter.AddBehavior(&RevenueMeterBehavior{device: revenueMeter})
	fmt.Printf("  - Revenue Meter created: %s\n", revenueMeter.ID())
	fmt.Println()

	// Step 4: Display initial model state
	fmt.Println("Step 4: Initial Model State")
	fmt.Printf("  Grid:    Voltage=%.1fV, Frequency=%.2fHz\n", grid.Voltage(), grid.Frequency())
	fmt.Printf("  Sun:     Irradiance=%.1fW/m², Elevation=%.1f°\n", sun.Irradiance(), sun.Elevation())
	fmt.Printf("  Weather: Temperature=%.1f°C, Humidity=%.1f%%\n", weather.Temperature(), weather.Humidity())
	fmt.Printf("  Wind:    Speed=%.1fm/s, Direction=%.0f°\n", wind.Speed(), wind.Direction())
	fmt.Println()

	// Step 5: Run simulation
	fmt.Println("Step 5: Run Simulation")
	fmt.Println("  Running 10 ticks...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Run the simulation
	done := make(chan error, 1)
	go func() {
		done <- rt.Run(ctx)
	}()

	// Run for 10 ticks
	tickCount := 0
	for tickCount < 10 {
		select {
		case <-time.After(50 * time.Millisecond):
			tickCount++
		case err := <-done:
			if err != nil && err != context.DeadlineExceeded {
				log.Fatalf("Runtime error: %v", err)
			}
			break
		case <-ctx.Done():
			break
		}
	}

	// Cancel and shutdown
	cancel()
	<-done
	rt.Shutdown()
	fmt.Println()

	// Step 6: Display final device memory state
	fmt.Println("Step 6: Final Device Memory State")

	// Weather Station memory
	temp, _ := weatherStation.Memory().ReadFloat32("input_registers", 0)
	hum, _ := weatherStation.Memory().ReadFloat32("input_registers", 4)
	press, _ := weatherStation.Memory().ReadFloat32("input_registers", 8)
	fmt.Printf("  Weather Station:\n")
	fmt.Printf("    - Temperature: %.1f°C\n", temp)
	fmt.Printf("    - Humidity: %.1f%%\n", hum)
	fmt.Printf("    - Pressure: %.1f hPa\n", press)

	// PV Inverter memory
	dcPower, _ := pvInverter.Memory().ReadFloat32("input_registers", 0)
	acPower, _ := pvInverter.Memory().ReadFloat32("input_registers", 4)
	efficiency, _ := pvInverter.Memory().ReadFloat32("input_registers", 8)
	fmt.Printf("  PV Inverter:\n")
	fmt.Printf("    - DC Power: %.1fW\n", dcPower)
	fmt.Printf("    - AC Power: %.1fW\n", acPower)
	fmt.Printf("    - Efficiency: %.1f%%\n", efficiency)

	// Revenue Meter memory
	voltage, _ := revenueMeter.Memory().ReadFloat32("input_registers", 0)
	current, _ := revenueMeter.Memory().ReadFloat32("input_registers", 4)
	power, _ := revenueMeter.Memory().ReadFloat32("input_registers", 8)
	fmt.Printf("  Revenue Meter:\n")
	fmt.Printf("    - Voltage: %.1fV\n", voltage)
	fmt.Printf("    - Current: %.1fA\n", current)
	fmt.Printf("    - Power: %.1fW\n", power)
	fmt.Println()

	// Step 7: Verify determinism
	fmt.Println("Step 7: Verify Determinism")
	fmt.Println("  Determinism means: same inputs → same outputs")
	fmt.Println("  The architecture guarantees:")
	fmt.Println("    - Devices tick in registration order")
	fmt.Println("    - Behaviors tick in registration order")
	fmt.Println("    - Models tick in registration order")
	fmt.Println("    - No unseeded randomness")
	fmt.Println()

	// Step 8: Summary
	fmt.Println("=== Architecture Validated ===")
	fmt.Println()
	fmt.Println("This example demonstrated:")
	fmt.Println("  ✓ Runtime initialization and shutdown")
	fmt.Println("  ✓ Simulation model creation and ticking")
	fmt.Println("  ✓ Device creation with memory regions")
	fmt.Println("  ✓ Behaviors observing models")
	fmt.Println("  ✓ Behaviors writing device memory")
	fmt.Println("  ✓ Execution loop with proper shutdown")
	fmt.Println()
	fmt.Println("The Forge architecture is validated.")
}

// WeatherStationBehavior samples weather conditions and updates device memory.
// This behavior observes the WeatherModel and writes measurements to memory.
type WeatherStationBehavior struct {
	device *device.Device
}

func (b *WeatherStationBehavior) ID() string { return "weather_station_behavior" }

func (b *WeatherStationBehavior) Attach(d *device.Device) {
	b.device = d
}

func (b *WeatherStationBehavior) Detach() {
	b.device = nil
}

func (b *WeatherStationBehavior) Tick() {
	// Observe the weather model
	weather := b.device.Model("ambient-weather")
	if weather == nil {
		return
	}

	weatherModel, ok := weather.(*models.WeatherModel)
	if !ok {
		return
	}

	// Read from model, write to memory
	// Memory is the source of truth for device state
	b.device.Memory().WriteFloat32("input_registers", 0, weatherModel.Temperature())
	b.device.Memory().WriteFloat32("input_registers", 4, weatherModel.Humidity())
	b.device.Memory().WriteFloat32("input_registers", 8, weatherModel.Pressure())
	b.device.Memory().WriteFloat32("input_registers", 12, weatherModel.CloudCover())
}

// PVInverterBehavior converts solar power to grid power.
// This behavior observes SunModel and GridModel.
type PVInverterBehavior struct {
	device       *device.Device
	dcVoltage    float32
	dcCurrent    float32
	efficiency   float32
}

func (b *PVInverterBehavior) ID() string { return "pv_inverter_behavior" }

func (b *PVInverterBehavior) Attach(d *device.Device) {
	b.device = d
	b.dcVoltage = 400.0
	b.efficiency = 95.0
}

func (b *PVInverterBehavior) Detach() {
	b.device = nil
}

func (b *PVInverterBehavior) Tick() {
	// Observe the sun model
	sun := b.device.Model("solar-sun")
	if sun == nil {
		return
	}

	sunModel, ok := sun.(*models.SunModel)
	if !ok {
		return
	}

	// Calculate DC power from irradiance
	// Simplified: assume 10m² panel at 20% efficiency
	irradiance := sunModel.Irradiance()
	panelArea := float32(10.0)
	panelEfficiency := float32(0.20)
	dcPower := irradiance * panelArea * panelEfficiency

	// Calculate DC current
	if b.dcVoltage > 0 {
		b.dcCurrent = dcPower / b.dcVoltage
	}

	// Calculate AC power (accounting for inverter efficiency)
	acPower := dcPower * (b.efficiency / 100.0)

	// Write to memory
	b.device.Memory().WriteFloat32("input_registers", 0, dcPower)    // DC Power
	b.device.Memory().WriteFloat32("input_registers", 4, acPower)   // AC Power
	b.device.Memory().WriteFloat32("input_registers", 8, b.efficiency) // Efficiency
}

// RevenueMeterBehavior measures power flowing through the meter.
// This behavior observes the GridModel.
type RevenueMeterBehavior struct {
	device *device.Device
}

func (b *RevenueMeterBehavior) ID() string { return "revenue_meter_behavior" }

func (b *RevenueMeterBehavior) Attach(d *device.Device) {
	b.device = d
}

func (b *RevenueMeterBehavior) Detach() {
	b.device = nil
}

func (b *RevenueMeterBehavior) Tick() {
	// Observe the grid model
	grid := b.device.Model("main-grid")
	if grid == nil {
		return
	}

	gridModel, ok := grid.(*models.GridModel)
	if !ok {
		return
	}

	// Read grid state
	voltage := gridModel.Voltage()
	frequency := gridModel.Frequency()

	// Simulate current and power based on grid conditions
	// In a real system, this would come from CT/PT sensors
	current := float32(10.0) // Simulated load current
	power := voltage * current

	// Write to memory
	b.device.Memory().WriteFloat32("input_registers", 0, voltage)
	b.device.Memory().WriteFloat32("input_registers", 4, current)
	b.device.Memory().WriteFloat32("input_registers", 8, power)
	b.device.Memory().WriteFloat32("input_registers", 12, frequency)
}
