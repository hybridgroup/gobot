package gpio

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*driver)(nil)

func initTestDriverWithStubbedAdaptor() (*driver, *gpioTestAdaptor) {
	a := newGpioTestAdaptor()
	d := newDriver(a, "GPIO_BASIC")
	return d, a
}

func initTestDriver() *driver {
	d, _ := initTestDriverWithStubbedAdaptor()
	return d
}

func TestNewDriver(t *testing.T) {
	// arrange
	a := newGpioTestAdaptor()
	// act
	d := newDriver(a, "GPIO_BASIC")
	// assert
	assert.IsType(t, &driver{}, d)
	assert.Contains(t, d.driverCfg.name, "GPIO_BASIC")
	assert.Equal(t, a, d.connection)
	require.NoError(t, d.afterStart())
	require.NoError(t, d.beforeHalt())
	assert.NotNil(t, d.Commander)
	assert.NotNil(t, d.mutex)
}

func Test_applyWithName(t *testing.T) {
	// arrange
	const name = "mybot"
	cfg := configuration{name: "oldname"}
	// act
	WithName(name).apply(&cfg)
	// assert
	assert.Equal(t, name, cfg.name)
}

func Test_applywithPin(t *testing.T) {
	// arrange
	const pin = "36"
	cfg := configuration{pin: "oldpin"}
	// act
	withPin(pin).apply(&cfg)
	// assert
	assert.Equal(t, pin, cfg.pin)
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
	require.EqualError(t, d.Start(), "after start error")
}

func TestHalt(t *testing.T) {
	// arrange
	d := initTestDriver()
	// act, assert
	require.NoError(t, d.Halt())
	// arrange after start function
	d.beforeHalt = func() error { return fmt.Errorf("before halt error") }
	// act, assert
	require.EqualError(t, d.Halt(), "before halt error")
}
