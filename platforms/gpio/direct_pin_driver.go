package gpio

import (
	"strconv"

	"github.com/hybridgroup/gobot"
)

var _ gobot.Driver = (*DirectPinDriver)(nil)

// DirectPinDriver represents a GPIO pin
type DirectPinDriver struct {
	name       string
	pin        string
	connection gobot.Connection
	gobot.Commander
}

// NewDirectPinDriver return a new DirectPinDriver given a Connection, name and pin.
//
// Adds the following API Commands:
// 	"DigitalRead" - See DirectPinDriver.DigitalRead
// 	"DigitalWrite" - See DirectPinDriver.DigitalWrite
// 	"AnalogRead" - See DirectPinDriver.AnalogRead
// 	"AnalogWrite" - See DirectPinDriver.AnalogWrite
// 	"PwmWrite" - See DirectPinDriver.PwmWrite
// 	"ServoWrite" - See DirectPinDriver.ServoWrite
func NewDirectPinDriver(a gobot.Connection, name string, pin string) *DirectPinDriver {
	d := &DirectPinDriver{
		name:       name,
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
	d.AddCommand("AnalogRead", func(params map[string]interface{}) interface{} {
		val, err := d.AnalogRead()
		return map[string]interface{}{"val": val, "err": err}
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

// Pin returns the DirectPinDrivers pin
func (d *DirectPinDriver) Pin() string { return d.pin }

// Connection returns the DirectPinDrivers Connection
func (d *DirectPinDriver) Connection() gobot.Connection { return d.connection }

// Start implements the Driver interface
func (d *DirectPinDriver) Start() (errs []error) { return }

// Halt implements the Driver interface
func (d *DirectPinDriver) Halt() (errs []error) { return }

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

// AnalogRead reads the current analog reading of the pin
func (d *DirectPinDriver) AnalogRead() (val int, err error) {
	if reader, ok := d.Connection().(AnalogReader); ok {
		return reader.AnalogRead(d.Pin())
	}
	err = ErrAnalogReadUnsupported
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
