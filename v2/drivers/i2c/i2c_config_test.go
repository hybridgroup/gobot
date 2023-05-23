package i2c

import (
	"testing"

	"gobot.io/x/gobot/v2/gobottest"
)

func TestNewConfig(t *testing.T) {
	// arrange, act
	ci := NewConfig()
	// assert
	c, ok := ci.(*i2cConfig)
	if !ok {
		t.Errorf("NewConfig() should have returned a *i2cConfig")
	}
	gobottest.Assert(t, c.bus, BusNotInitialized)
	gobottest.Assert(t, c.address, AddressNotInitialized)
}

func TestWithBus(t *testing.T) {
	// arrange
	c := NewConfig()
	// act
	c.SetBus(0x23)
	// assert
	gobottest.Assert(t, c.(*i2cConfig).bus, 0x23)
}

func TestWithAddress(t *testing.T) {
	// arrange
	c := NewConfig()
	// act
	c.SetAddress(0x24)
	// assert
	gobottest.Assert(t, c.(*i2cConfig).address, 0x24)
}

func TestGetBusOrDefaultWithBusOption(t *testing.T) {
	var tests = map[string]struct {
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
			gobottest.Assert(t, got, tc.want)
		})
	}
}

func TestGetAddressOrDefaultWithAddressOption(t *testing.T) {
	var tests = map[string]struct {
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
			gobottest.Assert(t, got, tc.want)
		})
	}
}
