package chip

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/adaptors"
	"gobot.io/x/gobot/system"
)

// Valids pins are the XIO-P0 through XIO-P7 pins from the
// extender (pins 13-20 on header 14), as well as the SoC pins
// aka all the other pins.
type sysfsPin struct {
	pin    int
	pwmPin int
}

// Adaptor represents a Gobot Adaptor for a C.H.I.P.
type Adaptor struct {
	name string
	sys  *system.Accesser
	*adaptors.DigitalPinsAdaptor
	*adaptors.PWMPinsAdaptor
	mutex    sync.Mutex
	pinmap   map[string]sysfsPin
	i2cBuses [3]i2c.I2cDevice
}

// NewAdaptor creates a C.H.I.P. Adaptor
func NewAdaptor() *Adaptor {
	sys := system.NewAccesser()
	c := &Adaptor{
		name: gobot.DefaultName("CHIP"),
		sys:  sys,
	}

	c.pinmap = chipPins
	baseAddr, _ := getXIOBase()
	for i := 0; i < 8; i++ {
		pin := fmt.Sprintf("XIO-P%d", i)
		c.pinmap[pin] = sysfsPin{pin: baseAddr + i, pwmPin: -1}
	}

	c.DigitalPinsAdaptor = adaptors.NewDigitalPinsAdaptor(sys, c.translateDigitalPin)
	c.PWMPinsAdaptor = adaptors.NewPWMPinsAdaptor(sys, c.translatePWMPin)
	return c
}

// NewProAdaptor creates a C.H.I.P. Pro Adaptor
func NewProAdaptor() *Adaptor {
	c := NewAdaptor()
	c.name = gobot.DefaultName("CHIP Pro")
	c.pinmap = chipProPins
	return c
}

// Name returns the name of the Adaptor
func (c *Adaptor) Name() string { return c.name }

// SetName sets the name of the Adaptor
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

// Finalize closes connection to board and pins
func (c *Adaptor) Finalize() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	err := c.DigitalPinsAdaptor.Finalize()

	if e := c.PWMPinsAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}

	for _, bus := range c.i2cBuses {
		if bus != nil {
			if e := bus.Close(); e != nil {
				err = multierror.Append(err, e)
			}
		}
	}
	return err
}

// GetConnection returns a connection to a device on a specified bus.
// Valid bus number is [0..2] which corresponds to /dev/i2c-0 through /dev/i2c-2.
func (c *Adaptor) GetConnection(address int, bus int) (connection i2c.Connection, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

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

func (c *Adaptor) translateDigitalPin(id string) (string, int, error) {
	if val, ok := c.pinmap[id]; ok {
		return "", val.pin, nil
	}
	return "", -1, fmt.Errorf("'%s' is not a valid id for a digital pin", id)
}

func (c *Adaptor) translatePWMPin(id string) (string, int, error) {
	sysPin, ok := c.pinmap[id]
	if !ok {
		return "", -1, fmt.Errorf("'%s' is not a valid id for a pin", id)
	}
	if sysPin.pwmPin == -1 {
		return "", -1, fmt.Errorf("'%s' is not a valid id for a PWM pin", id)
	}
	return "/sys/class/pwm/pwmchip0", sysPin.pwmPin, nil
}
