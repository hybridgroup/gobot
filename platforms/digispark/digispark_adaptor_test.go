package digispark

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestDigisparkAdaptor() *DigisparkAdaptor {
	a := NewDigisparkAdaptor("bot")
	a.connect = func(a *DigisparkAdaptor) {}
	return a
}

func TestDigisparkAdaptorFinalize(t *testing.T) {
	a := initTestDigisparkAdaptor()
	gobot.Assert(t, a.Finalize(), true)
}

func TestDigisparkAdaptorConnect(t *testing.T) {
	a := initTestDigisparkAdaptor()
	gobot.Assert(t, a.Connect(), true)
}

func TestDigisparkAdaptorDisconnect(t *testing.T) {
	a := initTestDigisparkAdaptor()
	gobot.Assert(t, a.Disconnect(), true)
}

func TestDigisparkAdaptorReconnect(t *testing.T) {
	a := initTestDigisparkAdaptor()
	gobot.Assert(t, a.Reconnect(), true)
}
