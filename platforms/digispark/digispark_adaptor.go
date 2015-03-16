package digispark

import (
	"errors"
	"strconv"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

var _ gobot.Adaptor = (*DigisparkAdaptor)(nil)

var _ gpio.DigitalWriter = (*DigisparkAdaptor)(nil)
var _ gpio.PwmWriter = (*DigisparkAdaptor)(nil)
var _ gpio.ServoWriter = (*DigisparkAdaptor)(nil)

// ErrConnection is the error resulting of a connection error with the digispark
var ErrConnection = errors.New("connection error")

// DigisparkAdaptor is the Gobot Adaptor for the Digispark
type DigisparkAdaptor struct {
	name       string
	littleWire lw
	servo      bool
	pwm        bool
	connect    func(*DigisparkAdaptor) (err error)
}

// NewDigisparkAdaptor returns a new DigisparkAdaptor with specified name
func NewDigisparkAdaptor(name string) *DigisparkAdaptor {
	return &DigisparkAdaptor{
		name: name,
		connect: func(d *DigisparkAdaptor) (err error) {
			d.littleWire = littleWireConnect()
			if d.littleWire.(*littleWire).lwHandle == nil {
				return ErrConnection
			}
			return
		},
	}
}

// Name returns the DigisparkAdaptors name
func (d *DigisparkAdaptor) Name() string { return d.name }

// Connect starts a connection to the digispark
func (d *DigisparkAdaptor) Connect() (errs []error) {
	if err := d.connect(d); err != nil {
		return []error{err}
	}
	return
}

// Finalize implements the Adaptor interface
func (d *DigisparkAdaptor) Finalize() (errs []error) { return }

// DigitalWrite writes a value to the pin. Acceptable values are 1 or 0.
func (d *DigisparkAdaptor) DigitalWrite(pin string, level byte) (err error) {
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
func (d *DigisparkAdaptor) PwmWrite(pin string, value byte) (err error) {
	if d.pwm == false {
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
func (d *DigisparkAdaptor) ServoWrite(pin string, angle uint8) (err error) {
	if d.servo == false {
		if err = d.littleWire.servoInit(); err != nil {
			return
		}
		d.servo = true
	}
	return d.littleWire.servoUpdateLocation(angle, angle)
}
