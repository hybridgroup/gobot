package bb8

import (
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/ble"
	"gobot.io/x/gobot/v2/platforms/sphero/ollie"
)

// Driver represents a Sphero BB-8
type BB8Driver struct {
	*ollie.Driver
}

// NewDriver creates a Driver for a Sphero BB-8
func NewDriver(a ble.BLEConnector) *BB8Driver {
	d := ollie.NewDriver(a)
	d.SetName(gobot.DefaultName("BB8"))

	return &BB8Driver{
		Driver: d,
	}
}
