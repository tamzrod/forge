package weatherstation

import (
	"testing"
	"time"

	"github.com/tamzrod/forge/internal/devices"
	"github.com/tamzrod/forge/internal/models/clock"
	"github.com/tamzrod/forge/internal/models/weather"
)

func testContext() (*devices.Context, *weather.Weather) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: time.Hour,
	})

	weatherModel := weather.New(weather.DefaultConfig(), simClock)

	return devices.NewContext(simClock, nil, weatherModel, nil), weatherModel
}

func TestNewStation(t *testing.T) {
	ctx, _ := testContext()
	cfg := DefaultConfig()
	cfg.Publishing.Enabled = false // Disable interface for tests

	station, err := NewStation(cfg, ctx)
	if err != nil {
		t.Fatalf("failed to create station: %v", err)
	}

	if station.ID() != cfg.ID {
		t.Errorf("expected ID %s, got %s", cfg.ID, station.ID())
	}

	if station.Type() != Type {
		t.Errorf("expected type %s, got %s", Type, station.Type())
	}

	if station.Name() != cfg.Name {
		t.Errorf("expected name %s, got %s", cfg.Name, station.Name())
	}
}

func TestNewStation_InvalidID(t *testing.T) {
	ctx, _ := testContext()
	cfg := Config{
		ID:   "",
		Name: "Test Station",
	}

	if _, err := NewStation(cfg, ctx); err == nil {
		t.Error("expected error for empty ID")
	}
}

func TestStation_Initialize(t *testing.T) {
	ctx, _ := testContext()
	cfg := DefaultConfig()
	cfg.Publishing.Enabled = false

	station, _ := NewStation(cfg, ctx)

	if err := station.Initialize(); err != nil {
		t.Fatalf("failed to initialize: %v", err)
	}

	if station.Status() != devices.StateInitialized {
		t.Errorf("expected state Initialized, got %s", station.Status())
	}

	// Check initial memory values
	if station.Temperature() != 20.0 {
		t.Errorf("expected initial temperature 20.0, got %f", station.Temperature())
	}
}

func TestStation_Tick(t *testing.T) {
	ctx, _ := testContext()
	cfg := DefaultConfig()
	cfg.Publishing.Enabled = false

	station, _ := NewStation(cfg, ctx)

	station.Initialize()

	// Tick and check observation
	station.Tick()

	if station.Status() != devices.StateRunning {
		t.Errorf("expected state Running, got %s", station.Status())
	}

	// Weather Station should have copied weather values
	temp := station.Temperature()
	if temp == 0 {
		t.Error("expected temperature to be observed")
	}
}

func TestStation_Shutdown(t *testing.T) {
	ctx, _ := testContext()
	cfg := DefaultConfig()
	cfg.Publishing.Enabled = false

	station, _ := NewStation(cfg, ctx)

	station.Initialize()
	station.Tick()

	if err := station.Shutdown(); err != nil {
		t.Fatalf("failed to shutdown: %v", err)
	}

	if station.Status() != devices.StateStopped {
		t.Errorf("expected state Stopped, got %s", station.Status())
	}
}

func TestStation_Memory(t *testing.T) {
	ctx, _ := testContext()
	cfg := DefaultConfig()
	cfg.Publishing.Enabled = false

	station, _ := NewStation(cfg, ctx)

	mem := station.Memory()
	if mem == nil {
		t.Error("expected non-nil memory")
	}
}

func TestStation_TemperatureConversion(t *testing.T) {
	ctx, _ := testContext()

	// Create station with Fahrenheit units
	cfg := Config{
		ID:   "station-001",
		Name: "Test Station",
		Units: Fahrenheit,
		Publishing: PublishingConfig{
			Enabled: false,
		},
	}

	station, _ := NewStation(cfg, ctx)
	station.Initialize()

	// Set weather to 20°C and observe
	// Fahrenheit should be 68°F
	station.Tick()

	// Temperature should be converted
	temp := station.Temperature()
	if temp < 67 || temp > 69 {
		t.Errorf("expected ~68°F, got %f", temp)
	}
}

func TestStation_TickCount(t *testing.T) {
	ctx, _ := testContext()
	cfg := DefaultConfig()
	cfg.Publishing.Enabled = false

	station, _ := NewStation(cfg, ctx)

	if station.TickCount() != 0 {
		t.Error("expected initial tick count 0")
	}

	station.Initialize()
	station.Tick()

	if station.TickCount() != 1 {
		t.Errorf("expected tick count 1, got %d", station.TickCount())
	}

	station.Tick()
	if station.TickCount() != 2 {
		t.Errorf("expected tick count 2, got %d", station.TickCount())
	}
}

func TestStation_State(t *testing.T) {
	ctx, _ := testContext()
	cfg := DefaultConfig()
	cfg.Publishing.Enabled = false

	station, _ := NewStation(cfg, ctx)

	station.Initialize()
	station.Tick()

	state := station.Status()

	if state.ID != "weather-station-001" {
		t.Errorf("expected ID weather-station-001, got %s", state.ID)
	}

	if state.Type != "weather_station" {
		t.Errorf("expected type weather_station, got %s", state.Type)
	}

	if state.DeviceState != devices.StateRunning {
		t.Errorf("expected device state Running, got %s", state.DeviceState)
	}

	if state.TickCount != 1 {
		t.Errorf("expected tick count 1, got %d", state.TickCount)
	}
}

func TestStation_PublishingState(t *testing.T) {
	ctx, _ := testContext()
	cfg := DefaultConfig()
	cfg.Publishing.Enabled = false

	station, _ := NewStation(cfg, ctx)

	pubState := station.PublishingState()
	if pubState.Enabled {
		t.Error("expected publishing to be disabled")
	}
}

func TestStation_PublishWithoutInterface(t *testing.T) {
	ctx, _ := testContext()
	cfg := DefaultConfig()
	cfg.Publishing.Enabled = false

	station, _ := NewStation(cfg, ctx)
	station.Initialize()

	// Should not panic even without interface
	station.Tick()

	if station.Status() != devices.StateRunning {
		t.Errorf("expected state Running, got %s", station.Status())
	}
}
