package devices

import (
	"github.com/tamzrod/forge/internal/models/clock"
	"github.com/tamzrod/forge/internal/models/grid"
	"github.com/tamzrod/forge/internal/models/sun"
	"github.com/tamzrod/forge/internal/models/weather"
)

// Context provides read-only access to Simulation Models.
//
// Virtual Devices use this context to observe the simulated world.
// Simulation Models are unaware of devices.
//
// Key principles:
// - Devices can only READ from models
// - Devices cannot modify model state directly
// - Models don't know devices exist
type Context struct {
	clock   *clock.Clock
	sun     *sun.Sun
	weather *weather.Weather
	grid    *grid.Grid
}

// NewContext creates a new simulation context.
func NewContext(simClock *clock.Clock, sunModel *sun.Sun, weatherModel *weather.Weather, gridModel *grid.Grid) *Context {
	return &Context{
		clock:   simClock,
		sun:     sunModel,
		weather: weatherModel,
		grid:    gridModel,
	}
}

// Clock returns the simulation clock (read-only).
func (c *Context) Clock() *clock.Clock {
	return c.clock
}

// Sun returns the sun model (read-only).
func (c *Context) Sun() *sun.Sun {
	return c.sun
}

// Weather returns the weather model (read-only).
func (c *Context) Weather() *weather.Weather {
	return c.weather
}

// Grid returns the grid model (read-only).
func (c *Context) Grid() *grid.Grid {
	return c.grid
}

// WeatherSnapshot represents a point-in-time snapshot of weather data.
type WeatherSnapshot struct {
	Temperature   float64
	Humidity      float64
	Pressure      float64
	CloudCover    float64
	WindSpeed     float64
	WindDirection float64
	IsRaining     bool
}

// ReadWeather copies the current weather values into a snapshot.
// This provides a consistent view of weather at a point in time.
func (c *Context) ReadWeather() WeatherSnapshot {
	return WeatherSnapshot{
		Temperature:   c.weather.Temperature(),
		Humidity:      c.weather.Humidity(),
		Pressure:      c.weather.Pressure(),
		CloudCover:    c.weather.CloudCover(),
		WindSpeed:     c.weather.WindSpeed(),
		WindDirection: c.weather.WindDirection(),
		IsRaining:     c.weather.IsRaining(),
	}
}

// GridSnapshot represents a point-in-time snapshot of grid data.
type GridSnapshot struct {
	Voltage          float64
	Frequency        float64
	VoltagePU        float64
	FrequencyPU      float64
	IsStable         bool
}

// ReadGrid copies the current grid values into a snapshot.
func (c *Context) ReadGrid() GridSnapshot {
	return GridSnapshot{
		Voltage:     c.grid.Voltage(),
		Frequency:   c.grid.Frequency(),
		VoltagePU:   c.grid.VoltagePU(),
		FrequencyPU: c.grid.FrequencyPU(),
		IsStable:    c.grid.IsStable(),
	}
}

// SunSnapshot represents a point-in-time snapshot of sun data.
type SunSnapshot struct {
	Elevation      float64
	Azimuth        float64
	Irradiance     float64
	DirectNormal   float64
	Diffuse        float64
	IsDaytime      bool
}

// ReadSun copies the current sun values into a snapshot.
func (c *Context) ReadSun() SunSnapshot {
	return SunSnapshot{
		Elevation:    c.sun.Elevation(),
		Azimuth:      c.sun.Azimuth(),
		Irradiance:   c.sun.Irradiance(),
		DirectNormal: c.sun.DirectNormalIrradiance(),
		Diffuse:      c.sun.DiffuseHorizontalIrradiance(),
		IsDaytime:    c.sun.IsDaytime(),
	}
}
