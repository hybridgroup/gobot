package edison

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/adaptors"
	"gobot.io/x/gobot/v2/system"
)

const (
	defaultI2cBusNumber      = 6
	defaultI2cBusNumberOther = 1
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
	*adaptors.PWMPinsAdaptor
	*adaptors.I2cBusAdaptor
	arduinoI2cInitialized bool
}

// NewAdaptor returns a new Edison Adaptor of the given type.
// Supported types are: "arduino", "miniboard", "sparkfun", an empty string defaults to "arduino"
func NewAdaptor(boardType ...string) *Adaptor {
	sys := system.NewAccesser()
	c := &Adaptor{
		name:   gobot.DefaultName("Edison"),
		board:  "arduino",
		sys:    sys,
		pinmap: arduinoPinMap,
	}
	if len(boardType) > 0 && boardType[0] != "" {
		c.board = boardType[0]
	}
	c.PWMPinsAdaptor = adaptors.NewPWMPinsAdaptor(sys, c.translateAndMuxPWMPin,
		adaptors.WithPWMPinInitializer(pwmPinInitializer))
	defI2cBusNr := defaultI2cBusNumber
	if c.board != "arduino" {
		defI2cBusNr = defaultI2cBusNumberOther
	}
	c.I2cBusAdaptor = adaptors.NewI2cBusAdaptor(sys, c.validateAndSetupI2cBusNumber, defI2cBusNr)
	return c
}

// Name returns the Adaptors name
func (c *Adaptor) Name() string { return c.name }

// SetName sets the Adaptors name
func (c *Adaptor) SetName(n string) { c.name = n }

// Connect initializes the Edison for use with the Arduino breakout board
func (c *Adaptor) Connect() error {
	c.digitalPins = make(map[int]gobot.DigitalPinner)

	if err := c.I2cBusAdaptor.Connect(); err != nil {
		return err
	}

	if err := c.PWMPinsAdaptor.Connect(); err != nil {
		return err
	}

	switch c.board {
	case "sparkfun":
		c.pinmap = sparkfunPinMap
	case "arduino":
		c.board = "arduino"
		c.pinmap = arduinoPinMap
		if err := c.arduinoSetup(); err != nil {
			return err
		}
	case "miniboard":
		c.pinmap = miniboardPinMap
	default:
		return fmt.Errorf("Unknown board type: %s", c.board)
	}

	return nil
}

// Finalize releases all i2c devices and exported analog, digital, pwm pins.
func (c *Adaptor) Finalize() (err error) {
	if c.tristate != nil {
		if errs := c.tristate.Unexport(); errs != nil {
			err = multierror.Append(err, errs)
		}
	}
	c.tristate = nil

	for _, pin := range c.digitalPins {
		if pin != nil {
			if errs := pin.Unexport(); errs != nil {
				err = multierror.Append(err, errs)
			}
		}
	}
	c.digitalPins = nil

	if e := c.PWMPinsAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}

	if e := c.I2cBusAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}
	c.arduinoI2cInitialized = false
	return
}

// DigitalRead reads digital value from pin
func (c *Adaptor) DigitalRead(pin string) (i int, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	sysPin, err := c.digitalPin(pin, system.WithPinDirectionInput())
	if err != nil {
		return
	}
	return sysPin.Read()
}

// DigitalWrite writes a value to the pin. Acceptable values are 1 or 0.
func (c *Adaptor) DigitalWrite(pin string, val byte) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.digitalWrite(pin, val)
}

// DigitalPin returns a digital pin. If the pin is initially acquired, it is an input.
// Pin direction and other options can be changed afterwards by pin.ApplyOptions() at any time.
func (c *Adaptor) DigitalPin(id string) (gobot.DigitalPinner, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.digitalPin(id)
}

// AnalogRead returns value from analog reading of specified pin
func (c *Adaptor) AnalogRead(pin string) (val int, err error) {
	buf, err := c.readFile("/sys/bus/iio/devices/iio:device1/in_voltage" + pin + "_raw")
	if err != nil {
		return
	}

	val, err = strconv.Atoi(string(buf[0 : len(buf)-1]))

	return val / 4, err
}

func (c *Adaptor) validateAndSetupI2cBusNumber(busNr int) error {
	// Valid bus number is 6 for "arduino", otherwise 1.
	if busNr == 6 && c.board == "arduino" {
		if !c.arduinoI2cInitialized {
			if err := c.arduinoI2CSetup(); err != nil {
				return err
			}
			c.arduinoI2cInitialized = true
			return nil
		}
		return nil
	}

	if busNr == 1 && c.board != "arduino" {
		return nil
	}

	return fmt.Errorf("Unsupported I2C bus '%d'", busNr)
}

// arduinoSetup does needed setup for the Arduino compatible breakout board
func (c *Adaptor) arduinoSetup() error {
	// TODO: also check to see if device labels for
	// /sys/class/gpio/gpiochip{200,216,232,248}/label == "pcal9555a"

	tpin, err := c.newExportedDigitalPin(214, system.WithPinDirectionOutput(system.LOW))
	if err != nil {
		return err
	}
	c.tristate = tpin

	for _, i := range []int{263, 262} {
		if err := c.newUnexportedDigitalPin(i, system.WithPinDirectionOutput(system.HIGH)); err != nil {
			return err
		}
	}

	for _, i := range []int{240, 241, 242, 243} {
		if err := c.newUnexportedDigitalPin(i, system.WithPinDirectionOutput(system.LOW)); err != nil {
			return err
		}
	}

	for _, i := range []string{"111", "115", "114", "109"} {
		if err := c.changePinMode(i, "1"); err != nil {
			return err
		}
	}

	for _, i := range []string{"131", "129", "40"} {
		if err := c.changePinMode(i, "0"); err != nil {
			return err
		}
	}

	return c.tristate.Write(system.HIGH)
}

func (c *Adaptor) arduinoI2CSetup() error {
	if c.tristate == nil {
		return fmt.Errorf("not connected")
	}

	if err := c.tristate.Write(system.LOW); err != nil {
		return err
	}

	for _, i := range []int{14, 165, 212, 213} {
		if err := c.newUnexportedDigitalPin(i, system.WithPinDirectionInput()); err != nil {
			return err
		}
	}

	for _, i := range []int{236, 237, 204, 205} {
		if err := c.newUnexportedDigitalPin(i, system.WithPinDirectionOutput(system.LOW)); err != nil {
			return err
		}
	}

	for _, i := range []string{"28", "27"} {
		if err := c.changePinMode(i, "1"); err != nil {
			return err
		}
	}

	return c.tristate.Write(system.HIGH)
}

func (c *Adaptor) readFile(path string) ([]byte, error) {
	file, err := c.sys.OpenFile(path, os.O_RDONLY, 0o644)
	defer file.Close() //nolint:staticcheck // for historical reasons
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

func (c *Adaptor) digitalPin(id string, o ...func(gobot.DigitalPinOptioner) bool) (gobot.DigitalPinner, error) {
	i := c.pinmap[id]

	err := c.ensureDigitalPin(i.pin, o...)
	if err != nil {
		return nil, err
	}
	pin := c.digitalPins[i.pin]
	vpin, ok := pin.(gobot.DigitalPinValuer)
	if !ok {
		return nil, fmt.Errorf("can not determine the direction behavior")
	}
	dir := vpin.DirectionBehavior()
	if i.resistor > 0 {
		rop := system.WithPinDirectionOutput(system.LOW)
		if dir == system.OUT {
			rop = system.WithPinDirectionInput()
		}
		if err := c.ensureDigitalPin(i.resistor, rop); err != nil {
			return nil, err
		}
	}

	if i.levelShifter > 0 {
		lop := system.WithPinDirectionOutput(system.LOW)
		if dir == system.OUT {
			lop = system.WithPinDirectionOutput(system.HIGH)
		}
		if err := c.ensureDigitalPin(i.levelShifter, lop); err != nil {
			return nil, err
		}
	}

	if len(i.mux) > 0 {
		for _, mux := range i.mux {
			if err := c.ensureDigitalPin(mux.pin, system.WithPinDirectionOutput(mux.value)); err != nil {
				return nil, err
			}
		}
	}

	return pin, nil
}

func (c *Adaptor) ensureDigitalPin(idx int, o ...func(gobot.DigitalPinOptioner) bool) error {
	pin := c.digitalPins[idx]
	var err error
	if pin == nil {
		pin, err = c.newExportedDigitalPin(idx, o...)
		if err != nil {
			return err
		}
		c.digitalPins[idx] = pin
	} else {
		if err := pin.ApplyOptions(o...); err != nil {
			return err
		}
	}
	return nil
}

func pwmPinInitializer(pin gobot.PWMPinner) error {
	if err := pin.Export(); err != nil {
		return err
	}
	return pin.SetEnabled(true)
}

func (c *Adaptor) translateAndMuxPWMPin(id string) (string, int, error) {
	sysPin, ok := c.pinmap[id]
	if !ok {
		return "", -1, fmt.Errorf("'%s' is not a valid id for a pin", id)
	}
	if sysPin.pwmPin == -1 {
		return "", -1, fmt.Errorf("'%s' is not a valid id for a PWM pin", id)
	}
	if err := c.digitalWrite(id, 1); err != nil {
		return "", -1, err
	}
	if err := c.changePinMode(strconv.Itoa(sysPin.pin), "1"); err != nil {
		return "", -1, err
	}
	return "/sys/class/pwm/pwmchip0", sysPin.pwmPin, nil
}

func (c *Adaptor) newUnexportedDigitalPin(i int, o ...func(gobot.DigitalPinOptioner) bool) error {
	io := c.sys.NewDigitalPin("", i, o...)
	if err := io.Export(); err != nil {
		return err
	}
	return io.Unexport()
}

func (c *Adaptor) newExportedDigitalPin(pin int, o ...func(gobot.DigitalPinOptioner) bool) (gobot.DigitalPinner, error) {
	sysPin := c.sys.NewDigitalPin("", pin, o...)
	err := sysPin.Export()
	return sysPin, err
}

// changePinMode writes pin mode to current_pinmux file
func (c *Adaptor) changePinMode(pin, mode string) error {
	file, err := c.sys.OpenFile("/sys/kernel/debug/gpio_debug/gpio"+pin+"/current_pinmux", os.O_WRONLY, 0o644)
	defer file.Close() //nolint:staticcheck // for historical reasons
	if err != nil {
		return err
	}
	_, err = file.Write([]byte("mode" + mode))
	return err
}

func (c *Adaptor) digitalWrite(pin string, val byte) (err error) {
	sysPin, err := c.digitalPin(pin, system.WithPinDirectionOutput(int(val)))
	if err != nil {
		return
	}
	return sysPin.Write(int(val))
}
