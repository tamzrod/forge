# Scenario Engine

## Philosophy

**Scenarios inject events into the simulation. Events affect devices.**

## Events

Scenarios inject events at specific times:

| Event | Effect |
|-------|--------|
| **Setpoint change** | Write to device memory |
| **Fault injection** | Add fault to device |
| **Fault removal** | Remove fault from device |
| **State change** | Change device state |

## Scenario Definition

```yaml
scenario:
  name: cloud_arrival
  events:
    - at: 10s
      setpoint:
        device: pv-001
        region: input_registers
        address: 0
        value: 800  # W/m²
    - at: 30s
      setpoint:
        device: pv-001
        region: input_registers
        address: 0
        value: 1000  # Full sun
```

## Event Execution

```go
func (e *ScenarioEngine) Tick() {
    e.simulationTime += e.tickInterval
    
    for _, event := range e.pendingEvents {
        if event.Time <= e.simulationTime {
            event.Execute(e.runtime)
            e.removeEvent(event)
        }
    }
}
```

## Device Communication

Events write to device memory:

```go
type SetpointEvent struct {
    DeviceID string
    Region  string
    Address uint32
    Value   []byte
}

func (e *SetpointEvent) Execute(rt *Runtime) {
    device := rt.Device(e.DeviceID)
    device.Memory().Write(e.Region, e.Address, e.Value)
}
```

## Scenarios and Faults

Scenarios can inject faults:

```yaml
events:
  - at: 30s
    inject_fault:
      device: meter-001
      fault: communication_loss
  - at: 60s
    remove_fault:
      device: meter-001
      fault: communication_loss
```
