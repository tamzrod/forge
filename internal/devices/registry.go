package devices

import (
	"fmt"
	"sync"
)

// Registry manages all virtual devices in the simulation.
type Registry struct {
	mu      sync.RWMutex
	devices map[DeviceID]Device
}

// NewRegistry creates a new device registry.
func NewRegistry() *Registry {
	return &Registry{
		devices: make(map[DeviceID]Device),
	}
}

// Register adds a device to the registry.
func (r *Registry) Register(device Device) error {
	if err := ValidateID(device.ID()); err != nil {
		return err
	}
	if err := ValidateType(device.Type()); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.devices[device.ID()]; exists {
		return fmt.Errorf("device already registered: %s", device.ID())
	}

	r.devices[device.ID()] = device
	return nil
}

// Unregister removes a device from the registry.
func (r *Registry) Unregister(id DeviceID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.devices[id]; !exists {
		return fmt.Errorf("device not found: %s", id)
	}

	delete(r.devices, id)
	return nil
}

// Device returns a device by ID.
func (r *Registry) Device(id DeviceID) (Device, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.devices[id]
	return d, ok
}

// Devices returns all registered devices.
func (r *Registry) Devices() []Device {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]Device, 0, len(r.devices))
	for _, d := range r.devices {
		result = append(result, d)
	}
	return result
}

// DevicesByType returns all devices of a specific type.
func (r *Registry) DevicesByType(typ DeviceType) []Device {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]Device, 0)
	for _, d := range r.devices {
		if d.Type() == typ {
			result = append(result, d)
		}
	}
	return result
}

// Initialize initializes all registered devices.
func (r *Registry) Initialize() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, d := range r.devices {
		if err := d.Initialize(); err != nil {
			return fmt.Errorf("failed to initialize device %s: %w", d.ID(), err)
		}
	}
	return nil
}

// Tick ticks all registered devices.
func (r *Registry) Tick() {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, d := range r.devices {
		d.Tick()
	}
}

// Shutdown shuts down all registered devices.
func (r *Registry) Shutdown() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var lastErr error
	for _, d := range r.devices {
		if err := d.Shutdown(); err != nil {
			lastErr = fmt.Errorf("failed to shutdown device %s: %w", d.ID(), err)
		}
	}
	return lastErr
}

// Count returns the number of registered devices.
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.devices)
}

// Clear removes all devices from the registry.
func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.devices = make(map[DeviceID]Device)
}
