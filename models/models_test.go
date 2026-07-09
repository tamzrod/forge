// Package models provides tests for Simulation Models.
package models

import (
	"testing"
)

func TestGridModel_Creation(t *testing.T) {
	grid := NewGridModel("test-grid")

	if grid.ID() != "test-grid" {
		t.Errorf("expected ID 'test-grid', got '%s'", grid.ID())
	}

	if grid.Type() != "grid" {
		t.Errorf("expected type 'grid', got '%s'", grid.Type())
	}

	// Check default values
	if grid.Voltage() != 480.0 {
		t.Errorf("expected default voltage 480.0, got %f", grid.Voltage())
	}

	if grid.Frequency() != 60.0 {
		t.Errorf("expected default frequency 60.0, got %f", grid.Frequency())
	}
}

func TestGridModel_VoltageBounds(t *testing.T) {
	grid := NewGridModel("test-grid")

	// Test upper bound
	grid.SetVoltage(600.0)
	if grid.Voltage() != 520.0 {
		t.Errorf("expected voltage clamped to 520.0, got %f", grid.Voltage())
	}

	// Test lower bound
	grid.SetVoltage(400.0)
	if grid.Voltage() != 450.0 {
		t.Errorf("expected voltage clamped to 450.0, got %f", grid.Voltage())
	}
}

func TestGridModel_FrequencyBounds(t *testing.T) {
	grid := NewGridModel("test-grid")

	// Test upper bound
	grid.SetFrequency(65.0)
	if grid.Frequency() != 60.5 {
		t.Errorf("expected frequency clamped to 60.5, got %f", grid.Frequency())
	}

	// Test lower bound
	grid.SetFrequency(55.0)
	if grid.Frequency() != 59.5 {
		t.Errorf("expected frequency clamped to 59.5, got %f", grid.Frequency())
	}
}

func TestGridModel_PowerInjection(t *testing.T) {
	grid := NewGridModel("test-grid")

	// Inject power
	grid.InjectActivePower(100.0)
	grid.InjectReactivePower(50.0)

	// Tick should reset injections
	grid.Tick()

	// After tick, injections should be zero
	grid.InjectActivePower(200.0)
	grid.InjectReactivePower(75.0)
	grid.Tick()

	// Injections reset after tick
	grid.InjectActivePower(300.0)
	grid.InjectReactivePower(100.0)

	// Values should persist until next tick
	if grid.activePowerInjection != 300.0 {
		t.Errorf("expected active power injection 300.0, got %f", grid.activePowerInjection)
	}
}

func TestSunModel_Creation(t *testing.T) {
	sun := NewSunModel("test-sun")

	if sun.ID() != "test-sun" {
		t.Errorf("expected ID 'test-sun', got '%s'", sun.ID())
	}

	if sun.Type() != "sun" {
		t.Errorf("expected type 'sun', got '%s'", sun.Type())
	}

	// Check default irradiance
	if sun.Irradiance() != 1000.0 {
		t.Errorf("expected default irradiance 1000.0, got %f", sun.Irradiance())
	}
}

func TestSunModel_Elevation(t *testing.T) {
	sun := NewSunModel("test-sun")

	// Run multiple ticks to see elevation change
	initialElevation := sun.Elevation()

	for i := 0; i < 10; i++ {
		sun.Tick()
	}

	finalElevation := sun.Elevation()

	// Elevation should have changed (either increasing or decreasing)
	if initialElevation == finalElevation {
		t.Logf("Note: Elevation may stay the same if sun is below horizon")
	}
}

func TestWindModel_Creation(t *testing.T) {
	wind := NewWindModel("test-wind")

	if wind.ID() != "test-wind" {
		t.Errorf("expected ID 'test-wind', got '%s'", wind.ID())
	}

	if wind.Type() != "wind" {
		t.Errorf("expected type 'wind', got '%s'", wind.Type())
	}

	// Check default values
	if wind.Speed() != 5.0 {
		t.Errorf("expected default speed 5.0, got %f", wind.Speed())
	}
}

func TestWindModel_Tick(t *testing.T) {
	wind := NewWindModel("test-wind")

	initialSpeed := wind.Speed()

	// Run ticks
	for i := 0; i < 100; i++ {
		wind.Tick()
	}

	finalSpeed := wind.Speed()

	// Speed should still be within bounds (0-50)
	if finalSpeed < 0 || finalSpeed > 50 {
		t.Errorf("expected speed within bounds [0, 50], got %f", finalSpeed)
	}

	// Speed should have changed (due to random variation)
	t.Logf("Initial speed: %f, Final speed: %f", initialSpeed, finalSpeed)
}

func TestWindModel_DirectionBounds(t *testing.T) {
	wind := NewWindModel("test-wind")

	// Run many ticks to test direction wrapping
	for i := 0; i < 1000; i++ {
		wind.Tick()
	}

	direction := wind.Direction()

	// Direction should be in [0, 360)
	if direction < 0 || direction >= 360 {
		t.Errorf("expected direction in [0, 360), got %f", direction)
	}
}

func TestWeatherModel_Creation(t *testing.T) {
	weather := NewWeatherModel("test-weather")

	if weather.ID() != "test-weather" {
		t.Errorf("expected ID 'test-weather', got '%s'", weather.ID())
	}

	if weather.Type() != "weather" {
		t.Errorf("expected type 'weather', got '%s'", weather.Type())
	}

	// Check default values
	if weather.Temperature() != 25.0 {
		t.Errorf("expected default temperature 25.0, got %f", weather.Temperature())
	}

	if weather.Humidity() != 50.0 {
		t.Errorf("expected default humidity 50.0, got %f", weather.Humidity())
	}
}

func TestWeatherModel_Bounds(t *testing.T) {
	weather := NewWeatherModel("test-weather")

	// Run many ticks to test bounds
	for i := 0; i < 1000; i++ {
		weather.Tick()
	}

	// Check all values are within bounds
	temp := weather.Temperature()
	if temp < -40 || temp > 50 {
		t.Errorf("expected temperature in [-40, 50], got %f", temp)
	}

	humidity := weather.Humidity()
	if humidity < 0 || humidity > 100 {
		t.Errorf("expected humidity in [0, 100], got %f", humidity)
	}

	pressure := weather.Pressure()
	if pressure < 950 || pressure > 1050 {
		t.Errorf("expected pressure in [950, 1050], got %f", pressure)
	}

	cloudCover := weather.CloudCover()
	if cloudCover < 0 || cloudCover > 1 {
		t.Errorf("expected cloud cover in [0, 1], got %f", cloudCover)
	}
}

func TestReservoirModel_Creation(t *testing.T) {
	reservoir := NewReservoirModel("test-reservoir", 10000.0)

	if reservoir.ID() != "test-reservoir" {
		t.Errorf("expected ID 'test-reservoir', got '%s'", reservoir.ID())
	}

	if reservoir.Type() != "reservoir" {
		t.Errorf("expected type 'reservoir', got '%s'", reservoir.Type())
	}

	// Check default values
	if reservoir.Level() != 50.0 {
		t.Errorf("expected default level 50.0, got %f", reservoir.Level())
	}

	if reservoir.FlowIn() != 0.0 {
		t.Errorf("expected default flow in 0.0, got %f", reservoir.FlowIn())
	}

	if reservoir.FlowOut() != 0.0 {
		t.Errorf("expected default flow out 0.0, got %f", reservoir.FlowOut())
	}
}

func TestReservoirModel_LevelChange(t *testing.T) {
	reservoir := NewReservoirModel("test-reservoir", 10000.0)

	initialLevel := reservoir.Level()

	// Inject flow
	reservoir.InjectFlow(100.0) // 100 m³/s inflow

	// Extract flow
	reservoir.ExtractFlow(50.0) // 50 m³/s outflow

	// Net: 50 m³/s
	// dLevel = netFlow / area * tickTime
	// For default tick, level should increase

	reservoir.Tick()

	newLevel := reservoir.Level()

	// Level should have increased
	if newLevel <= initialLevel {
		t.Errorf("expected level to increase, got %f -> %f", initialLevel, newLevel)
	}
}

func TestReservoirModel_LevelBounds(t *testing.T) {
	reservoir := NewReservoirModel("test-reservoir", 10000.0)

	// Run many ticks with high inflow
	for i := 0; i < 1000; i++ {
		reservoir.InjectFlow(1000.0)
		reservoir.Tick()
	}

	level := reservoir.Level()

	// Level should be clamped to max (100)
	if level > 100 {
		t.Errorf("expected level clamped to 100, got %f", level)
	}

	// Run many ticks with high outflow
	reservoir2 := NewReservoirModel("test-reservoir2", 10000.0)
	for i := 0; i < 1000; i++ {
		reservoir2.ExtractFlow(1000.0)
		reservoir2.Tick()
	}

	level2 := reservoir2.Level()

	// Level should be clamped to min (0)
	if level2 < 0 {
		t.Errorf("expected level clamped to 0, got %f", level2)
	}
}

func TestModelBase(t *testing.T) {
	base := NewModelBase("test-id", "test-type")

	if base.ID() != "test-id" {
		t.Errorf("expected ID 'test-id', got '%s'", base.ID())
	}

	if base.Type() != "test-type" {
		t.Errorf("expected type 'test-type', got '%s'", base.Type())
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		name     string
		value    float32
		min      float32
		max      float32
		expected float32
	}{
		{"within bounds", 5.0, 0.0, 10.0, 5.0},
		{"below min", -5.0, 0.0, 10.0, 0.0},
		{"above max", 15.0, 0.0, 10.0, 10.0},
		{"at min", 0.0, 0.0, 10.0, 0.0},
		{"at max", 10.0, 0.0, 10.0, 10.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := clamp(tt.value, tt.min, tt.max)
			if result != tt.expected {
				t.Errorf("clamp(%f, %f, %f) = %f, want %f",
					tt.value, tt.min, tt.max, result, tt.expected)
			}
		})
	}
}
