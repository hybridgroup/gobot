package joule

import (
	"errors"
	"os"
	"strconv"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/sysfs"
)

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

type sysfsPin struct {
	pin    int
	pwmPin int
}

// JouleAdaptor represents an Intel Joule
type JouleAdaptor struct {
	name        string
	digitalPins map[int]sysfs.DigitalPin
	pwmPins     map[int]*pwmPin
	i2cDevice   sysfs.I2cDevice
	connect     func(e *JouleAdaptor) (err error)
}

var sysfsPinMap = map[string]sysfsPin{
	// disabled
	"0": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	"1": sysfsPin{
		pin:    446,
		pwmPin: -1,
	},
	"2": sysfsPin{
		pin:    421,
		pwmPin: -1,
	},
	// disabled
	"3": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	"4": sysfsPin{
		pin:    422,
		pwmPin: -1,
	},
	"5": sysfsPin{
		pin:    356,
		pwmPin: -1,
	},
	"6": sysfsPin{
		pin:    417,
		pwmPin: -1,
	},
	// UART
	"7": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	"8": sysfsPin{
		pin:    419,
		pwmPin: -1,
	},
	// disabled
	"9": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	"10": sysfsPin{
		pin:    416,
		pwmPin: -1,
	},
	"11": sysfsPin{
		pin:    381,
		pwmPin: -1,
	},
	"13": sysfsPin{
		pin:    382,
		pwmPin: -1,
	},
	"15": sysfsPin{
		pin:    380,
		pwmPin: -1,
	},
	"17": sysfsPin{
		pin:    379,
		pwmPin: -1,
	},
	"19": sysfsPin{
		pin:    378,
		pwmPin: -1,
	},
	// UART
	"21": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	// disabled
	"22": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	// UART
	"23": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	// disabled
	"24": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	"25": sysfsPin{
		pin:    463,
		pwmPin: 0,
	},
	// low voltage should not use
	"26": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	"27": sysfsPin{
		pin:    464,
		pwmPin: 1,
	},
	// disabled
	"28": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	"29": sysfsPin{
		pin:    465,
		pwmPin: 2,
	},
	// disabled?
	"30": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	"31": sysfsPin{
		pin:    466,
		pwmPin: 3,
	},
	// disabled?
	"32": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	// disabled
	"33": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	"34": sysfsPin{
		pin:    393,
		pwmPin: -1,
	},
	// GND
	"35": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	// GND
	"36": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	// GND
	"37": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	// GND
	"38": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	// disabled
	"39": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	// GND
	"40": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},

	// Second header
	// GND
	"41": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	// 5V
	"42": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	// GND
	"43": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	// 5V
	"44": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	// GND
	"45": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	// 3.3V
	"46": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	// GND
	"47": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	// 3.3V
	"48": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	// GND
	"49": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	// 1.8V
	"50": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	// GPIO
	"51": sysfsPin{
		pin:    456,
		pwmPin: -1,
	},
	// 1.8V
	"52": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	// PANEL
	"53": sysfsPin{
		pin:    270,
		pwmPin: -1,
	},
	// GND
	"54": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	// PANEL
	"55": sysfsPin{
		pin:    271,
		pwmPin: -1,
	},
	// CAMERA
	"56": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	// PANEL
	"57": sysfsPin{
		pin:    272,
		pwmPin: -1,
	},
	// CAMERA
	"58": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	// SPP0FS0
	"59": sysfsPin{
		pin:    411,
		pwmPin: -1,
	},
	// CAMERA
	"60": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	// SPP0FS1
	"61": sysfsPin{
		pin:    412,
		pwmPin: -1,
	},
	// SPI_DAT
	"62": sysfsPin{
		pin:    -1,
		pwmPin: -1,
	},
	// SPP0FS2
	"63": sysfsPin{
		pin:    411,
		pwmPin: -1,
	},
	// SPICLKB
	"64": sysfsPin{
		pin:    384,
		pwmPin: -1,
	},
	// SPP0FS3
	"65": sysfsPin{
		pin:    410,
		pwmPin: -1,
	},
	// SPICLKA
	"66": sysfsPin{
		pin:    383,
		pwmPin: -1,
	},
	// SPP0TX
	"67": sysfsPin{
		pin:    414,
		pwmPin: -1,
	},
	// UART0RX
	"68": sysfsPin{
		pin:    467,
		pwmPin: -1,
	},
	// SPP0RX
	"69": sysfsPin{
		pin:    415,
		pwmPin: -1,
	},
	// UART0RT
	"70": sysfsPin{
		pin:    469,
		pwmPin: -1,
	},
	// I2C1SDA
	"71": sysfsPin{
		pin:    317,
		pwmPin: -1,
	},
	// UART0CT
	"72": sysfsPin{
		pin:    412,
		pwmPin: -1,
	},
	// I2C1SCL
	"73": sysfsPin{
		pin:    418,
		pwmPin: -1,
	},
	// UART1TX
	"74": sysfsPin{
		pin:    484,
		pwmPin: -1,
	},
	// I2C2SDA
	"75": sysfsPin{
		pin:    319,
		pwmPin: -1,
	},
	// UART1RX
	"76": sysfsPin{
		pin:    483,
		pwmPin: -1,
	},
	// I2C2SCL
	"77": sysfsPin{
		pin:    320,
		pwmPin: -1,
	},
	// UART1RT
	"78": sysfsPin{
		pin:    485,
		pwmPin: -1,
	},
	// RTC_CLK
	"79": sysfsPin{
		pin:    367,
		pwmPin: -1,
	},
	// UART1CT
	"80": sysfsPin{
		pin:    486,
		pwmPin: -1,
	},

	// Built-in LEDs
	// LED100
	"100": sysfsPin{
		pin:    337,
		pwmPin: -1,
	},
	// LED101
	"101": sysfsPin{
		pin:    338,
		pwmPin: -1,
	},
	// LED102
	"102": sysfsPin{
		pin:    339,
		pwmPin: -1,
	},
	// LED103
	"103": sysfsPin{
		pin:    340,
		pwmPin: -1,
	},
	// LEDWIFI
	"104": sysfsPin{
		pin:    438,
		pwmPin: -1,
	},
	// LEDBT
	"105": sysfsPin{
		pin:    439,
		pwmPin: -1,
	},
}

// NewJouleAdaptor returns a new JouleAdaptor with specified name
func NewJouleAdaptor(name string) *JouleAdaptor {
	return &JouleAdaptor{
		name: name,
		connect: func(e *JouleAdaptor) (err error) {
			return
		},
	}
}

// Name returns the JouleAdaptors name
func (e *JouleAdaptor) Name() string { return e.name }

// Connect initializes the Joule for use with the Arduino beakout board
func (e *JouleAdaptor) Connect() (errs []error) {
	e.digitalPins = make(map[int]sysfs.DigitalPin)
	e.pwmPins = make(map[int]*pwmPin)
	if err := e.connect(e); err != nil {
		return []error{err}
	}
	return
}

// Finalize releases all i2c devices and exported digital and pwm pins.
func (e *JouleAdaptor) Finalize() (errs []error) {
	for _, pin := range e.digitalPins {
		if pin != nil {
			if err := pin.Unexport(); err != nil {
				errs = append(errs, err)
			}
		}
	}
	for _, pin := range e.pwmPins {
		if pin != nil {
			if err := pin.enable("0"); err != nil {
				errs = append(errs, err)
			}
			if err := pin.unexport(); err != nil {
				errs = append(errs, err)
			}
		}
	}
	if e.i2cDevice != nil {
		if err := e.i2cDevice.Close(); errs != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// digitalPin returns matched digitalPin for specified values
func (e *JouleAdaptor) digitalPin(pin string, dir string) (sysfsPin sysfs.DigitalPin, err error) {
	i := sysfsPinMap[pin]
	if e.digitalPins[i.pin] == nil {
		e.digitalPins[i.pin] = sysfs.NewDigitalPin(i.pin)
		if err = e.digitalPins[i.pin].Export(); err != nil {
			// TODO: log error
			return
		}
	}

	if dir == "in" {
		if err = e.digitalPins[i.pin].Direction(sysfs.IN); err != nil {
			return
		}
	} else if dir == "out" {
		if err = e.digitalPins[i.pin].Direction(sysfs.OUT); err != nil {
			return
		}
	}
	return e.digitalPins[i.pin], nil
}

// DigitalRead reads digital value from pin
func (e *JouleAdaptor) DigitalRead(pin string) (i int, err error) {
	sysfsPin, err := e.digitalPin(pin, "in")
	if err != nil {
		return
	}
	return sysfsPin.Read()
}

// DigitalWrite writes a value to the pin. Acceptable values are 1 or 0.
func (e *JouleAdaptor) DigitalWrite(pin string, val byte) (err error) {
	sysfsPin, err := e.digitalPin(pin, "out")
	if err != nil {
		return
	}
	return sysfsPin.Write(int(val))
}

// PwmWrite writes the 0-254 value to the specified pin
func (e *JouleAdaptor) PwmWrite(pin string, val byte) (err error) {
	sysPin := sysfsPinMap[pin]
	if sysPin.pwmPin != -1 {
		if e.pwmPins[sysPin.pwmPin] == nil {
			if err = e.DigitalWrite(pin, 1); err != nil {
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

// I2cStart initializes i2c device for addresss
func (e *JouleAdaptor) I2cStart(address int) (err error) {
	if e.i2cDevice != nil {
		return
	}

	// TODO: handle the additional I2C buses
	e.i2cDevice, err = sysfs.NewI2cDevice("/dev/i2c-0", address)
	return
}

// I2cWrite writes data to i2c device
func (e *JouleAdaptor) I2cWrite(address int, data []byte) (err error) {
	if err = e.i2cDevice.SetAddress(address); err != nil {
		return err
	}
	_, err = e.i2cDevice.Write(data)
	return
}

// I2cRead returns size bytes from the i2c device
func (e *JouleAdaptor) I2cRead(address int, size int) (data []byte, err error) {
	data = make([]byte, size)
	if err = e.i2cDevice.SetAddress(address); err != nil {
		return
	}
	_, err = e.i2cDevice.Read(data)
	return
}
