// Package main provides an example of using the Simulation Inspector.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tamzrod/forge/internal/devices"
	"github.com/tamzrod/forge/internal/devices/weatherstation"
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

	// Create simulation context for devices
	simContext := devices.NewContext(simClock, sunModel, weatherModel, gridModel)

	// Create device registry
	registry := devices.NewRegistry()

	// Create Weather Station device with publishing enabled
	cfg := weatherstation.DefaultConfig()
	cfg.Publishing.Enabled = true
	cfg.Publishing.Host = "localhost"
	cfg.Publishing.Port = 500
	cfg.Publishing.UnitID = 1
	cfg.Publishing.Interval = 1 * time.Second

	ws, err := weatherstation.NewStation(cfg, simContext)
	if err != nil {
		log.Fatalf("Failed to create Weather Station: %v", err)
	}

	// Register device
	if err := registry.Register(ws); err != nil {
		log.Fatalf("Failed to register Weather Station: %v", err)
	}

	// Initialize device
	if err := ws.Initialize(); err != nil {
		log.Fatalf("Failed to initialize Weather Station: %v", err)
	}

	// Create inspector view
	view := inspector.NewView(simClock, sunModel, weatherModel, gridModel)
	view.SetRegistry(registry)

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
	simClock.Advance(12 * time.Hour)       // noon

	log.Println("Simulation Inspector running at http://localhost:8080")
	log.Println("Press Ctrl+C to stop")
	log.Println()
	log.Println("Simulation models: Clock, Sun, Weather, Grid")
	log.Println("Virtual devices: Weather Station #1")
	log.Println()
	log.Println("The Weather Station observes the Weather Model")
	log.Println("and copies values to its operational memory.")
	log.Println()
	log.Println("Publishing: Enabled (Raw Ingest)")
	log.Println("Target: localhost:500, UnitID: 1")
	log.Println()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	tickCount := 0
	for {
		select {
		case <-ctx.Done():
			ws.Shutdown()
			registry.Shutdown()
			server.Stop()
			return
		case <-ticker.C:
			// Tick simulation models
			simClock.Tick()
			sunModel.Tick()
			weatherModel.Tick()
			gridModel.Tick()

			// Tick devices (they observe models)
			registry.Tick()

			// Demonstrate the separation between model and device
			if tickCount%100 == 0 {
				log.Printf("Weather Model: %.1f°C, Weather Station: %.1f°C",
					weatherModel.Temperature(),
					ws.Temperature())
			}

			tickCount++

			// Stop after simulating some time
			if tickCount >= 86400 { // ~2 hours at 100ms tick
				log.Println("Simulation complete")
				cancel()
			}
		}
	}
}
