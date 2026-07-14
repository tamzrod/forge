# UX Decisions

**Document Status:** Frozen  
**Last Updated:** 2026-07-14  
**Vision:** Forge is the best engineering workbench for understanding and operating a utility-scale solar farm.

---

## Decision Criteria

All UX recommendations are evaluated against one criterion:

> **Does this move Forge closer to becoming the best engineering workbench for understanding and operating a utility-scale solar farm?**

Comparison is against:
- Forge's vision
- Real engineering workflows
- Operational requirements of solar farm operators

Comparison excludes:
- SCADA systems (different purpose)
- CAD tools (different purpose)
- Modern web application conventions (different domain)

---

## Design Constraints

These constraints define the boundaries of acceptable implementation:

### Plant Health Dashboard

**Constraint:** Must remain a high-level summary and must not evolve into a SCADA dashboard.

**Rationale:** Forge is an engineering workbench for understanding, not a real-time control system. SCADA dashboards belong in control rooms with 24/7 operators. Forge should surface insights and KPIs, not replicate the telemetry density of operational control systems.

**Boundaries:**
- Shows aggregates, deviations, and anomalies
- Does not show every sensor reading
- Does not support continuous real-time updates at sub-second intervals
- Does not replace operator station displays

### Power Flow Animation

**Constraint:** Must communicate engineering state, not provide visual decoration.

**Rationale:** Animation in Forge must serve diagnostic and educational purposes. Decorative animation distracts from understanding and can obscure important state changes.

**Requirements:**
- Animation speed reflects actual power flow magnitude
- Color coding reflects state (healthy/warning/fault)
- Animation pauses or dims during faults to draw attention
- No gratuitous motion that does not convey information

---

## Decisions

### P0 - Critical

| # | Recommendation | Decision | Reason | Impact | Priority |
|---|---------------|----------|--------|--------|----------|
| 1 | Implement IEEE Std 315 SLD symbols | **ACCEPT** | Engineering credibility requires industry-standard notation. Operators trained on IEEE symbols cannot work effectively with emoji representations. | High | P0 |
| 2 | Add Plant Health Dashboard | **ACCEPT** | Reduces operational cognitive load by surfacing plant-level KPIs (generation, capacity factor, MWh) in a single view. Directly serves the "operating" aspect of the vision. | Very High | P0 |
| 3 | Add capacity factor, daily MWh, PR metrics | **ACCEPT** | These are the primary KPIs solar farm operators manage. Without them, Forge cannot support real operational decision-making. | High | P0 |
| 4 | Add equipment designations (CB-101, TX-1) | **ACCEPT** | Real plants use designations for communication, maintenance, and fault isolation. Nameless equipment undermines operational credibility. | High | P0 |

### P1 - High

| # | Recommendation | Decision | Reason | Impact | Priority |
|---|---------------|----------|--------|--------|----------|
| 5 | Add power flow animation on SLD | **ACCEPT** | Static diagrams don't convey real-time state. Animation makes the simulation feel alive and helps operators understand power direction at a glance. | High | P1 |
| 6 | Replace emoji with industrial icons | **ACCEPT** | Emoji belong in consumer apps, not engineering software. Industrial icons convey meaning faster and with less ambiguity. | Medium | P1 |
| 7 | Add plant name header | **ACCEPT** | Operators work on named assets. "Utility-Scale Solar Farm" is a category, not an identity. Multi-plant deployments need differentiation. | Medium | P1 |
| 8 | Merge Equipment Details into collapsible panel | **ACCEPT** | Tab fragmentation (5 tabs) creates friction. A collapsible bottom panel keeps the SLD primary while providing detail on demand. | High | P1 |
| 9 | Add inverter clipping modeling | **ACCEPT** | Inverter clipping is a fundamental physical phenomenon in solar farms. Without it, the simulation produces incorrect power curves at high irradiance. | High | P1 |
| 9a | Equipment Operational State | **ACCEPT** | Operators need immediate visual feedback on equipment health. Healthy/Warning/Fault/Offline states are the foundation of operational awareness. | High | P1 |
| 9b | Cause Chain Visualization | **ACCEPT** | Solar farm operations require tracing effects upstream and downstream. Sun → Weather → PV → Inverter → Transformer → PCC forms the core diagnostic path. | High | P1 |

### P2 - Medium

| # | Recommendation | Decision | Reason | Impact | Priority |
|---|---------------|----------|--------|--------|----------|
| 10 | Add alarm/event persistent log | **ACCEPT** | Operators need event history to diagnose problems and maintain records. Current transient events don't serve this need. | High | P2 |
| 11 | Remove "Properties" tab from equipment | **DEFER** | Properties may be useful for training scenarios where users configure equipment. Reserve judgment until training mode is designed. | Low | P2 |
| 12 | Add actual vs. expected power to each PV | **ACCEPT** | Immediate deviation visibility is essential for fault detection. Operators should see expected output without manual calculation. | High | P2 |
| 13 | Replace "Explain" with "Performance" | **REJECT** | Explainability is a first-class Forge feature. The "Why?" panel demonstrates commitment to understanding over simple status display. Education and operations are complementary, not competing. | Medium | P2 |
| 14 | Add simulation clock as time-of-day | **ACCEPT** | Operators think in operational hours (ramp-up, peak, curtailment windows). Wall-clock HH:MM:SS is meaningless for solar operations. | Medium | P2 |

### P3 - Nice to Have

| # | Recommendation | Decision | Reason | Impact | Priority |
|---|---------------|----------|--------|--------|----------|
| 15 | Add reactive power display | **ACCEPT** | Grid operators care about VARs as much as MW. Reactive power is essential for understanding grid interactions. | Medium | P3 |
| 16 | Add curtailment modeling | **ACCEPT** | Curtailment represents real business impact in solar operations. Ignoring it creates a false picture of plant performance. | Medium | P3 |
| 17 | Add breaker open/close actions | **ACCEPT** | Control simulation is part of operating a plant. Users should experience the effect of isolation and re-energization. | Medium | P3 |
| 18 | Add SCADA polling visualization | **DEFER** | Advanced feature that requires Modbus/SCADA integration. Not core to the initial vision of understanding and operating a solar farm. | Low | P3 |

---

## Summary

| Decision | Count |
|----------|-------|
| **ACCEPT** | 18 |
| **REJECT** | 1 |
| **DEFER** | 2 |

### Rejection Rationale

1. **Replace "Explain" with "Performance" (REJECT)**: Explainability is a first-class Forge feature. The "Why?" panel demonstrates commitment to understanding over simple status display. Education and operations are complementary.

### Deferral Rationale

1. **Remove "Properties" tab (DEFER)**: May serve training/design mode where configuration is the learning objective.

2. **SCADA polling visualization (DEFER)**: Requires external integration complexity that exceeds the core vision scope.

---

## Decision Log

| Date | Decision | Rationale |
|------|----------|-----------|
| 2026-07-14 | Accept IEEE Std 315 symbols | Engineering credibility |
| 2026-07-14 | Accept Plant Health Dashboard | Reduces cognitive load, surfaces KPIs |
| 2026-07-14 | Accept capacity factor, daily MWh, PR | Primary operational KPIs |
| 2026-07-14 | Accept equipment designations | Operational communication requirement |
| 2026-07-14 | Accept power flow animation | Real-time state visualization |
| 2026-07-14 | Accept industrial icons | Credibility over aesthetics |
| 2026-07-14 | Accept plant name header | Basic operational context |
| 2026-07-14 | Accept collapsible Equipment Details | Reduces navigation friction |
| 2026-07-14 | Accept inverter clipping | Fundamental physics |
| 2026-07-14 | Accept Equipment Operational State | Foundation of operational awareness |
| 2026-07-14 | Accept Cause Chain Visualization | Core diagnostic path |
| 2026-07-14 | Accept persistent alarm log | Operational record-keeping |
| 2026-07-14 | Defer Properties tab removal | May serve training mode |
| 2026-07-14 | Accept actual vs. expected power | Immediate deviation visibility |
| 2026-07-14 | Reject Performance over Explain | Explainability is first-class |
| 2026-07-14 | Accept time-of-day clock | Operational context |
| 2026-07-14 | Accept reactive power display | Grid operator requirement |
| 2026-07-14 | Accept curtailment modeling | Real business impact |
| 2026-07-14 | Accept breaker control actions | Control simulation |
| 2026-07-14 | Defer SCADA visualization | Outside core vision scope |
| 2026-07-14 | Add Design Constraint: Plant Health Dashboard | Prevents SCADA feature creep |
| 2026-07-14 | Add Design Constraint: Power Flow Animation | Enforces engineering purpose |
| 2026-07-14 | Freeze UX_DECISIONS.md | Ready for implementation |

---

*This document is the authoritative UX decision log for Forge.*
