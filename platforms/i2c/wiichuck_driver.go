package i2c

import (
	"time"

	"github.com/hybridgroup/gobot"
)

var _ gobot.Driver = (*WiichuckDriver)(nil)

const wiichuckAddress = 0x52

type WiichuckDriver struct {
	name       string
	connection I2c
	interval   time.Duration
	gobot.Eventer
	joystick map[string]float64
	data     map[string]float64
}

// NewWiichuckDriver creates a WiichuckDriver with specified i2c interface and name.
//
// It adds the following events:
//	"z"- Get's triggered every interval amount of time if the z button is pressed
//	"c" - Get's triggered every interval amount of time if the c button is pressed
//	"joystick" - Get's triggered every "interval" amount of time if a joystick event occured, you can access values x, y
func NewWiichuckDriver(a I2c, name string, v ...time.Duration) *WiichuckDriver {
	w := &WiichuckDriver{
		name:       name,
		connection: a,
		interval:   10 * time.Millisecond,
		Eventer:    gobot.NewEventer(),
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
	}

	if len(v) > 0 {
		w.interval = v[0]
	}

	w.AddEvent(Z)
	w.AddEvent(C)
	w.AddEvent(Joystick)
	return w
}
func (w *WiichuckDriver) Name() string                 { return w.name }
func (w *WiichuckDriver) Connection() gobot.Connection { return w.connection.(gobot.Connection) }

// Start initilizes i2c and reads from adaptor
// using specified interval to update with new value
func (w *WiichuckDriver) Start() (errs []error) {
	if err := w.connection.I2cStart(wiichuckAddress); err != nil {
		return []error{err}
	}

	go func() {
		for {
			if err := w.connection.I2cWrite(wiichuckAddress, []byte{0x40, 0x00}); err != nil {
				gobot.Publish(w.Event(Error), err)
				continue
			}
			if err := w.connection.I2cWrite(wiichuckAddress, []byte{0x00}); err != nil {
				gobot.Publish(w.Event(Error), err)
				continue
			}
			newValue, err := w.connection.I2cRead(wiichuckAddress, 6)
			if err != nil {
				gobot.Publish(w.Event(Error), err)
				continue
			}
			if len(newValue) == 6 {
				if err = w.update(newValue); err != nil {
					gobot.Publish(w.Event(Error), err)
					continue
				}
			}
			<-time.After(w.interval)
		}
	}()
	return
}

// Halt returns true if driver is halted successfully
func (w *WiichuckDriver) Halt() (errs []error) { return }

// update parses value to update buttons and joystick.
// If value is encrypted, warning message is printed
func (w *WiichuckDriver) update(value []byte) (err error) {
	if w.isEncrypted(value) {
		return ErrEncryptedBytes
	} else {
		w.parse(value)
		w.adjustOrigins()
		w.updateButtons()
		w.updateJoystick()
	}
	return
}

// setJoystickDefaultValue sets default value if value is -1
func (w *WiichuckDriver) setJoystickDefaultValue(joystickAxis string, defaultValue float64) {
	if w.joystick[joystickAxis] == -1 {
		w.joystick[joystickAxis] = defaultValue
	}
}

// calculateJoystickValue returns distance between axis and origin
func (w *WiichuckDriver) calculateJoystickValue(axis float64, origin float64) float64 {
	return float64(axis - origin)
}

// isEncrypted returns true if value is encrypted
func (w *WiichuckDriver) isEncrypted(value []byte) bool {
	if value[0] == value[1] && value[2] == value[3] && value[4] == value[5] {
		return true
	}
	return false
}

// decode removes encoding from `x` byte
func (w *WiichuckDriver) decode(x byte) float64 {
	return float64((x ^ 0x17) + 0x17)
}

// adjustOrigins sets sy_origin and sx_origin with values from data
func (w *WiichuckDriver) adjustOrigins() {
	w.setJoystickDefaultValue("sy_origin", w.data["sy"])
	w.setJoystickDefaultValue("sx_origin", w.data["sx"])
}

// updateButtons publishes "c" and "x" events if present in data
func (w *WiichuckDriver) updateButtons() {
	if w.data["c"] == 0 {
		gobot.Publish(w.Event(C), true)
	}
	if w.data["z"] == 0 {
		gobot.Publish(w.Event(Z), true)
	}
}

// updateJoystick publishes event with current x and y values for joystick
func (w *WiichuckDriver) updateJoystick() {
	gobot.Publish(w.Event(Joystick), map[string]float64{
		"x": w.calculateJoystickValue(w.data["sx"], w.joystick["sx_origin"]),
		"y": w.calculateJoystickValue(w.data["sy"], w.joystick["sy_origin"]),
	})
}

// parse sets driver values based on parsed value
func (w *WiichuckDriver) parse(value []byte) {
	w.data["sx"] = w.decode(value[0])
	w.data["sy"] = w.decode(value[1])
	w.data["z"] = float64(uint8(w.decode(value[5])) & 0x01)
	w.data["c"] = float64(uint8(w.decode(value[5])) & 0x02)
}
