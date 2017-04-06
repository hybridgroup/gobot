package sphero

import (
	"bytes"
	"encoding/binary"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*SpheroDriver)(nil)

func initTestSpheroDriver() *SpheroDriver {
	a, _ := initTestSpheroAdaptor()
	a.Connect()
	return NewSpheroDriver(a)
}

func TestSpheroDriverName(t *testing.T) {
	d := initTestSpheroDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Sphero"), true)
	d.SetName("NewName")
	gobottest.Assert(t, d.Name(), "NewName")
}

func TestSpheroDriver(t *testing.T) {
	d := initTestSpheroDriver()
	var ret interface{}

	ret = d.Command("SetRGB")(
		map[string]interface{}{"r": 100.0, "g": 100.0, "b": 100.0},
	)
	gobottest.Assert(t, ret, nil)

	ret = d.Command("Roll")(
		map[string]interface{}{"speed": 100.0, "heading": 100.0},
	)
	gobottest.Assert(t, ret, nil)

	ret = d.Command("SetBackLED")(
		map[string]interface{}{"level": 100.0},
	)
	gobottest.Assert(t, ret, nil)

	ret = d.Command("ConfigureLocator")(
		map[string]interface{}{"Flags": 1.0, "X": 100.0, "Y": 100.0, "YawTare": 100.0},
	)
	gobottest.Assert(t, ret, nil)

	ret = d.Command("SetHeading")(
		map[string]interface{}{"heading": 100.0},
	)
	gobottest.Assert(t, ret, nil)

	ret = d.Command("SetRotationRate")(
		map[string]interface{}{"level": 100.0},
	)
	gobottest.Assert(t, ret, nil)

	ret = d.Command("SetStabilization")(
		map[string]interface{}{"enable": true},
	)
	gobottest.Assert(t, ret, nil)

	ret = d.Command("SetStabilization")(
		map[string]interface{}{"enable": false},
	)
	gobottest.Assert(t, ret, nil)

	ret = d.Command("Stop")(nil)
	gobottest.Assert(t, ret, nil)

	ret = d.Command("GetRGB")(nil)
	gobottest.Assert(t, ret.([]byte), []byte{})

	ret = d.Command("ReadLocator")(nil)
	gobottest.Assert(t, ret, []int16{})

	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Sphero"), true)
	gobottest.Assert(t, strings.HasPrefix(d.Connection().Name(), "Sphero"), true)
}

func TestSpheroDriverStart(t *testing.T) {
	d := initTestSpheroDriver()
	gobottest.Assert(t, d.Start(), nil)
}

func TestSpheroDriverHalt(t *testing.T) {
	d := initTestSpheroDriver()
	d.adaptor().connected = true
	gobottest.Assert(t, d.Halt(), nil)
}

func TestSpheroDriverSetDataStreaming(t *testing.T) {
	d := initTestSpheroDriver()
	d.SetDataStreaming(DefaultDataStreamingConfig())

	data := <-d.packetChannel

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, DefaultDataStreamingConfig())

	gobottest.Assert(t, data.body, buf.Bytes())

	ret := d.Command("SetDataStreaming")(
		map[string]interface{}{
			"N":     100.0,
			"M":     200.0,
			"Mask":  300.0,
			"Pcnt":  255.0,
			"Mask2": 400.0,
		},
	)
	gobottest.Assert(t, ret, nil)
	data = <-d.packetChannel

	dconfig := DataStreamingConfig{N: 100, M: 200, Mask: 300, Pcnt: 255, Mask2: 400}
	buf = new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, dconfig)

	gobottest.Assert(t, data.body, buf.Bytes())
}

func TestConfigureLocator(t *testing.T) {
	d := initTestSpheroDriver()
	d.ConfigureLocator(DefaultLocatorConfig())
	data := <-d.packetChannel

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, DefaultLocatorConfig())

	gobottest.Assert(t, data.body, buf.Bytes())

	ret := d.Command("ConfigureLocator")(
		map[string]interface{}{
			"Flags":   1.0,
			"X":       100.0,
			"Y":       100.0,
			"YawTare": 0.0,
		},
	)
	gobottest.Assert(t, ret, nil)
	data = <-d.packetChannel

	lconfig := LocatorConfig{Flags: 1, X: 100, Y: 100, YawTare: 0}
	buf = new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, lconfig)

	gobottest.Assert(t, data.body, buf.Bytes())
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
