package sphero

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

var s *SpheroDriver

func init() {
	a := NewSpheroAdaptor("bot", "/dev/null")
	a.sp = sp{}
	s = NewSpheroDriver(a, "bot")
}

func TestStart(t *testing.T) {
	gobot.Expect(t, s.Start(), true)
}

func TestHalt(t *testing.T) {
	gobot.Expect(t, s.Halt(), true)
}
