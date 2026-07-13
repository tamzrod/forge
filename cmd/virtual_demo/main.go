// Package main demonstrates virtual electrical entities.
package main

import (
	"fmt"
	"time"

	"github.com/tamzrod/forge/world"
)

func main() {
	fmt.Println("Forge Virtual Electrical Entities Demonstration")
	fmt.Println("============================================")
	fmt.Println()

	// Create multiple virtual generators
	gen1 := world.NewVirtualGenerator("gen-1", "Solar Plant A", 500) // 500 kW
	gen2 := world.NewVirtualGenerator("gen-2", "Solar Plant B", 300) // 300 kW
	gen3 := world.NewVirtualGenerator("gen-3", "Backup Generator", 200) // 200 kW

	// Create multiple virtual loads
	load1 := world.NewVirtualLoad("load-1", "Factory A", 400) // 400 kW
	load2 := world.NewVirtualLoad("load-2", "Factory B", 250) // 250 kW
	load3 := world.NewVirtualLoad("load-3", "Station Service", 50) // 50 kW

	// Create a virtual meter at the point of common coupling
	meter := world.NewVirtualMeter("pcc-meter", "PCC Meter")

	// Create a world to host entities
	w := world.NewWorld()

	// Add entities to the world
	w.AddEntity(gen1)
	w.AddEntity(gen2)
	w.AddEntity(gen3)
	w.AddEntity(load1)
	w.AddEntity(load2)
	w.AddEntity(load3)
	w.AddEntity(meter)

	fmt.Println("Initial Setup:")
	fmt.Println("--------------")
	fmt.Printf("Generators: %d total, rated capacity: %.0f kW\n",
		3, gen1.RatedCapacity()+gen2.RatedCapacity()+gen3.RatedCapacity())
	fmt.Printf("Loads: %d total, base demand: %.0f kW\n",
		3, load1.ActivePowerDemand()+load2.ActivePowerDemand()+load3.ActivePowerDemand())

	// Simulate a scenario
	fmt.Println("\n--- Scenario: Normal Operation ---")

	// Dispatch generators
	gen1.SetTargetPower(400) // Solar at 80%
	gen2.SetTargetPower(250) // Solar at 83%
	gen3.SetOnline(false)   // Backup offline

	// Simulate several ticks
	dt := 100 * time.Millisecond
	for i := 0; i < 5; i++ {
		w.Tick(dt)
	}

	// Calculate totals
	totalGen := gen1.ActivePower() + gen2.ActivePower() + gen3.ActivePower()
	totalLoad := load1.ActivePowerDemand() + load2.ActivePowerDemand() + load3.ActivePowerDemand()
	netPower := totalGen - totalLoad

	// Update meter
	meter.SetMeasurements(netPower, 0, 69000, 60)

	fmt.Printf("Total Generation: %.1f kW\n", totalGen)
	fmt.Printf("Total Load: %.1f kW\n", totalLoad)
	fmt.Printf("Net Power: %.1f kW (%s)\n", netPower, direction(netPower))

	fmt.Println("\nGenerator Status:")
	fmt.Printf("  %s: %.1f kW (online: %v)\n", gen1.Name(), gen1.ActivePower(), gen1.IsOnline())
	fmt.Printf("  %s: %.1f kW (online: %v)\n", gen2.Name(), gen2.ActivePower(), gen2.IsOnline())
	fmt.Printf("  %s: %.1f kW (online: %v)\n", gen3.Name(), gen3.ActivePower(), gen3.IsOnline())

	fmt.Println("\nLoad Status:")
	fmt.Printf("  %s: %.1f kW (connected: %v)\n", load1.Name(), load1.ActivePowerDemand(), load1.IsConnected())
	fmt.Printf("  %s: %.1f kW (connected: %v)\n", load2.Name(), load2.ActivePowerDemand(), load2.IsConnected())
	fmt.Printf("  %s: %.1f kW (connected: %v)\n", load3.Name(), load3.ActivePowerDemand(), load3.IsConnected())

	// Scenario: Cloud cover reduces generation
	fmt.Println("\n--- Scenario: Cloud Cover (Reduced Generation) ---")

	gen1.SetTargetPower(200) // Solar drops to 40%
	gen2.SetTargetPower(100) // Solar drops to 33%

	for i := 0; i < 5; i++ {
		w.Tick(dt)
	}

	totalGen = gen1.ActivePower() + gen2.ActivePower() + gen3.ActivePower()
	totalLoad = load1.ActivePowerDemand() + load2.ActivePowerDemand() + load3.ActivePowerDemand()
	netPower = totalGen - totalLoad

	// Update meter
	meter.SetMeasurements(netPower, 0, 69000, 60)

	fmt.Printf("Total Generation: %.1f kW\n", totalGen)
	fmt.Printf("Total Load: %.1f kW\n", totalLoad)
	fmt.Printf("Net Power: %.1f kW (%s)\n", netPower, direction(netPower))

	// Scenario: Start backup generator
	fmt.Println("\n--- Scenario: Backup Generator Starts ---")

	gen3.SetOnline(true)
	gen3.SetTargetPower(200)

	for i := 0; i < 5; i++ {
		w.Tick(dt)
	}

	totalGen = gen1.ActivePower() + gen2.ActivePower() + gen3.ActivePower()
	totalLoad = load1.ActivePowerDemand() + load2.ActivePowerDemand() + load3.ActivePowerDemand()
	netPower = totalGen - totalLoad

	// Update meter
	meter.SetMeasurements(netPower, 0, 69000, 60)

	fmt.Printf("Total Generation: %.1f kW\n", totalGen)
	fmt.Printf("Total Load: %.1f kW\n", totalLoad)
	fmt.Printf("Net Power: %.1f kW (%s)\n", netPower, direction(netPower))

	// Scenario: Load shed
	fmt.Println("\n--- Scenario: Load Shed (Priority Disconnection) ---")

	load3.SetConnected(false) // Shed station service

	for i := 0; i < 5; i++ {
		w.Tick(dt)
	}

	totalGen = gen1.ActivePower() + gen2.ActivePower() + gen3.ActivePower()
	totalLoad = load1.ActivePowerDemand() + load2.ActivePowerDemand() + load3.ActivePowerDemand()
	netPower = totalGen - totalLoad

	// Update meter
	meter.SetMeasurements(netPower, 0, 69000, 60)

	fmt.Printf("Total Generation: %.1f kW\n", totalGen)
	fmt.Printf("Total Load: %.1f kW\n", totalLoad)
	fmt.Printf("Net Power: %.1f kW (%s)\n", netPower, direction(netPower))

	// Meter measurements
	fmt.Println("\nMeter Measurements:")
	fmt.Printf("  Active Power: %.1f kW\n", meter.ActivePower())
	fmt.Printf("  Energy Import: %.2f kWh\n", meter.EnergyImport())
	fmt.Printf("  Energy Export: %.2f kWh\n", meter.EnergyExport())

	// Summary
	fmt.Println("\n============================================")
	fmt.Println("Summary: Virtual Entities Behavior")
	fmt.Println("============================================")
	fmt.Println()
	fmt.Println("Virtual Generators:")
	fmt.Println("  - Injected power based on target and ramp rate")
	fmt.Println("  - Respected available capacity limits")
	fmt.Println("  - Responded to online/offline status")
	fmt.Println()
	fmt.Println("Virtual Loads:")
	fmt.Println("  - Consumed power based on demand")
	fmt.Println("  - Responded to connect/disconnect status")
	fmt.Println("  - Priority-based shedding supported")
	fmt.Println()
	fmt.Println("Virtual Meter:")
	fmt.Println("  - Measured net power flow")
	fmt.Println("  - Accumulated import/export energy")
	fmt.Println("  - Positive = Export, Negative = Import")
	fmt.Println()
	fmt.Println("The World determines:")
	fmt.Println("  - Total Generation")
	fmt.Println("  - Total Consumption")
	fmt.Println("  - Net Power")
	fmt.Println("  - Import / Export")
	fmt.Println()
	fmt.Println("Without knowing underlying technology!")
}

func direction(power float32) string {
	if power > 0 {
		return "EXPORT"
	} else if power < 0 {
		return "IMPORT"
	}
	return "BALANCED"
}
