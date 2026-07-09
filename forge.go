// Package forge is the Industrial Simulation Runtime.
// A generic runtime for virtual industrial devices.
package forge

import (
	"github.com/tamzrod/forge/device"
	"github.com/tamzrod/forge/memory"
	"github.com/tamzrod/forge/runtime"
	"github.com/tamzrod/forge/scheduler"
)

// Re-export commonly used types.
type (
	// Device represents a virtual industrial device.
	Device = device.Device

	// DeviceID is a unique identifier for a device.
	DeviceID = device.DeviceID

	// Behavior is the interface for device behaviors.
	Behavior = device.Behavior

	// MemoryImage is the device memory.
	MemoryImage = memory.MemoryImage

	// Quality represents memory value quality.
	Quality = memory.Quality

	// Config holds runtime configuration.
	Config = runtime.Config

	// Scheduler advances simulation time.
	Scheduler = scheduler.Scheduler

	// SimulationClock tracks elapsed time.
	SimulationClock = scheduler.SimulationClock
)

// Quality constants.
const (
	QualityGood      = memory.QualityGood
	QualityUncertain = memory.QualityUncertain
	QualityBad       = memory.QualityBad
	QualityOffline   = memory.QualityOffline
)

// Memory errors.
var (
	ErrInvalidRegion  = memory.ErrInvalidRegion
	ErrInvalidAddress = memory.ErrInvalidAddress
	ErrInvalidSize    = memory.ErrInvalidSize
)

// NewRuntime creates a new Runtime.
var NewRuntime = runtime.New

// NewRuntimeFromFile creates a Runtime from a configuration file.
var NewRuntimeFromFile = runtime.NewFromFile

// LoadConfig loads configuration from a YAML file.
var LoadConfig = runtime.LoadConfig

// DefaultConfig returns the default configuration.
var DefaultConfig = runtime.DefaultConfig

// New creates a new device with the given ID, type, and memory regions.
var NewDevice = device.New

// NewMemory creates a new MemoryImage.
var NewMemory = memory.New
