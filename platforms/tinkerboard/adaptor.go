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
	pwmInvertedIdentifier = "inversed"

	defaultI2cBusNumber = 1

	defaultSpiBusNumber  = 0
	defaultSpiChipNumber = 0
	defaultSpiMode       = 0
	defaultSpiBitsNumber = 8
	defaultSpiMaxSpeed   = 500000
)

type pwmPinDefinition struct {
	channel   int
	dir       string
	dirRegexp string
}

// Adaptor represents a Gobot Adaptor for the ASUS Tinker Board
type Adaptor struct {
	name  string
	sys   *system.Accesser
	mutex sync.Mutex
	*adaptors.DigitalPinsAdaptor
	*adaptors.PWMPinsAdaptor
	*adaptors.I2cBusAdaptor
	*adaptors.SpiBusAdaptor
}

// NewAdaptor creates a Tinkerboard Adaptor
//
// Optional parameters:
//
//			adaptors.WithGpiodAccess():	use character device gpiod driver instead of sysfs (still used by default)
//			adaptors.WithSpiGpioAccess(sclk, nss, mosi, miso):	use GPIO's instead of /dev/spidev#.#
//	   adaptors.WithGpiosActiveLow(pin's): invert the pin behavior
//	   adaptors.WithGpiosPullUp/Down(pin's): sets the internal pull resistor
//
// note from RK3288 datasheet: "The pull direction (pullup or pulldown) for all of GPIOs are software-programmable", but
// the latter is not working for any pin (armbian 22.08.7)
func NewAdaptor(opts ...func(adaptors.Optioner)) *Adaptor {
	sys := system.NewAccesser(system.WithDigitalPinGpiodAccess())
	c := &Adaptor{
		name: gobot.DefaultName("Tinker Board"),
		sys:  sys,
	}
	c.DigitalPinsAdaptor = adaptors.NewDigitalPinsAdaptor(sys, c.translateDigitalPin, opts...)
	c.PWMPinsAdaptor = adaptors.NewPWMPinsAdaptor(sys, c.translatePWMPin,
		adaptors.WithPolarityInvertedIdentifier(pwmInvertedIdentifier))
	c.I2cBusAdaptor = adaptors.NewI2cBusAdaptor(sys, c.validateI2cBusNumber, defaultI2cBusNumber)
	c.SpiBusAdaptor = adaptors.NewSpiBusAdaptor(sys, c.validateSpiBusNumber, defaultSpiBusNumber, defaultSpiChipNumber,
		defaultSpiMode, defaultSpiBitsNumber, defaultSpiMaxSpeed)
	return c
}

// Name returns the name of the Adaptor
func (c *Adaptor) Name() string { return c.name }

// SetName sets the name of the Adaptor
func (c *Adaptor) SetName(n string) { c.name = n }

// Connect create new connection to board and pins.
func (c *Adaptor) Connect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if err := c.SpiBusAdaptor.Connect(); err != nil {
		return err
	}

	if err := c.I2cBusAdaptor.Connect(); err != nil {
		return err
	}

	if err := c.PWMPinsAdaptor.Connect(); err != nil {
		return err
	}
	return c.DigitalPinsAdaptor.Connect()
}

// Finalize closes connection to board, pins and bus
func (c *Adaptor) Finalize() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	err := c.DigitalPinsAdaptor.Finalize()

	if e := c.PWMPinsAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}

	if e := c.I2cBusAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}

	if e := c.SpiBusAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}
	return err
}

func (c *Adaptor) validateSpiBusNumber(busNr int) error {
	// Valid bus numbers are [0,2] which corresponds to /dev/spidev0.x, /dev/spidev2.x
	// x is the chip number <255
	if (busNr != 0) && (busNr != 2) {
		return fmt.Errorf("Bus number %d out of range", busNr)
	}
	return nil
}

func (c *Adaptor) validateI2cBusNumber(busNr int) error {
	// Valid bus number is [0..4] which corresponds to /dev/i2c-0 through /dev/i2c-4.
	// We don't support "/dev/i2c-6 DesignWare HDMI".
	if (busNr < 0) || (busNr > 4) {
		return fmt.Errorf("Bus number %d out of range", busNr)
	}
	return nil
}

func (c *Adaptor) translateDigitalPin(id string) (string, int, error) {
	pindef, ok := gpioPinDefinitions[id]
	if !ok {
		return "", -1, fmt.Errorf("'%s' is not a valid id for a digital pin", id)
	}
	if c.sys.IsSysfsDigitalPinAccess() {
		return "", pindef.sysfs, nil
	}
	chip := fmt.Sprintf("gpiochip%d", pindef.cdev.chip)
	line := int(pindef.cdev.line)
	return chip, line, nil
}

func (c *Adaptor) translatePWMPin(id string) (string, int, error) {
	pinInfo, ok := pwmPinDefinitions[id]
	if !ok {
		return "", -1, fmt.Errorf("'%s' is not a valid id for a PWM pin", id)
	}
	path, err := pinInfo.findPWMDir(c.sys)
	if err != nil {
		return "", -1, err
	}
	return path, pinInfo.channel, nil
}

func (p pwmPinDefinition) findPWMDir(sys *system.Accesser) (dir string, err error) {
	items, _ := sys.Find(p.dir, p.dirRegexp)
	if len(items) == 0 {
		return "", fmt.Errorf("No path found for PWM directory pattern, '%s' in path '%s'. See README.md for activation",
			p.dirRegexp, p.dir)
	}

	dir = items[0]
	info, err := sys.Stat(dir)
	if err != nil {
		return "", fmt.Errorf("Error (%v) on access '%s'", err, dir)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("The item '%s' is not a directory, which is not expected", dir)
	}

	return
}
