// Package inspector provides a development tool for visualizing simulation state.
//
// This is NOT a Virtual Device.
// This is NOT part of the Runtime.
// This is NOT connected to MMA2 or any protocols.
//
// The Simulation Inspector is a read-only developer tool that visualizes
// the internal state of Simulation Models in real time.
//
// Use cases:
// - Verify the world is alive during development
// - Debug simulation behavior
// - Understand model interactions
//
// The Inspector must NEVER:
// - Modify simulation state
// - Connect to MMA2
// - Expose protocols
// - Be used in production
package inspector

import (
	"time"

	"github.com/tamzrod/forge/internal/models/clock"
	"github.com/tamzrod/forge/internal/models/grid"
	"github.com/tamzrod/forge/internal/models/sun"
	"github.com/tamzrod/forge/internal/models/weather"
)

// View provides read-only access to simulation state.
// All methods return copies or are safe for concurrent reading.
type View struct {
	clock   *clock.Clock
	sun     *sun.Sun
	weather *weather.Weather
	grid    *grid.Grid
}

// NewView creates a new read-only view of the simulation state.
func NewView(simClock *clock.Clock, sunModel *sun.Sun, weatherModel *weather.Weather, gridModel *grid.Grid) *View {
	return &View{
		clock:   simClock,
		sun:     sunModel,
		weather: weatherModel,
		grid:    gridModel,
	}
}

// ClockState returns the current clock state.
func (v *View) ClockState() ClockState {
	return ClockState{
		Elapsed:    v.clock.Elapsed(),
		TickCount:  v.clock.TickCount(),
		Mode:       v.clock.Mode().String(),
		IsPaused:   v.clock.IsPaused(),
	}
}

// SunState returns the current sun model state.
func (v *View) SunState() SunState {
	return SunState{
		Elevation:        v.sun.Elevation(),
		Azimuth:          v.sun.Azimuth(),
		Irradiance:       v.sun.Irradiance(),
		DirectNormal:     v.sun.DirectNormalIrradiance(),
		Diffuse:          v.sun.DiffuseHorizontalIrradiance(),
		IsDaytime:        v.sun.IsDaytime(),
		Latitude:         v.sun.Latitude(),
		Longitude:        v.sun.Longitude(),
	}
}

// WeatherState returns the current weather model state.
func (v *View) WeatherState() WeatherState {
	return WeatherState{
		Temperature:   v.weather.Temperature(),
		Humidity:      v.weather.Humidity(),
		Pressure:      v.weather.Pressure(),
		CloudCover:    v.weather.CloudCover(),
		WindSpeed:     v.weather.WindSpeed(),
		WindDirection: v.weather.WindDirection(),
		IsRaining:    v.weather.IsRaining(),
	}
}

// GridState returns the current grid model state.
func (v *View) GridState() GridState {
	return GridState{
		Voltage:          v.grid.Voltage(),
		Frequency:        v.grid.Frequency(),
		VoltagePU:        v.grid.VoltagePU(),
		FrequencyPU:      v.grid.FrequencyPU(),
		ActiveBalance:    v.grid.ActivePowerBalance(),
		ReactiveBalance:  v.grid.ReactivePowerBalance(),
		IsStable:         v.grid.IsStable(),
		NominalVoltage:   v.grid.NominalVoltage(),
		NominalFrequency: v.grid.NominalFrequency(),
	}
}

// FullState returns a complete snapshot of all simulation state.
func (v *View) FullState() State {
	return State{
		Clock:   v.ClockState(),
		Sun:     v.SunState(),
		Weather: v.WeatherState(),
		Grid:    v.GridState(),
	}
}

// ClockState represents the simulation clock state.
type ClockState struct {
	Elapsed   time.Duration `json:"elapsed"`
	TickCount uint64        `json:"tick_count"`
	Mode      string        `json:"mode"`
	IsPaused  bool          `json:"is_paused"`
}

// SunState represents the sun model state.
type SunState struct {
	Elevation      float64 `json:"elevation"`       // degrees above horizon
	Azimuth        float64 `json:"azimuth"`         // degrees from north
	Irradiance     float64 `json:"irradiance"`      // W/m² GHI
	DirectNormal   float64 `json:"direct_normal"`   // W/m² DNI
	Diffuse        float64 `json:"diffuse"`         // W/m² diffuse
	IsDaytime      bool    `json:"is_daytime"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
}

// WeatherState represents the weather model state.
type WeatherState struct {
	Temperature   float64 `json:"temperature"`    // °C
	Humidity      float64 `json:"humidity"`       // %
	Pressure      float64 `json:"pressure"`       // hPa
	CloudCover    float64 `json:"cloud_cover"`    // fraction 0-1
	WindSpeed     float64 `json:"wind_speed"`     // m/s
	WindDirection float64 `json:"wind_direction"` // degrees from north
	IsRaining     bool    `json:"is_raining"`
}

// GridState represents the grid model state.
type GridState struct {
	Voltage          float64 `json:"voltage"`           // V
	Frequency        float64 `json:"frequency"`         // Hz
	VoltagePU        float64 `json:"voltage_pu"`       // per-unit
	FrequencyPU      float64 `json:"frequency_pu"`     // per-unit
	ActiveBalance    float64 `json:"active_balance"`   // MW
	ReactiveBalance  float64 `json:"reactive_balance"` // MVAr
	IsStable         bool    `json:"is_stable"`
	NominalVoltage   float64 `json:"nominal_voltage"`
	NominalFrequency float64 `json:"nominal_frequency"`
}

// State represents the complete simulation state.
type State struct {
	Clock   ClockState `json:"clock"`
	Sun     SunState   `json:"sun"`
	Weather WeatherState `json:"weather"`
	Grid    GridState  `json:"grid"`
}
