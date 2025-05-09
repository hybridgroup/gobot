package gpio

import (
	"fmt"

	"gobot.io/x/gobot/v2"
)

// ServoDriver Represents a Servo
type ServoDriver struct {
	*driver
	currentAngle byte
}

// NewServoDriver returns a new ServoDriver given a ServoWriter and pin.
//
// Supported options:
//
//	"WithName"
//
// Adds the following API Commands:
//
//	"Move" - See ServoDriver.Move
//	"Min" - See ServoDriver.ToMin
//	"Center" - See ServoDriver.ToCenter
//	"Max" - See ServoDriver.ToMax
func NewServoDriver(a ServoWriter, pin string, opts ...interface{}) *ServoDriver {
	//nolint:forcetypeassert // no error return value, so there is no better way
	d := &ServoDriver{
		driver: newDriver(a.(gobot.Connection), "Servo", append(opts, withPin(pin))...),
	}

	//nolint:forcetypeassert // ok here
	d.AddCommand("Move", func(params map[string]interface{}) interface{} {
		angle := byte(params["angle"].(float64))
		return d.Move(angle)
	})
	d.AddCommand("ToMin", func(_ map[string]interface{}) interface{} {
		return d.ToMin()
	})
	d.AddCommand("ToCenter", func(_ map[string]interface{}) interface{} {
		return d.ToCenter()
	})
	d.AddCommand("ToMax", func(_ map[string]interface{}) interface{} {
		return d.ToMax()
	})

	return d
}

// Move sets the servo to the specified angle. Acceptable angles are 0-180
func (d *ServoDriver) Move(angle uint8) error {
	if angle > 180 {
		return fmt.Errorf("servo angle (%d) must be between 0-180", angle)
	}
	d.currentAngle = angle
	return d.servoWrite(d.driverCfg.pin, angle)
}

// Min sets the servo to it's minimum position
func (d *ServoDriver) ToMin() error {
	return d.Move(0)
}

// Center sets the servo to it's center position
func (d *ServoDriver) ToCenter() error {
	return d.Move(90)
}

// Max sets the servo to its maximum position
func (d *ServoDriver) ToMax() error {
	return d.Move(180)
}

// Angle returns the current angle
func (d *ServoDriver) Angle() uint8 {
	return d.currentAngle
}
