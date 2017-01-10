package c1

import (
	"errors"
	"fmt"
	//"log"
	"os"
	"strconv"
	"strings"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/i2c"
	"github.com/hybridgroup/gobot/sysfs"
)

const (
	pwmBase    string = "/sys/devices/platform/pwm-ctrl"
	analogBase string = "/sys/class/saradc"
	i2cBase    string = "/sys/bus/i2c"
)

type ODroidC1Adaptor struct {
	name        string
	i2cLocation string
	digitalPins map[int]sysfs.DigitalPin
	pwmPins     map[int]*pwmPin
	i2cDevice   sysfs.I2cDevice
}

var _ gobot.Adaptor = (*ODroidC1Adaptor)(nil)

var _ gpio.DigitalReader = (*ODroidC1Adaptor)(nil)
var _ gpio.DigitalWriter = (*ODroidC1Adaptor)(nil)
var _ gpio.AnalogReader = (*ODroidC1Adaptor)(nil)
var _ gpio.PwmWriter = (*ODroidC1Adaptor)(nil)
var _ gpio.ServoWriter = (*ODroidC1Adaptor)(nil)

var _ i2c.I2c = (*ODroidC1Adaptor)(nil)

var pins = map[string]int{
	"3":  74,
	"5":  75,
	"7":  83,
	"8":  113,
	"10": 114,
	"11": 88,
	"12": 87,
	"13": 116,
	"15": 115,
	"16": 104,
	"18": 102,
	"21": 106,
	"22": 103,
	"23": 105,
	"24": 117,
	"26": 118,
	"27": 76,
	"28": 77,
	"29": 101,
	"31": 100,
	"32": 99,
	"35": 97,
	"36": 98,
}

// pwm 0 (pin 33) is enabled by default.  pwm 1 (19) can be enabled, but will disable spi as it shares a port
var pwmPins = map[string]map[string]int{
	"33": map[string]int{
		"gpioNum": 108,
		"pwmNum": 0,
	},
	"19": map[string]int{
		"gpioNum": 107,
		"pwmNum": 1,
	},
}

var analogPins = map[string]string{
	"37": "saradc_ch1",
	"40": "saradc_ch0",
}

// NewODroidC1Adaptor creates an ODroidC1Adaptor with specified name
func NewODroidC1Adaptor(name string) *ODroidC1Adaptor {
	o := &ODroidC1Adaptor{
		name:        name,
		digitalPins: make(map[int]sysfs.DigitalPin),
		pwmPins:     make(map[int]*pwmPin),
		i2cLocation: i2cBase,
	}

	return o
}

func (o *ODroidC1Adaptor) Name() string { return o.name }

// Connect starts conection with board and creates
// digitalPins and pwmPins adaptor maps
func (o *ODroidC1Adaptor) Connect() (errs []error) {
	return
}

// Finalize closes connection to board and pins
func (o *ODroidC1Adaptor) Finalize() (errs []error) {
	for _, pin := range o.digitalPins {
		if pin != nil {
			if err := pin.Unexport(); err != nil {
				errs = append(errs, err)
			}
		}
	}
	for _, pin := range o.pwmPins {
		if pin != nil {
			if err := pin.release(); err != nil {
				errs = append(errs, err)
			}
		}
	}
	if o.i2cDevice != nil {
		if err := o.i2cDevice.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// digitalPin returns matched digitalPin for specified values
func (o *ODroidC1Adaptor) digitalPin(pin string, dir string) (sysfsPin sysfs.DigitalPin, err error) {
	var i int

	if val, ok := pins[pin]; ok {
		i = val
	} else {
		err = errors.New("Not a valid pin")
		return
	}

	if o.digitalPins[i] == nil {
		o.digitalPins[i] = sysfs.NewDigitalPin(i)
		if err = o.digitalPins[i].Export(); err != nil {
			return
		}
	}

	if err = o.digitalPins[i].Direction(dir); err != nil {
		return
	}

	return o.digitalPins[i], nil
}

// DigitalRead reads digital value from pin
func (o *ODroidC1Adaptor) DigitalRead(pin string) (val int, err error) {
	sysfsPin, err := o.digitalPin(pin, sysfs.IN)
	if err != nil {
		return
	}
	return sysfsPin.Read()
}

// DigitalWrite writes digital value to specified pin
func (o *ODroidC1Adaptor) DigitalWrite(pin string, val byte) (err error) {
	sysfsPin, err := o.digitalPin(pin, sysfs.OUT)
	if err != nil {
		return err
	}
	return sysfsPin.Write(int(val))
}

// I2cStart starts a i2c device in specified address
func (o *ODroidC1Adaptor) I2cStart(address int) (err error) {
	if o.i2cDevice == nil {
		o.i2cDevice, err = sysfs.NewI2cDevice(o.i2cLocation, address)
	}
	return err
}

// I2CWrite writes data to i2c device
func (o *ODroidC1Adaptor) I2cWrite(address int, data []byte) (err error) {
	if err = o.i2cDevice.SetAddress(address); err != nil {
		return
	}
	_, err = o.i2cDevice.Write(data)
	return
}

// I2cRead returns value from i2c device using specified size
func (o *ODroidC1Adaptor) I2cRead(address int, size int) (data []byte, err error) {
	if err = o.i2cDevice.SetAddress(address); err != nil {
		return
	}
	data = make([]byte, size)
	_, err = o.i2cDevice.Read(data)
	return
}

// translatePwmPin converts pwm pin name to pin position
func (o *ODroidC1Adaptor) translatePwmPin(pin string) (gpioNum int, pwmNum int, err error) {
	//log.Printf("[translatePwmPin] pin: %s", pin)
	pwm := pwmPins[pin]
	if pwm == nil {
		err = errors.New("Not a valid pwm pin")
		return
	}

	//log.Printf("   pwm: %s", pwm)
	return pwm["gpioNum"], pwm["pwmNum"], nil

}

// pwmPin retrieves pwm pin value by name
func (o *ODroidC1Adaptor) pwmPin(pin string) (gpioNum int, pwmNum int, err error) {
	//log.Printf("[pwmPin] pin: %s", pin)
	
	gpioNum, pwmNum, err = o.translatePwmPin(pin)	
	if err != nil {
		//log.Printf("   err: %s", err)
		return
	}
	//log.Printf("   translated pin: %s to gpioNum: %v, pwmNum: %v", pin, gpioNum, pwmNum)

	if o.pwmPins[gpioNum] == nil {
		//log.Printf("   creating new pwmPin for gpioNum: %v", gpioNum)
		o.pwmPins[gpioNum], err = newPwmPin(pin, gpioNum, pwmNum, pwmBase)
		if err != nil {
			//log.Printf("   err: %s", err)
			return
		}
	}
	return
}



// PwmWrite writes the 0-255 value to the specified pin
func (o *ODroidC1Adaptor) PwmWrite(pin string, val byte) (err error) {
	return o.pwmWrite(pin, val)
}

// pwmWrite writes pwm value to specified pin
func (o *ODroidC1Adaptor) pwmWrite(pin string, val byte) (err error) {
	//log.Printf("[pwmWrite] pin: %s, val: %v", pin, val)
	
	gpioNum, _, err := o.pwmPin(pin)
	//log.Printf("[pwmWrite]   gpioNum: %s", strconv.Itoa(gpioNum));
	
	if err != nil {
		//log.Printf("[pwmWrite]   err: %s", err);
		return
	}
	freq := 100000.0
	duty := int( (1023.0 / 255.0) * float64(val))
	//log.Printf("[pwmWrite]   duty: %v", duty);
	return o.pwmPins[gpioNum].pwmWrite(strconv.Itoa(int(freq)), strconv.Itoa(duty))
}

// PwmWrite writes the 0-254 value to the specified pin
//func (o *ODroidC1Adaptor) PwmWrite(pin string, val byte) (err error) {
//	return o.pwmWrite(pin, val)
//}

// ServoWrite writes the 0-180 degree val to the specified pin.
func (o *ODroidC1Adaptor) ServoWrite(pin string, val byte) (err error) {
	gpioNum, _, err := o.pwmPin(pin)
	if err != nil {
		return err
	}
	freq := 16666666.0
	duty := (gobot.FromScale(float64(val), 0, 180.0) * 0.115) + 0.05
	return o.pwmPins[gpioNum].pwmWrite(strconv.Itoa(int(freq)), strconv.Itoa(int(freq*duty)))
}

// translateAnalogPin converts analog pin name to pin position
func (o *ODroidC1Adaptor) translateAnalogPin(pin string) (value string, err error) {
	for key, value := range analogPins {
		if key == pin {
			return value, nil
		}
	}
	err = errors.New("Not a valid pin")
	return
}

// AnalogRead returns an analog value from specified pin
func (o *ODroidC1Adaptor) AnalogRead(pin string) (val int, err error) {
	analogPin, err := o.translateAnalogPin(pin)
	if err != nil {
		return
	}
	fi, err := sysfs.OpenFile(fmt.Sprintf("%v/%v", analogBase, analogPin), os.O_RDONLY, 0644)
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
