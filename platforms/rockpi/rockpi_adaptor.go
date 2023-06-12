package rockpi

import (
	"errors"
	"fmt"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/adaptors"
	"gobot.io/x/gobot/v2/system"
)

const (
	procDeviceTreeModel = "/proc/device-tree/model"

	defaultI2cBusNumber = 7

	defaultSpiBusNumber  = 1
	defaultSpiChipNumber = 0
	defaultSpiMode       = 0
	defaultSpiBitsNumber = 8
	defaultSpiMaxSpeed   = 500000
)

// Adaptor is the Gobot Adaptor for Radxa's Rock Pi.
type Adaptor struct {
	name     string
	mutex    sync.Mutex
	sys      *system.Accesser
	revision string
	*adaptors.DigitalPinsAdaptor
	*adaptors.I2cBusAdaptor
	*adaptors.SpiBusAdaptor
}

// NewAdaptor creates a RockPi Adaptor
// Do not forget to enable the required overlays in /boot/hw_initfc.conf!
// See https://wiki.radxa.com/Rockpi4/dev/libmraa
//
// Optional parameters:
//
//	adaptors.WithGpiodAccess():	use character device gpiod driver instead of the default sysfs (does NOT work on RockPi4C+!)
//	adaptors.WithSpiGpioAccess(sclk, nss, mosi, miso):	use GPIO's instead of /dev/spidev#.#
//	adaptors.WithGpiosActiveLow(pin's): invert the pin behavior
func NewAdaptor(opts ...func(adaptors.Optioner)) *Adaptor {
	sys := system.NewAccesser()
	c := &Adaptor{
		name: gobot.DefaultName("RockPi"),
		sys:  sys,
	}
	c.DigitalPinsAdaptor = adaptors.NewDigitalPinsAdaptor(sys, c.getPinTranslatorFunction(), opts...)
	c.I2cBusAdaptor = adaptors.NewI2cBusAdaptor(sys, c.validateI2cBusNumber, defaultI2cBusNumber)
	c.SpiBusAdaptor = adaptors.NewSpiBusAdaptor(sys, c.validateSpiBusNumber, defaultSpiBusNumber, defaultSpiChipNumber,
		defaultSpiMode, defaultSpiBitsNumber, defaultSpiMaxSpeed)
	return c
}

// Name returns the Adaptor's name
func (c *Adaptor) Name() string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.name
}

// SetName sets the Adaptor's name
func (c *Adaptor) SetName(n string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.name = n
}

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

	return c.DigitalPinsAdaptor.Connect()
}

// Finalize closes connection to board and pins
func (c *Adaptor) Finalize() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	err := c.DigitalPinsAdaptor.Finalize()

	if e := c.I2cBusAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}

	if e := c.SpiBusAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}
	return err
}

// The RockPi4 has 2 SPI buses: 1, 2. See https://wiki.radxa.com/Rock4/hardware/gpio
// This could change in the future with other revisions!
func (c *Adaptor) validateSpiBusNumber(busNr int) error {
	if busNr != 1 && busNr != 2 {
		return fmt.Errorf("SPI Bus number %d invalid: only 1, 2 supported by current Rockchip.", busNr)
	}
	return nil
}

// The RockPi4 has 3 I2C buses: 2, 6, 7. See https://wiki.radxa.com/Rock4/hardware/gpio
// This could change in the future with other revisions!
func (c *Adaptor) validateI2cBusNumber(busNr int) error {
	if busNr != 2 && busNr != 6 && busNr != 7 {
		return fmt.Errorf("I2C Bus number %d invalid: only 2, 6, 7 supported by current Rockchip.", busNr)
	}
	return nil
}

func (c *Adaptor) getPinTranslatorFunction() func(string) (string, int, error) {
	return func(pin string) (chip string, line int, err error) {
		if val, ok := pins[pin][c.readRevision()]; ok {
			line = val
		} else if val, ok := pins[pin]["*"]; ok {
			line = val
		} else {
			err = errors.New("Not a valid pin")
			return
		}
		return "", line, nil
	}
}

func (c *Adaptor) readRevision() string {
	if c.revision == "" {
		content, err := c.sys.ReadFile(procDeviceTreeModel)
		if err != nil {
			return c.revision
		}
		model := string(content)
		switch model {
		case "Radxa ROCK 4":
			c.revision = "4"
		case "Radxa ROCK 4C+":
			c.revision = "4C+"
		default:
			c.revision = "4"
		}
	}

	return c.revision
}
