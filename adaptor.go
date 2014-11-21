package gobot

type Adaptor interface {
	Finalize() []error
	Connect() []error
	Name() string
	Port() string
	String() string
}
