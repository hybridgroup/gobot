package digispark

import (
	"errors"
	"strconv"

	"gobot.io/x/gobot"
)

// ErrConnection is the error resulting of a connection error with the digispark
var ErrConnection = errors.New("connection error")

// Adaptor is the Gobot Adaptor for the Digispark
type Adaptor struct {
	name       string
	littleWire lw
	servo      bool
	pwm        bool
	connect    func(*Adaptor) (err error)
}

// NewAdaptor returns a new Digispark Adaptor
func NewAdaptor() *Adaptor {
	return &Adaptor{
		name: gobot.DefaultName("Digispark"),
		connect: func(d *Adaptor) (err error) {
			d.littleWire = littleWireConnect()
			if d.littleWire.(*littleWire).lwHandle == nil {
				return ErrConnection
			}
			return
		},
	}
}

// Name returns the Digispark Adaptors name
func (d *Adaptor) Name() string { return d.name }

// SetName sets the Digispark Adaptors name
func (d *Adaptor) SetName(n string) { d.name = n }

// Connect starts a connection to the digispark
func (d *Adaptor) Connect() (err error) {
	err = d.connect(d)
	return
}

// Finalize implements the Adaptor interface
func (d *Adaptor) Finalize() (err error) { return }

// DigitalWrite writes a value to the pin. Acceptable values are 1 or 0.
func (d *Adaptor) DigitalWrite(pin string, level byte) (err error) {
	p, err := strconv.Atoi(pin)

	if err != nil {
		return
	}

	if err = d.littleWire.pinMode(uint8(p), 0); err != nil {
		return
	}

	return d.littleWire.digitalWrite(uint8(p), level)
}

// PwmWrite writes the 0-254 value to the specified pin
func (d *Adaptor) PwmWrite(pin string, value byte) (err error) {
	if !d.pwm {
		if err = d.littleWire.pwmInit(); err != nil {
			return
		}

		if err = d.littleWire.pwmUpdatePrescaler(1); err != nil {
			return
		}
		d.pwm = true
	}

	return d.littleWire.pwmUpdateCompare(value, value)
}

// ServoWrite writes the 0-180 degree val to the specified pin.
func (d *Adaptor) ServoWrite(pin string, angle uint8) (err error) {
	if !d.servo {
		if err = d.littleWire.servoInit(); err != nil {
			return
		}
		d.servo = true
	}
	return d.littleWire.servoUpdateLocation(angle, angle)
}
