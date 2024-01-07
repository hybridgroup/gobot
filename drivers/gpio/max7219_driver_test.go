package gpio

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/aio"
)

var _ gobot.Driver = (*MAX7219Driver)(nil)

func initTestMAX7219Driver() *MAX7219Driver {
	d, _ := initTestMAX7219DriverWithStubbedAdaptor()
	return d
}

func initTestMAX7219DriverWithStubbedAdaptor() (*MAX7219Driver, *gpioTestAdaptor) {
	a := newGpioTestAdaptor()
	return NewMAX7219Driver(a, "1", "2", "3", 1), a
}

func TestNewMAX7219Driver(t *testing.T) {
	// arrange
	a := newGpioTestAdaptor()
	// act
	d := NewMAX7219Driver(a, "1", "2", "3", 4)
	// assert
	assert.IsType(t, &MAX7219Driver{}, d)
	// assert: gpio.driver attributes
	require.NotNil(t, d.driver)
	assert.True(t, strings.HasPrefix(d.driverCfg.name, "MAX7219"))
	assert.Equal(t, a, d.connection)
	assert.NotNil(t, d.afterStart)
	assert.NotNil(t, d.beforeHalt)
	assert.NotNil(t, d.Commander)
	assert.NotNil(t, d.mutex)
	// assert: driver specific attributes
	assert.NotNil(t, d.pinClock)
	assert.NotNil(t, d.pinData)
	assert.NotNil(t, d.pinCS)
	assert.Equal(t, uint(4), d.count)
}

func TestNewMAX7219Driver_options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option, least one
	// option of this driver and one of another driver (which should lead to panic). Further tests for options can also
	// be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const (
		myName = "light chain 5"
	)
	panicFunc := func() {
		NewMAX7219Driver(newGpioTestAdaptor(), "1", "2", "3", 4, WithName("crazy"),
			aio.WithActuatorScaler(func(float64) int { return 0 }))
	}
	// act
	d := NewMAX7219Driver(newGpioTestAdaptor(), "1", "2", "3", 4, WithName(myName))
	// assert
	assert.Equal(t, myName, d.Name())
	assert.PanicsWithValue(t, "'scaler option for analog actuators' can not be applied on 'crazy'", panicFunc)
}

func TestMAX7219Start(t *testing.T) {
	// arrange
	d := initTestMAX7219Driver()
	// act & assert: tests also initialize()
	require.NoError(t, d.Start())
}
