package rockpi

import (
	"errors"
	"sync"

	multierror "github.com/hashicorp/go-multierror"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/adaptors"
	"gobot.io/x/gobot/v2/system"
)

const (
	procDeviceTreeModel = "/proc/device-tree/model"

	defaultI2cBusNumber = 7

	defaultSpiBusNumber  = 1
	defaultSpiChipNumber = 0
	defaultSpiMode       = 0
	defaultSpiBitsNumber = 8
	defaultSpiMaxSpeed   = 500000
)

// Adaptor is the Gobot Adaptor for Radxa's Rock Pi.
type Adaptor struct {
	name     string
	mutex    sync.Mutex
	sys      *system.Accesser
	revision string
	*adaptors.DigitalPinsAdaptor
	*adaptors.I2cBusAdaptor
	*adaptors.SpiBusAdaptor
}

// NewAdaptor creates a RockPi Adaptor
// Do not forget to enable the required overlays in /boot/hw_initfc.conf!
// See https://wiki.radxa.com/Rockpi4/dev/libmraa
//
// Optional parameters:
//
//	adaptors.WithGpiodAccess():	use character device gpiod driver instead of the default sysfs (NOT work on RockPi4C+!)
//	adaptors.WithSpiGpioAccess(sclk, ncs, sdo, sdi):	use GPIO's instead of /dev/spidev#.#
//	adaptors.WithGpiosActiveLow(pin's): invert the pin behavior
func NewAdaptor(opts ...func(adaptors.DigitalPinsOptioner)) *Adaptor {
	sys := system.NewAccesser()
	a := &Adaptor{
		name: gobot.DefaultName("RockPi"),
		sys:  sys,
	}

	// The RockPi4 has 3 I2C buses: 2, 6, 7. See https://wiki.radxa.com/Rock4/hardware/gpio
	// This could change in the future with other revisions!
	i2cBusNumberValidator := adaptors.NewBusNumberValidator([]int{2, 6, 7})
	// The RockPi4 has 2 SPI buses: 1, 2. See https://wiki.radxa.com/Rock4/hardware/gpio
	// This could change in the future with other revisions!
	spiBusNumberValidator := adaptors.NewBusNumberValidator([]int{1, 2})

	a.DigitalPinsAdaptor = adaptors.NewDigitalPinsAdaptor(sys, a.getPinTranslatorFunction(), opts...)
	a.I2cBusAdaptor = adaptors.NewI2cBusAdaptor(sys, i2cBusNumberValidator.Validate, defaultI2cBusNumber)
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

	return a.DigitalPinsAdaptor.Connect()
}

// Finalize closes connection to board and pins
func (a *Adaptor) Finalize() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	err := a.DigitalPinsAdaptor.Finalize()

	if e := a.I2cBusAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}

	if e := a.SpiBusAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}
	return err
}

func (a *Adaptor) getPinTranslatorFunction() func(string) (string, int, error) {
	return func(pin string) (string, int, error) {
		var line int
		if val, ok := pins[pin][a.readRevision()]; ok {
			line = val
		} else if val, ok := pins[pin]["*"]; ok {
			line = val
		} else {
			return "", 0, errors.New("Not a valid pin")
		}
		return "", line, nil
	}
}

func (a *Adaptor) readRevision() string {
	if a.revision == "" {
		content, err := a.sys.ReadFile(procDeviceTreeModel)
		if err != nil {
			return a.revision
		}
		model := string(content)
		switch model {
		case "Radxa ROCK 4":
			a.revision = "4"
		case "Radxa ROCK 4C+":
			a.revision = "4C+"
		default:
			a.revision = "4"
		}
	}

	return a.revision
}
