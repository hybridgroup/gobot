package joule

import (
	"errors"
	"os"
	"strconv"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot"
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
	var i = 0
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
	i2cDevice   sysfs.I2cDevice
	connect     func(e *Adaptor) (err error)
}

// NewAdaptor returns a new Joule Adaptor
func NewAdaptor() *Adaptor {
	return &Adaptor{
		name: "Joule",
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
	if e.i2cDevice != nil {
		if errs := e.i2cDevice.Close(); errs != nil {
			err = multierror.Append(err, errs)
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

// I2cStart initializes i2c device for addresss
func (e *Adaptor) I2cStart(address int) (err error) {
	if e.i2cDevice != nil {
		return
	}

	// TODO: handle the additional I2C buses
	e.i2cDevice, err = sysfs.NewI2cDevice("/dev/i2c-0", address)
	return
}

// I2cWrite writes data to i2c device
func (e *Adaptor) I2cWrite(address int, data []byte) (err error) {
	if err = e.i2cDevice.SetAddress(address); err != nil {
		return err
	}
	_, err = e.i2cDevice.Write(data)
	return
}

// I2cRead returns size bytes from the i2c device
func (e *Adaptor) I2cRead(address int, size int) (data []byte, err error) {
	data = make([]byte, size)
	if err = e.i2cDevice.SetAddress(address); err != nil {
		return
	}
	_, err = e.i2cDevice.Read(data)
	return
}
