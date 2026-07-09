package memory

import (
	"testing"
)

func TestNew(t *testing.T) {
	regions := map[string]uint32{
		"holding_registers": 100,
		"input_registers":  200,
		"coils":           20,
	}

	m := New(regions)

	if len(m.regions) != 3 {
		t.Errorf("expected 3 regions, got %d", len(m.regions))
	}

	// Check region sizes
	if m.regions["holding_registers"].Size != 100 {
		t.Errorf("expected holding_registers size 100, got %d", m.regions["holding_registers"].Size)
	}
	if m.regions["input_registers"].Size != 200 {
		t.Errorf("expected input_registers size 200, got %d", m.regions["input_registers"].Size)
	}
}

func TestReadWrite(t *testing.T) {
	m := New(map[string]uint32{
		"input_registers": 10,
	})

	// Write
	data := []byte{0x00, 0x01, 0x00, 0x02}
	err := m.Write("input_registers", 0, data)
	if err != nil {
		t.Fatalf("write failed: %v", err)
	}

	// Read
	result, err := m.Read("input_registers", 0, 4)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}

	if string(result) != string(data) {
		t.Errorf("expected %v, got %v", data, result)
	}
}

func TestReadWriteUint16(t *testing.T) {
	m := New(map[string]uint32{
		"input_registers": 10,
	})

	// Write
	err := m.WriteUint16("input_registers", 0, 0x1234)
	if err != nil {
		t.Fatalf("write failed: %v", err)
	}

	// Read
	val, err := m.ReadUint16("input_registers", 0)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}

	if val != 0x1234 {
		t.Errorf("expected 0x1234, got 0x%04x", val)
	}
}

func TestReadWriteFloat32(t *testing.T) {
	m := New(map[string]uint32{
		"input_registers": 10,
	})

	// Write
	err := m.WriteFloat32("input_registers", 0, 3.14159)
	if err != nil {
		t.Fatalf("write failed: %v", err)
	}

	// Read
	val, err := m.ReadFloat32("input_registers", 0)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}

	if val != 3.14159 {
		t.Errorf("expected 3.14159, got %f", val)
	}
}

func TestInvalidRegion(t *testing.T) {
	m := New(map[string]uint32{
		"input_registers": 10,
	})

	_, err := m.Read("nonexistent", 0, 4)
	if err == nil {
		t.Error("expected error for invalid region")
	}
}

func TestInvalidAddress(t *testing.T) {
	m := New(map[string]uint32{
		"input_registers": 10,
	})

	_, err := m.Read("input_registers", 8, 4) // 8 + 4 > 10
	if err == nil {
		t.Error("expected error for invalid address")
	}
}

func TestQuality(t *testing.T) {
	m := New(map[string]uint32{
		"input_registers": 10,
	})

	// Check initial quality
	q, err := m.Quality("input_registers", 0)
	if err != nil {
		t.Fatalf("quality failed: %v", err)
	}
	if q != QualityGood {
		t.Errorf("expected QualityGood, got %v", q)
	}

	// Set quality
	err = m.SetQuality("input_registers", 0, QualityOffline)
	if err != nil {
		t.Fatalf("set quality failed: %v", err)
	}

	// Check updated quality
	q, err = m.Quality("input_registers", 0)
	if err != nil {
		t.Fatalf("quality failed: %v", err)
	}
	if q != QualityOffline {
		t.Errorf("expected QualityOffline, got %v", q)
	}
}

func TestRegions(t *testing.T) {
	regions := map[string]uint32{
		"holding_registers": 100,
		"input_registers":  200,
	}

	m := New(regions)

	names := m.Regions()
	if len(names) != 2 {
		t.Errorf("expected 2 regions, got %d", len(names))
	}
}
