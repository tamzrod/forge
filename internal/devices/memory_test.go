package devices

import (
	"testing"
)

func TestDeviceMemory_SetAndGet(t *testing.T) {
	m := NewDeviceMemory()

	m.Set("temperature", 25.5)
	m.Set("humidity", 60.0)

	if v, ok := m.Get("temperature"); !ok || v != 25.5 {
		t.Errorf("expected temperature 25.5, got %f, ok=%v", v, ok)
	}

	if v, ok := m.Get("humidity"); !ok || v != 60.0 {
		t.Errorf("expected humidity 60.0, got %f, ok=%v", v, ok)
	}
}

func TestDeviceMemory_GetOrDefault(t *testing.T) {
	m := NewDeviceMemory()
	m.Set("existing", 42.0)

	if v := m.GetOrDefault("existing", 0); v != 42.0 {
		t.Errorf("expected 42.0, got %f", v)
	}

	if v := m.GetOrDefault("missing", 99.0); v != 99.0 {
		t.Errorf("expected default 99.0, got %f", v)
	}
}

func TestDeviceMemory_Values(t *testing.T) {
	m := NewDeviceMemory()
	m.Set("a", 1.0)
	m.Set("b", 2.0)

	vals := m.Values()

	if len(vals) != 2 {
		t.Errorf("expected 2 values, got %d", len(vals))
	}

	if vals["a"] != 1.0 || vals["b"] != 2.0 {
		t.Errorf("unexpected values: %v", vals)
	}
}

func TestDeviceMemory_Contains(t *testing.T) {
	m := NewDeviceMemory()
	m.Set("exists", 1.0)

	if !m.Contains("exists") {
		t.Error("expected exists to be true")
	}

	if m.Contains("missing") {
		t.Error("expected missing to be false")
	}
}

func TestDeviceMemory_Reset(t *testing.T) {
	m := NewDeviceMemory()
	m.Set("temp", 25.0)
	m.Reset()

	if m.Contains("temp") {
		t.Error("expected temp to be cleared after reset")
	}
}

func TestMemoryRegion(t *testing.T) {
	r := MemoryRegion{
		Name: "registers",
		Base: 0,
		Size: 100,
	}

	if r.Name != "registers" {
		t.Errorf("expected name registers, got %s", r.Name)
	}
	if r.Base != 0 {
		t.Errorf("expected base 0, got %d", r.Base)
	}
	if r.Size != 100 {
		t.Errorf("expected size 100, got %d", r.Size)
	}
}

func TestMemoryMap(t *testing.T) {
	m := NewMemoryMap()
	m.AddRegion("input", 0, 10)
	m.AddRegion("holding", 100, 100)

	r, ok := m.Region("input")
	if !ok {
		t.Error("expected region input to exist")
	}
	if r.Base != 0 || r.Size != 10 {
		t.Errorf("unexpected region: %+v", r)
	}

	_, ok = m.Region("nonexistent")
	if ok {
		t.Error("expected nonexistent to not exist")
	}
}

func TestMemoryMap_ValidateAddress(t *testing.T) {
	m := NewMemoryMap()
	m.AddRegion("registers", 0, 100)

	if err := m.ValidateAddress("registers", 50); err != nil {
		t.Errorf("expected valid address, got: %v", err)
	}

	if err := m.ValidateAddress("registers", 100); err == nil {
		t.Error("expected error for out of bounds address")
	}

	if err := m.ValidateAddress("nonexistent", 0); err == nil {
		t.Error("expected error for unknown region")
	}
}
