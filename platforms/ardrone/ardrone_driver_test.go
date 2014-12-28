package ardrone

import (
	"testing"

	"github.com/hybridgroup/gobot"
)

func initTestArdroneDriver() *ArdroneDriver {
	a := NewArdroneAdaptor("drone")
	a.connect = func(a *ArdroneAdaptor) (drone, error) {
		return &testDrone{}, nil
	}
	d := NewArdroneDriver(a, "drone")
	a.Connect()
	return d
}

func TestArdroneDriver(t *testing.T) {
	d := initTestArdroneDriver()
	gobot.Assert(t, d.Name(), "drone")
	gobot.Assert(t, d.Connection().Name(), "drone")
}

func TestArdroneDriverStart(t *testing.T) {
	d := initTestArdroneDriver()
	gobot.Assert(t, len(d.Start()), 0)
}

func TestArdroneDriverHalt(t *testing.T) {
	d := initTestArdroneDriver()
	gobot.Assert(t, len(d.Halt()), 0)
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
