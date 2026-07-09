# Memory Model Comparison: MMA2 vs RMA

## Executive Summary

This document compares two memory appliance architectures—MMA2 and RMA—for use as the Simulation Runtime's operational memory integration point.

**Recommendation: MMA2 remains the better choice for the Simulation Runtime.**

While RMA offers more flexible memory primitives (arbitrary bit widths), this flexibility comes with unnecessary complexity for the simulation use case. MMA2's fixed Modbus memory model provides better ecosystem compatibility, simpler integration, and clearer semantics for industrial device simulation.

---

## 1. Overview

### MMA2 (Modbus Memory Appliance 2.0)

A deterministic, minimal, opinionated Modbus TCP memory appliance focused on storing and serving raw Modbus memory correctly, predictably, and safely.

**Design Philosophy:** "Boring, strict, predictable"

### RMA (Real-Time Memory Appliance)

A protocol-agnostic, variable-width real-time memory core designed for industrial telemetry and infrastructure systems.

**Design Philosophy:** "Memory is the product. Determinism over convenience."

---

## 2. Memory Model Comparison

### MMA2 Memory Model

```
Memory Types: FIXED (4 areas)
├── Coils              (1-bit, boolean)
├── Discrete Inputs     (1-bit, boolean)
├── Holding Registers   (16-bit unsigned)
└── Input Registers    (16-bit unsigned)

Identity: (Port, UnitID)
```

### RMA Memory Model

```
Memory Types: VARIABLE
├── Banks with arbitrary unit_width_bits
├── Examples: 1, 5, 8, 16, 128 bits
└── Flexible unit counts

Identity: Bank ID
```

### Analysis

| Aspect | MMA2 | RMA |
|--------|------|-----|
| **Type System** | Fixed 4 types | Arbitrary bit widths |
| **Flexibility** | Low | High |
| **Complexity** | Low | Medium |
| **Simplicity** | High | Medium |
| **Predictability** | Very High | High |

**Key Finding:** MMA2's fixed memory model is a strength, not a limitation, for industrial device simulation. The four Modbus memory types (Coils, Discrete Inputs, Holding Registers, Input Registers) map naturally to industrial device semantics.

---

## 3. Device Simulation Analysis

### Weather Station

| Model | Support | Notes |
|-------|---------|-------|
| MMA2 | ✓ | Input Registers for irradiance, temperature; Coils for status |
| RMA | ✓ | Custom widths possible, but unnecessary |

### Revenue Meter

| Model | Support | Notes |
|-------|---------|-------|
| MMA2 | ✓ | Input Registers for energy, power; Holding Registers for config |
| RMA | ✓ | Same capability |

### Relay

| Model | Support | Notes |
|-------|---------|-------|
| MMA2 | ✓ | Coils for control, Discrete Inputs for status |
| RMA | ✓ | Same capability |

### PLC

| Model | Support | Notes |
|-------|---------|-------|
| MMA2 | ✓ | All four areas used for different access patterns |
| RMA | ✓ | Custom widths add no value |

### Camera

| Model | Support | Notes |
|-------|---------|-------|
| MMA2 | Limited | No native support for image/frame data |
| RMA | Limited | Variable widths don't solve this either |

### GPS

| Model | Support | Notes |
|-------|---------|-------|
| MMA2 | ✓ | Registers for lat/lon/altitude |
| RMA | ✓ | Same |

### Battery

| Model | Support | Notes |
|-------|---------|-------|
| MMA2 | ✓ | Registers for SOC, voltage, current |
| RMA | ✓ | Same |

### SmartLogger

| Model | Support | Notes |
|-------|---------|-------|
| MMA2 | ✓ | All register types |
| RMA | ✓ | Same |

**Key Finding:** Both models can represent all standard industrial device types. RMA's variable-width memory provides no meaningful advantage for these use cases.

---

## 4. Protocol Adaptation

### MMA2 Protocol Support

| Protocol | Native | Notes |
|----------|--------|-------|
| Modbus TCP | ✓ | Native - built-in |
| Raw Ingest | ✓ | Built-in write-only ingest |
| REST | Adapter | External adapter |
| MQTT | Adapter | External adapter |

### RMA Protocol Support

| Protocol | Native | Notes |
|----------|--------|-------|
| Raw TCP | ✓ | Minimal transport |
| Modbus TCP | Adapter | Must be implemented |
| REST | Adapter | Must be implemented |
| MQTT | Adapter | Must be implemented |

### Analysis

| Aspect | MMA2 | RMA |
|--------|------|-----|
| Modbus TCP | Native | Requires adapter |
| Protocol Maturity | Proven | Newer |
| Ecosystem Fit | Excellent | Unknown |
| Integration Effort | Low | Higher |

**Key Finding:** MMA2 is already integrated with the Atlas ecosystem. RMA would require building protocol adapters from scratch.

---

## 5. Performance

### MMA2 Performance

- **Memory Layout:** Sequential, fixed-size regions
- **Cache Locality:** Excellent for sequential access
- **Allocation:** Static at startup, no runtime allocation
- **Determinism:** Guaranteed

### RMA Performance

- **Memory Layout:** Flexible, variable-sized banks
- **Cache Locality:** Depends on implementation
- **Allocation:** Static at startup (pre-allocated)
- **Determinism:** Guaranteed

**Key Finding:** Performance characteristics are similar. Both are static, pre-allocated memory systems.

---

## 6. Complexity

### API Complexity

| Metric | MMA2 | RMA |
|--------|------|-----|
| Memory Types | 4 fixed | Arbitrary |
| Configuration | Simple | More flexible |
| Learning Curve | Low | Medium |
| Error Surface | Small | Larger |

### Implementation Complexity

| Component | MMA2 | RMA |
|-----------|------|-----|
| Core Memory | Simple | Medium |
| Configuration | Simple | More flexible |
| State Sealing | Explicit flag | Implicit via init |
| Error Handling | 1-byte codes | Categorized errors |

### Maintenance Cost

| Factor | MMA2 | RMA |
|--------|------|-----|
| Documentation | Comprehensive | Good |
| Code Maturity | Mature | Newer |
| Community | Established | Developing |
| Support Burden | Lower | Higher |

**Key Finding:** MMA2 has lower complexity and maintenance burden.

---

## 7. Atlas Ecosystem Compatibility

### Current Ecosystem

```
Atlas-PPC → MMA2 → Modbus TCP
Replicator → MMA2 → Raw Ingest
Simulation Runtime → MMA2 → Raw Ingest
```

### MMA2 Compatibility

| Component | Compatibility |
|-----------|---------------|
| Atlas-PPC | ✓ Native Modbus TCP |
| Replicator | ✓ Native Raw Ingest |
| Simulation Runtime | ✓ Native Raw Ingest |
| Other Projects | ✓ Established pattern |

### RMA Compatibility

| Component | Compatibility |
|-----------|---------------|
| Atlas-PPC | Requires Modbus adapter |
| Replicator | Requires new ingest protocol |
| Simulation Runtime | Requires new integration |
| Other Projects | Unknown |

**Key Finding:** MMA2 is the established standard in the Atlas ecosystem. Switching to RMA would break existing integrations.

---

## 8. Long-Term Extensibility

### When RMA Excels

RMA's variable-width memory model provides advantages for:
- Non-standard device types (custom sensors with 5-bit enums)
- Specialized protocols with non-16-bit word sizes
- Research/prototyping scenarios

### When MMA2 Excels

MMA2's fixed model provides advantages for:
- Standard industrial devices (the majority)
- Ecosystem integration
- Predictability and simplicity
- Long-term maintainability

### Analysis

The Simulation Runtime's primary use case is simulating standard industrial devices (weather stations, meters, relays, inverters, etc.). These devices universally use standard Modbus semantics.

**Adding RMA would introduce complexity for edge cases that don't exist in the simulation domain.**

---

## 9. Side-by-Side Comparison

| Criterion | MMA2 | RMA | Winner |
|-----------|------|-----|--------|
| **Memory Model** | Fixed 4 types | Variable widths | MMA2 (simplicity) |
| **Device Mapping** | Natural for industrial | Overly flexible | MMA2 |
| **Protocol Support** | Modbus native | Adapter required | MMA2 |
| **Ecosystem Fit** | Established | New | MMA2 |
| **Configuration** | Simple | More complex | MMA2 |
| **Determinism** | Guaranteed | Guaranteed | Tie |
| **Performance** | Excellent | Excellent | Tie |
| **Flexibility** | Limited | High | RMA |
| **Edge Cases** | Limited | Good | RMA |

---

## 10. Architectural Conflicts

### Potential Issues with MMA2

1. **Fixed memory types** - Cannot represent non-standard data directly
2. **16-bit only** - Requires encoding for 32-bit values
3. **No native variable-width support** - Must work within constraints

### Potential Issues with RMA

1. **No native Modbus** - Requires adapter layer
2. **Arbitrary widths** - Can lead to over-engineering
3. **Ecosystem mismatch** - Doesn't fit Atlas standards
4. **Higher complexity** - More error cases, larger API surface

### Conflict Resolution

**The conflicts identified for MMA2 (fixed types, 16-bit limitation) are resolved by the Simulation Runtime architecture:**

```
Simulation Runtime
├── Device Memory (internal, variable types)
├── Encoding (device responsibility)
└── MMA2 (raw Modbus, 16-bit)
```

Devices can use any internal representation. Only the MMA2 interface requires 16-bit encoding.

**The conflicts identified for RMA (ecosystem mismatch, complexity) are not resolvable without abandoning Atlas ecosystem compatibility.**

---

## 11. Recommended Direction

### Recommendation: **Continue with MMA2**

### Justification

1. **Ecosystem Compatibility:** MMA2 is the established standard. Atlas-PPC, Replicator, and other projects use MMA2. The Simulation Runtime should integrate with the existing ecosystem.

2. **Simplicity:** MMA2's fixed memory model is a feature, not a limitation. The four Modbus areas map naturally to industrial device semantics.

3. **Device Simulation Fit:** Standard industrial devices (weather stations, meters, relays, inverters, PLCs) all use Modbus semantics. MMA2 is purpose-built for this domain.

4. **Encoding Separation:** The Simulation Runtime architecture already handles encoding at the device layer. Devices can use any internal representation and encode to 16-bit for MMA2.

5. **Lower Risk:** Using the established MMA2 reduces integration risk and maintenance burden.

### When to Reconsider RMA

Consider RMA only if:
- The simulation domain expands to include non-Modbus devices
- Ecosystem integration requirements change
- A specific use case requires variable-width memory that cannot be encoded into 16-bit registers

### Hybrid Architecture (Not Recommended)

A hybrid approach (MMA2 for Modbus devices, RMA for others) would:
- Add complexity
- Create two integration points
- Increase maintenance burden
- Provide no clear benefit for the current use case

**Recommendation:** Reject hybrid architecture. Standardize on MMA2.

---

## 12. Summary

| Question | Answer |
|----------|--------|
| Which model is simpler? | MMA2 |
| Which model fits industrial devices? | MMA2 |
| Which model integrates with Atlas ecosystem? | MMA2 |
| Which model is more flexible? | RMA |
| Which model should the Simulation Runtime use? | **MMA2** |

### Final Verdict

**MMA2 is the correct choice for the Simulation Runtime.**

MMA2's fixed Modbus memory model is purpose-built for industrial device simulation. RMA's flexibility (variable-width memory) addresses edge cases that don't exist in the standard industrial device domain. The Simulation Runtime architecture already handles encoding at the device layer, making MMA2's 16-bit limitation irrelevant.

The Simulation Runtime should:
1. Continue using MMA2 as the operational memory integration point
2. Implement Raw Ingest for publishing simulation data
3. Let devices handle encoding from internal representations to 16-bit registers
4. Focus engineering effort on device simulation, not memory infrastructure

---

*Document Version: 1.0*
*Date: 2026-07-08*
*Analysis Basis: MMA2 and RMA public documentation*
