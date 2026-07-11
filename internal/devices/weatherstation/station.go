package weatherstation

import (
	"fmt"
	"sync"
	"time"

	"github.com/tamzrod/forge/internal/devices"
	"github.com/tamzrod/forge/internal/publishers/rawingest"
)

// Type is the device type identifier.
const Type = devices.DeviceType("weather_station")

// Station represents the virtual firmware for a Weather Station device.
//
// The Weather Station firmware samples the Weather Model (external world)
// and updates its Device Memory with observations. This mirrors how
// real embedded weather station firmware samples sensors.
//
// Architecture:
//   Weather Model (external world)
//           ↓
//   Weather Station Firmware (samples, updates memory)
//           ↓
//   Device Memory (firmware-owned)
//           ↓
//   Raw Ingest Interface (serializes memory)
//           ↓
//   MMA2
type Station struct {
	*devices.BaseDevice

	config    Config
	ctx       *devices.Context
	memory    *devices.DeviceMemory
	publisher *rawingest.Interface

	mu          sync.RWMutex
	tickCount   uint64
	lastPublish time.Time
}

// NewStation creates a new Weather Station device.
func NewStation(cfg Config, ctx *devices.Context) (*Station, error) {
	if err := devices.ValidateID(cfg.ID); err != nil {
		return nil, err
	}

	station := &Station{
		BaseDevice: devices.NewBaseDevice(cfg.ID, Type, cfg.Name),
		config:     cfg,
		ctx:        ctx,
		memory:     devices.NewDeviceMemory(),
	}

	// Create Raw Ingest interface if enabled
	if cfg.Publishing.Enabled {
		ifaceCfg := rawingest.Config{
			Host:     cfg.Publishing.Host,
			Port:     cfg.Publishing.Port,
			UnitID:   cfg.Publishing.UnitID,
			Enabled:  true,
			Interval: cfg.Publishing.Interval,
		}
	 iface, err := rawingest.NewInterface(ifaceCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create interface: %w", err)
		}
		station.publisher = iface
	}

	return station, nil
}

// Initialize prepares the Weather Station firmware for operation.
func (s *Station) Initialize() error {
	if s.BaseDevice.State() != devices.StateCreated {
		return fmt.Errorf("cannot initialize device in state %s", s.BaseDevice.State())
	}

	// Initialize device memory with default values
	s.memory.Set(MemoryTemperature, 20.0)
	s.memory.Set(MemoryHumidity, 50.0)
	s.memory.Set(MemoryPressure, 1013.25)
	s.memory.Set(MemoryCloudCover, 0.0)
	s.memory.Set(MemoryWindSpeed, 0.0)
	s.memory.Set(MemoryWindDirection, 0.0)
	s.memory.Set(MemoryRainStatus, 0.0)
	s.memory.Set(MemoryStatus, 1.0) // 1 = OK
	s.memory.Set(MemoryTickCount, 0.0)

	// Start interface if configured
	if s.publisher != nil {
		if err := s.publisher.Start(); err != nil {
			return fmt.Errorf("failed to start interface: %w", err)
		}
	}

	s.BaseDevice.SetState(devices.StateInitialized)
	return nil
}

// Tick advances the Weather Station firmware by one simulation step.
// The firmware samples the Weather Model and updates Device Memory.
func (s *Station) Tick() {
	if s.BaseDevice.State() != devices.StateInitialized && s.BaseDevice.State() != devices.StateRunning {
		return
	}

	if s.BaseDevice.State() == devices.StateInitialized {
		s.BaseDevice.SetState(devices.StateRunning)
	}

	// Sample weather model (external world)
	weather := s.ctx.ReadWeather()

	s.mu.Lock()

	// Update device memory with observations
	s.memory.Set(MemoryTemperature, s.convertTemperature(weather.Temperature))
	s.memory.Set(MemoryHumidity, weather.Humidity)
	s.memory.Set(MemoryPressure, weather.Pressure)
	s.memory.Set(MemoryCloudCover, weather.CloudCover*100) // Store as percentage
	s.memory.Set(MemoryWindSpeed, weather.WindSpeed)
	s.memory.Set(MemoryWindDirection, weather.WindDirection)

	if weather.IsRaining {
		s.memory.Set(MemoryRainStatus, 1.0)
	} else {
		s.memory.Set(MemoryRainStatus, 0.0)
	}

	s.tickCount++
	s.memory.Set(MemoryTickCount, float64(s.tickCount))

	s.mu.Unlock()

	// Push device memory through interface
	s.publish()
}

// publish sends current device memory through the communication interface.
func (s *Station) publish() {
	if s.publisher == nil {
		return
	}

	values := s.memory.Values()
	if err := s.publisher.Publish(values); err != nil {
		// Publishing failure is logged but doesn't stop the firmware
		return
	}

	s.mu.Lock()
	s.lastPublish = time.Now()
	s.mu.Unlock()
}

// convertTemperature converts temperature to the configured unit.
func (s *Station) convertTemperature(celsius float64) float64 {
	switch s.config.Units {
	case Fahrenheit:
		return celsius*9/5 + 32
	default:
		return celsius
	}
}

// Shutdown stops the Weather Station firmware.
func (s *Station) Shutdown() error {
	s.BaseDevice.SetState(devices.StateStopped)

	if s.publisher != nil {
		s.publisher.Stop()
	}

	s.memory.Reset()
	return nil
}

// Memory returns the device memory owned by this firmware.
func (s *Station) Memory() *devices.DeviceMemory {
	return s.memory
}

// Interface returns the Raw Ingest communication interface.
func (s *Station) Interface() *rawingest.Interface {
	return s.publisher
}

// Temperature returns the current temperature from device memory.
func (s *Station) Temperature() float64 {
	return s.memory.GetOrDefault(MemoryTemperature, 0)
}

// Humidity returns the current humidity from device memory.
func (s *Station) Humidity() float64 {
	return s.memory.GetOrDefault(MemoryHumidity, 0)
}

// Pressure returns the current pressure from device memory.
func (s *Station) Pressure() float64 {
	return s.memory.GetOrDefault(MemoryPressure, 0)
}

// CloudCover returns the current cloud cover from device memory.
func (s *Station) CloudCover() float64 {
	return s.memory.GetOrDefault(MemoryCloudCover, 0)
}

// WindSpeed returns the current wind speed from device memory.
func (s *Station) WindSpeed() float64 {
	return s.memory.GetOrDefault(MemoryWindSpeed, 0)
}

// WindDirection returns the current wind direction from device memory.
func (s *Station) WindDirection() float64 {
	return s.memory.GetOrDefault(MemoryWindDirection, 0)
}

// TickCount returns the number of ticks since initialization.
func (s *Station) TickCount() uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tickCount
}

// PublishingState returns the communication interface status.
func (s *Station) PublishingState() PublishingState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state := PublishingState{
		Enabled:     s.publisher != nil,
		LastPublish: s.lastPublish,
	}

	if s.publisher != nil {
		state.Connected = s.publisher.IsConnected()
		state.PacketsSent = s.publisher.Stats().PacketsSent
		state.Errors = s.publisher.Stats().Errors
		state.LastError = s.publisher.Stats().LastError
	}

	return state
}

// Status returns a snapshot of all measurements.
func (s *Station) Status() WeatherStationState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	pubState := s.PublishingState()

	return WeatherStationState{
		ID:               s.ID(),
		Name:             s.Name(),
		Type:             s.Type(),
		DeviceState:      s.BaseDevice.State(),
		Temperature:      s.memory.GetOrDefault(MemoryTemperature, 0),
		Humidity:         s.memory.GetOrDefault(MemoryHumidity, 0),
		Pressure:         s.memory.GetOrDefault(MemoryPressure, 0),
		CloudCover:       s.memory.GetOrDefault(MemoryCloudCover, 0),
		WindSpeed:        s.memory.GetOrDefault(MemoryWindSpeed, 0),
		WindDirection:    s.memory.GetOrDefault(MemoryWindDirection, 0),
		RainStatus:       s.memory.GetOrDefault(MemoryRainStatus, 0) == 1.0,
		TickCount:        s.tickCount,
		PublishingState:  pubState,
	}
}

// PublishingState represents the Raw Ingest interface status.
type PublishingState struct {
	Enabled     bool      `json:"enabled"`
	Connected   bool      `json:"connected"`
	PacketsSent uint64    `json:"packets_sent"`
	Errors      uint64    `json:"errors"`
	LastPublish time.Time `json:"last_publish"`
	LastError   string    `json:"last_error,omitempty"`
}

// WeatherStationState represents the current state of the Weather Station.
type WeatherStationState struct {
	ID              devices.DeviceID   `json:"id"`
	Name            string             `json:"name"`
	Type            devices.DeviceType `json:"type"`
	DeviceState     devices.State      `json:"device_state"`
	Temperature     float64           `json:"temperature"`
	Humidity        float64           `json:"humidity"`
	Pressure        float64           `json:"pressure"`
	CloudCover      float64           `json:"cloud_cover"`
	WindSpeed       float64           `json:"wind_speed"`
	WindDirection   float64           `json:"wind_direction"`
	RainStatus      bool              `json:"rain_status"`
	TickCount       uint64            `json:"tick_count"`
	PublishingState PublishingState   `json:"publishing_state"`
}
