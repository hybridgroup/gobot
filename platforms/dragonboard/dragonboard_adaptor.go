package dragonboard

import (
	"errors"
	"fmt"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/sysfs"
)

// Adaptor represents a Gobot Adaptor for a DragonBoard 410c
type Adaptor struct {
	name        string
	digitalPins map[int]*sysfs.DigitalPin
	pinMap      map[string]int
	i2cBuses    [3]i2c.I2cDevice
	mutex       *sync.Mutex
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
		name:  gobot.DefaultName("DragonBoard"),
		mutex: &sync.Mutex{},
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

// DigitalPin returns matched digitalPin for specified values
func (c *Adaptor) DigitalPin(pin string, dir string) (sysfsPin *sysfs.DigitalPin, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	i, err := c.translatePin(pin)

	if err != nil {
		return
	}

	if c.digitalPins[i] == nil {
		c.digitalPins[i] = sysfs.NewDigitalPin(i)
		if err = c.digitalPins[i].Export(); err != nil {
			return
		}
	}

	if err = c.digitalPins[i].Direction(dir); err != nil {
		return
	}

	return c.digitalPins[i], nil
}

// DigitalRead reads digital value to the specified pin.
// Valids pins are the GPIO_A through GPIO_L pins from the
// extender (pins 23-34 on header J8), as well as the SoC pins
// aka all the other pins, APQ GPIO_0-GPIO_122 and PM_MPP_0-4.
func (c *Adaptor) DigitalRead(pin string) (val int, err error) {
	sysfsPin, err := c.DigitalPin(pin, sysfs.IN)
	if err != nil {
		return
	}
	return sysfsPin.Read()
}

// DigitalWrite writes digital value to the specified pin.
// Valids pins are the GPIO_A through GPIO_L pins from the
// extender (pins 23-34 on header J8), as well as the SoC pins
// aka all the other pins, APQ GPIO_0-GPIO_122 and PM_MPP_0-4.
func (c *Adaptor) DigitalWrite(pin string, val byte) (err error) {
	sysfsPin, err := c.DigitalPin(pin, sysfs.OUT)
	if err != nil {
		return err
	}
	return sysfsPin.Write(int(val))
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
		c.i2cBuses[bus], err = sysfs.NewI2cDevice(fmt.Sprintf("/dev/i2c-%d", bus))
	}
	return i2c.NewConnection(c.i2cBuses[bus], address), err
}

// GetDefaultBus returns the default i2c bus for this platform
func (c *Adaptor) GetDefaultBus() int {
	return 0
}

func (c *Adaptor) setPins() {
	c.digitalPins = make(map[int]*sysfs.DigitalPin)
	c.pinMap = fixedPins
	for i := 0; i < 122; i++ {
		pin := fmt.Sprintf("GPIO_%d", i)
		c.pinMap[pin] = i
	}
}

func (c *Adaptor) translatePin(pin string) (i int, err error) {
	if val, ok := c.pinMap[pin]; ok {
		i = val
	} else {
		err = errors.New("Not a valid pin")
	}
	return
}
