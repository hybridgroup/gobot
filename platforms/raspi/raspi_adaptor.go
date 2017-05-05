package raspi

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
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
	digitalPins   map[int]sysfs.DigitalPin
	pwmPins       []int
	i2cDefaultBus int
	i2cBuses      [2]sysfs.I2cDevice
}

// NewAdaptor creates a Raspi Adaptor
func NewAdaptor() *Adaptor {
	r := &Adaptor{
		mutex:       &sync.Mutex{},
		name:        gobot.DefaultName("RaspberryPi"),
		digitalPins: make(map[int]sysfs.DigitalPin),
		pwmPins:     []int{},
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
		if perr := r.piBlaster(fmt.Sprintf("release %v\n", pin)); err != nil {
			err = multierror.Append(err, perr)
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

func (r *Adaptor) pwmPin(pin string) (i int, err error) {
	i, err = r.translatePin(pin)
	if err != nil {
		return
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, pin := range r.pwmPins {
		if i == pin {
			return
		}
	}

	r.pwmPins = append(r.pwmPins, i)

	return
}

// digitalPin returns matched digitalPin for specified values
func (r *Adaptor) digitalPin(pin string, dir string) (sysfsPin sysfs.DigitalPin, err error) {
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

func (r *Adaptor) getExportedDigitalPin(translatedPin int, dir string) (sysfsPin sysfs.DigitalPin, err error) {
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
	sysfsPin, err := r.digitalPin(pin, sysfs.IN)
	if err != nil {
		return
	}
	return sysfsPin.Read()
}

// DigitalWrite writes digital value to specified pin
func (r *Adaptor) DigitalWrite(pin string, val byte) (err error) {
	sysfsPin, err := r.digitalPin(pin, sysfs.OUT)
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

// PwmWrite writes a PWM signal to the specified pin
func (r *Adaptor) PwmWrite(pin string, val byte) (err error) {
	sysfsPin, err := r.pwmPin(pin)
	if err != nil {
		return err
	}
	return r.piBlaster(fmt.Sprintf("%v=%v\n", sysfsPin, gobot.FromScale(float64(val), 0, 255)))
}

// ServoWrite writes a servo signal to the specified pin
func (r *Adaptor) ServoWrite(pin string, angle byte) (err error) {
	sysfsPin, err := r.pwmPin(pin)
	if err != nil {
		return err
	}

	val := (gobot.ToScale(gobot.FromScale(float64(angle), 0, 180), 0, 200) / 1000.0) + 0.05

	return r.piBlaster(fmt.Sprintf("%v=%v\n", sysfsPin, val))
}

func (r *Adaptor) piBlaster(data string) (err error) {
	fi, err := sysfs.OpenFile("/dev/pi-blaster", os.O_WRONLY|os.O_APPEND, 0644)
	defer fi.Close()

	if err != nil {
		return err
	}

	_, err = fi.WriteString(data)
	return
}
