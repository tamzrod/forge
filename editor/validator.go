// Package editor provides the Forge Editor for creating and editing simulation models.
package editor

import (
	"fmt"
)

// ValidationResult represents the result of a validation check.
type ValidationResult struct {
	Valid   bool     `json:"valid"`
	Message string   `json:"message,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

// ConnectionValidator validates connections based on topology rules.
type ConnectionValidator struct {
	entities    []*CanvasEntity
	connections []*Connection
}

// NewConnectionValidator creates a new connection validator.
func NewConnectionValidator(entities []*CanvasEntity, connections []*Connection) *ConnectionValidator {
	return &ConnectionValidator{
		entities:    entities,
		connections: connections,
	}
}

// CanConnect checks if two entities can be connected.
func (v *ConnectionValidator) CanConnect(fromID, toID ID, fromTerminal, toTerminal string) *ValidationResult {
	fromEntity := v.getEntity(fromID)
	toEntity := v.getEntity(toID)

	if fromEntity == nil {
		return &ValidationResult{Valid: false, Message: fmt.Sprintf("Source entity %s not found", fromID)}
	}
	if toEntity == nil {
		return &ValidationResult{Valid: false, Message: fmt.Sprintf("Destination entity %s not found", toID)}
	}

	// Validation rules based on entity types and roles

	// Rule 1: Bus can connect to any entity (through terminal on entity side)
	if fromEntity.EntityType == EntityTypeBus && isEntityWithTerminal(toEntity.EntityType) {
		return &ValidationResult{Valid: true}
	}
	if toEntity.EntityType == EntityTypeBus && isEntityWithTerminal(fromEntity.EntityType) {
		return &ValidationResult{Valid: true}
	}

	// Rule 2: Bus to Bus connection (through breaker or transformer)
	if fromEntity.EntityType == EntityTypeBus && toEntity.EntityType == EntityTypeBus {
		return v.validateBusToBus(fromEntity, toEntity)
	}

	// Rule 3: Grid connects to Bus
	if fromEntity.EntityType == EntityTypeGrid && toEntity.EntityType == EntityTypeBus {
		return v.validateGridToBus(fromEntity, toEntity)
	}
	if toEntity.EntityType == EntityTypeGrid && fromEntity.EntityType == EntityTypeBus {
		return v.validateGridToBus(toEntity, fromEntity)
	}

	// Rule 4: Generator connects to Bus
	if fromEntity.EntityType == EntityTypeGenerator && toEntity.EntityType == EntityTypeBus {
		return v.validateGeneratorToBus(fromEntity, toEntity)
	}
	if toEntity.EntityType == EntityTypeGenerator && fromEntity.EntityType == EntityTypeBus {
		return v.validateGeneratorToBus(toEntity, fromEntity)
	}

	// Rule 5: Load connects to Bus
	if fromEntity.EntityType == EntityTypeLoad && toEntity.EntityType == EntityTypeBus {
		return v.validateLoadToBus(fromEntity, toEntity)
	}
	if toEntity.EntityType == EntityTypeLoad && fromEntity.EntityType == EntityTypeBus {
		return v.validateLoadToBus(toEntity, fromEntity)
	}

	// Rule 6: Meter connects to Bus (observation)
	if fromEntity.EntityType == EntityTypeMeter && toEntity.EntityType == EntityTypeBus {
		return v.validateMeterToBus(fromEntity, toEntity)
	}
	if toEntity.EntityType == EntityTypeMeter && fromEntity.EntityType == EntityTypeBus {
		return v.validateMeterToBus(toEntity, fromEntity)
	}

	// Rule 7: Breaker connects Bus to Bus
	if fromEntity.EntityType == EntityTypeBreaker && toEntity.EntityType == EntityTypeBus {
		return v.validateBreakerToBus(fromEntity, toEntity)
	}

	// Rule 8: Transformer connects Bus to Bus
	if fromEntity.EntityType == EntityTypeTransformer && toEntity.EntityType == EntityTypeBus {
		return v.validateTransformerToBus(fromEntity, toEntity)
	}

	return &ValidationResult{Valid: true}
}

// validateBusToBus validates connection between two buses.
func (v *ConnectionValidator) validateBusToBus(from, to *CanvasEntity) *ValidationResult {
	// Buses can only connect through a transformer or breaker
	return &ValidationResult{
		Valid:   false,
		Message: "Buses must be connected through a breaker or transformer",
	}
}

// validateGridToBus validates utility grid connection.
func (v *ConnectionValidator) validateGridToBus(grid, bus *CanvasEntity) *ValidationResult {
	// Grid should connect to HV bus
	return &ValidationResult{Valid: true}
}

// validateGeneratorToBus validates generator connection.
func (v *ConnectionValidator) validateGeneratorToBus(gen, bus *CanvasEntity) *ValidationResult {
	// Generators typically connect to LV buses
	return &ValidationResult{Valid: true}
}

// validateLoadToBus validates load connection.
func (v *ConnectionValidator) validateLoadToBus(load, bus *CanvasEntity) *ValidationResult {
	return &ValidationResult{Valid: true}
}

// validateMeterToBus validates meter observation connection.
func (v *ConnectionValidator) validateMeterToBus(meter, bus *CanvasEntity) *ValidationResult {
	return &ValidationResult{Valid: true}
}

// validateBreakerToBus validates breaker connection.
func (v *ConnectionValidator) validateBreakerToBus(breaker, bus *CanvasEntity) *ValidationResult {
	return &ValidationResult{Valid: true}
}

// validateTransformerToBus validates transformer connection.
func (v *ConnectionValidator) validateTransformerToBus(tx, bus *CanvasEntity) *ValidationResult {
	return &ValidationResult{Valid: true}
}

// ValidateVoltageMatch checks if voltage levels match for connection.
func (v *ConnectionValidator) ValidateVoltageMatch(entityID ID, voltage float64) *ValidationResult {
	entity := v.getEntity(entityID)
	if entity == nil {
		return &ValidationResult{Valid: false, Message: "Entity not found"}
	}

	// Get entity's nominal voltage
	entityVoltage, ok := entity.Properties["nominal_voltage"].Value.(float64)
	if !ok {
		// Try hv_voltage for transformers
		entityVoltage, ok = entity.Properties["hv_voltage"].Value.(float64)
		if !ok {
			return &ValidationResult{Valid: true} // No voltage to validate
		}
	}

	// Check if voltages match
	if entityVoltage != voltage {
		return &ValidationResult{
			Valid:   false,
			Message: fmt.Sprintf("Voltage mismatch: %.0fV vs %.0fV", entityVoltage, voltage),
		}
	}

	return &ValidationResult{Valid: true}
}

// ValidateConnection validates an existing connection.
func (v *ConnectionValidator) ValidateConnection(conn *Connection) *ValidationResult {
	return v.CanConnect(conn.FromEntity, conn.ToEntity, conn.FromTerminal, conn.ToTerminal)
}

// IsCyclic checks if adding a connection would create a cycle.
func (v *ConnectionValidator) IsCyclic(fromID, toID ID) bool {
	// Build adjacency list
	adj := make(map[ID][]ID)
	for _, conn := range v.connections {
		adj[conn.FromEntity] = append(adj[conn.FromEntity], conn.ToEntity)
	}

	// Add proposed connection
	adj[fromID] = append(adj[fromID], toID)

	// DFS to detect cycle
	visited := make(map[ID]bool)
	recStack := make(map[ID]bool)

	var dfs func(node ID) bool
	dfs = func(node ID) bool {
		visited[node] = true
		recStack[node] = true

		for _, neighbor := range adj[node] {
			if !visited[neighbor] {
				if dfs(neighbor) {
					return true
				}
			} else if recStack[neighbor] {
				return true
			}
		}

		recStack[node] = false
		return false
	}

	for entity := range adj {
		if !visited[entity] {
			if dfs(entity) {
				return true
			}
		}
	}

	return false
}

// getEntity returns an entity by ID.
func (v *ConnectionValidator) getEntity(id ID) *CanvasEntity {
	for _, e := range v.entities {
		if e.ID == id {
			return e
		}
	}
	return nil
}

// getTerminalRole returns the role of a terminal based on entity type.
func (v *ConnectionValidator) getTerminalRole(entityType EntityType, terminalName string) string {
	switch entityType {
	case EntityTypeGrid:
		return "source"
	case EntityTypeGenerator:
		return "source"
	case EntityTypeLoad:
		return "destination"
	case EntityTypeMeter:
		return "observation"
	case EntityTypeBus:
		return "through"
	case EntityTypeBreaker:
		return "through"
	case EntityTypeTransformer:
		return "through"
	default:
		return "unknown"
	}
}

// isEntityWithTerminal returns true if entity type has terminals.
func isEntityWithTerminal(entityType EntityType) bool {
	switch entityType {
	case EntityTypeGrid, EntityTypeGenerator, EntityTypeLoad, EntityTypeMeter,
		EntityTypeBreaker, EntityTypeTransformer:
		return true
	default:
		return false
	}
}

// GetValidConnections returns all valid connection targets for an entity.
func (v *ConnectionValidator) GetValidConnections(entityID ID) []*CanvasEntity {
	entity := v.getEntity(entityID)
	if entity == nil {
		return nil
	}

	var validTargets []*CanvasEntity

	for _, e := range v.entities {
		if e.ID == entityID {
			continue
		}

		// Check if this is a valid target
		result := v.CanConnect(entityID, e.ID, "", "")
		if result.Valid {
			validTargets = append(validTargets, e)
		}
	}

	return validTargets
}
