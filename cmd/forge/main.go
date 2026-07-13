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
	flagMaxDevices   = flag.Int("devices", 100, "Maximum number of devices")
	flagDuration     = flag.Duration("duration", 0, "Run duration (0 = indefinite)")
	flagVerbose      = flag.Bool("v", false, "Verbose output")
)

func main() {
	flag.Parse()

	if *flagVerbose {
		fmt.Println("Forge Industrial Simulation Runtime")
		fmt.Printf("Tick Interval: %v\n", *flagTickInterval)
		fmt.Printf("Max Devices: %d\n", *flagMaxDevices)
		fmt.Println()
	}

	cfg := runtime.Config{
		TickInterval: *flagTickInterval,
		MaxDevices:  *flagMaxDevices,
	}

	rt := runtime.New(cfg)

	// Environment models
	_ = rt.CreateSunModel("solar-sun")
	_ = rt.CreateWeatherModel("ambient-weather")
	_ = rt.CreateWindModel("wind-farm")

	// Electrical system: Grid -> Bus -> Transformer -> Bus -> Breaker -> Load
	grid := rt.CreateGridModel("utility-grid")
	_ = rt.CreateBusModel("bus-grid", 480)
	_ = rt.CreateTransformerModel("xfmr-main", "bus-grid", "bus-feeder")
	feederBus := rt.CreateBusModel("bus-feeder", 480)
	feederBreaker := rt.CreateBreakerModel("breaker-feeder", "bus-feeder", "bus-load")

	// Loads on feeder bus (in kW for 480V simulation)
	_ = rt.CreateLoadModel("load-building", "bus-feeder", 10.0)   // 10 kW base load
	_ = rt.CreateLoadModel("load-industrial", "bus-feeder", 20.0) // 20 kW base load

	if *flagVerbose {
		fmt.Println("Electrical System:")
		fmt.Println("  Utility Grid -> Bus Grid -> Xfmr Main -> Bus PV")
		fmt.Println("    -> Xfmr Feeder -> Bus Feeder -> Breaker -> Loads")
		fmt.Println()
	}

	// Weather station
	weatherStation := rt.CreateDevice("ws-001", "weather_station", map[string]uint32{
		"sensors": 64, "computed": 64, "status": 16,
	})
	weatherStation.AddBehavior(&weatherBehavior{device: weatherStation})

	// PV inverter
	pvInverter := rt.CreateDevice("pv-001", "pv_inverter", map[string]uint32{
		"input": 32, "output": 32, "config": 32, "status": 16,
	})
	pvInverter.AddBehavior(&pvBehavior{device: pvInverter, pvBus: feederBus})

	// Revenue meter
	meter := rt.CreateDevice("meter-001", "revenue_meter", map[string]uint32{
		"input_registers": 200, "status": 16,
	})
	meter.AddBehavior(&meterBehavior{
		device: meter, grid: grid, feederBus: feederBus, feederBreaker: feederBreaker,
	})

	// Breaker control
	breakerCtrl := rt.CreateDevice("breaker-ctrl", "breaker_controller", map[string]uint32{
		"control": 16, "status": 16,
	})
	breakerCtrl.AddBehavior(&breakerControl{device: breakerCtrl, feederBreaker: feederBreaker})

	if *flagVerbose {
		fmt.Printf("Models: %d, Devices: %d\n\n", len(rt.Models()), len(rt.Devices()))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\nShutting down...")
		cancel()
	}()

	done := make(chan error, 1)
	go func() { done <- rt.Run(ctx) }()

	tickCount := 0
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	fmt.Println("Simulation running... Press Ctrl+C to stop")
	fmt.Println()

	runDuration := 10 * time.Second
	if *flagDuration > 0 {
		runDuration = *flagDuration
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
				sun := rt.Scheduler().Model("solar-sun").(*models.SunModel)
				weather := rt.Scheduler().Model("ambient-weather").(*models.WeatherModel)

				pvV, _ := meter.Memory().ReadFloat32("input_registers", 0)
				pvP, _ := meter.Memory().ReadFloat32("input_registers", 4)
				loadP, _ := meter.Memory().ReadFloat32("input_registers", 8)
				netP, _ := meter.Memory().ReadFloat32("input_registers", 12)
				breakerOpen, _ := meter.Memory().ReadFloat32("input_registers", 16)

				breakerStatus := "CLOSED"
				if breakerOpen > 0 {
					breakerStatus = "OPEN"
				}

				fmt.Printf("[%3d] t=%v\n", tickCount, clock.Elapsed())
				fmt.Printf("  Sun: %.0fW/m² @ %.0f°el | Weather: %.1f°C %.0f%%rh\n",
					sun.Irradiance(), sun.Elevation(), weather.Temperature(), weather.Humidity())
				fmt.Printf("  Grid: %.1fV @ %.2fHz | Feeder V: %.1fV\n",
					grid.Voltage(), grid.Frequency(), pvV)
				fmt.Printf("  PV: %.1fkW | Load: %.1fkW | Net: %.1fkW | Breaker: %s\n",
					pvP, loadP, netP, breakerStatus)
				fmt.Println()
			} else {
				fmt.Printf("\rRunning... Tick %d", tickCount)
			}
		}
	}

shutdown:
	<-done
	rt.Shutdown()
	fmt.Println("\nSimulation complete")
}

type weatherBehavior struct {
	device *device.Device
}

func (b *weatherBehavior) ID() string                             { return "weather_sampling" }
func (b *weatherBehavior) Attach(d *device.Device)              { b.device = d }
func (b *weatherBehavior) Detach()                               { b.device = nil }
func (b *weatherBehavior) Tick() {
	if b.device == nil {
		return
	}
	wm, ok := b.device.Model("ambient-weather").(*models.WeatherModel)
	if !ok {
		return
	}
	b.device.Memory().WriteFloat32("sensors", 0, wm.Temperature())
	b.device.Memory().WriteFloat32("sensors", 4, wm.Humidity())
	b.device.Memory().WriteFloat32("sensors", 8, wm.Pressure())
}

type pvBehavior struct {
	device *device.Device
	pvBus  *models.BusModel
}

func (b *pvBehavior) ID() string { return "pv_power_calc" }
func (b *pvBehavior) Attach(d *device.Device) { b.device = d }
func (b *pvBehavior) Detach()     { b.device = nil }
func (b *pvBehavior) Tick() {
	if b.device == nil || b.pvBus == nil {
		return
	}
	sm, ok := b.device.Model("solar-sun").(*models.SunModel)
	if !ok {
		return
	}
	irradiance := sm.Irradiance()
	panelArea := float32(100.0)
	efficiency := float32(0.18)
	powerKW := irradiance * panelArea * efficiency / 1000.0

	b.device.Memory().WriteFloat32("input", 0, irradiance)
	b.device.Memory().WriteFloat32("output", 0, powerKW)

	// Inject PV power to bus (positive = generation)
	b.pvBus.InjectPower(powerKW)
}

type meterBehavior struct {
	device        *device.Device
	grid          *models.GridModel
	feederBus     *models.BusModel
	feederBreaker *models.BreakerModel
}

func (b *meterBehavior) ID() string { return "meter_measurement" }
func (b *meterBehavior) Attach(d *device.Device) { b.device = d }
func (b *meterBehavior) Detach()     { b.device = nil }
func (b *meterBehavior) Tick() {
	if b.device == nil || b.feederBus == nil {
		return
	}

	// Get loads
	loadB := b.device.Model("load-building").(*models.LoadModel)
	loadI := b.device.Model("load-industrial").(*models.LoadModel)
	totalLoad := loadB.CurrentLoad() + loadI.CurrentLoad()

	// Withdraw load from bus only if breaker is closed
	if !b.feederBreaker.IsOpen() {
		b.feederBus.WithdrawPower(totalLoad)
	}

	// Read values from bus
	pvGen := b.feederBus.PowerInjection()
	feederV := b.feederBus.ActualVoltage()

	// Net power (positive = consumption from grid, negative = export)
	netPower := totalLoad - pvGen

	// Write measurements
	b.device.Memory().WriteFloat32("input_registers", 0, feederV)        // Feeder voltage
	b.device.Memory().WriteFloat32("input_registers", 4, pvGen)          // PV generation
	b.device.Memory().WriteFloat32("input_registers", 8, totalLoad)      // Total load
	b.device.Memory().WriteFloat32("input_registers", 12, netPower)      // Net power
	var breakerStatus float32 = 0
	if b.feederBreaker.IsOpen() {
		breakerStatus = 1
	}
	b.device.Memory().WriteFloat32("input_registers", 16, breakerStatus) // Breaker status
}

type breakerControl struct {
	device        *device.Device
	feederBreaker *models.BreakerModel
	tickCounter   int
}

func (b *breakerControl) ID() string { return "breaker_control" }
func (b *breakerControl) Attach(d *device.Device) { b.device = d }
func (b *breakerControl) Detach()     { b.device = nil }
func (b *breakerControl) Tick() {
	if b.device == nil || b.feederBreaker == nil {
		return
	}
	b.tickCounter++

	if b.tickCounter == 30 && !b.feederBreaker.IsOpen() {
		b.feederBreaker.Open()
		fmt.Println("\n*** BREAKER TRIPPED ***")
	}
	if b.tickCounter == 50 && b.feederBreaker.IsOpen() {
		b.feederBreaker.Close()
		fmt.Println("\n*** BREAKER RECLOSED ***")
	}

	var status float32 = 0
	if b.feederBreaker.IsOpen() {
		status = 1
	}
	b.device.Memory().WriteFloat32("status", 0, status)
}
