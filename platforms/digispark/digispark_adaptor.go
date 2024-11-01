//go:build libusb
// +build libusb

package digispark

import (
	"errors"
	"fmt"
	"strconv"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/i2c"
)

// ErrConnection is the error resulting of a connection error with the digispark
var ErrConnection = errors.New("connection error")

// Adaptor is the Gobot Adaptor for the Digispark
type Adaptor struct {
	name       string
	littleWire lw
	servo      bool
	pwm        bool
	i2c        bool
	connect    func(*Adaptor) error
}

// NewAdaptor returns a new Digispark Adaptor
func NewAdaptor() *Adaptor {
	return &Adaptor{
		name: gobot.DefaultName("Digispark"),
		connect: func(d *Adaptor) error {
			d.littleWire = littleWireConnect()
			//nolint:forcetypeassert // ok here
			if d.littleWire.(*littleWire).lwHandle == nil {
				return ErrConnection
			}
			return nil
		},
	}
}

// Name returns the Digispark adaptors name
func (d *Adaptor) Name() string { return d.name }

// SetName sets the Digispark adaptors name
func (d *Adaptor) SetName(n string) { d.name = n }

// Connect starts a connection to the digispark
func (d *Adaptor) Connect() error {
	return d.connect(d)
}

// Finalize implements the Adaptor interface
func (d *Adaptor) Finalize() error { return nil }

// DigitalWrite writes a value to the pin. Acceptable values are 1 or 0.
func (d *Adaptor) DigitalWrite(pin string, level byte) error {
	p, err := strconv.Atoi(pin)
	if err != nil {
		return err
	}

	//nolint:gosec // TODO: fix later
	if err := d.littleWire.pinMode(uint8(p), 0); err != nil {
		return err
	}

	//nolint:gosec // TODO: fix later
	return d.littleWire.digitalWrite(uint8(p), level)
}

// PwmWrite writes the 0-254 value to the specified pin
func (d *Adaptor) PwmWrite(pin string, value byte) error {
	if !d.pwm {
		if err := d.littleWire.pwmInit(); err != nil {
			return err
		}

		if err := d.littleWire.pwmUpdatePrescaler(1); err != nil {
			return err
		}
		d.pwm = true
	}

	return d.littleWire.pwmUpdateCompare(value, value)
}

// ServoWrite writes the 0-180 degree val to the specified pin.
func (d *Adaptor) ServoWrite(pin string, angle uint8) error {
	if !d.servo {
		if err := d.littleWire.servoInit(); err != nil {
			return err
		}
		d.servo = true
	}
	return d.littleWire.servoUpdateLocation(angle, angle)
}

// GetI2cConnection returns an i2c connection to a device on a specified bus.
// Only supports bus number 0
func (d *Adaptor) GetI2cConnection(address int, bus int) (i2c.Connection, error) {
	if bus != 0 {
		return nil, fmt.Errorf("Invalid bus number %d, only 0 is supported", bus)
	}
	//nolint:gosec // TODO: fix later
	c := NewDigisparkI2cConnection(d, uint8(address))
	if err := c.Init(); err != nil {
		return nil, err
	}
	return i2c.Connection(c), nil
}

// DefaultI2cBus returns the default i2c bus for this platform
func (d *Adaptor) DefaultI2cBus() int {
	return 0
}
