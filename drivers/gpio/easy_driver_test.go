package gpio

import (
	"strings"
	"testing"
	"time"

	"gobot.io/x/gobot/v2/gobottest"
)

const (
	stepAngle   = 0.5 // use non int step angle to check int math
	stepsPerRev = 720
)

var adapter *gpioTestAdaptor

func initEasyDriver() *EasyDriver {
	adapter = newGpioTestAdaptor()
	return NewEasyDriver(adapter, stepAngle, "1", "2", "3", "4")
}

func TestEasyDriver_Connection(t *testing.T) {
	d := initEasyDriver()
	gobottest.Assert(t, d.Connection(), adapter)
}

func TestEasyDriverDefaultName(t *testing.T) {
	d := initEasyDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "EasyDriver"), true)
}

func TestEasyDriverSetName(t *testing.T) {
	d := initEasyDriver()
	d.SetName("OtherDriver")
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "OtherDriver"), true)
}

func TestEasyDriverStart(t *testing.T) {
	d := initEasyDriver()
	_ = d.Start()
	// noop - no error occurred
}

func TestEasyDriverHalt(t *testing.T) {
	d := initEasyDriver()
	_ = d.Run()
	gobottest.Assert(t, d.IsMoving(), true)
	_ = d.Halt()
	gobottest.Assert(t, d.IsMoving(), false)
}

func TestEasyDriverMove(t *testing.T) {
	d := initEasyDriver()
	_ = d.Move(2)
	time.Sleep(2 * time.Millisecond)
	gobottest.Assert(t, d.GetCurrentStep(), 4)
	gobottest.Assert(t, d.IsMoving(), false)
}

func TestEasyDriverRun(t *testing.T) {
	d := initEasyDriver()
	_ = d.Run()
	gobottest.Assert(t, d.IsMoving(), true)
	_ = d.Run()
	gobottest.Assert(t, d.IsMoving(), true)
}

func TestEasyDriverStop(t *testing.T) {
	d := initEasyDriver()
	_ = d.Run()
	gobottest.Assert(t, d.IsMoving(), true)
	_ = d.Stop()
	gobottest.Assert(t, d.IsMoving(), false)
}

func TestEasyDriverStep(t *testing.T) {
	d := initEasyDriver()
	_ = d.Step()
	gobottest.Assert(t, d.GetCurrentStep(), 1)
	_ = d.Step()
	_ = d.Step()
	_ = d.Step()
	gobottest.Assert(t, d.GetCurrentStep(), 4)
	_ = d.SetDirection("ccw")
	_ = d.Step()
	gobottest.Assert(t, d.GetCurrentStep(), 3)
}

func TestEasyDriverSetDirection(t *testing.T) {
	d := initEasyDriver()
	gobottest.Assert(t, d.dir, int8(1))
	_ = d.SetDirection("cw")
	gobottest.Assert(t, d.dir, int8(1))
	_ = d.SetDirection("ccw")
	gobottest.Assert(t, d.dir, int8(-1))
	_ = d.SetDirection("nothing")
	gobottest.Assert(t, d.dir, int8(1))
}

func TestEasyDriverSetDirectionNoPin(t *testing.T) {
	d := initEasyDriver()
	d.dirPin = ""
	err := d.SetDirection("cw")
	gobottest.Refute(t, err, nil)
}

func TestEasyDriverSetSpeed(t *testing.T) {
	d := initEasyDriver()
	gobottest.Assert(t, d.rpm, uint(stepsPerRev/4)) // default speed of 720/4
	_ = d.SetSpeed(0)
	gobottest.Assert(t, d.rpm, uint(1))
	_ = d.SetSpeed(200)
	gobottest.Assert(t, d.rpm, uint(200))
	_ = d.SetSpeed(1000)
	gobottest.Assert(t, d.rpm, uint(stepsPerRev))
}

func TestEasyDriverGetMaxSpeed(t *testing.T) {
	d := initEasyDriver()
	gobottest.Assert(t, d.GetMaxSpeed(), uint(stepsPerRev))
}

func TestEasyDriverSleep(t *testing.T) {
	// let's test basic functionality
	d := initEasyDriver()
	_ = d.Sleep()
	gobottest.Assert(t, d.IsSleeping(), true)

	// let's make sure it stops first
	d = initEasyDriver()
	_ = d.Run()
	_ = d.Sleep()
	gobottest.Assert(t, d.IsSleeping(), true)
	gobottest.Assert(t, d.IsMoving(), false)
}

func TestEasyDriverSleepNoPin(t *testing.T) {
	d := initEasyDriver()
	d.sleepPin = ""
	err := d.Sleep()
	gobottest.Refute(t, err, nil)
	err = d.Wake()
	gobottest.Refute(t, err, nil)
}

func TestEasyDriverWake(t *testing.T) {
	// let's test basic functionality
	d := initEasyDriver()
	_ = d.Sleep()
	gobottest.Assert(t, d.IsSleeping(), true)
	_ = d.Wake()
	gobottest.Assert(t, d.IsSleeping(), false)
}

func TestEasyDriverEnable(t *testing.T) {
	// let's test basic functionality
	d := initEasyDriver()
	_ = d.Disable()
	gobottest.Assert(t, d.IsEnabled(), false)
	_ = d.Enable()
	gobottest.Assert(t, d.IsEnabled(), true)
}

func TestEasyDriverEnableNoPin(t *testing.T) {
	d := initEasyDriver()
	d.enPin = ""
	err := d.Disable()
	gobottest.Refute(t, err, nil)
	err = d.Enable()
	gobottest.Refute(t, err, nil)
}

func TestEasyDriverDisable(t *testing.T) {
	// let's test basic functionality
	d := initEasyDriver()
	_ = d.Disable()
	gobottest.Assert(t, d.IsEnabled(), false)

	// let's make sure it stops first
	d = initEasyDriver()
	_ = d.Run()
	_ = d.Disable()
	gobottest.Assert(t, d.IsEnabled(), false)
	gobottest.Assert(t, d.IsMoving(), false)
}
