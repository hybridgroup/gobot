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
// Optional parameters:
//
//	adaptors.WithGpioCdevAccess():	use character device driver instead of sysfs
//	adaptors.WithSpiGpioAccess(sclk, ncs, sdo, sdi):	use GPIO's instead of /dev/spidev#.#
//
//	Optional parameters for PWM, see [adaptors.NewPWMPinsAdaptor]
func NewPocketBeagleAdaptor(opts ...interface{}) *PocketBeagleAdaptor {
	a := NewAdaptor(opts...)
	a.SetName(gobot.DefaultName("PocketBeagle"))

	analogPinTranslator := adaptors.NewAnalogPinTranslator(a.sys, pocketBeagleAnalogPinMap)
	pwmPinTranslator := adaptors.NewPWMPinTranslator(a.sys, pocketBeaglePwmPinMap)

	a.AnalogPinsAdaptor = adaptors.NewAnalogPinsAdaptor(a.sys, analogPinTranslator.Translate)
	a.pinMap = pocketBeaglePinMap
	a.pwmPinTranslate = pwmPinTranslator.Translate

	return &PocketBeagleAdaptor{
		Adaptor: a,
	}
}
