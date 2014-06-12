package digispark

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

var d *DigisparkAdaptor

func init() {
	d = NewDigisparkAdaptor("bot")
	d.connect = func(d *DigisparkAdaptor) {}
}

func TestFinalize(t *testing.T) {
	gobot.Expect(t, d.Finalize(), true)
}
func TestConnect(t *testing.T) {
	gobot.Expect(t, d.Connect(), true)
}
func TestDisconnect(t *testing.T) {
	gobot.Expect(t, d.Disconnect(), true)
}
func TestReconnect(t *testing.T) {
	gobot.Expect(t, d.Reconnect(), true)
}
