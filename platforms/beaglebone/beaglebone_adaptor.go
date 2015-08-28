package beaglebone

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/i2c"
	"github.com/hybridgroup/gobot/sysfs"
)

var _ gobot.Adaptor = (*BeagleboneAdaptor)(nil)

var _ gpio.DigitalReader = (*BeagleboneAdaptor)(nil)
var _ gpio.DigitalWriter = (*BeagleboneAdaptor)(nil)
var _ gpio.AnalogReader = (*BeagleboneAdaptor)(nil)
var _ gpio.PwmWriter = (*BeagleboneAdaptor)(nil)
var _ gpio.ServoWriter = (*BeagleboneAdaptor)(nil)

var _ i2c.I2c = (*BeagleboneAdaptor)(nil)

var slots = "/sys/devices/bone_capemgr.*"
var ocp = "/sys/devices/ocp.*"
var usrLed = "/sys/devices/ocp.3/gpio-leds.8/leds/beaglebone:green:"

var glob = func(pattern string) (matches []string, err error) {
	return filepath.Glob(pattern)
}

var pins = map[string]int{
	"P8_3":  38,
	"P8_4":  39,
	"P8_5":  34,
	"P8_6":  35,
	"P8_7":  66,
	"P8_8":  67,
	"P8_9":  69,
	"P8_10": 68,
	"P8_11": 45,
	"P8_12": 44,
	"P8_13": 23,
	"P8_14": 26,
	"P8_15": 47,
	"P8_16": 46,
	"P8_17": 27,
	"P8_18": 65,
	"P8_19": 22,
	"P8_20": 63,
	"P8_21": 62,
	"P8_22": 37,
	"P8_23": 36,
	"P8_24": 33,
	"P8_25": 32,
	"P8_26": 61,
	"P8_27": 86,
	"P8_28": 88,
	"P8_29": 87,
	"P8_30": 89,
	"P8_31": 10,
	"P8_32": 11,
	"P8_33": 9,
	"P8_34": 81,
	"P8_35": 8,
	"P8_36": 80,
	"P8_37": 78,
	"P8_38": 79,
	"P8_39": 76,
	"P8_40": 77,
	"P8_41": 74,
	"P8_42": 75,
	"P8_43": 72,
	"P8_44": 73,
	"P8_45": 70,
	"P8_46": 71,
	"P9_11": 30,
	"P9_12": 60,
	"P9_13": 31,
	"P9_14": 50,
	"P9_15": 48,
	"P9_16": 51,
	"P9_17": 5,
	"P9_18": 4,
	"P9_19": 13,
	"P9_20": 12,
	"P9_21": 3,
	"P9_22": 2,
	"P9_23": 49,
	"P9_24": 15,
	"P9_25": 117,
	"P9_26": 14,
	"P9_27": 115,
	"P9_28": 113,
	"P9_29": 111,
	"P9_30": 112,
	"P9_31": 110,
}

var pwmPins = map[string]string{
	"P9_14": "P9_14",
	"P9_21": "P9_21",
	"P9_22": "P9_22",
	"P9_29": "P9_29",
	"P9_42": "P9_42",
	"P8_13": "P8_13",
	"P8_34": "P8_34",
	"P8_45": "P8_45",
	"P8_46": "P8_46",
}

var analogPins = map[string]string{
	"P9_39": "AIN0",
	"P9_40": "AIN1",
	"P9_37": "AIN2",
	"P9_38": "AIN3",
	"P9_33": "AIN4",
	"P9_36": "AIN5",
	"P9_35": "AIN6",
}

// BeagleboneAdaptor is the gobot.Adaptor representation for the Beaglebone
type BeagleboneAdaptor struct {
	name        string
	digitalPins []sysfs.DigitalPin
	pwmPins     map[string]*pwmPin
	i2cDevice   sysfs.I2cDevice
	ocp         string
	helper      string
	slots       string
}

// NewBeagleboneAdaptor returns a new BeagleboneAdaptor with specified name
func NewBeagleboneAdaptor(name string) *BeagleboneAdaptor {
	b := &BeagleboneAdaptor{
		name:        name,
		digitalPins: make([]sysfs.DigitalPin, 120),
		pwmPins:     make(map[string]*pwmPin),
	}

	g, _ := glob(ocp)
	b.ocp = g[0]
	g, _ = glob(slots)
	b.slots = fmt.Sprintf("%v/slots", g[0])

	return b
}

// Name returns the BeagleboneAdaptors name
func (b *BeagleboneAdaptor) Name() string { return b.name }

// Connect initializes the pwm and analog dts.
func (b *BeagleboneAdaptor) Connect() (errs []error) {
	if err := ensureSlot(b.slots, "cape-bone-iio"); err != nil {
		return []error{err}
	}

	if err := ensureSlot(b.slots, "am33xx_pwm"); err != nil {
		return []error{err}
	}

	g, err := glob(fmt.Sprintf("%v/helper.*", b.ocp))
	if err != nil {
		return []error{err}
	}
	b.helper = g[0]

	return
}

// Finalize releases all i2c devices and exported analog, digital, pwm pins.
func (b *BeagleboneAdaptor) Finalize() (errs []error) {
	for _, pin := range b.pwmPins {
		if pin != nil {
			if err := pin.release(); err != nil {
				errs = append(errs, err)
			}
		}
	}
	for _, pin := range b.digitalPins {
		if pin != nil {
			if err := pin.Unexport(); err != nil {
				errs = append(errs, err)
			}
		}
	}
	if b.i2cDevice != nil {
		if err := b.i2cDevice.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return
}

// PwmWrite writes the 0-254 value to the specified pin
func (b *BeagleboneAdaptor) PwmWrite(pin string, val byte) (err error) {
	return b.pwmWrite(pin, val)
}

// ServoWrite writes the 0-180 degree val to the specified pin.
func (b *BeagleboneAdaptor) ServoWrite(pin string, val byte) (err error) {
	i, err := b.pwmPin(pin)
	if err != nil {
		return err
	}
	period := 16666666.0
	duty := (gobot.FromScale(float64(val), 0, 180.0) * 0.115) + 0.05
	return b.pwmPins[i].pwmWrite(strconv.Itoa(int(period)), strconv.Itoa(int(period*duty)))
}

// DigitalRead returns a digital value from specified pin
func (b *BeagleboneAdaptor) DigitalRead(pin string) (val int, err error) {
	sysfsPin, err := b.digitalPin(pin, sysfs.IN)
	if err != nil {
		return
	}
	return sysfsPin.Read()
}

// DigitalWrite writes a digital value to specified pin.
// valid usr pin values are usr0, usr1, usr2 and usr3
func (b *BeagleboneAdaptor) DigitalWrite(pin string, val byte) (err error) {
	if strings.Contains(pin, "usr") {
		fi, err := sysfs.OpenFile(usrLed+pin+"/brightness", os.O_WRONLY|os.O_APPEND, 0666)
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
func (b *BeagleboneAdaptor) AnalogRead(pin string) (val int, err error) {
	analogPin, err := b.translateAnalogPin(pin)
	if err != nil {
		return
	}
	fi, err := sysfs.OpenFile(fmt.Sprintf("%v/%v", b.helper, analogPin), os.O_RDONLY, 0644)
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

// I2cStart starts a i2c device in specified address on i2c bus /dev/i2c-1
func (b *BeagleboneAdaptor) I2cStart(address int) (err error) {
	if b.i2cDevice == nil {
		b.i2cDevice, err = sysfs.NewI2cDevice("/dev/i2c-1", address)
	}
	return
}

// I2cWrite writes data to i2c device
func (b *BeagleboneAdaptor) I2cWrite(address int, data []byte) (err error) {
	if err = b.i2cDevice.SetAddress(address); err != nil {
		return
	}
	_, err = b.i2cDevice.Write(data)
	return
}

// I2cRead returns size bytes from the i2c device
func (b *BeagleboneAdaptor) I2cRead(address int, size int) (data []byte, err error) {
	if err = b.i2cDevice.SetAddress(address); err != nil {
		return
	}
	data = make([]byte, size)
	_, err = b.i2cDevice.Read(data)
	return
}

// translatePin converts digital pin name to pin position
func (b *BeagleboneAdaptor) translatePin(pin string) (value int, err error) {
	for key, value := range pins {
		if key == pin {
			return value, nil
		}
	}
	err = errors.New("Not a valid pin")
	return
}

// translatePwmPin converts pwm pin name to pin position
func (b *BeagleboneAdaptor) translatePwmPin(pin string) (value string, err error) {
	for key, value := range pwmPins {
		if key == pin {
			return value, nil
		}
	}
	err = errors.New("Not a valid pin")
	return
}

// translateAnalogPin converts analog pin name to pin position
func (b *BeagleboneAdaptor) translateAnalogPin(pin string) (value string, err error) {
	for key, value := range analogPins {
		if key == pin {
			return value, nil
		}
	}
	err = errors.New("Not a valid pin")
	return
}

// digitalPin retrieves digital pin value by name
func (b *BeagleboneAdaptor) digitalPin(pin string, dir string) (sysfsPin sysfs.DigitalPin, err error) {
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
func (b *BeagleboneAdaptor) pwmPin(pin string) (i string, err error) {
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
func (b *BeagleboneAdaptor) pwmWrite(pin string, val byte) (err error) {
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
