# ADR-003: Memory Model Design

**ADR ID:** ADR-003  
**Title:** Device Memory Model with Quality Flags  
**Date:** 2026-07-11  
**Status:** Accepted  
**Deciders:** Engineering Team  
**Repository:** https://github.com/tamzrod/forge

---

## Context

Virtual devices need to store data internally, similar to PLC memory:
- Named memory regions (Holding Registers, Input Registers, etc.)
- Typed access (uint16, float32, uint32)
- Quality tracking (Good, Uncertain, Bad, Offline)
- Atomic read/write operations

---

## Decision

We adopt a **Region-Based Memory Model** with quality tracking:

```go
type Memory struct {
    regions map[string]*Region
    mu      sync.RWMutex
}

type Region struct {
    name    string
    base    uint16
    size    uint16
    data    []byte
    quality []Quality
    access  Access
}
```

### Memory Regions

| Region | Access | Purpose |
|--------|--------|---------|
| Holding Registers | Read/Write | Device configuration and control |
| Input Registers | Read Only | Device measurements and status |
| Coils | Read/Write | Binary control bits |
| Discrete Inputs | Read Only | Binary status bits |

---

## Implementation Details

### Quality Flags

```go
type Quality int

const (
    QualityGood     Quality = 0
    QualityUncertain Quality = 1
    QualityBad       Quality = 2
    QualityOffline   Quality = 3
)
```

### Read Operations

```go
func (r *Region) ReadFloat32(offset uint16) (float32, error)
func (r *Region) ReadUint16(offset uint16) (uint16, error)
func (r *Region) ReadUint32(offset uint16) (uint32, error)
func (r *Region) Quality(offset uint16) Quality
```

### Write Operations

```go
func (r *Region) WriteFloat32(offset uint16, value float32) error
func (r *Region) WriteUint16(offset uint16, value uint16) error
func (r *Region) WriteUint32(offset uint16, value uint32) error
func (r *Region) SetQuality(offset uint16, quality Quality)
```

### Device Memory API

```go
type Device interface {
    Memory() *Memory
    AddMemoryRegion(name string, base, size uint16, access Access) error
}
```

---

## Quality Propagation

When reading, quality is checked at the memory level:

```go
func (r *Region) ReadFloat32(offset uint16) (float32, error) {
    if r.quality[offset] != QualityGood {
        return 0, ErrQualityNotGood
    }
    // ... return value
}
```

### Default Quality

- All memory initialized with `QualityGood`
- Quality persists until explicitly changed
- Use `SetQuality()` to mark bad data

---

## Consequences

### Positive
- PLC-like memory model familiar to industrial developers
- Quality flags enable fault simulation
- Thread-safe with mutex protection
- Extensible region types

### Negative
- Fixed-size regions
- No atomic multi-register operations
- Quality is per-register, not per-region

### Risks
- Large memory regions may impact performance
- Quality not automatically propagated from simulation models

---

## Alternatives Considered

### Alternative 1: Map-Based Memory
```go
memory := map[string]interface{}{}
```
**Rejected**: No bounds checking, poor cache locality.

### Alternative 2: No Quality Tracking
**Rejected**: Cannot simulate bad sensor data, wire faults.

---

## References

- [Memory Model](docs/architecture/memory-model.md)
- [Memory Model Comparison](docs/architecture/memory-model-comparison.md)
- [Protocol Architecture](docs/architecture/protocol-architecture.md)

---

## Related ADRs

- ADR-001: Runtime Architecture
- ADR-002: Behavior Model Design

---

## Testing

Memory model tests cover:
- Read/write operations (uint16, float32, uint32)
- Bounds checking
- Quality flag operations
- Concurrent access
- Region management

---

## Milestone Traceability

| Milestone | Status |
|-----------|--------|
| 1.3.1 Memory Core | ✅ Complete |
| 1.3.2 Quality Flags | ✅ Complete |
| 1.3.3 Region Management | ✅ Complete |
| 1.3.4 Memory Testing | ✅ Complete |
| 1.3.5 Protocol Integration | ⏳ Pending |
