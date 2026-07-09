// Package sun provides a simple solar position and irradiance model.
//
// This model calculates solar position and irradiance based on:
// - Simulation clock time
// - Geographic location (latitude, longitude)
//
// The model uses simplified algorithms suitable for industrial simulation.
// For software testing purposes, we prioritize believable behavior over
// astronomical precision.
//
// References:
// - NOAA Solar Calculator (simplified)
// - SPA algorithm concepts (for more complex implementations)
package sun

import (
	"math"

	"github.com/tamzrod/forge/internal/models/clock"
)

const (
	// DegreesToRadians converts degrees to radians.
	DegreesToRadians = math.Pi / 180.0
	// RadiansToDegrees converts radians to degrees.
	RadiansToDegrees = 180.0 / math.Pi

	// SolarConstant is the solar constant in W/m².
	SolarConstant = 1361.0

	// MaxIrradiance is the maximum expected irradiance at sea level.
	MaxIrradiance = 1000.0
)

// Sun represents the sun's position and irradiance.
type Sun struct {
	clock    *clock.Clock
	latitude float64  // degrees
	longitude float64 // degrees

	// Cached values (updated each tick)
	elevation  float64 // degrees above horizon
	azimuth    float64 // degrees from north
	irradiance float64 // W/m²
}

// Config holds sun model configuration.
type Config struct {
	// Latitude in degrees (positive = north)
	Latitude float64

	// Longitude in degrees (positive = east)
	Longitude float64
}

// New creates a new sun model.
func New(cfg Config, simClock *clock.Clock) *Sun {
	return &Sun{
		clock:     simClock,
		latitude:  cfg.Latitude,
		longitude: cfg.Longitude,
		elevation: 0,
		azimuth:   0,
		irradiance: 0,
	}
}

// Elevation returns the solar elevation angle in degrees.
// This is the angle above the horizon.
func (s *Sun) Elevation() float64 {
	return s.elevation
}

// Azimuth returns the solar azimuth angle in degrees.
// This is the compass direction (0 = north, 90 = east, 180 = south, 270 = west).
func (s *Sun) Azimuth() float64 {
	return s.azimuth
}

// Irradiance returns the global horizontal irradiance in W/m².
func (s *Sun) Irradiance() float64 {
	return s.irradiance
}

// DirectNormalIrradiance returns the direct normal irradiance in W/m².
func (s *Sun) DirectNormalIrradiance() float64 {
	if s.elevation <= 0 {
		return 0
	}
	// Approximate using cosine of zenith angle
	zenith := 90.0 - s.elevation
	cosZenith := math.Max(0, math.Cos(zenith*DegreesToRadians))
	return s.irradiance / math.Max(0.01, cosZenith)
}

// DiffuseHorizontalIrradiance returns the diffuse horizontal irradiance in W/m².
func (s *Sun) DiffuseHorizontalIrradiance() float64 {
	if s.elevation <= 0 {
		return 0
	}
	// Simple model: diffuse is a fraction of global
	clearness := 0.7 + 0.3*(s.elevation/90.0)
	return s.irradiance * (1.0 - clearness)
}

// Tick updates the sun position based on simulation time.
func (s *Sun) Tick() {
	s.update()
}

// Latitude returns the model's latitude.
func (s *Sun) Latitude() float64 {
	return s.latitude
}

// Longitude returns the model's longitude.
func (s *Sun) Longitude() float64 {
	return s.longitude
}

// update recalculates sun position and irradiance.
func (s *Sun) update() {
	elapsed := s.clock.Elapsed()

	// Convert elapsed time to hours
	hours := elapsed.Hours()

	// Calculate day of year (simplified: assume starting at midnight Jan 1)
	// In a full implementation, you'd track calendar date from the clock
	dayOfYear := 1 + int(hours/24.0)
	hourOfDay := math.Mod(hours, 24.0)

	// Solar declination (simplified)
	declination := 23.45 * math.Sin(DegreesToRadians*360.0*(float64(dayOfYear)-81.0)/365.0)

	// Hour angle (0 at solar noon, +/- 15° per hour)
	hourAngle := 15.0 * (hourOfDay - 12.0)

	// Solar altitude (elevation)
	sinAlt := math.Sin(s.latitude*DegreesToRadians)*math.Sin(declination*DegreesToRadians) +
		math.Cos(s.latitude*DegreesToRadians)*math.Cos(declination*DegreesToRadians)*math.Cos(hourAngle*DegreesToRadians)
	s.elevation = math.Asin(sinAlt) * RadiansToDegrees

	// Solar azimuth
	cosAz := (math.Sin(s.latitude*DegreesToRadians)*sinAlt - math.Sin(declination*DegreesToRadians)) /
		(math.Cos(s.latitude*DegreesToRadians) * math.Max(0.01, math.Cos(s.elevation*DegreesToRadians)))
	cosAz = math.Max(-1, math.Min(1, cosAz))
	azFromEast := math.Acos(cosAz) * RadiansToDegrees

	// Convert azimuth to compass direction (0 = north)
	if hourAngle > 0 {
		s.azimuth = 180.0 - azFromEast
	} else {
		s.azimuth = 180.0 + azFromEast
	}

	// Calculate irradiance
	if s.elevation <= 0 {
		// Sun below horizon
		s.irradiance = 0
	} else {
		// Air mass (simplified)
		airMass := 1.0 / math.Max(0.01, math.Sin(s.elevation*DegreesToRadians))

		// Atmospheric transmission (simplified)
		transmission := math.Exp(-0.0001 * airMass)

		// Extraterrestrial irradiance
		extraterrestrial := SolarConstant * (1.0 + 0.033*math.Cos(360.0*float64(dayOfYear)*DegreesToRadians/365.0))

		// Ground reflection (albedo effect)
		albedo := 0.2

		// Direct irradiance
		directHI := extraterrestrial * transmission * math.Sin(s.elevation*DegreesToRadians)

		// Diffuse irradiance (simplified sky model)
		diffuseHI := extraterrestrial * 0.3 * (1.0 - transmission) * math.Sin(s.elevation*DegreesToRadians)

		// Ground reflected
		reflected := extraterrestrial * albedo * math.Sin(s.elevation*DegreesToRadians) * 0.5

		s.irradiance = directHI + diffuseHI + reflected

		// Cap at physically reasonable maximum
		if s.irradiance > MaxIrradiance {
			s.irradiance = MaxIrradiance
		}
	}
}

// IsDaytime returns true if the sun is above the horizon.
func (s *Sun) IsDaytime() bool {
	return s.elevation > 0
}
