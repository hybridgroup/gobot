package gpio

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, adapter, d.Connection())
}

func TestEasyDriverDefaultName(t *testing.T) {
	d := initEasyDriver()
	assert.True(t, strings.HasPrefix(d.Name(), "EasyDriver"))
}

func TestEasyDriverSetName(t *testing.T) {
	d := initEasyDriver()
	d.SetName("OtherDriver")
	assert.True(t, strings.HasPrefix(d.Name(), "OtherDriver"))
}

func TestEasyDriverStart(t *testing.T) {
	d := initEasyDriver()
	_ = d.Start()
	// noop - no error occurred
}

func TestEasyDriverHalt(t *testing.T) {
	d := initEasyDriver()
	_ = d.Run()
	assert.True(t, d.IsMoving())
	_ = d.Halt()
	assert.False(t, d.IsMoving())
}

func TestEasyDriverMove(t *testing.T) {
	d := initEasyDriver()
	_ = d.Move(2)
	time.Sleep(2 * time.Millisecond)
	assert.Equal(t, 4, d.GetCurrentStep())
	assert.False(t, d.IsMoving())
}

func TestEasyDriverRun(t *testing.T) {
	d := initEasyDriver()
	_ = d.Run()
	assert.True(t, d.IsMoving())
	_ = d.Run()
	assert.True(t, d.IsMoving())
}

func TestEasyDriverStop(t *testing.T) {
	d := initEasyDriver()
	_ = d.Run()
	assert.True(t, d.IsMoving())
	_ = d.Stop()
	assert.False(t, d.IsMoving())
}

func TestEasyDriverStep(t *testing.T) {
	d := initEasyDriver()
	_ = d.Step()
	assert.Equal(t, 1, d.GetCurrentStep())
	_ = d.Step()
	_ = d.Step()
	_ = d.Step()
	assert.Equal(t, 4, d.GetCurrentStep())
	_ = d.SetDirection("ccw")
	_ = d.Step()
	assert.Equal(t, 3, d.GetCurrentStep())
}

func TestEasyDriverSetDirection(t *testing.T) {
	d := initEasyDriver()
	assert.Equal(t, int8(1), d.dir)
	_ = d.SetDirection("cw")
	assert.Equal(t, int8(1), d.dir)
	_ = d.SetDirection("ccw")
	assert.Equal(t, int8(-1), d.dir)
	_ = d.SetDirection("nothing")
	assert.Equal(t, int8(1), d.dir)
}

func TestEasyDriverSetDirectionNoPin(t *testing.T) {
	d := initEasyDriver()
	d.dirPin = ""
	err := d.SetDirection("cw")
	assert.NotNil(t, err)
}

func TestEasyDriverSetSpeed(t *testing.T) {
	d := initEasyDriver()
	assert.Equal(t, uint(stepsPerRev/4), d.rpm) // default speed of 720/4
	_ = d.SetSpeed(0)
	assert.Equal(t, uint(1), d.rpm)
	_ = d.SetSpeed(200)
	assert.Equal(t, uint(200), d.rpm)
	_ = d.SetSpeed(1000)
	assert.Equal(t, uint(stepsPerRev), d.rpm)
}

func TestEasyDriverGetMaxSpeed(t *testing.T) {
	d := initEasyDriver()
	assert.Equal(t, uint(stepsPerRev), d.GetMaxSpeed())
}

func TestEasyDriverSleep(t *testing.T) {
	// let's test basic functionality
	d := initEasyDriver()
	_ = d.Sleep()
	assert.True(t, d.IsSleeping())

	// let's make sure it stops first
	d = initEasyDriver()
	_ = d.Run()
	_ = d.Sleep()
	assert.True(t, d.IsSleeping())
	assert.False(t, d.IsMoving())
}

func TestEasyDriverSleepNoPin(t *testing.T) {
	d := initEasyDriver()
	d.sleepPin = ""
	err := d.Sleep()
	assert.NotNil(t, err)
	err = d.Wake()
	assert.NotNil(t, err)
}

func TestEasyDriverWake(t *testing.T) {
	// let's test basic functionality
	d := initEasyDriver()
	_ = d.Sleep()
	assert.True(t, d.IsSleeping())
	_ = d.Wake()
	assert.False(t, d.IsSleeping())
}

func TestEasyDriverEnable(t *testing.T) {
	// let's test basic functionality
	d := initEasyDriver()
	_ = d.Disable()
	assert.False(t, d.IsEnabled())
	_ = d.Enable()
	assert.True(t, d.IsEnabled())
}

func TestEasyDriverEnableNoPin(t *testing.T) {
	d := initEasyDriver()
	d.enPin = ""
	err := d.Disable()
	assert.NotNil(t, err)
	err = d.Enable()
	assert.NotNil(t, err)
}

func TestEasyDriverDisable(t *testing.T) {
	// let's test basic functionality
	d := initEasyDriver()
	_ = d.Disable()
	assert.False(t, d.IsEnabled())

	// let's make sure it stops first
	d = initEasyDriver()
	_ = d.Run()
	_ = d.Disable()
	assert.False(t, d.IsEnabled())
	assert.False(t, d.IsMoving())
}
