package dragonboard

import (
	"fmt"
	"sync"

	multierror "github.com/hashicorp/go-multierror"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/adaptors"
	"gobot.io/x/gobot/v2/system"
)

const defaultI2cBusNumber = 0

// Adaptor represents a Gobot Adaptor for a DragonBoard 410c
type Adaptor struct {
	name   string
	sys    *system.Accesser // used for unit tests only
	mutex  sync.Mutex
	pinMap map[string]int
	*adaptors.DigitalPinsAdaptor
	*adaptors.I2cBusAdaptor
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
func NewAdaptor(opts ...adaptors.DigitalPinsOptionApplier) *Adaptor {
	sys := system.NewAccesser(system.WithDigitalPinSysfsAccess())
	c := &Adaptor{
		name: gobot.DefaultName("DragonBoard"),
		sys:  sys,
	}

	// Valid bus numbers are [0,1] which corresponds to /dev/i2c-0 through /dev/i2c-1.
	i2cBusNumberValidator := adaptors.NewBusNumberValidator([]int{0, 1})

	c.DigitalPinsAdaptor = adaptors.NewDigitalPinsAdaptor(sys, c.translateDigitalPin, opts...)
	c.I2cBusAdaptor = adaptors.NewI2cBusAdaptor(sys, i2cBusNumberValidator.Validate, defaultI2cBusNumber)
	c.pinMap = fixedPins
	for i := 0; i < 122; i++ {
		pin := fmt.Sprintf("GPIO_%d", i)
		c.pinMap[pin] = i
	}
	return c
}

// Name returns the name of the Adaptor
func (c *Adaptor) Name() string { return c.name }

// SetName sets the name of the Adaptor
func (c *Adaptor) SetName(n string) { c.name = n }

// Connect create new connection to board and pins.
func (c *Adaptor) Connect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if err := c.I2cBusAdaptor.Connect(); err != nil {
		return err
	}

	return c.DigitalPinsAdaptor.Connect()
}

// Finalize closes connection to board and pins
func (c *Adaptor) Finalize() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	err := c.DigitalPinsAdaptor.Finalize()

	if e := c.I2cBusAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}

	return err
}

func (c *Adaptor) translateDigitalPin(id string) (string, int, error) {
	if line, ok := c.pinMap[id]; ok {
		return "", line, nil
	}
	return "", -1, fmt.Errorf("'%s' is not a valid id for a digital pin", id)
}
