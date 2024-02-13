//nolint:forcetypeassert // ok here
package i2c

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

// this ensures that the implementation is based on i2c.Driver, which implements the gobot.Driver
// and tests all implementations, so no further tests needed here for gobot.Driver interface
var _ gobot.Driver = (*BlinkMDriver)(nil)

func initTestBlinkMDriverWithStubbedAdaptor() (*BlinkMDriver, *i2cTestAdaptor) {
	a := newI2cTestAdaptor()
	d := NewBlinkMDriver(a)
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d, a
}

func TestNewBlinkMDriver(t *testing.T) {
	var di interface{} = NewBlinkMDriver(newI2cTestAdaptor())
	d, ok := di.(*BlinkMDriver)
	if !ok {
		require.Fail(t, "NewBlinkMDriver() should have returned a *BlinkMDriver")
	}
	assert.NotNil(t, d.Driver)
	assert.True(t, strings.HasPrefix(d.Name(), "BlinkM"))
	assert.Equal(t, 0x09, d.defaultAddress)
}

func TestBlinkMOptions(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewBlinkMDriver(newI2cTestAdaptor(), WithBus(2))
	assert.Equal(t, 2, d.GetBusOrDefault(1))
}

func TestBlinkMStart(t *testing.T) {
	d := NewBlinkMDriver(newI2cTestAdaptor())
	require.NoError(t, d.Start())
}

func TestBlinkMHalt(t *testing.T) {
	d, _ := initTestBlinkMDriverWithStubbedAdaptor()
	require.NoError(t, d.Halt())
}

// Commands
func TestNewBlinkMDriverCommands_Rgb(t *testing.T) {
	d, _ := initTestBlinkMDriverWithStubbedAdaptor()

	result := d.Command("Rgb")(rgb)
	assert.Nil(t, result)
}

func TestNewBlinkMDriverCommands_Fade(t *testing.T) {
	d, _ := initTestBlinkMDriverWithStubbedAdaptor()

	result := d.Command("Fade")(rgb)
	assert.Nil(t, result)
}

func TestNewBlinkMDriverCommands_FirmwareVersion(t *testing.T) {
	d, a := initTestBlinkMDriverWithStubbedAdaptor()
	param := make(map[string]interface{})
	// When len(data) is 2
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{99, 1})
		return 2, nil
	}

	result := d.Command("FirmwareVersion")(param)

	version, _ := d.FirmwareVersion()
	assert.Equal(t, version, result.(map[string]interface{})["version"].(string))

	// When len(data) is not 2
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{99})
		return 1, nil
	}
	result = d.Command("FirmwareVersion")(param)

	version, _ = d.FirmwareVersion()
	assert.Equal(t, version, result.(map[string]interface{})["version"].(string))
}

func TestNewBlinkMDriverCommands_Color(t *testing.T) {
	d, _ := initTestBlinkMDriverWithStubbedAdaptor()
	param := make(map[string]interface{})

	result := d.Command("Color")(param)

	color, _ := d.Color()
	assert.Equal(t, color, result.(map[string]interface{})["color"].([]byte))
}

func TestBlinkMFirmwareVersion(t *testing.T) {
	d, a := initTestBlinkMDriverWithStubbedAdaptor()
	// when len(data) is 2
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{99, 1})
		return 2, nil
	}

	version, _ := d.FirmwareVersion()
	assert.Equal(t, "99.1", version)

	// when len(data) is not 2
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{99})
		return 1, nil
	}

	version, _ = d.FirmwareVersion()
	assert.Equal(t, "", version)

	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}

	_, err := d.FirmwareVersion()
	require.ErrorContains(t, err, "write error")
}

func TestBlinkMColor(t *testing.T) {
	d, a := initTestBlinkMDriverWithStubbedAdaptor()
	// when len(data) is 3
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{99, 1, 2})
		return 3, nil
	}

	color, _ := d.Color()
	assert.Equal(t, []byte{99, 1, 2}, color)

	// when len(data) is not 3
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{99})
		return 1, nil
	}

	color, _ = d.Color()
	assert.Equal(t, []byte{}, color)

	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}

	_, err := d.Color()
	require.ErrorContains(t, err, "write error")
}

func TestBlinkMFade(t *testing.T) {
	d, a := initTestBlinkMDriverWithStubbedAdaptor()
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}

	err := d.Fade(100, 100, 100)
	require.ErrorContains(t, err, "write error")
}

func TestBlinkMRGB(t *testing.T) {
	d, a := initTestBlinkMDriverWithStubbedAdaptor()
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}

	err := d.Rgb(100, 100, 100)
	require.ErrorContains(t, err, "write error")
}
