package dragonboard

import (
	"fmt"
	"sync"

	multierror "github.com/hashicorp/go-multierror"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/adaptors"
	"gobot.io/x/gobot/v2/system"
)

const (
	defaultI2cBusNumber = 0

	defaultSpiMaxSpeed = 5000 // 5kHz (more than 15kHz not possible with SPI over GPIO)
)

// Adaptor represents a Gobot Adaptor for a DragonBoard 410c
type Adaptor struct {
	name   string
	sys    *system.Accesser // used for unit tests only
	mutex  sync.Mutex
	pinMap map[string]int
	*adaptors.DigitalPinsAdaptor
	*adaptors.I2cBusAdaptor
	*adaptors.SpiBusAdaptor // for usage of "adaptors.WithSpiGpioAccess()"
}

// Valid pins are the GPIO_A through GPIO_L pins from the
// extender (pins 23-34 on header J8), as well as the SoC pins
// aka all the other pins, APQ GPIO_0-GPIO_122 and PM_MPP_0-4.
var fixedPins = map[string]int{
	"GPIO_A": 36,
	"GPIO_B": 12,
	"GPIO_C": 13,
	"GPIO_D": 69,
	"GPIO_E": 115,
	"GPIO_F": 507,
	"GPIO_G": 24,
	"GPIO_H": 25,
	"GPIO_I": 35,
	"GPIO_J": 34,
	"GPIO_K": 28,
	"GPIO_L": 33,

	"LED_1": 21,
	"LED_2": 120,
}

// NewAdaptor creates a DragonBoard 410c Adaptor
//
// Optional parameters:
//
//	adaptors.WithGpioCdevAccess():	use character device driver instead of sysfs
//	adaptors.WithSpiGpioAccess(sclk, ncs, sdo, sdi):	use GPIO's instead of /dev/spidev#.#
func NewAdaptor(opts ...interface{}) *Adaptor {
	sys := system.NewAccesser(system.WithDigitalPinSysfsAccess())
	a := &Adaptor{
		name: gobot.DefaultName("DragonBoard"),
		sys:  sys,
	}

	var digitalPinsOpts []adaptors.DigitalPinsOptionApplier
	var spiBusOpts []adaptors.SpiBusOptionApplier
	for _, opt := range opts {
		switch o := opt.(type) {
		case adaptors.DigitalPinsOptionApplier:
			digitalPinsOpts = append(digitalPinsOpts, o)
		case adaptors.SpiBusOptionApplier:
			spiBusOpts = append(spiBusOpts, o)
		default:
			panic(fmt.Sprintf("'%s' can not be applied on adaptor '%s'", opt, a.name))
		}
	}

	// Valid bus numbers are [0,1] which corresponds to /dev/i2c-0 through /dev/i2c-1.
	i2cBusNumberValidator := adaptors.NewBusNumberValidator([]int{0, 1})

	a.DigitalPinsAdaptor = adaptors.NewDigitalPinsAdaptor(sys, a.translateDigitalPin, digitalPinsOpts...)
	a.I2cBusAdaptor = adaptors.NewI2cBusAdaptor(sys, i2cBusNumberValidator.Validate, defaultI2cBusNumber)

	// SPI is only supported when "adaptors.WithSpiGpioAccess()" is given
	if len(spiBusOpts) > 0 {
		a.SpiBusAdaptor = adaptors.NewSpiBusAdaptor(sys, func(int) error { return nil }, 0, 0, 0, 0, defaultSpiMaxSpeed,
			a.DigitalPinsAdaptor, spiBusOpts...)
	}

	a.pinMap = fixedPins
	for i := 0; i < 122; i++ {
		pin := fmt.Sprintf("GPIO_%d", i)
		a.pinMap[pin] = i
	}
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

	return err
}

func (a *Adaptor) translateDigitalPin(id string) (string, int, error) {
	if line, ok := a.pinMap[id]; ok {
		return "", line, nil
	}
	return "", -1, fmt.Errorf("'%s' is not a valid id for a digital pin", id)
}
