// Package main runs the world simulation example.
package main

import (
	"fmt"
	"time"

	"github.com/tamzrod/forge/world"
	"github.com/tamzrod/forge/world/electrical"
	"github.com/tamzrod/forge/world/environmental"
)

func main() {
	fmt.Println("Forge World Simulation")
	fmt.Println("=====================")
	fmt.Println()

	// Create a new simulation world
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

	// Connect PV arrays to sun for irradiance
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

	// Simulation loop
	tickInterval := 100 * time.Millisecond
	tickCount := 0
	duration := 10 * time.Second
	endTime := time.Now().Add(duration)

	fmt.Println("Starting simulation...")
	fmt.Println()

	for time.Now().Before(endTime) {
		// Tick the world
		w.Tick(tickInterval)
		tickCount++

		// Manual power flow calculations (simplified)
		// In a real simulation, this would be handled by entity connections

		// Get measurements
		pv1Power := pv1PowerOutput(w, "pv-array-1")
		pv2Power := pv1PowerOutput(w, "pv-array-2")
		totalPV := pv1Power + pv2Power
		auxPower := getLoadPower(w, "aux-load")
		stationPower := getLoadPower(w, "station-load")
		totalLoad := auxPower + stationPower
		netPower := totalPV - totalLoad

		// Apply breaker state
		breakerOpen := isBreakerOpen(w, "grid-breaker")
		if breakerOpen {
			netPower = 0
		}

		// Update meter
		pccMeter.InjectInput("active_power", netPower)

		// Print status every second
		if tickCount%10 == 0 {
			sunMeas := w.Measurement(sun.ID(), "irradiance")
			weatherMeas := w.Measurement(weather.ID(), "temperature")
			voltageMeas := w.Measurement(grid.ID(), "voltage")

			fmt.Printf("[%3d] t=%v\n", tickCount/10, w.Time().Format("15:04:05"))
			fmt.Printf("  Sun: %.0f W/m2 | Weather: %.1f C\n",
				toFloat(sunMeas.Value), toFloat(weatherMeas.Value))
			fmt.Printf("  69kV Grid: %.0f V @ %.2f Hz\n",
				toFloat(voltageMeas.Value), 60.0)
			fmt.Printf("  PV Generation: %.1f kW | Load: %.1f kW\n",
				totalPV, totalLoad)
			fmt.Printf("  Net Power: %+.1f kW | Direction: %s\n",
				netPower, direction(netPower))
			if breakerOpen {
				fmt.Printf("  Grid Breaker: OPEN (Islanded)\n")
			} else {
				fmt.Printf("  Grid Breaker: CLOSED (Grid Connected)\n")
			}
			fmt.Println()

			// Publish events
			if tickCount == 200 {
				fmt.Println("*** Publishing: Grid Breaker Open Event ***")
				w.PublishEvent(world.Event{
					Type: "breaker_open",
					Source: "grid-breaker",
					Data: map[string]interface{}{"breaker_id": "grid-breaker"},
				})
				gridBreaker.Open()
			}
			if tickCount == 350 {
				fmt.Println("*** Publishing: Grid Breaker Close Event ***")
				w.PublishEvent(world.Event{
					Type: "breaker_close",
					Source: "grid-breaker",
					Data: map[string]interface{}{"breaker_id": "grid-breaker"},
				})
				gridBreaker.Close()
			}
		}

		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Simulation complete")
	fmt.Println()

	// Print summary
	fmt.Println("Final Measurements:")
	for _, m := range w.AllMeasurements() {
		if m.Value != nil && m.Value != 0 {
			fmt.Printf("  %s.%s = %v %s\n", m.EntityID, m.Name, m.Value, m.Unit)
		}
	}
}

// Helper functions

func pv1PowerOutput(w world.World, id string) float32 {
	m := w.Measurement(world.EntityID(id), "power")
	return toFloat(m.Value)
}

func getLoadPower(w world.World, id string) float32 {
	m := w.Measurement(world.EntityID(id), "power")
	return toFloat(m.Value)
}

func isBreakerOpen(w world.World, id string) bool {
	m := w.Measurement(world.EntityID(id), "state")
	return toString(m.Value) == "OPEN"
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
	case int64:
		return float32(val)
	default:
		return 0
	}
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

func direction(power float32) string {
	if power > 0 {
		return "EXPORT"
	} else if power < 0 {
		return "IMPORT"
	}
	return "ZERO"
}
