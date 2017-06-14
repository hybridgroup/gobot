package firmata

import (
	"fmt"
	"io"
	"strconv"
	"time"

	serial "go.bug.st/serial.v1"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/firmata/client"
)

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
	ServoConfig(int, int, int) error
	WriteSysex(data []byte) error
	gobot.Eventer
}

type FirmataAdaptor interface {
	Connect() (err error)
	Finalize() (err error)
	Name() string
	SetName(n string)
	WriteSysex(data []byte) error
	gobot.Eventer
}

// Adaptor is the Gobot Adaptor for Firmata based boards
type Adaptor struct {
	name       string
	port       string
	Board      firmataBoard
	conn       io.ReadWriteCloser
	PortOpener func(port string) (io.ReadWriteCloser, error)
	gobot.Eventer
}

// NewAdaptor returns a new Firmata Adaptor which optionally accepts:
//
//	string: port the Adaptor uses to connect to a serial port with a baude rate of 57600
//	io.ReadWriteCloser: connection the Adaptor uses to communication with the hardware
//
// If an io.ReadWriteCloser is not supplied, the Adaptor will open a connection
// to a serial port with a baude rate of 57600. If an io.ReadWriteCloser
// is supplied, then the Adaptor will use the provided io.ReadWriteCloser and use the
// string port as a label to be displayed in the log and api.
func NewAdaptor(args ...interface{}) *Adaptor {
	f := &Adaptor{
		name:  gobot.DefaultName("Firmata"),
		port:  "",
		conn:  nil,
		Board: client.New(),
		PortOpener: func(port string) (io.ReadWriteCloser, error) {
			return serial.Open(port, &serial.Mode{BaudRate: 57600})
		},
		Eventer: gobot.NewEventer(),
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
func (f *Adaptor) Connect() (err error) {
	if f.conn == nil {
		sp, e := f.PortOpener(f.Port())
		if e != nil {
			return e
		}
		f.conn = sp
	}
	if err = f.Board.Connect(f.conn); err != nil {
		return err
	}

	f.Board.On("SysexResponse", func(data interface{}) {
		f.Publish("SysexResponse", data)
	})

	return
}

// Disconnect closes the io connection to the Board
func (f *Adaptor) Disconnect() (err error) {
	if f.Board != nil {
		return f.Board.Disconnect()
	}
	return nil
}

// Finalize terminates the firmata connection
func (f *Adaptor) Finalize() (err error) {
	err = f.Disconnect()
	return err
}

// Port returns the Firmata Adaptors port
func (f *Adaptor) Port() string { return f.port }

// Name returns the Firmata Adaptors name
func (f *Adaptor) Name() string { return f.name }

// SetName sets the Firmata Adaptors name
func (f *Adaptor) SetName(n string) { f.name = n }

// ServoConfig sets the pulse width in microseconds for a pin attached to a servo
func (f *Adaptor) ServoConfig(pin string, min, max int) error {
	p, err := strconv.Atoi(pin)
	if err != nil {
		return err
	}

	return f.Board.ServoConfig(p, max, min)
}

// ServoWrite writes the 0-180 degree angle to the specified pin.
func (f *Adaptor) ServoWrite(pin string, angle byte) (err error) {
	p, err := strconv.Atoi(pin)
	if err != nil {
		return err
	}

	if f.Board.Pins()[p].Mode != client.Servo {
		err = f.Board.SetPinMode(p, client.Servo)
		if err != nil {
			return err
		}
	}
	err = f.Board.AnalogWrite(p, int(angle))
	return
}

// PwmWrite writes the 0-254 value to the specified pin
func (f *Adaptor) PwmWrite(pin string, level byte) (err error) {
	p, err := strconv.Atoi(pin)
	if err != nil {
		return err
	}

	if f.Board.Pins()[p].Mode != client.Pwm {
		err = f.Board.SetPinMode(p, client.Pwm)
		if err != nil {
			return err
		}
	}
	err = f.Board.AnalogWrite(p, int(level))
	return
}

// DigitalWrite writes a value to the pin. Acceptable values are 1 or 0.
func (f *Adaptor) DigitalWrite(pin string, level byte) (err error) {
	p, err := strconv.Atoi(pin)
	if err != nil {
		return
	}

	if f.Board.Pins()[p].Mode != client.Output {
		err = f.Board.SetPinMode(p, client.Output)
		if err != nil {
			return
		}
	}

	err = f.Board.DigitalWrite(p, int(level))
	return
}

// DigitalRead retrieves digital value from specified pin.
// Returns -1 if the response from the board has timed out
func (f *Adaptor) DigitalRead(pin string) (val int, err error) {
	p, err := strconv.Atoi(pin)
	if err != nil {
		return
	}

	if f.Board.Pins()[p].Mode != client.Input {
		if err = f.Board.SetPinMode(p, client.Input); err != nil {
			return
		}
		if err = f.Board.ReportDigital(p, 1); err != nil {
			return
		}
		<-time.After(10 * time.Millisecond)
	}

	return f.Board.Pins()[p].Value, nil
}

// AnalogRead retrieves value from analog pin.
// Returns -1 if the response from the board has timed out
func (f *Adaptor) AnalogRead(pin string) (val int, err error) {
	p, err := strconv.Atoi(pin)
	if err != nil {
		return
	}

	p = f.digitalPin(p)

	if f.Board.Pins()[p].Mode != client.Analog {
		if err = f.Board.SetPinMode(p, client.Analog); err != nil {
			return
		}

		if err = f.Board.ReportAnalog(p, 1); err != nil {
			return
		}
		<-time.After(10 * time.Millisecond)
	}

	return f.Board.Pins()[p].Value, nil
}

func (f *Adaptor) WriteSysex(data []byte) error {
	return f.Board.WriteSysex(data)
}

// digitalPin converts pin number to digital mapping
func (f *Adaptor) digitalPin(pin int) int {
	return pin + 14
}

// GetConnection returns an i2c connection to a device on a specified bus.
// Only supports bus number 0
func (f *Adaptor) GetConnection(address int, bus int) (connection i2c.Connection, err error) {
	if bus != 0 {
		return nil, fmt.Errorf("Invalid bus number %d, only 0 is supported", bus)
	}
	err = f.Board.I2cConfig(0)
	return NewFirmataI2cConnection(f, address), err
}

// GetDefaultBus returns the default i2c bus for this platform
func (f *Adaptor) GetDefaultBus() int {
	return 0
}
