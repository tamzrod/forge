# Architectural Decision Records (ADRs)

This directory contains architectural decision records for the Industrial Simulation Runtime.

## Index

| ID | Title | Status | Date |
|----|-------|--------|------|
| [ADR-001](001-runtime-architecture.md) | Runtime Architecture | Accepted | 2026-07-11 |
| [ADR-002](002-behavior-model.md) | Behavior Model Design | Accepted | 2026-07-11 |
| [ADR-003](003-memory-model.md) | Memory Model Design | Accepted | 2026-07-11 |
| [ADR-004](004-simulation-models.md) | Simulation Models Design | Accepted | 2026-07-11 |
| [ADR-005](005-plugin-architecture-audit.md) | Plugin Architecture Audit | Accepted | 2026-07-13 |
| [ADR-006](006-plugin-architecture.md) | Plugin Architecture | Accepted | 2026-07-13 |

---

## ADR Template

Each ADR follows this structure:

1. **Context**: Problem statement
2. **Decision**: The chosen solution
3. **Consequences**: Benefits, drawbacks, risks
4. **Alternatives**: Options that were considered
5. **References**: Related documentation
6. **Milestone Traceability**: Link to roadmap milestones

---

## Creating New ADRs

1. Copy the template from [ADR-TEMPLATE.md](ADR-TEMPLATE.md)
2. Name the file: `###-title-lowercase-with-dashes.md`
3. Update the index above
4. Commit with message: `docs: Add ADR-###: Title`

---

## Status Definitions

| Status | Meaning |
|--------|---------|
| Proposed | Under review |
| Accepted | Approved for implementation |
| Deprecated | Superseded by another ADR |
| Rejected | Not adopted |

---

## Related Documentation

- [Architecture Overview](../architecture/overview.md)
- [Engineering Model](../../.kdse/docs/004-engineering-model.md)
- [Glossary](../architecture/GLOSSARY.md)
