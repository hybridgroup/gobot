package leap

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

var l *LeapMotionAdaptor

func init() {
	l = NewLeapMotionAdaptor("bot", "/dev/null")
}

func TestFinalize(t *testing.T) {
	t.SkipNow()
	gobot.Expect(t, l.Finalize(), true)
}
func TestConnect(t *testing.T) {
	t.SkipNow()
	gobot.Expect(t, l.Connect(), true)
}
