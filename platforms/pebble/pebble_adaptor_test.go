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
	gobot.Assert(t, a.Connect(), nil)
}

func TestPebbleAdaptorFinalize(t *testing.T) {
	a := initTestPebbleAdaptor()
	gobot.Assert(t, a.Finalize(), nil)
}
