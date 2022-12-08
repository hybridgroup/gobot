package raspi

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/spi"
	"gobot.io/x/gobot/platforms/adaptors"
	"gobot.io/x/gobot/system"
)

const infoFile = "/proc/cpuinfo"

// Adaptor is the Gobot Adaptor for the Raspberry Pi
type Adaptor struct {
	name     string
	mutex    sync.Mutex
	sys      *system.Accesser
	revision string
	pwmPins  map[string]gobot.PWMPinner
	*adaptors.DigitalPinsAdaptor
	*adaptors.I2cBusAdaptor
	spiDevices         [2]spi.Connection
	spiDefaultMaxSpeed int64
	PiBlasterPeriod    uint32
}

// NewAdaptor creates a Raspi Adaptor
func NewAdaptor() *Adaptor {
	sys := system.NewAccesser("cdev")
	c := &Adaptor{
		name:            gobot.DefaultName("RaspberryPi"),
		sys:             sys,
		PiBlasterPeriod: 10000000,
	}
	c.DigitalPinsAdaptor = adaptors.NewDigitalPinsAdaptor(sys, c.getPinTranslatorFunction())
	c.I2cBusAdaptor = adaptors.NewI2cBusAdaptor(sys, c.validateI2cBusNumber, 1)
	return c
}

// Name returns the Adaptor's name
func (c *Adaptor) Name() string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.name
}

// SetName sets the Adaptor's name
func (c *Adaptor) SetName(n string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.name = n
}

// Connect create new connection to board and pins.
func (c *Adaptor) Connect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if err := c.I2cBusAdaptor.Connect(); err != nil {
		return err
	}

	c.pwmPins = make(map[string]gobot.PWMPinner)
	return c.DigitalPinsAdaptor.Connect()
}

// Finalize closes connection to board and pins
func (c *Adaptor) Finalize() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	err := c.DigitalPinsAdaptor.Finalize()

	for _, pin := range c.pwmPins {
		if pin != nil {
			if perr := pin.Unexport(); err != nil {
				err = multierror.Append(err, perr)
			}
		}
	}
	c.pwmPins = nil

	if e := c.I2cBusAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}

	for _, dev := range c.spiDevices {
		if dev != nil {
			if e := dev.Close(); e != nil {
				err = multierror.Append(err, e)
			}
		}
	}
	return err
}

// GetDefaultBus returns the default i2c bus for this platform.
// This overrides the base function due to the revision dependency.
func (c *Adaptor) GetDefaultBus() int {
	rev := c.readRevision()
	if rev == "2" || rev == "3" {
		return 1
	}
	return 0
}

// GetSpiConnection returns an spi connection to a device on a specified bus.
// Valid bus number is [0..1] which corresponds to /dev/spidev0.0 through /dev/spidev0.1.
func (c *Adaptor) GetSpiConnection(busNum, chipNum, mode, bits int, maxSpeed int64) (connection spi.Connection, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if (busNum < 0) || (busNum > 1) {
		return nil, fmt.Errorf("Bus number %d out of range", busNum)
	}

	if c.spiDevices[busNum] == nil {
		c.spiDevices[busNum], err = spi.GetSpiConnection(busNum, chipNum, mode, bits, maxSpeed)
	}

	return c.spiDevices[busNum], err
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

// GetSpiDefaultMaxSpeed returns the default spi bus for this platform.
func (c *Adaptor) GetSpiDefaultMaxSpeed() int64 {
	return 500000
}

// PWMPin returns a raspi.PWMPin which provides the gobot.PWMPinner interface
func (c *Adaptor) PWMPin(id string) (gobot.PWMPinner, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.pwmPin(id)
}

// PwmWrite writes a PWM signal to the specified pin
func (c *Adaptor) PwmWrite(pin string, val byte) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	sysPin, err := c.pwmPin(pin)
	if err != nil {
		return err
	}

	duty := uint32(gobot.FromScale(float64(val), 0, 255) * float64(c.PiBlasterPeriod))
	return sysPin.SetDutyCycle(duty)
}

// ServoWrite writes a servo signal to the specified pin
func (c *Adaptor) ServoWrite(pin string, angle byte) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	sysPin, err := c.pwmPin(pin)
	if err != nil {
		return err
	}

	duty := uint32(gobot.FromScale(float64(angle), 0, 180) * float64(c.PiBlasterPeriod))
	return sysPin.SetDutyCycle(duty)
}

func (c *Adaptor) validateI2cBusNumber(busNr int) error {
	// Valid bus number is [0..1] which corresponds to /dev/i2c-0 through /dev/i2c-1.
	if (busNr < 0) || (busNr > 1) {
		return fmt.Errorf("Bus number %d out of range", busNr)
	}
	return nil
}

func (c *Adaptor) getPinTranslatorFunction() func(string) (string, int, error) {
	return func(pin string) (chip string, line int, err error) {
		if val, ok := pins[pin][c.readRevision()]; ok {
			line = val
		} else if val, ok := pins[pin]["*"]; ok {
			line = val
		} else {
			err = errors.New("Not a valid pin")
			return
		}
		// TODO: Pi1 model B has only this single "gpiochip0", a change of the translator is needed,
		// to support different chips with different revisions
		return "gpiochip0", line, nil
	}
}

func (c *Adaptor) readRevision() string {
	if c.revision == "" {
		c.revision = "0"
		content, err := c.sys.ReadFile(infoFile)
		if err != nil {
			return c.revision
		}
		for _, v := range strings.Split(string(content), "\n") {
			if strings.Contains(v, "Revision") {
				s := strings.Split(string(v), " ")
				version, _ := strconv.ParseInt("0x"+s[len(s)-1], 0, 64)
				if version <= 3 {
					c.revision = "1"
				} else if version <= 15 {
					c.revision = "2"
				} else {
					c.revision = "3"
				}
			}
		}
	}

	return c.revision
}

func (c *Adaptor) pwmPin(id string) (gobot.PWMPinner, error) {
	pin := c.pwmPins[id]

	if pin == nil {
		tf := c.getPinTranslatorFunction()
		_, i, err := tf(id)
		if err != nil {
			return nil, err
		}
		pin = NewPWMPin(c.sys, "/dev/pi-blaster", strconv.Itoa(i))
		pin.SetPeriod(c.PiBlasterPeriod)
		c.pwmPins[id] = pin
	}

	return pin, nil
}
