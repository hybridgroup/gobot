package ble

import (
	"testing"

	"gobot.io/x/gobot/gobottest"
)

func initTestDeviceInformationDriver() *DeviceInformationDriver {
	d := NewDeviceInformationDriver(NewClientAdaptor("D7:99:5A:26:EC:38"))
	return d
}

func TestDeviceInformationDriver(t *testing.T) {
	d := initTestDeviceInformationDriver()
	gobottest.Assert(t, d.Name(), "DeviceInformation")
}
