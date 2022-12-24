package adaptors

// Optioner is the interface for adaptors options. This provides the possibility for change the platform behavior
// by the user when creating the platform, e.g. by "NewAdaptor()".
type Optioner interface {
	digitalPinsOptioner
}
