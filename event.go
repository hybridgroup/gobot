package gobot

// Event represents when something asynchronous happens in a Driver
// or Adaptor
type Event struct {
	Name string
	Data interface{}
}

// NewEvent returns a new Event and its associated data.
func NewEvent(name string, data interface{}) *Event {
	return &Event{Name: name, Data: data}
}
