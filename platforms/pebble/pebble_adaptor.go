package pebble

type Adaptor struct {
	name string
}

// NewAdaptor creates a new pebble adaptor
func NewAdaptor() *Adaptor {
	return &Adaptor{
		name: "Pebble",
	}
}

func (a *Adaptor) Name() string     { return a.name }
func (a *Adaptor) SetName(n string) { a.name = n }

// Connect returns true if connection to pebble is established successfully
func (a *Adaptor) Connect() (err error) {
	return
}

// Finalize returns true if connection to pebble is finalized successfully
func (a *Adaptor) Finalize() (err error) {
	return
}
