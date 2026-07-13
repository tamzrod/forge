# Event Ownership Model

## Design Principle

**Forge simulates causes first and consequences second.**

```
Scenario creates condition
        ↓
World responds (physical effects)
        ↓
Protection makes decisions
        ↓
Equipment acts
        ↓
Measurements observe
```

## Event Ownership

### 1. Scenario - Creates Engineering Conditions

Scenarios describe **what happens** in the simulated world, not **how equipment responds**.

**Scenarios do NOT:**
- Directly command breakers
- Set protection relay states
- Manipulate equipment internals

**Scenarios DO:**
- Publish engineering conditions
- Schedule when conditions occur
- Describe the experiment timeline

**Condition Events (Scenario → World):**
```go
// Grid conditions
type GridCondition struct {
    Type "voltage_sag" | "frequency_disturbance" | "utility_outage"
    Duration time.Duration
    Severity float32  // 0.0 to 1.0
}

// Fault conditions
type FaultCondition struct {
    Type "single_phase" | "two_phase" | "three_phase"
    Location string  // bus ID
    Impedance float32
}

// Environmental conditions
type WeatherCondition struct {
    Type "cloud_cover" | "wind_gust" | "fog"
    Intensity float32
}

// Test conditions
type TestCondition struct {
    Type "island_command" | "sync_check_fail"
    Reason string
}
```

### 2. World/Electrical - Physical Response

The electrical network responds to conditions with physical effects.

**World does:**
- Calculates power flow changes
- Applies fault impedances
- Tracks voltage/frequency deviations
- Detects fault currents

**World does NOT:**
- Command breakers
- Make protection decisions
- Override equipment settings

### 3. Protection - Decision Making

Protection systems observe conditions and make decisions.

**Protection does:**
- Monitors measurements
- Applies protection logic (overcurrent, undervoltage, etc.)
- Issues trip/close commands based on settings
- Respects coordination time delays

**Protection does NOT:**
- Know about scenarios
- Directly open breakers
- Override equipment

**Protection Events:**
```go
type ProtectionCommand struct {
    Type "trip_command" | "close_command" | "alarm"
    Source string  // relay ID
    Target string  // breaker ID
    Reason string  // "overcurrent", "undervoltage", etc.
    Delay time.Duration  // intentional delay for coordination
}
```

### 4. Equipment - Acts on Commands

Equipment responds to protection commands according to its characteristics.

**Equipment does:**
- Receives trip/close commands
- Performs operating actions (open/close)
- Respects operating times
- Reports state changes

**Equipment does NOT:**
- Know about scenarios
- Make protection decisions
- Create conditions

**Equipment Events:**
```go
type EquipmentState struct {
    Type "breaker_state" | "switch_state" | "generator_state"
    ID string
    State string  // "open", "closed", "tripped"
    Timestamp time.Time
}
```

### 5. Measurements - Observes Results

Measurements are passive observers of system state.

**Measurements do:**
- Report current values
- Record timestamps
- Include quality indicators

**Measurements do NOT:**
- Cause changes
- Make decisions
- Create events

## Event Flow Examples

### Example 1: Three-Phase Fault

```
Scenario: "three_phase_fault"
    └── publishes: {Type: "fault", Phases: 3, Location: "bus-69kV"}

World: electrical network
    └── calculates: fault currents, voltage collapse

Protection: relay "r-69kV-1"
    └── detects: overcurrent (50x rated)
    └── waits: 0.3s coordination delay
    └── issues: {Type: "trip_command", Target: "breaker-69kV"}

Equipment: breaker "breaker-69kV"
    └── receives: trip_command
    └── opens: 0.05s operating time
    └── state: "open"

Measurements:
    └── I-69kV: 0 A (breaker open)
    └── V-69kV: 0 V (no voltage)
```

### Example 2: Grid Voltage Sag

```
Scenario: "voltage_sag_test"
    └── publishes: {Type: "voltage_sag", Duration: 5s, Depth: 0.8}

World: grid model
    └── applies: 80% voltage for 5 seconds

Protection: relay "v-69kV-1"
    └── detects: undervoltage (80% for 3s)
    └── settings: UV delay = 2s
    └── issues: {Type: "trip_command", Target: "breaker-pcc"}

Equipment: breaker "breaker-pcc"
    └── opens after 2s delay

Measurements:
    └── V-69kV: 55200 V (80%)
    └── P-69kV: 0 kW (islanded)
```

### Example 3: Islanding Test

```
Scenario: "islanding_test"
    └── publishes: {Type: "island_command", Reason: "planned_test"}

Protection: sync check relay
    └── verifies: frequency match, voltage match, phase match
    └── issues: {Type: "close_command", Target: "breaker-pcc"}

Equipment: breaker "breaker-pcc"
    └── closes

Measurements:
    └── Frequency: 60.00 Hz (synced)
    └── Power: variable (load dependent)
```

## Event Type Categories

### Conditions (Scenario → World)

| Event | Parameters | Effect |
|-------|------------|--------|
| `fault` | phases, location, impedance | Current spikes, voltage collapse |
| `voltage_sag` | depth, duration | Grid voltage drops |
| `voltage_swell` | magnitude, duration | Grid voltage rises |
| `frequency_disturbance` | frequency, duration | Grid frequency deviation |
| `utility_outage` | duration | Complete loss of grid |
| `cloud_cover` | coverage (0-1) | Irradiance reduction |
| `wind_gust` | speed, direction | Mechanical effects |
| `island_command` | reason | Initiate islanding |
| `reconnect_command` | reason | Attempt grid reconnect |

### Commands (Protection → Equipment)

| Event | Parameters | Effect |
|-------|------------|--------|
| `trip_command` | breaker_id, reason, delay | Open breaker |
| `close_command` | breaker_id, delay | Close breaker |
| `alarm` | relay_id, reason | Alert operators |
| `trip_inhibit` | breaker_id | Block tripping |

### State Changes (Equipment → World)

| Event | Parameters | Effect |
|-------|------------|--------|
| `breaker_opened` | breaker_id, timestamp | Network topology changes |
| `breaker_closed` | breaker_id, timestamp | Network topology changes |
| `breaker_tripped` | breaker_id, fault | Protection trip occurred |
| `generator_started` | gen_id | Source added |
| `generator_stopped` | gen_id | Source removed |

### Measurements (World → External)

| Event | Parameters |
|-------|------------|
| `voltage` | entity_id, value, unit |
| `current` | entity_id, value, unit |
| `frequency` | entity_id, value |
| `power` | entity_id, active, reactive |
| `energy` | entity_id, import, export |

## Refactoring Guidelines

### Before (Wrong)

```go
// Scenario directly commands breaker
scenario.AddAction(EventAction{
    Time: 10 * time.Second,
    Type: "breaker_open",
    Data: map[string]interface{}{"breaker_id": "grid-breaker"},
})

// Example directly manipulates breaker
func handleEvent(evt world.Event) {
    if evt.Type == "breaker_open" {
        breaker.Open()  // Direct manipulation
    }
}
```

### After (Correct)

```go
// Scenario creates engineering condition
scenario.AddAction(EventAction{
    Time: 10 * time.Second,
    Type: "utility_outage",
    Data: map[string]interface{}{"duration": 30 * time.Second},
})

// Protection logic decides what to do
func (r *Relay) HandleEvent(evt world.Event) {
    switch evt.Type {
    case "utility_outage":
        // Protection logic
        if r.detectLossOfVoltage() {
            r.scheduleTrip("breaker-pcc", 0.1*time.Second)
        }
    }
}

// Equipment responds to protection command
func (b *Breaker) HandleEvent(evt world.Event) {
    switch evt.Type {
    case "trip_command":
        b.Open()
    }
}
```

## Summary

| Owner | Responsibility | Creates |
|-------|---------------|---------|
| Scenario | Engineering conditions | Conditions |
| World | Physical response | Physical effects |
| Protection | Decision making | Commands |
| Equipment | Actions | State changes |
| Measurements | Observation | Reports |

**Causality Chain:**
```
Scenario → Condition → Physical Effect → Protection Decision → Command → Equipment Action → Measurement
```
