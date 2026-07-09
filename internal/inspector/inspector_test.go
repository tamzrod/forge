package inspector

import (
	"testing"
	"time"

	"github.com/tamzrod/forge/internal/models/clock"
	"github.com/tamzrod/forge/internal/models/grid"
	"github.com/tamzrod/forge/internal/models/sun"
	"github.com/tamzrod/forge/internal/models/weather"
)

func TestView_Creation(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: time.Hour,
	})

	sunModel := sun.New(sun.Config{
		Latitude:  40.0,
		Longitude: -105.0,
	}, simClock)

	weatherModel := weather.New(weather.DefaultConfig(), simClock)

	gridModel := grid.New(grid.DefaultConfig(), simClock)

	view := NewView(simClock, sunModel, weatherModel, gridModel)

	if view == nil {
		t.Error("expected non-nil view")
	}
}

func TestView_ClockState(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: time.Hour,
	})

	sunModel := sun.New(sun.Config{Latitude: 40.0, Longitude: -105.0}, simClock)
	weatherModel := weather.New(weather.DefaultConfig(), simClock)
	gridModel := grid.New(grid.DefaultConfig(), simClock)

	view := NewView(simClock, sunModel, weatherModel, gridModel)

	// Advance time
	simClock.Advance(5 * time.Hour)
	simClock.Tick()

	state := view.ClockState()

	if state.Elapsed != 5*time.Hour {
		t.Errorf("expected elapsed 5h, got %v", state.Elapsed)
	}

	if state.TickCount != 1 {
		t.Errorf("expected tick count 1, got %d", state.TickCount)
	}

	if state.Mode != "Manual" {
		t.Errorf("expected mode Manual, got %s", state.Mode)
	}
}

func TestView_SunState(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: time.Hour,
	})

	sunModel := sun.New(sun.Config{
		Latitude:  40.0,
		Longitude: -105.0,
	}, simClock)

	weatherModel := weather.New(weather.DefaultConfig(), simClock)
	gridModel := grid.New(grid.DefaultConfig(), simClock)

	view := NewView(simClock, sunModel, weatherModel, gridModel)

	// Set to noon
	simClock.Advance(80 * 24 * time.Hour) // March 21
	simClock.Advance(12 * time.Hour)       // noon
	sunModel.Tick()

	state := view.SunState()

	if state.Latitude != 40.0 {
		t.Errorf("expected latitude 40.0, got %f", state.Latitude)
	}

	if state.Longitude != -105.0 {
		t.Errorf("expected longitude -105.0, got %f", state.Longitude)
	}

	if state.Elevation < 0 {
		t.Errorf("expected positive elevation at noon, got %f", state.Elevation)
	}

	if !state.IsDaytime {
		t.Error("expected daytime at noon")
	}
}

func TestView_WeatherState(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: time.Hour,
	})

	sunModel := sun.New(sun.Config{Latitude: 40.0, Longitude: -105.0}, simClock)
	weatherModel := weather.New(weather.DefaultConfig(), simClock)
	gridModel := grid.New(grid.DefaultConfig(), simClock)

	view := NewView(simClock, sunModel, weatherModel, gridModel)

	weatherModel.Tick()

	state := view.WeatherState()

	if state.Temperature < -50 || state.Temperature > 60 {
		t.Errorf("temperature out of reasonable range: %f", state.Temperature)
	}

	if state.Humidity < 0 || state.Humidity > 100 {
		t.Errorf("humidity out of range: %f", state.Humidity)
	}

	if state.CloudCover < 0 || state.CloudCover > 1 {
		t.Errorf("cloud cover out of range: %f", state.CloudCover)
	}

	if state.WindSpeed < 0 {
		t.Errorf("wind speed should not be negative: %f", state.WindSpeed)
	}
}

func TestView_GridState(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: time.Hour,
	})

	sunModel := sun.New(sun.Config{Latitude: 40.0, Longitude: -105.0}, simClock)
	weatherModel := weather.New(weather.DefaultConfig(), simClock)
	gridModel := grid.New(grid.DefaultConfig(), simClock)

	view := NewView(simClock, sunModel, weatherModel, gridModel)

	state := view.GridState()

	if state.Voltage != 480.0 {
		t.Errorf("expected nominal voltage 480.0, got %f", state.Voltage)
	}

	if state.Frequency != 60.0 {
		t.Errorf("expected nominal frequency 60.0, got %f", state.Frequency)
	}

	if !state.IsStable {
		t.Error("expected grid to be stable at nominal")
	}
}

func TestView_FullState(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: time.Hour,
	})

	sunModel := sun.New(sun.Config{Latitude: 40.0, Longitude: -105.0}, simClock)
	weatherModel := weather.New(weather.DefaultConfig(), simClock)
	gridModel := grid.New(grid.DefaultConfig(), simClock)

	view := NewView(simClock, sunModel, weatherModel, gridModel)

	state := view.FullState()

	// All sub-states should be populated
	if state.Clock.Elapsed != 0 {
		t.Error("clock should start at 0")
	}

	if state.Sun.Latitude != 40.0 {
		t.Error("sun state not populated")
	}

	if state.Weather.Temperature == 0 && state.Weather.Humidity == 0 {
		// Temperature might be 0 in some configs
		t.Log("weather state appears empty")
	}

	if state.Grid.Voltage != 480.0 {
		t.Error("grid state not populated")
	}
}

func TestView_InjectPower(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: 100 * time.Millisecond,
	})

	sunModel := sun.New(sun.Config{Latitude: 40.0, Longitude: -105.0}, simClock)
	weatherModel := weather.New(weather.DefaultConfig(), simClock)
	gridModel := grid.New(grid.DefaultConfig(), simClock)

	view := NewView(simClock, sunModel, weatherModel, gridModel)

	// Inject power
	gridModel.InjectActivePower(100.0)
	gridModel.InjectReactivePower(50.0)
	gridModel.Tick()

	state := view.GridState()

	// After tick, balance should be reset
	if state.ActiveBalance != 0 {
		t.Errorf("expected active balance 0 after tick, got %f", state.ActiveBalance)
	}

	if state.ReactiveBalance != 0 {
		t.Errorf("expected reactive balance 0 after tick, got %f", state.ReactiveBalance)
	}

	// But frequency and voltage should have changed
	if state.Frequency <= 60.0 {
		t.Errorf("expected frequency > 60 after injection, got %f", state.Frequency)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Address != "localhost:8080" {
		t.Errorf("expected default address localhost:8080, got %s", cfg.Address)
	}

	if cfg.UpdateInterval != 100 {
		t.Errorf("expected default update interval 100ms, got %d", cfg.UpdateInterval)
	}
}
