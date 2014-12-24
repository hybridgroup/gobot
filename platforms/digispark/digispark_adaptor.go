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

var ErrConnection = errors.New("connection error")

type DigisparkAdaptor struct {
	name       string
	littleWire lw
	servo      bool
	pwm        bool
	connect    func(*DigisparkAdaptor) (err error)
}

// NewDigisparkAdaptor create a Digispark adaptor with specified name
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

func (d *DigisparkAdaptor) Name() string { return d.name }

// Connect starts connection to digispark, returns true if successful
func (d *DigisparkAdaptor) Connect() (errs []error) {
	if err := d.connect(d); err != nil {
		return []error{err}
	}
	return
}

// Finalize returns true if finalization is successful
func (d *DigisparkAdaptor) Finalize() (errs []error) { return }

// DigitalWrite writes level to specified pin using littlewire
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

// PwmWrite updates pwm pin with sent value
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

// ServoWrite updates servo location with specified angle
func (d *DigisparkAdaptor) ServoWrite(pin string, angle uint8) (err error) {
	if d.servo == false {
		if err = d.littleWire.servoInit(); err != nil {
			return
		}
		d.servo = true
	}
	return d.littleWire.servoUpdateLocation(angle, angle)
}
