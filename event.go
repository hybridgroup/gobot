package gobot

type callback struct {
	f    func(interface{})
	once bool
}

type Event struct {
	Chan      chan interface{}
	Callbacks []callback
}

// NewEvent generates a new event by making a channel
// and start reading from it
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

// Writes sends event data to channel
func (e *Event) Write(data interface{}) {
	select {
	case e.Chan <- data:
	default:
	}
}

// Read waits data from channel and execute callbacks
// for each event when received
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
