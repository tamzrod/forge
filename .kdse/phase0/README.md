# Phase 0: Runtime Initialization

**Document Version:** 1.0  
**Type:** Runtime Implementation  
**Effective Date:** 2026-07-11  

---

## Purpose

Phase 0 ensures the KDSE Runtime automatically loads the KDSE methodology into AI working context before any engineering activity begins.

**Engineering Principle:**
> "The Runtime—not the operator prompt—shall own AI initialization."

---

## Bootstrap Sequence

```
┌─────────────────────────────────────────────────────────────┐
│                    PHASE 0 BOOTSTRAP                         │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Step 1. Discover Installation                               │
│  Step 2. Load Manifest                                      │
│  Step 3. Verify Versions                                     │
│  Step 4. Load Knowledge                                      │
│  Step 5. Verify Integrity                                    │
│  Step 6. Discover Capabilities                               │
│  Step 7. Generate AI Context                                 │
│  Step 8. Produce Initialization Summary                      │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## Step Details

### Step 1: Discover Installation

Verifies KDSE Runtime is properly installed.

**Checks:**
- `.kdse/` directory exists
- `.kdse/knowledge/manifest.yaml` exists
- `.kdse/config.sh` exists (if installed)

**Failure Mode:**
```
ERROR: KDSE Runtime not installed
Hint: Run ./runtime/install/install.sh
```

### Step 2: Load Manifest

Parses the Knowledge Manifest to understand required knowledge.

**Actions:**
- Read `.kdse/knowledge/manifest.yaml`
- Validate manifest schema
- Extract required/optional knowledge lists
- Extract capability definitions

**Failure Mode:**
```
ERROR: Invalid Knowledge Manifest
Hint: manifest.yaml is corrupted or missing required fields
```

### Step 3: Verify Versions

Ensures Runtime and Standard versions are compatible.

**Checks:**
- Runtime version from manifest
- Standard version compatibility
- Minimum version requirements

**Failure Mode:**
```
ERROR: Version incompatibility detected
Hint: Run 'kdse update' to sync with compatible version
```

### Step 4: Load Knowledge

Loads all required knowledge documents in specified order.

**Loading Order:**
| Order | Document | Source |
|-------|----------|--------|
| 1 | Core Principles | docs/foundation/003-core-principles.md |
| 2 | Engineering Model | docs/foundation/004-engineering-model.md |
| 3 | Chain of Authority | docs/foundation/006-chain-of-authority.md |
| 4 | Glossary | docs/foundation/007-glossary.md |
| 5 | Session Protocol | runtime/SESSION_PROTOCOL.md |
| 6 | Command Registry | runtime/install/commands.yaml |
| 7 | Runtime Configuration | runtime/VERSIONING.md |

**Failure Mode:**
```
ERROR: Required knowledge missing
  Missing: docs/foundation/003-core-principles.md
Hint: Restore missing file from KDSE repository
```

### Step 5: Verify Integrity

Generates and verifies knowledge fingerprint.

**Fingerprint Calculation:**
```
Fingerprint = SHA256(
  sorted([
    "source:path" + SHA256(content),
    ...
  ])
)
```

### Step 6: Discover Capabilities

Builds capability registry from loaded knowledge.

**Capabilities:**
- Assessment
- Architecture
- Implementation
- Verification
- Evolution
- Feedback

### Step 7: Generate AI Context

Creates machine-readable AI knowledge artifact.

**Output:** `.kdse/knowledge/kdse-ai.json`

### Step 8: Produce Initialization Summary

Displays human-readable initialization report.

---

## Initialization Summary Format

```
═══════════════════════════════════════════════════════════════
                    KDSE Runtime Initialized
═══════════════════════════════════════════════════════════════

Runtime Version:    1.0.0
Knowledge Version:  1.0.0
Knowledge Fingerprint: a3f2b8c1d4e5f6...

Capabilities Loaded:
  ✓ Assessment
  ✓ Architecture
  ✓ Verification
  ✓ Evolution
  ✓ Feedback

Repository:
  Path: /workspace/project/KDSE
  Lifecycle: Active

Status: READY

═══════════════════════════════════════════════════════════════
```

---

## Knowledge Manifest

The Knowledge Manifest (`.kdse/knowledge/manifest.yaml`) defines:

- **Required Knowledge:** Must be loaded before engineering activity
- **Optional Knowledge:** Should be loaded for complete context
- **Loading Order:** Explicit sequence for loading
- **Capability Definitions:** Available capabilities after initialization
- **Workflow Entrypoints:** Standard workflow commands

---

## AI Knowledge Artifact

The AI Knowledge Artifact (`.kdse/knowledge/kdse-ai.json`) contains:

- Machine-readable engineering knowledge
- Capability registry
- Command definitions
- Core principles
- Authority hierarchy
- Lifecycle stages

---

## Knowledge Fingerprint

The Knowledge Fingerprint is a SHA-256 hash that verifies knowledge integrity.

**Purpose:**
- Detects unauthorized knowledge changes
- Ensures reproducibility of initialization
- Provides audit trail

**Format:**
```
Knowledge Fingerprint: sha256:a3f2b8c1d4e5f6...
```

---

## Capability Discovery

After Phase 0 initialization, the following capabilities are available:

| Capability | Description | Entrypoint |
|------------|-------------|------------|
| Assessment | Repository compliance assessment | Run KDSE |
| Architecture | Architecture design and review | Engineering Phase |
| Implementation | Implementation guidance | Engineering Phase |
| Verification | Implementation verification | Verification Phase |
| Evolution | Methodology evolution | Evolution Phase |
| Feedback | Feedback collection | Session Protocol |

---

## Failure Modes

### Missing Installation

```
ERROR: KDSE Runtime not installed
Hint: Run ./runtime/install/install.sh
```

### Invalid Manifest

```
ERROR: Invalid Knowledge Manifest
Hint: manifest.yaml is corrupted or missing required fields
```

### Version Incompatibility

```
ERROR: Version incompatibility detected
  Runtime Version: 1.0.0
  Standard Version: 0.9.0
  Required: >= 1.0.0
Hint: Run 'kdse update' to sync with compatible version
```

### Missing Required Knowledge

```
ERROR: Required knowledge missing
  Missing: docs/foundation/003-core-principles.md
  Required by: manifest.yaml
Hint: Restore missing file from KDSE repository
```

---

## Files

| File | Purpose |
|------|---------|
| `phase0-init.sh` | Main Phase 0 bootstrap script |
| `runtime-state.sh` | Runtime state management |
| `load-knowledge.sh` | Knowledge loading script |
| `generate-fingerprint.sh` | Fingerprint generation |

---

## Usage

### Automatic Initialization

Phase 0 runs automatically when starting a KDSE session:

```bash
Run KDSE
```

### Manual Initialization

To manually run Phase 0:

```bash
./.kdse/phase0/phase0-init.sh --verbose
```

### Check Status

```bash
cat .kdse/runtime/state.json
```

---

*This document defines Phase 0: Runtime Initialization for the KDSE Runtime.*
