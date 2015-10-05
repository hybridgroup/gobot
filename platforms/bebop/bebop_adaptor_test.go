package bebop

import (
	"errors"
	"testing"

	"github.com/hybridgroup/gobot"
)

func initTestBebopAdaptor() *BebopAdaptor {
	a := NewBebopAdaptor("bot")
	a.connect = func(b *BebopAdaptor) (err error) {
		b.drone = &testDrone{}
		return nil
	}
	return a
}

func TestBebopAdaptorConnect(t *testing.T) {
	a := initTestBebopAdaptor()
	gobot.Assert(t, len(a.Connect()), 0)

	a.connect = func(a *BebopAdaptor) (error) {
		return errors.New("connection error")
	}
	gobot.Assert(t, a.Connect()[0], errors.New("connection error"))
}

func TestBebopAdaptorFinalize(t *testing.T) {
	a := initTestBebopAdaptor()
	a.Connect()
	gobot.Assert(t, len(a.Finalize()), 0)
}
