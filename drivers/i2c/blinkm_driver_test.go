package i2c

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
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
		t.Errorf("NewBlinkMDriver() should have returned a *BlinkMDriver")
	}
	gobottest.Refute(t, d.Driver, nil)
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "BlinkM"), true)
	gobottest.Assert(t, d.defaultAddress, 0x09)
}

func TestBlinkMOptions(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithBus() option and
	// least one of this driver. Further tests for options can also be done by call of "WithOption(val)(d)".
	d := NewBlinkMDriver(newI2cTestAdaptor(), WithBus(2))
	gobottest.Assert(t, d.GetBusOrDefault(1), 2)
}

func TestBlinkMStart(t *testing.T) {
	d := NewBlinkMDriver(newI2cTestAdaptor())
	gobottest.Assert(t, d.Start(), nil)
}

func TestBlinkMHalt(t *testing.T) {
	d, _ := initTestBlinkMDriverWithStubbedAdaptor()
	gobottest.Assert(t, d.Halt(), nil)
}

// Commands
func TestNewBlinkMDriverCommands_Rgb(t *testing.T) {
	d, _ := initTestBlinkMDriverWithStubbedAdaptor()

	result := d.Command("Rgb")(rgb)
	gobottest.Assert(t, result, nil)
}

func TestNewBlinkMDriverCommands_Fade(t *testing.T) {
	d, _ := initTestBlinkMDriverWithStubbedAdaptor()

	result := d.Command("Fade")(rgb)
	gobottest.Assert(t, result, nil)
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
	gobottest.Assert(t, result.(map[string]interface{})["version"].(string), version)

	// When len(data) is not 2
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{99})
		return 1, nil
	}
	result = d.Command("FirmwareVersion")(param)

	version, _ = d.FirmwareVersion()
	gobottest.Assert(t, result.(map[string]interface{})["version"].(string), version)
}

func TestNewBlinkMDriverCommands_Color(t *testing.T) {
	d, _ := initTestBlinkMDriverWithStubbedAdaptor()
	param := make(map[string]interface{})

	result := d.Command("Color")(param)

	color, _ := d.Color()
	gobottest.Assert(t, result.(map[string]interface{})["color"].([]byte), color)
}

func TestBlinkMFirmwareVersion(t *testing.T) {
	d, a := initTestBlinkMDriverWithStubbedAdaptor()
	// when len(data) is 2
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{99, 1})
		return 2, nil
	}

	version, _ := d.FirmwareVersion()
	gobottest.Assert(t, version, "99.1")

	// when len(data) is not 2
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{99})
		return 1, nil
	}

	version, _ = d.FirmwareVersion()
	gobottest.Assert(t, version, "")

	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}

	version, err := d.FirmwareVersion()
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestBlinkMColor(t *testing.T) {
	d, a := initTestBlinkMDriverWithStubbedAdaptor()
	// when len(data) is 3
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{99, 1, 2})
		return 3, nil
	}

	color, _ := d.Color()
	gobottest.Assert(t, color, []byte{99, 1, 2})

	// when len(data) is not 3
	a.i2cReadImpl = func(b []byte) (int, error) {
		copy(b, []byte{99})
		return 1, nil
	}

	color, _ = d.Color()
	gobottest.Assert(t, color, []byte{})

	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}

	color, err := d.Color()
	gobottest.Assert(t, err, errors.New("write error"))

}

func TestBlinkMFade(t *testing.T) {
	d, a := initTestBlinkMDriverWithStubbedAdaptor()
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}

	err := d.Fade(100, 100, 100)
	gobottest.Assert(t, err, errors.New("write error"))

}

func TestBlinkMRGB(t *testing.T) {
	d, a := initTestBlinkMDriverWithStubbedAdaptor()
	a.i2cWriteImpl = func([]byte) (int, error) {
		return 0, errors.New("write error")
	}

	err := d.Rgb(100, 100, 100)
	gobottest.Assert(t, err, errors.New("write error"))

}
