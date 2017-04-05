package ble

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*GenericAccessDriver)(nil)

func initTestGenericAccessDriver() *GenericAccessDriver {
	d := NewGenericAccessDriver(NewBleTestAdaptor())
	return d
}

func TestGenericAccessDriver(t *testing.T) {
	d := initTestGenericAccessDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "GenericAccess"), true)
	d.SetName("NewName")
	gobottest.Assert(t, d.Name(), "NewName")
}

func TestGenericAccessDriverStartAndHalt(t *testing.T) {
	d := initTestGenericAccessDriver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Halt(), nil)
}

func TestGenericAccessDriverGetDeviceName(t *testing.T) {
	a := NewBleTestAdaptor()
	d := NewGenericAccessDriver(a)
	a.TestReadCharacteristic(func(cUUID string) ([]byte, error) {
		return []byte("TestDevice"), nil
	})

	gobottest.Assert(t, d.GetDeviceName(), "TestDevice")
}

func TestGenericAccessDriverGetAppearance(t *testing.T) {
	a := NewBleTestAdaptor()
	d := NewGenericAccessDriver(a)
	a.TestReadCharacteristic(func(cUUID string) ([]byte, error) {
		return []byte{128, 0}, nil
	})

	gobottest.Assert(t, d.GetAppearance(), "Generic Computer")
}
