package up2

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/spi"
	"gobot.io/x/gobot/platforms/adaptors"
	"gobot.io/x/gobot/system"
)

const (
	// LEDRed is the built-in red LED.
	LEDRed = "red"
	// LEDBlue is the built-in blue LED.
	LEDBlue = "blue"
	// LEDGreen is the built-in green LED.
	LEDGreen = "green"
	// LEDYellow is the built-in yellow LED.
	LEDYellow = "yellow"

	defaultI2cBusNumber = 5
)

type sysfsPin struct {
	pin    int
	pwmPin int
}

// Adaptor represents a Gobot Adaptor for the Upboard UP2
type Adaptor struct {
	name    string
	sys     *system.Accesser
	mutex   sync.Mutex
	pinmap  map[string]sysfsPin
	ledPath string
	*adaptors.DigitalPinsAdaptor
	*adaptors.PWMPinsAdaptor
	*adaptors.I2cBusAdaptor
	spiBuses [2]spi.Connection
}

// NewAdaptor creates a UP2 Adaptor
func NewAdaptor() *Adaptor {
	sys := system.NewAccesser()
	c := &Adaptor{
		name:    gobot.DefaultName("UP2"),
		sys:     sys,
		ledPath: "/sys/class/leds/upboard:%s:/brightness",
		pinmap:  fixedPins,
	}
	c.DigitalPinsAdaptor = adaptors.NewDigitalPinsAdaptor(sys, c.translateDigitalPin)
	c.PWMPinsAdaptor = adaptors.NewPWMPinsAdaptor(sys, c.translatePWMPin)
	c.I2cBusAdaptor = adaptors.NewI2cBusAdaptor(sys, c.validateI2cBusNumber, defaultI2cBusNumber)
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

	if err := c.I2cBusAdaptor.Connect(); err != nil {
		return err
	}

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

	if e := c.I2cBusAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}

	for _, bus := range c.spiBuses {
		if bus != nil {
			if e := bus.Close(); e != nil {
				err = multierror.Append(err, e)
			}
		}
	}

	return err
}

// DigitalWrite writes digital value to the specified pin.
func (c *Adaptor) DigitalWrite(id string, val byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// is it one of the built-in LEDs?
	if id == LEDRed || id == LEDBlue || id == LEDGreen || id == LEDYellow {
		pinPath := fmt.Sprintf(c.ledPath, id)
		fi, err := c.sys.OpenFile(pinPath, os.O_WRONLY|os.O_APPEND, 0666)
		defer fi.Close()
		if err != nil {
			return err
		}
		_, err = fi.WriteString(strconv.Itoa(int(val)))
		return err
	}

	return c.DigitalPinsAdaptor.DigitalWrite(id, val)
}

// GetSpiConnection returns an spi connection to a device on a specified bus.
// Valid bus number is [0..1] which corresponds to /dev/spidev0.0 through /dev/spidev0.1.
func (c *Adaptor) GetSpiConnection(busNum, chipNum, mode, bits int, maxSpeed int64) (connection spi.Connection, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if (busNum < 0) || (busNum > 1) {
		return nil, fmt.Errorf("Bus number %d out of range", busNum)
	}

	if c.spiBuses[busNum] == nil {
		c.spiBuses[busNum], err = spi.GetSpiConnection(busNum, chipNum, mode, bits, maxSpeed)
	}

	return c.spiBuses[busNum], err
}

// GetSpiDefaultBus returns the default spi bus for this platform.
func (c *Adaptor) GetSpiDefaultBus() int {
	return 0
}

// GetSpiDefaultChip returns the default spi chip for this platform.
func (c *Adaptor) GetSpiDefaultChip() int {
	return 0
}

// GetSpiDefaultMode returns the default spi mode for this platform.
func (c *Adaptor) GetSpiDefaultMode() int {
	return 0
}

// GetSpiDefaultBits returns the default spi number of bits for this platform.
func (c *Adaptor) GetSpiDefaultBits() int {
	return 8
}

// GetSpiDefaultMaxSpeed returns the default spi max speed for this platform.
func (c *Adaptor) GetSpiDefaultMaxSpeed() int64 {
	return 500000
}

func (c *Adaptor) validateI2cBusNumber(busNr int) error {
	// Valid bus number is [5..6] which corresponds to /dev/i2c-5 through /dev/i2c-6.
	if (busNr < 5) || (busNr > 6) {
		return fmt.Errorf("Bus number %d out of range", busNr)
	}
	return nil
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
