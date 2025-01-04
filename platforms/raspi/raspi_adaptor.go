package raspi

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	multierror "github.com/hashicorp/go-multierror"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/adaptors"
	"gobot.io/x/gobot/v2/system"
)

const (
	infoFile = "/proc/cpuinfo"

	defaultSpiBusNumber  = 0
	defaultSpiChipNumber = 0
	defaultSpiMode       = 0
	defaultSpiBitsNumber = 8
	defaultSpiMaxSpeed   = 500000
)

// Adaptor is the Gobot Adaptor for the Raspberry Pi
type Adaptor struct {
	name     string
	mutex    sync.Mutex
	sys      *system.Accesser
	revision string
	*adaptors.AnalogPinsAdaptor
	*adaptors.DigitalPinsAdaptor
	*adaptors.PWMPinsAdaptor
	*adaptors.I2cBusAdaptor
	*adaptors.SpiBusAdaptor
}

// NewAdaptor creates a Raspi Adaptor
//
// Optional parameters:
//
//	adaptors.WithGpioSysfsAccess():	use legacy sysfs driver instead of default character device driver
//	adaptors.WithSpiGpioAccess(sclk, ncs, sdo, sdi):	use GPIO's instead of /dev/spidev#.#
//	adaptors.WithGpiosActiveLow(pin's): invert the pin behavior
//	adaptors.WithGpiosPullUp/Down(pin's): sets the internal pull resistor
//	adaptors.WithGpiosOpenDrain/Source(pin's): sets the output behavior
//	adaptors.WithGpioDebounce(pin, period): sets the input debouncer
//	adaptors.WithGpioEventOnFallingEdge/RaisingEdge/BothEdges(pin, handler): activate edge detection
func NewAdaptor(opts ...interface{}) *Adaptor {
	sys := system.NewAccesser()
	a := &Adaptor{
		name: gobot.DefaultName("RaspberryPi"),
		sys:  sys,
	}

	var digitalPinsOpts []adaptors.DigitalPinsOptionApplier
	var pwmPinsOpts []adaptors.PwmPinsOptionApplier
	for _, opt := range opts {
		switch o := opt.(type) {
		case adaptors.DigitalPinsOptionApplier:
			digitalPinsOpts = append(digitalPinsOpts, o)
		case adaptors.PwmPinsOptionApplier:
			pwmPinsOpts = append(pwmPinsOpts, o)
		default:
			panic(fmt.Sprintf("'%s' can not be applied on adaptor '%s'", opt, a.name))
		}
	}

	analogPinTranslator := adaptors.NewAnalogPinTranslator(sys, analogPinDefinitions)
	// Valid bus numbers are [0,1] which corresponds to /dev/i2c-0 through /dev/i2c-1.
	i2cBusNumberValidator := adaptors.NewBusNumberValidator([]int{0, 1})
	// Valid bus numbers are [0,1] which corresponds to /dev/spidev0.x through /dev/spidev1.x.
	// x is the chip number <255
	spiBusNumberValidator := adaptors.NewBusNumberValidator([]int{0, 1})

	a.AnalogPinsAdaptor = adaptors.NewAnalogPinsAdaptor(sys, analogPinTranslator.Translate)
	a.DigitalPinsAdaptor = adaptors.NewDigitalPinsAdaptor(sys, a.getPinTranslatorFunction(), digitalPinsOpts...)
	a.PWMPinsAdaptor = adaptors.NewPWMPinsAdaptor(sys, a.getPinTranslatorFunction(), pwmPinsOpts...)
	a.I2cBusAdaptor = adaptors.NewI2cBusAdaptor(sys, i2cBusNumberValidator.Validate, 1)
	a.SpiBusAdaptor = adaptors.NewSpiBusAdaptor(sys, spiBusNumberValidator.Validate, defaultSpiBusNumber,
		defaultSpiChipNumber, defaultSpiMode, defaultSpiBitsNumber, defaultSpiMaxSpeed)
	return a
}

// Name returns the adaptors name
func (a *Adaptor) Name() string {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.name
}

// SetName sets the adaptors name
func (a *Adaptor) SetName(n string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.name = n
}

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

// Finalize closes connection to board and pins
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

// DefaultI2cBus returns the default i2c bus for this platform.
// This overrides the base function due to the revision dependency.
func (a *Adaptor) DefaultI2cBus() int {
	rev := a.readRevision()
	if rev == "2" || rev == "3" {
		return 1
	}
	return 0
}

// getPinTranslatorFunction returns a function to be able to translate GPIO and PWM pins.
// This means for pi-blaster usage, each pin can be used and therefore the pin is given as number, like a GPIO pin.
// For sysfs-PWM usage, the pin will be given as "pwm0" or "pwm1", because the real pin number depends on the user
// configuration in "/boot/config.txt". For further details, see "/boot/overlays/README".
func (a *Adaptor) getPinTranslatorFunction() func(string) (string, int, error) {
	return func(pin string) (string, int, error) {
		var line int
		if val, ok := pins[pin][a.readRevision()]; ok {
			line = val
		} else if val, ok := pins[pin]["*"]; ok {
			line = val
		} else {
			return "", 0, fmt.Errorf("'%s' is not a valid pin id for raspi revision %s", pin, a.revision)
		}
		// We always use "gpiochip0", because currently all pins are available with this approach. A change of the
		// translator would be needed to support different chips (e.g. gpiochip1) with different revisions.
		path := "gpiochip0"
		if strings.HasPrefix(pin, "pwm") {
			path = "/sys/class/pwm/pwmchip0"
		}

		return path, line, nil
	}
}

func (a *Adaptor) readRevision() string {
	if a.revision == "" {
		a.revision = "0"
		content, err := a.sys.ReadFile(infoFile)
		if err != nil {
			return a.revision
		}
		for _, v := range strings.Split(string(content), "\n") {
			if strings.Contains(v, "Revision") {
				s := strings.Split(v, " ")
				version, _ := strconv.ParseInt("0x"+s[len(s)-1], 0, 64)
				switch {
				case version <= 3:
					a.revision = "1"
				case version <= 15:
					a.revision = "2"
				default:
					a.revision = "3"
				}
			}
		}
	}

	return a.revision
}
