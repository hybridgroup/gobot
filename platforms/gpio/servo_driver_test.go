package gpio

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestServoDriver() *ServoDriver {
	return NewServoDriver(newGpioTestAdaptor("adaptor"), "bot", "1")
}

func TestServoDriverStart(t *testing.T) {
	d := initTestServoDriver()
	gobot.Assert(t, d.Start(), true)
}

func TestServoDriverHalt(t *testing.T) {
	d := initTestServoDriver()
	gobot.Assert(t, d.Halt(), true)
}

func TestServoDriverInit(t *testing.T) {
	d := initTestServoDriver()
	gobot.Assert(t, d.Init(), true)
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
