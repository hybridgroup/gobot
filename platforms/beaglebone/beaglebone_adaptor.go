package beaglebone

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/drivers/spi"
	"gobot.io/x/gobot/platforms/adaptors"
	"gobot.io/x/gobot/system"
)

type pwmPinData struct {
	channel int
	path    string
}

const pwmDefaultPeriod = 500000

// Adaptor is the gobot.Adaptor representation for the Beaglebone Black/Green
type Adaptor struct {
	name  string
	sys   *system.Accesser
	mutex sync.Mutex
	*adaptors.DigitalPinsAdaptor
	pwmPins      map[string]gobot.PWMPinner
	i2cBuses     map[int]i2c.I2cDevice
	usrLed       string
	analogPath   string
	pinMap       map[string]int
	pwmPinMap    map[string]pwmPinData
	analogPinMap map[string]string
	findPin      func(pinPath string) (string, error)
	spiBuses     [2]spi.Connection
}

// NewAdaptor returns a new Beaglebone Black/Green Adaptor
func NewAdaptor() *Adaptor {
	sys := system.NewAccesser()
	c := &Adaptor{
		name:         gobot.DefaultName("BeagleboneBlack"),
		sys:          sys,
		pwmPins:      make(map[string]gobot.PWMPinner),
		i2cBuses:     make(map[int]i2c.I2cDevice),
		pinMap:       bbbPinMap,
		pwmPinMap:    bbbPwmPinMap,
		analogPinMap: bbbAnalogPinMap,
		findPin: func(pinPath string) (string, error) {
			files, err := filepath.Glob(pinPath)
			return files[0], err
		},
		usrLed:     "/sys/class/leds/beaglebone:green:",
		analogPath: "/sys/bus/iio/devices/iio:device0",
	}

	c.DigitalPinsAdaptor = adaptors.NewDigitalPinsAdaptor(sys, c.translateAndMuxDigitalPin)
	return c
}

// Name returns the Adaptor name
func (c *Adaptor) Name() string { return c.name }

// SetName sets the Adaptor name
func (c *Adaptor) SetName(n string) { c.name = n }

// Connect create new connection to board and pins.
func (c *Adaptor) Connect() error {
	err := c.DigitalPinsAdaptor.Connect()
	return err
}

// Finalize releases all i2c devices and exported analog, digital, pwm pins.
func (c *Adaptor) Finalize() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	err := c.DigitalPinsAdaptor.Finalize()

	for _, pin := range c.pwmPins {
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
	for _, bus := range c.spiBuses {
		if bus != nil {
			if e := bus.Close(); e != nil {
				err = multierror.Append(err, e)
			}
		}
	}
	return err
}

// PwmWrite writes the 0-254 value to the specified pin
func (c *Adaptor) PwmWrite(pin string, val byte) (err error) {
	pwmPin, err := c.PWMPin(pin)
	if err != nil {
		return
	}
	period, err := pwmPin.Period()
	if err != nil {
		return err
	}
	duty := gobot.FromScale(float64(val), 0, 255.0)
	return pwmPin.SetDutyCycle(uint32(float64(period) * duty))
}

// ServoWrite writes a servo signal to the specified pin
func (c *Adaptor) ServoWrite(pin string, angle byte) (err error) {
	pwmPin, err := c.PWMPin(pin)
	if err != nil {
		return
	}

	// TODO: take into account the actual period setting, not just assume default
	const minDuty = 100 * 0.0005 * pwmDefaultPeriod
	const maxDuty = 100 * 0.0020 * pwmDefaultPeriod
	duty := uint32(gobot.ToScale(gobot.FromScale(float64(angle), 0, 180), minDuty, maxDuty))
	return pwmPin.SetDutyCycle(duty)
}

// DigitalWrite writes a digital value to specified pin.
// valid usr pin values are usr0, usr1, usr2 and usr3
func (c *Adaptor) DigitalWrite(id string, val byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if strings.Contains(id, "usr") {
		fi, e := c.sys.OpenFile(c.usrLed+id+"/brightness", os.O_WRONLY|os.O_APPEND, 0666)
		defer fi.Close()
		if e != nil {
			return e
		}
		_, err := fi.WriteString(strconv.Itoa(int(val)))
		return err
	}

	return c.DigitalPinsAdaptor.DigitalWrite(id, val)
}

// PWMPin returns matched pwmPin for specified pin number
func (c *Adaptor) PWMPin(pin string) (gobot.PWMPinner, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	pinInfo, err := c.translatePwmPin(pin)
	if err != nil {
		return nil, err
	}

	if c.pwmPins[pin] == nil {
		newPath, err := c.findPin(pinInfo.path)
		if err != nil {
			return nil, err
		}
		newPin := c.sys.NewPWMPin(newPath, pinInfo.channel)
		if err := c.muxPin(pin, "pwm"); err != nil {
			return nil, err
		}
		if err = newPin.Export(); err != nil {
			return nil, err
		}
		if err = newPin.SetPeriod(pwmDefaultPeriod); err != nil {
			return nil, err
		}
		if err = newPin.Enable(true); err != nil {
			return nil, err
		}
		c.pwmPins[pin] = newPin
	}

	return c.pwmPins[pin], nil
}

// AnalogRead returns an analog value from specified pin
func (c *Adaptor) AnalogRead(pin string) (val int, err error) {
	analogPin, err := c.translateAnalogPin(pin)
	if err != nil {
		return
	}
	fi, err := c.sys.OpenFile(fmt.Sprintf("%v/%v", c.analogPath, analogPin), os.O_RDONLY, 0644)
	defer fi.Close()

	if err != nil {
		return
	}

	var buf = make([]byte, 1024)
	_, err = fi.Read(buf)
	if err != nil {
		return
	}

	val, _ = strconv.Atoi(strings.Split(string(buf), "\n")[0])
	return
}

// GetConnection returns a connection to a device on a specified bus.
// Valid bus number is either 0 or 2 which corresponds to /dev/i2c-0 or /dev/i2c-2.
func (c *Adaptor) GetConnection(address int, bus int) (connection i2c.Connection, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if (bus != 0) && (bus != 2) {
		return nil, fmt.Errorf("Bus number %d out of range", bus)
	}
	if c.i2cBuses[bus] == nil {
		c.i2cBuses[bus], err = c.sys.NewI2cDevice(fmt.Sprintf("/dev/i2c-%d", bus))
	}
	return i2c.NewConnection(c.i2cBuses[bus], address), err
}

// GetDefaultBus returns the default i2c bus for this platform
func (c *Adaptor) GetDefaultBus() int {
	return 2
}

// GetSpiConnection returns an spi connection to a device on a specified bus.
// Valid bus number is [0..1] which corresponds to /dev/spidev0.0 through /dev/spidev0.1.
func (c *Adaptor) GetSpiConnection(busNum, chipNum, mode, bits int, maxSpeed int64) (connection spi.Connection, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if (busNum < 0) || (busNum > 1) {
		return nil, fmt.Errorf("Bus number %d out of range", busNum)
	}

	if c.spiBuses[busNum] == nil {
		c.spiBuses[busNum], err = spi.GetSpiConnection(busNum, chipNum, mode, bits, maxSpeed)
	}

	return c.spiBuses[busNum], err
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
	return 500000
}

// translatePin converts digital pin name to pin position
func (c *Adaptor) translateAndMuxDigitalPin(id string) (string, int, error) {
	line, ok := c.pinMap[id]
	if !ok {
		return "", -1, fmt.Errorf("'%s' is not a valid id for a digital pin", id)
	}
	// mux is done by id, not by line
	if err := c.muxPin(id, "gpio"); err != nil {
		return "", -1, err
	}
	return "", line, nil
}

func (c *Adaptor) translatePwmPin(pin string) (pwmPinData, error) {
	if val, ok := c.pwmPinMap[pin]; ok {
		return val, nil
	}
	return pwmPinData{}, errors.New("Not a valid PWM pin")
}

// translateAnalogPin converts analog pin name to pin position
func (c *Adaptor) translateAnalogPin(pin string) (string, error) {
	if val, ok := c.analogPinMap[pin]; ok {
		return val, nil
	}
	return "", errors.New("Not a valid analog pin")
}

func (c *Adaptor) muxPin(pin, cmd string) error {
	path := fmt.Sprintf("/sys/devices/platform/ocp/ocp:%s_pinmux/state", pin)
	fi, e := c.sys.OpenFile(path, os.O_WRONLY, 0666)
	defer fi.Close()
	if e != nil {
		return e
	}
	_, e = fi.WriteString(cmd)
	return e
}
