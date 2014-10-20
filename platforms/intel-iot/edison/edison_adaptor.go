package edison

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/hybridgroup/gobot"
)

type mux struct {
	pin   int
	value int
}
type sysfsPin struct {
	pin          int
	resistor     int
	levelShifter int
	pwmPin       int
	mux          []mux
}

type EdisonAdaptor struct {
	gobot.Adaptor
	tristate    *digitalPin
	digitalPins map[int]*digitalPin
	pwmPins     map[int]*pwmPin
	i2cDevice   *i2cDevice
	connect     func(e *EdisonAdaptor)
}

var sysfsPinMap = map[string]sysfsPin{
	"0": sysfsPin{
		pin:          130,
		resistor:     216,
		levelShifter: 248,
		pwmPin:       -1,
		mux:          []mux{},
	},
	"1": sysfsPin{
		pin:          131,
		resistor:     217,
		levelShifter: 249,
		pwmPin:       -1,
		mux:          []mux{},
	},
	"2": sysfsPin{
		pin:          128,
		resistor:     218,
		levelShifter: 250,
		pwmPin:       -1,
		mux:          []mux{},
	},
	"3": sysfsPin{
		pin:          12,
		resistor:     219,
		levelShifter: 251,
		pwmPin:       0,
		mux:          []mux{},
	},

	"4": sysfsPin{
		pin:          129,
		resistor:     220,
		levelShifter: 252,
		pwmPin:       -1,
		mux:          []mux{},
	},
	"5": sysfsPin{
		pin:          13,
		resistor:     221,
		levelShifter: 253,
		pwmPin:       1,
		mux:          []mux{},
	},
	"6": sysfsPin{
		pin:          182,
		resistor:     222,
		levelShifter: 254,
		pwmPin:       2,
		mux:          []mux{},
	},
	"7": sysfsPin{
		pin:          48,
		resistor:     223,
		levelShifter: 255,
		pwmPin:       -1,
		mux:          []mux{},
	},
	"8": sysfsPin{
		pin:          49,
		resistor:     224,
		levelShifter: 256,
		pwmPin:       -1,
		mux:          []mux{},
	},
	"9": sysfsPin{
		pin:          183,
		resistor:     225,
		levelShifter: 257,
		pwmPin:       3,
		mux:          []mux{},
	},
	"10": sysfsPin{
		pin:          41,
		resistor:     226,
		levelShifter: 258,
		pwmPin:       4,
		mux: []mux{
			mux{263, 1},
			mux{240, 0},
		},
	},
	"11": sysfsPin{
		pin:          43,
		resistor:     227,
		levelShifter: 259,
		pwmPin:       5,
		mux: []mux{
			mux{262, 1},
			mux{241, 0},
		},
	},
	"12": sysfsPin{
		pin:          42,
		resistor:     228,
		levelShifter: 260,
		pwmPin:       -1,
		mux: []mux{
			mux{242, 0},
		},
	},
	"13": sysfsPin{
		pin:          40,
		resistor:     229,
		levelShifter: 261,
		pwmPin:       -1,
		mux: []mux{
			mux{243, 0},
		},
	},
}

// writeFile validates file existence and writes data into it
func writeFile(name, data string) error {
	if _, err := os.Stat(name); err == nil {
		err := ioutil.WriteFile(
			name,
			[]byte(data),
			0644,
		)
		if err != nil {
			return err
		}
	} else {
		return errors.New("File doesn't exist: " + name)
	}
	return nil
}

// changePinMode writes pin mode to current_pinmux file
func changePinMode(pin, mode string) {
	err := writeFile(
		"/sys/kernel/debug/gpio_debug/gpio"+pin+"/current_pinmux",
		"mode"+mode,
	)
	if err != nil {
		panic(err)
	}
}

// NewEditionAdaptor creates a EdisonAdaptor with specified name and
// creates connect function
func NewEdisonAdaptor(name string) *EdisonAdaptor {
	return &EdisonAdaptor{
		Adaptor: *gobot.NewAdaptor(
			name,
			"EdisonAdaptor",
		),
		connect: func(e *EdisonAdaptor) {
			e.tristate = newDigitalPin(214)
			e.tristate.setDir("out")
			e.tristate.digitalWrite("0")

			for _, i := range []int{263, 262} {
				io := newDigitalPin(i)
				io.setDir("out")
				io.digitalWrite("1")
				io.close()
			}

			for _, i := range []int{240, 241, 242, 243} {
				io := newDigitalPin(i)
				io.setDir("out")
				io.digitalWrite("0")
				io.close()
			}

			for _, i := range []string{"111", "115", "114", "109"} {
				changePinMode(i, "1")
			}

			for _, i := range []string{"131", "129", "40"} {
				changePinMode(i, "0")
			}

			e.tristate.digitalWrite("1")
		},
	}
}

// Connect starts conection with board and creates
// digitalPins and pwmPins adaptor maps
func (e *EdisonAdaptor) Connect() bool {
	e.digitalPins = make(map[int]*digitalPin)
	e.pwmPins = make(map[int]*pwmPin)
	e.connect(e)
	return true
}

// Finalize closes connection to board and pins
func (e *EdisonAdaptor) Finalize() bool {
	e.tristate.close()
	for _, pin := range e.digitalPins {
		if pin != nil {
			pin.close()
		}
	}
	for _, pin := range e.pwmPins {
		if pin != nil {
			pin.close()
		}
	}
	if e.i2cDevice != nil {
		e.i2cDevice.file.Close()
	}
	return true
}

// Reconnect retries connection to edison board
func (e *EdisonAdaptor) Reconnect() bool { return true }

// Disconnect returns true if connection to edison board is finished successfully
func (e *EdisonAdaptor) Disconnect() bool { return true }

// digitalPin returns matched digitalPin for specified values
func (e *EdisonAdaptor) digitalPin(pin string, dir string) *digitalPin {
	i := sysfsPinMap[pin]
	if e.digitalPins[i.pin] == nil {
		e.digitalPins[i.pin] = newDigitalPin(i.pin)
		e.digitalPins[i.resistor] = newDigitalPin(i.resistor)
		e.digitalPins[i.levelShifter] = newDigitalPin(i.levelShifter)
		if len(i.mux) > 0 {
			for _, mux := range i.mux {
				e.digitalPins[mux.pin] = newDigitalPin(mux.pin)
				e.digitalPins[mux.pin].setDir("out")
				e.digitalPins[mux.pin].digitalWrite(strconv.Itoa(mux.value))
			}
		}
	}

	if dir == "in" && e.digitalPins[i.pin].dir != "in" {
		e.digitalPins[i.pin].setDir("in")
		e.digitalPins[i.resistor].setDir("out")
		e.digitalPins[i.resistor].digitalWrite("0")
		e.digitalPins[i.levelShifter].setDir("out")
		e.digitalPins[i.levelShifter].digitalWrite("0")
	} else if dir == "out" && e.digitalPins[i.pin].dir != "out" {
		e.digitalPins[i.pin].setDir("out")
		e.digitalPins[i.resistor].setDir("in")
		e.digitalPins[i.levelShifter].setDir("out")
		e.digitalPins[i.levelShifter].digitalWrite("1")
	}
	return e.digitalPins[i.pin]
}

// DigitalRead reads digital value from pin
func (e *EdisonAdaptor) DigitalRead(pin string) int {
	return e.digitalPin(pin, "in").digitalRead()
}

// DigitalWrite writes digital value to specified pin
func (e *EdisonAdaptor) DigitalWrite(pin string, val byte) {
	e.digitalPin(pin, "out").digitalWrite(strconv.Itoa(int(val)))
}

// PwmWrite writes scaled pwm value to specified pin
func (e *EdisonAdaptor) PwmWrite(pin string, val byte) {
	sysPin := sysfsPinMap[pin]
	if sysPin.pwmPin != -1 {
		if e.pwmPins[sysPin.pwmPin] == nil {
			e.DigitalWrite(pin, 1)
			changePinMode(strconv.Itoa(int(sysPin.pin)), "1")
			e.pwmPins[sysPin.pwmPin] = newPwmPin(sysPin.pwmPin)
		}
		period, err := strconv.Atoi(e.pwmPins[sysPin.pwmPin].period())
		if err != nil {
			panic(err)
		}
		duty := gobot.FromScale(float64(val), 0, 255.0)
		e.pwmPins[sysPin.pwmPin].writeDuty(strconv.Itoa(int(float64(period) * duty)))
	} else {
		fmt.Println("Not a PWM pin")
	}
}

// AnalogRead returns value from analog reading of specified pin
func (e *EdisonAdaptor) AnalogRead(pin string) int {
	buf, err := ioutil.ReadFile(
		"/sys/bus/iio/devices/iio:device1/in_voltage" + pin + "_raw",
	)
	if err != nil {
		panic(err)
	}
	val, err := strconv.Atoi(string(buf[0 : len(buf)-1]))
	if err != nil {
		panic(err)
	}
	return val
}

// I2cStart initializes i2c device for addresss
func (e *EdisonAdaptor) I2cStart(address byte) {
	e.tristate.digitalWrite("0")

	for _, i := range []int{14, 165, 212, 213} {
		io := newDigitalPin(i)
		io.setDir("in")
		io.close()
	}

	for _, i := range []int{236, 237, 204, 205} {
		io := newDigitalPin(i)
		io.setDir("out")
		io.digitalWrite("0")
		io.close()
	}

	for _, i := range []string{"28", "27"} {
		changePinMode(i, "1")
	}

	e.tristate.digitalWrite("1")

	e.i2cDevice = newI2cDevice(address)
	e.i2cDevice.start()
}

// I2cWrite writes data to i2cDevice
func (e *EdisonAdaptor) I2cWrite(data []byte) {
	e.i2cDevice.write(data)
}

// I2cRead reads data from i2cDevice
func (e *EdisonAdaptor) I2cRead(size uint) []byte {
	return e.i2cDevice.read(size)
}
