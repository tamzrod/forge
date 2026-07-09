// Package runtime provides the simulation runtime.
// The runtime hosts devices. That's all.
package runtime

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/tamzrod/forge/device"
	"github.com/tamzrod/forge/scheduler"
	"gopkg.in/yaml.v3"
)

// Config holds runtime configuration.
type Config struct {
	TickInterval  time.Duration `yaml:"tick_interval"`
	MaxDevices   int           `yaml:"max_devices"`
}

// DefaultConfig returns the default configuration.
func DefaultConfig() Config {
	return Config{
		TickInterval: 100 * time.Millisecond,
		MaxDevices:   1000,
	}
}

// LoadConfig loads configuration from a YAML file.
func LoadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("failed to parse config: %w", err)
	}

	// Apply defaults
	if cfg.TickInterval == 0 {
		cfg.TickInterval = DefaultConfig().TickInterval
	}
	if cfg.MaxDevices == 0 {
		cfg.MaxDevices = DefaultConfig().MaxDevices
	}

	return cfg, nil
}

// Runtime is the simulation runtime.
// It hosts devices and provides scheduling infrastructure.
type Runtime struct {
	config    Config
	sched     *scheduler.Scheduler
	devices   map[device.DeviceID]*device.Device
}

// New creates a new Runtime.
func New(cfg Config) *Runtime {
	return &Runtime{
		config:  cfg,
		sched:   scheduler.New(cfg.TickInterval),
		devices: make(map[device.DeviceID]*device.Device),
	}
}

// NewFromFile creates a new Runtime from a configuration file.
func NewFromFile(configPath string) (*Runtime, error) {
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}
	return New(cfg), nil
}

// CreateDevice creates a new device and adds it to the runtime.
func (r *Runtime) CreateDevice(id device.DeviceID, typeName string, memRegions map[string]uint32) *device.Device {
	d := device.New(id, typeName, memRegions)
	r.devices[id] = d
	r.sched.AddDevice(d)
	return d
}

// Device returns a device by ID.
func (r *Runtime) Device(id device.DeviceID) *device.Device {
	return r.devices[id]
}

// Devices returns all devices.
func (r *Runtime) Devices() []*device.Device {
	result := make([]*device.Device, 0, len(r.devices))
	for _, d := range r.devices {
		result = append(result, d)
	}
	return result
}

// Start starts all devices.
func (r *Runtime) Start() {
	for _, d := range r.devices {
		d.Start()
	}
}

// Stop stops all devices.
func (r *Runtime) Stop() {
	for _, d := range r.devices {
		d.Stop()
	}
}

// Run starts the runtime.
// It runs until the context is cancelled.
func (r *Runtime) Run(ctx context.Context) error {
	r.Start()
	return r.sched.Run(ctx)
}

// Shutdown stops the runtime and cleans up.
func (r *Runtime) Shutdown() error {
	r.Stop()
	r.sched.Stop()
	return nil
}

// Config returns the runtime configuration.
func (r *Runtime) Config() Config {
	return r.config
}

// Scheduler returns the runtime scheduler.
func (r *Runtime) Scheduler() *scheduler.Scheduler {
	return r.sched
}
