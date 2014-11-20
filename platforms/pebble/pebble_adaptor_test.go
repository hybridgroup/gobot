package pebble

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func initTestPebbleAdaptor() *PebbleAdaptor {
	return NewPebbleAdaptor("pebble")
}

func TestPebbleAdaptorConnect(t *testing.T) {
	a := initTestPebbleAdaptor()
	gobot.Assert(t, len(a.Connect()), 0)
}

func TestPebbleAdaptorFinalize(t *testing.T) {
	a := initTestPebbleAdaptor()
	gobot.Assert(t, len(a.Finalize()), 0)
}
