package onewire

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func initTestDriver() *driver {
	d, _ := initDriverWithStubbedAdaptor()
	return d
}

func initDriverWithStubbedAdaptor() (*driver, *oneWireAdaptorMock) {
	a := newOneWireTestAdaptor()
	return newDriver(a, "name", 28, 9876), a
}

func Test_newDriver(t *testing.T) {
	// arrange
	const (
		familyCode   = 99
		serialNumber = 1234567890
	)
	// act
	d := newDriver(newOneWireTestAdaptor(), "name", familyCode, serialNumber)
	// assert
	assert.IsType(t, &driver{}, d)
	assert.NotNil(t, d.driverCfg)
	assert.NotNil(t, d.connector)
	assert.NotNil(t, d.afterStart)
	assert.NotNil(t, d.beforeHalt)
	assert.NotNil(t, d.Commander)
	assert.NotNil(t, d.mutex)
}

func TestConnection(t *testing.T) {
	// arrange
	d := initTestDriver()
	require.NoError(t, d.Start())
	// act, assert
	assert.NotNil(t, d.Connection())
}

func TestStart(t *testing.T) {
	// arrange
	d := initTestDriver()
	// act, assert
	require.NoError(t, d.Start())
}

func TestStartConnectError(t *testing.T) {
	// arrange
	d, c := initDriverWithStubbedAdaptor()
	c.retErr = true
	// act, assert
	require.ErrorContains(t, d.Start(), "GetOneWireConnection error")
}

func TestHalt(t *testing.T) {
	// arrange
	d := initTestDriver()
	// act, assert
	require.NoError(t, d.Halt())
}

func TestWithName(t *testing.T) {
	// This is a general test, that options are applied by using the WithName() option.
	// All other configuration options can also be tested by With..(val).apply(cfg).
	// arrange
	const newName = "new name"
	a := newOneWireTestAdaptor()
	// act
	d := newDriver(a, "name", 28, 9876, WithName(newName))
	// assert
	assert.Equal(t, newName, d.driverCfg.name)
	assert.Equal(t, newName, d.Name())
}
