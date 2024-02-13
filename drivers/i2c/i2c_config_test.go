//nolint:forcetypeassert // ok here
package i2c

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	// arrange, act
	ci := NewConfig()
	// assert
	c, ok := ci.(*i2cConfig)
	if !ok {
		require.Fail(t, "NewConfig() should have returned a *i2cConfig")
	}
	assert.Equal(t, BusNotInitialized, c.bus)
	assert.Equal(t, AddressNotInitialized, c.address)
}

func TestWithBus(t *testing.T) {
	// arrange
	c := NewConfig()
	// act
	c.SetBus(0x23)
	// assert
	assert.Equal(t, 0x23, c.(*i2cConfig).bus)
}

func TestWithAddress(t *testing.T) {
	// arrange
	c := NewConfig()
	// act
	c.SetAddress(0x24)
	// assert
	assert.Equal(t, 0x24, c.(*i2cConfig).address)
}

func TestGetBusOrDefaultWithBusOption(t *testing.T) {
	tests := map[string]struct {
		init int
		bus  int
		want int
	}{
		"not_initialized": {init: -1, bus: 0x25, want: 0x25},
		"initialized":     {init: 0x26, bus: 0x27, want: 0x26},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			c := NewConfig()
			// act
			WithBus(tc.init)(c)
			got := c.GetBusOrDefault(tc.bus)
			// assert
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGetAddressOrDefaultWithAddressOption(t *testing.T) {
	tests := map[string]struct {
		init    int
		address int
		want    int
	}{
		"not_initialized": {init: -1, address: 0x28, want: 0x28},
		"initialized":     {init: 0x29, address: 0x2A, want: 0x29},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			c := NewConfig()
			// act
			WithAddress(tc.init)(c)
			got := c.GetAddressOrDefault(tc.address)
			// assert
			assert.Equal(t, tc.want, got)
		})
	}
}
