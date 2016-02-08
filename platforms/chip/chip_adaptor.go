package chip

import (
	"errors"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/i2c"
	"github.com/hybridgroup/gobot/sysfs"
)

var _ gobot.Adaptor = (*ChipAdaptor)(nil)

var _ gpio.DigitalReader = (*ChipAdaptor)(nil)
var _ gpio.DigitalWriter = (*ChipAdaptor)(nil)

var _ i2c.I2c = (*ChipAdaptor)(nil)

type ChipAdaptor struct {
	name        string
	digitalPins map[int]sysfs.DigitalPin
	i2cDevice   sysfs.I2cDevice
}

var pins = map[string]int{
	"XIO-P0": 408,
	"XIO-P1": 409,
	"XIO-P2": 410,
	"XIO-P3": 411,
	"XIO-P4": 412,
	"XIO-P5": 413,
	"XIO-P6": 414,
	"XIO-P7": 415,
}

// NewChipAdaptor creates a ChipAdaptor with the specified name
func NewChipAdaptor(name string) *ChipAdaptor {
	c := &ChipAdaptor{
		name:        name,
		digitalPins: make(map[int]sysfs.DigitalPin),
	}
	return c
}

// Name returns the name of the ChipAdaptor
func (c *ChipAdaptor) Name() string { return c.name }

// Connect initializes the board
func (c *ChipAdaptor) Connect() (errs []error) {
	return
}

// Finalize closes connection to board and pins
func (c *ChipAdaptor) Finalize() (errs []error) {
	for _, pin := range c.digitalPins {
		if pin != nil {
			if err := pin.Unexport(); err != nil {
				errs = append(errs, err)
			}
		}
	}
	if c.i2cDevice != nil {
		if err := c.i2cDevice.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (c *ChipAdaptor) translatePin(pin string) (i int, err error) {
	if val, ok := pins[pin]; ok {
		i = val
	} else {
		err = errors.New("Not a valid pin")
	}
	return
}

// digitalPin returns matched digitalPin for specified values
func (c *ChipAdaptor) digitalPin(pin string, dir string) (sysfsPin sysfs.DigitalPin, err error) {
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
// Valids pins are XIO-P0 through XIO-P7 (pins 13-20 on header 14).
func (c *ChipAdaptor) DigitalRead(pin string) (val int, err error) {
	sysfsPin, err := c.digitalPin(pin, sysfs.IN)
	if err != nil {
		return
	}
	return sysfsPin.Read()
}

// DigitalWrite writes digital value to the specified pin.
// Valids pins are XIO-P0 through XIO-P7 (pins 13-20 on header 14).
func (c *ChipAdaptor) DigitalWrite(pin string, val byte) (err error) {
	sysfsPin, err := c.digitalPin(pin, sysfs.OUT)
	if err != nil {
		return err
	}
	return sysfsPin.Write(int(val))
}

// I2cStart starts an i2c device in specified address.
// This assumes that the bus used is /dev/i2c-1, which corresponds to
// pins labeled TWI1-SDA and TW1-SCK (pins 9 and 11 on header 13).
func (c *ChipAdaptor) I2cStart(address int) (err error) {
	if c.i2cDevice == nil {
		c.i2cDevice, err = sysfs.NewI2cDevice("/dev/i2c-1", address)
	}
	return err
}

// I2cWrite writes data to i2c device
func (c *ChipAdaptor) I2cWrite(address int, data []byte) (err error) {
	if err = c.i2cDevice.SetAddress(address); err != nil {
		return
	}
	_, err = c.i2cDevice.Write(data)
	return
}

// I2cRead returns value from i2c device using specified size
func (c *ChipAdaptor) I2cRead(address int, size int) (data []byte, err error) {
	if err = c.i2cDevice.SetAddress(address); err != nil {
		return
	}
	data = make([]byte, size)
	_, err = c.i2cDevice.Read(data)
	return
}
