package zero

import (
	"fmt"
	"sync"

	multierror "github.com/hashicorp/go-multierror"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/adaptors"
	"gobot.io/x/gobot/v2/system"
)

const (
	defaultI2cBusNumber = 3

	defaultSpiBusNumber  = 0
	defaultSpiChipNumber = 0
	defaultSpiMode       = 0
	defaultSpiBitsNumber = 8
	defaultSpiMaxSpeed   = 10000000
)

// Adaptor represents a Gobot Adaptor for the Radxa Zero
type Adaptor struct {
	name  string
	sys   *system.Accesser // used for unit tests only
	mutex *sync.Mutex
	*adaptors.AnalogPinsAdaptor
	*adaptors.DigitalPinsAdaptor
	*adaptors.PWMPinsAdaptor
	*adaptors.I2cBusAdaptor
	*adaptors.SpiBusAdaptor
	*adaptors.OneWireBusAdaptor
}

// NewAdaptor creates a Zero Adaptor
//
// Optional parameters:
//
//	adaptors.WithGpioSysfsAccess():	use legacy sysfs driver instead of default character device driver
//	adaptors.WithSpiGpioAccess(sclk, ncs, sdo, sdi):	use GPIO's instead of /dev/spidev#.#
//	adaptors.WithGpiosActiveLow(pin's): invert the pin behavior
//	adaptors.WithGpiosPullUp/Down(pin's): sets the internal pull resistor
//	adaptors.WithGpiosOpenDrain/Source(pin's): sets the output behavior
//	adaptors.WithGpioEventOnFallingEdge/RaisingEdge/BothEdges(pin, handler): activate edge detection
//
//	Optional parameters for PWM, see [adaptors.NewPWMPinsAdaptor]
func NewAdaptor(opts ...interface{}) *Adaptor {
	sys := system.NewAccesser()
	a := &Adaptor{
		name:  gobot.DefaultName("Zero"),
		sys:   sys,
		mutex: &sync.Mutex{},
	}

	var digitalPinsOpts []adaptors.DigitalPinsOptionApplier
	var pwmPinsOpts []adaptors.PwmPinsOptionApplier
	var spiBusOpts []adaptors.SpiBusOptionApplier
	for _, opt := range opts {
		switch o := opt.(type) {
		case adaptors.DigitalPinsOptionApplier:
			digitalPinsOpts = append(digitalPinsOpts, o)
		case adaptors.PwmPinsOptionApplier:
			pwmPinsOpts = append(pwmPinsOpts, o)
		case adaptors.SpiBusOptionApplier:
			spiBusOpts = append(spiBusOpts, o)
		default:
			panic(fmt.Sprintf("'%s' can not be applied on adaptor '%s'", opt, a.name))
		}
	}

	analogPinTranslator := adaptors.NewAnalogPinTranslator(sys, analogPinDefinitions)
	digitalPinTranslator := adaptors.NewDigitalPinTranslator(sys, gpioPinDefinitions)
	pwmPinTranslator := adaptors.NewPWMPinTranslator(sys, pwmPinDefinitions)
	// Valid bus numbers are [1,3,4] which corresponds to /dev/i2c-1, /dev/i2c-3, /dev/i2c-4
	// We don't support /dev/i2c-5 (DesignWare HDMI)
	i2cBusNumberValidator := adaptors.NewBusNumberValidator([]int{1, 3, 4})
	// Valid bus numbers are [0,1] which corresponds to /dev/spidev0.x, /dev/spidev1.x
	// x is the chip number <255
	spiBusNumberValidator := adaptors.NewBusNumberValidator([]int{0, 1})

	a.AnalogPinsAdaptor = adaptors.NewAnalogPinsAdaptor(sys, analogPinTranslator.Translate)
	a.DigitalPinsAdaptor = adaptors.NewDigitalPinsAdaptor(sys, digitalPinTranslator.Translate, digitalPinsOpts...)
	a.PWMPinsAdaptor = adaptors.NewPWMPinsAdaptor(sys, pwmPinTranslator.Translate, pwmPinsOpts...)
	a.I2cBusAdaptor = adaptors.NewI2cBusAdaptor(sys, i2cBusNumberValidator.Validate, defaultI2cBusNumber)
	a.SpiBusAdaptor = adaptors.NewSpiBusAdaptor(sys, spiBusNumberValidator.Validate, defaultSpiBusNumber,
		defaultSpiChipNumber, defaultSpiMode, defaultSpiBitsNumber, defaultSpiMaxSpeed, a.DigitalPinsAdaptor, spiBusOpts...)
	// pin ?? needs to be activated by DT-overlay w1-gpio
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
