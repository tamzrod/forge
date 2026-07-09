// Package weather provides a simple weather simulation model.
//
// The weather model simulates ambient atmospheric conditions based on:
// - Simulation clock time
// - Sun model (for day/night effects)
//
// Weather evolves deterministically with a seeded random number generator.
// This ensures reproducible weather patterns for testing.
//
// The model is intentionally simple—believable behavior over physical accuracy.
package weather

import (
	"math"
	"math/rand"

	"github.com/tamzrod/forge/internal/models/clock"
)

// Weather represents the ambient weather conditions.
type Weather struct {
	clock *clock.Clock
	seed  int64

	// State
	temperature  float64 // °C
	humidity     float64 // % relative humidity (0-100)
	pressure     float64 // hPa
	cloudCover   float64 // fraction (0-1)
	windSpeed    float64 // m/s
	windDirection float64 // degrees from north

	// Internal RNG (seeded for determinism)
	rng *rand.Rand

	// Baseline values (can be configured)
	baseTemp float64
	tempRange float64
	baseHumidity float64
}

// Config holds weather model configuration.
type Config struct {
	// Base temperature in °C
	BaseTemperature float64

	// Daily temperature range in °C
	TemperatureRange float64

	// Base humidity in %
	BaseHumidity float64

	// Base pressure in hPa
	BasePressure float64

	// Seed for deterministic randomness (0 = use wall clock)
	Seed int64
}

// DefaultConfig returns a reasonable default configuration.
func DefaultConfig() Config {
	return Config{
		BaseTemperature: 20.0,
		TemperatureRange: 10.0,
		BaseHumidity:     50.0,
		BasePressure:     1013.25,
		Seed:             42, // Deterministic by default
	}
}

// New creates a new weather model.
func New(cfg Config, simClock *clock.Clock) *Weather {
	seed := cfg.Seed
	if seed == 0 {
		// Use a fixed seed based on clock for reproducibility
		seed = 12345
	}

	return &Weather{
		clock:         simClock,
		seed:          seed,
		temperature:   cfg.BaseTemperature,
		humidity:     cfg.BaseHumidity,
		pressure:      cfg.BasePressure,
		cloudCover:    0.3,
		windSpeed:     5.0,
		windDirection: 0,
		baseTemp:      cfg.BaseTemperature,
		tempRange:      cfg.TemperatureRange,
		baseHumidity:   cfg.BaseHumidity,
		rng:           rand.New(rand.NewSource(seed)),
	}
}

// Temperature returns the ambient temperature in °C.
func (w *Weather) Temperature() float64 {
	return w.temperature
}

// Humidity returns the relative humidity in %.
func (w *Weather) Humidity() float64 {
	return w.humidity
}

// Pressure returns the atmospheric pressure in hPa.
func (w *Weather) Pressure() float64 {
	return w.pressure
}

// CloudCover returns the cloud cover fraction (0-1).
func (w *Weather) CloudCover() float64 {
	return w.cloudCover
}

// WindSpeed returns the wind speed in m/s.
func (w *Weather) WindSpeed() float64 {
	return w.windSpeed
}

// WindDirection returns the wind direction in degrees from north.
func (w *Weather) WindDirection() float64 {
	return w.windDirection
}

// Tick advances the weather model by one simulation step.
func (w *Weather) Tick() {
	w.evolve()
}

// evolve updates weather based on time and sun position.
// Uses deterministic randomness for reproducibility.
func (w *Weather) evolve() {
	elapsed := w.clock.Elapsed()
	hours := elapsed.Hours()

	// Daily temperature cycle (sine wave)
	dayHours := math.Mod(hours, 24.0)
	tempPhase := 2 * math.Pi * (dayHours - 6.0) / 24.0 // Peak at 6 PM
	dailyTemp := w.baseTemp + w.tempRange*math.Sin(tempPhase)

	// Add slow random variation
	dailyVariation := w.noise(0.1) * 3.0

	// Combine
	w.temperature = dailyTemp + dailyVariation
	w.temperature = clamp(w.temperature, -40, 50)

	// Humidity inversely correlated with temperature
	humidityBase := w.baseHumidity
	humidityTempEffect := -(w.temperature - w.baseTemp) * 2.0
	humidityNoise := w.noise(0.2) * 10.0
	w.humidity = humidityBase + humidityTempEffect + humidityNoise
	w.humidity = clamp(w.humidity, 0, 100)

	// Pressure slowly varies (high pressure systems)
	pressureNoise := w.noise(0.05) * 5.0
	w.pressure = 1013.25 + pressureNoise
	w.pressure = clamp(w.pressure, 950, 1050)

	// Cloud cover evolves slowly
	cloudNoise := w.noise(0.3) * 0.2
	w.cloudCover += cloudNoise
	w.cloudCover = clamp(w.cloudCover, 0, 1)

	// Wind speed with gusts
	baseWind := 5.0 + w.noise(0.4)*3.0
	gust := 0.0
	if w.rng.Float64() > 0.9 {
		gust = w.rng.Float64() * 10.0
	}
	w.windSpeed = baseWind + gust
	w.windSpeed = clamp(w.windSpeed, 0, 50)

	// Wind direction slowly drifts
	dirDrift := w.noise(0.1) * 10.0
	w.windDirection += dirDrift
	if w.windDirection < 0 {
		w.windDirection += 360
	}
	if w.windDirection >= 360 {
		w.windDirection -= 360
	}
}

// noise returns a pseudo-random value between -1 and 1.
// Uses tick count for determinism.
func (w *Weather) noise(frequency float64) float64 {
	// Combine multiple sine waves for "noise-like" behavior
	tick := float64(w.clock.TickCount())
	val := math.Sin(tick*frequency) +
		0.5*math.Sin(tick*frequency*1.7+1.0) +
		0.25*math.Sin(tick*frequency*2.3+2.0)
	return val / 1.75 // Normalize to approximately [-1, 1]
}

// Seed returns the model's random seed.
func (w *Weather) Seed() int64 {
	return w.seed
}

// Reset resets the weather to initial conditions.
func (w *Weather) Reset() {
	w.rng = rand.New(rand.NewSource(w.seed))
	w.temperature = w.baseTemp
	w.humidity = w.baseHumidity
	w.pressure = 1013.25
	w.cloudCover = 0.3
	w.windSpeed = 5.0
	w.windDirection = 0
}

// clamp constrains a value to a range.
func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// IsRaining returns true if rain is occurring.
func (w *Weather) IsRaining() bool {
	return w.cloudCover > 0.7 && w.rng.Float64() > 0.95
}

// RainRate returns the rainfall rate in mm/h.
func (w *Weather) RainRate() float64 {
	if !w.IsRaining() {
		return 0
	}
	return w.rng.Float64() * 10.0
}
