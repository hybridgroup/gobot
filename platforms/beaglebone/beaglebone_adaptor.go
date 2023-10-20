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

type pwmPinData struct {
	channel   int
	dir       string
	dirRegexp string
}

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
	mutex sync.Mutex
	*adaptors.DigitalPinsAdaptor
	*adaptors.PWMPinsAdaptor
	*adaptors.I2cBusAdaptor
	*adaptors.SpiBusAdaptor
	usrLed       string
	analogPath   string
	pinMap       map[string]int
	pwmPinMap    map[string]pwmPinData
	analogPinMap map[string]string
}

// NewAdaptor returns a new Beaglebone Black/Green Adaptor
//
// Optional parameters:
//
//	adaptors.WithGpiodAccess():	use character device gpiod driver instead of sysfs
//	adaptors.WithSpiGpioAccess(sclk, nss, mosi, miso):	use GPIO's instead of /dev/spidev#.#
func NewAdaptor(opts ...func(adaptors.Optioner)) *Adaptor {
	sys := system.NewAccesser()
	c := &Adaptor{
		name:         gobot.DefaultName("BeagleboneBlack"),
		sys:          sys,
		pinMap:       bbbPinMap,
		pwmPinMap:    bbbPwmPinMap,
		analogPinMap: bbbAnalogPinMap,
		usrLed:       "/sys/class/leds/beaglebone:green:",
		analogPath:   "/sys/bus/iio/devices/iio:device0",
	}

	c.DigitalPinsAdaptor = adaptors.NewDigitalPinsAdaptor(sys, c.translateAndMuxDigitalPin, opts...)
	c.PWMPinsAdaptor = adaptors.NewPWMPinsAdaptor(sys, c.translateAndMuxPWMPin,
		adaptors.WithPWMPinDefaultPeriod(pwmPeriodDefault))
	c.I2cBusAdaptor = adaptors.NewI2cBusAdaptor(sys, c.validateI2cBusNumber, defaultI2cBusNumber)
	c.SpiBusAdaptor = adaptors.NewSpiBusAdaptor(sys, c.validateSpiBusNumber, defaultSpiBusNumber, defaultSpiChipNumber,
		defaultSpiMode, defaultSpiBitsNumber, defaultSpiMaxSpeed)
	return c
}

// Name returns the Adaptor name
func (c *Adaptor) Name() string { return c.name }

// SetName sets the Adaptor name
func (c *Adaptor) SetName(n string) { c.name = n }

// Connect create new connection to board and pins.
func (c *Adaptor) Connect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if err := c.SpiBusAdaptor.Connect(); err != nil {
		return err
	}

	if err := c.I2cBusAdaptor.Connect(); err != nil {
		return err
	}

	if err := c.PWMPinsAdaptor.Connect(); err != nil {
		return err
	}
	return c.DigitalPinsAdaptor.Connect()
}

// Finalize releases all i2c devices and exported analog, digital, pwm pins.
func (c *Adaptor) Finalize() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	err := c.DigitalPinsAdaptor.Finalize()

	if e := c.PWMPinsAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}

	if e := c.I2cBusAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}

	if e := c.SpiBusAdaptor.Finalize(); e != nil {
		err = multierror.Append(err, e)
	}
	return err
}

// DigitalWrite writes a digital value to specified pin.
// valid usr pin values are usr0, usr1, usr2 and usr3
func (c *Adaptor) DigitalWrite(id string, val byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if strings.Contains(id, "usr") {
		fi, e := c.sys.OpenFile(c.usrLed+id+"/brightness", os.O_WRONLY|os.O_APPEND, 0o666)
		defer fi.Close() //nolint:staticcheck // for historical reasons
		if e != nil {
			return e
		}
		_, err := fi.WriteString(strconv.Itoa(int(val)))
		return err
	}

	return c.DigitalPinsAdaptor.DigitalWrite(id, val)
}

// AnalogRead returns an analog value from specified pin
func (c *Adaptor) AnalogRead(pin string) (val int, err error) {
	analogPin, err := c.translateAnalogPin(pin)
	if err != nil {
		return
	}
	fi, err := c.sys.OpenFile(fmt.Sprintf("%v/%v", c.analogPath, analogPin), os.O_RDONLY, 0o644)
	defer fi.Close() //nolint:staticcheck // for historical reasons

	if err != nil {
		return
	}

	buf := make([]byte, 1024)
	_, err = fi.Read(buf)
	if err != nil {
		return
	}

	val, _ = strconv.Atoi(strings.Split(string(buf), "\n")[0])
	return
}

func (c *Adaptor) validateSpiBusNumber(busNr int) error {
	// Valid bus numbers are [0,1] which corresponds to /dev/spidev0.x through /dev/spidev1.x.
	// x is the chip number <255
	if (busNr < 0) || (busNr > 1) {
		return fmt.Errorf("Bus number %d out of range", busNr)
	}
	return nil
}

func (c *Adaptor) validateI2cBusNumber(busNr int) error {
	// Valid bus number is either 0 or 2 which corresponds to /dev/i2c-0 or /dev/i2c-2.
	if (busNr != 0) && (busNr != 2) {
		return fmt.Errorf("Bus number %d out of range", busNr)
	}
	return nil
}

// translateAnalogPin converts analog pin name to pin position
func (c *Adaptor) translateAnalogPin(pin string) (string, error) {
	if val, ok := c.analogPinMap[pin]; ok {
		return val, nil
	}
	return "", fmt.Errorf("Not a valid analog pin")
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

func (c *Adaptor) translateAndMuxPWMPin(id string) (string, int, error) {
	pinInfo, ok := c.pwmPinMap[id]
	if !ok {
		return "", -1, fmt.Errorf("'%s' is not a valid id for a PWM pin", id)
	}

	path, err := pinInfo.findPWMDir(c.sys)
	if err != nil {
		return "", -1, err
	}

	if err := c.muxPin(id, "pwm"); err != nil {
		return "", -1, err
	}

	return path, pinInfo.channel, nil
}

func (p pwmPinData) findPWMDir(sys *system.Accesser) (dir string, err error) {
	items, _ := sys.Find(p.dir, p.dirRegexp)
	if len(items) == 0 {
		return "", fmt.Errorf("No path found for PWM directory pattern, '%s' in path '%s'", p.dirRegexp, p.dir)
	}

	dir = items[0]
	info, err := sys.Stat(dir)
	if err != nil {
		return "", fmt.Errorf("Error (%v) on access '%s'", err, dir)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("The item '%s' is not a directory, which is not expected", dir)
	}

	return
}

func (c *Adaptor) muxPin(pin, cmd string) error {
	path := fmt.Sprintf("/sys/devices/platform/ocp/ocp:%s_pinmux/state", pin)
	fi, e := c.sys.OpenFile(path, os.O_WRONLY, 0o666)
	defer fi.Close() //nolint:staticcheck // for historical reasons
	if e != nil {
		return e
	}
	_, e = fi.WriteString(cmd)
	return e
}
