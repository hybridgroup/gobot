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
	pinMap      map[string]sysfsPin
	tristate    gobot.DigitalPinner
	digitalPins map[int]gobot.DigitalPinner
	*adaptors.AnalogPinsAdaptor
	*adaptors.PWMPinsAdaptor
	*adaptors.I2cBusAdaptor
	arduinoI2cInitialized bool
}

// NewAdaptor returns a new Edison Adaptor of the given type.
// Supported types are: "arduino", "miniboard", "sparkfun", an empty string defaults to "arduino"
//
//	Optional parameters for PWM, see [adaptors.NewPWMPinsAdaptor]
func NewAdaptor(opts ...interface{}) *Adaptor {
	sys := system.NewAccesser(system.WithDigitalPinSysfsAccess())
	sys.AddDigitalPinSupport()
	a := &Adaptor{
		name:   gobot.DefaultName("Edison"),
		board:  "arduino",
		sys:    sys,
		pinMap: arduinoPinMap,
	}

	pwmPinsOpts := []adaptors.PwmPinsOptionApplier{adaptors.WithPWMPinInitializer(pwmPinInitializer)}
	for _, opt := range opts {
		switch o := opt.(type) {
		case string:
			if o != "" {
				a.board = o
			}
		case adaptors.PwmPinsOptionApplier:
			pwmPinsOpts = append(pwmPinsOpts, o)
		default:
			panic(fmt.Sprintf("'%s' can not be applied on adaptor '%s'", opt, a.name))
		}
	}

	a.AnalogPinsAdaptor = adaptors.NewAnalogPinsAdaptor(sys, a.translateAnalogPin)
	a.PWMPinsAdaptor = adaptors.NewPWMPinsAdaptor(sys, a.translateAndMuxPWMPin, pwmPinsOpts...)
	defI2cBusNr := defaultI2cBusNumber
	if a.board != "arduino" {
		defI2cBusNr = defaultI2cBusNumberOther
	}
	a.I2cBusAdaptor = adaptors.NewI2cBusAdaptor(sys, a.validateAndSetupI2cBusNumber, defI2cBusNr)
	return a
}

// Name returns the adaptors name
func (a *Adaptor) Name() string { return a.name }

// SetName sets the adaptors name
func (a *Adaptor) SetName(n string) { a.name = n }

// Connect initializes the Edison for use with the Arduino breakout board
func (a *Adaptor) Connect() error {
	a.digitalPins = make(map[int]gobot.DigitalPinner)

	if err := a.I2cBusAdaptor.Connect(); err != nil {
		return err
	}

	if err := a.AnalogPinsAdaptor.Connect(); err != nil {
		return err
	}

	if err := a.PWMPinsAdaptor.Connect(); err != nil {
		return err
	}

	switch a.board {
	case "sparkfun":
		a.pinMap = sparkfunPinMap
	case "arduino":
		a.board = "arduino"
		a.pinMap = arduinoPinMap
		if err := a.arduinoSetup(); err != nil {
			return err
		}
	case "miniboard":
		a.pinMap = miniboardPinMap
	default:
		return fmt.Errorf("Unknown board type: %s", a.board)
	}

	return nil
}

// Finalize releases all i2c devices and exported analog, digital, pwm pins.
func (a *Adaptor) Finalize() error {
	var err error
	if a.tristate != nil {
		if errs := a.tristate.Unexport(); errs != nil {
			err = multierror.Append(err, errs)
		}
	}
	a.tristate = nil

	for _, pin := range a.digitalPins {
		if pin != nil {
			if errs := pin.Unexport(); errs != nil {
				err = multierror.Append(err, errs)
			}
		}
	}
	a.digitalPins = nil

	if e := a.PWMPinsAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}

	if e := a.AnalogPinsAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}

	if e := a.I2cBusAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}
	a.arduinoI2cInitialized = false
	return err
}

// DigitalRead reads digital value from pin
func (a *Adaptor) DigitalRead(pin string) (int, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	sysPin, err := a.digitalPin(pin, system.WithPinDirectionInput())
	if err != nil {
		return 0, err
	}
	return sysPin.Read()
}

// DigitalWrite writes a value to the pin. Acceptable values are 1 or 0.
func (a *Adaptor) DigitalWrite(pin string, val byte) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.digitalWrite(pin, val)
}

// DigitalPin returns a digital pin. If the pin is initially acquired, it is an input.
// Pin direction and other options can be changed afterwards by pin.ApplyOptions() at any time.
func (a *Adaptor) DigitalPin(id string) (gobot.DigitalPinner, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.digitalPin(id)
}

// AnalogRead returns value from analog reading of specified pin
func (a *Adaptor) AnalogRead(pin string) (int, error) {
	rawRead, err := a.AnalogPinsAdaptor.AnalogRead(pin)
	if err != nil {
		return 0, err
	}

	return rawRead / 4, err
}

func (a *Adaptor) validateAndSetupI2cBusNumber(busNr int) error {
	// Valid bus number is 6 for "arduino", otherwise 1.
	if busNr == 6 && a.board == "arduino" {
		if !a.arduinoI2cInitialized {
			if err := a.arduinoI2CSetup(); err != nil {
				return err
			}
			a.arduinoI2cInitialized = true
			return nil
		}
		return nil
	}

	if busNr == 1 && a.board != "arduino" {
		return nil
	}

	return fmt.Errorf("Unsupported I2C bus '%d'", busNr)
}

// arduinoSetup does needed setup for the Arduino compatible breakout board
func (a *Adaptor) arduinoSetup() error {
	// TODO: also check to see if device labels for
	// /sys/class/gpio/gpiochip{200,216,232,248}/label == "pcal9555a"

	tpin, err := a.newExportedDigitalPin(214, system.WithPinDirectionOutput(system.LOW))
	if err != nil {
		return err
	}
	a.tristate = tpin

	for _, i := range []int{263, 262} {
		if err := a.newUnexportedDigitalPin(i, system.WithPinDirectionOutput(system.HIGH)); err != nil {
			return err
		}
	}

	for _, i := range []int{240, 241, 242, 243} {
		if err := a.newUnexportedDigitalPin(i, system.WithPinDirectionOutput(system.LOW)); err != nil {
			return err
		}
	}

	for _, i := range []string{"111", "115", "114", "109"} {
		if err := a.changePinMode(i, "1"); err != nil {
			return err
		}
	}

	for _, i := range []string{"131", "129", "40"} {
		if err := a.changePinMode(i, "0"); err != nil {
			return err
		}
	}

	return a.tristate.Write(system.HIGH)
}

func (a *Adaptor) arduinoI2CSetup() error {
	if a.tristate == nil {
		return fmt.Errorf("not connected")
	}

	if err := a.tristate.Write(system.LOW); err != nil {
		return err
	}

	for _, i := range []int{14, 165, 212, 213} {
		if err := a.newUnexportedDigitalPin(i, system.WithPinDirectionInput()); err != nil {
			return err
		}
	}

	for _, i := range []int{236, 237, 204, 205} {
		if err := a.newUnexportedDigitalPin(i, system.WithPinDirectionOutput(system.LOW)); err != nil {
			return err
		}
	}

	for _, i := range []string{"28", "27"} {
		if err := a.changePinMode(i, "1"); err != nil {
			return err
		}
	}

	return a.tristate.Write(system.HIGH)
}

func (a *Adaptor) digitalPin(id string, o ...func(gobot.DigitalPinOptioner) bool) (gobot.DigitalPinner, error) {
	i := a.pinMap[id]

	err := a.ensureDigitalPin(i.pin, o...)
	if err != nil {
		return nil, err
	}
	pin := a.digitalPins[i.pin]
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
		if err := a.ensureDigitalPin(i.resistor, rop); err != nil {
			return nil, err
		}
	}

	if i.levelShifter > 0 {
		lop := system.WithPinDirectionOutput(system.LOW)
		if dir == system.OUT {
			lop = system.WithPinDirectionOutput(system.HIGH)
		}
		if err := a.ensureDigitalPin(i.levelShifter, lop); err != nil {
			return nil, err
		}
	}

	if len(i.mux) > 0 {
		for _, mux := range i.mux {
			if err := a.ensureDigitalPin(mux.pin, system.WithPinDirectionOutput(mux.value)); err != nil {
				return nil, err
			}
		}
	}

	return pin, nil
}

func (a *Adaptor) ensureDigitalPin(idx int, o ...func(gobot.DigitalPinOptioner) bool) error {
	pin := a.digitalPins[idx]
	var err error
	if pin == nil {
		pin, err = a.newExportedDigitalPin(idx, o...)
		if err != nil {
			return err
		}
		a.digitalPins[idx] = pin
	} else {
		if err := pin.ApplyOptions(o...); err != nil {
			return err
		}
	}
	return nil
}

func pwmPinInitializer(_ string, pin gobot.PWMPinner) error {
	if err := pin.Export(); err != nil {
		return err
	}
	return pin.SetEnabled(true)
}

func (a *Adaptor) translateAnalogPin(pin string) (string, bool, uint16, error) {
	path := fmt.Sprintf("/sys/bus/iio/devices/iio:device1/in_voltage%s_raw", pin)
	const (
		write      = false
		readBufLen = 200
	)
	return path, write, readBufLen, nil
}

func (a *Adaptor) translateAndMuxPWMPin(id string) (string, int, error) {
	sysPin, ok := a.pinMap[id]
	if !ok {
		return "", -1, fmt.Errorf("'%s' is not a valid id for a pin", id)
	}
	if sysPin.pwmPin == -1 {
		return "", -1, fmt.Errorf("'%s' is not a valid id for a PWM pin", id)
	}
	if err := a.digitalWrite(id, 1); err != nil {
		return "", -1, err
	}
	if err := a.changePinMode(strconv.Itoa(sysPin.pin), "1"); err != nil {
		return "", -1, err
	}
	return "/sys/class/pwm/pwmchip0", sysPin.pwmPin, nil
}

func (a *Adaptor) newUnexportedDigitalPin(i int, o ...func(gobot.DigitalPinOptioner) bool) error {
	io := a.sys.NewDigitalPin("", i, o...)
	if err := io.Export(); err != nil {
		return err
	}
	return io.Unexport()
}

func (a *Adaptor) newExportedDigitalPin(
	pin int,
	o ...func(gobot.DigitalPinOptioner) bool,
) (gobot.DigitalPinner, error) {
	sysPin := a.sys.NewDigitalPin("", pin, o...)
	err := sysPin.Export()
	return sysPin, err
}

// changePinMode writes pin mode to current_pinmux file
func (a *Adaptor) changePinMode(pin, mode string) error {
	file, err := a.sys.OpenFile("/sys/kernel/debug/gpio_debug/gpio"+pin+"/current_pinmux", os.O_WRONLY, 0o644)
	defer file.Close() //nolint:staticcheck // for historical reasons
	if err != nil {
		return err
	}
	_, err = file.Write([]byte("mode" + mode))
	return err
}

func (a *Adaptor) digitalWrite(pin string, val byte) error {
	sysPin, err := a.digitalPin(pin, system.WithPinDirectionOutput(int(val)))
	if err != nil {
		return err
	}
	return sysPin.Write(int(val))
}
