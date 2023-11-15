package gpio

import (
	"log"

	"gobot.io/x/gobot/v2"
)

// ServoDriver Represents a Servo
type ServoDriver struct {
	name       string
	pin        string
	connection ServoWriter
	gobot.Commander
	CurrentAngle byte
}

// NewServoDriver returns a new ServoDriver given a ServoWriter and pin.
//
// Adds the following API Commands:
//
//	"Move" - See ServoDriver.Move
//		"Min" - See ServoDriver.Min
//		"Center" - See ServoDriver.Center
//		"Max" - See ServoDriver.Max
func NewServoDriver(a ServoWriter, pin string) *ServoDriver {
	d := &ServoDriver{
		name:         gobot.DefaultName("Servo"),
		connection:   a,
		pin:          pin,
		Commander:    gobot.NewCommander(),
		CurrentAngle: 0,
	}

	d.AddCommand("Move", func(params map[string]interface{}) interface{} {
		angle := byte(params["angle"].(float64)) //nolint:forcetypeassert // ok here
		return d.Move(angle)
	})
	d.AddCommand("Min", func(params map[string]interface{}) interface{} {
		return d.Min()
	})
	d.AddCommand("Center", func(params map[string]interface{}) interface{} {
		return d.Center()
	})
	d.AddCommand("Max", func(params map[string]interface{}) interface{} {
		return d.Max()
	})

	return d
}

// Name returns the ServoDrivers name
func (d *ServoDriver) Name() string { return d.name }

// SetName sets the ServoDrivers name
func (d *ServoDriver) SetName(n string) { d.name = n }

// Pin returns the ServoDrivers pin
func (d *ServoDriver) Pin() string { return d.pin }

// Connection returns the ServoDrivers connection
func (d *ServoDriver) Connection() gobot.Connection {
	if conn, ok := d.connection.(gobot.Connection); ok {
		return conn
	}

	log.Printf("%s has no gobot connection\n", d.name)
	return nil
}

// Start implements the Driver interface
func (d *ServoDriver) Start() error { return nil }

// Halt implements the Driver interface
func (d *ServoDriver) Halt() error { return nil }

// Move sets the servo to the specified angle. Acceptable angles are 0-180
func (d *ServoDriver) Move(angle uint8) error {
	if angle > 180 {
		return ErrServoOutOfRange
	}
	d.CurrentAngle = angle
	return d.connection.ServoWrite(d.Pin(), angle)
}

// Min sets the servo to it's minimum position
func (d *ServoDriver) Min() error {
	return d.Move(0)
}

// Center sets the servo to it's center position
func (d *ServoDriver) Center() error {
	return d.Move(90)
}

// Max sets the servo to its maximum position
func (d *ServoDriver) Max() error {
	return d.Move(180)
}
