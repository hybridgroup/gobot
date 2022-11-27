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
	"gobot.io/x/gobot/system"
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
	sys         *system.Accesser
	mutex       sync.Mutex
	pinmap      map[string]sysfsPin
	tristate    gobot.DigitalPinner
	digitalPins map[int]gobot.DigitalPinner
	pwmPins     map[int]gobot.PWMPinner
	i2cBus      i2c.I2cDevice
}

// NewAdaptor returns a new Edison Adaptor
func NewAdaptor() *Adaptor {
	return &Adaptor{
		name:   gobot.DefaultName("Edison"),
		sys:    system.NewAccesser(),
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

// Connect initializes the Edison for use with the Arduino breakout board
func (e *Adaptor) Connect() error {
	e.digitalPins = make(map[int]gobot.DigitalPinner)
	e.pwmPins = make(map[int]gobot.PWMPinner)

	switch e.board {
	case "sparkfun":
		e.pinmap = sparkfunPinMap
	case "arduino", "":
		e.board = "arduino"
		e.pinmap = arduinoPinMap
		if err := e.arduinoSetup(); err != nil {
			return err
		}
	case "miniboard":
		e.pinmap = miniboardPinMap
	default:
		return errors.New("Unknown board type: " + e.Board())
	}
	return nil
}

// Finalize releases all i2c devices and exported analog, digital, pwm pins.
func (e *Adaptor) Finalize() (err error) {
	if e.tristate != nil {
		if errs := e.tristate.Unexport(); errs != nil {
			err = multierror.Append(err, errs)
		}
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
	e.mutex.Lock()
	defer e.mutex.Unlock()

	sysPin, err := e.digitalPin(pin, system.WithDirectionInput())
	if err != nil {
		return
	}
	return sysPin.Read()
}

// DigitalWrite writes a value to the pin. Acceptable values are 1 or 0.
func (e *Adaptor) DigitalWrite(pin string, val byte) (err error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	sysPin, err := e.digitalPin(pin, system.WithDirectionOutput(int(val)))
	if err != nil {
		return
	}
	return sysPin.Write(int(val))
}

// DigitalPin returns a digital pin. If the pin is initially acquired, it is an input.
// Pin direction and other options can be changed afterwards by pin.ApplyOptions() at any time.
func (e *Adaptor) DigitalPin(id string) (gobot.DigitalPinner, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return e.digitalPin(id)
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
	buf, err := e.readFile("/sys/bus/iio/devices/iio:device1/in_voltage" + pin + "_raw")
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
		e.i2cBus, err = e.sys.NewI2cDevice(fmt.Sprintf("/dev/i2c-%d", bus))
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

// PWMPin returns a system.PWMPin
func (e *Adaptor) PWMPin(pin string) (gobot.PWMPinner, error) {
	sysPin := e.pinmap[pin]
	if sysPin.pwmPin != -1 {
		if e.pwmPins[sysPin.pwmPin] == nil {
			if err := e.DigitalWrite(pin, 1); err != nil {
				return nil, err
			}
			if err := e.changePinMode(strconv.Itoa(int(sysPin.pin)), "1"); err != nil {
				return nil, err
			}
			e.mutex.Lock()
			defer e.mutex.Unlock()

			e.pwmPins[sysPin.pwmPin] = e.sys.NewPWMPin("/sys/class/pwm/pwmchip0", sysPin.pwmPin)
			if err := e.pwmPins[sysPin.pwmPin].Export(); err != nil {
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

func (e *Adaptor) newUnexportedDigitalPin(i int, o ...func(gobot.DigitalPinOptioner) bool) error {
	io := e.sys.NewDigitalPin("", i, o...)
	if err := io.Export(); err != nil {
		return err
	}
	return io.Unexport()
}

func (e *Adaptor) newExportedDigitalPin(pin int, o ...func(gobot.DigitalPinOptioner) bool) (gobot.DigitalPinner, error) {
	sysPin := e.sys.NewDigitalPin("", pin, o...)
	err := sysPin.Export()
	return sysPin, err
}

// arduinoSetup does needed setup for the Arduino compatible breakout board
func (e *Adaptor) arduinoSetup() error {
	// TODO: also check to see if device labels for
	// /sys/class/gpio/gpiochip{200,216,232,248}/label == "pcal9555a"

	tpin, err := e.newExportedDigitalPin(214, system.WithDirectionOutput(system.LOW))
	if err != nil {
		return err
	}
	e.tristate = tpin

	for _, i := range []int{263, 262} {
		if err := e.newUnexportedDigitalPin(i, system.WithDirectionOutput(system.HIGH)); err != nil {
			return err
		}
	}

	for _, i := range []int{240, 241, 242, 243} {
		if err := e.newUnexportedDigitalPin(i, system.WithDirectionOutput(system.LOW)); err != nil {
			return err
		}
	}

	for _, i := range []string{"111", "115", "114", "109"} {
		if err := e.changePinMode(i, "1"); err != nil {
			return err
		}
	}

	for _, i := range []string{"131", "129", "40"} {
		if err := e.changePinMode(i, "0"); err != nil {
			return err
		}
	}

	return e.tristate.Write(system.HIGH)
}

func (e *Adaptor) arduinoI2CSetup() error {
	if err := e.tristate.Write(system.LOW); err != nil {
		return err
	}

	for _, i := range []int{14, 165, 212, 213} {
		if err := e.newUnexportedDigitalPin(i, system.WithDirectionInput()); err != nil {
			return err
		}
	}

	for _, i := range []int{236, 237, 204, 205} {
		if err := e.newUnexportedDigitalPin(i, system.WithDirectionOutput(system.LOW)); err != nil {
			return err
		}
	}

	for _, i := range []string{"28", "27"} {
		if err := e.changePinMode(i, "1"); err != nil {
			return err
		}
	}

	return e.tristate.Write(system.HIGH)
}

func (e *Adaptor) writeFile(path string, data []byte) (i int, err error) {
	file, err := e.sys.OpenFile(path, os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return
	}

	return file.Write(data)
}

func (e *Adaptor) readFile(path string) ([]byte, error) {
	file, err := e.sys.OpenFile(path, os.O_RDONLY, 0644)
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
func (e *Adaptor) changePinMode(pin, mode string) error {
	_, err := e.writeFile("/sys/kernel/debug/gpio_debug/gpio"+pin+"/current_pinmux", []byte("mode"+mode))
	return err
}

func (e *Adaptor) digitalPin(id string, o ...func(gobot.DigitalPinOptioner) bool) (gobot.DigitalPinner, error) {
	i := e.pinmap[id]

	err := e.ensureDigitalPin(i.pin, o...)
	if err != nil {
		return nil, err
	}
	pin := e.digitalPins[i.pin]
	vpin, ok := pin.(gobot.DigitalPinValuer)
	if !ok {
		return nil, fmt.Errorf("can not determine the direction behavior")
	}
	dir := vpin.DirectionBehavior()
	if i.resistor > 0 {
		rop := system.WithDirectionOutput(system.LOW)
		if dir == system.OUT {
			rop = system.WithDirectionInput()
		}
		if err := e.ensureDigitalPin(i.resistor, rop); err != nil {
			return nil, err
		}
	}

	if i.levelShifter > 0 {
		lop := system.WithDirectionOutput(system.LOW)
		if dir == system.OUT {
			lop = system.WithDirectionOutput(system.HIGH)
		}
		if err := e.ensureDigitalPin(i.levelShifter, lop); err != nil {
			return nil, err
		}
	}

	if len(i.mux) > 0 {
		for _, mux := range i.mux {
			if err := e.ensureDigitalPin(mux.pin, system.WithDirectionOutput(mux.value)); err != nil {
				return nil, err
			}
		}
	}

	return pin, nil
}

func (e *Adaptor) ensureDigitalPin(idx int, o ...func(gobot.DigitalPinOptioner) bool) error {
	pin := e.digitalPins[idx]
	var err error
	if pin == nil {
		pin, err = e.newExportedDigitalPin(idx, o...)
		if err != nil {
			return err
		}
		e.digitalPins[idx] = pin
	} else {
		if err := pin.ApplyOptions(o...); err != nil {
			return err
		}
	}
	return nil
}
