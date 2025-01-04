package beaglebone

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2/platforms/adaptors"
)

func TestNewPocketBeagleAdaptor(t *testing.T) {
	// arrange & act
	a := NewPocketBeagleAdaptor()
	// assert
	assert.IsType(t, &PocketBeagleAdaptor{}, a)
	assert.True(t, strings.HasPrefix(a.Name(), "PocketBeagle"))
	assert.NotNil(t, a.sys)
	assert.NotNil(t, a.AnalogPinsAdaptor)
	assert.NotNil(t, a.pwmPinTranslate)
	assert.Equal(t, pocketBeaglePinMap, a.pinMap)
	assert.Equal(t, "/sys/class/leds/beaglebone:green:", a.usrLed)
	assert.True(t, a.sys.IsSysfsDigitalPinAccess())
}

func TestNewPocketBeagleAdaptorWithOption(t *testing.T) {
	// arrange & act
	a := NewPocketBeagleAdaptor(adaptors.WithGpioCdevAccess())
	// we have to mock the fs at this point to ensure the option can be applied on each test environment
	a.sys.UseMockFilesystem([]string{"/dev/gpiochip0"})
	// assert
	require.NoError(t, a.Connect())
	assert.True(t, a.sys.IsCdevDigitalPinAccess())
}
