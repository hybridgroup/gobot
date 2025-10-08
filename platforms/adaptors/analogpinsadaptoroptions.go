package adaptors

// AnalogPinsOptionApplier is the interface for analog pins adaptor options. This provides the possibility for change
// the platform behavior by the user when creating the platform, e.g. by "NewAdaptor()".
// The interface needs to be implemented by each configurable option type.
type AnalogPinsOptionApplier interface {
	apply(cfg *analogPinsConfiguration)
}

// analogPinsDebugOption is the type to switch on analog pins related debug messages.
type analogPinsDebugOption bool

func (o analogPinsDebugOption) String() string {
	return "switch on debugging for analog pins option"
}

func (o analogPinsDebugOption) apply(cfg *analogPinsConfiguration) {
	cfg.debug = bool(o)
}
