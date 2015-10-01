package gobot

import "sync"

type callback struct {
	f    func(interface{})
	once bool
}

// Event executes the list of Callbacks when Chan is written to.
type Event struct {
	sync.Mutex
	Callbacks []callback
}

// NewEvent returns a new Event which is now listening for data.
func NewEvent() *Event {
	return &Event{}
}

// Write writes data to the Event, it will not block and will not buffer if there
// are no active subscribers to the Event.
func (e *Event) Write(data interface{}) {
	e.Lock()
	defer e.Unlock()

	tmp := []callback{}
	for _, cb := range e.Callbacks {
		go cb.f(data)
		if !cb.once {
			tmp = append(tmp, cb)
		}
	}
	e.Callbacks = tmp
}
