package sphero

import (
	"bytes"
	"encoding/binary"
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

	ret = d.Command("ConfigureLocator")(
		map[string]interface{}{"Flags": 1.0, "X": 100.0, "Y": 100.0, "YawTare": 100.0},
	)
	gobot.Assert(t, ret, nil)

	ret = d.Command("SetHeading")(
		map[string]interface{}{"heading": 100.0},
	)
	gobot.Assert(t, ret, nil)

	ret = d.Command("SetRotationRate")(
		map[string]interface{}{"level": 100.0},
	)
	gobot.Assert(t, ret, nil)

	ret = d.Command("SetStabilization")(
		map[string]interface{}{"enable": true},
	)
	gobot.Assert(t, ret, nil)

	ret = d.Command("SetStabilization")(
		map[string]interface{}{"enable": false},
	)
	gobot.Assert(t, ret, nil)

	ret = d.Command("Stop")(nil)
	gobot.Assert(t, ret, nil)

	ret = d.Command("GetRGB")(nil)
	gobot.Assert(t, ret.([]byte), []byte{})

	ret = d.Command("ReadLocator")(nil)
	gobot.Assert(t, ret, []int16{})

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

func TestSpheroDriverSetDataStreaming(t *testing.T) {
	d := initTestSpheroDriver()
	d.SetDataStreaming(DefaultDataStreamingConfig())

	data := <-d.packetChannel

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, DefaultDataStreamingConfig())

	gobot.Assert(t, data.body, buf.Bytes())

	ret := d.Command("SetDataStreaming")(
		map[string]interface{}{
			"N":     100.0,
			"M":     200.0,
			"Mask":  300.0,
			"Pcnt":  255.0,
			"Mask2": 400.0,
		},
	)
	gobot.Assert(t, ret, nil)
	data = <-d.packetChannel

	dconfig := DataStreamingConfig{N: 100, M: 200, Mask: 300, Pcnt: 255, Mask2: 400}
	buf = new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, dconfig)

	gobot.Assert(t, data.body, buf.Bytes())
}

func TestConfigureLocator(t *testing.T) {
	d := initTestSpheroDriver()
	d.ConfigureLocator(DefaultLocatorConfig())
	data := <-d.packetChannel

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, DefaultLocatorConfig())

	gobot.Assert(t, data.body, buf.Bytes())

	ret := d.Command("ConfigureLocator")(
		map[string]interface{}{
			"Flags":   1.0,
			"X":       100.0,
			"Y":       100.0,
			"YawTare": 0.0,
		},
	)
	gobot.Assert(t, ret, nil)
	data = <-d.packetChannel

	lconfig := LocatorConfig{Flags: 1, X: 100, Y: 100, YawTare: 0}
	buf = new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, lconfig)

	gobot.Assert(t, data.body, buf.Bytes())
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
