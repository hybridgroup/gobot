package adaptors

// OneWireBusOptionApplier is the interface for 1-wire bus adaptors options. This provides the possibility for change
// the platform behavior by the user when creating the platform, e.g. by "NewAdaptor()".
// The interface needs to be implemented by each configurable option type.
type OneWireBusOptionApplier interface {
	apply(cfg *oneWireBusConfiguration)
}

// oneWireBusDebugOption is the type to switch on 1-wire related debug messages.
type oneWireBusDebugOption bool

func (o oneWireBusDebugOption) String() string {
	return "switch on debugging for 1-wire option"
}

func (o oneWireBusDebugOption) apply(cfg *oneWireBusConfiguration) {
	cfg.debug = bool(o)
}
