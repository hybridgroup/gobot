package pebble

import (
	"github.com/hybridgroup/gobot"
)

type PebbleAdaptor struct {
	gobot.Adaptor
}

func (me *PebbleAdaptor) Connect() bool {
  return true
}

func (me *PebbleAdaptor) Reconnect() bool {
  return true
}

func (me *PebbleAdaptor) Disconnect() bool {
  return true
}

func (me *PebbleAdaptor) Finalize() bool {
  return true
}
