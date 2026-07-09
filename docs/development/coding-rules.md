# Coding Rules

## Philosophy

Write code that reflects the architecture: simple, memory-driven, device-centric.

## General

- Keep code small
- Name things clearly
- Handle errors explicitly
- Test behavior, not implementation

## Memory Access

Behaviors read and write device memory:

```go
// Write
device.Memory().Write("input_registers", addr, value)

// Read
val := device.Memory().Read("input_registers", addr)
```

## Error Handling

Always handle errors:

```go
value, err := device.Memory().Read(...)
if err != nil {
    return fmt.Errorf("behavior %s: %w", b.ID(), err)
}
```

## Ownership

Code reflects ownership:

- `device.Memory()` - Device owns memory
- `device.AddFault()` - Device owns faults
- `device.ExposeProtocol()` - Device owns protocols
- `runtime.LoadPlugins()` - Runtime owns plugin loading

## Testing

Test behavior:

```go
func TestPVModel_CalculatesPower(t *testing.T) {
    device := setupDevice()
    behavior := NewPVModelBehavior()
    behavior.Attach(device)
    
    behavior.Tick()
    
    power := device.Memory().ReadFloat32("input_registers", 8)
    assert.InDelta(t, 10000, power, 1)
}
```

## Determinism

Behaviors must be deterministic:

- No `time.Now()`
- No unseeded random
- No external system calls

## Naming

| Thing | Style |
|-------|-------|
| Variables | camelCase |
| Functions | PascalCase |
| Constants | PascalCase |
| Packages | lowercase |

## No Premature Abstraction

Don't add layers that don't exist in the architecture:

- No "managers" unless genuinely needed
- No "controllers" unless genuinely needed
- No "services" unless genuinely needed

The architecture is simple. Code should be too.
