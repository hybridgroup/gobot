package gpio

import (
	"github.com/hybridgroup/gobot"
)

var _ gobot.Driver = (*GroveTouchDriver)(nil)
var _ gobot.Driver = (*GroveSoundSensorDriver)(nil)
var _ gobot.Driver = (*GroveButtonDriver)(nil)
var _ gobot.Driver = (*GroveBuzzerDriver)(nil)
var _ gobot.Driver = (*GroveLightSensorDriver)(nil)
var _ gobot.Driver = (*GrovePiezoVibrationSensorDriver)(nil)
var _ gobot.Driver = (*GroveLedDriver)(nil)
var _ gobot.Driver = (*GroveRotaryDriver)(nil)
var _ gobot.Driver = (*GroveRelayDriver)(nil)
