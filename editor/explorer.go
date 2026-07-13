// Package editor provides the Forge Editor for creating and editing simulation models.
package editor

import "github.com/tamzrod/forge/world"

// Explorer builds the project explorer tree.
type Explorer struct {
	project *Project
	world   *world.World
}

// NewExplorer creates a new explorer.
func NewExplorer(project *Project) *Explorer {
	return &Explorer{
		project: project,
	}
}

// BuildTree builds the project explorer tree.
func (e *Explorer) BuildTree() *TreeNode {
	root := &TreeNode{
		ID:    "project",
		Label: e.project.Name,
		Icon:  "📁",
		Type:  "project",
		Children: []*TreeNode{
			e.buildWorldNode(),
			e.buildTopologyNode(),
			e.buildEntitiesNode(),
			e.buildScenariosNode(),
			e.buildSimulationNode(),
		},
		Expanded: true,
	}
	return root
}

func (e *Explorer) buildWorldNode() *TreeNode {
	return &TreeNode{
		ID:    "world",
		Label: "World",
		Icon:  "🌍",
		Type:  "world",
		Children: []*TreeNode{
			{
				ID:    "world-clock",
				Label: "Clock",
				Icon:  "⏱️",
				Type:  "clock",
			},
			{
				ID:    "world-solver",
				Label: "Solver",
				Icon:  "⚙️",
				Type:  "solver",
			},
		},
		Expanded: false,
	}
}

func (e *Explorer) buildTopologyNode() *TreeNode {
	children := make([]*TreeNode, 0)

	// Group entities by type
	buses := make([]*TreeNode, 0)
	branches := make([]*TreeNode, 0)
	switches := make([]*TreeNode, 0)

	for _, entity := range e.project.Entities {
		node := &TreeNode{
			ID:    string(entity.ID),
			Label: entity.Name,
			Icon:  e.getEntityIcon(entity.EntityType),
			Type:  string(entity.EntityType),
			Data:  entity,
		}

		switch entity.EntityType {
		case EntityTypeBus:
			buses = append(buses, node)
		case EntityTypeBreaker, EntityTypeTransformer:
			branches = append(branches, node)
		case EntityTypeGrid:
			switches = append(switches, node)
		}
	}

	if len(buses) > 0 {
		children = append(children, &TreeNode{
			ID:       "topology-buses",
			Label:    "Buses",
			Icon:     "⚫",
			Type:     "category",
			Children: buses,
		})
	}

	if len(branches) > 0 {
		children = append(children, &TreeNode{
			ID:       "topology-branches",
			Label:    "Branches",
			Icon:     "🔗",
			Type:     "category",
			Children: branches,
		})
	}

	if len(switches) > 0 {
		children = append(children, &TreeNode{
			ID:       "topology-switches",
			Label:    "Grid Connections",
			Icon:     "🔌",
			Type:     "category",
			Children: switches,
		})
	}

	return &TreeNode{
		ID:       "topology",
		Label:    "Topology",
		Icon:     "🔗",
		Type:     "topology",
		Children: children,
	}
}

func (e *Explorer) buildEntitiesNode() *TreeNode {
	children := make([]*TreeNode, 0)

	// Group entities by type
	generators := make([]*TreeNode, 0)
	loads := make([]*TreeNode, 0)
	meters := make([]*TreeNode, 0)

	for _, entity := range e.project.Entities {
		node := &TreeNode{
			ID:    string(entity.ID),
			Label: entity.Name,
			Icon:  e.getEntityIcon(entity.EntityType),
			Type:  string(entity.EntityType),
			Data:  entity,
		}

		switch entity.EntityType {
		case EntityTypeGenerator:
			generators = append(generators, node)
		case EntityTypeLoad:
			loads = append(loads, node)
		case EntityTypeMeter:
			meters = append(meters, node)
		}
	}

	if len(generators) > 0 {
		children = append(children, &TreeNode{
			ID:       "entities-generators",
			Label:    "Generators",
			Icon:     "☀️",
			Type:     "category",
			Children: generators,
		})
	}

	if len(loads) > 0 {
		children = append(children, &TreeNode{
			ID:       "entities-loads",
			Label:    "Loads",
			Icon:     "🏭",
			Type:     "category",
			Children: loads,
		})
	}

	if len(meters) > 0 {
		children = append(children, &TreeNode{
			ID:       "entities-meters",
			Label:    "Meters",
			Icon:     "📊",
			Type:     "category",
			Children: meters,
		})
	}

	return &TreeNode{
		ID:       "entities",
		Label:    "Entities",
		Icon:     "📦",
		Type:     "entities",
		Children: children,
	}
}

func (e *Explorer) buildScenariosNode() *TreeNode {
	return &TreeNode{
		ID:    "scenarios",
		Label: "Scenarios",
		Icon:  "🎬",
		Type:  "scenarios",
		Children: []*TreeNode{
			{
				ID:    "scenarios-default",
				Label: "Default Scenario",
				Icon:  "🎭",
				Type:  "scenario",
			},
		},
	}
}

func (e *Explorer) buildSimulationNode() *TreeNode {
	return &TreeNode{
		ID:    "simulation",
		Label: "Simulation",
		Icon:  "▶️",
		Type:  "simulation",
		Children: []*TreeNode{
			{
				ID:    "sim-controls",
				Label: "Controls",
				Icon:  "🎛️",
				Type:  "controls",
			},
			{
				ID:    "sim-clock",
				Label: "Clock",
				Icon:  "⏰",
				Type:  "clock",
			},
		},
	}
}

func (e *Explorer) getEntityIcon(entityType EntityType) string {
	switch entityType {
	case EntityTypeGrid:
		return "🔌"
	case EntityTypeBus:
		return "⚫"
	case EntityTypeBreaker:
		return "🔀"
	case EntityTypeTransformer:
		return "🔄"
	case EntityTypeGenerator:
		return "☀️"
	case EntityTypeLoad:
		return "🏭"
	case EntityTypeMeter:
		return "📊"
	case EntityTypeSun:
		return "🌞"
	case EntityTypeWeather:
		return "🌤️"
	case EntityTypeWind:
		return "💨"
	case EntityTypeScenario:
		return "🎬"
	case EntityTypeClock:
		return "⏱️"
	default:
		return "📦"
	}
}

// GetEntityNode returns the tree node for an entity.
func (e *Explorer) GetEntityNode(entityID ID) *TreeNode {
	for _, entity := range e.project.Entities {
		if entity.ID == entityID {
			return &TreeNode{
				ID:    string(entity.ID),
				Label: entity.Name,
				Icon:  e.getEntityIcon(entity.EntityType),
				Type:  string(entity.EntityType),
				Data:  entity,
			}
		}
	}
	return nil
}
