package gobot

import "fmt"

type Adaptor struct {
	name        string
	port        string
	connected   bool
	adaptorType string
}

type AdaptorInterface interface {
	Finalize() bool
	Connect() bool
	Port() string
	Name() string
	Type() string
	Connected() bool
	SetConnected(bool)
	SetName(string)
	SetPort(string)
	ToJSON() *JSONConnection
}

func NewAdaptor(name string, adaptorType string, v ...interface{}) *Adaptor {
	if name == "" {
		name = fmt.Sprintf("%X", Rand(int(^uint(0)>>1)))
	}

	a := &Adaptor{
		adaptorType: adaptorType,
		name:        name,
		port:        "",
	}

	for i := range v {
		switch v[i].(type) {
		case string:
			a.port = v[i].(string)
		}
	}

	return a
}

func (a *Adaptor) Port() string {
	return a.port
}

func (a *Adaptor) SetPort(s string) {
	a.port = s
}

func (a *Adaptor) Name() string {
	return a.name
}

func (a *Adaptor) SetName(s string) {
	a.name = s
}

func (a *Adaptor) Type() string {
	return a.adaptorType
}

func (a *Adaptor) Connected() bool {
	return a.connected
}

func (a *Adaptor) SetConnected(b bool) {
	a.connected = b
}

func (a *Adaptor) ToJSON() *JSONConnection {
	return &JSONConnection{
		Name:    a.Name(),
		Adaptor: a.Type(),
	}
}
