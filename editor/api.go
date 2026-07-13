// Package editor provides the Forge Editor for creating and editing simulation models.
package editor

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Handler handles editor HTTP requests.
type Handler struct {
	editor *Editor
}

// Editor is the main editor state manager.
type Editor struct {
	project  *Project
	palette  *Palette
	explorer *Explorer
}

// NewEditor creates a new editor instance.
func NewEditor() *Editor {
	return &Editor{
		project: NewProject("New Project"),
		palette: DefaultPalette,
		explorer: NewExplorer(nil),
	}
}

// NewHandler creates a new editor handler.
func NewHandler(e *Editor) *Handler {
	e.explorer = NewExplorer(e.project)
	return &Handler{editor: e}
}

// ServeHTTP handles editor requests.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.URL.Path {
	case "/api/editor/state":
		h.handleGetState(w, r)
	case "/api/editor/project":
		h.handleProject(w, r)
	case "/api/editor/entity":
		h.handleEntity(w, r)
	case "/api/editor/connection":
		h.handleConnection(w, r)
	case "/api/editor/palette":
		h.handlePalette(w, r)
	case "/api/editor/explorer":
		h.handleExplorer(w, r)
	case "/api/editor/simulation":
		h.handleSimulation(w, r)
	default:
		http.NotFound(w, r)
	}
}

// handleGetState returns the full editor state.
func (h *Handler) handleGetState(w http.ResponseWriter, r *http.Request) {
	state := &EditorState{
		Project:    h.editor.project,
		Selection:  &Selection{},
		Inspector:  &InspectorState{},
		IsModified: false,
		IsRunning:  false,
		IsPaused:   false,
		Speed:      1.0,
		CurrentTime: "00:00:00",
	}

	h.writeJSON(w, state)
}

// handleProject handles project operations.
func (h *Handler) handleProject(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.writeJSON(w, h.editor.project)
	case http.MethodPost:
		var req struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		h.editor.project = NewProject(req.Name)
		h.editor.explorer = NewExplorer(h.editor.project)
		h.writeJSON(w, h.editor.project)
	case http.MethodPut:
		var project *Project
		if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		h.editor.project = project
		h.editor.explorer = NewExplorer(h.editor.project)
		h.writeJSON(w, h.editor.project)
	}
}

// handleEntity handles entity operations.
func (h *Handler) handleEntity(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id := r.URL.Query().Get("id")
		if id == "" {
			h.writeJSON(w, h.editor.project.Entities)
			return
		}
		entity := h.editor.project.GetEntity(ID(id))
		if entity == nil {
			http.NotFound(w, r)
			return
		}
		h.writeJSON(w, entity)

	case http.MethodPost:
		var req struct {
			EntityType EntityType `json:"entity_type"`
			Name       string     `json:"name"`
			Position   Point      `json:"position"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		entity := NewCanvasEntity(req.EntityType, req.Name, req.Position)
		h.editor.project.AddEntity(entity)
		h.writeJSON(w, entity)

	case http.MethodPut:
		var entity *CanvasEntity
		if err := json.NewDecoder(r.Body).Decode(&entity); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		existing := h.editor.project.GetEntity(entity.ID)
		if existing == nil {
			http.NotFound(w, r)
			return
		}
		*existing = *entity
		h.writeJSON(w, existing)

	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "id required", http.StatusBadRequest)
			return
		}
		h.editor.project.RemoveEntity(ID(id))
		h.writeJSON(w, map[string]bool{"deleted": true})
	}
}

// handleConnection handles connection operations.
func (h *Handler) handleConnection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.writeJSON(w, h.editor.project.Connections)

	case http.MethodPost:
		var conn Connection
		if err := json.NewDecoder(r.Body).Decode(&conn); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Validate connection
		validator := NewConnectionValidator(h.editor.project.Entities, h.editor.project.Connections)
		result := validator.CanConnect(conn.FromEntity, conn.ToEntity, conn.FromTerminal, conn.ToTerminal)
		if !result.Valid {
			http.Error(w, result.Message, http.StatusBadRequest)
			return
		}

		conn.ID = ID(generateID())
		h.editor.project.AddConnection(&conn)
		h.writeJSON(w, &conn)

	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "id required", http.StatusBadRequest)
			return
		}
		h.editor.project.RemoveConnection(ID(id))
		h.writeJSON(w, map[string]bool{"deleted": true})
	}
}

// handlePalette returns the editor palette.
func (h *Handler) handlePalette(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, h.editor.palette)
}

// handleExplorer returns the project explorer tree.
func (h *Handler) handleExplorer(w http.ResponseWriter, r *http.Request) {
	tree := h.editor.explorer.BuildTree()
	h.writeJSON(w, tree)
}

// handleSimulation handles simulation control operations.
func (h *Handler) handleSimulation(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.writeJSON(w, map[string]interface{}{
			"is_running":   false,
			"is_paused":    false,
			"speed":        1.0,
			"current_time": "00:00:00",
		})

	case http.MethodPost:
		var req struct {
			Action string  `json:"action"`
			Speed  float64 `json:"speed,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		switch req.Action {
		case "run":
			h.writeJSON(w, map[string]interface{}{
				"action":    "run",
				"is_running": true,
				"is_paused":  false,
			})
		case "pause":
			h.writeJSON(w, map[string]interface{}{
				"action":    "pause",
				"is_running": true,
				"is_paused":  true,
			})
		case "stop":
			h.writeJSON(w, map[string]interface{}{
				"action":     "stop",
				"is_running": false,
				"is_paused":  false,
			})
		case "reset":
			h.writeJSON(w, map[string]interface{}{
				"action":      "reset",
				"is_running":  false,
				"is_paused":   false,
				"current_time": "00:00:00",
			})
		case "speed":
			h.writeJSON(w, map[string]interface{}{
				"action": "speed",
				"speed":  req.Speed,
			})
		default:
			http.Error(w, fmt.Sprintf("Unknown action: %s", req.Action), http.StatusBadRequest)
		}
	}
}

// writeJSON writes a JSON response.
func (h *Handler) writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// generateID generates a unique ID.
func generateID() string {
	return fmt.Sprintf("id-%d", timeNow().UnixNano())
}

// timeNow returns the current time (injectable for testing).
var timeNow = func() interface{ UnixNano() int64 } {
	return &fakeTime{}
}

type fakeTime struct{}

func (f *fakeTime) UnixNano() int64 {
	return 0
}
