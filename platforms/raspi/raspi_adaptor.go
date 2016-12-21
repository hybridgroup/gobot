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
	"gobot.io/x/gobot/sysfs"
)

var readFile = func() ([]byte, error) {
	return ioutil.ReadFile("/proc/cpuinfo")
}

// Adaptor is the Gobot Adaptor for the Raspberry Pi
type Adaptor struct {
	name        string
	revision    string
	i2cLocation string
	digitalPins map[int]sysfs.DigitalPin
	pwmPins     []int
	i2cDevice   sysfs.I2cDevice
}

// NewAdaptor creates a Raspi Adaptor
func NewAdaptor() *Adaptor {
	r := &Adaptor{
		name:        "RaspberryPi",
		digitalPins: make(map[int]sysfs.DigitalPin),
		pwmPins:     []int{},
	}
	content, _ := readFile()
	for _, v := range strings.Split(string(content), "\n") {
		if strings.Contains(v, "Revision") {
			s := strings.Split(string(v), " ")
			version, _ := strconv.ParseInt("0x"+s[len(s)-1], 0, 64)
			r.i2cLocation = "/dev/i2c-1"
			if version <= 3 {
				r.revision = "1"
				r.i2cLocation = "/dev/i2c-0"
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
func (r *Adaptor) Name() string { return r.name }

// SetName sets the Adaptor's name
func (r *Adaptor) SetName(n string) { r.name = n }

// Connect starts connection with board and creates
// digitalPins and pwmPins adaptor maps
func (r *Adaptor) Connect() (err error) {
	return
}

// Finalize closes connection to board and pins
func (r *Adaptor) Finalize() (err error) {
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
	if r.i2cDevice != nil {
		if perr := r.i2cDevice.Close(); err != nil {
			err = multierror.Append(err, perr)
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

	newPin := true
	for _, pin := range r.pwmPins {
		if i == pin {
			newPin = false
			return
		}
	}

	if newPin {
		r.pwmPins = append(r.pwmPins, i)
	}

	return
}

// digitalPin returns matched digitalPin for specified values
func (r *Adaptor) digitalPin(pin string, dir string) (sysfsPin sysfs.DigitalPin, err error) {
	i, err := r.translatePin(pin)

	if err != nil {
		return
	}

	if r.digitalPins[i] == nil {
		r.digitalPins[i] = sysfs.NewDigitalPin(i)
		if err = r.digitalPins[i].Export(); err != nil {
			return
		}
	}

	if err = r.digitalPins[i].Direction(dir); err != nil {
		return
	}

	return r.digitalPins[i], nil
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

// I2cStart starts a i2c device in specified address
func (r *Adaptor) I2cStart(address int) (err error) {
	if r.i2cDevice == nil {
		r.i2cDevice, err = sysfs.NewI2cDevice(r.i2cLocation, address)
	}
	return err
}

// I2CWrite writes data to i2c device
func (r *Adaptor) I2cWrite(address int, data []byte) (err error) {
	if err = r.i2cDevice.SetAddress(address); err != nil {
		return
	}
	_, err = r.i2cDevice.Write(data)
	return
}

// I2cRead returns value from i2c device using specified size
func (r *Adaptor) I2cRead(address int, size int) (data []byte, err error) {
	if err = r.i2cDevice.SetAddress(address); err != nil {
		return
	}
	data = make([]byte, size)
	_, err = r.i2cDevice.Read(data)
	return
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
