package gpio

import (
	"strconv"

	"gobot.io/x/gobot"
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
// 	"DigitalRead" - See DirectPinDriver.DigitalRead
// 	"DigitalWrite" - See DirectPinDriver.DigitalWrite
// 	"AnalogWrite" - See DirectPinDriver.AnalogWrite
// 	"PwmWrite" - See DirectPinDriver.PwmWrite
// 	"ServoWrite" - See DirectPinDriver.ServoWrite
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
func (d *DirectPinDriver) Start() (err error) { return }

// Halt implements the Driver interface
func (d *DirectPinDriver) Halt() (err error) { return }

// Turn Off pin
func (d *DirectPinDriver) Off() (err error) {
	if writer, ok := d.Connection().(DigitalWriter); ok {
		return writer.DigitalWrite(d.Pin(), byte(0))
	}
	err = ErrDigitalWriteUnsupported
	return
}

// Turn On pin
func (d *DirectPinDriver) On() (err error) {
	if writer, ok := d.Connection().(DigitalWriter); ok {
		return writer.DigitalWrite(d.Pin(), byte(1))
	}
	err = ErrDigitalWriteUnsupported
	return
}

// DigitalRead returns the current digital state of the pin
func (d *DirectPinDriver) DigitalRead() (val int, err error) {
	if reader, ok := d.Connection().(DigitalReader); ok {
		return reader.DigitalRead(d.Pin())
	}
	err = ErrDigitalReadUnsupported
	return
}

// DigitalWrite writes to the pin. Acceptable values are 1 or 0
func (d *DirectPinDriver) DigitalWrite(level byte) (err error) {
	if writer, ok := d.Connection().(DigitalWriter); ok {
		return writer.DigitalWrite(d.Pin(), level)
	}
	err = ErrDigitalWriteUnsupported
	return
}

// PwmWrite writes the 0-254 value to the specified pin
func (d *DirectPinDriver) PwmWrite(level byte) (err error) {
	if writer, ok := d.Connection().(PwmWriter); ok {
		return writer.PwmWrite(d.Pin(), level)
	}
	err = ErrPwmWriteUnsupported
	return
}

// ServoWrite writes value to the specified pin
func (d *DirectPinDriver) ServoWrite(level byte) (err error) {
	if writer, ok := d.Connection().(ServoWriter); ok {
		return writer.ServoWrite(d.Pin(), level)
	}
	err = ErrServoWriteUnsupported
	return
}
