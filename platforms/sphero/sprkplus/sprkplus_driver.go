package sprkplus

import (
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/ble"
	"gobot.io/x/gobot/v2/platforms/sphero/ollie"
)

// Driver represents a Sphero SPRK+
type SPRKPlusDriver struct {
	*ollie.Driver
}

// NewDriver creates a Driver for a Sphero SPRK+
func NewDriver(a ble.BLEConnector) *SPRKPlusDriver {
	d := ollie.NewDriver(a)
	d.SetName(gobot.DefaultName("SPRKPlus"))

	return &SPRKPlusDriver{
		Driver: d,
	}
}
