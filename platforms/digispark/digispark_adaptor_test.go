package digispark

import (
	"testing"

	"github.com/hybridgroup/gobot"
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
func (l *mock) pwmInit() error { return l.error() }
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
func (l *mock) error() error { return nil }

func initTestDigisparkAdaptor() *DigisparkAdaptor {
	a := NewDigisparkAdaptor("bot")
	a.connect = func(a *DigisparkAdaptor) (err error) { return nil }
	a.littleWire = new(mock)
	return a
}

func TestDigisparkAdaptorFinalize(t *testing.T) {
	a := initTestDigisparkAdaptor()
	gobot.Assert(t, a.Finalize(), nil)
}

func TestDigisparkAdaptorConnect(t *testing.T) {
	a := initTestDigisparkAdaptor()
	gobot.Assert(t, a.Connect(), nil)
}

func TestDigisparkAdaptorIO(t *testing.T) {
	a := initTestDigisparkAdaptor()
	a.InitServo()
	a.DigitalRead("1")
	a.DigitalWrite("0", uint8(1))
	gobot.Assert(t, a.littleWire.(*mock).pin, uint8(0))
	gobot.Assert(t, a.littleWire.(*mock).state, uint8(1))
	a.PwmWrite("1", uint8(100))
	gobot.Assert(t, a.littleWire.(*mock).pwmChannelA, uint8(100))
	gobot.Assert(t, a.littleWire.(*mock).pwmChannelB, uint8(100))
	a.ServoWrite("2", uint8(80))
	gobot.Assert(t, a.littleWire.(*mock).locationA, uint8(80))
	gobot.Assert(t, a.littleWire.(*mock).locationB, uint8(80))
}
