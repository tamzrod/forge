// Package electrical provides entity models for electrical power systems.
package electrical

import (
	"time"

	"github.com/tamzrod/forge/world"
)

// GridEntity represents the utility grid (infinite bus).
type GridEntity struct {
	world.BaseEntity
	voltage   float32
	frequency float32
}

// NewGrid creates a new grid entity.
func NewGrid(id world.EntityID, voltage, frequency float32) *GridEntity {
	e := &GridEntity{
		voltage:   voltage,
		frequency: frequency,
	}
	e.BaseEntity = world.NewBaseEntity(id, "grid")
	return e
}

// Tick updates grid state.
func (e *GridEntity) Tick(dt time.Duration) {
	// Grid is an infinite bus - voltage and frequency are fixed
	e.SetOutput("voltage", e.voltage)
	e.SetOutput("frequency", e.frequency)
}

// Measurements returns grid measurements.
func (e *GridEntity) Measurements() []world.Measurement {
	return []world.Measurement{
		{EntityID: e.ID(), Name: "voltage", Value: e.voltage, Unit: "V", Timestamp: time.Now()},
		{EntityID: e.ID(), Name: "frequency", Value: e.frequency, Unit: "Hz", Timestamp: time.Now()},
	}
}

// BusEntity represents an electrical bus.
type BusEntity struct {
	world.BaseEntity
	nominalVoltage float32
	voltage       float32
	powerIn       float32
	powerOut      float32
}

// NewBus creates a new bus entity.
func NewBus(id world.EntityID, nominalVoltage float32) *BusEntity {
	e := &BusEntity{
		nominalVoltage: nominalVoltage,
		voltage:        nominalVoltage,
	}
	e.BaseEntity = world.NewBaseEntity(id, "bus")
	return e
}

// Tick updates bus state based on power flow.
func (e *BusEntity) Tick(dt time.Duration) {
	// Get inputs
	if v := e.GetInput("power_injection"); v != nil {
		e.powerIn = v.(float32)
	}
	if v := e.GetInput("power_withdrawal"); v != nil {
		e.powerOut = v.(float32)
	}

	// Voltage droops with load (simplified model)
	netPower := e.powerIn - e.powerOut
	voltageDelta := netPower * 0.0001
	e.voltage += voltageDelta
	e.voltage += (e.nominalVoltage - e.voltage) * 0.05

	// Clamp voltage
	if e.voltage < e.nominalVoltage*0.97 {
		e.voltage = e.nominalVoltage * 0.97
	}
	if e.voltage > e.nominalVoltage*1.03 {
		e.voltage = e.nominalVoltage * 1.03
	}

	// Set outputs
	e.SetOutput("voltage", e.voltage)
	e.SetOutput("power_injection", e.powerIn)
	e.SetOutput("power_withdrawal", e.powerOut)
	e.SetOutput("net_power", netPower)

	// Reset accumulators
	e.powerIn = 0
	e.powerOut = 0
}

// Measurements returns bus measurements.
func (e *BusEntity) Measurements() []world.Measurement {
	return []world.Measurement{
		{EntityID: e.ID(), Name: "voltage", Value: e.voltage, Unit: "V", Timestamp: time.Now()},
		{EntityID: e.ID(), Name: "power_injection", Value: e.powerIn, Unit: "kW", Timestamp: time.Now()},
		{EntityID: e.ID(), Name: "power_withdrawal", Value: e.powerOut, Unit: "kW", Timestamp: time.Now()},
	}
}

// InjectPower adds power injection to the bus.
func (e *BusEntity) InjectPower(kw float32) {
	e.powerIn += kw
}

// WithdrawPower adds power withdrawal from the bus.
func (e *BusEntity) WithdrawPower(kw float32) {
	e.powerOut += kw
}

// BreakerEntity represents a circuit breaker.
type BreakerEntity struct {
	world.BaseEntity
	isOpen     bool
	tripCount  int
	closeCount int
}

// NewBreaker creates a new breaker entity.
func NewBreaker(id world.EntityID) *BreakerEntity {
	e := &BreakerEntity{}
	e.BaseEntity = world.NewBaseEntity(id, "breaker")
	return e
}

// Open opens the breaker.
func (e *BreakerEntity) Open() {
	if !e.isOpen {
		e.isOpen = true
		e.tripCount++
	}
}

// Close closes the breaker.
func (e *BreakerEntity) Close() {
	if e.isOpen {
		e.isOpen = false
		e.closeCount++
	}
}

// IsOpen returns true if the breaker is open.
func (e *BreakerEntity) IsOpen() bool {
	return e.isOpen
}

// Tick updates breaker state.
func (e *BreakerEntity) Tick(dt time.Duration) {
	e.SetOutput("state", map[string]interface{}{
		"open":        e.isOpen,
		"trip_count":  e.tripCount,
		"close_count": e.closeCount,
	})
}

// Measurements returns breaker measurements.
func (e *BreakerEntity) Measurements() []world.Measurement {
	state := "CLOSED"
	if e.isOpen {
		state = "OPEN"
	}
	return []world.Measurement{
		{EntityID: e.ID(), Name: "state", Value: state, Unit: "", Timestamp: time.Now()},
		{EntityID: e.ID(), Name: "trip_count", Value: e.tripCount, Unit: "", Timestamp: time.Now()},
	}
}

// LoadEntity represents an electrical load.
type LoadEntity struct {
	world.BaseEntity
	power       float32
	basePower   float32
	variableLoad bool
}

// NewLoad creates a new load entity.
func NewLoad(id world.EntityID, powerKW float32) *LoadEntity {
	e := &LoadEntity{
		power:     powerKW,
		basePower: powerKW,
	}
	e.BaseEntity = world.NewBaseEntity(id, "load")
	return e
}

// SetVariableLoad enables variable load behavior.
func (e *LoadEntity) SetVariableLoad(enabled bool) {
	e.variableLoad = enabled
}

// Tick updates load state.
func (e *LoadEntity) Tick(dt time.Duration) {
	// Variable load oscillates slightly
	if e.variableLoad {
		e.power = e.basePower * (0.9 + 0.2*float32(int(time.Now().UnixNano())%1000)/1000.0)
	}

	e.SetOutput("power", e.power)
}

// Measurements returns load measurements.
func (e *LoadEntity) Measurements() []world.Measurement {
	return []world.Measurement{
		{EntityID: e.ID(), Name: "power", Value: e.power, Unit: "kW", Timestamp: time.Now()},
	}
}

// GeneratorEntity represents a power generator.
type GeneratorEntity struct {
	world.BaseEntity
	powerOutput float32
	ratedPower  float32
	isRunning   bool
}

// NewGenerator creates a new generator entity.
func NewGenerator(id world.EntityID, ratedPowerKW float32) *GeneratorEntity {
	e := &GeneratorEntity{
		powerOutput: 0,
		ratedPower:  ratedPowerKW,
		isRunning:   false,
	}
	e.BaseEntity = world.NewBaseEntity(id, "generator")
	return e
}

// Start starts the generator.
func (e *GeneratorEntity) Start() {
	e.isRunning = true
}

// Stop stops the generator.
func (e *GeneratorEntity) Stop() {
	e.isRunning = false
	e.powerOutput = 0
}

// SetPowerOutput sets the power output.
func (e *GeneratorEntity) SetPowerOutput(kw float32) {
	if kw > e.ratedPower {
		kw = e.ratedPower
	}
	e.powerOutput = kw
}

// Tick updates generator state.
func (e *GeneratorEntity) Tick(dt time.Duration) {
	if !e.isRunning {
		e.powerOutput = 0
	}

	e.SetOutput("power", e.powerOutput)
	e.SetOutput("status", e.isRunning)
}

// Measurements returns generator measurements.
func (e *GeneratorEntity) Measurements() []world.Measurement {
	status := "STOPPED"
	if e.isRunning {
		status = "RUNNING"
	}
	return []world.Measurement{
		{EntityID: e.ID(), Name: "power", Value: e.powerOutput, Unit: "kW", Timestamp: time.Now()},
		{EntityID: e.ID(), Name: "status", Value: status, Unit: "", Timestamp: time.Now()},
		{EntityID: e.ID(), Name: "rated_power", Value: e.ratedPower, Unit: "kW", Timestamp: time.Now()},
	}
}

// TransformerEntity represents a power transformer.
type TransformerEntity struct {
	world.BaseEntity
	primaryVoltage   float32
	secondaryVoltage float32
	ratio            float32
	loading          float32
}

// NewTransformer creates a new transformer entity.
func NewTransformer(id world.EntityID, primaryV, secondaryV float32) *TransformerEntity {
	e := &TransformerEntity{
		primaryVoltage:   primaryV,
		secondaryVoltage: secondaryV,
		ratio:            secondaryV / primaryV,
		loading:          0,
	}
	e.BaseEntity = world.NewBaseEntity(id, "transformer")
	return e
}

// Tick updates transformer state.
func (e *TransformerEntity) Tick(dt time.Duration) {
	// Get input power
	if v := e.GetInput("power"); v != nil {
		e.loading = v.(float32) / 1000.0 // Assume 1000kVA base
	}

	e.SetOutput("primary_voltage", e.primaryVoltage)
	e.SetOutput("secondary_voltage", e.secondaryVoltage)
	e.SetOutput("loading", e.loading)
}

// Measurements returns transformer measurements.
func (e *TransformerEntity) Measurements() []world.Measurement {
	return []world.Measurement{
		{EntityID: e.ID(), Name: "primary_voltage", Value: e.primaryVoltage, Unit: "V", Timestamp: time.Now()},
		{EntityID: e.ID(), Name: "secondary_voltage", Value: e.secondaryVoltage, Unit: "V", Timestamp: time.Now()},
		{EntityID: e.ID(), Name: "loading", Value: e.loading, Unit: "pu", Timestamp: time.Now()},
	}
}

// MeterEntity represents a revenue/power meter.
type MeterEntity struct {
	world.BaseEntity
	voltage        float32
	frequency      float32
	activePower    float32
	reactivePower  float32
	powerFactor    float32
	energyExport   float32
	energyImport   float32
	tickHours      float32
}

// NewMeter creates a new meter entity.
func NewMeter(id world.EntityID) *MeterEntity {
	e := &MeterEntity{
		voltage:     69000,
		frequency:   60.0,
		powerFactor: 1.0,
	}
	e.BaseEntity = world.NewBaseEntity(id, "meter")
	return e
}

// InjectInput sets an input value directly.
func (e *MeterEntity) InjectInput(name string, value interface{}) {
	e.SetInput(name, value)
}

// Tick updates meter state.
func (e *MeterEntity) Tick(dt time.Duration) {
	// Get inputs
	if v := e.GetInput("voltage"); v != nil {
		e.voltage = v.(float32)
	}
	if v := e.GetInput("frequency"); v != nil {
		e.frequency = v.(float32)
	}
	if v := e.GetInput("active_power"); v != nil {
		e.activePower = v.(float32)
	}
	if v := e.GetInput("reactive_power"); v != nil {
		e.reactivePower = v.(float32)
	}

	// Calculate power factor
	if e.activePower != 0 {
		apparentPower := float32(0)
		if e.activePower > 0 {
			apparentPower = e.activePower / 0.95 // Simplified
		}
		if apparentPower > 0 {
			e.powerFactor = e.activePower / apparentPower
		}
	}

	// Calculate apparent power
	apparentPower := float32(0)
	if e.powerFactor > 0 {
		apparentPower = e.activePower / e.powerFactor
	}

	// Accumulate energy (kWh)
	e.tickHours = float32(dt.Seconds()) / 3600.0
	if e.activePower > 0 {
		e.energyExport += e.activePower * e.tickHours
	} else if e.activePower < 0 {
		e.energyImport += (-e.activePower) * e.tickHours
	}

	// Set outputs
	e.SetOutput("voltage", e.voltage)
	e.SetOutput("frequency", e.frequency)
	e.SetOutput("active_power", e.activePower)
	e.SetOutput("reactive_power", e.reactivePower)
	e.SetOutput("apparent_power", apparentPower)
	e.SetOutput("power_factor", e.powerFactor)
	e.SetOutput("energy_export", e.energyExport)
	e.SetOutput("energy_import", e.energyImport)
	e.SetOutput("direction", e.activePower)
}

// Measurements returns meter measurements.
func (e *MeterEntity) Measurements() []world.Measurement {
	direction := "EXPORT"
	if e.activePower < 0 {
		direction = "IMPORT"
	}

	return []world.Measurement{
		{EntityID: e.ID(), Name: "voltage", Value: e.voltage, Unit: "V", Timestamp: time.Now()},
		{EntityID: e.ID(), Name: "frequency", Value: e.frequency, Unit: "Hz", Timestamp: time.Now()},
		{EntityID: e.ID(), Name: "active_power", Value: e.activePower, Unit: "kW", Timestamp: time.Now()},
		{EntityID: e.ID(), Name: "reactive_power", Value: e.reactivePower, Unit: "kVAr", Timestamp: time.Now()},
		{EntityID: e.ID(), Name: "apparent_power", Value: float32(0), Unit: "kVA", Timestamp: time.Now()},
		{EntityID: e.ID(), Name: "power_factor", Value: e.powerFactor, Unit: "", Timestamp: time.Now()},
		{EntityID: e.ID(), Name: "energy_export", Value: e.energyExport, Unit: "kWh", Timestamp: time.Now()},
		{EntityID: e.ID(), Name: "energy_import", Value: e.energyImport, Unit: "kWh", Timestamp: time.Now()},
		{EntityID: e.ID(), Name: "direction", Value: direction, Unit: "", Timestamp: time.Now()},
	}
}
