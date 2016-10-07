package edison

import (
	"errors"
	"os"
	"strconv"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/sysfs"
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

type mux struct {
	pin   int
	value int
}
type sysfsPin struct {
	pin          int
	resistor     int
	levelShifter int
	pwmPin       int
	mux          []mux
}

// Adaptor represents a Gobot Adaptor for an Intel Edison
type Adaptor struct {
	name        string
	board       string
	pinmap      map[string]sysfsPin
	tristate    sysfs.DigitalPin
	digitalPins map[int]sysfs.DigitalPin
	pwmPins     map[int]*pwmPin
	i2cDevice   sysfs.I2cDevice
	connect     func(e *Adaptor) (err error)
}

// changePinMode writes pin mode to current_pinmux file
func changePinMode(pin, mode string) (err error) {
	_, err = writeFile(
		"/sys/kernel/debug/gpio_debug/gpio"+pin+"/current_pinmux",
		[]byte("mode"+mode),
	)
	return
}

// NewAdaptor returns a new Edison Adaptor
func NewAdaptor() *Adaptor {
	return &Adaptor{
		name:   "Edison",
		board:  "arduino",
		pinmap: arduinoPinMap,
	}
}

// Name returns the Adaptors name
func (e *Adaptor) Name() string { return e.name }

// SetName sets the Adaptors name
func (e *Adaptor) SetName(n string) { e.name = n }

// Board returns the Adaptors board name
func (e *Adaptor) Board() string { return e.board }

// SetBoard sets the Adaptors name
func (e *Adaptor) SetBoard(n string) { e.board = n }

// Connect initializes the Edison for use with the Arduino beakout board
func (e *Adaptor) Connect() (errs []error) {
	e.digitalPins = make(map[int]sysfs.DigitalPin)
	e.pwmPins = make(map[int]*pwmPin)
	if e.board != "arduino" {
		e.pinmap = miniboardPinMap
		if err := e.miniboardSetup(); err != nil {
			return []error{err}
		}
		return
	}
	if err := e.arduinoSetup(); err != nil {
		return []error{err}
	}
	return
}

// Finalize releases all i2c devices and exported analog, digital, pwm pins.
func (e *Adaptor) Finalize() (errs []error) {
	if err := e.tristate.Unexport(); err != nil {
		errs = append(errs, err)
	}
	for _, pin := range e.digitalPins {
		if pin != nil {
			if err := pin.Unexport(); err != nil {
				errs = append(errs, err)
			}
		}
	}
	for _, pin := range e.pwmPins {
		if pin != nil {
			if err := pin.enable("0"); err != nil {
				errs = append(errs, err)
			}
			if err := pin.unexport(); err != nil {
				errs = append(errs, err)
			}
		}
	}
	if e.i2cDevice != nil {
		if err := e.i2cDevice.Close(); errs != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// arduinoSetup does needed setup for the Arduino compatible breakout board
func (e *Adaptor) arduinoSetup() (err error) {
	e.tristate = sysfs.NewDigitalPin(214)
	if err = e.tristate.Export(); err != nil {
		return err
	}
	if err = e.tristate.Direction(sysfs.OUT); err != nil {
		return err
	}
	if err = e.tristate.Write(sysfs.LOW); err != nil {
		return err
	}

	for _, i := range []int{263, 262} {
		io := sysfs.NewDigitalPin(i)
		if err = io.Export(); err != nil {
			return err
		}
		if err = io.Direction(sysfs.OUT); err != nil {
			return err
		}
		if err = io.Write(sysfs.HIGH); err != nil {
			return err
		}
		if err = io.Unexport(); err != nil {
			return err
		}
	}

	for _, i := range []int{240, 241, 242, 243} {
		io := sysfs.NewDigitalPin(i)
		if err = io.Export(); err != nil {
			return err
		}
		if err = io.Direction(sysfs.OUT); err != nil {
			return err
		}
		if err = io.Write(sysfs.LOW); err != nil {
			return err
		}
		if err = io.Unexport(); err != nil {
			return err
		}

	}

	for _, i := range []string{"111", "115", "114", "109"} {
		if err = changePinMode(i, "1"); err != nil {
			return err
		}
	}

	for _, i := range []string{"131", "129", "40"} {
		if err = changePinMode(i, "0"); err != nil {
			return err
		}
	}

	err = e.tristate.Write(sysfs.HIGH)
	return
}

func (e *Adaptor) arduinoI2CSetup() (err error) {
	if err = e.tristate.Write(sysfs.LOW); err != nil {
		return
	}

	for _, i := range []int{14, 165, 212, 213} {
		io := sysfs.NewDigitalPin(i)
		if err = io.Export(); err != nil {
			return
		}
		if err = io.Direction(sysfs.IN); err != nil {
			return
		}
		if err = io.Unexport(); err != nil {
			return
		}
	}

	for _, i := range []int{236, 237, 204, 205} {
		io := sysfs.NewDigitalPin(i)
		if err = io.Export(); err != nil {
			return
		}
		if err = io.Direction(sysfs.OUT); err != nil {
			return
		}
		if err = io.Write(sysfs.LOW); err != nil {
			return
		}
		if err = io.Unexport(); err != nil {
			return
		}
	}

	for _, i := range []string{"28", "27"} {
		if err = changePinMode(i, "1"); err != nil {
			return
		}
	}

	if err = e.tristate.Write(sysfs.HIGH); err != nil {
		return
	}

	return
}

// miniboardSetup does needed setup for Edison minibpard and other compatible boards
func (e *Adaptor) miniboardSetup() (err error) {
	return
}

// digitalPin returns matched digitalPin for specified values
func (e *Adaptor) digitalPin(pin string, dir string) (sysfsPin sysfs.DigitalPin, err error) {
	i := e.pinmap[pin]
	if e.digitalPins[i.pin] == nil {
		e.digitalPins[i.pin] = sysfs.NewDigitalPin(i.pin)
		if err = e.digitalPins[i.pin].Export(); err != nil {
			return
		}

		if i.resistor > 0 {
			e.digitalPins[i.resistor] = sysfs.NewDigitalPin(i.resistor)
			if err = e.digitalPins[i.resistor].Export(); err != nil {
				return
			}
		}

		if i.levelShifter > 0 {
			e.digitalPins[i.levelShifter] = sysfs.NewDigitalPin(i.levelShifter)
			if err = e.digitalPins[i.levelShifter].Export(); err != nil {
				return
			}
		}

		if len(i.mux) > 0 {
			for _, mux := range i.mux {
				e.digitalPins[mux.pin] = sysfs.NewDigitalPin(mux.pin)
				if err = e.digitalPins[mux.pin].Export(); err != nil {
					return
				}

				if err = e.digitalPins[mux.pin].Direction(sysfs.OUT); err != nil {
					return
				}

				if err = e.digitalPins[mux.pin].Write(mux.value); err != nil {
					return
				}

			}
		}
	}

	if dir == "in" {
		if err = e.digitalPins[i.pin].Direction(sysfs.IN); err != nil {
			return
		}

		if i.resistor > 0 {
			if err = e.digitalPins[i.resistor].Direction(sysfs.OUT); err != nil {
				return
			}

			if err = e.digitalPins[i.resistor].Write(sysfs.LOW); err != nil {
				return
			}
		}

		if i.levelShifter > 0 {
			if err = e.digitalPins[i.levelShifter].Direction(sysfs.OUT); err != nil {
				return
			}

			if err = e.digitalPins[i.levelShifter].Write(sysfs.LOW); err != nil {
				return
			}
		}
	} else if dir == "out" {
		if err = e.digitalPins[i.pin].Direction(sysfs.OUT); err != nil {
			return
		}

		if i.resistor > 0 {
			if err = e.digitalPins[i.resistor].Direction(sysfs.IN); err != nil {
				return
			}
		}

		if i.levelShifter > 0 {
			if err = e.digitalPins[i.levelShifter].Direction(sysfs.OUT); err != nil {
				return
			}

			if err = e.digitalPins[i.levelShifter].Write(sysfs.HIGH); err != nil {
				return
			}
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
	sysPin := e.pinmap[pin]
	if sysPin.pwmPin != -1 {
		if e.pwmPins[sysPin.pwmPin] == nil {
			if err = e.DigitalWrite(pin, 1); err != nil {
				return
			}
			if err = changePinMode(strconv.Itoa(int(sysPin.pin)), "1"); err != nil {
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

// AnalogRead returns value from analog reading of specified pin
func (e *Adaptor) AnalogRead(pin string) (val int, err error) {
	buf, err := readFile(
		"/sys/bus/iio/devices/iio:device1/in_voltage" + pin + "_raw",
	)
	if err != nil {
		return
	}

	val, err = strconv.Atoi(string(buf[0 : len(buf)-1]))

	return val / 4, err
}

// I2cStart initializes i2c device for addresss
func (e *Adaptor) I2cStart(address int) (err error) {
	if e.i2cDevice != nil {
		return
	}

	// most board use I2C bus 1
	bus := "/dev/i2c-1"

	// except for Arduino which uses bus 6
	if e.board == "arduino" {
		bus = "/dev/i2c-6"
		e.arduinoI2CSetup()
	}

	e.i2cDevice, err = sysfs.NewI2cDevice(bus, address)
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
