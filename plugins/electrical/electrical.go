// Package electrical provides the Electrical Domain Plugin for Forge.
package electrical

import (
	"fmt"

	"github.com/tamzrod/forge/plugin"
	"github.com/tamzrod/forge/plugin/validators"
	"github.com/tamzrod/forge/registry"
	"github.com/tamzrod/forge/registry/types"
)

// Plugin implements the Forge Plugin interface for the electrical domain.
type Plugin struct {
	ctx plugin.Context
}

// Ensure Plugin implements plugin.Plugin at compile time.
var _ plugin.Plugin = (*Plugin)(nil)

// ID implements plugin.Plugin.
func (p *Plugin) ID() string {
	return "forge-electrical"
}

// Name implements plugin.Plugin.
func (p *Plugin) Name() string {
	return "Electrical Plugin"
}

// Version implements plugin.Plugin.
func (p *Plugin) Version() string {
	return "1.0.0"
}

// Description implements plugin.Plugin.
func (p *Plugin) Description() string {
	return "Electrical power distribution domain components, solvers, and validators"
}

// Dependencies implements plugin.Plugin.
func (p *Plugin) Dependencies() []string {
	return nil
}

// OnInit implements plugin.Plugin.
func (p *Plugin) OnInit(ctx plugin.Context) error {
	p.ctx = ctx

	// Register validators with the registry
	catalog := ctx.ComponentCatalog()

	// Register categories
	if err := catalog.RegisterCategory(toComponentCategory(Category)); err != nil {
		return fmt.Errorf("failed to register electrical category: %w", err)
	}

	// Register components
	for _, comp := range Components {
		if err := catalog.Register(toComponentDescriptor(comp)); err != nil {
			return fmt.Errorf("failed to register component %s: %w", comp.ID, err)
		}
	}

	// Register the electrical validator
	if vr, ok := catalog.(interface {
		RegisterValidator(v registry.ConnectionValidator)
	}); ok {
		vr.RegisterValidator(validators.NewElectricalValidator())
	}

	return nil
}

// OnShutdown implements plugin.Plugin.
func (p *Plugin) OnShutdown() error {
	return nil
}

// Components implements plugin.Plugin.
func (p *Plugin) Components() []*plugin.ComponentDescriptor {
	descriptors := make([]*plugin.ComponentDescriptor, len(Components))
	for i, comp := range Components {
		descriptors[i] = toComponentDescriptor(comp)
	}
	return descriptors
}

// Categories implements plugin.Plugin.
func (p *Plugin) Categories() []*plugin.ComponentCategory {
	return []*plugin.ComponentCategory{toComponentCategory(Category)}
}

// Validators implements plugin.Plugin.
func (p *Plugin) Validators() []plugin.ConnectionValidator {
	return []plugin.ConnectionValidator{
		validators.NewElectricalValidator(),
	}
}

// RegisterEntities implements plugin.Plugin.
func (p *Plugin) RegisterEntities(registry interface{}) {
	// Entity registration would happen here when world entities are integrated
}

// Category is the electrical category definition.
var Category = types.CategoryDef{
	ID:         "electrical",
	Name:       "Electrical",
	Icon:       "⚡",
	Order:      1,
	Domain:     "forge-electrical",
	Expandable: true,
}

// Components contains all electrical component definitions.
var Components = []types.Component{
	{
		ID:          "forge-electrical:grid",
		Name:        "Utility Grid",
		Category:    "electrical",
		Icon:        "🔌",
		Description: "Utility grid connection point",
		Properties: []types.PropertyDef{
			{Key: "name", Label: "Name", Type: "string", Default: "Utility Grid", Required: true},
			{Key: "nominal_voltage", Label: "Nominal Voltage", Type: "number", Default: float64(69000), Unit: "V"},
			{Key: "nominal_frequency", Label: "Frequency", Type: "number", Default: float64(60), Unit: "Hz", Options: []string{"50", "60"}},
		},
		Terminals: []types.TerminalDef{
			{ID: "output", Name: "Output", Role: "source", Voltage: types.FloatPtr(69000), Direction: "output"},
		},
		Width:  80,
		Height: 60,
	},
	{
		ID:          "forge-electrical:bus",
		Name:        "Bus",
		Category:    "electrical",
		Icon:        "⚫",
		Description: "Electrical bus node",
		Properties: []types.PropertyDef{
			{Key: "name", Label: "Name", Type: "string", Default: "New Bus", Required: true},
			{Key: "nominal_voltage", Label: "Nominal Voltage", Type: "number", Default: float64(480), Unit: "V"},
		},
		Terminals: []types.TerminalDef{
			{ID: "input", Name: "Input", Role: "through", Direction: "input"},
			{ID: "output", Name: "Output", Role: "through", Direction: "output"},
		},
		Width:  60,
		Height: 60,
	},
	{
		ID:          "forge-electrical:breaker",
		Name:        "Breaker",
		Category:    "electrical",
		Icon:        "🔀",
		Description: "Circuit breaker switch",
		Properties: []types.PropertyDef{
			{Key: "name", Label: "Name", Type: "string", Default: "Circuit Breaker", Required: true},
			{Key: "is_open", Label: "Open", Type: "boolean", Default: false},
			{Key: "rating", Label: "Rating", Type: "number", Default: float64(1200), Unit: "A"},
		},
		Terminals: []types.TerminalDef{
			{ID: "input", Name: "Input", Role: "through", Direction: "input"},
			{ID: "output", Name: "Output", Role: "through", Direction: "output"},
		},
		Width:  50,
		Height: 50,
	},
	{
		ID:          "forge-electrical:transformer",
		Name:        "Transformer",
		Category:    "electrical",
		Icon:        "🔄",
		Description: "Power transformer",
		Properties: []types.PropertyDef{
			{Key: "name", Label: "Name", Type: "string", Default: "Transformer", Required: true},
			{Key: "hv_voltage", Label: "HV Voltage", Type: "number", Default: float64(69000), Unit: "V"},
			{Key: "lv_voltage", Label: "LV Voltage", Type: "number", Default: float64(480), Unit: "V"},
			{Key: "rating", Label: "Rating", Type: "number", Default: float64(1000), Unit: "kVA"},
			{Key: "tap_position", Label: "Tap Position", Type: "number", Default: float64(0)},
		},
		Terminals: []types.TerminalDef{
			{ID: "hv", Name: "HV", Role: "through", Direction: "input"},
			{ID: "lv", Name: "LV", Role: "through", Direction: "output"},
		},
		Width:  80,
		Height: 60,
	},
	{
		ID:          "forge-electrical:generator",
		Name:        "Virtual Generator",
		Category:    "electrical",
		Icon:        "☀️",
		Description: "Solar or wind generator",
		Properties: []types.PropertyDef{
			{Key: "name", Label: "Name", Type: "string", Default: "Solar Generator", Required: true},
			{Key: "rated_capacity", Label: "Rated Capacity", Type: "number", Default: float64(500), Unit: "kW"},
			{Key: "available_capacity", Label: "Available Capacity", Type: "number", Default: float64(500), Unit: "kW"},
			{Key: "is_online", Label: "Online", Type: "boolean", Default: true},
			{Key: "is_dispatchable", Label: "Dispatchable", Type: "boolean", Default: true},
		},
		Terminals: []types.TerminalDef{
			{ID: "output", Name: "Output", Role: "source", Direction: "output"},
		},
		Width:  80,
		Height: 80,
	},
	{
		ID:          "forge-electrical:load",
		Name:        "Virtual Load",
		Category:    "electrical",
		Icon:        "🏭",
		Description: "Factory or facility load",
		Properties: []types.PropertyDef{
			{Key: "name", Label: "Name", Type: "string", Default: "Factory Load", Required: true},
			{Key: "active_power_demand", Label: "Active Power", Type: "number", Default: float64(400), Unit: "kW"},
			{Key: "power_factor", Label: "Power Factor", Type: "number", Default: float64(0.9)},
			{Key: "is_connected", Label: "Connected", Type: "boolean", Default: true},
		},
		Terminals: []types.TerminalDef{
			{ID: "input", Name: "Input", Role: "destination", Direction: "input"},
		},
		Width:  80,
		Height: 80,
	},
	{
		ID:          "forge-electrical:meter",
		Name:        "Meter",
		Category:    "electrical",
		Icon:        "📊",
		Description: "Power measurement meter",
		Properties: []types.PropertyDef{
			{Key: "name", Label: "Name", Type: "string", Default: "PCC Meter", Required: true},
			{Key: "meter_type", Label: "Type", Type: "enum", Default: "pcc", Options: []string{"pcc", "array", "feeder"}},
		},
		Terminals: []types.TerminalDef{
			{ID: "observation", Name: "Observation", Role: "observation", Direction: "bidirectional"},
		},
		Width:  70,
		Height: 70,
	},
}

func toComponentCategory(cat types.CategoryDef) *plugin.ComponentCategory {
	return &plugin.ComponentCategory{
		ID:         cat.ID,
		Name:       cat.Name,
		Icon:       cat.Icon,
		Order:      cat.Order,
		Domain:     cat.Domain,
		Expandable: cat.Expandable,
	}
}

func toComponentDescriptor(comp types.Component) *plugin.ComponentDescriptor {
	props := make([]plugin.PropertyDescriptor, len(comp.Properties))
	for i, p := range comp.Properties {
		props[i] = plugin.PropertyDescriptor{
			Key:      p.Key,
			Label:    p.Label,
			Type:     plugin.PropertyType(p.Type),
			Default:  p.Default,
			Unit:     p.Unit,
			Min:      p.Min,
			Max:      p.Max,
			Options:  p.Options,
			ReadOnly: p.ReadOnly,
			Required: p.Required,
		}
	}

	terms := make([]plugin.TerminalDescriptor, len(comp.Terminals))
	for i, t := range comp.Terminals {
		terms[i] = plugin.TerminalDescriptor{
			ID:        t.ID,
			Name:      t.Name,
			Role:      plugin.TerminalRole(t.Role),
			Voltage:   t.Voltage,
			Direction: plugin.TerminalDirection(t.Direction),
		}
	}

	return &plugin.ComponentDescriptor{
		ID:          comp.ID,
		Name:        comp.Name,
		Category:    comp.Category,
		Icon:        comp.Icon,
		Description: comp.Description,
		Properties:  props,
		Terminals:   terms,
		Width:       comp.Width,
		Height:      comp.Height,
		Domain:      fmt.Sprintf("forge-%s", comp.Category),
	}
}
