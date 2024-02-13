package beaglebone

import (
	"gobot.io/x/gobot/v2"
)

// PocketBeagleAdaptor is the Gobot Adaptor for the PocketBeagle
// For more information check out:
//
//	http://beagleboard.org/pocket
type PocketBeagleAdaptor struct {
	*Adaptor
}

// NewPocketBeagleAdaptor creates a new Adaptor for the PocketBeagle
// Optional parameters:
//
//	adaptors.WithGpiodAccess():	use character device gpiod driver instead of sysfs
//	adaptors.WithSpiGpioAccess(sclk, ncs, sdo, sdi):	use GPIO's instead of /dev/spidev#.#
//
//	Optional parameters for PWM, see [adaptors.NewPWMPinsAdaptor]
func NewPocketBeagleAdaptor(opts ...interface{}) *PocketBeagleAdaptor {
	a := NewAdaptor(opts...)
	a.SetName(gobot.DefaultName("PocketBeagle"))
	a.pinMap = pocketBeaglePinMap
	a.pwmPinMap = pocketBeaglePwmPinMap
	a.analogPinMap = pocketBeagleAnalogPinMap

	return &PocketBeagleAdaptor{
		Adaptor: a,
	}
}
