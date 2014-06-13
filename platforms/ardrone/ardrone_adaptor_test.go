package ardrone

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestArdroneAdaptor() (*ArdroneAdaptor, *testDrone) {
	d := &testDrone{}
	a := NewArdroneAdaptor("drone")
	a.connect = func(a *ArdroneAdaptor) {
		a.drone = d
	}
	return a, d
}

func TestFinalize(t *testing.T) {
	a, _ := initTestArdroneAdaptor()
	gobot.Expect(t, a.Finalize(), true)
}
func TestConnect(t *testing.T) {
	a, _ := initTestArdroneAdaptor()
	gobot.Expect(t, a.Connect(), true)
}
func TestDisconnect(t *testing.T) {
	a, _ := initTestArdroneAdaptor()
	gobot.Expect(t, a.Disconnect(), true)
}

func TestReconnect(t *testing.T) {
	a, _ := initTestArdroneAdaptor()
	gobot.Expect(t, a.Reconnect(), true)
}

func TestDrone(t *testing.T) {
	a, d := initTestArdroneAdaptor()
	a.Connect()
	gobot.Expect(t, a.Drone(), d)
}
