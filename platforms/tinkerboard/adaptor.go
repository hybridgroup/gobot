package tinkerboard

import (
	"errors"
	"fmt"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/sysfs"
)

// Adaptor represents a Gobot Adaptor for Tinkerboard
type Adaptor struct {
	name        string
	digitalPins map[int]sysfs.DigitalPin
	pinMap      map[string]int
	i2cBuses    [2]sysfs.I2cDevice
	//	pwm         *pwmControl
}

// NewAdaptor creates a Tinkerboard Adaptor
func NewAdaptor() *Adaptor {
	c := &Adaptor{
		name:        gobot.DefaultName("Tinkerboard"),
		digitalPins: make(map[int]sysfs.DigitalPin),
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
	// if c.pwm != nil {
	// 	if e := c.closePWM(); e != nil {
	// 		err = multierror.Append(err, e)
	// 	}
	// 	c.pwm = nil
	// }
	for _, pin := range c.digitalPins {
		if pin != nil {
			if e := pin.Unexport(); e != nil {
				err = multierror.Append(err, e)
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

func (c *Adaptor) setPins() {
	c.pinMap = fixedPins
}

func (c *Adaptor) translatePin(pin string) (i int, err error) {
	if val, ok := c.pinMap[pin]; ok {
		i = val
	} else {
		err = errors.New("Not a valid pin")
	}
	return
}

// digitalPin returns matched digitalPin for specified values
func (c *Adaptor) digitalPin(pin string, dir string) (sysfsPin sysfs.DigitalPin, err error) {
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

// DigitalRead reads digital value from the specified pin.
func (c *Adaptor) DigitalRead(pin string) (val int, err error) {
	sysfsPin, err := c.digitalPin(pin, sysfs.IN)
	if err != nil {
		return
	}
	return sysfsPin.Read()
}

// DigitalWrite writes digital value to the specified pin.
func (c *Adaptor) DigitalWrite(pin string, val byte) (err error) {
	sysfsPin, err := c.digitalPin(pin, sysfs.OUT)
	if err != nil {
		return err
	}
	return sysfsPin.Write(int(val))
}

// GetConnection returns a connection to a device on a specified bus.
// Valid bus number is [0..1] which corresponds to /dev/i2c-0 through /dev/i2c-1.
func (c *Adaptor) GetConnection(address int, bus int) (connection i2c.Connection, err error) {
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
