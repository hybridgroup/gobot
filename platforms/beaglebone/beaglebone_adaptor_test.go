package beaglebone

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestBeagleboneAdaptor() *BeagleboneAdaptor {
	return NewBeagleboneAdaptor("bot")
}

func TestBeagleboneAdaptorFinalize(t *testing.T) {
	a := initTestBeagleboneAdaptor()
	gobot.Expect(t, a.Finalize(), true)
}
func TestBeagleboneAdaptorConnect(t *testing.T) {
	a := initTestBeagleboneAdaptor()
	gobot.Expect(t, a.Connect(), true)
}
func TestBeagleboneAdaptorDisconnect(t *testing.T) {
	a := initTestBeagleboneAdaptor()
	gobot.Expect(t, a.Disconnect(), true)
}
func TestBeagleboneAdaptorReconnect(t *testing.T) {
	a := initTestBeagleboneAdaptor()
	gobot.Expect(t, a.Reconnect(), true)
}
