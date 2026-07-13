// Package main demonstrates the Solver architecture.
package main

import (
	"fmt"
	"time"

	"github.com/tamzrod/forge/solver"
	"github.com/tamzrod/forge/world"
)

func main() {
	fmt.Println("Forge Solver Architecture Demonstration")
	fmt.Println("======================================")
	fmt.Println()

	// Create a world
	w := world.NewWorld()

	// Create a solver
	elecSolver := solver.NewElectricalSolver()

	// Attach solver to world
	w.SetSolver(elecSolver)

	// Create electrical entities
	gen1 := world.NewVirtualGenerator("gen-1", "Solar Plant A", 500)
	gen2 := world.NewVirtualGenerator("gen-2", "Solar Plant B", 300)
	gen3 := world.NewVirtualGenerator("gen-3", "Backup Generator", 200)

	load1 := world.NewVirtualLoad("load-1", "Factory A", 400)
	load2 := world.NewVirtualLoad("load-2", "Factory B", 250)
	load3 := world.NewVirtualLoad("load-3", "Station Service", 50)

	meter := world.NewVirtualMeter("pcc-meter", "PCC Meter")

	// Add entities to world
	w.AddEntity(gen1)
	w.AddEntity(gen2)
	w.AddEntity(gen3)
	w.AddEntity(load1)
	w.AddEntity(load2)
	w.AddEntity(load3)
	w.AddEntity(meter)

	fmt.Println("Setup:")
	fmt.Printf("  World: %d entities\n", len(w.Entities()))
	fmt.Printf("  Solver: %s (%s)\n", elecSolver.Name(), elecSolver.Type())
	fmt.Printf("  Generators: %d (total capacity: %.0f kW)\n",
		elecSolver.GeneratorCount(),
		gen1.RatedCapacity()+gen2.RatedCapacity()+gen3.RatedCapacity())
	fmt.Printf("  Loads: %d (total demand: %.0f kW)\n",
		elecSolver.LoadCount(),
		load1.ActivePowerDemand()+load2.ActivePowerDemand()+load3.ActivePowerDemand())

	// Scenario: Normal Operation
	fmt.Println("\n--- Scenario: Normal Operation ---")

	gen1.SetTargetPower(400)
	gen2.SetTargetPower(250)
	gen3.SetOnline(false)

	dt := 100 * time.Millisecond
	for i := 0; i < 5; i++ {
		w.Tick(dt)
	}

	fmt.Printf("Solver State: %s\n", elecSolver)
	fmt.Printf("  Total Generation: %.1f kW\n", elecSolver.TotalGeneration())
	fmt.Printf("  Total Consumption: %.1f kW\n", elecSolver.TotalConsumption())
	fmt.Printf("  Net Power: %.1f kW\n", elecSolver.NetPower())
	fmt.Printf("  Capacity: %.0f kW (utilization: %.1f%%)\n",
		elecSolver.TotalCapacity(),
		elecSolver.CapacityUtilization()*100)

	// Scenario: Cloud Cover
	fmt.Println("\n--- Scenario: Cloud Cover ---")

	gen1.SetTargetPower(200)
	gen2.SetTargetPower(100)

	for i := 0; i < 5; i++ {
		w.Tick(dt)
	}

	fmt.Printf("Solver State: %s\n", elecSolver)
	fmt.Printf("  Total Generation: %.1f kW\n", elecSolver.TotalGeneration())
	fmt.Printf("  Total Consumption: %.1f kW\n", elecSolver.TotalConsumption())
	fmt.Printf("  Net Power: %.1f kW (%s)\n", elecSolver.NetPower(), direction(elecSolver.NetPower()))
	fmt.Printf("  Excess Capacity: %.1f kW\n", elecSolver.ExcessCapacity())

	// Scenario: Backup Generator
	fmt.Println("\n--- Scenario: Backup Generator Starts ---")

	gen3.SetOnline(true)
	gen3.SetTargetPower(200)

	for i := 0; i < 5; i++ {
		w.Tick(dt)
	}

	fmt.Printf("Solver State: %s\n", elecSolver)
	fmt.Printf("  Total Generation: %.1f kW\n", elecSolver.TotalGeneration())
	fmt.Printf("  Total Consumption: %.1f kW\n", elecSolver.TotalConsumption())
	fmt.Printf("  Net Power: %.1f kW (%s)\n", elecSolver.NetPower(), direction(elecSolver.NetPower()))
	fmt.Printf("  Load Served: %.1f%%\n", elecSolver.LoadServed()*100)

	// Scenario: Load Shed
	fmt.Println("\n--- Scenario: Load Shed ---")

	load3.SetConnected(false)

	for i := 0; i < 5; i++ {
		w.Tick(dt)
	}

	fmt.Printf("Solver State: %s\n", elecSolver)
	fmt.Printf("  Total Generation: %.1f kW\n", elecSolver.TotalGeneration())
	fmt.Printf("  Total Consumption: %.1f kW\n", elecSolver.TotalConsumption())
	fmt.Printf("  Net Power: %.1f kW (%s)\n", elecSolver.NetPower(), direction(elecSolver.NetPower()))
	fmt.Printf("  Load Served: %.1f%%\n", elecSolver.LoadServed()*100)

	// Meter energy
	fmt.Println("\nMeter Energy:")
	fmt.Printf("  Import: %.2f kWh\n", meter.EnergyImport())
	fmt.Printf("  Export: %.2f kWh\n", meter.EnergyExport())

	// Summary
	fmt.Println("\n======================================")
	fmt.Println("Solver Architecture Summary")
	fmt.Println("======================================")
	fmt.Println()
	fmt.Println("The Solver owns:")
	fmt.Println("  - Evaluation order")
	fmt.Println("  - Dependency resolution")
	fmt.Println("  - State propagation")
	fmt.Println("  - Network traversal")
	fmt.Println("  - Simulation iteration")
	fmt.Println()
	fmt.Println("The Solver does NOT own:")
	fmt.Println("  - Time (World clock)")
	fmt.Println("  - Topology (Network)")
	fmt.Println("  - Entity behavior (Entities)")
	fmt.Println("  - Events (World)")
	fmt.Println("  - Measurements (Entities)")
	fmt.Println()
	fmt.Println("World delegates simulation evolution to Solver.")
	fmt.Println("The Solver becomes the engineering engine of Forge.")
}

func direction(power float32) string {
	if power > 0.1 {
		return "EXPORT"
	} else if power < -0.1 {
		return "IMPORT"
	}
	return "BALANCED"
}
