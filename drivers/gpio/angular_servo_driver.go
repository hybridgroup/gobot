package gpio

import "gobot.io/x/gobot"

// AngularServoDriver Represents a Servo
type AngularServoDriver struct {
	name       string
	pin        string
	connection PwmWriter
	gobot.Commander
	CurrentAngle float64
	MinAngle     float64
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
func NewAngularServoDriver(a PwmWriter, pin string, options ...func(*AngularServoDriver)) *AngularServoDriver {
	s := &AngularServoDriver{
		name:         gobot.DefaultName("AngularServo"),
		connection:   a,
		pin:          pin,
		Commander:    gobot.NewCommander(),
		CurrentAngle: 0.0,
		MinAngle:     0.0,
		MaxAngle:     180.0,
		minPeriod:    1.0 / 1000.0,
		maxPeriod:    2.0 / 1000.0,
	}

	for _, option := range options {
		option(s)
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

// WithServoAngles sets the range of acceptable angles.
// Note that the min can be less than the max to reverse servo operation.
func (s *AngularServoDriver) WithServoAngles(minAngle float64, maxAngle float64) {
	s.MinAngle = minAngle
	s.MaxAngle = maxAngle
}

// WithServoAngles sets the range of acceptable angles.
func WithServoAngles(minAngle float64, maxAngle float64) func(*AngularServoDriver) {
	return func(s *AngularServoDriver) {
		s.WithServoAngles(minAngle, maxAngle)
	}
}

// WithPeriodRange sets the range of acceptable periods.
func (s *AngularServoDriver) WithPeriodRange(minPeriod float64, maxPeriod float64) {
	s.minPeriod = minPeriod
	s.maxPeriod = maxPeriod
}

// WithPeriodRange sets the range of acceptable periods.
func WithPeriodRange(minPeriod float64, maxPeriod float64) func(*AngularServoDriver) {
	return func(s *AngularServoDriver) {
		s.WithPeriodRange(minPeriod, maxPeriod)
	}
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

// Move sets the servo to the specified angle. Acceptable angles are defined at setup time
func (s *AngularServoDriver) Move(angle float64) (err error) {
	if !((angle >= s.MinAngle && angle <= s.MaxAngle) ||
		(angle <= s.MinAngle && angle >= s.MaxAngle)) {
		// We need to support the case were minAngle>maxAngle, to reverse the servo
		return ErrServoOutOfRange
	}
	s.CurrentAngle = angle
	period := s.minPeriod + (s.maxPeriod-s.minPeriod)*(angle-s.MinAngle)/(s.MaxAngle-s.MinAngle)
	return s.connection.SetPwmPeriod(s.Pin(), period)
}

// Min sets the servo to it's minimum position
func (s *AngularServoDriver) Min() (err error) {
	return s.Move(s.MinAngle)
}

// Center sets the servo to it's center position
func (s *AngularServoDriver) Center() (err error) {
	return s.Move((s.MaxAngle - s.MinAngle) / 2.0)
}

// Max sets the servo to its maximum position
func (s *AngularServoDriver) Max() (err error) {
	return s.Move(s.MaxAngle)
}
