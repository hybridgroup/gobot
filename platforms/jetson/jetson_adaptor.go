package jetson

import (
	"errors"
	"fmt"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/drivers/spi"
	"gobot.io/x/gobot/platforms/adaptors"
	"gobot.io/x/gobot/system"
)

const pwmPeriodDefault = 3000000 // 3 ms = 333 Hz

// Adaptor is the Gobot adaptor for the Jetson Nano
type Adaptor struct {
	name  string
	sys   *system.Accesser
	mutex sync.Mutex
	*adaptors.DigitalPinsAdaptor
	pwmPins    map[string]gobot.PWMPinner
	i2cBuses   [2]i2c.I2cDevice
	spiDevices [2]spi.Connection
}

// NewAdaptor creates a Jetson Nano adaptor
func NewAdaptor() *Adaptor {
	sys := system.NewAccesser()
	c := &Adaptor{
		name: gobot.DefaultName("JetsonNano"),
		sys:  sys,
	}
	c.DigitalPinsAdaptor = adaptors.NewDigitalPinsAdaptor(sys, c.translateDigitalPin)
	return c
}

// Name returns the Adaptor's name
func (c *Adaptor) Name() string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.name
}

// SetName sets the Adaptor's name
func (c *Adaptor) SetName(n string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.name = n
}

// Connect create new connection to board and pins.
func (c *Adaptor) Connect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.pwmPins = make(map[string]gobot.PWMPinner)
	return c.DigitalPinsAdaptor.Connect()
}

// Finalize closes connection to board and pins
func (c *Adaptor) Finalize() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	err := c.DigitalPinsAdaptor.Finalize()

	for _, pin := range c.pwmPins {
		if pin != nil {
			if perr := pin.Unexport(); err != nil {
				err = multierror.Append(err, perr)
			}
		}
	}
	c.pwmPins = nil

	for _, bus := range c.i2cBuses {
		if bus != nil {
			if e := bus.Close(); e != nil {
				err = multierror.Append(err, e)
			}
		}
	}
	for _, dev := range c.spiDevices {
		if dev != nil {
			if e := dev.Close(); e != nil {
				err = multierror.Append(err, e)
			}
		}
	}
	return err
}

// GetConnection returns an i2c connection to a device on a specified bus.
// Valid bus number is [0..1] which corresponds to /dev/i2c-0 through /dev/i2c-1.
func (c *Adaptor) GetConnection(address int, bus int) (connection i2c.Connection, err error) {
	if (bus < 0) || (bus > 1) {
		return nil, fmt.Errorf("Bus number %d out of range", bus)
	}

	device, err := c.getI2cBus(bus)

	return i2c.NewConnection(device, address), err
}

// GetDefaultBus returns the default i2c bus for this platform
func (c *Adaptor) GetDefaultBus() int {
	return 1
}

// GetSpiConnection returns an spi connection to a device on a specified bus.
// Valid bus number is [0..1] which corresponds to /dev/spidev0.0 through /dev/spidev0.1.
func (c *Adaptor) GetSpiConnection(busNum, chipNum, mode, bits int, maxSpeed int64) (connection spi.Connection, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if (busNum < 0) || (busNum > 1) {
		return nil, fmt.Errorf("Bus number %d out of range", busNum)
	}

	if c.spiDevices[busNum] == nil {
		c.spiDevices[busNum], err = spi.GetSpiConnection(busNum, chipNum, mode, bits, maxSpeed)
	}

	return c.spiDevices[busNum], err
}

// GetSpiDefaultBus returns the default spi bus for this platform.
func (c *Adaptor) GetSpiDefaultBus() int {
	return 0
}

// GetSpiDefaultChip returns the default spi chip for this platform.
func (c *Adaptor) GetSpiDefaultChip() int {
	return 0
}

// GetSpiDefaultMode returns the default spi mode for this platform.
func (c *Adaptor) GetSpiDefaultMode() int {
	return 0
}

// GetSpiDefaultBits returns the default spi number of bits for this platform.
func (c *Adaptor) GetSpiDefaultBits() int {
	return 8
}

// GetSpiDefaultMaxSpeed returns the default spi bus for this platform.
func (c *Adaptor) GetSpiDefaultMaxSpeed() int64 {
	return 10000000
}

//PWMPin returns a Jetson Nano. PWMPin which provides the gobot.PWMPinner interface
func (c *Adaptor) PWMPin(pin string) (gobot.PWMPinner, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.pwmPin(pin)
}

// PwmWrite writes a PWM signal to the specified pin
func (c *Adaptor) PwmWrite(pin string, val byte) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	sysPin, err := c.pwmPin(pin)
	if err != nil {
		return err
	}

	duty := uint32(gobot.FromScale(float64(val), 0, 255) * float64(pwmPeriodDefault))
	return sysPin.SetDutyCycle(duty)
}

// ServoWrite writes a servo signal to the specified pin
func (c *Adaptor) ServoWrite(pin string, angle byte) (err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	sysPin, err := c.pwmPin(pin)
	if err != nil {
		return err
	}

	duty := uint32(gobot.FromScale(float64(angle), 0, 180) * float64(pwmPeriodDefault))
	return sysPin.SetDutyCycle(duty)
}

func (c *Adaptor) getI2cBus(bus int) (_ i2c.I2cDevice, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.i2cBuses[bus] == nil {
		c.i2cBuses[bus], err = c.sys.NewI2cDevice(fmt.Sprintf("/dev/i2c-%d", bus))
	}

	return c.i2cBuses[bus], err
}

func (c *Adaptor) pwmPin(pin string) (gobot.PWMPinner, error) {
	if c.pwmPins == nil {
		return nil, fmt.Errorf("not connected")
	}

	if c.pwmPins[pin] != nil {
		return c.pwmPins[pin], nil
	}

	fn, err := c.translatePwmPin(pin)
	if err != nil {
		return nil, err
	}

	c.pwmPins[pin] = NewPWMPin(c.sys, "/sys/class/pwm/pwmchip0", fn)
	c.pwmPins[pin].Export()
	c.pwmPins[pin].SetPeriod(pwmPeriodDefault)
	c.pwmPins[pin].Enable(true)

	return c.pwmPins[pin], nil
}

func (c *Adaptor) translateDigitalPin(id string) (string, int, error) {
	if line, ok := gpioPins[id]; ok {
		return "", line, nil
	}
	return "", -1, fmt.Errorf("'%s' is not a valid id for a digital pin", id)
}

func (c *Adaptor) translatePwmPin(pin string) (fn string, err error) {
	if fn, ok := pwmPins[pin]; ok {
		return fn, nil
	}
	return "", errors.New("Not a valid pin")
}
