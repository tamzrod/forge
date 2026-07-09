# Plugin Architecture

## Philosophy

**Plugins provide device types. The runtime knows only the generic Device interface.**

## Role

A plugin provides:

- Device type definitions
- Behavior implementations
- Memory layouts

## Plugin Interface

```go
type Plugin interface {
    ID() string
    DeviceTypes() []DeviceType
}
```

## Example

```go
type EnergyPlugin struct{}

func (p *EnergyPlugin) DeviceTypes() []DeviceType {
    return []DeviceType{
        NewRevenueMeterType(),
        NewWeatherStationType(),
        NewPVInverterType(),
    }
}
```

## Device Types from Plugins

```
Energy Plugin
├── revenue_meter
├── weather_station
├── pv_inverter
├── relay
└── grid

Water Plugin
├── pump
├── valve
├── tank
└── flow_meter
```

## Runtime Knowledge

The runtime knows only:

```go
type Device struct {
    id        DeviceID
    memory    *MemoryImage
    behaviors []Behavior
    protocols []Protocol
}
```

The runtime does not know about power, flow, or any domain concept.

## Adding New Domains

Adding a new domain requires only new plugins:

1. Create plugin with device types
2. Load plugin
3. Create devices from types

No runtime changes.

## Configuration

```yaml
plugins:
  paths:
    - ./plugins/energy
    - ./plugins/water
```
