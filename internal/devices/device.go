// Package devices provides the Virtual Firmware framework.
//
// Each Virtual Device represents the firmware running inside an industrial device.
// Virtual Firmware:
// - Samples Simulation Models (the external world)
// - Updates its Device Memory with observations
// - Exposes Device Memory through Communication Interfaces
//
// Architecture:
//
//   Simulation Models (external physical world)
//           ↓
//   Virtual Firmware (samples models, owns device memory)
//           ↓
//   Device Memory (firmware-owned internal state)
//           ↓
//   Communication Interfaces (serialize memory for protocols)
//           ↓
//   External Systems (MMA2, SCADA, etc.)
//
// The Virtual Firmware model mirrors real embedded systems:
// - A weather station firmware samples environmental sensors
// - Updates its internal memory with readings
// - Exposes memory through Modbus, Raw Ingest, etc.
// - Never exposes raw sensor data directly
package devices

import (
	"fmt"
	"sync"
)

// DeviceID uniquely identifies a device within the simulation.
type DeviceID string

// DeviceType categorizes devices (e.g., "weather_station", "revenue_meter").
type DeviceType string

// State represents the current operational state of the firmware.
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

// Device represents virtual firmware running inside an industrial device.
//
// Each virtual firmware:
// - Has a unique identity (ID, Name, Type)
// - Owns Device Memory
// - Samples Simulation Models through a context
// - Exposes memory through Communication Interfaces
// - Implements lifecycle methods (Initialize, Tick, Shutdown)
type Device interface {
	// Identity returns the device's unique identifier.
	ID() DeviceID

	// Type returns the device type.
	Type() DeviceType

	// Name returns the device's human-readable name.
	Name() string

	// State returns the current firmware state.
	State() State

	// Initialize prepares the firmware for operation.
	// Called once before the first Tick.
	Initialize() error

	// Tick advances the firmware by one simulation step.
	// The firmware samples models and updates Device Memory.
	Tick()

	// Shutdown stops the firmware and releases resources.
	Shutdown() error
}

// BaseDevice provides common functionality for all virtual firmware.
type BaseDevice struct {
	mu    sync.RWMutex
	id    DeviceID
	typ   DeviceType
	name  string
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

func (d *BaseDevice) ID() DeviceID   { return d.id }
func (d *BaseDevice) Type() DeviceType { return d.typ }
func (d *BaseDevice) Name() string { return d.name }

func (d *BaseDevice) State() State {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.state
}

func (d *BaseDevice) SetState(state State) {
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
