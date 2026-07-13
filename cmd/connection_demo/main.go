// Package main demonstrates the electrical connection topology.
package main

import (
	"fmt"

	"github.com/tamzrod/forge/topology"
)

func main() {
	fmt.Println("Forge Electrical Connection Topology Demonstration")
	fmt.Println("=================================================")
	fmt.Println()

	// Create network
	net := topology.NewNetwork()

	// Add buses at different voltage levels
	net.AddBus(topology.NewBus("hv-bus", "69kV PCC", 69000))
	net.AddBus(topology.NewBus("lv-bus", "480V Collector", 480))
	net.AddBus(topology.NewBus("aux-bus", "Auxiliary Bus", 208))

	fmt.Println("Buses Created:")
	fmt.Println("--------------")
	for _, bus := range net.Buses() {
		fmt.Printf("  %s: %.0fV (%s)\n", bus.Name, bus.NominalVoltage, bus.ID)
	}

	// Add branches
	net.AddBranch(topology.NewBranch("tx", "Transformer", net.Bus("hv-bus"), net.Bus("lv-bus")))
	net.AddBranch(topology.NewBranch("aux-tx", "Aux Transformer", net.Bus("lv-bus"), net.Bus("aux-bus")))

	// Add switches
	gridBreaker := topology.NewSwitch("grid-breaker", "Grid Breaker", topology.SwitchTypeBreaker)
	gridBreaker.SetBranch(net.Branch("tx"))
	net.AddSwitch(gridBreaker)

	fmt.Println("\nBranches Created:")
	fmt.Println("-----------------")
	for _, branch := range net.Branches() {
		fmt.Printf("  %s: %s -> %s\n", branch.Name, branch.FromBus.ID, branch.ToBus.ID)
	}

	// Create terminals for generators (Source role)
	gen1Terminal := topology.NewSourceTerminal("solar-1-output", "solar-gen-1", "output", 480)
	gen2Terminal := topology.NewSourceTerminal("solar-2-output", "solar-gen-2", "output", 480)
	gridTerminal := topology.NewSourceTerminal("grid-output", "utility-grid", "grid", 69000)

	// Create terminal for loads (Destination role)
	load1Terminal := topology.NewDestinationTerminal("factory-input", "factory-1", "input", 480)
	load2Terminal := topology.NewDestinationTerminal("station-input", "station-service", "input", 208)

	// Create terminal for meters (Observation role)
	pccMeterTerminal := topology.NewObservationTerminal("pcc-meter-t", "pcc-meter", "meter", 69000)
	collectorMeterTerminal := topology.NewObservationTerminal("collector-meter-t", "collector-meter", "meter", 480)

	fmt.Println("\nTerminals Created:")
	fmt.Println("-------------------")
	terminals := []struct {
		t    *topology.Terminal
		desc string
	}{
		{gen1Terminal, "Solar Generator 1"},
		{gen2Terminal, "Solar Generator 2"},
		{gridTerminal, "Utility Grid"},
		{load1Terminal, "Factory Load"},
		{load2Terminal, "Station Service Load"},
		{pccMeterTerminal, "PCC Meter"},
		{collectorMeterTerminal, "Collector Meter"},
	}
	for _, item := range terminals {
		fmt.Printf("  %s: Role=%s, Type=%s, Voltage=%.0fV\n",
			item.desc, item.t.Role, item.t.Type, item.t.Voltage)
	}

	// Add terminals to network
	net.AddTerminal(gen1Terminal)
	net.AddTerminal(gen2Terminal)
	net.AddTerminal(gridTerminal)
	net.AddTerminal(load1Terminal)
	net.AddTerminal(load2Terminal)
	net.AddTerminal(pccMeterTerminal)
	net.AddTerminal(collectorMeterTerminal)

	// Connect terminals to buses
	fmt.Println("\n--- Connection Validation ---")
	fmt.Println()

	// Valid connections
	validConnections := []struct {
		terminal *topology.Terminal
		busID    topology.ID
		desc     string
	}{
		{gridTerminal, "hv-bus", "Grid -> HV Bus"},
		{gen1Terminal, "lv-bus", "Solar 1 -> LV Bus"},
		{gen2Terminal, "lv-bus", "Solar 2 -> LV Bus"},
		{load1Terminal, "lv-bus", "Factory -> LV Bus"},
		{load2Terminal, "aux-bus", "Station -> Aux Bus"},
		{pccMeterTerminal, "hv-bus", "PCC Meter -> HV Bus"},
		{collectorMeterTerminal, "lv-bus", "Collector Meter -> LV Bus"},
	}

	for _, conn := range validConnections {
		err := net.ConnectTerminal(conn.terminal.ID, conn.busID)
		if err != nil {
			fmt.Printf("  ✗ %s: ERROR - %v\n", conn.desc, err)
		} else {
			fmt.Printf("  ✓ %s: Connected successfully\n", conn.desc)
		}
	}

	// Invalid connection: LV terminal to HV bus
	fmt.Println()
	fmt.Println("--- Invalid Connection Test ---")
	invalidTerminal := topology.NewSourceTerminal("invalid-output", "invalid-gen", "output", 480)
	net.AddTerminal(invalidTerminal)
	err := net.ConnectTerminal(invalidTerminal.ID, "hv-bus")
	if err != nil {
		fmt.Printf("  ✗ LV generator -> HV Bus: ERROR - %v\n", err)
	} else {
		fmt.Printf("  ✓ LV generator -> HV Bus: Connected (should have failed!)\n")
	}

	// Print connection summary
	fmt.Println()
	fmt.Println("--- Connection Summary ---")
	fmt.Println()
	fmt.Printf("Total Terminals: %d\n", len(net.Terminals()))
	fmt.Printf("Source Terminals (Generators): %d\n", len(net.SourceTerminals()))
	fmt.Printf("Destination Terminals (Loads): %d\n", len(net.DestinationTerminals()))
	fmt.Printf("Observation Terminals (Meters): %d\n", len(net.ObservationTerminals()))
	fmt.Printf("Through Terminals (Transformers): %d\n", len(net.ThroughTerminals()))

	// Print bus connections
	fmt.Println()
	fmt.Println("--- Bus Connections ---")
	for _, bus := range net.Buses() {
		fmt.Printf("\nBus: %s (%.0fV)\n", bus.Name, bus.NominalVoltage)
		for _, t := range bus.Terminals() {
			fmt.Printf("  ├── %s (%s, %s) -> Entity: %s\n",
				t.Name, t.Role, t.Type, t.EntityID)
		}
	}

	// Network diagram
	fmt.Println()
	fmt.Println("=================================================")
	fmt.Println("Connection Diagram")
	fmt.Println("=================================================")
	fmt.Println(`
                    ┌─────────────────────────────────────┐
                    │         UTILITY GRID                │
                    │  (Source Terminal, HV, 69kV)         │
                    └──────────────┬──────────────────────┘
                                   │
                    ┌──────────────▼──────────────────────┐
                    │         69kV PCC BUS                 │
                    │  ┌─ Grid Terminal (Source)           │
                    │  └─ PCC Meter Terminal (Observation) │
                    └──────────────┬──────────────────────┘
                                   │
                    ┌──────────────▼──────────────────────┐
                    │         TRANSFORMER                  │
                    │  HV Terminal ──> LV Terminal         │
                    └──────────────┬──────────────────────┘
                                   │
                    ┌──────────────▼──────────────────────┐
                    │         480V COLLECTOR BUS          │
                    │                                      │
                    │  ┌─ Solar 1 Terminal (Source)       │
                    │  ├─ Solar 2 Terminal (Source)       │
                    │  ├─ Factory Terminal (Destination)   │
                    │  └─ Collector Meter (Observation)   │
                    └──────────────┬──────────────────────┘
                                   │
                    ┌──────────────▼──────────────────────┐
                    │         AUX TRANSFORMER             │
                    └──────────────┬──────────────────────┘
                                   │
                    ┌──────────────▼──────────────────────┐
                    │         208V AUX BUS               │
                    │  └─ Station Service Terminal (Dest)  │
                    └─────────────────────────────────────┘`)

	fmt.Println()
	fmt.Println("=================================================")
	fmt.Println("Summary: Electrical Connection Topology")
	fmt.Println("=================================================")
	fmt.Println()
	fmt.Println("Terminal Roles:")
	fmt.Println("  Source:       Injects power (generators)")
	fmt.Println("  Destination:  Withdraws power (loads)")
	fmt.Println("  Observation:  Measures without affecting (meters)")
	fmt.Println("  Through:      Passes power through (transformers)")
	fmt.Println()
	fmt.Println("Terminal Types:")
	fmt.Println("  HV:   High voltage (> 1000V)")
	fmt.Println("  LV:   Low voltage (<= 1000V)")
	fmt.Println("  Grid: Utility grid connection")
	fmt.Println()
	fmt.Println("Connection Validation:")
	fmt.Println("  - Terminal voltage must match bus voltage")
	fmt.Println("  - Observation terminals can connect to any bus")
	fmt.Println("  - Each terminal connects to exactly one bus")
}
