package gobot

type eventChannel chan *Event

type eventer struct {
	// map of valid Event names
	eventnames map[string]string

	// new events get put in to the event channel
	in eventChannel

	// slice of out channels used by subscribers
	outs []eventChannel
}

// Eventer is the interface which describes how a Driver or Adaptor
// handles events.
type Eventer interface {
	// Events returns the map of valid Event names.
	Events() (eventnames map[string]string)
	// Event returns the map of valid Event names.
	Event(name string) string
	// AddEvent registers a new Event name.
	AddEvent(name string)
	// Publish new events to anyone listening
	Publish(name string, data interface{})
	// Subscribe to any events from this eventer
	Subscribe() (events eventChannel)
}

// NewEventer returns a new Eventer.
func NewEventer() Eventer {
	evtr := &eventer{
		eventnames: make(map[string]string),
		in:         make(eventChannel, 1),
		outs:       make([]eventChannel, 1),
	}

	// goroutine to cascade in events to all out event channels
	go func() {
		for {
			select {
			case evt := <-evtr.in:
				for _, out := range evtr.outs[1:] {
					out <- evt
				}
			}
		}
	}()

	return evtr
}

func (e *eventer) Events() map[string]string {
	return e.eventnames
}

func (e *eventer) Event(name string) string {
	return e.eventnames[name]
}

func (e *eventer) AddEvent(name string) {
	e.eventnames[name] = name
}

func (e *eventer) Publish(name string, data interface{}) {
	evt := NewEvent(name, data)
	e.in <- evt
}

func (e *eventer) Subscribe() eventChannel {
	out := make(eventChannel)
	e.outs = append(e.outs, out)
	return out
}

// On executes f when e is Published to. Returns ErrUnknownEvent if Event
// does not exist.
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

// Once is similar to On except that it only executes f one time. Returns
//ErrUnknownEvent if Event does not exist.
func (e *eventer) Once(n string, f func(s interface{})) (err error) {
	out := e.Subscribe()
	go func() {
		for {
			select {
			case evt := <-out:
				if evt.Name == n {
					f(evt.Data)
					break
				}
			}
		}
	}()

	return
}
