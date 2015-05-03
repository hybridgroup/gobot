package c1

import (
	"errors"
	"io"
	"log"
	//"io/ioutil"
	//"strconv"
	//"strings"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/i2c"
	"github.com/hybridgroup/gobot/sysfs"
)

var _ gobot.Adaptor = (*ODroidC1Adaptor)(nil)

var _ gpio.DigitalReader = (*ODroidC1Adaptor)(nil)
var _ gpio.DigitalWriter = (*ODroidC1Adaptor)(nil)

var _ i2c.I2c = (*ODroidC1Adaptor)(nil)

type ODroidC1Adaptor struct {
	name        string
	//revision    string
	i2cLocation string
	digitalPins map[int]sysfs.DigitalPin
	i2cDevice   io.ReadWriteCloser
}

var pins = map[string]int{
	"3": 74,
	"5": 75,
	"7": 83,
	"8": 113,
	"10": 114,
	"11": 88,
	"12": 87,
	"13": 116,
	"15": 115,
	"16": 104,
	"18": 102,
	"19": 107,
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
	"33": 108,
	"35": 97,
	"36": 98,
}

// NewODroidC1Adaptor creates an ODroidC1Adaptor with specified name
func NewODroidC1Adaptor(name string) *ODroidC1Adaptor {
	r := &ODroidC1Adaptor{
		name:        name,
		digitalPins: make(map[int]sysfs.DigitalPin),
		i2cLocation: "/sys/bus/i2c",
	}

	return r
}

func (r *ODroidC1Adaptor) Name() string { return r.name }

// Connect starts conection with board and creates
// digitalPins and pwmPins adaptor maps
func (r *ODroidC1Adaptor) Connect() (errs []error) {
	return
}

// Finalize closes connection to board and pins
func (r *ODroidC1Adaptor) Finalize() (errs []error) {
	for _, pin := range r.digitalPins {
		if pin != nil {
			if err := pin.Unexport(); err != nil {
				errs = append(errs, err)
			}
		}
	}
	if r.i2cDevice != nil {
		if err := r.i2cDevice.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// digitalPin returns matched digitalPin for specified values
func (r *ODroidC1Adaptor) digitalPin(pin string, dir string) (sysfsPin sysfs.DigitalPin, err error) {
	log.Println("Looking for digitalPin " + pin + " in dir " + dir)

	var i int

	if val, ok := pins[pin]; ok {
		i = val
	} else {
		err = errors.New("Not a valid pin")
		return
	}

	if r.digitalPins[i] == nil {
		r.digitalPins[i] = sysfs.NewDigitalPin(i)
		if err = r.digitalPins[i].Export(); err != nil {
			log.Println("Err with Export")
			log.Fatal(err)
			return
		}
	}

	if err = r.digitalPins[i].Direction(dir); err != nil {
			log.Println("Err with Direction")
			log.Fatal(err)
		return
	}

	//log.Println("Returning pin: " + r.digitalPins[i].pin + ", label: " + r.digitalPins[i].label)
	return r.digitalPins[i], nil
}

// DigitalRead reads digital value from pin
func (r *ODroidC1Adaptor) DigitalRead(pin string) (val int, err error) {
	sysfsPin, err := r.digitalPin(pin, sysfs.IN)
	if err != nil {
		return
	}
	return sysfsPin.Read()
}

// DigitalWrite writes digital value to specified pin
func (r *ODroidC1Adaptor) DigitalWrite(pin string, val byte) (err error) {
	log.Println("DigitalWrite " + pin )

	sysfsPin, err := r.digitalPin(pin, sysfs.OUT)
	if err != nil {
		return err
	}
	return sysfsPin.Write(int(val))
}

// I2cStart starts a i2c device in specified address
func (r *ODroidC1Adaptor) I2cStart(address byte) (err error) {
	r.i2cDevice, err = sysfs.NewI2cDevice(r.i2cLocation, address)
	return err
}

// I2CWrite writes data to i2c device
func (r *ODroidC1Adaptor) I2cWrite(data []byte) (err error) {
	_, err = r.i2cDevice.Write(data)
	return
}

// I2cRead returns value from i2c device using specified size
func (r *ODroidC1Adaptor) I2cRead(size uint) (data []byte, err error) {
	data = make([]byte, size)
	_, err = r.i2cDevice.Read(data)
	return
}
