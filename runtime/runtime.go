// Package runtime provides the simulation runtime.
// The runtime hosts simulation models and devices.
package runtime

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/tamzrod/forge/device"
	"github.com/tamzrod/forge/models"
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
// It hosts simulation models and devices, and provides scheduling infrastructure.
type Runtime struct {
	config    Config
	sched     *scheduler.Scheduler
	models    map[models.ModelID]models.Model
	devices   map[device.DeviceID]*device.Device
}

// New creates a new Runtime.
func New(cfg Config) *Runtime {
	return &Runtime{
		config:  cfg,
		sched:   scheduler.New(cfg.TickInterval),
		models:  make(map[models.ModelID]models.Model),
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
	d.SetModelGetter(func(modelID models.ModelID) models.Model {
		return r.Model(modelID)
	})
	r.devices[id] = d
	r.sched.AddDevice(d)
	return d
}

// CreateModel creates a new simulation model and adds it to the runtime.
func (r *Runtime) CreateModel(m models.Model) {
	r.models[m.ID()] = m
	r.sched.AddModel(m)
}

// CreateGridModel creates a new Grid model.
func (r *Runtime) CreateGridModel(id models.ModelID) *models.GridModel {
	m := models.NewGridModel(id)
	r.CreateModel(m)
	return m
}

// CreateSunModel creates a new Sun model.
func (r *Runtime) CreateSunModel(id models.ModelID) *models.SunModel {
	m := models.NewSunModel(id)
	r.CreateModel(m)
	return m
}

// CreateWindModel creates a new Wind model.
func (r *Runtime) CreateWindModel(id models.ModelID) *models.WindModel {
	m := models.NewWindModel(id)
	r.CreateModel(m)
	return m
}

// CreateWeatherModel creates a new Weather model.
func (r *Runtime) CreateWeatherModel(id models.ModelID) *models.WeatherModel {
	m := models.NewWeatherModel(id)
	r.CreateModel(m)
	return m
}

// CreateBusModel creates a new Bus model.
func (r *Runtime) CreateBusModel(id models.ModelID, nominalV float32) *models.BusModel {
	m := models.NewBusModel(id, nominalV)
	r.CreateModel(m)
	return m
}

// CreateTransformerModel creates a new Transformer model.
func (r *Runtime) CreateTransformerModel(id models.ModelID, from, to models.ModelID) *models.TransformerModel {
	m := models.NewTransformerModel(id, from, to)
	r.CreateModel(m)
	return m
}

// CreateLoadModel creates a new Load model.
func (r *Runtime) CreateLoadModel(id models.ModelID, bus models.ModelID, baseLoad float32) *models.LoadModel {
	m := models.NewLoadModel(id, bus, baseLoad)
	r.CreateModel(m)
	return m
}

// CreateBreakerModel creates a new Breaker model.
func (r *Runtime) CreateBreakerModel(id models.ModelID, bus1, bus2 models.ModelID) *models.BreakerModel {
	m := models.NewBreakerModel(id, bus1, bus2)
	r.CreateModel(m)
	return m
}

// CreatePVArrayModel creates a new PV array model.
func (r *Runtime) CreatePVArrayModel(id models.ModelID, bus models.ModelID, ratedPower float32) *models.PVArrayModel {
	m := models.NewPVArrayModel(id, bus, ratedPower)
	r.CreateModel(m)
	return m
}

// CreateReservoirModel creates a new Reservoir model.
func (r *Runtime) CreateReservoirModel(id models.ModelID, area float32) *models.ReservoirModel {
	m := models.NewReservoirModel(id, area)
	r.CreateModel(m)
	return m
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

// Model returns a simulation model by ID.
func (r *Runtime) Model(id models.ModelID) models.Model {
	return r.models[id]
}

// Models returns all simulation models.
func (r *Runtime) Models() []models.Model {
	result := make([]models.Model, 0, len(r.models))
	for _, m := range r.models {
		result = append(result, m)
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
