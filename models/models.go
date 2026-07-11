// Package models provides Simulation Models that represent the physical world.
// Simulation Models own physical state, physical behavior, and mathematical models.
// They do NOT own protocols, device identity, MMA2 publishing, or industrial communications.
package models

import (
	"math"
	"math/rand"
	"time"
)

// ModelID is a unique identifier for a simulation model.
type ModelID string

// Model is the interface for all simulation models.
// Models represent the physical world and are observed by virtual devices.
type Model interface {
	// ID returns the model identifier.
	ID() ModelID

	// Type returns the model type name.
	Type() string

	// Tick advances the model state by one simulation step.
	// Models evolve their internal state based on physics.
	Tick()
}

// ModelBase provides common functionality for all models.
type ModelBase struct {
	id       ModelID
	typeName string
}

// ID returns the model identifier.
func (m *ModelBase) ID() ModelID {
	return m.id
}

// Type returns the model type name.
func (m *ModelBase) Type() string {
	return m.typeName
}

// NewModelBase creates a new model base.
func NewModelBase(id ModelID, typeName string) ModelBase {
	return ModelBase{
		id:       id,
		typeName: typeName,
	}
}

// GridModel represents an electrical grid.
// It models voltage, frequency, Thevenin impedance, and power flow.
type GridModel struct {
	ModelBase

	// Physical state (stored directly in RAM)
	voltage         float32 // Volts
	frequency       float32 // Hz
	theveninZReal   float32 // Ohms (real component)
	theveninZImag   float32 // Ohms (imaginary component)
	shortCircuitMVA float32 // MVA
	reactiveSens    float32 // PU MVAr per PU voltage

	// Power injection from devices
	activePowerInjection   float32 // MW
	reactivePowerInjection float32 // MVAr
}

// NewGridModel creates a new Grid model.
func NewGridModel(id ModelID) *GridModel {
	return &GridModel{
		ModelBase: NewModelBase(id, "grid"),
		// Default grid conditions
		voltage:         480.0,
		frequency:       60.0,
		theveninZReal:   0.1,
		theveninZImag:   0.1,
		shortCircuitMVA: 1000.0,
		reactiveSens:    0.2,
	}
}

// Voltage returns the grid voltage in volts.
func (g *GridModel) Voltage() float32 {
	return g.voltage
}

// Frequency returns the grid frequency in Hz.
func (g *GridModel) Frequency() float32 {
	return g.frequency
}

// TheveninImpedance returns the Thevenin impedance components.
func (g *GridModel) TheveninImpedance() (real, imag float32) {
	return g.theveninZReal, g.theveninZImag
}

// ShortCircuitMVA returns the short circuit level in MVA.
func (g *GridModel) ShortCircuitMVA() float32 {
	return g.shortCircuitMVA
}

// InjectActivePower records active power injection from a device.
func (g *GridModel) InjectActivePower(mw float32) {
	g.activePowerInjection += mw
}

// InjectReactivePower records reactive power injection from a device.
func (g *GridModel) InjectReactivePower(mvar float32) {
	g.reactivePowerInjection += mvar
}

// Tick advances the grid model state.
func (g *GridModel) Tick() {
	// Reset power injections for this tick
	g.activePowerInjection = 0
	g.reactivePowerInjection = 0

	// Grid voltage regulation (simplified model)
	// Voltage responds to reactive power imbalance
	voltageDelta := -g.reactivePowerInjection * g.reactiveSens * 0.01
	g.voltage += voltageDelta

	// Keep voltage within bounds
	g.voltage = clamp(g.voltage, 450.0, 520.0)

	// Frequency regulation (simplified model)
	// Frequency responds to active power imbalance
	frequencyDelta := -g.activePowerInjection * 0.001
	g.frequency += frequencyDelta

	// Keep frequency within bounds
	g.frequency = clamp(g.frequency, 59.5, 60.5)
}

// SetVoltage sets the grid voltage directly (for testing/scenarios).
func (g *GridModel) SetVoltage(v float32) {
	g.voltage = clamp(v, 0, 10000)
}

// SetFrequency sets the grid frequency directly (for testing/scenarios).
func (g *GridModel) SetFrequency(f float32) {
	g.frequency = clamp(f, 0, 100)
}

// SunModel represents the sun.
// It models irradiance based on position and time of day.
type SunModel struct {
	ModelBase

	// Position
	azimuth   float32 // degrees from north
	elevation float32 // degrees above horizon

	// Irradiance
	irradiance float32 // W/m² (global horizontal irradiance)
	directNI   float32 // W/m² (direct normal irradiance)
	diffuseHI float32 // W/m² (diffuse horizontal irradiance)

	// Solar time
	solarHourAngle float32 // degrees
}

// NewSunModel creates a new Sun model.
func NewSunModel(id ModelID) *SunModel {
	return &SunModel{
		ModelBase: NewModelBase(id, "sun"),
		irradiance: 1000.0,
		directNI:   900.0,
		diffuseHI:  100.0,
		elevation:  45.0,
		azimuth:    180.0,
	}
}

// Irradiance returns the global horizontal irradiance in W/m².
func (s *SunModel) Irradiance() float32 {
	return s.irradiance
}

// DirectNormalIrradiance returns the direct normal irradiance in W/m².
func (s *SunModel) DirectNormalIrradiance() float32 {
	return s.directNI
}

// DiffuseHorizontalIrradiance returns the diffuse horizontal irradiance in W/m².
func (s *SunModel) DiffuseHorizontalIrradiance() float32 {
	return s.diffuseHI
}

// Elevation returns the sun elevation angle in degrees.
func (s *SunModel) Elevation() float32 {
	return s.elevation
}

// Azimuth returns the sun azimuth angle in degrees from north.
func (s *SunModel) Azimuth() float32 {
	return s.azimuth
}

// Tick advances the sun model state based on simulation time.
func (s *SunModel) Tick() {
	// Sun moves based on simulation time (simplified)
	// In a real implementation, this would use the scheduler clock
	s.solarHourAngle += 0.25 // degrees per tick (simplified)
	if s.solarHourAngle > 180 {
		s.solarHourAngle = -180
	}

	// Calculate elevation from hour angle (simplified)
	s.elevation = float32(45.0 * math.Sin(float64(s.solarHourAngle)*math.Pi/180.0))
	if s.elevation < 0 {
		s.elevation = 0
		s.irradiance = 0
		s.directNI = 0
		s.diffuseHI = 0
		return
	}

	// Calculate irradiance based on elevation (simplified)
	clearness := 0.7 + 0.3*math.Sin(float64(s.elevation)*math.Pi/180.0)
	s.irradiance = 1000.0 * float32(clearness)
	s.directNI = s.irradiance * 0.9 * float32(math.Sin(float64(s.elevation)*math.Pi/180.0))
	s.diffuseHI = s.irradiance * 0.1
}

// WindModel represents wind conditions.
type WindModel struct {
	ModelBase

	speed     float32 // m/s
	direction float32 // degrees from north
	gusts     float32 // m/s
	turbulence float32 // dimensionless intensity
}

// NewWindModel creates a new Wind model.
func NewWindModel(id ModelID) *WindModel {
	return &WindModel{
		ModelBase:  NewModelBase(id, "wind"),
		speed:      5.0,
		direction:  0.0,
		gusts:      0.0,
		turbulence: 0.1,
	}
}

// Speed returns the wind speed in m/s.
func (w *WindModel) Speed() float32 {
	return w.speed
}

// Direction returns the wind direction in degrees from north.
func (w *WindModel) Direction() float32 {
	return w.direction
}

// Gusts returns the gust speed in m/s.
func (w *WindModel) Gusts() float32 {
	return w.gusts
}

// TurbulenceIntensity returns the turbulence intensity.
func (w *WindModel) TurbulenceIntensity() float32 {
	return w.turbulence
}

// Tick advances the wind model state.
func (w *WindModel) Tick() {
	// Simplified wind model with random variations
	// In production, this would use a more sophisticated model
	delta := (float32(random()) - 0.5) * 0.5
	w.speed += delta
	w.speed = clamp(w.speed, 0, 50)

	// Direction slowly drifts
	dirDelta := (float32(random()) - 0.5) * 2.0
	w.direction += dirDelta
	if w.direction < 0 {
		w.direction += 360
	}
	if w.direction >= 360 {
		w.direction -= 360
	}

	// Gusts are intermittent
	if float32(random()) > 0.9 {
		w.gusts = w.speed * (0.2 + 0.1*float32(random()))
	} else {
		w.gusts *= 0.9
	}
}

// WeatherModel represents ambient weather conditions.
type WeatherModel struct {
	ModelBase

	temperature float32 // °C
	humidity    float32 // % relative humidity
	pressure    float32 // hPa
	cloudCover  float32 // 0-1 fraction
	rainRate    float32 // mm/h
}

// NewWeatherModel creates a new Weather model.
func NewWeatherModel(id ModelID) *WeatherModel {
	return &WeatherModel{
		ModelBase: NewModelBase(id, "weather"),
		temperature: 25.0,
		humidity:    50.0,
		pressure:    1013.25,
		cloudCover:  0.3,
		rainRate:    0.0,
	}
}

// Temperature returns the ambient temperature in °C.
func (w *WeatherModel) Temperature() float32 {
	return w.temperature
}

// Humidity returns the relative humidity in %.
func (w *WeatherModel) Humidity() float32 {
	return w.humidity
}

// Pressure returns the atmospheric pressure in hPa.
func (w *WeatherModel) Pressure() float32 {
	return w.pressure
}

// CloudCover returns the cloud cover fraction (0-1).
func (w *WeatherModel) CloudCover() float32 {
	return w.cloudCover
}

// RainRate returns the rainfall rate in mm/h.
func (w *WeatherModel) RainRate() float32 {
	return w.rainRate
}

// Tick advances the weather model state.
func (w *WeatherModel) Tick() {
	// Simplified weather model
	// Temperature varies with sun elevation
	tempDelta := (float32(random()) - 0.5) * 0.2
	w.temperature += tempDelta
	w.temperature = clamp(w.temperature, -40, 50)

	// Humidity inversely correlated with temperature
	humDelta := -tempDelta * 2.0 + (float32(random())-0.5)*0.5
	w.humidity += humDelta
	w.humidity = clamp(w.humidity, 0, 100)

	// Pressure slowly varies
	pressDelta := (float32(random()) - 0.5) * 0.1
	w.pressure += pressDelta
	w.pressure = clamp(w.pressure, 950, 1050)

	// Cloud cover evolves slowly
	cloudDelta := (float32(random()) - 0.5) * 0.05
	w.cloudCover += cloudDelta
	w.cloudCover = clamp(w.cloudCover, 0, 1)

	// Rain when cloud cover is high
	if w.cloudCover > 0.7 && float32(random()) > 0.95 {
		w.rainRate = 5.0 * float32(random())
	} else {
		w.rainRate *= 0.95
	}
}

// ReservoirModel represents a water reservoir.
type ReservoirModel struct {
	ModelBase

	level     float32 // m (water level above bottom)
	flowIn    float32 // m³/s (inflow)
	flowOut   float32 // m³/s (outflow)
	area      float32 // m² (surface area)
	temp      float32 // °C (water temperature)
}

// NewReservoirModel creates a new Reservoir model.
func NewReservoirModel(id ModelID, area float32) *ReservoirModel {
	return &ReservoirModel{
		ModelBase: NewModelBase(id, "reservoir"),
		level:     50.0,
		area:      area,
		flowIn:    10.0,
		flowOut:   5.0,
		temp:      15.0,
	}
}

// Level returns the water level in meters.
func (r *ReservoirModel) Level() float32 {
	return r.level
}

// FlowIn returns the inflow rate in m³/s.
func (r *ReservoirModel) FlowIn() float32 {
	return r.flowIn
}

// FlowOut returns the outflow rate in m³/s.
func (r *ReservoirModel) FlowOut() float32 {
	return r.flowOut
}

// Temperature returns the water temperature in °C.
func (r *ReservoirModel) Temperature() float32 {
	return r.temp
}

// InjectFlow adds inflow to the reservoir.
func (r *ReservoirModel) InjectFlow(m3s float32) {
	r.flowIn += m3s
}

// ExtractFlow removes outflow from the reservoir.
func (r *ReservoirModel) ExtractFlow(m3s float32) {
	r.flowOut += m3s
}

// Tick advances the reservoir model state.
func (r *ReservoirModel) Tick() {
	// Calculate net flow
	netFlow := r.flowIn - r.flowOut

	// Update level (Volume = Area * Level, dV/dt = A * dL/dt)
	// So dL/dt = netFlow / Area
	levelDelta := netFlow / r.area
	r.level += levelDelta
	r.level = clamp(r.level, 0, 100)

	// Reset flow for next tick
	r.flowIn = 0
	r.flowOut = 0
}

// Helper function to clamp values
func clamp(v, min, max float32) float32 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// seededRand is a deterministic random number generator for model behavior.
// Using a seeded source ensures reproducible simulations.
var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// random returns a pseudo-random float64 between 0 and 1.
func random() float64 {
	return seededRand.Float64()
}
