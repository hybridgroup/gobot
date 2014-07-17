package ardrone

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestArdroneDriver() *ArdroneDriver {
	a := NewArdroneAdaptor("drone")
	a.connect = func(a *ArdroneAdaptor) {
		a.drone = &testDrone{}
	}
	d := NewArdroneDriver(a, "drone")
	a.Connect()
	return d
}

func TestArdroneDriverStart(t *testing.T) {
	d := initTestArdroneDriver()
	gobot.Assert(t, d.Start(), true)
}

func TestArdroneDriverHalt(t *testing.T) {
	d := initTestArdroneDriver()
	gobot.Assert(t, d.Halt(), true)
}
func TestArdroneDriverTakeOff(t *testing.T) {
	d := initTestArdroneDriver()
	d.TakeOff()
}

func TestArdroneDriverand(t *testing.T) {
	d := initTestArdroneDriver()
	d.Land()
}

func TestArdroneDriverUp(t *testing.T) {
	d := initTestArdroneDriver()
	d.Up(1)
}

func TestArdroneDriverDown(t *testing.T) {
	d := initTestArdroneDriver()
	d.Down(1)
}

func TestArdroneDriverLeft(t *testing.T) {
	d := initTestArdroneDriver()
	d.Left(1)
}

func TestArdroneDriverRight(t *testing.T) {
	d := initTestArdroneDriver()
	d.Right(1)
}

func TestArdroneDriverForward(t *testing.T) {
	d := initTestArdroneDriver()
	d.Forward(1)
}

func TestArdroneDriverackward(t *testing.T) {
	d := initTestArdroneDriver()
	d.Backward(1)
}

func TestArdroneDriverClockwise(t *testing.T) {
	d := initTestArdroneDriver()
	d.Clockwise(1)
}

func TestArdroneDriverCounterClockwise(t *testing.T) {
	d := initTestArdroneDriver()
	d.CounterClockwise(1)
}

func TestArdroneDriverHover(t *testing.T) {
	d := initTestArdroneDriver()
	d.Hover()
}
