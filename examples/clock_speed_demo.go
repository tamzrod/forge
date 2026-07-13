// Package main demonstrates simulation clock speed control.
package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/tamzrod/forge/simulation"
)

func main() {
	speed := flag.Float64("speed", 1.0, "Simulation speed multiplier")
	mode := flag.String("mode", "realtime", "Simulation mode (realtime, simulated)")
	flag.Parse()

	fmt.Println("Forge Simulation Clock Demonstration")
	fmt.Println("====================================")
	fmt.Println()

	// Create clock
	clock := simulation.NewClock()

	// Set mode
	switch *mode {
	case "realtime":
		clock.SetMode(simulation.ModeRealtime)
	case "simulated":
		clock.SetMode(simulation.ModeSimulated)
	case "manual":
		clock.SetMode(simulation.ModeManual)
	default:
		clock.SetMode(simulation.ModeRealtime)
	}

	// Set speed
	clock.SetSpeed(*speed)

	// Start at a specific datetime
	startTime := time.Date(2026, 6, 21, 12, 0, 0, 0, time.UTC)
	fmt.Printf("Start Time: %s\n", startTime.Format(time.RFC3339))
	fmt.Printf("Mode: %s\n", clock.Mode())
	fmt.Printf("Speed: %.1fx\n", clock.Speed())
	fmt.Println()

	// Start the clock
	clock.Start(startTime)

	// Simulate ticks
	ticks := 50

	fmt.Println("Simulation Progress:")
	fmt.Println("-------------------")
	start := time.Now()

	for i := 0; i < ticks; i++ {
		clock.Update()
		
		elapsed := time.Since(start)
		simTime := clock.Now()
		
		// Print every 10 ticks
		if i%10 == 0 {
			fmt.Printf("Tick %3d | Wall: %-8s | Sim: %s | Speed: %.1fx\n",
				i, elapsed.Round(100*time.Millisecond), 
				simTime.Format("15:04:05.000"),
				clock.Speed())
		}

		// Change speed at tick 25
		if i == 25 {
			clock.SetSpeed(5.0)
			fmt.Println("  [Speed changed to 5x]")
		}

		// Sleep to control wall time
		if clock.Mode() == simulation.ModeRealtime {
			time.Sleep(10 * time.Millisecond)
		}
	}

	fmt.Println()
	fmt.Println("Final State:")
	fmt.Printf("  Ticks: %d\n", clock.Tick())
	fmt.Printf("  Sim Time: %s\n", clock.Now().Format(time.RFC3339))
	fmt.Printf("  Elapsed: %v\n", clock.Elapsed())
	fmt.Printf("  Wall Time: %v\n", time.Since(start))
	fmt.Printf("  Speed: %.1fx\n", clock.Speed())

	// Demonstrate different speeds
	fmt.Println()
	fmt.Println("Speed Comparison:")
	fmt.Println("-----------------")
	speeds := []float64{0.5, 1.0, 2.0, 5.0, 10.0, 20.0, 100.0}
	
	for _, s := range speeds {
		testClock := simulation.NewClock()
		testClock.SetMode(simulation.ModeSimulated)
		testClock.SetSpeed(s)
		testClock.Start(time.Date(2026, 6, 21, 12, 0, 0, 0, time.UTC))
		
		wallStart := time.Now()
		for i := 0; i < 100; i++ {
			testClock.Update()
		}
		wallTime := time.Since(wallStart)
		
		fmt.Printf("  %5.1fx: %v wall time → %v sim time\n", 
			s, wallTime.Round(time.Millisecond), testClock.Elapsed().Round(time.Millisecond))
	}
}
