package gpio

import (
	"github.com/hybridgroup/gobot"
)

type ServoDriver struct {
	gobot.Driver
	Adaptor      Servo
	CurrentAngle byte
}

func NewServoDriver(a Servo) *ServoDriver {
	return &ServoDriver{
		Driver: gobot.Driver{
			Commands: []string{
				"MoveC",
				"MinC",
				"CenterC",
				"MaxC",
			},
		},
		CurrentAngle: 0,
		Adaptor:      a,
	}
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
