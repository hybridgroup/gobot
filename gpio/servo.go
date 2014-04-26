package gobotGPIO

import (
	"github.com/hybridgroup/gobot"
)

type ServoInterface interface {
	InitServo()
	ServoWrite(string, byte)
}

type Servo struct {
	gobot.Driver
	Adaptor      ServoInterface
	CurrentAngle byte
}

func NewServo(a ServoInterface) *Servo {
	s := new(Servo)
	s.CurrentAngle = 0
	s.Adaptor = a
	s.Commands = []string{
		"MoveC",
		"MinC",
		"CenterC",
		"MaxC",
	}
	return s
}

func (s *Servo) Start() bool { return true }
func (s *Servo) Halt() bool  { return true }
func (s *Servo) Init() bool  { return true }

func (s *Servo) InitServo() {
	s.Adaptor.InitServo()
}

func (s *Servo) Move(angle uint8) {
	if !(angle >= 0 && angle <= 180) {
		panic("Servo angle must be an integer between 0-180")
	}
	s.CurrentAngle = angle
	s.Adaptor.ServoWrite(s.Pin, s.angleToSpan(angle))
}

func (s *Servo) Min() {
	s.Move(0)
}

func (s *Servo) Center() {
	s.Move(90)
}

func (s *Servo) Max() {
	s.Move(180)
}

func (s *Servo) angleToSpan(angle byte) byte {
	return byte(angle * (255 / 180))
}
