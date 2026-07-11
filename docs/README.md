# Forge Documentation

> Virtual Industrial Laboratory for industrial software development.

---

## Start Here

| Document | Purpose |
|----------|---------|
| **[Architecture Glossary](architecture/GLOSSARY.md)** | **Single source of truth for project terminology.** Read this first to understand the project's language. |
| [Vision](architecture/vision.md) | Project purpose, audience, and philosophy |
| [Architecture Overview](architecture/overview.md) | System architecture and design principles |

---

## Architecture

### Core Concepts

| Document | Description |
|----------|-------------|
| **[Architecture Glossary](architecture/GLOSSARY.md)** | Authoritative terminology definitions |
| [Overview](architecture/overview.md) | System architecture, layers, and data flow |
| [Vision](architecture/vision.md) | Project purpose and philosophy |

### Architecture Details

| Document | Description |
|----------|-------------|
| [Simulation Models](architecture/simulation-models.md) | Physical world representation |
| [Runtime](architecture/runtime.md) | Runtime hosting and coordination |
| [Execution Model](architecture/execution-model.md) | End-to-end execution flow |

### Device Architecture

| Document | Description |
|----------|-------------|
| [Device Model](architecture/device-model.md) | Virtual device structure |
| [Memory Model](architecture/memory-model.md) | Device memory ownership |
| [Behavior Model](architecture/behavior-model.md) | Device-owned logic |

### Communication

| Document | Description |
|----------|-------------|
| [Protocol Architecture](architecture/protocol-architecture.md) | External memory views |
| [MMA2 Integration](architecture/mma2-integration.md) | MMA2 operational telemetry |

---

## Development

| Document | Description |
|----------|-------------|
| [Design Principles](development/design-principles.md) | Development philosophy |
| [Coding Rules](development/coding-rules.md) | Code style and conventions |
| [Roadmap](development/roadmap.md) | Project milestones |

---

## Examples

| Document | Description |
|----------|-------------|
| [End-to-End Example](examples/end-to-end.md) | Complete architecture validation |

---

## Quick Reference

### Architecture Layers

```
Simulation Runtime (hosts, schedules)
        ↓
Simulation Models (physics)
        ↓
Virtual Firmware (samples, owns memory)
        ↓
Communication Interfaces (serialize)
        ↓
MMA2 / SCADA / External Systems
```

### Key Terms

| Term | Definition |
|------|------------|
| **Simulation Model** | Represents physical world (Weather, Grid, Sun) |
| **Virtual Firmware** | Software running inside a Virtual Device |
| **Device Memory** | Internal RAM owned by firmware |
| **Communication Interface** | Serializes Device Memory for external systems |

See [Architecture Glossary](architecture/GLOSSARY.md) for complete terminology.

---

## Contributing

When adding new concepts:

1. Check [Architecture Glossary](architecture/GLOSSARY.md) for existing terms
2. Reuse existing terminology whenever possible
3. If a new concept is required, add it to the glossary
4. Update documentation to reference the glossary

---

*For questions about architecture, see the [Architecture Glossary](architecture/GLOSSARY.md).*
