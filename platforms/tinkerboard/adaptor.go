package tinkerboard

import (
	"errors"
	"fmt"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/sysfs"
)

type sysfsPin struct {
	pin    int
	pwmPin int
}

// Adaptor represents a Gobot Adaptor for the ASUS Tinker Board
type Adaptor struct {
	name        string
	pinmap      map[string]sysfsPin
	digitalPins map[int]*sysfs.DigitalPin
	pwmPins     map[int]*sysfs.PWMPin
	i2cBuses    [2]i2c.I2cDevice
	mutex       *sync.Mutex
}

// NewAdaptor creates a Tinkerboard Adaptor
func NewAdaptor() *Adaptor {
	c := &Adaptor{
		name:  gobot.DefaultName("Tinker Board"),
		mutex: &sync.Mutex{},
	}

	c.setPins()
	return c
}

// Name returns the name of the Adaptor
func (c *Adaptor) Name() string { return c.name }

// SetName sets the name of the Adaptor
func (c *Adaptor) SetName(n string) { c.name = n }

// Connect initializes the board
func (c *Adaptor) Connect() (err error) {
	return nil
}

// Finalize closes connection to board and pins
func (c *Adaptor) Finalize() (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, pin := range c.digitalPins {
		if pin != nil {
			if e := pin.Unexport(); e != nil {
				err = multierror.Append(err, e)
			}
		}
	}
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
			if e := bus.Close(); e != nil {
				err = multierror.Append(err, e)
			}
		}
	}
	return
}

// DigitalRead reads digital value from the specified pin.
func (c *Adaptor) DigitalRead(pin string) (val int, err error) {
	sysfsPin, err := c.DigitalPin(pin, sysfs.IN)
	if err != nil {
		return
	}
	return sysfsPin.Read()
}

// DigitalWrite writes digital value to the specified pin.
func (c *Adaptor) DigitalWrite(pin string, val byte) (err error) {
	sysfsPin, err := c.DigitalPin(pin, sysfs.OUT)
	if err != nil {
		return err
	}
	return sysfsPin.Write(int(val))
}

// PwmWrite writes a PWM signal to the specified pin
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

// TODO: take into account the actual period setting, not just assume default
const pwmPeriod = 10000000

// ServoWrite writes a servo signal to the specified pin
func (c *Adaptor) ServoWrite(pin string, angle byte) (err error) {
	pwmPin, err := c.PWMPin(pin)
	if err != nil {
		return
	}

	// 0.5 ms => -90
	// 1.5 ms =>   0
	// 2.0 ms =>  90
	const minDuty = 100 * 0.0005 * pwmPeriod
	const maxDuty = 100 * 0.0020 * pwmPeriod
	duty := uint32(gobot.ToScale(gobot.FromScale(float64(angle), 0, 180), minDuty, maxDuty))
	return pwmPin.SetDutyCycle(duty)
}

// DigitalPin returns matched digitalPin for specified values
func (c *Adaptor) DigitalPin(pin string, dir string) (sysfsPin sysfs.DigitalPinner, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	i, err := c.translatePin(pin)

	if err != nil {
		return
	}

	if c.digitalPins[i] == nil {
		c.digitalPins[i] = sysfs.NewDigitalPin(i)
		if err = c.digitalPins[i].Export(); err != nil {
			return
		}
	}

	if err = c.digitalPins[i].Direction(dir); err != nil {
		return
	}

	return c.digitalPins[i], nil
}

// PWMPin returns matched pwmPin for specified pin number
func (c *Adaptor) PWMPin(pin string) (sysfsPin sysfs.PWMPinner, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	i, err := c.translatePwmPin(pin)
	if err != nil {
		return nil, err
	}
	if i == -1 {
		return nil, errors.New("Not a PWM pin")
	}

	if c.pwmPins[i] == nil {
		newPin := sysfs.NewPWMPin(i)
		if err = newPin.Export(); err != nil {
			return
		}
		// Make sure pwm is disabled when setting polarity
		if err = newPin.Enable(false); err != nil {
			return
		}
		if err = newPin.InvertPolarity(false); err != nil {
			return
		}
		if err = newPin.Enable(true); err != nil {
			return
		}
		if err = newPin.SetPeriod(10000000); err != nil {
			return
		}
		c.pwmPins[i] = newPin
	}

	sysfsPin = c.pwmPins[i]
	return
}

// GetConnection returns a connection to a device on a specified bus.
// Valid bus number is [0..1] which corresponds to /dev/i2c-0 through /dev/i2c-1.
func (c *Adaptor) GetConnection(address int, bus int) (connection i2c.Connection, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if (bus < 0) || (bus > 1) {
		return nil, fmt.Errorf("Bus number %d out of range", bus)
	}
	if c.i2cBuses[bus] == nil {
		c.i2cBuses[bus], err = sysfs.NewI2cDevice(fmt.Sprintf("/dev/i2c-%d", bus))
	}
	return i2c.NewConnection(c.i2cBuses[bus], address), err
}

// GetDefaultBus returns the default i2c bus for this platform
func (c *Adaptor) GetDefaultBus() int {
	return 1
}

func (c *Adaptor) setPins() {
	c.digitalPins = make(map[int]*sysfs.DigitalPin)
	c.pwmPins = make(map[int]*sysfs.PWMPin)
	c.pinmap = fixedPins
}

func (c *Adaptor) translatePin(pin string) (i int, err error) {
	if val, ok := c.pinmap[pin]; ok {
		i = val.pin
	} else {
		err = errors.New("Not a valid pin")
	}
	return
}

func (c *Adaptor) translatePwmPin(pin string) (i int, err error) {
	if val, ok := c.pinmap[pin]; ok {
		i = val.pwmPin
	} else {
		err = errors.New("Not a valid pin")
	}
	return
}
