// Package forge is the Industrial Simulation Runtime.
// A generic runtime for virtual industrial devices.
package forge

import (
	"github.com/tamzrod/forge/device"
	"github.com/tamzrod/forge/memory"
	"github.com/tamzrod/forge/runtime"
	"github.com/tamzrod/forge/scheduler"
	"github.com/tamzrod/forge/simulation"
	"github.com/tamzrod/forge/topology"
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

	// SimClock is the simulation clock.
	SimClock = simulation.SimClock

	// Clock is the simulation clock interface.
	Clock = simulation.Clock

	// Mode is the simulation mode.
	Mode = simulation.Mode

	// Network is the electrical network topology.
	Network = topology.Network

	// Bus is an electrical node.
	Bus = topology.Bus

	// Branch is a connection between buses.
	Branch = topology.Branch

	// Terminal is a connection point on an entity.
	Terminal = topology.Terminal

	// Switch is a switching device in a branch.
	Switch = topology.Switch
)

// Quality constants.
const (
	QualityGood      = memory.QualityGood
	QualityUncertain = memory.QualityUncertain
	QualityBad       = memory.QualityBad
	QualityOffline   = memory.QualityOffline
)

// Simulation modes.
const (
	ModeRealtime  = simulation.ModeRealtime
	ModeSimulated = simulation.ModeSimulated
	ModeManual    = simulation.ModeManual
	ModeReplay    = simulation.ModeReplay
)

// Simulation speeds.
const (
	SpeedRealtime   = simulation.SpeedRealtime
	SpeedSlow      = simulation.SpeedSlow
	SpeedFast      = simulation.SpeedFast
	SpeedVeryFast  = simulation.SpeedVeryFast
	SpeedUltraFast = simulation.SpeedUltraFast
	SpeedExtreme   = simulation.SpeedExtreme
	SpeedLudicrous = simulation.SpeedLudicrous
	SpeedPlaid     = simulation.SpeedPlaid
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
