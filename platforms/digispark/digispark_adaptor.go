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

func (d *DigisparkAdaptor) Connect() bool {
	d.connect(d)
	d.SetConnected(true)
	return true
}

func (d *DigisparkAdaptor) Reconnect() bool {
	return d.Connect()
}

func (d *DigisparkAdaptor) Finalize() bool   { return true }
func (d *DigisparkAdaptor) Disconnect() bool { return true }

func (d *DigisparkAdaptor) DigitalWrite(pin string, level byte) {
	p, _ := strconv.Atoi(pin)

	d.littleWire.PinMode(uint8(p), 0)
	d.littleWire.DigitalWrite(uint8(p), level)
}
func (d *DigisparkAdaptor) DigitalRead(pin string, level byte) {}
func (d *DigisparkAdaptor) PwmWrite(pin string, value byte) {
	if d.pwm == false {
		d.littleWire.PwmInit()
		d.littleWire.PwmUpdatePrescaler(1)
		d.pwm = true
	}
	d.littleWire.PwmUpdateCompare(value, value)
}
func (d *DigisparkAdaptor) AnalogRead(string) int { return -1 }

func (d *DigisparkAdaptor) InitServo() {}
func (d *DigisparkAdaptor) ServoWrite(pin string, angle uint8) {
	if d.servo == false {
		d.littleWire.ServoInit()
		d.servo = true
	}
	d.littleWire.ServoUpdateLocation(angle, angle)
}

func (d *DigisparkAdaptor) I2cStart(byte)           {}
func (d *DigisparkAdaptor) I2cRead(uint16) []uint16 { return make([]uint16, 0) }
func (d *DigisparkAdaptor) I2cWrite([]uint16)       {}
