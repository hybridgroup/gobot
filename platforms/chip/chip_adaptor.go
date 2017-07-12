package chip

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/sysfs"
)

type sysfsPin struct {
	pin    int
	pwmPin int
}

// Adaptor represents a Gobot Adaptor for a C.H.I.P.
type Adaptor struct {
	name        string
	board       string
	pinmap      map[string]sysfsPin
	digitalPins map[int]*sysfs.DigitalPin
	pwmPins     map[int]*sysfs.PWMPin
	i2cBuses    [3]i2c.I2cDevice
	mutex       *sync.Mutex
}

// NewAdaptor creates a C.H.I.P. Adaptor
func NewAdaptor() *Adaptor {
	c := &Adaptor{
		name:  gobot.DefaultName("CHIP"),
		board: "chip",
		mutex: &sync.Mutex{},
	}

	c.setPins()
	return c
}

// NewAdaptor creates a C.H.I.P. Pro Adaptor
func NewProAdaptor() *Adaptor {
	c := &Adaptor{
		name:  gobot.DefaultName("CHIP Pro"),
		board: "pro",
		mutex: &sync.Mutex{},
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
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, pin := range c.digitalPins {
		if pin != nil {
			if e := pin.Unexport(); e != nil {
				err = multierror.Append(err, e)
			}
		}
	}
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
	return
}

// DigitalRead reads digital value from the specified pin.
// Valids pins are the XIO-P0 through XIO-P7 pins from the
// extender (pins 13-20 on header 14), as well as the SoC pins
// aka all the other pins.
func (c *Adaptor) DigitalRead(pin string) (val int, err error) {
	sysfsPin, err := c.DigitalPin(pin, sysfs.IN)
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
	sysfsPin, err := c.DigitalPin(pin, sysfs.OUT)
	if err != nil {
		return err
	}
	return sysfsPin.Write(int(val))
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
		c.i2cBuses[bus], err = sysfs.NewI2cDevice(fmt.Sprintf("/dev/i2c-%d", bus))
	}
	return i2c.NewConnection(c.i2cBuses[bus], address), err
}

// GetDefaultBus returns the default i2c bus for this platform
func (c *Adaptor) GetDefaultBus() int {
	return 1
}

// digitalPin returns matched digitalPin for specified values
func (c *Adaptor) DigitalPin(pin string, dir string) (sysfsPin sysfs.DigitalPinner, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

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

// pwmPin returns matched pwmPin for specified pin number
func (c *Adaptor) PWMPin(pin string) (sysfsPin sysfs.PWMPinner, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	sysPin := c.pinmap[pin]
	if sysPin.pwmPin != -1 {
		if c.pwmPins[sysPin.pwmPin] == nil {
			newPin := sysfs.NewPWMPin(sysPin.pwmPin)
			if err = newPin.Export(); err != nil {
				return
			}
			// Make sure pwm is disabled when setting polarity
			if err = newPin.Enable(false); err != nil {
				return
			}
			if err = newPin.InvertPolarity(false); err != nil {
				return
			}
			if err = newPin.Enable(true); err != nil {
				return
			}
			if err = newPin.SetPeriod(10000000); err != nil {
				return
			}
			c.pwmPins[sysPin.pwmPin] = newPin
		}

		sysfsPin = c.pwmPins[sysPin.pwmPin]
		return
	}
	err = errors.New("Not a PWM pin")
	return
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
	//
	// Duty cycle is in nanos
	const minDuty = 0.0005 * 1e9
	const maxDuty = 0.0020 * 1e9
	duty := uint32(gobot.ToScale(gobot.FromScale(float64(angle), 0, 180), minDuty, maxDuty))
	return pwmPin.SetDutyCycle(duty)
}

// SetBoard sets the name of the type of board
func (c *Adaptor) SetBoard(n string) (err error) {
	if n == "chip" || n == "pro" {
		c.board = n
		c.setPins()
		return
	}
	return errors.New("Invalid board type")
}

func (c *Adaptor) setPins() {
	c.digitalPins = make(map[int]*sysfs.DigitalPin)
	c.pwmPins = make(map[int]*sysfs.PWMPin)

	if c.board == "pro" {
		c.pinmap = chipProPins
		return
	}
	// otherwise, original CHIP
	c.pinmap = chipPins
	baseAddr, _ := getXIOBase()
	for i := 0; i < 8; i++ {
		pin := fmt.Sprintf("XIO-P%d", i)
		c.pinmap[pin] = sysfsPin{pin: baseAddr + i, pwmPin: -1}
	}
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

func (c *Adaptor) translatePin(pin string) (i int, err error) {
	if val, ok := c.pinmap[pin]; ok {
		i = val.pin
	} else {
		err = errors.New("Not a valid pin")
	}
	return
}
