package gobotDigispark

import (
	"github.com/hybridgroup/gobot"
	"strconv"
)

type DigisparkAdaptor struct {
	gobot.Adaptor
	LittleWire *LittleWire
	servo      bool
	pwm        bool
}

var connect = func() *LittleWire {
	return LittleWireConnect()
}

func (da *DigisparkAdaptor) Connect() bool {
	da.LittleWire = connect()
	da.Connected = true
	return true
}

func (da *DigisparkAdaptor) Reconnect() bool {
	return da.Connect()
}

func (da *DigisparkAdaptor) Finalize() bool   { return true }
func (da *DigisparkAdaptor) Disconnect() bool { return true }

func (da *DigisparkAdaptor) DigitalWrite(pin string, level byte) {
	p, _ := strconv.Atoi(pin)

	da.LittleWire.PinMode(uint8(p), 0)
	da.LittleWire.DigitalWrite(uint8(p), level)
}
func (da *DigisparkAdaptor) DigitalRead(pin string, level byte) {}
func (da *DigisparkAdaptor) PwmWrite(pin string, value byte) {
	if da.pwm == false {
		da.LittleWire.PwmInit()
		da.LittleWire.PwmUpdatePrescaler(1)
		da.pwm = true
	}
	da.LittleWire.PwmUpdateCompare(value, value)
}
func (da *DigisparkAdaptor) AnalogRead(string) int { return -1 }

func (da *DigisparkAdaptor) InitServo() {}
func (da *DigisparkAdaptor) ServoWrite(pin string, angle uint8) {
	if da.servo == false {
		da.LittleWire.ServoInit()
		da.servo = true
	}
	da.LittleWire.ServoUpdateLocation(angle, angle)
}

func (da *DigisparkAdaptor) I2cStart(byte)           {}
func (da *DigisparkAdaptor) I2cRead(uint16) []uint16 { return make([]uint16, 0) }
func (da *DigisparkAdaptor) I2cWrite([]uint16)       {}
