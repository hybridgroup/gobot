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
	"gobot.io/x/gobot/sysfs"
)

type pwmPinData struct {
	channel int
	path    string
}

const pwmDefaultPeriod = 500000

// Adaptor is the gobot.Adaptor representation for the Beaglebone Black/Green
type Adaptor struct {
	name         string
	sysfs        *sysfs.Accesser
	mutex        sync.Mutex
	digitalPins  []sysfs.DigitalPinner
	pwmPins      map[string]sysfs.PWMPinner
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
	b := &Adaptor{
		name:         gobot.DefaultName("BeagleboneBlack"),
		sysfs:        sysfs.NewAccesser(),
		digitalPins:  make([]sysfs.DigitalPinner, 120),
		pwmPins:      make(map[string]sysfs.PWMPinner),
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

	return b
}

// Name returns the Adaptor name
func (b *Adaptor) Name() string { return b.name }

// SetName sets the Adaptor name
func (b *Adaptor) SetName(n string) { b.name = n }

// Connect do nothing at the moment
func (b *Adaptor) Connect() error { return nil }

// Finalize releases all i2c devices and exported analog, digital, pwm pins.
func (b *Adaptor) Finalize() (err error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for _, pin := range b.digitalPins {
		if pin != nil {
			if e := pin.Unexport(); e != nil {
				err = multierror.Append(err, e)
			}
		}
	}
	for _, pin := range b.pwmPins {
		if pin != nil {
			if e := pin.Unexport(); e != nil {
				err = multierror.Append(err, e)
			}
		}
	}
	for _, bus := range b.i2cBuses {
		if bus != nil {
			if e := bus.Close(); e != nil {
				err = multierror.Append(err, e)
			}
		}
	}
	for _, bus := range b.spiBuses {
		if bus != nil {
			if e := bus.Close(); e != nil {
				err = multierror.Append(err, e)
			}
		}
	}
	return
}

// PwmWrite writes the 0-254 value to the specified pin
func (b *Adaptor) PwmWrite(pin string, val byte) (err error) {
	pwmPin, err := b.PWMPin(pin)
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
func (b *Adaptor) ServoWrite(pin string, angle byte) (err error) {
	pwmPin, err := b.PWMPin(pin)
	if err != nil {
		return
	}

	// TODO: take into account the actual period setting, not just assume default
	const minDuty = 100 * 0.0005 * pwmDefaultPeriod
	const maxDuty = 100 * 0.0020 * pwmDefaultPeriod
	duty := uint32(gobot.ToScale(gobot.FromScale(float64(angle), 0, 180), minDuty, maxDuty))
	return pwmPin.SetDutyCycle(duty)
}

// DigitalRead returns a digital value from specified pin
func (b *Adaptor) DigitalRead(pin string) (val int, err error) {
	sysfsPin, err := b.DigitalPin(pin, sysfs.IN)
	if err != nil {
		return
	}
	return sysfsPin.Read()
}

// DigitalWrite writes a digital value to specified pin.
// valid usr pin values are usr0, usr1, usr2 and usr3
func (b *Adaptor) DigitalWrite(pin string, val byte) (err error) {
	if strings.Contains(pin, "usr") {
		fi, e := b.sysfs.OpenFile(b.usrLed+pin+"/brightness", os.O_WRONLY|os.O_APPEND, 0666)
		defer fi.Close()
		if e != nil {
			return e
		}
		_, err = fi.WriteString(strconv.Itoa(int(val)))
		return err
	}
	sysfsPin, err := b.DigitalPin(pin, sysfs.OUT)
	if err != nil {
		return err
	}
	return sysfsPin.Write(int(val))
}

// DigitalPin retrieves digital pin value by name
func (b *Adaptor) DigitalPin(pin string, dir string) (sysfsPin sysfs.DigitalPinner, err error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	i, err := b.translatePin(pin)
	if err != nil {
		return
	}
	if b.digitalPins[i] == nil {
		b.digitalPins[i] = b.sysfs.NewDigitalPin(i)
		if err = b.muxPin(pin, "gpio"); err != nil {
			return
		}

		err := b.digitalPins[i].Export()
		if err != nil {
			return nil, err
		}
	}
	if err = b.digitalPins[i].Direction(dir); err != nil {
		return
	}
	return b.digitalPins[i], nil
}

// PWMPin returns matched pwmPin for specified pin number
func (b *Adaptor) PWMPin(pin string) (sysfs.PWMPinner, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	pinInfo, err := b.translatePwmPin(pin)
	if err != nil {
		return nil, err
	}

	if b.pwmPins[pin] == nil {
		newPath, err := b.findPin(pinInfo.path)
		if err != nil {
			return nil, err
		}
		newPin := b.sysfs.NewPWMPin(newPath, pinInfo.channel)
		if err := b.muxPin(pin, "pwm"); err != nil {
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
		b.pwmPins[pin] = newPin
	}

	return b.pwmPins[pin], nil
}

// AnalogRead returns an analog value from specified pin
func (b *Adaptor) AnalogRead(pin string) (val int, err error) {
	analogPin, err := b.translateAnalogPin(pin)
	if err != nil {
		return
	}
	fi, err := b.sysfs.OpenFile(fmt.Sprintf("%v/%v", b.analogPath, analogPin), os.O_RDONLY, 0644)
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
func (b *Adaptor) GetConnection(address int, bus int) (connection i2c.Connection, err error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if (bus != 0) && (bus != 2) {
		return nil, fmt.Errorf("Bus number %d out of range", bus)
	}
	if b.i2cBuses[bus] == nil {
		b.i2cBuses[bus], err = b.sysfs.NewI2cDevice(fmt.Sprintf("/dev/i2c-%d", bus))
	}
	return i2c.NewConnection(b.i2cBuses[bus], address), err
}

// GetDefaultBus returns the default i2c bus for this platform
func (b *Adaptor) GetDefaultBus() int {
	return 2
}

// GetSpiConnection returns an spi connection to a device on a specified bus.
// Valid bus number is [0..1] which corresponds to /dev/spidev0.0 through /dev/spidev0.1.
func (b *Adaptor) GetSpiConnection(busNum, chipNum, mode, bits int, maxSpeed int64) (connection spi.Connection, err error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if (busNum < 0) || (busNum > 1) {
		return nil, fmt.Errorf("Bus number %d out of range", busNum)
	}

	if b.spiBuses[busNum] == nil {
		b.spiBuses[busNum], err = spi.GetSpiConnection(busNum, chipNum, mode, bits, maxSpeed)
	}

	return b.spiBuses[busNum], err
}

// GetSpiDefaultBus returns the default spi bus for this platform.
func (b *Adaptor) GetSpiDefaultBus() int {
	return 0
}

// GetSpiDefaultChip returns the default spi chip for this platform.
func (b *Adaptor) GetSpiDefaultChip() int {
	return 0
}

// GetSpiDefaultMode returns the default spi mode for this platform.
func (b *Adaptor) GetSpiDefaultMode() int {
	return 0
}

// GetSpiDefaultBits returns the default spi number of bits for this platform.
func (b *Adaptor) GetSpiDefaultBits() int {
	return 8
}

// GetSpiDefaultMaxSpeed returns the default spi bus for this platform.
func (b *Adaptor) GetSpiDefaultMaxSpeed() int64 {
	return 500000
}

// translatePin converts digital pin name to pin position
func (b *Adaptor) translatePin(pin string) (value int, err error) {
	if val, ok := b.pinMap[pin]; ok {
		value = val
	} else {
		err = errors.New("Not a valid pin")
	}
	return
}

func (b *Adaptor) translatePwmPin(pin string) (p pwmPinData, err error) {
	if val, ok := b.pwmPinMap[pin]; ok {
		p = val
	} else {
		err = errors.New("Not a valid PWM pin")
	}
	return
}

// translateAnalogPin converts analog pin name to pin position
func (b *Adaptor) translateAnalogPin(pin string) (value string, err error) {
	if val, ok := b.analogPinMap[pin]; ok {
		value = val
	} else {
		err = errors.New("Not a valid analog pin")
	}
	return
}

func (b *Adaptor) muxPin(pin, cmd string) error {
	path := fmt.Sprintf("/sys/devices/platform/ocp/ocp:%s_pinmux/state", pin)
	fi, e := b.sysfs.OpenFile(path, os.O_WRONLY, 0666)
	defer fi.Close()
	if e != nil {
		return e
	}
	_, e = fi.WriteString(cmd)
	return e
}
