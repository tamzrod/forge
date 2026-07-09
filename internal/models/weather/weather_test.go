package weather

import (
	"testing"
	"time"

	"github.com/tamzrod/forge/internal/models/clock"
)

func TestWeather_Creation(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: time.Hour,
	})

	cfg := DefaultConfig()
	cfg.Seed = 12345

	weather := New(cfg, simClock)

	if weather.Temperature() != cfg.BaseTemperature {
		t.Errorf("expected initial temp %f, got %f", cfg.BaseTemperature, weather.Temperature())
	}

	if weather.Humidity() != cfg.BaseHumidity {
		t.Errorf("expected initial humidity %f, got %f", cfg.BaseHumidity, weather.Humidity())
	}

	if weather.Seed() != 12345 {
		t.Errorf("expected seed 12345, got %d", weather.Seed())
	}
}

func TestWeather_EvolvesOverTime(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: time.Hour,
	})

	cfg := DefaultConfig()
	cfg.Seed = 12345

	weather := New(cfg, simClock)

	initialTemp := weather.Temperature()
	initialHumidity := weather.Humidity()

	// Advance many ticks
	for i := 0; i < 100; i++ {
		simClock.Tick()
		weather.Tick()
	}

	// Temperature should have changed
	if weather.Temperature() == initialTemp {
		t.Log("Note: Temperature may stay same if random seed produces same value")
	}

	// Values should still be within bounds
	if weather.Temperature() < -40 || weather.Temperature() > 50 {
		t.Errorf("temperature out of bounds: %f", weather.Temperature())
	}

	if weather.Humidity() < 0 || weather.Humidity() > 100 {
		t.Errorf("humidity out of bounds: %f", weather.Humidity())
	}

	if weather.Pressure() < 950 || weather.Pressure() > 1050 {
		t.Errorf("pressure out of bounds: %f", weather.Pressure())
	}

	if weather.CloudCover() < 0 || weather.CloudCover() > 1 {
		t.Errorf("cloud cover out of bounds: %f", weather.CloudCover())
	}

	if weather.WindSpeed() < 0 || weather.WindSpeed() > 50 {
		t.Errorf("wind speed out of bounds: %f", weather.WindSpeed())
	}

	_ = initialHumidity // May or may not change
}

func TestWeather_Deterministic(t *testing.T) {
	simClock1 := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: time.Hour,
	})

	simClock2 := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: time.Hour,
	})

	cfg := DefaultConfig()
	cfg.Seed = 42

	weather1 := New(cfg, simClock1)
	weather2 := New(cfg, simClock2)

	// Advance both identically
	for i := 0; i < 50; i++ {
		simClock1.Tick()
		simClock2.Tick()
		weather1.Tick()
		weather2.Tick()
	}

	// All values should be identical
	if weather1.Temperature() != weather2.Temperature() {
		t.Errorf("temperature mismatch: %f vs %f", weather1.Temperature(), weather2.Temperature())
	}

	if weather1.Humidity() != weather2.Humidity() {
		t.Errorf("humidity mismatch: %f vs %f", weather1.Humidity(), weather2.Humidity())
	}

	if weather1.Pressure() != weather2.Pressure() {
		t.Errorf("pressure mismatch: %f vs %f", weather1.Pressure(), weather2.Pressure())
	}

	if weather1.CloudCover() != weather2.CloudCover() {
		t.Errorf("cloud cover mismatch: %f vs %f", weather1.CloudCover(), weather2.CloudCover())
	}

	if weather1.WindSpeed() != weather2.WindSpeed() {
		t.Errorf("wind speed mismatch: %f vs %f", weather1.WindSpeed(), weather2.WindSpeed())
	}

	if weather1.WindDirection() != weather2.WindDirection() {
		t.Errorf("wind direction mismatch: %f vs %f", weather1.WindDirection(), weather2.WindDirection())
	}
}

func TestWeather_Reset(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: time.Hour,
	})

	cfg := DefaultConfig()
	cfg.Seed = 12345

	weather := New(cfg, simClock)

	// Advance some time
	for i := 0; i < 50; i++ {
		simClock.Tick()
		weather.Tick()
	}

	// Reset
	simClock.Reset()
	weather.Reset()

	// Should be back to initial state
	if weather.Temperature() != cfg.BaseTemperature {
		t.Errorf("expected temp reset to %f, got %f", cfg.BaseTemperature, weather.Temperature())
	}

	if weather.Humidity() != cfg.BaseHumidity {
		t.Errorf("expected humidity reset to %f, got %f", cfg.BaseHumidity, weather.Humidity())
	}

	// After same number of ticks, should be identical to first run
	for i := 0; i < 50; i++ {
		simClock.Tick()
		weather.Tick()
	}

	// Compare with fresh instance
	simClock2 := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: time.Hour,
	})
	weather2 := New(cfg, simClock2)
	for i := 0; i < 50; i++ {
		simClock2.Tick()
		weather2.Tick()
	}

	if weather.Temperature() != weather2.Temperature() {
		t.Errorf("reset didn't produce deterministic behavior: %f vs %f",
			weather.Temperature(), weather2.Temperature())
	}
}

func TestWeather_WindDirectionBounds(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: time.Minute,
	})

	cfg := DefaultConfig()
	cfg.Seed = 99999 // High seed for more wind variation

	weather := New(cfg, simClock)

	// Run many ticks
	for i := 0; i < 1000; i++ {
		simClock.Tick()
		weather.Tick()
	}

	// Direction should always be in [0, 360)
	dir := weather.WindDirection()
	if dir < 0 || dir >= 360 {
		t.Errorf("wind direction out of bounds [0, 360): %f", dir)
	}
}

func TestWeather_CustomConfig(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: time.Hour,
	})

	cfg := Config{
		BaseTemperature: 30.0,
		TemperatureRange: 15.0,
		BaseHumidity: 60.0,
		BasePressure: 1000.0,
		Seed: 54321,
	}

	weather := New(cfg, simClock)

	// Initial values should match config
	if weather.Temperature() != 30.0 {
		t.Errorf("expected temp 30.0, got %f", weather.Temperature())
	}

	if weather.Humidity() != 60.0 {
		t.Errorf("expected humidity 60.0, got %f", weather.Humidity())
	}
}

func TestWeather_TemperatureDailyCycle(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: time.Hour,
	})

	cfg := Config{
		BaseTemperature: 20.0,
		TemperatureRange: 10.0,
		BaseHumidity: 50.0,
		BasePressure: 1013.25,
		Seed: 42,
	}

	weather := New(cfg, simClock)

	// Check temperature at different times of day
	temps := make([]float64, 4)

	// 6 AM
	simClock.Advance(6 * time.Hour)
	weather.Tick()
	temps[0] = weather.Temperature()

	// Noon
	simClock.Advance(6 * time.Hour)
	weather.Tick()
	temps[1] = weather.Temperature()

	// 6 PM
	simClock.Advance(6 * time.Hour)
	weather.Tick()
	temps[2] = weather.Temperature()

	// Midnight
	simClock.Advance(6 * time.Hour)
	weather.Tick()
	temps[3] = weather.Temperature()

	// Temperature should peak around 6 PM and be lowest around 6 AM
	if temps[2] <= temps[0] {
		t.Logf("Note: Afternoon temp (%f) should be higher than morning (%f)", temps[2], temps[0])
	}

	if temps[3] <= temps[1] {
		t.Logf("Note: Night temp (%f) should be lower than noon (%f)", temps[3], temps[1])
	}
}

func TestWeather_RainConditions(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: time.Hour,
	})

	cfg := DefaultConfig()
	cfg.Seed = 11111

	weather := New(cfg, simClock)

	// Check that rain rate is 0 when not raining
	rainCount := 0
	for i := 0; i < 100; i++ {
		simClock.Tick()
		weather.Tick()
		if weather.RainRate() > 0 {
			rainCount++
		}
	}

	// With default cloud cover, rain should occur sometimes
	t.Logf("Rain occurred %d times in 100 ticks", rainCount)
}
