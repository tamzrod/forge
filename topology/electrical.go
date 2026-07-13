// Package topology provides electrical network topology modeling.
// Topology owns connectivity - entities own behavior.
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

// String returns the terminal role as a string.
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

// String returns the terminal type as a string.
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

// VoltageLevelForTerminalType returns the voltage level for a terminal type.
func VoltageLevelForTerminalType(t TerminalType) VoltageLevel {
	switch t {
	case TerminalTypeLV:
		return VoltageLevelLow
	case TerminalTypeHV:
		return VoltageLevelMedium // Treat HV as medium for simplicity
	case TerminalTypeGrid:
		return VoltageLevelHigh // Grid is typically high voltage
	default:
		return VoltageLevelLow
	}
}

// ID uniquely identifies a topology element.
type ID string

// Bus represents an electrical node where conductors connect.
// Multiple branches and equipment terminals connect at a bus.
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
// Each entity can have multiple terminals (e.g., a transformer has HV and LV).
type Terminal struct {
	ID       ID
	EntityID world.EntityID
	Name     string // e.g., "primary", "secondary", "output", "input"
	Role     TerminalRole  // Source, Destination, Through, Observation
	Type     TerminalType  // HV, LV, Grid
	Voltage  float32       // V - nominal voltage
	bus      *Bus
}

// NewTerminal creates a new terminal with role and type.
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

// IsSource returns true if this is a source terminal.
func (t *Terminal) IsSource() bool {
	return t.Role == TerminalRoleSource
}

// IsDestination returns true if this is a destination terminal.
func (t *Terminal) IsDestination() bool {
	return t.Role == TerminalRoleDestination
}

// IsThrough returns true if this is a through terminal.
func (t *Terminal) IsThrough() bool {
	return t.Role == TerminalRoleThrough
}

// IsObservation returns true if this is an observation terminal.
func (t *Terminal) IsObservation() bool {
	return t.Role == TerminalRoleObservation
}

// CanConnectTo returns true if this terminal can connect to the given bus.
func (t *Terminal) CanConnectTo(bus *Bus) bool {
	// Observation terminals can connect to any voltage
	if t.Role == TerminalRoleObservation {
		return true
	}
	// Through terminals (transformers) can connect to matching voltage
	if t.Role == TerminalRoleThrough {
		return t.Voltage == bus.NominalVoltage
	}
	// Source and Destination terminals must match voltage
	return t.Voltage == bus.NominalVoltage
}

// voltageToTerminalType converts voltage to terminal type.
func voltageToTerminalType(voltage float32) TerminalType {
	if voltage >= 1000 {
		return TerminalTypeHV
	}
	return TerminalTypeLV
}

// Branch represents a connection between two buses.
// Branches contain switching devices (breakers) that can interrupt flow.
type Branch struct {
	ID             ID
	Name           string
	FromBus       *Bus
	ToBus         *Bus
	FromTerminal  *Terminal // Terminal on source entity
	ToTerminal    *Terminal // Terminal on destination entity
	switchDevice   *Switch  // The switching device in this branch
	isEnergized    bool
}

// NewBranch creates a new branch between two buses.
func NewBranch(id ID, name string, from, to *Bus) *Branch {
	return &Branch{
		ID:       id,
		Name:     name,
		FromBus: from,
		ToBus:   to,
	}
}

// SetSwitch assigns a switching device to this branch.
func (br *Branch) SetSwitch(sw *Switch) {
	br.switchDevice = sw
}

// SetEnergized sets whether this branch is energized.
func (br *Branch) SetEnergized(energized bool) {
	br.isEnergized = energized
}

// IsEnergized returns whether this branch is energized.
func (br *Branch) IsEnergized() bool {
	return br.isEnergized
}

// SwitchDevice returns the switch in this branch.
func (br *Branch) SwitchDevice() *Switch {
	return br.switchDevice
}

// OtherBus returns the bus on the other end of the branch.
func (br *Branch) OtherBus(b *Bus) *Bus {
	if br.FromBus == b {
		return br.ToBus
	}
	return br.FromBus
}

// SwitchType represents the type of switching device.
type SwitchType int

const (
	SwitchTypeBreaker SwitchType = iota
	SwitchTypeRecloser
	SwitchTypeSectionalizer
	SwitchTypeDisconnector
)

// Switch represents a switching device in a branch.
type Switch struct {
	ID            ID
	Name          string
	Type          SwitchType
	branch        *Branch
	isOpen        bool
	operatingTime float32 // seconds
}

// NewSwitch creates a new switch.
func NewSwitch(id ID, name string, switchType SwitchType) *Switch {
	return &Switch{
		ID:            id,
		Name:          name,
		Type:          switchType,
		isOpen:        false,
		operatingTime:  0.05, // 50ms default
	}
}

// Open opens the switch.
func (s *Switch) Open() {
	s.isOpen = true
}

// Close closes the switch.
func (s *Switch) Close() {
	s.isOpen = false
}

// IsOpen returns true if the switch is open.
func (s *Switch) IsOpen() bool {
	return s.isOpen
}

// SetBranch assigns this switch to a branch.
func (s *Switch) SetBranch(br *Branch) {
	s.branch = br
}

// Branch returns the branch this switch is in.
func (s *Switch) Branch() *Branch {
	return s.branch
}

// VoltageLevel represents a voltage classification.
type VoltageLevel int

const (
	VoltageLevelLow VoltageLevel = iota // < 1kV
	VoltageLevelMedium                  // 1kV - 35kV
	VoltageLevelHigh                    // 35kV - 230kV
	VoltageLevelExtraHigh               // > 230kV
)

// VoltageLevel returns the voltage level classification.
func VoltageLevelFor(voltage float32) VoltageLevel {
	switch {
	case voltage < 1000:
		return VoltageLevelLow
	case voltage < 35000:
		return VoltageLevelMedium
	case voltage < 230000:
		return VoltageLevelHigh
	default:
		return VoltageLevelExtraHigh
	}
}

// Network represents the complete electrical network topology.
type Network struct {
	mu      sync.RWMutex
	buses   map[ID]*Bus
	branches map[ID]*Branch
	switches map[ID]*Switch
	terminals map[ID]*Terminal
}

// NewNetwork creates a new electrical network.
func NewNetwork() *Network {
	return &Network{
		buses:    make(map[ID]*Bus),
		branches: make(map[ID]*Branch),
		switches: make(map[ID]*Switch),
		terminals: make(map[ID]*Terminal),
	}
}

// AddBus adds a bus to the network.
func (n *Network) AddBus(b *Bus) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.buses[b.ID] = b
}

// RemoveBus removes a bus from the network.
func (n *Network) RemoveBus(id ID) {
	n.mu.Lock()
	defer n.mu.Unlock()
	delete(n.buses, id)
}

// Bus returns a bus by ID.
func (n *Network) Bus(id ID) *Bus {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.buses[id]
}

// AddBranch adds a branch to the network.
func (n *Network) AddBranch(br *Branch) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.branches[br.ID] = br
	br.FromBus.AddBranch(br)
	br.ToBus.AddBranch(br)
}

// RemoveBranch removes a branch from the network.
func (n *Network) RemoveBranch(id ID) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if br, ok := n.branches[id]; ok {
		br.FromBus.RemoveBranch(id)
		br.ToBus.RemoveBranch(id)
	}
	delete(n.branches, id)
}

// Branch returns a branch by ID.
func (n *Network) Branch(id ID) *Branch {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.branches[id]
}

// AddSwitch adds a switch to the network.
func (n *Network) AddSwitch(s *Switch) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.switches[s.ID] = s
}

// RemoveSwitch removes a switch from the network.
func (n *Network) RemoveSwitch(id ID) {
	n.mu.Lock()
	defer n.mu.Unlock()
	delete(n.switches, id)
}

// Switch returns a switch by ID.
func (n *Network) Switch(id ID) *Switch {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.switches[id]
}

// AddTerminal adds a terminal to the network (not connected).
func (n *Network) AddTerminal(t *Terminal) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.terminals[t.ID] = t
}

// RemoveTerminal removes a terminal from the network.
func (n *Network) RemoveTerminal(id ID) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if t, ok := n.terminals[id]; ok {
		if t.bus != nil {
			t.bus.RemoveTerminal(id)
		}
	}
	delete(n.terminals, id)
}

// Terminal returns a terminal by ID.
func (n *Network) Terminal(id ID) *Terminal {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.terminals[id]
}

// Terminals returns all terminals.
func (n *Network) Terminals() []*Terminal {
	n.mu.RLock()
	defer n.mu.RUnlock()
	result := make([]*Terminal, 0, len(n.terminals))
	for _, t := range n.terminals {
		result = append(result, t)
	}
	return result
}

// ConnectTerminal connects a terminal to a bus with validation.
// Returns error if connection is invalid.
func (n *Network) ConnectTerminal(terminalID ID, busID ID) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	terminal, ok := n.terminals[terminalID]
	if !ok {
		return fmt.Errorf("terminal %s not found", terminalID)
	}

	bus, ok := n.buses[busID]
	if !ok {
		return fmt.Errorf("bus %s not found", busID)
	}

	// Check if terminal is already connected
	if terminal.bus != nil {
		return fmt.Errorf("terminal %s is already connected to bus %s", terminalID, terminal.bus.ID)
	}

	// Validate connection
	if !terminal.CanConnectTo(bus) {
		return fmt.Errorf("terminal %s (%.0fV, %s) cannot connect to bus %s (%.0fV)",
			terminalID, terminal.Voltage, terminal.Role, busID, bus.NominalVoltage)
	}

	// Make connection
	terminal.bus = bus
	bus.AddTerminal(terminal)

	return nil
}

// DisconnectTerminal disconnects a terminal from its bus.
func (n *Network) DisconnectTerminal(terminalID ID) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	terminal, ok := n.terminals[terminalID]
	if !ok {
		return fmt.Errorf("terminal %s not found", terminalID)
	}

	if terminal.bus == nil {
		return fmt.Errorf("terminal %s is not connected", terminalID)
	}

	// Remove from bus
	terminal.bus.RemoveTerminal(terminalID)
	terminal.bus = nil

	return nil
}

// ConnectTerminalToBus connects a terminal directly to a bus (unsafe, no validation).
func (n *Network) ConnectTerminalToBus(terminal *Terminal, bus *Bus) {
	n.mu.Lock()
	defer n.mu.Unlock()
	terminal.bus = bus
	bus.AddTerminal(terminal)
}

// SourceTerminals returns all source terminals (generators).
func (n *Network) SourceTerminals() []*Terminal {
	n.mu.RLock()
	defer n.mu.RUnlock()
	result := make([]*Terminal, 0)
	for _, t := range n.terminals {
		if t.Role == TerminalRoleSource {
			result = append(result, t)
		}
	}
	return result
}

// DestinationTerminals returns all destination terminals (loads).
func (n *Network) DestinationTerminals() []*Terminal {
	n.mu.RLock()
	defer n.mu.RUnlock()
	result := make([]*Terminal, 0)
	for _, t := range n.terminals {
		if t.Role == TerminalRoleDestination {
			result = append(result, t)
		}
	}
	return result
}

// ObservationTerminals returns all observation terminals (meters).
func (n *Network) ObservationTerminals() []*Terminal {
	n.mu.RLock()
	defer n.mu.RUnlock()
	result := make([]*Terminal, 0)
	for _, t := range n.terminals {
		if t.Role == TerminalRoleObservation {
			result = append(result, t)
		}
	}
	return result
}

// ThroughTerminals returns all through terminals (transformers).
func (n *Network) ThroughTerminals() []*Terminal {
	n.mu.RLock()
	defer n.mu.RUnlock()
	result := make([]*Terminal, 0)
	for _, t := range n.terminals {
		if t.Role == TerminalRoleThrough {
			result = append(result, t)
		}
	}
	return result
}

// TerminalsByBus returns all terminals connected to a bus.
func (n *Network) TerminalsByBus(busID ID) []*Terminal {
	n.mu.RLock()
	defer n.mu.RUnlock()

	bus, ok := n.buses[busID]
	if !ok {
		return nil
	}

	return bus.Terminals()
}

// TerminalsByEntity returns all terminals belonging to an entity.
func (n *Network) TerminalsByEntity(entityID world.EntityID) []*Terminal {
	n.mu.RLock()
	defer n.mu.RUnlock()

	result := make([]*Terminal, 0)
	for _, t := range n.terminals {
		if t.EntityID == entityID {
			result = append(result, t)
		}
	}
	return result
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

// ConnectedTo returns all buses directly connected to this bus.
func (n *Network) ConnectedTo(bus *Bus) []*Bus {
	bus.mu.RLock()
	defer bus.mu.RUnlock()
	result := make([]*Bus, 0)
	for _, br := range bus.branches {
		if other := br.OtherBus(bus); other != nil {
			result = append(result, other)
		}
	}
	return result
}

// Upstream returns all buses upstream (toward sources) from this bus.
func (n *Network) Upstream(from *Bus) []*Bus {
	visited := make(map[ID]bool)
	return n.upstreamDFS(from, visited)
}

func (n *Network) upstreamDFS(bus *Bus, visited map[ID]bool) []*Bus {
	if visited[bus.ID] {
		return nil
	}
	visited[bus.ID] = true

	result := []*Bus{bus}
	for _, connected := range n.ConnectedTo(bus) {
		result = append(result, n.upstreamDFS(connected, visited)...)
	}
	return result
}

// Downstream returns all buses downstream (away from sources) from this bus.
func (n *Network) Downstream(from *Bus) []*Bus {
	visited := make(map[ID]bool)
	return n.downstreamDFS(from, visited)
}

func (n *Network) downstreamDFS(bus *Bus, visited map[ID]bool) []*Bus {
	if visited[bus.ID] {
		return nil
	}
	visited[bus.ID] = true

	result := make([]*Bus, 0)
	for _, connected := range n.ConnectedTo(bus) {
		result = append(result, bus)
		result = append(result, n.downstreamDFS(connected, visited)...)
	}
	return result
}

// EntitiesConnectedTo returns all entity IDs connected to this bus.
func (n *Network) EntitiesConnectedTo(bus *Bus) []world.EntityID {
	return bus.ConnectedEntities()
}

// EntitiesUpstream returns all entities upstream from this bus.
func (n *Network) EntitiesUpstream(bus *Bus) []world.EntityID {
	buses := n.Upstream(bus)
	return n.entitiesOnBuses(buses)
}

// EntitiesDownstream returns all entities downstream from this bus.
func (n *Network) EntitiesDownstream(bus *Bus) []world.EntityID {
	buses := n.Downstream(bus)
	return n.entitiesOnBuses(buses)
}

func (n *Network) entitiesOnBuses(buses []*Bus) []world.EntityID {
	busSet := make(map[ID]*Bus)
	for _, b := range buses {
		busSet[b.ID] = b
	}

	result := make([]world.EntityID, 0)
	seen := make(map[world.EntityID]bool)

	n.mu.RLock()
	for _, t := range n.terminals {
		if busSet[t.Bus().ID] != nil && !seen[t.EntityID] {
			result = append(result, t.EntityID)
			seen[t.EntityID] = true
		}
	}
	n.mu.RUnlock()

	return result
}

// IsolatedIf returns which entities become isolated if the specified switch opens.
func (n *Network) IsolatedIf(switchID ID) []world.EntityID {
	sw := n.Switch(switchID)
	if sw == nil || sw.branch == nil {
		return nil
	}

	br := sw.branch

	// Temporarily open the switch to calculate islands
	wasOpen := sw.IsOpen()
	if !wasOpen {
		sw.isOpen = true // Temporarily open
	}

	// Calculate islands with switch open
	islands := n.Islands()

	// Restore switch state
	sw.isOpen = wasOpen

	// Find which island contains the switch's from bus
	var switchIsland *Island
	for _, island := range islands {
		for _, bus := range island.Buses {
			if bus.ID == br.FromBus.ID {
				switchIsland = island
				break
			}
		}
		if switchIsland != nil {
			break
		}
	}

	// If switch island doesn't have source, those entities are isolated
	// For simplicity, entities not in the switch's island are isolated
	var isolatedEntities []world.EntityID
	if switchIsland != nil {
		// Entities in the same island as the from bus are not isolated
		// Entities in other islands are isolated
		isolatedBuses := make(map[ID]*Bus)
		for _, island := range islands {
			if island.ID != switchIsland.ID {
				for _, bus := range island.Buses {
					isolatedBuses[bus.ID] = bus
				}
			}
		}
		isolatedEntities = n.entitiesOnBusesFromMap(isolatedBuses)
	}

	return isolatedEntities
}

func (n *Network) entitiesOnBusesFromMap(buses map[ID]*Bus) []world.EntityID {
	result := make([]world.EntityID, 0)
	seen := make(map[world.EntityID]bool)

	n.mu.RLock()
	for _, t := range n.terminals {
		if buses[t.Bus().ID] != nil && !seen[t.EntityID] {
			result = append(result, t.EntityID)
			seen[t.EntityID] = true
		}
	}
	n.mu.RUnlock()

	return result
}

// Island represents a connected subgraph of the network.
type Island struct {
	ID    ID
	Buses []*Bus
}

// Islands returns all islands in the network.
func (n *Network) Islands() []*Island {
	n.mu.Lock()
	defer n.mu.Unlock()

	visited := make(map[ID]bool)
	islands := make([]*Island, 0)

	for _, bus := range n.buses {
		if !visited[bus.ID] {
			islandBuses := n.collectIsland(bus, visited)
			island := &Island{
				ID:    ID(fmt.Sprintf("island-%d", len(islands))),
				Buses: islandBuses,
			}
			islands = append(islands, island)
		}
	}

	return islands
}

func (n *Network) collectIsland(bus *Bus, visited map[ID]bool) []*Bus {
	if visited[bus.ID] {
		return nil
	}
	visited[bus.ID] = true

	// Check if this bus is connected to any other bus through closed switches
	hasConnection := false
	for _, br := range bus.branches {
		if sw := br.SwitchDevice(); sw != nil && !sw.IsOpen() {
			otherBus := br.OtherBus(bus)
			if otherBus != nil && !visited[otherBus.ID] {
				hasConnection = true
			}
		}
	}

	result := []*Bus{bus}

	if hasConnection {
		for _, br := range bus.branches {
			if sw := br.SwitchDevice(); sw != nil && !sw.IsOpen() {
				otherBus := br.OtherBus(bus)
				if otherBus != nil {
					result = append(result, n.collectIsland(otherBus, visited)...)
				}
			}
		}
	}

	return result
}

// IslandFor returns the island that contains this bus.
func (n *Network) IslandFor(bus *Bus) *Island {
	islands := n.Islands()
	for _, island := range islands {
		for _, b := range island.Buses {
			if b.ID == bus.ID {
				return island
			}
		}
	}
	return nil
}

// IsBusEnergized returns true if the bus has voltage.
func (n *Network) IsBusEnergized(bus *Bus) bool {
	// A bus is energized if it's connected to a source through closed paths
	// For now, assume grid buses are always energized
	
	// Check if any upstream path is through a closed breaker
	islands := n.Islands()
	for _, island := range islands {
		for _, b := range island.Buses {
			if b.ID == bus.ID {
				// This bus is in an island - check if it has a source
				for _, connected := range island.Buses {
					for _, br := range connected.branches {
						// Check if branch leads to a source (no switch or closed switch)
						if br.SwitchDevice() == nil || !br.SwitchDevice().IsOpen() {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

// PathBetween finds a path between two buses.
func (n *Network) PathBetween(from, to *Bus) []*Bus {
	visited := make(map[ID]bool)
	path := make([]*Bus, 0)
	if n.findPathDFS(from, to, visited, &path) {
		return path
	}
	return nil // No path found
}

func (n *Network) findPathDFS(current, target *Bus, visited map[ID]bool, path *[]*Bus) bool {
	if visited[current.ID] {
		return false
	}
	visited[current.ID] = true
	*path = append(*path, current)

	if current.ID == target.ID {
		return true
	}

	for _, connected := range n.ConnectedTo(current) {
		if n.findPathDFS(connected, target, visited, path) {
			return true
		}
	}

	*path = (*path)[:len(*path)-1]
	return false
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
