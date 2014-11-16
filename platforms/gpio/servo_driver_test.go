package gpio

import (
	"testing"

	"github.com/hybridgroup/gobot"
)

func initTestServoDriver() *ServoDriver {
	return NewServoDriver(newGpioTestAdaptor("adaptor"), "bot", "1")
}

func TestServoDriverStart(t *testing.T) {
	d := initTestServoDriver()
	gobot.Assert(t, d.Start(), nil)
}

func TestServoDriverHalt(t *testing.T) {
	d := initTestServoDriver()
	gobot.Assert(t, d.Halt(), nil)
}

func TestServoDriverMove(t *testing.T) {
	d := initTestServoDriver()
	d.Move(100)
	gobot.Assert(t, d.CurrentAngle, uint8(100))
}

func TestServoDriverMin(t *testing.T) {
	d := initTestServoDriver()
	d.Min()
	gobot.Assert(t, d.CurrentAngle, uint8(0))
}

func TestServoDriverMax(t *testing.T) {
	d := initTestServoDriver()
	d.Max()
	gobot.Assert(t, d.CurrentAngle, uint8(180))
}

func TestServoDriverCenter(t *testing.T) {
	d := initTestServoDriver()
	d.Center()
	gobot.Assert(t, d.CurrentAngle, uint8(90))
}

func TestServoDriverInitServo(t *testing.T) {
	d := initTestServoDriver()
	d.InitServo()
}
