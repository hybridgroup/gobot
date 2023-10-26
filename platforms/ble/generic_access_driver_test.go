package ble

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*GenericAccessDriver)(nil)

func initTestGenericAccessDriver() *GenericAccessDriver {
	d := NewGenericAccessDriver(NewBleTestAdaptor())
	return d
}

func TestGenericAccessDriver(t *testing.T) {
	d := initTestGenericAccessDriver()
	assert.True(t, strings.HasPrefix(d.Name(), "GenericAccess"))
	d.SetName("NewName")
	assert.Equal(t, "NewName", d.Name())
}

func TestGenericAccessDriverStartAndHalt(t *testing.T) {
	d := initTestGenericAccessDriver()
	assert.NoError(t, d.Start())
	assert.NoError(t, d.Halt())
}

func TestGenericAccessDriverGetDeviceName(t *testing.T) {
	a := NewBleTestAdaptor()
	d := NewGenericAccessDriver(a)
	a.TestReadCharacteristic(func(cUUID string) ([]byte, error) {
		return []byte("TestDevice"), nil
	})

	assert.Equal(t, "TestDevice", d.GetDeviceName())
}

func TestGenericAccessDriverGetAppearance(t *testing.T) {
	a := NewBleTestAdaptor()
	d := NewGenericAccessDriver(a)
	a.TestReadCharacteristic(func(cUUID string) ([]byte, error) {
		return []byte{128, 0}, nil
	})

	assert.Equal(t, "Generic Computer", d.GetAppearance())
}
