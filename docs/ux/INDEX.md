# UX Design Package: Solar Farm Reference World Engineering Workbench

**Package ID:** UX-PKG-SOLAR-002  
**Version:** 2.0 (Optimized)  
**Date:** 2026-07-13  
**Status:** Ready for Implementation  

---

## Package Overview

This is an **optimized** design package focused solely on the Solar Farm Reference World MVP. The goal is to provide the shortest path from launching Forge to understanding how a solar farm works.

### Key Optimization Principle

> **Forge's first deliverable is a realistic Solar Farm Reference World. Everything else is deferred to post-MVP.**

---

## Package Contents

| Document | Purpose | Version | Status |
|----------|---------|---------|--------|
| **OPTIMIZATION_AUDIT.md** | Complete audit with rationale for all changes | 2.0 | Complete |
| **UX_SPECIFICATION_V2.md** | Optimized UX specification | 2.0 | Complete |
| **COMPONENT_WIREFRAMES_V2.md** | Simplified wireframes | 2.0 | Complete |
| **INDEX.md** | This file | 2.0 | Complete |

---

## Quick Start: Understanding the UX

### The Optimized Workflow

```
LAUNCH FORGE
     ↓
LOAD SOLAR FARM REFERENCE WORLD (PRIMARY PATH)
     ↓
RUN SIMULATION
     ↓
UNDERSTAND PLANT BEHAVIOR (Why? tab)
     ↓
MODIFY EQUIPMENT
     ↓
OBSERVE RESULTS
     ↓
SAVE PROJECT
```

### What Changed from V1?

| Metric | V1 | V2 | Change |
|--------|-----|-----|--------|
| Screens | 8 | 4 | -50% |
| Menu Items | 42 | 16 | -62% |
| Palette Items | 23 | 12 | -48% |
| Keyboard Shortcuts | 31 | 14 | -55% |
| Debug Concepts | 5 | 1 (Why?) | -80% |

---

## Document Summary

### 1. OPTIMIZATION_AUDIT.md

Comprehensive audit of the V1 UX specification with rationale for all changes.

**Contains:**
- Phase 1: Mission Alignment Audit (screens, panels, workflows, menus)
- Phase 2: Workflow Optimization (Reference World as primary path)
- Phase 3: Screen Reduction (merged Timeline/Events/Trace/Watch → Analysis)
- Phase 4: Equipment Palette Optimization (reduced from 23 to 12 items)
- Phase 5: Project Explorer Redesign (World → Plant hierarchy)
- Phase 6: Engineering Explainability (Why? replaces Debug)
- Phase 7: Learning Optimization (Observe First principle)
- Phase 8: UX Simplification (menus, panels, shortcuts)
- MVP Scope Definition
- Deferred Features List
- Final Recommendations

### 2. UX_SPECIFICATION_V2.md

The optimized UX specification for implementation.

**Key Sections:**
1. Mission Statement (Solar Farm Reference World only)
2. User Profile (Learning-focused engineer)
3. Core Principle: Observe First
4. Welcome Screen (Reference World prominent)
5. Main Editor Layout (simplified)
6. Equipment Palette (12 items only)
7. Project Explorer (Plant-focused)
8. Analysis Panel (merged Timeline/Events/Why?)
9. Inspector (with Why? explainability)
10. Simplified Menus (16 items)
11. Keyboard Shortcuts (14 items)
12. User Journey (optimized for first-time)
13. MVP Scope Definition

### 3. COMPONENT_WIREFRAMES_V2.md

Simplified wireframes for implementation.

**Contains:**
- Welcome Screen
- Reference World Selector
- Main Editor Layout
- Equipment Palette (Plant + Environment)
- Canvas Equipment (PV Block, Meter, Breaker)
- Project Explorer (Plant hierarchy)
- Analysis Panel (Timeline/Events/Why? tabs)
- Simplified Menus
- Inspector States

---

## Key Design Decisions (V2)

### 1. Reference World is PRIMARY

```
V1:  [New Project]  [Open]  [Load Reference World]
      (equal weight)

V2:  [LOAD SOLAR FARM REFERENCE WORLD] (prominent)
      [Create New Project] (subdued, "for experts")
```

### 2. Engineering Explainability replaces Debug

```
V1:  [Signal Trace] [Data Watch] [Breakpoints] [Memory View]

V2:  [Why?] - Click any value to see why it is what it is
```

**Example:**
```
WHY IS PCC EXPORT = 48.95 MW?

Because:
  PV Block 1:    4.90 MW
  PV Block 2:    4.90 MW
  ...
  Total Gross:  49.00 MW
  Auxiliary:    -0.05 MW
  ─────────────────────────
  Net Export:   48.95 MW
```

### 3. Analysis Panel (Merged)

Single panel with tabs instead of separate screens:

```
┌─────────────────────────────────────────────────────┐
│ ANALYSIS                                            │
├─────────────────────────────────────────────────────┤
│  [Timeline] [Events] [Why?]                        │
│  ────────────────────────────────────────────────   │
│                                                     │
│  (Content changes based on selected tab)           │
│                                                     │
└─────────────────────────────────────────────────────┘
```

### 4. Simplified Equipment Palette

**V1:** 23 items across 6 categories (Substation, Collection, Generation, Protection, Environment, Simulation)

**V2:** 12 items across 2 categories (Plant, Environment)

**NEW CONCEPT: PV Block**
A self-contained generation unit that includes PV Array, Combiner, Inverter, and Transformer internally.

### 5. Plant-Focused Explorer

```
V1:  Project → World → Topology → Entities → Simulation

V2:  Solar Farm Project → Plant → Grid Interconnection,
                                  Collection System,
                                  Auxiliary Systems,
                                  Environment
```

---

## MVP Scope

### IN SCOPE (Implement)

| Feature | Priority | Notes |
|---------|----------|-------|
| Welcome Screen | High | Entry point |
| Reference World Selection | High | Primary workflow |
| Main Editor | High | Core workspace |
| Equipment Palette (12 items) | High | Solar farm only |
| Canvas with Single-Line Diagram | High | Primary visualization |
| Inspector Panel | High | Equipment details |
| Project Explorer | High | Plant hierarchy |
| Analysis Panel (Timeline/Events/Why?) | High | Understanding |
| Simulation Controls | High | Run/Pause/Stop/Reset |
| Engineering Explainability | High | Learning focus |
| Project Save/Load | Medium | Persistence |
| Real-time Measurements | Medium | Observation |

### OUT OF SCOPE (Post-MVP)

| Feature | Priority | Rationale |
|---------|----------|-----------|
| Empty Project Creation | Low | Reference World preferred |
| Export functionality | Low | Not core to learning |
| Undo/Redo | Medium | Complexity vs. benefit |
| Step-by-Step Execution | Low | Not needed |
| Speed Control | Low | Default works |
| Zoom Controls | Low | Basic zoom sufficient |
| Advanced Protection | Medium | Post-solar basics |
| Battery Storage | Medium | Post-solar basics |
| IEC 61850 Support | Low | Post-MVP |
| Multi-User Collaboration | Low | Post-MVP |
| Report Generation | Low | Post-MVP |

---

## File Structure

```
docs/ux/
├── INDEX.md                        # This file (overview)
├── OPTIMIZATION_AUDIT.md          # Complete audit with rationale
├── UX_SPECIFICATION_V2.md          # Optimized UX specification
├── COMPONENT_WIREFRAMES_V2.md      # Simplified wireframes
├── UX_SPECIFICATION.md             # Original V1 (archived)
└── COMPONENT_WIREFRAMES.md         # Original V1 (archived)
```

---

## Review Checklist

Before implementation, verify:

- [ ] UX aligns with "Observe First" learning principle
- [ ] Reference World is the primary path from Welcome
- [ ] Engineering Explainability (Why?) replaces debugging
- [ ] Equipment palette is scoped to solar farm only
- [ ] Project Explorer reflects engineering mental model
- [ ] Menus and shortcuts are simplified
- [ ] MVP scope is clear and achievable

---

## Next Steps

1. **Review**: Present optimization audit to stakeholders
2. **Approve**: Get sign-off on V2 approach
3. **Plan**: Create implementation plan based on MVP scope
4. **Implement**: Begin with Welcome Screen and Reference World
5. **Iterate**: Build iteratively, testing with users early

---

*Package Version: 2.0 (Optimized)*  
*Status: Ready for Implementation Planning*
*Mission: Solar Farm Reference World only*
