# Generic Inspector Framework

> A data-driven inspection framework for the Engineering Workbench

## Overview

The Generic Inspector is a reusable framework that provides read-only visibility into any object within Forge. It replaces the hardcoded inspection logic with a structured, extensible data model.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    Engineering Workbench                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────────┐        ┌────────────────────────────────┐ │
│  │  World Explorer  │───────▶│     Generic Inspector          │ │
│  │  (Navigation)    │        │                                │ │
│  └─────────────────┘        │  ┌──────────────────────────┐  │ │
│                             │  │ Section: Identity         │  │ │
│                             │  │ Section: Overview         │  │ │
│                             │  │ Section: State           │  │ │
│                             │  │ Section: Configuration    │  │ │
│                             │  │ Section: Diagnostics      │  │ │
│                             │  │ Section: Communications   │  │ │
│                             │  │ Section: Memory           │  │ │
│                             │  │ Section: Children        │  │ │
│                             │  └──────────────────────────┘  │ │
│                             └────────────────────────────────┘ │
│                                           ▲                     │
│                                           │                     │
│                             ┌─────────────┴─────────────┐       │
│                             │     REST API               │       │
│                             │  GET /api/inspect/{id}    │       │
│                             └─────────────┬─────────────┘       │
└───────────────────────────────────────────┼─────────────────────┘
                                            │
┌───────────────────────────────────────────┼─────────────────────┐
│                    Backend                 │                      │
│                                           ▼                      │
│  ┌────────────────────────────────────────────────────────────────┐│
│  │                    Inspector View                               ││
│  │                                                                 ││
│  │   Simulation Models      Virtual Devices      Other Objects   ││
│  │   ┌─────────┐          ┌─────────┐          ┌─────────┐       ││
│  │   │ Clock   │          │ Weather │          │ Clock   │       ││
│  │   │ Sun     │          │ Station │          │ Driver  │       ││
│  │   │ Weather │          │ Revenue │          │ ...     │       ││
│  │   │ Grid    │          │ Meter   │          └─────────┘       ││
│  │   └─────────┘          └─────────┘                            ││
│  └────────────────────────────────────────────────────────────────┘│
└───────────────────────────────────────────────────────────────────┘
```

## Supported Objects

### Simulation Models
- **Clock**: Simulation time management
- **Sun**: Solar position and irradiance
- **Weather**: Atmospheric conditions
- **Grid**: Power grid state

### Virtual Devices
- **Weather Station**: Environmental monitoring device
- Future: Revenue Meter, PV Inverter, Battery, Relay, PLC, IED, RTU

### Nested Objects
- **Virtual Firmware**: Device control logic
- **Device Memory**: Internal state storage
- **Communication Interface**: External data transmission

## Inspector Sections

Every object exposes sections based on available data:

| Section | Description | Visibility |
|---------|-------------|------------|
| **Identity** | Name, Type, ID | Always shown |
| **Overview** | High-level summary | Always shown |
| **State** | Current operational values | When applicable |
| **Configuration** | Setup parameters | When applicable |
| **Diagnostics** | Health and error info | When applicable |
| **Communications** | Interface statistics | Devices with interfaces |
| **Memory** | Device Memory contents | Devices with memory |
| **Children** | Nested inspectable objects | When applicable |

Sections with no content are automatically hidden.

## Property Types

The framework supports the following property types for rendering:

| Type | Description | Example |
|------|-------------|---------|
| `text` | Plain string values | Device name |
| `number` | Numeric values with optional unit | Temperature: 25.5 °C |
| `boolean` | Yes/No values | Is Paused: Yes |
| `status` | Operational status with color | Running, Warning, Fault |
| `timestamp` | Date/time values | Last Update: 2026-07-09 12:30 |
| `duration` | Time durations | Elapsed: 2h 15m 30s |
| `quality` | Data quality indicators | Good, Uncertain, Bad |
| `enum` | Enumeration values | Mode: Realtime |
| `nested` | Grouped properties | Measurements { ... } |
| `list` | Lists of items | Memory Regions: 4 |
| `angle` | Angular degrees | Azimuth: 180° |
| `percentage` | Percentage values | Humidity: 65% |

## API

### Endpoint

```
GET /api/inspect/{objectId}
```

### Response Format

```json
{
  "object": {
    "id": "weather-station-001",
    "type": "device",
    "name": "Weather Station 001"
  },
  "sections": [
    {
      "id": "identity",
      "title": "Identity",
      "icon": "tag",
      "properties": [
        { "name": "ID", "value": "weather-station-001", "type": "text" },
        { "name": "Type", "value": "weather_station", "type": "text" },
        { "name": "Name", "value": "Weather Station 001", "type": "text" }
      ]
    },
    {
      "id": "overview",
      "title": "Overview",
      "icon": "eye",
      "properties": [
        { "name": "Status", "value": "Running", "type": "status" },
        { "name": "Temperature", "value": 25.5, "type": "number", "unit": "°C" }
      ]
    },
    {
      "id": "memory",
      "title": "Device Memory",
      "icon": "cpu",
      "properties": [
        { "name": "Measurements", "type": "nested", "children": [...] }
      ]
    },
    {
      "id": "children",
      "title": "Children",
      "icon": "layers",
      "children": [
        { "id": "device-001-firmware", "type": "firmware", "name": "Virtual Firmware" },
        { "id": "device-001-memory", "type": "memory", "name": "Device Memory" },
        { "id": "device-001-interface", "type": "interface", "name": "Communication Interface" }
      ]
    }
  ]
}
```

### Object IDs

| Object | ID | Description |
|--------|-----|-------------|
| Simulation World | `world` | Root world object |
| Clock Model | `clock` | Simulation clock |
| Sun Model | `sun` | Solar model |
| Weather Model | `weather` | Weather model |
| Grid Model | `grid` | Power grid model |
| Device | `device-{id}` | Specific device by ID |
| Firmware | `device-{id}-firmware` | Device's virtual firmware |
| Memory | `device-{id}-memory` | Device's memory space |
| Interface | `device-{id}-interface` | Device's communication interface |

## Backend Implementation

### Data Model

The backend defines structured types in `internal/inspector/generic.go`:

```go
type ObjectType string
type SectionID string
type PropertyType string

type GenericInspectorData struct {
    Object   ObjectIdentity `json:"object"`
    Sections []*Section     `json:"sections"`
}

type Section struct {
    ID         SectionID    `json:"id"`
    Title      string       `json:"title"`
    Icon       string       `json:"icon,omitempty"`
    Properties []*Property  `json:"properties,omitempty"`
    Children   []*ObjectRef `json:"children,omitempty"`
}

type Property struct {
    Name       string       `json:"name"`
    Value      interface{}  `json:"value"`
    Type       PropertyType `json:"type"`
    Unit       string       `json:"unit,omitempty"`
    Quality    Quality      `json:"quality,omitempty"`
    Precision  int          `json:"precision,omitempty"`
    Children   []*Property  `json:"children,omitempty"`
    Items      []*Property  `json:"items,omitempty"`
}
```

### Generator

The `internal/inspector/generator.go` file provides the `Generator` struct that builds inspection data for each object type:

```go
type Generator struct {
    view *View
}

func (g *Generator) Inspect(objectID string) (*GenericInspectorData, error)
func (g *Generator) InspectWorld() (*GenericInspectorData, error)
func (g *Generator) InspectClock() (*GenericInspectorData, error)
func (g *Generator) InspectSun() (*GenericInspectorData, error)
func (g *Generator) InspectWeather() (*GenericInspectorData, error)
func (g *Generator) InspectGrid() (*GenericInspectorData, error)
func (g *Generator) InspectDevice(deviceID string) (*GenericInspectorData, error)
func (g *Generator) InspectFirmware(deviceID string) (*GenericInspectorData, error)
func (g *Generator) InspectMemory(deviceID string) (*GenericInspectorData, error)
func (g *Generator) InspectInterface(deviceID string) (*GenericInspectorData, error)
```

## Frontend Implementation

### Components

```
ui/src/components/inspector/
├── index.ts              # Exports
├── GenericInspector.tsx  # Main inspector component
├── SectionCard.tsx      # Section display component
└── PropertyValue.tsx     # Property rendering component
```

### GenericInspector

The main component that:
1. Fetches inspection data from the API
2. Displays available tabs based on section data
3. Renders sections with automatic visibility
4. Handles child navigation

### PropertyValue

Renders individual property values based on type:
- Color coding for temperatures, voltages, etc.
- Status badges with health indication
- Unit display for numeric values
- Quality indicators

## Extensibility

### Adding New Object Types

1. **Backend**: Add inspection method to `Generator`:
   ```go
   func (g *Generator) InspectNewObject(id string) (*GenericInspectorData, error)
   ```

2. **Update** `Inspect()` to route to new method:
   ```go
   switch objectID {
   case "new-object":
       return g.InspectNewObject()
   }
   ```

3. **Frontend**: No changes required! The Generic Inspector automatically picks up the new object.

### Adding New Property Types

1. **Backend**: Define type constant in `generic.go`:
   ```go
   PropertyTypeNewType PropertyType = "new_type"
   ```

2. **Frontend**: Add rendering in `PropertyValue.tsx`:
   ```go
   case 'new_type':
       return <span>{formatNewType(value)}</span>
   ```

3. **Types**: Update `ui/src/types/inspector.ts`

### Adding New Sections

1. **Backend**: Define section ID in `generic.go`:
   ```go
   SectionNewSection SectionID = "new_section"
   ```

2. **Frontend**: Add to `SECTION_META` in `ui/src/types/inspector.ts`:
   ```ts
   new_section: { id: 'new_section', label: 'New Section', icon: 'folder', order: 9 }
   ```

3. **Section to Tab Mapping**: Update `sectionToTab` in `GenericInspector.tsx`

## Terminology

This framework follows the [Architecture Glossary](../architecture/GLOSSARY.md) terminology:

| Glossary Term | Inspector Usage |
|--------------|-----------------|
| Simulation Model | Object type `simulation` |
| Virtual Device | Object type `device` |
| Virtual Firmware | Object type `firmware` |
| Device Memory | Object type `memory` |
| Communication Interface | Object type `interface` |

## Testing

### Backend Tests

```bash
go test ./internal/inspector/...
```

### Frontend Tests

```bash
cd ui
npm test
```

## Future Enhancements

- [ ] Search/filter within sections
- [ ] Property grouping with collapsible sections
- [ ] Historical data visualization
- [ ] Export inspection data (JSON, CSV)
- [ ] Comparison mode (compare two objects)
- [ ] Custom property formatters per device type

---

*Last Updated: 2026-07-09*
*Part of: UI Milestone 1 - Generic Inspector Framework*
