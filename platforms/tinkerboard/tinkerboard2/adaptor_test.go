package tinkerboard2

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
	// act & assert
	a.SetName("NewName")
	assert.Equal(t, "NewName", a.Name())
}
