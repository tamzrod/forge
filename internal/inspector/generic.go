// Package inspector provides a development tool for visualizing simulation state.
//
// The Generic Inspector is a data-driven inspection framework that supports
// any inspectable object in Forge through a unified data model.
//
// Supported Objects:
//   - Simulation Models (Clock, Sun, Weather, Grid, Wind)
//   - Virtual Devices (Weather Station, Revenue Meter, PV Inverter, etc.)
//   - Virtual Firmware
//   - Device Memory
//   - Communication Interfaces
//
// Inspector Sections:
//   - Identity: Name, Type, ID
//   - Overview: High-level summary
//   - State: Current operational state
//   - Configuration: Setup parameters
//   - Diagnostics: Health and error information
//   - Communications: Interface statistics
//   - Children: Nested inspectable objects
//
// Property Types:
//   - Text: Plain string values
//   - Number: Numeric values with optional unit
//   - Boolean: True/false values
//   - Status: Operational status with health indication
//   - Timestamp: Date/time values
//   - Duration: Time durations
//   - Quality: Data quality indicators
//   - Enum: Enumeration values
//   - Nested: Nested inspectable objects
//   - List: Lists of values or objects
package inspector

import (
	"encoding/json"
	"time"
)

// ObjectType represents the category of an inspectable object.
type ObjectType string

const (
	ObjectTypeSimulation ObjectType = "simulation" // Simulation Models
	ObjectTypeDevice     ObjectType = "device"     // Virtual Devices
	ObjectTypeFirmware   ObjectType = "firmware"  // Virtual Firmware
	ObjectTypeMemory     ObjectType = "memory"    // Device Memory
	ObjectTypeInterface  ObjectType = "interface" // Communication Interface
	ObjectTypeUnknown    ObjectType = "unknown"
)

// SectionID represents an inspector section.
type SectionID string

const (
	SectionIdentity      SectionID = "identity"
	SectionOverview      SectionID = "overview"
	SectionState         SectionID = "state"
	SectionConfiguration SectionID = "configuration"
	SectionDiagnostics   SectionID = "diagnostics"
	SectionCommunications SectionID = "communications"
	SectionChildren      SectionID = "children"
	SectionMemory        SectionID = "memory"
)

// PropertyType defines how a property should be rendered.
type PropertyType string

const (
	PropertyTypeText       PropertyType = "text"
	PropertyTypeNumber     PropertyType = "number"
	PropertyTypeBoolean    PropertyType = "boolean"
	PropertyTypeStatus     PropertyType = "status"
	PropertyTypeTimestamp  PropertyType = "timestamp"
	PropertyTypeDuration  PropertyType = "duration"
	PropertyTypeQuality    PropertyType = "quality"
	PropertyTypeEnum       PropertyType = "enum"
	PropertyTypeNested     PropertyType = "nested"
	PropertyTypeList       PropertyType = "list"
	PropertyTypeAngle      PropertyType = "angle"
	PropertyTypePercentage PropertyType = "percentage"
)

// Quality represents data quality.
type Quality int

const (
	QualityGood      Quality = 0
	QualityUncertain Quality = 1
	QualityBad       Quality = 2
	QualityOffline   Quality = 3
)

func (q Quality) String() string {
	switch q {
	case QualityGood:
		return "Good"
	case QualityUncertain:
		return "Uncertain"
	case QualityBad:
		return "Bad"
	case QualityOffline:
		return "Offline"
	default:
		return "Unknown"
	}
}

// ObjectIdentity provides basic identifying information.
type ObjectIdentity struct {
	ID   string     `json:"id"`
	Type ObjectType `json:"type"`
	Name string     `json:"name"`
}

// Property represents a single inspectable property.
type Property struct {
	Name       string       `json:"name"`       // Display name
	Value      interface{}  `json:"value"`      // The actual value
	Type       PropertyType `json:"type"`       // How to render
	Unit       string       `json:"unit,omitempty"`
	Quality    Quality      `json:"quality,omitempty"`
	Precision  int          `json:"precision,omitempty"`  // Decimal places for numbers
	Options    []string     `json:"options,omitempty"`   // For enum type
	Children   []*Property  `json:"children,omitempty"`   // For nested type
	Items      []*Property  `json:"items,omitempty"`      // For list type
	SortOrder  int          `json:"sort_order,omitempty"` // Display order
	ColorFunc  string       `json:"color_func,omitempty"` // Color function name
}

// Section represents an inspector section with properties.
type Section struct {
	ID         SectionID    `json:"id"`
	Title      string       `json:"title"`
	Icon       string       `json:"icon,omitempty"`
	Properties []*Property `json:"properties,omitempty"`
	Children   []*ObjectRef `json:"children,omitempty"`
}

// ObjectRef references another inspectable object.
type ObjectRef struct {
	ID    string     `json:"id"`
	Type  ObjectType `json:"type"`
	Name  string     `json:"name"`
	Path  string     `json:"path,omitempty"` // Navigation path
}

// GenericInspectorData is the complete inspection data for an object.
type GenericInspectorData struct {
	Object   ObjectIdentity `json:"object"`
	Sections []*Section     `json:"sections"`
}

// HasSection checks if a section exists with content.
func (d *GenericInspectorData) HasSection(sectionID SectionID) bool {
	for _, s := range d.Sections {
		if s.ID == sectionID {
			return len(s.Properties) > 0 || len(s.Children) > 0
		}
	}
	return false
}

// GetSection returns a section by ID.
func (d *GenericInspectorData) GetSection(sectionID SectionID) *Section {
	for _, s := range d.Sections {
		if s.ID == sectionID {
			return s
		}
	}
	return nil
}

// ToJSON returns the JSON representation.
func (d *GenericInspectorData) ToJSON() ([]byte, error) {
	return json.Marshal(d)
}

// NewProperty creates a new property with defaults.
func NewProperty(name string, value interface{}, propType PropertyType) *Property {
	return &Property{
		Name: name,
		Value: value,
		Type: propType,
	}
}

// WithUnit sets the unit.
func (p *Property) WithUnit(unit string) *Property {
	p.Unit = unit
	return p
}

// WithPrecision sets the decimal precision.
func (p *Property) WithPrecision(precision int) *Property {
	p.Precision = precision
	return p
}

// WithQuality sets the quality.
func (p *Property) WithQuality(quality Quality) *Property {
	p.Quality = quality
	return p
}

// WithChildren sets nested children.
func (p *Property) WithChildren(children []*Property) *Property {
	p.Children = children
	return p
}

// WithItems sets list items.
func (p *Property) WithItems(items []*Property) *Property {
	p.Items = items
	return p
}

// WithOptions sets enum options.
func (p *Property) WithOptions(options []string) *Property {
	p.Options = options
	return p
}

// WithColorFunc sets a color function.
func (p *Property) WithColorFunc(colorFunc string) *Property {
	p.ColorFunc = colorFunc
	return p
}

// WithSortOrder sets the display order.
func (p *Property) WithSortOrder(order int) *Property {
	p.SortOrder = order
	return p
}

// Property helpers for common types
func TextProperty(name, value string) *Property {
	return NewProperty(name, value, PropertyTypeText)
}

func NumberProperty(name string, value float64, unit string) *Property {
	return NewProperty(name, value, PropertyTypeNumber).WithUnit(unit)
}

func IntProperty(name string, value int64) *Property {
	return NewProperty(name, value, PropertyTypeNumber)
}

func BoolProperty(name string, value bool) *Property {
	return NewProperty(name, value, PropertyTypeBoolean)
}

func StatusProperty(name, value string) *Property {
	return NewProperty(name, value, PropertyTypeStatus)
}

func TimestampProperty(name string, value time.Time) *Property {
	return NewProperty(name, value.Unix(), PropertyTypeTimestamp)
}

func DurationProperty(name string, value time.Duration) *Property {
	return NewProperty(name, value.Nanoseconds(), PropertyTypeDuration)
}

func QualityProperty(name string, value Quality) *Property {
	return NewProperty(name, value, PropertyTypeQuality)
}

func EnumProperty(name, value string, options []string) *Property {
	return NewProperty(name, value, PropertyTypeEnum).WithOptions(options)
}

func NestedProperty(name string, children []*Property) *Property {
	return NewProperty(name, nil, PropertyTypeNested).WithChildren(children)
}

func ListProperty(name string, items []*Property) *Property {
	return NewProperty(name, len(items), PropertyTypeList).WithItems(items)
}

func AngleProperty(name string, degrees float64) *Property {
	return NewProperty(name, degrees, PropertyTypeAngle).WithUnit("°")
}

func PercentageProperty(name string, value float64) *Property {
	return NewProperty(name, value*100, PropertyTypePercentage).WithUnit("%")
}

// NumberWithColor creates a number property with conditional coloring based on value ranges.
func NumberWithColor(name string, value float64, unit string, goodRange, warnRange [2]float64, colorFunc string) *Property {
	p := NumberProperty(name, value, unit)
	p.ColorFunc = colorFunc
	return p
}
