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

const (
	wiichuckDefaultAddress = 0x52
	wiichuckStopTimeout    = 2 * time.Second
)

// WiichuckDriver contains the attributes for the i2c driver
type WiichuckDriver struct {
	*Driver
	gobot.Eventer

	interval       time.Duration
	pauseTime      time.Duration
	mtx            sync.Mutex
	doneChan       chan struct{}
	reallyDoneChan chan struct{}
	joystick       map[string]float64
	data           map[string]float64
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
	d := &WiichuckDriver{
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
	d.afterStart = d.initialize
	d.beforeHalt = d.shutdown

	for _, option := range options {
		option(d)
	}

	d.AddEvent(Z)
	d.AddEvent(C)
	d.AddEvent(Joystick)
	d.AddEvent(Error)

	return d
}

// Joystick returns the current value for the joystick
func (d *WiichuckDriver) Joystick() map[string]float64 {
	val := make(map[string]float64)
	d.mtx.Lock()
	defer d.mtx.Unlock()
	val["sx_origin"] = d.joystick["sx_origin"]
	val["sy_origin"] = d.joystick["sy_origin"]
	return val
}

// update parses value to update buttons and joystick.
// If value is encrypted, warning message is printed
func (d *WiichuckDriver) update(value []byte) error {
	if d.isEncrypted(value) {
		return fmt.Errorf("encrypted bytes")
	}

	d.parse(value)
	d.adjustOrigins()
	d.updateButtons()
	d.updateJoystick()
	return nil
}

// setJoystickDefaultValue sets default value if value is -1
func (d *WiichuckDriver) setJoystickDefaultValue(joystickAxis string, defaultValue float64) {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	if d.joystick[joystickAxis] == -1 {
		d.joystick[joystickAxis] = defaultValue
	}
}

// calculateJoystickValue returns distance between axis and origin
func (d *WiichuckDriver) calculateJoystickValue(axis float64, origin float64) float64 {
	return axis - origin
}

// isEncrypted returns true if value is encrypted
func (d *WiichuckDriver) isEncrypted(value []byte) bool {
	if value[0] == value[1] && value[2] == value[3] && value[4] == value[5] {
		return true
	}
	return false
}

// decode removes encoding from `x` byte
func (d *WiichuckDriver) decode(x byte) float64 {
	return float64((x ^ 0x17) + 0x17)
}

// adjustOrigins sets sy_origin and sx_origin with values from data
func (d *WiichuckDriver) adjustOrigins() {
	d.setJoystickDefaultValue("sy_origin", d.data["sy"])
	d.setJoystickDefaultValue("sx_origin", d.data["sx"])
}

// updateButtons publishes "c" and "x" events if present in data
func (d *WiichuckDriver) updateButtons() {
	if d.data["c"] == 0 {
		d.Publish(d.Event(C), true)
	}
	if d.data["z"] == 0 {
		d.Publish(d.Event(Z), true)
	}
}

// updateJoystick publishes event with current x and y values for joystick
func (d *WiichuckDriver) updateJoystick() {
	joy := d.Joystick()
	d.Publish(d.Event(Joystick), map[string]float64{
		"x": d.calculateJoystickValue(d.data["sx"], joy["sx_origin"]),
		"y": d.calculateJoystickValue(d.data["sy"], joy["sy_origin"]),
	})
}

// parse sets driver values based on parsed value
func (d *WiichuckDriver) parse(value []byte) {
	d.data["sx"] = d.decode(value[0])
	d.data["sy"] = d.decode(value[1])
	d.data["z"] = float64(uint8(d.decode(value[5])) & 0x01)
	d.data["c"] = float64(uint8(d.decode(value[5])) & 0x02)
}

// reads from adaptor using specified interval to update with new value
func (d *WiichuckDriver) initialize() error {
	d.doneChan = make(chan struct{})
	d.reallyDoneChan = make(chan struct{})

	go func() {
		for {
			select {
			case <-d.doneChan:
				if d.reallyDoneChan != nil {
					close(d.reallyDoneChan)
				}
				return
			default:
				if _, err := d.connection.Write([]byte{0x40, 0x00}); err != nil {
					d.Publish(d.Event(Error), err)
					continue
				}
				time.Sleep(d.pauseTime)
				if _, err := d.connection.Write([]byte{0x00}); err != nil {
					d.Publish(d.Event(Error), err)
					continue
				}
				time.Sleep(d.pauseTime)
				newValue := make([]byte, 6)
				bytesRead, err := d.connection.Read(newValue)
				if err != nil {
					d.Publish(d.Event(Error), err)
					continue
				}
				if bytesRead == 6 {
					if err = d.update(newValue); err != nil {
						d.Publish(d.Event(Error), err)
						continue
					}
				}
				time.Sleep(d.interval)
			}
		}
	}()
	return nil
}

func (d *WiichuckDriver) shutdown() error {
	if d.doneChan != nil {
		close(d.doneChan)
		// wait until go routine is really finished, so connection can be set to nil safely
		select {
		case <-time.After(wiichuckStopTimeout):
			return fmt.Errorf("go routine not finished in %s", wiichuckStopTimeout)
		case <-d.reallyDoneChan:
		}
	}

	return nil
}
