package sphero

import (
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble"
	"gobot.io/x/gobot/v2/drivers/common/spherocommon"
)

// BB8Driver represents a Sphero BB-8
type BB8Driver struct {
	*OllieDriver
}

// NewBB8Driver creates a driver for a Sphero BB-8
func NewBB8Driver(a gobot.BLEConnector, opts ...ble.OptionApplier) *BB8Driver {
	return &BB8Driver{OllieDriver: newOllieBaseDriver(a, "BB8", bb8DefaultCollisionConfig(), opts...)}
}

// bb8DefaultCollisionConfig returns a CollisionConfig with sensible collision defaults
func bb8DefaultCollisionConfig() spherocommon.CollisionConfig {
	return spherocommon.CollisionConfig{
		Method: 0x01,
		Xt:     0x20,
		Yt:     0x20,
		Xs:     0x20,
		Ys:     0x20,
		Dead:   0x01,
	}
}
