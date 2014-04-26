package gobotBeaglebone

import (
	"github.com/hybridgroup/gobot"
	"strconv"
)

const SLOTS = "/sys/devices/bone_capemgr.*"
const OCP = "/sys/devices/ocp.*"
const I2C_LOCATION = "/dev/i2c-1"

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

type Beaglebone struct {
	gobot.Adaptor
	digitalPins []*digitalPin
	pwmPins     map[string]*pwmPin
	analogPins  map[string]*analogPin
	i2cDevice   *i2cDevice
}

func (b *Beaglebone) Connect() bool {
	b.digitalPins = make([]*digitalPin, 120)
	b.pwmPins = make(map[string]*pwmPin)
	b.analogPins = make(map[string]*analogPin)
	return true
}

func (b *Beaglebone) Finalize() bool {
	for _, pin := range b.pwmPins {
		if pin != nil {
			pin.release()
		}
	}
	for _, pin := range b.digitalPins {
		if pin != nil {
			pin.close()
		}
	}
	if b.i2cDevice != nil {
		b.i2cDevice.i2cDevice.Close()
	}
	return true
}
func (b *Beaglebone) Reconnect() bool  { return true }
func (b *Beaglebone) Disconnect() bool { return true }

func (b *Beaglebone) PwmWrite(pin string, val byte) {
	i := b.pwmPin(pin)
	period := 500000.0
	duty := gobot.FromScale(float64(^val), 0, 255.0)
	b.pwmPins[i].pwmWrite(strconv.Itoa(int(period)), strconv.Itoa(int(period*duty)))
}

func (b *Beaglebone) InitServo() {}
func (b *Beaglebone) ServoWrite(pin string, val byte) {
	i := b.pwmPin(pin)
	period := 20000000.0
	duty := gobot.FromScale(float64(^val), 0, 180.0)
	b.pwmPins[i].pwmWrite(strconv.Itoa(int(period)), strconv.Itoa(int(period*duty)))
}

func (b *Beaglebone) DigitalWrite(pin string, val byte) {
	i := b.digitalPin(pin, "w")
	b.digitalPins[i].digitalWrite(strconv.Itoa(int(val)))
}

func (b *Beaglebone) AnalogRead(pin string) int {
	i := b.analogPin(pin)
	return b.analogPins[i].analogRead()
}

func (b *Beaglebone) I2cStart(address byte) {
	b.i2cDevice = newI2cDevice(I2C_LOCATION, address)
	b.i2cDevice.start()
}

func (b *Beaglebone) I2cWrite(data []byte) {
	b.i2cDevice.write(data)
}

func (b *Beaglebone) I2cRead(size byte) []byte {
	return b.i2cDevice.read(size)
}

func (b *Beaglebone) translatePin(pin string) int {
	for key, value := range pins {
		if key == pin {
			return value
		}
	}
	panic("Not a valid pin")
}

func (b *Beaglebone) translatePwmPin(pin string) string {
	for key, value := range pwmPins {
		if key == pin {
			return value
		}
	}
	panic("Not a valid pin")
}

func (b *Beaglebone) translateAnalogPin(pin string) string {
	for key, value := range analogPins {
		if key == pin {
			return value
		}
	}
	panic("Not a valid pin")
}

func (b *Beaglebone) analogPin(pin string) string {
	i := b.translateAnalogPin(pin)
	if b.analogPins[i] == nil {
		b.analogPins[i] = newAnalogPin(i)
	}
	return i
}

func (b *Beaglebone) digitalPin(pin string, mode string) int {
	i := b.translatePin(pin)
	if b.digitalPins[i] == nil || b.digitalPins[i].Mode != mode {
		b.digitalPins[i] = newDigitalPin(i, mode)
	}
	return i
}

func (b *Beaglebone) pwmPin(pin string) string {
	i := b.translatePwmPin(pin)
	if b.pwmPins[i] == nil {
		b.pwmPins[i] = newPwmPin(i)
	}
	return i
}
