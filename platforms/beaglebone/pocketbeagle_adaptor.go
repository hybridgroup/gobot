package beaglebone

import "gobot.io/x/gobot"

// PocketBeagleAdaptor is the Gobot Adaptor for the PocketBeagle
// For more information check out:
// 		http://beagleboard.org/pocket
//
type PocketBeagleAdaptor struct {
	*Adaptor
}

// NewPocketBeagleAdaptor creates a new Adaptor for the PocketBeagle
func NewPocketBeagleAdaptor() *PocketBeagleAdaptor {
	a := NewAdaptor()
	a.SetName(gobot.DefaultName("PocketBeagle"))
	a.pinMap = pocketBeaglePinMap
	a.pwmPinMap = pocketBeaglePwmPinMap
	a.analogPinMap = pocketBeagleAnalogPinMap

	return &PocketBeagleAdaptor{
		Adaptor: a,
	}
}
