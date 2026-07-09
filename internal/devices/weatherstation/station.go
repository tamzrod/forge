package weatherstation

import (
	"fmt"
	"sync"

	"github.com/tamzrod/forge/internal/devices"
)

// Type is the device type identifier.
const Type = devices.DeviceType("weather_station")

// Station represents a virtual Weather Station device.
//
// The Weather Station observes the Weather Model and maintains its own
// operational memory representing what it would report externally.
//
// Architecture:
//   Weather Model → Weather Station → Operational Memory → Protocols (future)
type Station struct {
	*devices.BaseDevice

	config  Config
	ctx     *devices.Context
	memory  *devices.Memory

	mu        sync.RWMutex
	tickCount uint64
}

// NewStation creates a new Weather Station device.
func NewStation(cfg Config, ctx *devices.Context) (*Station, error) {
	if err := devices.ValidateID(cfg.ID); err != nil {
		return nil, err
	}

	return &Station{
		BaseDevice: devices.NewBaseDevice(cfg.ID, Type, cfg.Name),
		config:     cfg,
		ctx:        ctx,
		memory:     devices.NewMemory(),
	}, nil
}

// Initialize prepares the Weather Station for operation.
func (s *Station) Initialize() error {
	if s.State() != devices.StateCreated {
		return fmt.Errorf("cannot initialize device in state %s", s.State())
	}

	// Initialize memory with default values
	s.memory.Set(MemoryTemperature, 20.0)
	s.memory.Set(MemoryHumidity, 50.0)
	s.memory.Set(MemoryPressure, 1013.25)
	s.memory.Set(MemoryCloudCover, 0.0)
	s.memory.Set(MemoryWindSpeed, 0.0)
	s.memory.Set(MemoryWindDirection, 0.0)
	s.memory.Set(MemoryRainStatus, 0.0)
	s.memory.Set(MemoryStatus, 1.0) // 1 = OK
	s.memory.Set(MemoryTickCount, 0.0)

	s.setState(devices.StateInitialized)
	return nil
}

// Tick updates the Weather Station's operational memory by observing
// the Weather Model.
func (s *Station) Tick() {
	if s.State() != devices.StateInitialized && s.State() != devices.StateRunning {
		return
	}

	if s.State() == devices.StateInitialized {
		s.setState(devices.StateRunning)
	}

	// Observe weather model
	weather := s.ctx.ReadWeather()

	s.mu.Lock()

	// Copy measurements to operational memory
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

// Shutdown stops the Weather Station.
func (s *Station) Shutdown() error {
	s.setState(devices.StateStopped)
	s.memory.Reset()
	return nil
}

// Memory returns the device's operational memory.
func (s *Station) Memory() *devices.Memory {
	return s.memory
}

// Temperature returns the current temperature measurement.
func (s *Station) Temperature() float64 {
	return s.memory.GetOrDefault(MemoryTemperature, 0)
}

// Humidity returns the current humidity measurement.
func (s *Station) Humidity() float64 {
	return s.memory.GetOrDefault(MemoryHumidity, 0)
}

// Pressure returns the current pressure measurement.
func (s *Station) Pressure() float64 {
	return s.memory.GetOrDefault(MemoryPressure, 0)
}

// CloudCover returns the current cloud cover measurement.
func (s *Station) CloudCover() float64 {
	return s.memory.GetOrDefault(MemoryCloudCover, 0)
}

// WindSpeed returns the current wind speed measurement.
func (s *Station) WindSpeed() float64 {
	return s.memory.GetOrDefault(MemoryWindSpeed, 0)
}

// WindDirection returns the current wind direction measurement.
func (s *Station) WindDirection() float64 {
	return s.memory.GetOrDefault(MemoryWindDirection, 0)
}

// TickCount returns the number of ticks since initialization.
func (s *Station) TickCount() uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tickCount
}

// State returns a snapshot of all measurements.
func (s *Station) State() WeatherStationState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return WeatherStationState{
		ID:             s.ID(),
		Name:           s.Name(),
		Type:           s.Type(),
		DeviceState:    s.BaseDevice.State(),
		Temperature:    s.memory.GetOrDefault(MemoryTemperature, 0),
		Humidity:       s.memory.GetOrDefault(MemoryHumidity, 0),
		Pressure:       s.memory.GetOrDefault(MemoryPressure, 0),
		CloudCover:     s.memory.GetOrDefault(MemoryCloudCover, 0),
		WindSpeed:      s.memory.GetOrDefault(MemoryWindSpeed, 0),
		WindDirection:  s.memory.GetOrDefault(MemoryWindDirection, 0),
		RainStatus:      s.memory.GetOrDefault(MemoryRainStatus, 0) == 1.0,
		TickCount:      s.tickCount,
	}
}

// WeatherStationState represents the current state of the Weather Station.
type WeatherStationState struct {
	ID             devices.DeviceID     `json:"id"`
	Name           string               `json:"name"`
	Type           devices.DeviceType   `json:"type"`
	DeviceState    devices.State        `json:"device_state"`
	Temperature    float64              `json:"temperature"`
	Humidity       float64             `json:"humidity"`
	Pressure       float64             `json:"pressure"`
	CloudCover     float64             `json:"cloud_cover"`
	WindSpeed      float64             `json:"wind_speed"`
	WindDirection  float64             `json:"wind_direction"`
	RainStatus     bool                `json:"rain_status"`
	TickCount      uint64              `json:"tick_count"`
}
