package digispark

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hybridgroup/gobot"
)

var _ gobot.AdaptorInterface = (*DigisparkAdaptor)(nil)

type DigisparkAdaptor struct {
	gobot.Adaptor
	littleWire lw
	servo      bool
	pwm        bool
	connect    func(*DigisparkAdaptor) (err error)
}

// NewDigisparkAdaptor create a Digispark adaptor with specified name
func NewDigisparkAdaptor(name string) *DigisparkAdaptor {
	return &DigisparkAdaptor{
		Adaptor: *gobot.NewAdaptor(
			name,
			"DigisparkAdaptor",
		),
		connect: func(d *DigisparkAdaptor) (err error) {
			d.littleWire = littleWireConnect()
			if d.littleWire.(*littleWire).lwHandle == nil {
				return errors.New(fmt.Sprintf("Error connecting to %s", d.Name()))
			}
			return
		},
	}
}

// Connect starts connection to digispark, returns true if successful
func (d *DigisparkAdaptor) Connect() error {
	return d.connect(d)
}

// Finalize returns true if finalization is successful
func (d *DigisparkAdaptor) Finalize() error { return nil }

// DigitalWrite writes level to specified pin using littlewire
func (d *DigisparkAdaptor) DigitalWrite(pin string, level byte) (err error) {
	p, err := strconv.Atoi(pin)

	if err != nil {
		return
	}

	err = d.littleWire.pinMode(uint8(p), 0)
	if err != nil {
		return
	}
	err = d.littleWire.digitalWrite(uint8(p), level)
	return
}

// DigitalRead (not yet implemented)
func (d *DigisparkAdaptor) DigitalRead(pin string) (val int, err error) {
	err = errors.New("DigitalRead is not yet implemented")
	return
}

// PwmWrite updates pwm pin with sent value
func (d *DigisparkAdaptor) PwmWrite(pin string, value byte) (err error) {
	if d.pwm == false {
		err = d.littleWire.pwmInit()
		if err != nil {
			return err
		}
		err = d.littleWire.pwmUpdatePrescaler(1)
		if err != nil {
			return err
		}
		d.pwm = true
	}
	err = d.littleWire.pwmUpdateCompare(value, value)
	return
}

// InitServo (not yet implemented)
func (d *DigisparkAdaptor) InitServo() (err error) {
	return errors.New("InitServo is not yet implemented")
}

// ServoWrite updates servo location with specified angle
func (d *DigisparkAdaptor) ServoWrite(pin string, angle uint8) (err error) {
	if d.servo == false {
		err = d.littleWire.servoInit()
		if err != nil {
			return err
		}
		d.servo = true
	}
	err = d.littleWire.servoUpdateLocation(angle, angle)
	return
}
