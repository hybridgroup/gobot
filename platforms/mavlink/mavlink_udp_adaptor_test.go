package mavlink

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Adaptor = (*UDPAdaptor)(nil)

func initTestMavlinkUDPAdaptor() *UDPAdaptor {
	m := NewUDPAdaptor(":14550")
	return m
}

func TestMavlinkUDPAdaptor(t *testing.T) {
	a := initTestMavlinkUDPAdaptor()
	gobottest.Assert(t, a.Port(), ":14550")
}

func TestMavlinkUDPAdaptorName(t *testing.T) {
	a := initTestMavlinkUDPAdaptor()
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "Mavlink"), true)
	a.SetName("NewName")
	gobottest.Assert(t, a.Name(), "NewName")
}

func TestMavlinkUDPAdaptorConnect(t *testing.T) {
	a := initTestMavlinkUDPAdaptor()
	gobottest.Assert(t, a.Connect(), nil)
}

func TestMavlinkUDPAdaptorFinalize(t *testing.T) {
	a := initTestMavlinkUDPAdaptor()
	gobottest.Assert(t, a.Finalize(), nil)
}
