package ble

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble/testutil"
)

var _ gobot.Driver = (*DeviceInformationDriver)(nil)

func TestNewDeviceInformationDriver(t *testing.T) {
	// arrange
	d := NewDeviceInformationDriver(testutil.NewBleTestAdaptor())
	// act & assert
	assert.True(t, strings.HasPrefix(d.Name(), "DeviceInformation"))
	assert.NotNil(t, d.Eventer)
}

func TestNewDeviceInformationDriverWithName(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option.	Further
	// tests for options can also be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const newName = "new name"
	a := testutil.NewBleTestAdaptor()
	// act
	d := NewDeviceInformationDriver(a, WithName(newName))
	// assert
	assert.Equal(t, newName, d.Name())
}

func TestDeviceInformationGetModelNumber(t *testing.T) {
	// arrange
	d := NewDeviceInformationDriver(testutil.NewBleTestAdaptor())
	// act
	modelNo, err := d.GetModelNumber()
	// assert
	require.NoError(t, err)
	assert.Equal(t, "2a24", modelNo)
}

func TestDeviceInformationGetFirmwareRevision(t *testing.T) {
	// arrange
	d := NewDeviceInformationDriver(testutil.NewBleTestAdaptor())
	// act
	fwRev, err := d.GetFirmwareRevision()
	// assert
	require.NoError(t, err)
	assert.Equal(t, "2a26", fwRev)
}

func TestDeviceInformationGetHardwareRevision(t *testing.T) {
	// arrange
	d := NewDeviceInformationDriver(testutil.NewBleTestAdaptor())
	// act
	hwRev, err := d.GetHardwareRevision()
	// assert
	require.NoError(t, err)
	assert.Equal(t, "2a27", hwRev)
}

func TestDeviceInformationGetManufacturerName(t *testing.T) {
	// arrange
	d := NewDeviceInformationDriver(testutil.NewBleTestAdaptor())
	// act
	manuName, err := d.GetManufacturerName()
	// assert
	require.NoError(t, err)
	assert.Equal(t, "2a29", manuName)
}

func TestDeviceInformationGetPnPId(t *testing.T) {
	// arrange
	d := NewDeviceInformationDriver(testutil.NewBleTestAdaptor())
	// act
	pid, err := d.GetPnPId()
	// assert
	require.NoError(t, err)
	assert.Equal(t, "2a50", pid)
}
