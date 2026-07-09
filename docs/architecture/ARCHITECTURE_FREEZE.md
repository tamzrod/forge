# Architecture Freeze

## Purpose

This document establishes the architectural laws and freezes the architecture from further redesign.

The architecture has reached a stable state. Future work prioritizes implementation over speculation.

## Why the Architecture Was Frozen

The architecture was refined through multiple iterations until it satisfied these criteria:

1. **Simple** - Few concepts, clear responsibilities
2. **Memory-centric** - Memory as single source of truth
3. **Device-owned** - Everything belongs to devices
4. **Runtime-minimal** - Infrastructure only, no business logic
5. **Deterministic** - Reproducible execution
6. **Extensible** - New domains through plugins, no runtime changes

The architecture is now considered a contract for implementation.

## Architectural Laws

These principles are no longer under discussion.

### 1. Device Definition

A virtual industrial device is fundamentally:
- Deterministic memory
- Executable behaviors
- Protocol interfaces

### 2. Memory as Source of Truth

Memory is the single source of truth. There is no state outside memory.

### 3. Behaviors Modify Memory

Behaviors read from and write to device memory. Behaviors never own state.

### 4. Protocols Expose Memory

Protocols expose device memory to external systems. Protocols never own state.

### 5. Protocols Never Own State

Protocols are external views. They never cache, synchronize, or maintain state.

### 6. Devices Never Communicate Directly

Devices communicate only by reading and writing memory. There are no direct device references, message buses, callbacks, or service calls between devices.

### 7. Runtime Provides Infrastructure

The runtime hosts devices. It provides scheduling, time advancement, and plugin loading. It contains no business logic.

### 8. Plugins Provide Domain Knowledge

Plugins provide device types. The runtime remains domain-independent. New domains require no runtime changes.

## Ownership Rules

```
Device owns:
├── Memory Image
├── Behaviors
├── Protocols
└── Faults

Runtime owns:
├── Scheduler
├── Simulation Clock
├── Device Registry
└── Plugin Loader
```

### What Devices Own

- Memory
- Behaviors
- Protocols
- Faults

### What Devices Do Not Own

- Scheduling
- Time management
- Plugin loading

### What the Runtime Owns

- Scheduling
- Time advancement
- Device lifecycle
- Plugin loading

### What the Runtime Does Not Own

- Memory
- Behaviors
- Protocols
- Domain logic

## When Architecture May Be Revisited

Future architectural modifications require evidence, not speculation.

### Valid Reasons to Revisit Architecture

1. **Repeated code duplication** - The same pattern appears in multiple plugins without a clean solution
2. **Performance bottlenecks** - The architecture prevents meeting measurable performance requirements
3. **Ownership confusion** - Unclear where a responsibility belongs despite clear rules
4. **Inability to support a real use case** - A legitimate use case cannot be modeled
5. **Architectural contradiction** - Implementation reveals an inherent conflict in the principles

### Invalid Reasons to Revisit Architecture

- Personal preference
- "This could be cleaner"
- "I would have done it differently"
- Hypothetical future requirements
- Speculation about scale

## Implementation Rule

**Implementation adapts to architecture. Architecture does not adapt to implementation.**

If implementation is difficult, consider:
1. Am I implementing correctly?
2. Is the problem in my code, not the architecture?
3. Can this be solved with a plugin, not a runtime change?

Only genuine architectural limitations justify revisiting the architecture.

## Future Work

Priority order:

1. **Implement one complete device** - A working revenue meter with behaviors and protocols
2. **Validate execution model** - Does the tick loop work as specified?
3. **Validate protocol adapters** - Can Modbus TCP expose device memory?
4. **Measure performance** - What are actual tick times?
5. **Improve implementation** - Refactor based on working code

Avoid architecture discussions unless implementation demonstrates a real problem.

## Project Principles

- **Prefer implementation over speculation** - Write code that works
- **Prefer measured evidence over assumptions** - Profile before optimizing
- **Prefer small evolutionary improvements over rewrites** - Incrementally improve working code

## Summary

The architecture is complete. It is a contract for implementation.

The goal is now to prove the architecture through working software rather than continuing to redesign it.

---

**Last Updated:** 2024
**Status:** Frozen
**Rationale:** Architecture satisfies simplicity, memory-centricity, device-ownership, and extensibility requirements.
