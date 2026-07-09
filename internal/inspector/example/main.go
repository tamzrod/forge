// Package main provides an example of using the Simulation Inspector.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tamzrod/forge/internal/inspector"
	"github.com/tamzrod/forge/internal/models/clock"
	"github.com/tamzrod/forge/internal/models/grid"
	"github.com/tamzrod/forge/internal/models/sun"
	"github.com/tamzrod/forge/internal/models/weather"
)

func main() {
	// Create simulation clock
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: 100 * time.Millisecond,
	})

	// Create simulation models
	sunModel := sun.New(sun.Config{
		Latitude:  40.0,  // Denver
		Longitude: -105.0,
	}, simClock)

	weatherModel := weather.New(weather.DefaultConfig(), simClock)

	gridModel := grid.New(grid.DefaultConfig(), simClock)

	// Create inspector view
	view := inspector.NewView(simClock, sunModel, weatherModel, gridModel)

	// Create inspector server
	server := inspector.NewServer(view, inspector.DefaultConfig())

	// Start server in background
	go func() {
		if err := server.Start(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Run simulation loop
	ctx, cancel := context.WithCancel(context.Background())

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down...")
		cancel()
	}()

	// Advance simulation time
	// Start at solar noon on a spring day
	simClock.Advance(80 * 24 * time.Hour) // ~March 21
	simClock.Advance(12 * time.Hour)      // noon

	log.Println("Simulation Inspector running at http://localhost:8080")
	log.Println("Press Ctrl+C to stop")
	log.Println()
	log.Println("Simulating 24 hours of sun movement...")

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	tickCount := 0
	for {
		select {
		case <-ctx.Done():
			server.Stop()
			return
		case <-ticker.C:
			// Tick all models
			simClock.Tick()
			sunModel.Tick()
			weatherModel.Tick()
			gridModel.Tick()

			// Inject some power to demonstrate grid response
			if tickCount%100 == 0 {
				gridModel.InjectActivePower(50.0)
				gridModel.InjectReactivePower(25.0)
			}

			tickCount++

			// Stop after simulating a full day (every tick = 100ms)
			if tickCount >= 864000 { // 24 hours * 60 min * 60 sec * 10 ticks
				log.Println("Simulation complete")
				cancel()
			}
		}
	}
}
