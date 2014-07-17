package beaglebone

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestBeagleboneAdaptor() *BeagleboneAdaptor {
	b := NewBeagleboneAdaptor("bot")
	b.connect = func() {}
	return b
}

func TestBeagleboneAdaptorFinalize(t *testing.T) {
	a := initTestBeagleboneAdaptor()
	gobot.Assert(t, a.Finalize(), true)
}
func TestBeagleboneAdaptorConnect(t *testing.T) {
	a := initTestBeagleboneAdaptor()
	gobot.Assert(t, a.Connect(), true)
}
func TestBeagleboneAdaptorDisconnect(t *testing.T) {
	a := initTestBeagleboneAdaptor()
	gobot.Assert(t, a.Disconnect(), true)
}
func TestBeagleboneAdaptorReconnect(t *testing.T) {
	a := initTestBeagleboneAdaptor()
	gobot.Assert(t, a.Reconnect(), true)
}
