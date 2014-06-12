package ardrone

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

var adaptor *ArdroneAdaptor
var d *testDrone

func init() {
	d = &testDrone{}
	adaptor = NewArdroneAdaptor("drone")
	adaptor.connect = func(a *ArdroneAdaptor) {
		a.drone = d
	}
}

func TestFinalize(t *testing.T) {
	gobot.Expect(t, adaptor.Finalize(), true)
}
func TestConnect(t *testing.T) {
	gobot.Expect(t, adaptor.Connect(), true)
}
func TestDisconnect(t *testing.T) {
	gobot.Expect(t, adaptor.Disconnect(), true)
}

func TestReconnect(t *testing.T) {
	gobot.Expect(t, adaptor.Reconnect(), true)
}

func TestDrone(t *testing.T) {
	adaptor.Connect()
	gobot.Expect(t, adaptor.Drone(), d)
}
