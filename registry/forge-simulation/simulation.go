// Package forgesimulation provides simulation components for Forge.
package forgesimulation

import "github.com/tamzrod/forge/registry/types"

// Components returns all simulation components.
var Components = []types.Component{
	{
		ID:          "forge-simulation:scenario",
		Name:        "Scenario",
		Category:    "simulation",
		Icon:        "🎬",
		Description: "Test scenario",
		Properties: []types.PropertyDef{
			{Key: "name", Label: "Name", Type: "string", Default: "Test Scenario", Required: true},
			{Key: "duration", Label: "Duration", Type: "number", Default: float64(3600), Unit: "s"},
			{Key: "description", Label: "Description", Type: "string", Default: ""},
		},
		Terminals: []types.TerminalDef{},
		Width:      80,
		Height:     60,
	},
	{
		ID:          "forge-simulation:clock",
		Name:        "Simulation Clock",
		Category:    "simulation",
		Icon:        "⏱️",
		Description: "Simulation time control",
		Properties: []types.PropertyDef{
			{Key: "name", Label: "Name", Type: "string", Default: "Clock", Required: true},
			{Key: "start_time", Label: "Start Time", Type: "string", Default: "2024-01-01T08:00:00Z"},
			{Key: "end_time", Label: "End Time", Type: "string", Default: "2024-01-01T20:00:00Z"},
			{Key: "time_step", Label: "Time Step", Type: "number", Default: float64(100), Unit: "ms"},
		},
		Terminals: []types.TerminalDef{
			{ID: "output", Name: "Time", Role: "source", Direction: "output"},
		},
		Width:  60,
		Height: 60,
	},
}

// Category is the simulation category.
var Category = types.CategoryDef{
	ID:       "simulation",
	Name:     "Simulation",
	Icon:     "🎬",
	Order:    3,
	Domain:   "forge-simulation",
}
