package gpio

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*RgbLedDriver)(nil)

func initTestRgbLedDriver() *RgbLedDriver {
	a := newGpioTestAdaptor()
	a.digitalWriteFunc = func(string, byte) (err error) {
		return nil
	}
	a.pwmWriteFunc = func(string, byte) (err error) {
		return nil
	}
	return NewRgbLedDriver(a, "1", "2", "3")
}

func TestRgbLedDriver(t *testing.T) {
	var err interface{}

	a := newGpioTestAdaptor()
	d := NewRgbLedDriver(a, "1", "2", "3")

	assert.Equal(t, "r=1, g=2, b=3", d.Pin())
	assert.Equal(t, "1", d.RedPin())
	assert.Equal(t, "2", d.GreenPin())
	assert.Equal(t, "3", d.BluePin())
	assert.NotNil(t, d.Connection())

	a.digitalWriteFunc = func(string, byte) (err error) {
		return errors.New("write error")
	}
	a.pwmWriteFunc = func(string, byte) (err error) {
		return errors.New("pwm error")
	}

	err = d.Command("Toggle")(nil)
	assert.ErrorContains(t, err.(error), "pwm error")

	err = d.Command("On")(nil)
	assert.ErrorContains(t, err.(error), "pwm error")

	err = d.Command("Off")(nil)
	assert.ErrorContains(t, err.(error), "pwm error")

	err = d.Command("SetRGB")(map[string]interface{}{"r": 0xff, "g": 0xff, "b": 0xff})
	assert.ErrorContains(t, err.(error), "pwm error")
}

func TestRgbLedDriverStart(t *testing.T) {
	d := initTestRgbLedDriver()
	assert.NoError(t, d.Start())
}

func TestRgbLedDriverHalt(t *testing.T) {
	d := initTestRgbLedDriver()
	assert.NoError(t, d.Halt())
}

func TestRgbLedDriverToggle(t *testing.T) {
	d := initTestRgbLedDriver()
	_ = d.Off()
	_ = d.Toggle()
	assert.True(t, d.State())
	_ = d.Toggle()
	assert.False(t, d.State())
}

func TestRgbLedDriverSetLevel(t *testing.T) {
	a := newGpioTestAdaptor()
	d := NewRgbLedDriver(a, "1", "2", "3")
	assert.NoError(t, d.SetLevel("1", 150))

	d = NewRgbLedDriver(a, "1", "2", "3")
	a.pwmWriteFunc = func(string, byte) (err error) {
		err = errors.New("pwm error")
		return
	}
	assert.ErrorContains(t, d.SetLevel("1", 150), "pwm error")
}

func TestRgbLedDriverDefaultName(t *testing.T) {
	a := newGpioTestAdaptor()
	d := NewRgbLedDriver(a, "1", "2", "3")
	assert.True(t, strings.HasPrefix(d.Name(), "RGB"))
}

func TestRgbLedDriverSetName(t *testing.T) {
	a := newGpioTestAdaptor()
	d := NewRgbLedDriver(a, "1", "2", "3")
	d.SetName("mybot")
	assert.Equal(t, "mybot", d.Name())
}
