package joule

import (
	"fmt"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/adaptors"
	"gobot.io/x/gobot/system"
)

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
	i2cBuses [3]i2c.I2cDevice
}

// NewAdaptor returns a new Joule Adaptor
func NewAdaptor() *Adaptor {
	sys := system.NewAccesser()
	c := &Adaptor{
		name: gobot.DefaultName("Joule"),
		sys:  sys,
	}
	c.DigitalPinsAdaptor = adaptors.NewDigitalPinsAdaptor(sys, c.translateDigitalPin)
	c.PWMPinsAdaptor = adaptors.NewPWMPinsAdaptor(sys, c.translatePWMPin, adaptors.WithPWMPinInitializer(pwmPinInitializer))
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

	for _, bus := range c.i2cBuses {
		if bus != nil {
			if errs := bus.Close(); errs != nil {
				err = multierror.Append(err, errs)
			}
		}
	}
	return err
}

// GetConnection returns an i2c connection to a device on a specified bus.
// Valid bus number is [0..2] which corresponds to /dev/i2c-0 through /dev/i2c-2.
func (c *Adaptor) GetConnection(address int, bus int) (connection i2c.Connection, err error) {
	if (bus < 0) || (bus > 2) {
		return nil, fmt.Errorf("Bus number %d out of range", bus)
	}
	if c.i2cBuses[bus] == nil {
		c.i2cBuses[bus], err = c.sys.NewI2cDevice(fmt.Sprintf("/dev/i2c-%d", bus))
	}
	return i2c.NewConnection(c.i2cBuses[bus], address), err
}

// GetDefaultBus returns the default i2c bus for this platform
func (c *Adaptor) GetDefaultBus() int {
	return 0
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
	return pin.Enable(true)
}
