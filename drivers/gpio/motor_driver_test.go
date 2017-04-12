package gpio

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*MotorDriver)(nil)

func initTestMotorDriver() *MotorDriver {
	return NewMotorDriver(newGpioTestAdaptor(), "1")
}

func TestMotorDriver(t *testing.T) {
	d := NewMotorDriver(newGpioTestAdaptor(), "1")
	gobottest.Refute(t, d.Connection(), nil)

}
func TestMotorDriverStart(t *testing.T) {
	d := initTestMotorDriver()
	gobottest.Assert(t, d.Start(), nil)
}

func TestMotorDriverHalt(t *testing.T) {
	d := initTestMotorDriver()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestMotorDriverIsOn(t *testing.T) {
	d := initTestMotorDriver()
	d.CurrentMode = "digital"
	d.CurrentState = 1
	gobottest.Assert(t, d.IsOn(), true)
	d.CurrentMode = "analog"
	d.CurrentSpeed = 100
	gobottest.Assert(t, d.IsOn(), true)
}

func TestMotorDriverIsOff(t *testing.T) {
	d := initTestMotorDriver()
	d.Off()
	gobottest.Assert(t, d.IsOff(), true)
}

func TestMotorDriverOn(t *testing.T) {
	d := initTestMotorDriver()
	d.CurrentMode = "digital"
	d.On()
	gobottest.Assert(t, d.CurrentState, uint8(1))
	d.CurrentMode = "analog"
	d.CurrentSpeed = 0
	d.On()
	gobottest.Assert(t, d.CurrentSpeed, uint8(255))
}

func TestMotorDriverOff(t *testing.T) {
	d := initTestMotorDriver()
	d.CurrentMode = "digital"
	d.Off()
	gobottest.Assert(t, d.CurrentState, uint8(0))
	d.CurrentMode = "analog"
	d.CurrentSpeed = 100
	d.Off()
	gobottest.Assert(t, d.CurrentSpeed, uint8(0))
}

func TestMotorDriverToggle(t *testing.T) {
	d := initTestMotorDriver()
	d.Off()
	d.Toggle()
	gobottest.Assert(t, d.IsOn(), true)
	d.Toggle()
	gobottest.Assert(t, d.IsOn(), false)
}

func TestMotorDriverMin(t *testing.T) {
	d := initTestMotorDriver()
	d.Min()
}

func TestMotorDriverMax(t *testing.T) {
	d := initTestMotorDriver()
	d.Max()
}

func TestMotorDriverSpeed(t *testing.T) {
	d := initTestMotorDriver()
	d.Speed(100)
}

func TestMotorDriverForward(t *testing.T) {
	d := initTestMotorDriver()
	d.Forward(100)
	gobottest.Assert(t, d.CurrentSpeed, uint8(100))
	gobottest.Assert(t, d.CurrentDirection, "forward")
}

func TestMotorDriverBackward(t *testing.T) {
	d := initTestMotorDriver()
	d.Backward(100)
	gobottest.Assert(t, d.CurrentSpeed, uint8(100))
	gobottest.Assert(t, d.CurrentDirection, "backward")
}

func TestMotorDriverDirection(t *testing.T) {
	d := initTestMotorDriver()
	d.Direction("none")
	d.DirectionPin = "2"
	d.Direction("forward")
	d.Direction("backward")
}

func TestMotorDriverDigital(t *testing.T) {
	d := initTestMotorDriver()
	d.SpeedPin = "" // Disable speed
	d.CurrentMode = "digital"
	d.ForwardPin = "2"
	d.BackwardPin = "3"

	d.On()
	gobottest.Assert(t, d.CurrentState, uint8(1))
	d.Off()
	gobottest.Assert(t, d.CurrentState, uint8(0))
}

func TestMotorDriverDefaultName(t *testing.T) {
	d := initTestMotorDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Motor"), true)
}

func TestMotorDriverSetName(t *testing.T) {
	d := initTestMotorDriver()
	d.SetName("mybot")
	gobottest.Assert(t, d.Name(), "mybot")
}
