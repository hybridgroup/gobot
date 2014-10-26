package arietta

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestAriettaAdaptor() *AriettaAdaptor {
	b := NewAriettaAdaptor("bot")
	return b
}

func TestAriettaAdaptorFinalize(t *testing.T) {
	a := initTestAriettaAdaptor()
	gobot.Assert(t, a.Finalize(), true)
}
func TestAriettaAdaptorConnect(t *testing.T) {
	a := initTestAriettaAdaptor()
	gobot.Assert(t, a.Connect(), true)
}
func TestAriettaAdaptorDisconnect(t *testing.T) {
	a := initTestAriettaAdaptor()
	gobot.Assert(t, a.Disconnect(), true)
}
func TestAriettaAdaptorReconnect(t *testing.T) {
	a := initTestAriettaAdaptor()
	gobot.Assert(t, a.Reconnect(), true)
}
