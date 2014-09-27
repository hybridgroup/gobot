package gpio

import (
	"github.com/hybridgroup/gobot"
)

// Represents a Servo
type ServoDriver struct {
	gobot.Driver
	CurrentAngle byte
}

// NewSerovDriver return a new ServoDriver  given a Servo, name and pin.
//
// Adds the following API Commands:
// 	"Move" - See ServoDriver.Move
//	"Min" - See ServoDriver.Min
//	"Center" - See ServoDriver.Center
//	"Max" - See ServoDriver.Max
func NewServoDriver(a Servo, name string, pin string) *ServoDriver {
	s := &ServoDriver{
		Driver: *gobot.NewDriver(
			name,
			"ServoDriver",
			a.(gobot.AdaptorInterface),
			pin,
		),
		CurrentAngle: 0,
	}

	s.AddCommand("Move", func(params map[string]interface{}) interface{} {
		angle := byte(params["angle"].(float64))
		s.Move(angle)
		return nil
	})
	s.AddCommand("Min", func(params map[string]interface{}) interface{} {
		s.Min()
		return nil
	})
	s.AddCommand("Center", func(params map[string]interface{}) interface{} {
		s.Center()
		return nil
	})
	s.AddCommand("Max", func(params map[string]interface{}) interface{} {
		s.Max()
		return nil
	})

	return s

}

func (s *ServoDriver) adaptor() Servo {
	return s.Adaptor().(Servo)
}

// Start starts the ServoDriver. Returns true on successful start of the driver.
func (s *ServoDriver) Start() bool { return true }

// Halt halts the ServoDriver. Returns true on successful halt of the driver.
func (s *ServoDriver) Halt() bool { return true }

// InitServo initializes the ServoDriver on platforms which require an explicit initialization.
func (s *ServoDriver) InitServo() {
	s.adaptor().InitServo()
}

// Move sets the servo to the specified angle
func (s *ServoDriver) Move(angle uint8) {
	if !(angle >= 0 && angle <= 180) {
		panic("Servo angle must be an integer between 0-180")
	}
	s.CurrentAngle = angle
	s.adaptor().ServoWrite(s.Pin(), s.angleToSpan(angle))
}

// Min sets the servo to it's minimum position
func (s *ServoDriver) Min() {
	s.Move(0)
}

// Center sets the servo to it's center position
func (s *ServoDriver) Center() {
	s.Move(90)
}

// Max sets the servo to its maximum position
func (s *ServoDriver) Max() {
	s.Move(180)
}

func (s *ServoDriver) angleToSpan(angle byte) byte {
	return byte(angle * (255 / 180))
}
