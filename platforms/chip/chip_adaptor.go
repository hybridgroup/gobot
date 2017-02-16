package chip

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/sysfs"
)

// Adaptor represents a Gobot Adaptor for a C.H.I.P.
type Adaptor struct {
	name        string
	digitalPins map[int]sysfs.DigitalPin
	pinMap      map[string]int
	i2cBuses    [3]sysfs.I2cDevice
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
		name:        gobot.DefaultName("CHIP"),
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
// Valids pins are the XIO-P0 through XIO-P7 pins from the
// extender (pins 13-20 on header 14), as well as the SoC pins
// aka all the other pins.
func (c *Adaptor) DigitalRead(pin string) (val int, err error) {
	sysfsPin, err := c.digitalPin(pin, sysfs.IN)
	if err != nil {
		return
	}
	return sysfsPin.Read()
}

// DigitalWrite writes digital value to the specified pin.
// Valids pins are the XIO-P0 through XIO-P7 pins from the
// extender (pins 13-20 on header 14), as well as the SoC pins
// aka all the other pins.
func (c *Adaptor) DigitalWrite(pin string, val byte) (err error) {
	sysfsPin, err := c.digitalPin(pin, sysfs.OUT)
	if err != nil {
		return err
	}
	return sysfsPin.Write(int(val))
}

// GetConnection returns a connection to a device on a specified bus.
// Valid bus number is [0..2] which corresponds to /dev/i2c-0 through /dev/i2c-2.
func (c *Adaptor) GetConnection(address int, bus int) (connection i2c.Connection, err error) {
	if (bus < 0) || (bus > 2) {
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
