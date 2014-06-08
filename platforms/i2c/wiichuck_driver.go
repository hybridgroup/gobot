package i2c

import (
	"fmt"
	"github.com/hybridgroup/gobot"
)

type WiichuckDriver struct {
	gobot.Driver
	Adaptor  I2cInterface
	joystick map[string]float64
	data     map[string]float64
}

func NewWiichuckDriver(a I2cInterface, name string) *WiichuckDriver {
	return &WiichuckDriver{
		Driver: gobot.Driver{
			Name: name,
			Events: map[string]chan interface{}{
				"z_button": make(chan interface{}),
				"c_button": make(chan interface{}),
				"joystick": make(chan interface{}),
			},
		},
		joystick: map[string]float64{
			"sy_origin": -1,
			"sx_origin": -1,
		},
		data: map[string]float64{
			"sx": 0,
			"sy": 0,
			"z":  0,
			"c":  0,
		},
		Adaptor: a,
	}
}

func (w *WiichuckDriver) Start() bool {
	w.Adaptor.I2cStart(0x52)
	gobot.Every(w.Interval, func() {
		w.Adaptor.I2cWrite([]byte{0x40, 0x00})
		w.Adaptor.I2cWrite([]byte{0x00})
		new_value := w.Adaptor.I2cRead(6)
		if len(new_value) == 6 {
			w.update(new_value)
		}
	})
	return true
}
func (w *WiichuckDriver) Init() bool { return true }
func (w *WiichuckDriver) Halt() bool { return true }

func (w *WiichuckDriver) update(value []byte) {
	if w.isEncrypted(value) {
		fmt.Println("Encrypted bytes from wii device!")
	} else {
		w.parse(value)
		w.adjustOrigins()
		w.updateButtons()
		w.updateJoystick()
	}
}

func (w *WiichuckDriver) setJoystickDefaultValue(joystick_axis string, default_value float64) {
	if w.joystick[joystick_axis] == -1 {
		w.joystick[joystick_axis] = default_value
	}
}

func (w *WiichuckDriver) calculateJoystickValue(axis float64, origin float64) float64 {
	return float64(axis - origin)
}

func (w *WiichuckDriver) isEncrypted(value []byte) bool {
	if value[0] == value[1] && value[2] == value[3] && value[4] == value[5] {
		return true
	} else {
		return false
	}
}

func (w *WiichuckDriver) decode(x byte) float64 {
	return float64((x ^ 0x17) + 0x17)
}

func (w *WiichuckDriver) adjustOrigins() {
	w.setJoystickDefaultValue("sy_origin", w.data["sy"])
	w.setJoystickDefaultValue("sx_origin", w.data["sx"])
}

func (w *WiichuckDriver) updateButtons() {
	if w.data["c"] == 0 {
		gobot.Publish(w.Events["c_button"], true)
	}
	if w.data["z"] == 0 {
		gobot.Publish(w.Events["z_button"], true)
	}
}

func (w *WiichuckDriver) updateJoystick() {
	gobot.Publish(w.Events["joystick"], map[string]float64{
		"x": w.calculateJoystickValue(w.data["sx"], w.joystick["sx_origin"]),
		"y": w.calculateJoystickValue(w.data["sy"], w.joystick["sy_origin"]),
	})
}

func (w *WiichuckDriver) parse(value []byte) {
	w.data["sx"] = w.decode(value[0])
	w.data["sy"] = w.decode(value[1])
	w.data["z"] = float64(uint8(w.decode(value[5])) & 0x01)
	w.data["c"] = float64(uint8(w.decode(value[5])) & 0x02)
}
