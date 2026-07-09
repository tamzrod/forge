# Colors

## Philosophy

Colors communicate **meaning**, not decoration. Every color has semantic purpose.

**Rule:** Never use colors for aesthetics alone. If a color doesn't convey information, don't use it.

## Semantic Color System

| Color | Hex | Usage | Examples |
|-------|-----|-------|----------|
| **Green** | #00ff88 | Healthy, Running, Connected, Good | Device online, healthy status, success |
| **Yellow** | #ffaa00 | Warning, Transition, Caution | Low memory, degraded performance |
| **Orange** | #ff8c00 | Environmental, Temperature, Solar | Temperature values, sun irradiance |
| **Blue** | #00d4ff | Selection, Navigation, Interactive | Selected items, active navigation, links |
| **Purple** | #b388ff | Engineering metadata | Pressure, flow rates, engineering values |
| **Gray** | #666666 | Disabled, Offline, Unknown | Inactive devices, unavailable features |
| **Red** | #ff4444 | Fault, Alarm, Critical, Error | Connection lost, fault condition, critical alert |

## Dark Theme Base Colors

### Background Layers

| Token | Hex | Usage |
|-------|-----|-------|
| `bg-primary` | #1a1a2e | Main background |
| `bg-secondary` | #16213e | Card backgrounds, panels |
| `bg-tertiary` | #0f3460 | Elevated elements, hover states |
| `bg-input` | #1e2a4a | Input fields, dropdowns |

### Text Colors

| Token | Hex | Usage |
|-------|-----|-------|
| `text-primary` | #ffffff | Headings, important labels |
| `text-secondary` | #cccccc | Body text |
| `text-muted` | #888888 | Hints, disabled text |
| `text-accent` | #00d4ff | Links, interactive elements |

### Border Colors

| Token | Hex | Usage |
|-------|-----|-------|
| `border-default` | rgba(255,255,255,0.1) | Default borders |
| `border-hover` | rgba(255,255,255,0.2) | Hover state borders |
| `border-active` | #00d4ff | Active/focused borders |

## Status Colors in Detail

### Green (#00ff88)
**Meaning:** Healthy, operational, positive state

| State | Example |
|-------|---------|
| Running | Simulation running |
| Connected | Device connected |
| Healthy | All systems nominal |
| Online | Service available |
| Success | Operation completed |

### Yellow (#ffaa00)
**Meaning:** Warning, degraded, attention needed

| State | Example |
|-------|---------|
| Warning | Approaching limit |
| Transition | State changing |
| Pending | Awaiting confirmation |
| Caution | User attention suggested |

### Orange (#ff8c00)
**Meaning:** Environmental or physical measurements

| State | Example |
|-------|---------|
| Temperature | Ambient, equipment temp |
| Solar | Irradiance, sun position |
| Heat | Thermal readings |
| Thermal | Temperature-related values |

This separates physical measurements from system warnings.

### Blue (#00d4ff)
**Meaning:** Navigation, selection, interactive elements

| State | Example |
|-------|---------|
| Selected | Current selection |
| Active | Active navigation item |
| Link | Clickable elements |
| Focus | Keyboard focus indicator |

### Purple (#b388ff)
**Meaning:** Engineering metadata, calculated values

| State | Example |
|-------|---------|
| Pressure | Fluid, atmospheric pressure |
| Flow | Flow rates |
| Engineering | Calculated metrics |
| Derived | Computed values |

### Gray (#666666)
**Meaning:** Disabled, inactive, unknown

| State | Example |
|-------|---------|
| Disabled | Cannot interact |
| Offline | Device disconnected |
| Inactive | Not running |
| Unknown | Status unclear |

### Red (#ff4444)
**Meaning:** Fault, critical, error

| State | Example |
|-------|---------|
| Fault | Device fault |
| Alarm | Critical alarm active |
| Error | Operation failed |
| Disconnected | Connection lost |
| Critical | Immediate attention required |

## Component Colors

### Status Badges

```css
.badge-healthy { background: #00ff88; color: #000; }
.badge-warning { background: #ffaa00; color: #000; }
.badge-fault { background: #ff4444; color: #fff; }
.badge-offline { background: #666; color: #fff; }
```

### Value Colors

| Component | Normal | Warning | Critical |
|-----------|--------|---------|----------|
| Voltage | #00ff88 | #ffaa00 | #ff4444 |
| Frequency | #00ff88 | #ffaa00 | #ff4444 |
| Temperature | #ff8c00 | #ff4444 | #ff4444 |
| Pressure | #b388ff | #ffaa00 | #ff4444 |

### Threshold Guidelines

| Metric | Good | Warning | Critical |
|--------|------|---------|----------|
| Voltage PU | 0.95-1.05 | 0.90-0.95 | <0.90 or >1.10 |
| Frequency PU | 0.995-1.005 | 0.99-0.995 | <0.98 or >1.02 |
| Temperature | 15-35°C | 35-45°C | >45°C or <5°C |

## Accessibility

### Contrast Requirements
- Minimum 4.5:1 for normal text
- Minimum 3:1 for large text and UI components
- Never rely on color alone for information

### Color Blind Considerations
- Pair colors with icons or labels
- Use patterns in addition to colors for charts
- Never use red-green alone for status

## Usage Rules

### DO
- Use semantic colors consistently
- Pair colors with text labels
- Provide icons for status
- Use muted versions for backgrounds

### DON'T
- Use color for decoration
- Create rainbow effects
- Use bright pastels on dark backgrounds
- Mix similar colors (green/cyan, yellow/orange)

---

*Color choices should be consistent across all Forge UI. When adding new colors, ensure they have clear semantic purpose.*
