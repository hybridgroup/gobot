package spherocommon

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCalculateChecksum(t *testing.T) {
	tests := []struct {
		data     []byte
		checksum byte
	}{
		{[]byte{0x00}, 0xff},
		{[]byte{0xf0, 0x0f}, 0x00},
	}

	for _, tt := range tests {
		actual := CalculateChecksum(tt.data)
		if actual != tt.checksum {
			require.Fail(t, fmt.Sprintf("Expected %x, got %x for data %x.", tt.checksum, actual, tt.data))
		}
	}
}

func TestFloatToBytes(t *testing.T) {
	tests := []struct {
		data  float32
		bytes []byte
	}{
		{0, []byte{0x00, 0x00, 0x00, 0x00}},
		{160, []byte{0x43, 0x20, 0x00, 0x00}},
	}

	for _, tt := range tests {
		actual := Float32ToBytes(tt.data)
		if !bytes.Equal(actual, tt.bytes) {
			require.Fail(t, fmt.Sprintf("Expected %x, got %x for data %v.", tt.bytes, actual, tt.data))
		}
	}
}

func TestIntToBytes(t *testing.T) {
	tests := []struct {
		data  uint16
		bytes []byte
	}{
		{0, []byte{0x00, 0x00}},
		{1609, []byte{0x06, 0x49}},
	}

	for _, tt := range tests {
		actual := Int16ToBytes(tt.data)
		if !bytes.Equal(actual, tt.bytes) {
			require.Fail(t, fmt.Sprintf("Expected %x, got %x for data %v.", tt.bytes, actual, tt.data))
		}
	}
}
