package beaglebone

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/sysfs"
)

const (
	Slots  = "/sys/devices/bone_capemgr.*"
	Ocp    = "/sys/devices/ocp.*"
	UsrLed = "/sys/devices/ocp.3/gpio-leds.8/leds/beaglebone:green:"
)

var i2cLocation = "/dev/i2c-1"

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
	"P8_36": "AIN5",
	"P8_35": "AIN6",
}

type BeagleboneAdaptor struct {
	gobot.Adaptor
	digitalPins []sysfs.DigitalPin
	pwmPins     map[string]*pwmPin
	analogPins  map[string]*analogPin
	i2cDevice   io.ReadWriteCloser
	connect     func()
}

// NewBeagleboneAdaptor returns a new beaglebone adaptor with specified name
func NewBeagleboneAdaptor(name string) *BeagleboneAdaptor {
	return &BeagleboneAdaptor{
		Adaptor: *gobot.NewAdaptor(
			name,
			"BeagleboneAdaptor",
		),
		connect: func() {
			ensureSlot("cape-bone-iio")
			ensureSlot("am33xx_pwm")
		},
	}
}

// Connect returns true on a succesful connection to beaglebone board.
// It initializes digital, pwm and analog pins
func (b *BeagleboneAdaptor) Connect() bool {
	b.digitalPins = make([]sysfs.DigitalPin, 120)
	b.pwmPins = make(map[string]*pwmPin)
	b.analogPins = make(map[string]*analogPin)
	b.connect()
	return true
}

// Finalize returns true when board connection is finalized correctly.
func (b *BeagleboneAdaptor) Finalize() bool {
	for _, pin := range b.pwmPins {
		if pin != nil {
			pin.release()
		}
	}
	for _, pin := range b.digitalPins {
		if pin != nil {
			pin.Unexport()
		}
	}
	if b.i2cDevice != nil {
		b.i2cDevice.Close()
	}
	return true
}

// PwmWrite writes value in specified pin
func (b *BeagleboneAdaptor) PwmWrite(pin string, val byte) {
	b.pwmWrite(pin, val)
}

// InitServo starts servo (not yet implemented)
func (b *BeagleboneAdaptor) InitServo() {}

// ServoWrite writes scaled value to servo in specified pin
func (b *BeagleboneAdaptor) ServoWrite(pin string, val byte) {
	i := b.pwmPin(pin)
	period := 20000000.0
	duty := gobot.FromScale(float64(^val), 0, 180.0)
	b.pwmPins[i].pwmWrite(strconv.Itoa(int(period)), strconv.Itoa(int(period*duty)))
}

// DigitalRead returns a digital value from specified pin
func (b *BeagleboneAdaptor) DigitalRead(pin string) (i int) {
	i, _ = b.digitalPin(pin, sysfs.IN).Read()
	return
}

// DigitalWrite writes a digital value to specified pin.
// valid usr pin values are usr0, usr1, usr2 and usr3
func (b *BeagleboneAdaptor) DigitalWrite(pin string, val byte) {
	if strings.Contains(pin, "usr") {
		fi, err := os.OpenFile(UsrLed+pin+"/brightness", os.O_WRONLY|os.O_APPEND, 0666)
		defer fi.Close()
		if err != nil {
			log.Fatal(err)
		}
		fi.WriteString(strconv.Itoa(int(val)))
	} else {
		b.digitalPin(pin, sysfs.OUT).Write(int(val))
	}
}

// AnalogRead returns an analog value from specified pin
func (b *BeagleboneAdaptor) AnalogRead(pin string) int {
	i := b.analogPin(pin)
	return b.analogPins[i].analogRead()
}

// AnalogWrite writes an analog value to specified pin
func (b *BeagleboneAdaptor) AnalogWrite(pin string, val byte) {
	b.pwmWrite(pin, val)
}

// I2cStart starts a i2c device in specified address
func (b *BeagleboneAdaptor) I2cStart(address byte) {
	b.i2cDevice, _ = sysfs.NewI2cDevice(i2cLocation, address)
}

// I2CWrite writes data to i2c device
func (b *BeagleboneAdaptor) I2cWrite(data []byte) {
	b.i2cDevice.Write(data)
}

// I2cRead returns value from i2c device using specified size
func (b *BeagleboneAdaptor) I2cRead(size uint) []byte {
	buf := make([]byte, size)
	b.i2cDevice.Read(buf)
	return buf
}

// translatePin converts digital pin name to pin position
func (b *BeagleboneAdaptor) translatePin(pin string) int {
	for key, value := range pins {
		if key == pin {
			return value
		}
	}
	panic("Not a valid pin")
}

// translatePwmPin converts pwm pin name to pin position
func (b *BeagleboneAdaptor) translatePwmPin(pin string) string {
	for key, value := range pwmPins {
		if key == pin {
			return value
		}
	}
	panic("Not a valid pin")
}

// translateAnalogPin converts analog pin name to pin position
func (b *BeagleboneAdaptor) translateAnalogPin(pin string) string {
	for key, value := range analogPins {
		if key == pin {
			return value
		}
	}
	panic("Not a valid pin")
}

// analogPin retrieves analog pin value by name
func (b *BeagleboneAdaptor) analogPin(pin string) string {
	i := b.translateAnalogPin(pin)
	if b.analogPins[i] == nil {
		b.analogPins[i] = newAnalogPin(i)
	}
	return i
}

// digitalPin retrieves digital pin value by name
func (b *BeagleboneAdaptor) digitalPin(pin string, dir string) sysfs.DigitalPin {
	i := b.translatePin(pin)
	if b.digitalPins[i] == nil {
		b.digitalPins[i] = sysfs.NewDigitalPin(i)
		err := b.digitalPins[i].Export()
		if err != nil {
			fmt.Println(err)
		}
	}
	b.digitalPins[i].Direction(dir)
	return b.digitalPins[i]
}

// pwPin retrieves pwm pin value by name
func (b *BeagleboneAdaptor) pwmPin(pin string) string {
	i := b.translatePwmPin(pin)
	if b.pwmPins[i] == nil {
		b.pwmPins[i] = newPwmPin(i)
	}
	return i
}

// pwmWrite writes pwm value to specified pin
func (b *BeagleboneAdaptor) pwmWrite(pin string, val byte) {
	i := b.pwmPin(pin)
	period := 500000.0
	duty := gobot.FromScale(float64(^val), 0, 255.0)
	b.pwmPins[i].pwmWrite(strconv.Itoa(int(period)), strconv.Itoa(int(period*duty)))
}

func ensureSlot(item string) {
	var err error
	var fi *os.File

	slot, err := filepath.Glob(Slots)
	if err != nil {
		panic(err)
	}
	fi, err = os.OpenFile(fmt.Sprintf("%v/slots", slot[0]), os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	// ensure the slot is not already written into the capemanager
	// (from: https://github.com/mrmorphic/hwio/blob/master/module_bb_pwm.go#L190)
	scanner := bufio.NewScanner(fi)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Index(line, item) > 0 {
			return
		}
	}

	fi.WriteString(item)
	fi.Sync()

	scanner = bufio.NewScanner(fi)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Index(line, item) > 0 {
			return
		}
	}
}
