# Infrastructure Model

> **Note:** This document has been superseded by [Simulation Models](simulation-models.md).
> The concept of "infrastructure" has been renamed to "Simulation Models" to better reflect
> that these components represent the physical world (physics), not infrastructure services.

## Migration Guide

If you have references to "infrastructure" in your code, update them to use "models":

| Old Concept | New Concept |
|-------------|-------------|
| `Infrastructure` | `SimulationModels` |
| `runtime.Infrastructure()` | `runtime.Models()` |
| `runtime.CreateInfrastructure()` | `runtime.CreateModel()` |
| `behavior.Infrastructure()` | `behavior.device.Model()` |

## Summary

Please refer to [Simulation Models](simulation-models.md) for the complete documentation of:

- Why Simulation Models exist
- Difference between Models and Devices
- Why Models are not exposed through protocols
- Why Models are not published to MMA2
- Why Models store state directly in RAM
- Model API and types
- Interaction patterns
- Execution order

---

*Last Updated: 2026-07-09*
*Superseded by: [Simulation Models](simulation-models.md)*
