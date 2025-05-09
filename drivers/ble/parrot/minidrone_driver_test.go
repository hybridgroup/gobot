package parrot

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble"
	"gobot.io/x/gobot/v2/drivers/ble/testutil"
)

var _ gobot.Driver = (*MinidroneDriver)(nil)

func initTestMinidroneDriver() *MinidroneDriver {
	d := NewMinidroneDriver(testutil.NewBleTestAdaptor())
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d
}

func TestNewMinidroneDriver(t *testing.T) {
	d := NewMinidroneDriver(testutil.NewBleTestAdaptor())
	assert.True(t, strings.HasPrefix(d.Name(), "Minidrone"))
	assert.NotNil(t, d.Eventer)
}

func TestNewMinidroneDriverWithName(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option.	Further
	// tests for options can also be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const newName = "new name"
	a := testutil.NewBleTestAdaptor()
	// act
	d := NewMinidroneDriver(a, ble.WithName(newName))
	// assert
	assert.Equal(t, newName, d.Name())
}

func TestMinidroneHalt(t *testing.T) {
	d := initTestMinidroneDriver()
	require.NoError(t, d.Halt())
}

func TestMinidroneTakeoff(t *testing.T) {
	d := initTestMinidroneDriver()
	require.NoError(t, d.TakeOff())
}

func TestMinidroneEmergency(t *testing.T) {
	d := initTestMinidroneDriver()
	require.NoError(t, d.Emergency())
}

func TestMinidroneTakePicture(t *testing.T) {
	d := initTestMinidroneDriver()
	require.NoError(t, d.TakePicture())
}

func TestMinidroneUp(t *testing.T) {
	d := initTestMinidroneDriver()
	require.NoError(t, d.Up(25))
}

func TestMinidroneUpTooFar(t *testing.T) {
	d := initTestMinidroneDriver()
	require.NoError(t, d.Up(125))
	require.NoError(t, d.Up(-50))
}

func TestMinidroneDown(t *testing.T) {
	d := initTestMinidroneDriver()
	require.NoError(t, d.Down(25))
}

func TestMinidroneForward(t *testing.T) {
	d := initTestMinidroneDriver()
	require.NoError(t, d.Forward(25))
}

func TestMinidroneBackward(t *testing.T) {
	d := initTestMinidroneDriver()
	require.NoError(t, d.Backward(25))
}

func TestMinidroneRight(t *testing.T) {
	d := initTestMinidroneDriver()
	require.NoError(t, d.Right(25))
}

func TestMinidroneLeft(t *testing.T) {
	d := initTestMinidroneDriver()
	require.NoError(t, d.Left(25))
}

func TestMinidroneClockwise(t *testing.T) {
	d := initTestMinidroneDriver()
	require.NoError(t, d.Clockwise(25))
}

func TestMinidroneCounterClockwise(t *testing.T) {
	d := initTestMinidroneDriver()
	require.NoError(t, d.CounterClockwise(25))
}

func TestMinidroneStop(t *testing.T) {
	d := initTestMinidroneDriver()
	require.NoError(t, d.Stop())
}

func TestMinidroneStartStopRecording(t *testing.T) {
	d := initTestMinidroneDriver()
	require.NoError(t, d.StartRecording())
	require.NoError(t, d.StopRecording())
}

func TestMinidroneHullProtectionOutdoor(t *testing.T) {
	d := initTestMinidroneDriver()
	require.NoError(t, d.HullProtection(true))
	require.NoError(t, d.Outdoor(true))
}

func TestMinidroneHullFlips(t *testing.T) {
	d := initTestMinidroneDriver()
	require.NoError(t, d.FrontFlip())
	require.NoError(t, d.BackFlip())
	require.NoError(t, d.RightFlip())
	require.NoError(t, d.LeftFlip())
}

func TestMinidroneLightControl(t *testing.T) {
	d := initTestMinidroneDriver()
	require.NoError(t, d.LightControl(0, LightBlinked, 25))
}

func TestMinidroneClawControl(t *testing.T) {
	d := initTestMinidroneDriver()
	require.NoError(t, d.ClawControl(0, ClawOpen))
}

func TestMinidroneGunControl(t *testing.T) {
	d := initTestMinidroneDriver()
	require.NoError(t, d.GunControl(0))
}

func TestMinidroneProcessFlightData(t *testing.T) {
	d := initTestMinidroneDriver()

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

	require.NoError(t, d.Stop())
}

func TestMinidroneValidatePitchWhenEqualOffset(t *testing.T) {
	assert.Equal(t, 100, ValidatePitch(32767.0, 32767.0))
}

func TestMinidroneValidatePitchWhenTiny(t *testing.T) {
	assert.Equal(t, 0, ValidatePitch(1.1, 32767.0))
}

func TestMinidroneValidatePitchWhenCentered(t *testing.T) {
	assert.Equal(t, 50, ValidatePitch(16383.5, 32767.0))
}
