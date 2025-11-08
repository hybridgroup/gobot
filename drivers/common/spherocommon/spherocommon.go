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

// Float32ToBytes splits the given 32 bit value in 4 bytes in big endian format without any further conversion of the IEEE 754 binary representation.
func Float32ToBytes(val float32) []byte {
	valBits := math.Float32bits(val)
	buf := make([]uint8, 4)
	binary.BigEndian.PutUint32(buf, valBits)
	return buf
}

// Uint16ToBytes splits the given 16 bit value in 2 bytes in big endian format.
func Uint16ToBytes(val uint16) []byte {
	buf := make([]uint8, 2)
	binary.BigEndian.PutUint16(buf, val)
	return buf
}
