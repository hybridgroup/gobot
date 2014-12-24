package pebble

import (
	"github.com/hybridgroup/gobot"
)

var _ gobot.Adaptor = (*PebbleAdaptor)(nil)

type PebbleAdaptor struct {
	name string
}

// NewPebbleAdaptor creates a new pebble adaptor with specified name
func NewPebbleAdaptor(name string) *PebbleAdaptor {
	return &PebbleAdaptor{
		name: name,
	}
}

func (a *PebbleAdaptor) Name() string { return a.name }

// Connect returns true if connection to pebble is established succesfully
func (a *PebbleAdaptor) Connect() (errs []error) {
	return
}

// Finalize returns true if connection to pebble is finalized succesfully
func (a *PebbleAdaptor) Finalize() (errs []error) {
	return
}
