package minidrone

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*Driver)(nil)

func initTestMinidroneDriver() *Driver {
	d := NewDriver(NewBleTestAdaptor())
	return d
}

func TestMinidroneDriver(t *testing.T) {
	d := initTestMinidroneDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Minidrone"), true)
	d.SetName("NewName")
	gobottest.Assert(t, d.Name(), "NewName")
}

func TestMinidroneDriverStartAndHalt(t *testing.T) {
	d := initTestMinidroneDriver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Halt(), nil)
}

func TestMinidroneTakeoff(t *testing.T) {
	d := initTestMinidroneDriver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.TakeOff(), nil)
}

func TestMinidroneEmergency(t *testing.T) {
	d := initTestMinidroneDriver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Emergency(), nil)
}

func TestMinidroneTakePicture(t *testing.T) {
	d := initTestMinidroneDriver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.TakePicture(), nil)
}

func TestMinidroneUp(t *testing.T) {
	d := initTestMinidroneDriver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Up(25), nil)
}

func TestMinidroneDown(t *testing.T) {
	d := initTestMinidroneDriver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Down(25), nil)
}

func TestMinidroneForward(t *testing.T) {
	d := initTestMinidroneDriver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Forward(25), nil)
}

func TestMinidroneBackward(t *testing.T) {
	d := initTestMinidroneDriver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Backward(25), nil)
}

func TestMinidroneRight(t *testing.T) {
	d := initTestMinidroneDriver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Right(25), nil)
}

func TestMinidroneLeft(t *testing.T) {
	d := initTestMinidroneDriver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Left(25), nil)
}

func TestMinidroneClockwise(t *testing.T) {
	d := initTestMinidroneDriver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Clockwise(25), nil)
}

func TestMinidroneCounterClockwise(t *testing.T) {
	d := initTestMinidroneDriver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.CounterClockwise(25), nil)
}

func TestMinidroneStartStopRecording(t *testing.T) {
	d := initTestMinidroneDriver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.StartRecording(), nil)
	gobottest.Assert(t, d.StopRecording(), nil)
}

func TestMinidroneHullProtectionOutdoor(t *testing.T) {
	d := initTestMinidroneDriver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.HullProtection(true), nil)
	gobottest.Assert(t, d.Outdoor(true), nil)
}

func TestMinidroneHullFlips(t *testing.T) {
	d := initTestMinidroneDriver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.FrontFlip(), nil)
	gobottest.Assert(t, d.BackFlip(), nil)
	gobottest.Assert(t, d.RightFlip(), nil)
	gobottest.Assert(t, d.LeftFlip(), nil)
}
