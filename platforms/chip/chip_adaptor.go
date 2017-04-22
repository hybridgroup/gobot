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

type sysfsPin struct {
	pin    int
	pwmPin int
}

// Adaptor represents a Gobot Adaptor for a C.H.I.P.
type Adaptor struct {
	name        string
	digitalPins map[int]sysfs.DigitalPin
	pinmap      map[string]sysfsPin
	i2cBuses    [3]sysfs.I2cDevice
	pwmPins     map[int]*sysfs.PWMPin
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
	c.pwmPins = make(map[int]*sysfs.PWMPin)
	return nil
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

func (c *Adaptor) setPins() {
	c.pinmap = fixedPins
	baseAddr, _ := getXIOBase()
	for i := 0; i < 8; i++ {
		pin := fmt.Sprintf("XIO-P%d", i)
		c.pinmap[pin] = sysfsPin{pin: baseAddr + i, pwmPin: -1}
	}
}

func (c *Adaptor) translatePin(pin string) (i int, err error) {
	if val, ok := c.pinmap[pin]; ok {
		i = val.pin
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

// PwmWrite writes a PWM signal to the specified pin
func (c *Adaptor) PwmWrite(pin string, val byte) (err error) {
	sysPin := c.pinmap[pin]
	if sysPin.pwmPin != -1 {
		if c.pwmPins[sysPin.pwmPin] == nil {
			c.pwmPins[sysPin.pwmPin] = sysfs.NewPWMPin(sysPin.pwmPin)
			if err = c.pwmPins[sysPin.pwmPin].Export(); err != nil {
				return
			}
			if err = c.pwmPins[sysPin.pwmPin].Enable(true); err != nil {
				return
			}
		}
		p, err := c.pwmPins[sysPin.pwmPin].Period()
		if err != nil {
			return err
		}
		period, err := strconv.Atoi(p)
		if err != nil {
			return err
		}
		duty := gobot.FromScale(float64(val), 0, 255.0)
		return c.pwmPins[sysPin.pwmPin].SetDutyCycle(uint32(float64(period) * duty))
	}
	return errors.New("Not a PWM pin")
}

const pwmFrequency = 100

// ServoWrite writes a servo signal to the specified pin
// func (c *Adaptor) ServoWrite(pin string, angle byte) (err error) {
// 	if pin != "PWM0" {
// 		return fmt.Errorf("Servo is only available on pin PWM0")
// 	}
// 	if c.pwm == nil {
// 		err = c.initPWM(pwmFrequency)
// 		if err != nil {
// 			return
// 		}
// 	}
// 	// 0.5 ms => -90
// 	// 1.5 ms =>   0
// 	// 2.0 ms =>  90
// 	const minDuty = 100 * 0.0005 * pwmFrequency
// 	const maxDuty = 100 * 0.0020 * pwmFrequency
// 	duty := gobot.ToScale(gobot.FromScale(float64(angle), 0, 180), minDuty, maxDuty)
// 	return c.pwm.setDutycycle(duty)
// }

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
