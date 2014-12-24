package gobot

type eventer struct {
	events map[string]*Event
}

type Eventer interface {
	Events() (events map[string]*Event)
	Event(name string) (event *Event)
	AddEvent(name string)
}

func NewEventer() Eventer {
	return &eventer{
		events: make(map[string]*Event),
	}
}

// Events returns driver events map
func (e *eventer) Events() map[string]*Event {
	return e.events
}

// Event returns an event by name if exists
func (e *eventer) Event(name string) (event *Event) {
	event, _ = e.events[name]
	return
}

// AddEvents adds a new event by name
func (e *eventer) AddEvent(name string) {
	e.events[name] = NewEvent()
}
