package digispark

import (
	"strconv"

	"github.com/hybridgroup/gobot"
)

var _ gobot.AdaptorInterface = (*DigisparkAdaptor)(nil)

type DigisparkAdaptor struct {
	gobot.Adaptor
	littleWire lw
	servo      bool
	pwm        bool
	connect    func(*DigisparkAdaptor)
}

// NewDigisparkAdaptor create a Digispark adaptor with specified name
func NewDigisparkAdaptor(name string) *DigisparkAdaptor {
	return &DigisparkAdaptor{
		Adaptor: *gobot.NewAdaptor(
			name,
			"DigisparkAdaptor",
		),
		connect: func(d *DigisparkAdaptor) {
			d.littleWire = littleWireConnect()
		},
	}
}

// Connect starts connection to digispark, returns true if successful
func (d *DigisparkAdaptor) Connect() error {
	d.connect(d)
	return nil
}

// Finalize returns true if finalization is successful
func (d *DigisparkAdaptor) Finalize() error { return nil }

// DigitalWrite writes level to specified pin using littlewire
func (d *DigisparkAdaptor) DigitalWrite(pin string, level byte) {
	p, _ := strconv.Atoi(pin)

	d.littleWire.pinMode(uint8(p), 0)
	d.littleWire.digitalWrite(uint8(p), level)
}

// DigitalRead (not yet implemented)
func (d *DigisparkAdaptor) DigitalRead(pin string) int {
	return -1
}

// PwmWrite updates pwm pin with sent value
func (d *DigisparkAdaptor) PwmWrite(pin string, value byte) {
	if d.pwm == false {
		d.littleWire.pwmInit()
		d.littleWire.pwmUpdatePrescaler(1)
		d.pwm = true
	}
	d.littleWire.pwmUpdateCompare(value, value)
}

// InitServo (not yet implemented)
func (d *DigisparkAdaptor) InitServo() {}

// ServoWrite updates servo location with specified angle
func (d *DigisparkAdaptor) ServoWrite(pin string, angle uint8) {
	if d.servo == false {
		d.littleWire.servoInit()
		d.servo = true
	}
	d.littleWire.servoUpdateLocation(angle, angle)
}
