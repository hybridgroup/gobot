package mavlink

import (
	"testing"

	"github.com/hybridgroup/gobot"
)

func initTestMavlinkAdaptor() *MavlinkAdaptor {
	return NewMavlinkAdaptor("myAdaptor", "/dev/null")
}

func TestMavlinkAdaptorConnect(t *testing.T) {
	t.SkipNow()
	a := initTestMavlinkAdaptor()
	gobot.Assert(t, a.Connect(), true)
}

func TestMavlinkAdaptorFinalize(t *testing.T) {
	t.SkipNow()
	a := initTestMavlinkAdaptor()
	gobot.Assert(t, a.Finalize(), true)
}
