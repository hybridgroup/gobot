package up2

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
	// LEDRed is the built-in red LED.
	LEDRed = "red"
	// LEDBlue is the built-in blue LED.
	LEDBlue = "blue"
	// LEDGreen is the built-in green LED.
	LEDGreen = "green"
	// LEDYellow is the built-in yellow LED.
	LEDYellow = "yellow"

	defaultI2cBusNumber = 5

	defaultSpiBusNumber  = 0
	defaultSpiChipNumber = 0
	defaultSpiMode       = 0
	defaultSpiBitsNumber = 8
	defaultSpiMaxSpeed   = 500000
)

type sysfsPin struct {
	pin    int
	pwmPin int
}

// Adaptor represents a Gobot Adaptor for the Upboard UP2
type Adaptor struct {
	name    string
	sys     *system.Accesser
	mutex   sync.Mutex
	pinmap  map[string]sysfsPin
	ledPath string
	*adaptors.DigitalPinsAdaptor
	*adaptors.PWMPinsAdaptor
	*adaptors.I2cBusAdaptor
	*adaptors.SpiBusAdaptor
}

// NewAdaptor creates a UP2 Adaptor
//
// Optional parameters:
//
//	adaptors.WithGpiodAccess():	use character device gpiod driver instead of sysfs
//	adaptors.WithSpiGpioAccess(sclk, ncs, sdo, sdi):	use GPIO's instead of /dev/spidev#.#
//
//	Optional parameters for PWM, see [adaptors.NewPWMPinsAdaptor]
func NewAdaptor(opts ...interface{}) *Adaptor {
	sys := system.NewAccesser()
	a := &Adaptor{
		name:    gobot.DefaultName("UP2"),
		sys:     sys,
		ledPath: "/sys/class/leds/upboard:%s:/brightness",
		pinmap:  fixedPins,
	}

	var digitalPinsOpts []func(adaptors.DigitalPinsOptioner)
	var pwmPinsOpts []adaptors.PwmPinsOptionApplier
	for _, opt := range opts {
		switch o := opt.(type) {
		case func(adaptors.DigitalPinsOptioner):
			digitalPinsOpts = append(digitalPinsOpts, o)
		case adaptors.PwmPinsOptionApplier:
			pwmPinsOpts = append(pwmPinsOpts, o)
		default:
			panic(fmt.Sprintf("'%s' can not be applied on adaptor '%s'", opt, a.name))
		}
	}

	// Valid bus numbers are [5,6] which corresponds to /dev/i2c-5 through /dev/i2c-6.
	i2cBusNumberValidator := adaptors.NewBusNumberValidator([]int{5, 6})
	// Valid bus numbers are [0,1] which corresponds to /dev/spidev0.x through /dev/spidev1.x.
	// x is the chip number <255
	spiBusNumberValidator := adaptors.NewBusNumberValidator([]int{0, 1})

	a.DigitalPinsAdaptor = adaptors.NewDigitalPinsAdaptor(sys, a.translateDigitalPin, digitalPinsOpts...)
	a.PWMPinsAdaptor = adaptors.NewPWMPinsAdaptor(sys, a.translatePWMPin, pwmPinsOpts...)
	a.I2cBusAdaptor = adaptors.NewI2cBusAdaptor(sys, i2cBusNumberValidator.Validate, defaultI2cBusNumber)
	a.SpiBusAdaptor = adaptors.NewSpiBusAdaptor(sys, spiBusNumberValidator.Validate, defaultSpiBusNumber,
		defaultSpiChipNumber, defaultSpiMode, defaultSpiBitsNumber, defaultSpiMaxSpeed)
	return a
}

// Name returns the name of the Adaptor
func (a *Adaptor) Name() string { return a.name }

// SetName sets the name of the Adaptor
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

	if err := a.PWMPinsAdaptor.Connect(); err != nil {
		return err
	}
	return a.DigitalPinsAdaptor.Connect()
}

// Finalize closes connection to board and pins
func (a *Adaptor) Finalize() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	err := a.DigitalPinsAdaptor.Finalize()

	if e := a.PWMPinsAdaptor.Finalize(); e != nil {
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

// DigitalWrite writes digital value to the specified pin.
func (a *Adaptor) DigitalWrite(id string, val byte) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	// is it one of the built-in LEDs?
	if id == LEDRed || id == LEDBlue || id == LEDGreen || id == LEDYellow {
		pinPath := fmt.Sprintf(a.ledPath, id)
		fi, err := a.sys.OpenFile(pinPath, os.O_WRONLY|os.O_APPEND, 0o666)
		defer fi.Close() //nolint:staticcheck // for historical reasons
		if err != nil {
			return err
		}
		_, err = fi.WriteString(strconv.Itoa(int(val)))
		return err
	}

	return a.DigitalPinsAdaptor.DigitalWrite(id, val)
}

func (a *Adaptor) translateDigitalPin(id string) (string, int, error) {
	if val, ok := a.pinmap[id]; ok {
		return "", val.pin, nil
	}
	return "", -1, fmt.Errorf("'%s' is not a valid id for a digital pin", id)
}

func (a *Adaptor) translatePWMPin(id string) (string, int, error) {
	sysPin, ok := a.pinmap[id]
	if !ok {
		return "", -1, fmt.Errorf("'%s' is not a valid id for a pin", id)
	}
	if sysPin.pwmPin == -1 {
		return "", -1, fmt.Errorf("'%s' is not a valid id for a PWM pin", id)
	}
	return "/sys/class/pwm/pwmchip0", sysPin.pwmPin, nil
}
