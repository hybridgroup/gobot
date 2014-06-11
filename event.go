package gobot

type Event struct {
	Chan      chan interface{}
	Callbacks []func(interface{})
}

func NewEvent() *Event {
	e := &Event{
		Chan:      make(chan interface{}, 1),
		Callbacks: []func(interface{}){},
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
		for _, f := range e.Callbacks {
			go f(s)
		}
	}
}
