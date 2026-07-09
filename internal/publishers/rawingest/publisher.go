package rawingest

import (
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// State represents the interface connection state.
type State int

const (
	StateDisconnected State = iota
	StateConnecting
	StateConnected
	StateError
)

func (s State) String() string {
	switch s {
	case StateDisconnected:
		return "Disconnected"
	case StateConnecting:
		return "Connecting"
	case StateConnected:
		return "Connected"
	case StateError:
		return "Error"
	default:
		return "Unknown"
	}
}

// Stats holds interface statistics.
type Stats struct {
	PacketsSent   uint64
	BytesSent     uint64
	Errors        uint64
	LastSent      time.Time
	LastError     string
	LastErrorTime time.Time
}

// Interface represents a Raw Ingest communication interface.
//
// The Interface is a communication channel that serializes Device Memory
// and sends it to MMA2. It is NOT responsible for:
// - Engineering calculations
// - Value scaling
// - Device state management
// - Accessing Simulation Models
//
// The Interface simply reads Device Memory and sends bytes.
//
// Architecture:
//   Virtual Firmware
//           ↓
//   Device Memory (owned by firmware)
//           ↓
//   Interface (serializes memory)
//           ↓
//   MMA2
type Interface struct {
	config   Config
	conn     net.Conn
	sequence uint16

	mu    sync.RWMutex
	state State
	stats Stats

	stopCh  chan struct{}
	doneCh  chan struct{}
	running atomic.Bool

	onStateChange func(State)
	onStatsUpdate func(Stats)
}

// NewInterface creates a new Raw Ingest interface.
func NewInterface(cfg Config) (*Interface, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &Interface{
		config: cfg,
		state:  StateDisconnected,
		stopCh: make(chan struct{}),
		doneCh: make(chan struct{}),
	}, nil
}

// Start begins the interface.
func (i *Interface) Start() error {
	if !i.config.Enabled {
		log.Printf("RawIngest Interface: disabled")
		return nil
	}

	if i.running.Load() {
		return fmt.Errorf("interface already running")
	}

	i.running.Store(true)
	go i.run()

	return nil
}

// Stop stops the interface.
func (i *Interface) Stop() error {
	if !i.running.Load() {
		return nil
	}

	close(i.stopCh)
	<-i.doneCh

	i.mu.Lock()
	if i.conn != nil {
		i.conn.Close()
		i.conn = nil
	}
	i.mu.Unlock()

	i.running.Store(false)
	return nil
}

// Publish serializes Device Memory and sends to MMA2.
func (i *Interface) Publish(values map[string]float64) error {
	if !i.config.Enabled {
		return nil
	}

	if !i.running.Load() {
		return fmt.Errorf("interface not running")
	}

	i.mu.RLock()
	state := i.state
	conn := i.conn
	i.mu.RUnlock()

	if state != StateConnected || conn == nil {
		return fmt.Errorf("not connected")
	}

	// Encode values (no engineering calculations here)
	data := EncodeMemoryValues(values)

	// Create packet
	seq := atomic.AddUint16(&i.sequence, 1)
	timestamp := uint64(time.Now().UnixNano())

	packet, err := EncodePacket(i.config.UnitID, seq, timestamp, TypeData, data)
	if err != nil {
		i.recordError(fmt.Errorf("failed to encode packet: %w", err))
		return err
	}

	// Send packet
	if err := conn.SetWriteDeadline(time.Now().Add(i.config.Timeout)); err != nil {
		i.recordError(fmt.Errorf("failed to set deadline: %w", err))
		return err
	}

	n, err := conn.Write(packet)
	if err != nil {
		i.recordError(fmt.Errorf("failed to write: %w", err))
		i.disconnect()
		return err
	}

	// Update stats
	i.mu.Lock()
	i.stats.PacketsSent++
	i.stats.BytesSent += uint64(n)
	i.stats.LastSent = time.Now()
	i.mu.Unlock()

	if i.onStatsUpdate != nil {
		i.onStatsUpdate(i.Stats())
	}

	return nil
}

// Stats returns a copy of current statistics.
func (i *Interface) Stats() Stats {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.stats
}

// State returns the current connection state.
func (i *Interface) State() State {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.state
}

// IsConnected returns true if connected.
func (i *Interface) IsConnected() bool {
	return i.State() == StateConnected
}

// OnStateChange sets a callback for state changes.
func (i *Interface) OnStateChange(f func(State)) {
	i.onStateChange = f
}

// OnStatsUpdate sets a callback for stats updates.
func (i *Interface) OnStatsUpdate(f func(Stats)) {
	i.onStatsUpdate = f
}

// run is the main interface loop.
func (i *Interface) run() {
	defer close(i.doneCh)

	ticker := time.NewTicker(i.config.Interval)
	defer ticker.Stop()

	retries := 0

	for {
		// Ensure connection
		if !i.IsConnected() {
			if err := i.connect(); err != nil {
				log.Printf("RawIngest Interface: Connection failed: %v", err)

				retries++
				if retries <= i.config.MaxRetries {
					delay := i.config.ReconnectDelay * time.Duration(retries)
					log.Printf("RawIngest Interface: Retrying in %v (attempt %d/%d)", delay, retries, i.config.MaxRetries)

					select {
					case <-i.stopCh:
						return
					case <-time.After(delay):
						continue
					}
				} else {
					log.Printf("RawIngest Interface: Max retries reached")
					retries = 0
				}
			} else {
				retries = 0
			}
		}

		select {
		case <-i.stopCh:
			return
		case <-ticker.C:
			// Ticker fires, but actual publishing is done externally
			// via Publish() calls. This keeps the connection alive.
		}
	}
}

// connect establishes connection to MMA2.
func (i *Interface) connect() error {
	i.setState(StateConnecting)

	address := i.config.Address()
	log.Printf("RawIngest Interface: Connecting to %s", address)

	conn, err := net.DialTimeout("tcp", address, i.config.Timeout)
	if err != nil {
		i.setState(StateError)
		return fmt.Errorf("dial failed: %w", err)
	}

	i.mu.Lock()
	i.conn = conn
	i.mu.Unlock()

	i.setState(StateConnected)
	log.Printf("RawIngest Interface: Connected to %s", address)

	return nil
}

// disconnect closes the connection.
func (i *Interface) disconnect() {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.conn != nil {
		i.conn.Close()
		i.conn = nil
	}

	i.setState(StateDisconnected)
}

// recordError records an error and triggers reconnect.
func (i *Interface) recordError(err error) {
	i.mu.Lock()
	i.stats.Errors++
	i.stats.LastError = err.Error()
	i.stats.LastErrorTime = time.Now()
	i.mu.Unlock()

	if i.onStatsUpdate != nil {
		i.onStatsUpdate(i.Stats())
	}
}

// setState updates the connection state.
func (i *Interface) setState(state State) {
	i.mu.Lock()
	i.state = state
	i.mu.Unlock()

	if i.onStateChange != nil {
		i.onStateChange(state)
	}
}
