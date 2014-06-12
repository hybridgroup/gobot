package beaglebone

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

var b *BeagleboneAdaptor

func init() {
	b = NewBeagleboneAdaptor("bot")
}

func TestFinalize(t *testing.T) {
	gobot.Expect(t, b.Finalize(), true)
}

func TestConnect(t *testing.T) {
	gobot.Expect(t, b.Connect(), true)
}
func TestDisconnect(t *testing.T) {
	gobot.Expect(t, b.Disconnect(), true)
}
func TestReconnect(t *testing.T) {
	gobot.Expect(t, b.Reconnect(), true)
}
