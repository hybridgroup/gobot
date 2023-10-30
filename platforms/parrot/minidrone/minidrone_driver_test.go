package minidrone

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*Driver)(nil)

func initTestMinidroneDriver() *Driver {
	d := NewDriver(NewBleTestAdaptor())
	return d
}

func TestMinidroneDriver(t *testing.T) {
	d := initTestMinidroneDriver()
	assert.True(t, strings.HasPrefix(d.Name(), "Minidrone"))
	d.SetName("NewName")
	assert.Equal(t, "NewName", d.Name())
}

func TestMinidroneDriverStartAndHalt(t *testing.T) {
	d := initTestMinidroneDriver()
	assert.NoError(t, d.Start())
	assert.NoError(t, d.Halt())
}

func TestMinidroneTakeoff(t *testing.T) {
	d := initTestMinidroneDriver()
	assert.NoError(t, d.Start())
	assert.NoError(t, d.TakeOff())
}

func TestMinidroneEmergency(t *testing.T) {
	d := initTestMinidroneDriver()
	assert.NoError(t, d.Start())
	assert.NoError(t, d.Emergency())
}

func TestMinidroneTakePicture(t *testing.T) {
	d := initTestMinidroneDriver()
	assert.NoError(t, d.Start())
	assert.NoError(t, d.TakePicture())
}

func TestMinidroneUp(t *testing.T) {
	d := initTestMinidroneDriver()
	assert.NoError(t, d.Start())
	assert.NoError(t, d.Up(25))
}

func TestMinidroneUpTooFar(t *testing.T) {
	d := initTestMinidroneDriver()
	assert.NoError(t, d.Start())
	assert.NoError(t, d.Up(125))
	assert.NoError(t, d.Up(-50))
}

func TestMinidroneDown(t *testing.T) {
	d := initTestMinidroneDriver()
	assert.NoError(t, d.Start())
	assert.NoError(t, d.Down(25))
}

func TestMinidroneForward(t *testing.T) {
	d := initTestMinidroneDriver()
	assert.NoError(t, d.Start())
	assert.NoError(t, d.Forward(25))
}

func TestMinidroneBackward(t *testing.T) {
	d := initTestMinidroneDriver()
	assert.NoError(t, d.Start())
	assert.NoError(t, d.Backward(25))
}

func TestMinidroneRight(t *testing.T) {
	d := initTestMinidroneDriver()
	assert.NoError(t, d.Start())
	assert.NoError(t, d.Right(25))
}

func TestMinidroneLeft(t *testing.T) {
	d := initTestMinidroneDriver()
	assert.NoError(t, d.Start())
	assert.NoError(t, d.Left(25))
}

func TestMinidroneClockwise(t *testing.T) {
	d := initTestMinidroneDriver()
	assert.NoError(t, d.Start())
	assert.NoError(t, d.Clockwise(25))
}

func TestMinidroneCounterClockwise(t *testing.T) {
	d := initTestMinidroneDriver()
	assert.NoError(t, d.Start())
	assert.NoError(t, d.CounterClockwise(25))
}

func TestMinidroneStop(t *testing.T) {
	d := initTestMinidroneDriver()
	assert.NoError(t, d.Start())
	assert.NoError(t, d.Stop())
}

func TestMinidroneStartStopRecording(t *testing.T) {
	d := initTestMinidroneDriver()
	assert.NoError(t, d.Start())
	assert.NoError(t, d.StartRecording())
	assert.NoError(t, d.StopRecording())
}

func TestMinidroneHullProtectionOutdoor(t *testing.T) {
	d := initTestMinidroneDriver()
	assert.NoError(t, d.Start())
	assert.NoError(t, d.HullProtection(true))
	assert.NoError(t, d.Outdoor(true))
}

func TestMinidroneHullFlips(t *testing.T) {
	d := initTestMinidroneDriver()
	assert.NoError(t, d.Start())
	assert.NoError(t, d.FrontFlip())
	assert.NoError(t, d.BackFlip())
	assert.NoError(t, d.RightFlip())
	assert.NoError(t, d.LeftFlip())
}

func TestMinidroneLightControl(t *testing.T) {
	d := initTestMinidroneDriver()
	assert.NoError(t, d.Start())
	assert.NoError(t, d.LightControl(0, LightBlinked, 25))
}

func TestMinidroneClawControl(t *testing.T) {
	d := initTestMinidroneDriver()
	assert.NoError(t, d.Start())
	assert.NoError(t, d.ClawControl(0, ClawOpen))
}

func TestMinidroneGunControl(t *testing.T) {
	d := initTestMinidroneDriver()
	assert.NoError(t, d.Start())
	assert.NoError(t, d.GunControl(0))
}

func TestMinidroneProcessFlightData(t *testing.T) {
	d := initTestMinidroneDriver()
	assert.NoError(t, d.Start())

	d.processFlightStatus([]byte{0x00, 0x00, 0x00})
	d.processFlightStatus([]byte{0x00, 0x00, 0x00, 0x00, 0x00})
	d.processFlightStatus([]byte{0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00})
	assert.False(t, d.flying)
	d.processFlightStatus([]byte{0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01})
	assert.False(t, d.flying)
	d.processFlightStatus([]byte{0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x02})
	assert.True(t, d.flying)
	d.processFlightStatus([]byte{0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x03})
	assert.True(t, d.flying)
	d.processFlightStatus([]byte{0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x04})
	d.processFlightStatus([]byte{0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x05})
	d.processFlightStatus([]byte{0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x06})

	assert.NoError(t, d.Stop())
}
