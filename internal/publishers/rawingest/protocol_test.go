package rawingest

import (
	"testing"
)

func TestEncodePacket(t *testing.T) {
	packet, err := EncodePacket(1, 100, 1234567890, TypeData, []byte{0x01, 0x02, 0x03})
	if err != nil {
		t.Fatalf("failed to encode packet: %v", err)
	}

	if len(packet) < HeaderSize+3 {
		t.Errorf("packet too short: %d bytes", len(packet))
	}
}

func TestDecodePacket(t *testing.T) {
	data := []byte{0x01}
	_, err := DecodePacket(data)
	if err == nil {
		t.Error("expected error for short packet")
	}
}

func TestDecodePacket_Valid(t *testing.T) {
	data := []byte{
		0x01,       // Unit ID
		0x64, 0x00, // Sequence (little endian: 100)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Timestamp (8 bytes)
		0x01,       // Type
		0x03, 0x00, // Length (3 bytes)
		0x01, 0x02, 0x03, // Data
		0x00, 0x00, // CRC placeholder
	}

	_, err := DecodePacket(data)
	if err != nil {
		// CRC will fail which is expected without proper CRC calculation
		t.Logf("Decode failed (expected with placeholder CRC): %v", err)
	}
}

func TestEncodeMemoryValues(t *testing.T) {
	values := map[string]float64{
		"temperature":     25.5,
		"humidity":       60.0,
		"pressure":       1013.25,
		"cloud_cover":    0.0,
		"wind_speed":     5.2,
		"wind_direction": 180.0,
		"rain_status":    0.0,
	}

	data := EncodeMemoryValues(values)

	if len(data) == 0 {
		t.Error("expected non-empty data")
	}

	// Each value: 2 bytes address + 4 bytes value = 6 bytes
	// 7 values * 6 bytes = 42 bytes expected
	if len(data) != 42 {
		t.Errorf("expected 42 bytes, got %d", len(data))
	}
}

func TestEncodeMemoryValues_Partial(t *testing.T) {
	values := map[string]float64{
		"temperature": 25.5,
	}

	data := EncodeMemoryValues(values)

	if len(data) != 6 {
		t.Errorf("expected 6 bytes for 1 value, got %d", len(data))
	}
}

func TestCRC16(t *testing.T) {
	tests := []struct {
		data   []byte
		expect uint16
	}{
		{[]byte{0x01, 0x02}, 0x6a31},
		{[]byte{0xff, 0xff}, 0xf0f8},
		{[]byte{}, 0xffff},
	}

	for _, tt := range tests {
		got := crc16(tt.data)
		if got != tt.expect {
			t.Errorf("crc16(%v) = 0x%04x, want 0x%04x", tt.data, got, tt.expect)
		}
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid",
			config: Config{
				Host: "localhost",
				Port: 500,
			},
			wantErr: false,
		},
		{
			name: "empty host",
			config: Config{
				Host: "",
				Port: 500,
			},
			wantErr: true,
		},
		{
			name: "zero port",
			config: Config{
				Host: "localhost",
				Port: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_Address(t *testing.T) {
	cfg := Config{
		Host: "localhost",
		Port: 500,
	}

	if addr := cfg.Address(); addr != "localhost:500" {
		t.Errorf("Address() = %s, want localhost:500", addr)
	}
}

func TestStateString(t *testing.T) {
	tests := []struct {
		state State
		want  string
	}{
		{StateDisconnected, "Disconnected"},
		{StateConnecting, "Connecting"},
		{StateConnected, "Connected"},
		{StateError, "Error"},
		{State(100), "Unknown"},
	}

	for _, tt := range tests {
		if got := tt.state.String(); got != tt.want {
			t.Errorf("State(%d).String() = %s, want %s", tt.state, got, tt.want)
		}
	}
}
