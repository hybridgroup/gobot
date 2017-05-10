package raspi

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/sysfs"
	"sync"
)

var readFile = func() ([]byte, error) {
	return ioutil.ReadFile("/proc/cpuinfo")
}

// Adaptor is the Gobot Adaptor for the Raspberry Pi
type Adaptor struct {
	mutex         *sync.Mutex
	name          string
	revision      string
	digitalPins   map[int]*sysfs.DigitalPin
	pwmPins       map[int]*PWMPin
	i2cDefaultBus int
	i2cBuses      [2]sysfs.I2cDevice
}

// NewAdaptor creates a Raspi Adaptor
func NewAdaptor() *Adaptor {
	r := &Adaptor{
		mutex:       &sync.Mutex{},
		name:        gobot.DefaultName("RaspberryPi"),
		digitalPins: make(map[int]*sysfs.DigitalPin),
		pwmPins:     make(map[int]*PWMPin),
	}
	content, _ := readFile()
	for _, v := range strings.Split(string(content), "\n") {
		if strings.Contains(v, "Revision") {
			s := strings.Split(string(v), " ")
			version, _ := strconv.ParseInt("0x"+s[len(s)-1], 0, 64)
			r.i2cDefaultBus = 1
			if version <= 3 {
				r.revision = "1"
				r.i2cDefaultBus = 0
			} else if version <= 15 {
				r.revision = "2"
			} else {
				r.revision = "3"
			}
		}
	}

	return r
}

// Name returns the Adaptor's name
func (r *Adaptor) Name() string {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.name
}

// SetName sets the Adaptor's name
func (r *Adaptor) SetName(n string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.name = n
}

// Connect starts connection with board and creates
// digitalPins and pwmPins adaptor maps
func (r *Adaptor) Connect() (err error) {
	return
}

// Finalize closes connection to board and pins
func (r *Adaptor) Finalize() (err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, pin := range r.digitalPins {
		if pin != nil {
			if perr := pin.Unexport(); err != nil {
				err = multierror.Append(err, perr)
			}
		}
	}
	for _, pin := range r.pwmPins {
		if pin != nil {
			if perr := pin.Unexport(); err != nil {
				err = multierror.Append(err, perr)
			}
		}
	}
	for _, bus := range r.i2cBuses {
		if bus != nil {
			if e := bus.Close(); e != nil {
				err = multierror.Append(err, e)
			}
		}
	}
	return
}

// DigitalPin returns matched digitalPin for specified values
func (r *Adaptor) DigitalPin(pin string, dir string) (sysfsPin sysfs.DigitalPinner, err error) {
	i, err := r.translatePin(pin)

	if err != nil {
		return
	}

	currentPin, err := r.getExportedDigitalPin(i, dir)

	if err != nil {
		return
	}

	if err = currentPin.Direction(dir); err != nil {
		return
	}

	return currentPin, nil
}

func (r *Adaptor) getExportedDigitalPin(translatedPin int, dir string) (sysfsPin sysfs.DigitalPinner, err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.digitalPins[translatedPin] == nil {
		r.digitalPins[translatedPin] = sysfs.NewDigitalPin(translatedPin)
		if err = r.digitalPins[translatedPin].Export(); err != nil {
			return
		}
	}

	return r.digitalPins[translatedPin], nil
}

// DigitalRead reads digital value from pin
func (r *Adaptor) DigitalRead(pin string) (val int, err error) {
	sysfsPin, err := r.DigitalPin(pin, sysfs.IN)
	if err != nil {
		return
	}
	return sysfsPin.Read()
}

// DigitalWrite writes digital value to specified pin
func (r *Adaptor) DigitalWrite(pin string, val byte) (err error) {
	sysfsPin, err := r.DigitalPin(pin, sysfs.OUT)
	if err != nil {
		return err
	}
	return sysfsPin.Write(int(val))
}

// GetConnection returns an i2c connection to a device on a specified bus.
// Valid bus number is [0..1] which corresponds to /dev/i2c-0 through /dev/i2c-1.
func (r *Adaptor) GetConnection(address int, bus int) (connection i2c.Connection, err error) {
	if (bus < 0) || (bus > 1) {
		return nil, fmt.Errorf("Bus number %d out of range", bus)
	}

	device, err := r.getI2cBus(bus)

	return i2c.NewConnection(device, address), err
}

func (r *Adaptor) getI2cBus(bus int) (_ sysfs.I2cDevice, err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.i2cBuses[bus] == nil {
		r.i2cBuses[bus], err = sysfs.NewI2cDevice(fmt.Sprintf("/dev/i2c-%d", bus))
	}

	return r.i2cBuses[bus], err
}

// GetDefaultBus returns the default i2c bus for this platform
func (r *Adaptor) GetDefaultBus() int {
	return r.i2cDefaultBus
}

// PWMPin returns a raspi.PWMPin which provides the sysfs.PWMPinner interface
func (r *Adaptor) PWMPin(pin string) (raspiPWMPin sysfs.PWMPinner, err error) {
	i, err := r.translatePin(pin)
	if err != nil {
		return
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.pwmPins[i] == nil {
		r.pwmPins[i] = NewPWMPin(strconv.Itoa(i))
	}

	return r.pwmPins[i], nil
}

// PwmWrite writes a PWM signal to the specified pin
func (r *Adaptor) PwmWrite(pin string, val byte) (err error) {
	sysfsPin, err := r.PWMPin(pin)
	if err != nil {
		return err
	}

	duty := uint32(gobot.FromScale(float64(val), 0, 255) * piBlasterPeriod)
	return sysfsPin.SetDutyCycle(duty)
}

// SetPwmPeriod writes a PWM signal of the given period to the specified pin
func (r *Adaptor) SetPwmPeriod(pin string, period float64) (err error) {
	sysfsPin, err := r.pwmPin(pin)
	if err != nil {
		return err
	}
	const maxPeriod float64 = 0.01 //10000us
	const minPeriod float64 = 0    //10us
	if period < minPeriod || period > maxPeriod {
		return errors.New("Out of bound: the PWM value must be between 0 and 1")
	}

	val := (period - minPeriod) / (maxPeriod - minPeriod)
	return r.piBlaster(fmt.Sprintf("%v=%v\n", sysfsPin, val))
}

// ServoWrite writes a servo signal to the specified pin
func (r *Adaptor) ServoWrite(pin string, angle byte) (err error) {
	sysfsPin, err := r.PWMPin(pin)
	if err != nil {
		return err
	}

	duty := uint32(gobot.FromScale(float64(angle), 0, 180) * piBlasterPeriod)
	return sysfsPin.SetDutyCycle(duty)
}

func (r *Adaptor) translatePin(pin string) (i int, err error) {
	if val, ok := pins[pin][r.revision]; ok {
		i = val
	} else if val, ok := pins[pin]["*"]; ok {
		i = val
	} else {
		err = errors.New("Not a valid pin")
		return
	}
	return
}
