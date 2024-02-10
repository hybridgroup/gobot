package sphero

import (
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble"
	"gobot.io/x/gobot/v2/drivers/common/sphero"
)

// SPRKPlusDriver represents a Sphero SPRK+
type SPRKPlusDriver struct {
	*OllieDriver
}

// NewSPRKPlusDriver creates a driver for a Sphero SPRK+
func NewSPRKPlusDriver(a gobot.BLEConnector, opts ...ble.OptionApplier) *SPRKPlusDriver {
	return &SPRKPlusDriver{OllieDriver: newOllieBaseDriver(a, "SPRKPlus", sprkplusDefaultCollisionConfig(), opts...)}
}

// sprkplusDefaultCollisionConfig returns a CollisionConfig with sensible collision defaults
func sprkplusDefaultCollisionConfig() sphero.CollisionConfig {
	return sphero.CollisionConfig{
		Method: 0x01,
		Xt:     0x20,
		Yt:     0x20,
		Xs:     0x20,
		Ys:     0x20,
		Dead:   0x01,
	}
}
