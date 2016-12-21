package i2c

import (
	"gobot.io/x/gobot"
)

var _ gobot.Driver = (*GroveLcdDriver)(nil)
var _ gobot.Driver = (*GroveAccelerometerDriver)(nil)
