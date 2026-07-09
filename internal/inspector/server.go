package inspector

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development tool
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Server provides the HTTP interface for the Simulation Inspector.
type Server struct {
	view   *View
	server *http.Server
	wsHub  *wsHub

	mu sync.RWMutex
}

// Config holds server configuration.
type Config struct {
	// Address to listen on (e.g., "localhost:8080")
	Address string

	// Update interval for WebSocket broadcasts
	UpdateInterval int // milliseconds
}

// DefaultConfig returns the default configuration.
func DefaultConfig() Config {
	return Config{
		Address:        "localhost:8080",
		UpdateInterval: 100, // 10 updates per second
	}
}

// NewServer creates a new inspector server.
func NewServer(view *View, cfg Config) *Server {
	hub := newWsHub()

	mux := http.NewServeMux()
	server := &Server{
		view:   view,
		wsHub:  hub,
	}

	mux.HandleFunc("/", server.handleIndex)
	mux.HandleFunc("/api/state", server.handleState)
	mux.HandleFunc("/api/inspect/", server.handleGenericInspect)
	mux.HandleFunc("/ws", server.handleWebSocket)

	server.server = &http.Server{
		Addr:    cfg.Address,
		Handler: mux,
	}

	// Start WebSocket broadcast goroutine
	go hub.run()
	go hub.broadcastLoop(view, cfg.UpdateInterval)

	return server
}

// Start starts the inspector server.
func (s *Server) Start() error {
	log.Printf("Simulation Inspector listening on http://%s", s.server.Addr)
	return s.server.ListenAndServe()
}

// Stop stops the inspector server.
func (s *Server) Stop() error {
	s.wsHub.close()
	return s.server.Close()
}

// handleIndex serves the dashboard HTML.
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("dashboard").Parse(dashboardHTML))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, nil)
}

// handleState returns the current state as JSON.
func (s *Server) handleState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	json.NewEncoder(w).Encode(s.view.FullState())
}

// handleGenericInspect returns generic inspection data for an object.
func (s *Server) handleGenericInspect(w http.ResponseWriter, r *http.Request) {
	// Extract object ID from path: /api/inspect/{objectID}
	path := r.URL.Path
	objectID := path[len("/api/inspect/"):]
	if objectID == "" {
		objectID = "world"
	}

	generator := NewGenerator(s.view)
	data, err := generator.Inspect(objectID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	json.NewEncoder(w).Encode(data)
}

// handleWebSocket upgrades the connection and registers the client.
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &wsClient{
		hub:  s.wsHub,
		conn: conn,
		send: make(chan []byte, 256),
	}

	s.wsHub.register <- client

	go client.writePump()
	go client.readPump()
}

// wsHub manages WebSocket clients and broadcasts messages.
type wsHub struct {
	clients    map[*wsClient]bool
	broadcast  chan []byte
	register   chan *wsClient
	unregister chan *wsClient
	done       chan struct{}
}

func newWsHub() *wsHub {
	return &wsHub{
		clients:    make(map[*wsClient]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *wsClient),
		unregister: make(chan *wsClient),
		done:       make(chan struct{}),
	}
}

func (h *wsHub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		case <-h.done:
			return
		}
	}
}

func (h *wsHub) broadcastLoop(view *View, intervalMs int) {
	ticker := time.NewTicker(time.Duration(intervalMs) * time.Millisecond)
	defer ticker.Stop()

	encoder := json.NewEncoder(nil)

	for {
		select {
		case <-ticker.C:
			state := view.FullState()
			encoder = json.NewEncoder(&broadcastWriter{h.broadcast})
			encoder.Encode(state)
		case <-h.done:
			return
		}
	}
}

type broadcastWriter struct {
	ch chan []byte
}

func (w *broadcastWriter) Write(p []byte) (int, error) {
	w.ch <- p
	return len(p), nil
}

func (h *wsHub) close() {
	close(h.done)
}

// wsClient handles WebSocket communication.
type wsClient struct {
	hub  *wsHub
	conn *websocket.Conn
	send chan []byte
}

func (c *wsClient) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (c *wsClient) writePump() {
	defer c.conn.Close()

	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			break
		}
	}
}
