package gobot

type Driver interface {
	Start() []error
	Halt() []error
	Name() string
	Connection() Connection
}

type Pinner interface {
	Pin() string
}
