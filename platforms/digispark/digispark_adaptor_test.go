//go:build libusb
// +build libusb

//nolint:forcetypeassert // ok here
package digispark

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
)

var _ gobot.Adaptor = (*Adaptor)(nil)

var (
	_ gpio.DigitalWriter = (*Adaptor)(nil)
	_ gpio.PwmWriter     = (*Adaptor)(nil)
	_ gpio.ServoWriter   = (*Adaptor)(nil)
)

type mock struct {
	locationA         uint8
	locationB         uint8
	pwmChannelA       uint8
	pwmChannelB       uint8
	pwmPrescalerValue uint
	pin               uint8
	mode              uint8
	state             uint8
}

// setup mock for GPIO, PWM and servo tests
func (l *mock) digitalWrite(pin uint8, state uint8) error {
	l.pin = pin
	l.state = state
	return l.error()
}

func (l *mock) pinMode(pin uint8, mode uint8) error {
	l.pin = pin
	l.mode = mode
	return l.error()
}

var pwmInitErrorFunc = func() error { return nil }

func (l *mock) pwmInit() error { return pwmInitErrorFunc() }
func (l *mock) pwmStop() error { return l.error() }
func (l *mock) pwmUpdateCompare(channelA uint8, channelB uint8) error {
	l.pwmChannelA = channelA
	l.pwmChannelB = channelB
	return l.error()
}

func (l *mock) pwmUpdatePrescaler(value uint) error {
	l.pwmPrescalerValue = value
	return l.error()
}
func (l *mock) servoInit() error { return l.error() }
func (l *mock) servoUpdateLocation(locationA uint8, locationB uint8) error {
	l.locationA = locationA
	l.locationB = locationB
	return l.error()
}

var errorFunc = func() error { return nil }

func (l *mock) error() error { return errorFunc() }

// i2c functions unused in this test scenarios
func (l *mock) i2cInit() error                                                  { return nil }
func (l *mock) i2cStart(address7bit uint8, direction uint8) error               { return nil }
func (l *mock) i2cWrite(sendBuffer []byte, length int, endWithStop uint8) error { return nil }
func (l *mock) i2cRead(readBuffer []byte, length int, endWithStop uint8) error  { return nil }
func (l *mock) i2cUpdateDelay(duration uint) error                              { return nil }

func initTestAdaptor() *Adaptor {
	a := NewAdaptor()
	a.connect = func(a *Adaptor) error { return nil }
	a.littleWire = new(mock)
	errorFunc = func() error { return nil }
	pwmInitErrorFunc = func() error { return nil }
	return a
}

func TestDigisparkAdaptorName(t *testing.T) {
	a := NewAdaptor()
	assert.True(t, strings.HasPrefix(a.Name(), "Digispark"))
	a.SetName("NewName")
	assert.Equal(t, "NewName", a.Name())
}

func TestDigisparkAdaptorConnect(t *testing.T) {
	a := initTestAdaptor()
	require.NoError(t, a.Connect())
}

func TestDigisparkAdaptorFinalize(t *testing.T) {
	a := initTestAdaptor()
	require.NoError(t, a.Finalize())
}

func TestDigisparkAdaptorDigitalWrite(t *testing.T) {
	a := initTestAdaptor()
	err := a.DigitalWrite("0", uint8(1))
	require.NoError(t, err)
	assert.Equal(t, uint8(0), a.littleWire.(*mock).pin)
	assert.Equal(t, uint8(1), a.littleWire.(*mock).state)

	err = a.DigitalWrite("?", uint8(1))
	require.Error(t, err)

	errorFunc = func() error { return errors.New("pin mode error") }
	err = a.DigitalWrite("0", uint8(1))
	require.ErrorContains(t, err, "pin mode error")
}

func TestDigisparkAdaptorServoWrite(t *testing.T) {
	a := initTestAdaptor()
	err := a.ServoWrite("2", uint8(80))
	require.NoError(t, err)
	assert.Equal(t, uint8(80), a.littleWire.(*mock).locationA)
	assert.Equal(t, uint8(80), a.littleWire.(*mock).locationB)

	a = initTestAdaptor()
	errorFunc = func() error { return errors.New("servo error") }
	err = a.ServoWrite("2", uint8(80))
	require.ErrorContains(t, err, "servo error")
}

func TestDigisparkAdaptorPwmWrite(t *testing.T) {
	a := initTestAdaptor()
	err := a.PwmWrite("1", uint8(100))
	require.NoError(t, err)
	assert.Equal(t, uint8(100), a.littleWire.(*mock).pwmChannelA)
	assert.Equal(t, uint8(100), a.littleWire.(*mock).pwmChannelB)

	a = initTestAdaptor()
	pwmInitErrorFunc = func() error { return errors.New("pwminit error") }
	err = a.PwmWrite("1", uint8(100))
	require.ErrorContains(t, err, "pwminit error")

	a = initTestAdaptor()
	errorFunc = func() error { return errors.New("pwm error") }
	err = a.PwmWrite("1", uint8(100))
	require.ErrorContains(t, err, "pwm error")
}
