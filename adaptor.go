package gobot

type Adaptor interface {
	Finalize() []error
	Connect() []error
	Name() string
}

type Porter interface {
	Port() string
}
