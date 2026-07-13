# ADR-007: KDSE Debug Runtime

**Status:** Proposed  
**Date:** 2026-07-11  
**Deciders:** KDSE Runtime Team  

---

## Context

During debugging sessions, the AI frequently enters long investigation loops:
1. Test → Hypothesis → Experiment → New Hypothesis → Experiment → Another Hypothesis → ...

Even after discovering a highly probable root cause (e.g., 92% confidence), the runtime continues searching for alternative explanations instead of fixing the issue. This wastes tool calls and delays implementation.

The KDSE Debug Runtime shall transform debugging from endless hypothesis generation into evidence-driven root cause analysis.

---

## Decision

Implement a KDSE Debug Runtime with the following characteristics:

### 1. Evidence-First Investigation

All hypotheses must be grounded in collected evidence. No speculation without evidence.

```
Evidence Collection → Hypothesis Generation → Evidence Evaluation → ...
```

### 2. Confidence-Driven Selection

Root cause selection occurs when confidence exceeds a configurable threshold (default: 90%).

```
confidence >= 90% → STOP searching → Select Root Cause → Implement
```

### 3. Loop Prevention

Detect and prevent repeated investigation of the same component.

### 4. Structured Artifacts

All debugging artifacts (evidence, hypotheses, reports) are stored as structured data.

### 5. Reusable Engine

The Debug Engine is a core runtime component reused by multiple commands:
- `kdse debug` - Primary debugging
- `kdse verify` - Verification failure investigation
- `kdse repair` - Automated repair
- `kdse audit` - Audit finding investigation

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    Debug Engine (Core)                          │
│                                                                 │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐                 │
│  │  Evidence  │  │ Hypothesis│  │ Confidence │                 │
│  │  Collector │  │  Manager  │  │  Tracker   │                 │
│  └────────────┘  └────────────┘  └────────────┘                 │
│                                                                 │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐                 │
│  │    Loop   │  │   Root    │  │  Report    │                 │
│  │  Detector │  │   Cause   │  │  Generator │                 │
│  └────────────┘  └────────────┘  └────────────┘                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## Debug Workflow

```
Failure Detected
       │
       ▼
Evidence Collection (exception, test_failure, log, source, state, config)
       │
       ▼
Hypothesis Generation (grounded in evidence)
       │
       ▼
Evidence Evaluation (supporting vs contradicting)
       │
       ▼
Confidence Assessment (evidence increases, contradictory decreases)
       │
       ├─── Confidence >= 90%? ───YES──▶ Root Cause Selection ──▶ Implementation ──▶ Verification ──▶ Regression ──▶ Report
       │                               │
       │                              NO
       │                               │
       └───────────────────────────────┘
```

---

## Data Structures

### Evidence

```yaml
evidence:
  id: "E-001"
  type: "exception|test_failure|log|source|state|config|dependency"
  timestamp: "ISO-8601"
  source: "file:line:function"
  content: "..."
  tags: ["tag1", "tag2"]
```

### Hypothesis

```yaml
hypothesis:
  id: "H-001"
  description: "..."
  status: "active|evaluated|rejected|selected"
  confidence:
    initial: 40
    current: 92
    threshold: 90
  supporting_evidence: ["E-001", "E-003"]
  contradicting_evidence: ["E-005"]
  affected_components: ["ComponentA", "ComponentB"]
  experiments: [...]
```

### Root Cause Report

```yaml
root_cause_report:
  id: "RC-001"
  session_id: "DEBUG-..."
  failure: {...}
  selected_hypothesis: "H-001"
  confidence: 92
  recommended_fix: {...}
  alternative_explanations: [...]
  operator_approved: true
```

---

## Configuration

```yaml
debug_runtime:
  confidence:
    threshold: 90          # Stop at 90% confidence
    minimum_initial: 20    # Minimum starting confidence
    maximum_initial: 60    # Maximum starting confidence
  
  hypothesis:
    max_active: 5          # Limit active hypotheses
    auto_reject_below: 20 # Auto-reject low-confidence
  
  loops:
    detection_enabled: true
    max_repetitions: 3
  
  implementation:
    require_operator_approval: true
```

---

## Consequences

### Positive

- Deterministic debugging process
- Reduced investigation loops
- Evidence-driven decisions
- Operator approval gate
- Structured debugging artifacts

### Negative

- Additional complexity in runtime
- Configuration requirements
- Training needed for operators

### Risks

- Confidence threshold may need tuning
- Loop detection may have false positives
- Evidence categorization may need refinement

---

## Alternatives Considered

### Alternative 1: Prompt-Based Debugging

Continue using prompt engineering to guide debugging.

**Rejected because:** Not deterministic, relies on AI reasoning rather than structured evidence.

### Alternative 2: Full Automated Fix

Let the runtime fix issues automatically without operator approval.

**Rejected because:** Violates KDSE principle of human approval for implementation.

### Alternative 3: Test-Then-Fix

Run tests, immediately try to fix first failure.

**Rejected because:** No evidence-driven root cause analysis, may fix symptoms instead of causes.

---

## Implementation Plan

### Phase 1: Core Engine
1. Create `runtime/debug/engine.sh`
2. Implement evidence collection
3. Implement hypothesis management
4. Implement confidence tracking

### Phase 2: Intelligence
5. Implement loop detection
6. Implement root cause selection
7. Implement report generation

### Phase 3: Integration
8. Add `kdse debug` command
9. Integrate with `kdse verify`
10. Integrate with `kdse repair`
11. Integrate with `kdse audit`

---

## References

- [DEBUG_RUNTIME.md](../runtime/DEBUG_RUNTIME.md) - Detailed architecture
- [ARCHITECTURE.md](../runtime/ARCHITECTURE.md) - Runtime architecture
- [SESSION_PROTOCOL.md](../runtime/SESSION_PROTOCOL.md) - Session lifecycle
