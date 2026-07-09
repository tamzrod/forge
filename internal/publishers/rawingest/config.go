package rawingest

import (
	"fmt"
	"time"
)

// Config holds the publisher configuration.
type Config struct {
	// MMA2 connection
	Host string
	Port uint16

	// Device identification
	UnitID uint8

	// Publish behavior
	Enabled        bool
	Interval       time.Duration
	ReconnectDelay time.Duration
	MaxRetries     int

	// Protocol settings
	Timeout time.Duration
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		Host:           "localhost",
		Port:           500,
		UnitID:         1,
		Enabled:        true,
		Interval:       1 * time.Second,
		ReconnectDelay: 5 * time.Second,
		MaxRetries:     3,
		Timeout:        5 * time.Second,
	}
}

// Validate checks the configuration.
func (c Config) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}
	if c.Port == 0 {
		return fmt.Errorf("port must be non-zero")
	}
	if c.Interval < 0 {
		return fmt.Errorf("interval cannot be negative")
	}
	return nil
}

// Address returns the connection address.
func (c Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
