// Package devices provides the Virtual Device framework.
//
// Virtual Devices observe Simulation Models and maintain their own
// operational memory. They are the bridge between the simulated world
// and the external applications that will eventually consume their data.
//
// Architecture:
//
//   Simulation Models (physics)
//           ↓
//   Virtual Devices (observation)
//           ↓
//   Operational Memory (device state)
//           ↓
//   Protocols (future)
//           ↓
//   MMA2 (future)
//
// Virtual Devices:
// - Observe Simulation Models through Simulation Context
// - Own operational memory
// - Know nothing about protocols
// - Know nothing about MMA2
package devices

import (
	"fmt"
	"sync"
)

// DeviceID uniquely identifies a device within the simulation.
type DeviceID string

// DeviceType categorizes devices (e.g., "weather_station", "revenue_meter").
type DeviceType string

// State represents the current operational state of a device.
type State int

const (
	StateCreated State = iota
	StateInitialized
	StateRunning
	StatePaused
	StateStopped
	StateFaulted
)

func (s State) String() string {
	switch s {
	case StateCreated:
		return "Created"
	case StateInitialized:
		return "Initialized"
	case StateRunning:
		return "Running"
	case StatePaused:
		return "Paused"
	case StateStopped:
		return "Stopped"
	case StateFaulted:
		return "Faulted"
	default:
		return "Unknown"
	}
}

// Device represents a virtual industrial device.
//
// Each device:
// - Has a unique identity (ID, Name, Type)
// - Owns operational memory
// - Observes Simulation Models through a context
// - Implements lifecycle methods (Initialize, Tick, Shutdown)
type Device interface {
	// Identity returns the device's unique identifier.
	ID() DeviceID

	// Type returns the device type.
	Type() DeviceType

	// Name returns the device's human-readable name.
	Name() string

	// State returns the current operational state.
	State() State

	// Initialize prepares the device for operation.
	// Called once before the first Tick.
	Initialize() error

	// Tick advances the device by one simulation step.
	// The device observes models and updates its operational memory.
	Tick()

	// Shutdown stops the device and releases resources.
	Shutdown() error
}

// BaseDevice provides common functionality for all devices.
type BaseDevice struct {
	mu   sync.RWMutex
	id   DeviceID
	typ  DeviceType
	name string
	state State
}

// NewBaseDevice creates a new base device with the given identity.
func NewBaseDevice(id DeviceID, typ DeviceType, name string) *BaseDevice {
	return &BaseDevice{
		id:    id,
		typ:   typ,
		name:  name,
		state: StateCreated,
	}
}

func (d *BaseDevice) ID() DeviceID { return d.id }
func (d *BaseDevice) Type() DeviceType { return d.typ }
func (d *BaseDevice) Name() string { return d.name }

func (d *BaseDevice) State() State {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.state
}

func (d *BaseDevice) setState(state State) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.state = state
}

// ValidateID checks if a device ID is valid.
func ValidateID(id DeviceID) error {
	if id == "" {
		return fmt.Errorf("device ID cannot be empty")
	}
	return nil
}

// ValidateType checks if a device type is valid.
func ValidateType(typ DeviceType) error {
	if typ == "" {
		return fmt.Errorf("device type cannot be empty")
	}
	return nil
}
