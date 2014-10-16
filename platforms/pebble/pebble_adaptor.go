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

// Reconnect retries connection to pebble. Returns true if succesfull
func (a *PebbleAdaptor) Reconnect() bool {
	return true
}

// Disconnect returns true if connection to pebble is closed succesfully
func (a *PebbleAdaptor) Disconnect() bool {
	return true
}

// Finalize returns true if connection to pebble is finalized succesfully
func (a *PebbleAdaptor) Finalize() bool {
	return true
}
