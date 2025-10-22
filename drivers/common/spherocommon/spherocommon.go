package spherocommon

import (
	"encoding/binary"
	"math"
)

const (
	// ErrorEvent event when error encountered
	ErrorEvent = "error"

	// SensorDataEvent event when sensor data is received
	SensorDataEvent = "sensordata"

	// CollisionEvent event when collision is detected
	CollisionEvent = "collision"
)

// DefaultDataStreamingConfig returns a config with a sampling rate of 40hz, 1 sample frame per package,
// unlimited streaming, and will stream all available sensor information
func DefaultDataStreamingConfig() DataStreamingConfig {
	return DataStreamingConfig{
		N:     10,
		M:     1,
		Mask:  4294967295,
		Pcnt:  0,
		Mask2: 4294967295,
	}
}

// CalculateChecksum calculates the checksum for Sphero packets
func CalculateChecksum(buf []byte) byte {
	var calculatedChecksum uint16
	for i := range buf {
		calculatedChecksum += uint16(buf[i])
	}
	return uint8(^(calculatedChecksum % 256)) //nolint:gosec // TODO: fix later
}

func FloatToBytes(flt float32) []byte {
	u := math.Float32bits(flt)
	buf := make([]uint8, 4)
	binary.BigEndian.PutUint32(buf, u)
	return buf
}

func IntToBytes(int uint16) []byte {
	buf := make([]uint8, 2)
	binary.BigEndian.PutUint16(buf, int)
	return buf
}
