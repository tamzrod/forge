# UI Design System

This is the authoritative reference for all Forge UI development.

## Purpose

Forge is an **engineering workbench** for building and operating Virtual Industrial Laboratories. The UI should resemble professional engineering tools:

- Visual Studio Code
- Unity Editor
- Unreal Editor
- JetBrains IDEs

Forge is **NOT**:
- A SCADA system
- A monitoring dashboard
- A consumer application

## Design Philosophy

The interface communicates:
- **Structure** - Clear organization of information
- **Hierarchy** - Visual precedence of elements
- **Clarity** - Unambiguous meaning
- **Stability** - Professional, trustworthy appearance

See [Design Language](DESIGN_LANGUAGE.md) for detailed philosophy.

## Document Index

| Document | Purpose |
|----------|---------|
| [Workspaces](WORKSPACES.md) | Functional areas and workspace hierarchy |
| [Design Language](DESIGN_LANGUAGE.md) | Overall design philosophy and principles |
| [Layout](LAYOUT.md) | Application layout, panels, and navigation |
| [Colors](COLORS.md) | Semantic color system |
| [Components](COMPONENTS.md) | Reusable UI components |
| [Generic Inspector](GENERIC_INSPECTOR.md) | Data-driven inspection framework |

## Quick Reference

### Primary Layout
```
┌─────────────────────────────────────────────────────────────┐
│ Toolbar                                                       │
├──────────────┬──────────────────────┬─────────────────────────┤
│ Navigation   │ World Explorer       │ Inspector               │
│              │                      │                         │
│              │                      │                         │
├──────────────┴──────────────────────┴─────────────────────────┤
│ Console / Logs / Events                                       │
└─────────────────────────────────────────────────────────────┘
```

### Color Semantics
| Color | Meaning |
|-------|---------|
| Green | Healthy, Running, Connected |
| Yellow | Warning, Transition |
| Orange | Environmental (Temperature, Solar) |
| Blue | Selection, Navigation |
| Purple | Engineering Metadata (Pressure) |
| Gray | Disabled, Offline, Unknown |
| Red | Fault, Alarm, Critical |

## Design Rules

### DO
- Use dark-first interface
- Prioritize information over decoration
- Maintain consistent spacing
- Keep animations minimal
- Use semantic colors

### DON'T
- Use colors for decoration
- Create SCADA-style graphics
- Build historical trends into the Inspector
- Add unnecessary visual flourishes

## Future Workspaces

All future pages should inherit this visual language:

- World Editor
- Device Editor
- Scenario Editor
- Network Editor
- Protocol Monitor
- Data Explorer
- Datasheet Importer
- Device Library
- Simulation Inspector

---

*This design system is the long-term UI reference for Forge.*
