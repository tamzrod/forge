# Design Principles

## Core Principles

### 1. Devices Are the System

The runtime hosts devices. Devices own memory, behaviors, and protocols.

### 2. Memory Is the Source of Truth

Everything reads from and writes to device memory. Nothing owns state outside memory.

### 3. Behaviors Own Logic

Behaviors are logic that reads and writes device memory. Behaviors never access other devices.

### 4. Protocols Are External Views

Protocols are external interfaces that expose device memory. Protocols never own state.

### 5. Devices Never Communicate Directly

Devices communicate only through the runtime's execution order. Behaviors read and write memory.

### 6. Runtime Is Infrastructure

The runtime provides hosting, scheduling, and plugin loading. It contains no domain knowledge.

### 7. Plugins Provide Domain Knowledge

New domains add device types through plugins. The runtime never changes.

## Ownership Model

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

## What Belongs Where

### Device

**Yes:**
- Memory
- Behaviors
- Protocols
- Faults

**No:**
- Scheduling
- Time management
- Plugin loading

### Runtime

**Yes:**
- Scheduling
- Time advancement
- Device lifecycle
- Plugin loading

**No:**
- Memory
- Behaviors
- Protocols
- Domain logic

## Anti-Patterns

### God Runtime

```go
// Bad
type BadRuntime struct {
    memoryManager
    behaviorExecutor
    protocolHost
    faultManager
    // ...
}
```

### Global Memory

```go
// Bad
runtime.GlobalMemory()
```

### Protocol Synchronization

```go
// Bad
protocolA.SyncWith(protocolB)
```

### Behavior Coupling

```go
// Bad
func (b *Behavior) AccessOtherDevice(other *Device) {
    other.Memory().Write(...)
}
```

## Determinism

Simulation is deterministic:

1. Devices tick in registration order
2. Behaviors tick in registration order
3. Same memory → same results
4. No randomness without seeded RNG

## Simplicity

Keep the architecture small:

- Runtime is intentionally minimal
- Devices are self-contained
- Memory is the only state
- No unnecessary layers

## Future Domains

Adding a new domain:

1. Create a new plugin
2. Define device types
3. Implement behaviors
4. No runtime changes required

The architecture naturally supports any industrial domain.
