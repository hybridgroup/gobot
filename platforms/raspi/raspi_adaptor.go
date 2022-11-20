package raspi

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/drivers/spi"
	"gobot.io/x/gobot/system"
)

const infoFile = "/proc/cpuinfo"

// Adaptor is the Gobot Adaptor for the Raspberry Pi
type Adaptor struct {
	name               string
	mutex              sync.Mutex
	sys                *system.Accesser
	revision           string
	digitalPins        map[int]gobot.DigitalPinner
	pwmPins            map[int]gobot.PWMPinner
	i2cBuses           [2]i2c.I2cDevice
	spiDevices         [2]spi.Connection
	spiDefaultMaxSpeed int64
	PiBlasterPeriod    uint32
}

// NewAdaptor creates a Raspi Adaptor
func NewAdaptor() *Adaptor {
	r := &Adaptor{
		name:            gobot.DefaultName("RaspberryPi"),
		sys:             system.NewAccesser(),
		digitalPins:     make(map[int]gobot.DigitalPinner),
		pwmPins:         make(map[int]gobot.PWMPinner),
		PiBlasterPeriod: 10000000,
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

// Connect do nothing at the moment
func (r *Adaptor) Connect() error { return nil }

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
	for _, dev := range r.spiDevices {
		if dev != nil {
			if e := dev.Close(); e != nil {
				err = multierror.Append(err, e)
			}
		}
	}
	return
}

// DigitalPin returns matched digitalPin for specified values
func (r *Adaptor) DigitalPin(pin string, dir string) (gobot.DigitalPinner, error) {
	i, err := r.translatePin(pin)

	if err != nil {
		return nil, err
	}

	currentPin, err := r.getExportedDigitalPin(i, dir)

	if err != nil {
		return nil, err
	}

	if err = currentPin.Direction(dir); err != nil {
		return nil, err
	}

	return currentPin, nil
}

func (r *Adaptor) getExportedDigitalPin(translatedPin int, dir string) (gobot.DigitalPinner, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.digitalPins[translatedPin] == nil {
		r.digitalPins[translatedPin] = r.sys.NewDigitalPin(translatedPin)
		if err := r.digitalPins[translatedPin].Export(); err != nil {
			return nil, err
		}
	}

	return r.digitalPins[translatedPin], nil
}

// DigitalRead reads digital value from pin
func (r *Adaptor) DigitalRead(pin string) (val int, err error) {
	sysPin, err := r.DigitalPin(pin, system.IN)
	if err != nil {
		return
	}
	return sysPin.Read()
}

// DigitalWrite writes digital value to specified pin
func (r *Adaptor) DigitalWrite(pin string, val byte) error {
	sysPin, err := r.DigitalPin(pin, system.OUT)
	if err != nil {
		return err
	}
	return sysPin.Write(int(val))
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

func (r *Adaptor) getI2cBus(bus int) (_ i2c.I2cDevice, err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.i2cBuses[bus] == nil {
		r.i2cBuses[bus], err = r.sys.NewI2cDevice(fmt.Sprintf("/dev/i2c-%d", bus))
	}

	return r.i2cBuses[bus], err
}

// GetDefaultBus returns the default i2c bus for this platform
func (r *Adaptor) GetDefaultBus() int {
	rev := r.readRevision()
	if rev == "2" || rev == "3" {
		return 1
	}
	return 0
}

// GetSpiConnection returns an spi connection to a device on a specified bus.
// Valid bus number is [0..1] which corresponds to /dev/spidev0.0 through /dev/spidev0.1.
func (r *Adaptor) GetSpiConnection(busNum, chipNum, mode, bits int, maxSpeed int64) (connection spi.Connection, err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if (busNum < 0) || (busNum > 1) {
		return nil, fmt.Errorf("Bus number %d out of range", busNum)
	}

	if r.spiDevices[busNum] == nil {
		r.spiDevices[busNum], err = spi.GetSpiConnection(busNum, chipNum, mode, bits, maxSpeed)
	}

	return r.spiDevices[busNum], err
}

// GetSpiDefaultBus returns the default spi bus for this platform.
func (r *Adaptor) GetSpiDefaultBus() int {
	return 0
}

// GetSpiDefaultChip returns the default spi chip for this platform.
func (r *Adaptor) GetSpiDefaultChip() int {
	return 0
}

// GetSpiDefaultMode returns the default spi mode for this platform.
func (r *Adaptor) GetSpiDefaultMode() int {
	return 0
}

// GetSpiDefaultBits returns the default spi number of bits for this platform.
func (r *Adaptor) GetSpiDefaultBits() int {
	return 8
}

// GetSpiDefaultMaxSpeed returns the default spi bus for this platform.
func (r *Adaptor) GetSpiDefaultMaxSpeed() int64 {
	return 500000
}

// PWMPin returns a raspi.PWMPin which provides the gobot.PWMPinner interface
func (r *Adaptor) PWMPin(pin string) (gobot.PWMPinner, error) {
	i, err := r.translatePin(pin)
	if err != nil {
		return nil, err
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.pwmPins[i] == nil {
		r.pwmPins[i] = NewPWMPin(r.sys, "/dev/pi-blaster", strconv.Itoa(i))
		r.pwmPins[i].SetPeriod(r.PiBlasterPeriod)
	}

	return r.pwmPins[i], nil
}

// PwmWrite writes a PWM signal to the specified pin
func (r *Adaptor) PwmWrite(pin string, val byte) (err error) {
	sysPin, err := r.PWMPin(pin)
	if err != nil {
		return err
	}

	duty := uint32(gobot.FromScale(float64(val), 0, 255) * float64(r.PiBlasterPeriod))
	return sysPin.SetDutyCycle(duty)
}

// ServoWrite writes a servo signal to the specified pin
func (r *Adaptor) ServoWrite(pin string, angle byte) (err error) {
	sysPin, err := r.PWMPin(pin)
	if err != nil {
		return err
	}

	duty := uint32(gobot.FromScale(float64(angle), 0, 180) * float64(r.PiBlasterPeriod))
	return sysPin.SetDutyCycle(duty)
}

func (r *Adaptor) translatePin(pin string) (i int, err error) {
	if val, ok := pins[pin][r.readRevision()]; ok {
		i = val
	} else if val, ok := pins[pin]["*"]; ok {
		i = val
	} else {
		err = errors.New("Not a valid pin")
		return
	}
	return
}

func (r *Adaptor) readRevision() string {
	if r.revision == "" {
		r.revision = "0"
		content, err := r.sys.ReadFile(infoFile)
		if err != nil {
			return r.revision
		}
		for _, v := range strings.Split(string(content), "\n") {
			if strings.Contains(v, "Revision") {
				s := strings.Split(string(v), " ")
				version, _ := strconv.ParseInt("0x"+s[len(s)-1], 0, 64)
				if version <= 3 {
					r.revision = "1"
				} else if version <= 15 {
					r.revision = "2"
				} else {
					r.revision = "3"
				}
			}
		}
	}

	return r.revision
}
