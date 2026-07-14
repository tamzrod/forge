# KDSE Evidence Directory

This directory stores evidence artifacts collected during KDSE sessions.

## Structure

```
.kdse/evidence/
├── screenshots/   # Visual evidence (UI screenshots, diagrams)
├── tests/        # Test results and logs
└── benchmarks/   # Performance data and metrics
```

## Usage

Evidence files are referenced in `.kdse/context.json` under the `evidence` field.

### Adding Evidence

```bash
# Add a screenshot
kdse context add-evidence .kdse/evidence/screenshots/dashboard.png

# Add test results
kdse context add-evidence .kdse/evidence/tests/unit-test-results.json

# Add benchmark data
kdse context add-evidence .kdse/evidence/benchmarks/performance-report.md
```

### Evidence Types

| Type | Purpose | File Types |
|------|---------|------------|
| screenshots | Visual validation evidence | PNG, JPG, GIF, SVG |
| tests | Test execution results | JSON, XML, HTML, MD |
| benchmarks | Performance measurements | JSON, CSV, MD |

## Guidelines

1. **Provenance**: Each evidence file should be traceable to its source
2. **Naming**: Use descriptive names with timestamps when applicable
3. **Organization**: Group related evidence by stage or feature
4. **Retention**: Keep evidence for audit purposes but clean up stale files
