package joule

import (
	"errors"
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
	pwmPins  map[int]gobot.PWMPinner
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
	return c
}

// Name returns the Adaptors name
func (c *Adaptor) Name() string { return c.name }

// SetName sets the Adaptors name
func (c *Adaptor) SetName(n string) { c.name = n }

// Connect initializes the Joule for use with the Arduino breakout board
func (c *Adaptor) Connect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.pwmPins = make(map[int]gobot.PWMPinner)
	return c.DigitalPinsAdaptor.Connect()
}

// Finalize releases all i2c devices and exported digital and pwm pins.
func (c *Adaptor) Finalize() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	err := c.DigitalPinsAdaptor.Finalize()

	for _, pin := range c.pwmPins {
		if pin != nil {
			if errs := pin.Enable(false); errs != nil {
				err = multierror.Append(err, errs)
			}
			if errs := pin.Unexport(); errs != nil {
				err = multierror.Append(err, errs)
			}
		}
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

// PwmWrite writes the 0-254 value to the specified pin
func (c *Adaptor) PwmWrite(pin string, val byte) (err error) {
	pwmPin, err := c.PWMPin(pin)
	if err != nil {
		return
	}
	period, err := pwmPin.Period()
	if err != nil {
		return err
	}
	duty := gobot.FromScale(float64(val), 0, 255.0)
	return pwmPin.SetDutyCycle(uint32(float64(period) * duty))
}

// PWMPin returns a PWM pin, implements the interface gobot.PWMPinnerProvider
func (c *Adaptor) PWMPin(pin string) (gobot.PWMPinner, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	sysPin, ok := sysfsPinMap[pin]
	if !ok {
		return nil, errors.New("Not a valid pin")
	}
	if sysPin.pwmPin != -1 {
		if c.pwmPins[sysPin.pwmPin] == nil {
			c.pwmPins[sysPin.pwmPin] = c.sys.NewPWMPin("/sys/class/pwm/pwmchip0", sysPin.pwmPin)
			if err := c.pwmPins[sysPin.pwmPin].Export(); err != nil {
				return nil, err
			}
			if err := c.pwmPins[sysPin.pwmPin].SetPeriod(10000000); err != nil {
				return nil, err
			}
			if err := c.pwmPins[sysPin.pwmPin].Enable(true); err != nil {
				return nil, err
			}
		}

		return c.pwmPins[sysPin.pwmPin], nil
	}

	return nil, errors.New("Not a PWM pin")
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
