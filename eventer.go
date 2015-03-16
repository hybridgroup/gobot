package gobot

type eventer struct {
	events map[string]*Event
}

// Eventer is the interface which describes behaviour for a Driver or Adaptor
// which uses events.
type Eventer interface {
	// Events returns the Event map.
	Events() (events map[string]*Event)
	// Event returns an Event by name. Returns nil if the Event is not found.
	Event(name string) (event *Event)
	// AddEvent adds a new Event given a name.
	AddEvent(name string)
}

// NewEventer returns a new Eventer.
func NewEventer() Eventer {
	return &eventer{
		events: make(map[string]*Event),
	}
}

func (e *eventer) Events() map[string]*Event {
	return e.events
}

func (e *eventer) Event(name string) (event *Event) {
	event, _ = e.events[name]
	return
}

func (e *eventer) AddEvent(name string) {
	e.events[name] = NewEvent()
}
