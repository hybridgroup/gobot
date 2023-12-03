package gpio

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/aio"
)

var _ gobot.Driver = (*TM1638Driver)(nil)

func initTestTM1638Driver() *TM1638Driver {
	d, _ := initTestTM1638DriverWithStubbedAdaptor()
	return d
}

func initTestTM1638DriverWithStubbedAdaptor() (*TM1638Driver, *gpioTestAdaptor) {
	a := newGpioTestAdaptor()
	return NewTM1638Driver(a, "1", "2", "3"), a
}

func TestNewTM1638Driver(t *testing.T) {
	// arrange
	a := newGpioTestAdaptor()
	// act
	d := NewTM1638Driver(a, "10", "20", "30")
	// assert
	assert.IsType(t, &TM1638Driver{}, d)
	// assert: gpio.driver attributes
	require.NotNil(t, d.driver)
	assert.True(t, strings.HasPrefix(d.driverCfg.name, "TM1638"))
	assert.Equal(t, a, d.connection)
	assert.NotNil(t, d.afterStart)
	assert.NotNil(t, d.beforeHalt)
	assert.NotNil(t, d.Commander)
	assert.NotNil(t, d.mutex)
	// assert: driver specific attributes
	assert.NotNil(t, d.pinClock)
	assert.NotNil(t, d.pinData)
	assert.NotNil(t, d.pinStrobe)
	assert.NotNil(t, d.fonts)
}

func TestNewTM1638Driver_options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option, least one
	// option of this driver and one of another driver (which should lead to panic). Further tests for options can also
	// be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const (
		myName = "show rotation count"
	)
	panicFunc := func() {
		NewTM1638Driver(newGpioTestAdaptor(), "1", "2", "3", WithName("crazy"),
			aio.WithActuatorScaler(func(float64) int { return 0 }))
	}
	// act
	d := NewTM1638Driver(newGpioTestAdaptor(), "1", "2", "3", WithName(myName))
	// assert
	assert.Equal(t, myName, d.Name())
	assert.PanicsWithValue(t, "'scaler option for analog actuators' can not be applied on 'crazy'", panicFunc)
}

func TestTM1638Start(t *testing.T) {
	// arrange
	d := initTestTM1638Driver()
	// act & assert: tests also initialize()
	require.NoError(t, d.Start())
}

func TestTM1638FromStringToByteArray(t *testing.T) {
	d := initTestTM1638Driver()
	data := d.fromStringToByteArray("Hello World")
	assert.Equal(t, []byte{0x76, 0x7B, 0x30, 0x30, 0x5C, 0x00, 0x1D, 0x5C, 0x50, 0x30, 0x5E}, data)
}

func TestTM1638AddFonts(t *testing.T) {
	d := initTestTM1638Driver()
	d.AddFonts(map[string]byte{"µ": 0x1C, "ß": 0x7F})
	data := d.fromStringToByteArray("µß")
	assert.Equal(t, []byte{0x1C, 0x7F}, data)
}

func TestTM1638ClearFonts(t *testing.T) {
	d := initTestTM1638Driver()
	d.ClearFonts()
	data := d.fromStringToByteArray("Hello World")
	assert.Equal(t, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, data)
}
