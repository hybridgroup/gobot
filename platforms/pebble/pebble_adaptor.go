package pebble

import (
	"github.com/hybridgroup/gobot"
)

type PebbleAdaptor struct {
	gobot.Adaptor
}

func NewPebbleAdaptor(name string) *PebbleAdaptor {
	return &PebbleAdaptor{
		Adaptor: *gobot.NewAdaptor(
			name,
			"PebbleAdaptor",
		),
	}
}

func (a *PebbleAdaptor) Connect() bool {
	return true
}

func (a *PebbleAdaptor) Reconnect() bool {
	return true
}

func (a *PebbleAdaptor) Disconnect() bool {
	return true
}

func (a *PebbleAdaptor) Finalize() bool {
	return true
}
