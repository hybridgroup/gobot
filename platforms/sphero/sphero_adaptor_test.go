package sphero

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

var a *SpheroAdaptor

func init() {
	a = NewSpheroAdaptor("bot", "/dev/null")
	a.sp = sp{}
	a.connect = func(a *SpheroAdaptor) {}
}

func TestFinalize(t *testing.T) {
	gobot.Expect(t, a.Finalize(), true)
}
func TestConnect(t *testing.T) {
	gobot.Expect(t, a.Connect(), true)
}
