package weatherstation

import (
	"time"

	"github.com/tamzrod/forge/internal/devices"
)

// Config holds Weather Station configuration.
type Config struct {
	// Device identity
	ID   devices.DeviceID
	Name string

	// Measurement configuration
	Units TemperatureUnit

	// Memory configuration
	RegisterBase uint16

	// Publishing configuration
	Publishing PublishingConfig
}

// PublishingConfig holds Raw Ingest publishing settings.
type PublishingConfig struct {
	Enabled     bool
	Host       string
	Port       uint16
	UnitID     uint8
	Interval   time.Duration
}

// TemperatureUnit specifies the temperature unit.
type TemperatureUnit int

const (
	Celsius TemperatureUnit = iota
	Fahrenheit
)

// DefaultConfig returns reasonable defaults.
func DefaultConfig() Config {
	return Config{
		ID:   "weather-station-001",
		Name: "Weather Station 001",
		Units: Celsius,
		RegisterBase: 0,
		Publishing: PublishingConfig{
			Enabled:   false,
			Host:     "localhost",
			Port:     500,
			UnitID:   1,
			Interval: 1 * time.Second,
		},
	}
}

// Memory register addresses.
const (
	RegisterTemperature   = 0
	RegisterHumidity     = 1
	RegisterPressure     = 2
	RegisterCloudCover   = 3
	RegisterWindSpeed    = 4
	RegisterWindDirection = 5
	RegisterRainStatus   = 6
	RegisterStatus       = 100
)

// Memory value names.
const (
	MemoryTemperature    = "temperature"
	MemoryHumidity      = "humidity"
	MemoryPressure      = "pressure"
	MemoryCloudCover    = "cloud_cover"
	MemoryWindSpeed     = "wind_speed"
	MemoryWindDirection = "wind_direction"
	MemoryRainStatus    = "rain_status"
	MemoryStatus        = "status"
	MemoryTickCount     = "tick_count"
)
