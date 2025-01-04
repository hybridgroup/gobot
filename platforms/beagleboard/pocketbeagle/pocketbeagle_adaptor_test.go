package pocketbeagle

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2/platforms/adaptors"
)

func TestNewAdaptor(t *testing.T) {
	// arrange & act
	a := NewAdaptor()
	// assert
	assert.IsType(t, &PocketBeagleAdaptor{}, a)
	assert.True(t, strings.HasPrefix(a.Name(), "PocketBeagle"))
	assert.NotNil(t, a.sys)
	assert.NotNil(t, a.AnalogPinsAdaptor)
	assert.NotNil(t, a.DigitalPinsAdaptor)
	assert.NotNil(t, a.PWMPinsAdaptor)
	assert.NotNil(t, a.I2cBusAdaptor)
	assert.NotNil(t, a.SpiBusAdaptor)
	assert.True(t, a.sys.IsCdevDigitalPinAccess())
}

func TestNewAdaptorWithOption(t *testing.T) {
	// arrange & act
	a := NewAdaptor(adaptors.WithGpioSysfsAccess())
	// assert
	require.NoError(t, a.Connect())
	assert.True(t, a.sys.IsSysfsDigitalPinAccess())
}

func TestDigitalIO(t *testing.T) {
	// some basic tests, further tests are done in "digitalpinsadaptor.go"
	// arrange
	a := NewAdaptor()
	require.NoError(t, a.Connect())
	dpa := a.sys.UseMockDigitalPinAccess()
	require.True(t, a.sys.IsCdevDigitalPinAccess())
	// act & assert write
	err := a.DigitalWrite("P1_02", 1)
	require.NoError(t, err)
	assert.Equal(t, []int{1}, dpa.Written("gpiochip2", "23"))
	// arrange, act & assert read
	dpa.UseValue("gpiochip3", "19", 2)
	i, err := a.DigitalRead("P2_34")
	require.NoError(t, err)
	assert.Equal(t, 2, i)
	// act and assert unknown pin
	require.ErrorContains(t, a.DigitalWrite("99", 1), "'99' is not a valid id for a digital pin")
	// act and assert finalize
	require.NoError(t, a.Finalize())
	assert.Equal(t, 0, dpa.Exported("gpiochip2", "23"))
	assert.Equal(t, 0, dpa.Exported("gpiochip3", "19"))
}
