// Package grid provides a simple electrical grid simulation model.
//
// The grid model simulates electrical grid conditions:
// - Voltage
// - Frequency
// - Power balance
//
// This is a simplified model suitable for industrial software testing.
// It does NOT implement power flow calculations or transient stability.
// The goal is believable behavior, not grid analysis accuracy.
//
// Key behaviors:
// - Voltage responds to reactive power imbalance
// - Frequency responds to active power imbalance
// - Both are self-regulating within bounds
package grid

import (
	"math"

	"github.com/tamzrod/forge/internal/models/clock"
)

const (
	// Nominal voltage in volts
	NominalVoltage = 480.0

	// Nominal frequency in Hz
	NominalFrequency = 60.0

	// Voltage bounds
	MinVoltage = 450.0
	MaxVoltage = 520.0

	// Frequency bounds
	MinFrequency = 59.5
	MaxFrequency = 60.5

	// Default grid strength (higher = more stable)
	DefaultStrength = 1000.0 // MVA short circuit

	// Voltage sensitivity (PU voltage change per PU reactive power)
	DefaultVoltageSensitivity = 0.05

	// Frequency sensitivity (Hz change per MW imbalance)
	DefaultFrequencySensitivity = 0.001
)

// Grid represents an electrical grid model.
type Grid struct {
	clock *clock.Clock

	// Grid parameters
	strength          float64 // Short circuit MVA
	voltageSensitivity float64 // PU voltage per PU reactive power
	frequencySensitivity float64 // Hz per MW imbalance

	// State
	voltage          float64 // Volts
	frequency        float64 // Hz
	activePowerBalance float64 // MW (positive = generation > load)
	reactivePowerBalance float64 // MVAr (positive = generation > load)

	// Nominal values
	nominalVoltage   float64
	nominalFrequency float64
}

// Config holds grid model configuration.
type Config struct {
	// Nominal voltage in volts
	NominalVoltage float64

	// Nominal frequency in Hz
	NominalFrequency float64

	// Short circuit strength in MVA
	// Higher = stronger grid, less voltage variation
	Strength float64

	// Voltage sensitivity (PU voltage change per PU reactive power)
	// Affects how much voltage changes with reactive power imbalance
	VoltageSensitivity float64

	// Frequency sensitivity (Hz per MW imbalance)
	// Affects how much frequency changes with active power imbalance
	FrequencySensitivity float64
}

// DefaultConfig returns a reasonable default configuration.
func DefaultConfig() Config {
	return Config{
		NominalVoltage:       NominalVoltage,
		NominalFrequency:     NominalFrequency,
		Strength:             DefaultStrength,
		VoltageSensitivity:   DefaultVoltageSensitivity,
		FrequencySensitivity: DefaultFrequencySensitivity,
	}
}

// New creates a new grid model.
func New(cfg Config, simClock *clock.Clock) *Grid {
	return &Grid{
		clock:               simClock,
		strength:            cfg.Strength,
		voltageSensitivity:   cfg.VoltageSensitivity,
		frequencySensitivity: cfg.FrequencySensitivity,
		voltage:             cfg.NominalVoltage,
		frequency:           cfg.NominalFrequency,
		nominalVoltage:      cfg.NominalVoltage,
		nominalFrequency:   cfg.NominalFrequency,
	}
}

// Voltage returns the grid voltage in volts.
func (g *Grid) Voltage() float64 {
	return g.voltage
}

// Frequency returns the grid frequency in Hz.
func (g *Grid) Frequency() float64 {
	return g.frequency
}

// ActivePowerBalance returns the active power imbalance in MW.
// Positive = generation exceeds load
// Negative = load exceeds generation
func (g *Grid) ActivePowerBalance() float64 {
	return g.activePowerBalance
}

// ReactivePowerBalance returns the reactive power imbalance in MVAr.
// Positive = generation exceeds load
// Negative = load exceeds generation
func (g *Grid) ReactivePowerBalance() float64 {
	return g.reactivePowerBalance
}

// Strength returns the short circuit strength in MVA.
func (g *Grid) Strength() float64 {
	return g.strength
}

// InjectActivePower records active power injection.
// Positive = device is injecting (generating) MW to the grid
// Negative = device is consuming (loading) MW from the grid
func (g *Grid) InjectActivePower(mw float64) {
	g.activePowerBalance += mw
}

// InjectReactivePower records reactive power injection.
// Positive = device is injecting (generating) MVAr to the grid
// Negative = device is consuming (absorbing) MVAr from the grid
func (g *Grid) InjectReactivePower(mvar float64) {
	g.reactivePowerBalance += mvar
}

// VoltagePU returns voltage in per-unit (relative to nominal).
func (g *Grid) VoltagePU() float64 {
	return g.voltage / g.nominalVoltage
}

// FrequencyPU returns frequency in per-unit (relative to nominal).
func (g *Grid) FrequencyPU() float64 {
	return g.frequency / g.nominalFrequency
}

// Tick advances the grid model by one simulation step.
// Grid self-regulates based on power imbalances.
func (g *Grid) Tick() {
	g.regulate()
}

// regulate adjusts voltage and frequency based on power balance.
// This is a simplified droop characteristic model.
func (g *Grid) regulate() {
	// Reactive power affects voltage (simplified)
	// Q imbalance causes voltage to change
	// Grid strength determines how much
	powerFactor := g.strength / (g.nominalVoltage * g.nominalVoltage / 1000.0) // Normalized strength
	qImbalancePU := g.reactivePowerBalance / (g.strength / 1000.0) // PU reactive power
	vDelta := -g.voltageSensitivity * qImbalancePU * powerFactor
	g.voltage += vDelta * g.nominalVoltage

	// Active power affects frequency (simplified droop)
	// P imbalance causes frequency to change
	pImbalancePU := g.activePowerBalance / 1000.0 // Convert MW to GW-like units
	fDelta := -g.frequencySensitivity * pImbalancePU
	g.frequency += fDelta

	// Clamp to operating bounds
	g.voltage = clamp(g.voltage, MinVoltage, MaxVoltage)
	g.frequency = clamp(g.frequency, MinFrequency, MaxFrequency)

	// Reset power balances for next tick
	g.activePowerBalance = 0
	g.reactivePowerBalance = 0
}

// IsStable returns true if grid is within normal operating bounds.
func (g *Grid) IsStable() bool {
	puV := g.VoltagePU()
	puF := g.FrequencyPU()

	// Normal operating range: 0.95 - 1.05 PU
	return puV >= 0.95 && puV <= 1.05 && puF >= 0.95 && puF <= 1.05
}

// IsUnderVoltage returns true if voltage is below 0.9 PU.
func (g *Grid) IsUnderVoltage() bool {
	return g.VoltagePU() < 0.9
}

// IsOverFrequency returns true if frequency is above 1.05 PU.
func (g *Grid) IsOverFrequency() bool {
	return g.FrequencyPU() > 1.05
}

// SetVoltage directly sets the grid voltage.
// This should only be used for testing or scenario injection.
func (g *Grid) SetVoltage(v float64) {
	g.voltage = clamp(v, MinVoltage, MaxVoltage)
}

// SetFrequency directly sets the grid frequency.
// This should only be used for testing or scenario injection.
func (g *Grid) SetFrequency(f float64) {
	g.frequency = clamp(f, MinFrequency, MaxFrequency)
}

// Reset restores the grid to nominal conditions.
func (g *Grid) Reset() {
	g.voltage = g.nominalVoltage
	g.frequency = g.nominalFrequency
	g.activePowerBalance = 0
	g.reactivePowerBalance = 0
}

// NominalVoltage returns the nominal grid voltage.
func (g *Grid) NominalVoltage() float64 {
	return g.nominalVoltage
}

// NominalFrequency returns the nominal grid frequency.
func (g *Grid) NominalFrequency() float64 {
	return g.nominalFrequency
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

// PowerFactor returns the apparent power factor based on power balances.
func (g *Grid) PowerFactor() float64 {
	apparent := math.Sqrt(g.activePowerBalance*g.activePowerBalance +
		g.reactivePowerBalance*g.reactivePowerBalance)
	if apparent == 0 {
		return 1.0
	}
	return g.activePowerBalance / apparent
}
