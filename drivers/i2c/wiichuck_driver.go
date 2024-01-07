package i2c

import (
	"fmt"
	"sync"
	"time"

	"gobot.io/x/gobot/v2"
)

const (
	// Joystick event when the Wiichuck joystick is moved
	Joystick = "joystick"

	// C event when the Wiichuck "C" button is pressed
	C = "c"

	// Z event when the Wiichuck "C" button is pressed
	Z = "z"
)

const wiichuckDefaultAddress = 0x52

// WiichuckDriver contains the attributes for the i2c driver
type WiichuckDriver struct {
	*Driver
	interval  time.Duration
	pauseTime time.Duration
	gobot.Eventer
	mtx      sync.Mutex
	joystick map[string]float64
	data     map[string]float64
}

// NewWiichuckDriver creates a WiichuckDriver with specified i2c interface.
//
// Params:
//
//	c Connector - the Adaptor to use with this Driver
//
// Optional params:
//
//	i2c.WithBus(int):	bus to use with this driver
//	i2c.WithAddress(int):	address to use with this driver
func NewWiichuckDriver(c Connector, options ...func(Config)) *WiichuckDriver {
	w := &WiichuckDriver{
		Driver:    NewDriver(c, "Wiichuck", wiichuckDefaultAddress),
		interval:  10 * time.Millisecond,
		pauseTime: 1 * time.Millisecond,
		Eventer:   gobot.NewEventer(),
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
	w.afterStart = w.initialize

	for _, option := range options {
		option(w)
	}

	w.AddEvent(Z)
	w.AddEvent(C)
	w.AddEvent(Joystick)
	w.AddEvent(Error)

	return w
}

// Joystick returns the current value for the joystick
func (w *WiichuckDriver) Joystick() map[string]float64 {
	val := make(map[string]float64)
	w.mtx.Lock()
	defer w.mtx.Unlock()
	val["sx_origin"] = w.joystick["sx_origin"]
	val["sy_origin"] = w.joystick["sy_origin"]
	return val
}

// update parses value to update buttons and joystick.
// If value is encrypted, warning message is printed
func (w *WiichuckDriver) update(value []byte) error {
	if w.isEncrypted(value) {
		return fmt.Errorf("Encrypted bytes")
	}

	w.parse(value)
	w.adjustOrigins()
	w.updateButtons()
	w.updateJoystick()
	return nil
}

// setJoystickDefaultValue sets default value if value is -1
func (w *WiichuckDriver) setJoystickDefaultValue(joystickAxis string, defaultValue float64) {
	w.mtx.Lock()
	defer w.mtx.Unlock()
	if w.joystick[joystickAxis] == -1 {
		w.joystick[joystickAxis] = defaultValue
	}
}

// calculateJoystickValue returns distance between axis and origin
func (w *WiichuckDriver) calculateJoystickValue(axis float64, origin float64) float64 {
	return axis - origin
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
		w.Publish(w.Event(C), true)
	}
	if w.data["z"] == 0 {
		w.Publish(w.Event(Z), true)
	}
}

// updateJoystick publishes event with current x and y values for joystick
func (w *WiichuckDriver) updateJoystick() {
	joy := w.Joystick()
	w.Publish(w.Event(Joystick), map[string]float64{
		"x": w.calculateJoystickValue(w.data["sx"], joy["sx_origin"]),
		"y": w.calculateJoystickValue(w.data["sy"], joy["sy_origin"]),
	})
}

// parse sets driver values based on parsed value
func (w *WiichuckDriver) parse(value []byte) {
	w.data["sx"] = w.decode(value[0])
	w.data["sy"] = w.decode(value[1])
	w.data["z"] = float64(uint8(w.decode(value[5])) & 0x01)
	w.data["c"] = float64(uint8(w.decode(value[5])) & 0x02)
}

// reads from adaptor using specified interval to update with new value
func (w *WiichuckDriver) initialize() error {
	go func() {
		for {
			if _, err := w.connection.Write([]byte{0x40, 0x00}); err != nil {
				w.Publish(w.Event(Error), err)
				continue
			}
			time.Sleep(w.pauseTime)
			if _, err := w.connection.Write([]byte{0x00}); err != nil {
				w.Publish(w.Event(Error), err)
				continue
			}
			time.Sleep(w.pauseTime)
			newValue := make([]byte, 6)
			bytesRead, err := w.connection.Read(newValue)
			if err != nil {
				w.Publish(w.Event(Error), err)
				continue
			}
			if bytesRead == 6 {
				if err = w.update(newValue); err != nil {
					w.Publish(w.Event(Error), err)
					continue
				}
			}
			time.Sleep(w.interval)
		}
	}()
	return nil
}
