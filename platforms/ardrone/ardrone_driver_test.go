package ardrone

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

var driver *ArdroneDriver

func init() {
	adaptor := NewArdroneAdaptor("drone")
	adaptor.connect = func(a *ArdroneAdaptor) {
		a.drone = &testDrone{}
	}
	driver = NewArdroneDriver(adaptor, "drone")
	adaptor.Connect()
}

func TestStart(t *testing.T) {
	gobot.Expect(t, driver.Start(), true)
}

func TestHalt(t *testing.T) {
	gobot.Expect(t, driver.Halt(), true)
}
func TestTakeOff(t *testing.T) {
	driver.TakeOff()
}

func TestLand(t *testing.T) {
	driver.Land()
}
func TestUp(t *testing.T) {
	driver.Up(1)
}

func TestDown(t *testing.T) {
	driver.Down(1)
}

func TestLeft(t *testing.T) {
	driver.Left(1)
}

func TestRight(t *testing.T) {
	driver.Right(1)
}

func TestForward(t *testing.T) {
	driver.Forward(1)
}

func TestBackward(t *testing.T) {
	driver.Backward(1)
}

func TestClockwise(t *testing.T) {
	driver.Clockwise(1)
}

func TestCounterClockwise(t *testing.T) {
	driver.CounterClockwise(1)
}

func TestHover(t *testing.T) {
	driver.Hover()
}
