package devices

import (
	"testing"
)

func TestBaseDevice(t *testing.T) {
	id := DeviceID("test-device-001")
	typ := DeviceType("test_device")
	name := "Test Device"

	device := NewBaseDevice(id, typ, name)

	if device.ID() != id {
		t.Errorf("expected ID %s, got %s", id, device.ID())
	}

	if device.Type() != typ {
		t.Errorf("expected type %s, got %s", typ, device.Type())
	}

	if device.Name() != name {
		t.Errorf("expected name %s, got %s", name, device.Name())
	}

	if device.State() != StateCreated {
		t.Errorf("expected state Created, got %s", device.State())
	}
}

func TestValidateID(t *testing.T) {
	tests := []struct {
		id    DeviceID
		valid bool
	}{
		{"device-001", true},
		{"", false},
	}

	for _, tt := range tests {
		err := ValidateID(tt.id)
		if tt.valid && err != nil {
			t.Errorf("expected valid ID, got error: %v", err)
		}
		if !tt.valid && err == nil {
			t.Errorf("expected invalid ID for %s", tt.id)
		}
	}
}

func TestValidateType(t *testing.T) {
	tests := []struct {
		typ   DeviceType
		valid bool
	}{
		{"weather_station", true},
		{"", false},
	}

	for _, tt := range tests {
		err := ValidateType(tt.typ)
		if tt.valid && err != nil {
			t.Errorf("expected valid type, got error: %v", err)
		}
		if !tt.valid && err == nil {
			t.Errorf("expected invalid type for %s", tt.typ)
		}
	}
}

func TestStateString(t *testing.T) {
	tests := []struct {
		state   State
		wanterr string
	}{
		{StateCreated, "Created"},
		{StateInitialized, "Initialized"},
		{StateRunning, "Running"},
		{StatePaused, "Paused"},
		{StateStopped, "Stopped"},
		{StateFaulted, "Faulted"},
		{State(100), "Unknown"},
	}

	for _, tt := range tests {
		if got := tt.state.String(); got != tt.wanterr {
			t.Errorf("State(%d).String() = %s, want %s", tt.state, got, tt.wanterr)
		}
	}
}
