// Package types provides shared types for the component registry.
package types

// PropertyDef defines a property.
type PropertyDef struct {
	Key      string
	Label    string
	Type     string
	Default  interface{}
	Unit     string
	Min      *float64
	Max      *float64
	Options  []string
	ReadOnly bool
	Required bool
}

// TerminalDef defines a terminal.
type TerminalDef struct {
	ID        string
	Name      string
	Role      string
	Voltage   *float64
	Direction string
}

// CategoryDef defines a category.
type CategoryDef struct {
	ID         string
	Name       string
	Icon       string
	Order      int
	Domain     string
	Expandable bool
}

// Component defines a component.
type Component struct {
	ID          string
	Name        string
	Category    string
	Icon        string
	Description string
	Properties  []PropertyDef
	Terminals   []TerminalDef
	Width       float64
	Height      float64
}

func FloatPtr(v float64) *float64 {
	return &v
}
