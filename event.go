package gobot

type callback struct {
	f    func(interface{})
	once bool
}

// Event executes the list of Callbacks when Chan is written to.
type Event struct {
	Chan      chan interface{}
	Callbacks []callback
}

// NewEvent returns a new Event which is now listening for data.
func NewEvent() *Event {
	e := &Event{
		Chan:      make(chan interface{}, 1),
		Callbacks: []callback{},
	}
	go func() {
		for {
			e.Read()
		}
	}()
	return e
}

// Write writes data to the Event, it will not block and will not buffer if there
// are no active subscribers to the Event.
func (e *Event) Write(data interface{}) {
	select {
	case e.Chan <- data:
	default:
	}
}

// Read executes all Callbacks when new data is available.
func (e *Event) Read() {
	for s := range e.Chan {
		tmp := []callback{}
		for i := range e.Callbacks {
			go e.Callbacks[i].f(s)
			if !e.Callbacks[i].once {
				tmp = append(tmp, e.Callbacks[i])
			}
		}
		e.Callbacks = tmp
	}
}
