package dragonboard

import (
	"errors"
	"fmt"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/system"
)

// Adaptor represents a Gobot Adaptor for a DragonBoard 410c
type Adaptor struct {
	name        string
	sys         *system.Accesser
	mutex       sync.Mutex
	digitalPins map[int]gobot.DigitalPinner
	pinMap      map[string]int
	i2cBuses    [3]i2c.I2cDevice
}

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
func NewAdaptor() *Adaptor {
	c := &Adaptor{
		name: gobot.DefaultName("DragonBoard"),
		sys:  system.NewAccesser(),
	}

	c.setPins()
	return c
}

// Name returns the name of the Adaptor
func (c *Adaptor) Name() string { return c.name }

// SetName sets the name of the Adaptor
func (c *Adaptor) SetName(n string) { c.name = n }

// Connect initializes the board
func (c *Adaptor) Connect() (err error) {
	return
}

// Finalize closes connection to board and pins
func (c *Adaptor) Finalize() (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, pin := range c.digitalPins {
		if pin != nil {
			if e := pin.Unexport(); e != nil {
				err = multierror.Append(err, e)
			}
		}
	}
	for _, bus := range c.i2cBuses {
		if bus != nil {
			if e := bus.Close(); e != nil {
				err = multierror.Append(err, e)
			}
		}
	}
	return
}

// DigitalRead reads digital value to the specified pin.
// Valids pins are the GPIO_A through GPIO_L pins from the
// extender (pins 23-34 on header J8), as well as the SoC pins
// aka all the other pins, APQ GPIO_0-GPIO_122 and PM_MPP_0-4.
func (c *Adaptor) DigitalRead(id string) (int, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	pin, err := c.digitalPin(id, system.WithDirectionInput())
	if err != nil {
		return 0, err
	}
	return pin.Read()
}

// DigitalWrite writes digital value to the specified pin.
// Valids pins are the GPIO_A through GPIO_L pins from the
// extender (pins 23-34 on header J8), as well as the SoC pins
// aka all the other pins, APQ GPIO_0-GPIO_122 and PM_MPP_0-4.
func (c *Adaptor) DigitalWrite(id string, val byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	pin, err := c.digitalPin(id, system.WithDirectionOutput(int(val)))
	if err != nil {
		return err
	}
	return pin.Write(int(val))
}

// DigitalPin returns a digital pin. If the pin is initially acquired, it is an input.
// Pin direction and other options can be changed afterwards by pin.ApplyOptions() at any time.
func (c *Adaptor) DigitalPin(id string) (gobot.DigitalPinner, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.digitalPin(id)
}

// GetConnection returns a connection to a device on a specified bus.
// Valid bus number is [0..1] which corresponds to /dev/i2c-0 through /dev/i2c-1.
func (c *Adaptor) GetConnection(address int, bus int) (connection i2c.Connection, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if (bus < 0) || (bus > 1) {
		return nil, fmt.Errorf("Bus number %d out of range", bus)
	}
	if c.i2cBuses[bus] == nil {
		c.i2cBuses[bus], err = c.sys.NewI2cDevice(fmt.Sprintf("/dev/i2c-%d", bus))
	}
	return i2c.NewConnection(c.i2cBuses[bus], address), err
}

// GetDefaultBus returns the default i2c bus for this platform
func (c *Adaptor) GetDefaultBus() int {
	return 0
}

func (c *Adaptor) setPins() {
	c.digitalPins = make(map[int]gobot.DigitalPinner)
	c.pinMap = fixedPins
	for i := 0; i < 122; i++ {
		pin := fmt.Sprintf("GPIO_%d", i)
		c.pinMap[pin] = i
	}
}

func (c *Adaptor) translateDigitalPin(id string) (i int, err error) {
	if val, ok := c.pinMap[id]; ok {
		i = val
	} else {
		err = errors.New("Not a valid pin")
	}
	return
}

func (c *Adaptor) digitalPin(id string, o ...func(gobot.DigitalPinOptioner) bool) (gobot.DigitalPinner, error) {
	i, err := c.translateDigitalPin(id)
	if err != nil {
		return nil, err
	}

	pin := c.digitalPins[i]
	if pin == nil {
		pin = c.sys.NewDigitalPin("", i, o...)
		if err = pin.Export(); err != nil {
			return nil, err
		}
		c.digitalPins[i] = pin
	} else {
		if err := pin.ApplyOptions(o...); err != nil {
			return nil, err
		}
	}

	return pin, nil
}
