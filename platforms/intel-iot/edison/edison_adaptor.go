package edison

import (
	"errors"
	"io"
	"os"
	"strconv"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/sysfs"
)

var _ gobot.AdaptorInterface = (*EdisonAdaptor)(nil)

func writeFile(path string, data []byte) (i int, err error) {
	file, err := sysfs.OpenFile(path, os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return
	}

	return file.Write(data)
}

func readFile(path string) ([]byte, error) {
	file, err := sysfs.OpenFile(path, os.O_RDONLY, 0644)
	defer file.Close()
	if err != nil {
		return make([]byte, 0), err
	}

	buf := make([]byte, 200)
	var i = 0
	i, err = file.Read(buf)
	if i == 0 {
		return buf, err
	}
	return buf[:i], err
}

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
	tristate    sysfs.DigitalPin
	digitalPins map[int]sysfs.DigitalPin
	pwmPins     map[int]*pwmPin
	i2cDevice   io.ReadWriteCloser
	connect     func(e *EdisonAdaptor) (err error)
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
			mux{263, sysfs.HIGH},
			mux{240, sysfs.LOW},
		},
	},
	"11": sysfsPin{
		pin:          43,
		resistor:     227,
		levelShifter: 259,
		pwmPin:       5,
		mux: []mux{
			mux{262, sysfs.HIGH},
			mux{241, sysfs.LOW},
		},
	},
	"12": sysfsPin{
		pin:          42,
		resistor:     228,
		levelShifter: 260,
		pwmPin:       -1,
		mux: []mux{
			mux{242, sysfs.LOW},
		},
	},
	"13": sysfsPin{
		pin:          40,
		resistor:     229,
		levelShifter: 261,
		pwmPin:       -1,
		mux: []mux{
			mux{243, sysfs.LOW},
		},
	},
}

// changePinMode writes pin mode to current_pinmux file
func changePinMode(pin, mode string) (err error) {
	_, err = writeFile(
		"/sys/kernel/debug/gpio_debug/gpio"+pin+"/current_pinmux",
		[]byte("mode"+mode),
	)
	return
}

// NewEditionAdaptor creates a EdisonAdaptor with specified name and
// creates connect function
func NewEdisonAdaptor(name string) *EdisonAdaptor {
	return &EdisonAdaptor{
		Adaptor: *gobot.NewAdaptor(
			name,
			"EdisonAdaptor",
		),
		connect: func(e *EdisonAdaptor) (err error) {
			e.tristate = sysfs.NewDigitalPin(214)
			if err = e.tristate.Export(); err != nil {
				return err
			}
			if err = e.tristate.Direction(sysfs.OUT); err != nil {
				return err
			}
			if err = e.tristate.Write(sysfs.LOW); err != nil {
				return err
			}

			for _, i := range []int{263, 262} {
				io := sysfs.NewDigitalPin(i)
				if err = io.Export(); err != nil {
					return err
				}
				if err = io.Direction(sysfs.OUT); err != nil {
					return err
				}
				if err = io.Write(sysfs.HIGH); err != nil {
					return err
				}
				if err = io.Unexport(); err != nil {
					return err
				}
			}

			for _, i := range []int{240, 241, 242, 243} {
				io := sysfs.NewDigitalPin(i)
				if err = io.Export(); err != nil {
					return err
				}
				if err = io.Direction(sysfs.OUT); err != nil {
					return err
				}
				if err = io.Write(sysfs.LOW); err != nil {
					return err
				}
				if err = io.Unexport(); err != nil {
					return err
				}

			}

			for _, i := range []string{"111", "115", "114", "109"} {
				if err = changePinMode(i, "1"); err != nil {
					return err
				}
			}

			for _, i := range []string{"131", "129", "40"} {
				if err = changePinMode(i, "0"); err != nil {
					return err
				}
			}

			err = e.tristate.Write(sysfs.HIGH)
			return
		},
	}
}

// Connect starts conection with board and creates
// digitalPins and pwmPins adaptor maps
func (e *EdisonAdaptor) Connect() error {
	e.digitalPins = make(map[int]sysfs.DigitalPin)
	e.pwmPins = make(map[int]*pwmPin)
	return e.connect(e)
}

// Finalize closes connection to board and pins
func (e *EdisonAdaptor) Finalize() error {
	e.tristate.Unexport()
	for _, pin := range e.digitalPins {
		if pin != nil {
			pin.Unexport()
		}
	}
	for _, pin := range e.pwmPins {
		if pin != nil {
			pin.enable("0")
			pin.unexport()
		}
	}
	if e.i2cDevice != nil {
		e.i2cDevice.Close()
	}
	return nil
}

// digitalPin returns matched digitalPin for specified values
func (e *EdisonAdaptor) digitalPin(pin string, dir string) (sysfsPin sysfs.DigitalPin, err error) {
	i := sysfsPinMap[pin]
	if e.digitalPins[i.pin] == nil {
		e.digitalPins[i.pin] = sysfs.NewDigitalPin(i.pin)
		if err = e.digitalPins[i.pin].Export(); err != nil {
			return
		}

		e.digitalPins[i.resistor] = sysfs.NewDigitalPin(i.resistor)
		if err = e.digitalPins[i.resistor].Export(); err != nil {
			return
		}

		e.digitalPins[i.levelShifter] = sysfs.NewDigitalPin(i.levelShifter)
		if err = e.digitalPins[i.levelShifter].Export(); err != nil {
			return
		}

		if len(i.mux) > 0 {
			for _, mux := range i.mux {
				e.digitalPins[mux.pin] = sysfs.NewDigitalPin(mux.pin)
				if err = e.digitalPins[mux.pin].Export(); err != nil {
					return
				}

				if err = e.digitalPins[mux.pin].Direction(sysfs.OUT); err != nil {
					return
				}

				if err = e.digitalPins[mux.pin].Write(mux.value); err != nil {
					return
				}

			}
		}
	}

	if dir == "in" {
		if err = e.digitalPins[i.pin].Direction(sysfs.IN); err != nil {
			return
		}

		if err = e.digitalPins[i.resistor].Direction(sysfs.OUT); err != nil {
			return
		}

		if err = e.digitalPins[i.resistor].Write(sysfs.LOW); err != nil {
			return
		}

		if err = e.digitalPins[i.levelShifter].Direction(sysfs.OUT); err != nil {
			return
		}

		if err = e.digitalPins[i.levelShifter].Write(sysfs.LOW); err != nil {
			return
		}

	} else if dir == "out" {
		if err = e.digitalPins[i.pin].Direction(sysfs.OUT); err != nil {
			return
		}

		if err = e.digitalPins[i.resistor].Direction(sysfs.IN); err != nil {
			return
		}

		if err = e.digitalPins[i.levelShifter].Direction(sysfs.OUT); err != nil {
			return
		}

		if err = e.digitalPins[i.levelShifter].Write(sysfs.HIGH); err != nil {
			return
		}

	}
	return e.digitalPins[i.pin], nil
}

// DigitalRead reads digital value from pin
func (e *EdisonAdaptor) DigitalRead(pin string) (i int, err error) {
	sysfsPin, err := e.digitalPin(pin, "in")
	if err != nil {
		return
	}
	return sysfsPin.Read()
}

// DigitalWrite writes digital value to specified pin
func (e *EdisonAdaptor) DigitalWrite(pin string, val byte) (err error) {
	sysfsPin, err := e.digitalPin(pin, "out")
	if err != nil {
		return
	}
	return sysfsPin.Write(int(val))
}

// PwmWrite writes scaled pwm value to specified pin
func (e *EdisonAdaptor) PwmWrite(pin string, val byte) (err error) {
	sysPin := sysfsPinMap[pin]
	if sysPin.pwmPin != -1 {
		if e.pwmPins[sysPin.pwmPin] == nil {
			if err = e.DigitalWrite(pin, 1); err != nil {
				return
			}
			if err = changePinMode(strconv.Itoa(int(sysPin.pin)), "1"); err != nil {
				return
			}
			e.pwmPins[sysPin.pwmPin] = newPwmPin(sysPin.pwmPin)
			if err = e.pwmPins[sysPin.pwmPin].export(); err != nil {
				return
			}
			if err = e.pwmPins[sysPin.pwmPin].enable("1"); err != nil {
				return
			}
		}
		p, err := e.pwmPins[sysPin.pwmPin].period()
		if err != nil {
			return err
		}
		period, err := strconv.Atoi(p)
		if err != nil {
			return err
		}
		duty := gobot.FromScale(float64(val), 0, 255.0)
		return e.pwmPins[sysPin.pwmPin].writeDuty(strconv.Itoa(int(float64(period) * duty)))
	}
	return errors.New("Not a PWM pin")
}

// AnalogWrite Not Implemented
func (e *EdisonAdaptor) AnalogWrite(string, byte) (err error) {
	return errors.New("AnalogWrite is not yet implemented")
}

// InitServo Not Implemented
func (e *EdisonAdaptor) InitServo() (err error) {
	return errors.New("InitServo is not yet implemented")
}

// ServoWrite Not Implemented
func (e *EdisonAdaptor) ServoWrite(string, byte) (err error) {
	return errors.New("ServoWrite is not yet implemented")
}

// AnalogRead returns value from analog reading of specified pin
func (e *EdisonAdaptor) AnalogRead(pin string) (val int, err error) {
	buf, err := readFile(
		"/sys/bus/iio/devices/iio:device1/in_voltage" + pin + "_raw",
	)
	if err != nil {
		return
	}
	return strconv.Atoi(string(buf[0 : len(buf)-1]))
}

// I2cStart initializes i2c device for addresss
func (e *EdisonAdaptor) I2cStart(address byte) {
	e.tristate.Write(sysfs.LOW)

	for _, i := range []int{14, 165, 212, 213} {
		io := sysfs.NewDigitalPin(i)
		io.Export()
		io.Direction(sysfs.IN)
		io.Unexport()
	}

	for _, i := range []int{236, 237, 204, 205} {
		io := sysfs.NewDigitalPin(i)
		io.Export()
		io.Direction(sysfs.OUT)
		io.Write(sysfs.LOW)
		io.Unexport()
	}

	for _, i := range []string{"28", "27"} {
		changePinMode(i, "1")
	}

	e.tristate.Write(sysfs.HIGH)

	e.i2cDevice, _ = sysfs.NewI2cDevice("/dev/i2c-6", address)
}

// I2cWrite writes data to i2cDevice
func (e *EdisonAdaptor) I2cWrite(data []byte) {
	e.i2cDevice.Write(data)
}

// I2cRead reads data from i2cDevice
func (e *EdisonAdaptor) I2cRead(size uint) []byte {
	b := make([]byte, size)
	e.i2cDevice.Read(b)
	return b
}
