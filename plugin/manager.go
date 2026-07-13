// Package plugin provides the Forge Plugin System.
//
// This file contains the Plugin Manager implementation.
package plugin

import (
	"fmt"
	"sort"
	"sync"
)

// PluginStatus represents the lifecycle status of a plugin.
type PluginStatus int

const (
	StatusDiscovered PluginStatus = iota
	StatusRegistered
	StatusInitialized
	StatusRunning
	StatusShutdown
	StatusUnregistered
)

func (s PluginStatus) String() string {
	switch s {
	case StatusDiscovered:
		return "discovered"
	case StatusRegistered:
		return "registered"
	case StatusInitialized:
		return "initialized"
	case StatusRunning:
		return "running"
	case StatusShutdown:
		return "shutdown"
	case StatusUnregistered:
		return "unregistered"
	default:
		return "unknown"
	}
}

// PluginInfo contains information about a loaded plugin.
type PluginInfo struct {
	Plugin Plugin
	Status PluginStatus
	Error  error
}

// Manager manages plugins in Forge.
type Manager struct {
	mu       sync.RWMutex
	plugins  map[string]*PluginInfo
	registry []Plugin // Registered plugins (static registration)
}

// NewManager creates a new plugin manager.
func NewManager() *Manager {
	return &Manager{
		plugins:  make(map[string]*PluginInfo),
		registry: make([]Plugin, 0),
	}
}

// Register registers a plugin for static loading.
//
// This supports the backwards-compatible init() pattern where
// plugins register themselves at package initialization.
func (m *Manager) Register(p Plugin) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := p.ID()
	if _, exists := m.plugins[id]; exists {
		// Plugin already registered
		return
	}

	m.plugins[id] = &PluginInfo{
		Plugin: p,
		Status: StatusRegistered,
	}
	m.registry = append(m.registry, p)
}

// Load initializes all registered plugins.
//
// This implements static loading. Dynamic loading can be added later
// without changing the contract.
func (m *Manager) Load(ctx Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Sort plugins by dependencies
	sorted, err := m.sortByDependencies()
	if err != nil {
		return fmt.Errorf("plugin dependency error: %w", err)
	}

	// Initialize each plugin
	for _, info := range sorted {
		if info.Status != StatusRegistered {
			continue
		}

		// Check dependencies are initialized
		for _, depID := range info.Plugin.Dependencies() {
			depInfo, exists := m.plugins[depID]
			if !exists {
				return fmt.Errorf("plugin %s requires missing dependency %s", info.Plugin.ID(), depID)
			}
			if depInfo.Status != StatusInitialized && depInfo.Status != StatusRunning {
				return fmt.Errorf("plugin %s requires uninitialized dependency %s", info.Plugin.ID(), depID)
			}
		}

		// Initialize plugin
		if err := info.Plugin.OnInit(ctx); err != nil {
			info.Status = StatusShutdown
			info.Error = fmt.Errorf("initialization failed: %w", err)
			return fmt.Errorf("failed to initialize plugin %s: %w", info.Plugin.ID(), err)
		}

		info.Status = StatusInitialized
	}

	// Mark all as running
	for _, info := range m.plugins {
		if info.Status == StatusInitialized {
			info.Status = StatusRunning
		}
	}

	return nil
}

// Unload shuts down and unregisters all plugins.
func (m *Manager) Unload() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Shutdown in reverse dependency order
	sorted := m.sortPluginsByDependencyOrder()

	for _, info := range sorted {
		if info.Status != StatusRunning && info.Status != StatusInitialized {
			continue
		}

		if err := info.Plugin.OnShutdown(); err != nil {
			info.Error = fmt.Errorf("shutdown failed: %w", err)
			// Continue shutting down other plugins
		}
		info.Status = StatusShutdown
	}

	return nil
}

// Get retrieves a plugin by ID.
func (m *Manager) Get(id string) Plugin {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if info, exists := m.plugins[id]; exists {
		return info.Plugin
	}
	return nil
}

// List returns all registered plugins.
func (m *Manager) List() []Plugin {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugins := make([]Plugin, 0, len(m.plugins))
	for _, info := range m.plugins {
		plugins = append(plugins, info.Plugin)
	}
	return plugins
}

// ListByStatus returns plugins filtered by status.
func (m *Manager) ListByStatus(status PluginStatus) []Plugin {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugins := make([]Plugin, 0)
	for _, info := range m.plugins {
		if info.Status == status {
			plugins = append(plugins, info.Plugin)
		}
	}
	return plugins
}

// Status returns the status of a plugin.
func (m *Manager) Status(id string) PluginStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if info, exists := m.plugins[id]; exists {
		return info.Status
	}
	return StatusUnregistered
}

// Info returns information about a plugin.
func (m *Manager) Info(id string) *PluginInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.plugins[id]
}

// sortByDependencies returns plugins sorted by dependency order.
func (m *Manager) sortByDependencies() ([]*PluginInfo, error) {
	// Build dependency graph
	graph := make(map[string][]string)
	allIDs := make([]string, 0)
	for id, info := range m.plugins {
		graph[id] = info.Plugin.Dependencies()
		allIDs = append(allIDs, id)
	}

	// Topological sort using Kahn's algorithm
	inDegree := make(map[string]int)
	for id := range m.plugins {
		inDegree[id] = 0
	}
	for _, deps := range graph {
		for _, dep := range deps {
			inDegree[dep]++
		}
	}

	// Find nodes with no incoming edges
	queue := make([]string, 0)
	for id, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, id)
		}
	}
	sort.Strings(queue)

	result := make([]*PluginInfo, 0, len(m.plugins))
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		result = append(result, m.plugins[id])

		for _, dep := range graph[id] {
			inDegree[dep]--
			if inDegree[dep] == 0 {
				queue = append(queue, dep)
				sort.Strings(queue)
			}
		}
	}

	// Check for cycles
	if len(result) != len(m.plugins) {
		return nil, fmt.Errorf("circular plugin dependency detected")
	}

	return result, nil
}

// sortPluginsByDependencyOrder returns plugins in reverse dependency order.
func (m *Manager) sortPluginsByDependencyOrder() []*PluginInfo {
	sorted, _ := m.sortByDependencies()

	// Reverse the order
	result := make([]*PluginInfo, len(sorted))
	for i, info := range sorted {
		result[len(sorted)-1-i] = info
	}
	return result
}
