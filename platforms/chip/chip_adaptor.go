package chip

import (
	"errors"
	"os/exec"
	"strings"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot/sysfs"
)

// Adaptor represents a Gobot Adaptor for a C.H.I.P.
type Adaptor struct {
	name        string
	digitalPins map[int]sysfs.DigitalPin
	pinMap      map[string]int
	i2cDevice   sysfs.I2cDevice
}

var pinsOriginal = map[string]int{
	"XIO-P0": 408,
	"XIO-P1": 409,
	"XIO-P2": 410,
	"XIO-P3": 411,
	"XIO-P4": 412,
	"XIO-P5": 413,
	"XIO-P6": 414,
	"XIO-P7": 415,
}

var pins44 = map[string]int{
	"XIO-P0": 1013,
	"XIO-P1": 1014,
	"XIO-P2": 1015,
	"XIO-P3": 1016,
	"XIO-P4": 1017,
	"XIO-P5": 1018,
	"XIO-P6": 1019,
	"XIO-P7": 1020,
}

// NewAdaptor creates a C.H.I.P. Adaptor
func NewAdaptor() *Adaptor {
	c := &Adaptor{
		name:        "CHIP",
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
	return
}

// Finalize closes connection to board and pins
func (c *Adaptor) Finalize() (err error) {
	for _, pin := range c.digitalPins {
		if pin != nil {
			if e := pin.Unexport(); e != nil {
				err = multierror.Append(err, e)
			}
		}
	}
	if c.i2cDevice != nil {
		if e := c.i2cDevice.Close(); e != nil {
			err = multierror.Append(err, e)
		}
	}
	return
}

func (c *Adaptor) setPins() {
	kernel := getKernel()
	if kernel[:3] == "4.3" {
		c.pinMap = pinsOriginal
	} else {
		c.pinMap = pins44
	}
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
// Valids pins are XIO-P0 through XIO-P7 (pins 13-20 on header 14).
func (c *Adaptor) DigitalRead(pin string) (val int, err error) {
	sysfsPin, err := c.digitalPin(pin, sysfs.IN)
	if err != nil {
		return
	}
	return sysfsPin.Read()
}

// DigitalWrite writes digital value to the specified pin.
// Valids pins are XIO-P0 through XIO-P7 (pins 13-20 on header 14).
func (c *Adaptor) DigitalWrite(pin string, val byte) (err error) {
	sysfsPin, err := c.digitalPin(pin, sysfs.OUT)
	if err != nil {
		return err
	}
	return sysfsPin.Write(int(val))
}

// I2cStart starts an i2c device in specified address.
// This assumes that the bus used is /dev/i2c-1, which corresponds to
// pins labeled TWI1-SDA and TW1-SCK (pins 9 and 11 on header 13).
func (c *Adaptor) I2cStart(address int) (err error) {
	if c.i2cDevice == nil {
		c.i2cDevice, err = sysfs.NewI2cDevice("/dev/i2c-1", address)
	}
	return err
}

// I2cWrite writes data to i2c device
func (c *Adaptor) I2cWrite(address int, data []byte) (err error) {
	if err = c.i2cDevice.SetAddress(address); err != nil {
		return
	}
	_, err = c.i2cDevice.Write(data)
	return
}

// I2cRead returns value from i2c device using specified size
func (c *Adaptor) I2cRead(address int, size int) (data []byte, err error) {
	if err = c.i2cDevice.SetAddress(address); err != nil {
		return
	}
	data = make([]byte, size)
	_, err = c.i2cDevice.Read(data)
	return
}

func getKernel() string {
	result, _ := exec.Command("uname", "-r").Output()

	return strings.TrimSpace(string(result))
}
