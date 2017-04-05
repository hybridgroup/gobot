package ble

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*DeviceInformationDriver)(nil)

func initTestDeviceInformationDriver() *DeviceInformationDriver {
	d := NewDeviceInformationDriver(NewBleTestAdaptor())
	return d
}

func TestDeviceInformationDriver(t *testing.T) {
	d := initTestDeviceInformationDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "DeviceInformation"), true)
	d.SetName("NewName")
	gobottest.Assert(t, d.Name(), "NewName")
}

func TestDeviceInformationDriverStartAndHalt(t *testing.T) {
	d := initTestDeviceInformationDriver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Halt(), nil)
}

func TestDeviceInformationDriverGetModelNumber(t *testing.T) {
	a := NewBleTestAdaptor()
	d := NewDeviceInformationDriver(a)
	a.TestReadCharacteristic(func(cUUID string) ([]byte, error) {
		return []byte("TestDevice"), nil
	})

	gobottest.Assert(t, d.GetModelNumber(), "TestDevice")
}

func TestDeviceInformationDriverGetFirmwareRevision(t *testing.T) {
	a := NewBleTestAdaptor()
	d := NewDeviceInformationDriver(a)
	a.TestReadCharacteristic(func(cUUID string) ([]byte, error) {
		return []byte("TestDevice"), nil
	})

	gobottest.Assert(t, d.GetFirmwareRevision(), "TestDevice")
}

func TestDeviceInformationDriverGetHardwareRevision(t *testing.T) {
	a := NewBleTestAdaptor()
	d := NewDeviceInformationDriver(a)
	a.TestReadCharacteristic(func(cUUID string) ([]byte, error) {
		return []byte("TestDevice"), nil
	})

	gobottest.Assert(t, d.GetHardwareRevision(), "TestDevice")
}

func TestDeviceInformationDriverGetManufacturerName(t *testing.T) {
	a := NewBleTestAdaptor()
	d := NewDeviceInformationDriver(a)
	a.TestReadCharacteristic(func(cUUID string) ([]byte, error) {
		return []byte("TestDevice"), nil
	})

	gobottest.Assert(t, d.GetManufacturerName(), "TestDevice")
}

func TestDeviceInformationDriverGetPnPId(t *testing.T) {
	a := NewBleTestAdaptor()
	d := NewDeviceInformationDriver(a)
	a.TestReadCharacteristic(func(cUUID string) ([]byte, error) {
		return []byte("TestDevice"), nil
	})

	gobottest.Assert(t, d.GetPnPId(), "TestDevice")
}
