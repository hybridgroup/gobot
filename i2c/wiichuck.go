package gobotI2C

import (
	"fmt"
	"github.com/hybridgroup/gobot"
)

type Wiichuck struct {
	gobot.Driver
	Adaptor  I2cInterface
	joystick map[string]float64
	data     map[string]float64
}

func NewWiichuck(a I2cInterface) *Wiichuck {
	w := new(Wiichuck)
	w.Adaptor = a
	w.Events = make(map[string]chan interface{})
	w.Events["z_button"] = make(chan interface{})
	w.Events["c_button"] = make(chan interface{})
	w.Events["joystick"] = make(chan interface{})
	w.joystick = map[string]float64{
		"sy_origin": -1,
		"sx_origin": -1,
	}
	w.data = map[string]float64{
		"sx": 0,
		"sy": 0,
		"z":  0,
		"c":  0,
	}
	return w
}

func (w *Wiichuck) Start() bool {
	w.Adaptor.I2cStart(byte(0x52))
	gobot.Every(w.Interval, func() {
		w.Adaptor.I2cWrite([]uint16{uint16(0x40), uint16(0x00)})
		w.Adaptor.I2cWrite([]uint16{uint16(0x00)})
		new_value := w.Adaptor.I2cRead(uint16(6))
		if len(new_value) == 6 {
			w.update(new_value)
		}
	})
	return true
}
func (w *Wiichuck) Init() bool { return true }
func (w *Wiichuck) Halt() bool { return true }

func (w *Wiichuck) update(value []uint16) {
	if w.isEncrypted(value) {
		fmt.Println("Encrypted bytes from wii device!")
	} else {
		w.parse(value)
		w.adjustOrigins()
		w.updateButtons()
		w.updateJoystick()
	}
}

func (w *Wiichuck) setJoystickDefaultValue(joystick_axis string, default_value float64) {
	if w.joystick[joystick_axis] == -1 {
		w.joystick[joystick_axis] = default_value
	}
}

func (w *Wiichuck) calculateJoystickValue(axis float64, origin float64) float64 {
	return float64(axis - origin)
}

func (w *Wiichuck) isEncrypted(value []uint16) bool {
	if value[0] == value[1] && value[2] == value[3] && value[4] == value[5] {
		return true
	} else {
		return false
	}
}

func (w *Wiichuck) decode(x uint16) float64 {
	return float64((x ^ 0x17) + 0x17)
}

func (w *Wiichuck) adjustOrigins() {
	w.setJoystickDefaultValue("sy_origin", w.data["sy"])
	w.setJoystickDefaultValue("sx_origin", w.data["sx"])
}

func (w *Wiichuck) updateButtons() {
	if w.data["c"] == 0 {
		w.Events["c_button"] <- ""
	}
	if w.data["z"] == 0 {
		w.Events["z_button"] <- ""
	}
}

func (w *Wiichuck) updateJoystick() {
	w.Events["joystick"] <- map[string]float64{
		"x": w.calculateJoystickValue(w.data["sx"], w.joystick["sx_origin"]),
		"y": w.calculateJoystickValue(w.data["sy"], w.joystick["sy_origin"]),
	}
}

func (w *Wiichuck) parse(value []uint16) {
	w.data["sx"] = w.decode(value[0])
	w.data["sy"] = w.decode(value[1])
	w.data["z"] = float64(uint8(w.decode(value[5])) & 0x01)
	w.data["c"] = float64(uint8(w.decode(value[5])) & 0x02)
}
