// Package main runs the scenario demonstration.
package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/tamzrod/forge/scenarios"
	"github.com/tamzrod/forge/world"
	"github.com/tamzrod/forge/world/electrical"
	"github.com/tamzrod/forge/world/environmental"
)

func main() {
	// Define flags first
	scenarioName := flag.String("scenario", "breaker-trip", "Scenario to run (normal, cloud, breaker-trip, islanding, storm)")
	flag.Parse()

	fmt.Println("Forge Scenario Demonstration")
	fmt.Println("============================")
	fmt.Println()

	// Create the simulation world
	w := world.NewWorld()
	defer w.Close()

	// Add environmental entities
	sun := environmental.NewSun("sun")
	weather := environmental.NewWeather("weather")
	pv1 := environmental.NewPVArray("pv-array-1", 50.0)
	pv2 := environmental.NewPVArray("pv-array-2", 50.0)
	w.AddEntity(sun)
	w.AddEntity(weather)
	w.AddEntity(pv1)
	w.AddEntity(pv2)

	// Connect PV arrays to sun
	pv1.Connect("irradiance", sun.ID(), "irradiance")
	pv2.Connect("irradiance", sun.ID(), "irradiance")

	// Add electrical entities
	grid := electrical.NewGrid("utility-grid", 69000, 60.0)
	collectorBus := electrical.NewBus("collector-bus", 480)
	gridBreaker := electrical.NewBreaker("grid-breaker")
	pccMeter := electrical.NewMeter("pcc-meter")
	auxLoad := electrical.NewLoad("aux-load", 5.0)
	stationLoad := electrical.NewLoad("station-load", 3.0)
	w.AddEntity(grid)
	w.AddEntity(collectorBus)
	w.AddEntity(gridBreaker)
	w.AddEntity(pccMeter)
	w.AddEntity(auxLoad)
	w.AddEntity(stationLoad)

	// Connect meter to grid
	pccMeter.Connect("voltage", grid.ID(), "voltage")
	pccMeter.Connect("frequency", grid.ID(), "frequency")

	// Create scenario runner
	runner := scenarios.NewScenarioRunner(w)

	// Select scenario based on flag or default
	var scenario scenarios.Scenario
	switch *scenarioName {
	case "normal":
		scenario = scenarios.NewNormalDay()
	case "cloud":
		scenario = scenarios.NewCloudPassing()
	case "voltage-sag":
		scenario = scenarios.NewGridVoltageSag()
	case "frequency":
		scenario = scenarios.NewFrequencyExcursion()
	case "breaker-trip", "":
		scenario = scenarios.NewBreakerTrip()
	case "islanding":
		scenario = scenarios.NewIslanding()
	case "load-step":
		scenario = scenarios.NewLoadStep()
	case "storm":
		scenario = scenarios.NewStorm()
	default:
		scenario = scenarios.NewBreakerTrip()
	}

	fmt.Printf("Loading Scenario: %s\n", scenario.Name())
	fmt.Printf("Status: %s\n", scenario.String())
	fmt.Printf("Duration: %v\n\n", scenario.Duration())

	// Load and start scenario
	if err := runner.Load(scenario); err != nil {
		fmt.Printf("Error loading scenario: %v\n", err)
		return
	}
	if err := runner.Start(); err != nil {
		fmt.Printf("Error starting scenario: %v\n", err)
		return
	}

	// Simulation loop
	tickInterval := 100 * time.Millisecond
	tickCount := 0
	printCount := 0

	fmt.Println("Starting simulation...")
	fmt.Println()

	for {
		// Tick the world
		w.Tick(tickInterval)
		tickCount++

		// Handle scenario events
		runner.Tick(tickInterval)

		// Check for world events
		for _, evt := range w.Events() {
			handleEvent(w, evt, grid, gridBreaker, auxLoad)
		}

		// Calculate power
		pv1Power := getPower(w, "pv-array-1")
		pv2Power := getPower(w, "pv-array-2")
		totalPV := pv1Power + pv2Power
		auxPower := getPower(w, "aux-load")
		stationPower := getPower(w, "station-load")
		totalLoad := auxPower + stationPower
		netPower := totalPV - totalLoad

		if gridBreaker.IsOpen() {
			netPower = 0
		}

		pccMeter.InjectInput("active_power", netPower)

		// Print status every second
		if tickCount%10 == 0 {
			printCount++
			sunMeas := w.Measurement(sun.ID(), "irradiance")
			weatherMeas := w.Measurement(weather.ID(), "temperature")
			voltageMeas := w.Measurement(grid.ID(), "voltage")

			direction := "EXPORT"
			if netPower < 0 {
				direction = "IMPORT"
			} else if netPower == 0 {
				direction = "ZERO"
			}

			fmt.Printf("[%3d] t=%-8v | %s\n", printCount, scenario.Elapsed().Round(time.Second), scenario.Name())
			fmt.Printf("       Sun: %5.0f W/m2 | Weather: %4.1f C\n",
				toFloat(sunMeas.Value), toFloat(weatherMeas.Value))
			fmt.Printf("       Grid: %.0f V @ %.2f Hz | PV: %5.1f kW | Load: %5.1f kW\n",
				toFloat(voltageMeas.Value), 60.0, totalPV, totalLoad)
			fmt.Printf("       Net: %+6.1f kW [%s] | Breaker: %s\n",
				netPower, direction, breakerState(gridBreaker))
			fmt.Println()
		}

		// Check scenario completion
		if scenario.IsComplete() || !runner.IsRunning() {
			fmt.Println("Scenario complete!")
			break
		}

		// Safety timeout
		if scenario.Elapsed() > scenario.Duration()+10*time.Second {
			fmt.Println("Timeout - forcing stop")
			break
		}

		time.Sleep(10 * time.Millisecond)
	}

	// Print scenario summary
	fmt.Println()
	fmt.Println("Scenario Summary:")
	fmt.Printf("  Name: %s\n", scenario.Name())
	fmt.Printf("  Duration: %v\n", scenario.Duration())
	fmt.Printf("  Elapsed: %v\n", scenario.Elapsed())
	fmt.Printf("  Events Published: %d\n", len(scenario.Events()))
	fmt.Println()
	fmt.Println("Events:")
	for _, evt := range scenario.Events() {
		fmt.Printf("  - [%v] %s: %v\n", evt.Time.Round(time.Second), evt.Type, evt.Data)
	}
}

// handleEvent processes simulation events.
func handleEvent(w world.World, evt world.Event, grid *electrical.GridEntity, breaker *electrical.BreakerEntity, load *electrical.LoadEntity) {
	switch evt.Type {
	case "breaker_open":
		fmt.Printf("\n*** EVENT: Breaker Open ***\n")
		breaker.Open()

	case "breaker_close":
		fmt.Printf("\n*** EVENT: Breaker Close ***\n")
		breaker.Close()

	case "cloud_cover":
		if coverage, ok := evt.Data["coverage"].(float32); ok {
			fmt.Printf("\n*** EVENT: Cloud Cover = %.0f%% ***\n", coverage*100)
			// Send to weather entity
			if weather := w.Entity("weather"); weather != nil {
				// Weather entity would handle this
			}
		}

	case "wind_gust":
		if speed, ok := evt.Data["speed"].(float32); ok {
			fmt.Printf("\n*** EVENT: Wind Gust = %.0f m/s ***\n", speed)
		}

	case "voltage_sag":
		if voltage, ok := evt.Data["voltage"].(float32); ok {
			fmt.Printf("\n*** EVENT: Grid Voltage Sag to %.0f V ***\n", voltage)
		}

	case "voltage_sag_end":
		fmt.Printf("\n*** EVENT: Grid Voltage Restored ***\n")

	case "frequency_sag":
		if freq, ok := evt.Data["frequency"].(float32); ok {
			fmt.Printf("\n*** EVENT: Frequency Excursion to %.1f Hz ***\n", freq)
		}

	case "frequency_sag_end":
		fmt.Printf("\n*** EVENT: Frequency Restored ***\n")

	case "load_increase":
		if power, ok := evt.Data["power"].(float32); ok {
			fmt.Printf("\n*** EVENT: Load Increase to %.1f kW ***\n", power)
			// In real scenario, would modify load entity
		}

	case "load_decrease":
		if power, ok := evt.Data["power"].(float32); ok {
			fmt.Printf("\n*** EVENT: Load Decrease to %.1f kW ***\n", power)
		}

	case "generator_trip":
		fmt.Printf("\n*** EVENT: Generator Trip ***\n")

	case "generator_start":
		fmt.Printf("\n*** EVENT: Generator Start ***\n")
	}
}

// Helper functions
func getPower(w world.World, id string) float32 {
	m := w.Measurement(world.EntityID(id), "power")
	return toFloat(m.Value)
}

func toFloat(v interface{}) float32 {
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case float32:
		return val
	case float64:
		return float32(val)
	case int:
		return float32(val)
	default:
		return 0
	}
}

func breakerState(b *electrical.BreakerEntity) string {
	if b.IsOpen() {
		return "OPEN"
	}
	return "CLOSED"
}
