// Package validators provides domain-specific connection validators.
package validators

import (
	"fmt"

	"github.com/tamzrod/forge/plugin"
)

// ElectricalValidator validates electrical connections.
//
// This validator implements domain-specific rules for electrical terminals.
// It checks voltage compatibility, terminal roles, and connection rules.
type ElectricalValidator struct{}

// NewElectricalValidator creates a new electrical validator.
func NewElectricalValidator() *ElectricalValidator {
	return &ElectricalValidator{}
}

// Domain implements ConnectionValidator.
func (v *ElectricalValidator) Domain() string {
	return "electrical"
}

// CanConnect validates an electrical connection between two terminals.
//
// Rules:
//   - Through terminals (buses) can connect to anything
//   - Observation terminals can connect to any voltage level
//   - Voltage levels must match for source/destination connections
//   - Source and destination cannot directly connect without through
func (v *ElectricalValidator) CanConnect(source, target *plugin.TerminalDescriptor) (bool, error) {
	// Through terminals can connect to anything
	if source.Role == plugin.TerminalRoleThrough || target.Role == plugin.TerminalRoleThrough {
		return true, nil
	}

	// Observation terminals can connect anywhere
	if source.Role == plugin.TerminalRoleObservation || target.Role == plugin.TerminalRoleObservation {
		return true, nil
	}

	// Check voltage compatibility for source/destination connections
	if source.Voltage != nil && target.Voltage != nil {
		if *source.Voltage != *target.Voltage {
			return false, fmt.Errorf(
				"voltage mismatch: source %.0fV cannot connect to target %.0fV",
				*source.Voltage,
				*target.Voltage,
			)
		}
	}

	return true, nil
}
