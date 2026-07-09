// Package device provides the device model.
// A device is a deterministic memory image with executable behavior and protocol interfaces.
package device

import (
	"github.com/tamzrod/forge/memory"
)

// DeviceID is a unique identifier for a device.
type DeviceID string

// Device represents a virtual industrial device.
// A device owns its memory, behaviors, and protocols.
type Device struct {
	id        DeviceID
	typeName  string
	mem       *memory.MemoryImage
	behaviors []Behavior
	running   bool
}

// New creates a new device with the given ID, type, and memory regions.
func New(id DeviceID, typeName string, memRegions map[string]uint32) *Device {
	return &Device{
		id:        id,
		typeName:  typeName,
		mem:       memory.New(memRegions),
		behaviors: make([]Behavior, 0),
		running:   false,
	}
}

// ID returns the device identifier.
func (d *Device) ID() DeviceID {
	return d.id
}

// Type returns the device type.
func (d *Device) Type() string {
	return d.typeName
}

// Memory returns the device memory image.
func (d *Device) Memory() *memory.MemoryImage {
	return d.mem
}

// AddBehavior adds a behavior to the device.
func (d *Device) AddBehavior(b Behavior) {
	b.Attach(d)
	d.behaviors = append(d.behaviors, b)
}

// Behaviors returns the device's behaviors.
func (d *Device) Behaviors() []Behavior {
	return d.behaviors
}

// Tick executes one simulation step.
// The device executes its behaviors in order.
func (d *Device) Tick() {
	for _, b := range d.behaviors {
		b.Tick()
	}
}

// Start puts the device in running state.
func (d *Device) Start() {
	d.running = true
}

// Stop puts the device in stopped state.
func (d *Device) Stop() {
	d.running = false
}

// Running returns true if the device is running.
func (d *Device) Running() bool {
	return d.running
}

// State represents the lifecycle state of a device.
type State int

const (
	StateCreated   State = iota
	StateInitialized
	StateRunning
	StateStopped
	StateDestroyed
)
