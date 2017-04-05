package ble

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*DeviceInformationDriver)(nil)

func initTestDeviceInformationDriver() *DeviceInformationDriver {
	d := NewDeviceInformationDriver(newBleTestAdaptor())
	return d
}

func TestDeviceInformationDriver(t *testing.T) {
	d := initTestDeviceInformationDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "DeviceInformation"), true)
}

func TestDeviceInformationDriverGetModelNumber(t *testing.T) {
	a := newBleTestAdaptor()
	d := NewDeviceInformationDriver(a)
	a.TestReadCharacteristic(func(cUUID string) ([]byte, error) {
		return []byte("TestDevice"), nil
	})

	gobottest.Assert(t, d.GetModelNumber(), "TestDevice")
}

func TestDeviceInformationDriverGetFirmwareRevision(t *testing.T) {
	a := newBleTestAdaptor()
	d := NewDeviceInformationDriver(a)
	a.TestReadCharacteristic(func(cUUID string) ([]byte, error) {
		return []byte("TestDevice"), nil
	})

	gobottest.Assert(t, d.GetFirmwareRevision(), "TestDevice")
}

func TestDeviceInformationDriverGetHardwareRevision(t *testing.T) {
	a := newBleTestAdaptor()
	d := NewDeviceInformationDriver(a)
	a.TestReadCharacteristic(func(cUUID string) ([]byte, error) {
		return []byte("TestDevice"), nil
	})

	gobottest.Assert(t, d.GetHardwareRevision(), "TestDevice")
}

func TestDeviceInformationDriverGetManufacturerName(t *testing.T) {
	a := newBleTestAdaptor()
	d := NewDeviceInformationDriver(a)
	a.TestReadCharacteristic(func(cUUID string) ([]byte, error) {
		return []byte("TestDevice"), nil
	})

	gobottest.Assert(t, d.GetManufacturerName(), "TestDevice")
}

func TestDeviceInformationDriverGetPnPId(t *testing.T) {
	a := newBleTestAdaptor()
	d := NewDeviceInformationDriver(a)
	a.TestReadCharacteristic(func(cUUID string) ([]byte, error) {
		return []byte("TestDevice"), nil
	})

	gobottest.Assert(t, d.GetPnPId(), "TestDevice")
}
