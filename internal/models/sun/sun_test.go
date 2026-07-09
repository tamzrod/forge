package sun

import (
	"math"
	"testing"
	"time"

	"github.com/tamzrod/forge/internal/models/clock"
)

func TestSun_Creation(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: time.Hour,
	})

	cfg := Config{
		Latitude:  40.0,  // Denver
		Longitude: -105.0,
	}

	sun := New(cfg, simClock)

	if sun.Latitude() != 40.0 {
		t.Errorf("expected latitude 40.0, got %f", sun.Latitude())
	}

	if sun.Longitude() != -105.0 {
		t.Errorf("expected longitude -105.0, got %f", sun.Longitude())
	}
}

func TestSun_NightProducesZeroIrradiance(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: time.Hour,
	})

	// Start at midnight
	cfg := Config{
		Latitude:  40.0,
		Longitude: -105.0,
	}

	sun := New(cfg, simClock)
	sun.Tick()

	if sun.Irradiance() != 0 {
		t.Errorf("expected midnight irradiance 0, got %f", sun.Irradiance())
	}

	if sun.IsDaytime() {
		t.Error("expected nighttime at midnight")
	}
}

func TestSun_NoonProducesMaximumIrradiance(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: time.Hour,
	})

	// Simulate solar noon on summer solstice
	// Midnight Jan 1 + 12 hours = noon Jan 1
	// Day 172 = summer solstice
	simClock.Advance(172*24*time.Hour + 12*time.Hour)

	cfg := Config{
		Latitude:  40.0,
		Longitude: -105.0,
	}

	sun := New(cfg, simClock)
	sun.Tick()

	if sun.Irradiance() <= 0 {
		t.Errorf("expected positive irradiance at noon, got %f", sun.Irradiance())
	}

	if sun.Irradiance() > MaxIrradiance {
		t.Errorf("expected irradiance <= %f at noon, got %f", MaxIrradiance, sun.Irradiance())
	}
}

func TestSun_ElevationRange(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: time.Hour,
	})

	cfg := Config{
		Latitude:  0.0,  // Equator
		Longitude: 0.0,
	}

	sun := New(cfg, simClock)

	// Check throughout a day
	maxElevation := -math.MaxFloat64
	minElevation := math.MaxFloat64

	for hour := 0.0; hour < 24.0; hour++ {
		simClock.Reset()
		simClock.Advance(time.Duration(hour) * time.Hour)
		sun.Tick()

		if sun.Elevation() > maxElevation {
			maxElevation = sun.Elevation()
		}
		if sun.Elevation() < minElevation {
			minElevation = sun.Elevation()
		}
	}

	// At equator, elevation should range from negative (night) to positive (day)
	if maxElevation < 60 {
		t.Errorf("expected max elevation >= 60 at equator, got %f", maxElevation)
	}

	if minElevation > 0 {
		t.Errorf("expected min elevation < 0 at equator (night), got %f", minElevation)
	}
}

func TestSun_AzimuthRange(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: time.Hour,
	})

	cfg := Config{
		Latitude:  40.0,
		Longitude: -105.0,
	}

	sun := New(cfg, simClock)

	// Azimuth should be in range [0, 360)
	for hour := 0.0; hour < 24.0; hour++ {
		simClock.Reset()
		simClock.Advance(time.Duration(hour) * time.Hour)
		sun.Tick()

		if sun.Azimuth() < 0 || sun.Azimuth() >= 360 {
			t.Errorf("azimuth out of range [0, 360): %f at hour %f", sun.Azimuth(), hour)
		}
	}
}

func TestSun_Deterministic(t *testing.T) {
	simClock1 := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: time.Hour,
	})

	simClock2 := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: time.Hour,
	})

	cfg := Config{
		Latitude:  40.0,
		Longitude: -105.0,
	}

	sun1 := New(cfg, simClock1)
	sun2 := New(cfg, simClock2)

	// Advance both clocks identically
	for i := 0; i < 100; i++ {
		simClock1.Tick()
		simClock2.Tick()
		sun1.Tick()
		sun2.Tick()
	}

	// Results should be identical
	if sun1.Elevation() != sun2.Elevation() {
		t.Errorf("elevation mismatch: %f vs %f", sun1.Elevation(), sun2.Elevation())
	}

	if sun1.Azimuth() != sun2.Azimuth() {
		t.Errorf("azimuth mismatch: %f vs %f", sun1.Azimuth(), sun2.Azimuth())
	}

	if sun1.Irradiance() != sun2.Irradiance() {
		t.Errorf("irradiance mismatch: %f vs %f", sun1.Irradiance(), sun2.Irradiance())
	}
}

func TestSun_DirectNormalIrradiance(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: time.Hour,
	})

	// Solar noon on equinox at equator
	simClock.Advance(80*24*time.Hour + 12*time.Hour) // ~March 21

	cfg := Config{
		Latitude:  0.0,
		Longitude: 0.0,
	}

	sun := New(cfg, simClock)
	sun.Tick()

	dni := sun.DirectNormalIrradiance()

	if dni <= 0 {
		t.Errorf("expected positive DNI at noon, got %f", dni)
	}

	if dni < sun.Irradiance() {
		t.Errorf("expected DNI >= GHI at noon, got DNI=%f, GHI=%f", dni, sun.Irradiance())
	}
}

func TestSun_NightDNI(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: time.Hour,
	})

	cfg := Config{
		Latitude:  40.0,
		Longitude: -105.0,
	}

	sun := New(cfg, simClock)
	// Midnight
	sun.Tick()

	if sun.DirectNormalIrradiance() != 0 {
		t.Errorf("expected DNI 0 at night, got %f", sun.DirectNormalIrradiance())
	}

	if sun.DiffuseHorizontalIrradiance() != 0 {
		t.Errorf("expected diffuse 0 at night, got %f", sun.DiffuseHorizontalIrradiance())
	}
}

func TestSun_DifferentLatitudes(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: time.Hour,
	})

	// Solar noon
	simClock.Advance(80*24*time.Hour + 12*time.Hour)

	testCases := []struct {
		latitude    float64
		minElevaton float64
	}{
		{0.0, 66.0},   // Equator: max elevation ~90 - |declination|
		{45.0, 21.0},  // 45°N: max elevation ~90 - |45 - declination|
		{-45.0, 21.0}, // -45°S
		{90.0, 0.0},   // North pole
		{-90.0, 0.0},  // South pole
	}

	for _, tc := range testCases {
		cfg := Config{
			Latitude:  tc.latitude,
			Longitude: 0.0,
		}

		sun := New(cfg, simClock)
		sun.Tick()

		if sun.Elevation() < tc.minElevaton {
			t.Errorf("latitude %f: expected elevation >= %f at noon, got %f",
				tc.latitude, tc.minElevaton, sun.Elevation())
		}
	}
}
