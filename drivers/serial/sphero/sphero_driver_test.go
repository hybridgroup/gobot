//nolint:forcetypeassert // ok here
package sphero

import (
	"bytes"
	"encoding/binary"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/common/spherocommon"
	"gobot.io/x/gobot/v2/drivers/serial"
	"gobot.io/x/gobot/v2/drivers/serial/testutil"
)

var _ gobot.Driver = (*SpheroDriver)(nil)

func initTestSpheroDriver() *SpheroDriver {
	a := testutil.NewSerialTestAdaptor()
	d := NewSpheroDriver(a)
	d.shutdownWaitTime = 0 // to speed up the tests
	return d
}

func TestNewSpheroDriver(t *testing.T) {
	d := initTestSpheroDriver()
	assert.True(t, strings.HasPrefix(d.Name(), "Sphero"))
	assert.NotNil(t, d.Eventer)
}

func TestNewSpheroDriverWithName(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option.	Further
	// tests for options can also be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const newName = "new name"
	a := testutil.NewSerialTestAdaptor()
	// act
	d := NewSpheroDriver(a, serial.WithName(newName))
	// assert
	assert.Equal(t, newName, d.Name())
}

func TestSpheroCommands(t *testing.T) {
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
}

func TestSpheroStart(t *testing.T) {
	d := initTestSpheroDriver()
	require.NoError(t, d.Start())
}

func TestSpheroHalt(t *testing.T) {
	a := testutil.NewSerialTestAdaptor()
	_ = a.Connect()
	d := NewSpheroDriver(a)
	d.shutdownWaitTime = 0 // to speed up the tests
	require.NoError(t, d.Halt())
}

func TestSpheroSetDataStreaming(t *testing.T) {
	d := initTestSpheroDriver()
	d.SetDataStreaming(spherocommon.DefaultDataStreamingConfig())

	data := <-d.packetChannel

	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, spherocommon.DefaultDataStreamingConfig())

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

	dconfig := spherocommon.DataStreamingConfig{N: 100, M: 200, Mask: 300, Pcnt: 255, Mask2: 400}
	buf = new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, dconfig)

	assert.Equal(t, buf.Bytes(), data.body)
}

func TestSpheroConfigureLocator(t *testing.T) {
	d := initTestSpheroDriver()
	d.ConfigureLocator(spheroDefaultLocatorConfig())
	data := <-d.packetChannel

	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, spheroDefaultLocatorConfig())

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

	lconfig := spherocommon.LocatorConfig{Flags: 1, X: 100, Y: 100, YawTare: 0}
	buf = new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, lconfig)

	assert.Equal(t, buf.Bytes(), data.body)
}
