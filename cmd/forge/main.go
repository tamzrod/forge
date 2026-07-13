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
	flagDuration    = flag.Duration("duration", 0, "Run duration (0 = indefinite)")
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
		MaxDevices:   *flagMaxDevices,
	}

	rt := runtime.New(cfg)

	// Environment models
	_ = rt.CreateSunModel("solar-sun")
	_ = rt.CreateWeatherModel("ambient-weather")
	_ = rt.CreateWindModel("wind-farm")

	// Power Generation Facility Architecture:
	// PV Arrays -> Collector Bus -> Transformer -> 69kV Bus -> Reference Meter -> Utility Grid
	// Auxiliary Loads

	// Utility Grid (infinite source/sink at 69kV)
	grid := rt.CreateGridModel("utility-grid")
	grid.SetVoltage(69000) // 69kV

	// Collector Bus (480V)
	collectorBus := rt.CreateBusModel("bus-collector", 480)

	// PV Arrays inject to collector bus (100kW total rated)
	_ = rt.CreatePVArrayModel("pv-array-1", "bus-collector", 50.0) // 50kW rated
	_ = rt.CreatePVArrayModel("pv-array-2", "bus-collector", 50.0) // 50kW rated

	// Auxiliary loads (consume power from collector)
	_ = rt.CreateLoadModel("load-aux", "bus-collector", 5.0)     // 5 kW auxiliary
	_ = rt.CreateLoadModel("load-station", "bus-collector", 3.0) // 3 kW station service

	// Grid breaker (connects plant to grid)
	gridBreaker := rt.CreateBreakerModel("breaker-grid", "bus-pcc", "bus-utility")

	if *flagVerbose {
		fmt.Println("Power Generation Facility:")
		fmt.Println("  PV Arrays (100kW rated)")
		fmt.Println("    -> Collector Bus (480V)")
		fmt.Println("    -> Step-up Transformer")
		fmt.Println("    -> 69kV Bus")
		fmt.Println("    -> Reference Meter (PCC)")
		fmt.Println("    -> Utility Grid (69kV infinite)")
		fmt.Println()
		fmt.Println("  Auxiliary Loads (8kW)")
		fmt.Println("    -> Collector Bus")
		fmt.Println()
	}

	// Weather station device
	weatherStation := rt.CreateDevice("ws-001", "weather_station", map[string]uint32{
		"sensors": 64, "computed": 64, "status": 16,
	})
	weatherStation.AddBehavior(&weatherBehavior{device: weatherStation})

	// 69kV Reference Meter
	refMeter := rt.CreateDevice("ref-meter-69kv", "reference_meter", map[string]uint32{
		"input_registers": 200,
		"status":         16,
	})
	refMeter.AddBehavior(&refMeterBehavior{
		device:      refMeter,
		grid:        grid,
		gridBreaker: gridBreaker,
	})

	// Plant Controller (manages power flow and breaker)
	plantController := rt.CreateDevice("plant-ctrl", "plant_controller", map[string]uint32{
		"input_registers": 100,
		"status":         16,
	})
	plantController.AddBehavior(&plantBehavior{
		device:       plantController,
		collectorBus: collectorBus,
		gridBreaker:  gridBreaker,
	})

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

				// Read 69kV meter values
				voltage, _ := refMeter.Memory().ReadFloat32("input_registers", 0)
				freq, _ := refMeter.Memory().ReadFloat32("input_registers", 4)
				activeP, _ := refMeter.Memory().ReadFloat32("input_registers", 8)
				reactiveQ, _ := refMeter.Memory().ReadFloat32("input_registers", 12)
				pf, _ := refMeter.Memory().ReadFloat32("input_registers", 16)
				energyExp, _ := refMeter.Memory().ReadFloat32("input_registers", 20)
				energyImp, _ := refMeter.Memory().ReadFloat32("input_registers", 24)
				direction, _ := refMeter.Memory().ReadFloat32("input_registers", 28)

				directionStr := "EXPORT"
				if direction < 0 {
					directionStr = "IMPORT"
				}

				fmt.Printf("[%3d] t=%v\n", tickCount, clock.Elapsed())
				fmt.Printf("  Sun: %.0fW/m2 @ %.0fdeg | Weather: %.1fC %.0f%%rh\n",
					sun.Irradiance(), sun.Elevation(), weather.Temperature(), weather.Humidity())
				fmt.Printf("  69kV PCC: %.0fV @ %.2fHz | P: %+.1fkW | Q: %+.1fkvar\n",
					voltage, freq, activeP, reactiveQ)
				fmt.Printf("  PF: %.3f | Energy Exp: %.4fkWh | Energy Imp: %.4fkWh\n",
					pf, energyExp, energyImp)
				fmt.Printf("  Power Flow: %s\n", directionStr)
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

// weatherBehavior observes weather model
type weatherBehavior struct {
	device *device.Device
}

func (b *weatherBehavior) ID() string                        { return "weather_sampling" }
func (b *weatherBehavior) Attach(d *device.Device)           { b.device = d }
func (b *weatherBehavior) Detach()                          { b.device = nil }
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

// refMeterBehavior measures power at the 69kV PCC
type refMeterBehavior struct {
	device      *device.Device
	grid        *models.GridModel
	gridBreaker *models.BreakerModel

	// Energy accumulation
	energyExport float32
	energyImport float32
	tickHours    float32
}

func (b *refMeterBehavior) ID() string                        { return "reference_meter" }
func (b *refMeterBehavior) Attach(d *device.Device)          { b.device = d }
func (b *refMeterBehavior) Detach()                          { b.device = nil }

func (b *refMeterBehavior) Tick() {
	if b.device == nil || b.grid == nil {
		return
	}

	// Get PV generation from sun model
	sunModel := b.device.Model("solar-sun").(*models.SunModel)
	irradiance := sunModel.Irradiance()

	// Calculate PV power
	area := float32(50.0)
	efficiency := float32(0.18)
	pvPower := irradiance * area * efficiency / 1000.0 * 2 // 2 arrays

	// Get auxiliary loads
	loadAux := b.device.Model("load-aux").(*models.LoadModel)
	loadStation := b.device.Model("load-station").(*models.LoadModel)
	totalLoad := loadAux.CurrentLoad() + loadStation.CurrentLoad()

	// Net power = generation - consumption
	// Positive = export to grid, Negative = import from grid
	netPower := pvPower - totalLoad

	// Apply breaker state
	if b.gridBreaker.IsOpen() {
		netPower = 0 // Islanded - no power flow to grid
	}

	// Grid is the infinite bus at 69kV
	pccVoltage := b.grid.Voltage()
	freq := b.grid.Frequency()

	// Reactive power (simplified)
	var reactiveQ float32 = 0.0

	// Calculate power factor
	var pf float32 = 1.0
	if netPower > 0.1 || netPower < -0.1 {
		pf = 0.95
	}

	// Accumulate energy (kWh)
	b.tickHours = 0.1 / 3600.0 // 100ms in hours
	if netPower > 0 {
		b.energyExport += netPower * b.tickHours
	} else {
		b.energyImport += (-netPower) * b.tickHours
	}

	// Write to memory
	b.device.Memory().WriteFloat32("input_registers", 0, pccVoltage)       // Voltage (69kV)
	b.device.Memory().WriteFloat32("input_registers", 4, freq)             // Frequency
	b.device.Memory().WriteFloat32("input_registers", 8, netPower)         // Active Power (signed)
	b.device.Memory().WriteFloat32("input_registers", 12, reactiveQ)       // Reactive Power
	b.device.Memory().WriteFloat32("input_registers", 16, pf)             // Power Factor
	b.device.Memory().WriteFloat32("input_registers", 20, b.energyExport) // Energy Export
	b.device.Memory().WriteFloat32("input_registers", 24, b.energyImport)  // Energy Import
	b.device.Memory().WriteFloat32("input_registers", 28, netPower)        // Direction (+export, -import)
}

// plantBehavior manages power flow through the plant
type plantBehavior struct {
	device       *device.Device
	collectorBus *models.BusModel
	gridBreaker  *models.BreakerModel
	tickCounter  int
}

func (b *plantBehavior) ID() string                        { return "plant_control" }
func (b *plantBehavior) Attach(d *device.Device)          { b.device = d }
func (b *plantBehavior) Detach()                          { b.device = nil }

func (b *plantBehavior) Tick() {
	if b.device == nil {
		return
	}

	b.tickCounter++

	// Get PV generation from sun model
	sunModel := b.device.Model("solar-sun").(*models.SunModel)
	irradiance := sunModel.Irradiance()

	// Calculate PV power
	area := float32(50.0)
	efficiency := float32(0.18)
	pv1Power := irradiance * area * efficiency / 1000.0
	pv2Power := irradiance * area * efficiency / 1000.0

	// Update PV array models
	pv1 := b.device.Model("pv-array-1").(*models.PVArrayModel)
	pv2 := b.device.Model("pv-array-2").(*models.PVArrayModel)
	pv1.SetPower(pv1Power)
	pv2.SetPower(pv2Power)
	totalPV := pv1Power + pv2Power

	// Get auxiliary loads
	loadAux := b.device.Model("load-aux").(*models.LoadModel)
	loadStation := b.device.Model("load-station").(*models.LoadModel)
	totalLoad := loadAux.CurrentLoad() + loadStation.CurrentLoad()

	// Inject PV to collector bus
	b.collectorBus.InjectPower(totalPV)

	// Withdraw auxiliary loads from collector bus
	b.collectorBus.WithdrawPower(totalLoad)

	// Auto breaker operations for demonstration
	if b.tickCounter == 20 && !b.gridBreaker.IsOpen() {
		b.gridBreaker.Open()
		fmt.Println("\n*** GRID BREAKER OPENED (Islanded) ***")
	}
	if b.tickCounter == 35 && b.gridBreaker.IsOpen() {
		b.gridBreaker.Close()
		fmt.Println("\n*** GRID BREAKER CLOSED (Grid Connected) ***")
	}

	// Write plant status
	var breakerStatus float32 = 0
	if b.gridBreaker.IsOpen() {
		breakerStatus = 1
	}
	b.device.Memory().WriteFloat32("status", 0, breakerStatus)
	b.device.Memory().WriteFloat32("input_registers", 0, totalPV)
	b.device.Memory().WriteFloat32("input_registers", 4, totalLoad)
	b.device.Memory().WriteFloat32("input_registers", 8, totalPV-totalLoad)
}
