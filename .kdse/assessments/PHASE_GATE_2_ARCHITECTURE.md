# Phase Gate 2 — Architecture Report

**Date:** 2026-07-11  
**Session:** KDSE-20260711-230631  
**Feature:** KDSE Debug Runtime

---

## Executive Summary

This document defines the architecture for the KDSE Debug Runtime—a deterministic engineering debugging system that transforms debugging from endless hypothesis generation into evidence-driven root cause analysis.

---

## Design Decisions

### 1. Evidence-First Investigation

All hypotheses MUST be grounded in collected evidence. No speculation without evidence.

**Rationale:** Prevents hypothesis generation without basis, ensuring every debugging path is evidence-driven.

### 2. Confidence-Driven Selection

Root cause selection occurs when confidence exceeds a configurable threshold (default: 90%).

**Rationale:** Provides a clear stopping criterion, preventing endless investigation loops.

### 3. Loop Prevention

Detect and prevent repeated investigation of the same component.

**Rationale:** Eliminates redundant work and wasted tool calls.

### 4. Structured Artifacts

All debugging artifacts are stored as structured data (JSON/YAML).

**Rationale:** Ensures traceability, reproducibility, and machine-readable output.

### 5. Reusable Engine

The Debug Engine is a core runtime component reused by multiple commands.

**Rationale:** Avoids code duplication and ensures consistent debugging behavior.

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                    Debug Engine (Core)                          │
│                                                                 │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐               │
│  │  Evidence  │  │ Hypothesis│  │ Confidence │               │
│  │  Collector │  │  Manager  │  │  Tracker   │               │
│  └────────────┘  └────────────┘  └────────────┘               │
│                                                                 │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐               │
│  │    Loop   │  │   Root    │  │  Report    │               │
│  │  Detector │  │   Cause   │  │  Generator │               │
│  └────────────┘  └────────────┘  └────────────┘               │
└─────────────────────────────────────────────────────────────────┘
```

---

## Components

### 1. Evidence Collector

Collects and categorizes evidence:

| Type | Description | Confidence Impact |
|------|-------------|------------------|
| `exception` | Exception messages and traces | +20% |
| `test_failure` | Test failure output | +15% |
| `log` | Runtime logs | +10% |
| `source` | Source code inspection | +10% |
| `state` | Repository state | +5% |
| `config` | Configuration values | +5% |
| `dependency` | Dependency versions | +5% |

### 2. Hypothesis Manager

Manages hypotheses with:
- ID, description, status
- Supporting/contradicting evidence
- Confidence score
- Affected components
- Experiment history

### 3. Confidence Tracker

Calculates confidence:
- Initial: 20-60% (configurable)
- Evidence increases confidence
- Contradicting evidence decreases confidence (-25%)
- Stop at 90% threshold

### 4. Loop Detector

Detects repeated investigation:
- File inspection patterns
- Module reload patterns
- Schema check patterns
- Database inspection patterns

### 5. Root Cause Selector

Selects root cause when:
- Confidence >= 90%
- Evidence evaluated
- Operator approval received

### 6. Report Generator

Produces structured Debug Session Reports including:
- Failure summary
- Evidence collected
- Hypotheses generated
- Confidence progression
- Implementation details
- Verification results

---

## Debug Workflow

```
Failure Detected
       │
       ▼
Evidence Collection
       │
       ▼
Hypothesis Generation
       │
       ▼
Evidence Evaluation
       │
       ▼
Confidence Assessment ──< 90%? ──> YES
       │                            │
       NO                           ▼
       │                     Root Cause Selection
       │                            │
       └────────────────────────────┘
                                    │
                             Operator Approval
                                    │
                                    ▼
                              Implementation
                                    │
                                    ▼
                              Verification
                                    │
                                    ▼
                           Regression Tests
                                    │
                                    ▼
                            Runtime Report
```

---

## Directory Structure

```
.dkdse/
├── debug/
│   ├── engine.sh              # Core debug engine
│   ├── evidence/
│   │   ├── collector.sh       # Evidence collection
│   │   └── store.json         # Evidence store
│   ├── hypotheses/
│   │   ├── manager.sh         # Hypothesis management
│   │   └── registry.json      # Hypothesis registry
│   ├── confidence/
│   │   └── tracker.sh         # Confidence calculation
│   ├── loops/
│   │   ├── detector.sh        # Loop detection
│   │   └── history.json       # Loop history
│   ├── reports/
│   │   ├── generator.sh       # Report generation
│   │   └── root-cause/        # Root cause reports
│   └── sessions/
│       └── DEBUG-*.json       # Debug session files
└── bootstrap/
    └── debug-config.yaml      # Debug configuration
```

---

## Command Integration

| Command | Purpose | Uses Debug Engine |
|---------|---------|------------------|
| `kdse debug` | Primary debugging | ✓ Full workflow |
| `kdse verify` | Verification failures | ✓ Full workflow |
| `kdse repair` | Automated fixes | ✓ Full workflow |
| `kdse audit` | Audit findings | ✓ Full workflow |

---

## Configuration

```yaml
debug_runtime:
  version: "1.0"
  
confidence:
  threshold: 90              # Stop at 90%
  minimum_initial: 20         # Min starting confidence
  maximum_initial: 60         # Max starting confidence

hypothesis:
  max_active: 5              # Limit active hypotheses
  auto_reject_below: 20     # Auto-reject low confidence

loops:
  detection_enabled: true
  max_repetitions: 3

implementation:
  require_operator_approval: true
  auto_backup: true
```

---

## Knowledge Manifest Update

### Required Knowledge

| # | ID | Name | Source |
|---|-----|------|--------|
| 1-7 | (existing) | Core KDSE knowledge | (existing) |

### Optional Knowledge (New)

| # | ID | Name | Source |
|---|-----|------|--------|
| 12 | `debug-runtime` | Debug Runtime | `runtime/DEBUG_RUNTIME.md` |
| 13 | `debug-adr` | Debug ADR | `.kdse/architecture/ADR-007-DEBUG_RUNTIME.md` |

---

## Capabilities Added

| Capability | Description |
|-----------|-------------|
| `debugging` | Evidence-driven debugging workflow |
| `evidence_collection` | Systematic evidence collection |
| `root_cause_analysis` | Confidence-driven root cause identification |

---

## Implementation Phases

### Phase 1: Core Engine
1. Create `runtime/debug/engine.sh`
2. Evidence collection with categorization
3. Hypothesis creation and management
4. Basic confidence tracking

### Phase 2: Intelligence
5. Loop detection
6. Evidence correlation
7. Root cause selection
8. Report generation

### Phase 3: Integration
9. `kdse debug` command
10. `kdse verify` integration
11. `kdse repair` integration
12. `kdse audit` integration

---

## Artifacts Created

| Artifact | Location | Purpose |
|----------|----------|---------|
| Architecture | `runtime/DEBUG_RUNTIME.md` | Detailed architecture specification |
| ADR | `.kdse/architecture/ADR-007-DEBUG_RUNTIME.md` | Architecture decision record |
| Knowledge Update | `.kdse/bootstrap/knowledge.yaml` | Added debug knowledge |
| Capabilities Update | `.kdse/bootstrap/capabilities.yaml` | Added debugging capabilities |
| Commands Update | `.kdse/bootstrap/commands.yaml` | Added debug commands |
| AI Context Update | `.kdse/bootstrap/kdse-ai.json` | Added debug capabilities |

---

## Success Criteria

| Criterion | Measurement |
|-----------|-------------|
| Evidence drives decisions | 0 hypotheses without evidence |
| Confidence-based selection | Root cause at >= 90% |
| Loop prevention | 0 repeated investigations |
| Operator approval | 100% of implementations approved |
| Verification required | 0 fixes without verification |

---

## Risks and Mitigations

| Risk | Mitigation |
|------|------------|
| Confidence threshold may need tuning | Make configurable |
| Loop detection may have false positives | Allow override with `--force` |
| Evidence categorization may need refinement | Support custom evidence types |

---

## Recommendation

**Proceed to Phase Gate 3: Implementation**

The architecture is complete and ready for implementation. All necessary components are defined, and the integration points are clear.

---

**Architecture Status:** ✅ Complete  
**Next Phase:** Phase Gate 3 — Implementation  
**Awaiting operator approval.**
