package tinkerboard2

import (
	"fmt"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/adaptors"
	"gobot.io/x/gobot/v2/platforms/tinkerboard"
	"gobot.io/x/gobot/v2/system"
)

const (
	defaultI2cBusNumber = 7 // i2c-7 (header pins 27, 28)

	defaultSpiBusNumber  = 1 // spidev1.x (header pins 19, 21, 23, 24)
	defaultSpiChipNumber = 0
	defaultSpiMode       = 0
	defaultSpiBitsNumber = 8
	defaultSpiMaxSpeed   = 500000
)

type Tinkerboard2Adaptor struct {
	*tinkerboard.Adaptor
	sys *system.Accesser // used for unit tests only
}

// NewAdaptor creates a Tinkerboard-2 Adaptor
//
// Optional parameters:
//
//	adaptors.WithSysfsAccess():	use legacy sysfs driver instead of default character device gpiod
//	adaptors.WithSpiGpioAccess(sclk, ncs, sdo, sdi):	use GPIO's instead of /dev/spidev#.#
//	adaptors.WithGpiosActiveLow(pin's): invert the pin behavior
//	adaptors.WithGpiosPullUp/Down(pin's): sets the internal pull resistor
//
//	Optional parameters for PWM, see [adaptors.NewPWMPinsAdaptor]
func NewAdaptor(opts ...interface{}) *Tinkerboard2Adaptor {
	sys := system.NewAccesser()
	a := tinkerboard.NewAdaptor(opts...)
	a.SetName(gobot.DefaultName("Tinker Board 2"))

	var digitalPinsOpts []adaptors.DigitalPinsOptionApplier
	var pwmPinsOpts []adaptors.PwmPinsOptionApplier
	for _, opt := range opts {
		switch o := opt.(type) {
		case adaptors.DigitalPinsOptionApplier:
			digitalPinsOpts = append(digitalPinsOpts, o)
		case adaptors.PwmPinsOptionApplier:
			pwmPinsOpts = append(pwmPinsOpts, o)
		default:
			panic(fmt.Sprintf("'%s' can not be applied on adaptor '%s'", opt, a.Name()))
		}
	}

	// note: only adaptors different from tinkerboard needs to be re-assigned
	digitalPinTranslator := adaptors.NewDigitalPinTranslator(sys, gpioPinDefinitions)
	pwmPinTranslator := adaptors.NewPWMPinTranslator(sys, pwmPinDefinitions)
	// Valid bus numbers are [6..8] which corresponds to /dev/i2c-6 through /dev/i2c-8.
	// We don't support "/dev/i2c-0, /dev/i2c-3, /dev/i2c-4".
	i2cBusNumberValidator := adaptors.NewBusNumberValidator([]int{6, 7, 8})
	spiBusNumberValidator := adaptors.NewBusNumberValidator([]int{1, 5})

	a.DigitalPinsAdaptor = adaptors.NewDigitalPinsAdaptor(sys, digitalPinTranslator.Translate, digitalPinsOpts...)
	a.PWMPinsAdaptor = adaptors.NewPWMPinsAdaptor(sys, pwmPinTranslator.Translate, pwmPinsOpts...)
	a.I2cBusAdaptor = adaptors.NewI2cBusAdaptor(sys, i2cBusNumberValidator.Validate, defaultI2cBusNumber)
	a.SpiBusAdaptor = adaptors.NewSpiBusAdaptor(sys, spiBusNumberValidator.Validate, defaultSpiBusNumber,
		defaultSpiChipNumber, defaultSpiMode, defaultSpiBitsNumber, defaultSpiMaxSpeed)

	return &Tinkerboard2Adaptor{Adaptor: a, sys: sys}
}
