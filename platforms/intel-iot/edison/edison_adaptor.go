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
	"gobot.io/x/gobot/platforms/adaptors"
	"gobot.io/x/gobot/system"
)

const pwmPeriodDefault = 10000000 // 10 ms = 100 Hz

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
	name  string
	board string
	sys   *system.Accesser
	*adaptors.PWMPinsAdaptor
	mutex       sync.Mutex
	pinmap      map[string]sysfsPin
	tristate    gobot.DigitalPinner
	digitalPins map[int]gobot.DigitalPinner
	i2cBus      i2c.I2cDevice
}

// NewAdaptor returns a new Edison Adaptor
func NewAdaptor() *Adaptor {
	sys := system.NewAccesser()
	c := &Adaptor{
		name:   gobot.DefaultName("Edison"),
		sys:    sys,
		pinmap: arduinoPinMap,
	}
	c.PWMPinsAdaptor = adaptors.NewPWMPinsAdaptor(sys, pwmPeriodDefault, c.translateAndMuxPWMPin,
		adaptors.WithPWMPinInitializer(pwmPinInitializer))
	return c
}

// Name returns the Adaptors name
func (c *Adaptor) Name() string { return c.name }

// SetName sets the Adaptors name
func (c *Adaptor) SetName(n string) { c.name = n }

// Board returns the Adaptors board name
func (c *Adaptor) Board() string { return c.board }

// SetBoard sets the Adaptors name
func (c *Adaptor) SetBoard(n string) { c.board = n }

// Connect initializes the Edison for use with the Arduino breakout board
func (c *Adaptor) Connect() error {
	c.digitalPins = make(map[int]gobot.DigitalPinner)

	if err := c.PWMPinsAdaptor.Connect(); err != nil {
		return err
	}

	switch c.board {
	case "sparkfun":
		c.pinmap = sparkfunPinMap
	case "arduino", "":
		c.board = "arduino"
		c.pinmap = arduinoPinMap
		if err := c.arduinoSetup(); err != nil {
			return err
		}
	case "miniboard":
		c.pinmap = miniboardPinMap
	default:
		return errors.New("Unknown board type: " + c.Board())
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
	for _, pin := range c.digitalPins {
		if pin != nil {
			if errs := pin.Unexport(); errs != nil {
				err = multierror.Append(err, errs)
			}
		}
	}
	if e := c.PWMPinsAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}
	if c.i2cBus != nil {
		if errs := c.i2cBus.Close(); errs != nil {
			err = multierror.Append(err, errs)
		}
	}
	return
}

// DigitalRead reads digital value from pin
func (c *Adaptor) DigitalRead(pin string) (i int, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	sysPin, err := c.digitalPin(pin, system.WithDirectionInput())
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

// GetConnection returns an i2c connection to a device on a specified bus.
// Valid bus numbers are 1 and 6 (arduino).
func (c *Adaptor) GetConnection(address int, bus int) (connection i2c.Connection, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !(bus == c.GetDefaultBus()) {
		return nil, errors.New("Unsupported I2C bus")
	}
	if c.i2cBus == nil {
		if bus == 6 && c.board == "arduino" {
			c.arduinoI2CSetup()
		}
		c.i2cBus, err = c.sys.NewI2cDevice(fmt.Sprintf("/dev/i2c-%d", bus))
	}
	return i2c.NewConnection(c.i2cBus, address), err
}

// GetDefaultBus returns the default i2c bus for this platform
func (c *Adaptor) GetDefaultBus() int {
	// Arduino uses bus 6
	if c.board == "arduino" {
		return 6
	}

	return 1
}

// arduinoSetup does needed setup for the Arduino compatible breakout board
func (c *Adaptor) arduinoSetup() error {
	// TODO: also check to see if device labels for
	// /sys/class/gpio/gpiochip{200,216,232,248}/label == "pcal9555a"

	tpin, err := c.newExportedDigitalPin(214, system.WithDirectionOutput(system.LOW))
	if err != nil {
		return err
	}
	c.tristate = tpin

	for _, i := range []int{263, 262} {
		if err := c.newUnexportedDigitalPin(i, system.WithDirectionOutput(system.HIGH)); err != nil {
			return err
		}
	}

	for _, i := range []int{240, 241, 242, 243} {
		if err := c.newUnexportedDigitalPin(i, system.WithDirectionOutput(system.LOW)); err != nil {
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
	if err := c.tristate.Write(system.LOW); err != nil {
		return err
	}

	for _, i := range []int{14, 165, 212, 213} {
		if err := c.newUnexportedDigitalPin(i, system.WithDirectionInput()); err != nil {
			return err
		}
	}

	for _, i := range []int{236, 237, 204, 205} {
		if err := c.newUnexportedDigitalPin(i, system.WithDirectionOutput(system.LOW)); err != nil {
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
	file, err := c.sys.OpenFile(path, os.O_RDONLY, 0644)
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
		rop := system.WithDirectionOutput(system.LOW)
		if dir == system.OUT {
			rop = system.WithDirectionInput()
		}
		if err := c.ensureDigitalPin(i.resistor, rop); err != nil {
			return nil, err
		}
	}

	if i.levelShifter > 0 {
		lop := system.WithDirectionOutput(system.LOW)
		if dir == system.OUT {
			lop = system.WithDirectionOutput(system.HIGH)
		}
		if err := c.ensureDigitalPin(i.levelShifter, lop); err != nil {
			return nil, err
		}
	}

	if len(i.mux) > 0 {
		for _, mux := range i.mux {
			if err := c.ensureDigitalPin(mux.pin, system.WithDirectionOutput(mux.value)); err != nil {
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
	return pin.Enable(true)
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
	if err := c.changePinMode(strconv.Itoa(int(sysPin.pin)), "1"); err != nil {
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
	file, err := c.sys.OpenFile("/sys/kernel/debug/gpio_debug/gpio"+pin+"/current_pinmux", os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return err
	}
	_, err = file.Write([]byte("mode" + mode))
	return err
}

func (c *Adaptor) digitalWrite(pin string, val byte) (err error) {
	sysPin, err := c.digitalPin(pin, system.WithDirectionOutput(int(val)))
	if err != nil {
		return
	}
	return sysPin.Write(int(val))
}
