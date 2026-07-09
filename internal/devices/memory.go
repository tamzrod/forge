package devices

import (
	"fmt"
	"sync"
)

// Memory represents a device's operational memory.
//
// Operational memory stores the values that the device would report to
// the outside world through protocols. This is distinct from:
// - Simulation Model state (physics)
// - Device internal calculation state
//
// The memory is organized as a collection of named values.
type Memory struct {
	mu    sync.RWMutex
	values map[string]float64
}

// NewMemory creates a new operational memory.
func NewMemory() *Memory {
	return &Memory{
		values: make(map[string]float64),
	}
}

// Set writes a value to memory.
func (m *Memory) Set(name string, value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.values[name] = value
}

// Get reads a value from memory.
func (m *Memory) Get(name string) (float64, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.values[name]
	return v, ok
}

// GetOrDefault reads a value or returns a default.
func (m *Memory) GetOrDefault(name string, defaultVal float64) float64 {
	if v, ok := m.Get(name); ok {
		return v
	}
	return defaultVal
}

// Values returns a copy of all values.
func (m *Memory) Values() map[string]float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make(map[string]float64, len(m.values))
	for k, v := range m.values {
		result[k] = v
	}
	return result
}

// Contains checks if a value exists.
func (m *Memory) Contains(name string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.values[name]
	return ok
}

// Reset clears all values.
func (m *Memory) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.values = make(map[string]float64)
}

// MemoryRegion represents a named region of memory.
type MemoryRegion struct {
	Name    string
	Base    uint16
	Size    uint16
}

// MemoryMap defines the memory layout for a device.
type MemoryMap struct {
	regions map[string]MemoryRegion
}

// NewMemoryMap creates a new memory map.
func NewMemoryMap() *MemoryMap {
	return &MemoryMap{
		regions: make(map[string]MemoryRegion),
	}
}

// AddRegion adds a memory region.
func (m *MemoryMap) AddRegion(name string, base uint16, size uint16) {
	m.regions[name] = MemoryRegion{
		Name: name,
		Base: base,
		Size: size,
	}
}

// Region returns a region by name.
func (m *MemoryMap) Region(name string) (MemoryRegion, bool) {
	r, ok := m.regions[name]
	return r, ok
}

// ValidateAddress checks if an address is valid for a region.
func (m *MemoryMap) ValidateAddress(region string, offset uint16) error {
	r, ok := m.regions[region]
	if !ok {
		return fmt.Errorf("unknown region: %s", region)
	}
	if offset >= r.Size {
		return fmt.Errorf("offset %d out of bounds for region %s (size %d)", offset, region, r.Size)
	}
	return nil
}

// RegisterMap creates a standard register-based memory map.
func StandardRegisterMap(baseAddress uint16) *MemoryMap {
	m := NewMemoryMap()
	m.AddRegion("registers", baseAddress, 100)
	return m
}
