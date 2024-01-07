package gpio

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/aio"
)

var _ gobot.Driver = (*AIP1640Driver)(nil)

func initTestAIP1640Driver() *AIP1640Driver {
	d, _ := initTestAIP1640DriverWithStubbedAdaptor()
	return d
}

func initTestAIP1640DriverWithStubbedAdaptor() (*AIP1640Driver, *gpioTestAdaptor) {
	a := newGpioTestAdaptor()
	return NewAIP1640Driver(a, "1", "2"), a
}

func TestNewAIP1640Driver(t *testing.T) {
	// arrange
	a := newGpioTestAdaptor()
	// act
	d := NewAIP1640Driver(a, "1", "2")
	// assert
	assert.IsType(t, &AIP1640Driver{}, d)
	// assert: gpio.driver attributes
	require.NotNil(t, d.driver)
	assert.True(t, strings.HasPrefix(d.driverCfg.name, "AIP1640"))
	assert.Equal(t, a, d.connection)
	assert.NotNil(t, d.afterStart)
	assert.NotNil(t, d.beforeHalt)
	assert.NotNil(t, d.Commander)
	assert.NotNil(t, d.mutex)
	// assert: driver specific attributes
	assert.NotNil(t, d.pinClock)
	assert.NotNil(t, d.pinData)
	assert.Equal(t, uint8(7), d.intensity)
}

func TestNewAIP1640Driver_options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option, least one
	// option of this driver and one of another driver (which should lead to panic). Further tests for options can also
	// be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const (
		myName     = "count up"
		cycReadDur = 30 * time.Millisecond
	)
	panicFunc := func() {
		NewAIP1640Driver(newGpioTestAdaptor(), "1", "2", WithName("crazy"),
			aio.WithActuatorScaler(func(float64) int { return 0 }))
	}
	// act
	d := NewAIP1640Driver(newGpioTestAdaptor(), "1", "2", WithName(myName))
	// assert
	assert.Equal(t, myName, d.Name())
	assert.PanicsWithValue(t, "'scaler option for analog actuators' can not be applied on 'crazy'", panicFunc)
}

func TestAIP1640Start(t *testing.T) {
	d := initTestAIP1640Driver()
	require.NoError(t, d.Start())
}

func TestAIP1640DrawPixel(t *testing.T) {
	d := initTestAIP1640Driver()
	d.DrawPixel(2, 3, true)
	d.DrawPixel(0, 3, true)
	assert.Equal(t, uint8(5), d.buffer[7-3])
}

func TestAIP1640DrawRow(t *testing.T) {
	d := initTestAIP1640Driver()
	d.DrawRow(4, 0x3C)
	assert.Equal(t, uint8(0x3C), d.buffer[7-4])
}

func TestAIP1640DrawMatrix(t *testing.T) {
	d := initTestAIP1640Driver()
	drawing := [8]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF}
	d.DrawMatrix(drawing)
	assert.Equal(t, [8]byte{0xEF, 0xCD, 0xAB, 0x89, 0x67, 0x45, 0x23, 0x01}, d.buffer)
}

func TestAIP1640Clear(t *testing.T) {
	d := initTestAIP1640Driver()
	drawing := [8]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF}
	d.DrawMatrix(drawing)
	assert.Equal(t, [8]byte{0xEF, 0xCD, 0xAB, 0x89, 0x67, 0x45, 0x23, 0x01}, d.buffer)
	d.Clear()
	assert.Equal(t, [8]byte{}, d.buffer)
}

func TestAIP1640SetIntensity(t *testing.T) {
	d := initTestAIP1640Driver()
	d.SetIntensity(3)
	assert.Equal(t, uint8(3), d.intensity)
}

func TestAIP1640SetIntensityHigherThan7(t *testing.T) {
	d := initTestAIP1640Driver()
	d.SetIntensity(19)
	assert.Equal(t, uint8(7), d.intensity)
}
