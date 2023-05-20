package bb8

import "gobot.io/x/gobot/v2/platforms/sphero"

// DefaultCollisionConfig returns a CollisionConfig with sensible collision defaults
func DefaultCollisionConfig() sphero.CollisionConfig {
	return sphero.CollisionConfig{
		Method: 0x01,
		Xt:     0x20,
		Yt:     0x20,
		Xs:     0x20,
		Ys:     0x20,
		Dead:   0x01,
	}
}
