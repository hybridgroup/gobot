//nolint:forcetypeassert // ok here
package gpio

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/aio"
)

var _ gobot.Driver = (*DirectPinDriver)(nil)

func initTestDirectPinDriver() *DirectPinDriver {
	a := newGpioTestAdaptor()
	a.digitalReadFunc = func(string) (int, error) {
		return 1, nil
	}
	a.digitalWriteFunc = func(string, byte) error {
		return errors.New("write error")
	}
	a.pwmWriteFunc = func(string, byte) error {
		return errors.New("write error")
	}
	a.servoWriteFunc = func(string, byte) error {
		return errors.New("write error")
	}
	return NewDirectPinDriver(a, "1")
}

func TestNewDirectPinDriver(t *testing.T) {
	// arrange
	a := newGpioTestAdaptor()
	// act
	d := NewDirectPinDriver(a, "10")
	// assert
	assert.IsType(t, &DirectPinDriver{}, d)
	// assert: gpio.driver attributes
	require.NotNil(t, d.driver)
	assert.True(t, strings.HasPrefix(d.driverCfg.name, "DirectPin"))
	assert.Equal(t, "10", d.driverCfg.pin)
	assert.Equal(t, a, d.connection)
	assert.NotNil(t, d.afterStart)
	assert.NotNil(t, d.beforeHalt)
	assert.NotNil(t, d.Commander)
	assert.NotNil(t, d.mutex)
}

func TestNewDirectPinDriver_options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option, least one
	// option of this driver and one of another driver (which should lead to panic). Further tests for options can also
	// be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const (
		myName = "count up"
	)
	panicFunc := func() {
		NewDirectPinDriver(newGpioTestAdaptor(), "1", WithName("crazy"),
			aio.WithActuatorScaler(func(float64) int { return 0 }))
	}
	// act
	d := NewDirectPinDriver(newGpioTestAdaptor(), "1", WithName(myName))
	// assert
	assert.Equal(t, myName, d.Name())
	assert.PanicsWithValue(t, "'scaler option for analog actuators' can not be applied on 'crazy'", panicFunc)
}

func TestDirectPin_Commands(t *testing.T) {
	var ret map[string]interface{}
	var err interface{}

	d := initTestDirectPinDriver()
	ret = d.Command("DigitalRead")(nil).(map[string]interface{})

	assert.Equal(t, 1, ret["val"].(int))
	assert.Nil(t, ret["err"])

	err = d.Command("DigitalWrite")(map[string]interface{}{"level": "1"})
	require.EqualError(t, err.(error), "write error")

	err = d.Command("PwmWrite")(map[string]interface{}{"level": "1"})
	require.EqualError(t, err.(error), "write error")

	err = d.Command("ServoWrite")(map[string]interface{}{"level": "1"})
	require.EqualError(t, err.(error), "write error")
}

func TestDirectPinOff(t *testing.T) {
	d := initTestDirectPinDriver()
	require.Error(t, d.Off())

	a := newGpioTestAdaptor()
	d = NewDirectPinDriver(a, "1")
	require.NoError(t, d.Off())
}

func TestDirectPinOffNotSupported(t *testing.T) {
	a := &gpioTestBareAdaptor{}
	d := NewDirectPinDriver(a, "1")
	require.EqualError(t, d.Off(), "DigitalWrite is not supported by this platform")
}

func TestDirectPinOn(t *testing.T) {
	a := newGpioTestAdaptor()
	d := NewDirectPinDriver(a, "1")
	require.NoError(t, d.On())
}

func TestDirectPinOnError(t *testing.T) {
	d := initTestDirectPinDriver()
	require.Error(t, d.On())
}

func TestDirectPinOnNotSupported(t *testing.T) {
	a := &gpioTestBareAdaptor{}
	d := NewDirectPinDriver(a, "1")
	require.EqualError(t, d.On(), "DigitalWrite is not supported by this platform")
}

func TestDirectPinDigitalWrite(t *testing.T) {
	adaptor := newGpioTestAdaptor()
	d := NewDirectPinDriver(adaptor, "1")
	require.NoError(t, d.DigitalWrite(1))
}

func TestDirectPinDigitalWriteNotSupported(t *testing.T) {
	a := &gpioTestBareAdaptor{}
	d := NewDirectPinDriver(a, "1")
	require.EqualError(t, d.DigitalWrite(1), "DigitalWrite is not supported by this platform")
}

func TestDirectPinDigitalWriteError(t *testing.T) {
	d := initTestDirectPinDriver()
	require.Error(t, d.DigitalWrite(1))
}

func TestDirectPinDigitalRead(t *testing.T) {
	d := initTestDirectPinDriver()
	ret, err := d.DigitalRead()
	assert.Equal(t, 1, ret)
	require.NoError(t, err)
}

func TestDirectPinDigitalReadNotSupported(t *testing.T) {
	a := &gpioTestBareAdaptor{}
	d := NewDirectPinDriver(a, "1")
	_, e := d.DigitalRead()
	require.EqualError(t, e, "DigitalRead is not supported by this platform")
}

func TestDirectPinPwmWrite(t *testing.T) {
	a := newGpioTestAdaptor()
	d := NewDirectPinDriver(a, "1")
	require.NoError(t, d.PwmWrite(1))
}

func TestDirectPinPwmWriteNotSupported(t *testing.T) {
	a := &gpioTestBareAdaptor{}
	d := NewDirectPinDriver(a, "1")
	require.EqualError(t, d.PwmWrite(1), "PwmWrite is not supported by this platform")
}

func TestDirectPinPwmWriteError(t *testing.T) {
	d := initTestDirectPinDriver()
	require.Error(t, d.PwmWrite(1))
}

func TestDirectPinServoWrite(t *testing.T) {
	a := newGpioTestAdaptor()
	d := NewDirectPinDriver(a, "1")
	require.NoError(t, d.ServoWrite(1))
}

func TestDirectPinServoWriteNotSupported(t *testing.T) {
	a := &gpioTestBareAdaptor{}
	d := NewDirectPinDriver(a, "1")
	require.EqualError(t, d.ServoWrite(1), "ServoWrite is not supported by this platform")
}

func TestDirectPinServoWriteError(t *testing.T) {
	d := initTestDirectPinDriver()
	require.Error(t, d.ServoWrite(1))
}
