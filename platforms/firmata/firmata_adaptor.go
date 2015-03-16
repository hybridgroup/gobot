package firmata

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/i2c"
	"github.com/tarm/goserial"
)

var _ gobot.Adaptor = (*FirmataAdaptor)(nil)

var _ gpio.DigitalReader = (*FirmataAdaptor)(nil)
var _ gpio.DigitalWriter = (*FirmataAdaptor)(nil)
var _ gpio.AnalogReader = (*FirmataAdaptor)(nil)
var _ gpio.PwmWriter = (*FirmataAdaptor)(nil)
var _ gpio.ServoWriter = (*FirmataAdaptor)(nil)

var _ i2c.I2c = (*FirmataAdaptor)(nil)

// FirmataAdaptor is the Gobot Adaptor for Firmata based boards
type FirmataAdaptor struct {
	name       string
	port       string
	board      *board
	i2cAddress byte
	conn       io.ReadWriteCloser
	connect    func(string) (io.ReadWriteCloser, error)
}

// NewFirmataAdaptor returns a new FirmataAdaptor with specified name and optionally accepts:
//
//	string: port the FirmataAdaptor uses to connect to a serial port with a baude rate of 57600
//	io.ReadWriteCloser: connection the FirmataAdaptor uses to communication with the hardware
//
// If an io.ReadWriteCloser is not supplied, the FirmataAdaptor will open a connection
// to a serial port with a baude rate of 57600. If an io.ReadWriteCloser
// is supplied, then the FirmataAdaptor will use the provided io.ReadWriteCloser and use the
// string port as a label to be displayed in the log and api.
func NewFirmataAdaptor(name string, args ...interface{}) *FirmataAdaptor {
	f := &FirmataAdaptor{
		name: name,
		port: "",
		conn: nil,
		connect: func(port string) (io.ReadWriteCloser, error) {
			return serial.OpenPort(&serial.Config{Name: port, Baud: 57600})
		},
	}

	for _, arg := range args {
		switch arg.(type) {
		case string:
			f.port = arg.(string)
		case io.ReadWriteCloser:
			f.conn = arg.(io.ReadWriteCloser)
		}
	}

	return f
}

// Connect starts a connection to the board.
func (f *FirmataAdaptor) Connect() (errs []error) {
	if f.conn == nil {
		sp, err := f.connect(f.Port())
		if err != nil {
			return []error{err}
		}
		f.conn = sp
	}
	f.board = newBoard(f.conn)
	f.board.connect()
	return
}

// Disconnect closes the io connection to the board
func (f *FirmataAdaptor) Disconnect() (err error) {
	if f.board != nil {
		return f.board.serial.Close()
	}
	return errors.New("no board connected")
}

// Finalize terminates the firmata connection
func (f *FirmataAdaptor) Finalize() (errs []error) {
	if err := f.Disconnect(); err != nil {
		return []error{err}
	}
	return
}

// Port returns the  FirmataAdaptors port
func (f *FirmataAdaptor) Port() string { return f.port }

// Name returns the  FirmataAdaptors name
func (f *FirmataAdaptor) Name() string { return f.name }

// ServoWrite writes the 0-180 degree angle to the specified pin.
func (f *FirmataAdaptor) ServoWrite(pin string, angle byte) (err error) {
	p, err := strconv.Atoi(pin)
	if err != nil {
		return err
	}

	err = f.board.setPinMode(byte(p), servo)
	if err != nil {
		return err
	}
	err = f.board.analogWrite(byte(p), angle)
	return
}

// PwmWrite writes the 0-254 value to the specified pin
func (f *FirmataAdaptor) PwmWrite(pin string, level byte) (err error) {
	p, err := strconv.Atoi(pin)
	if err != nil {
		return err
	}

	err = f.board.setPinMode(byte(p), pwm)
	if err != nil {
		return err
	}
	err = f.board.analogWrite(byte(p), level)
	return
}

// DigitalWrite writes a value to the pin. Acceptable values are 1 or 0.
func (f *FirmataAdaptor) DigitalWrite(pin string, level byte) (err error) {
	p, err := strconv.Atoi(pin)
	if err != nil {
		return
	}

	err = f.board.setPinMode(byte(p), output)
	if err != nil {
		return
	}

	err = f.board.digitalWrite(byte(p), level)
	return
}

// DigitalRead retrieves digital value from specified pin.
// Returns -1 if the response from the board has timed out
func (f *FirmataAdaptor) DigitalRead(pin string) (val int, err error) {
	ret := make(chan int)

	p, err := strconv.Atoi(pin)
	if err != nil {
		return
	}
	if err = f.board.setPinMode(byte(p), input); err != nil {
		return
	}
	if err = f.board.togglePinReporting(byte(p), high, reportDigital); err != nil {
		return
	}
	if err = f.board.readAndProcess(); err != nil {
		return
	}

	gobot.Once(f.board.events[fmt.Sprintf("digital_read_%v", pin)], func(data interface{}) {
		ret <- int(data.([]byte)[0])
	})

	select {
	case data := <-ret:
		return data, nil
	case <-time.After(10 * time.Millisecond):
	}
	return -1, nil
}

// AnalogRead retrieves value from analog pin.
// Returns -1 if the response from the board has timed out
func (f *FirmataAdaptor) AnalogRead(pin string) (val int, err error) {
	ret := make(chan int)

	p, err := strconv.Atoi(pin)
	if err != nil {
		return
	}
	// NOTE pins are numbered A0-A5, which translate to digital pins 14-19
	p = f.digitalPin(p)
	if err = f.board.setPinMode(byte(p), analog); err != nil {
		return
	}

	if err = f.board.togglePinReporting(byte(p), high, reportAnalog); err != nil {
		return
	}

	if err = f.board.readAndProcess(); err != nil {
		return
	}

	gobot.Once(f.board.events[fmt.Sprintf("analog_read_%v", pin)], func(data interface{}) {
		b := data.([]byte)
		ret <- int(uint(b[0])<<24 | uint(b[1])<<16 | uint(b[2])<<8 | uint(b[3]))
	})

	select {
	case data := <-ret:
		return data, nil
	case <-time.After(10 * time.Millisecond):
	}
	return -1, nil
}

// digitalPin converts pin number to digital mapping
func (f *FirmataAdaptor) digitalPin(pin int) int {
	return pin + 14
}

// I2cStart starts an i2c device at specified address
func (f *FirmataAdaptor) I2cStart(address byte) (err error) {
	f.i2cAddress = address
	return f.board.i2cConfig([]byte{0})
}

// I2cRead returns size bytes from the i2c device
// Returns an empty array if the response from the board has timed out
func (f *FirmataAdaptor) I2cRead(size uint) (data []byte, err error) {
	ret := make(chan []byte)
	if err = f.board.i2cReadRequest(f.i2cAddress, size); err != nil {
		return
	}

	if err = f.board.readAndProcess(); err != nil {
		return
	}

	gobot.Once(f.board.events["i2c_reply"], func(data interface{}) {
		ret <- data.(map[string][]byte)["data"]
	})

	select {
	case data := <-ret:
		return data, nil
	case <-time.After(10 * time.Millisecond):
	}
	return []byte{}, nil
}

// I2cWrite writes data to i2c device
func (f *FirmataAdaptor) I2cWrite(data []byte) (err error) {
	return f.board.i2cWriteRequest(f.i2cAddress, data)
}
