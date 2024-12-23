package tinkerboard2

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
	assert.IsType(t, &Tinkerboard2Adaptor{}, a)
	assert.True(t, strings.HasPrefix(a.Name(), "Tinker Board 2"))
	assert.NotNil(t, a.AnalogPinsAdaptor)
	assert.NotNil(t, a.DigitalPinsAdaptor)
	assert.NotNil(t, a.PWMPinsAdaptor)
	assert.NotNil(t, a.I2cBusAdaptor)
	assert.NotNil(t, a.SpiBusAdaptor)
	assert.True(t, a.sys.IsGpiodDigitalPinAccess())
	// act & assert
	a.SetName("NewName")
	assert.Equal(t, "NewName", a.Name())
}

func TestNewAdaptorWithOption(t *testing.T) {
	// arrange & act
	a := NewAdaptor(adaptors.WithGpiosActiveLow("1"), adaptors.WithSysfsAccess())
	// assert
	require.NoError(t, a.Connect())
	assert.True(t, a.sys.IsSysfsDigitalPinAccess())
}
