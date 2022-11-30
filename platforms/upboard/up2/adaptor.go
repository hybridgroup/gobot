package up2

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
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
)

// TODO: take into account the actual period setting, not just assume default
const pwmPeriod = 10000000

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
	pwmPins  map[int]gobot.PWMPinner
	i2cBuses [6]i2c.I2cDevice
	spiBuses [2]spi.Connection
}

// NewAdaptor creates a UP2 Adaptor
func NewAdaptor() *Adaptor {
	sys := system.NewAccesser()
	c := &Adaptor{
		name:    gobot.DefaultName("UP2"),
		sys:     sys,
		ledPath: "/sys/class/leds/upboard:%s:/brightness",
		pwmPins: make(map[int]gobot.PWMPinner),
		pinmap:  fixedPins,
	}
	c.DigitalPinsAdaptor = adaptors.NewDigitalPinsAdaptor(sys, c.translateDigitalPin)
	return c
}

// Name returns the name of the Adaptor
func (c *Adaptor) Name() string { return c.name }

// SetName sets the name of the Adaptor
func (c *Adaptor) SetName(n string) { c.name = n }

// Connect create new connection to board and pins.
func (c *Adaptor) Connect() error {
	err := c.DigitalPinsAdaptor.Connect()
	return err
}

// Finalize closes connection to board and pins
func (c *Adaptor) Finalize() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	err := c.DigitalPinsAdaptor.Finalize()

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

// PWMPin returns matched pwmPin for specified pin number
func (c *Adaptor) PWMPin(pin string) (gobot.PWMPinner, error) {
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
		newPin := c.sys.NewPWMPin("/sys/class/pwm/pwmchip0", i)
		if err = newPin.Export(); err != nil {
			return nil, err
		}
		// Make sure pwm is disabled when setting polarity
		if err = newPin.Enable(false); err != nil {
			return nil, err
		}
		if err = newPin.InvertPolarity(false); err != nil {
			return nil, err
		}
		if err = newPin.Enable(true); err != nil {
			return nil, err
		}
		if err = newPin.SetPeriod(10000000); err != nil {
			return nil, err
		}
		c.pwmPins[i] = newPin
	}

	return c.pwmPins[i], nil
}

// GetConnection returns a connection to a device on a specified bus.
// Valid bus number is [5..6] which corresponds to /dev/i2c-5 through /dev/i2c-6.
func (c *Adaptor) GetConnection(address int, bus int) (connection i2c.Connection, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if (bus < 5) || (bus > 6) {
		return nil, fmt.Errorf("Bus number %d out of range", bus)
	}
	if c.i2cBuses[bus] == nil {
		c.i2cBuses[bus], err = c.sys.NewI2cDevice(fmt.Sprintf("/dev/i2c-%d", bus))
	}
	return i2c.NewConnection(c.i2cBuses[bus], address), err
}

// GetDefaultBus returns the default i2c bus for this platform
func (c *Adaptor) GetDefaultBus() int {
	return 5
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

func (c *Adaptor) translateDigitalPin(id string) (string, int, error) {
	if val, ok := c.pinmap[id]; ok {
		return "", val.pin, nil
	}
	return "", -1, fmt.Errorf("'%s' is not a valid id for a digital pin", id)
}

func (c *Adaptor) translatePwmPin(pin string) (i int, err error) {
	if val, ok := c.pinmap[pin]; ok {
		i = val.pwmPin
	} else {
		err = errors.New("Not a valid pin")
	}
	return
}
