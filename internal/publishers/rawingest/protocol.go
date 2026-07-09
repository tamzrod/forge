package rawingest

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Protocol constants
const (
	ProtocolVersion = 1

	// Packet types
	TypeData = 0x01
	TypeHeartbeat = 0x02
	TypeConfig = 0x03

	// Packet sizes (fixed header + variable data)
	HeaderSize = 16 // unit_id(1) + seq(2) + timestamp(8) + type(1) + len(2) + crc(2)
)

// Packet represents a Raw Ingest protocol packet.
type Packet struct {
	UnitID    uint8
	Sequence  uint16
	Timestamp uint64
	Type      uint8
	Data      []byte
	CRC       uint16
}

// EncodePacket encodes a packet into bytes.
func EncodePacket(unitID uint8, seq uint16, timestamp uint64, packetType uint8, data []byte) ([]byte, error) {
	// Calculate total size
	totalSize := HeaderSize + len(data)

	// Build packet
	buf := new(bytes.Buffer)

	// Unit ID
	if err := binary.Write(buf, binary.LittleEndian, unitID); err != nil {
		return nil, fmt.Errorf("failed to write unit ID: %w", err)
	}

	// Sequence number
	if err := binary.Write(buf, binary.LittleEndian, seq); err != nil {
		return nil, fmt.Errorf("failed to write sequence: %w", err)
	}

	// Timestamp (Unix nanoseconds)
	if err := binary.Write(buf, binary.LittleEndian, timestamp); err != nil {
		return nil, fmt.Errorf("failed to write timestamp: %w", err)
	}

	// Packet type
	if err := binary.Write(buf, binary.LittleEndian, packetType); err != nil {
		return nil, fmt.Errorf("failed to write type: %w", err)
	}

	// Data length
	dataLen := uint16(len(data))
	if err := binary.Write(buf, binary.LittleEndian, dataLen); err != nil {
		return nil, fmt.Errorf("failed to write length: %w", err)
	}

	// Data payload
	if _, err := buf.Write(data); err != nil {
		return nil, fmt.Errorf("failed to write data: %w", err)
	}

	// Calculate CRC over header + data
	packet := buf.Bytes()
	crc := crc16(packet[:HeaderSize-2+len(data)])

	// CRC
	if err := binary.Write(buf, binary.LittleEndian, crc); err != nil {
		return nil, fmt.Errorf("failed to write CRC: %w", err)
	}

	return buf.Bytes(), nil
}

// DecodePacket decodes a packet from bytes.
func DecodePacket(data []byte) (*Packet, error) {
	if len(data) < HeaderSize {
		return nil, fmt.Errorf("packet too short: %d bytes", len(data))
	}

	buf := bytes.NewReader(data)

	var unitID uint8
	var seq uint16
	var timestamp uint64
	var packetType uint8
	var dataLen uint16

	// Read header fields
	if err := binary.Read(buf, binary.LittleEndian, &unitID); err != nil {
		return nil, fmt.Errorf("failed to read unit ID: %w", err)
	}
	if err := binary.Read(buf, binary.LittleEndian, &seq); err != nil {
		return nil, fmt.Errorf("failed to read sequence: %w", err)
	}
	if err := binary.Read(buf, binary.LittleEndian, &timestamp); err != nil {
		return nil, fmt.Errorf("failed to read timestamp: %w", err)
	}
	if err := binary.Read(buf, binary.LittleEndian, &packetType); err != nil {
		return nil, fmt.Errorf("failed to read type: %w", err)
	}
	if err := binary.Read(buf, binary.LittleEndian, &dataLen); err != nil {
		return nil, fmt.Errorf("failed to read length: %w", err)
	}

	// Validate data length
	expectedSize := HeaderSize + int(dataLen)
	if len(data) < expectedSize {
		return nil, fmt.Errorf("data length mismatch: expected %d, got %d", expectedSize, len(data))
	}

	// Read data
	payload := make([]byte, dataLen)
	if _, err := buf.Read(payload); err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}

	// Read and verify CRC
	var crc uint16
	if err := binary.Read(buf, binary.LittleEndian, &crc); err != nil {
		return nil, fmt.Errorf("failed to read CRC: %w", err)
	}

	// Verify CRC
	calculated := crc16(data[:HeaderSize-2+len(payload)])
	if calculated != crc {
		return nil, fmt.Errorf("CRC mismatch: expected 0x%04x, got 0x%04x", calculated, crc)
	}

	return &Packet{
		UnitID:    unitID,
		Sequence:  seq,
		Timestamp: timestamp,
		Type:      packetType,
		Data:      payload,
		CRC:       crc,
	}, nil
}

// EncodeMemoryValues encodes memory values into the data payload.
// Format: [register_addr(2), value(4)] repeated
func EncodeMemoryValues(values map[string]float64) []byte {
	// Register mapping for Weather Station
	registerMap := map[string]uint16{
		"temperature":     0,
		"humidity":       1,
		"pressure":       2,
		"cloud_cover":    3,
		"wind_speed":     4,
		"wind_direction": 5,
		"rain_status":    6,
	}

	buf := new(bytes.Buffer)

	for name, value := range registerMap {
		if v, ok := values[name]; ok {
			// Write register address
			binary.Write(buf, binary.LittleEndian, value)

			// Write value as fixed-point (multiply by 1000 for 3 decimal places)
			scaled := int32(v * 1000)
			binary.Write(buf, binary.LittleEndian, scaled)
		}
	}

	return buf.Bytes()
}

// crc16 calculates CRC-16/MODBUS checksum.
func crc16(data []byte) uint16 {
	var crc uint16 = 0xFFFF

	for _, b := range data {
		crc ^= uint16(b)
		for i := 0; i < 8; i++ {
			if crc&0x0001 != 0 {
				crc = (crc >> 1) ^ 0xA001
			} else {
				crc >>= 1
			}
		}
	}

	return crc
}
