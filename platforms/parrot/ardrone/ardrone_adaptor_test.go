package ardrone

import (
	"errors"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
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
	gobottest.Assert(t, a.Connect(), nil)

	a.connect = func(a *Adaptor) (drone, error) {
		return nil, errors.New("connection error")
	}
	gobottest.Assert(t, a.Connect(), errors.New("connection error"))
}

func TestArdroneAdaptorFinalize(t *testing.T) {
	a := initTestArdroneAdaptor()
	gobottest.Assert(t, a.Finalize(), nil)
}
