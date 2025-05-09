package chip

import "gobot.io/x/gobot/v2"

// NewProAdaptor creates a C.H.I.P. Pro Adaptor
func NewProAdaptor() *Adaptor {
	a := NewAdaptor()
	a.name = gobot.DefaultName("CHIP Pro")
	a.pinMap = chipProPins
	return a
}
