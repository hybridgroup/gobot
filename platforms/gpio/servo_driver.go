package gpio

import (
	"github.com/hybridgroup/gobot"
)

type ServoDriver struct {
	gobot.Driver
	Adaptor      Servo
	CurrentAngle byte
}

func NewServoDriver(a Servo, name string, pin string) *ServoDriver {
	s := &ServoDriver{
		Driver: gobot.Driver{
			Name:     name,
			Pin:      pin,
			Commands: make(map[string]func(map[string]interface{}) interface{}),
		},
		CurrentAngle: 0,
		Adaptor:      a,
	}

	s.Driver.AddCommand("Move", func(params map[string]interface{}) interface{} {
		angle := byte(params["angle"].(float64))
		s.Move(angle)
		return nil
	})
	s.Driver.AddCommand("Min", func(params map[string]interface{}) interface{} {
		s.Min()
		return nil
	})
	s.Driver.AddCommand("Center", func(params map[string]interface{}) interface{} {
		s.Center()
		return nil
	})
	s.Driver.AddCommand("Max", func(params map[string]interface{}) interface{} {
		s.Max()
		return nil
	})

	return s

}

func (s *ServoDriver) Start() bool { return true }
func (s *ServoDriver) Halt() bool  { return true }
func (s *ServoDriver) Init() bool  { return true }

func (s *ServoDriver) InitServo() {
	s.Adaptor.InitServo()
}

func (s *ServoDriver) Move(angle uint8) {
	if !(angle >= 0 && angle <= 180) {
		panic("Servo angle must be an integer between 0-180")
	}
	s.CurrentAngle = angle
	s.Adaptor.ServoWrite(s.Pin, s.angleToSpan(angle))
}

func (s *ServoDriver) Min() {
	s.Move(0)
}

func (s *ServoDriver) Center() {
	s.Move(90)
}

func (s *ServoDriver) Max() {
	s.Move(180)
}

func (s *ServoDriver) angleToSpan(angle byte) byte {
	return byte(angle * (255 / 180))
}
