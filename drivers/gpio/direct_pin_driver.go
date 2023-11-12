package gpio

import (
	"strconv"

	"gobot.io/x/gobot/v2"
)

// DirectPinDriver represents a GPIO pin
type DirectPinDriver struct {
	name       string
	pin        string
	connection gobot.Connection
	gobot.Commander
}

// NewDirectPinDriver return a new DirectPinDriver given a Connection and pin.
//
// Adds the following API Commands:
//
//	"DigitalRead" - See DirectPinDriver.DigitalRead
//	"DigitalWrite" - See DirectPinDriver.DigitalWrite
//	"PwmWrite" - See DirectPinDriver.PwmWrite
//	"ServoWrite" - See DirectPinDriver.ServoWrite
func NewDirectPinDriver(a gobot.Connection, pin string) *DirectPinDriver {
	d := &DirectPinDriver{
		name:       gobot.DefaultName("DirectPin"),
		connection: a,
		pin:        pin,
		Commander:  gobot.NewCommander(),
	}

	d.AddCommand("DigitalRead", func(params map[string]interface{}) interface{} {
		val, err := d.DigitalRead()
		return map[string]interface{}{"val": val, "err": err}
	})
	d.AddCommand("DigitalWrite", func(params map[string]interface{}) interface{} {
		level, _ := strconv.Atoi(params["level"].(string))
		return d.DigitalWrite(byte(level))
	})
	d.AddCommand("PwmWrite", func(params map[string]interface{}) interface{} {
		level, _ := strconv.Atoi(params["level"].(string))
		return d.PwmWrite(byte(level))
	})
	d.AddCommand("ServoWrite", func(params map[string]interface{}) interface{} {
		level, _ := strconv.Atoi(params["level"].(string))
		return d.ServoWrite(byte(level))
	})

	return d
}

// Name returns the DirectPinDrivers name
func (d *DirectPinDriver) Name() string { return d.name }

// SetName sets the DirectPinDrivers name
func (d *DirectPinDriver) SetName(n string) { d.name = n }

// Pin returns the DirectPinDrivers pin
func (d *DirectPinDriver) Pin() string { return d.pin }

// Connection returns the DirectPinDrivers Connection
func (d *DirectPinDriver) Connection() gobot.Connection { return d.connection }

// Start implements the Driver interface
func (d *DirectPinDriver) Start() error { return nil }

// Halt implements the Driver interface
func (d *DirectPinDriver) Halt() error { return nil }

// Off turn off pin
func (d *DirectPinDriver) Off() error {
	if writer, ok := d.Connection().(DigitalWriter); ok {
		return writer.DigitalWrite(d.Pin(), byte(0))
	}

	return ErrDigitalWriteUnsupported
}

// On turn on pin
func (d *DirectPinDriver) On() error {
	if writer, ok := d.Connection().(DigitalWriter); ok {
		return writer.DigitalWrite(d.Pin(), byte(1))
	}

	return ErrDigitalWriteUnsupported
}

// DigitalRead returns the current digital state of the pin
func (d *DirectPinDriver) DigitalRead() (int, error) {
	if reader, ok := d.Connection().(DigitalReader); ok {
		return reader.DigitalRead(d.Pin())
	}

	return 0, ErrDigitalReadUnsupported
}

// DigitalWrite writes to the pin. Acceptable values are 1 or 0
func (d *DirectPinDriver) DigitalWrite(level byte) error {
	if writer, ok := d.Connection().(DigitalWriter); ok {
		return writer.DigitalWrite(d.Pin(), level)
	}

	return ErrDigitalWriteUnsupported
}

// PwmWrite writes the 0-254 value to the specified pin
func (d *DirectPinDriver) PwmWrite(level byte) error {
	if writer, ok := d.Connection().(PwmWriter); ok {
		return writer.PwmWrite(d.Pin(), level)
	}

	return ErrPwmWriteUnsupported
}

// ServoWrite writes value to the specified pin
func (d *DirectPinDriver) ServoWrite(level byte) error {
	if writer, ok := d.Connection().(ServoWriter); ok {
		return writer.ServoWrite(d.Pin(), level)
	}

	return ErrServoWriteUnsupported
}
