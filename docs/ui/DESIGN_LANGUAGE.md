# Design Language

## Philosophy

Forge is an engineering workbench. The interface should feel like a professional development environment, not a consumer application.

### Core Principles

1. **Dark-First Interface**
   - Primary background: Deep charcoal (#1a1a2e or similar)
   - Reduces eye strain during extended use
   - Industry standard for developer tools
   - Consistent with modern IDE aesthetics

2. **Calm Professional Appearance**
   - No flashy graphics or gradients
   - Subtle shadows and borders
   - Muted accent colors
   - Professional, not playful

3. **Information Over Decoration**
   - Every visual element serves a purpose
   - No decorative graphics
   - Data density over whitespace
   - Compact but readable

4. **Engineering Workbench Philosophy**
   - Resembles Visual Studio Code, JetBrains IDEs
   - Multiple panels with clear boundaries
   - Tree-based navigation
   - Property editors
   - Console output

5. **Consistent Spacing**
   - 4px base unit
   - Standard padding: 8px, 12px, 16px, 24px
   - Consistent gap between elements
   - Predictable layouts

6. **Minimal Animations**
   - Subtle transitions only (150-200ms)
   - No bouncing, spinning, or attention-seeking effects
   - State changes are immediate
   - Loading indicators are functional, not decorative

7. **Rounded Cards**
   - Border radius: 4-8px
   - Soft, modern appearance
   - Not too playful (avoid 16px+ radius)

8. **Soft Contrast**
   - Avoid harsh white-on-black
   - Use mid-grays for text
   - Clear hierarchy without stark contrast
   - 4.5:1 minimum contrast ratio for accessibility

9. **Readability Before Aesthetics**
   - Monospace fonts for values and code
   - Sans-serif for labels and navigation
   - Clear typography hierarchy
   - Adequate line height

## Visual Language

### What We Communicate

The interface should communicate these qualities:

| Quality | How We Achieve It |
|---------|-------------------|
| **Structure** | Clear panel boundaries, consistent layout |
| **Hierarchy** | Font sizes, weights, and color distinguish importance |
| **Clarity** | Semantic colors, descriptive labels |
| **Stability** | Minimal motion, consistent behavior |

### What We Avoid

| Anti-Pattern | Why We Avoid It |
|--------------|-----------------|
| Industrial graphics | SCADA-style gauges and meters are inappropriate |
| Bright colors | Creates visual fatigue |
| Heavy shadows | Distracting, dated appearance |
| Gradient backgrounds | Decorative, not functional |
| Complex animations | Unprofessional, distracting |
| Consumer UI patterns | Mismatched with engineering context |

## Typography

### Font Families

- **Headings**: System UI or Inter (sans-serif)
- **Body**: System UI or Inter (sans-serif)
- **Values/Code**: JetBrains Mono, Fira Code, or Consolas (monospace)

### Type Scale

| Element | Size | Weight |
|---------|------|--------|
| Page Title | 20-24px | 600 |
| Section Header | 16px | 600 |
| Card Title | 14px | 600 |
| Body Text | 14px | 400 |
| Labels | 12-13px | 400-500 |
| Values | 13-14px | 400 (monospace) |
| Small Text | 11-12px | 400 |

## Spacing System

Based on 4px grid:

| Token | Value | Usage |
|-------|-------|-------|
| `xs` | 4px | Icon gaps, tight spacing |
| `sm` | 8px | Component internal padding |
| `md` | 12px | Card padding |
| `lg` | 16px | Section gaps |
| `xl` | 24px | Panel padding |
| `xxl` | 32px | Major section separation |

## Borders & Shadows

### Borders
- Color: rgba(255, 255, 255, 0.1) or #0f3460
- Width: 1px
- Use for panel separation and card boundaries

### Shadows
- Minimal and subtle
- `0 2px 4px rgba(0, 0, 0, 0.3)` for elevated cards
- Avoid heavy drop shadows

## Iconography

- Use consistent icon set (Lucide, Material Icons, etc.)
- Size: 16px for inline, 20-24px for standalone
- Stroke weight: 1.5-2px
- Color: inherit or muted gray

---

*Design decisions should align with this language. When in doubt, choose simpler, darker, and more professional.*
