package gobot

import "sync"

type eventChannel chan *Event

type eventer struct {
	// map of valid Event names
	eventnames map[string]string

	// new events get put in to the event channel
	in eventChannel

	// map of out channels used by subscribers
	outs map[eventChannel]eventChannel

	// mutex to protect the eventChannel map
	eventsMutex sync.Mutex
}

const eventChanBufferSize = 10

// Eventer is the interface which describes how a Driver or Adaptor
// handles events.
type Eventer interface {
	// Events returns the map of valid Event names.
	Events() (eventnames map[string]string)

	// Event returns an Event string from map of valid Event names.
	// Mostly used to validate that an Event name is valid.
	Event(name string) string

	// AddEvent registers a new Event name.
	AddEvent(name string)

	// DeleteEvent removes a previously registered Event name.
	DeleteEvent(name string)

	// Publish new events to any subscriber
	Publish(name string, data interface{})

	// Subscribe to events
	Subscribe() (events eventChannel)

	// Unsubscribe from an event channel
	Unsubscribe(events eventChannel)

	// Event handler
	On(name string, f func(s interface{})) (err error)

	// Event handler, only executes one time
	Once(name string, f func(s interface{})) (err error)
}

// NewEventer returns a new Eventer.
func NewEventer() Eventer {
	evtr := &eventer{
		eventnames: make(map[string]string),
		in:         make(eventChannel, eventChanBufferSize),
		outs:       make(map[eventChannel]eventChannel),
	}

	// goroutine to cascade "in" events to all "out" event channels
	go func() {
		for {
			select {
			case evt := <-evtr.in:
				evtr.eventsMutex.Lock()
				for _, out := range evtr.outs {
					out <- evt
				}
				evtr.eventsMutex.Unlock()
			}
		}
	}()

	return evtr
}

// Events returns the map of valid Event names.
func (e *eventer) Events() map[string]string {
	return e.eventnames
}

// Event returns an Event string from map of valid Event names.
// Mostly used to validate that an Event name is valid.
func (e *eventer) Event(name string) string {
	return e.eventnames[name]
}

// AddEvent registers a new Event name.
func (e *eventer) AddEvent(name string) {
	e.eventnames[name] = name
}

// DeleteEvent removes a previously registered Event name.
func (e *eventer) DeleteEvent(name string) {
	delete(e.eventnames, name)
}

// Publish new events to anyone that is subscribed
func (e *eventer) Publish(name string, data interface{}) {
	evt := NewEvent(name, data)
	e.in <- evt
}

// Subscribe to any events from this eventer
func (e *eventer) Subscribe() eventChannel {
	e.eventsMutex.Lock()
	defer e.eventsMutex.Unlock()
	out := make(eventChannel, eventChanBufferSize)
	e.outs[out] = out
	return out
}

// Unsubscribe from the event channel
func (e *eventer) Unsubscribe(events eventChannel) {
	e.eventsMutex.Lock()
	defer e.eventsMutex.Unlock()
	delete(e.outs, events)
}

// On executes the event handler f when e is Published to.
func (e *eventer) On(n string, f func(s interface{})) (err error) {
	out := e.Subscribe()
	go func() {
		for {
			select {
			case evt := <-out:
				if evt.Name == n {
					f(evt.Data)
				}
			}
		}
	}()

	return
}

// Once is similar to On except that it only executes f one time.
func (e *eventer) Once(n string, f func(s interface{})) (err error) {
	out := e.Subscribe()
	go func() {
	ProcessEvents:
		for evt := range out {
			if evt.Name == n {
				f(evt.Data)
				e.Unsubscribe(out)
				break ProcessEvents
			}
		}
	}()

	return
}
