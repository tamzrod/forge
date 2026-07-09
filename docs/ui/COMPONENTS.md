# Components

## Overview

Reusable UI components for Forge. Each component follows the design language and layout principles.

---

## Toolbar

**Purpose:** Global actions and navigation header

### Structure
```
┌─────────────────────────────────────────────────────────────┐
│ [Logo] Forge: Project Name    [Search]  [▶ Run] [⏹ Stop] ⚙│
└─────────────────────────────────────────────────────────────┘
```

### Elements
- Application logo (optional)
- Workspace/project name
- Global search input
- Simulation controls (Run/Stop/Pause)
- Settings button
- User menu

### Behavior
- Fixed position at top
- Always visible
- Responsive: collapses on narrow screens

### Variants
- **Default**: Full toolbar
- **Compact**: Icons only
- **Minimal**: Title + essential controls only

---

## Navigation Item

**Purpose:** Left sidebar navigation items

### Structure
```
┌─────────────────┐
│ 📊 Dashboard    │
│ 🌍 World        │ ← Active (highlighted)
│   Models        │ ← Expanded sub-item
│   Devices       │
└─────────────────┘
```

### States
| State | Appearance |
|-------|------------|
| Default | Icon + text, muted |
| Hover | Background lighten |
| Active | Blue left border, brighter text |
| Expanded | Shows children below |
| Collapsed | Shows parent only |

### Behavior
- Click: Navigate to section
- Arrow click: Expand/collapse children
- Drag: Reorder (future)

---

## Explorer Tree

**Purpose:** Hierarchical display of simulation objects

### Structure
```
🌐 Simulation World
├── 📊 Models
│   ├── ⏱️ Clock
│   └── ☀️ Sun
└── 📱 Devices
    └── 🌤️ Weather Station
```

### Elements
- Expand/collapse chevron
- Object icon
- Object name
- Selection indicator
- Context menu trigger

### States
| State | Appearance |
|-------|------------|
| Default | Normal text |
| Selected | Blue background highlight |
| Hover | Slight background change |
| Expanded | Chevron points down |
| Collapsed | Chevron points right |

### Behavior
- Single selection
- Double-click: Open in Inspector
- Right-click: Context menu
- Drag: Reorder (future)

---

## Inspector Card

**Purpose:** Container for Inspector panel content

### Structure
```
┌────────────────────────────────────┐
│ Title                          [⚙] │
├────────────────────────────────────┤
│ Content                            │
│                                    │
└────────────────────────────────────┘
```

### Elements
- Title (object name + type)
- Action button (settings, refresh)
- Content area

### Variants
- **Default**: Standard card
- **Section**: Compact, no shadow
- **Elevated**: More prominent (dialogs)

---

## Property Grid

**Purpose:** Display key-value property pairs

### Structure
```
┌────────────────────────────────────┐
│ Property Name          Value       │
│ Another Property       Value       │
│ Status              ● Running      │
└────────────────────────────────────┘
```

### Elements
- Property label (left-aligned, muted)
- Property value (right-aligned, monospace)
- Status indicator (optional colored dot)

### Alignment
- Labels: Left
- Values: Right
- Consistent spacing between rows

---

## Status Badge

**Purpose:** Compact status indicator

### Structure
```
┌──────────────────┐
│ ● Running        │
│ ● Warning        │
│ ● Fault          │
│ ● Offline        │
└──────────────────┘
```

### Variants
| Status | Color | Dot |
|--------|-------|-----|
| Running | Green #00ff88 | Yes |
| Warning | Yellow #ffaa00 | Yes |
| Fault | Red #ff4444 | Yes |
| Offline | Gray #666 | Yes |

### Sizing
- Small: 8px text
- Default: 12px text
- Large: 14px text

---

## Value Display

**Purpose:** Display a measurement or state value

### Structure
```
┌────────────────────────────────────┐
│ Label                              │
│                                    │
│ 480.5 V                           │
│                                   │
│ < Normal                          │
└────────────────────────────────────┘
```

### Elements
- Label (above, muted)
- Value (large, monospace, colored by status)
- Status indicator (below, optional)

### Value States
| State | Color |
|-------|-------|
| Good | Green |
| Warning | Yellow/Orange |
| Critical | Red |
| Neutral | White/Gray |

---

## Console

**Purpose:** Developer output panel

### Structure
```
┌─────────────────────────────────────────────┐
│ Console │ Logs │ Events │ Alarms        [×]│
├─────────────────────────────────────────────┤
│ [INFO] 10:23:45 Simulation started          │
│ [DEBUG] 10:23:46 Loading models...          │
│ [WARN]  10:23:47 Low memory warning         │
│ [ERROR] 10:23:48 Connection failed          │
│                                            │
└─────────────────────────────────────────────┘
```

### Elements
- Tab bar
- Clear button
- Filter input
- Output area with line numbers
- Timestamp
- Severity badge

### Severity Colors
| Level | Color |
|-------|-------|
| INFO | Default text |
| DEBUG | Muted gray |
| WARN | Yellow |
| ERROR | Red |

---

## Tabs

**Purpose:** Organize content within panels

### Structure
```
┌────────┬────────┬────────┐
│ Tab 1  │ Tab 2 ●│ Tab 3  │
└────────┴────────┴────────┘
```

### States
| State | Appearance |
|-------|------------|
| Inactive | Muted text, no underline |
| Active | Bright text, blue underline |
| Hover | Slight background change |
| Disabled | Grayed out |

### Variants
- **Underline**: Bottom border (default)
- **Pill**: Rounded backgrounds
- **Icon**: Icons only

---

## Button

**Purpose:** Trigger actions

### Variants
| Variant | Usage |
|---------|-------|
| Primary | Main actions (Run, Save) |
| Secondary | Alternative actions |
| Ghost | Subtle actions |
| Danger | Destructive actions |

### States
| State | Appearance |
|-------|------------|
| Default | Normal |
| Hover | Slightly brighter |
| Active | Pressed effect |
| Disabled | Grayed out |
| Loading | Spinner icon |

### Sizes
- Small: 28px height
- Default: 36px height
- Large: 44px height

---

## Input

**Purpose:** Text input fields

### Structure
```
┌─────────────────────────────────────┐
│ Label                               │
│ [Input field                     🔍] │
│ Hint text                          │
└─────────────────────────────────────┘
```

### Elements
- Label (above)
- Input field
- Hint text (below, muted)
- Icon (optional, right side)

### States
| State | Appearance |
|-------|------------|
| Default | Gray border |
| Focus | Blue border, glow |
| Error | Red border |
| Disabled | Grayed out |

---

## Dialog / Modal

**Purpose:** Focused interaction requiring user attention

### Structure
```
┌───────────────────────────────────────┐
│                           [×]         │
│ Title                                 │
├───────────────────────────────────────┤
│                                       │
│ Content                               │
│                                       │
├───────────────────────────────────────┤
│          [Cancel]  [Confirm]          │
└───────────────────────────────────────┘
```

### Elements
- Close button (top right)
- Title
- Content area
- Action buttons (bottom right)

### Behavior
- Overlay darkens background
- Focus trapped within dialog
- Escape closes (with confirmation for destructive)
- Click outside closes (unless critical)

### Variants
- **Default**: Standard modal
- **Confirm**: Simple yes/no
- **Alert**: Warning without actions
- **Critical**: Destructive confirmation required

---

## Tree View

See **Explorer Tree** above.

---

## Table

**Purpose:** Display structured data rows

### Structure
```
┌──────┬──────────┬─────────┬────────┐
│ Name │ Status   │ Value   │ Actions│
├──────┼──────────┼─────────┼────────┤
│ Dev1 │ ● Running│ 480.5 V │ [⚙][×]│
│ Dev2 │ ● Warning│ 445.2 V │ [⚙][×]│
│ Dev3 │ ● Fault  │ --      │ [⚙][×]│
└──────┴──────────┴─────────┴────────┘
```

### Elements
- Header row (bold, sortable)
- Data rows
- Selection checkbox (optional)
- Action buttons

### Features
- Column sorting
- Column resizing
- Row selection
- Pagination
- Empty state

---

## Card

**Purpose:** Group related content

### Structure
```
┌────────────────────────────────────┐
│ Header                       [⚙]  │
├────────────────────────────────────┤
│ Body                               │
│                                    │
│                                    │
└────────────────────────────────────┘
```

### Variants
| Variant | Usage |
|---------|-------|
| Default | Standard shadow |
| Flat | No shadow, border only |
| Elevated | Higher shadow |

### Use Cases
- Dashboard widgets
- Device overview
- Quick stats

---

## Component Patterns

### Composing Components

Components should be composable:

```
Card
├── Card.Header
│   ├── Title
│   └── ActionButton
├── PropertyGrid
│   ├── PropertyRow (label, value)
│   ├── PropertyRow (label, value)
│   └── PropertyRow (label, badge)
└── Card.Footer
    └── Button
```

### Responsive Behavior

| Screen Size | Panel Behavior |
|-------------|----------------|
| Large (>1400px) | All panels visible |
| Medium (1024-1400px) | Collapsible sidebars |
| Small (<1024px) | Tabbed interface |

---

*This component library should cover all common UI patterns in Forge. New components should follow these patterns and be added to this document.*
