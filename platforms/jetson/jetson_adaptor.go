package jetson

import (
	"errors"
	"fmt"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/drivers/spi"
	"gobot.io/x/gobot/system"
)

const pwmDefaultPeriod = 3000000

// Adaptor is the Gobot Adaptor for the Jetson Nano
type Adaptor struct {
	name        string
	sys         *system.Accesser
	mutex       sync.Mutex
	digitalPins map[int]gobot.DigitalPinner
	pwmPins     map[int]gobot.PWMPinner
	i2cBuses    [2]i2c.I2cDevice
	spiDevices  [2]spi.Connection
}

// NewAdaptor creates a Raspi Adaptor
func NewAdaptor() *Adaptor {
	j := &Adaptor{
		name:        gobot.DefaultName("JetsonNano"),
		sys:         system.NewAccesser(),
		digitalPins: make(map[int]gobot.DigitalPinner),
		pwmPins:     make(map[int]gobot.PWMPinner),
	}
	return j
}

// Name returns the Adaptor's name
func (j *Adaptor) Name() string {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	return j.name
}

// SetName sets the Adaptor's name
func (j *Adaptor) SetName(n string) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	j.name = n
}

// Connect do nothing at the moment
func (j *Adaptor) Connect() error { return nil }

// Finalize closes connection to board and pins
func (j *Adaptor) Finalize() (err error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	for _, pin := range j.digitalPins {
		if pin != nil {
			if perr := pin.Unexport(); err != nil {
				err = multierror.Append(err, perr)
			}
		}
	}

	for _, pin := range j.pwmPins {
		if pin != nil {
			if perr := pin.Unexport(); err != nil {
				err = multierror.Append(err, perr)
			}
		}
	}

	for _, bus := range j.i2cBuses {
		if bus != nil {
			if e := bus.Close(); e != nil {
				err = multierror.Append(err, e)
			}
		}
	}
	for _, dev := range j.spiDevices {
		if dev != nil {
			if e := dev.Close(); e != nil {
				err = multierror.Append(err, e)
			}
		}
	}
	return
}

// DigitalPin returns matched digitalPin for specified values
func (j *Adaptor) DigitalPin(pin string, dir string) (gobot.DigitalPinner, error) {
	i, err := j.translatePin(pin)

	if err != nil {
		return nil, err
	}

	currentPin, err := j.getExportedDigitalPin(i, dir)

	if err != nil {
		return nil, err
	}

	if err = currentPin.Direction(dir); err != nil {
		return nil, err
	}

	return currentPin, nil
}

// DigitalRead reads digital value from pin
func (j *Adaptor) DigitalRead(pin string) (val int, err error) {
	sysPin, err := j.DigitalPin(pin, system.IN)
	if err != nil {
		return
	}
	return sysPin.Read()
}

// DigitalWrite writes digital value to specified pin
func (j *Adaptor) DigitalWrite(pin string, val byte) (err error) {
	sysPin, err := j.DigitalPin(pin, system.OUT)
	if err != nil {
		return err
	}
	return sysPin.Write(int(val))
}

// GetConnection returns an i2c connection to a device on a specified bus.
// Valid bus number is [0..1] which corresponds to /dev/i2c-0 through /dev/i2c-1.
func (j *Adaptor) GetConnection(address int, bus int) (connection i2c.Connection, err error) {
	if (bus < 0) || (bus > 1) {
		return nil, fmt.Errorf("Bus number %d out of range", bus)
	}

	device, err := j.getI2cBus(bus)

	return i2c.NewConnection(device, address), err
}

// GetDefaultBus returns the default i2c bus for this platform
func (j *Adaptor) GetDefaultBus() int {
	return 1
}

// GetSpiConnection returns an spi connection to a device on a specified bus.
// Valid bus number is [0..1] which corresponds to /dev/spidev0.0 through /dev/spidev0.1.
func (j *Adaptor) GetSpiConnection(busNum, chipNum, mode, bits int, maxSpeed int64) (connection spi.Connection, err error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	if (busNum < 0) || (busNum > 1) {
		return nil, fmt.Errorf("Bus number %d out of range", busNum)
	}

	if j.spiDevices[busNum] == nil {
		j.spiDevices[busNum], err = spi.GetSpiConnection(busNum, chipNum, mode, bits, maxSpeed)
	}

	return j.spiDevices[busNum], err
}

// GetSpiDefaultBus returns the default spi bus for this platform.
func (j *Adaptor) GetSpiDefaultBus() int {
	return 0
}

// GetSpiDefaultChip returns the default spi chip for this platform.
func (j *Adaptor) GetSpiDefaultChip() int {
	return 0
}

// GetSpiDefaultMode returns the default spi mode for this platform.
func (j *Adaptor) GetSpiDefaultMode() int {
	return 0
}

// GetSpiDefaultBits returns the default spi number of bits for this platform.
func (j *Adaptor) GetSpiDefaultBits() int {
	return 8
}

// GetSpiDefaultMaxSpeed returns the default spi bus for this platform.
func (j *Adaptor) GetSpiDefaultMaxSpeed() int64 {
	return 10000000
}

//PWMPin returns a Jetson Nano. PWMPin which provides the gobot.PWMPinner interface
func (j *Adaptor) PWMPin(pin string) (gobot.PWMPinner, error) {
	i, err := j.translatePin(pin)
	if err != nil {
		return nil, err
	}
	j.mutex.Lock()
	defer j.mutex.Unlock()

	if j.pwmPins[i] != nil {
		return j.pwmPins[i], nil
	}

	j.pwmPins[i], err = NewPWMPin(j.sys, "/sys/class/pwm/pwmchip0", pin)
	if err != nil {
		return nil, err
	}
	j.pwmPins[i].Export()
	j.pwmPins[i].SetPeriod(pwmDefaultPeriod)
	j.pwmPins[i].Enable(true)

	return j.pwmPins[i], nil
}

// PwmWrite writes a PWM signal to the specified pin
func (j *Adaptor) PwmWrite(pin string, val byte) (err error) {
	sysPin, err := j.PWMPin(pin)
	if err != nil {
		return err
	}

	duty := uint32(gobot.FromScale(float64(val), 0, 255) * float64(pwmDefaultPeriod))
	return sysPin.SetDutyCycle(duty)
}

// ServoWrite writes a servo signal to the specified pin
func (j *Adaptor) ServoWrite(pin string, angle byte) (err error) {
	sysPin, err := j.PWMPin(pin)
	if err != nil {
		return err
	}

	duty := uint32(gobot.FromScale(float64(angle), 0, 180) * float64(pwmDefaultPeriod))
	return sysPin.SetDutyCycle(duty)
}

func (j *Adaptor) getExportedDigitalPin(translatedPin int, dir string) (gobot.DigitalPinner, error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	if j.digitalPins[translatedPin] == nil {
		j.digitalPins[translatedPin] = j.sys.NewDigitalPin(translatedPin)
		if err := j.digitalPins[translatedPin].Export(); err != nil {
			return nil, err
		}
	}

	return j.digitalPins[translatedPin], nil
}

func (j *Adaptor) getI2cBus(bus int) (_ i2c.I2cDevice, err error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	if j.i2cBuses[bus] == nil {
		j.i2cBuses[bus], err = j.sys.NewI2cDevice(fmt.Sprintf("/dev/i2c-%d", bus))
	}

	return j.i2cBuses[bus], err
}

func (j *Adaptor) translatePin(pin string) (i int, err error) {
	if val, ok := pins[pin]["*"]; ok {
		i = val
	} else {
		err = errors.New("Not a valid pin")
		return
	}
	return
}
