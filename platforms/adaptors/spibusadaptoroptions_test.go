package adaptors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2/system"
)

func TestNewSpiBusAdaptorWithSpiGpioAccess(t *testing.T) {
	// arrange
	const (
		sclkPin           = "1"
		ncsPin            = "2"
		sdoPin            = "3"
		sdiPin            = "4"
		sclkPinTranslated = "1"
		ncsPinTranslated  = "2"
		sdoPinTranslated  = "3"
		sdiPinTranslated  = "4"
	)
	sys := system.NewAccesser()
	dpa := sys.UseMockDigitalPinAccess()
	a := NewSpiBusAdaptor(sys, nil, 1, 2, 3, 4, 5, dpa, WithSpiGpioAccess(sclkPin, ncsPin, sdoPin, sdiPin))
	// act
	require.NoError(t, a.Connect())
	bus, err := a.sys.NewSpiDevice(0, 0, 0, 0, 1111)
	// assert
	require.NoError(t, err)
	assert.NotNil(t, bus)
	assert.True(t, a.sys.HasSpiGpioAccess())
	assert.Equal(t, 1, dpa.AppliedOptions("", sclkPinTranslated))
	assert.Equal(t, 1, dpa.AppliedOptions("", ncsPinTranslated))
	assert.Equal(t, 1, dpa.AppliedOptions("", sdoPinTranslated))
	assert.Equal(t, 0, dpa.AppliedOptions("", sdiPinTranslated)) // already input, so no option applied
}
