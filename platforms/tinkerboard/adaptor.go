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

type cdevPin struct {
	chip uint8
	line uint8
}

type gpioPinDefinition struct {
	sysfs int
	cdev  cdevPin
}

type analogPinDefinition struct {
	path   string
	r      bool // readable
	w      bool // writable
	bufLen uint16
}

type pwmPinDefinition struct {
	dir       string
	dirRegexp string
	channel   int
}

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

	a.AnalogPinsAdaptor = adaptors.NewAnalogPinsAdaptor(sys, a.translateAnalogPin)
	a.DigitalPinsAdaptor = adaptors.NewDigitalPinsAdaptor(sys, a.translateDigitalPin, digitalPinsOpts...)
	a.PWMPinsAdaptor = adaptors.NewPWMPinsAdaptor(sys, a.translatePWMPin, pwmPinsOpts...)
	a.I2cBusAdaptor = adaptors.NewI2cBusAdaptor(sys, a.validateI2cBusNumber, defaultI2cBusNumber)
	a.SpiBusAdaptor = adaptors.NewSpiBusAdaptor(sys, a.validateSpiBusNumber, defaultSpiBusNumber, defaultSpiChipNumber,
		defaultSpiMode, defaultSpiBitsNumber, defaultSpiMaxSpeed)
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

func (a *Adaptor) validateSpiBusNumber(busNr int) error {
	// Valid bus numbers are [0,2] which corresponds to /dev/spidev0.x, /dev/spidev2.x
	// x is the chip number <255
	if (busNr != 0) && (busNr != 2) {
		return fmt.Errorf("Bus number %d out of range", busNr)
	}
	return nil
}

func (a *Adaptor) validateI2cBusNumber(busNr int) error {
	// Valid bus number is [0..4] which corresponds to /dev/i2c-0 through /dev/i2c-4.
	// We don't support "/dev/i2c-6 DesignWare HDMI".
	if (busNr < 0) || (busNr > 4) {
		return fmt.Errorf("Bus number %d out of range", busNr)
	}
	return nil
}

func (a *Adaptor) translateAnalogPin(id string) (string, bool, bool, uint16, error) {
	pinInfo, ok := analogPinDefinitions[id]
	if !ok {
		return "", false, false, 0, fmt.Errorf("'%s' is not a valid id for a analog pin", id)
	}

	path := pinInfo.path
	info, err := a.sys.Stat(path)
	if err != nil {
		return "", false, false, 0, fmt.Errorf("Error (%v) on access '%s'", err, path)
	}
	if info.IsDir() {
		return "", false, false, 0, fmt.Errorf("The item '%s' is a directory, which is not expected", path)
	}

	return path, pinInfo.r, pinInfo.w, pinInfo.bufLen, nil
}

func (a *Adaptor) translateDigitalPin(id string) (string, int, error) {
	pindef, ok := gpioPinDefinitions[id]
	if !ok {
		return "", -1, fmt.Errorf("'%s' is not a valid id for a digital pin", id)
	}
	if a.sys.IsSysfsDigitalPinAccess() {
		return "", pindef.sysfs, nil
	}
	chip := fmt.Sprintf("gpiochip%d", pindef.cdev.chip)
	line := int(pindef.cdev.line)
	return chip, line, nil
}

func (a *Adaptor) translatePWMPin(id string) (string, int, error) {
	pinInfo, ok := pwmPinDefinitions[id]
	if !ok {
		return "", -1, fmt.Errorf("'%s' is not a valid id for a PWM pin", id)
	}
	path, err := pinInfo.findPWMDir(a.sys)
	if err != nil {
		return "", -1, err
	}
	return path, pinInfo.channel, nil
}

func (p pwmPinDefinition) findPWMDir(sys *system.Accesser) (string, error) {
	items, _ := sys.Find(p.dir, p.dirRegexp)
	if len(items) == 0 {
		return "", fmt.Errorf("No path found for PWM directory pattern, '%s' in path '%s'. See README.md for activation",
			p.dirRegexp, p.dir)
	}

	dir := items[0]
	info, err := sys.Stat(dir)
	if err != nil {
		return "", fmt.Errorf("Error (%v) on access '%s'", err, dir)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("The item '%s' is not a directory, which is not expected", dir)
	}

	return dir, nil
}
