package pebble

import (
	"github.com/hybridgroup/gobot"
)

type PebbleAdaptor struct {
	gobot.Adaptor
}

// NewPebbleAdaptor creates a new pebble adaptor with specified name
func NewPebbleAdaptor(name string) *PebbleAdaptor {
	return &PebbleAdaptor{
		Adaptor: *gobot.NewAdaptor(
			name,
			"PebbleAdaptor",
		),
	}
}

// Connect returns true if connection to pebble is established succesfully
func (a *PebbleAdaptor) Connect() bool {
	return true
}

// Finalize returns true if connection to pebble is finalized succesfully
func (a *PebbleAdaptor) Finalize() bool {
	return true
}
