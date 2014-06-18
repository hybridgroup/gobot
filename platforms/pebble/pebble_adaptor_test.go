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
	gobot.Expect(t, a.Connect(), true)
}

func TestPebbleAdaptorFinalize(t *testing.T) {
	a := initTestPebbleAdaptor()
	gobot.Expect(t, a.Finalize(), true)
}
