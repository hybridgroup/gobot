package ble

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*GenericAccessDriver)(nil)

func initTestGenericAccessDriver() *GenericAccessDriver {
	d := NewGenericAccessDriver(newBleTestAdaptor())
	return d
}

func TestGenericAccessDriver(t *testing.T) {
	d := initTestGenericAccessDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "GenericAccess"), true)
}

func TestGenericAccessDriverGetDeviceName(t *testing.T) {
	a := newBleTestAdaptor()
	d := NewGenericAccessDriver(a)
	a.TestReadCharacteristic(func(cUUID string) ([]byte, error) {
		return []byte("TestDevice"), nil
	})

	gobottest.Assert(t, d.GetDeviceName(), "TestDevice")
}

func TestGenericAccessDriverGetAppearance(t *testing.T) {
	a := newBleTestAdaptor()
	d := NewGenericAccessDriver(a)
	a.TestReadCharacteristic(func(cUUID string) ([]byte, error) {
		return []byte{128, 0}, nil
	})

	gobottest.Assert(t, d.GetAppearance(), "Generic Computer")
}
