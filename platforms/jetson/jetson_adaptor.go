package jetson

import (
	"errors"
	"fmt"

	"sync"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/drivers/spi"
	"gobot.io/x/gobot/sysfs"
)

// Adaptor is the Gobot Adaptor for the Jetson Nano
type Adaptor struct {
	mutex       *sync.Mutex
	name        string
	revision    string
	digitalPins map[int]*sysfs.DigitalPin
	//pwmPins            map[int]*PWMPin
	i2cDefaultBus      int
	i2cBuses           [2]i2c.I2cDevice
	spiDefaultBus      int
	spiDefaultChip     int
	spiDevices         [2]spi.Connection
	spiDefaultMode     int
	spiDefaultMaxSpeed int64
	//JSBlasterPeriod    uint32
}

// NewAdaptor creates a Raspi Adaptor
func NewAdaptor() *Adaptor {
	j := &Adaptor{
		mutex:       &sync.Mutex{},
		name:        gobot.DefaultName("JetsonNano"),
		digitalPins: make(map[int]*sysfs.DigitalPin),
		//pwmPins:         make(map[int]*PWMPin),
		//JSBlasterPeriod: 10000000,
	}

	j.i2cDefaultBus = 1
	j.spiDefaultBus = 0
	j.spiDefaultChip = 0
	j.spiDefaultMode = 0
	j.spiDefaultMaxSpeed = 10000000

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

// Connect starts connection with board and creates
// digitalPins and pwmPins adaptor maps
func (j *Adaptor) Connect() (err error) {
	return
}

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
	/*
		for _, pin := range r.pwmPins {
			if pin != nil {
				if perr := pin.Unexport(); err != nil {
					err = multierror.Append(err, perr)
				}
			}
		}
	*/
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
func (j *Adaptor) DigitalPin(pin string, dir string) (sysfsPin sysfs.DigitalPinner, err error) {
	i, err := j.translatePin(pin)

	if err != nil {
		return
	}

	currentPin, err := j.getExportedDigitalPin(i, dir)

	if err != nil {
		return
	}

	if err = currentPin.Direction(dir); err != nil {
		return
	}

	return currentPin, nil
}

func (j *Adaptor) getExportedDigitalPin(translatedPin int, dir string) (sysfsPin sysfs.DigitalPinner, err error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	if j.digitalPins[translatedPin] == nil {
		j.digitalPins[translatedPin] = sysfs.NewDigitalPin(translatedPin)
		if err = j.digitalPins[translatedPin].Export(); err != nil {
			return
		}
	}

	return j.digitalPins[translatedPin], nil
}

// DigitalRead reads digital value from pin
func (j *Adaptor) DigitalRead(pin string) (val int, err error) {
	sysfsPin, err := j.DigitalPin(pin, sysfs.IN)
	if err != nil {
		return
	}
	return sysfsPin.Read()
}

// DigitalWrite writes digital value to specified pin
func (j *Adaptor) DigitalWrite(pin string, val byte) (err error) {
	sysfsPin, err := j.DigitalPin(pin, sysfs.OUT)
	if err != nil {
		return err
	}
	return sysfsPin.Write(int(val))
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

func (j *Adaptor) getI2cBus(bus int) (_ i2c.I2cDevice, err error) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	if j.i2cBuses[bus] == nil {
		j.i2cBuses[bus], err = sysfs.NewI2cDevice(fmt.Sprintf("/dev/i2c-%d", bus))
	}

	return j.i2cBuses[bus], err
}

// GetDefaultBus returns the default i2c bus for this platform
func (j *Adaptor) GetDefaultBus() int {
	return j.i2cDefaultBus
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
	return j.spiDefaultBus
}

// GetSpiDefaultChip returns the default spi chip for this platform.
func (j *Adaptor) GetSpiDefaultChip() int {
	return j.spiDefaultChip
}

// GetSpiDefaultMode returns the default spi mode for this platform.
func (j *Adaptor) GetSpiDefaultMode() int {
	return j.spiDefaultMode
}

// GetSpiDefaultBits returns the default spi number of bits for this platform.
func (j *Adaptor) GetSpiDefaultBits() int {
	return 8
}

// GetSpiDefaultMaxSpeed returns the default spi bus for this platform.
func (j *Adaptor) GetSpiDefaultMaxSpeed() int64 {
	return j.spiDefaultMaxSpeed
}

func (j *Adaptor) translatePin(pin string) (i int, err error) {
	if val, ok := pins[pin][j.revision]; ok {
		i = val
	} else if val, ok := pins[pin]["*"]; ok {
		i = val
	} else {
		err = errors.New("Not a valid pin")
		return
	}
	return
}

/*
//PWMPin returns a jetson.PWMPin which provides the sysfs.PWMPinner interface
//
func (r *Adaptor) PWMPin(pin string) (raspiPWMPin sysfs.PWMPinner, err error) {
	i, err := r.translatePin(pin)
	if err != nil {
		return
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.pwmPins[i] == nil {
		r.pwmPins[i] = NewPWMPin(strconv.Itoa(i))
		r.pwmPins[i].SetPeriod(r.JSBlasterPeriod)
	}

	return r.pwmPins[i], nil
}

// PwmWrite writes a PWM signal to the specified pin
func (r *Adaptor) PwmWrite(pin string, val byte) (err error) {
	sysfsPin, err := r.PWMPin(pin)
	if err != nil {
		return err
	}

	duty := uint32(gobot.FromScale(float64(val), 0, 255) * float64(r.JSBlasterPeriod))
	return sysfsPin.SetDutyCycle(duty)
}

// ServoWrite writes a servo signal to the specified pin
func (r *Adaptor) ServoWrite(pin string, angle byte) (err error) {
	sysfsPin, err := r.PWMPin(pin)
	if err != nil {
		return err
	}

	duty := uint32(gobot.FromScale(float64(angle), 0, 180) * float64(r.JSBlasterPeriod))
	return sysfsPin.SetDutyCycle(duty)
}
*/
