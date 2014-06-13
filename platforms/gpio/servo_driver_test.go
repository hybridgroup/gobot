package gpio

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

var s *ServoDriver

func init() {
	s = NewServoDriver(TestAdaptor{}, "bot", "1")
}

func TestServoStart(t *testing.T) {
	gobot.Expect(t, l.Start(), true)
}

func TestServoHalt(t *testing.T) {
	gobot.Expect(t, l.Halt(), true)
}

func TestServoInit(t *testing.T) {
	gobot.Expect(t, l.Init(), true)
}

func TestServoMove(t *testing.T) {
	s.Move(100)
	gobot.Expect(t, s.CurrentAngle, uint8(100))
}

func TestServoMin(t *testing.T) {
	s.Min()
	gobot.Expect(t, s.CurrentAngle, uint8(0))
}

func TestServoMax(t *testing.T) {
	s.Max()
	gobot.Expect(t, s.CurrentAngle, uint8(180))
}

func TestServoCenter(t *testing.T) {
	s.Center()
	gobot.Expect(t, s.CurrentAngle, uint8(90))
}

func TestServoInitServo(t *testing.T) {
	s.InitServo()
}
