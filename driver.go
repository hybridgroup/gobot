package gobot

type Driver struct {
	Interval string
	Pin      string
	Name     string
	Params   map[string]string
	Commands []string
	Events   map[string]chan interface{} `json:"-"`
}

type DriverInterface interface {
	Start() bool
}
