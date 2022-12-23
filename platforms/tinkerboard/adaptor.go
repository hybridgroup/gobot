package tinkerboard

import (
	"fmt"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/adaptors"
	"gobot.io/x/gobot/system"
)

const (
	debug = false

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
func NewAdaptor() *Adaptor {
	sys := system.NewAccesser(system.WithDigitalPinAccess("cdev"))
	c := &Adaptor{
		name: gobot.DefaultName("Tinker Board"),
		sys:  sys,
	}
	c.DigitalPinsAdaptor = adaptors.NewDigitalPinsAdaptor(sys, c.translateDigitalPin)
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
	if items == nil || len(items) == 0 {
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
