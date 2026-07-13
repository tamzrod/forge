# Traceability Matrix

## Purpose

This matrix links requirements, architecture decisions, and code to ensure:
- Every requirement has an implementation
- Every decision is documented
- Changes can be analyzed for impact

---

## Requirements to Architecture

| Requirement | Source | Architecture Artifact | Status |
|-------------|--------|---------------------|--------|
| Deterministic execution | Vision | Runtime, Scheduler | ✓ |
| Memory as source of truth | Vision | Memory Model | ✓ |
| Virtual industrial devices | Vision | Device Model | ✓ |
| Physical world simulation | Vision | Simulation Models | ✓ |
| Protocol independence | Vision | MMA2 Integration | ✓ |

---

## Architecture to Implementation

| Architecture Decision | ADR | Implementation | Status |
|---------------------|-----|----------------|--------|
| Runtime hosts models and devices | ADR-001 | runtime/runtime.go | ✓ |
| Scheduler controls tick order | ADR-001 | scheduler/scheduler.go | ✓ |
| Device owns memory | ADR-001 | device/device.go, memory/memory.go | ✓ |
| Behaviors observe models | ADR-001 | device/behavior.go | ✓ |
| Models are physics-only | ADR-001 | models/*.go | ✓ |
| Raw Ingest to MMA2 | ADR-001 | publishers/rawingest/*.go | ✓ |

---

## Code to Tests

| Component | Source Files | Test Files | Coverage |
|-----------|--------------|------------|----------|
| Runtime | runtime/runtime.go | - | - |
| Scheduler | scheduler/scheduler.go | - | - |
| Device | device/device.go | device/device_test.go | - |
| Memory | memory/memory.go | memory/memory_test.go | - |
| Grid Model | models/grid.go | models/grid_test.go | - |
| Sun Model | models/sun.go | models/sun_test.go | - |
| Weather Model | models/weather.go | models/weather_test.go | - |
| Inspector | internal/inspector/*.go | - | - |
| Raw Ingest | internal/publishers/rawingest/*.go | protocol_test.go | - |

---

## Requirements to Tests

| Requirement | Test Coverage | Gaps |
|------------|--------------|------|
| Deterministic tick execution | device_test.go | Missing scheduler tick tests |
| Memory read/write | memory_test.go | ✓ |
| Model evolution | *_test.go | Limited |
| Device behavior | weatherstation_test.go | Few behaviors tested |

---

## Implementation Checklist

### Core Runtime (Must Have)

- [x] Runtime initialization
- [x] Device creation
- [x] Model creation
- [x] Scheduler tick loop
- [x] Clock advancement

### Memory System (Must Have)

- [x] MemoryImage struct
- [x] Read/Write operations
- [x] Region management
- [x] Quality flags

### Models (Must Have)

- [x] GridModel
- [x] SunModel
- [x] WindModel
- [x] WeatherModel
- [x] ReservoirModel

### Devices (Must Have)

- [x] Device struct
- [x] Behavior interface
- [x] WeatherStation device
- [ ] PV Inverter device
- [ ] Revenue Meter device
- [ ] Grid Proxy device

### Protocols (Must Have)

- [x] Raw Ingest client
- [ ] Modbus adapter
- [ ] DNP3 adapter

### Verification (Must Have)

- [ ] Unit tests for scheduler
- [x] Unit tests for memory
- [x] Unit tests for models
- [x] Unit tests for devices

---

## Gap Analysis

### Critical Gaps (Block Operation)

1. **No main/entry point** - Cannot build and run
2. **No CI/CD** - No automated testing
3. **No Go installation in test environment** - Cannot verify build

### Important Gaps (Should Have)

1. **Limited test coverage** - Most components untested
2. **Missing device types** - Only WeatherStation implemented
3. **No protocol adapters** - Only Raw Ingest implemented

### Nice to Have

1. Example configurations
2. CLI for simulation control
3. Web dashboard for monitoring

---

## Next Steps

### Priority 1: Make Forge Buildable

1. Add `cmd/forge/main.go` as entry point
2. Add `go.mod` if missing
3. Set up CI/CD pipeline

### Priority 2: Add Missing Device Types

1. PV Inverter device
2. Revenue Meter device
3. Grid Proxy device

### Priority 3: Add Protocol Adapters

1. Modbus TCP server
2. DNP3 server
3. REST API

---

## Related Documents

| Document | Purpose |
|----------|---------|
| [Vision](../architecture/vision.md) | Project requirements |
| [ADR-001](../architecture/adr-001-runtime-architecture.md) | Key architecture decisions |
| [Runtime](../architecture/runtime.md) | Runtime implementation |
| [Device Model](../architecture/device-model.md) | Device structure |

---

*Last Updated: 2026-07-13*  
*Type: Architecture Artifact*  
*Status: Initial*
