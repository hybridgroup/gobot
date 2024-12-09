package tinkerboard

import (
	"fmt"
	"sync"

	multierror "github.com/hashicorp/go-multierror"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/adaptors"
	"gobot.io/x/gobot/v2/system"
)

const (
	defaultI2cBusNumber = 1

	defaultSpiBusNumber  = 0
	defaultSpiChipNumber = 0
	defaultSpiMode       = 0
	defaultSpiBitsNumber = 8
	defaultSpiMaxSpeed   = 500000
)

// Adaptor represents a Gobot Adaptor for the ASUS Tinker Board
type Adaptor struct {
	name  string
	sys   *system.Accesser
	mutex *sync.Mutex
	*adaptors.AnalogPinsAdaptor
	*adaptors.DigitalPinsAdaptor
	*adaptors.PWMPinsAdaptor
	*adaptors.I2cBusAdaptor
	*adaptors.SpiBusAdaptor
	*adaptors.OneWireBusAdaptor
}

// NewAdaptor creates a Tinkerboard Adaptor
//
// Optional parameters:
//
//	adaptors.WithGpiodAccess():	use character device gpiod driver instead of sysfs (still used by default)
//	adaptors.WithSpiGpioAccess(sclk, ncs, sdo, sdi):	use GPIO's instead of /dev/spidev#.#
//	adaptors.WithGpiosActiveLow(pin's): invert the pin behavior
//	adaptors.WithGpiosPullUp/Down(pin's): sets the internal pull resistor
//
//	Optional parameters for PWM, see [adaptors.NewPWMPinsAdaptor]
//
// note from RK3288 datasheet: "The pull direction (pullup or pulldown) for all of GPIOs are software-programmable", but
// the latter is not working for any pin (armbian 22.08.7)
func NewAdaptor(opts ...interface{}) *Adaptor {
	sys := system.NewAccesser(system.WithDigitalPinGpiodAccess())
	a := &Adaptor{
		name:  gobot.DefaultName("Tinker Board"),
		sys:   sys,
		mutex: &sync.Mutex{},
	}

	var digitalPinsOpts []func(adaptors.DigitalPinsOptioner)
	var pwmPinsOpts []adaptors.PwmPinsOptionApplier
	for _, opt := range opts {
		switch o := opt.(type) {
		case func(adaptors.DigitalPinsOptioner):
			digitalPinsOpts = append(digitalPinsOpts, o)
		case adaptors.PwmPinsOptionApplier:
			pwmPinsOpts = append(pwmPinsOpts, o)
		default:
			panic(fmt.Sprintf("'%s' can not be applied on adaptor '%s'", opt, a.name))
		}
	}

	analogPinTranslator := adaptors.NewAnalogPinTranslator(sys, analogPinDefinitions)
	digitalPinTranslator := adaptors.NewDigitalPinTranslator(sys, gpioPinDefinitions)
	pwmPinTranslator := adaptors.NewPWMPinTranslator(sys, pwmPinDefinitions)
	// Valid bus numbers are [0..4] which corresponds to /dev/i2c-0 through /dev/i2c-4.
	// We don't support "/dev/i2c-6 DesignWare HDMI".
	i2cBusNumberValidator := adaptors.NewBusNumberValidator([]int{0, 1, 2, 3, 4})
	// Valid bus numbers are [0,2] which corresponds to /dev/spidev0.x, /dev/spidev2.x
	// x is the chip number <255
	spiBusNumberValidator := adaptors.NewBusNumberValidator([]int{0, 2})

	a.AnalogPinsAdaptor = adaptors.NewAnalogPinsAdaptor(sys, analogPinTranslator.Translate)
	a.DigitalPinsAdaptor = adaptors.NewDigitalPinsAdaptor(sys, digitalPinTranslator.Translate, digitalPinsOpts...)
	a.PWMPinsAdaptor = adaptors.NewPWMPinsAdaptor(sys, pwmPinTranslator.Translate, pwmPinsOpts...)
	a.I2cBusAdaptor = adaptors.NewI2cBusAdaptor(sys, i2cBusNumberValidator.Validate, defaultI2cBusNumber)
	a.SpiBusAdaptor = adaptors.NewSpiBusAdaptor(sys, spiBusNumberValidator.Validate, defaultSpiBusNumber,
		defaultSpiChipNumber, defaultSpiMode, defaultSpiBitsNumber, defaultSpiMaxSpeed)
	a.OneWireBusAdaptor = adaptors.NewOneWireBusAdaptor(sys)

	return a
}

// Name returns the name of the Adaptor
func (a *Adaptor) Name() string { return a.name }

// SetName sets the name of the Adaptor
func (a *Adaptor) SetName(n string) { a.name = n }

// Connect create new connection to board and pins.
func (a *Adaptor) Connect() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if err := a.OneWireBusAdaptor.Connect(); err != nil {
		return err
	}

	if err := a.SpiBusAdaptor.Connect(); err != nil {
		return err
	}

	if err := a.I2cBusAdaptor.Connect(); err != nil {
		return err
	}

	if err := a.AnalogPinsAdaptor.Connect(); err != nil {
		return err
	}

	if err := a.PWMPinsAdaptor.Connect(); err != nil {
		return err
	}
	return a.DigitalPinsAdaptor.Connect()
}

// Finalize closes connection to board, pins and bus
func (a *Adaptor) Finalize() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	err := a.DigitalPinsAdaptor.Finalize()

	if e := a.PWMPinsAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}

	if e := a.AnalogPinsAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}

	if e := a.I2cBusAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}

	if e := a.SpiBusAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}

	if e := a.OneWireBusAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}

	return err
}
