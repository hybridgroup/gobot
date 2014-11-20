package sphero

import (
	"testing"

	"github.com/hybridgroup/gobot"
)

func initTestSpheroDriver() *SpheroDriver {
	a := NewSpheroAdaptor("bot", "/dev/null")
	a.sp = gobot.NullReadWriteCloser{}
	return NewSpheroDriver(a, "bot")
}

func TestSpheroDriverStart(t *testing.T) {
	d := initTestSpheroDriver()
	gobot.Assert(t, len(d.Start()), 0)
}

func TestSpheroDriverHalt(t *testing.T) {
	d := initTestSpheroDriver()
	gobot.Assert(t, len(d.Halt()), 0)
}

func TestCalculateChecksum(t *testing.T) {
	tests := []struct {
		data     []byte
		checksum byte
	}{
		{[]byte{0x00}, 0xff},
		{[]byte{0xf0, 0x0f}, 0x00},
	}

	for _, tt := range tests {
		actual := calculateChecksum(tt.data)
		if actual != tt.checksum {
			t.Errorf("Expected %x, got %x for data %x.", tt.checksum, actual, tt.data)
		}
	}
}
