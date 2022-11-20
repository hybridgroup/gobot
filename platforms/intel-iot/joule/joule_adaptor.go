package joule

import (
	"errors"
	"fmt"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/system"
)

type sysfsPin struct {
	pin    int
	pwmPin int
}

// Adaptor represents an Intel Joule
type Adaptor struct {
	name        string
	sys         *system.Accesser
	mutex       sync.Mutex
	digitalPins map[int]gobot.DigitalPinner
	pwmPins     map[int]gobot.PWMPinner
	i2cBuses    [3]i2c.I2cDevice
}

// NewAdaptor returns a new Joule Adaptor
func NewAdaptor() *Adaptor {
	return &Adaptor{
		name: gobot.DefaultName("Joule"),
		sys:  system.NewAccesser(),
	}
}

// Name returns the Adaptors name
func (e *Adaptor) Name() string { return e.name }

// SetName sets the Adaptors name
func (e *Adaptor) SetName(n string) { e.name = n }

// Connect initializes the Joule for use with the Arduino beakout board
func (e *Adaptor) Connect() (err error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.digitalPins = make(map[int]gobot.DigitalPinner)
	e.pwmPins = make(map[int]gobot.PWMPinner)
	return
}

// Finalize releases all i2c devices and exported digital and pwm pins.
func (e *Adaptor) Finalize() (err error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	for _, pin := range e.digitalPins {
		if pin != nil {
			if errs := pin.Unexport(); errs != nil {
				err = multierror.Append(err, errs)
			}
		}
	}
	for _, pin := range e.pwmPins {
		if pin != nil {
			if errs := pin.Enable(false); errs != nil {
				err = multierror.Append(err, errs)
			}
			if errs := pin.Unexport(); errs != nil {
				err = multierror.Append(err, errs)
			}
		}
	}
	for _, bus := range e.i2cBuses {
		if bus != nil {
			if errs := bus.Close(); errs != nil {
				err = multierror.Append(err, errs)
			}
		}
	}
	return
}

// DigitalPin returns matched digitalPin for specified values
func (e *Adaptor) DigitalPin(pin string, dir string) (gobot.DigitalPinner, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	i := sysfsPinMap[pin]
	if e.digitalPins[i.pin] == nil {
		e.digitalPins[i.pin] = e.sys.NewDigitalPin(i.pin)
		if err := e.digitalPins[i.pin].Export(); err != nil {
			return nil, err
		}
	}

	if dir == "in" {
		if err := e.digitalPins[i.pin].Direction(system.IN); err != nil {
			return nil, err
		}
	} else if dir == "out" {
		if err := e.digitalPins[i.pin].Direction(system.OUT); err != nil {
			return nil, err
		}
	}
	return e.digitalPins[i.pin], nil
}

// DigitalRead reads digital value from pin
func (e *Adaptor) DigitalRead(pin string) (i int, err error) {
	sysPin, err := e.DigitalPin(pin, "in")
	if err != nil {
		return
	}
	return sysPin.Read()
}

// DigitalWrite writes a value to the pin. Acceptable values are 1 or 0.
func (e *Adaptor) DigitalWrite(pin string, val byte) (err error) {
	sysPin, err := e.DigitalPin(pin, "out")
	if err != nil {
		return
	}
	return sysPin.Write(int(val))
}

// PwmWrite writes the 0-254 value to the specified pin
func (e *Adaptor) PwmWrite(pin string, val byte) (err error) {
	pwmPin, err := e.PWMPin(pin)
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

// PWMPin returns a sys.PWMPin
func (e *Adaptor) PWMPin(pin string) (gobot.PWMPinner, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	sysPin, ok := sysfsPinMap[pin]
	if !ok {
		return nil, errors.New("Not a valid pin")
	}
	if sysPin.pwmPin != -1 {
		if e.pwmPins[sysPin.pwmPin] == nil {
			e.pwmPins[sysPin.pwmPin] = e.sys.NewPWMPin("/sys/class/pwm/pwmchip0", sysPin.pwmPin)
			if err := e.pwmPins[sysPin.pwmPin].Export(); err != nil {
				return nil, err
			}
			if err := e.pwmPins[sysPin.pwmPin].SetPeriod(10000000); err != nil {
				return nil, err
			}
			if err := e.pwmPins[sysPin.pwmPin].Enable(true); err != nil {
				return nil, err
			}
		}

		return e.pwmPins[sysPin.pwmPin], nil
	}

	return nil, errors.New("Not a PWM pin")
}

// GetConnection returns an i2c connection to a device on a specified bus.
// Valid bus number is [0..2] which corresponds to /dev/i2c-0 through /dev/i2c-2.
func (e *Adaptor) GetConnection(address int, bus int) (connection i2c.Connection, err error) {
	if (bus < 0) || (bus > 2) {
		return nil, fmt.Errorf("Bus number %d out of range", bus)
	}
	if e.i2cBuses[bus] == nil {
		e.i2cBuses[bus], err = e.sys.NewI2cDevice(fmt.Sprintf("/dev/i2c-%d", bus))
	}
	return i2c.NewConnection(e.i2cBuses[bus], address), err
}

// GetDefaultBus returns the default i2c bus for this platform
func (e *Adaptor) GetDefaultBus() int {
	return 0
}
