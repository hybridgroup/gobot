package sphero

import (
	"bytes"
	"encoding/binary"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*SpheroDriver)(nil)

func initTestSpheroDriver() *SpheroDriver {
	a, _ := initTestSpheroAdaptor()
	_ = a.Connect()
	return NewSpheroDriver(a)
}

func TestSpheroDriverName(t *testing.T) {
	d := initTestSpheroDriver()
	assert.True(t, strings.HasPrefix(d.Name(), "Sphero"))
	d.SetName("NewName")
	assert.Equal(t, "NewName", d.Name())
}

func TestSpheroDriver(t *testing.T) {
	d := initTestSpheroDriver()
	var ret interface{}

	ret = d.Command("SetRGB")(
		map[string]interface{}{"r": 100.0, "g": 100.0, "b": 100.0},
	)
	assert.Nil(t, ret)

	ret = d.Command("Roll")(
		map[string]interface{}{"speed": 100.0, "heading": 100.0},
	)
	assert.Nil(t, ret)

	ret = d.Command("SetBackLED")(
		map[string]interface{}{"level": 100.0},
	)
	assert.Nil(t, ret)

	ret = d.Command("ConfigureLocator")(
		map[string]interface{}{"Flags": 1.0, "X": 100.0, "Y": 100.0, "YawTare": 100.0},
	)
	assert.Nil(t, ret)

	ret = d.Command("SetHeading")(
		map[string]interface{}{"heading": 100.0},
	)
	assert.Nil(t, ret)

	ret = d.Command("SetRotationRate")(
		map[string]interface{}{"level": 100.0},
	)
	assert.Nil(t, ret)

	ret = d.Command("SetStabilization")(
		map[string]interface{}{"enable": true},
	)
	assert.Nil(t, ret)

	ret = d.Command("SetStabilization")(
		map[string]interface{}{"enable": false},
	)
	assert.Nil(t, ret)

	ret = d.Command("Stop")(nil)
	assert.Nil(t, ret)

	ret = d.Command("GetRGB")(nil)
	assert.Equal(t, []byte{}, ret.([]byte))

	ret = d.Command("ReadLocator")(nil)
	assert.Equal(t, []int16{}, ret)

	assert.True(t, strings.HasPrefix(d.Name(), "Sphero"))
	assert.True(t, strings.HasPrefix(d.Connection().Name(), "Sphero"))
}

func TestSpheroDriverStart(t *testing.T) {
	d := initTestSpheroDriver()
	assert.NoError(t, d.Start())
}

func TestSpheroDriverHalt(t *testing.T) {
	d := initTestSpheroDriver()
	d.adaptor().connected = true
	assert.NoError(t, d.Halt())
}

func TestSpheroDriverSetDataStreaming(t *testing.T) {
	d := initTestSpheroDriver()
	d.SetDataStreaming(DefaultDataStreamingConfig())

	data := <-d.packetChannel

	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, DefaultDataStreamingConfig())

	assert.Equal(t, buf.Bytes(), data.body)

	ret := d.Command("SetDataStreaming")(
		map[string]interface{}{
			"N":     100.0,
			"M":     200.0,
			"Mask":  300.0,
			"Pcnt":  255.0,
			"Mask2": 400.0,
		},
	)
	assert.Nil(t, ret)
	data = <-d.packetChannel

	dconfig := DataStreamingConfig{N: 100, M: 200, Mask: 300, Pcnt: 255, Mask2: 400}
	buf = new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, dconfig)

	assert.Equal(t, buf.Bytes(), data.body)
}

func TestConfigureLocator(t *testing.T) {
	d := initTestSpheroDriver()
	d.ConfigureLocator(DefaultLocatorConfig())
	data := <-d.packetChannel

	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, DefaultLocatorConfig())

	assert.Equal(t, buf.Bytes(), data.body)

	ret := d.Command("ConfigureLocator")(
		map[string]interface{}{
			"Flags":   1.0,
			"X":       100.0,
			"Y":       100.0,
			"YawTare": 0.0,
		},
	)
	assert.Nil(t, ret)
	data = <-d.packetChannel

	lconfig := LocatorConfig{Flags: 1, X: 100, Y: 100, YawTare: 0}
	buf = new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, lconfig)

	assert.Equal(t, buf.Bytes(), data.body)
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
