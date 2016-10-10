package ardrone

import (
	"errors"
	"testing"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/gobottest"
)

var _ gobot.Adaptor = (*Adaptor)(nil)

func initTestArdroneAdaptor() *Adaptor {
	a := NewAdaptor()
	a.connect = func(a *Adaptor) (drone, error) {
		return &testDrone{}, nil
	}
	return a
}

func TestArdroneAdaptor(t *testing.T) {
	a := NewAdaptor()
	gobottest.Assert(t, a.config.Ip, "192.168.1.1")

	a = NewAdaptor("192.168.100.100")
	gobottest.Assert(t, a.config.Ip, "192.168.100.100")
}

func TestArdroneAdaptorConnect(t *testing.T) {
	a := initTestArdroneAdaptor()
	gobottest.Assert(t, len(a.Connect()), 0)

	a.connect = func(a *Adaptor) (drone, error) {
		return nil, errors.New("connection error")
	}
	gobottest.Assert(t, a.Connect()[0], errors.New("connection error"))
}

func TestArdroneAdaptorFinalize(t *testing.T) {
	a := initTestArdroneAdaptor()
	gobottest.Assert(t, len(a.Finalize()), 0)
}
