package bebop

import (
	"errors"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
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
	gobottest.Assert(t, a.Connect(), nil)

	a.connect = func(a *Adaptor) error {
		return errors.New("connection error")
	}
	gobottest.Assert(t, a.Connect(), errors.New("connection error"))
}

func TestBebopAdaptorFinalize(t *testing.T) {
	a := initTestBebopAdaptor()
	a.Connect()
	gobottest.Assert(t, a.Finalize(), nil)
}
