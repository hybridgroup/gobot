package edison

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/sysfs"
)

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
	tristate    *sysfs.DigitalPin
	digitalPins map[int]*sysfs.DigitalPin
	pwmPins     map[int]*sysfs.PWMPin
	i2cBus      i2c.I2cDevice
	connect     func(e *Adaptor) (err error)
	writeFile   func(path string, data []byte) (i int, err error)
	readFile    func(path string) ([]byte, error)
	mutex       *sync.Mutex
}

// NewAdaptor returns a new Edison Adaptor
func NewAdaptor() *Adaptor {
	return &Adaptor{
		name:      gobot.DefaultName("Edison"),
		pinmap:    arduinoPinMap,
		writeFile: writeFile,
		readFile:  readFile,
		mutex:     &sync.Mutex{},
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
func (e *Adaptor) Connect() (err error) {
	e.digitalPins = make(map[int]*sysfs.DigitalPin)
	e.pwmPins = make(map[int]*sysfs.PWMPin)

	if e.Board() == "arduino" || e.Board() == "" {
		aerr := e.checkForArduino()
		if aerr != nil {
			return aerr
		}
		e.board = "arduino"
	}

	switch e.Board() {
	case "sparkfun":
		e.pinmap = sparkfunPinMap
		e.sparkfunSetup()
	case "arduino":
		e.pinmap = arduinoPinMap
		if errs := e.arduinoSetup(); errs != nil {
			err = multierror.Append(err, errs)
		}
	case "miniboard":
		e.pinmap = miniboardPinMap
		e.miniboardSetup()
	default:
		errs := errors.New("Unknown board type: " + e.Board())
		err = multierror.Append(err, errs)
	}
	return
}

// Finalize releases all i2c devices and exported analog, digital, pwm pins.
func (e *Adaptor) Finalize() (err error) {
	if errs := e.tristate.Unexport(); errs != nil {
		err = multierror.Append(err, errs)
	}
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
	if e.i2cBus != nil {
		if errs := e.i2cBus.Close(); errs != nil {
			err = multierror.Append(err, errs)
		}
	}
	return
}

// DigitalRead reads digital value from pin
func (e *Adaptor) DigitalRead(pin string) (i int, err error) {
	sysfsPin, err := e.DigitalPin(pin, "in")
	if err != nil {
		return
	}
	return sysfsPin.Read()
}

// DigitalWrite writes a value to the pin. Acceptable values are 1 or 0.
func (e *Adaptor) DigitalWrite(pin string, val byte) (err error) {
	sysfsPin, err := e.DigitalPin(pin, "out")
	if err != nil {
		return
	}
	return sysfsPin.Write(int(val))
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

// AnalogRead returns value from analog reading of specified pin
func (e *Adaptor) AnalogRead(pin string) (val int, err error) {
	buf, err := e.readFile(
		"/sys/bus/iio/devices/iio:device1/in_voltage" + pin + "_raw",
	)
	if err != nil {
		return
	}

	val, err = strconv.Atoi(string(buf[0 : len(buf)-1]))

	return val / 4, err
}

// GetConnection returns an i2c connection to a device on a specified bus.
// Valid bus numbers are 1 and 6 (arduino).
func (e *Adaptor) GetConnection(address int, bus int) (connection i2c.Connection, err error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if !(bus == e.GetDefaultBus()) {
		return nil, errors.New("Unsupported I2C bus")
	}
	if e.i2cBus == nil {
		if bus == 6 && e.board == "arduino" {
			e.arduinoI2CSetup()
		}
		e.i2cBus, err = sysfs.NewI2cDevice(fmt.Sprintf("/dev/i2c-%d", bus))
	}
	return i2c.NewConnection(e.i2cBus, address), err
}

// GetDefaultBus returns the default i2c bus for this platform
func (e *Adaptor) GetDefaultBus() int {
	// Arduino uses bus 6
	if e.board == "arduino" {
		return 6
	}

	return 1
}

// DigitalPin returns matched sysfs.DigitalPin for specified values
func (e *Adaptor) DigitalPin(pin string, dir string) (sysfsPin sysfs.DigitalPinner, err error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	i := e.pinmap[pin]
	if e.digitalPins[i.pin] == nil {
		e.digitalPins[i.pin], err = e.newExportedPin(i.pin)
		if err != nil {
			return
		}

		if i.resistor > 0 {
			e.digitalPins[i.resistor], err = e.newExportedPin(i.resistor)
			if err != nil {
				return
			}
		}

		if i.levelShifter > 0 {
			e.digitalPins[i.levelShifter], err = e.newExportedPin(i.levelShifter)
			if err != nil {
				return
			}
		}

		if len(i.mux) > 0 {
			for _, mux := range i.mux {
				e.digitalPins[mux.pin], err = e.newExportedPin(mux.pin)
				if err != nil {
					return
				}

				err = pinWrite(e.digitalPins[mux.pin], sysfs.OUT, mux.value)
				if err != nil {
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
			err = pinWrite(e.digitalPins[i.resistor], sysfs.OUT, sysfs.LOW)
			if err != nil {
				return
			}
		}

		if i.levelShifter > 0 {
			err = pinWrite(e.digitalPins[i.levelShifter], sysfs.OUT, sysfs.LOW)
			if err != nil {
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
			err = pinWrite(e.digitalPins[i.levelShifter], sysfs.OUT, sysfs.HIGH)
			if err != nil {
				return
			}
		}
	}
	return e.digitalPins[i.pin], nil
}

// PWMPin returns a sysfs.PWMPin
func (e *Adaptor) PWMPin(pin string) (sysfsPin sysfs.PWMPinner, err error) {
	sysPin := e.pinmap[pin]
	if sysPin.pwmPin != -1 {
		if e.pwmPins[sysPin.pwmPin] == nil {
			if err = e.DigitalWrite(pin, 1); err != nil {
				return
			}
			if err = changePinMode(e, strconv.Itoa(int(sysPin.pin)), "1"); err != nil {
				return
			}
			e.mutex.Lock()
			defer e.mutex.Unlock()

			e.pwmPins[sysPin.pwmPin] = sysfs.NewPWMPin(sysPin.pwmPin)
			if err = e.pwmPins[sysPin.pwmPin].Export(); err != nil {
				return
			}
			if err = e.pwmPins[sysPin.pwmPin].Enable(true); err != nil {
				return
			}
		}

		sysfsPin = e.pwmPins[sysPin.pwmPin]
		return
	}
	err = errors.New("Not a PWM pin")
	return
}

// TODO: also check to see if device labels for
// /sys/class/gpio/gpiochip{200,216,232,248}/label == "pcal9555a"
func (e *Adaptor) checkForArduino() error {
	if err := e.exportTristatePin(); err != nil {
		return err
	}
	return nil
}

func (e *Adaptor) newExportedPin(pin int) (sysfsPin *sysfs.DigitalPin, err error) {
	sysfsPin = sysfs.NewDigitalPin(pin)
	err = sysfsPin.Export()
	return
}

func (e *Adaptor) exportTristatePin() (err error) {
	e.tristate, err = e.newExportedPin(214)
	return
}

// arduinoSetup does needed setup for the Arduino compatible breakout board
func (e *Adaptor) arduinoSetup() (err error) {
	if err = e.exportTristatePin(); err != nil {
		return err
	}

	err = pinWrite(e.tristate, sysfs.OUT, sysfs.LOW)
	if err != nil {
		return
	}

	for _, i := range []int{263, 262} {
		if err = e.newDigitalPin(i, sysfs.HIGH); err != nil {
			return err
		}
	}

	for _, i := range []int{240, 241, 242, 243} {
		if err = e.newDigitalPin(i, sysfs.LOW); err != nil {
			return err
		}
	}

	for _, i := range []string{"111", "115", "114", "109"} {
		if err = changePinMode(e, i, "1"); err != nil {
			return err
		}
	}

	for _, i := range []string{"131", "129", "40"} {
		if err = changePinMode(e, i, "0"); err != nil {
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
		if err = e.newDigitalPin(i, sysfs.LOW); err != nil {
			return err
		}
	}

	for _, i := range []string{"28", "27"} {
		if err = changePinMode(e, i, "1"); err != nil {
			return
		}
	}

	if err = e.tristate.Write(sysfs.HIGH); err != nil {
		return
	}

	return
}

func (e *Adaptor) sparkfunSetup() (err error) {
	return
}

// miniboardSetup does needed setup for Edison minibpard and other compatible boards
func (e *Adaptor) miniboardSetup() (err error) {
	return
}

func (e *Adaptor) newDigitalPin(i int, level int) (err error) {
	io := sysfs.NewDigitalPin(i)
	if err = io.Export(); err != nil {
		return
	}
	if err = io.Direction(sysfs.OUT); err != nil {
		return
	}
	if err = io.Write(level); err != nil {
		return
	}
	err = io.Unexport()
	return
}

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

// changePinMode writes pin mode to current_pinmux file
func changePinMode(a *Adaptor, pin, mode string) (err error) {
	_, err = a.writeFile(
		"/sys/kernel/debug/gpio_debug/gpio"+pin+"/current_pinmux",
		[]byte("mode"+mode),
	)
	return
}

// pinWrite sets Direction and writes level for a specific pin
func pinWrite(pin *sysfs.DigitalPin, dir string, level int) (err error) {
	if err = pin.Direction(dir); err != nil {
		return
	}

	err = pin.Write(level)
	return
}
