package chip

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
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

var fixedPins = map[string]int{
	"PWM0":     34,
	"AP-EINT3": 35,

	"TWI1-SCK": 47,
	"TWI1-SDA": 48,
	"TWI2-SCK": 49,
	"TWI2-SDA": 50,

	"LCD-D2":    98,
	"LCD-D3":    99,
	"LCD-D4":    100,
	"LCD-D5":    101,
	"LCD-D6":    102,
	"LCD-D7":    103,
	"LCD-D10":   106,
	"LCD-D11":   107,
	"LCD-D12":   108,
	"LCD-D13":   109,
	"LCD-D14":   110,
	"LCD-D15":   111,
	"LCD-D18":   114,
	"LCD-D19":   115,
	"LCD-D20":   116,
	"LCD-D21":   117,
	"LCD-D22":   118,
	"LCD-D23":   119,
	"LCD-CLK":   120,
	"LCD-DE":    121,
	"LCD-HSYNC": 122,
	"LCD-VSYNC": 123,

	"CSIPCK":   128,
	"CSICK":    129,
	"CSIHSYNC": 130,
	"CSIVSYNC": 131,
	"CSID0":    132,
	"CSID1":    133,
	"CSID2":    134,
	"CSID3":    135,
	"CSID4":    136,
	"CSID5":    137,
	"CSID6":    138,
	"CSID7":    139,

	"AP-EINT1": 193,

	"UART1-TX": 195,
	"UART1-RX": 196,
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
	c.pinMap = fixedPins
	baseAddr, _ := getXIOBase()
	for i := 0; i < 8; i++ {
		pin := fmt.Sprintf("XIO-P%d", i)
		c.pinMap[pin] = baseAddr + i
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

func getXIOBase() (baseAddr int, err error) {
	// Default to original base from 4.3 kernel
	baseAddr = 408
	const expanderID = "pcf8574a"

	labels, err := filepath.Glob("/sys/class/gpio/*/label")
	if err != nil {
		return
	}

	for _, labelPath := range labels {
		label, err := ioutil.ReadFile(labelPath)
		if err != nil {
			return baseAddr, err
		}
		if strings.HasPrefix(string(label), expanderID) {
			expanderPath, _ := filepath.Split(labelPath)
			basePath := filepath.Join(expanderPath, "base")
			base, err := ioutil.ReadFile(basePath)
			if err != nil {
				return baseAddr, err
			}
			baseAddr, _ = strconv.Atoi(strings.TrimSpace(string(base)))
			break
		}
	}

	return baseAddr, nil
}
