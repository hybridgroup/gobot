package digispark

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Adaptor = (*Adaptor)(nil)

var _ gpio.DigitalWriter = (*Adaptor)(nil)
var _ gpio.PwmWriter = (*Adaptor)(nil)
var _ gpio.ServoWriter = (*Adaptor)(nil)

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

func initTestAdaptor() *Adaptor {
	a := NewAdaptor()
	a.connect = func(a *Adaptor) (err error) { return nil }
	a.littleWire = new(mock)
	errorFunc = func() error { return nil }
	pwmInitErrorFunc = func() error { return nil }
	return a
}

func TestDigisparkAdaptorName(t *testing.T) {
	a := NewAdaptor()
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "Digispark"), true)
	a.SetName("NewName")
	gobottest.Assert(t, a.Name(), "NewName")
}

func TestAdaptorConnect(t *testing.T) {
	a := NewAdaptor()
	gobottest.Assert(t, a.Connect(), ErrConnection)

	a = initTestAdaptor()
	gobottest.Assert(t, a.Connect(), nil)
}

func TestAdaptorFinalize(t *testing.T) {
	a := initTestAdaptor()
	gobottest.Assert(t, a.Finalize(), nil)
}

func TestAdaptorDigitalWrite(t *testing.T) {
	a := initTestAdaptor()
	err := a.DigitalWrite("0", uint8(1))
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, a.littleWire.(*mock).pin, uint8(0))
	gobottest.Assert(t, a.littleWire.(*mock).state, uint8(1))

	err = a.DigitalWrite("?", uint8(1))
	gobottest.Refute(t, err, nil)

	errorFunc = func() error { return errors.New("pin mode error") }
	err = a.DigitalWrite("0", uint8(1))
	gobottest.Assert(t, err, errors.New("pin mode error"))
}

func TestAdaptorServoWrite(t *testing.T) {
	a := initTestAdaptor()
	err := a.ServoWrite("2", uint8(80))
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, a.littleWire.(*mock).locationA, uint8(80))
	gobottest.Assert(t, a.littleWire.(*mock).locationB, uint8(80))

	a = initTestAdaptor()
	errorFunc = func() error { return errors.New("servo error") }
	err = a.ServoWrite("2", uint8(80))
	gobottest.Assert(t, err, errors.New("servo error"))
}

func TestAdaptorPwmWrite(t *testing.T) {
	a := initTestAdaptor()
	err := a.PwmWrite("1", uint8(100))
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, a.littleWire.(*mock).pwmChannelA, uint8(100))
	gobottest.Assert(t, a.littleWire.(*mock).pwmChannelB, uint8(100))

	a = initTestAdaptor()
	pwmInitErrorFunc = func() error { return errors.New("pwminit error") }
	err = a.PwmWrite("1", uint8(100))
	gobottest.Assert(t, err, errors.New("pwminit error"))

	a = initTestAdaptor()
	errorFunc = func() error { return errors.New("pwm error") }
	err = a.PwmWrite("1", uint8(100))
	gobottest.Assert(t, err, errors.New("pwm error"))
}
