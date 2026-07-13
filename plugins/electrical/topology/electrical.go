// Package topology provides electrical network topology modeling.
//
// This package is part of the electrical plugin and contains
// electrical-specific topology types.
package topology

import (
	"fmt"
	"sync"

	"github.com/tamzrod/forge/world"
)

// TerminalRole defines what a terminal does in the electrical network.
type TerminalRole int

const (
	// TerminalRoleSource injects power into the network (generators).
	TerminalRoleSource TerminalRole = iota
	// TerminalRoleDestination withdraws power from the network (loads).
	TerminalRoleDestination
	// TerminalRoleThrough passes power through (transformers, cables).
	TerminalRoleThrough
	// TerminalRoleObservation measures power without affecting the network (meters).
	TerminalRoleObservation
)

func (r TerminalRole) String() string {
	switch r {
	case TerminalRoleSource:
		return "Source"
	case TerminalRoleDestination:
		return "Destination"
	case TerminalRoleThrough:
		return "Through"
	case TerminalRoleObservation:
		return "Observation"
	default:
		return "Unknown"
	}
}

// TerminalType defines the voltage classification of a terminal.
type TerminalType int

const (
	// TerminalTypeHV is high voltage (> 1kV).
	TerminalTypeHV TerminalType = iota
	// TerminalTypeLV is low voltage (<= 1kV).
	TerminalTypeLV
	// TerminalTypeGrid is utility grid connection.
	TerminalTypeGrid
)

func (t TerminalType) String() string {
	switch t {
	case TerminalTypeHV:
		return "HV"
	case TerminalTypeLV:
		return "LV"
	case TerminalTypeGrid:
		return "Grid"
	default:
		return "Unknown"
	}
}

// ID uniquely identifies a topology element.
type ID string

// Bus represents an electrical node where conductors connect.
type Bus struct {
	ID             ID
	Name           string
	NominalVoltage float32 // V
	mu             sync.RWMutex
	terminals      map[ID]*Terminal
	branches       map[ID]*Branch
}

// NewBus creates a new bus.
func NewBus(id ID, name string, nominalVoltage float32) *Bus {
	return &Bus{
		ID:             id,
		Name:           name,
		NominalVoltage: nominalVoltage,
		terminals:      make(map[ID]*Terminal),
		branches:       make(map[ID]*Branch),
	}
}

// AddTerminal connects a terminal to this bus.
func (b *Bus) AddTerminal(t *Terminal) {
	b.mu.Lock()
	defer b.mu.Unlock()
	t.bus = b
	b.terminals[t.ID] = t
}

// RemoveTerminal disconnects a terminal from this bus.
func (b *Bus) RemoveTerminal(tid ID) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if t, ok := b.terminals[tid]; ok {
		t.bus = nil
	}
	delete(b.terminals, tid)
}

// Terminals returns all terminals connected to this bus.
func (b *Bus) Terminals() []*Terminal {
	b.mu.RLock()
	defer b.mu.RUnlock()
	result := make([]*Terminal, 0, len(b.terminals))
	for _, t := range b.terminals {
		result = append(result, t)
	}
	return result
}

// AddBranch connects a branch to this bus.
func (b *Bus) AddBranch(br *Branch) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.branches[br.ID] = br
}

// RemoveBranch disconnects a branch from this bus.
func (b *Bus) RemoveBranch(bid ID) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.branches, bid)
}

// Branches returns all branches connected to this bus.
func (b *Bus) Branches() []*Branch {
	b.mu.RLock()
	defer b.mu.RUnlock()
	result := make([]*Branch, 0, len(b.branches))
	for _, br := range b.branches {
		result = append(result, br)
	}
	return result
}

// ConnectedEntities returns all entities connected to this bus.
func (b *Bus) ConnectedEntities() []world.EntityID {
	b.mu.RLock()
	defer b.mu.RUnlock()
	result := make([]world.EntityID, 0)
	seen := make(map[world.EntityID]bool)
	for _, t := range b.terminals {
		if !seen[t.EntityID] {
			result = append(result, t.EntityID)
			seen[t.EntityID] = true
		}
	}
	return result
}

// Terminal represents a typed connection point on an entity.
type Terminal struct {
	ID       ID
	EntityID world.EntityID
	Name     string
	Role     TerminalRole
	Type     TerminalType
	Voltage  float32
	bus      *Bus
}

// NewTerminal creates a new terminal.
func NewTerminal(id ID, entityID world.EntityID, name string, role TerminalRole, terminalType TerminalType, voltage float32) *Terminal {
	return &Terminal{
		ID:       id,
		EntityID: entityID,
		Name:     name,
		Role:     role,
		Type:     terminalType,
		Voltage:  voltage,
	}
}

// NewSourceTerminal creates a new source terminal (for generators).
func NewSourceTerminal(id ID, entityID world.EntityID, name string, voltage float32) *Terminal {
	return NewTerminal(id, entityID, name, TerminalRoleSource, voltageToTerminalType(voltage), voltage)
}

// NewDestinationTerminal creates a new destination terminal (for loads).
func NewDestinationTerminal(id ID, entityID world.EntityID, name string, voltage float32) *Terminal {
	return NewTerminal(id, entityID, name, TerminalRoleDestination, voltageToTerminalType(voltage), voltage)
}

// NewObservationTerminal creates a new observation terminal (for meters).
func NewObservationTerminal(id ID, entityID world.EntityID, name string, voltage float32) *Terminal {
	return NewTerminal(id, entityID, name, TerminalRoleObservation, voltageToTerminalType(voltage), voltage)
}

// NewThroughTerminal creates a new through terminal (for transformers).
func NewThroughTerminal(id ID, entityID world.EntityID, name string, terminalType TerminalType, voltage float32) *Terminal {
	return NewTerminal(id, entityID, name, TerminalRoleThrough, terminalType, voltage)
}

// Bus returns the bus this terminal is connected to.
func (t *Terminal) Bus() *Bus {
	return t.bus
}

// IsConnected returns true if this terminal is connected to a bus.
func (t *Terminal) IsConnected() bool {
	return t.bus != nil
}

func voltageToTerminalType(voltage float32) TerminalType {
	if voltage >= 1000 {
		return TerminalTypeHV
	}
	return TerminalTypeLV
}

// Branch represents a connection between two buses.
type Branch struct {
	ID          ID
	Name        string
	FromBus     *Bus
	ToBus       *Bus
	FromTerminal *Terminal
	ToTerminal   *Terminal
	switchDevice *Switch
}

// NewBranch creates a new branch.
func NewBranch(id ID, name string, fromBus, toBus *Bus) *Branch {
	return &Branch{
		ID:      id,
		Name:    name,
		FromBus: fromBus,
		ToBus:   toBus,
	}
}

// OtherBus returns the other bus in the branch.
func (b *Branch) OtherBus(bus *Bus) *Bus {
	if b.FromBus == bus {
		return b.ToBus
	}
	return b.FromBus
}

// SetSwitch sets the switch device in this branch.
func (b *Branch) SetSwitch(sw *Switch) {
	b.switchDevice = sw
}

// SwitchDevice returns the switch in this branch.
func (b *Branch) SwitchDevice() *Switch {
	return b.switchDevice
}

// SwitchType defines the type of switch.
type SwitchType int

const (
	// SwitchTypeBreaker is a circuit breaker.
	SwitchTypeBreaker SwitchType = iota
	// SwitchTypeRecloser is a recloser.
	SwitchTypeRecloser
	// SwitchTypeSectionalizer is a sectionalizer.
	SwitchTypeSectionalizer
)

// Switch is a switching device that can interrupt flow.
type Switch struct {
	ID     ID
	Name   string
	Type   SwitchType
	isOpen bool
	mu     sync.RWMutex
	branch *Branch
}

// NewSwitch creates a new switch.
func NewSwitch(id ID, name string, switchType SwitchType) *Switch {
	return &Switch{
		ID:   id,
		Name: name,
		Type: switchType,
	}
}

// SetBranch sets the branch this switch is in.
func (s *Switch) SetBranch(branch *Branch) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.branch = branch
}

// Branch returns the branch this switch is in.
func (s *Switch) Branch() *Branch {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.branch
}

// Open opens the switch.
func (s *Switch) Open() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.isOpen = true
}

// Close closes the switch.
func (s *Switch) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.isOpen = false
}

// IsOpen returns true if the switch is open.
func (s *Switch) IsOpen() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isOpen
}

// Network represents an electrical network.
type Network struct {
	mu        sync.RWMutex
	buses     map[ID]*Bus
	branches  map[ID]*Branch
	terminals map[ID]*Terminal
	switches  map[ID]*Switch
}

// NewNetwork creates a new electrical network.
func NewNetwork() *Network {
	return &Network{
		buses:     make(map[ID]*Bus),
		branches:  make(map[ID]*Branch),
		terminals: make(map[ID]*Terminal),
		switches:  make(map[ID]*Switch),
	}
}

// AddBus adds a bus to the network.
func (n *Network) AddBus(bus *Bus) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.buses[bus.ID] = bus
}

// Bus returns a bus by ID.
func (n *Network) Bus(id ID) *Bus {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.buses[id]
}

// AddBranch adds a branch to the network.
func (n *Network) AddBranch(branch *Branch) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.branches[branch.ID] = branch
	branch.FromBus.AddBranch(branch)
	branch.ToBus.AddBranch(branch)
}

// Branch returns a branch by ID.
func (n *Network) Branch(id ID) *Branch {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.branches[id]
}

// AddTerminal adds a terminal to the network.
func (n *Network) AddTerminal(terminal *Terminal) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.terminals[terminal.ID] = terminal
}

// Terminal returns a terminal by ID.
func (n *Network) Terminal(id ID) *Terminal {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.terminals[id]
}

// AddSwitch adds a switch to the network.
func (n *Network) AddSwitch(sw *Switch) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.switches[sw.ID] = sw
}

// Switch returns a switch by ID.
func (n *Network) Switch(id ID) *Switch {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.switches[id]
}

// Buses returns all buses.
func (n *Network) Buses() []*Bus {
	n.mu.RLock()
	defer n.mu.RUnlock()
	result := make([]*Bus, 0, len(n.buses))
	for _, b := range n.buses {
		result = append(result, b)
	}
	return result
}

// Branches returns all branches.
func (n *Network) Branches() []*Branch {
	n.mu.RLock()
	defer n.mu.RUnlock()
	result := make([]*Branch, 0, len(n.branches))
	for _, br := range n.branches {
		result = append(result, br)
	}
	return result
}

// Switches returns all switches.
func (n *Network) Switches() []*Switch {
	n.mu.RLock()
	defer n.mu.RUnlock()
	result := make([]*Switch, 0, len(n.switches))
	for _, s := range n.switches {
		result = append(result, s)
	}
	return result
}

// String returns a string representation of the network.
func (n *Network) String() string {
	n.mu.RLock()
	defer n.mu.RUnlock()

	result := fmt.Sprintf("Network: %d buses, %d branches, %d switches\n",
		len(n.buses), len(n.branches), len(n.switches))

	for _, bus := range n.buses {
		connected := make([]string, 0)
		for _, br := range bus.branches {
			connected = append(connected, string(br.OtherBus(bus).ID))
		}
		result += fmt.Sprintf("  Bus %s (%s) -> [%s]\n", bus.ID, bus.Name, connected)
	}

	return result
}
