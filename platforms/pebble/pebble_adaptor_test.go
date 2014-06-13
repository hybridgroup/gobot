package pebble

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

var adaptor *PebbleAdaptor

func init() {
	adaptor = NewPebbleAdaptor("pebble")
}

func TestFinalize(t *testing.T) {
	gobot.Expect(t, adaptor.Finalize(), true)
}
func TestConnect(t *testing.T) {
	gobot.Expect(t, adaptor.Connect(), true)
}
