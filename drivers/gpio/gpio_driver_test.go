package gpio

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*Driver)(nil)

func initTestDriverWithStubbedAdaptor() (*Driver, *gpioTestAdaptor) {
	a := newGpioTestAdaptor()
	d := NewDriver(a, "GPIO_BASIC")
	return d, a
}

func initTestDriver() *Driver {
	d, _ := initTestDriverWithStubbedAdaptor()
	return d
}

func TestNewDriver(t *testing.T) {
	// arrange
	a := newGpioTestAdaptor()
	// act
	d := NewDriver(a, "GPIO_BASIC")
	// assert
	assert.IsType(t, &Driver{}, d)
	assert.Contains(t, d.name, "GPIO_BASIC")
	assert.Equal(t, a, d.connection)
	require.NoError(t, d.afterStart())
	require.NoError(t, d.beforeHalt())
	assert.NotNil(t, d.Commander)
	assert.NotNil(t, d.mutex)
}

func TestSetName(t *testing.T) {
	// arrange
	d := initTestDriver()
	// act
	d.SetName("TESTME")
	// assert
	assert.Equal(t, "TESTME", d.Name())
}

func TestConnection(t *testing.T) {
	// arrange
	d, a := initTestDriverWithStubbedAdaptor()
	// act, assert
	assert.Equal(t, a, d.Connection())
}

func TestStart(t *testing.T) {
	// arrange
	d := initTestDriver()
	// act, assert
	require.NoError(t, d.Start())
	// arrange after start function
	d.afterStart = func() error { return fmt.Errorf("after start error") }
	// act, assert
	require.ErrorContains(t, d.Start(), "after start error")
}

func TestHalt(t *testing.T) {
	// arrange
	d := initTestDriver()
	// act, assert
	require.NoError(t, d.Halt())
	// arrange after start function
	d.beforeHalt = func() error { return fmt.Errorf("before halt error") }
	// act, assert
	require.ErrorContains(t, d.Halt(), "before halt error")
}
