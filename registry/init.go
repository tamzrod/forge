// Package registry provides the Component Registry for Forge.
package registry

import (
	"fmt"

	"github.com/tamzrod/forge/registry/forge-electrical"
	"github.com/tamzrod/forge/registry/forge-environment"
	"github.com/tamzrod/forge/registry/forge-simulation"
	"github.com/tamzrod/forge/registry/types"
)

// init registers all built-in components.
func init() {
	r := GetRegistry()

	// Register categories
	r.RegisterCategory(categoryToDescriptor(forgeelectrical.Category))
	r.RegisterCategory(categoryToDescriptor(forgeenvironment.Category))
	r.RegisterCategory(categoryToDescriptor(forgesimulation.Category))

	// Register electrical components
	for _, comp := range forgeelectrical.Components {
		r.Register(componentToDescriptor(comp))
	}

	// Register environment components
	for _, comp := range forgeenvironment.Components {
		r.Register(componentToDescriptor(comp))
	}

	// Register simulation components
	for _, comp := range forgesimulation.Components {
		r.Register(componentToDescriptor(comp))
	}
}

func categoryToDescriptor(cat types.CategoryDef) *ComponentCategory {
	return &ComponentCategory{
		ID:         cat.ID,
		Name:       cat.Name,
		Icon:       cat.Icon,
		Order:      cat.Order,
		Domain:     cat.Domain,
		Expandable: cat.Expandable,
	}
}

func componentToDescriptor(comp types.Component) *ComponentDescriptor {
	props := make([]PropertyDescriptor, len(comp.Properties))
	for i, p := range comp.Properties {
		props[i] = PropertyDescriptor{
			Key:      p.Key,
			Label:    p.Label,
			Type:     PropertyType(p.Type),
			Default:  p.Default,
			Unit:     p.Unit,
			Min:      p.Min,
			Max:      p.Max,
			Options:  p.Options,
			ReadOnly: p.ReadOnly,
			Required: p.Required,
		}
	}

	terms := make([]TerminalDescriptor, len(comp.Terminals))
	for i, t := range comp.Terminals {
		terms[i] = TerminalDescriptor{
			ID:        t.ID,
			Name:      t.Name,
			Role:      TerminalRole(t.Role),
			Voltage:   t.Voltage,
			Direction: TerminalDirection(t.Direction),
		}
	}

	return &ComponentDescriptor{
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
