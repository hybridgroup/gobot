package sphero

import (
	"testing"

	"github.com/hybridgroup/gobot"
)

func initTestSpheroDriver() *SpheroDriver {
	a := NewSpheroAdaptor("bot", "/dev/null")
	a.sp = nullReadWriteCloser{}
	return NewSpheroDriver(a, "bot")
}

func TestSpheroDriver(t *testing.T) {
	d := initTestSpheroDriver()
	var ret interface{}

	ret = d.Command("SetRGB")(
		map[string]interface{}{"r": 100.0, "g": 100.0, "b": 100.0},
	)
	gobot.Assert(t, ret, nil)

	ret = d.Command("Roll")(
		map[string]interface{}{"speed": 100.0, "heading": 100.0},
	)
	gobot.Assert(t, ret, nil)

	ret = d.Command("SetBackLED")(
		map[string]interface{}{"level": 100.0},
	)
	gobot.Assert(t, ret, nil)

	ret = d.Command("SetHeading")(
		map[string]interface{}{"heading": 100.0},
	)
	gobot.Assert(t, ret, nil)

	ret = d.Command("SetStabilization")(
		map[string]interface{}{"enable": true},
	)
	gobot.Assert(t, ret, nil)

	ret = d.Command("Stop")(nil)
	gobot.Assert(t, ret, nil)

	ret = d.Command("GetRGB")(nil)
	gobot.Assert(t, ret.([]byte), []byte{})

	gobot.Assert(t, d.Name(), "bot")
	gobot.Assert(t, d.Connection().Name(), "bot")
}

func TestSpheroDriverStart(t *testing.T) {
	d := initTestSpheroDriver()
	gobot.Assert(t, len(d.Start()), 0)
}

func TestSpheroDriverHalt(t *testing.T) {
	d := initTestSpheroDriver()
	d.adaptor().connected = true
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
