package grid

import (
	"math"
	"testing"
	"time"

	"github.com/tamzrod/forge/internal/models/clock"
)

func TestGrid_Creation(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: 100 * time.Millisecond,
	})

	grid := New(DefaultConfig(), simClock)

	if grid.Voltage() != NominalVoltage {
		t.Errorf("expected nominal voltage %f, got %f", NominalVoltage, grid.Voltage())
	}

	if grid.Frequency() != NominalFrequency {
		t.Errorf("expected nominal frequency %f, got %f", NominalFrequency, grid.Frequency())
	}

	if !grid.IsStable() {
		t.Error("expected grid to be stable at nominal")
	}
}

func TestGrid_ReactivePowerInjectionChangesVoltage(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: 100 * time.Millisecond,
	})

	grid := New(DefaultConfig(), simClock)

	initialVoltage := grid.Voltage()

	// Inject reactive power (generating VARs to the grid)
	grid.InjectReactivePower(100.0) // 100 MVAr injection

	// Tick should cause voltage to rise
	grid.Tick()

	// Voltage should increase from injection
	if grid.Voltage() <= initialVoltage {
		t.Errorf("expected voltage to rise from reactive injection, got %f", grid.Voltage())
	}

	// Test the opposite
	grid.SetVoltage(NominalVoltage)
	grid.InjectReactivePower(-100.0) // 100 MVAr absorption

	grid.Tick()

	if grid.Voltage() >= NominalVoltage {
		t.Errorf("expected voltage to drop from reactive absorption, got %f", grid.Voltage())
	}
}

func TestGrid_ActivePowerInjectionChangesFrequency(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: 100 * time.Millisecond,
	})

	grid := New(DefaultConfig(), simClock)

	initialFrequency := grid.Frequency()

	// Inject active power (generating MW to the grid)
	grid.InjectActivePower(100.0) // 100 MW injection

	// Tick should cause frequency to rise
	grid.Tick()

	// Frequency should increase from injection
	if grid.Frequency() <= initialFrequency {
		t.Errorf("expected frequency to rise from active injection, got %f", grid.Frequency())
	}

	// Test the opposite
	grid.SetFrequency(NominalFrequency)
	grid.InjectActivePower(-100.0) // 100 MW load

	grid.Tick()

	if grid.Frequency() >= NominalFrequency {
		t.Errorf("expected frequency to drop from active load, got %f", grid.Frequency())
	}
}

func TestGrid_PowerBalanceResetsEachTick(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: 100 * time.Millisecond,
	})

	grid := New(DefaultConfig(), simClock)

	// Inject power
	grid.InjectActivePower(50.0)
	grid.InjectReactivePower(30.0)

	// Balance should be recorded
	if grid.ActivePowerBalance() != 50.0 {
		t.Errorf("expected active balance 50.0, got %f", grid.ActivePowerBalance())
	}

	// Tick resets balances
	grid.Tick()

	// After tick, balances should be zero
	if grid.ActivePowerBalance() != 0 {
		t.Errorf("expected active balance 0 after tick, got %f", grid.ActivePowerBalance())
	}

	if grid.ReactivePowerBalance() != 0 {
		t.Errorf("expected reactive balance 0 after tick, got %f", grid.ReactivePowerBalance())
	}
}

func TestGrid_VoltageBounds(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: 100 * time.Millisecond,
	})

	grid := New(DefaultConfig(), simClock)

	// Try to inject massive reactive power
	for i := 0; i < 1000; i++ {
		grid.InjectReactivePower(1000.0)
		grid.Tick()
	}

	// Voltage should be clamped to max
	if grid.Voltage() > MaxVoltage {
		t.Errorf("expected voltage <= %f, got %f", MaxVoltage, grid.Voltage())
	}

	// Try to drain all reactive power
	for i := 0; i < 1000; i++ {
		grid.InjectReactivePower(-1000.0)
		grid.Tick()
	}

	// Voltage should be clamped to min
	if grid.Voltage() < MinVoltage {
		t.Errorf("expected voltage >= %f, got %f", MinVoltage, grid.Voltage())
	}
}

func TestGrid_FrequencyBounds(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: 100 * time.Millisecond,
	})

	grid := New(DefaultConfig(), simClock)

	// Try to inject massive active power
	for i := 0; i < 1000; i++ {
		grid.InjectActivePower(5000.0)
		grid.Tick()
	}

	// Frequency should be clamped to max
	if grid.Frequency() > MaxFrequency {
		t.Errorf("expected frequency <= %f, got %f", MaxFrequency, grid.Frequency())
	}

	// Try to create massive load
	for i := 0; i < 1000; i++ {
		grid.InjectActivePower(-5000.0)
		grid.Tick()
	}

	// Frequency should be clamped to min
	if grid.Frequency() < MinFrequency {
		t.Errorf("expected frequency >= %f, got %f", MinFrequency, grid.Frequency())
	}
}

func TestGrid_IsStable(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: 100 * time.Millisecond,
	})

	grid := New(DefaultConfig(), simClock)

	// At nominal, should be stable
	if !grid.IsStable() {
		t.Error("expected grid to be stable at nominal")
	}

	// Slight deviation should still be stable
	grid.SetVoltage(455.0) // Slightly low but within 0.95 PU
	if !grid.IsStable() {
		t.Error("expected grid to be stable at 455V")
	}

	// Large deviation should be unstable
	grid.SetVoltage(400.0) // Below 0.9 PU
	if grid.IsStable() {
		t.Error("expected grid to be unstable at 400V")
	}
}

func TestGrid_IsUnderVoltage(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: 100 * time.Millisecond,
	})

	grid := New(DefaultConfig(), simClock)

	// Normal voltage
	if grid.IsUnderVoltage() {
		t.Error("expected no undervoltage at nominal")
	}

	// Below 0.9 PU
	grid.SetVoltage(430.0) // 430/480 = 0.896 PU
	if !grid.IsUnderVoltage() {
		t.Error("expected undervoltage at 430V")
	}
}

func TestGrid_IsOverFrequency(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: 100 * time.Millisecond,
	})

	grid := New(DefaultConfig(), simClock)

	// Normal frequency
	if grid.IsOverFrequency() {
		t.Error("expected no overfrequency at nominal")
	}

	// Above 1.05 PU (63 Hz)
	grid.SetFrequency(63.0)
	if !grid.IsOverFrequency() {
		t.Error("expected overfrequency at 63 Hz")
	}
}

func TestGrid_Reset(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: 100 * time.Millisecond,
	})

	grid := New(DefaultConfig(), simClock)

	// Deviate from nominal
	grid.SetVoltage(400.0)
	grid.SetFrequency(59.0)
	grid.InjectActivePower(100.0)
	grid.InjectReactivePower(50.0)

	// Reset
	grid.Reset()

	if grid.Voltage() != NominalVoltage {
		t.Errorf("expected voltage reset to %f, got %f", NominalVoltage, grid.Voltage())
	}

	if grid.Frequency() != NominalFrequency {
		t.Errorf("expected frequency reset to %f, got %f", NominalFrequency, grid.Frequency())
	}

	if grid.ActivePowerBalance() != 0 {
		t.Errorf("expected active balance reset to 0, got %f", grid.ActivePowerBalance())
	}
}

func TestGrid_VoltagePU(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: 100 * time.Millisecond,
	})

	grid := New(DefaultConfig(), simClock)

	// At nominal, PU should be 1.0
	if grid.VoltagePU() != 1.0 {
		t.Errorf("expected voltage PU 1.0 at nominal, got %f", grid.VoltagePU())
	}

	// At 480V
	grid.SetVoltage(480.0)
	if math.Abs(grid.VoltagePU()-1.0) > 0.001 {
		t.Errorf("expected voltage PU 1.0 at 480V, got %f", grid.VoltagePU())
	}

	// At 240V (half)
	grid.SetVoltage(240.0)
	if math.Abs(grid.VoltagePU()-0.5) > 0.001 {
		t.Errorf("expected voltage PU 0.5 at 240V, got %f", grid.VoltagePU())
	}
}

func TestGrid_FrequencyPU(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: 100 * time.Millisecond,
	})

	grid := New(DefaultConfig(), simClock)

	// At nominal, PU should be 1.0
	if grid.FrequencyPU() != 1.0 {
		t.Errorf("expected frequency PU 1.0 at nominal, got %f", grid.FrequencyPU())
	}
}

func TestGrid_PowerFactor(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: 100 * time.Millisecond,
	})

	grid := New(DefaultConfig(), simClock)

	// Pure active power (PF = 1)
	grid.InjectActivePower(100.0)
	pf := grid.PowerFactor()
	if math.Abs(pf-1.0) > 0.001 {
		t.Errorf("expected PF 1.0 for pure active, got %f", pf)
	}

	// Reset and try pure reactive
	grid.Reset()
	grid.InjectReactivePower(100.0)
	pf = grid.PowerFactor()
	if math.Abs(pf-0.0) > 0.001 {
		t.Errorf("expected PF 0.0 for pure reactive, got %f", pf)
	}

	// Equal P and Q (PF = 0.707)
	grid.Reset()
	grid.InjectActivePower(100.0)
	grid.InjectReactivePower(100.0)
	pf = grid.PowerFactor()
	expected := 1.0 / math.Sqrt(2.0)
	if math.Abs(pf-expected) > 0.001 {
		t.Errorf("expected PF %f for 45°, got %f", expected, pf)
	}
}

func TestGrid_CumulativeInjection(t *testing.T) {
	simClock := clock.New(clock.Config{
		Mode:        clock.ModeManual,
		TickInterval: 100 * time.Millisecond,
	})

	grid := New(DefaultConfig(), simClock)

	// Multiple injections in same tick
	grid.InjectActivePower(10.0)
	grid.InjectActivePower(20.0)
	grid.InjectActivePower(30.0)

	// Should accumulate
	if grid.ActivePowerBalance() != 60.0 {
		t.Errorf("expected balance 60.0, got %f", grid.ActivePowerBalance())
	}

	// Tick resets
	grid.Tick()

	if grid.ActivePowerBalance() != 0.0 {
		t.Errorf("expected balance 0 after tick, got %f", grid.ActivePowerBalance())
	}
}
