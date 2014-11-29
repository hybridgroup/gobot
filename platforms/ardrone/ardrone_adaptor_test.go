package ardrone

import (
	"testing"

	"github.com/hybridgroup/gobot"
)

func initTestArdroneAdaptor() *ArdroneAdaptor {
	a := NewArdroneAdaptor("drone")
	a.connect = func(a *ArdroneAdaptor) error {
		a.drone = &testDrone{}
		return nil
	}
	return a
}

func TestConnect(t *testing.T) {
	a := initTestArdroneAdaptor()
	gobot.Assert(t, len(a.Connect()), 0)
}

func TestFinalize(t *testing.T) {
	a := initTestArdroneAdaptor()
	gobot.Assert(t, len(a.Finalize()), 0)
}
