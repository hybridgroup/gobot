package digispark

import (
	"github.com/hybridgroup/gobot"
	"strconv"
)

type DigisparkAdaptor struct {
	gobot.Adaptor
	littleWire *LittleWire
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
			d.littleWire = LittleWireConnect()
		},
	}
}

// Connect starts connection to digispark, returns true if successful
func (d *DigisparkAdaptor) Connect() bool {
	d.connect(d)
	d.SetConnected(true)
	return true
}

// Reconnect retries connection to digispark, returns true if successful
func (d *DigisparkAdaptor) Reconnect() bool {
	return d.Connect()
}

// Finalize returns true if finalization is successful
func (d *DigisparkAdaptor) Finalize() bool { return true }

// Disconnect returns true if connection to digispark is ended successfully
func (d *DigisparkAdaptor) Disconnect() bool { return true }

// DigitalWrite writes level to specified pin using littlewire
func (d *DigisparkAdaptor) DigitalWrite(pin string, level byte) {
	p, _ := strconv.Atoi(pin)

	d.littleWire.PinMode(uint8(p), 0)
	d.littleWire.DigitalWrite(uint8(p), level)
}

// DigitalRead (not yet implemented)
func (d *DigisparkAdaptor) DigitalRead(pin string, level byte) {}

// PwmWrite updates pwm pin with sent value
func (d *DigisparkAdaptor) PwmWrite(pin string, value byte) {
	if d.pwm == false {
		d.littleWire.PwmInit()
		d.littleWire.PwmUpdatePrescaler(1)
		d.pwm = true
	}
	d.littleWire.PwmUpdateCompare(value, value)
}

// AnalogRead (not yet implemented)
func (d *DigisparkAdaptor) AnalogRead(string) int { return -1 }

// InitServo (not yet implemented)
func (d *DigisparkAdaptor) InitServo() {}

// ServoWrite updates servo location with specified angle
func (d *DigisparkAdaptor) ServoWrite(pin string, angle uint8) {
	if d.servo == false {
		d.littleWire.ServoInit()
		d.servo = true
	}
	d.littleWire.ServoUpdateLocation(angle, angle)
}

// I2cStart (not yet implemented)
func (d *DigisparkAdaptor) I2cStart(byte) {}

// I2cRead (not yet implemented)
func (d *DigisparkAdaptor) I2cRead(uint16) []uint16 { return make([]uint16, 0) }

// I2cWrite (not yet implemented)
func (d *DigisparkAdaptor) I2cWrite([]uint16) {}
