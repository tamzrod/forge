# Fault Model

## Philosophy

**Faults modify how device memory behaves. Faults belong to devices.**

## Fault Principles

1. **Faults modify memory access** - They intercept reads and writes
2. **Faults never own state** - They only modify behavior
3. **Faults belong to devices** - The runtime knows nothing about faults

## Fault Types

| Type | Effect |
|------|--------|
| **Frozen** | Writes are ignored |
| **Noise** | Reads return corrupted values |
| **Offline** | All locations set to offline quality |
| **Bad Quality** | Quality flag set to bad |

## Device Owns Faults

```go
device.AddFault(NewFrozenValuesFault())
device.AddFault(NewOfflineFault())
```

## Memory Access with Faults

```go
func (m *Memory) Read(region string, addr uint32) []byte {
    for _, fault := range m.faults {
        if fault.ModifiesRead() {
            return fault.ModifyRead(m.rawMemory, region, addr)
        }
    }
    return m.rawMemory.Read(region, addr)
}
```

## Quality Flags

Faults set quality flags:

```go
const (
    QualityGood     Quality = 0x00
    QualityUncertain Quality = 0x40
    QualityBad      Quality = 0x80
    QualityOffline  Quality = 0x84
)
```

## Example

```go
// Frozen fault - reads return last value, writes ignored
type FrozenFault struct {
    frozenValues map[string][]byte
}

func (f *FrozenFault) ModifyRead(mem *Memory, region string, addr uint32) []byte {
    return f.frozenValues[region+string(addr)]
}

func (f *FrozenFault) ModifyWrite(...) error {
    return nil  // Ignored
}
```

## Configuration

```yaml
device:
  faults:
    - type: frozen
      region: input_registers
      duration: 60s
```
