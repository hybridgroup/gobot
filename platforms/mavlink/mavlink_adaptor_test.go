package mavlink

import (
	"testing"

	"github.com/hybridgroup/gobot"
)

func initTestMavlinkAdaptor() *MavlinkAdaptor {
	m := NewMavlinkAdaptor("myAdaptor", "/dev/null")
	m.sp = gobot.NullReadWriteCloser{}
	m.connect = func(a *MavlinkAdaptor) {}
	return m
}

func TestMavlinkAdaptorConnect(t *testing.T) {
	a := initTestMavlinkAdaptor()
	gobot.Assert(t, a.Connect(), true)
}

func TestMavlinkAdaptorFinalize(t *testing.T) {
	a := initTestMavlinkAdaptor()
	gobot.Assert(t, a.Finalize(), true)
}
