package gpio

import (
	"strings"
	"testing"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/gobottest"
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
	_ = d.Off()
	gobottest.Assert(t, d.IsOff(), true)
}

func TestMotorDriverOn(t *testing.T) {
	d := initTestMotorDriver()
	d.CurrentMode = "digital"
	_ = d.On()
	gobottest.Assert(t, d.CurrentState, uint8(1))
	d.CurrentMode = "analog"
	d.CurrentSpeed = 0
	_ = d.On()
	gobottest.Assert(t, d.CurrentSpeed, uint8(255))
}

func TestMotorDriverOff(t *testing.T) {
	d := initTestMotorDriver()
	d.CurrentMode = "digital"
	_ = d.Off()
	gobottest.Assert(t, d.CurrentState, uint8(0))
	d.CurrentMode = "analog"
	d.CurrentSpeed = 100
	_ = d.Off()
	gobottest.Assert(t, d.CurrentSpeed, uint8(0))
}

func TestMotorDriverToggle(t *testing.T) {
	d := initTestMotorDriver()
	_ = d.Off()
	_ = d.Toggle()
	gobottest.Assert(t, d.IsOn(), true)
	_ = d.Toggle()
	gobottest.Assert(t, d.IsOn(), false)
}

func TestMotorDriverMin(t *testing.T) {
	d := initTestMotorDriver()
	_ = d.Min()
}

func TestMotorDriverMax(t *testing.T) {
	d := initTestMotorDriver()
	_ = d.Max()
}

func TestMotorDriverSpeed(t *testing.T) {
	d := initTestMotorDriver()
	_ = d.Speed(100)
}

func TestMotorDriverForward(t *testing.T) {
	d := initTestMotorDriver()
	_ = d.Forward(100)
	gobottest.Assert(t, d.CurrentSpeed, uint8(100))
	gobottest.Assert(t, d.CurrentDirection, "forward")
}

func TestMotorDriverBackward(t *testing.T) {
	d := initTestMotorDriver()
	_ = d.Backward(100)
	gobottest.Assert(t, d.CurrentSpeed, uint8(100))
	gobottest.Assert(t, d.CurrentDirection, "backward")
}

func TestMotorDriverDirection(t *testing.T) {
	d := initTestMotorDriver()
	_ = d.Direction("none")
	d.DirectionPin = "2"
	_ = d.Direction("forward")
	_ = d.Direction("backward")
}

func TestMotorDriverDigital(t *testing.T) {
	d := initTestMotorDriver()
	d.SpeedPin = "" // Disable speed
	d.CurrentMode = "digital"
	d.ForwardPin = "2"
	d.BackwardPin = "3"

	_ = d.On()
	gobottest.Assert(t, d.CurrentState, uint8(1))
	_ = d.Off()
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
