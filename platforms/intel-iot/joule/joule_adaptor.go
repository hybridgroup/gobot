package joule

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/sysfs"
)

func writeFile(path string, data []byte) (i int, err error) {
	file, err := sysfs.OpenFile(path, os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return
	}

	return file.Write(data)
}

func readFile(path string) ([]byte, error) {
	file, err := sysfs.OpenFile(path, os.O_RDONLY, 0644)
	defer file.Close()
	if err != nil {
		return make([]byte, 0), err
	}

	buf := make([]byte, 200)
	var i int
	i, err = file.Read(buf)
	if i == 0 {
		return buf, err
	}
	return buf[:i], err
}

type sysfsPin struct {
	pin    int
	pwmPin int
}

// Adaptor represents an Intel Joule
type Adaptor struct {
	name        string
	digitalPins map[int]sysfs.DigitalPin
	pwmPins     map[int]*pwmPin
	i2cBuses    [3]sysfs.I2cDevice
	connect     func(e *Adaptor) (err error)
}

// NewAdaptor returns a new Joule Adaptor
func NewAdaptor() *Adaptor {
	return &Adaptor{
		name: gobot.DefaultName("Joule"),
		connect: func(e *Adaptor) (err error) {
			return
		},
	}
}

// Name returns the Adaptors name
func (e *Adaptor) Name() string { return e.name }

// SetName sets the Adaptors name
func (e *Adaptor) SetName(n string) { e.name = n }

// Connect initializes the Joule for use with the Arduino beakout board
func (e *Adaptor) Connect() (err error) {
	e.digitalPins = make(map[int]sysfs.DigitalPin)
	e.pwmPins = make(map[int]*pwmPin)
	err = e.connect(e)
	return
}

// Finalize releases all i2c devices and exported digital and pwm pins.
func (e *Adaptor) Finalize() (err error) {
	for _, pin := range e.digitalPins {
		if pin != nil {
			if errs := pin.Unexport(); errs != nil {
				err = multierror.Append(err, errs)
			}
		}
	}
	for _, pin := range e.pwmPins {
		if pin != nil {
			if errs := pin.enable("0"); errs != nil {
				err = multierror.Append(err, errs)
			}
			if errs := pin.unexport(); errs != nil {
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

// digitalPin returns matched digitalPin for specified values
func (e *Adaptor) digitalPin(pin string, dir string) (sysfsPin sysfs.DigitalPin, err error) {
	i := sysfsPinMap[pin]
	if e.digitalPins[i.pin] == nil {
		e.digitalPins[i.pin] = sysfs.NewDigitalPin(i.pin)
		if err = e.digitalPins[i.pin].Export(); err != nil {
			// TODO: log error
			return
		}
	}

	if dir == "in" {
		if err = e.digitalPins[i.pin].Direction(sysfs.IN); err != nil {
			return
		}
	} else if dir == "out" {
		if err = e.digitalPins[i.pin].Direction(sysfs.OUT); err != nil {
			return
		}
	}
	return e.digitalPins[i.pin], nil
}

// DigitalRead reads digital value from pin
func (e *Adaptor) DigitalRead(pin string) (i int, err error) {
	sysfsPin, err := e.digitalPin(pin, "in")
	if err != nil {
		return
	}
	return sysfsPin.Read()
}

// DigitalWrite writes a value to the pin. Acceptable values are 1 or 0.
func (e *Adaptor) DigitalWrite(pin string, val byte) (err error) {
	sysfsPin, err := e.digitalPin(pin, "out")
	if err != nil {
		return
	}
	return sysfsPin.Write(int(val))
}

// PwmWrite writes the 0-254 value to the specified pin
func (e *Adaptor) PwmWrite(pin string, val byte) (err error) {
	sysPin := sysfsPinMap[pin]
	if sysPin.pwmPin != -1 {
		if e.pwmPins[sysPin.pwmPin] == nil {
			if err = e.DigitalWrite(pin, 1); err != nil {
				return
			}
			e.pwmPins[sysPin.pwmPin] = newPwmPin(sysPin.pwmPin)
			if err = e.pwmPins[sysPin.pwmPin].export(); err != nil {
				return
			}
			if err = e.pwmPins[sysPin.pwmPin].enable("1"); err != nil {
				return
			}
		}
		p, err := e.pwmPins[sysPin.pwmPin].period()
		if err != nil {
			return err
		}
		period, err := strconv.Atoi(p)
		if err != nil {
			return err
		}
		duty := gobot.FromScale(float64(val), 0, 255.0)
		return e.pwmPins[sysPin.pwmPin].writeDuty(strconv.Itoa(int(float64(period) * duty)))
	}
	return errors.New("Not a PWM pin")
}

// GetConnection returns an i2c connection to a device on a specified bus.
// Valid bus number is [0..2] which corresponds to /dev/i2c-0 through /dev/i2c-2.
func (e *Adaptor) GetConnection(address int, bus int) (connection i2c.Connection, err error) {
	if (bus < 0) || (bus > 2) {
		return nil, fmt.Errorf("Bus number %d out of range", bus)
	}
	if e.i2cBuses[bus] == nil {
		e.i2cBuses[bus], err = sysfs.NewI2cDevice(fmt.Sprintf("/dev/i2c-%d", bus))
	}
	return i2c.NewConnection(e.i2cBuses[bus], address), err
}

// GetDefaultBus returns the default i2c bus for this platform
func (e *Adaptor) GetDefaultBus() int {
	return 0
}
