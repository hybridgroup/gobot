package beaglebone

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	multierror "github.com/hashicorp/go-multierror"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/adaptors"
	"gobot.io/x/gobot/v2/system"
)

const (
	pwmPeriodDefault = 500000 // 0.5 ms = 2 kHz

	defaultI2cBusNumber = 2

	defaultSpiBusNumber  = 0
	defaultSpiChipNumber = 0
	defaultSpiMode       = 0
	defaultSpiBitsNumber = 8
	defaultSpiMaxSpeed   = 500000
)

// Adaptor is the gobot.Adaptor representation for the Beaglebone Black/Green
type Adaptor struct {
	name  string
	sys   *system.Accesser
	mutex *sync.Mutex
	*adaptors.AnalogPinsAdaptor
	*adaptors.DigitalPinsAdaptor
	*adaptors.PWMPinsAdaptor
	*adaptors.I2cBusAdaptor
	*adaptors.SpiBusAdaptor
	usrLed string
}

// NewAdaptor returns a new Beaglebone Black/Green Adaptor
//
// Optional parameters:
//
//	adaptors.WithGpioCdevAccess():	use character device driver instead of sysfs
//	adaptors.WithSpiGpioAccess(sclk, ncs, sdo, sdi):	use GPIO's instead of /dev/spidev#.#
//
//	Optional parameters for PWM, see [adaptors.NewPWMPinsAdaptor]
func NewAdaptor(opts ...interface{}) *Adaptor {
	sys := system.NewAccesser(system.WithDigitalPinSysfsAccess())
	pwmPinTranslator := adaptors.NewPWMPinTranslator(sys, bbbPwmPinMap)
	a := &Adaptor{
		name:   gobot.DefaultName("BeagleboneBlack"),
		sys:    sys,
		mutex:  &sync.Mutex{},
		usrLed: "/sys/class/leds/beaglebone:green:",
	}

	var digitalPinsOpts []adaptors.DigitalPinsOptionApplier
	pwmPinsOpts := []adaptors.PwmPinsOptionApplier{adaptors.WithPWMDefaultPeriod(pwmPeriodDefault)}
	var spiBusOpts []adaptors.SpiBusOptionApplier
	for _, opt := range opts {
		switch o := opt.(type) {
		case adaptors.DigitalPinsOptionApplier:
			digitalPinsOpts = append(digitalPinsOpts, o)
		case adaptors.PwmPinsOptionApplier:
			pwmPinsOpts = append(pwmPinsOpts, o)
		case adaptors.SpiBusOptionApplier:
			spiBusOpts = append(spiBusOpts, o)
		default:
			panic(fmt.Sprintf("'%s' can not be applied on adaptor '%s'", opt, a.name))
		}
	}

	analogPinTranslator := adaptors.NewAnalogPinTranslator(sys, bbbAnalogPinMap)
	// Valid bus number is either 0 or 2 which corresponds to /dev/i2c-0 or /dev/i2c-2.
	i2cBusNumberValidator := adaptors.NewBusNumberValidator([]int{0, 2})
	// Valid bus numbers are [0,1] which corresponds to /dev/spidev0.x through /dev/spidev1.x.
	// x is the chip number <255
	spiBusNumberValidator := adaptors.NewBusNumberValidator([]int{0, 1})

	a.AnalogPinsAdaptor = adaptors.NewAnalogPinsAdaptor(sys, analogPinTranslator.Translate)
	a.DigitalPinsAdaptor = adaptors.NewDigitalPinsAdaptor(sys, a.translateAndMuxDigitalPin, digitalPinsOpts...)
	a.PWMPinsAdaptor = adaptors.NewPWMPinsAdaptor(sys, a.getTranslateAndMuxPWMPinFunc(pwmPinTranslator.Translate),
		pwmPinsOpts...)
	a.I2cBusAdaptor = adaptors.NewI2cBusAdaptor(sys, i2cBusNumberValidator.Validate, defaultI2cBusNumber)
	a.SpiBusAdaptor = adaptors.NewSpiBusAdaptor(sys, spiBusNumberValidator.Validate, defaultSpiBusNumber,
		defaultSpiChipNumber, defaultSpiMode, defaultSpiBitsNumber, defaultSpiMaxSpeed, a.DigitalPinsAdaptor, spiBusOpts...)
	return a
}

// Name returns the Adaptor name
func (a *Adaptor) Name() string { return a.name }

// SetName sets the Adaptor name
func (a *Adaptor) SetName(n string) { a.name = n }

// Connect create new connection to board and pins.
func (a *Adaptor) Connect() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if err := a.SpiBusAdaptor.Connect(); err != nil {
		return err
	}

	if err := a.I2cBusAdaptor.Connect(); err != nil {
		return err
	}

	if err := a.AnalogPinsAdaptor.Connect(); err != nil {
		return err
	}

	if err := a.PWMPinsAdaptor.Connect(); err != nil {
		return err
	}
	return a.DigitalPinsAdaptor.Connect()
}

// Finalize releases all i2c devices and exported analog, digital, pwm pins.
func (a *Adaptor) Finalize() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	err := a.DigitalPinsAdaptor.Finalize()

	if e := a.PWMPinsAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}

	if e := a.AnalogPinsAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}

	if e := a.I2cBusAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}

	if e := a.SpiBusAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}
	return err
}

// DigitalWrite writes a digital value to specified pin.
// valid usr pin values are usr0, usr1, usr2 and usr3
func (a *Adaptor) DigitalWrite(id string, val byte) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if strings.Contains(id, "usr") {
		fi, e := a.sys.OpenFile(a.usrLed+id+"/brightness", os.O_WRONLY|os.O_APPEND, 0o666)
		defer fi.Close() //nolint:staticcheck // for historical reasons
		if e != nil {
			return e
		}
		_, err := fi.WriteString(strconv.Itoa(int(val)))
		return err
	}

	return a.DigitalPinsAdaptor.DigitalWrite(id, val)
}

// translatePin converts digital pin name to pin position
func (a *Adaptor) translateAndMuxDigitalPin(id string) (string, int, error) {
	line, ok := bbbPinMap[id]
	if !ok {
		return "", -1, fmt.Errorf("'%s' is not a valid id for a digital pin", id)
	}
	// mux is done by id, not by line
	if err := a.muxPin(id, "gpio"); err != nil {
		return "", -1, err
	}
	return "", line, nil
}

func (a *Adaptor) getTranslateAndMuxPWMPinFunc(
	pwmPinTranslate func(id string) (string, int, error),
) func(id string) (string, int, error) {
	return func(id string) (string, int, error) {
		path, channel, err := pwmPinTranslate(id)
		if err != nil {
			return path, channel, err
		}

		if err := a.muxPin(id, "pwm"); err != nil {
			return "", -1, err
		}

		return path, channel, nil
	}
}

func (a *Adaptor) muxPin(pin, cmd string) error {
	path := fmt.Sprintf("/sys/devices/platform/ocp/ocp:%s_pinmux/state", pin)
	fi, e := a.sys.OpenFile(path, os.O_WRONLY, 0o666)
	defer fi.Close() //nolint:staticcheck // for historical reasons
	if e != nil {
		return e
	}
	_, e = fi.WriteString(cmd)
	return e
}
