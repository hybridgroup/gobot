package gobot

import "fmt"

type Adaptor struct {
	Name      string
	Port      string
	Connected bool
	Params    map[string]interface{}
	Type      string
}

type AdaptorInterface interface {
	Finalize() bool
	Connect() bool
	port() string
	name() string
	setName(string)
	params() map[string]interface{}
	ToJSON() *JSONConnection
}

func (a *Adaptor) port() string {
	return a.Port
}

func (a *Adaptor) name() string {
	return a.Name
}

func (a *Adaptor) setName(s string) {
	a.Name = s
}

func (a *Adaptor) params() map[string]interface{} {
	return a.Params
}

func (a *Adaptor) ToJSON() *JSONConnection {
	return &JSONConnection{
		Name:    a.Name,
		Port:    a.Port,
		Adaptor: a.Type,
	}
}

func NewAdaptor(name, port, t string) *Adaptor {
	if name == "" {
		name = fmt.Sprintf("%X", Rand(int(^uint(0)>>1)))
	}
	return &Adaptor{
		Type: t,
		Name: name,
		Port: port,
	}
}
