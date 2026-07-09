// Package memory provides the memory model for virtual devices.
// Memory is the source of truth for device state.
package memory

import (
	"encoding/binary"
	"errors"
	"fmt"
	"sync"
)

// Quality represents the quality of a memory value.
type Quality uint8

const (
	QualityGood      Quality = 0x00
	QualityUncertain Quality = 0x40
	QualityBad       Quality = 0x80
	QualityOffline   Quality = 0x84
)

// ErrInvalidRegion is returned when a region does not exist.
var ErrInvalidRegion = errors.New("invalid region")

// ErrInvalidAddress is returned when an address is out of bounds.
var ErrInvalidAddress = errors.New("invalid address")

// ErrInvalidSize is returned when a read/write size is invalid.
var ErrInvalidSize = errors.New("invalid size")

// Region represents a named memory region.
type Region struct {
	Name   string
	Size   uint32
	Values []byte
	Quality []Quality
}

// MemoryImage represents the memory of a device.
// It is the single source of truth for device state.
type MemoryImage struct {
	mu      sync.RWMutex
	regions map[string]*Region
}

// New creates a new MemoryImage with the given region definitions.
func New(regionDefs map[string]uint32) *MemoryImage {
	m := &MemoryImage{
		regions: make(map[string]*Region),
	}
	for name, size := range regionDefs {
		m.regions[name] = &Region{
			Name:    name,
			Size:    size,
			Values:  make([]byte, size),
			Quality: make([]Quality, size),
		}
		// Initialize all quality to Good
		for i := range m.regions[name].Quality {
			m.regions[name].Quality[i] = QualityGood
		}
	}
	return m
}

// Read reads bytes from a memory region.
func (m *MemoryImage) Read(region string, address uint32, size uint32) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	r, ok := m.regions[region]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrInvalidRegion, region)
	}
	if address+size > r.Size {
		return nil, fmt.Errorf("%w: address %d + size %d > region size %d",
			ErrInvalidAddress, address, size, r.Size)
	}

	result := make([]byte, size)
	copy(result, r.Values[address:address+size])
	return result, nil
}

// Write writes bytes to a memory region.
func (m *MemoryImage) Write(region string, address uint32, data []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	r, ok := m.regions[region]
	if !ok {
		return fmt.Errorf("%w: %s", ErrInvalidRegion, region)
	}
	if address+uint32(len(data)) > r.Size {
		return fmt.Errorf("%w: address %d + size %d > region size %d",
			ErrInvalidAddress, address, len(data), r.Size)
	}

	copy(r.Values[address:], data)
	return nil
}

// ReadUint16 reads a 16-bit unsigned integer.
func (m *MemoryImage) ReadUint16(region string, address uint32) (uint16, error) {
	data, err := m.Read(region, address, 2)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(data), nil
}

// WriteUint16 writes a 16-bit unsigned integer.
func (m *MemoryImage) WriteUint16(region string, address uint32, value uint16) error {
	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data, value)
	return m.Write(region, address, data)
}

// ReadFloat32 reads a 32-bit float.
func (m *MemoryImage) ReadFloat32(region string, address uint32) (float32, error) {
	data, err := m.Read(region, address, 4)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(data), nil
}

// WriteFloat32 writes a 32-bit float.
func (m *MemoryImage) WriteFloat32(region string, address uint32, value float32) error {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, value)
	return m.Write(region, address, data)
}

// Quality returns the quality flag for a memory location.
func (m *MemoryImage) Quality(region string, address uint32) (Quality, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	r, ok := m.regions[region]
	if !ok {
		return QualityBad, fmt.Errorf("%w: %s", ErrInvalidRegion, region)
	}
	if address >= r.Size {
		return QualityBad, fmt.Errorf("%w: %d >= %d", ErrInvalidAddress, address, r.Size)
	}
	return r.Quality[address], nil
}

// SetQuality sets the quality flag for a memory location.
func (m *MemoryImage) SetQuality(region string, address uint32, quality Quality) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	r, ok := m.regions[region]
	if !ok {
		return fmt.Errorf("%w: %s", ErrInvalidRegion, region)
	}
	if address >= r.Size {
		return fmt.Errorf("%w: %d >= %d", ErrInvalidAddress, address, r.Size)
	}
	r.Quality[address] = quality
	return nil
}

// RegionSize returns the size of a region.
func (m *MemoryImage) RegionSize(region string) (uint32, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	r, ok := m.regions[region]
	if !ok {
		return 0, fmt.Errorf("%w: %s", ErrInvalidRegion, region)
	}
	return r.Size, nil
}

// Regions returns the names of all regions.
func (m *MemoryImage) Regions() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.regions))
	for name := range m.regions {
		names = append(names, name)
	}
	return names
}
