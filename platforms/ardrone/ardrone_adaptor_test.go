package ardrone

import (
	"errors"
	"testing"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/gobottest"
)

var _ gobot.Adaptor = (*ArdroneAdaptor)(nil)

func initTestArdroneAdaptor() *ArdroneAdaptor {
	a := NewArdroneAdaptor("drone")
	a.connect = func(a *ArdroneAdaptor) (drone, error) {
		return &testDrone{}, nil
	}
	return a
}

func TestArdroneAdaptor(t *testing.T) {
	a := NewArdroneAdaptor("drone")
	gobottest.Assert(t, a.Name(), "drone")
	gobottest.Assert(t, a.config.Ip, "192.168.1.1")

	a = NewArdroneAdaptor("drone", "192.168.100.100")
	gobottest.Assert(t, a.config.Ip, "192.168.100.100")
}

func TestArdroneAdaptorConnect(t *testing.T) {
	a := initTestArdroneAdaptor()
	gobottest.Assert(t, len(a.Connect()), 0)

	a.connect = func(a *ArdroneAdaptor) (drone, error) {
		return nil, errors.New("connection error")
	}
	gobottest.Assert(t, a.Connect()[0], errors.New("connection error"))
}

func TestArdroneAdaptorFinalize(t *testing.T) {
	a := initTestArdroneAdaptor()
	gobottest.Assert(t, len(a.Finalize()), 0)
}
