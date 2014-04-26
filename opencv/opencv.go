package gobotOpencv

import (
	"github.com/hybridgroup/gobot"
)

type Opencv struct {
	gobot.Adaptor
}

func (me *Opencv) Connect() bool {
	return true
}

func (me *Opencv) Reconnect() bool {
	return true
}

func (me *Opencv) Disconnect() bool {
	return true
}

func (me *Opencv) Finalize() bool {
	return true
}
