package joule

import (
	"fmt"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/adaptors"
	"gobot.io/x/gobot/v2/system"
)

const defaultI2cBusNumber = 0

type sysfsPin struct {
	pin    int
	pwmPin int
}

// Adaptor represents an Intel Joule
type Adaptor struct {
	name  string
	sys   *system.Accesser
	mutex sync.Mutex
	*adaptors.DigitalPinsAdaptor
	*adaptors.PWMPinsAdaptor
	*adaptors.I2cBusAdaptor
}

// NewAdaptor returns a new Joule Adaptor
//
// Optional parameters:
//
//	adaptors.WithGpiodAccess():	use character device gpiod driver instead of sysfs
//	adaptors.WithSpiGpioAccess(sclk, nss, mosi, miso):	use GPIO's instead of /dev/spidev#.#
func NewAdaptor(opts ...func(adaptors.Optioner)) *Adaptor {
	sys := system.NewAccesser()
	c := &Adaptor{
		name: gobot.DefaultName("Joule"),
		sys:  sys,
	}
	c.DigitalPinsAdaptor = adaptors.NewDigitalPinsAdaptor(sys, c.translateDigitalPin, opts...)
	c.PWMPinsAdaptor = adaptors.NewPWMPinsAdaptor(sys, c.translatePWMPin, adaptors.WithPWMPinInitializer(pwmPinInitializer))
	c.I2cBusAdaptor = adaptors.NewI2cBusAdaptor(sys, c.validateI2cBusNumber, defaultI2cBusNumber)
	return c
}

// Name returns the Adaptors name
func (c *Adaptor) Name() string { return c.name }

// SetName sets the Adaptors name
func (c *Adaptor) SetName(n string) { c.name = n }

// Connect create new connection to board and pins.
func (c *Adaptor) Connect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if err := c.I2cBusAdaptor.Connect(); err != nil {
		return err
	}

	if err := c.PWMPinsAdaptor.Connect(); err != nil {
		return err
	}
	return c.DigitalPinsAdaptor.Connect()
}

// Finalize releases all i2c devices and exported digital and pwm pins.
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

	return err
}

func (c *Adaptor) validateI2cBusNumber(busNr int) error {
	// Valid bus number is [0..2] which corresponds to /dev/i2c-0 through /dev/i2c-2.
	if (busNr < 0) || (busNr > 2) {
		return fmt.Errorf("Bus number %d out of range", busNr)
	}
	return nil
}

func (c *Adaptor) translateDigitalPin(id string) (string, int, error) {
	if val, ok := sysfsPinMap[id]; ok {
		return "", val.pin, nil
	}
	return "", -1, fmt.Errorf("'%s' is not a valid id for a digital pin", id)
}

func (c *Adaptor) translatePWMPin(id string) (string, int, error) {
	sysPin, ok := sysfsPinMap[id]
	if !ok {
		return "", -1, fmt.Errorf("'%s' is not a valid id for a pin", id)
	}
	if sysPin.pwmPin == -1 {
		return "", -1, fmt.Errorf("'%s' is not a valid id for a PWM pin", id)
	}
	return "/sys/class/pwm/pwmchip0", sysPin.pwmPin, nil
}

func pwmPinInitializer(pin gobot.PWMPinner) error {
	if err := pin.Export(); err != nil {
		return err
	}
	if err := pin.SetPeriod(10000000); err != nil {
		return err
	}
	return pin.SetEnabled(true)
}
