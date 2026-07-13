// Package topology provides electrical network topology modeling.
// Topology owns connectivity - entities own behavior.
package topology

import (
	"github.com/tamzrod/forge/world"
)

// Builder constructs electrical networks.
type Builder struct {
	network *Network
}

// NewBuilder creates a new topology builder.
func NewBuilder() *Builder {
	return &Builder{
		network: NewNetwork(),
	}
}

// Network returns the constructed network.
func (b *Builder) Network() *Network {
	return b.network
}

// AddBus adds a bus to the network.
func (b *Builder) AddBus(id ID, name string, voltage float32) *Builder {
	b.network.AddBus(NewBus(id, name, voltage))
	return b
}

// AddBranch adds a branch between two buses.
func (b *Builder) AddBranch(id ID, name string, fromID, toID ID) *Builder {
	from := b.network.Bus(fromID)
	to := b.network.Bus(toID)
	if from == nil || to == nil {
		return b
	}
	b.network.AddBranch(NewBranch(id, name, from, to))
	return b
}

// AddBreaker adds a breaker in a branch.
func (b *Builder) AddBreaker(breakerID, branchID, name string) *Builder {
	breaker := NewSwitch(ID(breakerID), name, SwitchTypeBreaker)
	b.network.AddSwitch(breaker)

	branch := b.network.Branch(ID(branchID))
	if branch != nil {
		breaker.SetBranch(branch)
		branch.SetSwitch(breaker)
	}
	return b
}

// AddTransformer adds a transformer between two buses.
// Transformers have high and low voltage sides.
func (b *Builder) AddTransformer(id ID, name string, highBusID, lowBusID ID, ratio float32) *Builder {
	highBus := b.network.Bus(highBusID)
	lowBus := b.network.Bus(lowBusID)
	if highBus == nil || lowBus == nil {
		return b
	}

	// Create branch for transformer
	branch := NewBranch(id, name, highBus, lowBus)
	b.network.AddBranch(branch)

	// Transformers are not typically switched, so no switch needed
	return b
}

// ConnectEntity connects an entity to a bus with a terminal.
func (b *Builder) ConnectEntity(entityID world.EntityID, busID ID, terminalName string) *Builder {
	bus := b.network.Bus(busID)
	if bus == nil {
		return b
	}

	terminalID := ID(string(entityID) + "-" + terminalName)
	terminal := NewTerminal(terminalID, entityID, terminalName)
	terminal.bus = bus
	bus.AddTerminal(terminal)
	b.network.AddTerminal(terminal)

	return b
}

// ConnectEntityToBreaker connects an entity through a breaker.
func (b *Builder) ConnectEntityToBreaker(entityID world.EntityID, breakerID ID, busID ID, terminalName string, side string) *Builder {
	bus := b.network.Bus(busID)
	breaker := b.network.Switch(breakerID)
	if bus == nil || breaker == nil || breaker.branch == nil {
		return b
	}

	terminalID := ID(string(entityID) + "-" + terminalName + "-" + side)
	terminal := NewTerminal(terminalID, entityID, terminalName+"-"+side)

	// Determine which bus the terminal connects to based on side
	var targetBus *Bus
	if side == "from" && breaker.branch.FromBus == bus {
		targetBus = bus
	} else if side == "to" && breaker.branch.ToBus == bus {
		targetBus = bus
	} else {
		// Connect to whichever bus matches
		targetBus = bus
	}

	terminal.bus = targetBus
	targetBus.AddTerminal(terminal)
	b.network.AddTerminal(terminal)

	return b
}

// Build creates a simple radial network:
// Grid -> Main Breaker -> HV Bus -> Transformer -> LV Bus -> CB Breaker -> Collector Bus -> PV/Meter/Load
func (b *Builder) BuildSimpleRadial() *Builder {
	// Buses
	b.AddBus("hv-bus", "69kV High Voltage Bus", 69000)
	b.AddBus("lv-bus", "480V Low Voltage Bus", 480)
	b.AddBus("collector-bus", "Collector Bus", 480)

	// Branches
	b.AddBranch("main-tx", "Main Transformer", "hv-bus", "lv-bus")
	b.AddBranch("cb-tx", "Circuit Breaker Transformer", "lv-bus", "collector-bus")

	// Switches (breakers)
	b.AddBreaker("grid-breaker", "main-tx", "Grid Breaker")
	b.AddBreaker("cb-breaker", "cb-tx", "Circuit Breaker")

	// Connect entities
	b.ConnectEntity("utility-grid", "hv-bus", "grid")
	b.ConnectEntity("transformer", "hv-bus", "primary")
	b.ConnectEntity("transformer", "lv-bus", "secondary")
	b.ConnectEntity("pcc-meter", "hv-bus", "meter")
	b.ConnectEntity("collector-bus", "lv-bus", "bus")
	b.ConnectEntity("pv-array-1", "collector-bus", "ac")
	b.ConnectEntity("pv-array-2", "collector-bus", "ac")
	b.ConnectEntity("aux-load", "collector-bus", "load")
	b.ConnectEntity("station-load", "collector-bus", "load")

	return b
}

// Example: Build a mesh network with multiple sources
func (b *Builder) BuildMeshNetwork() *Builder {
	// Buses
	b.AddBus("bus-a", "Bus A", 138000)
	b.AddBus("bus-b", "Bus B", 138000)
	b.AddBus("bus-c", "Bus C", 69000)

	// Branches (mesh)
	b.AddBranch("line-1", "Transmission Line 1", "bus-a", "bus-b")
	b.AddBranch("line-2", "Transmission Line 2", "bus-a", "bus-b")
	b.AddBranch("tx-ab-c", "Transformer AB-C", "bus-b", "bus-c")

	// Breakers
	b.AddBreaker("cb-a1", "line-1", "Breaker A1")
	b.AddBreaker("cb-a2", "line-2", "Breaker A2")
	b.AddBreaker("cb-c", "tx-ab-c", "Breaker C")

	// Sources
	b.ConnectEntity("grid-1", "bus-a", "grid")
	b.ConnectEntity("grid-2", "bus-b", "grid")

	// Load
	b.ConnectEntity("load-1", "bus-c", "load")

	return b
}
