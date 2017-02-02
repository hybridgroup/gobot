package ble

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*DeviceInformationDriver)(nil)

func initTestDeviceInformationDriver() *DeviceInformationDriver {
	d := NewDeviceInformationDriver(NewClientAdaptor("D7:99:5A:26:EC:38"))
	return d
}

func TestDeviceInformationDriver(t *testing.T) {
	d := initTestDeviceInformationDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "DeviceInformation"), true)
}
