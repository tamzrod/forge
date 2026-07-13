// Package main provides the topology demonstration.
package main

import (
	"fmt"

	"github.com/tamzrod/forge/topology"
)

func main() {
	fmt.Println("Forge Electrical Topology Demonstration")
	fmt.Println("======================================")
	fmt.Println()

	// Build a simple radial network
	builder := topology.NewBuilder()
	builder.BuildSimpleRadial()
	net := builder.Network()

	fmt.Println("Network Structure:")
	fmt.Println("------------------")
	fmt.Println(net.String())

	// Query examples
	fmt.Println("Topology Queries:")
	fmt.Println("----------------")

	// Get a bus
	hvBus := net.Bus("hv-bus")
	if hvBus != nil {
		fmt.Printf("\nBus: %s\n", hvBus.Name)
		fmt.Printf("  Nominal Voltage: %.0f V\n", hvBus.NominalVoltage)
		fmt.Printf("  Connected Entities: %v\n", net.EntitiesConnectedTo(hvBus))
	}

	// Get connected buses
	lvBus := net.Bus("lv-bus")
	if lvBus != nil {
		connected := net.ConnectedTo(lvBus)
		fmt.Printf("\nBuses connected to %s:\n", lvBus.Name)
		for _, bus := range connected {
			fmt.Printf("  - %s\n", bus.Name)
		}
	}

	// Get upstream/downstream
	fmt.Printf("\nUpstream from collector-bus:\n")
	collectorBus := net.Bus("collector-bus")
	if collectorBus != nil {
		upstream := net.EntitiesUpstream(collectorBus)
		fmt.Printf("  Entities: %v\n", upstream)
	}

	fmt.Printf("\nDownstream from hv-bus:\n")
	if hvBus != nil {
		downstream := net.EntitiesDownstream(hvBus)
		fmt.Printf("  Entities: %v\n", downstream)
	}

	// Check islands
	fmt.Println("\nIslands:")
	islands := net.Islands()
	for _, island := range islands {
		fmt.Printf("  %s: %d buses\n", island.ID, len(island.Buses))
		for _, bus := range island.Buses {
			fmt.Printf("    - %s\n", bus.Name)
		}
	}

	// Simulate breaker operations
	fmt.Println("\nBreaker Operations:")
	fmt.Println("------------------")

	gridBreaker := net.Switch("grid-breaker")
	if gridBreaker != nil {
		fmt.Printf("\nGrid Breaker: %s\n", gridBreaker.Name)
		fmt.Printf("  Initial State: CLOSED=%v (OPEN=%v)\n", !gridBreaker.IsOpen(), gridBreaker.IsOpen())

		// Open breaker
		gridBreaker.Open()
		fmt.Printf("  After Open: CLOSED=%v (OPEN=%v)\n", !gridBreaker.IsOpen(), gridBreaker.IsOpen())

		// Check islands when breaker is open
		islands = net.Islands()
		fmt.Printf("  Islands when breaker is open: %d\n", len(islands))

		// Check isolation
		isolated := net.IsolatedIf("grid-breaker")
		fmt.Printf("  Entities isolated when grid breaker opens: %v\n", isolated)

		// Close breaker
		gridBreaker.Close()
		fmt.Printf("  After Close: CLOSED=%v (OPEN=%v)\n", !gridBreaker.IsOpen(), gridBreaker.IsOpen())
	}

	// Test path finding
	fmt.Println("\nPath Finding:")
	fmt.Println("-------------")
	path := net.PathBetween(net.Bus("collector-bus"), net.Bus("hv-bus"))
	if path != nil {
		fmt.Printf("Path from collector-bus to hv-bus:\n")
		for i, bus := range path {
			fmt.Printf("  %d. %s\n", i+1, bus.Name)
		}
	}

	// Close breakers and check full connectivity
	fmt.Println("\nFull Network (all breakers closed):")
	fmt.Println("------------------------------------")
	gridBreaker.Close()
	cbBreaker := net.Switch("cb-breaker")
	if cbBreaker != nil {
		cbBreaker.Close()
	}

	islands = net.Islands()
	fmt.Printf("Number of islands: %d\n", len(islands))
	for _, island := range islands {
		fmt.Printf("Island %s: %d buses - ", island.ID, len(island.Buses))
		for i, bus := range island.Buses {
			if i > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("%s", bus.Name)
		}
		fmt.Println()
	}

	fmt.Println("\nTopology Demonstration Complete")
}
