package gobot

type callback struct {
	f    func(interface{})
	once bool
}

type Event struct {
	Chan      chan interface{}
	Callbacks []callback
}

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

func (e *Event) Write(data interface{}) {
	select {
	case e.Chan <- data:
	default:
	}
}

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
