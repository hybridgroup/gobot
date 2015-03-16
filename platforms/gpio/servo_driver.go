package gpio

import "github.com/hybridgroup/gobot"

var _ gobot.Driver = (*ServoDriver)(nil)

// ServoDriver Represents a Servo
type ServoDriver struct {
	name       string
	pin        string
	connection ServoWriter
	gobot.Commander
	CurrentAngle byte
}

// NewServoDriver returns a new ServoDriver given a ServoWriter, name and pin.
//
// Adds the following API Commands:
// 	"Move" - See ServoDriver.Move
//	"Min" - See ServoDriver.Min
//	"Center" - See ServoDriver.Center
//	"Max" - See ServoDriver.Max
func NewServoDriver(a ServoWriter, name string, pin string) *ServoDriver {
	s := &ServoDriver{
		name:         name,
		connection:   a,
		pin:          pin,
		Commander:    gobot.NewCommander(),
		CurrentAngle: 0,
	}

	s.AddCommand("Move", func(params map[string]interface{}) interface{} {
		angle := byte(params["angle"].(float64))
		return s.Move(angle)
	})
	s.AddCommand("Min", func(params map[string]interface{}) interface{} {
		return s.Min()
	})
	s.AddCommand("Center", func(params map[string]interface{}) interface{} {
		return s.Center()
	})
	s.AddCommand("Max", func(params map[string]interface{}) interface{} {
		return s.Max()
	})

	return s

}

// Name returns the ServoDrivers name
func (s *ServoDriver) Name() string { return s.name }

// Pin returns the ServoDrivers pin
func (s *ServoDriver) Pin() string { return s.pin }

// Connection returns the ServoDrivers connection
func (s *ServoDriver) Connection() gobot.Connection { return s.connection.(gobot.Connection) }

// Start implements the Driver interface
func (s *ServoDriver) Start() (errs []error) { return }

// Halt implements the Driver interface
func (s *ServoDriver) Halt() (errs []error) { return }

// Move sets the servo to the specified angle. Acceptable angles are 0-180
func (s *ServoDriver) Move(angle uint8) (err error) {
	if !(angle >= 0 && angle <= 180) {
		return ErrServoOutOfRange
	}
	s.CurrentAngle = angle
	return s.connection.ServoWrite(s.Pin(), s.angleToSpan(angle))
}

// Min sets the servo to it's minimum position
func (s *ServoDriver) Min() (err error) {
	return s.Move(0)
}

// Center sets the servo to it's center position
func (s *ServoDriver) Center() (err error) {
	return s.Move(90)
}

// Max sets the servo to its maximum position
func (s *ServoDriver) Max() (err error) {
	return s.Move(180)
}

func (s *ServoDriver) angleToSpan(angle byte) byte {
	return byte(angle * (255 / 180))
}
