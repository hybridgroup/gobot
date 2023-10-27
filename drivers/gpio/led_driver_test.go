package gpio

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*LedDriver)(nil)

func initTestLedDriver() *LedDriver {
	a := newGpioTestAdaptor()
	a.digitalWriteFunc = func(string, byte) (err error) {
		return nil
	}
	a.pwmWriteFunc = func(string, byte) (err error) {
		return nil
	}
	return NewLedDriver(a, "1")
}

func TestLedDriver(t *testing.T) {
	var err interface{}
	a := newGpioTestAdaptor()
	d := NewLedDriver(a, "1")

	assert.Equal(t, "1", d.Pin())
	assert.NotNil(t, d.Connection())

	a.digitalWriteFunc = func(string, byte) (err error) {
		return errors.New("write error")
	}
	a.pwmWriteFunc = func(string, byte) (err error) {
		return errors.New("pwm error")
	}

	err = d.Command("Toggle")(nil)
	assert.ErrorContains(t, err.(error), "write error")

	err = d.Command("On")(nil)
	assert.ErrorContains(t, err.(error), "write error")

	err = d.Command("Off")(nil)
	assert.ErrorContains(t, err.(error), "write error")

	err = d.Command("Brightness")(map[string]interface{}{"level": 100.0})
	assert.ErrorContains(t, err.(error), "pwm error")
}

func TestLedDriverStart(t *testing.T) {
	d := initTestLedDriver()
	assert.NoError(t, d.Start())
}

func TestLedDriverHalt(t *testing.T) {
	d := initTestLedDriver()
	assert.NoError(t, d.Halt())
}

func TestLedDriverToggle(t *testing.T) {
	d := initTestLedDriver()
	_ = d.Off()
	_ = d.Toggle()
	assert.True(t, d.State())
	_ = d.Toggle()
	assert.False(t, d.State())
}

func TestLedDriverBrightness(t *testing.T) {
	a := newGpioTestAdaptor()
	d := NewLedDriver(a, "1")
	a.pwmWriteFunc = func(string, byte) (err error) {
		err = errors.New("pwm error")
		return
	}
	assert.ErrorContains(t, d.Brightness(150), "pwm error")
}

func TestLEDDriverDefaultName(t *testing.T) {
	a := newGpioTestAdaptor()
	d := NewLedDriver(a, "1")
	assert.True(t, strings.HasPrefix(d.Name(), "LED"))
}

func TestLEDDriverSetName(t *testing.T) {
	a := newGpioTestAdaptor()
	d := NewLedDriver(a, "1")
	d.SetName("mybot")
	assert.Equal(t, "mybot", d.Name())
}
