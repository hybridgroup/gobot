package firmata

import (
	"io"
	"strconv"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata/client"
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

type firmataBoard interface {
	Connect(io.ReadWriteCloser) error
	Disconnect() error
	Pins() []client.Pin
	AnalogWrite(int, int) error
	SetPinMode(int, int) error
	ReportAnalog(int, int) error
	ReportDigital(int, int) error
	DigitalWrite(int, int) error
	I2cRead(int, int) error
	I2cWrite(int, []byte) error
	I2cConfig(int) error
	Event(string) *gobot.Event
}

// FirmataAdaptor is the Gobot Adaptor for Firmata based boards
type FirmataAdaptor struct {
	name   string
	port   string
	board  firmataBoard
	conn   io.ReadWriteCloser
	openSP func(port string) (io.ReadWriteCloser, error)
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
		name:  name,
		port:  "",
		conn:  nil,
		board: client.New(),
		openSP: func(port string) (io.ReadWriteCloser, error) {
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
		sp, err := f.openSP(f.Port())
		if err != nil {
			return []error{err}
		}
		f.conn = sp
	}
	if err := f.board.Connect(f.conn); err != nil {
		return []error{err}
	}
	return
}

// Disconnect closes the io connection to the board
func (f *FirmataAdaptor) Disconnect() (err error) {
	if f.board != nil {
		return f.board.Disconnect()
	}
	return nil
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

	if f.board.Pins()[p].Mode != client.Servo {
		err = f.board.SetPinMode(p, client.Servo)
		if err != nil {
			return err
		}
	}
	err = f.board.AnalogWrite(p, int(angle))
	return
}

// PwmWrite writes the 0-254 value to the specified pin
func (f *FirmataAdaptor) PwmWrite(pin string, level byte) (err error) {
	p, err := strconv.Atoi(pin)
	if err != nil {
		return err
	}

	if f.board.Pins()[p].Mode != client.Pwm {
		err = f.board.SetPinMode(p, client.Pwm)
		if err != nil {
			return err
		}
	}
	err = f.board.AnalogWrite(p, int(level))
	return
}

// DigitalWrite writes a value to the pin. Acceptable values are 1 or 0.
func (f *FirmataAdaptor) DigitalWrite(pin string, level byte) (err error) {
	p, err := strconv.Atoi(pin)
	if err != nil {
		return
	}

	if f.board.Pins()[p].Mode != client.Output {
		err = f.board.SetPinMode(p, client.Output)
		if err != nil {
			return
		}
	}

	err = f.board.DigitalWrite(p, int(level))
	return
}

// DigitalRead retrieves digital value from specified pin.
// Returns -1 if the response from the board has timed out
func (f *FirmataAdaptor) DigitalRead(pin string) (val int, err error) {
	p, err := strconv.Atoi(pin)
	if err != nil {
		return
	}

	if f.board.Pins()[p].Mode != client.Input {
		if err = f.board.SetPinMode(p, client.Input); err != nil {
			return
		}
		if err = f.board.ReportDigital(p, 1); err != nil {
			return
		}
		<-time.After(10 * time.Millisecond)
	}

	return f.board.Pins()[p].Value, nil
}

// AnalogRead retrieves value from analog pin.
// Returns -1 if the response from the board has timed out
func (f *FirmataAdaptor) AnalogRead(pin string) (val int, err error) {
	p, err := strconv.Atoi(pin)
	if err != nil {
		return
	}

	p = f.digitalPin(p)

	if f.board.Pins()[p].Mode != client.Analog {
		if err = f.board.SetPinMode(p, client.Analog); err != nil {
			return
		}

		if err = f.board.ReportAnalog(p, 1); err != nil {
			return
		}
		<-time.After(10 * time.Millisecond)
	}

	return f.board.Pins()[p].Value, nil
}

// digitalPin converts pin number to digital mapping
func (f *FirmataAdaptor) digitalPin(pin int) int {
	return pin + 14
}

// I2cStart starts an i2c device at specified address
func (f *FirmataAdaptor) I2cStart(address int) (err error) {
	return f.board.I2cConfig(0)
}

// I2cRead returns size bytes from the i2c device
// Returns an empty array if the response from the board has timed out
func (f *FirmataAdaptor) I2cRead(address int, size int) (data []byte, err error) {
	ret := make(chan []byte)

	if err = f.board.I2cRead(address, size); err != nil {
		return
	}

	gobot.Once(f.board.Event("I2cReply"), func(data interface{}) {
		ret <- data.(client.I2cReply).Data
	})

	data = <-ret

	return
}

// I2cWrite writes data to i2c device
func (f *FirmataAdaptor) I2cWrite(address int, data []byte) (err error) {
	return f.board.I2cWrite(address, data)
}
