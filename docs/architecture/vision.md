# Vision

## What This Project Is

**A Virtual Industrial Laboratory for industrial software development.**

The simulator provides a realistic industrial environment for:

- Software development
- Controller development
- SCADA development
- Protocol integration
- Factory Acceptance Testing (FAT)
- Commissioning
- Training
- Education
- Demonstrations

Simulation is the enabling technology. The laboratory is the product.

---

## What Problems It Solves

Industrial software is difficult to develop and test without physical hardware.

Real hardware is:
- Expensive to acquire
- Difficult to set up
- Risky to experiment with
- Unavailable during commissioning delays
- Not portable between locations

This project provides a deterministic virtual industrial environment where industrial software behaves the same way it would when connected to a real installation.

---

## Who It's For

| User | Use Case |
|------|----------|
| **Software Developers** | Develop and test industrial applications without hardware |
| **Control Engineers** | Validate controller logic before deployment |
| **SCADA Engineers** | Configure and test human-machine interfaces |
| **Integration Teams** | Verify protocol implementations |
| **Commissioning Teams** | Prepare and practice before site visits |
| **Training Teams** | Train operators in a safe, repeatable environment |
| **Technical Educators** | Teach industrial automation concepts |

---

## What It Deliberately Does Not Attempt

This project is **not** intended to become:

| Not A | Why Not |
|--------|---------|
| Power system analysis package | Focus is on software behavior, not power flow studies |
| Electromagnetic transient simulator | Not needed for industrial software development |
| Finite element solver | Out of scope for software testing |
| CFD package | Not relevant to industrial protocols |
| Generic physics engine | Domain-specific models are sufficient |
| Digital twin platform | Focus is on software integration, not plant fidelity |

These may inspire future plugins but are not the mission of the Runtime.

---

## Design Philosophy

### Believe Before Sophisticate

Models should be **credible** before they become **sophisticated**.

Simple deterministic models are preferred over highly accurate but complex models unless additional fidelity clearly benefits industrial software development.

### Fitness for Purpose

Every feature should be evaluated against this question:

> *"Does this improve the ability to develop, test, commission, or train industrial software?"*

If the answer is no, it should probably not be part of the Runtime.

### Deterministic Execution

The simulation must be deterministic. Same inputs produce same outputs, every time.

This enables:
- Reproducible test results
- Deterministic CI/CD pipelines
- Reliable debugging
- Repeatable training scenarios

---

## Architectural Principles

### The Laboratory Metaphor

Think of the Runtime as a complete industrial plant in software:

```
┌─────────────────────────────────────────────────────────────────┐
│                    Virtual Industrial Laboratory                     │
│                                                                 │
│  Simulation Models ←→ Virtual Devices ←→ MMA2 ←→ Applications  │
│                                                                 │
│  Physics            Equipment         Telemetry      Control SW   │
└─────────────────────────────────────────────────────────────────┘
```

- **Simulation Models** represent the physical world (grid voltage, sun irradiance, etc.)
- **Virtual Devices** represent industrial equipment (meters, inverters, PLCs)
- **MMA2** represents the shared operational view of the plant
- **Applications** are the software under development/test

### Key Constraints

1. **Simulation Models never expose protocols** - Physics doesn't have Modbus
2. **Virtual Devices never access other devices** - Equipment doesn't talk directly
3. **Applications never access internal state** - Real systems use published telemetry
4. **MMA2 is the boundary** - Everything external connects through operational memory

---

## Success Criteria

A successful simulation means:

| Criterion | Meaning |
|-----------|---------|
| **Deterministic** | Same scenario produces same results every time |
| **Believable** | Industrial software behaves identically to real deployment |
| **Portable** | Runs anywhere without physical hardware |
| **Sustainable** | Simple enough to maintain and extend |
| **Focused** | Addresses industrial software development, not physics accuracy |

---

## The Message

> **This project provides a deterministic Virtual Industrial Laboratory for developing, testing, commissioning, and training industrial software through realistic virtual industrial environments.**

---

*Last Updated: 2026-07-09*
