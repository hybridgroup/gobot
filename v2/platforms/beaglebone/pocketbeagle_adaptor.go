package beaglebone

import (
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/adaptors"
)

// PocketBeagleAdaptor is the Gobot Adaptor for the PocketBeagle
// For more information check out:
//
//	http://beagleboard.org/pocket
type PocketBeagleAdaptor struct {
	*Adaptor
}

// NewPocketBeagleAdaptor creates a new Adaptor for the PocketBeagle
func NewPocketBeagleAdaptor(opts ...func(adaptors.Optioner)) *PocketBeagleAdaptor {
	a := NewAdaptor(opts...)
	a.SetName(gobot.DefaultName("PocketBeagle"))
	a.pinMap = pocketBeaglePinMap
	a.pwmPinMap = pocketBeaglePwmPinMap
	a.analogPinMap = pocketBeagleAnalogPinMap

	return &PocketBeagleAdaptor{
		Adaptor: a,
	}
}
