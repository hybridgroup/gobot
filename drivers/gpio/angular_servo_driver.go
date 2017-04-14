package gpio

import "gobot.io/x/gobot"

// AngularServoDriver Represents a Servo
type AngularServoDriver struct {
	name       string
	pin        string
	connection PwmWriter
	gobot.Commander
	CurrentAngle float64
	MaxAngle     float64
	minPeriod    float64
	maxPeriod    float64
}

// NewAngularServoDriver returns a new ServoDriver given a ServoWriter and pin.
//
// Adds the following API Commands:
// 	"Move" - See ServoDriver.Move
//	"Min" - See ServoDriver.Min
//	"Center" - See ServoDriver.Center
//	"Max" - See ServoDriver.Max
func NewAngularServoDriver(a PwmWriter, pin string, maxAngle float64, minPeriod float64, maxPeriod float64) *AngularServoDriver {
	s := &AngularServoDriver{
		name:         gobot.DefaultName("AngularServo"),
		connection:   a,
		pin:          pin,
		Commander:    gobot.NewCommander(),
		CurrentAngle: 0,
		MaxAngle:     maxAngle,
		minPeriod:    minPeriod,
		maxPeriod:    maxPeriod,
	}

	s.AddCommand("Move", func(params map[string]interface{}) interface{} {
		angle := params["angle"].(float64)
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
func (s *AngularServoDriver) Name() string { return s.name }

// SetName sets the ServoDrivers name
func (s *AngularServoDriver) SetName(n string) { s.name = n }

// Pin returns the ServoDrivers pin
func (s *AngularServoDriver) Pin() string { return s.pin }

// Connection returns the ServoDrivers connection
func (s *AngularServoDriver) Connection() gobot.Connection { return s.connection.(gobot.Connection) }

// Start implements the Driver interface
func (s *AngularServoDriver) Start() (err error) { return }

// Halt implements the Driver interface
func (s *AngularServoDriver) Halt() (err error) { return }

// Move sets the servo to the specified angle. Acceptable angles are 0-180
func (s *AngularServoDriver) Move(angle float64) (err error) {
	if !(angle >= 0 && angle <= s.MaxAngle) {
		return ErrServoOutOfRange
	}
	s.CurrentAngle = angle
	period := s.minPeriod + (s.maxPeriod-s.minPeriod)*angle/s.MaxAngle
	return s.connection.SetPwmPeriod(s.Pin(), period)
}

// Min sets the servo to it's minimum position
func (s *AngularServoDriver) Min() (err error) {
	return s.Move(0)
}

// Center sets the servo to it's center position
func (s *AngularServoDriver) Center() (err error) {
	return s.Move(s.MaxAngle / 2.0)
}

// Max sets the servo to its maximum position
func (s *AngularServoDriver) Max() (err error) {
	return s.Move(s.MaxAngle)
}
