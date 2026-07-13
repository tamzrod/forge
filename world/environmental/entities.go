// Package environmental provides entity models for environmental conditions.
package environmental

import (
	"math"
	"time"

	"github.com/tamzrod/forge/world"
)

// SunEntity represents the sun and solar irradiance.
type SunEntity struct {
	world.BaseEntity
	azimuth   float32 // degrees from north
	elevation float32 // degrees above horizon
	irradiance float32 // W/m^2
	timeOfDay float32 // hours
}

// NewSun creates a new sun entity.
func NewSun(id world.EntityID) *SunEntity {
	e := &SunEntity{
		azimuth:    180,
		elevation:  0,
		irradiance: 0,
		timeOfDay:  6, // Start at 6 AM
	}
	e.BaseEntity = world.NewBaseEntity(id, "sun")
	return e
}

// Tick updates sun position and irradiance.
func (e *SunEntity) Tick(dt time.Duration) {
	// Advance time of day
	e.timeOfDay += float32(dt.Seconds()) / 3600.0
	if e.timeOfDay >= 24 {
		e.timeOfDay -= 24
	}

	// Calculate solar position
	// Simplified: peak at noon, zero at night
	hour := e.timeOfDay
	if hour < 6 || hour > 18 {
		e.elevation = 0
		e.irradiance = 0
	} else {
		// Elevation: peaks at 90 degrees at noon
		noonOffset := hour - 12
		e.elevation = float32(90 * math.Cos(float64(noonOffset)*math.Pi/12.0))
		if e.elevation < 0 {
			e.elevation = 0
		}

		// Irradiance: peaks at 1000 W/m^2 at noon
		if e.elevation > 0 {
			e.irradiance = float32(1000 * math.Sqrt(float64(e.elevation)/90.0))
		} else {
			e.irradiance = 0
		}
	}

	// Azimuth follows sun path
	e.azimuth = 180 - (hour-6)*15
	if e.azimuth < 0 {
		e.azimuth += 360
	}
	if e.azimuth > 360 {
		e.azimuth -= 360
	}

	e.SetOutput("azimuth", e.azimuth)
	e.SetOutput("elevation", e.elevation)
	e.SetOutput("irradiance", e.irradiance)
}

// Measurements returns sun measurements.
func (e *SunEntity) Measurements() []world.Measurement {
	return []world.Measurement{
		{EntityID: e.ID(), Name: "azimuth", Value: e.azimuth, Unit: "deg", Timestamp: time.Now()},
		{EntityID: e.ID(), Name: "elevation", Value: e.elevation, Unit: "deg", Timestamp: time.Now()},
		{EntityID: e.ID(), Name: "irradiance", Value: e.irradiance, Unit: "W/m2", Timestamp: time.Now()},
	}
}

// HandleEvent handles cloud cover events.
func (e *SunEntity) HandleEvent(evt world.Event) {
	if evt.Type == "cloud_cover" {
		if coverage, ok := evt.Data["coverage"].(float32); ok {
			e.irradiance *= (1.0 - coverage*0.8) // Clouds reduce irradiance
		}
	}
}

// WeatherEntity represents ambient weather conditions.
type WeatherEntity struct {
	world.BaseEntity
	temperature float32 // Celsius
	humidity    float32 // Percent
	pressure    float32 // hPa
	windSpeed   float32 // m/s
	cloudCover  float32 // 0-1
}

// NewWeather creates a new weather entity.
func NewWeather(id world.EntityID) *WeatherEntity {
	e := &WeatherEntity{
		temperature: 25,
		humidity:    50,
		pressure:    1013.25,
		windSpeed:   5,
		cloudCover:  0,
	}
	e.BaseEntity = world.NewBaseEntity(id, "weather")
	return e
}

// Tick updates weather conditions.
func (e *WeatherEntity) Tick(dt time.Duration) {
	// Slight random variations
	noise := func(base, range_ float32) float32 {
		t := float32(time.Now().UnixNano() % 1000)
		return base + (t/1000.0-0.5)*range_
	}

	e.temperature = noise(e.temperature, 0.5)
	e.humidity = noise(e.humidity, 2)
	e.windSpeed = noise(e.windSpeed, 1)

	e.SetOutput("temperature", e.temperature)
	e.SetOutput("humidity", e.humidity)
	e.SetOutput("pressure", e.pressure)
	e.SetOutput("wind_speed", e.windSpeed)
	e.SetOutput("cloud_cover", e.cloudCover)
}

// Measurements returns weather measurements.
func (e *WeatherEntity) Measurements() []world.Measurement {
	return []world.Measurement{
		{EntityID: e.ID(), Name: "temperature", Value: e.temperature, Unit: "C", Timestamp: time.Now()},
		{EntityID: e.ID(), Name: "humidity", Value: e.humidity, Unit: "%", Timestamp: time.Now()},
		{EntityID: e.ID(), Name: "pressure", Value: e.pressure, Unit: "hPa", Timestamp: time.Now()},
		{EntityID: e.ID(), Name: "wind_speed", Value: e.windSpeed, Unit: "m/s", Timestamp: time.Now()},
		{EntityID: e.ID(), Name: "cloud_cover", Value: e.cloudCover, Unit: "pu", Timestamp: time.Now()},
	}
}

// HandleEvent handles weather events.
func (e *WeatherEntity) HandleEvent(evt world.Event) {
	switch evt.Type {
	case "temperature_change":
		if temp, ok := evt.Data["temperature"].(float32); ok {
			e.temperature = temp
		}
	case "wind_gust":
		if speed, ok := evt.Data["speed"].(float32); ok {
			e.windSpeed = speed
		}
	case "cloud_cover":
		if cover, ok := evt.Data["coverage"].(float32); ok {
			e.cloudCover = cover
		}
	}
}

// PVArrayEntity represents a solar PV array.
type PVArrayEntity struct {
	world.BaseEntity
	irradiance  float32
	powerOutput float32
	ratedPower  float32
	panelArea   float32
	efficiency  float32
}

// NewPVArray creates a new PV array entity.
func NewPVArray(id world.EntityID, ratedPowerKW float32) *PVArrayEntity {
	e := &PVArrayEntity{
		powerOutput: 0,
		ratedPower:  ratedPowerKW,
		panelArea:   ratedPowerKW * 10, // 10 m^2 per kW
		efficiency:  0.18,
	}
	e.BaseEntity = world.NewBaseEntity(id, "pv_array")
	return e
}

// Tick updates PV array power output.
func (e *PVArrayEntity) Tick(dt time.Duration) {
	// Get irradiance from sun
	if v := e.GetInput("irradiance"); v != nil {
		e.irradiance = v.(float32)
	}

	// Calculate power output
	// Power = irradiance * area * efficiency
	e.powerOutput = e.irradiance * e.panelArea * e.efficiency / 1000.0 // kW

	// Limit to rated power
	if e.powerOutput > e.ratedPower {
		e.powerOutput = e.ratedPower
	}

	e.SetOutput("power", e.powerOutput)
}

// Measurements returns PV array measurements.
func (e *PVArrayEntity) Measurements() []world.Measurement {
	return []world.Measurement{
		{EntityID: e.ID(), Name: "power", Value: e.powerOutput, Unit: "kW", Timestamp: time.Now()},
		{EntityID: e.ID(), Name: "irradiance", Value: e.irradiance, Unit: "W/m2", Timestamp: time.Now()},
		{EntityID: e.ID(), Name: "rated_power", Value: e.ratedPower, Unit: "kW", Timestamp: time.Now()},
	}
}
