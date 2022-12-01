package tinkerboard

import (
	"fmt"
	"log"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/adaptors"
	"gobot.io/x/gobot/system"
)

const debug = false

const (
	pwmNormal        = "normal"
	pwmInverted      = "inversed"
	pwmPeriodDefault = 10000000 // 10ms = 100Hz
)

type pwmPinDefinition struct {
	channel   int
	dir       string
	dirRegexp string
}

// Adaptor represents a Gobot Adaptor for the ASUS Tinker Board
type Adaptor struct {
	name  string
	sys   *system.Accesser
	mutex sync.Mutex
	*adaptors.DigitalPinsAdaptor
	pwmPins  map[string]gobot.PWMPinner
	i2cBuses [5]i2c.I2cDevice
}

// NewAdaptor creates a Tinkerboard Adaptor
func NewAdaptor() *Adaptor {
	sys := system.NewAccesser("cdev")
	c := &Adaptor{
		name:    gobot.DefaultName("Tinker Board"),
		sys:     sys,
		pwmPins: make(map[string]gobot.PWMPinner),
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
	c.mutex.Lock()
	defer c.mutex.Unlock()

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
	return err
}

// PwmWrite writes a PWM signal to the specified pin.
func (c *Adaptor) PwmWrite(pin string, val byte) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	pwmPin, err := c.pwmPin(pin)
	if err != nil {
		return
	}
	period, err := pwmPin.Period()
	if err != nil {
		return err
	}
	duty := gobot.FromScale(float64(val), 0, 255.0)
	if debug {
		log.Printf("Tinkerboard PwmWrite - raw: %d, period: %d, duty: %.2f %%", val, period, duty*100)
	}
	return pwmPin.SetDutyCycle(uint32(float64(period) * duty))
}

// ServoWrite writes a servo signal to the specified pin.
func (c *Adaptor) ServoWrite(pin string, angle byte) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	pwmPin, err := c.pwmPin(pin)
	if err != nil {
		return
	}
	period, err := pwmPin.Period()
	if err != nil {
		return err
	}

	// 0.5 ms => -90
	// 1.5 ms =>   0
	// 2.0 ms =>  90
	minDuty := 100 * 0.0005 * float64(period)
	maxDuty := 100 * 0.0020 * float64(period)
	duty := uint32(gobot.ToScale(gobot.FromScale(float64(angle), 0, 180), minDuty, maxDuty))
	return pwmPin.SetDutyCycle(duty)
}

// SetPeriod adjusts the period of the specified PWM pin.
// If duty cycle is already set, also this value will be adjusted in the same ratio.
func (c *Adaptor) SetPeriod(pin string, period uint32) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	pwmPin, err := c.pwmPin(pin)
	if err != nil {
		return err
	}
	return setPeriod(pwmPin, period)
}

// PWMPin initializes the pin for PWM and returns matched pwmPin for specified pin number.
// It implements the PWMPinnerProvider interface.
func (c *Adaptor) PWMPin(pin string) (gobot.PWMPinner, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.pwmPin(pin)
}

// GetConnection returns a connection to a device on a specified i2c bus.
// Valid bus number is [0..4] which corresponds to /dev/i2c-0 through /dev/i2c-4.
// We don't support "/dev/i2c-6 DesignWare HDMI".
func (c *Adaptor) GetConnection(address int, bus int) (connection i2c.Connection, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if (bus < 0) || (bus > 4) {
		return nil, fmt.Errorf("Bus number %d out of range", bus)
	}
	if c.i2cBuses[bus] == nil {
		c.i2cBuses[bus], err = c.sys.NewI2cDevice(fmt.Sprintf("/dev/i2c-%d", bus))
	}
	return i2c.NewConnection(c.i2cBuses[bus], address), err
}

// GetDefaultBus returns the default i2c bus for this platform.
func (c *Adaptor) GetDefaultBus() int {
	return 1
}

// pwmPin initializes the pin for PWM and returns matched pwmPin for specified pin number.
func (c *Adaptor) pwmPin(pin string) (gobot.PWMPinner, error) {
	var pwmPinData pwmPinDefinition
	pwmPinData, err := c.translatePwmPin(pin)
	if err != nil {
		return nil, err
	}

	if c.pwmPins[pin] == nil {
		var path string
		path, err := pwmPinData.findPwmDir(*c.sys)
		if err != nil {
			return nil, err
		}
		newPin := c.sys.NewPWMPin(path, pwmPinData.channel)
		if err := newPin.Export(); err != nil {
			return nil, err
		}
		// Make sure pwm is disabled before change anything
		if err := newPin.Enable(false); err != nil {
			return nil, err
		}
		if err := setPeriod(newPin, pwmPeriodDefault); err != nil {
			return nil, err
		}
		if err := newPin.SetPolarity(pwmNormal); err != nil {
			return nil, err
		}
		if err := newPin.Enable(true); err != nil {
			return nil, err
		}
		if debug {
			log.Printf("New PWMPin created for %s\n", pin)
		}
		c.pwmPins[pin] = newPin
	}

	return c.pwmPins[pin], nil
}

// setPeriod adjusts the PWM period of the given pin.
// If duty cycle is already set, also this value will be adjusted in the same ratio.
// The order in which the values are written must be observed, otherwise an error occur "write error: Invalid argument".
func setPeriod(pwmPin gobot.PWMPinner, period uint32) error {
	var errorBase = fmt.Sprintf("tinkerboard.setPeriod(%v, %d) failed", pwmPin, period)
	oldDuty, err := pwmPin.DutyCycle()
	if err != nil {
		return fmt.Errorf("%s with '%v'", errorBase, err)
	}

	if oldDuty == 0 {
		if err := pwmPin.SetPeriod(period); err != nil {
			log.Println(1, period)
			return fmt.Errorf("%s with '%v'", errorBase, err)
		}
	} else {
		// adjust duty cycle in the same ratio
		oldPeriod, err := pwmPin.Period()
		if err != nil {
			return fmt.Errorf("%s with '%v'", errorBase, err)
		}
		duty := uint32(uint64(oldDuty) * uint64(period) / uint64(oldPeriod))
		if debug {
			log.Printf("oldPeriod: %d, oldDuty: %d, new period: %d, new duty: %d", oldPeriod, oldDuty, period, duty)
		}

		// the order depends on value (duty must not be bigger than period in any situation)
		if duty > oldPeriod {
			if err := pwmPin.SetPeriod(period); err != nil {
				log.Println(2, period)
				return fmt.Errorf("%s with '%v'", errorBase, err)
			}
			if err := pwmPin.SetDutyCycle(uint32(duty)); err != nil {
				log.Println(2, duty)
				return fmt.Errorf("%s with '%v'", errorBase, err)
			}
		} else {
			if err := pwmPin.SetDutyCycle(uint32(duty)); err != nil {
				log.Println(3, duty)
				return fmt.Errorf("%s with '%v'", errorBase, err)
			}
			if err := pwmPin.SetPeriod(period); err != nil {
				log.Println(3, period)
				return fmt.Errorf("%s with '%v'", errorBase, err)
			}
		}
	}
	return nil
}

func (c *Adaptor) translateDigitalPin(id string) (string, int, error) {
	pindef, ok := gpioPinDefinitions[id]
	if !ok {
		return "", -1, fmt.Errorf("'%s' is not a valid id for a digital pin", id)
	}
	if c.sys.IsSysfsDigitalPinAccess() {
		return "", pindef.sysfs, nil
	}
	chip := fmt.Sprintf("gpiochip%d", pindef.cdev.chip)
	line := int(pindef.cdev.line)
	return chip, line, nil
}

func (c *Adaptor) translatePwmPin(pin string) (pwmPin pwmPinDefinition, err error) {
	var ok bool
	if pwmPin, ok = pwmPinDefinitions[pin]; !ok {
		err = fmt.Errorf("Not a valid PWM pin")
	}
	return
}

func (p pwmPinDefinition) findPwmDir(sys system.Accesser) (dir string, err error) {
	items, _ := sys.Find(p.dir, p.dirRegexp)
	if items == nil || len(items) == 0 {
		return "", fmt.Errorf("No path found for PWM directory pattern, '%s' in path '%s'. See README.md for activation", p.dirRegexp, p.dir)
	}

	dir = items[0]
	info, err := sys.Stat(dir)
	if err != nil {
		return "", fmt.Errorf("Error (%v) on access '%s'", err, dir)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("The item '%s' is not a directory, which is not expected", dir)
	}

	return
}
