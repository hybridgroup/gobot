package ble

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble/testutil"
)

var _ gobot.Driver = (*DeviceInformationDriver)(nil)

func TestNewDeviceInformationDriver(t *testing.T) {
	d := NewDeviceInformationDriver(testutil.NewBleTestAdaptor())
	assert.True(t, strings.HasPrefix(d.Name(), "DeviceInformation"))
	assert.NotNil(t, d.Eventer)
}

func TestDeviceInformationGetModelNumber(t *testing.T) {
	d := NewDeviceInformationDriver(testutil.NewBleTestAdaptor())
	assert.Equal(t, "2a24", d.GetModelNumber())
}

func TestDeviceInformationGetFirmwareRevision(t *testing.T) {
	d := NewDeviceInformationDriver(testutil.NewBleTestAdaptor())
	assert.Equal(t, "2a26", d.GetFirmwareRevision())
}

func TestDeviceInformationGetHardwareRevision(t *testing.T) {
	d := NewDeviceInformationDriver(testutil.NewBleTestAdaptor())
	assert.Equal(t, "2a27", d.GetHardwareRevision())
}

func TestDeviceInformationGetManufacturerName(t *testing.T) {
	d := NewDeviceInformationDriver(testutil.NewBleTestAdaptor())
	assert.Equal(t, "2a29", d.GetManufacturerName())
}

func TestDeviceInformationGetPnPId(t *testing.T) {
	d := NewDeviceInformationDriver(testutil.NewBleTestAdaptor())
	assert.Equal(t, "2a50", d.GetPnPId())
}
