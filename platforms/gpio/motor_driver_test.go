package gpio

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

var m *MotorDriver

func init() {
	m = NewMotorDriver(TestAdaptor{}, "bot", "1")
}

func TestMotorStart(t *testing.T) {
	gobot.Expect(t, m.Start(), true)
}

func TestMotorHalt(t *testing.T) {
	gobot.Expect(t, m.Halt(), true)
}

func TestMotorInit(t *testing.T) {
	gobot.Expect(t, m.Init(), true)
}

func TestMotorIsOn(t *testing.T) {
	m.CurrentMode = "digital"
	m.CurrentState = 1
	gobot.Expect(t, m.IsOn(), true)
	m.CurrentMode = "analog"
	m.CurrentSpeed = 100
	gobot.Expect(t, m.IsOn(), true)
}

func TestMotorIsOff(t *testing.T) {
	m.Off()
	gobot.Expect(t, m.IsOff(), true)
}

func TestMotorOn(t *testing.T) {
	m.CurrentMode = "digital"
	m.On()
	gobot.Expect(t, m.CurrentState, uint8(1))
	m.CurrentMode = "analog"
	m.CurrentSpeed = 0
	m.On()
	gobot.Expect(t, m.CurrentSpeed, uint8(255))
}

func TestMotorOff(t *testing.T) {
	m.CurrentMode = "digital"
	m.Off()
	gobot.Expect(t, m.CurrentState, uint8(0))
	m.CurrentMode = "analog"
	m.CurrentSpeed = 100
	m.Off()
	gobot.Expect(t, m.CurrentSpeed, uint8(0))
}

func TestMotorToggle(t *testing.T) {
	m.Off()
	m.Toggle()
	gobot.Expect(t, m.IsOn(), true)
	m.Toggle()
	gobot.Expect(t, m.IsOn(), false)
}

func TestMotorMin(t *testing.T) {
	m.Min()
}

func TestMotorMax(t *testing.T) {
	m.Max()
}

func TestMotorSpeed(t *testing.T) {
	m.Speed(100)
}

func TestMotorForward(t *testing.T) {
	m.Forward(100)
	gobot.Expect(t, m.CurrentSpeed, uint8(100))
	gobot.Expect(t, m.CurrentDirection, "forward")
}
func TestMotorBackward(t *testing.T) {
	m.Backward(100)
	gobot.Expect(t, m.CurrentSpeed, uint8(100))
	gobot.Expect(t, m.CurrentDirection, "backward")
}

func TestMotorDirection(t *testing.T) {
	m.Direction("none")
	m.DirectionPin = "2"
	m.Direction("forward")
	m.Direction("backward")
}
