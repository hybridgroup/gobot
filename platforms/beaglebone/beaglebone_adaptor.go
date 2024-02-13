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

type pwmPinDefinition struct {
	channel   int
	dir       string
	dirRegexp string
}

type analogPinDefinition struct {
	path   string
	r      bool // readable
	w      bool // writable
	bufLen uint16
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
	mutex *sync.Mutex
	*adaptors.AnalogPinsAdaptor
	*adaptors.DigitalPinsAdaptor
	*adaptors.PWMPinsAdaptor
	*adaptors.I2cBusAdaptor
	*adaptors.SpiBusAdaptor
	usrLed       string
	pinMap       map[string]int
	pwmPinMap    map[string]pwmPinDefinition
	analogPinMap map[string]analogPinDefinition
}

// NewAdaptor returns a new Beaglebone Black/Green Adaptor
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
		name:         gobot.DefaultName("BeagleboneBlack"),
		sys:          sys,
		mutex:        &sync.Mutex{},
		pinMap:       bbbPinMap,
		pwmPinMap:    bbbPwmPinMap,
		analogPinMap: bbbAnalogPinMap,
		usrLed:       "/sys/class/leds/beaglebone:green:",
	}

	var digitalPinsOpts []func(adaptors.DigitalPinsOptioner)
	pwmPinsOpts := []adaptors.PwmPinsOptionApplier{adaptors.WithPWMDefaultPeriod(pwmPeriodDefault)}
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

	a.AnalogPinsAdaptor = adaptors.NewAnalogPinsAdaptor(sys, a.translateAnalogPin)
	a.DigitalPinsAdaptor = adaptors.NewDigitalPinsAdaptor(sys, a.translateAndMuxDigitalPin, digitalPinsOpts...)
	a.PWMPinsAdaptor = adaptors.NewPWMPinsAdaptor(sys, a.translateAndMuxPWMPin, pwmPinsOpts...)
	a.I2cBusAdaptor = adaptors.NewI2cBusAdaptor(sys, a.validateI2cBusNumber, defaultI2cBusNumber)
	a.SpiBusAdaptor = adaptors.NewSpiBusAdaptor(sys, a.validateSpiBusNumber, defaultSpiBusNumber, defaultSpiChipNumber,
		defaultSpiMode, defaultSpiBitsNumber, defaultSpiMaxSpeed)
	return a
}

// Name returns the Adaptor name
func (a *Adaptor) Name() string { return a.name }

// SetName sets the Adaptor name
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

	if err := a.AnalogPinsAdaptor.Connect(); err != nil {
		return err
	}

	if err := a.PWMPinsAdaptor.Connect(); err != nil {
		return err
	}
	return a.DigitalPinsAdaptor.Connect()
}

// Finalize releases all i2c devices and exported analog, digital, pwm pins.
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

// DigitalWrite writes a digital value to specified pin.
// valid usr pin values are usr0, usr1, usr2 and usr3
func (a *Adaptor) DigitalWrite(id string, val byte) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if strings.Contains(id, "usr") {
		fi, e := a.sys.OpenFile(a.usrLed+id+"/brightness", os.O_WRONLY|os.O_APPEND, 0o666)
		defer fi.Close() //nolint:staticcheck // for historical reasons
		if e != nil {
			return e
		}
		_, err := fi.WriteString(strconv.Itoa(int(val)))
		return err
	}

	return a.DigitalPinsAdaptor.DigitalWrite(id, val)
}

func (a *Adaptor) validateSpiBusNumber(busNr int) error {
	// Valid bus numbers are [0,1] which corresponds to /dev/spidev0.x through /dev/spidev1.x.
	// x is the chip number <255
	if (busNr < 0) || (busNr > 1) {
		return fmt.Errorf("Bus number %d out of range", busNr)
	}
	return nil
}

func (a *Adaptor) validateI2cBusNumber(busNr int) error {
	// Valid bus number is either 0 or 2 which corresponds to /dev/i2c-0 or /dev/i2c-2.
	if (busNr != 0) && (busNr != 2) {
		return fmt.Errorf("Bus number %d out of range", busNr)
	}
	return nil
}

// translateAnalogPin converts analog pin name to pin position
func (a *Adaptor) translateAnalogPin(pin string) (string, bool, bool, uint16, error) {
	pinInfo, ok := a.analogPinMap[pin]
	if !ok {
		return "", false, false, 0, fmt.Errorf("Not a valid analog pin")
	}

	return pinInfo.path, pinInfo.r, pinInfo.w, pinInfo.bufLen, nil
}

// translatePin converts digital pin name to pin position
func (a *Adaptor) translateAndMuxDigitalPin(id string) (string, int, error) {
	line, ok := a.pinMap[id]
	if !ok {
		return "", -1, fmt.Errorf("'%s' is not a valid id for a digital pin", id)
	}
	// mux is done by id, not by line
	if err := a.muxPin(id, "gpio"); err != nil {
		return "", -1, err
	}
	return "", line, nil
}

func (a *Adaptor) translateAndMuxPWMPin(id string) (string, int, error) {
	pinInfo, ok := a.pwmPinMap[id]
	if !ok {
		return "", -1, fmt.Errorf("'%s' is not a valid id for a PWM pin", id)
	}

	path, err := pinInfo.findPWMDir(a.sys)
	if err != nil {
		return "", -1, err
	}

	if err := a.muxPin(id, "pwm"); err != nil {
		return "", -1, err
	}

	return path, pinInfo.channel, nil
}

func (p pwmPinDefinition) findPWMDir(sys *system.Accesser) (string, error) {
	items, _ := sys.Find(p.dir, p.dirRegexp)
	if len(items) == 0 {
		return "", fmt.Errorf("No path found for PWM directory pattern, '%s' in path '%s'", p.dirRegexp, p.dir)
	}

	dir := items[0]
	info, err := sys.Stat(dir)
	if err != nil {
		return "", fmt.Errorf("Error (%v) on access '%s'", err, dir)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("The item '%s' is not a directory, which is not expected", dir)
	}

	return dir, nil
}

func (a *Adaptor) muxPin(pin, cmd string) error {
	path := fmt.Sprintf("/sys/devices/platform/ocp/ocp:%s_pinmux/state", pin)
	fi, e := a.sys.OpenFile(path, os.O_WRONLY, 0o666)
	defer fi.Close() //nolint:staticcheck // for historical reasons
	if e != nil {
		return e
	}
	_, e = fi.WriteString(cmd)
	return e
}
