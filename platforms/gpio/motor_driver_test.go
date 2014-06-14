package gpio

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestMotorDriver() *MotorDriver {
	return NewMotorDriver(TestAdaptor{}, "bot", "1")
}

func TestMotorDriverStart(t *testing.T) {
	d := initTestMotorDriver()
	gobot.Expect(t, d.Start(), true)
}

func TestMotorDriverHalt(t *testing.T) {
	d := initTestMotorDriver()
	gobot.Expect(t, d.Halt(), true)
}

func TestMotorDriverInit(t *testing.T) {
	d := initTestMotorDriver()
	gobot.Expect(t, d.Init(), true)
}

func TestMotorDriverIsOn(t *testing.T) {
	d := initTestMotorDriver()
	d.CurrentMode = "digital"
	d.CurrentState = 1
	gobot.Expect(t, d.IsOn(), true)
	d.CurrentMode = "analog"
	d.CurrentSpeed = 100
	gobot.Expect(t, d.IsOn(), true)
}

func TestMotorDriverIsOff(t *testing.T) {
	d := initTestMotorDriver()
	d.Off()
	gobot.Expect(t, d.IsOff(), true)
}

func TestMotorDriverOn(t *testing.T) {
	d := initTestMotorDriver()
	d.CurrentMode = "digital"
	d.On()
	gobot.Expect(t, d.CurrentState, uint8(1))
	d.CurrentMode = "analog"
	d.CurrentSpeed = 0
	d.On()
	gobot.Expect(t, d.CurrentSpeed, uint8(255))
}

func TestMotorDriverOff(t *testing.T) {
	d := initTestMotorDriver()
	d.CurrentMode = "digital"
	d.Off()
	gobot.Expect(t, d.CurrentState, uint8(0))
	d.CurrentMode = "analog"
	d.CurrentSpeed = 100
	d.Off()
	gobot.Expect(t, d.CurrentSpeed, uint8(0))
}

func TestMotorDriverToggle(t *testing.T) {
	d := initTestMotorDriver()
	d.Off()
	d.Toggle()
	gobot.Expect(t, d.IsOn(), true)
	d.Toggle()
	gobot.Expect(t, d.IsOn(), false)
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
	gobot.Expect(t, d.CurrentSpeed, uint8(100))
	gobot.Expect(t, d.CurrentDirection, "forward")
}
func TestMotorDriverBackward(t *testing.T) {
	d := initTestMotorDriver()
	d.Backward(100)
	gobot.Expect(t, d.CurrentSpeed, uint8(100))
	gobot.Expect(t, d.CurrentDirection, "backward")
}

func TestMotorDriverDirection(t *testing.T) {
	d := initTestMotorDriver()
	d.Direction("none")
	d.DirectionPin = "2"
	d.Direction("forward")
	d.Direction("backward")
}
