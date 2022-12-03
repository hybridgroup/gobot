package tinkerboard

import (
	"fmt"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/adaptors"
	"gobot.io/x/gobot/system"
)

const debug = false

const (
	pwmNormal        = "normal"
	pwmInverted      = "inversed"
	pwmPeriodDefault = 10000000 // 10000000ns = 10ms = 100Hz
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
	i2cBuses [5]i2c.I2cDevice
}

// NewAdaptor creates a Tinkerboard Adaptor
func NewAdaptor() *Adaptor {
	sys := system.NewAccesser("cdev")
	c := &Adaptor{
		name: gobot.DefaultName("Tinker Board"),
		sys:  sys,
	}
	c.DigitalPinsAdaptor = adaptors.NewDigitalPinsAdaptor(sys, c.translateDigitalPin)
	c.PWMPinsAdaptor = adaptors.NewPWMPinsAdaptor(sys, pwmPeriodDefault, c.translatePWMPin)
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

	if err := c.DigitalPinsAdaptor.Connect(); err != nil {
		return err
	}

	return c.PWMPinsAdaptor.Connect()
}

// Finalize closes connection to board, pins and bus
func (c *Adaptor) Finalize() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	err := c.DigitalPinsAdaptor.Finalize()

	if e := c.PWMPinsAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}

	for _, bus := range c.i2cBuses {
		if bus != nil {
			if e := bus.Close(); e != nil {
				err = multierror.Append(err, e)
			}
		}
	}
	return err
}

// GetConnection returns a connection to a device on a specified i2c bus.
// Valid bus number is [0..4] which corresponds to /dev/i2c-0 through /dev/i2c-4.
// We don't support "/dev/i2c-6 DesignWare HDMI".
func (c *Adaptor) GetConnection(address int, bus int) (connection i2c.Connection, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if (bus < 0) || (bus > 4) {
		return nil, fmt.Errorf("Bus number %d out of range", bus)
	}
	if c.i2cBuses[bus] == nil {
		c.i2cBuses[bus], err = c.sys.NewI2cDevice(fmt.Sprintf("/dev/i2c-%d", bus))
	}
	return i2c.NewConnection(c.i2cBuses[bus], address), err
}

// GetDefaultBus returns the default i2c bus for this platform.
func (c *Adaptor) GetDefaultBus() int {
	return 1
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

// TODO: test for this function:
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
		return "", fmt.Errorf("No path found for PWM directory pattern, '%s' in path '%s'. See README.md for activation", p.dirRegexp, p.dir)
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
