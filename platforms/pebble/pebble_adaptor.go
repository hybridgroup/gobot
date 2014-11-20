package pebble

import (
	"github.com/hybridgroup/gobot"
)

var _ gobot.AdaptorInterface = (*PebbleAdaptor)(nil)

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
func (a *PebbleAdaptor) Connect() (errs []error) {
	return
}

// Finalize returns true if connection to pebble is finalized succesfully
func (a *PebbleAdaptor) Finalize() (errs []error) {
	return
}
