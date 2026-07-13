// Package editor provides the Forge Editor for creating and editing simulation models.
// The Editor is responsible for editing the simulation world, not executing simulation logic.
package editor

import (
	"fmt"

	"github.com/tamzrod/forge/world"
)

// ID uniquely identifies editor elements.
type ID string

// Point represents a 2D coordinate.
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Size represents dimensions.
type Size struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// Bounds represents a rectangular area.
type Bounds struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// EntityType represents the type of entity.
type EntityType string

const (
	EntityTypeGrid         EntityType = "grid"
	EntityTypeBus          EntityType = "bus"
	EntityTypeBreaker      EntityType = "breaker"
	EntityTypeTransformer  EntityType = "transformer"
	EntityTypeGenerator    EntityType = "generator"
	EntityTypeLoad         EntityType = "load"
	EntityTypeMeter        EntityType = "meter"
	EntityTypeSun          EntityType = "sun"
	EntityTypeWeather       EntityType = "weather"
	EntityTypeWind         EntityType = "wind"
	EntityTypeScenario     EntityType = "scenario"
	EntityTypeClock        EntityType = "clock"
)

// EntityCategory categorizes entities in the palette.
type EntityCategory string

const (
	CategoryElectrical   EntityCategory = "electrical"
	CategoryEnvironment EntityCategory = "environment"
	CategorySimulation  EntityCategory = "simulation"
)

// CanvasEntity represents an entity on the canvas.
type CanvasEntity struct {
	ID          ID         `json:"id"`
	EntityType  EntityType `json:"entity_type"`
	Name        string     `json:"name"`
	Position    Point      `json:"position"`
	Size        Size       `json:"size"`
	WorldID     world.EntityID `json:"world_id"`
	Properties  Properties `json:"properties"`
}

// Properties represents editable entity properties.
type Properties map[string]PropertyValue

// PropertyValue represents a single property value.
type PropertyValue struct {
	Value    interface{} `json:"value"`
	Type     string      `json:"type"`
	ReadOnly bool        `json:"readonly"`
	Unit     string      `json:"unit,omitempty"`
	Min      *float64    `json:"min,omitempty"`
	Max      *float64    `json:"max,omitempty"`
	Options  []string    `json:"options,omitempty"`
}

// Terminal represents a connection point on an entity.
type Terminal struct {
	ID       ID     `json:"id"`
	Name     string `json:"name"`
	Position Point  `json:"position"`
	Role     string `json:"role"`
	Type     string `json:"type"`
}

// Connection represents a visual connection between entities.
type Connection struct {
	ID            ID     `json:"id"`
	FromEntity    ID     `json:"from_entity"`
	FromTerminal  string `json:"from_terminal"`
	ToEntity      ID     `json:"to_entity"`
	ToTerminal    string `json:"to_terminal"`
	BusID         string `json:"bus_id"`
}

// CanvasState represents the canvas viewport state.
type CanvasState struct {
	Zoom       float64 `json:"zoom"`
	PanX       float64 `json:"pan_x"`
	PanY       float64 `json:"pan_y"`
	GridVisible bool    `json:"grid_visible"`
	SnapToGrid  bool    `json:"snap_to_grid"`
	GridSize    float64 `json:"grid_size"`
}

// Selection represents the current selection.
type Selection struct {
	EntityIDs []ID `json:"entity_ids"`
	Anchor    *ID  `json:"anchor,omitempty"`
}

// Project represents an editor project.
type Project struct {
	ID        ID              `json:"id"`
	Name      string          `json:"name"`
	Entities  []*CanvasEntity `json:"entities"`
	Connections []*Connection `json:"connections"`
	Canvas    *CanvasState    `json:"canvas"`
	Metadata  ProjectMetadata `json:"metadata"`
}

// ProjectMetadata contains project metadata.
type ProjectMetadata struct {
	CreatedAt   string `json:"created_at"`
	ModifiedAt  string `json:"modified_at"`
	Author      string `json:"author"`
	Description string `json:"description"`
}

// EditorState represents the current editor state.
type EditorState struct {
	Project     *Project    `json:"project"`
	Selection   *Selection  `json:"selection"`
	Inspector   *InspectorState `json:"inspector"`
	IsModified  bool        `json:"is_modified"`
	IsRunning   bool        `json:"is_running"`
	IsPaused    bool        `json:"is_paused"`
	Speed       float64     `json:"speed"`
	CurrentTime string      `json:"current_time"`
}

// InspectorState represents the inspector panel state.
type InspectorState struct {
	SelectedEntity *CanvasEntity `json:"selected_entity,omitempty"`
	Sections       []InspectorSection `json:"sections"`
}

// InspectorSection represents a section in the inspector.
type InspectorSection struct {
	Title      string             `json:"title"`
	Properties []InspectorProperty `json:"properties"`
}

// InspectorProperty represents a property in the inspector.
type InspectorProperty struct {
	Key       string         `json:"key"`
	Label     string         `json:"label"`
	Value     interface{}    `json:"value"`
	Type      string         `json:"type"`
	ReadOnly  bool           `json:"readonly"`
	Unit      string         `json:"unit,omitempty"`
	Min       *float64       `json:"min,omitempty"`
	Max       *float64       `json:"max,omitempty"`
	Options   []string       `json:"options,omitempty"`
}

// PaletteItem represents an item in the palette.
type PaletteItem struct {
	ID         ID             `json:"id"`
	Name       string         `json:"name"`
	Category   EntityCategory `json:"category"`
	EntityType EntityType     `json:"entity_type"`
	Icon       string         `json:"icon"`
	DefaultProperties Properties `json:"default_properties,omitempty"`
	Description string        `json:"description,omitempty"`
}

// TreeNode represents a node in the project explorer tree.
type TreeNode struct {
	ID       string      `json:"id"`
	Label    string      `json:"label"`
	Icon     string      `json:"icon,omitempty"`
	Type     string      `json:"type"`
	Children []*TreeNode `json:"children,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	Expanded bool        `json:"expanded,omitempty"`
}

// NewProject creates a new project.
func NewProject(name string) *Project {
	return &Project{
		ID:   ID(fmt.Sprintf("proj-%s", generateID())),
		Name: name,
		Entities: make([]*CanvasEntity, 0),
		Connections: make([]*Connection, 0),
		Canvas: &CanvasState{
			Zoom:       1.0,
			PanX:       0,
			PanY:       0,
			GridVisible: true,
			SnapToGrid:  true,
			GridSize:    20,
		},
		Metadata: ProjectMetadata{
			CreatedAt:  "",
			ModifiedAt: "",
		},
	}
}

// NewCanvasEntity creates a new canvas entity.
func NewCanvasEntity(entityType EntityType, name string, position Point) *CanvasEntity {
	return &CanvasEntity{
		ID:         ID(generateID()),
		EntityType: entityType,
		Name:       name,
		Position:   position,
		Size:       getDefaultSize(entityType),
		Properties: getDefaultProperties(entityType),
	}
}

// AddEntity adds an entity to the project.
func (p *Project) AddEntity(entity *CanvasEntity) {
	p.Entities = append(p.Entities, entity)
}

// RemoveEntity removes an entity from the project.
func (p *Project) RemoveEntity(id ID) {
	for i, e := range p.Entities {
		if e.ID == id {
			p.Entities = append(p.Entities[:i], p.Entities[i+1:]...)
			return
		}
	}
}

// GetEntity returns an entity by ID.
func (p *Project) GetEntity(id ID) *CanvasEntity {
	for _, e := range p.Entities {
		if e.ID == id {
			return e
		}
	}
	return nil
}

// AddConnection adds a connection to the project.
func (p *Project) AddConnection(conn *Connection) {
	p.Connections = append(p.Connections, conn)
}

// RemoveConnection removes a connection from the project.
func (p *Project) RemoveConnection(id ID) {
	for i, c := range p.Connections {
		if c.ID == id {
			p.Connections = append(p.Connections[:i], p.Connections[i+1:]...)
			return
		}
	}
}

// getDefaultSize returns the default size for an entity type.
func getDefaultSize(entityType EntityType) Size {
	switch entityType {
	case EntityTypeBus:
		return Size{Width: 60, Height: 60}
	case EntityTypeGrid:
		return Size{Width: 80, Height: 60}
	case EntityTypeBreaker:
		return Size{Width: 50, Height: 50}
	case EntityTypeTransformer:
		return Size{Width: 80, Height: 60}
	case EntityTypeGenerator:
		return Size{Width: 80, Height: 80}
	case EntityTypeLoad:
		return Size{Width: 80, Height: 80}
	case EntityTypeMeter:
		return Size{Width: 70, Height: 70}
	default:
		return Size{Width: 60, Height: 60}
	}
}

// getDefaultProperties returns default properties for an entity type.
func getDefaultProperties(entityType EntityType) Properties {
	switch entityType {
	case EntityTypeBus:
		return Properties{
			"name":            {Value: "New Bus", Type: "string"},
			"nominal_voltage": {Value: float64(480), Type: "number", Unit: "V"},
		}
	case EntityTypeGrid:
		return Properties{
			"name":              {Value: "Utility Grid", Type: "string"},
			"nominal_voltage":   {Value: float64(69000), Type: "number", Unit: "V"},
			"nominal_frequency": {Value: float64(60), Type: "number", Unit: "Hz"},
		}
	case EntityTypeBreaker:
		return Properties{
			"name":    {Value: "Breaker", Type: "string"},
			"is_open": {Value: false, Type: "boolean"},
		}
	case EntityTypeTransformer:
		return Properties{
			"name":           {Value: "Transformer", Type: "string"},
			"hv_voltage":     {Value: float64(69000), Type: "number", Unit: "V"},
			"lv_voltage":     {Value: float64(480), Type: "number", Unit: "V"},
			"rating":         {Value: float64(1000), Type: "number", Unit: "kVA"},
			"tap_position":   {Value: float64(0), Type: "number"},
		}
	case EntityTypeGenerator:
		return Properties{
			"name":               {Value: "Solar Generator", Type: "string"},
			"rated_capacity":     {Value: float64(500), Type: "number", Unit: "kW", Min: floatPtr(0)},
			"available_capacity": {Value: float64(500), Type: "number", Unit: "kW", Min: floatPtr(0)},
			"is_online":          {Value: true, Type: "boolean"},
			"is_dispatchable":    {Value: true, Type: "boolean"},
		}
	case EntityTypeLoad:
		return Properties{
			"name":                {Value: "Factory Load", Type: "string"},
			"active_power_demand": {Value: float64(400), Type: "number", Unit: "kW", Min: floatPtr(0)},
			"power_factor":        {Value: float64(0.9), Type: "number", Min: floatPtr(0), Max: floatPtr(1)},
			"is_connected":        {Value: true, Type: "boolean"},
		}
	case EntityTypeMeter:
		return Properties{
			"name":   {Value: "PCC Meter", Type: "string"},
			"type":   {Value: "pcc", Type: "enum", Options: []string{"pcc", "array", "feeder"}},
		}
	default:
		return Properties{}
	}
}

func floatPtr(v float64) *float64 {
	return &v
}

func generateEntityID() string {
	// Simple ID generation - uses time-based prefix for uniqueness
	return fmt.Sprintf("entity-%d", timeNow().UnixNano())
}
