package gobot

type Driver struct {
	Interval string                      `json:"interval"`
	Pin      string                      `json:"pin"`
	Name     string                      `json:"name"`
	Commands []string                    `json:"commands"`
	Events   map[string]chan interface{} `json:"-"`
}

type DriverInterface interface {
	Start() bool
	Halt() bool
}
