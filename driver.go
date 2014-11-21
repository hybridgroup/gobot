package gobot

type Driver interface {
	Start() []error
	Halt() []error
	Name() string
	Pin() string
	String() string
	Connection() Connection
}
