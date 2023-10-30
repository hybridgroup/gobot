package raspi

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/adaptors"
	"gobot.io/x/gobot/v2/system"
)

const (
	infoFile = "/proc/cpuinfo"

	defaultSpiBusNumber  = 0
	defaultSpiChipNumber = 0
	defaultSpiMode       = 0
	defaultSpiBitsNumber = 8
	defaultSpiMaxSpeed   = 500000
)

// Adaptor is the Gobot Adaptor for the Raspberry Pi
type Adaptor struct {
	name     string
	mutex    sync.Mutex
	sys      *system.Accesser
	revision string
	pwmPins  map[string]gobot.PWMPinner
	*adaptors.DigitalPinsAdaptor
	*adaptors.I2cBusAdaptor
	*adaptors.SpiBusAdaptor
	PiBlasterPeriod uint32
}

// NewAdaptor creates a Raspi Adaptor
//
// Optional parameters:
//
//			adaptors.WithGpiodAccess():	use character device gpiod driver instead of sysfs (still used by default)
//			adaptors.WithSpiGpioAccess(sclk, nss, mosi, miso):	use GPIO's instead of /dev/spidev#.#
//	   adaptors.WithGpiosActiveLow(pin's): invert the pin behavior
//	   adaptors.WithGpiosPullUp/Down(pin's): sets the internal pull resistor
//	   adaptors.WithGpiosOpenDrain/Source(pin's): sets the output behavior
//	   adaptors.WithGpioDebounce(pin, period): sets the input debouncer
//	   adaptors.WithGpioEventOnFallingEdge/RaisingEdge/BothEdges(pin, handler): activate edge detection
func NewAdaptor(opts ...func(adaptors.Optioner)) *Adaptor {
	sys := system.NewAccesser(system.WithDigitalPinGpiodAccess())
	c := &Adaptor{
		name:            gobot.DefaultName("RaspberryPi"),
		sys:             sys,
		PiBlasterPeriod: 10000000,
	}
	c.DigitalPinsAdaptor = adaptors.NewDigitalPinsAdaptor(sys, c.getPinTranslatorFunction(), opts...)
	c.I2cBusAdaptor = adaptors.NewI2cBusAdaptor(sys, c.validateI2cBusNumber, 1)
	c.SpiBusAdaptor = adaptors.NewSpiBusAdaptor(sys, c.validateSpiBusNumber, defaultSpiBusNumber, defaultSpiChipNumber,
		defaultSpiMode, defaultSpiBitsNumber, defaultSpiMaxSpeed)
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

	if err := c.SpiBusAdaptor.Connect(); err != nil {
		return err
	}

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

	if e := c.SpiBusAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}
	return err
}

// DefaultI2cBus returns the default i2c bus for this platform.
// This overrides the base function due to the revision dependency.
func (c *Adaptor) DefaultI2cBus() int {
	rev := c.readRevision()
	if rev == "2" || rev == "3" {
		return 1
	}
	return 0
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

func (c *Adaptor) validateSpiBusNumber(busNr int) error {
	// Valid bus numbers are [0,1] which corresponds to /dev/spidev0.x through /dev/spidev1.x.
	// x is the chip number <255
	if (busNr < 0) || (busNr > 1) {
		return fmt.Errorf("Bus number %d out of range", busNr)
	}
	return nil
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
				s := strings.Split(v, " ")
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
		if err := pin.SetPeriod(c.PiBlasterPeriod); err != nil {
			return nil, err
		}
		c.pwmPins[id] = pin
	}

	return pin, nil
}
