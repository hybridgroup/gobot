package bebop

import (
	"errors"
	"testing"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/gobottest"
)

var _ gobot.Adaptor = (*Adaptor)(nil)

func initTestBebopAdaptor() *Adaptor {
	a := NewAdaptor()
	a.connect = func(b *Adaptor) (err error) {
		b.drone = &testDrone{}
		return nil
	}
	return a
}

func TestBebopAdaptorConnect(t *testing.T) {
	a := initTestBebopAdaptor()
	gobottest.Assert(t, len(a.Connect()), 0)

	a.connect = func(a *Adaptor) error {
		return errors.New("connection error")
	}
	gobottest.Assert(t, a.Connect()[0], errors.New("connection error"))
}

func TestBebopAdaptorFinalize(t *testing.T) {
	a := initTestBebopAdaptor()
	a.Connect()
	gobottest.Assert(t, len(a.Finalize()), 0)
}
