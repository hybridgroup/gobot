package firmata

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/tarm/goserial"
)

var _ gobot.AdaptorInterface = (*FirmataAdaptor)(nil)

type FirmataAdaptor struct {
	gobot.Adaptor
	board      *board
	i2cAddress byte
	connect    func(*FirmataAdaptor)
}

// NewFirmataAdaptor returns a new firmata adaptor with specified name
// Generates a connect function that opens serial communication in specified port
func NewFirmataAdaptor(name, port string) *FirmataAdaptor {
	return &FirmataAdaptor{
		Adaptor: *gobot.NewAdaptor(
			name,
			"FirmataAdaptor",
			port,
		),
		connect: func(f *FirmataAdaptor) {
			sp, err := serial.OpenPort(&serial.Config{Name: f.Port(), Baud: 57600})
			if err != nil {
				panic(err)
			}
			f.board = newBoard(sp)
		},
	}
}

// Connect returns true if connection to board is succesfull
func (f *FirmataAdaptor) Connect() error {
	f.connect(f)
	f.board.connect()
	f.SetConnected(true)
	return nil
}

// close finishes connection to serial port
// Prints error message on error
func (f *FirmataAdaptor) Disconnect() error {
	err := f.board.serial.Close()
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

// Finalize disconnects firmata adaptor
func (f *FirmataAdaptor) Finalize() error { return f.Disconnect() }

// InitServo (not yet implemented)
func (f *FirmataAdaptor) InitServo() {}

// ServoWrite sets angle form 0 to 360 to specified servo pin
func (f *FirmataAdaptor) ServoWrite(pin string, angle byte) {
	p, _ := strconv.Atoi(pin)

	f.board.setPinMode(byte(p), servo)
	f.board.analogWrite(byte(p), angle)
}

// PwmWrite writes analog value to specified pin
func (f *FirmataAdaptor) PwmWrite(pin string, level byte) {
	p, _ := strconv.Atoi(pin)

	f.board.setPinMode(byte(p), pwm)
	f.board.analogWrite(byte(p), level)
}

// DigitalWrite writes digital values to specified pin
func (f *FirmataAdaptor) DigitalWrite(pin string, level byte) {
	p, _ := strconv.Atoi(pin)

	f.board.setPinMode(byte(p), output)
	f.board.digitalWrite(byte(p), level)
}

// DigitalRead retrieves digital value from specified pin
// Returns -1 if response from board is timed out
func (f *FirmataAdaptor) DigitalRead(pin string) int {
	ret := make(chan int)

	p, _ := strconv.Atoi(pin)
	f.board.setPinMode(byte(p), input)
	f.board.togglePinReporting(byte(p), high, reportDigital)
	f.board.readAndProcess()

	gobot.Once(f.board.events[fmt.Sprintf("digital_read_%v", pin)], func(data interface{}) {
		ret <- int(data.([]byte)[0])
	})

	select {
	case data := <-ret:
		return data
	case <-time.After(10 * time.Millisecond):
	}
	return -1
}

// AnalogRead retrieves value from analog pin.
// NOTE pins are numbered A0-A5, which translate to digital pins 14-19
func (f *FirmataAdaptor) AnalogRead(pin string) int {
	ret := make(chan int)

	p, _ := strconv.Atoi(pin)
	p = f.digitalPin(p)
	f.board.setPinMode(byte(p), analog)
	f.board.togglePinReporting(byte(p), high, reportAnalog)
	f.board.readAndProcess()

	gobot.Once(f.board.events[fmt.Sprintf("analog_read_%v", pin)], func(data interface{}) {
		b := data.([]byte)
		ret <- int(uint(b[0])<<24 | uint(b[1])<<16 | uint(b[2])<<8 | uint(b[3]))
	})

	select {
	case data := <-ret:
		return data
	case <-time.After(10 * time.Millisecond):
	}
	return -1
}

// AnalogWrite writes value to ananlog pin
func (f *FirmataAdaptor) AnalogWrite(pin string, level byte) {
	f.PwmWrite(pin, level)
}

// digitalPin converts pin number to digital mapping
func (f *FirmataAdaptor) digitalPin(pin int) int {
	return pin + 14
}

// I2cStart initializes board with i2c configuration
func (f *FirmataAdaptor) I2cStart(address byte) {
	f.i2cAddress = address
	f.board.i2cConfig([]byte{0})
}

// I2cRead reads from I2c specified size
// Returns empty byte array if response is timed out
func (f *FirmataAdaptor) I2cRead(size uint) []byte {
	ret := make(chan []byte)
	f.board.i2cReadRequest(f.i2cAddress, size)

	f.board.readAndProcess()

	gobot.Once(f.board.events["i2c_reply"], func(data interface{}) {
		ret <- data.(map[string][]byte)["data"]
	})

	select {
	case data := <-ret:
		return data
	case <-time.After(10 * time.Millisecond):
	}
	return []byte{}
}

// I2cWrite retrieves i2c data
func (f *FirmataAdaptor) I2cWrite(data []byte) {
	f.board.i2cWriteRequest(f.i2cAddress, data)
}
