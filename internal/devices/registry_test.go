package devices

import (
	"testing"
)

// mockDevice implements Device for testing
type mockDevice struct {
	*BaseDevice
}

func newMockDevice(id DeviceID, typ DeviceType, name string) *mockDevice {
	return &mockDevice{
		BaseDevice: NewBaseDevice(id, typ, name),
	}
}

func (m *mockDevice) Initialize() error {
	m.setState(StateInitialized)
	return nil
}

func (m *mockDevice) Tick() {
	if m.State() == StateInitialized {
		m.setState(StateRunning)
	}
}

func (m *mockDevice) Shutdown() error {
	m.setState(StateStopped)
	return nil
}

func TestRegistry_Register(t *testing.T) {
	r := NewRegistry()
	device := newMockDevice("device-001", "test", "Test Device")

	if err := r.Register(device); err != nil {
		t.Errorf("expected successful registration, got: %v", err)
	}

	if r.Count() != 1 {
		t.Errorf("expected count 1, got %d", r.Count())
	}
}

func TestRegistry_RegisterDuplicate(t *testing.T) {
	r := NewRegistry()
	device := newMockDevice("device-001", "test", "Test Device")

	r.Register(device)

	if err := r.Register(device); err == nil {
		t.Error("expected error for duplicate registration")
	}
}

func TestRegistry_Unregister(t *testing.T) {
	r := NewRegistry()
	device := newMockDevice("device-001", "test", "Test Device")
	r.Register(device)

	if err := r.Unregister("device-001"); err != nil {
		t.Errorf("expected successful unregister, got: %v", err)
	}

	if r.Count() != 0 {
		t.Errorf("expected count 0, got %d", r.Count())
	}
}

func TestRegistry_UnregisterMissing(t *testing.T) {
	r := NewRegistry()

	if err := r.Unregister("nonexistent"); err == nil {
		t.Error("expected error for unregistering missing device")
	}
}

func TestRegistry_Device(t *testing.T) {
	r := NewRegistry()
	device := newMockDevice("device-001", "test", "Test Device")
	r.Register(device)

	d, ok := r.Device("device-001")
	if !ok {
		t.Error("expected device to be found")
	}
	if d.ID() != "device-001" {
		t.Errorf("expected ID device-001, got %s", d.ID())
	}

	_, ok = r.Device("nonexistent")
	if ok {
		t.Error("expected device to not be found")
	}
}

func TestRegistry_Devices(t *testing.T) {
	r := NewRegistry()
	r.Register(newMockDevice("device-001", "test", "Device 1"))
	r.Register(newMockDevice("device-002", "test", "Device 2"))

	devices := r.Devices()
	if len(devices) != 2 {
		t.Errorf("expected 2 devices, got %d", len(devices))
	}
}

func TestRegistry_DevicesByType(t *testing.T) {
	r := NewRegistry()
	r.Register(newMockDevice("device-001", "weather_station", "Station 1"))
	r.Register(newMockDevice("device-002", "weather_station", "Station 2"))
	r.Register(newMockDevice("device-003", "revenue_meter", "Meter 1"))

	stations := r.DevicesByType("weather_station")
	if len(stations) != 2 {
		t.Errorf("expected 2 weather stations, got %d", len(stations))
	}

	meters := r.DevicesByType("revenue_meter")
	if len(meters) != 1 {
		t.Errorf("expected 1 revenue meter, got %d", len(meters))
	}
}

func TestRegistry_Initialize(t *testing.T) {
	r := NewRegistry()
	device := newMockDevice("device-001", "test", "Test Device")
	r.Register(device)

	if err := r.Initialize(); err != nil {
		t.Errorf("expected successful initialization, got: %v", err)
	}

	if device.State() != StateInitialized {
		t.Errorf("expected device state Initialized, got %s", device.State())
	}
}

func TestRegistry_Tick(t *testing.T) {
	r := NewRegistry()
	device := newMockDevice("device-001", "test", "Test Device")
	r.Register(device)

	r.Initialize()
	r.Tick()

	if device.State() != StateRunning {
		t.Errorf("expected device state Running, got %s", device.State())
	}
}

func TestRegistry_Shutdown(t *testing.T) {
	r := NewRegistry()
	device := newMockDevice("device-001", "test", "Test Device")
	r.Register(device)

	r.Initialize()
	r.Tick()

	if err := r.Shutdown(); err != nil {
		t.Errorf("expected successful shutdown, got: %v", err)
	}

	if device.State() != StateStopped {
		t.Errorf("expected device state Stopped, got %s", device.State())
	}
}

func TestRegistry_Clear(t *testing.T) {
	r := NewRegistry()
	r.Register(newMockDevice("device-001", "test", "Device 1"))
	r.Register(newMockDevice("device-002", "test", "Device 2"))

	r.Clear()

	if r.Count() != 0 {
		t.Errorf("expected count 0 after clear, got %d", r.Count())
	}
}
