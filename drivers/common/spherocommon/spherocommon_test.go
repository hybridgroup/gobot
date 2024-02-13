package spherocommon

import (
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
			require.Fail(t, "Expected %x, got %x for data %x.", tt.checksum, actual, tt.data)
		}
	}
}
