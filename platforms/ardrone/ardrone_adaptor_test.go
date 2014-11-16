package ardrone

import (
	"github.com/hybridgroup/gobot"
	"testing"
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
	gobot.Assert(t, a.Connect(), nil)
}

func TestFinalize(t *testing.T) {
	a := initTestArdroneAdaptor()
	gobot.Assert(t, a.Finalize(), nil)
}
