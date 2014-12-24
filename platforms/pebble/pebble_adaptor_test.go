package pebble

import (
	"testing"

	"github.com/hybridgroup/gobot"
)

func initTestPebbleAdaptor() *PebbleAdaptor {
	return NewPebbleAdaptor("pebble")
}

func TestPebbleAdaptor(t *testing.T) {
	a := initTestPebbleAdaptor()
	gobot.Assert(t, a.Name(), "pebble")
}
func TestPebbleAdaptorConnect(t *testing.T) {
	a := initTestPebbleAdaptor()
	gobot.Assert(t, len(a.Connect()), 0)
}

func TestPebbleAdaptorFinalize(t *testing.T) {
	a := initTestPebbleAdaptor()
	gobot.Assert(t, len(a.Finalize()), 0)
}
