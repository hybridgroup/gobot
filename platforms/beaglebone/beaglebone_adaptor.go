package beaglebone

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	multierror "github.com/hashicorp/go-multierror"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/sysfs"
)

var glob = func(pattern string) (matches []string, err error) {
	return filepath.Glob(pattern)
}

// Adaptor is the gobot.Adaptor representation for the Beaglebone
type Adaptor struct {
	name         string
	kernel       string
	digitalPins  []sysfs.DigitalPin
	pwmPins      map[string]*pwmPin
	i2cBuses     map[int]sysfs.I2cDevice
	usrLed       string
	ocp          string
	analogPath   string
	analogPinMap map[string]string
	slots        string
}

// NewAdaptor returns a new Beaglebone Adaptor
func NewAdaptor() *Adaptor {
	b := &Adaptor{
		name:        gobot.DefaultName("Beaglebone"),
		digitalPins: make([]sysfs.DigitalPin, 120),
		pwmPins:     make(map[string]*pwmPin),
		i2cBuses:    make(map[int]sysfs.I2cDevice),
	}

	b.setSlots()
	return b
}

func (b *Adaptor) setSlots() {
	ocp := "/sys/devices/ocp.*"
	slots := "/sys/devices/bone_capemgr.*"

	b.kernel = getKernel()
	if b.kernel[:1] == "4" {
		ocp = "/sys/devices/platform/ocp/ocp*"
		slots = "/sys/devices/platform/bone_capemgr"
	}

	b.usrLed = "/sys/class/leds/beaglebone:green:"

	g, _ := glob(ocp)
	b.ocp = g[0]

	g, _ = glob(slots)
	b.slots = fmt.Sprintf("%v/slots", g[0])
}

// Name returns the Adaptor name
func (b *Adaptor) Name() string { return b.name }

// SetName sets the Adaptor name
func (b *Adaptor) SetName(n string) { b.name = n }

// Kernel returns the Linux kernel version for the BeagleBone
func (b *Adaptor) Kernel() string { return b.kernel }

// Connect initializes the pwm and analog dts.
func (b *Adaptor) Connect() error {
	// enable analog
	if b.kernel[:1] == "4" {
		if err := ensureSlot(b.slots, "BB-ADC"); err != nil {
			return err
		}

		b.analogPath = "/sys/bus/iio/devices/iio:device0"
		b.analogPinMap = analogPins44
	} else {
		if err := ensureSlot(b.slots, "cape-bone-iio"); err != nil {
			return err
		}

		g, err := glob(fmt.Sprintf("%v/helper.*", b.ocp))
		if err != nil {
			return err
		}
		b.analogPath = g[0]
		b.analogPinMap = analogPins3
	}

	// enable pwm
	if err := ensureSlot(b.slots, "am33xx_pwm"); err != nil {
		return err
	}

	return nil
}

// Finalize releases all i2c devices and exported analog, digital, pwm pins.
func (b *Adaptor) Finalize() (err error) {
	for _, pin := range b.pwmPins {
		if pin != nil {
			if e := pin.release(); e != nil {
				err = multierror.Append(err, e)
			}
		}
	}
	for _, pin := range b.digitalPins {
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
	return
}

// PwmWrite writes the 0-254 value to the specified pin
func (b *Adaptor) PwmWrite(pin string, val byte) (err error) {
	return b.pwmWrite(pin, val)
}

// ServoWrite writes the 0-180 degree val to the specified pin.
func (b *Adaptor) ServoWrite(pin string, val byte) (err error) {
	i, err := b.pwmPin(pin)
	if err != nil {
		return err
	}
	period := 16666666.0
	duty := (gobot.FromScale(float64(val), 0, 180.0) * 0.115) + 0.05
	return b.pwmPins[i].pwmWrite(strconv.Itoa(int(period)), strconv.Itoa(int(period*duty)))
}

// DigitalRead returns a digital value from specified pin
func (b *Adaptor) DigitalRead(pin string) (val int, err error) {
	sysfsPin, err := b.digitalPin(pin, sysfs.IN)
	if err != nil {
		return
	}
	return sysfsPin.Read()
}

// DigitalWrite writes a digital value to specified pin.
// valid usr pin values are usr0, usr1, usr2 and usr3
func (b *Adaptor) DigitalWrite(pin string, val byte) (err error) {
	if strings.Contains(pin, "usr") {
		fi, err := sysfs.OpenFile(b.usrLed+pin+"/brightness", os.O_WRONLY|os.O_APPEND, 0666)
		defer fi.Close()
		if err != nil {
			return err
		}
		_, err = fi.WriteString(strconv.Itoa(int(val)))
		return err
	}
	sysfsPin, err := b.digitalPin(pin, sysfs.OUT)
	if err != nil {
		return err
	}
	return sysfsPin.Write(int(val))
}

// AnalogRead returns an analog value from specified pin
func (b *Adaptor) AnalogRead(pin string) (val int, err error) {
	analogPin, err := b.translateAnalogPin(pin)
	if err != nil {
		return
	}
	fi, err := sysfs.OpenFile(fmt.Sprintf("%v/%v", b.analogPath, analogPin), os.O_RDONLY, 0644)
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
	if (bus != 0) && (bus != 2) {
		return nil, fmt.Errorf("Bus number %d out of range", bus)
	}
	if b.i2cBuses[bus] == nil {
		b.i2cBuses[bus], err = sysfs.NewI2cDevice(fmt.Sprintf("/dev/i2c-%d", bus))
	}
	return i2c.NewConnection(b.i2cBuses[bus], address), err
}

// GetDefaultBus returns the default i2c bus for this platform
func (b *Adaptor) GetDefaultBus() int {
	return 2
}

// translatePin converts digital pin name to pin position
func (b *Adaptor) translatePin(pin string) (value int, err error) {
	for key, value := range pins {
		if key == pin {
			return value, nil
		}
	}
	err = errors.New("Not a valid pin")
	return
}

// translatePwmPin converts pwm pin name to pin position
func (b *Adaptor) translatePwmPin(pin string) (value string, err error) {
	for key, value := range pwmPins {
		if key == pin {
			return value, nil
		}
	}
	err = errors.New("Not a valid pin")
	return
}

// translateAnalogPin converts analog pin name to pin position
func (b *Adaptor) translateAnalogPin(pin string) (value string, err error) {
	for key, value := range b.analogPinMap {
		if key == pin {
			return value, nil
		}
	}
	err = errors.New("Not a valid pin")
	return
}

// digitalPin retrieves digital pin value by name
func (b *Adaptor) digitalPin(pin string, dir string) (sysfsPin sysfs.DigitalPin, err error) {
	i, err := b.translatePin(pin)
	if err != nil {
		return
	}
	if b.digitalPins[i] == nil {
		b.digitalPins[i] = sysfs.NewDigitalPin(i)
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

// pwPin retrieves pwm pin value by name
func (b *Adaptor) pwmPin(pin string) (i string, err error) {
	i, err = b.translatePwmPin(pin)
	if err != nil {
		return
	}
	if b.pwmPins[i] == nil {
		err = ensureSlot(b.slots, fmt.Sprintf("bone_pwm_%v", pin))
		if err != nil {
			return
		}
		b.pwmPins[i], err = newPwmPin(i, b.ocp)
		if err != nil {
			return
		}
	}
	return
}

// pwmWrite writes pwm value to specified pin
func (b *Adaptor) pwmWrite(pin string, val byte) (err error) {
	i, err := b.pwmPin(pin)
	if err != nil {
		return
	}
	period := 500000.0
	duty := gobot.FromScale(float64(val), 0, 255.0)
	return b.pwmPins[i].pwmWrite(strconv.Itoa(int(period)), strconv.Itoa(int(period*duty)))
}

func ensureSlot(slots, item string) (err error) {
	fi, err := sysfs.OpenFile(slots, os.O_RDWR|os.O_APPEND, 0666)
	defer fi.Close()
	if err != nil {
		return
	}

	// ensure the slot is not already written into the capemanager
	// (from: https://github.com/mrmorphic/hwio/blob/master/module_bb_pwm.go#L190)
	scanner := bufio.NewScanner(fi)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Index(line, item) > 0 {
			return
		}
	}

	_, err = fi.WriteString(item)
	if err != nil {
		return err
	}
	fi.Sync()

	scanner = bufio.NewScanner(fi)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Index(line, item) > 0 {
			return
		}
	}
	return
}

func getKernel() string {
	result, _ := exec.Command("uname", "-r").Output()

	return strings.TrimSpace(string(result))
}
