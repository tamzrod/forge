// Package forgeelectrical provides electrical components for Forge.
package forgeelectrical

import "github.com/tamzrod/forge/registry/types"

// Components returns all electrical components.
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

// Category is the electrical category.
var Category = types.CategoryDef{
	ID:       "electrical",
	Name:     "Electrical",
	Icon:     "⚡",
	Order:    1,
	Domain:   "forge-electrical",
}
