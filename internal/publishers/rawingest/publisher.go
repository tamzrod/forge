package rawingest

import (
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// State represents the publisher connection state.
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

// Stats holds publishing statistics.
type Stats struct {
	PacketsSent   uint64
	BytesSent     uint64
	Errors        uint64
	LastSent      time.Time
	LastError     string
	LastErrorTime time.Time
}

// Publisher publishes operational memory to MMA2 using Raw Ingest.
type Publisher struct {
	config   Config
	conn     net.Conn
	sequence uint16

	mu       sync.RWMutex
	state    State
	stats    Stats

	stopCh   chan struct{}
	doneCh   chan struct{}
	running  atomic.Bool

	onStateChange func(State)
	onStatsUpdate func(Stats)
}

// NewPublisher creates a new Raw Ingest publisher.
func NewPublisher(cfg Config) (*Publisher, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &Publisher{
		config: cfg,
		state:  StateDisconnected,
		stopCh: make(chan struct{}),
		doneCh: make(chan struct{}),
	}, nil
}

// Start begins the publishing loop.
func (p *Publisher) Start() error {
	if !p.config.Enabled {
		log.Printf("RawIngest: Publishing disabled")
		return nil
	}

	if p.running.Load() {
		return fmt.Errorf("publisher already running")
	}

	p.running.Store(true)
	go p.run()

	return nil
}

// Stop stops the publisher.
func (p *Publisher) Stop() error {
	if !p.running.Load() {
		return nil
	}

	close(p.stopCh)
	<-p.doneCh

	p.mu.Lock()
	if p.conn != nil {
		p.conn.Close()
		p.conn = nil
	}
	p.mu.Unlock()

	p.running.Store(false)
	return nil
}

// Publish sends operational memory to MMA2.
func (p *Publisher) Publish(values map[string]float64) error {
	if !p.config.Enabled {
		return nil
	}

	if !p.running.Load() {
		return fmt.Errorf("publisher not running")
	}

	p.mu.RLock()
	state := p.state
	conn := p.conn
	p.mu.RUnlock()

	if state != StateConnected || conn == nil {
		return fmt.Errorf("not connected")
	}

	// Encode values
	data := EncodeMemoryValues(values)

	// Create packet
	seq := atomic.AddUint16(&p.sequence, 1)
	timestamp := uint64(time.Now().UnixNano())

	packet, err := EncodePacket(p.config.UnitID, seq, timestamp, TypeData, data)
	if err != nil {
		p.recordError(fmt.Errorf("failed to encode packet: %w", err))
		return err
	}

	// Send packet
	if err := conn.SetWriteDeadline(time.Now().Add(p.config.Timeout)); err != nil {
		p.recordError(fmt.Errorf("failed to set deadline: %w", err))
		return err
	}

	n, err := conn.Write(packet)
	if err != nil {
		p.recordError(fmt.Errorf("failed to write: %w", err))
		p.disconnect()
		return err
	}

	// Update stats
	p.mu.Lock()
	p.stats.PacketsSent++
	p.stats.BytesSent += uint64(n)
	p.stats.LastSent = time.Now()
	p.mu.Unlock()

	if p.onStatsUpdate != nil {
		p.onStatsUpdate(p.Stats())
	}

	return nil
}

// Stats returns a copy of current statistics.
func (p *Publisher) Stats() Stats {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.stats
}

// State returns the current connection state.
func (p *Publisher) State() State {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.state
}

// IsConnected returns true if connected.
func (p *Publisher) IsConnected() bool {
	return p.State() == StateConnected
}

// OnStateChange sets a callback for state changes.
func (p *Publisher) OnStateChange(f func(State)) {
	p.onStateChange = f
}

// OnStatsUpdate sets a callback for stats updates.
func (p *Publisher) OnStatsUpdate(f func(Stats)) {
	p.onStatsUpdate = f
}

// run is the main publishing loop.
func (p *Publisher) run() {
	defer close(p.doneCh)

	ticker := time.NewTicker(p.config.Interval)
	defer ticker.Stop()

	retries := 0

	for {
		// Ensure connection
		if !p.IsConnected() {
			if err := p.connect(); err != nil {
				log.Printf("RawIngest: Connection failed: %v", err)

				retries++
				if retries <= p.config.MaxRetries {
					delay := p.config.ReconnectDelay * time.Duration(retries)
					log.Printf("RawIngest: Retrying in %v (attempt %d/%d)", delay, retries, p.config.MaxRetries)

					select {
					case <-p.stopCh:
						return
					case <-time.After(delay):
						continue
					}
				} else {
					log.Printf("RawIngest: Max retries reached, waiting for manual reconnect")
					retries = 0
				}
			} else {
				retries = 0
			}
		}

		select {
		case <-p.stopCh:
			return
		case <-ticker.C:
			// Ticker fires, but actual publishing is done externally
			// via Publish() calls. This keeps the connection alive.
		}
	}
}

// connect establishes connection to MMA2.
func (p *Publisher) connect() error {
	p.setState(StateConnecting)

	address := p.config.Address()
	log.Printf("RawIngest: Connecting to %s", address)

	conn, err := net.DialTimeout("tcp", address, p.config.Timeout)
	if err != nil {
		p.setState(StateError)
		return fmt.Errorf("dial failed: %w", err)
	}

	p.mu.Lock()
	p.conn = conn
	p.mu.Unlock()

	p.setState(StateConnected)
	log.Printf("RawIngest: Connected to %s", address)

	return nil
}

// disconnect closes the connection.
func (p *Publisher) disconnect() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.conn != nil {
		p.conn.Close()
		p.conn = nil
	}

	p.setState(StateDisconnected)
}

// recordError records an error and triggers reconnect.
func (p *Publisher) recordError(err error) {
	p.mu.Lock()
	p.stats.Errors++
	p.stats.LastError = err.Error()
	p.stats.LastErrorTime = time.Now()
	p.mu.Unlock()

	if p.onStatsUpdate != nil {
		p.onStatsUpdate(p.Stats())
	}
}

// setState updates the connection state.
func (p *Publisher) setState(state State) {
	p.mu.Lock()
	p.state = state
	p.mu.Unlock()

	if p.onStateChange != nil {
		p.onStateChange(state)
	}
}
