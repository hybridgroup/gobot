package ardrone

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestArdroneAdaptor() *ArdroneAdaptor {
	a := NewArdroneAdaptor("drone")
	a.connect = func(a *ArdroneAdaptor) {
		a.drone = &testDrone{}
	}
	return a
}

func TestConnect(t *testing.T) {
	a := initTestArdroneAdaptor()
	gobot.Assert(t, a.Connect(), true)
}

func TestFinalize(t *testing.T) {
	a := initTestArdroneAdaptor()
	gobot.Assert(t, a.Finalize(), true)
}
