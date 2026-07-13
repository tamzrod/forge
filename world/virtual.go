// Package world provides simulation entities representing the physical world.
package world

import (
	"math"
	"sync"
	"time"
)

// VirtualGenerator represents any source capable of injecting electrical power.
// It is a behavioral placeholder - technology-specific generators extend this.
type VirtualGenerator struct {
	BaseEntity

	mu sync.RWMutex

	// Name for display
	name string

	// Power output
	activePower   float32 // kW
	reactivePower float32 // kVAr

	// Capacity
	ratedCapacity    float32 // kW
	availableCapacity float32 // kW

	// Status
	isOnline      bool
	isDispatchable bool
	rampRate      float32 // kW per minute

	// Internal state
	targetPower   float32 // kW - dispatch target
	minPower      float32 // kW - minimum output
	maxPower      float32 // kW - maximum output
	efficiency    float32 // 0.0 - 1.0
}

// NewVirtualGenerator creates a new virtual generator.
func NewVirtualGenerator(id EntityID, name string, ratedCapacity float32) *VirtualGenerator {
	e := &VirtualGenerator{
		name:              name,
		ratedCapacity:     ratedCapacity,
		availableCapacity: ratedCapacity,
		maxPower:         ratedCapacity,
		minPower:         0,
		efficiency:       1.0,
		rampRate:         ratedCapacity, // Can ramp to full in 1 minute
		isOnline:         true,
		isDispatchable:   true,
	}
	e.BaseEntity = NewBaseEntity(id, "virtual-generator")
	return e
}

// ActivePower returns the current active power output in kW.
func (g *VirtualGenerator) ActivePower() float32 {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.activePower
}

// ReactivePower returns the current reactive power output in kVAr.
func (g *VirtualGenerator) ReactivePower() float32 {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.reactivePower
}

// RatedCapacity returns the rated capacity in kW.
func (g *VirtualGenerator) RatedCapacity() float32 {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.ratedCapacity
}

// AvailableCapacity returns the currently available capacity in kW.
func (g *VirtualGenerator) AvailableCapacity() float32 {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.availableCapacity
}

// IsOnline returns true if the generator is online.
func (g *VirtualGenerator) IsOnline() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.isOnline
}

// IsDispatchable returns true if the generator can be dispatched.
func (g *VirtualGenerator) IsDispatchable() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.isDispatchable
}

// RampRate returns the ramp rate in kW per minute.
func (g *VirtualGenerator) RampRate() float32 {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.rampRate
}

// SetOnline sets the online status.
func (g *VirtualGenerator) SetOnline(online bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.isOnline = online
	if !online {
		g.activePower = 0
		g.reactivePower = 0
		g.availableCapacity = 0
	}
}

// SetTargetPower sets the dispatch target power in kW.
func (g *VirtualGenerator) SetTargetPower(kw float32) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.targetPower = clamp(kw, g.minPower, g.maxPower)
}

// SetAvailableCapacity sets the available capacity in kW.
func (g *VirtualGenerator) SetAvailableCapacity(kw float32) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.availableCapacity = clamp(kw, 0, g.ratedCapacity)
}

// SetReactivePower sets the reactive power output in kVAr.
func (g *VirtualGenerator) SetReactivePower(kvar float32) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.reactivePower = kvar
}

// SetRampRate sets the ramp rate in kW per minute.
func (g *VirtualGenerator) SetRampRate(rate float32) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.rampRate = rate
}

// SetEfficiency sets the conversion efficiency (0.0-1.0).
func (g *VirtualGenerator) SetEfficiency(eff float32) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.efficiency = clamp(eff, 0, 1)
}

// Tick advances the generator state.
func (g *VirtualGenerator) Tick(dt time.Duration) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.isOnline {
		g.activePower = 0
		g.reactivePower = 0
		g.availableCapacity = 0
		return
	}

	// Ramp toward target power
	ramp := g.rampRate * float32(dt.Minutes())
	diff := g.targetPower - g.activePower

	if diff > ramp {
		g.activePower += ramp
	} else if diff < -ramp {
		g.activePower -= ramp
	} else {
		g.activePower = g.targetPower
	}

	// Clamp to available capacity
	if g.activePower > g.availableCapacity {
		g.activePower = g.availableCapacity
	}

	// Set outputs
	g.SetOutput("active_power", g.activePower)
	g.SetOutput("reactive_power", g.reactivePower)
	g.SetOutput("rated_capacity", g.ratedCapacity)
	g.SetOutput("available_capacity", g.availableCapacity)
	g.SetOutput("is_online", g.isOnline)
	g.SetOutput("is_dispatchable", g.isDispatchable)
}

// Measurements returns generator measurements.
func (g *VirtualGenerator) Measurements() []Measurement {
	g.mu.RLock()
	defer g.mu.RUnlock()

	now := time.Now()
	return []Measurement{
		{EntityID: g.ID(), Name: "active_power", Value: g.activePower, Unit: "kW", Timestamp: now},
		{EntityID: g.ID(), Name: "reactive_power", Value: g.reactivePower, Unit: "kVAr", Timestamp: now},
		{EntityID: g.ID(), Name: "rated_capacity", Value: g.ratedCapacity, Unit: "kW", Timestamp: now},
		{EntityID: g.ID(), Name: "available_capacity", Value: g.availableCapacity, Unit: "kW", Timestamp: now},
		{EntityID: g.ID(), Name: "is_online", Value: boolToFloat(g.isOnline), Unit: "", Timestamp: now},
		{EntityID: g.ID(), Name: "is_dispatchable", Value: boolToFloat(g.isDispatchable), Unit: "", Timestamp: now},
	}
}

// InjectPower is called by the network to record power injection.
// Returns the actual power that can be injected.
func (g *VirtualGenerator) InjectPower(kw float32) float32 {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.isOnline {
		return 0
	}

	// Apply capacity limits
	injectable := min(kw, g.availableCapacity)
	injectable = clamp(injectable, 0, g.maxPower)

	g.activePower = injectable
	g.SetOutput("active_power", injectable)

	return injectable
}

// VirtualLoad represents any consumer of electrical power.
// It is a behavioral placeholder - technology-specific loads extend this.
type VirtualLoad struct {
	BaseEntity

	mu sync.RWMutex

	// Name for display
	name string

	// Power demand
	activePowerDemand   float32 // kW
	reactivePowerDemand float32 // kVAr

	// Status
	isConnected bool
	priority    int // 1-10, higher = more important

	// Internal state
	basePower  float32 // kW - baseline demand
	powerScale float32 // 0.0-1.0 - scaling factor
}

// NewVirtualLoad creates a new virtual load.
func NewVirtualLoad(id EntityID, name string, basePower float32) *VirtualLoad {
	e := &VirtualLoad{
		name:                 name,
		activePowerDemand:    basePower,
		reactivePowerDemand:  0,
		basePower:            basePower,
		powerScale:          1.0,
		isConnected:          true,
		priority:             5,
	}
	e.BaseEntity = NewBaseEntity(id, "virtual-load")
	return e
}

// ActivePowerDemand returns the current active power demand in kW.
func (l *VirtualLoad) ActivePowerDemand() float32 {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.activePowerDemand
}

// ReactivePowerDemand returns the current reactive power demand in kVAr.
func (l *VirtualLoad) ReactivePowerDemand() float32 {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.reactivePowerDemand
}

// IsConnected returns true if the load is connected.
func (l *VirtualLoad) IsConnected() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.isConnected
}

// Priority returns the load priority (1-10).
func (l *VirtualLoad) Priority() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.priority
}

// SetConnected sets the connected status.
func (l *VirtualLoad) SetConnected(connected bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.isConnected = connected
	if !connected {
		l.activePowerDemand = 0
		l.reactivePowerDemand = 0
	}
}

// SetPowerDemand sets the power demand in kW.
func (l *VirtualLoad) SetPowerDemand(kw float32) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.activePowerDemand = max(0, kw)
	l.basePower = l.activePowerDemand
}

// SetPowerScale sets the scaling factor (0.0-1.0).
func (l *VirtualLoad) SetPowerScale(scale float32) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.powerScale = clamp(scale, 0, 1)
	l.activePowerDemand = l.basePower * l.powerScale
}

// SetPriority sets the load priority (1-10).
func (l *VirtualLoad) SetPriority(p int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.priority = clampInt(p, 1, 10)
}

// SetReactivePowerDemand sets the reactive power demand in kVAr.
func (l *VirtualLoad) SetReactivePowerDemand(kvar float32) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.reactivePowerDemand = kvar
}

// Tick advances the load state.
func (l *VirtualLoad) Tick(dt time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.isConnected {
		l.activePowerDemand = 0
		l.reactivePowerDemand = 0
	} else {
		l.activePowerDemand = l.basePower * l.powerScale
	}

	// Set outputs
	l.SetOutput("active_power_demand", l.activePowerDemand)
	l.SetOutput("reactive_power_demand", l.reactivePowerDemand)
	l.SetOutput("is_connected", l.isConnected)
	l.SetOutput("priority", float32(l.priority))
}

// Measurements returns load measurements.
func (l *VirtualLoad) Measurements() []Measurement {
	l.mu.RLock()
	defer l.mu.RUnlock()

	now := time.Now()
	return []Measurement{
		{EntityID: l.ID(), Name: "active_power_demand", Value: l.activePowerDemand, Unit: "kW", Timestamp: now},
		{EntityID: l.ID(), Name: "reactive_power_demand", Value: l.reactivePowerDemand, Unit: "kVAr", Timestamp: now},
		{EntityID: l.ID(), Name: "is_connected", Value: boolToFloat(l.isConnected), Unit: "", Timestamp: now},
		{EntityID: l.ID(), Name: "priority", Value: float32(l.priority), Unit: "", Timestamp: now},
	}
}

// WithdrawPower is called by the network to record power withdrawal.
// Returns the actual power that can be withdrawn.
func (l *VirtualLoad) WithdrawPower(kw float32) float32 {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.isConnected {
		return 0
	}

	// Load can withdraw up to its demand
	withdrawable := min(kw, l.activePowerDemand)
	withdrawable = max(0, withdrawable)

	return withdrawable
}

// VirtualMeter represents a virtual power meter.
// It measures power flow at a point in the network.
type VirtualMeter struct {
	BaseEntity

	mu sync.RWMutex

	// Name for display
	name string

	// Measured values
	activePower    float32 // kW - positive = import, negative = export
	reactivePower  float32 // kVAr
	apparentPower float32 // kVA
	powerFactor   float32

	// Energy
	energyImport float32 // kWh
	energyExport float32 // kWh

	// Voltage and current
	voltage float32 // V
	current float32 // A
	frequency float32 // Hz
}

// NewVirtualMeter creates a new virtual meter.
func NewVirtualMeter(id EntityID, name string) *VirtualMeter {
	e := &VirtualMeter{
		name:       name,
		voltage:    480,
		frequency: 60,
		powerFactor: 1.0,
	}
	e.BaseEntity = NewBaseEntity(id, "virtual-meter")
	return e
}

// ActivePower returns the measured active power in kW.
func (m *VirtualMeter) ActivePower() float32 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.activePower
}

// ReactivePower returns the measured reactive power in kVAr.
func (m *VirtualMeter) ReactivePower() float32 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.reactivePower
}

// ApparentPower returns the measured apparent power in kVA.
func (m *VirtualMeter) ApparentPower() float32 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.apparentPower
}

// PowerFactor returns the measured power factor.
func (m *VirtualMeter) PowerFactor() float32 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.powerFactor
}

// EnergyImport returns the total imported energy in kWh.
func (m *VirtualMeter) EnergyImport() float32 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.energyImport
}

// EnergyExport returns the total exported energy in kWh.
func (m *VirtualMeter) EnergyExport() float32 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.energyExport
}

// Voltage returns the measured voltage in V.
func (m *VirtualMeter) Voltage() float32 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.voltage
}

// Frequency returns the measured frequency in Hz.
func (m *VirtualMeter) Frequency() float32 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.frequency
}

// SetMeasurements updates the meter measurements.
func (m *VirtualMeter) SetMeasurements(p, q, v, f float32) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.activePower = p
	m.reactivePower = q
	m.voltage = v
	m.frequency = f

	// Calculate apparent power
	m.apparentPower = sqrt(p*p + q*q)

	// Calculate power factor
	if m.apparentPower > 0 {
		m.powerFactor = p / m.apparentPower
	} else {
		m.powerFactor = 1.0
	}

	// Set outputs
	m.SetOutput("active_power", m.activePower)
	m.SetOutput("reactive_power", m.reactivePower)
	m.SetOutput("apparent_power", m.apparentPower)
	m.SetOutput("power_factor", m.powerFactor)
	m.SetOutput("voltage", m.voltage)
	m.SetOutput("frequency", m.frequency)
}

// RecordEnergy accumulates energy based on power and time step.
func (m *VirtualMeter) RecordEnergy(dt time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Energy in kWh = power in kW * time in hours
	energy := m.activePower * float32(dt.Hours())

	if energy > 0 {
		m.energyImport += energy
	} else {
		m.energyExport += -energy
	}

	m.SetOutput("energy_import", m.energyImport)
	m.SetOutput("energy_export", m.energyExport)
}

// Tick advances the meter state.
func (m *VirtualMeter) Tick(dt time.Duration) {
	m.RecordEnergy(dt)
}

// Measurements returns meter measurements.
func (m *VirtualMeter) Measurements() []Measurement {
	m.mu.RLock()
	defer m.mu.RUnlock()

	now := time.Now()
	return []Measurement{
		{EntityID: m.ID(), Name: "active_power", Value: m.activePower, Unit: "kW", Timestamp: now},
		{EntityID: m.ID(), Name: "reactive_power", Value: m.reactivePower, Unit: "kVAr", Timestamp: now},
		{EntityID: m.ID(), Name: "apparent_power", Value: m.apparentPower, Unit: "kVA", Timestamp: now},
		{EntityID: m.ID(), Name: "power_factor", Value: m.powerFactor, Unit: "", Timestamp: now},
		{EntityID: m.ID(), Name: "voltage", Value: m.voltage, Unit: "V", Timestamp: now},
		{EntityID: m.ID(), Name: "frequency", Value: m.frequency, Unit: "Hz", Timestamp: now},
		{EntityID: m.ID(), Name: "energy_import", Value: m.energyImport, Unit: "kWh", Timestamp: now},
		{EntityID: m.ID(), Name: "energy_export", Value: m.energyExport, Unit: "kWh", Timestamp: now},
	}
}

// Helper functions

func clamp(v, min, max float32) float32 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func clampInt(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func boolToFloat(b bool) float32 {
	if b {
		return 1
	}
	return 0
}

func sqrt(x float32) float32 {
	return float32(math.Sqrt(float64(x)))
}

// Name returns the entity name.
func (g *VirtualGenerator) Name() string {
return g.name
}

// Name returns the entity name.
func (l *VirtualLoad) Name() string {
return l.name
}

// Name returns the entity name.
func (m *VirtualMeter) Name() string {
return m.name
}
