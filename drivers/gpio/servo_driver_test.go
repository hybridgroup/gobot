//nolint:forcetypeassert // ok here
package gpio

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*ServoDriver)(nil)

func initTestServoDriver() *ServoDriver {
	return NewServoDriver(newGpioTestAdaptor(), "1")
}

func TestServoDriver(t *testing.T) {
	var err interface{}

	a := newGpioTestAdaptor()
	d := NewServoDriver(a, "1")

	assert.Equal(t, "1", d.Pin())
	assert.NotNil(t, d.Connection())

	a.servoWriteFunc = func(string, byte) error {
		return errors.New("pwm error")
	}

	err = d.Command("Min")(nil)
	require.ErrorContains(t, err.(error), "pwm error")

	err = d.Command("Center")(nil)
	require.ErrorContains(t, err.(error), "pwm error")

	err = d.Command("Max")(nil)
	require.ErrorContains(t, err.(error), "pwm error")

	err = d.Command("Move")(map[string]interface{}{"angle": 100.0})
	require.ErrorContains(t, err.(error), "pwm error")
}

func TestServoDriverStart(t *testing.T) {
	d := initTestServoDriver()
	require.NoError(t, d.Start())
}

func TestServoDriverHalt(t *testing.T) {
	d := initTestServoDriver()
	require.NoError(t, d.Halt())
}

func TestServoDriverMove(t *testing.T) {
	d := initTestServoDriver()
	_ = d.Move(100)
	assert.Equal(t, uint8(100), d.CurrentAngle)
	err := d.Move(200)
	assert.Equal(t, ErrServoOutOfRange, err)
}

func TestServoDriverMin(t *testing.T) {
	d := initTestServoDriver()
	_ = d.Min()
	assert.Equal(t, uint8(0), d.CurrentAngle)
}

func TestServoDriverMax(t *testing.T) {
	d := initTestServoDriver()
	_ = d.Max()
	assert.Equal(t, uint8(180), d.CurrentAngle)
}

func TestServoDriverCenter(t *testing.T) {
	d := initTestServoDriver()
	_ = d.Center()
	assert.Equal(t, uint8(90), d.CurrentAngle)
}

func TestServoDriverDefaultName(t *testing.T) {
	d := initTestServoDriver()
	assert.True(t, strings.HasPrefix(d.Name(), "Servo"))
}

func TestServoDriverSetName(t *testing.T) {
	d := initTestServoDriver()
	d.SetName("mybot")
	assert.Equal(t, "mybot", d.Name())
}
