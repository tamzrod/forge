// Package forgeenvironment provides environment components for Forge.
package forgeenvironment

import "github.com/tamzrod/forge/registry/types"

// Components returns all environment components.
var Components = []types.Component{
	{
		ID:          "forge-environment:sun",
		Name:        "Sun",
		Category:    "environment",
		Icon:        "🌞",
		Description: "Solar position and irradiance",
		Properties: []types.PropertyDef{
			{Key: "name", Label: "Name", Type: "string", Default: "Sun", Required: true},
			{Key: "latitude", Label: "Latitude", Type: "number", Default: float64(35.2271)},
			{Key: "longitude", Label: "Longitude", Type: "number", Default: float64(-80.8431)},
			{Key: "tilt", Label: "Panel Tilt", Type: "number", Default: float64(20), Unit: "°"},
			{Key: "azimuth", Label: "Azimuth", Type: "number", Default: float64(180), Unit: "°"},
		},
		Terminals: []types.TerminalDef{
			{ID: "output", Name: "Irradiance", Role: "source", Direction: "output"},
		},
		Width:  60,
		Height: 60,
	},
	{
		ID:          "forge-environment:weather",
		Name:        "Weather",
		Category:    "environment",
		Icon:        "🌤️",
		Description: "Weather conditions",
		Properties: []types.PropertyDef{
			{Key: "name", Label: "Name", Type: "string", Default: "Weather", Required: true},
			{Key: "temperature", Label: "Temperature", Type: "number", Default: float64(25), Unit: "°C"},
			{Key: "humidity", Label: "Humidity", Type: "number", Default: float64(50), Unit: "%"},
			{Key: "cloud_cover", Label: "Cloud Cover", Type: "number", Default: float64(0), Unit: "%"},
		},
		Terminals: []types.TerminalDef{
			{ID: "output", Name: "Conditions", Role: "observation", Direction: "output"},
		},
		Width:  60,
		Height: 60,
	},
	{
		ID:          "forge-environment:wind",
		Name:        "Wind",
		Category:    "environment",
		Icon:        "💨",
		Description: "Wind conditions",
		Properties: []types.PropertyDef{
			{Key: "name", Label: "Name", Type: "string", Default: "Wind", Required: true},
			{Key: "speed", Label: "Speed", Type: "number", Default: float64(5), Unit: "m/s"},
			{Key: "direction", Label: "Direction", Type: "number", Default: float64(0), Unit: "°"},
		},
		Terminals: []types.TerminalDef{
			{ID: "output", Name: "Wind Data", Role: "source", Direction: "output"},
		},
		Width:  60,
		Height: 60,
	},
}

// Category is the environment category.
var Category = types.CategoryDef{
	ID:       "environment",
	Name:     "Environment",
	Icon:     "🌤️",
	Order:    2,
	Domain:   "forge-environment",
}
